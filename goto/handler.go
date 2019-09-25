package function

import (
	"io/ioutil"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	// http.Redirect(w, r, "https://google.com/", http.StatusTemporaryRedirect)
	apiKey, err := ioutil.ReadFile("/var/openfaas/secrets/api-key")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(apiKey))
}
