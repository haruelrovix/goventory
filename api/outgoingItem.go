package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type outgoingItem struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	Price     float64   `json:"price"`
	Total     float64   `json:"total"`
	Note      string    `json:"note"`
}

type outgoingItems struct {
	Items []outgoingItem `json:"outgoingitems"`
}

// OutgoingItemsHandleFunc to be used as http.HandleFunc for Outgoing Item API
func OutgoingItemsHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		OutgoingItems := outgoingItems{Items: getOutgoingItems()}
		WriteJSON(w, OutgoingItems)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		ot := fromJSONOutgoingItem(body)
		id, created := createOutgoingItem(ot)
		if created {
			w.Header().Add("Location", "/api/barangkeluar/"+strconv.Itoa(id))
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	default:
		WriteDefaultResponse(w)
	}
}

// OutgoingItemHandleFunc to be used as http.HandleFunc for an Outgoing Item API
func OutgoingItemHandleFunc(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/barangkeluar/"):]
	id, err := strconv.Atoi(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	switch method := r.Method; method {
	case http.MethodGet:
		ot, found := getOutgoingItem(id)
		if found {
			WriteJSON(w, ot)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		ot := fromJSONOutgoingItem(body)
		exists := updateOutgoingItem(id, ot)
		if exists {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		ok := deleteOutgoingItem(id)
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

// GetOutgoingItems from database
func getOutgoingItems() []outgoingItem {
	items := []outgoingItem{}
	rows, _ := DB.Query(`
		SELECT timestamp, sku, name, amount, price,
					 ( amount * price ) AS total, note
		FROM transactions t
		LEFT JOIN items i
		ON t.transaction_sku = i.sku
		WHERE transaction_code = 'BK';
	`)

	for rows.Next() {
		row := outgoingItem{}
		rows.Scan(
			&row.Timestamp, &row.SKU, &row.Name, &row.Amount,
			&row.Price, &row.Total, &row.Note,
		)
		items = append(items, row)
	}
	rows.Close()

	return items
}

// createOutgoingItem creates a new Outgoing Item if it does not exist
func createOutgoingItem(ot outgoingItem) (int, bool) {
	result, err := DB.Exec(
		`INSERT INTO transactions (amount, note, price, timestamp, transaction_code,
			transaction_sku) VALUES (?, ?, ?, ?, ?, ?)`,
		ot.Amount, ot.Note, ot.Price, convertDate(ot.Timestamp), "BK", ot.SKU,
	)
	if err != nil {
		return 0, false
	}

	id, _ := result.LastInsertId()

	// check note
	pesanan := ot.Note[:7]
	if pesanan == "Pesanan" {
		result, err = DB.Exec(
			`INSERT INTO OutgoingTransactions (transaction_id, order_id) VALUES (?, ?)`,
			id, ot.Note[8:len(ot.Note)],
		)
		if err != nil {
			return int(id), false
		}
	}
	return int(id), true
}

// fromJSONOutgoingItem to be used for unmarshalling of OutgoingItem type
func fromJSONOutgoingItem(data []byte) outgoingItem {
	ot := outgoingItem{}
	err := json.Unmarshal(data, &ot)
	if err != nil {
		panic(err)
	}

	return ot
}

// there is no timezone in Toko Ijah universe, remove the T and Z
func convertDate(timeStamp time.Time) string {
	return timeStamp.Format("2006-01-02 15:04:05")
}

// getOutgoingItem returns the outgoing item for a given id
func getOutgoingItem(id int) (outgoingItem, bool) {
	ot := outgoingItem{ID: id}
	err := DB.QueryRow(`
		SELECT amount, note, price, timestamp, sku, name, (amount * price) AS total
		FROM transactions t
		LEFT JOIN items i
		ON t.transaction_sku = i.sku
		WHERE id = ?`, id,
	).Scan(&ot.Amount, &ot.Note, &ot.Price, &ot.Timestamp, &ot.SKU, &ot.Name, &ot.Total)
	if err != nil {
		return outgoingItem{}, false
	}

	return ot, true
}

// updateOutgoingItem updates an existing outgoing item
func updateOutgoingItem(id int, ot outgoingItem) bool {
	sid := strconv.Itoa(id)
	_, err := DB.Exec(
		"UPDATE transactions SET amount = " + strconv.Itoa(ot.Amount) +
			", transaction_sku = '" + ot.SKU + "', note = '" + ot.Note +
			"', timestamp = '" + convertDate(ot.Timestamp) +
			"', price = '" + strconv.FormatFloat(ot.Price, 'f', 0, 64) +
			"' WHERE id = " + sid + "",
	)
	if err != nil {
		return false
	}

	// check note
	pesanan := ot.Note[:7]
	if pesanan == "Pesanan" {
		_, err = DB.Exec(
			"UPDATE OutgoingTransactions SET order_id = '" + ot.Note[8:len(ot.Note)] +
				"'WHERE transaction_id = '" + sid + "'",
		)
		if err != nil {
			return false
		}
	} else {
		return deleteOutgoingTransaction(sid)
	}
	return true
}

// deleteOutgoingItem removes an outgoing item from transactions table
func deleteOutgoingItem(id int) bool {
	sid := strconv.Itoa(id)
	_, err := DB.Exec("DELETE FROM transactions WHERE id = '" + sid + "'")
	if err != nil {
		return false
	}
	return deleteOutgoingTransaction(sid)
}

// delete correspondence transaction
func deleteOutgoingTransaction(sid string) bool {
	_, err := DB.Exec(
		"DELETE FROM OutgoingTransactions WHERE transaction_id = '" +
			sid + "'",
	)
	if err != nil {
		return false
	}
	return true
}
