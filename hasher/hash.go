package hasher

import "strings"

type Hasher interface {
	Hash(str string) uint64
}

type Bkdr struct {
}

func (b *Bkdr) Hash(str string) uint64 {
	seed := uint64(131) // 31 131 1313 13131 131313 etc..
	hash := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = (hash * seed) + uint64(str[i])
	}
	return hash
}

const (
	// offset64 FNVa offset basis. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	offset64 = 14695981039346656037
	// prime64 FNVa prime value. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	prime64 = 1099511628211
)

type Fnv struct {
}

func (b *Fnv) Hash(key string) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}
	return hash
}

func NewHasher(name string) Hasher {
	var enc Hasher
	encoder := strings.ToLower(name)
	if encoder == "fnv" {
		enc = &Fnv{}
	} else {
		enc = &Bkdr{}
	}
	return enc
}
