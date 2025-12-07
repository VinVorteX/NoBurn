package models

import (
	"time"
	"gorm.io/gorm"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StringArray []string

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type Survey struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CompanyID   uint           `json:"company_id"`
	Company     Company        `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Title       string         `json:"title" gorm:"not null"`
	Questions   StringArray    `json:"questions" gorm:"type:jsonb"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	Responses   []SurveyResponse `json:"responses,omitempty" gorm:"foreignKey:SurveyID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type SurveyResponse struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	SurveyID   uint           `json:"survey_id"`
	Survey     Survey         `json:"survey,omitempty" gorm:"foreignKey:SurveyID"`
	UserID     uint           `json:"user_id"`
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Responses  StringArray    `json:"responses" gorm:"type:jsonb"`
	Sentiment  float64        `json:"sentiment" gorm:"default:0"` // -1 to 1
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type AttritionRisk struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id"`
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	RiskScore  float64        `json:"risk_score" gorm:"default:0"` // 0 to 1
	Factors    StringArray    `json:"factors" gorm:"type:jsonb"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}