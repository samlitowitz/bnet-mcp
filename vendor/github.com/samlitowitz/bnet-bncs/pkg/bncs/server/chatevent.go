package server

// ChatEventResponse is the structure of a SID_CHATEVENT response
type ChatEvent struct {
	EventID               uint32
	UserFlags             uint32
	Ping                  uint32
	IP                    uint32 // Defunct
	AccountNumber         uint32 // Defunct
	RegistrationAuthority uint32 // Defunct
	Username              string
	Text                  string // Max length 254. Max length for official clients is 223
}
