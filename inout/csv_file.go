package inout

import (
	"encoding/csv"
	"os"
	"strings"
)

type CsvFileConfig struct {
	Flag      string
	Filename  string
	WriteSize int
}

func NewCsvFileConfig(flag string, filename string, writeSize int) *CsvFileConfig {
	return &CsvFileConfig{Flag: flag, Filename: filename, WriteSize: writeSize}
}

type CsvFile struct {
	f *os.File
	w *csv.Writer
}

func (x *CsvFile) Close() error {
	x.w.Flush()
	return x.f.Close()
}

func NewCsvFile(cfg *CsvFileConfig) *CsvFile {
	ret := &CsvFile{}
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
	ret.w = csv.NewWriter(obj)
	return ret
}

func (x *CsvFile) Exec(records ...string) error {
	return x.w.Write(records)
}
