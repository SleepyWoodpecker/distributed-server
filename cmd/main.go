package main

import (
	"distfileserver/package/p2p"
	"fmt"
)

const PORT = ":3000"

func main() {
	tr := p2p.NewTCPTransport(PORT)
	tr.ListenAndAccept()

	fmt.Printf("Starting our TCP server at %s\n", PORT)

	// introduce a blocking loop
	select {}
}