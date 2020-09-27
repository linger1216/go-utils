package encode

import (
	"fmt"
	"strings"
)

type Encode interface {
	Marshal(message interface{}) ([]byte, error)
	Unmarshal(buf []byte, obj interface{}) error
}

func NewEncode(name string) Encode {
	var enc Encode
	encoder := strings.ToLower(name)
	switch encoder {
	case "json":
		enc = &JsonEncoder{}
	case "pb", "protobuf", "proto":
		enc = &ProtoBufEncoder{}
	case "pb64":
		enc = &ProtoBufBase64Encoder{}
	default:
		panic(fmt.Sprintf("unsupported encode:%s", name))
	}
	return enc
}
