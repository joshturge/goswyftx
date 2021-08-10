package swyftx

import "errors"

var (
	errAssetCode = errors.New("asset code was not set")
)

// Error is an swyftx API error response
type Error struct {
	Summary string `json:"error"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Summary + ": " + e.Message
}
