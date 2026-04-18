package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	http_api "github.com/Sugyk/auth_service/internal/api/http"
	"github.com/Sugyk/auth_service/internal/api/http/handlers"
	"github.com/Sugyk/auth_service/internal/config"
	"github.com/Sugyk/auth_service/internal/pkg/hasher"
	"github.com/Sugyk/auth_service/internal/repository"
	"github.com/Sugyk/auth_service/internal/service"
	"github.com/Sugyk/auth_service/pkg/logger"
	"github.com/Sugyk/auth_service/pkg/postgres"
)

const LOGLEVEL = "info"

// Struct that representing whole application
type Application struct {
	logger logger.Logger
	db     *postgres.Provider

	dbCfg     *config.PgConfig
	hasherCfg *config.HasherConfig

	repository *repository.Repository

	service *service.Service

	handler *handlers.Handler

	router *http_api.Router
}

func (a *Application) Init(ctx context.Context) {
	// Init logger
	if err := a.InitLogger(LOGLEVEL); err != nil {
		log.Fatalln("Init application logger error:", err)
	}
	// Load config
	if err := a.LoadConfigs(); err != nil {
		log.Fatalln("Application config loading error:", err)
	}
	// Init DB connection
	if err := a.InitDB(ctx); err != nil {
		log.Fatalln("Init application DB error:", err)
	}
	// Init Repo layer
	if err := a.InitRepository(); err != nil {
		log.Fatalln("Init application repository error:", err)
	}
	// Init Service layer
	if err := a.InitService(); err != nil {
		log.Fatalln("Init application service error:", err)
	}
	// Init Handlers layer
	if err := a.InitHandler(); err != nil {
		log.Fatalln("Init application handler error:", err)
	}
	// Init Router
	if err := a.InitRouter(); err != nil {
		log.Fatalln("Init application router error:", err)
	}
	a.logger.Info(ctx, "App initialisation completed successfully")
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) InitLogger(level string) error {
	a.logger = logger.New(level)

	return nil
}

func (a *Application) LoadConfigs() error {
	dbCfg, hasherCfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	a.dbCfg = dbCfg
	a.hasherCfg = hasherCfg

	return nil
}

func (a *Application) InitDB(ctx context.Context) error {
	provider := postgres.NewProvider(
		a.logger,
		a.dbCfg.ConnStr,
		a.dbCfg.MaxConns,
		a.dbCfg.MinConns,
		a.dbCfg.MaxConnLifetime,
		a.dbCfg.MaxConnIdleTime,
	)
	if err := provider.Open(ctx); err != nil {
		return err
	}

	a.db = provider

	return nil
}

func (a *Application) InitRepository() error {
	a.repository = repository.NewRepository()

	return nil
}

func (a *Application) InitService() error {
	txManager := postgres.NewTxManager(a.db.DB())
	passwordHasher := hasher.NewPasswordHasher(a.hasherCfg.Cost)

	a.service = service.NewService(
		a.repository,
		txManager,
		passwordHasher,
	)

	return nil
}

func (a *Application) InitHandler() error {
	a.handler = handlers.NewHandler(a.service)

	return nil
}

func (a *Application) InitRouter() error {
	a.router = http_api.NewRouter(a.handler)

	return nil
}

func (a *Application) Start(ctx context.Context) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	errchan := make(chan error, 1)

	go func() {
		errchan <- a.router.Start()
	}()

	var err error

	select {
	case <-sigChan:
		a.logger.Info(ctx, "got signal to shutdown")
		return nil
	case err = <-errchan:
		return fmt.Errorf("server crashed with error: %w", err)
	}
}

func (a *Application) Shutdown(ctx context.Context) {
	a.logger.Info(ctx, "starting gracefull shutdown")

	if err := a.router.Shutdown(ctx); err != nil {
		a.logger.Info(ctx, "server closed with error", "error", err)
	}

	a.db.Close()
	a.logger.Info(ctx, "db connection closed")

	a.logger.Info(ctx, "gracefull shutdown completed without error")
}
