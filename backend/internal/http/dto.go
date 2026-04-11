package httpapi

import "dialecticlab/backend/internal/domain"

type EvaluateRequest struct {
	Nodes  []domain.Node `json:"nodes"`
	Edges  []domain.Edge `json:"edges"`
	Params ParamsInput   `json:"params"`
}

type ParamsInput struct {
	Damping       *float64 `json:"damping"`
	Epsilon       *float64 `json:"epsilon"`
	MaxIterations *int     `json:"maxIterations"`
	Baseline      *float64 `json:"baseline"`
}

func (in ParamsInput) ToDomain() domain.Params {
	params := domain.Params{
		Damping:       0.35,
		Epsilon:       0.000001,
		MaxIterations: 250,
		Baseline:      0.5,
	}

	if in.Damping != nil {
		params.Damping = *in.Damping
	}
	if in.Epsilon != nil {
		params.Epsilon = *in.Epsilon
	}
	if in.MaxIterations != nil {
		params.MaxIterations = *in.MaxIterations
	}
	if in.Baseline != nil {
		params.Baseline = *in.Baseline
	}
	return params
}

type APIErrorResponse struct {
	Message     string `json:"message"`
	Diagnostics struct {
		Warnings []string `json:"warnings"`
		Errors   []string `json:"errors"`
	} `json:"diagnostics"`
}

type HealthResponse struct {
	Status string `json:"status"`
}
