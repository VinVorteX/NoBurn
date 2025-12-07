package worker

import (

	"github.com/hibiken/asynq"
	"github.com/VinVorteX/NoBurn/internal/repository"
	"github.com/VinVorteX/NoBurn/internal/services"
	"github.com/VinVorteX/NoBurn/pkg/logger"
)

type WorkerServer struct {
	server  *asynq.Server
	mux     *asynq.ServeMux
	handler *TaskHandler
}

func NewWorkerServer(redisAddr string, analyticsService *services.AnalyticsService, userRepo *repository.UserRepository, surveyRepo *repository.SurveyRepository) *WorkerServer {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	handler := NewTaskHandler(analyticsService, userRepo, surveyRepo)

	// Register task handlers
	mux.HandleFunc(TypeProcessSurvey, handler.HandleProcessSurvey)
	mux.HandleFunc(TypeCalculateChurn, handler.HandleCalculateChurn)
	mux.HandleFunc(TypeSendNotification, handler.HandleSendNotification)
	mux.HandleFunc(TypeSurveyInvitation, handler.HandleSurveyInvitation)

	return &WorkerServer{
		server:  server,
		mux:     mux,
		handler: handler,
	}
}

func (ws *WorkerServer) Start() error {
	logger.Log.Info("Starting Asynq worker server")
	return ws.server.Run(ws.mux)
}

func (ws *WorkerServer) Stop() {
	logger.Log.Info("Stopping Asynq worker server")
	ws.server.Shutdown()
}