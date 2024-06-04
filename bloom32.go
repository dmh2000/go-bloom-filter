package bloom

import (
	"hash/crc32"
	"sync"
)

/*
Conveniently, hash/crc32 provides three different polynomials that
can be used to generate a hash of a byte slice. This implementation
uses all three polynomials to generate three different hashes for
each byte slice that is added to the filter.

This implementation uses a map of empty structs to make lookups cheap. The
map uses empty structs because they use less memory than bools.

This implementation uses a Mutex to protect the sequence of 3 hash
computations and lookups. This is necessary because the map is not
safe for concurrent access. A plain Mutex is preferred over RWMutex
because most operations will be Add


If a 32 bit hash is not enough for your use case, you can add an
implementation that uses other hash functions. Typically a bloom filter
wants at least K = 3 hash functions, but you can use more if you want.
*/

// use map of empty structs to store the hashes because it uses less memory than bool's
type BloomFilter32 struct {
	set map[uint32]struct{}
	mtx sync.Mutex
}

func NewBloomFilter32() *BloomFilter32 {
	return &BloomFilter32{
		set: make(map[uint32]struct{}),
		mtx: sync.Mutex{},
	}
}

func NewBloomFilter32Size(size int) *BloomFilter32 {
	return &BloomFilter32{
		set: make(map[uint32]struct{}, size),
		mtx: sync.Mutex{},
	}
}
func (bf *BloomFilter32) test(a, b, c uint32) bool {
	// must be called with lock aready held
	_, a_ok := bf.set[a]
	_, b_ok := bf.set[b]
	_, c_ok := bf.set[c]

	return a_ok && b_ok && c_ok
}

func (bf *BloomFilter32) Add(id []byte) bool {
	// generate three different hashes for the byte slice
	a := crc32ieee(id)
	b := crc32castagnoli(id)
	c := crc32koopman(id)

	bf.mtx.Lock()
	defer bf.mtx.Unlock()

	// if the byte slice is already in the filter, return true
	if bf.test(a, b, c) {
		return true
	}

	// add the byte slice to the filter
	bf.set[a] = struct{}{}
	bf.set[b] = struct{}{}
	bf.set[c] = struct{}{}

	// return false to indicate that the byte slice was not already in the filter
	return false
}

func (bf *BloomFilter32) Exists(id []byte) bool {
	a := crc32ieee(id)
	b := crc32castagnoli(id)
	c := crc32koopman(id)

	bf.mtx.Lock()
	defer bf.mtx.Unlock()

	return bf.test(a, b, c)
}

func (bf *BloomFilter32) AddString(id string) bool {
	return bf.Add([]byte(id))
}

func (bf *BloomFilter32) ExistsString(id string) bool {
	s := []byte(id)
	return bf.Exists(s)
}

func (bf *BloomFilter32) Len() int {
	bf.mtx.Lock()
	defer bf.mtx.Unlock()

	return len(bf.set)
}

// ----------------
// HASH FUNCTIONS
// ----------------

// precompute the tables for the three polynomials
var ieee *crc32.Table = crc32.MakeTable(crc32.IEEE)
var castagnoli *crc32.Table = crc32.MakeTable(crc32.Castagnoli)
var koopman *crc32.Table = crc32.MakeTable(crc32.Koopman)

// generate the IEEE CRC32 hash of a byte slice
func crc32ieee(b []byte) uint32 {
	return crc32.Checksum(b, ieee)
}

// generate the Castagnoli CRC32 hash of a byte slice
func crc32castagnoli(b []byte) uint32 {
	return crc32.Checksum(b, castagnoli)
}

// generate the Koopman CRC32 hash of a byte slice
func crc32koopman(b []byte) uint32 {
	return crc32.Checksum(b, koopman)
}
