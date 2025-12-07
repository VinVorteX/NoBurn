package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/VinVorteX/NoBurn/internal/repository"
	"github.com/VinVorteX/NoBurn/internal/services"
	"github.com/VinVorteX/NoBurn/pkg/logger"
	"go.uber.org/zap"
)

type TaskHandler struct {
	analyticsService *services.AnalyticsService
	userRepo         *repository.UserRepository
	surveyRepo       *repository.SurveyRepository
}

func NewTaskHandler(analyticsService *services.AnalyticsService, userRepo *repository.UserRepository, surveyRepo *repository.SurveyRepository) *TaskHandler {
	return &TaskHandler{
		analyticsService: analyticsService,
		userRepo:         userRepo,
		surveyRepo:       surveyRepo,
	}
}

func (h *TaskHandler) HandleProcessSurvey(ctx context.Context, t *asynq.Task) error {
	var payload SurveyPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	logger.Log.Info("Processing survey response",
		zap.Uint("response_id", payload.ResponseID),
		zap.Uint("user_id", payload.UserID),
		zap.String("language", payload.Language),
	)

	// Process sentiment analysis in background
	// TODO: Get actual response from database and process
	
	// Trigger churn calculation after processing
	churnTask, err := NewCalculateChurnTask(payload.UserID, 1) // TODO: Get actual company ID
	if err != nil {
		logger.Log.Error("Failed to create churn task", zap.Error(err))
		return err
	}

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
	defer client.Close()

	_, err = client.Enqueue(churnTask)
	if err != nil {
		logger.Log.Error("Failed to enqueue churn task", zap.Error(err))
	}

	return nil
}

func (h *TaskHandler) HandleCalculateChurn(ctx context.Context, t *asynq.Task) error {
	log.Printf("🔄 WORKER: Starting churn calculation task")
	
	var payload ChurnPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("❌ WORKER: Failed to unmarshal payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	log.Printf("🔄 WORKER: Processing churn for user %d, company %d", payload.UserID, payload.CompanyID)

	// Calculate churn risk using ML
	risk, err := h.analyticsService.AnalyzeUserChurnRisk(payload.UserID)
	if err != nil {
		log.Printf("❌ WORKER: Failed to analyze churn risk: %v", err)
		return fmt.Errorf("failed to analyze churn risk: %v", err)
	}

	log.Printf("📊 WORKER: Churn risk calculated: %.2f for user %d", risk.RiskScore, payload.UserID)

	// Save risk to database
	attritionRepo := repository.NewAttritionRepository()
	if err := attritionRepo.Create(risk); err != nil {
		log.Printf("❌ WORKER: Failed to save risk: %v", err)
		return fmt.Errorf("failed to save risk: %v", err)
	}
	log.Printf("✅ WORKER: Risk saved to database")

	// If high risk, send notification
	if risk.RiskScore > 0.7 {
		log.Printf("🚨 WORKER: High risk detected (%.2f), sending notification", risk.RiskScore)
		notificationTask, err := NewSendNotificationTask(
			payload.UserID,
			"high_churn_risk",
			fmt.Sprintf("Employee has high churn risk: %.2f", risk.RiskScore),
		)
		if err != nil {
			return err
		}

		client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
		defer client.Close()

		_, err = client.Enqueue(notificationTask)
		if err != nil {
			log.Printf("❌ WORKER: Failed to enqueue notification: %v", err)
		} else {
			log.Printf("✅ WORKER: Notification task enqueued")
		}
	} else {
		log.Printf("ℹ️ WORKER: Risk score %.2f below threshold, no notification sent", risk.RiskScore)
	}

	log.Printf("✅ WORKER: Churn calculation completed successfully")
	return nil
}

func (h *TaskHandler) HandleSendNotification(ctx context.Context, t *asynq.Task) error {
	var payload NotificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	logger.Log.Info("Sending notification",
		zap.Uint("user_id", payload.UserID),
		zap.String("type", payload.Type),
		zap.String("message", payload.Message),
	)

	// Get user details
	user, err := h.userRepo.GetByID(payload.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	// Initialize notification service with config from env
	notifService := services.NewNotificationService(
		"", // Slack webhook from env
		"smtp.gmail.com",
		587,
		"", // SMTP user from env
		"", // SMTP password from env
	)
	log.Printf("📧 Sending notifications for user: %s", user.Email)

	// Send Slack alert
	log.Printf("🔔 Attempting Slack alert...")
	if err := notifService.SendSlackAlert(user.Name, payload.Message, 0.85); err != nil {
		log.Printf("❌ Slack alert failed: %v", err)
	} else {
		log.Printf("✅ Slack alert sent successfully")
	}

	// Send Email alert to configured alert email
	log.Printf("📧 Attempting email alert...")
	if err := notifService.SendEmailAlert("", user.Name, payload.Message, 0.85); err != nil {
		log.Printf("❌ Email alert FAILED: %v", err)
	} else {
		log.Printf("✅ Email alert sent successfully")
	}

	log.Printf("✅ NOTIFICATION COMPLETED: %s - %s", payload.Type, payload.Message)
	return nil
}

func (h *TaskHandler) HandleSurveyInvitation(ctx context.Context, t *asynq.Task) error {
	var payload SurveyInvitationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	log.Printf("📧 WORKER: Sending survey invitation to %s", payload.Email)

	// Get survey details
	survey, err := h.surveyRepo.GetByID(payload.SurveyID)
	if err != nil {
		log.Printf("❌ WORKER: Failed to get survey: %v", err)
		return fmt.Errorf("failed to get survey: %v", err)
	}

	// Get company SMTP settings
	companyRepo := repository.NewCompanyRepository()
	company, err := companyRepo.GetByID(survey.CompanyID)
	if err != nil {
		log.Printf("❌ WORKER: Failed to get company: %v", err)
		return fmt.Errorf("failed to get company: %v", err)
	}

	// Use company SMTP or fallback to env
	smtpHost := company.SMTPHost
	smtpPort := company.SMTPPort
	smtpUser := company.SMTPUser
	smtpPassword := company.SMTPPassword

	if smtpUser == "" {
		// Fallback to env variables
		smtpHost = os.Getenv("SMTP_HOST")
		smtpPort = 587
		smtpUser = os.Getenv("SMTP_USER")
		smtpPassword = os.Getenv("SMTP_PASSWORD")
	}

	// Generate survey link with token
	surveyLink := fmt.Sprintf("http://localhost:3002/survey/%d?token=%d", payload.SurveyID, payload.UserID)

	// Create email content
	emailSubject := fmt.Sprintf("New Survey: %s", survey.Title)
	emailBody := fmt.Sprintf(`
Hi there!

You have been invited to participate in a new survey: "%s"

Please click the link below to complete the survey:
%s

This survey will help us improve your work experience.

Thank you!
NoBurn HR Team
`, survey.Title, surveyLink)

	// Initialize notification service
	notifService := services.NewNotificationService(
		"", // No Slack for survey invitations
		smtpHost,
		smtpPort,
		smtpUser,
		smtpPassword,
	)

	// Send email invitation
	log.Printf("📧 WORKER: Sending survey email to %s", payload.Email)
	if err := notifService.SendSurveyInvitation(payload.Email, emailSubject, emailBody); err != nil {
		log.Printf("❌ WORKER: Failed to send survey invitation: %v", err)
		return fmt.Errorf("failed to send survey invitation: %v", err)
	}

	log.Printf("✅ WORKER: Survey invitation sent successfully to %s", payload.Email)
	return nil
}