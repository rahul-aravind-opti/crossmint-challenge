package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/crossmint/megaverse-challenge/internal/application/strategies"
)

// NewPhase2Command returns the command that executes Phase 2 of the challenge.
func NewPhase2Command(deps *Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phase2",
		Short: "Render the Phase 2 megaverse logo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if deps.Service == nil || deps.Repository == nil {
				return fmt.Errorf("dependencies not initialised for Phase 2")
			}

			ctx, cancel := withTimeout(context.Background(), deps)
			defer cancel()

			strategy := strategies.NewLogoPatternStrategy(deps.Repository)
			if err := deps.Service.ExecuteStrategy(ctx, strategy); err != nil {
				return fmt.Errorf("failed to execute Phase 2 strategy: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Phase 2 logo created successfully")
			return nil
		},
	}
	return cmd
}
