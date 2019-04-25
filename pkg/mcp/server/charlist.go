package server

// CharListResponseCharacter is the character structure of a MCP_CHARLIST response
type CharListCharacter struct {
	Name       string
	Statstring string
}

// CharListResponse is the structure of a MCP_CHARLIST response
type CharList struct {
	RequestCount  uint16
	ExistCount    uint32
	ReturnedCount uint16 `bnet:"save-CLReturned"`

	Characters []CharListCharacter `bnet:"len-CLReturned"`
}
