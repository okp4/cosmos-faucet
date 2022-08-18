package cmd

import (
	"okp4/cosmos-faucet/pkg/actor/message"
	"okp4/cosmos-faucet/pkg/cosmos"
	"okp4/cosmos-faucet/pkg/faucet"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewSendCommand returns a CLI command to interactively send amount token(s) to given address.
func NewSendCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send <address>",
		Short: "Send tokens to a given address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			conf := types.GetConfig()
			conf.SetBech32PrefixForAccount(prefix, prefix)

			privKey, err := cosmos.ParseMnemonic(mnemonic)
			if err != nil {
				log.Panic().Err(err).Msg("❌ Could not parse mnemonic")
			}

			toAddress, err := types.GetFromBech32(args[0], prefix)
			if err != nil {
				log.Panic().Err(err).Str("toAddress", args[0]).Msg("❌ Could not parse address")
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

			wg := sync.WaitGroup{}
			wg.Add(1)
			subPID := actorCTX.Spawn(actor.PropsFromFunc(func(c actor.Context) {
				switch c.Message().(type) {
				case message.BroadcastTxResponse:
					wg.Done()
					c.Stop(c.Self())
				}
			}))

			actorCTX.Send(faucetPID, message.RequestFunds{
				Address:      toAddress,
				TxSubscriber: subPID,
			})
			actorCTX.Send(faucetPID, message.TriggerTx{
				Deadline:  time.Now().Add(txTimeout),
				Memo:      memo,
				GasLimit:  gasLimit,
				FeeAmount: types.NewCoins(types.NewInt64Coin(denom, feeAmount)),
			})

			wg.Wait()
			actorCTX.Stop(faucetPID)
		},
	}

	return sendCmd
}

func init() {
	rootCmd.AddCommand(NewSendCommand())
}
