package eval

import (
	"math"
	"strings"
	"testing"

	"dialecticlab/backend/internal/domain"
)

func TestSupportEdgeRaisesTargetScore(t *testing.T) {
	engine := NewEngine()
	graph := domain.Graph{
		Nodes: []domain.Node{
			{ID: "A", Label: "Claim A", Kind: domain.NodeKindClaim},
			{ID: "B", Label: "Claim B", Kind: domain.NodeKindClaim},
		},
		Edges: []domain.Edge{
			{From: "A", To: "B", Relation: domain.RelationSupport, Weight: 1},
		},
	}
	params := domain.Params{Damping: 0.4, Epsilon: 1e-7, MaxIterations: 400, Baseline: 0.5}

	result, err := engine.Evaluate(graph, params)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	scoreMap := toMap(result.Scores)
	if scoreMap["B"] <= scoreMap["A"] {
		t.Fatalf("expected B > A with support edge, got B=%f A=%f", scoreMap["B"], scoreMap["A"])
	}
}

func TestStrongAttackLowersTargetScore(t *testing.T) {
	engine := NewEngine()
	graph := domain.Graph{
		Nodes: []domain.Node{
			{ID: "A", Label: "Counter", Kind: domain.NodeKindClaim},
			{ID: "B", Label: "Target", Kind: domain.NodeKindClaim},
		},
		Edges: []domain.Edge{
			{From: "A", To: "B", Relation: domain.RelationAttack, Weight: 1},
		},
	}
	params := domain.Params{Damping: 0.6, Epsilon: 1e-7, MaxIterations: 400, Baseline: 0.5}

	result, err := engine.Evaluate(graph, params)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	scoreMap := toMap(result.Scores)
	if scoreMap["B"] >= 0.5 {
		t.Fatalf("expected B < baseline 0.5 with strong attacker, got B=%f", scoreMap["B"])
	}
}

func TestCycleConvergesWithDamping(t *testing.T) {
	engine := NewEngine()
	graph := domain.Graph{
		Nodes: []domain.Node{
			{ID: "A", Label: "A", Kind: domain.NodeKindClaim},
			{ID: "B", Label: "B", Kind: domain.NodeKindClaim},
		},
		Edges: []domain.Edge{
			{From: "A", To: "B", Relation: domain.RelationSupport, Weight: 0.8},
			{From: "B", To: "A", Relation: domain.RelationAttack, Weight: 0.6},
		},
	}
	params := domain.Params{Damping: 0.35, Epsilon: 1e-6, MaxIterations: 1000, Baseline: 0.5}

	result, err := engine.Evaluate(graph, params)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !result.Meta.Converged {
		t.Fatalf("expected converged=true, got false; iterations=%d, maxDelta=%f", result.Meta.Iterations, result.Meta.MaxDelta)
	}
	if result.Meta.Iterations >= params.MaxIterations {
		t.Fatalf("expected to converge before max iterations, got iterations=%d", result.Meta.Iterations)
	}

	for _, score := range result.Scores {
		if score.Score < 0 || score.Score > 1 {
			t.Fatalf("expected score in [0,1], got %f for node %s", score.Score, score.NodeID)
		}
	}
}

func TestInvalidGraphRejected(t *testing.T) {
	engine := NewEngine()
	graph := domain.Graph{
		Nodes: []domain.Node{
			{ID: "A", Label: "A", Kind: domain.NodeKindClaim},
		},
		Edges: []domain.Edge{
			{From: "A", To: "A", Relation: domain.RelationSupport, Weight: 1},
			{From: "missing", To: "A", Relation: domain.RelationAttack, Weight: 0.9},
		},
	}
	params := domain.Params{Damping: 0.3, Epsilon: 1e-6, MaxIterations: 200, Baseline: 0.5}

	_, err := engine.Evaluate(graph, params)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	validationErr, ok := IsValidationError(err)
	if !ok {
		t.Fatalf("expected validation error type, got: %T", err)
	}
	if len(validationErr.Messages) == 0 {
		t.Fatal("expected validation messages to be non-empty")
	}
}

func TestInvalidNodeKindRejected(t *testing.T) {
	engine := NewEngine()
	graph := domain.Graph{
		Nodes: []domain.Node{
			{ID: "A", Label: "A", Kind: domain.NodeKind("unknown")},
		},
		Edges: []domain.Edge{},
	}
	params := domain.Params{Damping: 0.3, Epsilon: 1e-6, MaxIterations: 200, Baseline: 0.5}

	_, err := engine.Evaluate(graph, params)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	validationErr, ok := IsValidationError(err)
	if !ok {
		t.Fatalf("expected validation error type, got: %T", err)
	}

	found := false
	for _, message := range validationErr.Messages {
		if strings.Contains(message, "kind must be claim or evidence") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected kind validation message, got %v", validationErr.Messages)
	}
}

func TestNonFiniteEdgeWeightRejected(t *testing.T) {
	engine := NewEngine()
	graph := domain.Graph{
		Nodes: []domain.Node{
			{ID: "A", Label: "A", Kind: domain.NodeKindClaim},
			{ID: "B", Label: "B", Kind: domain.NodeKindClaim},
		},
		Edges: []domain.Edge{
			{From: "A", To: "B", Relation: domain.RelationSupport, Weight: math.NaN()},
		},
	}
	params := domain.Params{Damping: 0.3, Epsilon: 1e-6, MaxIterations: 200, Baseline: 0.5}

	_, err := engine.Evaluate(graph, params)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	validationErr, ok := IsValidationError(err)
	if !ok {
		t.Fatalf("expected validation error type, got: %T", err)
	}
	found := false
	for _, message := range validationErr.Messages {
		if strings.Contains(message, "edge[0] weight must be finite") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected finite weight message, got %v", validationErr.Messages)
	}
}

func TestNonFiniteParamsRejected(t *testing.T) {
	engine := NewEngine()
	graph := domain.Graph{
		Nodes: []domain.Node{
			{ID: "A", Label: "A", Kind: domain.NodeKindClaim},
		},
		Edges: []domain.Edge{},
	}
	params := domain.Params{
		Damping:       math.Inf(1),
		Epsilon:       1e-6,
		MaxIterations: 200,
		Baseline:      0.5,
	}

	_, err := engine.Evaluate(graph, params)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	validationErr, ok := IsValidationError(err)
	if !ok {
		t.Fatalf("expected validation error type, got: %T", err)
	}
	found := false
	for _, message := range validationErr.Messages {
		if strings.Contains(message, "params.damping must be finite") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected finite damping message, got %v", validationErr.Messages)
	}
}

func toMap(scores []NodeScore) map[string]float64 {
	out := make(map[string]float64, len(scores))
	for _, score := range scores {
		out[score.NodeID] = score.Score
	}
	return out
}
