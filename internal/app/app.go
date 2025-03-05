package app

import (
	"final/internal/config"
	"final/internal/repository"
	"final/internal/service"
	"final/internal/transport/gen"
	grpc2 "final/internal/transport/grpc"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

const (
	PostgresDriver = "postgres"
	TCPNetwork     = "tcp"
)

// App is a struct with Run, Shutdown methods, which represents application.
type App struct {
	l          *zap.SugaredLogger
	db         *sqlx.DB
	grpcServer *grpc.Server
	cfg        *config.Config
}

// New creates connection to db, registers grpc endpoints and returns new App instance.
// Returns error if failed to connect to db.
func New(cfg *config.Config, l *zap.SugaredLogger) (*App, error) {

	l.Debugf("starting app with config: %v", *cfg)

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)

	db, err := sqlx.Connect(PostgresDriver, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	g := grpc.NewServer()

	rateRepo := repository.NewRateRepository(db)
	rateFetcher := service.NewGarantexFetcher()

	rateService := service.NewRateService(rateRepo, rateFetcher, l)

	rateServiceServer := grpc2.NewRateServiceServer(rateService)

	gen.RegisterRateServiceServer(g, rateServiceServer)

	healthServiceServer := grpc2.NewHealthServiceServer()

	gen.RegisterHealthServiceServer(g, healthServiceServer)

	app := &App{
		l:          l,
		db:         db,
		grpcServer: g,
		cfg:        cfg,
	}

	return app, nil
}

// Run starts the app.
// Returns an error if failed to listen port or failed to serve.
func (a *App) Run() error {

	addr := fmt.Sprintf("%s:%s", a.cfg.AppIP, a.cfg.AppPort)

	lis, err := net.Listen(TCPNetwork, addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	a.l.Infof("server is listening on %s", addr)

	if err := a.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down db connection and grpc server.
// Returns an error if failed to close db connection.
func (a *App) Shutdown() error {

	a.l.Infoln("shutting down grpc server")
	a.grpcServer.GracefulStop()

	a.l.Infoln("shutting down db connection")
	if err := a.db.Close(); err != nil {
		return fmt.Errorf("failed to close db connection: %w", err)
	}

	return nil
}
