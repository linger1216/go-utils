package jelly

import (
	fmt "fmt"
	"github.com/gogo/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Object struct {
	Values               map[string]*Value `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Object) Reset()         { *m = Object{} }
func (m *Object) String() string { return proto.CompactTextString(m) }
func (*Object) ProtoMessage()    {}
func (*Object) Descriptor() ([]byte, []int) {
	return fileDescriptor_afb3bfaba1b12cb1, []int{0}
}
func (m *Object) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Object.Unmarshal(m, b)
}
func (m *Object) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Object.Marshal(b, m, deterministic)
}
func (m *Object) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object.Merge(m, src)
}
func (m *Object) XXX_Size() int {
	return xxx_messageInfo_Object.Size(m)
}
func (m *Object) XXX_DiscardUnknown() {
	xxx_messageInfo_Object.DiscardUnknown(m)
}

var xxx_messageInfo_Object proto.InternalMessageInfo

func (m *Object) GetValues() map[string]*Value {
	if m != nil {
		return m.Values
	}
	return nil
}

type ArrayValue struct {
	Values               []*Value `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ArrayValue) Reset()         { *m = ArrayValue{} }
func (m *ArrayValue) String() string { return proto.CompactTextString(m) }
func (*ArrayValue) ProtoMessage()    {}
func (*ArrayValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_afb3bfaba1b12cb1, []int{1}
}
func (m *ArrayValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ArrayValue.Unmarshal(m, b)
}
func (m *ArrayValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ArrayValue.Marshal(b, m, deterministic)
}
func (m *ArrayValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ArrayValue.Merge(m, src)
}
func (m *ArrayValue) XXX_Size() int {
	return xxx_messageInfo_ArrayValue.Size(m)
}
func (m *ArrayValue) XXX_DiscardUnknown() {
	xxx_messageInfo_ArrayValue.DiscardUnknown(m)
}

var xxx_messageInfo_ArrayValue proto.InternalMessageInfo

func (m *ArrayValue) GetValues() []*Value {
	if m != nil {
		return m.Values
	}
	return nil
}

