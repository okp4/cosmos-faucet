package cmd

import (
    "github.com/spf13/cobra"
)

// NewStartCommand returns a CLI command to start the REST api allowing to send tokens.
func NewStartCommand() *cobra.Command {
    startCmd := &cobra.Command{
        Use:   "start",
        Short: "Start the REST api",
        RunE: func(cmd *cobra.Command, args []string) error {
            return nil
        },
    }

    return startCmd
}

func init() {
    rootCmd.AddCommand(NewStartCommand())
}
