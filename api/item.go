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

// ItemHandleFunc to be used as http.HandleFunc for an Item API
func ItemHandleFunc(w http.ResponseWriter, r *http.Request) {
	sku := r.URL.Path[len("/api/items/"):]

	switch method := r.Method; method {
	case http.MethodGet:
		item, found := getItem(sku)
		if found {
			WriteJSON(w, item)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		item := fromJSON(body)
		exists := updateItem(sku, item)
		if exists {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		ok := deleteItem(sku)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
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

// getItem returns the item for a given SKU
func getItem(sku string) (item, bool) {
	var name string
	err := DB.QueryRow("SELECT name FROM items WHERE sku = ?", sku).Scan(&name)
	if err != nil {
		return item{}, false
	}

	return item{SKU: sku, Name: name}, true
}

// deleteItem removes an item from items table
func deleteItem(sku string) bool {
	_, err := DB.Exec("DELETE FROM items WHERE sku = '" + sku + "'")
	if err != nil {
		return false
	}

	return true
}

// updateItem updates an existing item
func updateItem(sku string, it item) bool {
	_, err := DB.Exec(
		"UPDATE items SET sku = '" + sku + "', name = '" +
			it.Name + "' WHERE sku = '" + sku + "'",
	)
	if err != nil {
		return false
	}
	return true
}
