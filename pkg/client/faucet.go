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

type Faucet interface {
	GetConfig() pkg.Config
	GetFromAddr() types.AccAddress
	SendTxMsg(ctx context.Context, addr string) (*types.TxResponse, error)
	Close() error
}

type faucet struct {
	config      pkg.Config
	grpcConn    *grpc.ClientConn
	fromAddr    types.AccAddress
	fromPrivKey crypto.PrivKey
	txConfig    client.TxConfig
}

func (f faucet) GetFromAddr() types.AccAddress {
	return f.fromAddr
}

func (f faucet) GetConfig() pkg.Config {
	return f.config
}

func NewFaucet(config pkg.Config) (Faucet, error) {
	conf := types.GetConfig()
	conf.SetBech32PrefixForAccount(config.Prefix, config.Prefix)

	grpcConn, err := grpc.Dial(
		config.GrpcAddress,
		grpc.WithTransportCredentials(getTransportCredentials(config)),
	)
	if err != nil {
		return nil, err
	}

	fromPrivKey, err := cosmos.GeneratePrivateKey(config.Mnemonic)
	if err != nil {
		return nil, err
	}

	fromAddr := types.AccAddress(fromPrivKey.PubKey().Address())

	return &faucet{
		config:      config,
		grpcConn:    grpcConn,
		fromAddr:    fromAddr,
		fromPrivKey: fromPrivKey,
		txConfig:    simapp.MakeTestEncodingConfig().TxConfig,
	}, nil
}

func (f faucet) SendTxMsg(ctx context.Context, addr string) (*types.TxResponse, error) {
	toAddr, err := types.GetFromBech32(addr, f.config.Prefix)
	if err != nil {
		return nil, err
	}

	txBuilder, err := cosmos.BuildUnsignedTx(f.config, f.txConfig, f.fromAddr, toAddr)
	if err != nil {
		return nil, err
	}

	account, err := cosmos.GetAccount(ctx, f.grpcConn, f.fromAddr.String())
	if err != nil {
		return nil, err
	}

	signerData := signing.SignerData{
		ChainID:       f.config.ChainID,
		AccountNumber: account.GetAccountNumber(),
		Sequence:      account.GetSequence(),
	}

	err = cosmos.SignTx(f.fromPrivKey, signerData, f.txConfig, txBuilder)
	if err != nil {
		return nil, err
	}

	txBytes, err := f.txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return cosmos.BroadcastTx(ctx, f.grpcConn, txBytes)
}

func (f faucet) Close() error {
	return f.grpcConn.Close()
}

func getTransportCredentials(config pkg.Config) credentials.TransportCredentials {
	switch {
	case config.NoTLS:
		return insecure.NewCredentials()
	case config.TLSSkipVerify:
		return credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}) // #nosec G402 : skip lint since it's an optional flag
	default:
		return credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}
}
