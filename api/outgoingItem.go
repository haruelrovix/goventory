package api

import "net/http"

// Barang Keluar
type OutgoingItem struct {
	Timestamp string `json:"timestamp"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Amount    string `json:"amount"`
	Price     string `json:"price"`
	Total     string `json:"total"`
	Note      string `json:"note"`
}

// OutgoingItemsHandleFunc to be used as http.HandleFunc for Outgoing Item API
func OutgoingItemsHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		OutgoingItems := GetOutgoingItems()
		writeJSON(w, OutgoingItems)
	default:
		writeDefaultResponse(w)
	}
}

// Get Outgoing Items from database
func GetOutgoingItems() []OutgoingItem {
	OutgoingItems := []OutgoingItem{}
	rows, _ := DB.Query(`
		SELECT timestamp, sku, name, amount, price,
					 ( amount * price ) AS total, note
		FROM transactions t
		LEFT JOIN items i
		ON t.transaction_sku = i.sku
		WHERE transaction_code = 'BK';
	`)

	for rows.Next() {
		row := OutgoingItem{}
		rows.Scan(
			&row.Timestamp, &row.SKU, &row.Name, &row.Amount,
			&row.Price, &row.Total, &row.Note,
		)
		OutgoingItems = append(OutgoingItems, row)
	}

	return OutgoingItems
}
