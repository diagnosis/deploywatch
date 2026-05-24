package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Database DatabaseConfig
	App      AppConfig
	JWT      JWTConfig
	GitHub   GitHubConfig
}

type DatabaseConfig struct {
	Dsn               string
	MinConns          int32
	MaxConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnectTimeout    time.Duration
}

type AppConfig struct {
	Env         string
	BaseURL     string
	FrontendURL string
	Port        string
	Name        string
	Platform    string
}

type JWTConfig struct {
	AccessSecret       string
	RefreshSecret      string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
	Audience           string
}
type GitHubConfig struct {
	GitHubAppID             string
	GitHubAppPrivateKey     []byte
	GitHubWebhookSecret     string
	GitHubOauthClientID     string
	GitHubOauthClientSecret string
}

func Load() (*Config, error) {
	dsn := getEnv("DATABASE_URL", "")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	accessSecret := getEnv("JWT_ACCESS_SECRET", "")
	if accessSecret == "" {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET is required")
	}

	refreshSecret := getEnv("JWT_REFRESH_SECRET", "")
	if refreshSecret == "" {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET is required")
	}

	appID := getEnv("GITHUB_APP_ID", "")
	if appID == "" {
		return nil, fmt.Errorf("GITHUB_APP_ID is required")
	}

	privateKeyPath := getEnv("GITHUB_APP_PRIVATE_KEY_PATH", "")
	if privateKeyPath == "" {
		return nil, fmt.Errorf("GITHUB_APP_PRIVATE_KEY_PATH is required")
	}
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	webhookSecret := getEnv("GITHUB_WEBHOOK_SECRET", "")
	if webhookSecret == "" {
		return nil, fmt.Errorf("GITHUB_WEBHOOK_SECRET is required")
	}

	oauthClientID := getEnv("GITHUB_OAUTH_CLIENT_ID", "")
	if oauthClientID == "" {
		return nil, fmt.Errorf("GITHUB_OAUTH_CLIENT_ID is required")
	}

	oauthClientSecret := getEnv("GITHUB_OAUTH_CLIENT_SECRET", "")
	if oauthClientSecret == "" {
		return nil, fmt.Errorf("GITHUB_OAUTH_CLIENT_SECRET is required")
	}

	return &Config{
		Database: DatabaseConfig{
			Dsn:               dsn,
			MinConns:          getEnvInt32("DB_MIN_CONNS", 2),
			MaxConns:          getEnvInt32("DB_MAX_CONNS", 10),
			MaxConnLifetime:   getEnvDuration("DB_MAX_CONN_LIFETIME", 1*time.Hour),
			MaxConnIdleTime:   getEnvDuration("DB_MAX_CONN_IDLE_TIME", 30*time.Minute),
			HealthCheckPeriod: getEnvDuration("DB_HEALTH_CHECK_PERIOD", 1*time.Minute),
			ConnectTimeout:    getEnvDuration("DB_CONNECT_TIMEOUT", 5*time.Second),
		},
		App: AppConfig{
			Env:         getEnv("APP_ENV", "dev"),
			BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
			Port:        getEnv("PORT", "8080"),
			Name:        getEnv("APP_NAME", "deploywatch"),
			Platform:    getEnv("PLATFORM", "web"),
		},
		JWT: JWTConfig{
			AccessSecret:       accessSecret,
			RefreshSecret:      refreshSecret,
			AccessTokenExpiry:  getEnvDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			Issuer:             getEnv("JWT_ISSUER", "deploywatch"),
			Audience:           getEnv("JWT_AUDIENCE", "deploywatch-api"),
		},
		GitHub: GitHubConfig{
			GitHubAppID:             appID,
			GitHubAppPrivateKey:     privateKey,
			GitHubWebhookSecret:     webhookSecret,
			GitHubOauthClientID:     oauthClientID,
			GitHubOauthClientSecret: oauthClientSecret,
		},
	}, nil
}

// helper functions
func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getEnvInt32(key string, def int32) int32 {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return def
	}
	return int32(i)
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return def
	}
	return d
}

func getEnvBool(key string, def bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}
	return b
}

func getEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}
