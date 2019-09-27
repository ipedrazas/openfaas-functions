package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// PlusOne contains 3 attributes: userid (hash), the topic and the counter
type PlusOne struct {
	userID  string
	topic   string
	counter int64
}

func setup() error {
	endpoint, err := ioutil.ReadFile("/var/openfaas/secrets/redis-endpoint")
	if err != nil {
		return err
	}
	pwd, err := ioutil.ReadFile("/var/openfaas/secrets/redis-password")
	if err != nil {
		return err
	}
	initialize(string(endpoint), string(pwd))
	return nil
}

// Handle the requests
func Handle(w http.ResponseWriter, r *http.Request) {

	var userid string
	err := setup()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
		entry.counter, err = increaseTopic(userid, topic)
		if err = r.ParseForm(); err != nil {
			fmt.Fprintf(w, "redis.incr error: %v", err)
			return
		}
		m := []byte("topic: " + userid + ", " + strconv.Itoa(int(entry.counter)))
		w.Write(m)

	}
	msg := []byte("userid: " + userid)
	w.Write(msg)

	keys, err := getKeys(userid+"*", client)
	for _, k := range keys {
		w.Write([]byte(k))
	}

	entries := getEntries(keys, client)

	byteSlice, err := json.Marshal(entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(byteSlice)

}
