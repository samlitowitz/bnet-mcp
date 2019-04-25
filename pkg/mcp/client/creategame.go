package client

// CreateGameRequest is structure of a MCP_CREATEGAME request
type CreateGame struct {
	RequestID        uint16
	Difficulty       uint32
	Unknown          uint8
	LevelRestriction uint8
	MaxPlayers       uint8
	Name             string
	Password         string
	Description      string
}
