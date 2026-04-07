package domain

type Relation string

const (
	RelationSupport Relation = "support"
	RelationAttack  Relation = "attack"
)

type NodeKind string

const (
	NodeKindClaim    NodeKind = "claim"
	NodeKindEvidence NodeKind = "evidence"
)

type Node struct {
	ID    string   `json:"id"`
	Label string   `json:"label"`
	Kind  NodeKind `json:"kind"`
}

type Edge struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Relation Relation `json:"relation"`
	Weight   float64  `json:"weight"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Params struct {
	Damping       float64 `json:"damping"`
	Epsilon       float64 `json:"epsilon"`
	MaxIterations int     `json:"maxIterations"`
	Baseline      float64 `json:"baseline"`
}
