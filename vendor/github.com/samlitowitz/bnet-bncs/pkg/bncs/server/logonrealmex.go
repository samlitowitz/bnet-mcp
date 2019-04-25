package server

// LogonRealmExResponse is the structure of a SID_LOGONREALMEX response
type LogonRealmEx struct {
	MCPCookie uint32
	MCPStatus uint32
	Chunk1    [2]uint32
	IP        [4]uint8
	Port      uint16 `bnet:"bigendian"`
	Unknown1  [2]uint8
	Chunk2    [12]uint32
	Name      string
}
