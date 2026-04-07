package main

import (
	"log"
	"net/http"
	"os"

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
		Addr:    ":" + port,
		Handler: httpapi.NewRouter(engine),
	}

	log.Printf("Dialectic Lab backend running on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
