package main

import (
	"fmt"
	"net/http"
	"study/transport"
	"time"
)

func main() {
	fmt.Println("Запуск сервера..")
	time.Sleep(2 * time.Second)

	http.HandleFunc("GET /node", transport.GetNode)
	http.HandleFunc("GET /node/{id}", transport.GetByID)
	http.HandleFunc("POST /node", transport.PostNode)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
