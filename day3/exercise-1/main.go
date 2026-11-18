package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func Workind_count1(ctx context.Context) {
	defer wg.Done()
	user := 1
	i := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Job Done! 1")
			return
		case <-time.After(1 * time.Second):
			fmt.Println("я рабочий:", user, "я складываю каждое число на 1:", i)
			i++
		}
	}
}
func Workind_count2(ctx context.Context) {
	defer wg.Done()
	user := 2
	i := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Job Done! 2")
			return
		case <-time.After(1 * time.Second):
			fmt.Println("я рабочий:", user, "я вычитаю каждое число на 1:", i)
			i--
		}
	}
}
func Workind_count3(ctx context.Context) {
	defer wg.Done()
	user := 3
	i := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Job Done! 3")
			return
		case <-time.After(1 * time.Second):
			fmt.Println("я рабочий:", user, "я умножаю каждое число на 1:", i)
			i *= i
		}
	}

}
func Workind_count4(ctx context.Context) {
	defer wg.Done()
	user := 4
	i := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Job Done!")
			return
		case <-time.After(1 * time.Second):
			fmt.Println("я рабочий:", user, "я делю каждое число на 1:", i)
			i /= i
		}
	}
}
func main() {
	fmt.Println("PID", os.Getpid())
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	fmt.Println("запуск всех горутин")
	wg.Add(4)
	go Workind_count1(ctx)
	go Workind_count2(ctx)
	go Workind_count3(ctx)
	go Workind_count4(ctx)
	<-ctx.Done()
	wg.Wait()

	fmt.Println("job application correctly!")

}
