package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"net/http"
	"net/http/httputil"
)

func TxnsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetTxnsHandler(w, r)
	} else if r.Method == "POST" {
		UploadTxnsHandler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetTxnsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

func UploadTxnsHandler(w http.ResponseWriter, r *http.Request) {
	requestBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		LogError(err)
	}
	fmt.Fprintln(w, string(requestBytes))
}

func IdentifyRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "RequestID", UUID())
		handler.ServeHTTP(w, r)
	})
}

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logData := map[string]string{
			"eventType": "IncomingRequest",
			"requestID": context.Get(r, "RequestID").(string),
			"remoteHost": r.RemoteAddr,
			"httpMethod": r.Method,
			"resource": r.URL.String(),
		}
		logJSON, err := json.Marshal(logData)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(logJSON))
		handler.ServeHTTP(w, r)
	})
}

func LogError(e error) {
	logData := map[string]string{
		"eventType": "Error",
		"errorMessage": e.Error(),
	}
	logJSON, err := json.Marshal(logData)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(logJSON))
}

func UUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func main() {
	http.HandleFunc("/txns", TxnsHandler)
	requestHandler := LogRequest(http.DefaultServeMux)
	requestHandler = IdentifyRequest(requestHandler)
	http.ListenAndServe("localhost:4000", requestHandler)
}
