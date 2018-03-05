package api

import (
	"net/http"
	"time"
)

// Laporan Nilai Barang
type ItemReport struct {
	SKU    string  `json:"sku"`
	Name   string  `json:"name"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
	Value  float64 `json:"value"`
}

type Summary struct {
	PrintDate   time.Time `json:"printdate"`
	TotalSKU    int       `json:"totalsku"`
	TotalAmount int       `json:"totalamount"`
	TotalValue  float64   `json:"totalvalue"`
}

type Report struct {
	Items   []ItemReport `json:"items"`
	Summary Summary      `json:"summary"`
}

// ItemReportHandleFunc to be used as http.HandleFunc for Report Item API
func ItemReportHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		report := CreateReport()
		writeJSON(w, report)
	default:
		writeDefaultResponse(w)
	}
}

// Create Report
func CreateReport() Report {
	itemReport := []ItemReport{}
	rows, _ := DB.Query(`
		SELECT sku, name, (SUM(booking * price) * 1.0 / SUM(t.amount)) AS price,
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

	summary := Summary{}
	for rows.Next() {
		// ItemReport
		row := ItemReport{}
		rows.Scan(&row.SKU, &row.Name, &row.Price, &row.Amount)
		row.Value = row.Price * float64(row.Amount)
		itemReport = append(itemReport, row)

		// Summary
		summary.TotalAmount += row.Amount
		summary.TotalValue += row.Value
	}
	summary.TotalSKU = len(itemReport)
	summary.PrintDate = time.Now()

	// Report
	report := Report{
		Items:   itemReport,
		Summary: summary,
	}

	return report
}
