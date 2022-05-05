package client

import (
	"testing"

	"okp4/cosmos-faucet/pkg"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewFaucet(t *testing.T) {
	Convey("Given a good configuration with grpc address", t, func() {
		grpcAddre := "127.0.0.1:9090"
		config := pkg.Config{
			Prefix:      "okp4",
			GrpcAddress: grpcAddre,
			Mnemonic:    "nasty random alter chronic become keen stadium test chaos fashion during claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard",
		}

		Convey("When creating the new faucet", func() {
			faucet, err := NewFaucet(config)

			Convey("Faucet should be successfully created with given configuration", func() {
				So(faucet, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(faucet.Config, ShouldResemble, config)
			})

			Convey("Grpc connection should be target the good address", func() {
				So(faucet.GRPCConn.Target(), ShouldEqual, grpcAddre)
			})

			Convey("Faucet should be set with a from private key and from address", func() {
				So(faucet.FromPrivKey, ShouldNotBeNil)
				So(faucet.FromAddr, ShouldNotBeNil)
			})
		})
	})
}

func TestGetTransportCredentials(t *testing.T) {
	Convey("Given a config without specifying TLS", t, func() {
		config := pkg.Config{}

		Convey("When get transport credentials options", func() {
			opts := getTransportCredentials(config)

			Convey("Transport credentials should be set by default on TLS", func() {
				So(opts.Info().SecurityProtocol, ShouldEqual, "tls")
			})
		})
	})

	Convey("Given a config specifying the no-tls option", t, func() {
		config := pkg.Config{NoTLS: true}

		Convey("When get transport credentials options", func() {
			opts := getTransportCredentials(config)

			Convey("Transport credentials should be insecure", func() {
				So(opts.Info().SecurityProtocol, ShouldEqual, "insecure")
			})
		})
	})

	Convey("Given a config specifying the tls-skip-verify option", t, func() {
		config := pkg.Config{TLSSkipVerify: true}

		Convey("When get transport credentials options", func() {
			opts := getTransportCredentials(config)

			Convey("Transport credentials should be set on TLS", func() {
				So(opts.Info().SecurityProtocol, ShouldEqual, "tls")
			})
		})
	})
}
