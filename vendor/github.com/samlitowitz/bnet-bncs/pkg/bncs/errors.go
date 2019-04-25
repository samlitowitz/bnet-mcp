package bncs

import "fmt"

type InvalidProtocolError struct {
	Protocol byte
}

func (e *InvalidProtocolError) Error() string {
	return fmt.Sprintf("Invalid protocol %d", e.Protocol)
}
