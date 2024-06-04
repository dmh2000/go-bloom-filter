// Handles BL32 server not running or unreachable
package bloom

import (
	"testing"
)

func newBL32(*testing.T) *BloomFilter32 {
	return NewBloomFilter32()
}

func newBL32B(*testing.B) *BloomFilter32 {
	return NewBloomFilter32()
}

// Verify Add returns false when adding a new byte slice
func TestBL32AddNewByteSlice(t *testing.T) {
	bf := newBL32(t)

	var collisions = 0
	input := []byte("test data 2")
	result := bf.Add(input)
	if result == MAYBE_IN_FILTER {
		collisions++
	}
	t.Log("collisions:", collisions)
}

// Verify Exists returns true when adding an existing byte slice
func TestBL32ExistsByteSlice(t *testing.T) {
	bf := newBL32(t)

	input := []byte("test data 2")
	bf.Add(input)
	result := bf.Exists(input)
	if result == NOT_IN_FILTER {
		t.Errorf("Expected Exists to return true for existing, got %v", result)
	}
}

// Verify AddString returns false when adding a new string
func TestBL32AddString(t *testing.T) {
	bf := newBL32(t)

	var collisions = 0
	input := "test data 3"
	result := bf.AddString(input)
	if result == MAYBE_IN_FILTER {
		collisions++
	}
	t.Log("collisions:", collisions)
}

// Verify ExistsString returns true when adding an existing string
func TestBL32ExistsString(t *testing.T) {
	bf := newBL32(t)

	input := "test data 3"
	bf.AddString(input)
	result := bf.ExistsString(input)
	if result == NOT_IN_FILTER {
		t.Errorf("Expected ExistsString to return true for existing, got %v", result)
	}
}

func TestBL32GenerateInstanceID(t *testing.T) {
	bf := newBL32(t)
	const iterations = 1000000
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

func BenchmarkBL32GenerateInstanceID0001(b *testing.B) {
	bf := newBL32B(b)

	const iterations = 1000000
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
