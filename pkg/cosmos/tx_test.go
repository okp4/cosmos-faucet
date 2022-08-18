package cosmos

import (
	"testing"

	"okp4/cosmos-faucet/pkg"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wingyplus/must"
)

func newTxConfig() client.TxConfig {
	enc := simapp.MakeTestEncodingConfig()
	return enc.TxConfig
}

func TestBuildUnsignedTx(t *testing.T) {
	Convey("Given a config file and an address", t, func() {
		config := pkg.Config{
			Denom:      "know",
			Prefix:     "okp4",
			FeeAmount:  20,
			AmountSend: 10,
			Memo:       "memo",
			GasLimit:   2000,
		}

		Convey("When building an unsigned transaction", func() {
			unsignedTx, err := BuildUnsignedTx(config, newTxConfig(), nil)

			Convey("Then the transaction should be successfully built", func() {
				So(err, ShouldBeNil)
				So(unsignedTx, ShouldNotBeNil)
			})

			Convey("And the transaction should be correctly configured", func() {
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
		privKey := must.Must(ParseMnemonic("nasty random alter chronic become keen stadium test chaos fashion during claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard"))
		config := pkg.Config{
			Denom:      "know",
			Prefix:     "okp4",
			FeeAmount:  20,
			AmountSend: 10,
			Memo:       "memo",
			GasLimit:   2000,
		}

		Convey("When signing the transaction", func() {
			signerData := signing.SignerData{
				ChainID:       "chain-id",
				AccountNumber: 10,
				Sequence:      54,
			}
			tx, _ := BuildUnsignedTx(config, newTxConfig(), nil)
			err := SignTx(privKey, signerData, newTxConfig(), tx)

			Convey("Then the transaction should be signed", func() {
				So(err, ShouldBeNil)

				signatures, err := tx.GetTx().GetSignaturesV2()
				So(signatures, ShouldNotBeEmpty)
				So(err, ShouldBeNil)

				Convey("And the signature should be the same than the signer public key", func() {
					So(signatures[0].PubKey.String(), ShouldEqual, privKey.PubKey().String())
				})

				Convey("And the account sequence and number should be correctly set", func() {
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
