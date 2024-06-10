package bloom

import (
	"hash/fnv"
	"math"
)

// BloomFilterBits represents a Bloom filter data structure.
type BloomFilterBits struct {
	bits      []uint64 // Bit array for storing filter data
	numHashes int      // Number of hash functions
}

// NewBloomFilterBits creates a new Bloom filter with the specified capacity and false positive rate.
func NewBloomFilterBits(capacity int, falsePositiveRate float64) *BloomFilterBits {
	// Calculate the optimal number of bits and hash functions
	numBits := int(math.Ceil(-float64(capacity) * math.Log(falsePositiveRate) / math.Pow(math.Log(2), 2)))
	numHashes := int(math.Ceil(float64(numBits) / float64(capacity) * math.Log(2)))

	// Create the bit array
	bits := make([]uint64, (numBits+63)/64) // Allocate enough 64-bit words

	return &BloomFilterBits{
		bits:      bits,
		numHashes: numHashes,
	}
}

// Add adds a byte slice to the Bloom filter.
func (bf *BloomFilterBits) Add(data []byte) bool {

	// Check if the element is already in the filter
	if bf.Exists(data) {
		return true // Indicate that the element was already added
	}

	// Calculate hash values
	hashes := bf.calculateHashes(data)

	// Set bits in the bit array
	for _, hash := range hashes {
		index := hash % uint64(len(bf.bits)-1)
		if index >= uint64(len(bf.bits)) {
			panic(index)
		}

		bf.bits[index] |= 1 << (hash % 64) // Set the corresponding bit
	}

	return false // Indicate that the element was added
}

// Exists checks if a byte slice might be in the Bloom filter.
func (bf *BloomFilterBits) Exists(data []byte) bool {
	// Calculate hash values
	hashes := bf.calculateHashes(data)

	// Check if all bits are set
	for _, hash := range hashes {
		index := hash % uint64(len(bf.bits)-1)
		if bf.bits[index]&(1<<(hash%64)) == 0 { // Check if the bit is not set
			return false // Definitely not in the filter
		}
	}

	return true // Possibly in the filter (could be a false positive)
}

// calculateHashes calculates the hash values for a given byte slice.
func (bf *BloomFilterBits) calculateHashes(data []byte) []uint64 {
	hashes := make([]uint64, bf.numHashes)
	for i := 0; i < bf.numHashes; i++ {
		h := fnv.New32a() // Use FNV-1a hash function
		h.Write(data)
		hashes[i] = uint64(h.Sum32())
	}
	return hashes
}

// AddString adds a string to the Bloom filter.
func (bf *BloomFilterBits) AddString(s string) bool {
	return bf.Add([]byte(s))
}

// ExistsString checks if a string might be in the Bloom filter.
func (bf *BloomFilterBits) ExistsString(s string) bool {
	return bf.Exists([]byte(s))
}
