package survey

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/VinVorteX/NoBurn/internal/config"
	"github.com/VinVorteX/NoBurn/internal/models"
	"github.com/VinVorteX/NoBurn/internal/repository"
	"github.com/VinVorteX/NoBurn/internal/services/sentiment"
	"github.com/VinVorteX/NoBurn/internal/worker"
	"github.com/VinVorteX/NoBurn/internal/utils"
)

type CreateSurveyRequest struct {
	Title     string   `json:"title"`
	Questions []string `json:"questions"`
}

type SurveyResponseRequest struct {
	SurveyID  uint     `json:"survey_id"`
	Responses []string `json:"responses"`
}

func CreateSurvey(w http.ResponseWriter, r *http.Request) {
	var req CreateSurveyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	userID := r.Context().Value("userID").(uint)
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not found")
		return
	}

	survey := &models.Survey{
		CompanyID: user.CompanyID,
		Title:     req.Title,
		Questions: models.StringArray(req.Questions),
		IsActive:  true,
	}

	surveyRepo := repository.NewSurveyRepository()
	if err := surveyRepo.Create(survey); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create survey")
		return
	}

	// Send survey invitations to employees only (not admins)
	employees, err := userRepo.GetByCompanyID(user.CompanyID)
	if err == nil {
		workerClient := getWorkerClient()
		if workerClient != nil {
			sentCount := 0
			for _, employee := range employees {
				// Skip admins
				if employee.Role == "hr_admin" {
					continue
				}
				log.Printf("📧 Sending survey invitation to %s", employee.Email)
				if err := workerClient.EnqueueSurveyInvitation(survey.ID, employee.ID, employee.Email); err != nil {
					log.Printf("❌ Failed to enqueue survey invitation for %s: %v", employee.Email, err)
				} else {
					sentCount++
				}
			}
			log.Printf("✅ Survey invitations enqueued for %d employees", sentCount)
		}
	}

	w.WriteHeader(http.StatusCreated)
	utils.WriteSuccess(w, survey)
}

func SubmitResponse(w http.ResponseWriter, r *http.Request) {
	var req SurveyResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	userID := r.Context().Value("userID").(uint)
	
	// Get user's company language
	userRepo := repository.NewUserRepository()
	user, _ := userRepo.GetByID(userID)
	companyRepo := repository.NewCompanyRepository()
	company, _ := companyRepo.GetByID(user.CompanyID)
	language := "en"
	if company != nil && company.Language != "" {
		language = company.Language
	}

	// Combine responses for sentiment analysis
	combinedText := strings.Join(req.Responses, " ")

	// Analyze sentiment using IndicBERT (with fallback to rule-based)
	log.Printf("Analyzing: '%s' in language: %s", combinedText, language)
	mlService := sentiment.NewMLService(config.AppConfig.HuggingFaceToken)
	sentimentScore, err := mlService.AnalyzeSentiment(combinedText, language)
	if err != nil {
		log.Printf("⚠️ ML error: %v, using rule-based", err)
		sentimentScore = sentiment.AnalyzeSentiment(combinedText, language)
	}
	log.Printf("✅ Final Sentiment: %f", sentimentScore)

	// Save response with sentiment
	response := &models.SurveyResponse{
		SurveyID:  req.SurveyID,
		UserID:    userID,
		Responses: models.StringArray(req.Responses),
		Sentiment: sentimentScore,
	}

	surveyRepo := repository.NewSurveyRepository()
	if err := surveyRepo.CreateResponse(response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to save response")
		return
	}
	
	// Enqueue background job for churn calculation
	workerClient := getWorkerClient()
	if workerClient != nil {
		log.Printf("📤 Enqueueing churn calculation for user %d", userID)
		if err := workerClient.EnqueueChurnCalculation(userID, user.CompanyID); err != nil {
			log.Printf("❌ Failed to enqueue churn calculation: %v", err)
		} else {
			log.Printf("✅ Churn calculation job enqueued successfully")
		}
	} else {
		log.Printf("⚠️ Worker client is nil, skipping background job")
	}

	w.WriteHeader(http.StatusCreated)
	utils.WriteSuccess(w, response)
}

