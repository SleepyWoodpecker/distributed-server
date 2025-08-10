package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, msg *Message) error {
	return gob.NewDecoder(r).Decode(&msg.Payload)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *Message) error {
	// make a buffer to read the incoming message
	// if the current message overflows the buffer, the reader will just stop reading and decode
	// this in the next iteration of the loop
	buf := make([]byte, 1024)

	n, err := r.Read(buf)

	if err != nil {
		return err
	}

	// store the message in the message struct
	msg.Payload = buf[:n]
	return nil
}
