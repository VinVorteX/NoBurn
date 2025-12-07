package repository

import (
	"github.com/VinVorteX/NoBurn/internal/database"
	"github.com/VinVorteX/NoBurn/internal/models"
)

type AttritionRepository struct{}

func NewAttritionRepository() *AttritionRepository {
	return &AttritionRepository{}
}

func (r *AttritionRepository) Create(risk *models.AttritionRisk) error {
	return database.DB.Create(risk).Error
}

func (r *AttritionRepository) GetByUserID(userID uint) (*models.AttritionRisk, error) {
	var risk models.AttritionRisk
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").First(&risk).Error
	return &risk, err
}

func (r *AttritionRepository) GetHighRiskUsers(companyID uint, threshold float64) ([]models.AttritionRisk, error) {
	var risks []models.AttritionRisk
	err := database.DB.
		Joins("JOIN users ON users.id = attrition_risks.user_id").
		Where("users.company_id = ? AND attrition_risks.risk_score >= ?", companyID, threshold).
		Preload("User").
		Order("attrition_risks.risk_score DESC").
		Find(&risks).Error
	return risks, err
}