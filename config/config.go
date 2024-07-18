package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort string
	LogLevel   string
	LogFile    string
}

func Load() (*Config, error) {
	viper.SetDefault("SERVER_PORT", ":8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FILE", "server.log")

	viper.AutomaticEnv()

	return &Config{
		ServerPort: viper.GetString("SERVER_PORT"),
		LogLevel:   viper.GetString("LOG_LEVEL"),
		LogFile:    viper.GetString("LOG_FILE"),
	}, nil
}
