package server

// JoinGameResponse is structure of a MCP_JOINGAME response
type JoinGame struct {
	RequestID    uint16
	GameToken    uint16
	Unknown      uint16
	GameServerIP [4]uint8
	GameHash     uint32
	Result       uint32
}
