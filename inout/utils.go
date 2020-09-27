package inout

import (
	"io/ioutil"
	"os"
)

func ReadFileContent(filename string) ([]byte, error) {
	obj, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(obj)
	_ = obj.Close()
	return buf, err
}
