// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

// Represent the actual server configuration
type Configuration struct {
	// Amount value of token to send
	AmountSend int64 `json:"amountSend"`
	// The network chain ID
	ChainID string `json:"chainId"`
	// Token denom
	Denom string `json:"denom"`
	// Fee amount allowed
	FeeAmount int64 `json:"feeAmount"`
	// Gas limit allowed on transaction
	GasLimit uint64 `json:"gasLimit"`
	// Memo used when send transaction
	Memo string `json:"memo"`
	// Address prefix
	Prefix string `json:"prefix"`
}

// All inputs needed to send token to a given address
type SendInput struct {
	// Captcha token
	CaptchaToken *string `json:"captchaToken"`
	// Address where to send token(s)
	ToAddress string `json:"toAddress"`
}

// Represent a transaction response
type TxResponse struct {
	// Return the result code of transaction.
	// See code correspondence error : https://github.com/cosmos/cosmos-sdk/blob/main/types/errors/errors.go
	Code int `json:"code"`
	// Transaction gas used
	GasUsed int64 `json:"gasUsed"`
	// Transaction gas wanted
	GasWanted int64 `json:"gasWanted"`
	// Corresponding to the transaction hash.
	Hash string `json:"hash"`
	// Description of error if available.
	RawLog *string `json:"rawLog"`
}
