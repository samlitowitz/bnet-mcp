package main

import (
	"fmt"
	"log"

	"github.com/samlitowitz/bnet-encoding/pkg/encoding/bnet"
	"github.com/samlitowitz/bnet-mcp/pkg/mcp"
	"github.com/samlitowitz/bnet-mcp/pkg/mcp/client"
)

func main() {
	buffer := make([]byte, 9)

	header := &mcp.Header{
		Length:    0x09,
		MessageID: mcp.McpCharLogon,
	}

	data, err := bnet.Marshal(&header)
	if err != nil {
		log.Fatal(err)
	}

	if copy(buffer, data) != 3 {
		log.Fatal("Expected 4 bytes copied.")
	}

	charlogon := &client.CharLogon{CharacterName: "Conan\n"}
	data, err = bnet.Marshal(&charlogon)
	if err != nil {
		log.Fatal(err)
	}

	if copy(buffer[3:], data) != 6 {
		log.Fatal("Expected 4 bytes copied.")
	}

	fmt.Printf("%+v\n", buffer)
}
