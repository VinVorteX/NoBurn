package services

import (
	"github.com/VinVorteX/NoBurn/internal/models"
	"github.com/VinVorteX/NoBurn/internal/repository"
	"github.com/VinVorteX/NoBurn/internal/services/sentiment"
)

type AnalyticsService struct {
	userRepo      *repository.UserRepository
	surveyRepo    *repository.SurveyRepository
	mlService     *sentiment.MLService
	churnPredictor *sentiment.ChurnPredictor
}

func NewAnalyticsService(userRepo *repository.UserRepository, surveyRepo *repository.SurveyRepository, mlAPIURL string) *AnalyticsService {
	return &AnalyticsService{
		userRepo:      userRepo,
		surveyRepo:    surveyRepo,
		mlService:     sentiment.NewMLService(mlAPIURL),
		churnPredictor: sentiment.NewChurnPredictor(),
	}
}

func (s *AnalyticsService) AnalyzeUserChurnRisk(userID uint) (*models.AttritionRisk, error) {
	// Get user responses
	responses, err := s.surveyRepo.GetResponsesByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Calculate features
	features := s.calculateChurnFeatures(responses)
	
	// Predict risk
	riskScore := s.churnPredictor.PredictChurnRisk(features)
	factors := s.churnPredictor.GetRiskFactors(features)

	return &models.AttritionRisk{
		UserID:    userID,
		RiskScore: riskScore,
		Factors:   models.StringArray(factors),
	}, nil
}

func (s *AnalyticsService) GenerateRetentionSuggestions(userID uint, language string) ([]RetentionSuggestion, error) {
	responses, err := s.surveyRepo.GetResponsesByUserID(userID)
	if err != nil {
		return nil, err
	}

	features := s.calculateChurnFeatures(responses)
	riskFactors := s.churnPredictor.GetRiskFactors(features)
	
	// Use AI service for intelligent suggestions
	aiService := NewAIService("") // TODO: Get HF token from config
	return aiService.GenerateRetentionSuggestions(riskFactors, language)
}

func (s *AnalyticsService) ProcessSurveyResponse(response *models.SurveyResponse, language string) error {
	// Analyze sentiment using ML
	totalText := ""
	for _, resp := range response.Responses {
		totalText += resp + " "
	}

	sentimentScore, err := s.mlService.AnalyzeSentiment(totalText, language)
	if err != nil {
		// Fallback to rule-based
		sentimentScore = sentiment.AnalyzeSentiment(totalText, language)
	}

	response.Sentiment = sentimentScore
	return s.surveyRepo.CreateResponse(response)
}

func (s *AnalyticsService) calculateChurnFeatures(responses []models.SurveyResponse) sentiment.ChurnFeatures {
	if len(responses) == 0 {
		return sentiment.ChurnFeatures{}
	}

	totalSentiment := 0.0
	negativeCount := 0

	for _, resp := range responses {
		totalSentiment += resp.Sentiment
		if resp.Sentiment < -0.1 {
			negativeCount++
		}
	}

	return sentiment.ChurnFeatures{
		AvgSentiment:      totalSentiment / float64(len(responses)),
		ResponseRate:      1.0, // TODO: Calculate actual rate
		DaysInactive:      0,   // TODO: Calculate from last activity
		NegativeResponses: negativeCount,
		TotalResponses:    len(responses),
		LastLoginDays:     0, // TODO: Calculate from user activity
	}
}