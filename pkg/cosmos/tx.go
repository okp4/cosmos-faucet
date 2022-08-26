package cosmos

import (
	"fmt"
	"okp4/cosmos-faucet/pkg/actor/message"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	sdk "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/rs/zerolog/log"
)

// TxHandler represents an actor in charge of building and signing transactions.
type TxHandler struct {
	privKey           crypto.PrivKey
	address           string
	chainID           string
	config            sdk.TxConfig
	signMode          signing.SignMode
	cosmosClientProps *actor.Props
	cosmosClient      *actor.PID
}

func NewTxHandler(opts ...Option) *TxHandler {
	handler := &TxHandler{}
	for _, opt := range opts {
		opt(handler)
	}

	return handler
}

type Option func(handler *TxHandler)

func WithMnemonicMust(mnemonic string) Option {
	privateKey, err := ParseMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}

	return WithPrivateKey(privateKey)
}

func WithPrivateKey(privateKey crypto.PrivKey) Option {
	return func(handler *TxHandler) {
		handler.privKey = privateKey
		handler.address = types.AccAddress(privateKey.PubKey().Address()).String()
	}
}

func WithChainID(chainID string) Option {
	return func(handler *TxHandler) {
		handler.chainID = chainID
	}
}

func WithTxConfig(config sdk.TxConfig) Option {
	return func(handler *TxHandler) {
		handler.config = config
		handler.signMode = config.SignModeHandler().DefaultMode()
	}
}

func WithCosmosClientProps(props *actor.Props) Option {
	return func(handler *TxHandler) {
		handler.cosmosClientProps = props
	}
}

// nolint: funlen
func (handler *TxHandler) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		handler.cosmosClient = ctx.Spawn(handler.cosmosClientProps)

	case *message.MakeTx:
		if time.Now().After(msg.Deadline) {
			log.Warn().Msg("😞 Deadline exceeded, ignore transaction.")
			break
		}

		unsignedTx, err := handler.BuildUnsignedTx(msg.Msgs, msg.Memo, msg.GasLimit, msg.FeeAmount)
		if err != nil {
			log.Panic().Err(err).Msg("❌ Could not build transaction.")
		}

		accountResp, err := ctx.RequestFuture(
			handler.cosmosClient,
			&message.GetAccount{Deadline: msg.Deadline, Address: handler.address},
			time.Until(msg.Deadline),
		).Result()
		if err != nil {
			log.Panic().Err(err).Msg("❌ Could not get account information.")
		}

		var account *auth.BaseAccount
		switch resp := accountResp.(type) {
		case *message.GetAccountResponse:
			account = resp.Account
		default:
			log.Panic().Err(fmt.Errorf("wrong response message")).Msg("❌ Could not get account information.")
		}

		signedTx, err := handler.SignTx(
			unsignedTx, authsigning.SignerData{
				ChainID:       handler.chainID,
				AccountNumber: account.GetAccountNumber(),
				Sequence:      account.GetSequence(),
			},
		)
		if err != nil {
			log.Panic().Err(err).Msg("❌ Could not sign transaction.")
		}

		tx, err := handler.EncodeTx(signedTx)
		if err != nil {
			log.Panic().Err(err).Msg("❌ Could not encode transaction.")
		}

		txResp, err := ctx.RequestFuture(
			handler.cosmosClient,
			&message.BroadcastTx{Deadline: msg.Deadline, Tx: tx},
			time.Until(msg.Deadline),
		).Result()
		if err != nil {
			log.Panic().Err(err).Msg("❌ Could not broadcast transaction.")
		}

		switch resp := txResp.(type) {
		case *message.BroadcastTxResponse:
			ctx.Send(msg.TxSubscriber, txResp)
			if resp.TxResponse.Code != 0 {
				log.Warn().
					Int("messageCount", len(msg.Msgs)).
					Interface("tx", resp.TxResponse).
					Msg("😞 Transaction submitted with non 0 code")
			} else {
				log.Info().
					Int("messageCount", len(msg.Msgs)).
					Str("txHash", resp.TxResponse.TxHash).
					Uint32("txCode", resp.TxResponse.Code).
					Msg("🚀 Successfully submit transaction")
			}
		default:
			log.Panic().Err(fmt.Errorf("wrong response message")).Msg("❌ Could not broadcast transaction.")
		}
	}
}

func (handler *TxHandler) BuildUnsignedTx(
	msgs []types.Msg,
	memo string,
	gasLimit uint64,
	feeAmount types.Coins,
) (sdk.TxBuilder, error) {
	txBuilder := handler.config.NewTxBuilder()

	err := txBuilder.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}

	txBuilder.SetMemo(memo)
	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(feeAmount)

	return txBuilder, nil
}

func (handler *TxHandler) SignTx(txBuilder sdk.TxBuilder, signerData authsigning.SignerData) (authsigning.Tx, error) {
	sig := signing.SignatureV2{
		PubKey: handler.privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  handler.signMode,
			Signature: nil,
		},
		Sequence: signerData.Sequence,
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	bytesToSign, err := handler.config.SignModeHandler().GetSignBytes(handler.signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	sigBytes, err := handler.privKey.Sign(bytesToSign)
	if err != nil {
		return nil, err
	}

	sig.Data = &signing.SingleSignatureData{
		SignMode:  handler.signMode,
		Signature: sigBytes,
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	return txBuilder.GetTx(), nil
}

func (handler *TxHandler) EncodeTx(tx authsigning.Tx) ([]byte, error) {
	txBytes, err := handler.config.TxEncoder()(tx)
	if err != nil {
		return nil, err
	}

	return txBytes, err
}

func ParseMnemonic(mnemonic string) (crypto.PrivKey, error) {
	algo, err := keyring.NewSigningAlgoFromString("secp256k1", keyring.SigningAlgoList{hd.Secp256k1})
	if err != nil {
		return nil, err
	}

	hdPath := hd.CreateHDPath(118, 0, 0).String()

	// create master key and derive first key for keyring
	derivedPriv, err := algo.Derive()(mnemonic, "", hdPath)
	if err != nil {
		return nil, err
	}

	return algo.Generate()(derivedPriv), nil
}
