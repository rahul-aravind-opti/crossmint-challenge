package cli

import (
    "context"
    "fmt"
    "time"

    "github.com/spf13/cobra"
)

// NewStatusCommand returns a command that prints current megaverse metadata.
func NewStatusCommand(deps *Dependencies) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "status",
        Short: "Display current configuration and goal map dimensions",
        RunE: func(cmd *cobra.Command, args []string) error {
            if deps.Config == nil {
                return fmt.Errorf("configuration not loaded")
            }

            fmt.Fprintf(cmd.OutOrStdout(), "Candidate ID: %s\n", deps.Config.API.CandidateID)
            fmt.Fprintf(cmd.OutOrStdout(), "API Base URL: %s\n", deps.Config.API.BaseURL)

            if deps.Repository != nil {
                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                defer cancel()

                goal, err := deps.Repository.GetGoalMap(ctx)
                if err == nil && goal != nil && len(goal.Goal) > 0 {
                    rows := len(goal.Goal)
                    cols := len(goal.Goal[0])
                    fmt.Fprintf(cmd.OutOrStdout(), "Goal map dimensions: %dx%d\n", cols, rows)
                } else if err != nil {
                    fmt.Fprintf(cmd.ErrOrStderr(), "Warning: unable to fetch goal map: %v\n", err)
                }
            }

            return nil
        },
    }

    return cmd
}
