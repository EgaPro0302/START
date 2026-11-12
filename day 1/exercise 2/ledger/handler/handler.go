package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"study/ledger/models"
	"sync"
)

// без бз было тяжело, пришлось сделать так! не сдить других способов не знаю
var paypal = []models.Transaction{}
var mtx sync.Mutex

// проверить статус
func HendlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error Method", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := models.Health{
		Status: "ok",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "error Encoder", http.StatusInternalServerError)
		return
	}
	fmt.Println("на сервере все хорошо!")

}

// операции покупок
// Смотрим
func HandlerTransactionGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error Method,", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paypal); err != nil {
		http.Error(w, "error Encoder", http.StatusInternalServerError)
		return
	}

}

func HandlerTransactionPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error Method", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var tp models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tp); err != nil {
		http.Error(w, "Error Decoder", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	mtx.Lock()
	paypal = append(paypal, tp)
	mtx.Unlock()
	if err := json.NewEncoder(w).Encode(paypal); err != nil {
		http.Error(w, "error Encoder", http.StatusInternalServerError)
		return
	}
}
