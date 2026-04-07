package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"dialecticlab/backend/internal/domain"
	"dialecticlab/backend/internal/eval"
)

type Handlers struct {
	engine *eval.Engine
}

func NewHandlers(engine *eval.Engine) *Handlers {
	return &Handlers{engine: engine}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (h *Handlers) Evaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	var req EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body", []string{err.Error()})
		return
	}

	result, err := h.engine.Evaluate(domain.Graph{
		Nodes: req.Nodes,
		Edges: req.Edges,
	}, req.Params.ToDomain())
	if err != nil {
		if validationErr, ok := eval.IsValidationError(err); ok {
			writeError(w, http.StatusBadRequest, "graph validation failed", validationErr.Messages)
			return
		}

		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			writeError(w, http.StatusBadRequest, "invalid JSON syntax", []string{fmt.Sprintf("offset: %d", syntaxErr.Offset)})
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to evaluate graph", []string{err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeMethodNotAllowed(w http.ResponseWriter, allowedMethod string) {
	w.Header().Set("Allow", allowedMethod)
	writeError(w, http.StatusMethodNotAllowed, "method not allowed", []string{fmt.Sprintf("use %s", allowedMethod)})
}

func writeError(w http.ResponseWriter, statusCode int, message string, errors []string) {
	resp := APIErrorResponse{
		Message: message,
	}
	resp.Diagnostics.Warnings = []string{}
	resp.Diagnostics.Errors = errors
	writeJSON(w, statusCode, resp)
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}
