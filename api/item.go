package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

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
		WriteJSON(w, items{Items: getItems()})
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		item := fromJSON(body)
		sku, created := createItem(item)
		if created {
			w.Header().Add("Location", "/api/items/"+sku)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
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

// fromJSON to be used for unmarshalling of Item type
func fromJSON(data []byte) item {
	item := item{}
	err := json.Unmarshal(data, &item)
	if err != nil {
		panic(err)
	}

	return item
}

// createItem creates a new Item if it does not exist
func createItem(it item) (string, bool) {
	_, err := DB.Exec(
		"INSERT INTO items (sku, name) VALUES (?, ?)",
		it.SKU, it.Name,
	)

	if err != nil {
		return "", false
	}

	return it.SKU, true
}
