package main

import (
	"log"

	"go-oj/internal/bootstrap"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("failed to bootstrap app: %v", err)
	}

	log.Printf("starting %s on :%s", app.Config.AppName, app.Config.HTTPPort)
	if err := app.Router.Run(":" + app.Config.HTTPPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
