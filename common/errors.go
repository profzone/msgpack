package common

import "errors"

var (
	ErrMsgType   = errors.New("cannot read message type")
	ErrMsgLength = errors.New("message not complete")
)
