package cmd

import (
	"okp4/cosmos-faucet/graph"
	"okp4/cosmos-faucet/graph/model"
	"okp4/cosmos-faucet/internal/server"
	"okp4/cosmos-faucet/pkg/actor/message"
	"okp4/cosmos-faucet/pkg/captcha"
	"okp4/cosmos-faucet/pkg/cosmos"
	"okp4/cosmos-faucet/pkg/faucet"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	FlagAddress       = "address"
	FlagBatchWindow   = "batch-window"
	FlagMetrics       = "metrics"
	FlagHealth        = "health"
	FlagCaptchaSecret = "captcha-secret"
	FlagCaptchaURL    = "captcha-verify-url"
	FlagCaptchaScore  = "captcha-min-score"
	FlagEnableCaptcha = "captcha"
)

// NewStartCommand returns a CLI command to start the REST api allowing to send tokens.
// nolint: funlen
func NewStartCommand() *cobra.Command {
	var addr string
	var batchWindow time.Duration
	var metrics bool
	var health bool
	var captchaConf captcha.ResolverConfig

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the GraphQL api",
		Run: func(cmd *cobra.Command, args []string) {
			conf := types.GetConfig()
			conf.SetBech32PrefixForAccount(prefix, prefix)

			privKey, err := cosmos.ParseMnemonic(mnemonic)
			if err != nil {
				log.Panic().Err(err).Msg("❌ Could not parse mnemonic")
			}

			cosmosClientProps := actor.PropsFromProducer(func() actor.Actor {
				grpcClient, err := cosmos.NewGrpcClient(grpcAddress, getTransportCredentials())
				if err != nil {
					log.Panic().Err(err).Msg("❌ Could not create grpc client")
				}

				return grpcClient
			})

			txHandlerProps := actor.PropsFromProducer(func() actor.Actor {
				return cosmos.NewTxHandler(
					cosmos.WithChainID(chainID),
					cosmos.WithPrivateKey(privKey),
					cosmos.WithTxConfig(simapp.MakeTestEncodingConfig().TxConfig),
					cosmos.WithCosmosClientProps(cosmosClientProps),
				)
			})

			actorCTX := actor.NewActorSystem().Root
			faucetPID := actorCTX.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return faucet.NewFaucet(
					faucet.WithChainID(chainID),
					faucet.WithAmount(types.NewCoins(types.NewInt64Coin(denom, amountSend))),
					faucet.WithAddress(types.AccAddress(privKey.PubKey().Address())),
					faucet.WithTxHandlerProps(txHandlerProps),
				)
			}))

			graphqlResolver := &graph.Resolver{
				Faucet:          faucetPID,
				Context:         actorCTX,
				AddressPrefix:   prefix,
				CaptchaResolver: captcha.NewCaptchaResolver(captchaConf),
				Config: &model.Configuration{
					AmountSend: amountSend,
					ChainID:    chainID,
					Denom:      denom,
					FeeAmount:  feeAmount,
					GasLimit:   gasLimit,
					Memo:       memo,
					Prefix:     prefix,
				},
			}

			go func() {
				for range time.Tick(batchWindow) {
					actorCTX.Send(faucetPID, message.TriggerTx{
						Deadline:  time.Now().Add(txTimeout),
						Memo:      memo,
						GasLimit:  gasLimit,
						FeeAmount: types.NewCoins(types.NewInt64Coin(denom, feeAmount)),
					})
				}
			}()

			server.NewServer(graphqlResolver, health, metrics).Start(addr)
		},
	}

	startCmd.Flags().StringVar(&addr, FlagAddress, ":8080", "graphql api address")
	startCmd.Flags().DurationVar(
		&batchWindow,
		FlagBatchWindow,
		8*time.Second,
		"Batch temporal window, can be seen a the minimum duration between too transactions.",
	)
	startCmd.Flags().BoolVar(&metrics, FlagMetrics, false, "enable metrics endpoint")
	startCmd.Flags().BoolVar(&health, FlagHealth, false, "enable health endpoint")
	startCmd.Flags().BoolVar(
		&captchaConf.Enable,
		FlagEnableCaptcha,
		false,
		"enable captcha verification",
	)
	startCmd.Flags().StringVar(
		&captchaConf.Secret,
		FlagCaptchaSecret,
		"",
		"set Captcha secret",
	)
	startCmd.Flags().StringVar(
		&captchaConf.VerifyURL,
		FlagCaptchaURL,
		"https://www.google.com/recaptcha/api/siteverify",
		"set Captcha verify URL",
	)
	startCmd.Flags().Float64Var(
		&captchaConf.MinScore,
		FlagCaptchaScore,
		0.5,
		"set Captcha min score",
	)

	return startCmd
}

func init() {
	rootCmd.AddCommand(NewStartCommand())
}
