package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"okp4/cosmos-faucet/graph/generated"
	"okp4/cosmos-faucet/graph/model"

	"github.com/rs/zerolog/log"
)

// Send is the resolver for the send field.
func (r *mutationResolver) Send(ctx context.Context, input model.SendInput) (*string, error) {
	if err := r.CaptchaResolver.CheckRecaptcha(ctx, input.CaptchaToken); err != nil {
		return nil, err
	}

	if err := r.Faucet.Send(input.ToAddress); err != nil {
		log.Err(err).Str("toAddress", input.ToAddress).Msg("Could not register send request.")
		return nil, err
	}

	log.Info().Str("toAddress", input.ToAddress).Msg("Register send request")
	return nil, nil
}

// Configuration is the resolver for the configuration field.
func (r *queryResolver) Configuration(ctx context.Context) (*model.Configuration, error) {
	return &model.Configuration{
		ChainID:    r.Faucet.GetConfig().ChainID,
		Denom:      r.Faucet.GetConfig().Denom,
		Prefix:     r.Faucet.GetConfig().Prefix,
		AmountSend: r.Faucet.GetConfig().AmountSend,
		FeeAmount:  r.Faucet.GetConfig().FeeAmount,
		Memo:       r.Faucet.GetConfig().Memo,
		GasLimit:   r.Faucet.GetConfig().GasLimit,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
