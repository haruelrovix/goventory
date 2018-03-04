package api

import "net/http"

// Laporan Nilai Barang
type ItemReport struct {
	SKU    string `json:"sku"`
	Name   string `json:"name"`
	Amount string `json:"amount"`
	Total  string `json:"total"`
}

// ItemReportHandleFunc to be used as http.HandleFunc for Incoming Item API
func ItemReportHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		itemReport := GetItemReport()
		writeJSON(w, itemReport)
	default:
		writeDefaultResponse(w)
	}
}

// Get Item Report from database
func GetItemReport() []ItemReport {
	itemReport := []ItemReport{}
	rows, _ := DB.Query(`
		SELECT sku, name, (SUM(booking * price) * 1.0 / SUM(t.amount)) AS total,
					(SELECT SUM(s.amount)
					 FROM stock s
					 WHERE s.item_sku = t.transaction_sku) AS amount
		FROM transactions t
			INNER JOIN IncomingTransactions it
			ON t.id = it.transaction_id
		LEFT JOIN items i
		ON t.transaction_sku = i.sku
		WHERE transaction_code = 'BM'
		GROUP BY transaction_sku;
	`)

	for rows.Next() {
		row := ItemReport{}
		// rows.Scan(&row.SKU, &row.Name, &row.Amount, &row.Total)
		rows.Scan(&row.SKU, &row.Name, &row.Total, &row.Amount)
		itemReport = append(itemReport, row)
	}

	return itemReport
}
