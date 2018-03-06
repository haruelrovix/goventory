package api

import "net/http"

// item type with SKU, Name and Total
type item struct {
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type items struct {
	Items []item `json:"items"`
}

// ItemsHandleFunc to be used as http.HandleFunc for Item API
func ItemsHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		writeJSON(w, items{Items: getItems()})
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}

// Get Items from database
func getItems() []item {
	it := []item{}
	rows, _ := DB.Query(`
		SELECT sku, name, COALESCE(SUM(amount),0) AS total
		FROM items i LEFT JOIN stock s 
		ON i.sku = s.item_sku GROUP BY sku;
	`)

	for rows.Next() {
		row := item{}
		rows.Scan(&row.SKU, &row.Name, &row.Total)
		it = append(it, row)
	}

	return it
}
