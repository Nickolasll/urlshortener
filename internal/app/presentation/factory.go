package presentation

import "net/http"

func MuxFactory() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	return mux
}
