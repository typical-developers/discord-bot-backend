package handlers

import "errors"

var (
	ErrGatewayTimeout     = errors.New("gateway timeout")
	ErrInternalError      = errors.New("internal error")
	ErrInvalidRequestBody = errors.New("malformed request body")
)
