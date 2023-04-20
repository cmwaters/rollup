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
	Execute(tx []byte) (cursor uint64, err error)

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
// by the network as well as managing the keys/signatures of the user.
type Publisher interface {
	Reader
	Writer
}

type Reader interface {
	// Get fetches and verifies a range of data from a namespace using a cursor and limit. If the cursor is 0
	// the latest data from the namespace is returned. If limit is 0, the server decides how many blobs of data
	// to return. The caller should continually track the cursor to know which data is missing.
	Get(ctx context.Context, namespace []byte, cursor uint64, limit uint64) (data [][]byte, err error)

	// Has checks if a blob of data exists in the namespace. The method returns the cursor position of the data
	// This is similar to Get but avoids needing to download the actual data.
	Has(ctx context.Context, namespace []byte, data []byte) (cursor uint64, err error)
}

type Writer interface {
	// Set publishes data to a namespace. The method is fire and forget. Use `Has` to verify the data was published.
	Set(ctx context.Context, namespace []byte, data []byte) error
}
