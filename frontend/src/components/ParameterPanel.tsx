import type { EvaluationParams } from '../types'

type ParameterPanelProps = {
  params: EvaluationParams
  onChange: (next: EvaluationParams) => void
}

function numberValue(value: string, fallback: number): number {
  const parsed = Number(value)
  if (!Number.isFinite(parsed)) {
    return fallback
  }
  return parsed
}

export function ParameterPanel({ params, onChange }: ParameterPanelProps) {
  return (
    <section className="panel">
      <h2>Evaluation Parameters</h2>
      <p className="hint">
        Tune convergence speed and sensitivity. Higher damping usually stabilizes
        dense cyclical graphs.
      </p>
      <div className="grid two">
        <label>
          damping (0–1)
          <input
            type="number"
            min={0}
            max={1}
            step={0.01}
            value={params.damping}
            onChange={(event) =>
              onChange({
                ...params,
                damping: numberValue(event.target.value, params.damping),
              })
            }
          />
        </label>
        <label>
          baseline (0–1)
          <input
            type="number"
            min={0}
            max={1}
            step={0.01}
            value={params.baseline}
            onChange={(event) =>
              onChange({
                ...params,
                baseline: numberValue(event.target.value, params.baseline),
              })
            }
          />
        </label>
        <label>
          epsilon (&gt;0)
          <input
            type="number"
            min={0.00001}
            step={0.0001}
            value={params.epsilon}
            onChange={(event) =>
              onChange({
                ...params,
                epsilon: numberValue(event.target.value, params.epsilon),
              })
            }
          />
        </label>
        <label>
          maxIterations (1–2000)
          <input
            type="number"
            min={1}
            max={2000}
            step={1}
            value={params.maxIterations}
            onChange={(event) =>
              onChange({
                ...params,
                maxIterations: numberValue(event.target.value, params.maxIterations),
              })
            }
          />
        </label>
      </div>
    </section>
  )
}
