package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Item type with SKU and Name
type Item struct {
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Total string `json:"total"`
}

type IncomingItem struct {
	Timestamp string `json:"timestamp"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Booking   string `json:"booking"`
	Amount    string `json:"amount"`
	Price     string `json:"price"`
	Total     string `json:"total"`
	Receipt   string `json:"receipt"`
	Note      string `json:"note"`
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

// ItemsHandleFunc to be used as http.HandleFunc for Incoming Item API
func IncomingItemsHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		incomingItems := GetIncomingItems()
		writeJSON(w, incomingItems)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}

// Get Items from database
func GetItems() []Item {
	items := []Item{}
	rows, _ := DB.Query(`
		SELECT sku, name, COALESCE(SUM(amount),0) AS total
		FROM items i LEFT JOIN stock s 
		ON i.sku = s.item_sku GROUP BY sku;
	`)

	for rows.Next() {
		row := Item{}
		rows.Scan(&row.SKU, &row.Name, &row.Total)
		items = append(items, row)
	}

	return items
}

// Get Incoming Items from database
func GetIncomingItems() []IncomingItem {
	incomingItems := []IncomingItem{}
	rows, _ := DB.Query(`
		SELECT timestamp, sku, name, booking, amount, price,
					 ( booking * price ) AS total, receipt, note
		FROM transactions t
			INNER JOIN IncomingTransactions it
			ON t.id = it.transaction_id
		LEFT JOIN items i
		ON t.transaction_sku = i.sku
		WHERE transaction_code = 'BM';
	`)

	for rows.Next() {
		row := IncomingItem{}
		rows.Scan(
			&row.Timestamp, &row.SKU, &row.Name, &row.Booking, &row.Amount,
			&row.Price, &row.Total, &row.Receipt, &row.Note,
		)
		incomingItems = append(incomingItems, row)
	}

	return incomingItems
}

func writeJSON(w http.ResponseWriter, i interface{}) {
	result, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(result)
}
