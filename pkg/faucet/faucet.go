package faucet

import (
	"okp4/cosmos-faucet/pkg/actor/message"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/router"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"
)

type Faucet struct {
	address        types.AccAddress
	amount         types.Coins
	txHandlerProps *actor.Props
	txHandler      *actor.PID
	msgs           []types.Msg
	txSubscribers  []*actor.PID
}

func NewFaucet(opts ...Option) *Faucet {
	faucet := &Faucet{}
	for _, opt := range opts {
		opt(faucet)
	}

	return faucet
}

type Option func(faucet *Faucet)

func WithAddress(address types.AccAddress) Option {
	return func(faucet *Faucet) {
		faucet.address = address
	}
}

func WithAmount(amount types.Coins) Option {
	return func(faucet *Faucet) {
		faucet.amount = amount
	}
}

func WithTxHandlerProps(props *actor.Props) Option {
	return func(faucet *Faucet) {
		faucet.txHandlerProps = props
	}
}

func (faucet *Faucet) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		faucet.txHandler = ctx.Spawn(faucet.txHandlerProps)

	case *message.RequestFunds:
		faucet.msgs = append(faucet.msgs, faucet.MakeSendMsg(msg.Address))
		if msg.TxSubscriber != nil {
			faucet.txSubscribers = append(faucet.txSubscribers, msg.TxSubscriber)
		}
		log.Info().Str("address", msg.Address.String()).Msg("✍️  Register fund request")

	case *message.TriggerTx:
		if len(faucet.msgs) == 0 {
			log.Info().Msg("😥 Ignore transaction trigger, no message to submit")
			break
		}

		log.Info().Time("deadline", msg.Deadline).Msg("🔥 Trigger new transaction")
		ctx.Send(faucet.txHandler, &message.MakeTx{
			Deadline:     msg.Deadline,
			TxSubscriber: ctx.Spawn(router.NewBroadcastGroup(faucet.txSubscribers...)),
			Msgs:         faucet.msgs,
			Memo:         msg.Memo,
			GasLimit:     msg.GasLimit,
			FeeAmount:    msg.FeeAmount,
		})
		faucet.msgs = faucet.msgs[:0]
		faucet.txSubscribers = faucet.txSubscribers[:0]
	}
}

func (faucet *Faucet) MakeSendMsg(addr types.AccAddress) types.Msg {
	return banktypes.NewMsgSend(
		faucet.address,
		addr,
		faucet.amount,
	)
}
