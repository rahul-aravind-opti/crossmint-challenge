package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/crossmint/megaverse-challenge/internal/application/strategies"
)

// NewPhase1Command returns the command that executes Phase 1 of the challenge.
func NewPhase1Command(deps *Dependencies) *cobra.Command {
	var shouldValidate bool

	cmd := &cobra.Command{
		Use:   "phase1",
		Short: "Create the Phase 1 POLYanet cross",
		RunE: func(cmd *cobra.Command, args []string) error {
			if deps.Service == nil {
				return fmt.Errorf("service dependency not initialised")
			}

			ctx, cancel := withTimeout(context.Background(), deps)
			defer cancel()

			strategy := strategies.NewCrossPatternStrategy()
			if err := deps.Service.ExecuteStrategy(ctx, strategy); err != nil {
				return fmt.Errorf("failed to execute Phase 1 strategy: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Phase 1 cross created successfully")

			if shouldValidate {
				if err := deps.Service.ValidateMegaverse(ctx, strategy); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Phase 1 validation completed")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&shouldValidate, "validate", false, "Validate the megaverse after creation")

	return cmd
}
