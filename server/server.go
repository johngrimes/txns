package main

import (
	"fmt"
	"net/http"
)

func getTxnsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(404)
	} else {
		fmt.Fprint(w, "Hello!")
	}
}

func uploadTxnsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
	} else {

	}
}

func main() {
	http.HandleFunc("/txns", getTxnsHandler)
	http.ListenAndServe("localhost:4000", nil)
}
