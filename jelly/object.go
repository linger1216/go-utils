package jelly

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	jsoniter.RegisterTypeDecoder("jelly.Object", &ObjectCodec{})
	jsoniter.RegisterTypeEncoder("jelly.Object", &ObjectCodec{})

	jsoniter.RegisterTypeDecoder("jelly.Value", &ValueCodec{})
	jsoniter.RegisterTypeEncoder("jelly.Value", &ValueCodec{})
}

func (o Object) GetValue(key string) *Value {
	if o.Values != nil {
		return o.Values[key]
	}
	return nil
}

func (o Object) GetInt32Array(key string) []int32 {
	if o.Values != nil {
		return o.Values[key].GetInt32Array()
	}
	return []int32{}
}

func (o Object) GetInt64Array(key string) []int64 {
	if o.Values != nil {
		return o.Values[key].GetInt64Array()
	}
	return []int64{}
}

func (o Object) GetUint32Array(key string) []uint32 {
	if o.Values != nil {
		return o.Values[key].GetUint32Array()
	}
	return []uint32{}
}

func (o Object) GetUint64Array(key string) []uint64 {
	if o.Values != nil {
		return o.Values[key].GetUint64Array()
	}
	return []uint64{}
}

func (o Object) GetStringArray(key string) []string {
	if o.Values != nil {
		return o.Values[key].GetStringArray()
	}
	return []string{}
}

func (o Object) GetObjectArray(key string) []*Object {
	if o.Values != nil {
		return o.Values[key].GetObjectArray()
	}
	return []*Object{}
}

func (o *Object) SetValue(key string, v *Value) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = v
	return o
}

func (o *Object) SetBool(key string, v bool) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewBoolValue(v)
	return o
}

func (o *Object) SetInt(key string, v int) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewIntValue(v)
	return o
}

func (o *Object) SetInt32(key string, v int32) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewInt32Value(v)
	return o
}

func (o *Object) SetInt64(key string, v int64) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewInt64Value(v)
	return o
}

func (o *Object) SetUint32(key string, v uint32) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewUint32Value(v)
	return o
}

func (o *Object) SetUint64(key string, v uint64) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewUint64Value(v)
	return o
}

func (o *Object) SetFloat32(key string, v float32) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewFloat32Value(v)
	return o
}

func (o *Object) SetFloat64(key string, v float64) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewFloat64Value(v)
	return o
}

func (o *Object) SetString(key string, v string) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewStringValue(v)
	return o
}

func (o *Object) SetObject(key string, v *Object) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewObjectValue(v)
	return o
}

func (o *Object) SetInt32Array(key string, vals ...int32) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewInt32ArrayValue(vals...)
	return o
}

func (o *Object) SetInt64Array(key string, vals ...int64) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewInt64ArrayValue(vals...)
	return o
}

func (o *Object) SetUint32Array(key string, vals ...uint32) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewUint32ArrayValue(vals...)
	return o
}

func (o *Object) SetUint64Array(key string, vals ...uint64) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewUint64ArrayValue(vals...)
	return o
}

func (o *Object) SetStringArray(key string, vals ...string) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewStringArrayValue(vals...)
	return o
}

func (o *Object) SetObjectArray(key string, vals ...*Object) *Object {
	if o.Values == nil {
		o.Values = make(map[string]*Value)
	}
	o.Values[key] = NewObjectArrayValue(vals...)
	return o
}

type ObjectCodec struct {
}

func (codec *ObjectCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	any := iter.ReadAny()

	if any.ValueType() == jsoniter.ObjectValue {
		obj := (*Object)(ptr)
		obj.Values = make(map[string]*Value)
		for _, key := range any.Keys() {
			v := any.Get(key)
			switch v.ValueType() {
			case jsoniter.BoolValue:
				obj.Values[key] = NewBoolValue(v.ToBool())
			case jsoniter.NumberValue:
				int64V := any.ToInt64()
				floatV := any.ToFloat64()
				if int64(floatV) == int64V {
					obj.Values[key] = NewInt64Value(int64V)
				} else {
					obj.Values[key] = NewFloat64Value(floatV)
				}
				obj.Values[key] = NewUint64Value(v.ToUint64())
			case jsoniter.StringValue:
				obj.Values[key] = NewStringValue(v.ToString())
			case jsoniter.ArrayValue:
				values := make([]*Value, 0, v.Size())
				v.ToVal(&values)
				obj.Values[key] = NewArrayValue(values...)
			case jsoniter.ObjectValue:
				object := &Object{}
				v.ToVal(object)
				obj.Values[key] = NewObjectValue(object)
			}
		}
	}
}

