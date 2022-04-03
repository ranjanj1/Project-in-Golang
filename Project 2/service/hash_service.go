package service

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/aclements/go-gg/generic/slice"
)

type HashService struct{}

func ProvideHashService() HashService {
	return HashService{}
}

func (hs HashService) FindArgMin(name string, open, close uint64) uint64 {
	hashStore := make([]uint64, 0)

	for i := open; i < close; i++ {
		hashStore = append(hashStore, hs.hash(name, i))
	}

	idx := slice.ArgMin(hashStore)
	return hashStore[idx]
}

// hash concatenates a message and a nonce and generates a hash value.
func (hs HashService) hash(name string, nonce uint64) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}
