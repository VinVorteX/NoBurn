package employee

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/VinVorteX/NoBurn/internal/models"
	"github.com/VinVorteX/NoBurn/internal/repository"
	"github.com/VinVorteX/NoBurn/internal/utils"
)

type AddEmployeeRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// AddEmployee - Admin adds single employee
func AddEmployee(w http.ResponseWriter, r *http.Request) {
	var req AddEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	userID := r.Context().Value("userID").(uint)
	userRepo := repository.NewUserRepository()
	admin, err := userRepo.GetByID(userID)
	if err != nil || admin.Role != "hr_admin" {
		utils.WriteError(w, http.StatusForbidden, "Only admins can add employees")
		return
	}

	tempPassword, _ := utils.HashPassword("temp123")
	
	employee := &models.User{
		Email:     req.Email,
		Password:  tempPassword,
		Name:      req.Name,
		Role:      "employee",
		CompanyID: admin.CompanyID,
	}

	if err := userRepo.Create(employee); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create employee")
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.WriteSuccess(w, map[string]interface{}{
		"message": "Employee added successfully",
		"employee": employee,
	})
}

// BulkUploadEmployees - Admin uploads CSV with employees
func BulkUploadEmployees(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	userRepo := repository.NewUserRepository()
	admin, err := userRepo.GetByID(userID)
	if err != nil || admin.Role != "hr_admin" {
		utils.WriteError(w, http.StatusForbidden, "Only admins can upload employees")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid CSV file")
		return
	}

	if len(records) < 2 {
		utils.WriteError(w, http.StatusBadRequest, "CSV must have header and at least one employee")
		return
	}

	tempPassword, _ := utils.HashPassword("temp123")
	successCount := 0
	failedEmails := []string{}

	for i, record := range records[1:] {
		if len(record) < 2 {
			continue
		}

		email := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])

		if email == "" || name == "" {
			failedEmails = append(failedEmails, email)
			continue
		}

		employee := &models.User{
			Email:     email,
			Password:  tempPassword,
			Name:      name,
			Role:      "employee",
			CompanyID: admin.CompanyID,
		}

		if err := userRepo.Create(employee); err != nil {
			failedEmails = append(failedEmails, email)
		} else {
			successCount++
		}

		if i >= 999 {
			break
		}
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"message":       "Bulk upload completed",
		"success_count": successCount,
		"failed_count":  len(failedEmails),
		"failed_emails": failedEmails,
	})
}

// GetEmployees - List all employees in company
func GetEmployees(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not found")
		return
	}

	employees, err := userRepo.GetByCompanyID(user.CompanyID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch employees")
		return
	}

	utils.WriteSuccess(w, employees)
}
