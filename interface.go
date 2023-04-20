package rollup

import (
	"context"
)

// A sequencer is a process that is connected with other sequencers in a network.
// Collectively, they follow a protocol whereby concurrently submitted transactions
// across multiple processes are serialized. This could be single leader, 
// mulit-leader, or rely on some other mechanism. A sequencer is paired with an 
// Executor or execution environment that applies the transactions that have been
// agreed upon. One or more sequencers are responsible for publishing the data to
// a data availability layer which cannonicalizes the data, ensuring stateless
// participants in a network can still verify inclusion of the data.
type Sequencer interface {
	Start(context.Context) error
	Stop() error
	Write(tx []byte) error
}

// Executor is the gateway to an arbitrary state machine. Transactions
// are always processed strictly serially. They must be deterministic. The state
// is represented as a monotonically increasing number known as the cursor.
type Executor interface {
	// Execute executes a transaction. Invalid transactions are ignored and
	// the cursor is not updated. The state machine must be crash safe. No error
	// means the state transition was successful. The executor always returns
	// the current state.
	Execeute(tx []byte) (cursor uint64, err error)

	// Query provides read access to the state machine. The executor returns the
	// response and the cursor representing the state of the machine at the time
	// of the query. The query may be integrated with a verification scheme.
	// An executor is not responsible for historical queries.
	Query(query []byte) (response []byte, cursor uint64, err error)

	// Cursor returns the current state.
	Cursor() uint64
}

// Generic client interface for data availability layers. The client is resposible
// for running the underlying verification scheme that proves the data was published
// by the network.
type Publisher interface {
	Get(ctx context.Context, namespace []byte, cursor uint64, limit uint64) (data [][]byte, err error)
	Set(ctx context.Context, namespace []byte, data []byte) error
}
