package api

import (
	"fmt"
	"net/http"
	"time"

	strftime "github.com/jehiah/go-strftime"
)

type transactionReport struct {
	OrderID   string  `json:"orderid"`
	TimeStamp string  `json:"timestamp"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Amount    int     `json:"amount"`
	Price     float64 `json:"price"`
	Purchase  float64 `json:"purchase"`
	Omzet     float64 `json:"omzet"`  // Amount * Price
	Profit    float64 `json:"profit"` // Omzet - Harga Beli * Jumlah
}

type salesSummary struct {
	PrintDate   string  `json:"printdate"`
	StartDate   string  `json:"startdate"`
	EndDate     string  `json:"enddate"`
	TotalSales  int     `json:"totalsales"`
	TotalAmount int     `json:"totalamount"`
	TotalOmzet  float64 `json:"totalomzet"`
	TotalProfit float64 `json:"totalprofit"`
}

type salesReport struct {
	Items   []transactionReport `json:"items"`
	Summary salesSummary        `json:"summary"`
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
		startDate := r.FormValue("startdate")
		from, ok := validateDate(startDate)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Should have startdate params in YYYY-MM-DD format.")
			return
		}
		endDate := r.FormValue("enddate")
		to, ok := validateDate(endDate)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Should have enddate params in YYYY-MM-DD format.")
			return
		}
		report := createSalesReport(startDate, endDate)
		report.Summary.StartDate = strftime.Format("%d %B %Y", from)
		report.Summary.EndDate = strftime.Format("%d %B %Y", to)
		WriteJSON(w, report)
	default:
		WriteDefaultResponse(w)
	}
}

// CreateSalesReport creates Transaction Report
func createSalesReport(startDate string, endDate string) salesReport {
	transaction := []transactionReport{}
	rows, _ := DB.Query(`
		SELECT order_id, timestamp, i.sku, name, amount, price,
					(amount * price) AS omzet, purchase,
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
			AND timestamp BETWEEN '` + startDate + ` 00:00:00' 
				AND '` + endDate + ` 23:59:59'
		ORDER BY timestamp;
	`)

	summary := salesSummary{}
	i := 0
	for rows.Next() {
		// TransactionReport
		row := transactionReport{}
		rows.Scan(
			&row.OrderID, &row.TimeStamp, &row.SKU, &row.Name, &row.Amount,
			&row.Price, &row.Omzet, &row.Purchase, &row.Profit,
		)
		row.Omzet = row.Price * float64(row.Amount)
		transaction = append(transaction, row)

		// Summary
		summary.TotalAmount += row.Amount
		summary.TotalOmzet += row.Omzet
		summary.TotalProfit += row.Profit

		if i > 0 && transaction[i].OrderID != transaction[i-1].OrderID {
			summary.TotalSales++
		} else if i == 0 {
			summary.TotalSales = 1
		}
		i++
	}
	summary.PrintDate = strftime.Format("%d %B %Y", time.Now())

	// Report
	report := salesReport{
		Items:   transaction,
		Summary: summary,
	}
	rows.Close()

	return report
}
