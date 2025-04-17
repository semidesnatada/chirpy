package main

import (
	"net/http"
)


func healthzHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	w.Write([]byte(http.StatusText(200)))
}

