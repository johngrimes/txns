package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

var (
	mutex sync.RWMutex
	data = make(map[*http.Request]map[interface{}]interface{})
	datat = make(map[*http.Request]int64)
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
	for k, _ := range r.Header {
		canonicalKey := http.CanonicalHeaderKey(k)
		fmt.Fprintf(w, "%s: %s\n", canonicalKey, r.Header.Get(canonicalKey))
	}
	requestBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		LogError(err)
	}
	fmt.Fprintln(w, string(requestBytes))
}

func IdentifyRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Set(r, "RequestID", UUID())
		handler.ServeHTTP(w, r)
	})
}

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logData := map[string]string{
			"eventType": "IncomingRequest",
			"requestID": Get(r, "RequestID").(string),
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

func Set(r *http.Request, key, val interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	if data[r] == nil {
		data[r] = make(map[interface{}]interface{})
		datat[r] = time.Now().Unix()
	}
	data[r][key] = val
}

func Get(r *http.Request, key interface{}) interface{} {
	mutex.RLock()
	if ctx := data[r]; ctx != nil {
		value := ctx[key]
		mutex.RUnlock()
		return value
	}
	mutex.RUnlock()
	return nil
}

func main() {
	http.HandleFunc("/txns", TxnsHandler)
	requestHandler := LogRequest(http.DefaultServeMux)
	requestHandler = IdentifyRequest(requestHandler)
	http.ListenAndServe("localhost:4000", requestHandler)
}
