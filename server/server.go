package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"io"
	"net/http"
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
	reader, err := r.MultipartReader()
	CheckError(w, err)

	form, err := reader.ReadForm(1024 ^ 2)
	CheckError(w, err)

	file_header := form.File["txn_file"][0]
	file, err := file_header.Open()
	CheckError(w, err)

	fmt.Println(file)
	written, err := io.Copy(w, file)
	CheckError(w, err)

	LogDebugMessage(fmt.Sprintf("%d bytes were written.", written))
}

func IdentifyRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "RequestID", UUID())
		handler.ServeHTTP(w, r)
	})
}

func CheckError(w http.ResponseWriter, err error) {
	if err != nil {
		LogError(err)
		http.Error(w, "", http.StatusInternalServerError)
	}
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

func LogDebugMessage(m string) {
	logData := map[string]string{
		"eventType": "DebugMessage",
		"debugMessage": m,
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
