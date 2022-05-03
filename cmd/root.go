package cmd

import (
	"errors"
	"fmt"
	"os"

	"okp4/cosmos-faucet/pkg"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConfigFilename = "config"
	envPrefix             = "FAUCET"
)

var config pkg.Config

// NewRootCommand returns the root CLI command with persistent flag handling.
var rootCmd = &cobra.Command{
	Use:   "cosmos-faucet",
	Short: "A CØSMOS Faucet",
	Long:  "CØSMOS Faucet",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ReadPersistentFlags(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentFlags().StringVar(&config.Mnemonic, FlagMnemonic, "", "")
	rootCmd.PersistentFlags().StringVar(&config.ChainID, FlagChainID, "localnet-okp4-1", "The network chain ID")
	rootCmd.PersistentFlags().StringVar(&config.GrpcAddress, FlagGrpcAddress, "127.0.0.1:9090", "The grpc okp4 server url")
	rootCmd.PersistentFlags().StringVar(&config.Denom, FlagDenom, "know", "Token denom")
	rootCmd.PersistentFlags().StringVar(&config.Prefix, FlagPrefix, "okp4", "Address prefix")
	rootCmd.PersistentFlags().Int64Var(&config.FeeAmount, FlagFeeAmount, 0, "Fee amount")
	rootCmd.PersistentFlags().Int64Var(&config.AmountSend, FlagAmountSend, 1, "Amount send value")
	rootCmd.PersistentFlags().StringVar(&config.Memo, FlagMemo, "Sent by økp4 faucet", "The memo description")
	rootCmd.PersistentFlags().Uint64Var(&config.GasLimit, FlagGasLimit, 200000, "Gas limit")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ReadPersistentFlags(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) {
			return err
		}
	}
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			var _ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	return nil
}
