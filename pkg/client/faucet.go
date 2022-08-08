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
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Faucet interface {
	Start()
	GetConfig() pkg.Config
	GetFromAddr() types.AccAddress
	SubmitTx(ctx context.Context) (*types.TxResponse, error)
	Send(addr string) error
	Close() error
}

type faucet struct {
	config      pkg.Config
	grpcConn    *grpc.ClientConn
	fromAddr    types.AccAddress
	fromPrivKey crypto.PrivKey
	txConfig    client.TxConfig
	buffer      []types.Msg
	triggerTx   <-chan bool
}

func (f *faucet) GetFromAddr() types.AccAddress {
	return f.fromAddr
}

func (f *faucet) GetConfig() pkg.Config {
	return f.config
}

func NewFaucet(config pkg.Config, triggerTxChan <-chan bool) (Faucet, error) {
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
		triggerTx:   triggerTxChan,
	}, nil
}

func (f *faucet) Start() {
	go func() {
		for range f.triggerTx {
			msgCount := len(f.buffer)
			resp, err := f.SubmitTx(context.Background())
			if err != nil {
				log.Err(err).Int("msgCount", msgCount).Msg("Could not submit transaction")
			} else if resp != nil {
				log.Info().
					Int("messageCount", msgCount).
					Str("txHash", resp.TxHash).
					Uint32("txCode", resp.Code).
					Msg("Successfully submit transaction")
			} else {
				log.Info().Msg("No message to submit")
			}
		}
		log.Info().Msg("Stopping submit routine")
	}()
}

func (f *faucet) SubmitTx(ctx context.Context) (*types.TxResponse, error) {
	if len(f.buffer) == 0 {
		return nil, nil
	}

	defer func() {
		f.buffer = f.buffer[:0]
	}()

	txBuilder, err := cosmos.BuildUnsignedTx(f.config, f.txConfig, f.buffer)
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

func (f *faucet) Send(addr string) error {
	toAddr, err := types.GetFromBech32(addr, f.config.Prefix)
	if err != nil {
		return err
	}

	f.buffer = append(
		f.buffer,
		banktypes.NewMsgSend(
			f.fromAddr,
			toAddr,
			types.NewCoins(types.NewInt64Coin(f.config.Denom, f.config.AmountSend)),
		),
	)
	return nil
}

func (f *faucet) Close() error {
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
