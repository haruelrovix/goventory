package api

import (
	"fmt"
	"net/http"
	"time"

	strftime "github.com/jehiah/go-strftime"
)

// Laporan Penjualan
type TransactionReport struct {
	OrderID   string  `json:"orderid"`
	TimeStamp string  `json:"timestamp"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Amount    int     `json:"amount"`
	Price     float64 `json:"price"`
	Purchase  float64 `json:"purchase"`
	Value     float64 `json:"value"`  // Amount * Price
	Profit    float64 `json:"profit"` // Value - Harga Beli * Jumlah
}

type SalesSummary struct {
	PrintDate   string  `json:"printdate"`
	TimeFrom    string  `json:"timefrom"`
	TimeTo      string  `json:"timeto"`
	TotalSales  int     `json:"totalsales"`
	TotalAmount int     `json:"totalamount"`
	TotalValue  float64 `json:"totalvalue"`
}

type SalesReport struct {
	Items   []TransactionReport `json:"items"`
	Summary SalesSummary        `json:"summary"`
}

func validateDate(date string) (time.Time, bool) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return t, false
	}
	return t, true
}

// TransactionReportHandleFunc to be used as http.HandleFunc for Transaction Report API
func TransactionReportHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		timeFrom := r.FormValue("timefrom")
		from, ok := validateDate(timeFrom)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Beginning date should be on YYYY-MM-DD format.")
			return
		}
		timeTo := r.FormValue("timeto")
		to, ok := validateDate(timeTo)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "End date should be on YYYY-MM-DD format.")
			return
		}
		report := CreateSalesReport(timeFrom, timeTo)
		report.Summary.TimeFrom = strftime.Format("%d %B %Y", from)
		report.Summary.TimeTo = strftime.Format("%d %B %Y", to)
		writeJSON(w, report)
	default:
		writeDefaultResponse(w)
	}
}

// Create Transaction Report
func CreateSalesReport(timeFrom string, timeTo string) SalesReport {
	transactionReport := []TransactionReport{}
	rows, _ := DB.Query(`
		SELECT order_id, timestamp, i.sku, name, amount, price,
					(amount * price) AS value, purchase,
					((amount * price) - (purchase * amount)) AS profit
		FROM transactions t
			INNER JOIN OutgoingTransactions ot
			ON t.id = ot.transaction_id
			INNER JOIN items i
			ON t.transaction_sku = i.sku
			INNER JOIN (
				SELECT i.sku, (SUM(booking * price) * 1.0 / SUM(t.amount)) AS purchase
				FROM transactions t
					INNER JOIN IncomingTransactions it
					ON t.id = it.transaction_id
				LEFT JOIN items i
				ON t.transaction_sku = i.sku
				WHERE transaction_code = 'BM'
				GROUP BY i.sku
			) AS dt
			ON dt.sku = t.transaction_sku
		WHERE transaction_code = 'BK'
			AND timestamp BETWEEN '` + timeFrom + ` 00:00:00' 
				AND '` + timeTo + ` 23:59:59'
		ORDER BY timestamp;
	`)

	summary := SalesSummary{}
	i := 0
	for rows.Next() {
		// TransactionReport
		row := TransactionReport{}
		rows.Scan(
			&row.OrderID, &row.TimeStamp, &row.SKU, &row.Name, &row.Amount,
			&row.Price, &row.Value, &row.Purchase, &row.Profit,
		)
		row.Value = row.Price * float64(row.Amount)
		transactionReport = append(transactionReport, row)

		// Summary
		summary.TotalAmount += row.Amount
		summary.TotalValue += row.Value

		if i > 0 && transactionReport[i].OrderID != transactionReport[i-1].OrderID {
			summary.TotalSales += 1
		} else if i == 0 {
			summary.TotalSales = 1
		}
		i += 1
	}
	summary.PrintDate = strftime.Format("%d %B %Y", time.Now())

	// Report
	report := SalesReport{
		Items:   transactionReport,
		Summary: summary,
	}

	return report
}
