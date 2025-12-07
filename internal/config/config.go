package config

import (
	"fmt"

	"github.com/spf13/viper"

	"time"
) 

type Config struct {
	Port string `mapstructure:"PORT"`
	Env string `mapstructure:"ENV"`
	DbUrl string `mapstructure:"DB_URL"`
	JwtSecret string `mapstructure:"JWT_SECRET"`
	JwtExpiresIn time.Duration `mapstructure:"JWT_EXPIRES_IN"`
	HuggingFaceToken string `mapstructure:"HUGGING_FACE_TOKEN"`
	RedisURL string `mapstructure:"REDIS_URL"`
	SlackWebhookURL string `mapstructure:"SLACK_WEBHOOK_URL"`
	SMTPHost string `mapstructure:"SMTP_HOST"`
	SMTPPort int `mapstructure:"SMTP_PORT"`
	SMTPUser string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	AlertEmail string `mapstructure:"ALERT_EMAIL"`
}

var AppConfig *Config

func LoadConfig() error{
	// setting default values to initial configurations
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("JWT_EXPIRES_IN", "24h")

	viper.SetConfigName(".env") // name of config file
	viper.SetConfigType("env") // type of config file
	viper.AddConfigPath(".") // adding config path
	viper.AddConfigPath("./config")

	// Explicitly bind environment variables
	viper.BindEnv("PORT")
	viper.BindEnv("ENV")
	viper.BindEnv("DB_URL")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("JWT_EXPIRES_IN")
	viper.BindEnv("HUGGING_FACE_TOKEN")
	viper.BindEnv("REDIS_URL")
	viper.BindEnv("SLACK_WEBHOOK_URL")
	viper.BindEnv("SMTP_HOST")
	viper.BindEnv("SMTP_PORT")
	viper.BindEnv("SMTP_USER")
	viper.BindEnv("SMTP_PASSWORD")
	viper.BindEnv("ALERT_EMAIL")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil{
		if _, ok:= err.(viper.ConfigFileNotFoundError); ok{
			// Config file not found; using environment variables only
		} else {
			return err
		}
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil{
		return err
	}

	if AppConfig.DbUrl == ""{
		return fmt.Errorf("DB_URL field is required")
	}

	if AppConfig.JwtSecret == ""{
		return fmt.Errorf("JWT_SECRET field is required")
	}

	if AppConfig.Env == "development"{
		viper.Debug()
	}

	return nil
}