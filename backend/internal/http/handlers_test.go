package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dialecticlab/backend/internal/eval"
)

func TestHealthEndpoint(t *testing.T) {
	router := NewRouter(eval.NewEngine())
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
	router := NewRouter(eval.NewEngine())

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
	router := NewRouter(eval.NewEngine())
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
