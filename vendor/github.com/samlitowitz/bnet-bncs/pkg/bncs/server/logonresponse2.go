package server

// LogonResponse2Response is the structure of a SID_LOGONRESPONSE2 response
type LogonResponse2 struct {
	Status uint32
	//	Info   string // not always present
}
