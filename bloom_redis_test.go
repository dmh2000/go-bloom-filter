// Handles Redis server not running or unreachable
package bloom

import (
	"context"
	"testing"
)

// singleton
var xbf *BloomRedis

// create a singleton instance of bloom filter with redist client
func newRedis(t *testing.T) *BloomRedis {

	if xbf != nil {
		return xbf
	}
	xf, err := NewBloomFilterRedis("bf_key", 0.001, 1000000)
	if err != nil {
		t.Log(err)
		return xbf
	}
	xbf = xf

	return xbf
}

// create a singleton instance of bloom filter with redist client
func newRedisB(t *testing.B) *BloomRedis {
	if xbf != nil {
		return xbf
	}
	xf, err := NewBloomFilterRedis("bf_key", 0.001, 1000000)
	if err != nil {
		t.Log(err)
		return xbf
	}
	xbf = xf

	return xbf
}

func TestRedisNewClientPing(t *testing.T) {
	bf := newRedis(t)

	ctx := context.Background()
	_, err := bf.client.Ping(ctx).Result()
	if err != nil {
		t.Error(err)
	}
}

// Verify Add returns false when adding a new byte slice
func TestRedisAddNewByteSlice(t *testing.T) {
	bf := newRedis(t)

	var collisions = 0
	input := []byte("test data 2")
	result := bf.Add(input)
	if result == MAYBE_IN_FILTER {
		collisions++
	}
	t.Log("collisions:", collisions)
}

// Verify Exists returns true when adding an existing byte slice
func TestRedisExistsByteSlice(t *testing.T) {
	bf := newRedis(t)

	input := []byte("test data 2")
	bf.Add(input)
	result := bf.Exists(input)
	if result == NOT_IN_FILTER {
		t.Errorf("Expected Exists to return true for existing, got %v", result)
	}
}

// Verify AddString returns false when adding a new string
func TestRedisAddString(t *testing.T) {
	bf := newRedis(t)

	var collisions = 0
	input := "test data 3"
	result := bf.AddString(input)
	if result == MAYBE_IN_FILTER {
		collisions++
	}
	t.Log("collisions:", collisions)
}

// Verify ExistsString returns true when adding an existing string
func TestRedisExistsString(t *testing.T) {
	bf := newRedis(t)

	input := "test data 3"
	bf.AddString(input)
	result := bf.ExistsString(input)
	if result == NOT_IN_FILTER {
		t.Errorf("Expected ExistsString to return true for existing, got %v", result)
	}
}

// Test Add with an empty byte slice
func TestRedisAddEmptyByteSlice(t *testing.T) {
	bf := newRedis(t)

	var collisions = 0
	input := []byte{}
	result := bf.Add(input)
	if result == MAYBE_IN_FILTER {
		collisions++
	}
	t.Log("collisions:", collisions)
}

// Test Exists with an empty byte slice
func TestRedisExistsEmptyByteSlice(t *testing.T) {
	bf := newRedis(t)

	input := []byte{}
	result := bf.Exists(input)
	if result == NOT_IN_FILTER {
		t.Errorf("Expected Exists to return true for existing, got %v", result)
	}
}

func TestRedisGenerateInstanceID(t *testing.T) {
	bf := newRedis(t)
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

func BenchmarkRedisGenerateInstanceID0001(b *testing.B) {
	bf := newRedisB(b)

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
