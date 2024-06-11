package bloom

import (
	"crypto/sha256"
	"hash"
	"math"
	"math/big"
	"math/rand"
)

type BloomFilterChatGPT struct {
	bitset  []bool
	k       int
	hashers []hash.Hash
	salt    []int64
}

func NewBloomFilterChatGPT(capacity int, falsePositiveRate float64) *BloomFilterChatGPT {
	m := int(math.Ceil(-float64(capacity) * math.Log(falsePositiveRate) / math.Pow(math.Log(2), 2)))
	k := int(math.Ceil(float64(m) / float64(capacity) * math.Log(2)))

	salt := make([]int64, k)
	for i := 0; i < k; i++ {
		salt[i] = rand.Int63()
	}

	bf := &BloomFilterChatGPT{
		bitset:  make([]bool, m),
		k:       k,
		hashers: make([]hash.Hash, k),
		salt:    salt,
	}

	for i := 0; i < k; i++ {
		bf.hashers[i] = sha256.New()
	}

	return bf
}

func (bf *BloomFilterChatGPT) getIndices(data []byte) []int {
	indices := make([]int, bf.k)
	for i := 0; i < bf.k; i++ {
		bf.hashers[i].Write(data)
		bf.hashers[i].Write([]byte{byte(bf.salt[i])})
		hashBytes := bf.hashers[i].Sum(nil)
		bf.hashers[i].Reset()
		indices[i] = int(new(big.Int).SetBytes(hashBytes).Uint64() % uint64(len(bf.bitset)))
	}
	return indices
}

func (bf *BloomFilterChatGPT) Add(data []byte) bool {
	indices := bf.getIndices(data)
	alreadyPresent := true
	for _, idx := range indices {
		if !bf.bitset[idx] {
			alreadyPresent = false
			bf.bitset[idx] = true
		}
	}
	return alreadyPresent
}

func (bf *BloomFilterChatGPT) AddString(s string) bool {
	return bf.Add([]byte(s))
}

func (bf *BloomFilterChatGPT) Exists(data []byte) bool {
	indices := bf.getIndices(data)
	for _, idx := range indices {
		if !bf.bitset[idx] {
			return false
		}
	}
	return true
}

func (bf *BloomFilterChatGPT) ExistsString(s string) bool {
	return bf.Exists([]byte(s))
}
