import type { ArgumentNode, EvaluationResponse } from '../types'

type ResultPanelProps = {
  result: EvaluationResponse | null
  nodes: ArgumentNode[]
  isLoading: boolean
  errorMessage: string | null
}

export function ResultPanel({
  result,
  nodes,
  isLoading,
  errorMessage,
}: ResultPanelProps) {
  const labelMap = new Map(nodes.map((node) => [node.id, node.label]))

  return (
    <section className="panel">
      <h2>Evaluation Result</h2>

      {isLoading && (
        <p className="status loading" role="status">
          Evaluating graph...
        </p>
      )}

      {errorMessage && (
        <p className="status error" role="alert">
          {errorMessage}
        </p>
      )}

      {!isLoading && !errorMessage && !result && (
        <p className="status neutral">
          No result yet. Build a graph and run evaluation.
        </p>
      )}

      {!isLoading && !errorMessage && result && (
        <>
          <div className="meta-grid">
            <div>
              <span className="meta-label">Converged</span>
              <span className="meta-value">
                {result.meta.converged ? 'Yes' : 'No'}
              </span>
            </div>
            <div>
              <span className="meta-label">Iterations</span>
              <span className="meta-value">{result.meta.iterations}</span>
            </div>
            <div>
              <span className="meta-label">Max delta</span>
              <span className="meta-value">
                {result.meta.maxDelta.toFixed(6)}
              </span>
            </div>
          </div>

          <h3>Claim Ranking</h3>
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Rank</th>
                  <th>Node</th>
                  <th>Score</th>
                </tr>
              </thead>
              <tbody>
                {result.scores.map((entry, index) => (
                  <tr key={entry.nodeId}>
                    <td>{index + 1}</td>
                    <td>{labelMap.get(entry.nodeId) ?? entry.nodeId}</td>
                    <td>{entry.score.toFixed(4)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className="diagnostics-grid">
            <div>
              <h3>Warnings</h3>
              {result.diagnostics.warnings.length === 0 ? (
                <p className="status neutral">No warnings.</p>
              ) : (
                <ul>
                  {result.diagnostics.warnings.map((warning, idx) => (
                    <li key={`${warning}-${idx}`}>{warning}</li>
                  ))}
                </ul>
              )}
            </div>
            <div>
              <h3>Errors</h3>
              {result.diagnostics.errors.length === 0 ? (
                <p className="status neutral">No errors.</p>
              ) : (
                <ul>
                  {result.diagnostics.errors.map((err, idx) => (
                    <li key={`${err}-${idx}`}>{err}</li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        </>
      )}
    </section>
  )
}
