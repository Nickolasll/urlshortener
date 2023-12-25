package presentation

import "net/http"

func mainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		postHandler(res, req)
	} else if req.Method == http.MethodGet {
		expandHandler(res, req)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}
