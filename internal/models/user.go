package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Name      string         `json:"name" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:employee"` // hr_admin, employee
	CompanyID uint           `json:"company_id"`
	Company   Company        `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Company struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null"`
	Plan         string         `json:"plan" gorm:"default:free"` // free, premium
	Language     string         `json:"language" gorm:"default:en"` // en, hi, ta
	SMTPHost     string         `json:"smtp_host,omitempty"`
	SMTPPort     int            `json:"smtp_port,omitempty" gorm:"default:587"`
	SMTPUser     string         `json:"smtp_user,omitempty"`
	SMTPPassword string         `json:"-" gorm:"column:smtp_password"`
	Users        []User         `json:"users,omitempty" gorm:"foreignKey:CompanyID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}