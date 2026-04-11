package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"dialecticlab/backend/internal/domain"
	"dialecticlab/backend/internal/eval"
)

// ErrNilEvaluator is returned by NewHandlers when the evaluator argument is nil.
var ErrNilEvaluator = errors.New("httpapi: evaluator must not be nil")

const maxEvaluateRequestBodyBytes int64 = 1 << 20

type Evaluator interface {
	Evaluate(graph domain.Graph, params domain.Params) (eval.Result, error)
}

type Handlers struct {
	evaluator Evaluator
}

func NewHandlers(evaluator Evaluator) (*Handlers, error) {
	if evaluator == nil {
		return nil, ErrNilEvaluator
	}
	// Typed nil pointers (e.g. var e *eval.Engine; NewHandlers(e)) are non-nil as interfaces.
	ev := reflect.ValueOf(evaluator)
	if (ev.Kind() == reflect.Pointer || ev.Kind() == reflect.Interface) && ev.IsNil() {
		return nil, ErrNilEvaluator
	}
	return &Handlers{evaluator: evaluator}, nil
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (h *Handlers) Evaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxEvaluateRequestBodyBytes)
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req EvaluateRequest
	if err := decoder.Decode(&req); err != nil {
		writeJSONDecodeError(w, err)
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid JSON body", []string{"request body must contain a single JSON object"})
		return
	}

	result, err := h.evaluator.Evaluate(domain.Graph{
		Nodes: req.Nodes,
		Edges: req.Edges,
	}, req.Params.ToDomain())
	if err != nil {
		if validationErr, ok := eval.IsValidationError(err); ok {
			writeError(w, http.StatusBadRequest, "graph validation failed", validationErr.Messages)
			return
		}

		log.Printf("evaluation failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to evaluate graph", []string{})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeJSONDecodeError(w http.ResponseWriter, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		writeError(
			w,
			http.StatusRequestEntityTooLarge,
			"request body too large",
			[]string{fmt.Sprintf("max request size is %d bytes", maxEvaluateRequestBodyBytes)},
		)
		return
	}

	writeError(w, http.StatusBadRequest, "invalid JSON body", []string{err.Error()})
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
