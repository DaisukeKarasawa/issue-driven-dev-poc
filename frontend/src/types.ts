export type NodeKind = 'claim' | 'evidence'
export type Relation = 'support' | 'attack'

export interface ArgumentNode {
  id: string
  label: string
  kind: NodeKind
}

export interface ArgumentEdge {
  id: string
  from: string
  to: string
  relation: Relation
  weight: number
}

export interface EvaluationParams {
  damping: number
  epsilon: number
  maxIterations: number
  baseline: number
}

export interface EvaluateRequest {
  nodes: ArgumentNode[]
  edges: Array<Omit<ArgumentEdge, 'id'>>
  params: EvaluationParams
}

export interface NodeScore {
  nodeId: string
  score: number
}

export interface EvaluationResponse {
  scores: NodeScore[]
  meta: {
    converged: boolean
    iterations: number
    maxDelta: number
  }
  diagnostics: {
    warnings: string[]
    errors: string[]
  }
}

export interface APIErrorResponse {
  message: string
  diagnostics?: {
    warnings?: string[]
    errors?: string[]
  }
}
