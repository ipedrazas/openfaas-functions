package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// PlusOne contains 3 attributes: userid (hash), the topic and the counter
type PlusOne struct {
	UserID  string `json:"userid"`
	Topic   string `json:"topic"`
	Counter int64  `json:"counter"`
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
		entry.UserID = r.URL.Query().Get("uid")
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		entry = &PlusOne{
			UserID: r.FormValue("userid"),
			Topic:  r.FormValue("topic"),
		}

		err = increaseTopic(entry)
		if err != nil {
			fmt.Fprintf(w, "redis.incr error: %v\n", err)
			return
		}
	}
	keys, err := getKeys(entry.UserID+"*", client)
	if err != nil {
		fmt.Fprintf(w, "redis.getKeys error: %v\n", err)
		return
	}

	entries, err := getEntries(keys, client)
	if err != nil {
		fmt.Fprintf(w, "redis.getEntries: %v\n", err)
		return
	}

	byteSlice, err := json.Marshal(entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(byteSlice)

}
