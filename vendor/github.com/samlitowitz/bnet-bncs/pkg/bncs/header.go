package bncs

// Header is the structure of a BNCS header
type Header struct {
	Fixed     uint8
	MessageID MessageID
	Length    uint16
}
