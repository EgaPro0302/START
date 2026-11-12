package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"study/ledger/models"
)

func HedlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error Method", 405)
		return
	}
	w.Header().Set("Content-Type", "aplication/json")
	response := models.Health{
		Status: "ok",
	}
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "error json", http.StatusInternalServerError)
		return
	}
	fmt.Println("на сервере все хорошо!")

}
