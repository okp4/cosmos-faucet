package client

import (
	"sync"

	"github.com/cosmos/cosmos-sdk/types"
)

// MessagePool atomically manage a pool of messages waiting to be sent through a transaction. It also allows to
// subscribe to the transaction response once submitted to receive it through a channel.
type MessagePool struct {
	mut         *sync.Mutex
	txFunc      TxSubmitter
	msgs        []*types.Msg
	subscribers []chan *types.TxResponse
}

// TxSubmitter shall implement all the logic to build, sign and submit a transaction containing all the messages of the
// pool.
type TxSubmitter func([]*types.Msg) (*types.TxResponse, error)

// MessagePoolOption allow to configure a MessagePool.
type MessagePoolOption func(pool *MessagePool)

// NewMessagePool initialize a new MessagePool instance, and configure it with the provided options.
func NewMessagePool(opts ...MessagePoolOption) *MessagePool {
	pool := &MessagePool{
		mut: &sync.Mutex{},
	}

	for _, opt := range opts {
		opt(pool)
	}

	return pool
}

// RegisterMsg atomically add the message in the pool.
func (pool *MessagePool) RegisterMsg(msg *types.Msg) {
	pool.lock()
	defer pool.unlock()

	pool.msgs = append(pool.msgs, msg)
}

// SubscribeMsg atomically add the message in the pool, it also returns a channel on which will be sent the
// corresponding transaction response, see Submit for more information on how the channel is managed.
func (pool *MessagePool) SubscribeMsg(msg *types.Msg) <-chan *types.TxResponse {
	pool.lock()
	defer pool.unlock()

	pool.msgs = append(pool.msgs, msg)

	txResponseChan := make(chan *types.TxResponse)
	pool.subscribers = append(pool.subscribers, txResponseChan)

	return txResponseChan
}

// Submit atomically submit a new transaction using configured the TxSubmitter embedding all the pooled messages and
// empty the pool. If an error occur there is no retry on the concerned messages.
//
// The subscribed channels will be closed following the transaction response, if an error occur submitting the
// transaction the channels are closed without any data.
//
// Warning: To avoid locking the MessagePool, the channels are closed in separate routines which can lead to goroutine
// leak if they are not consumed.
func (pool *MessagePool) Submit() error {
	pool.lock()
	defer func() {
		pool.flush()
		pool.unlock()
	}()

	resp, err := pool.txFunc(pool.msgs)
	if err != nil {
		return err
	}

	for _, subscriber := range pool.subscribers {
		subscriber <- resp
	}

	return nil
}

func (pool *MessagePool) lock() {
	pool.mut.Lock()
}

func (pool *MessagePool) unlock() {
	pool.mut.Unlock()
}

func (pool *MessagePool) flush() {
	for _, subscriber := range pool.subscribers {
		subscriber := subscriber
		go func() {
			close(subscriber)
		}()
	}

	pool.msgs = pool.msgs[:0]
	pool.subscribers = pool.subscribers[:0]
}
