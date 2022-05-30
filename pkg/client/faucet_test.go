package client

import (
	"crypto/tls"
	"errors"
	"testing"

	"okp4/cosmos-faucet/pkg"

	. "github.com/smartystreets/goconvey/convey"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNewFaucet(t *testing.T) {
	Convey("Given a correct configuration with grpc address", t, func() {
		grpcAddre := "127.0.0.1:9090"
		config := pkg.Config{
			Prefix:      "okp4",
			GrpcAddress: grpcAddre,
			Mnemonic:    "nasty random alter chronic become keen stadium test chaos fashion during claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard",
		}

		Convey("When creating the faucet service", func() {
			object, err := NewFaucet(config)
			faucet, ok := object.(*faucet)

			Convey("Then the faucet should be successfully created with the provided configuration", func() {
				So(faucet, ShouldNotBeNil)
				So(ok, ShouldBeTrue)
				So(err, ShouldBeNil)
				So(faucet.config, ShouldResemble, config)
			})

			Convey("And the GRPC connection should target the expected address", func() {
				So(faucet.grpcConn.Target(), ShouldEqual, grpcAddre)
			})

			Convey("And the faucet should be configured with a 'from' private key and 'from' address", func() {
				So(faucet.fromPrivKey, ShouldNotBeNil)
				So(faucet.fromAddr, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a configuration with a bad mnemonic", t, func() {
		grpcAddre := "127.0.0.1:9090"
		config := pkg.Config{
			Prefix:      "okp4",
			GrpcAddress: grpcAddre,
			Mnemonic:    "nasty random alter chronic become keen stadium test chaos fashion  rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard",
		}

		Convey("When creating the faucet service", func() {
			faucet, err := NewFaucet(config)

			Convey("Then the faucet creation should fail", func() {
				So(faucet, ShouldBeNil)
				So(err, ShouldResemble, errors.New("Invalid mnemonic"))
			})
		})
	})
}

func TestGetTransportCredentials(t *testing.T) {
	Convey("Given a configuration without specifying the 'tls' option", t, func() {
		config := pkg.Config{}

		Convey("When getting the transport credentials option", func() {
			opt := getTransportCredentials(config)

			Convey("Then the transport credentials should be set by default on TLS", func() {
				So(opt, ShouldResemble, credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12}))
			})
		})
	})

	Convey("Given a configuration specifying the 'no-tls' option", t, func() {
		config := pkg.Config{NoTLS: true}

		Convey("When getting the transport credentials option", func() {
			opt := getTransportCredentials(config)

			Convey("Then the transport credentials should be insecure", func() {
				So(opt, ShouldResemble, insecure.NewCredentials())
			})
		})
	})

	Convey("Given a configuration specifying the 'tls-skip-verify' option", t, func() {
		config := pkg.Config{TLSSkipVerify: true}

		Convey("When getting the transport credentials option", func() {
			opts := getTransportCredentials(config)

			Convey("Then the transport credentials should be set on TLS", func() {
				So(opts, ShouldResemble, credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})) // #nosec G402
			})
		})
	})
}
