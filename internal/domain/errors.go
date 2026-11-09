package domain

import "errors"

var (
	ErrBinNotFound     = errors.New("bin not found")
	ErrRequestNotFound = errors.New("request not found")
)
