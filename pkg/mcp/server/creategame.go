package server

// CreateGameResponse is structure of a MCP_CREATEGAME response
type CreateGame struct {
	RequestID uint16
	GameToken uint16
	Unknown   uint16
	Result    uint32
}
