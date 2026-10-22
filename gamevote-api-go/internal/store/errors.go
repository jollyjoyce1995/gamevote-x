package store

import "fmt"

// Kind classifies an error so the HTTP layer can pick a status code.
type Kind int

const (
	KindInvalid     Kind = iota // bad input -> 400
	KindNotFound                // 404
	KindConflict                // 409
	KindUnavailable             // upstream dependency failed -> 502
)

// Error is an expected, user-facing failure. Anything else is a 500.
type Error struct {
	Kind Kind
	Msg  string
}

func (e *Error) Error() string { return e.Msg }

func Invalid(format string, a ...any) *Error {
	return &Error{KindInvalid, fmt.Sprintf(format, a...)}
}

func NotFound(format string, a ...any) *Error {
	return &Error{KindNotFound, fmt.Sprintf(format, a...)}
}

func Conflict(format string, a ...any) *Error {
	return &Error{KindConflict, fmt.Sprintf(format, a...)}
}

func Unavailable(format string, a ...any) *Error {
	return &Error{KindUnavailable, fmt.Sprintf(format, a...)}
}
