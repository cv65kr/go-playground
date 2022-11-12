package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func registerSigHandler(signals ...os.Signal) <-chan struct{} {
	stopCh := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, signals...)

	go func() {
		sig := <-sigCh
		log.Println("Received signal:", sig)
		signal.Reset(signals...)
		close(stopCh)
	}()

	return stopCh
}

func Task(stopCh <-chan struct{}) {
	mainCh := make(chan int, 1)
	go func() {
		log.Println("START")
		time.Sleep(20 * time.Second)
		log.Println("END")
		mainCh <- 1
	}()

	for {
		select {
		case <-stopCh:
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			<-ctx.Done()
			log.Printf("Closed because of context", ctx.Err())
			return
		case <-mainCh:
			log.Printf("Process finished")
			return
		}
	}
}

func main() {
	log.Println("PID:", os.Getpid())
	stopCh := registerSigHandler(syscall.SIGINT, syscall.SIGTERM)
	Task(stopCh)
	log.Println("Program exited")
}
