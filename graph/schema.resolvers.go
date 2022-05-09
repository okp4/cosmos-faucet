package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"okp4/cosmos-faucet/graph/generated"
	"okp4/cosmos-faucet/graph/model"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog/log"
)

func (r *mutationResolver) Send(ctx context.Context, input model.SendInput) (*model.TxResponse, error) {
	resp, err := r.Faucet.SendTxMsg(ctx, input.ToAddress)

	if err != nil {
		log.Err(err).Msg("Transaction failed.")
		return nil, err
	}

	log.Info().
		Str("toAddress", input.ToAddress).
		Str("fromAddress", r.Faucet.FromAddr.String()).
		Msgf("Send %d%s to %s...", r.Faucet.Config.AmountSend, r.Faucet.Config.Denom, input.ToAddress)

	response := &model.TxResponse{
		Hash:      resp.TxHash,
		Code:      int(resp.Code),
		RawLog:    &resp.RawLog,
		GasWanted: resp.GasWanted,
		GasUsed:   resp.GasUsed,
	}

	log.Debug().
		Interface("response", response).
		Msg("Transaction has been broadcast")

	if resp.Code != 0 {
		graphql.AddErrorf(ctx, "transaction is not successful: %s (code: %d)", resp.RawLog, resp.Code)
		log.Error().Str("log", resp.RawLog).
			Uint32("code", resp.Code).
			Msgf("transaction is not successful")
	}

	return response, nil
}

func (r *queryResolver) Configuration(ctx context.Context) (*model.Configuration, error) {
	return &model.Configuration{
		ChainID:    r.Faucet.Config.ChainID,
		Denom:      r.Faucet.Config.Denom,
		Prefix:     r.Faucet.Config.Prefix,
		AmountSend: r.Faucet.Config.AmountSend,
		FeeAmount:  r.Faucet.Config.FeeAmount,
		Memo:       r.Faucet.Config.Memo,
		GasLimit:   r.Faucet.Config.GasLimit,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
