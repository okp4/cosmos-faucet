package cosmos

import (
	"testing"

	"okp4/cosmos-faucet/pkg"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
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
