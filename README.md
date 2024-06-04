# A Simple Bloom Filter in Go

A Bloom filter lets you check if an element is 'maybe' in a set, or 'positivly' not in the set.

Bloom filters were created by Burton Bloom in 1970. Back then, memory was scarce and expensive, creating an incentive to find ways to conserve it. A Bloom filter does that for 'set' data structures vs a complete hash map.

[Read all about the math in the Wikipedia article](https://en.wikipedia.org/wiki/Bloom_filter)

[Or an even better general description in the Redis documentation](https://redis.io/docs/latest/develop/data-types/probabilistic/bloom-filter/)

Bottom line, a Bloom filter might help avoid expensive operations while being memory efficient. If I needed a robust and tunable Bloom filter, I would use Redis if it was available in the app architecture.

## Simple Go Version

Without doing the math , I see a lot of examples using K=3, where K is the number of hash functions used. It turns out in the Go standard library there is a convenient hash type (hash/crc32) that uses 32 bit CRC's with 3 different polynomials. CRC's have a complexity of O(n) and the computations are pretty simple. Anyway, it makes for a decent example implementation.

1. bloom.go : defines an interface for the Bloom filter

```go
type BloomFilter interface {
	// if the byte slice is already in the filter or a false positive, returns true immediately
	// if the byte slice is not in the filter, adds it and returns false
	Add([]byte) bool
	// if the string is already in the filter or a false positive, returns true immediately
	// if the string is not in the filter, adds it and returns false
	AddString(string) bool
	// checks if a byte slice is in the filter
	// returns true if the byte slice is in the filter or a false positive
	// return false if the byte slice is not in the filter
	Exists([]byte) bool
	// checks if a string is in the filter
	// returns true if the string is in the filter or a false positive
	// return false if the string is not in the filter
	ExistsString(string) bool
}
```

This interface is general enough that implementations could use any hash algorithms and values of K.

2. bloom32.go implements the interface methods that use the CRC32 hashes.

[The standard library for CRC32 has polynomials for IEEE, Castagnoli and Koopman](https://en.wikipedia.org/wiki/cyclic_redundancy_check#Standards_and_common_use)

The code is in [bloom32.go]().

3. bloom_test.go implements a set of tests using the standard Go test framework.

## Redis version (for comparison)

[Run Redis in the official docker container](https://redis.io/learn/operate/orchestration/docker)

- docker run -d --name redis-stack -p 6379:6379 redis/redis-stack:latest

the code is in [bloom_redis.go]() and the tests are in [bloom_redis_test.go]()

You will need to restart the docker instance for repeated tests so you don't get false positives for old data.

### Redis interface workaround

Redis 'add' functions return true if item added, false if its already there. That's the opposite of the logic used in this package. That's because I wanted Add and Exists to return the same value for 'not already there : false' and 'already there (maybe) : true'

The Redis approach takes a lot longer per operation than the local 'set' implementation. That's to be expected. The value of the Redis approach would be as a microservice in a system that had multiple separate clients. The Redis Bloom filter is also much more tunable than the local 'set'.
