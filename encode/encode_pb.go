package encode

import (
	"errors"
	// "github.com/golang/protobuf/proto"

	// "github.com/frankee/protobuf/proto"

	"github.com/gogo/protobuf/proto"
)

type ProtoBufEncoder struct {
}

func (x *ProtoBufEncoder) Marshal(message interface{}) ([]byte, error) {
	if v, ok := message.(proto.Message); ok {
		buf, err := proto.Marshal(v)
		if err != nil {
			return nil, err
		}
		return buf, nil
	}
	return nil, errors.New("not proto message")
}

func (x *ProtoBufEncoder) Unmarshal(buf []byte, obj interface{}) error {
	if v, ok := obj.(proto.Message); ok {
		err := proto.Unmarshal(buf, v)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("not proto message")
}
