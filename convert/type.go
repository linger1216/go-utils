package convert

import (
	"bytes"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

func UInt64ToString(n uint64) string {
	return strconv.FormatUint(uint64(n), 10)
}

func Decimal(value float64) float64 {
	return math.Round(value*1000000) / 1000000
}
func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

func StringToUint64(s string) uint64 {
	ret, _ := strconv.ParseUint(s, 10, 64)
	return ret
}

func StringToInt64(s string) int64 {
	ret, _ := strconv.ParseInt(s, 10, 64)
	return ret
}

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 6, 64)
}

func StringToFloat(s string) float64 {
	ret, _ := strconv.ParseFloat(s, 64)
	return ret
}

func ToInt64(v interface{}) int64 {
	if ret, ok := v.(int64); ok {
		return ret
	}
	return 0
}

func ToFloat64(v interface{}) float64 {
	if ret, ok := v.(float64); ok {
		return ret
	}
	return 0
}

func ToString(v interface{}) string {
	if ret, ok := v.(string); ok {
		return ret
	}

	if ret, ok := v.([]byte); ok {
		return string(ret)
	}
	return ""
}

// 驼峰式写法转为下划线写法
func Camel2Case(name string) string {
	var buffer bytes.Buffer
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteByte('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteByte(byte(r))
		}
	}
	return buffer.String()
}

// 下划线写法转为驼峰写法
// 大小大小
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// 首字母大写
func UpperFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// 首字母小写
func LowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
