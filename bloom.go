package bloom

/*
A Bloom filter is a probabilistic data structure that is used to test
whether an element is MAYBE a member of a set.

- a bloom filter helps when accessing a member of a set/database is expensive,
where checking with the bloom filter is not a member of the set is cheap.
- A true result means the element is probably a member of the set but might not be.
- A false test means the element is definitely not a member of the set.
- False positives are possible, but false negatives are not.
- Elements can be added to the bloom filter, but not removed, since there can be
collisions between elements.
- The more elements that are added to the bloom filter, the larger the probability of false positives.
- If the input data is short, there will be more false positives.
*/

const (
	NOT_IN_FILTER   = false
	MAYBE_IN_FILTER = true
)

type BloomFilter interface {
	// returns true immediately if the byte slice is already in the filter,
	// if the byte slice is not in the filter, adds it and returns false
	Add([]byte) bool
	// returns true immediately if the string is already in the filter,
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
