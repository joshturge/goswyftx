package goswyftx

import "errors"

var (
	errAssetCode = errors.New("asset code was not set")
)

type Error struct {
	Summary string `json:"error"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Summary + ": " + e.Message
}
