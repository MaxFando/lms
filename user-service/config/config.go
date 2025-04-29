package config

import (
	"time"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string
	Env         string
	LogLevel    string
	GRPCPort    string
	DatabaseDSN string
	TracerDSN   string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func Load() *Config {
	viper.SetConfigName("env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	return &Config{
		ServiceName: viper.GetString("SERVICE_NAME"),
		Env:         viper.GetString("APP_ENV"),
		LogLevel:    viper.GetString("LOG_LEVEL"),
		GRPCPort:    viper.GetString("GRPC_PORT"),
		DatabaseDSN: viper.GetString("DATABASE_DSN"),
		TracerDSN:   viper.GetString("TRACER_DSN"),
		JWTSecret:       viper.GetString("JWT_SECRET"),
		AccessTokenTTL:  viper.GetDuration("ACCESS_TOKEN_TTL"),
		RefreshTokenTTL: viper.GetDuration("REFRESH_TOKEN_TTL"),
	}
}
