// Handles  server not running or unreachable
package bloom

import (
	"testing"
)

// Verify Add returns false when adding a new byte slice
func AddNewByteSlice(t *testing.T, bf BloomFilter) {

	var collisions = 0
	input := []byte("test data 2")
	result := bf.Add(input)
	if result == MAYBE_IN_FILTER {
		t.Errorf("Expected Add to return false for added new, got %v", result)
	}
	t.Log("collisions:", collisions)

	// check that is now found in the bitmask
	result = bf.Exists(input)
	if result == NOT_IN_FILTER {
		t.Errorf("Expected Exists to return true for existing, got %v", result)
	}
}

// Verify Exists returns true when adding an existing byte slice
func ExistsByteSlice(t *testing.T, bf BloomFilter) {

	input := []byte("test data 2")
	result := bf.Exists(input)
	if result == NOT_IN_FILTER {
		t.Errorf("Expected Exists to return true for existing, got %v", result)
	}
}

// Verify AddString returns false when adding a new string
func AddString(t *testing.T, bf BloomFilter) {
	var collisions = 0
	input := "test data 3"
	result := bf.AddString(input)
	if result == MAYBE_IN_FILTER {
		collisions++
	}
	t.Log("collisions:", collisions)
}

// Verify ExistsString returns true when adding an existing string
func ExistsString(t *testing.T, bf BloomFilter) {
	input := "test data 3"
	bf.AddString(input)
	result := bf.ExistsString(input)
	if result == NOT_IN_FILTER {
		t.Errorf("Expected ExistsString to return true for existing, got %v", result)
	}
}

// Test Add with an empty byte slice
func AddEmptyByteSlice(t *testing.T, bf BloomFilter) {
	var collisions = 0
	input := []byte{}
	result := bf.Add(input)
	if result == MAYBE_IN_FILTER {
		collisions++
	}
	t.Log("collisions:", collisions)
}

func GenerateInstanceID(t *testing.T, bf BloomFilter) {
	const iterations = 10000
	var collisions = 0
	ids := make([]string, 0, iterations)
	for i := 0; i < iterations; i++ {
		id := generateInstanceID(16)
		ids = append(ids, id)
		result := bf.AddString(id)
		if result == MAYBE_IN_FILTER {
			collisions++
		}
	}
	for _, id := range ids {
		result := bf.ExistsString(id)
		if result == NOT_IN_FILTER {
			t.Errorf("Expected ExistsString to return true for existing, got %v", result)
		}
	}
	t.Log("collisions:", collisions)

}

func GenerateInstanceID0001(b *testing.B, bf BloomFilter) {
	const iterations = 10000
	var collisions = 0
	for i := 0; i < iterations; i++ {
		id := generateInstanceID(32)
		result := bf.AddString(id)
		if result == MAYBE_IN_FILTER {
			collisions++
		}
	}
	b.Log("collisions:", collisions)
}

func TestBloomFilter32(t *testing.T) {

	bf := NewBloomFilter32(10*1024*1024, .001)
	AddNewByteSlice(t, bf)
	ExistsByteSlice(t, bf)
	AddString(t, bf)
	ExistsString(t, bf)
	AddEmptyByteSlice(t, bf)
	GenerateInstanceID(t, bf)
}

func TestBloomFilterBits(t *testing.T) {

	bf := NewBloomFilterBits(10*1024*1024, .001)
	AddNewByteSlice(t, bf)
	ExistsByteSlice(t, bf)
	AddString(t, bf)
	ExistsString(t, bf)
	AddEmptyByteSlice(t, bf)
	GenerateInstanceID(t, bf)
}

func TestBloomFilterRedis(t *testing.T) {

	key := generateInstanceID(8)
	t.Log(key)
	bf, err := NewBloomFilterRedis(key, 10*1024*1024, .001)
	if err != nil {
		t.Error(err)
	}
	AddNewByteSlice(t, bf)
	ExistsByteSlice(t, bf)
	AddString(t, bf)
	ExistsString(t, bf)
	AddEmptyByteSlice(t, bf)
	GenerateInstanceID(t, bf)
}

func TestBloomFilter(t *testing.T) {
	key := generateInstanceID(8)
	b32 := NewBloomFilter32(10*1024*1024, .001)
	bits := NewBloomFilterBits(10*1024*1024, .001)
	redis, err := NewBloomFilterRedis(key, 10*1024*1024, .001)
	if err != nil {
		t.Error(err)
	}
	testCases := []struct {
		name string
		bf   BloomFilter
	}{
		{
			"b32",
			b32,
		},
		{
			"bits",
			bits,
		},
		{
			"redis",
			redis,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			AddNewByteSlice(t, tc.bf)
			ExistsByteSlice(t, tc.bf)
			AddString(t, tc.bf)
			ExistsString(t, tc.bf)
			AddEmptyByteSlice(t, tc.bf)
			GenerateInstanceID(t, tc.bf)
		})
	}
}

func RunGenerateInstanceID0001(b *testing.B, bf BloomFilter) {
	const iterations = 10000
	var collisions = 0
	for i := 0; i < iterations; i++ {
		id := generateInstanceID(32)
		result := bf.AddString(id)
		if result == MAYBE_IN_FILTER {
			collisions++
		}
	}
	b.Log("collisions:", collisions)
}

func BenchmarkBloomFilter(b *testing.B) {
	key := generateInstanceID(8)
	b32 := NewBloomFilter32(10*1024*1024, .001)
	bits := NewBloomFilterBits(10*1024*1024, .001)
	redis, err := NewBloomFilterRedis(key, 10*1024*1024, .001)
	if err != nil {
		b.Error(err)
	}
	testCases := []struct {
		name string
		bf   BloomFilter
	}{
		{
			"b32",
			b32,
		},
		{
			"bits",
			bits,
		},
		{
			"redis",
			redis,
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			RunGenerateInstanceID0001(b, tc.bf)
		})
	}
}
