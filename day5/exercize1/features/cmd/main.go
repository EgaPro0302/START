package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"study/core"
	"study/features/repositry"
	"study/features/service"
	"study/features/transport"

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

	if err := pool.Ping(ctx); err != nil {
		fmt.Println("Error pool Ping:", err)
		return
	}
	//если все хорошо то пишем success
	logger.Info("SUCCESS")
	fmt.Println("DB OK!")

	repo := repositry.NewTaskRepository(pool)
	serv := service.NewTaskService(repo)
	h := transport.NewTaskHandler(serv)
	r.HandleFunc("/task/{id}/completed", h.UpdateStatus).Methods("PATCH") //работает оишбок не вижу
	r.HandleFunc("/task/{id}", h.UPDATECONTENT).Methods("PUT")            //работает оишбок не вижу
	r.HandleFunc("/task/{id}", h.DELETE).Methods("DELETE")                //работает оишбок не вижу
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
