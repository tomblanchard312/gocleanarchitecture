package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort       string
	DBPath           string
	LogLevel         string
	LogFile          string
	DBType           string // "sqlite", "supabase", "inmemory"
	SupabaseURL      string
	SupabaseKey      string
	SupabaseUser     string // Database username for Supabase
	SupabasePass     string // Database password for Supabase
	JWTSecret        string
	JWTTokenDuration time.Duration
}

func Load() (*Config, error) {
	viper.SetDefault("SERVER_PORT", ":8080")
	viper.SetDefault("DB_PATH", "./blog.db")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FILE", "server.log")
	viper.SetDefault("DB_TYPE", "sqlite")
	viper.SetDefault("JWT_SECRET", "your-secret-key-change-this-in-production")
	viper.SetDefault("JWT_TOKEN_DURATION_HOURS", 24)

	viper.AutomaticEnv()

	return &Config{
		ServerPort:       viper.GetString("SERVER_PORT"),
		DBPath:           viper.GetString("DB_PATH"),
		LogLevel:         viper.GetString("LOG_LEVEL"),
		LogFile:          viper.GetString("LOG_FILE"),
		DBType:           viper.GetString("DB_TYPE"),
		SupabaseURL:      viper.GetString("SUPABASE_URL"),
		SupabaseKey:      viper.GetString("SUPABASE_KEY"),
		SupabaseUser:     viper.GetString("SUPABASE_USER"),
		SupabasePass:     viper.GetString("SUPABASE_PASS"),
		JWTSecret:        viper.GetString("JWT_SECRET"),
		JWTTokenDuration: time.Duration(viper.GetInt("JWT_TOKEN_DURATION_HOURS")) * time.Hour,
	}, nil
}
