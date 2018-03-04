package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Item type with SKU and Name
type Item struct {
	SKU  string `json:"sku"`
	Name string `json:"name"`
}

var DB *sql.DB

// ItemsHandleFunc to be used as http.HandleFunc for Item API
func ItemsHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		items := GetItems()
		writeJSON(w, items)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}

// Get Items from database
func GetItems() []Item {
	items := []Item{}
	rows, _ := DB.Query("SELECT sku, name FROM items")
	for rows.Next() {
		row := Item{}
		rows.Scan(&row.SKU, &row.Name)
		items = append(items, row)
	}

	return items
}

func writeJSON(w http.ResponseWriter, i interface{}) {
	result, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(result)
}
