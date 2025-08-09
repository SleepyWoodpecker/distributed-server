package p2p

import (
	"fmt"
	"io"
	"net"
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

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr string
	HandshakeFunc
	Decoder
}

type TCPTransport struct {
	TCPTransportOpts

	listener   net.Listener
	msgChannel chan Message
}

// chose to return a TCPTransport rather than a transport here
// becuase it would make accessing its struct members easier (no need to do a type assertion)
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		msgChannel:       make(chan Message),
	}
}

func (t *TCPTransport) Consume() <-chan Message {
	return t.msgChannel
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
	defer peer.Close()

	// prints structs in a human readable way
	fmt.Printf("Incoming connection from %+v\n", peer)

	// test a handshake first, close the connection if it does not succeed
	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("TCP Handshake error: %+v\n", err)
		return
	}

	// receive incoming messages
	msg := Message{}
	for {
		// all errors to be handled at the top most component here
		if err := t.Decode(conn, &msg); err != nil {
			if err == io.EOF {
				fmt.Println("Closing connection")
				return
			}

			fmt.Printf("TCP Decode erorr: %+v\n", err)
			return
		}

		// record the sender address so you can send back a message later
		msg.From = conn.RemoteAddr()
		t.msgChannel <- msg
	}
}
