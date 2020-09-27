package jelly

import (
	"github.com/json-iterator/go"
	"time"
	"unsafe"

	"github.com/araddon/dateparse"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

var loc, _ = time.LoadLocation("Asia/Shanghai")
var normalFormatLen = len("2006-01-02 15:04:05")

func init() {
	jsoniter.RegisterTypeDecoder("jelly.Timestamp", &TimestampCodec{})
	jsoniter.RegisterTypeEncoder("jelly.Timestamp", &TimestampCodec{})
}

func (m Timestamp) ToTime() time.Time {
	return time.Unix(m.Seconds, int64(m.Nanos)).In(loc)
}

func (m *Timestamp) ParseTime(t time.Time) {
	m.Seconds = t.Unix()
	m.Nanos = int32(t.UnixNano() - m.Seconds*1000000000)
}

type TimestampCodec struct {
}

func (codec *TimestampCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	any := iter.ReadAny()
	if any.ValueType() == jsoniter.NumberValue {
		number := any.ToInt64()
		ts := (*Timestamp)(ptr)
		if number < int64(MaxInt) {
			ts.Seconds = number
		} else {
			ts.Seconds = number / 1000
			ts.Nanos = int32((number - ts.Seconds*1000) * 1000000)
		}
	} else if any.ValueType() == jsoniter.StringValue {
		str := any.ToString()

		// ts has timezone info, like "2006-01-02 15:04:05+0800"
		// since '+' will be replaced by space in url, we restore it to '+' if possible
		if len(str) > normalFormatLen && str[normalFormatLen] == ' ' {
			str = str[:normalFormatLen] + "+" + str[normalFormatLen+1:]
		}

		t, err := dateparse.ParseIn(str, loc)
		if err != nil {
		}

		ts := (*Timestamp)(ptr)
		ts.ParseTime(t)
	}
}

func (codec *TimestampCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return ((*Timestamp)(ptr)).Seconds == 0
}

func (codec *TimestampCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := (*Timestamp)(ptr)
	stream.WriteString(ts.ToTime().Format(time.RFC3339))
}
