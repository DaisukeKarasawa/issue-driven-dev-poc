package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"dialecticlab/backend/internal/eval"
	httpapi "dialecticlab/backend/internal/http"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	engine := eval.NewEngine()
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           httpapi.NewRouter(engine),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Dialectic Lab backend running on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}