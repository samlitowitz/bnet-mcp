package server

// EnterChatResponse is the structure of a SID_ENTERCHAT response
type EnterChat struct {
	UniqueName  string
	Statstring  string
	AccountName string
}
