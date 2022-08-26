package cmd

import (
	"okp4/cosmos-faucet/pkg/actor/message"
	"okp4/cosmos-faucet/pkg/actor/system"
	"okp4/cosmos-faucet/pkg/cosmos"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
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

			actorCTX, faucetPID := system.BootstrapActors(
				chainID,
				privKey,
				types.NewCoins(types.NewInt64Coin(denom, amountSend)),
				grpcAddress,
				getTransportCredentials(),
			)

			wg := sync.WaitGroup{}
			wg.Add(1)
			subPID := actorCTX.Spawn(actor.PropsFromFunc(func(c actor.Context) {
				if _, ok := c.Message().(*message.BroadcastTxResponse); ok {
					wg.Done()
					c.Stop(c.Self())
				}
			}))

			actorCTX.Send(faucetPID, &message.RequestFunds{
				Address:      toAddress,
				TxSubscriber: subPID,
			})
			actorCTX.Send(faucetPID, &message.TriggerTx{
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
