package server

import (
	"net/http"

	"okp4/cosmos-faucet/pkg/client"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

// NewSendRequestHandlerFn returns an HTTP REST handler for make transaction to a given address.
func NewSendRequestHandlerFn(faucet *client.Faucet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		err := faucet.SendTxMsg(bech32Addr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
	}
}
