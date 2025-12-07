package health

import (
	"net/http"
	"time"

	"github.com/VinVorteX/NoBurn/internal/database"
	"github.com/VinVorteX/NoBurn/internal/utils"
	"github.com/VinVorteX/NoBurn/pkg/logger"
	"go.uber.org/zap"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	UpTime    string `json:"uptime"`
	TimeStamp int64  `json:"timestamp"`
	Services  struct {
		Database bool `json:"database"`
	} `json:"services"`
}

var startTime = time.Now()

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbHealthy := true

	// Check database health
	if database.DB != nil {
		sqlDB, err := database.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbHealthy = false
		}
	} else {
		dbHealthy = false
	}

	// Check overall health
	status := "ok"
	statusCode := http.StatusOK
	if !dbHealthy {
		status = "degraded"
		statusCode = http.StatusServiceUnavailable
	}

	// Return health response
	resp := HealthResponse{
		Status:    status,
		Version:   "1.0.0",
		UpTime:    time.Since(startTime).Round(time.Second).String(),
		TimeStamp: time.Now().Unix(),
	}
	resp.Services.Database = dbHealthy

	logger.Log.Info("Health check", 
		zap.String("status", status), 
		zap.Bool("database", dbHealthy),
	)

	utils.WriteJSON(w, statusCode, resp)
}