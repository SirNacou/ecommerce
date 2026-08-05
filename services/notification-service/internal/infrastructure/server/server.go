package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SirNacou/ecommerce/pkg/auth"
	"github.com/SirNacou/ecommerce/pkg/eventbus"
	notificationv1 "github.com/SirNacou/ecommerce/services/notification-service/gen/v1/notificationv1connect"
	"github.com/SirNacou/ecommerce/services/notification-service/internal/app"
	"github.com/SirNacou/ecommerce/services/notification-service/internal/infrastructure/config"
	"github.com/SirNacou/ecommerce/services/notification-service/internal/infrastructure/persistence/postgres"
	grpcport "github.com/SirNacou/ecommerce/services/notification-service/internal/ports/grpc"
)

type Server struct {
	cfg        config.EnvConfig
	pool       *pgxpool.Pool
	server     *http.Server
	bus        *eventbus.Bus
	dispatcher *eventbus.Dispatcher
	consumer   *app.OrderCreatedConsumer
}

func New(ctx context.Context, cfg config.EnvConfig) (*Server, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	uow := postgres.NewUnitOfWork(pool)
	sendCmd := app.NewSendNotificationCommandHandler(uow)
	getQry := app.NewGetNotificationQueryHandler(uow)
	consumer := app.NewOrderCreatedConsumer(sendCmd)

	validator := auth.NewJWTValidator(cfg.JWTSecret)
	publicEndpoints := map[string]bool{}
	authInterceptor := auth.NewConnectInterceptor(validator, publicEndpoints)

	handler := grpcport.NewNotificationHandler(sendCmd, getQry)

	mux := http.NewServeMux()
	path, connectHandler := notificationv1.NewNotificationServiceHandler(
		handler,
		connect.WithInterceptors(authInterceptor),
	)
	mux.Handle(path, connectHandler)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      mux,
		Protocols:    protocols,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	srv := &Server{
		cfg:      cfg,
		pool:     pool,
		server:   httpServer,
		consumer: consumer,
	}

	// Event bus + outbox dispatcher + OrderCreated consumer
	if cfg.NATSURL != "" {
		bus, err := eventbus.New(cfg.NATSURL)
		if err != nil {
			pool.Close()
			return nil, fmt.Errorf("eventbus connect failed: %w", err)
		}
		reader := postgres.NewOutboxReader(pool)
		srv.bus = bus
		srv.dispatcher = eventbus.NewDispatcher(bus, reader, "notifications", "notifications")
	}

	return srv, nil
}

func (s *Server) Run(ctx context.Context) error {
	defer s.pool.Close()
	if s.bus != nil {
		defer s.bus.Close()
	}

	if s.dispatcher != nil {
		if err := s.dispatcher.Start(ctx); err != nil {
			return fmt.Errorf("failed to start outbox dispatcher: %w", err)
		}
	}

	if s.bus != nil {
		if err := s.bus.Subscribe(ctx, "orders", "orders.OrderCreated", "notification-order-created", s.consumer.Handle); err != nil {
			return fmt.Errorf("failed to subscribe to order events: %w", err)
		}
	}

	shutdownErr := make(chan error, 1)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr <- s.server.Shutdown(shutdownCtx)
	}()

	log.Printf("notification-service listening on port :%s ...", s.cfg.Port)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server error: %w", err)
	}

	return <-shutdownErr
}
