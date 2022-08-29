package message

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// RequestFunds represents a message to request funds.
type RequestFunds struct {
	// Address on which to send requested funds.
	Address types.AccAddress

	// TxSubscriber denotes an actor on which to forward the response of the submitted transaction containing the
	// associated send message (i.e. BroadcastTxResponse).
	TxSubscriber *actor.PID
}

// TriggerTx represents a message trigger to process of submitting a transaction to the blockchain.
type TriggerTx struct {
	// Deadline the deadline before which the transaction shall be submitted.
	Deadline time.Time

	// Memo is the 'memo' field content of the transaction.
	Memo string

	// GasLimit allowed on the transaction.
	GasLimit uint64

	// FeeAmount to set on the transaction.
	FeeAmount types.Coins
}

// MakeTx represents a message to build, sign and submit a transaction.
type MakeTx struct {
	// Deadline the deadline before which the transaction shall be submitted.
	Deadline time.Time

	// TxSubscriberPID denotes an actor on which to forward the response of the submitted transaction containing the
	// associated send message (i.e. BroadcastTxResponse).
	TxSubscriber *actor.PID

	// Msgs contains the messages to embed in the transaction.
	Msgs []types.Msg

	// Memo is the 'memo' field content of the transaction.
	Memo string

	// GasLimit allowed on the transaction.
	GasLimit uint64

	// FeeAmount to set on the transaction.
	FeeAmount types.Coins
}

type GetAccount struct {
	// Deadline the deadline before which the account shall be retrieved.
	Deadline time.Time

	// Address on which to retrieve the account.
	Address string
}

type GetAccountResponse struct {
	// Account is the account response.
	Account *auth.BaseAccount
}

// BroadcastTx represents a message to submit a transaction against the blockchain.
type BroadcastTx struct {
	// Deadline the deadline before which the transaction shall be submitted.
	Deadline time.Time

	// Tx the raw signed transaction.
	Tx []byte
}

// BroadcastTxResponse represents a message emitted in response to BroadcastTx containing the transaction response, if
// successful.
type BroadcastTxResponse struct {
	// TxResponse is the submitted transaction response.
	TxResponse *types.TxResponse
}
