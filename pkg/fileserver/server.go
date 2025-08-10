package fileserver

import (
	"distfileserver/pkg/p2p"
	"distfileserver/pkg/store"
	"fmt"
	"sync"
)

type FileServerOpts struct {
	BootstrapNodes []string
}

type FileServer struct {
	FileServerOpts FileServerOpts
	Store          *store.Store
	Transport      *p2p.TCPTransport
	quitch         chan struct{}

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	// track whether the server's set up is complete
	ServerDoneChan chan struct{}
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
		peers:  make(map[string]p2p.Peer),
		ServerDoneChan: make(chan struct{}),
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()
	s.ServerDoneChan<-struct{}{}
	
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
		case msg := <-s.Transport.Consume():
			fmt.Printf("Incoming message: %+v\n", msg)
		case <-s.quitch:
			return
		}
	}
}

// connect the current file server to other file servers on the network
func (s *FileServer) bootstrapNetwork() error {
	if len(s.FileServerOpts.BootstrapNodes) == 0 {
		return nil
	}

	for _, node := range s.FileServerOpts.BootstrapNodes {
		if err := s.Transport.Dial(node); err != nil {
			fmt.Printf("file server bootstrap error: %v\n", err)
		}
	}

	return nil
}

// function to run once the peer is created
func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p
	return nil
}
