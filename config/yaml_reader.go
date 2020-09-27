package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"strings"
)

type YamlReader struct {
}

func (y *YamlReader) ScanKey(key string, v interface{}) error {
	return config.BindStruct(key, v)
}

func (y *YamlReader) GetString(path ...string) string {
	return config.String(strings.Join(path, "."), "")
}

func (y *YamlReader) GetInt64(path ...string) int64 {
	return config.Int64(strings.Join(path, "."), 0)
}

func NewYamlReader(filename string) *YamlReader {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles(filename)
	if err != nil {
		panic(err)
	}
	return &YamlReader{}
}
