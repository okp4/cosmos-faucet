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
	SubmitTx(ctx context.Context) (*types.TxResponse, error)
	Send(addr string) error
	Subscribe(addr string) (<-chan *types.TxResponse, error)
	Close() error
}

type faucet struct {
	config        pkg.Config
	fromAddress   types.AccAddress
	amount        int64
	denom         string
	accountPrefix string
	grpcConn      *grpc.ClientConn
	triggerTx     <-chan bool
	pool          *MessagePool
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

	fromAddress := types.AccAddress(fromPrivKey.PubKey().Address())

	return &faucet{
		config:        config,
		fromAddress:   fromAddress,
		amount:        config.AmountSend,
		denom:         config.Denom,
		accountPrefix: config.Prefix,
		grpcConn:      grpcConn,
		triggerTx:     triggerTxChan,
		pool: NewMessagePool(
			WithTxSubmitter(
				makeTxSubmitter(config, simapp.MakeTestEncodingConfig().TxConfig, grpcConn, fromPrivKey, fromAddress),
			),
		),
	}, nil
}

func (f *faucet) Start() {
	go func() {
		log.Info().Msg("Starting submit routine")
		for range f.triggerTx {
			msgCount := f.pool.Size()

			resp, err := f.SubmitTx(context.Background())
			if err != nil {
				log.Err(err).Int("msgCount", msgCount).Msg("Could not submit transaction")
			} else if resp != nil {
				if resp.Code != 0 {
					log.Warn().
						Int("messageCount", msgCount).
						Interface("tx", resp).
						Msg("Transaction submitted with non 0 code")

				} else {
					log.Info().
						Int("messageCount", msgCount).
						Str("txHash", resp.TxHash).
						Uint32("txCode", resp.Code).
						Msg("Successfully submit transaction")
				}
			} else {
				log.Info().Msg("No message to submit")
			}
		}

		log.Info().Msg("Stopping submit routine")
	}()
}

func (f *faucet) GetConfig() pkg.Config {
	return f.config
}

func (f *faucet) SubmitTx(ctx context.Context) (*types.TxResponse, error) {
	return f.pool.Submit()
}

func (f *faucet) Send(addr string) error {
	msgSend, err := f.makeSendMsg(addr)
	if err != nil {
		return err
	}

	f.pool.RegisterMsg(msgSend)
	return nil
}

func (f *faucet) Subscribe(addr string) (<-chan *types.TxResponse, error) {
	msgSend, err := f.makeSendMsg(addr)
	if err != nil {
		return nil, err
	}

	return f.pool.SubscribeMsg(msgSend), nil
}

func (f *faucet) Close() error {
	return f.grpcConn.Close()
}

func (f *faucet) makeSendMsg(addr string) (types.Msg, error) {
	toAddr, err := types.GetFromBech32(addr, f.accountPrefix)
	if err != nil {
		return nil, err
	}

	return banktypes.NewMsgSend(
		f.fromAddress,
		toAddr,
		types.NewCoins(types.NewInt64Coin(f.denom, f.amount)),
	), nil
}

func makeTxSubmitter(config pkg.Config, txConfig client.TxConfig, grpcConn *grpc.ClientConn, privKey crypto.PrivKey, addr types.AccAddress) TxSubmitter {
	return func(msgs []types.Msg) (*types.TxResponse, error) {
		txBuilder, err := cosmos.BuildUnsignedTx(config, txConfig, msgs)
		if err != nil {
			return nil, err
		}

		account, err := cosmos.GetAccount(context.Background(), grpcConn, addr.String())
		if err != nil {
			return nil, err
		}

		signerData := signing.SignerData{
			ChainID:       config.ChainID,
			AccountNumber: account.GetAccountNumber(),
			Sequence:      account.GetSequence(),
		}

		err = cosmos.SignTx(privKey, signerData, txConfig, txBuilder)
		if err != nil {
			return nil, err
		}

		txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
		if err != nil {
			return nil, err
		}

		return cosmos.BroadcastTx(context.Background(), grpcConn, txBytes)
	}
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
