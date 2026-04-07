import type { EvaluationParams, Relation, ArgumentNode } from '../types'

type SampleLoaderProps = {
  onLoad: (payload: {
    nodes: ArgumentNode[]
    edges: Array<{
      from: string
      to: string
      relation: Relation
      weight: number
    }>
    params: EvaluationParams
  }) => void
}

const sample = {
  nodes: [
    {
      id: 'claim_core',
      label: 'Restrict private car inflow to the city center',
      kind: 'claim' as const,
    },
    {
      id: 'evidence_noise',
      label: 'Noise and PM2.5 correlate with measurable health harm',
      kind: 'evidence' as const,
    },
    {
      id: 'evidence_access',
      label: 'Logistics delays can increase small-retailer spoilage',
      kind: 'evidence' as const,
    },
    {
      id: 'claim_alt',
      label: 'Invest in public transit first, then regulate inflow',
      kind: 'claim' as const,
    },
    {
      id: 'evidence_budget',
      label: 'Fiscal slack is limited for parallel infrastructure programs',
      kind: 'evidence' as const,
    },
    {
      id: 'evidence_transition',
      label: 'Phased rollout can mitigate freight-side disruption',
      kind: 'evidence' as const,
    },
  ],
  edges: [
    { from: 'evidence_noise', to: 'claim_core', relation: 'support' as const, weight: 0.82 },
    { from: 'evidence_access', to: 'claim_core', relation: 'attack' as const, weight: 0.71 },
    { from: 'claim_alt', to: 'claim_core', relation: 'attack' as const, weight: 0.55 },
    { from: 'evidence_budget', to: 'claim_alt', relation: 'attack' as const, weight: 0.68 },
    { from: 'evidence_transition', to: 'claim_core', relation: 'support' as const, weight: 0.63 },
    { from: 'evidence_transition', to: 'evidence_access', relation: 'attack' as const, weight: 0.42 },
  ],
  params: {
    damping: 0.72,
    baseline: 0.5,
    epsilon: 0.0005,
    maxIterations: 180,
  },
}

export function SampleLoader({ onLoad }: SampleLoaderProps) {
  return (
    <section className="panel">
      <h2>Quick Start</h2>
      <p className="muted">
        Load a non-trivial civic policy debate graph with competing support and
        attack pathways.
      </p>
      <button className="button" type="button" onClick={() => onLoad(sample)}>
        Load sample debate graph
      </button>
    </section>
  )
}
