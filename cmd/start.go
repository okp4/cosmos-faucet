package cmd

import (
	"okp4/cosmos-faucet/internal/server"
	"okp4/cosmos-faucet/pkg/client"
	"time"

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

var serverConfig server.Config

// NewStartCommand returns a CLI command to start the REST api allowing to send tokens.
func NewStartCommand() *cobra.Command {
	var addr string
	var batchWindow time.Duration

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the GraphQL api",
		Run: func(cmd *cobra.Command, args []string) {
			triggerTxChan := make(chan *client.TriggerTx)
			go func() {
				for range time.Tick(batchWindow) {
					triggerTxChan <- client.MakeTriggerTx(client.WithDeadline(time.Now().Add(config.TxTimeout)))
				}
			}()

			faucet, err := client.NewFaucet(config, triggerTxChan)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed create a new faucet instance")
			}

			defer func(faucet *client.Faucet) {
				_ = faucet.Close()
				log.Info().Msg("Server stopped")
			}(faucet)

			serverConfig.Faucet = faucet
			server.NewServer(serverConfig).Start(addr)
		},
	}

	startCmd.Flags().StringVar(&addr, FlagAddress, ":8080", "graphql api address")
	startCmd.Flags().DurationVar(
		&batchWindow,
		FlagBatchWindow,
		8*time.Second,
		"Batch temporal window, can be seen a the minimum duration between too transactions.",
	)
	startCmd.Flags().BoolVar(&serverConfig.EnableMetrics, FlagMetrics, false, "enable metrics endpoint")
	startCmd.Flags().BoolVar(&serverConfig.EnableHealth, FlagHealth, false, "enable health endpoint")
	startCmd.Flags().BoolVar(
		&serverConfig.CaptchaConf.Enable,
		FlagEnableCaptcha,
		false,
		"enable captcha verification",
	)
	startCmd.Flags().StringVar(
		&serverConfig.CaptchaConf.Secret,
		FlagCaptchaSecret,
		"",
		"set Captcha secret",
	)
	startCmd.Flags().StringVar(
		&serverConfig.CaptchaConf.VerifyURL,
		FlagCaptchaURL,
		"https://www.google.com/recaptcha/api/siteverify",
		"set Captcha verify URL",
	)
	startCmd.Flags().Float64Var(
		&serverConfig.CaptchaConf.MinScore,
		FlagCaptchaScore,
		0.5,
		"set Captcha min score",
	)

	return startCmd
}

func init() {
	rootCmd.AddCommand(NewStartCommand())
}
