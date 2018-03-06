package api

import (
	"net/http"
	"time"
)

type incomingItem struct {
	Timestamp time.Time `json:"timestamp"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Booking   string    `json:"booking"`
	Amount    int       `json:"amount"`
	Price     float64   `json:"price"`
	Total     float64   `json:"total"`
	Receipt   string    `json:"receipt"`
	Note      string    `json:"note"`
}

type incomingItems struct {
	Items []incomingItem `json:"incomingitems"`
}

// IncomingItemsHandleFunc to be used as http.HandleFunc for Incoming Item API
func IncomingItemsHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		WriteJSON(w, incomingItems{Items: getIncomingItems()})
	default:
		WriteDefaultResponse(w)
	}
}

// GetIncomingItems from database
func getIncomingItems() []incomingItem {
	items := []incomingItem{}
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
		row := incomingItem{}
		rows.Scan(
			&row.Timestamp, &row.SKU, &row.Name, &row.Booking, &row.Amount,
			&row.Price, &row.Total, &row.Receipt, &row.Note,
		)
		items = append(items, row)
	}
	rows.Close()

	return items
}
