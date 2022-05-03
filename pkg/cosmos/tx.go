package cosmos

import (
	"okp4/cosmos-faucet/pkg"

	sdk "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func BuildUnsignedTx(config pkg.Config, txConfig sdk.TxConfig, fromAddr, toAddr types.AccAddress) (sdk.TxBuilder, error) {
	msg := bank.NewMsgSend(fromAddr, toAddr, types.NewCoins(types.NewInt64Coin(config.Denom, config.AmountSend)))

	txBuilder := txConfig.NewTxBuilder()

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}
	txBuilder.SetGasLimit(config.GasLimit)
	txBuilder.SetMemo(config.Memo)
	txBuilder.SetFeeAmount(types.NewCoins(types.NewInt64Coin(config.Denom, config.FeeAmount)))

	return txBuilder, nil
}

func SignTx(
	config pkg.Config, fromPrivKey crypto.PrivKey,
	account *auth.BaseAccount, txConfig sdk.TxConfig, txBuilder sdk.TxBuilder) error {
	signMode := txConfig.SignModeHandler().DefaultMode()

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

	if err := txBuilder.SetSignatures(sig); err != nil {
		return err
	}

	// Generate the bytes to be signed.
	bytesToSign, err := txConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
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

	if err := txBuilder.SetSignatures(sig); err != nil {
		return err
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
