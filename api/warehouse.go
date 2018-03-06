package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// warehouse type with ID and Description
type warehouse struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type warehouses struct {
	Warehouses []warehouse `json:"warehouses"`
}

// WarehousesHandleFunc to be used as http.HandleFunc for Warehouse API
func WarehousesHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		WriteJSON(w, warehouses{Warehouses: getWarehouses()})
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		wh := fromJSONWarehouse(body)
		id, created := createWarehouse(wh)
		if created {
			w.Header().Add("Location", "/api/warehouses/"+id)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}

// WarehouseHandleFunc to be used as http.HandleFunc for an Warehouse API
func WarehouseHandleFunc(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/warehouses/"):]

	switch method := r.Method; method {
	case http.MethodGet:
		warehouse, found := getWarehouse(id)
		if found {
			WriteJSON(w, warehouse)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		wh := fromJSONWarehouse(body)
		exists := updateWarehouse(id, wh)
		if exists {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		ok := deleteWarehouse(id)
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

// Get Warehouses from database
func getWarehouses() []warehouse {
	wh := []warehouse{}
	rows, _ := DB.Query(`
		SELECT id, description
		FROM warehouses;
	`)

	for rows.Next() {
		row := warehouse{}
		rows.Scan(&row.ID, &row.Description)
		wh = append(wh, row)
	}

	return wh
}

// fromJSON to be used for unmarshalling of Warehouse type
func fromJSONWarehouse(data []byte) warehouse {
	warehouse := warehouse{}
	err := json.Unmarshal(data, &warehouse)
	if err != nil {
		panic(err)
	}

	return warehouse
}

// createWarehouse creates a new Warehouse if wh does not exist
func createWarehouse(wh warehouse) (string, bool) {
	_, err := DB.Exec(
		"INSERT INTO warehouses (id, description) VALUES (?, ?)",
		wh.ID, wh.Description,
	)

	if err != nil {
		return "", false
	}

	return wh.ID, true
}

// getWarehouse returns the warehouse for a given ID
func getWarehouse(id string) (warehouse, bool) {
	var description string
	err := DB.QueryRow("SELECT description FROM warehouses WHERE id = ?", id).Scan(&description)
	if err != nil {
		return warehouse{}, false
	}

	return warehouse{ID: id, Description: description}, true
}

// deleteWarehouse removes an warehouse from warehouses table
func deleteWarehouse(id string) bool {
	_, err := DB.Exec("DELETE FROM warehouses WHERE id = '" + id + "'")
	if err != nil {
		return false
	}
	return true
}

// updateWarehouse updates an existing warehouse
func updateWarehouse(id string, wh warehouse) bool {
	_, err := DB.Exec(
		"UPDATE warehouses SET id = '" + id + "', description = '" +
			wh.Description + "' WHERE id = '" + id + "'",
	)
	if err != nil {
		return false
	}
	return true
}
