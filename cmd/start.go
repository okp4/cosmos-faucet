package cmd

import (
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"net/http"
	"okp4/cosmos-faucet/rest"
)

// NewStartCommand returns a CLI command to start the REST api allowing to send tokens.
func NewStartCommand() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the REST api",
		RunE: func(cmd *cobra.Command, args []string) error {
			router := mux.NewRouter().StrictSlash(true)
			router.HandleFunc("/send/{address}", rest.NewSendRequestHandlerFn(config)).Methods("GET")

			return http.ListenAndServe(":10000", router)
		},
	}

	return startCmd
}

func init() {
	rootCmd.AddCommand(NewStartCommand())
}
