package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Server
	ServerPort string `mapstructure:"SERVER_PORT"`
	ServerMode string `mapstructure:"SERVER_MODE"`
	AppName    string `mapstructure:"APP_NAME"`
	AppEnv     string `mapstructure:"APP_ENV"`

	// Database
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBTimezone string `mapstructure:"DB_TIMEZONE"`

	// Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	// JWT
	JWTSecret           string        `mapstructure:"JWT_SECRET"`
	JWTAccessExpiry     time.Duration `mapstructure:"JWT_ACCESS_TOKEN_EXPIRY"`
	JWTRefreshExpiry    time.Duration `mapstructure:"JWT_REFRESH_TOKEN_EXPIRY"`

	// OAuth Google
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`

	// OAuth Facebook
	FacebookAppID     string `mapstructure:"FACEBOOK_APP_ID"`
	FacebookAppSecret string `mapstructure:"FACEBOOK_APP_SECRET"`
	FacebookRedirectURL string `mapstructure:"FACEBOOK_REDIRECT_URL"`

	// Midtrans
	MidtransServerKey   string `mapstructure:"MIDTRANS_SERVER_KEY"`
	MidtransClientKey   string `mapstructure:"MIDTRANS_CLIENT_KEY"`
	MidtransEnvironment string `mapstructure:"MIDTRANS_ENVIRONMENT"`

	// RajaOngkir
	RajaOngkirAPIKey string `mapstructure:"RAJAONGKIR_API_KEY"`
	RajaOngkirBaseURL string `mapstructure:"RAJAONGKIR_BASE_URL"`

	// TikTok API
	TikTokAppID     string `mapstructure:"TIKTOK_APP_ID"`
	TikTokAppSecret string `mapstructure:"TIKTOK_APP_SECRET"`
	TikTokAccessToken string `mapstructure:"TIKTOK_ACCESS_TOKEN"`

	// SerpAPI (Google Trends)
	SerpAPIKey string `mapstructure:"SERPAPI_KEY"`

	// Email
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`

	// CORS
	CORSAllowedOrigins  string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CORSAllowedMethods  string `mapstructure:"CORS_ALLOWED_METHODS"`
	CORSAllowedHeaders  string `mapstructure:"CORS_ALLOWED_HEADERS"`

	// Rate Limiting
	RateLimitRequests   int `mapstructure:"RATE_LIMIT_REQUESTS"`
	RateLimitDuration   int `mapstructure:"RATE_LIMIT_DURATION"`

	// Upload
	UploadMaxSize      int    `mapstructure:"UPLOAD_MAX_SIZE"`
	UploadAllowedTypes string `mapstructure:"UPLOAD_ALLOWED_TYPES"`
	UploadPath         string `mapstructure:"UPLOAD_PATH"`

	// Logging
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	LogFilePath string `mapstructure:"LOG_FILE_PATH"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_MODE", "debug")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("JWT_ACCESS_TOKEN_EXPIRY", 3600)
	viper.SetDefault("JWT_REFRESH_TOKEN_EXPIRY", 604800)
	viper.SetDefault("RATE_LIMIT_REQUESTS", 100)
	viper.SetDefault("RATE_LIMIT_DURATION", 60)
	viper.SetDefault("UPLOAD_MAX_SIZE", 5242880)
	viper.SetDefault("UPLOAD_PATH", "./uploads")
	viper.SetDefault("LOG_LEVEL", "debug")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Convert timeouts from seconds to time.Duration
	config.JWTAccessExpiry = time.Duration(config.JWTAccessExpiry) * time.Second
	config.JWTRefreshExpiry = time.Duration(config.JWTRefreshExpiry) * time.Second

	return &config, nil
}

// GetDSN returns MySQL DSN string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBTimezone)
}