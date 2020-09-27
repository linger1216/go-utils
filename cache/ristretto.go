package cache

import (
	"github.com/dgraph-io/ristretto"
	"time"
)

type Ristretto struct {
	proxy *ristretto.Cache
}

func (x *Ristretto) Set(k string, v interface{}) error {
	x.proxy.Set(k, v, 1)
	return nil
}

func (x *Ristretto) SetWithExpire(k string, v interface{}, t time.Duration) error {
	x.proxy.SetWithTTL(k, v, 1, t)
	return nil
}

func (x *Ristretto) Get(k string) (interface{}, error) {
	v, _ := x.proxy.Get(k)
	return v, nil
}

func (x *Ristretto) Del(k string) error {
	x.proxy.Del(k)
	return nil
}

func NewRistretto(maxCost int64) *Ristretto {
	ret := &Ristretto{}
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     maxCost, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	ret.proxy = cache
	return ret
}
