package graph

import (
	"testing"

	"okp4/cosmos-faucet/graph/generated"
	"okp4/cosmos-faucet/graph/model"
	"okp4/cosmos-faucet/pkg"
	"okp4/cosmos-faucet/pkg/client"

	gql "github.com/99designs/gqlgen/client"

	"github.com/99designs/gqlgen/graphql/handler"
	. "github.com/smartystreets/goconvey/convey"
)

func TestQueryResolver_Configuration(t *testing.T) {
	Convey("Given a faucet configuration context to the resolver", t, func() {
		config := pkg.Config{
			GrpcAddress: "127.0.0.1:9090",
			ChainID:     "my-chain",
			Denom:       "denom",
			AmountSend:  58,
			FeeAmount:   74,
			Memo:        "my memo",
			Prefix:      "okp4",
			GasLimit:    42,
			Mnemonic:    "nasty random alter chronic become keen stadium test chaos fashion during claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard",
		}

		Convey("When create query context with faucet and configuration", func() {
			faucet, err := client.NewFaucet(config)

			srv := gql.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{Faucet: faucet}})))

			var result struct {
				Configuration model.Configuration
			}
			q := `
                query {
                    configuration {
                        chainId
                        denom
                        prefix
                        amountSend
                        feeAmount
                        memo
                        gasLimit
                    }
                }
                `
			srv.MustPost(q, &result)

			Convey("Configuration should be the same as the given server initialisation", func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				So(result.Configuration.ChainID, ShouldEqual, config.ChainID)
				So(result.Configuration.Denom, ShouldEqual, config.Denom)
				So(result.Configuration.AmountSend, ShouldEqual, config.AmountSend)
				So(result.Configuration.FeeAmount, ShouldEqual, config.FeeAmount)
				So(result.Configuration.Memo, ShouldEqual, config.Memo)
				So(result.Configuration.Prefix, ShouldEqual, config.Prefix)
				So(result.Configuration.GasLimit, ShouldEqual, config.GasLimit)
			})
		})
	})
}
