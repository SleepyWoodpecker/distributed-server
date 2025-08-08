package p2p

import (
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddr string
	listener net.Listener

	// mutexes are usually placed above the thing they are meant to guard
	mu sync.RWMutex
	// net address represents a network connection and a destination address
	peerMap map[net.Addr]Peer 
}

// chose to return a TCPTransport rather than a transport here 
// becuase it would make accessing its struct members easier (no need to do a type assertion)
func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: listenAddr,
	}
}