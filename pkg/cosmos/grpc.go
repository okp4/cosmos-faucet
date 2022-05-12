package cosmos

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"
)

func GetAccount(context context.Context, grpcConn *grpc.ClientConn, address string) (*auth.BaseAccount, error) {
	authClient := auth.NewQueryClient(grpcConn)
	query, err := authClient.Account(context, &auth.QueryAccountRequest{Address: address})
	if err != nil {
		return nil, err
	}

	var account auth.BaseAccount
	if err := account.Unmarshal(query.GetAccount().Value); err != nil {
		return nil, err
	}

	return &account, nil
}

func BroadcastTx(context context.Context, grpcConn *grpc.ClientConn, txBytes []byte) (*types.TxResponse, error) {
	txClient := tx.NewServiceClient(grpcConn)
	grpcRes, err := txClient.BroadcastTx(
		context,
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return nil, err
	}

	return grpcRes.TxResponse, nil
}
