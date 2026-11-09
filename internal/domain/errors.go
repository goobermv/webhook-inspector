package domain

import "errors"

var (
	ErrBinNotFound     = errors.New("bin not found")
	ErrRequestNotFound = errors.New("request not found")

	ErrFailedToDeleteBin = errors.New("failed to delete bin")
	ErrBinFull           = errors.New("bin is full")
	ErrInvalidSignature  = errors.New("invalid hmac signature")
	ErrBodyTooLarge      = errors.New("request body too large")
	ErrInvalidTargetURL  = errors.New("invalid target url")
	ErrForbiddenTarget   = errors.New("target url is not allowed")
	ErrReplayFailed      = errors.New("replay failed")
	ErrBinExpired        = errors.New("bin expired")
)
