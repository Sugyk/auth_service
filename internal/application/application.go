package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	grpc_api "github.com/Sugyk/auth_service/internal/api/grpc"
	http_api "github.com/Sugyk/auth_service/internal/api/http"
	"github.com/Sugyk/auth_service/internal/api/http/handlers"
	"github.com/Sugyk/auth_service/internal/config"
	"github.com/Sugyk/auth_service/internal/pkg/hasher"
	"github.com/Sugyk/auth_service/internal/pkg/jwt_manager"
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

	cfg *config.AppConfig

	repository *repository.Repository

	service *service.Service

	handler *handlers.Handler

	router     *http_api.Router
	grpcRouter *grpc_api.Router
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
	// Init gRPC server
	if err := a.InitGRPCServer(); err != nil {
		log.Fatalln("Init application gRPC server error:", err)
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
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	a.cfg = cfg

	return nil
}

func (a *Application) InitDB(ctx context.Context) error {
	provider := postgres.NewProvider(
		a.logger,
		a.cfg.DBCfg.ConnStr,
		a.cfg.DBCfg.MaxConns,
		a.cfg.DBCfg.MinConns,
		a.cfg.DBCfg.MaxConnLifetime,
		a.cfg.DBCfg.MaxConnIdleTime,
	)
	if err := provider.Open(ctx); err != nil {
		return err
	}

	a.db = provider

	return nil
}

func (a *Application) InitRepository() error {
	a.repository = repository.NewRepository(a.db.DB())

	return nil
}

func (a *Application) InitService() error {
	txManager := postgres.NewTxManager(a.db.DB())
	passwordHasher := hasher.NewPasswordHasher(a.cfg.HasherCfg.Cost)
	jwtManager, err := jwt_manager.NewJWTManager(
		[]byte(a.cfg.JWTConfig.Secret),
		a.cfg.JWTConfig.TTL,
	)
	if err != nil {
		return err
	}

	a.service = service.NewService(
		a.repository,
		txManager,
		passwordHasher,
		jwtManager,
	)

	return nil
}

func (a *Application) InitHandler() error {
	a.handler = handlers.NewHandler(a.service, a.logger)

	return nil
}

func (a *Application) InitRouter() error {
	a.router = http_api.NewRouter(a.handler)

	return nil
}

func (a *Application) InitGRPCServer() error {
	authServer := grpc_api.NewServer(a.service, a.logger)

	router, err := grpc_api.NewRouter(a.cfg.GRPCConfig.Addr, authServer)
	if err != nil {
		return err
	}

	a.grpcRouter = router

	return nil
}

func (a *Application) Start(ctx context.Context) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	errchan := make(chan error, 2)

	go func() {
		errchan <- a.router.Start()
	}()

	go func() {
		errchan <- a.grpcRouter.Start()
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
		a.logger.Info(ctx, "http server closed with error", "error", err)
	} else {
		a.logger.Info(ctx, "http server closed with no errors")
	}

	if err := a.grpcRouter.Shutdown(ctx); err != nil {
		a.logger.Info(ctx, "grpc server closed with error", "error", err)
	} else {
		a.logger.Info(ctx, "grpc server closed with no errors")
	}

	a.db.Close()
	a.logger.Info(ctx, "db connection closed")

	a.logger.Info(ctx, "gracefull shutdown completed without errors")
}
