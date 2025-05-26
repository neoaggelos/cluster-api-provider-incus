package utils

import "errors"

type terminalError struct {
	error
}

// TerminalError is used to indicate a non-retriable error.
func TerminalError(err error) error {
	return terminalError{err}
}

// IsTerminalError checks whether the error is a terminalError.
// These are returned to indicate non-retriable errors.
func IsTerminalError(err error) bool {
	return errors.As(err, &terminalError{})
}
