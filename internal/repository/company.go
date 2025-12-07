package repository

import (
	"github.com/VinVorteX/NoBurn/internal/database"
	"github.com/VinVorteX/NoBurn/internal/models"
)

type CompanyRepository struct{}

func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{}
}

func (r *CompanyRepository) Create(company *models.Company) error {
	return database.DB.Create(company).Error
}

func (r *CompanyRepository) GetByID(id uint) (*models.Company, error) {
	var company models.Company
	err := database.DB.First(&company, id).Error
	return &company, err
}

func (r *CompanyRepository) Update(company *models.Company) error {
	return database.DB.Save(company).Error
}