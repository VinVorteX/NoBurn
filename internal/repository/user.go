package repository

import (
	"github.com/VinVorteX/NoBurn/internal/database"
	"github.com/VinVorteX/NoBurn/internal/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Company").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByCompanyID(companyID uint) ([]models.User, error) {
	var users []models.User
	err := database.DB.Where("company_id = ?", companyID).Find(&users).Error
	return users, err
}