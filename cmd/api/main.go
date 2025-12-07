package main

import (
	"log"
	"net/http"
	"github.com/VinVorteX/NoBurn/internal/database"
	"github.com/VinVorteX/NoBurn/internal/config"
	"github.com/VinVorteX/NoBurn/internal/server"
	"github.com/VinVorteX/NoBurn/pkg/logger"
)

func main() {
	logger.Init()

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router := server.New()

	log.Printf("✅ NoBurn HR Analytics server starting in %s mode on port %s", config.AppConfig.Env, config.AppConfig.Port)
	
	if err := http.ListenAndServe(":"+config.AppConfig.Port, router); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}