package api

import (
	"net/http"
	"time"

	strftime "github.com/jehiah/go-strftime"
)

type itemReport struct {
	SKU    string  `json:"sku"`
	Name   string  `json:"name"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
	Value  float64 `json:"value"`
}

type summary struct {
	PrintDate   string  `json:"printdate"`
	TotalSKU    int     `json:"totalsku"`
	TotalAmount int     `json:"totalamount"`
	TotalValue  float64 `json:"totalvalue"`
}

type report struct {
	Items   []itemReport `json:"items"`
	Summary summary      `json:"summary"`
}

// ItemReportHandleFunc to be used as http.HandleFunc for Report Item API
func ItemReportHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		report := createReport()
		writeJSON(w, report)
	default:
		writeDefaultResponse(w)
	}
}

// CreateReport produces Report for Nilai Barang
func createReport() report {
	items := []itemReport{}
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

	summary := summary{}
	for rows.Next() {
		// ItemReport
		row := itemReport{}
		rows.Scan(&row.SKU, &row.Name, &row.Price, &row.Amount)
		row.Value = row.Price * float64(row.Amount)
		items = append(items, row)

		// Summary
		summary.TotalAmount += row.Amount
		summary.TotalValue += row.Value
	}
	summary.TotalSKU = len(items)
	summary.PrintDate = strftime.Format("%d %B %Y", time.Now())

	// Report
	report := report{
		Items:   items,
		Summary: summary,
	}
	rows.Close()

	return report
}
