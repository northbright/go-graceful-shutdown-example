package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
)

func home(w http.ResponseWriter, r *http.Request) {
	i := 0
	// Use a timeout to emulate a long time work.
	ch := time.After(time.Second * 10)

	for {
		select {
		case <-ctx.Done():
			log.Printf("home() stopped")
			return
		case <-ch:
			log.Printf("home() finished(timeout), i: %v", i)
			w.Write([]byte(fmt.Sprintf("i: %v", i)))
			return
		default:
			log.Printf("i: %v", i)
			time.Sleep(time.Second * 1)
			i++
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	server := &http.Server{Addr: "localhost:8080", Handler: mux}
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)

		// Wait interrupt signal
		<-sigint
		log.Printf("os.Interrupt received")

		// Call cancel func to stop all worker goroutines.
		cancel()

		// Use a new context to shutdown the server.
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}

	// Wait for idle connections closed(shutdown goroutine exited).
	<-idleConnsClosed
	log.Printf("main() exited")
}
