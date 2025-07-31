// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domain

// Error implements constant errors
type Error string

// Error implements the Errors interface
func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotImplemented = Error("not implemented")
)
