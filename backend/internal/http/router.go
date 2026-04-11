package httpapi

import (
	"fmt"
	"net/http"

	"dialecticlab/backend/internal/eval"
)

func NewRouter(engine *eval.Engine) (http.Handler, error) {
	handlers, err := NewHandlers(engine)
	if err != nil {
		return nil, fmt.Errorf("new router: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handlers.Health)
	mux.HandleFunc("/api/evaluate", handlers.Evaluate)
	return mux, nil
}
