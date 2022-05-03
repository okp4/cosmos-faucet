package cmd

import (
	"net/http"

	"okp4/cosmos-faucet/rest"

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
			router := mux.NewRouter().StrictSlash(true)
			router.HandleFunc("/send/{address}", rest.NewSendRequestHandlerFn(config)).Methods("GET")

			return http.ListenAndServe(addr, router)
		},
	}

	startCmd.Flags().StringVar(&addr, FlagAddress, ":8080", "rest api address")

	return startCmd
}

func init() {
	rootCmd.AddCommand(NewStartCommand())
}
