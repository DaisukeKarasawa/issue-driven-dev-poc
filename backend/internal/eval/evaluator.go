package eval

import (
	"errors"
	"fmt"
	"math"
	"slices"

	"dialecticlab/backend/internal/domain"
)

type NodeScore struct {
	NodeID string  `json:"nodeId"`
	Score  float64 `json:"score"`
}

type Meta struct {
	Converged  bool    `json:"converged"`
	Iterations int     `json:"iterations"`
	MaxDelta   float64 `json:"maxDelta"`
}

type Diagnostics struct {
	Warnings []string `json:"warnings"`
	Errors   []string `json:"errors"`
}

type Result struct {
	Scores      []NodeScore `json:"scores"`
	Meta        Meta        `json:"meta"`
	Diagnostics Diagnostics `json:"diagnostics"`
}

type ValidationError struct {
	Messages []string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("graph validation failed (%d issues)", len(e.Messages))
}

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Evaluate(graph domain.Graph, params domain.Params) (Result, error) {
	validation := validateInput(graph, params)
	if len(validation) > 0 {
		return Result{}, ValidationError{Messages: validation}
	}

	nodeByID := make(map[string]domain.Node, len(graph.Nodes))
	incoming := make(map[string][]domain.Edge, len(graph.Nodes))
	for _, n := range graph.Nodes {
		nodeByID[n.ID] = n
		incoming[n.ID] = nil
	}
	for _, edge := range graph.Edges {
		incoming[edge.To] = append(incoming[edge.To], edge)
	}

	prev := make(map[string]float64, len(graph.Nodes))
	next := make(map[string]float64, len(graph.Nodes))
	for _, n := range graph.Nodes {
		prev[n.ID] = params.Baseline
		next[n.ID] = params.Baseline
	}

	var maxDelta float64
	converged := false
	iterations := 0

	for i := 1; i <= params.MaxIterations; i++ {
		maxDelta = 0
		for _, n := range graph.Nodes {
			edgeSet := incoming[n.ID]
			supportInfluence := 0.0
			attackInfluence := 0.0

			for _, edge := range edgeSet {
				sourceScore := prev[edge.From]
				contribution := sourceScore * edge.Weight
				if edge.Relation == domain.RelationSupport {
					supportInfluence += contribution
					continue
				}
				attackInfluence += contribution
			}

			raw := params.Baseline + supportInfluence - attackInfluence
			blended := (1-params.Damping)*prev[n.ID] + params.Damping*raw
			next[n.ID] = clamp01(blended)

			delta := math.Abs(next[n.ID] - prev[n.ID])
			maxDelta = math.Max(maxDelta, delta)
		}

		iterations = i
		if maxDelta < params.Epsilon {
			converged = true
			break
		}

		for _, n := range graph.Nodes {
			prev[n.ID] = next[n.ID]
		}
	}

	scores := make([]NodeScore, 0, len(graph.Nodes))
	for _, n := range graph.Nodes {
		scores = append(scores, NodeScore{NodeID: n.ID, Score: round(next[n.ID], 6)})
	}
	slices.SortFunc(scores, func(a, b NodeScore) int {
		if b.Score == a.Score {
			return cmpString(a.NodeID, b.NodeID)
		}
		if b.Score > a.Score {
			return 1
		}
		return -1
	})

	return Result{
		Scores: scores,
		Meta: Meta{
			Converged:  converged,
			Iterations: iterations,
			MaxDelta:   round(maxDelta, 8),
		},
		Diagnostics: Diagnostics{
			Warnings: []string{},
			Errors:   []string{},
		},
	}, nil
}

func validateInput(graph domain.Graph, params domain.Params) []string {
	var errs []string
	if len(graph.Nodes) == 0 {
		errs = append(errs, "at least one node is required")
	}

	nodeIDs := map[string]struct{}{}
	for i, n := range graph.Nodes {
		if n.ID == "" {
			errs = append(errs, fmt.Sprintf("node[%d] id must not be empty", i))
			continue
		}
		if _, exists := nodeIDs[n.ID]; exists {
			errs = append(errs, fmt.Sprintf("duplicate node id: %s", n.ID))
			continue
		}
		nodeIDs[n.ID] = struct{}{}

		if n.Label == "" {
			errs = append(errs, fmt.Sprintf("node[%d] (%s) label must not be empty", i, n.ID))
		}
		if n.Kind != domain.NodeKindClaim && n.Kind != domain.NodeKindEvidence {
			errs = append(errs, fmt.Sprintf("node[%d] (%s) kind must be claim or evidence", i, n.ID))
		}
	}

	for i, edge := range graph.Edges {
		if edge.From == "" || edge.To == "" {
			errs = append(errs, fmt.Sprintf("edge[%d] endpoints must not be empty", i))
		}
		if edge.From == edge.To {
			errs = append(errs, fmt.Sprintf("edge[%d] self-loop is not allowed (%s -> %s)", i, edge.From, edge.To))
		}

		if _, ok := nodeIDs[edge.From]; !ok {
			errs = append(errs, fmt.Sprintf("edge[%d] source node not found: %s", i, edge.From))
		}
		if _, ok := nodeIDs[edge.To]; !ok {
			errs = append(errs, fmt.Sprintf("edge[%d] target node not found: %s", i, edge.To))
		}

		if edge.Relation != domain.RelationSupport && edge.Relation != domain.RelationAttack {
			errs = append(errs, fmt.Sprintf("edge[%d] relation must be support or attack", i))
		}
		if edge.Weight < 0 || edge.Weight > 1 {
			errs = append(errs, fmt.Sprintf("edge[%d] weight must be in [0,1]", i))
		}
	}

	errs = append(errs, validateParams(params)...)
	return errs
}

func validateParams(params domain.Params) []string {
	var errs []string

	if params.Damping < 0 || params.Damping > 1 {
		errs = append(errs, "params.damping must be in [0,1]")
	}
	if params.Epsilon <= 0 || params.Epsilon > 1 {
		errs = append(errs, "params.epsilon must be > 0 and <= 1")
	}
	if params.MaxIterations <= 0 || params.MaxIterations > 10000 {
		errs = append(errs, "params.maxIterations must be in [1, 10000]")
	}
	if params.Baseline < 0 || params.Baseline > 1 {
		errs = append(errs, "params.baseline must be in [0,1]")
	}
	return errs
}

func IsValidationError(err error) (ValidationError, bool) {
	var validationErr ValidationError
	if errors.As(err, &validationErr) {
		return validationErr, true
	}
	return ValidationError{}, false
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func round(v float64, places int) float64 {
	pow := math.Pow10(places)
	return math.Round(v*pow) / pow
}

func cmpString(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
