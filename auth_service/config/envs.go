package config

import (
	"github.com/spf13/viper"
)

type EnvsSchema struct {
	HOST                  string
	PORT                  int
	LOG_LEVEL             string
	JWT_SECRET_KEY        string
	JWT_EXP_HOURS         int
	JWT_REFRESH_EXP_HOURS int
	POSTGRESQL_HOST       string
	POSTGRESQL_PORT       int
	POSTGRESQL_USER       string
	POSTGRESQL_PASSWORD   string
	POSTGRESQL_DB         string
}

var Envs *EnvsSchema

func envInitiator() {
	Envs = &EnvsSchema{
		HOST:                  viper.GetString("HOST"),
		PORT:                  viper.GetInt("PORT"),
		LOG_LEVEL:             viper.GetString("LOG_LEVEL"),
		JWT_SECRET_KEY:        viper.GetString("JWT_SECRET_KEY"),
		JWT_EXP_HOURS:         viper.GetInt("JWT_EXP_HOURS"),
		JWT_REFRESH_EXP_HOURS: viper.GetInt("JWT_REFRESH_EXP_HOURS"),
		POSTGRESQL_HOST:       viper.GetString("POSTGRESQL_HOST"),
		POSTGRESQL_PORT:       viper.GetInt("POSTGRESQL_PORT"),
		POSTGRESQL_USER:       viper.GetString("POSTGRESQL_USER"),
		POSTGRESQL_PASSWORD:   viper.GetString("POSTGRESQL_PASSWORD"),
		POSTGRESQL_DB:         viper.GetString("POSTGRESQL_DB"),
	}
}

func InitEnv(filepath string) {
	viper.SetConfigType("env")
	viper.SetConfigFile(filepath)
	if err := viper.ReadInConfig(); err != nil {
		logger.Warningf("error loading environment variables from %s: %w", filepath, err)
	}
	viper.AutomaticEnv()
	envInitiator()
}
