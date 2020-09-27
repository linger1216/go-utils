package encode

import (
	"encoding/base64"
	"errors"
	"github.com/gogo/protobuf/proto"
)

type ProtoBufBase64Encoder struct {
}

func (x *ProtoBufBase64Encoder) Marshal(message interface{}) ([]byte, error) {
	if v, ok := message.(proto.Message); ok {
		buf, err := proto.Marshal(v)
		if err != nil {
			return nil, err
		}
		str := base64.StdEncoding.EncodeToString(buf)
		return []byte(str), nil
	}
	return nil, errors.New("not proto message")
}

func (x *ProtoBufBase64Encoder) Unmarshal(buf []byte, obj interface{}) error {
	if v, ok := obj.(proto.Message); ok {

		protoBuf, err := base64.StdEncoding.DecodeString(string(buf))
		if err != nil {
			return err
		}

		err = proto.Unmarshal(protoBuf, v)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("not proto message")
}
