package cosmos

import (
	"context"
	"okp4/cosmos-faucet/pkg/actor/message"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcClient struct {
	grpcConn *grpc.ClientConn
}

func NewGrpcClient(address string, transportCreds credentials.TransportCredentials) (*GrpcClient, error) {
	grpcConn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(transportCreds),
	)
	if err != nil {
		return nil, err
	}

	return &GrpcClient{grpcConn: grpcConn}, nil
}

func (client *GrpcClient) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case message.GetAccount:
		goCTX, cancelFunc := context.WithDeadline(context.Background(), msg.Deadline)
		defer cancelFunc()

		account, err := client.GetAccount(goCTX, msg.Address)
		if err != nil {
			panic(err)
		}
		ctx.Respond(message.GetAccountResponse{
			Account: account,
		})

	case message.BroadcastTx:
		goCTX, cancelFunc := context.WithDeadline(context.Background(), msg.Deadline)
		defer cancelFunc()

		resp, err := client.BroadcastTx(goCTX, msg.Tx)
		if err != nil {
			panic(err)
		}
		ctx.Respond(message.BroadcastTxResponse{
			TxResponse: resp,
		})
	}
}

func (client *GrpcClient) GetAccount(context context.Context, address string) (*auth.BaseAccount, error) {
	authClient := auth.NewQueryClient(client.grpcConn)
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

func (client *GrpcClient) BroadcastTx(context context.Context, txBytes []byte) (*types.TxResponse, error) {
	txClient := tx.NewServiceClient(client.grpcConn)
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

func (client *GrpcClient) Close() error {
	return client.grpcConn.Close()
}