type Value struct {
	// Types that are valid to be assigned to Val:
	//	*Value_StringVal
	//	*Value_DoubleVal
	//	*Value_PosIntVal
	//	*Value_NegIntVal
	//	*Value_BoolVal
	//	*Value_ArrayVal
	//	*Value_ObjectVal
	Val                  isValue_Val `protobuf_oneof:"val"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_afb3bfaba1b12cb1, []int{2}
}
func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (m *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(m, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

type isValue_Val interface {
	isValue_Val()
}

type Value_StringVal struct {
	StringVal string `protobuf:"bytes,1,opt,name=string_val,json=stringVal,proto3,oneof"`
}
type Value_DoubleVal struct {
	DoubleVal float64 `protobuf:"fixed64,2,opt,name=double_val,json=doubleVal,proto3,oneof"`
}
type Value_PosIntVal struct {
	PosIntVal uint64 `protobuf:"varint,3,opt,name=pos_int_val,json=posIntVal,proto3,oneof"`
}
type Value_NegIntVal struct {
	NegIntVal uint64 `protobuf:"varint,4,opt,name=neg_int_val,json=negIntVal,proto3,oneof"`
}
type Value_BoolVal struct {
	BoolVal bool `protobuf:"varint,5,opt,name=bool_val,json=boolVal,proto3,oneof"`
}
type Value_ArrayVal struct {
	ArrayVal *ArrayValue `protobuf:"bytes,13,opt,name=array_val,json=arrayVal,proto3,oneof"`
}
type Value_ObjectVal struct {
	ObjectVal *Object `protobuf:"bytes,14,opt,name=object_val,json=objectVal,proto3,oneof"`
}

func (*Value_StringVal) isValue_Val() {}
func (*Value_DoubleVal) isValue_Val() {}
func (*Value_PosIntVal) isValue_Val() {}
func (*Value_NegIntVal) isValue_Val() {}
func (*Value_BoolVal) isValue_Val()   {}
func (*Value_ArrayVal) isValue_Val()  {}
func (*Value_ObjectVal) isValue_Val() {}

func (m *Value) GetVal() isValue_Val {
	if m != nil {
		return m.Val
	}
	return nil
}

func (m *Value) GetStringVal() string {
	if x, ok := m.GetVal().(*Value_StringVal); ok {
		return x.StringVal
	}
	return ""
}

func (m *Value) GetDoubleVal() float64 {
	if x, ok := m.GetVal().(*Value_DoubleVal); ok {
		return x.DoubleVal
	}
	return 0
}

func (m *Value) GetPosIntVal() uint64 {
	if x, ok := m.GetVal().(*Value_PosIntVal); ok {
		return x.PosIntVal
	}
	return 0
}

func (m *Value) GetNegIntVal() uint64 {
	if x, ok := m.GetVal().(*Value_NegIntVal); ok {
		return x.NegIntVal
	}
	return 0
}

func (m *Value) GetBoolVal() bool {
	if x, ok := m.GetVal().(*Value_BoolVal); ok {
		return x.BoolVal
	}
	return false
}

func (m *Value) GetArrayVal() *ArrayValue {
	if x, ok := m.GetVal().(*Value_ArrayVal); ok {
		return x.ArrayVal
	}
	return nil
}

func (m *Value) GetObjectVal() *Object {
	if x, ok := m.GetVal().(*Value_ObjectVal); ok {
		return x.ObjectVal
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Value) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Value_OneofMarshaler, _Value_OneofUnmarshaler, _Value_OneofSizer, []interface{}{
		(*Value_StringVal)(nil),
		(*Value_DoubleVal)(nil),
		(*Value_PosIntVal)(nil),
		(*Value_NegIntVal)(nil),
		(*Value_BoolVal)(nil),
		(*Value_ArrayVal)(nil),
		(*Value_ObjectVal)(nil),
	}
}

func _Value_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Value)
	// val
	switch x := m.Val.(type) {
	case *Value_StringVal:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.StringVal)
	case *Value_DoubleVal:
		_ = b.EncodeVarint(2<<3 | proto.WireFixed64)
		_ = b.EncodeFixed64(math.Float64bits(x.DoubleVal))
	case *Value_PosIntVal:
		_ = b.EncodeVarint(3<<3 | proto.WireVarint)
		_ = b.EncodeVarint(uint64(x.PosIntVal))
	case *Value_NegIntVal:
		_ = b.EncodeVarint(4<<3 | proto.WireVarint)
		_ = b.EncodeVarint(uint64(x.NegIntVal))
	case *Value_BoolVal:
		t := uint64(0)
		if x.BoolVal {
			t = 1
		}
		_ = b.EncodeVarint(5<<3 | proto.WireVarint)
		_ = b.EncodeVarint(t)
	case *Value_ArrayVal:
		_ = b.EncodeVarint(13<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ArrayVal); err != nil {
			return err
		}
	case *Value_ObjectVal:
		_ = b.EncodeVarint(14<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ObjectVal); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Value.Val has unexpected type %T", x)
	}
	return nil
}

func _Value_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Value)
	switch tag {
	case 1: // val.string_val
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Val = &Value_StringVal{x}
		return true, err
	case 2: // val.double_val
		if wire != proto.WireFixed64 {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeFixed64()
		m.Val = &Value_DoubleVal{math.Float64frombits(x)}
		return true, err
	case 3: // val.pos_int_val
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Val = &Value_PosIntVal{x}
		return true, err
	case 4: // val.neg_int_val
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Val = &Value_NegIntVal{x}
		return true, err
	case 5: // val.bool_val
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Val = &Value_BoolVal{x != 0}
		return true, err
	case 13: // val.array_val
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ArrayValue)
		err := b.DecodeMessage(msg)
		m.Val = &Value_ArrayVal{msg}
		return true, err
	case 14: // val.object_val
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Object)
		err := b.DecodeMessage(msg)
		m.Val = &Value_ObjectVal{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Value_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Value)
	// val
	switch x := m.Val.(type) {
	case *Value_StringVal:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.StringVal)))
		n += len(x.StringVal)
	case *Value_DoubleVal:
		n += 1 // tag and wire
		n += 8
	case *Value_PosIntVal:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(x.PosIntVal))
	case *Value_NegIntVal:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(x.NegIntVal))
	case *Value_BoolVal:
		n += 1 // tag and wire
		n += 1
	case *Value_ArrayVal:
		s := proto.Size(x.ArrayVal)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Value_ObjectVal:
		s := proto.Size(x.ObjectVal)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

var fileDescriptor_afb3bfaba1b12cb1 = []byte{
	// 386 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0x5d, 0xeb, 0xd3, 0x30,
	0x18, 0xc5, 0x97, 0x75, 0x9d, 0xeb, 0x53, 0x94, 0x59, 0x50, 0xc6, 0x44, 0x2c, 0xbb, 0x90, 0xde,
	0x98, 0xc2, 0x7c, 0x41, 0xdc, 0x95, 0x83, 0xc1, 0xc4, 0x0b, 0xa1, 0x17, 0x43, 0xbc, 0x19, 0x69,
	0x17, 0x4a, 0x67, 0x96, 0xa7, 0xa4, 0xd9, 0xa4, 0x7e, 0x06, 0xbf, 0x81, 0x5f, 0x56, 0x92, 0xd4,
	0x6e, 0xfc, 0x77, 0x53, 0x9a, 0x73, 0x7e, 0x27, 0x9c, 0x24, 0x0f, 0x3c, 0x3f, 0xe1, 0x11, 0xd3,
	0x02, 0x15, 0x4f, 0x31, 0x3f, 0xf2, 0x42, 0xd3, 0x5a, 0xa1, 0xc6, 0x28, 0x30, 0x3a, 0x35, 0xfa,
	0x7c, 0x7e, 0x45, 0xb8, 0x52, 0xa8, 0xf6, 0x05, 0x1e, 0xb8, 0xc3, 0x16, 0x7f, 0x08, 0x8c, 0xbf,
	0xd9, 0x5c, 0xf4, 0x1e, 0xc6, 0x17, 0x26, 0xce, 0xbc, 0x99, 0x91, 0xd8, 0x4b, 0xc2, 0xe5, 0x4b,
	0xda, 0x6f, 0x41, 0x1d, 0x42, 0x77, 0xd6, 0xdf, 0x48, 0xad, 0xda, 0xac, 0x83, 0xe7, 0x5f, 0x21,
	0xbc, 0x91, 0xa3, 0x29, 0x78, 0x3f, 0x79, 0x3b, 0x23, 0x31, 0x49, 0x82, 0xcc, 0xfc, 0x46, 0xaf,
	0xc1, 0xb7, 0xe8, 0x6c, 0x18, 0x93, 0x24, 0x5c, 0x4e, 0x6f, 0xb6, 0xb5, 0xc1, 0xcc, 0xd9, 0x9f,
	0x86, 0x1f, 0xc9, 0xe2, 0x03, 0xc0, 0x67, 0xa5, 0x58, 0x6b, 0x8d, 0x28, 0x79, 0xd0, 0xe8, 0x3e,
	0xda, 0xf9, 0x8b, 0xbf, 0x43, 0xf0, 0x5d, 0xe6, 0x15, 0x40, 0xa3, 0x55, 0x25, 0xcb, 0xfd, 0x85,
	0x09, 0x57, 0x63, 0x3b, 0xc8, 0x02, 0xa7, 0xed, 0x98, 0x30, 0xc0, 0x01, 0xcf, 0xb9, 0xe0, 0x16,
	0x30, 0x9d, 0x88, 0x01, 0x9c, 0x66, 0x80, 0x18, 0xc2, 0x1a, 0x9b, 0x7d, 0x25, 0xb5, 0x25, 0xbc,
	0x98, 0x24, 0x23, 0x43, 0xd4, 0xd8, 0x7c, 0x91, 0xba, 0x23, 0x24, 0x2f, 0x7b, 0x62, 0xf4, 0x9f,
	0x90, 0xbc, 0xec, 0x88, 0x17, 0x30, 0xc9, 0x11, 0x85, 0xb5, 0xfd, 0x98, 0x24, 0x93, 0xed, 0x20,
	0x7b, 0x64, 0x14, 0x63, 0xbe, 0x83, 0x80, 0x99, 0x43, 0x5a, 0xf7, 0xb1, 0xbd, 0x94, 0x67, 0x37,
	0x27, 0xbb, 0x5e, 0xc0, 0x76, 0x90, 0x4d, 0x58, 0xb7, 0x8a, 0x96, 0x00, 0xee, 0x81, 0x6d, 0xec,
	0x89, 0x8d, 0x3d, 0xbd, 0x7b, 0x22, 0x53, 0xc3, 0x61, 0x3b, 0x26, 0xd6, 0x3e, 0x78, 0x17, 0x26,
	0xd6, 0xdf, 0x21, 0x2a, 0xf0, 0x44, 0x7f, 0xb1, 0xf6, 0xf7, 0x35, 0xb0, 0xf6, 0x37, 0x66, 0x18,
	0x7e, 0xac, 0xca, 0x4a, 0x53, 0x56, 0x9d, 0x58, 0x4d, 0x2b, 0x4c, 0x05, 0x16, 0x4c, 0x57, 0x28,
	0xd3, 0x92, 0x4b, 0x3b, 0x21, 0xbd, 0xf2, 0x86, 0xd5, 0x55, 0x93, 0xf6, 0xa3, 0xb4, 0x32, 0x9f,
	0x7c, 0x6c, 0x99, 0xb7, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x38, 0xf5, 0xb0, 0xaf, 0x86, 0x02,
	0x00, 0x00,
}
