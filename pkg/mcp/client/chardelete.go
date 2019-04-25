package client

// CharDeleteRequest is structure of a MCP_CHARDELETE request
type CharDelete struct {
	Unknown       uint16 // Cookie?
	CharacterName string
}
