package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // Easy connect to SQLite3
)

// DB driver
var DB *sql.DB

func writeDefaultResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Unsupported request method."))
}

func writeJSON(w http.ResponseWriter, i interface{}) {
	result, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(result)
}
