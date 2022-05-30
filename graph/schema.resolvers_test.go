package graph

import (
	"context"
	"errors"
	"fmt"
	"okp4/cosmos-faucet/pkg/captcha"
	"testing"

	"okp4/cosmos-faucet/graph/generated"
	"okp4/cosmos-faucet/graph/model"
	"okp4/cosmos-faucet/pkg"
	"okp4/cosmos-faucet/pkg/client"

	gql "github.com/99designs/gqlgen/client"
	"github.com/cosmos/cosmos-sdk/types"

	"github.com/99designs/gqlgen/graphql/handler"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/wingyplus/must"
)

const (
	graphQLRequestSendWithIncorrectAddress = `
                mutation {
                    send(input: {
                        captchaToken: "token"
                        toAddress: "wrong formatted address"
                    }) {
                        hash
                    }
                }
                `
	graphQLRequestSend = `
                mutation {
                    send(input: {
                        captchaToken: "token"
                        toAddress: "okp41jse8senm9hcvydhl8v9x47kfe5z82zmwtw8jvj"
                    }) {
                        hash
                        code
                        rawLog
                        gasWanted
                        gasUsed
                    }
                }`
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

type mockFaucet struct {
	config    pkg.Config
	withError bool
}

func (f mockFaucet) GetConfig() pkg.Config {
	return config
}

func (f mockFaucet) GetFromAddr() types.AccAddress {
	bech32, _ := types.AccAddressFromBech32("okp41jse8senm9hcvydhl8v9x47kfe5z82zmwtw8jvj")
	return bech32
}

func (f mockFaucet) Close() error {
	panic("implement me")
}

type mockCaptchaResolver struct{}

func newMockCaptchaResolver() captcha.Resolver {
	return mockCaptchaResolver{}
}

func (r mockCaptchaResolver) CheckRecaptcha(_ context.Context, _ *string) error {
	return nil
}

func (f mockFaucet) SendTxMsg(_ context.Context, _ string) (*types.TxResponse, error) {
	var code uint32
	if f.withError {
		code = 12
	} else {
		code = 0
	}

	return &types.TxResponse{
		Height:    0,
		TxHash:    "HASH",
		Codespace: "",
		Code:      code,
		Data:      "",
		RawLog:    "",
		Logs:      nil,
		Info:      "",
		GasWanted: 10,
		GasUsed:   20,
		Tx:        nil,
		Timestamp: "",
		Events:    nil,
	}, nil
}

func TestMutationResolver_Send(t *testing.T) {
	cases := []struct {
		name              string
		request           string
		faucet            client.Faucet
		expectedError     string
		expectedCode      int
		expectedHash      string
		expectedGasUsed   int64
		expectedGasWanted int64
	}{
		{
			name:          "a Faucet service configured to succeed and a graphQL 'send' mutation request with an incorrect address",
			request:       graphQLRequestSendWithIncorrectAddress,
			faucet:        must.Must(client.NewFaucet(config)),
			expectedError: "decoding bech32 failed: invalid character in string: ' '",
		},
		{
			name:              "a Faucet service configured to succeed and a correct graphQL 'send' mutation",
			request:           graphQLRequestSend,
			faucet:            mockFaucet{config: config},
			expectedHash:      "HASH",
			expectedGasUsed:   20,
			expectedGasWanted: 10,
		},
		{
			name:              "a Faucet service configured to fail and a correct graphQL 'send' mutation",
			request:           graphQLRequestSend,
			faucet:            mockFaucet{config: config, withError: true},
			expectedError:     `transaction is not successful:  (code: 12)","path":["send"]}]`,
			expectedCode:      12,
			expectedHash:      "HASH",
			expectedGasUsed:   20,
			expectedGasWanted: 10,
		},
	}

	for n, c := range cases {
		Convey(fmt.Sprintf("Given %s (case %d)", c.name, n), t, func() {
			srv := gql.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{Faucet: c.faucet, CaptchaResolver: newMockCaptchaResolver()}})))

			Convey("When posting the graphQL request", func() {
				var result struct {
					Send model.TxResponse
				}
				err := srv.Post(c.request, &result)

				Convey("Then the request should meet expectations", func() {
					if c.expectedError != "" {
						So(err, ShouldNotBeNil)

						var jsonError gql.RawJsonError
						So(errors.As(err, &jsonError), ShouldBeTrue)
						So(jsonError.Error(), ShouldContainSubstring, c.expectedError)
					} else {
						So(err, ShouldBeNil)
					}

					So(result.Send.Code, ShouldEqual, c.expectedCode)
					So(result.Send.Hash, ShouldEqual, c.expectedHash)
					So(result.Send.GasUsed, ShouldEqual, c.expectedGasUsed)
					So(result.Send.GasWanted, ShouldEqual, c.expectedGasWanted)
				})
			})
		})
	}
}

func TestQueryResolver_Configuration(t *testing.T) {
	Convey("Given a faucet service configured to succeed", t, func() {
		faucet := must.Must(client.NewFaucet(config))

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
