package server

import (
	"fmt"
	"net/http"
	"study/ledger/handler"
)

func LaunchServer() error {

	http.HandleFunc("/health", handler.HedlerGet)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("Сервер успешно закрыт!")
			return nil
		}
		return err
	}
	return nil
}
