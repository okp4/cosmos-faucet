package client

import (
	"context"
	"crypto/tls"
	"time"

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

// TriggerTx is the message sent to the faucet to trigger a transaction submit, it conveys parameters related to the
// ongoing transaction.
type TriggerTx struct {
	// Deadline specify a time which the transaction execute shall not exceed.
	Deadline *time.Time
}

// TriggerTxOption is used to configure a TriggerTx.
type TriggerTxOption func(msg *TriggerTx)

// MakeTriggerTx create a new TriggerTx configured through the provided options.
func MakeTriggerTx(opts ...TriggerTxOption) *TriggerTx {
	msg := &TriggerTx{}
	for _, opt := range opts {
		opt(msg)
	}

	return msg
}

// WithDeadline configure a deadline on a TriggerTx.
func WithDeadline(deadline time.Time) TriggerTxOption {
	return func(msg *TriggerTx) {
		msg.Deadline = &deadline
	}
}

func (trigger TriggerTx) toCtx() (context.Context, context.CancelFunc) {
	if trigger.Deadline != nil {
		return context.WithDeadline(context.Background(), *trigger.Deadline)
	}
	return context.Background(), nil
}

type Faucet struct {
	config      pkg.Config
	fromAddress types.AccAddress
	grpcConn    *grpc.ClientConn
	triggerTx   <-chan *TriggerTx
	pool        *MessagePool
}

func NewFaucet(config pkg.Config, triggerTxChan <-chan *TriggerTx) (*Faucet, error) {
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

	faucet := &Faucet{
		config:      config,
		fromAddress: fromAddress,
		grpcConn:    grpcConn,
		triggerTx:   triggerTxChan,
		pool: NewMessagePool(
			WithTxSubmitter(
				makeTxSubmitter(config, simapp.MakeTestEncodingConfig().TxConfig, grpcConn, fromPrivKey, fromAddress),
			),
		),
	}

	faucet.start()
	return faucet, nil
}

func (f *Faucet) start() {
	go func() {
		log.Info().Msg("Starting submit routine")
		for trigger := range f.triggerTx {
			if trigger.Deadline.After(time.Now()) {
				f.handleTriggerTx(trigger)
			}
		}

		log.Info().Msg("Stopping submit routine")
	}()
}

func (f *Faucet) handleTriggerTx(trigger *TriggerTx) {
	ctx, cancelFunc := trigger.toCtx()
	defer cancelFunc()

	msgCount := f.pool.Size()
	resp, err := f.pool.Submit(ctx)
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

func (f *Faucet) GetConfig() pkg.Config {
	return f.config
}

func (f *Faucet) Send(addr string) error {
	msgSend, err := f.makeSendMsg(addr)
	if err != nil {
		return err
	}

	f.pool.RegisterMsg(msgSend)
	return nil
}

func (f *Faucet) Subscribe(addr string) (<-chan *types.TxResponse, error) {
	msgSend, err := f.makeSendMsg(addr)
	if err != nil {
		return nil, err
	}

	return f.pool.SubscribeMsg(msgSend), nil
}

func (f *Faucet) Close() error {
	return f.grpcConn.Close()
}

func (f *Faucet) makeSendMsg(addr string) (types.Msg, error) {
	toAddr, err := types.GetFromBech32(addr, f.config.Prefix)
	if err != nil {
		return nil, err
	}

	return banktypes.NewMsgSend(
		f.fromAddress,
		toAddr,
		types.NewCoins(types.NewInt64Coin(f.config.Denom, f.config.AmountSend)),
	), nil
}

func makeTxSubmitter(config pkg.Config, txConfig client.TxConfig, grpcConn *grpc.ClientConn, privKey crypto.PrivKey, addr types.AccAddress) TxSubmitter {
	return func(ctx context.Context, msgs []types.Msg) (*types.TxResponse, error) {
		txBuilder, err := cosmos.BuildUnsignedTx(config, txConfig, msgs)
		if err != nil {
			return nil, err
		}

		account, err := cosmos.GetAccount(ctx, grpcConn, addr.String())
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

		return cosmos.BroadcastTx(ctx, grpcConn, txBytes)
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
