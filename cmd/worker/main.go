package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/VinVorteX/NoBurn/internal/config"
	"github.com/VinVorteX/NoBurn/internal/database"
	"github.com/VinVorteX/NoBurn/internal/repository"
	"github.com/VinVorteX/NoBurn/internal/services"
	"github.com/VinVorteX/NoBurn/internal/worker"
	"github.com/VinVorteX/NoBurn/pkg/logger"
)

func main() {
	logger.Init()

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.Connect(); err != nil {
		logger.Log.Fatal("Failed to connect to database")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	surveyRepo := repository.NewSurveyRepository()

	// Initialize services
	analyticsService := services.NewAnalyticsService(userRepo, surveyRepo, "http://localhost:5000")

	// Initialize worker server
	// Parse Redis URL to get host:port (Asynq expects "host:port" not "redis://host:port")
	redisAddr := strings.TrimPrefix(config.AppConfig.RedisURL, "redis://")
	workerServer := worker.NewWorkerServer(redisAddr, analyticsService, userRepo, surveyRepo)

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down worker server...")
		workerServer.Stop()
		os.Exit(0)
	}()

	log.Println("NoBurn worker server starting...")
	if err := workerServer.Start(); err != nil {
		logger.Log.Fatal("Failed to start worker server")
	}
}