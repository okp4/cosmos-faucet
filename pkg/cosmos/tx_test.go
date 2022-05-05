package cosmos

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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
