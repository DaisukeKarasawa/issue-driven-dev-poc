import { afterEach, describe, expect, it, vi } from 'vitest'
import { evaluateGraph } from './api'
import type { EvaluateRequest } from './types'

const payload: EvaluateRequest = {
  nodes: [{ id: 'n1', label: 'Claim', kind: 'claim' }],
  edges: [],
  params: {
    damping: 0.6,
    epsilon: 0.0001,
    maxIterations: 50,
    baseline: 0.5,
  },
}

describe('evaluateGraph', () => {
  afterEach(() => vi.unstubAllGlobals())

  it('throws API error message from backend', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 400,
        json: vi.fn().mockResolvedValue({
          message: 'invalid payload',
          diagnostics: {
            warnings: [],
            errors: ['invalid payload'],
          },
        }),
      }),
    )

    await expect(evaluateGraph(payload)).rejects.toThrow('invalid payload')
  })
})