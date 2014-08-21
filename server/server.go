package main

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
	"net/http"
	"regexp"
	"time"
)

func TxnsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		UploadTxnsHandler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func UploadTxnsHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	CheckError(w, err)

	form, err := reader.ReadForm(1024 ^ 2)
	CheckError(w, err)

	file_header := form.File["txn_file"][0]
	file, err := file_header.Open()
	CheckError(w, err)

	db, err := sql.Open("postgres", "user=johngrimes dbname=prec sslmode=disable")
	CheckError(w, err)

	csv := csv.NewReader(file)
	for i := 0; true; i++ {
		row, err := csv.Read()
		if i == 0 { continue }
		if row == nil { break }
		concatRow := strings.Join(row, "")
		hasher := sha512.New()
		hasher.Write([]byte(concatRow))
		rowHash := hex.EncodeToString(hasher.Sum(nil))
		dbDate, err := DBifyDateString(row[0])
		if CheckError(w, err) { break }
		debitCents, err := ParseCurrencyString(row[2], true)
		if CheckError(w, err) { break }
		creditCents, err := ParseCurrencyString(row[3], true)
		if CheckError(w, err) { break }
		balanceCents, err := ParseCurrencyString(row[4], false)
		if CheckError(w, err) { break }
		if debitCents < 0 && creditCents < 0 {
			BadRequest(w, "Bad statement line encountered: non-zero debit and credit values.")
		}
		_, err = db.Query(`INSERT INTO txns (hash, date, description, debit_cents, credit_cents, balance_cents)
	SELECT $1, $2, $3, $4, $5, $6
	WHERE NOT EXISTS (
		SELECT id FROM txns WHERE hash = cast($1 AS varchar)
	);`, rowHash, dbDate, row[1], debitCents, creditCents, balanceCents)
		if CheckError(w, err) { break }
	}

	db.Close()
}

func ParseCurrencyString(s string, convertToUnsigned bool) (int64, error) {
	if len(s) == 0 { return 0, nil }
	decimalRe := regexp.MustCompile("\\.")
	negativeZeroRe := regexp.MustCompile("-0")
	s = decimalRe.ReplaceAllString(s, "")
	s = negativeZeroRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "0 ")
	amount, err := strconv.ParseInt(s, 10, 64)
	if convertToUnsigned && amount < 0 {
		return amount * -1, err
	} else {
		return amount, err
	}
}

func DBifyDateString(s string) (string, error) {
	parsedTime, err := time.Parse("2/01/2006", s)
	if err != nil { return "", err }
	return parsedTime.Format("2006-01-02"), nil
}

func IdentifyRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "RequestID", UUID())
		handler.ServeHTTP(w, r)
	})
}

func BadRequest(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}

func CheckError(w http.ResponseWriter, err error) bool {
	if err != nil {
		LogError(err)
		http.Error(w, "", http.StatusInternalServerError)
		return true
	} else {
		return false
	}
}

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logData := map[string]string{
			"eventType": "Request",
			"requestID": context.Get(r, "RequestID").(string),
			"remoteHost": r.RemoteAddr,
			"httpMethod": r.Method,
			"resource": r.URL.String(),
		}
		LogEvent(logData)
		handler.ServeHTTP(w, r)
	})
}

func LogError(e error) {
	logData := map[string]string{
		"eventType": "Error",
		"errorMessage": e.Error(),
	}
	LogEvent(logData)
}

func LogDebugMessage(m string) {
	logData := map[string]string{
		"eventType": "DebugMessage",
		"debugMessage": m,
	}
	LogEvent(logData)
}

func LogEvent(event map[string]string) {
	event["loggedAt"] = time.Now().Format(time.RFC3339Nano)
	logJSON, err := json.Marshal(event)
	if err != nil { panic(err) }
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
