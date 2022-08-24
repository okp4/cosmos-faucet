package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"okp4/cosmos-faucet/graph/generated"
	"okp4/cosmos-faucet/graph/model"
	"okp4/cosmos-faucet/pkg/actor/message"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
)

// Send is the resolver for the send field.
func (r *mutationResolver) Send(ctx context.Context, input model.SendInput) (*string, error) {
	addr, err := types.GetFromBech32(input.ToAddress, r.AddressPrefix)
	if err != nil {
		log.Err(err).Str("toAddress", input.ToAddress).Msg("❌ Could not serve send mutation")
		return nil, err
	}

	if err = r.CaptchaResolver.CheckRecaptcha(ctx, input.CaptchaToken); err != nil {
		log.Err(err).Str("toAddress", input.ToAddress).Msg("❌ Could not serve send mutation")
		return nil, err
	}

	r.Context.Send(r.Faucet, &message.RequestFunds{Address: addr})
	return nil, nil
}

// Configuration is the resolver for the configuration field.
func (r *queryResolver) Configuration(ctx context.Context) (*model.Configuration, error) {
	return r.Config, nil
}

// Send is the resolver for the send field.
func (r *subscriptionResolver) Send(ctx context.Context, input model.SendInput) (<-chan *model.TxResponse, error) {
	addr, err := types.GetFromBech32(input.ToAddress, r.AddressPrefix)
	if err != nil {
		log.Err(err).Str("toAddress", input.ToAddress).Msg("❌ Could not serve send mutation")
		return nil, err
	}

	if err := r.CaptchaResolver.CheckRecaptcha(ctx, input.CaptchaToken); err != nil {
		log.Err(err).Str("toAddress", input.ToAddress).Msg("❌ Could not serve send mutation")
		return nil, err
	}

	txResponseChan := make(chan *model.TxResponse)
	r.Context.Send(
		r.Faucet,
		&message.RequestFunds{
			Address: addr,
			TxSubscriber: r.Context.Spawn(
				actor.PropsFromFunc(
					func(c actor.Context) {
						switch msg := c.Message().(type) {
						case *message.BroadcastTxResponse:
							txResponseChan <- &model.TxResponse{
								Hash:      msg.TxResponse.TxHash,
								Code:      int(msg.TxResponse.Code),
								RawLog:    &msg.TxResponse.RawLog,
								GasWanted: msg.TxResponse.GasWanted,
								GasUsed:   msg.TxResponse.GasUsed,
							}
							close(txResponseChan)
							c.Stop(c.Self())
						}
					},
				),
			),
		},
	)

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
