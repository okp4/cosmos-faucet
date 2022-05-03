package client

import (
	"context"

	"okp4/cosmos-faucet/pkg"
	"okp4/cosmos-faucet/pkg/cosmos"

	"github.com/cosmos/cosmos-sdk/client"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Faucet struct {
	Config      pkg.Config
	GRPCConn    *grpc.ClientConn
	FromAddr    types.AccAddress
	FromPrivKey crypto.PrivKey
	Account     *auth.BaseAccount
	TxConfig    client.TxConfig
}

func NewFaucet(ctx context.Context, config pkg.Config) (*Faucet, error) {
	conf := types.GetConfig()
	conf.SetBech32PrefixForAccount(config.Prefix, config.Prefix)

	grpcConn, err := grpc.Dial(
		config.GrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	fromPrivKey, err := cosmos.GeneratePrivateKey(config.Mnemonic)
	if err != nil {
		return nil, err
	}

	fromAddr := types.AccAddress(fromPrivKey.PubKey().Address())
	account, err := cosmos.GetAccount(ctx, grpcConn, fromAddr.String())
	if err != nil {
		return nil, err
	}

	return &Faucet{
		Config:      config,
		GRPCConn:    grpcConn,
		FromAddr:    fromAddr,
		FromPrivKey: fromPrivKey,
		Account:     account,
		TxConfig:    simapp.MakeTestEncodingConfig().TxConfig,
	}, nil
}

func (f *Faucet) SendTxMsg(ctx context.Context, addr string) error {
	toAddr, err := types.GetFromBech32(addr, f.Config.Prefix)
	if err != nil {
		return err
	}

	txBuilder, err := cosmos.BuildUnsignedTx(f.Config, f.TxConfig, f.FromAddr, toAddr)
	if err != nil {
		return err
	}

	err = cosmos.SignTx(f.Config, f.FromPrivKey, f.Account, f.TxConfig, txBuilder)
	if err != nil {
		return err
	}

	txBytes, err := f.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	return cosmos.BroadcastTx(ctx, f.GRPCConn, txBytes)
}

func (f *Faucet) Close() error {
	return f.GRPCConn.Close()
}
