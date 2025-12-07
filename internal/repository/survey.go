package repository

import (
	"github.com/VinVorteX/NoBurn/internal/database"
	"github.com/VinVorteX/NoBurn/internal/models"
)

type SurveyRepository struct{}

func NewSurveyRepository() *SurveyRepository {
	return &SurveyRepository{}
}

func (r *SurveyRepository) Create(survey *models.Survey) error {
	return database.DB.Create(survey).Error
}

func (r *SurveyRepository) GetByID(id uint) (*models.Survey, error) {
	var survey models.Survey
	err := database.DB.First(&survey, id).Error
	return &survey, err
}

func (r *SurveyRepository) GetByCompanyID(companyID uint) ([]models.Survey, error) {
	var surveys []models.Survey
	err := database.DB.Where("company_id = ? AND is_active = ?", companyID, true).Find(&surveys).Error
	return surveys, err
}

func (r *SurveyRepository) CreateResponse(response *models.SurveyResponse) error {
	return database.DB.Create(response).Error
}

func (r *SurveyRepository) GetResponsesByUserID(userID uint) ([]models.SurveyResponse, error) {
	var responses []models.SurveyResponse
	err := database.DB.Preload("Survey").Where("user_id = ?", userID).Find(&responses).Error
	return responses, err
}