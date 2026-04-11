import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { ResultPanel } from './ResultPanel'
import type { EvaluationResponse } from '../types'

describe('ResultPanel', () => {
  it('shows empty state when no result exists', () => {
    render(
      <ResultPanel
        result={null}
        nodes={[]}
        isLoading={false}
        errorMessage={null}
      />,
    )
    expect(
      screen.getByText('No result yet. Build a graph and run evaluation.'),
    ).toBeInTheDocument()
  })

  it('renders sorted scores and diagnostics', () => {
    const result: EvaluationResponse = {
      scores: [
        { nodeId: 'a', score: 0.8 },
        { nodeId: 'b', score: 0.3 },
      ],
      meta: {
        converged: true,
        iterations: 9,
        maxDelta: 0.0002,
      },
      diagnostics: {
        warnings: ['minor cycle'],
        errors: [],
      },
    }

    render(
      <ResultPanel
        result={result}
        nodes={[
          { id: 'a', label: 'Claim A', kind: 'claim' },
          { id: 'b', label: 'Claim B', kind: 'claim' },
        ]}
        isLoading={false}
        errorMessage={null}
      />,
    )

    expect(screen.getByText('Converged')).toBeInTheDocument()
    expect(screen.getByText('minor cycle')).toBeInTheDocument()
    expect(screen.getByText('Claim A')).toBeInTheDocument()
    expect(screen.getByText(/0.8000/)).toBeInTheDocument()
  })
})
