package server

// GameInfoResponse is structure of a MCP_GAMEINFO response
type GameInfo struct {
	RequestID                  uint16
	Status                     uint32
	Uptime                     uint32 // Seconds
	LevelRestrictionLevel      uint8
	LevelRestrictionDifference uint8
	MaxPlayers                 uint8
	CharacterCount             uint8
	CharacterClasses           [16]uint8
	CharacterLevels            [16]uint8
	Description                string
	CharacterNames             string
}
