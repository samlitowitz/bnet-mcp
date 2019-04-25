package server

// CharList2ResponseCharacter is the character structure of a MCP_CHARLIST2 response
type CharList2Character struct {
	ExpirationDate uint32 // Unix time
	Name           string
	Statstring     string
}

// CharList2Response is the structure of a MCP_CHARLIST2 response
type CharList2 struct {
	RequestCount  uint16
	ExistCount    uint32
	ReturnedCount uint16 `bnet:"save-CL2Returned"`

	Characters []CharList2Character `bnet:"len-CL2Returned"`
}
