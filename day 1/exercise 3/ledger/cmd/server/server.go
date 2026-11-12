package server

import (
	"fmt"
	"net/http"
	"study/ledger/handler"
)

func LaunchServer() error {

	http.HandleFunc("/health", handler.HendlerGet)
	http.HandleFunc("POST /transactions", handler.HandlerTransactionPost) //луна зарпетила использовать сторонние библиотеке пришлось писать так, по другому хз, хотлось бы через mux сделать
	http.HandleFunc("GET /transactions", handler.HandlerTransactionGet)
	http.HandleFunc("GET /transactions/{id}", handler.HandlerTRansactionGETID)
	http.HandleFunc("DELETE /transactions/{id}", handler.HandlerTransactionDELETE)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		if err == http.ErrServerClosed {
			//когда я нажимаю ctrl + C пишет ошибку, я ее обработал
			fmt.Println("Сервер успешно закрыт!")
			return nil
		}
		return err
	}
	return nil
}
