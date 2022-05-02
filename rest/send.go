package rest

import (
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
	"okp4/cosmos-faucet/pkg/send"
	"okp4/cosmos-faucet/util"
)

// NewSendRequestHandlerFn returns an HTTP REST handler for make transaction to a given address.
func NewSendRequestHandlerFn(config util.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		err := send.SendTx(config, bech32Addr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
	}
}
