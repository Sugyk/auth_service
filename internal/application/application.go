package application

import (
	http_api "github.com/Sugyk/auth_service/internal/api/http"
	"github.com/Sugyk/auth_service/internal/repository"
	"github.com/Sugyk/auth_service/internal/service"
	pgprovider "github.com/Sugyk/auth_service/pkg/postgres"
)

// Struct that representing whole application
type Application struct {
	db pgprovider.Provider

	repository repository.Repository
	service    service.Service
	router     http_api.Router
}

func (a *Application) Init() {
	// Init logger
	// Init DB connection
	// Init Repo layer
	// Init Service layer
	// Init Handlers layer
	// Init Router
}

func NewApplication() *Application {
	return &Application{}
}
