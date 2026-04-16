package application

import (
	"context"
	"log"

	http_api "github.com/Sugyk/auth_service/internal/api/http"
	"github.com/Sugyk/auth_service/internal/config"
	"github.com/Sugyk/auth_service/internal/repository"
	"github.com/Sugyk/auth_service/internal/service"
	"github.com/Sugyk/auth_service/pkg/logger"
	"github.com/Sugyk/auth_service/pkg/postgres"
	pgprovider "github.com/Sugyk/auth_service/pkg/postgres"
)

const LOGLEVEL = "warn"

// Struct that representing whole application
type Application struct {
	logger logger.Logger
	db     *pgprovider.Provider

	dbCfg     *config.PgConfig
	hasherCfg *config.HasherConfig

	repository *repository.Repository
	service    *service.Service
	router     *http_api.Router
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
	// Init Service layer
	// Init Handlers layer
	// Init Router
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
	return nil
}
