package function

import (
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://google.com/", http.StatusTemporaryRedirect)
}
