package api

import "net/http"

// Barang Masuk
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

// IncomingItemsHandleFunc to be used as http.HandleFunc for Incoming Item API
func IncomingItemsHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		incomingItems := GetIncomingItems()
		writeJSON(w, incomingItems)
	default:
		writeDefaultResponse(w)
	}
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
