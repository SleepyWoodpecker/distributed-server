package p2p

import (
	"fmt"
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

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	
	t.listener, err = net.Listen("tcp", t.listenAddr)

	if err != nil {
		fmt.Printf("TCP Listen error: %v\n", err)
		return err
	}

	go t.listenLoop()
	return nil
}

func (t *TCPTransport) listenLoop() {
	for {
		conn, err := t.listener.Accept()

		if err != nil {
			fmt.Printf("TCP Accept error: %v\n", err)
		}

		// actually, is there a need for there to be more than one acceptor?
		go t.acceptConn(conn)
	}
}

func (t *TCPTransport) acceptConn(conn net.Conn) {
	// prints structs in a human readable way
	fmt.Printf("Incoming connection from %+v\n", conn)
}