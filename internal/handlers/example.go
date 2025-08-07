package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/phizhang/mini-tiktok/internal/db"
	"github.com/phizhang/mini-tiktok/internal/models"
)

// GetExampleHandler returns an example record from Cassandra
func GetExampleHandler(w http.ResponseWriter, r *http.Request) {
	var result models.ExampleModel
	if db.Session == nil {
		http.Error(w, "Cassandra not connected", http.StatusInternalServerError)
		return
	}
	iter := db.Session.Query("SELECT id, name FROM example_table LIMIT 1").Iter()
	for iter.Scan(&result.ID, &result.Name) {
		break
	}
	if err := iter.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
