package logger

import (
	"github.com/VinVorteX/NoBurn/internal/config"
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(){
	var err error

	// Default to development if config not loaded yet
	env := "development"
	if config.AppConfig != nil {
		env = config.AppConfig.Env
	}

	if env == "development"{
		Log, err = zap.NewDevelopment()
	} else {
		Log, err = zap.NewProduction()
	}
	
	if err != nil{
		panic(err)
	}

	Log.Info("Logger Initialized", zap.String("env", env))
}

