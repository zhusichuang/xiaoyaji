package main

import (
	"log"

	"wxcloudrun-golang/internal/app"
	"wxcloudrun-golang/internal/db"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatalf("init mysql failed: %v", err)
	}

	router := app.NewRouter()
	if err := router.Run(":80"); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
