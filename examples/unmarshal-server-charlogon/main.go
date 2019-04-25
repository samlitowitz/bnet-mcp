package main

import (
	"fmt"
	"log"

	"github.com/samlitowitz/bnet-encoding/pkg/encoding/bnet"
	"github.com/samlitowitz/bnet-mcp/pkg/mcp"
	"github.com/samlitowitz/bnet-mcp/pkg/mcp/server"
)

func main() {
	data := []byte{0x09, 0x00, 0x07, 0x01, 0x00, 0x00, 0x00}

	var header mcp.Header
	err := bnet.Unmarshal(data[:3], &header)
	if err != nil {
		log.Fatal(err)
	}

	var charlogon server.CharLogon
	err = bnet.Unmarshal(data[3:], &charlogon)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n%+v\n", header, charlogon)
}
