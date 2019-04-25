package client

// CharCreateRequest is structure of a MCP_CHARCREATE request
type CharCreate struct {
	Class uint32
	Flags uint16
	Name  string
}