func (codec *ObjectCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return len(((*Object)(ptr)).Values) == 0
}

func (codec *ObjectCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	object := (*Object)(ptr)
	stream.WriteVal(object.Values)
}

func NewObjectValue(obj *Object) *Value {
	return &Value{Val: &Value_ObjectVal{ObjectVal: obj}}
}

func NewArrayValue(values ...*Value) *Value {
	return &Value{Val: &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}}
}

func NewBoolValue(val bool) *Value {
	return &Value{Val: &Value_BoolVal{BoolVal: val}}
}

func NewIntValue(val int) *Value {
	return NewInt64Value(int64(val))
}

func NewInt32Value(val int32) *Value {
	return NewInt64Value(int64(val))
}

func NewInt64Value(val int64) *Value {
	if val < 0 {
		return &Value{Val: &Value_NegIntVal{NegIntVal: uint64(-val)}}
	} else {
		return &Value{Val: &Value_PosIntVal{PosIntVal: uint64(val)}}
	}
}

func NewUintValue(val uint) *Value {
	return NewUint64Value(uint64(val))
}

func NewUint32Value(val uint32) *Value {
	return NewUint64Value(uint64(val))
}

func NewUint64Value(val uint64) *Value {
	return &Value{Val: &Value_PosIntVal{PosIntVal: val}}
}

func NewFloat32Value(val float32) *Value {
	return NewFloat64Value(float64(val))
}

func NewFloat64Value(val float64) *Value {
	return &Value{Val: &Value_DoubleVal{DoubleVal: val}}
}

func NewStringValue(val string) *Value {
	return &Value{Val: &Value_StringVal{StringVal: val}}
}

func NewIntArrayValue(vals ...int) *Value {
	values := make([]*Value, 0, len(vals))
	for _, v := range vals {
		values = append(values, NewIntValue(v))
	}
	return &Value{Val: &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}}
}

func NewInt32ArrayValue(vals ...int32) *Value {
	values := make([]*Value, 0, len(vals))
	for _, v := range vals {
		values = append(values, NewInt32Value(v))
	}
	return &Value{Val: &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}}
}

func NewInt64ArrayValue(vals ...int64) *Value {
	values := make([]*Value, 0, len(vals))
	for _, v := range vals {
		values = append(values, NewInt64Value(v))
	}
	return &Value{Val: &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}}
}

func NewUint32ArrayValue(vals ...uint32) *Value {
	values := make([]*Value, 0, len(vals))
	for _, v := range vals {
		values = append(values, NewUint32Value(v))
	}
	return &Value{Val: &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}}
}

func NewUint64ArrayValue(vals ...uint64) *Value {
	values := make([]*Value, 0, len(vals))
	for _, v := range vals {
		values = append(values, NewUint64Value(v))
	}
	return &Value{Val: &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}}
}

func NewFloat32ArrayValue(vals ...float32) *Value {
	values := make([]*Value, 0, len(vals))
	for _, v := range vals {
		values = append(values, NewFloat32Value(v))
	}
	return &Value{Val: &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}}
}

func NewFloat64ArrayValue(vals ...float64) *Value {
	values := make([]*Value, 0, len(vals))
	for _, v := range vals {
		values = append(values, NewFloat64Value(v))
	}
	return &Value{Val: &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}}
}

func NewStringArrayValue(vals ...string) *Value {
	values := make([]*Value, 0, len(vals))
	for _, v := range vals {
		values = append(values, NewStringValue(v))
	}
	return &Value{Val: &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}}
}

func NewObjectArrayValue(vals ...*Object) *Value {
	values := make([]*Value, 0, len(vals))
	for _, v := range vals {
		values = append(values, NewObjectValue(v))
	}
	return &Value{Val: &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}}
}

func (v Value) GetBool() bool {
	return v.GetBoolVal()
}

func (v Value) GetInt() int {
	return int(v.GetInt64())
}

func (v Value) GetInt32() int32 {
	return int32(v.GetInt64())
}

func (v Value) GetInt64() int64 {
	pos := v.GetPosIntVal()
	if pos == 0 {
		return -int64(v.GetNegIntVal())
	} else {
		return int64(pos)
	}
}

func (v Value) GetUint() uint {
	return uint(v.GetUint64())
}

