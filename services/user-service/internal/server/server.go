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

	userv1 "github.com/SirNacou/ecommerce/services/user-service/gen/v1/userv1connect"
	"github.com/SirNacou/ecommerce/services/user-service/internal/app"
	"github.com/SirNacou/ecommerce/services/user-service/internal/infrastructure/config"
	"github.com/SirNacou/ecommerce/services/user-service/internal/infrastructure/persistence/postgres"
	"github.com/SirNacou/ecommerce/services/user-service/internal/infrastructure/security"
	grpcport "github.com/SirNacou/ecommerce/services/user-service/internal/ports/grpc"
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

	// 1. Adapters
	uow := postgres.NewUnitOfWork(pool)
	hasher := security.NewBcryptHasher(10)
	tokenProvider := security.NewJWTProvider([]byte(cfg.JWTSecret), 15*time.Minute, 7*24*time.Hour)

	// 2. Use Cases
	registerCmd := app.NewRegisterUserCommandHandler(uow, hasher)
	loginCmd := app.NewLoginUserCommandHandler(uow, hasher, tokenProvider)
	getUserQry := app.NewGetUserQueryHandler(uow)

	// 3. Transport Handler & Routing
	userHandler := grpcport.NewUserHandler(registerCmd, loginCmd, getUserQry)
	authInterceptor := grpcport.NewAuthInterceptor(tokenProvider)
	mux := http.NewServeMux()
	path, handler := userv1.NewUserServiceHandler(userHandler, connect.WithInterceptors(authInterceptor))
	mux.Handle(path, handler)

	// 4. HTTP Protocols setup
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

	log.Printf("user-service listening on port :%s ...", s.cfg.Port)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server error: %w", err)
	}

	return <-shutdownErr
}
