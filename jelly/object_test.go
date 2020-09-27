package jelly

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

const ObjectValue1 = `{"key":"value","key2":12}`
const ObjectValue2 = `{"ok":{"okk1":"v1","okk2":123}, "key2":true}`
const ObjectValue3 = `{"key":["val1","val2","val3"]}`

func TestObjectCodec_Decode(t *testing.T) {
	obj := &Object{}

	err := jsoniter.ConfigDefault.Unmarshal([]byte(ObjectValue1), obj)
	assert.Nil(t, err)
	assert.Equal(t, "value", obj.GetValue("key").GetString())

	err = jsoniter.ConfigDefault.Unmarshal([]byte(ObjectValue2), obj)
	assert.Nil(t, err)

	err = jsoniter.ConfigDefault.Unmarshal([]byte(ObjectValue3), obj)
	assert.Nil(t, err)
	assert.Equal(t, "val1", obj.GetStringArray("key")[0])
}

func TestObjectCodec_Encode(t *testing.T) {
	obj := &Object{}
	obj.SetString("key", "value")
	obj.SetInt("key2", 12)

	str, err := jsoniter.ConfigDefault.MarshalToString(obj)
	assert.Nil(t, err)
	assert.Contains(t, str, `"key":"value"`, `"key2":12`)
}

func TestObjectCodec_Encode2(t *testing.T) {
	obj := &Object{}
	obj.SetString("okk1", "v1")
	obj.SetInt("okk2", 123)

	o := &Object{}
	o.SetObject("ok", obj)
	o.SetBool("key2", true)

	str, err := jsoniter.ConfigDefault.MarshalToString(o)
	assert.Nil(t, err)
	assert.Contains(t, str, `"okk1":"v1"`, `"okk2":123`, `"key2":true`)
}
