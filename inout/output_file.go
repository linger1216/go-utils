package inout

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"sync"
)

type WriterFileConfig struct {
	// os.O_RDWR|os.O_CREATE|os.O_TRUNC
	// os.O_RDWR|os.O_CREATE|O_APPEND
	Flag      string
	Filename  string
	WriteSize int
}

func NewWriterFileConfig(flag string, filename string, writeSize int) *WriterFileConfig {
	return &WriterFileConfig{Flag: flag, Filename: filename, WriteSize: writeSize}
}

type WriterFile struct {
	sync.RWMutex
	f *os.File
	w *bufio.Writer
}

func (x *WriterFile) Close() error {
	x.Lock()
	defer x.Unlock()
	_ = x.w.Flush()
	return x.f.Close()
}

func NewWriterFile(cfg *WriterFileConfig) *WriterFile {
	ret := &WriterFile{}
	flag := 0
	if strings.ToLower(cfg.Flag) == "append" {
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	} else if strings.ToLower(cfg.Flag) == "trunc" {
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	}

	writeSize := cfg.WriteSize
	if writeSize == 0 {
		writeSize = 4096
	}

	obj, err := os.OpenFile(cfg.Filename, flag, 0644)
	if err != nil {
		panic(err)
	}
	ret.f = obj
	ret.w = bufio.NewWriterSize(ret.f, writeSize)
	return ret
}

func (x *WriterFile) Exec(args ...[]interface{}) (interface{}, error) {
	x.Lock()
	defer x.Unlock()
	for _, v := range args {
		var content bytes.Buffer
		for _, unit := range v {
			switch x := unit.(type) {
			case []byte:
				content.Write(x)
			case string:
				content.WriteString(x)
			case byte:
				content.WriteByte(x)
			}
		}
		_, err := x.w.Write(content.Bytes())
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
