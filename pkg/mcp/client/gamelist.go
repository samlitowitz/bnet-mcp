package client

// GameListRequest is structure of a MCP_GAMELIST request
type GameList struct {
	RequestID uint16
	Unknown   uint32
	Search    string
}
