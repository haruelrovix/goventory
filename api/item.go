package api

import "net/http"

// Item type with SKU and Name
type Item struct {
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type Items struct {
	Items []Item `json:"items"`
}

// ItemsHandleFunc to be used as http.HandleFunc for Item API
func ItemsHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		items := Items{Items: GetItems()}
		writeJSON(w, items)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}

// Get Items from database
func GetItems() []Item {
	items := []Item{}
	rows, _ := DB.Query(`
		SELECT sku, name, COALESCE(SUM(amount),0) AS total
		FROM items i LEFT JOIN stock s 
		ON i.sku = s.item_sku GROUP BY sku;
	`)

	for rows.Next() {
		row := Item{}
		rows.Scan(&row.SKU, &row.Name, &row.Total)
		items = append(items, row)
	}

	return items
}
