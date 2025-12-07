package database

import (
	"github.com/VinVorteX/NoBurn/internal/config"
	"github.com/VinVorteX/NoBurn/internal/models"
	"github.com/VinVorteX/NoBurn/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	var err error

	DB, err = gorm.Open(postgres.Open(config.AppConfig.DbUrl), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	// AutoMigrate for development only
	if config.AppConfig.Env == "development" {
		if err := DB.AutoMigrate(
			&models.Company{},
			&models.User{},
			&models.Survey{},
			&models.SurveyResponse{},
			&models.AttritionRisk{},
		); err != nil {
			logger.Log.Warn("AutoMigrate failed, use migrations instead")
		}
		logger.Log.Info("AutoMigrate completed (dev mode)")
	}
	// Production: Use 'make migrate-up' instead

	logger.Log.Info("Connected to database and migrated models")
	return nil
}

func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}