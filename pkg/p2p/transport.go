package p2p

import "net"

// A peer represents a remote node
type Peer interface {
	RemoteAddr() net.Addr
	Close() error
}

// A transport is anything that handles the
// communication between 2 nodes
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan Message
	Close()
}
