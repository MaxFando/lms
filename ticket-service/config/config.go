package config

import "github.com/spf13/viper"

type Config struct {
	ServiceName         string
	Env                 string
	LogLevel            string
	GRPCPort            string
	DatabaseDSN         string
	TracerDSN           string
	RedisDSN            string
	RedisDrawChannel    string
	RedisInvoiceChannel string
}

func Load() *Config {
	viper.SetConfigName("env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	return &Config{
		ServiceName:         viper.GetString("SERVICE_NAME"),
		Env:                 viper.GetString("APP_ENV"),
		LogLevel:            viper.GetString("LOG_LEVEL"),
		GRPCPort:            viper.GetString("GRPC_PORT"),
		DatabaseDSN:         viper.GetString("DATABASE_DSN"),
		TracerDSN:           viper.GetString("TRACER_DSN"),
		RedisDSN:            viper.GetString("REDIS_DSN"),
		RedisDrawChannel:    viper.GetString("REDIS_DRAW_CHANNEL"),
		RedisInvoiceChannel: viper.GetString("REDIS_INVOICE_CHANNEL"),
	}
}
