package p2p

// A peer represents a remote node
type Peer interface {
}

// A transport is anything that handles the
// communication between 2 nodes
type Transport interface {
	ListenAndAccept() error
}
