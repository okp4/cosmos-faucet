package faucet

import (
	"okp4/cosmos-faucet/pkg/actor/message"
	"okp4/cosmos-faucet/test/mock"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/stretchr/testify/mock"
)

var (
	amount   = types.NewCoins(types.NewInt64Coin("uknow", 1000000))
	fromAddr = types.AccAddress("from")
	toAddr   = types.AccAddress("to")
)

func TestOptions(t *testing.T) {
	Convey("Given a set of options", t, func() {
		opts := []Option{
			WithAddress(fromAddr),
			WithAmount(amount),
			WithTxHandlerProps(&actor.Props{}),
		}

		Convey("When creating the faucet with the options", func() {
			faucet := NewFaucet(opts...)

			Convey("Then the returned faucet should be configured accordingly", func() {
				So(faucet.address, ShouldResemble, fromAddr)
				So(faucet.amount, ShouldResemble, amount)
				So(faucet.txHandlerProps, ShouldResemble, &actor.Props{})
				So(faucet.txHandler, ShouldBeNil)
				So(faucet.msgs, ShouldBeNil)
				So(faucet.txSubscribers, ShouldBeNil)
			})
		})
	})
}

func TestStarted(t *testing.T) {
	Convey("Given a faucet actor", t, func() {
		faucet := NewFaucet(WithTxHandlerProps(actor.PropsFromFunc(func(c actor.Context) {})))

		Convey("When receiving a Started message", func() {
			mockedContext := &mock.ActorContext{}
			mockedContext.On("Message").Return(&actor.Started{})
			mockedContext.On("Spawn", Anything).Return(&actor.PID{})
			faucet.Receive(mockedContext)

			Convey("Then the TxHandler actor should be spawned", func() {
				mockedContext.AssertCalled(t, "Message")
				mockedContext.AssertCalled(t, "Spawn", Anything)
				So(faucet.txHandler, ShouldNotBeNil)
			})
		})
	})
}

func TestRequestFundsWithoutSubscriber(t *testing.T) {
	Convey("Given a faucet actor", t, func() {
		faucet := NewFaucet(
			WithAddress(fromAddr),
			WithAmount(amount),
			WithTxHandlerProps(actor.PropsFromFunc(func(c actor.Context) {})),
		)

		Convey("When receiving a RequestFunds message without subscriber", func() {
			mockedContext := &mock.ActorContext{}
			mockedContext.On("Message").Return(&message.RequestFunds{Address: toAddr})
			faucet.Receive(mockedContext)

			Convey("Then the send msg should be in the pool with no subscriber", func() {
				mockedContext.AssertCalled(t, "Message")
				So(len(faucet.msgs), ShouldEqual, 1)
				So(len(faucet.txSubscribers), ShouldEqual, 0)
				So(faucet.msgs[0], ShouldResemble, banktypes.NewMsgSend(fromAddr, toAddr, amount))
			})
		})
	})
}

func TestRequestFundsWithSubscriber(t *testing.T) {
	Convey("Given a faucet actor", t, func() {
		faucet := &Faucet{address: fromAddr, amount: amount}

		Convey("When receiving a RequestFunds message with subscriber", func() {
			mockedContext := &mock.ActorContext{}
			mockedContext.On("Message").Return(&message.RequestFunds{Address: toAddr, TxSubscriber: &actor.PID{}})
			faucet.Receive(mockedContext)

			Convey("Then the send msg should be in the pool with a subscriber", func() {
				mockedContext.AssertCalled(t, "Message")
				So(len(faucet.msgs), ShouldEqual, 1)
				So(len(faucet.txSubscribers), ShouldEqual, 1)
				So(faucet.msgs[0], ShouldResemble, banktypes.NewMsgSend(fromAddr, toAddr, amount))
			})
		})
	})
}

func TestTriggerTxWithoutMsgs(t *testing.T) {
	Convey("Given a faucet actor", t, func() {
		faucet := &Faucet{}

		Convey("When receiving a TriggerTx message", func() {
			triggerMsg := message.TriggerTx{
				Deadline:  time.Now(),
				Memo:      "Sent from tests",
				GasLimit:  200000,
				FeeAmount: amount,
			}
			mockedContext := &mock.ActorContext{}
			mockedContext.On("Message").Return(&triggerMsg)
			faucet.Receive(mockedContext)

			Convey("Then the send msg should be in the pool with no subscriber", func() {
				mockedContext.AssertCalled(t, "Message")
				mockedContext.AssertNotCalled(t, "Spawn", Anything)
				mockedContext.AssertNotCalled(t, "Send", Anything, Anything)
				So(len(faucet.msgs), ShouldEqual, 0)
				So(len(faucet.txSubscribers), ShouldEqual, 0)
			})
		})
	})
}

func TestTriggerTxWithMsgs(t *testing.T) {
	Convey("Given a faucet actor", t, func() {
		txMsgs := []types.Msg{
			banktypes.NewMsgSend(fromAddr, toAddr, amount),
			banktypes.NewMsgSend(fromAddr, toAddr, amount),
		}
		faucet := &Faucet{
			msgs:          txMsgs,
			txSubscribers: []*actor.PID{{}, {}},
			txHandler:     &actor.PID{Id: "txHandler"},
		}

		Convey("When receiving a TriggerTx message", func() {
			var messageSent interface{}
			triggerMsg := message.TriggerTx{
				Deadline:  time.Now(),
				Memo:      "Sent from tests",
				GasLimit:  200000,
				FeeAmount: amount,
			}
			mockedContext := &mock.ActorContext{}
			mockedContext.On("Message").Return(&triggerMsg)
			mockedContext.On("Spawn", Anything).Return(&actor.PID{Id: "subscriber"})
			mockedContext.On("Send", Anything, Anything).Run(func(args Arguments) {
				messageSent = args.Get(1)
			}).Return()
			faucet.Receive(mockedContext)

			Convey("Then the send msg should be in the pool with no subscriber", func() {
				mockedContext.AssertCalled(t, "Message")
				mockedContext.AssertCalled(t, "Spawn", Anything)
				mockedContext.AssertCalled(t, "Send", &actor.PID{Id: "txHandler"}, Anything)
				So(len(faucet.msgs), ShouldEqual, 0)
				So(len(faucet.txSubscribers), ShouldEqual, 0)
				So(messageSent, ShouldHaveSameTypeAs, &message.MakeTx{})
				So(messageSent.(*message.MakeTx).Deadline, ShouldResemble, triggerMsg.Deadline)
				So(messageSent.(*message.MakeTx).TxSubscriber.Id, ShouldEqual, "subscriber")
				So(messageSent.(*message.MakeTx).Msgs, ShouldResemble, txMsgs)
				So(messageSent.(*message.MakeTx).Memo, ShouldResemble, "Sent from tests")
				So(messageSent.(*message.MakeTx).GasLimit, ShouldEqual, 200000)
				So(messageSent.(*message.MakeTx).FeeAmount, ShouldResemble, amount)
			})
		})
	})
}
