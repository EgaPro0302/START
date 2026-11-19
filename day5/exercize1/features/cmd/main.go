package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"study/core"
	repositry "study/features/repository"
	"study/features/service"
	"study/features/transport"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	//создаем нужные подключения
	var ctx = context.Background()
	r := mux.NewRouter()

	//подключаем логгирование
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()

	//делаем подлкючение
	fmt.Println("Создаем подключение..")
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error load .nv file")
	}
	Database := os.Getenv("DB")
	pool, err := pgxpool.New(ctx, Database)
	if err != nil {
		fmt.Println("Error poolNew:", err)
		return
	}
	defer pool.Close()

	PingCtx, close := context.WithTimeout(ctx, 5*time.Second)
	defer close()
	if err := pool.Ping(PingCtx); err != nil {
		fmt.Println("Error pool Ping:", err)
		return
	}
	//если все хорошо то пишем success
	logger.Info("SUCCESS")
	fmt.Println("DB OK!")

	repo := repositry.NewTaskRepository(pool, logger)
	serv := service.NewTaskService(repo, logger)
	h := transport.NewTaskHandler(serv, logger)
	r.HandleFunc("/task/{id}/completed", h.UpdateStatus).Methods("PATCH") //работает оишбок не вижу
	r.HandleFunc("/task/{id}", h.UpdateContent).Methods("PUT")            //работает оишбок не вижу
	r.HandleFunc("/task/{id}", h.Delete).Methods("DELETE")                //работает оишбок не вижу
	r.HandleFunc("/task/{id}", h.GetbyID).Methods("GET")                  //работает оишбок не вижу
	r.HandleFunc("/task", h.Create).Methods("POST")                       //работает оишбок не вижу
	r.HandleFunc("/task", h.GetAll).Methods("GET")                        //работает оишбок не вижу
	if err := http.ListenAndServe(":8080", r); err != nil {
		if err == http.ErrServerClosed {
			logger.Info("Server Closed")
			fmt.Println("До скорого!")
			return
		}
		logger.Error(fmt.Sprintf("%sError:", err))
		return
	}

}
