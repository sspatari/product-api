package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sspatari/product-api/handlers"
)

func main() {
	port := 9090
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)

	ph := handlers.NewProduct(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				l.Fatal(err)
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(tc); err != nil {
		l.Fatal(err)
	}

	l.Println("Shutdown completed")
}
