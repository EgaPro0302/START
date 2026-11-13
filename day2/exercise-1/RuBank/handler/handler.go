package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"study/RuBank/models"
	"study/RuBank/service"
)

type UserHandler struct {
	serv *service.UserService
}

func NewUserHandler(serv *service.UserService) *UserHandler {
	return &UserHandler{serv: serv}
}

var text = "User Not Found"

func (h *UserHandler) GetId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	RealId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		_, err1 := w.Write([]byte(text))
		if err1 != nil {
			http.Error(w, "Error on Server", 500)
			return
		}
		return
	}
	errortext := models.ErrorText{
		Error: "User Not Found",
	}
	if RealId <= 0 {
		http.Error(w, "Bad Request", 400)
		return
	}
	realUser, err := h.serv.GetId(r.Context(), RealId)
	if err != nil {
		if errors.Is(err, models.ErrUsersNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			if err := json.NewEncoder(w).Encode(errortext); err != nil {
				http.Error(w, "Error Encoder", 500)
				return
			}
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(realUser); err != nil {
		http.Error(w, "Error Encoder", 500)
		return
	}

}

func (h *UserHandler) CreateAcc(w http.ResponseWriter, r *http.Request) {
	var RealUser models.Users
	if err := json.NewDecoder(r.Body).Decode(&RealUser); err != nil {
		http.Error(w, "Error Dcoder", 500)
		return
	}
	if err := h.serv.CreateAcc(r.Context(), &RealUser); err != nil {
		http.Error(w, "Failed Create Account", 500)
		return
	}

	if err := json.NewEncoder(w).Encode(RealUser); err != nil {
		http.Error(w, "Error Encoder", 500)
		return
	}

}

func (h *UserHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var realCheck models.Check
	if err := json.NewDecoder(r.Body).Decode(&realCheck); err != nil {
		http.Error(w, "Error Dcoder", 500)
		return
	}
	if err := h.serv.Transfer(r.Context(), realCheck.From_account_id, realCheck.To_account_id, realCheck.Amount); err != nil {
		if errors.Is(err, models.ErrYourself) {
			http.Error(w, "Check Your Self", 400)
			return
		}
		if errors.Is(err, models.ErrSmallAmount) {
			http.Error(w, "amount small", 400)
			return

		}
		if errors.Is(err, models.ErrInvalidTransfer) {
			http.Error(w, "Bad Request", 400)
			return
		}
		return
	}

}
