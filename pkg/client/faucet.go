package client

import (
	"context"
	"crypto/tls"

	"okp4/cosmos-faucet/pkg"
	"okp4/cosmos-faucet/pkg/cosmos"

	"github.com/cosmos/cosmos-sdk/client"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Faucet struct {
	Config      pkg.Config
	GRPCConn    *grpc.ClientConn
	FromAddr    types.AccAddress
	FromPrivKey crypto.PrivKey
	TxConfig    client.TxConfig
}

func NewFaucet(config pkg.Config) (*Faucet, error) {
	conf := types.GetConfig()
	conf.SetBech32PrefixForAccount(config.Prefix, config.Prefix)

	var opts credentials.TransportCredentials
	switch {
	case config.NoTLS:
		opts = insecure.NewCredentials()
	case config.TLSSkipVerify:
		credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}) // #nosec G402 : skip lint since it's an optional flag
	default:
		credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}

	grpcConn, err := grpc.Dial(
		config.GrpcAddress,
		grpc.WithTransportCredentials(opts),
	)
	if err != nil {
		return nil, err
	}

	fromPrivKey, err := cosmos.GeneratePrivateKey(config.Mnemonic)
	if err != nil {
		return nil, err
	}

	fromAddr := types.AccAddress(fromPrivKey.PubKey().Address())

	return &Faucet{
		Config:      config,
		GRPCConn:    grpcConn,
		FromAddr:    fromAddr,
		FromPrivKey: fromPrivKey,
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

	account, err := cosmos.GetAccount(ctx, f.GRPCConn, f.FromAddr.String())
	if err != nil {
		return err
	}

	signerData := signing.SignerData{
		ChainID:       f.Config.ChainID,
		AccountNumber: account.GetAccountNumber(),
		Sequence:      account.GetSequence(),
	}

	err = cosmos.SignTx(f.FromPrivKey, signerData, f.TxConfig, txBuilder)
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
