package function

import (
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

var (
	client = &redisClient{}
)

const separator = "|=|"

type redisClient struct {
	c *redis.Client
}

func initialize(addr string, pwd string) *redisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0,
	})

	if err := c.Ping().Err(); err != nil {
		panic("Unable to connect to redis " + err.Error())
	}
	client.c = c
	return client
}

func increaseTopic(entry *PlusOne) error {
	result, err := client.c.Incr(entry.userID + separator + entry.topic).Result()
	if err != nil {
		return err
	}
	entry.counter = result
	return nil
}

func getKeys(pattern string, client *redisClient) ([]string, error) {

	var cursor uint64
	var keys []string
	for {
		var ks []string
		var err error
		ks, cursor, err = client.c.Scan(cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		if cursor == 0 {
			break
		}
		keys = append(keys, ks...)
	}
	return keys, nil
}

func getEntries(keys []string, client *redisClient) ([]PlusOne, error) {
	var entries []PlusOne

	for _, key := range keys {
		res := strings.Split(key, separator)
		value := client.c.Get(key)
		if value.Err() != nil {
			return nil, value.Err()
		}

		count, _ := strconv.Atoi(value.Val())

		p := PlusOne{
			userID:  res[0],
			topic:   res[1],
			counter: int64(count),
		}
		entries = append(entries, p)

	}
	return entries, nil
}
