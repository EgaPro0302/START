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

	text := "Операций не было, своершите хотя бы 1 что бы увидить историю покупок"
	mtx.Lock()
	defer mtx.Unlock()
	if len(paypal) == 0 {
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte(text))
		if err != nil {
			http.Error(w, "Error Print", http.StatusInternalServerError)
			return
		}
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
		http.Error(w, "Error Decoder", http.StatusBadRequest)
		return
	}
	mtx.Lock()
	paypal = append(paypal, tp)
	mtx.Unlock()
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(paypal); err != nil {
		http.Error(w, "error Encoder", http.StatusInternalServerError)
		return
	}

}

// найти по айди
func HandlerTRansactionGET(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error Method", http.StatusMethodNotAllowed)
		return
	}
}
