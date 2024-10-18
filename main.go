package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	server := &http.Server{Addr: ":8080"}
	pgURL := os.Getenv("POSTGRES_URL")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if pgURL == "" {
			io.WriteString(w, "hello world\n")
			return
		}

		ctx := r.Context()
		conn, err := pgx.Connect(ctx, pgURL)
		if err != nil {
			io.WriteString(w, fmt.Sprintf("postgres connection failed: %v\n", err))
			return
		}
		defer conn.Close(ctx)

		err = conn.Ping(ctx)
		if err != nil {
			io.WriteString(w, fmt.Sprintf("postgres ping failed: %v\n", err))
			return
		}

		io.WriteString(w, "hello from postgres\n")

	})

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown(): %v", err)
	}
}
