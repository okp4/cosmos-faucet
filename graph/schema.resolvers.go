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
		log.Err(err).Str("toAddress", input.ToAddress).Msg("Could not register send request")
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

// Send is the resolver for the send field.
func (r *subscriptionResolver) Send(ctx context.Context, input model.SendInput) (<-chan *model.TxResponse, error) {
	if err := r.CaptchaResolver.CheckRecaptcha(ctx, input.CaptchaToken); err != nil {
		return nil, err
	}

	rawTxResponseChan, err := r.Faucet.Subscribe(input.ToAddress)
	if err != nil {
		log.Err(err).Str("toAddress", input.ToAddress).Msg("Could not register send request")
		return nil, err
	}

	log.Info().Str("toAddress", input.ToAddress).Msg("Register send request")
	txResponseChan := make(chan *model.TxResponse)
	go func() {
		for rawTx := range rawTxResponseChan {
			txResponseChan <- &model.TxResponse{
				Hash:      rawTx.TxHash,
				Code:      int(rawTx.Code),
				RawLog:    &rawTx.RawLog,
				GasWanted: rawTx.GasWanted,
				GasUsed:   rawTx.GasUsed,
			}
		}
		close(txResponseChan)
	}()

	return txResponseChan, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
