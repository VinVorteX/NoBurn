package auth

import (
	"encoding/json"
	"net/http"

	"github.com/VinVorteX/NoBurn/internal/models"
	"github.com/VinVorteX/NoBurn/internal/repository"
	"github.com/VinVorteX/NoBurn/internal/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	CompanyName string `json:"company_name"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByEmail(req.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email)
	
	utils.WriteSuccess(w, map[string]interface{}{
		"token": token,
		"user": user,
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create company
	companyRepo := repository.NewCompanyRepository()
	company := &models.Company{
		Name:     req.CompanyName,
		Plan:     "free",
		Language: "en",
	}
	if err := companyRepo.Create(company); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create company")
		return
	}

	// Create user
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		Name:      req.Name,
		Role:      "hr_admin",
		CompanyID: company.ID,
	}
	if err := userRepo.Create(user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email)
	
	w.WriteHeader(http.StatusCreated)
	utils.WriteSuccess(w, map[string]interface{}{
		"token": token,
		"message": "User registered successfully",
	})
}