package server

// AuthCheckResponse is the structure of a SID_AUTH_CHECK response
type AuthCheck struct {
	Result uint32
	Info   string
}
