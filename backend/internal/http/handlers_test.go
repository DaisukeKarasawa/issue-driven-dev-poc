package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dialecticlab/backend/internal/domain"
	"dialecticlab/backend/internal/eval"
)

type stubEvaluator struct {
	result eval.Result
	err    error
}

func (s stubEvaluator) Evaluate(_ domain.Graph, _ domain.Params) (eval.Result, error) {
	return s.result, s.err
}

func mustRouter(t *testing.T, engine *eval.Engine) http.Handler {
	t.Helper()
	router, err := NewRouter(engine)
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}
	return router
}

func TestNewHandlersRejectsNilEvaluator(t *testing.T) {
	_, err := NewHandlers(nil)
	if err == nil {
		t.Fatal("expected error for nil evaluator")
	}
	if !errors.Is(err, ErrNilEvaluator) {
		t.Fatalf("expected ErrNilEvaluator, got %v", err)
	}
}

func TestNewHandlersRejectsTypedNilEnginePointer(t *testing.T) {
	var engine *eval.Engine
	_, err := NewHandlers(engine)
	if err == nil {
		t.Fatal("expected error for typed-nil *eval.Engine")
	}
	if !errors.Is(err, ErrNilEvaluator) {
		t.Fatalf("expected ErrNilEvaluator, got %v", err)
	}
}

func TestHealthEndpoint(t *testing.T) {
	router := mustRouter(t, eval.NewEngine())
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status=ok, got %q", body["status"])
	}
}

func TestEvaluateEndpointSuccess(t *testing.T) {
	router := mustRouter(t, eval.NewEngine())

	payload := map[string]any{
		"nodes": []map[string]any{
			{"id": "a", "label": "Main claim", "kind": "claim"},
			{"id": "b", "label": "Evidence", "kind": "evidence"},
		},
		"edges": []map[string]any{
			{"from": "b", "to": "a", "relation": "support", "weight": 0.8},
		},
		"params": map[string]any{
			"damping":       0.35,
			"epsilon":       0.000001,
			"maxIterations": 200,
			"baseline":      0.5,
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/evaluate", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if _, ok := resp["scores"]; !ok {
		t.Fatalf("expected scores field in response, got %v", resp)
	}
	if _, ok := resp["meta"]; !ok {
		t.Fatalf("expected meta field in response, got %v", resp)
	}
}

func TestEvaluateEndpointValidationFailure(t *testing.T) {
	router := mustRouter(t, eval.NewEngine())
	payload := map[string]any{
		"nodes": []map[string]any{
			{"id": "a", "label": "Node A", "kind": "claim"},
		},
		"edges": []map[string]any{
			{"from": "missing", "to": "a", "relation": "attack", "weight": 0.7},
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/evaluate", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d body=%s", rec.Code, rec.Body.String())
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	diagnostics, ok := resp["diagnostics"].(map[string]any)
	if !ok {
		t.Fatalf("expected diagnostics object, got %v", resp["diagnostics"])
	}

	errorsField, ok := diagnostics["errors"].([]any)
	if !ok || len(errorsField) == 0 {
		t.Fatalf("expected non-empty diagnostics.errors, got %v", diagnostics["errors"])
	}
}

func TestEvaluateEndpointRejectsUnknownFields(t *testing.T) {
	router := mustRouter(t, eval.NewEngine())
	body := []byte(`{"nodes":[{"id":"a","label":"Node A","kind":"claim"}],"edges":[],"params":{},"extra":"nope"}`)

	req := httptest.NewRequest(http.MethodPost, "/api/evaluate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d body=%s", rec.Code, rec.Body.String())
	}

	var resp APIErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if resp.Message != "invalid JSON body" {
		t.Fatalf("expected invalid JSON message, got %q", resp.Message)
	}
	if len(resp.Diagnostics.Errors) == 0 || !strings.Contains(resp.Diagnostics.Errors[0], "unknown field") {
		t.Fatalf("expected unknown field diagnostics, got %v", resp.Diagnostics.Errors)
	}
}

func TestEvaluateEndpointRejectsOversizedBody(t *testing.T) {
	router := mustRouter(t, eval.NewEngine())
	largeLabel := strings.Repeat("x", int(maxEvaluateRequestBodyBytes))
	payload := map[string]any{
		"nodes": []map[string]any{
			{"id": "a", "label": largeLabel, "kind": "claim"},
		},
		"edges": []map[string]any{},
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}
	if int64(len(bodyBytes)) <= maxEvaluateRequestBodyBytes {
		t.Fatalf("expected oversized payload, got %d bytes", len(bodyBytes))
	}

	req := httptest.NewRequest(http.MethodPost, "/api/evaluate", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestEvaluateEndpointDoesNotLeakInternalErrors(t *testing.T) {
	handlers, err := NewHandlers(stubEvaluator{err: errors.New("internal failure detail")})
	if err != nil {
		t.Fatalf("NewHandlers: %v", err)
	}

	body := []byte(`{"nodes":[{"id":"a","label":"Node A","kind":"claim"}],"edges":[],"params":{}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/evaluate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.Evaluate(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d body=%s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "internal failure detail") {
		t.Fatalf("response leaked internal error details: %s", rec.Body.String())
	}

	var resp APIErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if resp.Message != "failed to evaluate graph" {
		t.Fatalf("unexpected message: %q", resp.Message)
	}
	if len(resp.Diagnostics.Errors) != 0 {
		t.Fatalf("expected no diagnostics errors for internal failure, got %v", resp.Diagnostics.Errors)
	}
}
