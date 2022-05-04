package cmd

import (
	"github.com/rs/zerolog/log"
	"okp4/cosmos-faucet/pkg/client"
	"okp4/cosmos-faucet/pkg/server"

	"github.com/spf13/cobra"
)

const (
	FlagAddress = "address"
)

// NewStartCommand returns a CLI command to start the REST api allowing to send tokens.
func NewStartCommand() *cobra.Command {
	var addr string

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the REST api",
		Run: func(cmd *cobra.Command, args []string) {
			faucet, err := client.NewFaucet(config)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed create a new faucet instance")
			}

			defer func(faucet *client.Faucet) {
				_ = faucet.Close()
				log.Info().Msg("Server stopped")
			}(faucet)

			httpServer := server.NewServer(faucet)
			httpServer.Start(addr)
		},
	}

	startCmd.Flags().StringVar(&addr, FlagAddress, ":8080", "rest api address")

	return startCmd
}

func init() {
	rootCmd.AddCommand(NewStartCommand())
}
