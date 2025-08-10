package fileserver

import (
	"distfileserver/pkg/p2p"
	"distfileserver/pkg/store"
	"fmt"
)

type FileServerOpts struct {
}

type FileServer struct {
	FileServerOpts FileServerOpts
	Store          *store.Store
	Transport      *p2p.TCPTransport
	quitch chan struct{}
}

func NewFileServer(
	serverOpts FileServerOpts,
	storeOpts store.StoreOpts,
	tcpTransportOpts p2p.TCPTransportOpts,
) *FileServer {
	return &FileServer{
		FileServerOpts: serverOpts,
		Store: store.NewStore(
			storeOpts,
		),
		Transport: p2p.NewTCPTransport(
			tcpTransportOpts,
		),
		quitch: make(chan struct{}),
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.loop()

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) loop() {
	defer func() {
		fmt.Println("Stopping server...")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <- s.Transport.Consume():
			fmt.Printf("Incoming message: %+v\n", msg)
		case <-s.quitch:
			return
		}
	}
}