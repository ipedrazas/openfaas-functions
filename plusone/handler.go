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

func setup() {
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
	setup()

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

	keys, err := getKeys(userid+"*", client)
	entries := getEntries(keys, client)
	byteSlice, err := json.Marshal(entries)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(byteSlice)

}
