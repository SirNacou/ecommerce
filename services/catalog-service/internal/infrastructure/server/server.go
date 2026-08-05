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
	catalogv1 "github.com/SirNacou/ecommerce/services/catalog-service/gen/v1/catalogv1connect"
	"github.com/SirNacou/ecommerce/services/catalog-service/internal/app"
	"github.com/SirNacou/ecommerce/services/catalog-service/internal/infrastructure/config"
	"github.com/SirNacou/ecommerce/services/catalog-service/internal/infrastructure/persistence/postgres"
	grpcport "github.com/SirNacou/ecommerce/services/catalog-service/internal/ports/grpc"
)

type Server struct {
	cfg        config.EnvConfig
	pool       *pgxpool.Pool
	server     *http.Server
	bus        *eventbus.Bus
	dispatcher *eventbus.Dispatcher
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

	// 1. Adapters & Unit of Work
	uow := postgres.NewUnitOfWork(pool)

	// 2. CQRS Command Handlers
	createProductCmd := app.NewCreateProductCommandHandler(uow)
	getProductQry := app.NewGetProductQueryHandler(uow)
	getProductsByIdsQry := app.NewGetProductsByIdsQueryHandler(uow)
	listProductsQry := app.NewListProductsQueryHandler(uow)
	createCategoryCmd := app.NewCreateCategoryCommandHandler(uow)
	listCategoriesQry := app.NewListCategoriesQueryHandler(uow)
	getStockQry := app.NewGetStockQueryHandler(uow)
	setStockCmd := app.NewSetStockCommandHandler(uow)
	reserveCmd := app.NewReserveStockCommandHandler(uow)
	releaseCmd := app.NewReleaseStockCommandHandler(uow)

	// 3. Auth Interceptor Setup (using shared pkg/auth)
	validator := auth.NewJWTValidator(cfg.JWTSecret)
	publicEndpoints := map[string]bool{
		"/catalog.v1.CatalogService/ListProducts":     true,
		"/catalog.v1.CatalogService/GetProduct":       true,
		"/catalog.v1.CatalogService/GetProductsByIds": true,
		"/catalog.v1.CatalogService/ListCategories":   true,
		"/catalog.v1.CatalogService/GetStock":         true,
	}
	authInterceptor := auth.NewConnectInterceptor(validator, publicEndpoints)

	// 4. Transport Handler & ServeMux Routing
	catalogHandler := grpcport.NewCatalogHandler(
		createProductCmd, getProductQry, getProductsByIdsQry, listProductsQry, createCategoryCmd, listCategoriesQry,
		getStockQry, setStockCmd, reserveCmd, releaseCmd,
	)

	mux := http.NewServeMux()
	path, handler := catalogv1.NewCatalogServiceHandler(
		catalogHandler,
		connect.WithInterceptors(authInterceptor),
	)
	mux.Handle(path, handler)

	// Health Check Endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// 5. HTTP Protocol Configuration (Go 1.22+)
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
		cfg:    cfg,
		pool:   pool,
		server: httpServer,
	}

	// Event bus + outbox dispatcher
	if cfg.NATSURL != "" {
		bus, err := eventbus.New(cfg.NATSURL)
		if err != nil {
			pool.Close()
			return nil, fmt.Errorf("eventbus connect failed: %w", err)
		}
		reader := postgres.NewOutboxReader(pool)
		srv.bus = bus
		srv.dispatcher = eventbus.NewDispatcher(bus, reader, "inventory", "inventory")
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

	shutdownErr := make(chan error, 1)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr <- s.server.Shutdown(shutdownCtx)
	}()

	log.Printf("catalog-service listening on port :%s ...", s.cfg.Port)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server error: %w", err)
	}

	return <-shutdownErr
}
