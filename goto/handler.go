package function

import (
	"io/ioutil"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	// http.Redirect(w, r, "https://google.com/", http.StatusTemporaryRedirect)
	apiKey, _ := ioutil.ReadFile("/var/openfaas/secrets/api-key")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(apiKey))
}
