package cosmos

import (
	"testing"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	signing2 "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mnemonic = "orient few plate lawsuit pumpkin shallow wild relax mosquito ill stool pitch orbit develop ginger bachelor slot useful worry arena bamboo combine arch envelope"
	chainID  = "testnet-okp4-1"
	txConfig = simapp.MakeTestEncodingConfig().TxConfig

	msgs = []types.Msg{
		banktypes.NewMsgSend(
			types.AccAddress("from"),
			types.AccAddress("to"),
			types.NewCoins(types.NewInt64Coin("uknow", 1000000)),
		),
		banktypes.NewMsgSend(
			types.AccAddress("from"),
			types.AccAddress("to"),
			types.NewCoins(types.NewInt64Coin("uknow", 5000000)),
		),
	}
	memo       = "Sent from tests"
	gasLimit   = uint64(200000)
	feeAmount  = types.NewCoins(types.NewInt64Coin("uknow", 50000))
	signerData = signing.SignerData{
		ChainID:       chainID,
		AccountNumber: 10,
		Sequence:      54,
	}
	signData = &signing2.SingleSignatureData{
		SignMode:  txConfig.SignModeHandler().DefaultMode(),
		Signature: []uint8{119, 228, 243, 201, 202, 94, 114, 31, 59, 53, 109, 198, 174, 101, 91, 211, 189, 4, 30, 39, 142, 187, 123, 199, 255, 129, 19, 91, 228, 229, 123, 202, 9, 154, 140, 241, 253, 186, 250, 153, 62, 178, 183, 126, 241, 8, 71, 161, 53, 153, 62, 202, 241, 22, 67, 132, 37, 78, 240, 41, 83, 157, 161, 16},
	}
)

func TestOptions(t *testing.T) {
	Convey("Given a set of options", t, func() {
		privKey, err := ParseMnemonic(mnemonic)
		if err != nil {
			panic(err)
		}
		opts := []Option{
			WithMnemonicMust(mnemonic),
			WithChainID(chainID),
			WithTxConfig(txConfig),
			WithCosmosClientProps(&actor.Props{}),
		}

		Convey("When creating the TxHandler with the options", func() {
			txHandler := NewTxHandler(opts...)

			Convey("Then the returned txHandler should be configured accordingly", func() {
				So(txHandler.privKey, ShouldResemble, privKey)
				So(txHandler.address, ShouldEqual, types.AccAddress(privKey.PubKey().Address()).String())
				So(txHandler.chainID, ShouldEqual, chainID)
				So(txHandler.config, ShouldResemble, txConfig)
				So(txHandler.signMode, ShouldResemble, txConfig.SignModeHandler().DefaultMode())
				So(txHandler.cosmosClientProps, ShouldResemble, &actor.Props{})
				So(txHandler.cosmosClient, ShouldBeNil)
			})
		})
	})
}

func TestBuildUnsignedTx(t *testing.T) {
	Convey("Given a TxHandler and some transaction parameters", t, func() {
		txHandler := NewTxHandler(
			WithMnemonicMust(mnemonic),
			WithChainID(chainID),
			WithTxConfig(txConfig),
		)

		Convey("When building an unsigned transaction", func() {
			unsignedTx, err := txHandler.BuildUnsignedTx(msgs, memo, gasLimit, feeAmount)

			Convey("Then the transaction should be successfully built", func() {
				So(err, ShouldBeNil)
				So(unsignedTx, ShouldNotBeNil)
			})

			Convey("And the transaction should be correctly configured", func() {
				So(unsignedTx.GetTx().GetGas(), ShouldEqual, gasLimit)
				So(unsignedTx.GetTx().GetFee().String(), ShouldEqual, feeAmount.String())
				So(unsignedTx.GetTx().GetMemo(), ShouldEqual, memo)
				So(unsignedTx.GetTx().GetMsgs(), ShouldResemble, msgs)
			})
		})
	})
}

func TestSignTx(t *testing.T) {
	Convey("Given a TxHandler and an unsigned transaction", t, func() {
		txHandler := NewTxHandler(
			WithMnemonicMust(mnemonic),
			WithChainID(chainID),
			WithTxConfig(txConfig),
		)
		unsignedTx, err := txHandler.BuildUnsignedTx(msgs, memo, gasLimit, feeAmount)
		if err != nil {
			panic(err)
		}

		Convey("When signing the transaction", func() {
			signedTx, err := txHandler.SignTx(unsignedTx, signerData)

			Convey("Then the transaction should be signed", func() {
				So(err, ShouldBeNil)

				signatures, err := signedTx.GetSignaturesV2()
				So(signatures, ShouldNotBeEmpty)
				So(err, ShouldBeNil)

				Convey("And the signature should be the same than the signer public key", func() {
					So(signatures[0].PubKey.String(), ShouldEqual, txHandler.privKey.PubKey().String())
				})

				Convey("And the account sequence should be correctly set", func() {
					So(signatures[0].Sequence, ShouldEqual, signerData.Sequence)
				})

				Convey("And the signature data shall correspond to the tx", func() {
					So(signatures[0].Data, ShouldResemble, signData)
				})
			})
		})
	})
}

func TestParseMnemonic(t *testing.T) {
	Convey("Given a well formed mnemonic", t, func() {
		mnemonic := "nasty random alter chronic become keen stadium test chaos fashion during claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard"

		Convey("When generating its corresponding private key", func() {
			privKey, err := ParseMnemonic(mnemonic)

			Convey("Then The private key have been successfully decoded", func() {
				So(privKey, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given a malformed mnemonic", t, func() {
		mnemonic := "nasty random alter chronic become keen stadium test chaos fashion durin claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard"

		Convey("When generating its corresponding private key", func() {
			privKey, err := ParseMnemonic(mnemonic)

			Convey("Then the private key shall not be decoded", func() {
				So(privKey, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err, ShouldBeError, "Invalid mnemonic")
			})
		})
	})
}
