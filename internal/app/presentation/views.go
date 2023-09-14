package presentation

import "net/http"

func mainPage(res http.ResponseWriter, req *http.Request) {
	PostHandler(res, req)
	GetHandler(res, req)
	res.WriteHeader(http.StatusBadRequest)
}
