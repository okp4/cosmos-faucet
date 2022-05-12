package graph

import "okp4/cosmos-faucet/pkg/client"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Faucet        client.Faucet
	CaptchaSecret string
}
