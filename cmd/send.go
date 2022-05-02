package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"okp4/cosmos-faucet/pkg/send"
	"okp4/cosmos-faucet/util"
)

const (
	defaultConfigFilename = "config"
	envPrefix             = "FAUCET"
)

// NewSendCommand returns a CLI command to interactively send amount token(s) to given address.
func NewSendCommand() *cobra.Command {
	var config util.Config

	sendCmd := &cobra.Command{
		Use:   "send <address>",
		Short: "Send tokens to a given address",
		Args:  cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return send.SendTx(config, args[0])
		},
	}

	sendCmd.Flags().StringVar(&config.Mnemonic, "mnemonic", "", "")
	sendCmd.Flags().StringVar(&config.ChainId, "chain-id", "okp4", "The network chain ID")
	sendCmd.Flags().StringVar(&config.GrpcAddress, "grpcAddress", "127.0.0.1:9090", "The grpc okp4 server url")
	sendCmd.Flags().StringVar(&config.Denom, "denom", "know", "Token denom")
	sendCmd.Flags().StringVar(&config.Prefix, "prefix", "okp4", "Address prefix")
	sendCmd.Flags().Int64Var(&config.FeeAmount, "fee-amount", 1000, "Fee amount") // TODO: Determine the default value
	sendCmd.Flags().Int64Var(&config.AmountSend, "amount-send", 1, "Amount send value")
	sendCmd.Flags().StringVar(&config.Memo, "memo", "Sent by økp4 faucet", "The memo description")
	sendCmd.Flags().Uint64Var(&config.GasLimit, "gas-limit", 200000, "Gas limit")

	return sendCmd
}

func init() {
	rootCmd.AddCommand(NewSendCommand())
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	return nil
}
