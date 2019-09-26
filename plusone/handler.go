package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

// PlusOne contains 3 attributes: userid (hash), the topic and the counter
type PlusOne struct {
	userID  string
	topic   string
	counter int64
}

func init() {
	endpoint, err := ioutil.ReadFile("/var/openfaas/secrets/redis-endpoint")
	if err != nil {
		panic(err)
	}
	pwd, err := ioutil.ReadFile("/var/openfaas/secrets/redis-password")
	if err != nil {
		panic(err)
	}
	initialize(string(endpoint), string(pwd))
}

// Handle the requests
func Handle(w http.ResponseWriter, r *http.Request) {

	var userid string

	if r.Method == "GET" {
		userid = r.URL.Query().Get("uid")
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		userid = r.FormValue("userid")
		topic := r.FormValue("topic")
		entry := &PlusOne{
			userID: userid,
			topic:  topic,
		}
		entry.counter = increaseTopic(userid, topic)

	}

	keys, err := getKeys(userid + "*")
	entries := getEntries(keys)
	byteSlice, err := json.Marshal(entries)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(byteSlice)

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

func increaseTopic(userid string, topic string) int64 {
	result, err := client.c.Incr(userid + separator + topic).Result()
	if err != nil {
		panic(err)
	}
	return result
}

func getKeys(pattern string) ([]string, error) {

	var cursor uint64
	var keys []string
	for {
		var ks []string
		var err error
		ks, cursor, err = client.c.Scan(cursor, pattern, 10).Result()
		if err != nil {
			panic(err)
		}
		if cursor == 0 {
			break
		}
		keys = append(keys, ks...)
	}
	return keys, nil
}

func getEntries(keys []string) []PlusOne {
	var entries []PlusOne

	for _, key := range keys {
		res := strings.Split(key, separator)
		count, _ := strconv.Atoi(client.c.Get(key).Val())

		p := PlusOne{
			userID:  res[0],
			topic:   res[1],
			counter: int64(count),
		}
		entries = append(entries, p)

	}
	return entries
}
