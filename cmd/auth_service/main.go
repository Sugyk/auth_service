package main

import (
	"context"
	"log"

	"github.com/Sugyk/auth_service/internal/application"
)

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
