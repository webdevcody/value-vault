package handlers

import (
	"fmt"
	"key-value-app/persistence"
	"net/http"
)

func DeleteKey(w http.ResponseWriter, r *http.Request) {
	context := GetRequestContext(r)
	key := r.PathValue("key")

	Log(context, fmt.Sprintf("DELETE /keys/%s", key))

	if err := persistence.DeleteKey(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
