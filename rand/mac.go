package rand

import (
	"fmt"
	"math/rand"
	"time"
)

type Mac [6]byte

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func (m Mac) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", m[0], m[1], m[2], m[3], m[4], m[5])
}

func NewRandomMac() Mac {
	var m [6]byte
	for i := 0; i < 6; i++ {
		macByte := r.Intn(256)
		m[i] = byte(macByte)
	}
	return Mac(m)
}
