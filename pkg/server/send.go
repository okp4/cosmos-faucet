package server

import (
	"net/http"

	"okp4/cosmos-faucet/pkg/client"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// NewSendRequestHandlerFn returns an HTTP REST handler for make transaction to a given address.
func NewSendRequestHandlerFn(faucet *client.Faucet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		log.Info().Str("toAddress", bech32Addr).
			Str("fromAddress", faucet.FromAddr.String()).
			Msgf("Send %d%s to %s...", faucet.Config.AmountSend, faucet.Config.Denom, bech32Addr)

		err := faucet.SendTxMsg(r.Context(), bech32Addr)

		if err != nil {
			log.Err(err).Msg("Transaction failed.")
		}

		rest.CheckBadRequestError(w, err)
	}
}
