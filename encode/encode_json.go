package encode

import (
	"github.com/json-iterator/go"
)

type JsonEncoder struct {
}

func (x *JsonEncoder) Marshal(message interface{}) ([]byte, error) {
	return jsoniter.ConfigFastest.Marshal(message)
}

func (x *JsonEncoder) Unmarshal(buf []byte, obj interface{}) error {
	return jsoniter.ConfigFastest.Unmarshal(buf, obj)
}
