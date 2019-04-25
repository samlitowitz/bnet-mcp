package mcp

// Header is the structure of a MCP header
type Header struct {
	Length    uint16
	MessageID MessageID
}
