package json

import "encoding/json"

type JsonMessagePack struct{}

func NewJsonMessagePacker() *JsonMessagePack {
	return &JsonMessagePack{}
}

func (j *JsonMessagePack) DecodeMessage(data []byte, message interface{}) error {
	return json.Unmarshal(data, message)
}

func (j *JsonMessagePack) EncodeMessage(message interface{}) ([]byte, error) {
	return json.Marshal(message)
}
