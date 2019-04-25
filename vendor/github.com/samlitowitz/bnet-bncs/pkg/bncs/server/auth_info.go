package server

// AuthInfoResponse is the structure of a SID_AUTH_INFO response
type AuthInfo struct {
	LogonType   uint32
	ServerToken uint32
	UDPValue    uint32
	MPQFiletime uint64
	MPQFilename string
	Value       string
}
