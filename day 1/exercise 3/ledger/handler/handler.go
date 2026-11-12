package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"study/ledger/models"
	"sync"
)

// без бз было тяжело, пришлось сделать так! не сдить других способов не знаю
var paypal = []models.Transaction{}
var mtx sync.Mutex
var text = "Операций не было, своершите хотя бы 1 что бы увидить историю покупок"

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

// создаем операцию
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
	defer mtx.Unlock()
	for _, v := range paypal {
		if tp.ID == v.ID {
			http.Error(w, "Операция с таким айди уже создана!", http.StatusBadRequest)
			return
		}
	}
	paypal = append(paypal, tp)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(paypal); err != nil {
		http.Error(w, "error Encoder", http.StatusInternalServerError)
		return
	}

}

// удаляем операцию
func HandlerTransactionDELETE(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Error Method", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	RealID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Error ID", http.StatusBadRequest)
		return
	}

	mtx.Lock()
	defer mtx.Unlock()
	var found1 = false
	for _, v := range paypal {
		if RealID == v.ID {
			found1 = true
		}
	}
	if !found1 {
		http.Error(w, "NOT FOUND", http.StatusNotFound)
		return
	}
	if len(paypal) == 0 {

		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte(text))
		if err != nil {
			http.Error(w, "Error Print", http.StatusInternalServerError)
			return
		}
		return
	}
	var paypalFAKE []models.Transaction
	var found bool = false
	for _, v := range paypal {
		if RealID == v.ID {
			fmt.Println("пользователь удалил значение")
			found = true

		} else {
			paypalFAKE = append(paypalFAKE, v)
		}
	}
	if !found {
		http.Error(w, "NOT FOUND", http.StatusNotFound)
		return
	}
	paypal = paypalFAKE
	w.WriteHeader(http.StatusNoContent)

}

// найти по айди
func HandlerTRansactionGETID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error Method", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	RealID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Error ID", http.StatusBadRequest)
		return
	}
	mtx.Lock()
	defer mtx.Unlock()
	var found1 bool = false
	for _, v := range paypal {
		if RealID == v.ID {
			found1 = true
		}
	}
	if !found1 {
		http.Error(w, "NOT FOUND", http.StatusNotFound)
		return
	}
	if len(paypal) == 0 {
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte(text))
		if err != nil {
			http.Error(w, "Error Print", http.StatusInternalServerError)
			return
		}
		return
	}
	var paypalFAKE models.Transaction
	var found bool = false
	for _, v := range paypal {
		if RealID == v.ID {
			paypalFAKE = v
			found = true
		}
	}
	if !found {
		http.Error(w, "NOT FOUND", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paypalFAKE); err != nil {
		http.Error(w, "error Encoder", http.StatusInternalServerError)
		return
	}

}
