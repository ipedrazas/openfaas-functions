package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	entry := &PlusOne{}

	err := setup()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		entry.userID = r.URL.Query().Get("uid")
		fmt.Fprintf(w, "r.URL.Query().Get(\"uid\"): %s \n", r.URL.Query().Get("uid"))
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		entry = &PlusOne{
			userID: r.FormValue("userid"),
			topic:  r.FormValue("topic"),
		}

		err = increaseTopic(entry)
		if err != nil {
			fmt.Fprintf(w, "redis.incr error: %v\n", err)
			return
		}
		fmt.Fprintf(w, "topic: %s, %d\n", entry.userID, entry.counter)

	}

	fmt.Fprintf(w, "userid: %s\n", entry.userID)

	keys, err := getKeys(entry.userID+"*", client)
	if err != nil {
		fmt.Fprintf(w, "redis.incr error: %v\n", err)
		return
	}
	for _, k := range keys {
		fmt.Fprintf(w, "keys: %s\n", k)
	}

	entries, err := getEntries(keys, client)
	if err != nil {
		fmt.Fprintf(w, "redis.getEntries: %v\n", err)
		return
	}
	for _, e := range entries {
		fmt.Fprintf(w, "entry uid: %s\n", e.userID)
		fmt.Fprintf(w, "entry topic: %s\n", e.topic)
		fmt.Fprintf(w, "entry counter: %d\n", e.counter)
	}

	byteSlice, err := json.Marshal(entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "entries: %s\n", string(byteSlice))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(byteSlice)

}
