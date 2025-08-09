package p2p

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn net.Conn

	// a connection is considered to be outbound if it is the one making the connection request
	// a connection is considered to be inbound if it is the one receiving the connection request
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr string
	HandshakeFunc
	Decoder
}

type TCPTransport struct {
	TCPTransportOpts

	listener net.Listener

	// mutexes are usually placed above the thing they are meant to guard
	mu sync.RWMutex
	// net address represents a network connection and a destination address
	peerMap map[net.Addr]Peer
}

// chose to return a TCPTransport rather than a transport here
// becuase it would make accessing its struct members easier (no need to do a type assertion)
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)

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
	peer := NewTCPPeer(conn, false)
	defer peer.conn.Close()

	// prints structs in a human readable way
	fmt.Printf("Incoming connection from %+v\n", peer)

	// test a handshake first, close the connection if it does not succeed
	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("TCP Handshake error: %+v\n", err)
		return
	}

	// receive incoming messages
	buffer := new(bytes.Buffer)
	for {
		if err := t.Decode(conn, buffer); err != nil {
			fmt.Printf("TCP Decode erorr: %+v\n", err)
			return
		}

		fmt.Printf("Msg: %s", buffer)
	}
}
