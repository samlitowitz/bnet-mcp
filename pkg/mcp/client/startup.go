package client

// StartupRequest is structure of a MCP_STARTUP request
type Startup struct {
	MCPCookie  uint32
	MCPStatus  uint32
	Chunk1     [2]uint32
	Chunk2     [12]uint32
	UniqueName string
}
