package worker

import (
	"encoding/json"
	"testing"
)

func TestNewProcessSurveyTask(t *testing.T) {
	task, err := NewProcessSurveyTask(1, 5, "en")
	if err != nil {
		t.Fatalf("NewProcessSurveyTask failed: %v", err)
	}

	if task.Type() != TypeProcessSurvey {
		t.Errorf("Expected task type %s, got %s", TypeProcessSurvey, task.Type())
	}

	var payload SurveyPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		t.Fatalf("Failed to unmarshal payload: %v", err)
	}

	if payload.ResponseID != 1 {
		t.Errorf("Expected ResponseID 1, got %d", payload.ResponseID)
	}

	if payload.UserID != 5 {
		t.Errorf("Expected UserID 5, got %d", payload.UserID)
	}

	if payload.Language != "en" {
		t.Errorf("Expected Language en, got %s", payload.Language)
	}
}

func TestNewCalculateChurnTask(t *testing.T) {
	task, err := NewCalculateChurnTask(5, 1)
	if err != nil {
		t.Fatalf("NewCalculateChurnTask failed: %v", err)
	}

	if task.Type() != TypeCalculateChurn {
		t.Errorf("Expected task type %s, got %s", TypeCalculateChurn, task.Type())
	}

	var payload ChurnPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		t.Fatalf("Failed to unmarshal payload: %v", err)
	}

	if payload.UserID != 5 {
		t.Errorf("Expected UserID 5, got %d", payload.UserID)
	}

	if payload.CompanyID != 1 {
		t.Errorf("Expected CompanyID 1, got %d", payload.CompanyID)
	}
}

func TestNewSendNotificationTask(t *testing.T) {
	task, err := NewSendNotificationTask(5, "high_risk", "Test message")
	if err != nil {
		t.Fatalf("NewSendNotificationTask failed: %v", err)
	}

	if task.Type() != TypeSendNotification {
		t.Errorf("Expected task type %s, got %s", TypeSendNotification, task.Type())
	}

	var payload NotificationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		t.Fatalf("Failed to unmarshal payload: %v", err)
	}

	if payload.UserID != 5 {
		t.Errorf("Expected UserID 5, got %d", payload.UserID)
	}

	if payload.Type != "high_risk" {
		t.Errorf("Expected Type high_risk, got %s", payload.Type)
	}

	if payload.Message != "Test message" {
		t.Errorf("Expected Message 'Test message', got %s", payload.Message)
	}
}