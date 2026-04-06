package main

import "github.com/Sugyk/auth_service/internal/application"

func main() {
	// Create Application
	app := application.NewApplication()
	// Init Application
	app.Init()
	// Start Application
}
