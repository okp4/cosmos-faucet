package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"okp4/cosmos-faucet/graph/generated"
	"okp4/cosmos-faucet/graph/model"
)

func (r *mutationResolver) Send(ctx context.Context, input model.SendInput) (*model.TxResponse, error) {
	resp, err := r.Faucet.SendTxMsg(ctx, input.ToAddress)

	if err != nil {
		return nil, err
	}

	return &model.TxResponse{
		Hash:      resp.TxHash,
		Code:      int(resp.Code),
		RawLog:    &resp.RawLog,
		GasWanted: resp.GasWanted,
		GasUsed:   resp.GasUsed,
	}, nil
}

func (r *queryResolver) Health(ctx context.Context) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
