package cosmos

import (
	"context"
	"errors"

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
	err = account.Unmarshal(query.GetAccount().Value)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func BroadcastTx(context context.Context, grpcConn *grpc.ClientConn, txBytes []byte) error {
	txClient := tx.NewServiceClient(grpcConn)
	grpcRes, err := txClient.BroadcastTx(
		context,
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return err
	}

	if grpcRes.TxResponse.Code != 0 {
		return errors.New(grpcRes.TxResponse.RawLog)
	}
	return nil
}
