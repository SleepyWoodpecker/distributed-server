package main

import (
	"distfileserver/pkg/fileserver"
	"distfileserver/pkg/p2p"
	"distfileserver/pkg/store"
	"fmt"
	"os"
)

const PORT = ":3000"

func main() {
	transportOpts := p2p.TCPTransportOpts{
		ListenAddr:    PORT,
		HandshakeFunc: p2p.NOPHandshake,
		Decoder:       p2p.DefaultDecoder{},
	}

	storeOpts := store.StoreOpts{
		PathTransformFunc: store.CASPathTransformFunc,
	}

	server := fileserver.NewFileServer(
		fileserver.FileServerOpts{},
		storeOpts,
		transportOpts,
	)

	if err := server.Start(); err != nil {
		fmt.Printf("Error starting server: %v", err)
		os.Exit(1)
	}
}
