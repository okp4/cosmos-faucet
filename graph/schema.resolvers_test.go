package graph

import (
	"context"
	"okp4/cosmos-faucet/pkg/captcha"
	"testing"

	"okp4/cosmos-faucet/graph/generated"
	"okp4/cosmos-faucet/graph/model"
	"okp4/cosmos-faucet/pkg"
	"okp4/cosmos-faucet/pkg/client"

	gql "github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/wingyplus/must"
)

var config = pkg.Config{
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

type mockCaptchaResolver struct{}

func newMockCaptchaResolver() captcha.Resolver {
	return mockCaptchaResolver{}
}

func (r mockCaptchaResolver) CheckRecaptcha(_ context.Context, _ *string) error {
	return nil
}

func TestQueryResolver_Configuration(t *testing.T) {
	Convey("Given a faucet service configured to succeed", t, func() {
		faucet := must.Must(client.NewFaucet(config, nil))

		srv := gql.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{Faucet: faucet, CaptchaResolver: newMockCaptchaResolver()}})))

		Convey("And a a graphQL 'configuration' query request", func() {
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
			Convey("When posting the graphQL request", func() {
				var result struct {
					Configuration model.Configuration
				}

				err := srv.Post(q, &result)

				Convey("Then the returned configuration should be the same as the given server initialisation", func() {
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
	})
}