func (v Value) GetUint32() uint32 {
	return uint32(v.GetPosIntVal())
}

func (v Value) GetUint64() uint64 {
	pos := v.GetPosIntVal()
	if pos == 0 {
		return uint64(-int64(v.GetNegIntVal()))
	} else {
		return pos
	}
}

func (v Value) GetFloat32() float32 {
	return float32(v.GetFloat64())
}

func (v Value) GetFloat64() float64 {
	return v.GetFloat64()
}

func (v Value) GetString() string {
	return v.GetStringVal()
}

func (v Value) GetObject() *Object {
	return v.GetObjectVal()
}

func (v Value) GetArray() []*Value {
	return v.GetArrayVal().Values
}

func (v Value) GetObjectArray() []*Object {
	values := v.GetArrayVal().Values
	vals := make([]*Object, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.GetObject())
	}
	return vals
}

func (v Value) GetStringArray() []string {
	values := v.GetArrayVal().Values
	vals := make([]string, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.GetString())
	}
	return vals
}

func (v Value) GetIntArray() []int {
	values := v.GetArrayVal().Values
	vals := make([]int, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.GetInt())
	}
	return vals
}

func (v Value) GetInt32Array() []int32 {
	values := v.GetArrayVal().Values
	vals := make([]int32, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.GetInt32())
	}
	return vals
}

func (v Value) GetInt64Array() []int64 {
	values := v.GetArrayVal().Values
	vals := make([]int64, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.GetInt64())
	}
	return vals
}

func (v Value) GetUintArray() []uint {
	values := v.GetArrayVal().Values
	vals := make([]uint, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.GetUint())
	}
	return vals
}

func (v Value) GetUint32Array() []uint32 {
	values := v.GetArrayVal().Values
	vals := make([]uint32, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.GetUint32())
	}
	return vals
}

func (v Value) GetUint64Array() []uint64 {
	values := v.GetArrayVal().Values
	vals := make([]uint64, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.GetUint64())
	}
	return vals
}

func (v Value) GetFloat32Array() []float32 {
	values := v.GetArrayVal().Values
	vals := make([]float32, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.GetFloat32())
	}
	return vals
}

func (v Value) GetFloat64Array() []float64 {
	values := v.GetArrayVal().Values
	vals := make([]float64, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.GetFloat64())
	}
	return vals
}

type ValueCodec struct {
}

func (codec *ValueCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	any := iter.ReadAny()
	val := (*Value)(ptr)
	switch any.ValueType() {
	case jsoniter.BoolValue:
		val.Val = &Value_BoolVal{BoolVal: any.ToBool()}
	case jsoniter.NumberValue:
		int64V := any.ToInt64() // if integer overflow will change pos to neg int
		floatV := any.ToFloat64()
		if floatV == float64(int64V) { // [-2^53, 2^53]
			if int64V < 0 {
				val.Val = &Value_NegIntVal{NegIntVal: uint64(-int64V)}
			} else {
				val.Val = &Value_PosIntVal{PosIntVal: uint64(int64V)}
			}
		} else {
			val.Val = &Value_DoubleVal{DoubleVal: floatV}
		}
	case jsoniter.StringValue:
		val.Val = &Value_StringVal{StringVal: any.ToString()}
	case jsoniter.ArrayValue:
		values := make([]*Value, 0, any.Size())
		any.ToVal(&values)
		val.Val = &Value_ArrayVal{ArrayVal: &ArrayValue{Values: values}}
	case jsoniter.ObjectValue:
		obj := &Object{}
		any.ToVal(obj)
		val.Val = &Value_ObjectVal{ObjectVal: obj}
	}
}

func (codec *ValueCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return ((*Value)(ptr)) == nil
}

func (codec *ValueCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	value := (*Value)(ptr)
	switch val := value.Val.(type) {
	case *Value_BoolVal:
		stream.WriteBool(val.BoolVal)
	case *Value_NegIntVal:
		stream.WriteInt64(int64(-val.NegIntVal))
	case *Value_PosIntVal:
		stream.WriteUint64(val.PosIntVal)
	case *Value_DoubleVal:
		stream.WriteFloat64Lossy(val.DoubleVal)
	case *Value_StringVal:
		stream.WriteString(val.StringVal)
	case *Value_ArrayVal:
		stream.WriteVal(val.ArrayVal.Values)
	case *Value_ObjectVal:
		stream.WriteVal(val.ObjectVal.Values)
	}
}
