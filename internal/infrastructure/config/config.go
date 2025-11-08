package config

import (
	"fmt"
	"os"
	"time"

	"github.com/crossmint/megaverse-challenge/pkg/retry"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	API       APIConfig       `mapstructure:"api"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Execution ExecutionConfig `mapstructure:"execution"`
}

// APIConfig contains API-related configuration
type APIConfig struct {
	BaseURL         string          `mapstructure:"base_url"`
	CandidateID     string          `mapstructure:"candidate_id"`
	Timeout         time.Duration   `mapstructure:"timeout"`
	RetryConfig     RetryConfig     `mapstructure:"retry"`
	RateLimitConfig RateLimitConfig `mapstructure:"rate_limit"`
}

// RetryConfig contains retry-related configuration
type RetryConfig struct {
	MaxAttempts  int           `mapstructure:"max_attempts"`
	InitialDelay time.Duration `mapstructure:"initial_delay"`
	MaxDelay     time.Duration `mapstructure:"max_delay"`
	Multiplier   float64       `mapstructure:"multiplier"`
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond float64 `mapstructure:"requests_per_second"`
}

// LoggingConfig contains logging-related configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// ExecutionConfig contains execution-related configuration
type ExecutionConfig struct {
	MaxWorkers int           `mapstructure:"max_workers"`
	BatchSize  int           `mapstructure:"batch_size"`
	Timeout    time.Duration `mapstructure:"timeout"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			BaseURL:     "https://challenge.crossmint.io/api",
			CandidateID: "",
			Timeout:     30 * time.Second,
			RetryConfig: RetryConfig{
				MaxAttempts:  6,
				InitialDelay: 1 * time.Second,
				MaxDelay:     30 * time.Second,
				Multiplier:   2.0,
			},
			RateLimitConfig: RateLimitConfig{
				RequestsPerSecond: 2.0,
			},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Execution: ExecutionConfig{
			MaxWorkers: 5,
			BatchSize:  5,
			Timeout:    5 * time.Minute,
		},
	}
}

// Load loads the configuration from various sources
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Set up viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.crossmint")

	// Environment variables
	viper.SetEnvPrefix("CROSSMINT")
	viper.AutomaticEnv()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is not an error, we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal configuration
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Override with environment variables if set
	if candidateID := os.Getenv("CROSSMINT_CANDIDATE_ID"); candidateID != "" {
		cfg.API.CandidateID = candidateID
	}

	if baseURL := os.Getenv("CROSSMINT_API_URL"); baseURL != "" {
		cfg.API.BaseURL = baseURL
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// LoadFromFile loads configuration from a specific file
func LoadFromFile(path string) (*Config, error) {
	cfg := DefaultConfig()

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Override with environment variables
	if candidateID := os.Getenv("CROSSMINT_CANDIDATE_ID"); candidateID != "" {
		cfg.API.CandidateID = candidateID
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.API.BaseURL == "" {
		return fmt.Errorf("API base URL is required")
	}

	if c.API.CandidateID == "" {
		return fmt.Errorf("candidate ID is required - set CROSSMINT_CANDIDATE_ID environment variable or add to config file")
	}

	if c.API.Timeout <= 0 {
		return fmt.Errorf("API timeout must be positive")
	}

	if c.API.RetryConfig.MaxAttempts <= 0 {
		return fmt.Errorf("retry max attempts must be positive")
	}

	if c.API.RateLimitConfig.RequestsPerSecond <= 0 {
		return fmt.Errorf("rate limit must be positive")
	}

	if c.Execution.Timeout <= 0 {
		return fmt.Errorf("execution timeout must be positive")
	}

	return nil
}

// ToRetryConfig converts the retry configuration to the retry package format
func (r RetryConfig) ToRetryConfig() retry.Config {
	return retry.Config{
		MaxAttempts:  r.MaxAttempts,
		InitialDelay: r.InitialDelay,
		MaxDelay:     r.MaxDelay,
		Multiplier:   r.Multiplier,
	}
}

// Save saves the configuration to a file
func (c *Config) Save(path string) error {
	viper.Set("api", c.API)
	viper.Set("logging", c.Logging)
	viper.Set("execution", c.Execution)

	return viper.WriteConfigAs(path)
}
