package msgpack

import (
	"bytes"
	"encoding/binary"
	"github.com/profzone/msgpack/common"
	"io"
	"reflect"
)

type MessagePack interface {
	DecodeMessage(data []byte, message interface{}) error
	EncodeMessage(message interface{}) ([]byte, error)
}

type MessageController struct {
	packer MessagePack
}

func (mc *MessageController) ReadMessage(reader io.Reader, message interface{}) error {
	if mc.packer == nil {
		panic("message packer not initialized")
	}

	msgType := reflect.TypeOf(message)
	if msgType.Kind() != reflect.Ptr {
		panic("message must be a pointer")
	}

	// read type
	typeBuf := make([]byte, 1)
	_, err := reader.Read(typeBuf)
	if err != nil {
		return err
	}

	// read length
	var length uint64
	err = binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// read data
	dataBuf := make([]byte, length)
	count, err := io.ReadFull(reader, dataBuf)
	if err != nil {
		return err
	}

	if uint64(count) != length {
		return common.ErrMsgLength
	}

	return mc.packer.DecodeMessage(dataBuf, message)
}

func (mc *MessageController) WriteMessage(writer io.Writer, message interface{}) error {
	if mc.packer == nil {
		panic("message packer not initialized")
	}

	buf := bytes.NewBuffer([]byte{})

	// write type
	_, err := buf.Write([]byte{1})
	if err != nil {
		return err
	}

	dataBuf, err := mc.packer.EncodeMessage(message)
	if err != nil {
		return err
	}

	// write length
	length := uint64(len(dataBuf))
	err = binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return err
	}

	// write data
	_, err = buf.Write(dataBuf)
	if err != nil {
		return err
	}

	_, err = writer.Write(buf.Bytes())
	return err
}
