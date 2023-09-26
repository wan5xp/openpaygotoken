package openpaygotoken

import "fmt"

// ErrInvalidTokenBase is returned when the token base is invalid.
type ErrInvalidTokenBase struct {
	Value int
}

func (e *ErrInvalidTokenBase) Error() string {
	return fmt.Sprintf("Invalid token base %d", e.Value)
}

// ErrInvalidTokenValue is returned when the token value is invalid.
type ErrValidOlderToken struct {
}

func (e *ErrValidOlderToken) Error() string {
	return "Valid older token"
}

// ErrInvalidToken is returned when the token is invalid.
type ErrInvalidToken struct {
}

func (e *ErrInvalidToken) Error() string {
	return "Invalid token"
}
