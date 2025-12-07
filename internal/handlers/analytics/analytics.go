package analytics

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/VinVorteX/NoBurn/internal/models"
	"github.com/VinVorteX/NoBurn/internal/repository"
	"github.com/VinVorteX/NoBurn/internal/services"
	"github.com/VinVorteX/NoBurn/internal/utils"
)

type DashboardData struct {
	TotalEmployees   int                    `json:"total_employees"`
	AtRiskEmployees  int                    `json:"at_risk_employees"`
	AvgSentiment     float64                `json:"avg_sentiment"`
	ChurnRate        float64                `json:"churn_rate"`
	TopRiskFactors   []string               `json:"top_risk_factors"`
	AttritionRisks   []models.AttritionRisk `json:"attrition_risks"`
}

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not found")
		return
	}

	// Get all company employees
	employees, _ := userRepo.GetByCompanyID(user.CompanyID)
	totalEmployees := len(employees)

	// Get high-risk employees
	attritionRepo := repository.NewAttritionRepository()
	atRiskUsers, _ := attritionRepo.GetHighRiskUsers(user.CompanyID, 0.7)
	atRiskCount := len(atRiskUsers)

	// Calculate average sentiment
	surveyRepo := repository.NewSurveyRepository()
	avgSentiment := 0.0
	totalResponses := 0
	for _, emp := range employees {
		responses, _ := surveyRepo.GetResponsesByUserID(emp.ID)
		for _, resp := range responses {
			avgSentiment += resp.Sentiment
			totalResponses++
		}
	}
	if totalResponses > 0 {
		avgSentiment /= float64(totalResponses)
	}

	churnRate := 0.0
	if totalEmployees > 0 {
		churnRate = (float64(atRiskCount) / float64(totalEmployees)) * 100
	}

	dashboard := DashboardData{
		TotalEmployees:  totalEmployees,
		AtRiskEmployees: atRiskCount,
		AvgSentiment:    avgSentiment,
		ChurnRate:       churnRate,
		TopRiskFactors:  []string{"Work-life balance", "Career growth", "Compensation"},
		AttritionRisks:  atRiskUsers,
	}

	utils.WriteSuccess(w, dashboard)
}

func GetAttritionRisks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not found")
		return
	}

	attritionRepo := repository.NewAttritionRepository()
	risks, err := attritionRepo.GetHighRiskUsers(user.CompanyID, 0.5)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch risks")
		return
	}

	utils.WriteSuccess(w, risks)
}

func GetRetentionSuggestions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	language := r.URL.Query().Get("language")
	if language == "" {
		language = "en"
	}
	
	// TODO: Convert userID string to uint and use actual service
	// For now, return AI-generated suggestions
	aiService := services.NewAIService("") // TODO: Get from config
	riskFactors := []string{"Low sentiment scores", "Poor survey participation"}
	
	suggestions, err := aiService.GenerateRetentionSuggestions(riskFactors, language)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to generate suggestions")
		return
	}

	response := map[string]interface{}{
		"user_id":     userID,
		"suggestions": suggestions,
		"language":    language,
	}

	utils.WriteSuccess(w, response)
}