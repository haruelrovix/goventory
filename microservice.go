package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"goventory/api"
)

// Config
const driver = "sqlite3"
const database = "inventory.db"

var db *sql.DB

func main() {
	// Connect to database
	db, err := sql.Open(driver, database)
	if err != nil {
		panic(err)
	}
	api.DB = db

	// Handle HTTP request
	http.HandleFunc("/", index)
	http.HandleFunc("/api/items", api.ItemsHandleFunc)
	http.HandleFunc("/api/barangmasuk", api.IncomingItemsHandleFunc)
	http.HandleFunc("/api/barangkeluar", api.OutgoingItemsHandleFunc)
	http.HandleFunc("/api/nilaibarang", api.ItemReportHandleFunc)

	http.ListenAndServe(port(), nil)
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Inventory REST API.")
}
