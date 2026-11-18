package main

import (
	"errors"
	"fmt"
	"net/http"
	createlogs "study/Server/createLOGS"
	"study/Server/handler"
)

func main() {
	logger, LogFileClose, err := createlogs.NewLOgger("DEBUG")
	if err != nil {
		panic(err)
	}
	defer LogFileClose()
	http.HandleFunc("POST /strings", handler.HandlerPostString)
	http.HandleFunc("GET /strings", handler.HanderGetString)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Success!")
			return
		}
		logger.Debug("ОШИБКА СЕРВЕРА!")
		return
	}

}
