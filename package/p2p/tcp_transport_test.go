package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	var listenAddr string = ":3000"
	opts := TCPTransportOpts{ListenAddr: listenAddr, HandshakeFunc: NOPHandshake}

	tcpTransport := NewTCPTransport(opts)

	assert.Equal(t, tcpTransport.ListenAddr, listenAddr)

	// Test the listening
	assert.Nil(t, tcpTransport.ListenAndAccept())
}