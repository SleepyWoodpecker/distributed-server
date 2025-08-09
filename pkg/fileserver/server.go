package fileserver

import (
	"distfileserver/pkg/p2p"
	"distfileserver/pkg/store"
)

type FileServerOpts struct {
}

type FileServer struct {
	FileServerOpts FileServerOpts
	Store          *store.Store
	Transport      *p2p.TCPTransport
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
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	return nil
}
