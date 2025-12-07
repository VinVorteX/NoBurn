package worker

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeProcessSurvey      = "survey:process"
	TypeCalculateChurn     = "churn:calculate"
	TypeSendNotification   = "notification:send"
	TypeSurveyInvitation   = "survey:invitation"
)

type SurveyPayload struct {
	ResponseID uint   `json:"response_id"`
	UserID     uint   `json:"user_id"`
	Language   string `json:"language"`
}

type ChurnPayload struct {
	UserID    uint `json:"user_id"`
	CompanyID uint `json:"company_id"`
}

type NotificationPayload struct {
	UserID  uint   `json:"user_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type SurveyInvitationPayload struct {
	SurveyID uint   `json:"survey_id"`
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
}

func NewProcessSurveyTask(responseID, userID uint, language string) (*asynq.Task, error) {
	payload, err := json.Marshal(SurveyPayload{
		ResponseID: responseID,
		UserID:     userID,
		Language:   language,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProcessSurvey, payload), nil
}

func NewCalculateChurnTask(userID, companyID uint) (*asynq.Task, error) {
	payload, err := json.Marshal(ChurnPayload{
		UserID:    userID,
		CompanyID: companyID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCalculateChurn, payload), nil
}

func NewSendNotificationTask(userID uint, notificationType, message string) (*asynq.Task, error) {
	payload, err := json.Marshal(NotificationPayload{
		UserID:  userID,
		Type:    notificationType,
		Message: message,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSendNotification, payload), nil
}

func NewSurveyInvitationTask(surveyID, userID uint, email string) (*asynq.Task, error) {
	payload, err := json.Marshal(SurveyInvitationPayload{
		SurveyID: surveyID,
		UserID:   userID,
		Email:    email,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSurveyInvitation, payload), nil
}