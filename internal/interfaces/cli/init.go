package cli

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"

    cfgpkg "github.com/crossmint/megaverse-challenge/internal/infrastructure/config"
)

// NewInitCommand creates the CLI command for bootstrapping configuration.
func NewInitCommand(deps *Dependencies) *cobra.Command {
    var candidateID string
    var baseURL string

    cmd := &cobra.Command{
        Use:   "init",
        Short: "Initialise the megaverse configuration",
        Long:  "Persist the Crossmint candidate ID and optional API base URL to the configuration file.",
        RunE: func(cmd *cobra.Command, args []string) error {
            if candidateID == "" {
                return fmt.Errorf("candidate ID is required (use --candidate)")
            }

            cfg := deps.Config
            if cfg == nil {
                cfg = cfgpkg.DefaultConfig()
            }

            cfg.API.CandidateID = candidateID
            if baseURL != "" {
                cfg.API.BaseURL = baseURL
            }

            if err := cfg.Validate(); err != nil {
                return err
            }

            configPath := deps.ConfigPath
            if configPath == "" {
                configPath = filepath.Join("config", "config.yaml")
            }

            if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
                return fmt.Errorf("failed to create config directory: %w", err)
            }

            if err := cfg.Save(configPath); err != nil {
                return fmt.Errorf("failed to save configuration: %w", err)
            }

            deps.Config = cfg
            deps.ConfigPath = configPath

            fmt.Fprintf(cmd.OutOrStdout(), "Configuration saved to %s\n", configPath)
            return nil
        },
    }

    cmd.Flags().StringVar(&candidateID, "candidate", "", "Crossmint candidate ID")
    cmd.Flags().StringVar(&baseURL, "base-url", "", "Override API base URL")

    return cmd
}
