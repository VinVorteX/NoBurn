package settings

import (
	"encoding/json"
	"net/http"

	"github.com/VinVorteX/NoBurn/internal/repository"
	"github.com/VinVorteX/NoBurn/internal/utils"
)

type UpdateSMTPRequest struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUser     string `json:"smtp_user"`
	SMTPPassword string `json:"smtp_password"`
}

func UpdateSMTPSettings(w http.ResponseWriter, r *http.Request) {
	var req UpdateSMTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	userID := r.Context().Value("userID").(uint)
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(userID)
	if err != nil || user.Role != "hr_admin" {
		utils.WriteError(w, http.StatusForbidden, "Only admins can update SMTP settings")
		return
	}

	companyRepo := repository.NewCompanyRepository()
	company, err := companyRepo.GetByID(user.CompanyID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Company not found")
		return
	}

	company.SMTPHost = req.SMTPHost
	company.SMTPPort = req.SMTPPort
	company.SMTPUser = req.SMTPUser
	if req.SMTPPassword != "" {
		company.SMTPPassword = req.SMTPPassword
	}

	if err := companyRepo.Update(company); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update settings")
		return
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"message": "SMTP settings updated successfully",
		"smtp_configured": company.SMTPUser != "",
	})
}

func GetSMTPSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not found")
		return
	}

	companyRepo := repository.NewCompanyRepository()
	company, err := companyRepo.GetByID(user.CompanyID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Company not found")
		return
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"smtp_host": company.SMTPHost,
		"smtp_port": company.SMTPPort,
		"smtp_user": company.SMTPUser,
		"smtp_configured": company.SMTPUser != "",
	})
}
