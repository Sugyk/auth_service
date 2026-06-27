package main

import (
	"context"
	"log"

	"github.com/Sugyk/auth_service/internal/application"
)

// @title           Auth_service API
// @version         1.0.0
// @description     This service provide API for using JWT authentification in your service.

// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT

// @host            localhost:8080
// @BasePath        /api/v1

// @schemes         http https

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Токен в формате: Bearer <token>

func main() {
	ctx := context.Background()
	// Create Application
	app := application.NewApplication()

	// Init Application
	app.Init(ctx)

	// Start Application
	if err := app.Start(ctx); err != nil {
		log.Fatalln("all systems ended with error:", err)
	}
	app.Shutdown(ctx)
}
