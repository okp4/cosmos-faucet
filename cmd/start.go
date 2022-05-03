package cmd

import (
	"context"
	"net/http"

	"okp4/cosmos-faucet/pkg/client"
	"okp4/cosmos-faucet/pkg/server"

	"github.com/gorilla/mux"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			faucet, err := client.NewFaucet(context.Background(), config)
			if err != nil {
				return err
			}

			defer func(faucet *client.Faucet) {
				_ = faucet.Close()
			}(faucet)

			router := mux.NewRouter().StrictSlash(true)
			router.HandleFunc("/send/{address}", server.NewSendRequestHandlerFn(faucet)).Methods("GET")

			return http.ListenAndServe(addr, router)
		},
	}

	startCmd.Flags().StringVar(&addr, FlagAddress, ":8080", "rest api address")

	return startCmd
}

func init() {
	rootCmd.AddCommand(NewStartCommand())
}
