package util

type Config struct {
	GrpcAddress string `mapstructure:"grpc-address"`
	Mnemonic    string `mapstructure:"mnemonic"`
	ChainId     string `mapstructure:"chain-id"`
	Denom       string `mapstructure:"denom"`
	Prefix      string `mapstructure:"prefix"`
	FeeAmount   int64  `mapstructure:"fee-amount"`
	AmountSend  int64  `mapstructure:"amount-send"`
	Memo        string `mapstructure:"memo"`
	GasLimit    uint64 `mapstructure:"gas-limit"`
}
