package main

import (
	"distfileserver/pkg/fileserver"
	"distfileserver/pkg/p2p"
	"distfileserver/pkg/store"
	"log"
)

func main() {
	server1 := makeServer(":3000", []string{})
	server2 := makeServer(":4000", []string{":3000"})

	go func() { log.Fatal(server1.Start()) }()
	go func() {
		<-server1.ServerDoneChan
		log.Fatal(server2.Start())
	}()

	select {}
}

func makeServer(port string, bootstrapNodes []string) *fileserver.FileServer {
	transportOpts := p2p.TCPTransportOpts{
		ListenAddr:    port,
		HandshakeFunc: p2p.NOPHandshake,
		Decoder:       p2p.DefaultDecoder{},
	}

	storeOpts := store.StoreOpts{
		PathTransformFunc: store.CASPathTransformFunc,
	}

	server := fileserver.NewFileServer(
		fileserver.FileServerOpts{
			BootstrapNodes: bootstrapNodes,
		},
		storeOpts,
		transportOpts,
	)

	server.Transport.OnPeer = server.OnPeer
	return server
}
