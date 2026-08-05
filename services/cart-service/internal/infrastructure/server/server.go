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
	cartv1 "github.com/SirNacou/ecommerce/services/cart-service/gen/v1/cartv1connect"
	"github.com/SirNacou/ecommerce/services/cart-service/internal/app"
	"github.com/SirNacou/ecommerce/services/cart-service/internal/infrastructure/clients"
	"github.com/SirNacou/ecommerce/services/cart-service/internal/infrastructure/config"
	"github.com/SirNacou/ecommerce/services/cart-service/internal/infrastructure/persistence/postgres"
	grpcport "github.com/SirNacou/ecommerce/services/cart-service/internal/ports/grpc"
)

type Server struct {
	cfg    config.EnvConfig
	pool   *pgxpool.Pool
	server *http.Server
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

	// 1. Adapters & UoW
	uow := postgres.NewUnitOfWork(pool)

	// Catalog client for authoritative product prices
	catalogClient := clients.NewCatalogClient(cfg.CatalogURL)

	// 2. Handlers
	getCartQry := app.NewGetCartQueryHandler(uow)
	addItemCmd := app.NewAddItemCommandHandler(uow, getCartQry, catalogClient)
	updateQtyCmd := app.NewUpdateQuantityCommandHandler(uow, getCartQry)
	removeItemCmd := app.NewRemoveItemCommandHandler(uow, getCartQry)
	clearCartCmd := app.NewClearCartCommandHandler(uow, getCartQry)

	// 3. Auth Interceptor (pkg/auth)
	validator := auth.NewJWTValidator(cfg.JWTSecret)
	publicEndpoints := map[string]bool{} // All Cart operations require authentication
	authInterceptor := auth.NewConnectInterceptor(validator, publicEndpoints)

	// 4. Transport Handler
	handler := grpcport.NewCartHandler(getCartQry, addItemCmd, updateQtyCmd, removeItemCmd, clearCartCmd)

	mux := http.NewServeMux()
	path, connectHandler := cartv1.NewCartServiceHandler(
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

	return &Server{
		cfg:    cfg,
		pool:   pool,
		server: httpServer,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	defer s.pool.Close()

	shutdownErr := make(chan error, 1)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr <- s.server.Shutdown(shutdownCtx)
	}()

	log.Printf("cart-service listening on port :%s ...", s.cfg.Port)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server error: %w", err)
	}

	return <-shutdownErr
}
