package bloom

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type BloomRedis struct {
	client *redis.Client
	key    string
}

/*
 ADD returns true if the item was not previously in the filter, and false if it was.
*/

func (br *BloomRedis) Add(id []byte) bool {
	ctx := context.Background()
	cmd := br.client.BFAdd(ctx, br.key, id)
	inserted, err := cmd.Result()
	if err != nil {
		panic(err)
	}
	return !inserted
}

// returns true is added or false if already in the filter
func (br *BloomRedis) AddString(id string) bool {
	return br.Add([]byte(id))
}

// checks if a byte slice is in the filter
// returns true if the byte slice is in the filter or a false positive
// return false if the byte slice is not in the filter
func (br *BloomRedis) Exists(id []byte) bool {
	ctx := context.Background()
	cmd := br.client.BFExists(ctx, "bf_key", id)
	exists, err := cmd.Result()
	if err != nil {
		panic(err)
	}

	return exists
}

// checks if a string is in the filter
// returns true if the string is in the filter or a false positive
// return false if the string is not in the filter
func (br *BloomRedis) ExistsString(id string) bool {
	return br.Add([]byte(id))
}

func NewBloomFilterRedis(key string, rate float64, cap int64) (br *BloomRedis, err error) {

	ctx := context.Background()
	client := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		},
	)
	if client == nil {
		return nil, errors.New("failed to create redis client")
	}

	// check if the key already exists
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %s", err)
	}

	cmd := client.BFReserve(ctx, key, rate, cap)
	if cmd.Err() != nil {
		// if it doesn't create the error
		if cmd.Err().Error() != "ERR item exists" {
			return nil, fmt.Errorf("failed to create bloom filter: %s", cmd.Err())
		}
	}

	return &BloomRedis{client, key}, nil
}
