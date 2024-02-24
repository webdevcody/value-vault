package handlers

import (
	"key-value-app/persistence"
	"net/http"
)

func DeleteKey(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	if err := persistence.DeleteKey(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
