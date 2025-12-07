package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/VinVorteX/NoBurn/internal/utils"
	"github.com/VinVorteX/NoBurn/internal/worker"
)

type SlackEventPayload struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
	Event     struct {
		Type    string `json:"type"`
		User    string `json:"user"`
		Text    string `json:"text"`
		Channel string `json:"channel"`
	} `json:"event"`
}

type SurveyWebhookPayload struct {
	UserID    uint     `json:"user_id"`
	Responses []string `json:"responses"`
	Source    string   `json:"source"`
}

func HandleSlackWebhook(w http.ResponseWriter, r *http.Request) {
	var payload SlackEventPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if payload.Type == "url_verification" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(payload.Challenge))
		return
	}

	if payload.Event.Type == "message" {
		workerClient := worker.NewClient("localhost:6379")
		defer workerClient.Close()
		workerClient.EnqueueSurveyProcessing(0, 0, "en")
	}

	utils.WriteSuccess(w, map[string]string{"status": "ok"})
}

func HandleSurveyWebhook(w http.ResponseWriter, r *http.Request) {
	var payload SurveyWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	workerClient := worker.NewClient("localhost:6379")
	defer workerClient.Close()

	if err := workerClient.EnqueueSurveyProcessing(0, payload.UserID, "en"); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to process survey")
		return
	}

	utils.WriteSuccess(w, map[string]string{
		"status":  "accepted",
		"message": "Survey response queued for processing",
	})
}

func HandleEmailWebhook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	utils.WriteSuccess(w, map[string]string{"status": "ok"})
}