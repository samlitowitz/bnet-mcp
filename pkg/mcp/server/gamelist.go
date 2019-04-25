package server

// GameListResponse is structure of a MCP_GAMELIST response
type GameList struct {
	RequestID   uint16
	Index       uint32
	PlayerCount uint8
	Status      uint32
	Name        string
	Description string
}
