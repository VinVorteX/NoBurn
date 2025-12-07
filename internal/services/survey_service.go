package services

import (
	"strings"

	"github.com/VinVorteX/NoBurn/internal/repository"
	"github.com/VinVorteX/NoBurn/internal/services/sentiment"
)

type SurveyService struct {
	surveyRepo *repository.SurveyRepository
	mlService  *sentiment.MLService
}

func NewSurveyService(hfToken string) *SurveyService {
	return &SurveyService{
		surveyRepo: repository.NewSurveyRepository(),
		mlService:  sentiment.NewMLService(hfToken),
	}
}

func (s *SurveyService) ProcessSurveyResponse(responseID uint, language string) error {
	// Get response from database
	responses, err := s.surveyRepo.GetResponsesByUserID(responseID)
	if err != nil || len(responses) == 0 {
		return err
	}

	response := &responses[len(responses)-1] // Get latest response

	// Combine all response texts
	combinedText := strings.Join(response.Responses, " ")

	// Analyze sentiment using IndicBERT
	sentimentScore, err := s.mlService.AnalyzeSentiment(combinedText, language)
	if err != nil {
		// Fallback to rule-based
		sentimentScore = sentiment.AnalyzeSentiment(combinedText, language)
	}

	// Update response with sentiment score
	response.Sentiment = sentimentScore
	
	// Save updated response
	return s.surveyRepo.CreateResponse(response)
}