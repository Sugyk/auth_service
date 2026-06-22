package main

import (
	"context"
	"log"

	"github.com/Sugyk/auth_service/internal/application"
)

// @title           Название API
// @version         1.0.0
// @description     Описание вашего API. Можно многострочное.

// @contact.name    Имя контактного лица
// @contact.url     https://example.com/support
// @contact.email   support@example.com

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
