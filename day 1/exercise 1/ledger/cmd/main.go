package main

import (
	"fmt"
	"study/ledger/cmd/server"
	"time"
)

func main() {
	//было созданы time для хоть какого то оживляения
	fmt.Println("Запуск сервера!")
	time.Sleep(2 * time.Second)
	fmt.Println("Сервер запущен")
	err := server.LaunchServer()
	fmt.Println(err)

}
