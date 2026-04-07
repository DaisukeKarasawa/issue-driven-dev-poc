package httpapi

import (
	"net/http"

	"dialecticlab/backend/internal/eval"
)

func NewRouter(engine *eval.Engine) http.Handler {
	handlers := NewHandlers(engine)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handlers.Health)
	mux.HandleFunc("/api/evaluate", handlers.Evaluate)
	return mux
}
