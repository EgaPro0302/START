package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	repository "study/RuBank/Repository"
	"study/RuBank/handler"
	"study/RuBank/service"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// создаем нужны праметра
	ctx := context.Background()
	r := mux.NewRouter()

	// делаем ОБЯЗАТЕЛЬНУЮ проверку env
	if err := godotenv.Load(); err != nil {
		fmt.Println("error Load env", err)
		return
	}

	//читаем ее
	DB_URL := os.Getenv("DB_url")
	fmt.Println(DB_URL)

	//подключаем pgx
	pool, err := pgxpool.New(ctx, DB_URL)
	if err != nil {
		fmt.Println("Error connect pgxpool:", err)
		return
	}

	//проверяем
	if err := pool.Ping(ctx); err != nil {
		fmt.Println("Error poolPing:", err)
		return
	}
	fmt.Println("BD OK!")

	//можно сказать делаем кабель подлкючений
	repo := repository.NewUserRepository(pool)
	serv := service.NewUserService(repo)
	hand := handler.NewUserHandler(serv)

	//создаем патерны
	r.HandleFunc("/accounts", hand.CreateAcc).Methods("POST")
	r.HandleFunc("/accounts/{id}", hand.GetId).Methods("GET")
	r.HandleFunc("/transfers", hand.Transfer).Methods("POST")
	//делвем поделючение на сервер
	if err := http.ListenAndServe(":8080", r); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server Success Closed!")
			return
		}
		fmt.Println("Error:", err)
		return
	}
}