// TODO: Replace with proper dependency injection
func getWorkerClient() *worker.Client {
	redisAddr := strings.TrimPrefix(config.AppConfig.RedisURL, "redis://")
	return worker.NewClient(redisAddr)
}

func GetSurveys(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not found")
		return
	}

	surveyRepo := repository.NewSurveyRepository()
	surveys, err := surveyRepo.GetByCompanyID(user.CompanyID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch surveys")
		return
	}

	utils.WriteSuccess(w, surveys)
}

func GetPublicSurvey(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "surveyID")
	if surveyID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Survey ID required")
		return
	}

	surveyRepo := repository.NewSurveyRepository()
	survey, err := surveyRepo.GetByID(uint(parseUint(surveyID)))
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Survey not found")
		return
	}

	if !survey.IsActive {
		utils.WriteError(w, http.StatusGone, "Survey is no longer active")
		return
	}

	utils.WriteSuccess(w, survey)
}

type PublicSurveyResponseRequest struct {
	SurveyID  uint     `json:"survey_id"`
	UserToken uint     `json:"user_token"`
	Responses []string `json:"responses"`
}

func SubmitPublicResponse(w http.ResponseWriter, r *http.Request) {
	var req PublicSurveyResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Get user from token (user ID)
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(req.UserToken)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid user token")
		return
	}

	// Get company language
	companyRepo := repository.NewCompanyRepository()
	company, _ := companyRepo.GetByID(user.CompanyID)
	language := "en"
	if company != nil && company.Language != "" {
		language = company.Language
	}

	// Combine responses for sentiment analysis
	combinedText := strings.Join(req.Responses, " ")

	// Analyze sentiment
	mlService := sentiment.NewMLService(config.AppConfig.HuggingFaceToken)
	sentimentScore, err := mlService.AnalyzeSentiment(combinedText, language)
	if err != nil {
		log.Printf("⚠️ ML error: %v, using rule-based", err)
		sentimentScore = sentiment.AnalyzeSentiment(combinedText, language)
	}

	// Save response
	response := &models.SurveyResponse{
		SurveyID:  req.SurveyID,
		UserID:    req.UserToken,
		Responses: models.StringArray(req.Responses),
		Sentiment: sentimentScore,
	}

	surveyRepo := repository.NewSurveyRepository()
	if err := surveyRepo.CreateResponse(response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to save response")
		return
	}

	// Enqueue churn calculation
	workerClient := getWorkerClient()
	if workerClient != nil {
		workerClient.EnqueueChurnCalculation(req.UserToken, user.CompanyID)
	}

	utils.WriteSuccess(w, map[string]string{"message": "Response submitted successfully"})
}

func parseUint(s string) uint64 {
	val, _ := strconv.ParseUint(s, 10, 32)
	return val
}

func GetSurveyResponses(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "surveyID")
	if surveyID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Survey ID required")
		return
	}

	surveyRepo := repository.NewSurveyRepository()
	responses, err := surveyRepo.GetResponsesBySurveyID(uint(parseUint(surveyID)))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch responses")
		return
	}

	// Enrich with user details
	userRepo := repository.NewUserRepository()
	type ResponseWithUser struct {
		ID        uint     `json:"id"`
		UserID    uint     `json:"user_id"`
		UserName  string   `json:"user_name"`
		UserEmail string   `json:"user_email"`
		Responses []string `json:"responses"`
		Sentiment float64  `json:"sentiment"`
		CreatedAt string   `json:"created_at"`
	}

	result := []ResponseWithUser{}
	for _, resp := range responses {
		user, _ := userRepo.GetByID(resp.UserID)
		result = append(result, ResponseWithUser{
			ID:        resp.ID,
			UserID:    resp.UserID,
			UserName:  user.Name,
			UserEmail: user.Email,
			Responses: resp.Responses,
			Sentiment: resp.Sentiment,
			CreatedAt: resp.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	utils.WriteSuccess(w, result)
}