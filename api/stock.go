package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// stock type with SKU, Warehouse and Amount
type stock struct {
	ID        int    `json:"id"`
	SKU       string `json:"sku"`
	Warehouse int    `json:"warehouse"`
	Amount    int    `json:"total"`
}

type stocks struct {
	Stocks []stock `json:"stocks"`
}

// StocksHandleFunc to be used as http.HandleFunc for Stock API
func StocksHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		WriteJSON(w, stocks{Stocks: getStocks()})
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		stock := fromJSONStock(body)
		created := createStock(stock)
		if created {
			w.Header().Add("Location", "/api/stocks/")
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}

// StockHandleFunc to be used as http.HandleFunc for an Stock API
func StockHandleFunc(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/stocks/"):]
	stockID, err := strconv.Atoi(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	switch method := r.Method; method {
	case http.MethodGet:
		stock, found := getStock(stockID)
		if found {
			WriteJSON(w, stock)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		stock := fromJSONStock(body)
		exists := updateStock(stockID, stock)
		if exists {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		ok := deleteStock(stockID)
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

// Get Stocks from database
func getStocks() []stock {
	st := []stock{}
	rows, _ := DB.Query(`
		SELECT s.id, sku, w.id, amount
		FROM stock s
		INNER JOIN items i
		ON s.item_sku = i.sku
		INNER JOIN warehouses w
		ON s.warehouse_id = w.id;
	`)

	for rows.Next() {
		row := stock{}
		rows.Scan(&row.ID, &row.SKU, &row.Warehouse, &row.Amount)
		st = append(st, row)
	}

	return st
}

// fromJSONStock to be used for unmarshalling of Stock type
func fromJSONStock(data []byte) stock {
	stock := stock{}
	err := json.Unmarshal(data, &stock)
	if err != nil {
		panic(err)
	}

	return stock
}

// createStock creates a new Stock if it does not exist
func createStock(it stock) bool {
	_, err := DB.Exec(
		"INSERT INTO stock (item_sku, warehouse_id, amount) VALUES (?, ?, ?)",
		it.SKU, it.Warehouse, it.Amount,
	)

	if err != nil {
		return false
	}

	return true
}

// getStock returns the stock for a given ID
func getStock(stockID int) (stock, bool) {
	var warehouse int
	var sku string
	var amount int
	err := DB.QueryRow(`
		SELECT item_sku, warehouse_id, amount
		FROM stock WHERE id = ?`, stockID).Scan(&sku, &warehouse, &amount)
	if err != nil {
		return stock{}, false
	}

	return stock{ID: stockID, SKU: sku, Warehouse: warehouse, Amount: amount}, true
}

// deleteStock removes an stock from stock table
func deleteStock(stockID int) bool {
	_, err := DB.Exec("DELETE FROM stock WHERE id = '" + strconv.Itoa(stockID) + "'")
	if err != nil {
		return false
	}
	return true
}

// updateStock updates an existing stock
func updateStock(stockID int, it stock) bool {
	_, err := DB.Exec(
		"UPDATE stock SET item_sku = '" + it.SKU + "', warehouse_id = '" +
			strconv.Itoa(it.Warehouse) + "', amount = '" + strconv.Itoa(it.Amount) +
			"' WHERE id = '" + strconv.Itoa(stockID) + "'",
	)
	if err != nil {
		return false
	}
	return true
}
