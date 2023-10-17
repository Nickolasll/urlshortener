package presentation

import "net/http"

func mainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		PostHandler(res, req)
	} else if req.Method == http.MethodGet {
		GetHandler(res, req)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}
