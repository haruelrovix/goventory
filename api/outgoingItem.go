package api

import (
	"net/http"
	"time"
)

type outgoingItem struct {
	Timestamp time.Time `json:"timestamp"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	Price     float64   `json:"price"`
	Total     float64   `json:"total"`
	Note      string    `json:"note"`
}

type outgoingItems struct {
	Items []outgoingItem `json:"outgoingitems"`
}

// OutgoingItemsHandleFunc to be used as http.HandleFunc for Outgoing Item API
func OutgoingItemsHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		OutgoingItems := outgoingItems{Items: getOutgoingItems()}
		writeJSON(w, OutgoingItems)
	default:
		writeDefaultResponse(w)
	}
}

// GetOutgoingItems from database
func getOutgoingItems() []outgoingItem {
	items := []outgoingItem{}
	rows, _ := DB.Query(`
		SELECT timestamp, sku, name, amount, price,
					 ( amount * price ) AS total, note
		FROM transactions t
		LEFT JOIN items i
		ON t.transaction_sku = i.sku
		WHERE transaction_code = 'BK';
	`)

	for rows.Next() {
		row := outgoingItem{}
		rows.Scan(
			&row.Timestamp, &row.SKU, &row.Name, &row.Amount,
			&row.Price, &row.Total, &row.Note,
		)
		items = append(items, row)
	}
	rows.Close()

	return items
}
