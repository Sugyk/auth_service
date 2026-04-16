package main

import (
	"context"

	"github.com/Sugyk/auth_service/internal/application"
)

func main() {
	ctx := context.Background()
	// Create Application
	app := application.NewApplication()

	// Init Application
	app.Init(ctx)

	// Start Application
}
