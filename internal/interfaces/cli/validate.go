package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/crossmint/megaverse-challenge/internal/application/strategies"
)

// NewValidateCommand returns a command that validates the megaverse.
func NewValidateCommand(deps *Dependencies) *cobra.Command {
	var phase string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the megaverse against a target phase",
		RunE: func(cmd *cobra.Command, args []string) error {
			if deps.Service == nil {
				return fmt.Errorf("service dependency not initialised")
			}

			ctx, cancel := withTimeout(context.Background(), deps)
			defer cancel()

			phase = strings.ToLower(phase)
			var strategy strategies.PatternStrategy

			switch phase {
			case "phase1", "cross", "x":
				strategy = strategies.NewCrossPatternStrategy()
			case "phase2", "logo":
				if deps.Repository == nil {
					return fmt.Errorf("repository dependency not initialised for Phase 2 validation")
				}
				logoStrategy := strategies.NewLogoPatternStrategy(deps.Repository)
				if _, err := logoStrategy.GeneratePlan(ctx); err != nil {
					return err
				}
				strategy = logoStrategy
			default:
				return fmt.Errorf("unknown phase '%s' (expected phase1 or phase2)", phase)
			}

			if err := deps.Service.ValidateMegaverse(ctx, strategy); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Validation for %s completed\n", phase)
			return nil
		},
	}

	cmd.Flags().StringVar(&phase, "phase", "phase1", "Target phase to validate (phase1 or phase2)")

	return cmd
}
