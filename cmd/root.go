package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultConfigFilename = "config"
	envPrefix             = "FAUCET"
)

var (
	mnemonic      string
	chainID       string
	grpcAddress   string
	denom         string
	prefix        string
	feeAmount     int64
	amountSend    int64
	memo          string
	gasLimit      uint64
	noTLS         bool
	tlsSkipVerify bool
	txTimeout     time.Duration
)

// NewRootCommand returns the root CLI command with persistent flag handling.
var rootCmd = &cobra.Command{
	Use:   "cosmos-faucet",
	Short: "A CØSMOS Faucet",
	Long:  "CØSMOS Faucet",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: cmd.OutOrStdout()})
		return ReadPersistentFlags(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentFlags().StringVar(&mnemonic, FlagMnemonic, "", "")
	rootCmd.PersistentFlags().StringVar(&chainID, FlagChainID, "localnet-okp4-1", "The network chain ID")
	rootCmd.PersistentFlags().StringVar(&grpcAddress, FlagGrpcAddress, "127.0.0.1:9090", "The grpc okp4 server url")
	rootCmd.PersistentFlags().StringVar(&denom, FlagDenom, "know", "Token denom")
	rootCmd.PersistentFlags().StringVar(&prefix, FlagPrefix, "okp4", "Address prefix")
	rootCmd.PersistentFlags().Int64Var(&feeAmount, FlagFeeAmount, 0, "Fee amount")
	rootCmd.PersistentFlags().Int64Var(&amountSend, FlagAmountSend, 1, "Amount send value")
	rootCmd.PersistentFlags().StringVar(&memo, FlagMemo, "Sent by økp4 faucet", "The memo description")
	rootCmd.PersistentFlags().Uint64Var(&gasLimit, FlagGasLimit, 200000, "Gas limit")
	rootCmd.PersistentFlags().BoolVar(&noTLS, FlagNoTLS, false, "No encryption with the GRPC endpoint")
	rootCmd.PersistentFlags().BoolVar(&tlsSkipVerify,
		FlagTLSSkipVerify,
		false,
		"Encryption with the GRPC endpoint but skip certificates verification")
	rootCmd.PersistentFlags().DurationVar(&txTimeout, FlagTxTimeout, 5*time.Second, "Transaction timeout")

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
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			var _ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	return nil
}

func getTransportCredentials() credentials.TransportCredentials {
	switch {
	case noTLS:
		return insecure.NewCredentials()
	case tlsSkipVerify:
		return credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}) // #nosec G402 : skip lint since it's an optional flag
	default:
		return credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}
}
