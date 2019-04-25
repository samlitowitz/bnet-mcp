package server

// GetFiletimeResponse is the structure of a SID_GETFILETIME response
type GetFiletime struct {
	RequestID uint32
	Unknown   uint32
	Last      uint64 // Filetime
	Filename  string
}
