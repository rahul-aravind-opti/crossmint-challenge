package main

import (
	"log"
	"os"

	"github.com/crossmint/megaverse-challenge/internal/application"
	"github.com/crossmint/megaverse-challenge/internal/infrastructure/api"
	cfgpkg "github.com/crossmint/megaverse-challenge/internal/infrastructure/config"
	"github.com/crossmint/megaverse-challenge/internal/interfaces/cli"
	"github.com/crossmint/megaverse-challenge/pkg/ratelimit"
)

func main() {
	deps := &cli.Dependencies{}
	defaultConfigPath := "config/config.yaml"

	// Attempt to load configuration; log but don't exit if missing (init command handles writing it)
	if _, err := os.Stat(defaultConfigPath); err == nil {
		cfg, err := cfgpkg.LoadFromFile(defaultConfigPath)
		if err != nil {
			log.Printf("warning: failed to load configuration: %v", err)
		} else {
			deps.Config = cfg
		}
	} else {
		// Fall back to default config for init command
		deps.Config = cfgpkg.DefaultConfig()
	}
	deps.ConfigPath = defaultConfigPath

	if deps.Config != nil && deps.Config.API.CandidateID != "" {
		retryCfg := deps.Config.API.RetryConfig.ToRetryConfig()
		client := api.NewClient(api.ClientConfig{
			BaseURL:           deps.Config.API.BaseURL,
			CandidateID:       deps.Config.API.CandidateID,
			Timeout:           deps.Config.API.Timeout,
			RetryConfig:       retryCfg,
			RequestsPerSecond: deps.Config.API.RateLimitConfig.RequestsPerSecond,
		})

		repository := api.NewRepository(client)
		deps.Repository = repository

		logger := log.New(os.Stdout, "[megaverse] ", log.LstdFlags)
		rps := deps.Config.API.RateLimitConfig.RequestsPerSecond
		if rps <= 0 {
			rps = 2.0
		}
		limiter := ratelimit.NewLimiter(rps)

		deps.Service = application.NewMegaverseService(repository, logger, limiter)
	}

	rootCmd := cli.NewRootCommand(deps)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
