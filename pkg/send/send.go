package send

import (
	"context"
	"errors"
	"fmt"

	"okp4/cosmos-faucet/util"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
)

func Send(config util.Config, address string) error {
	conf := types.GetConfig()
	conf.SetBech32PrefixForAccount(config.Prefix, config.Prefix)

	fromPrivKey, err := GeneratePrivateKey(config.Mnemonic)
	if err != nil {
		return err
	}

	fromAddr := types.AccAddress(fromPrivKey.PubKey().Address())

	toAddr, err := types.GetFromBech32(address, config.Prefix)
	if err != nil {
		return err
	}

	msg := bank.NewMsgSend(fromAddr, toAddr, types.NewCoins(types.NewInt64Coin(config.Denom, config.AmountSend)))

	encCfg := simapp.MakeTestEncodingConfig() // TODO:
	txBuilder := encCfg.TxConfig.NewTxBuilder()

	err = txBuilder.SetMsgs(msg)
	if err != nil {
		return err
	}
	txBuilder.SetGasLimit(config.GasLimit)
	txBuilder.SetMemo(config.Memo)
	txBuilder.SetFeeAmount(types.NewCoins(types.NewInt64Coin(config.Denom, config.FeeAmount)))

	grpcConn, _ := grpc.Dial(
		config.GrpcAddress,
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()

	account, err := GetAccount(grpcConn, fromAddr.String())
	if err != nil {
		return err
	}

	signMode := encCfg.TxConfig.SignModeHandler().DefaultMode()

	pubKey := fromPrivKey.PubKey()
	signerData := authsigning.SignerData{
		ChainID:       config.ChainID,
		AccountNumber: account.GetAccountNumber(),
		Sequence:      account.GetSequence(),
	}

	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: account.Sequence,
	}
	var prevSignatures []signing.SignatureV2

	if err := txBuilder.SetSignatures(sig); err != nil {
		return err
	}

	// Generate the bytes to be signed.
	bytesToSign, err := encCfg.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return err
	}

	// Sign those bytes
	sigBytes, err := fromPrivKey.Sign(bytesToSign)
	if err != nil {
		return err
	}

	// Construct the SignatureV2 struct
	sigData = signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}
	sig = signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: account.Sequence,
	}

	prevSignatures = append(prevSignatures, sig)
	if err := txBuilder.SetSignatures(prevSignatures...); err != nil {
		return err
	}

	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	txJSONBytes, err := encCfg.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}
	txJSON := string(txJSONBytes)
	fmt.Println(txJSON)

	txClient := tx.NewServiceClient(grpcConn)
	grpcRes, err := txClient.BroadcastTx(
		context.Background(),
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

func GeneratePrivateKey(mnemonic string) (crypto.PrivKey, error) {
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

func GetAccount(grpcConn *grpc.ClientConn, address string) (*auth.BaseAccount, error) {
	authClient := auth.NewQueryClient(grpcConn)
	query, err := authClient.Account(context.Background(), &auth.QueryAccountRequest{Address: address})
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
