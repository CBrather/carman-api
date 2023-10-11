package main

import (
	"github.com/vrischmann/envconfig"
	"go.uber.org/zap"

	"github.com/CBrather/carman-api/pkg/log"
)

var env *EnvConfig = nil

func getEnvironment() *EnvConfig {
	if env == nil {
		initConfig()
	}

	return env
}

func initConfig() {
	logger, err := log.GetLoggerWithLevel("info")
	if err != nil {
		zap.L().Fatal("Failed to get a new logger", zap.Error(err))
	}

	env = new(EnvConfig)

	if err = envconfig.Init(env); err != nil {
		logger.Fatal("Failed initializing the app config", zap.Error(err))
	}
}

type EnvConfig struct {
	Auth struct {
		Domain   string `envconfig:"AUTH_DOMAIN"`
		Audience string `envconfig:"AUTH_AUDIENCE"`
	}

	DBConnectionString string `envconfig:"DB_CONNECTION_STRING"`

	LogLevel string `envconfig:"LOGLEVEL"`
}
