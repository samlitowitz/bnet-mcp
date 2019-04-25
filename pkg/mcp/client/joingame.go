package client

// JoinGameRequest is structure of a MCP_JOINGAME request
type JoinGame struct {
	RequestID uint16
	Name      string
	Password  string
}
