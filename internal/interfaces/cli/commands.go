package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/crossmint/megaverse-challenge/internal/application"
	"github.com/crossmint/megaverse-challenge/internal/domain"
	cfgpkg "github.com/crossmint/megaverse-challenge/internal/infrastructure/config"
)

// Dependencies bundles the services required by CLI commands.
type Dependencies struct {
	Config     *cfgpkg.Config
	ConfigPath string
	Service    *application.MegaverseService
	Repository domain.MegaverseRepository
}

// NewRootCommand creates the root cobra command and registers all subcommands.
func NewRootCommand(deps *Dependencies) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "megaverse",
		Short: "Command line tools for mastering the Crossmint megaverse",
		Long:  "Megaverse CLI allows you to initialise, render and validate Crossmint megaverses programmatically.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "init" {
				return nil
			}
			if deps.Config == nil {
				return fmt.Errorf("configuration not loaded")
			}
			if deps.Config.API.CandidateID == "" {
				return fmt.Errorf("candidate ID missing; run 'megaverse init --candidate <id>'")
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&deps.ConfigPath, "config", deps.ConfigPath, "Path to configuration file")

	rootCmd.AddCommand(NewInitCommand(deps))
	rootCmd.AddCommand(NewPhase1Command(deps))
	rootCmd.AddCommand(NewPhase2Command(deps))
	rootCmd.AddCommand(NewStatusCommand(deps))

	return rootCmd
}

// withTimeout returns a context with timeout suitable for API calls.
func withTimeout(parent context.Context, deps *Dependencies) (context.Context, context.CancelFunc) {
	timeout := 2 * time.Minute
	if deps != nil && deps.Config != nil && deps.Config.Execution.Timeout > 0 {
		timeout = deps.Config.Execution.Timeout
	}
	return context.WithTimeout(parent, timeout)
}
