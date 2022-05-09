package cosmos

import (
	"testing"

	"okp4/cosmos-faucet/pkg"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	. "github.com/smartystreets/goconvey/convey"
)

func newTxConfig() client.TxConfig {
	enc := simapp.MakeTestEncodingConfig()
	return enc.TxConfig
}

func TestBuildUnsignedTx(t *testing.T) {
	Convey("Given a config file and address", t, func() {
		config := pkg.Config{
			Denom:      "know",
			Prefix:     "okp4",
			FeeAmount:  20,
			AmountSend: 10,
			Memo:       "memo",
			GasLimit:   2000,
		}
		fromAddr := types.AccAddress("okp4AAAAA")
		toAddr := types.AccAddress("okp4BBB")

		Convey("When build unsigned transaction", func() {
			unsignedTx, err := BuildUnsignedTx(config, newTxConfig(), fromAddr, toAddr)

			Convey("Unsigned transaction should be successfully build", func() {
				So(err, ShouldBeNil)
				So(unsignedTx, ShouldNotBeNil)
			})

			Convey("Transaction should be correctly configured", func() {
				So(unsignedTx.GetTx().GetGas(), ShouldEqual, config.GasLimit)
				So(unsignedTx.GetTx().GetFee().String(),
					ShouldEqual,
					types.NewCoins(types.NewInt64Coin(config.Denom, config.FeeAmount)).String())
				So(unsignedTx.GetTx().GetMemo(), ShouldEqual, config.Memo)
			})
		})
	})
}

func TestSignTx(t *testing.T) {
	Convey("Given a private key, account and message", t, func() {
		mnemonic := "nasty random alter chronic become keen stadium test chaos fashion during claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard"
		privKey, _ := GeneratePrivateKey(mnemonic)
		config := pkg.Config{
			Denom:      "know",
			Prefix:     "okp4",
			FeeAmount:  20,
			AmountSend: 10,
			Memo:       "memo",
			GasLimit:   2000,
		}
		fromAddr := types.AccAddress("okp4AAAAA")
		toAddr := types.AccAddress("okp4BBB")
		signerData := signing.SignerData{
			ChainID:       "chain-id",
			AccountNumber: 10,
			Sequence:      54,
		}
		tx, _ := BuildUnsignedTx(config, newTxConfig(), fromAddr, toAddr)

		Convey("When sign transaction", func() {
			err := SignTx(privKey, signerData, newTxConfig(), tx)

			Convey("Sign transaction should be successful", func() {
				So(err, ShouldBeNil)
				signatures, err := tx.GetTx().GetSignaturesV2()
				So(signatures, ShouldNotBeEmpty)
				So(err, ShouldBeNil)

				Convey("Signature ppublic key should be the same as the signer", func() {
					So(signatures[0].PubKey.String(), ShouldEqual, privKey.PubKey().String())
				})

				Convey("Account sequence and number should be correctly set", func() {
					So(signatures[0].Sequence, ShouldEqual, signerData.Sequence)
				})
			})
		})
	})
}

func TestGeneratePrivateKey(t *testing.T) {
	Convey("Given a well formed mnemonic", t, func() {
		mnemonic := "nasty random alter chronic become keen stadium test chaos fashion during claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard"

		Convey("When generating its corresponding private key", func() {
			privKey, err := GeneratePrivateKey(mnemonic)

			Convey("The private key have been successfully decoded", func() {
				So(privKey, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given a malformed mnemonic", t, func() {
		mnemonic := "nasty random alter chronic become keen stadium test chaos fashion durin claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard"

		Convey("When generating its corresponding private key", func() {
			privKey, err := GeneratePrivateKey(mnemonic)

			Convey("The private key shall not be decoded", func() {
				So(privKey, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err, ShouldBeError, "Invalid mnemonic")
			})
		})
	})
}
