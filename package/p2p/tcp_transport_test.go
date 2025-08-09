package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4123"
	tcpTransport := NewTCPTransport(listenAddr)

	assert.Equal(t, tcpTransport.listenAddr, listenAddr)

	// Test the listening
	assert.Nil(t, tcpTransport.ListenAndAccept())
}