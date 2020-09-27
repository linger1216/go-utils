package config

type Reader interface {
	ScanKey(key string, v interface{}) error
	GetString(path ...string) string
	GetInt64(path ...string) int64
}
