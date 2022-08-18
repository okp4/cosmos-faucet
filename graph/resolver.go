package graph

import (
	"okp4/cosmos-faucet/graph/model"
	"okp4/cosmos-faucet/pkg/captcha"

	"github.com/asynkron/protoactor-go/actor"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Faucet          *actor.PID
	Context         *actor.RootContext
	AddressPrefix   string
	CaptchaResolver captcha.Resolver
	Config          *model.Configuration
}
