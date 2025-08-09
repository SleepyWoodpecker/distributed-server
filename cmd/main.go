package main

import (
	"distfileserver/package/p2p"
	"fmt"
)

const PORT = ":3000"

func main() {
	opts := p2p.TCPTransportOpts{
		ListenAddr:    PORT,
		HandshakeFunc: p2p.NOPHandshake,
		Decoder:       p2p.GOBDecoder{},
	}

	tr := p2p.NewTCPTransport(opts)
	tr.ListenAndAccept()

	fmt.Printf("Starting our TCP server at %s\n", PORT)

	// introduce a blocking loop
	select {}
}
