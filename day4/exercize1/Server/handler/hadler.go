package handler

import (
	"encoding/json"
	"net/http"
	createlogs "study/Server/createLOGS"
	"study/Server/models"
)

var FileString []models.Strings

func HandlerPostString(w http.ResponseWriter, r *http.Request) {
	logger, logFileClose, err := createlogs.NewLOgger("DEBUG")
	if err != nil {
		panic(err)
	}
	defer logFileClose()

	var minifile models.Strings
	if r.Method != http.MethodPost {
		logger.Info("пользоватлеь, выбрал неправильный метод")
		http.Error(w, "Error Methods!", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&minifile); err != nil {
		logger.Info("User Bad request")
		http.Error(w, "Bad Request", 400)
		return
	}
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(minifile); err != nil {
		logger.Debug("error Encoder")
		return
	}

	FileString = append(FileString, minifile)
}

func HanderGetString(w http.ResponseWriter, r *http.Request) {
	logger, logFileClose, err := createlogs.NewLOgger("DEBUG")
	if err != nil {
		panic(err)
	}
	defer logFileClose()
	if r.Method != http.MethodGet {
		logger.Info("пользователь выбрал неверный method!")
		http.Error(w, "Error Method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FileString)
}
