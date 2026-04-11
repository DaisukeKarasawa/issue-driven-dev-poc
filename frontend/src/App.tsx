import { useRef, useState } from 'react'
import './App.css'
import { EvaluateAPIError, evaluateGraph } from './api'
import { EdgeEditor } from './components/EdgeEditor'
import { NodeEditor } from './components/NodeEditor'
import { ParameterPanel } from './components/ParameterPanel'
import { ResultPanel } from './components/ResultPanel'
import { SampleLoader } from './components/SampleLoader'
import type {
  ArgumentEdge,
  ArgumentNode,
  EvaluationParams,
  EvaluationResponse,
} from './types'

function App() {
  const [nodes, setNodes] = useState<ArgumentNode[]>([])
  const [edges, setEdges] = useState<ArgumentEdge[]>([])
  const [params, setParams] = useState<EvaluationParams>({
    damping: 0.7,
    epsilon: 0.0001,
    maxIterations: 200,
    baseline: 0.5,
  })
  const [result, setResult] = useState<EvaluationResponse | null>(null)
  const [requestError, setRequestError] = useState<string | null>(null)
  const [isEvaluating, setIsEvaluating] = useState(false)
  const nodeCounterRef = useRef(1)

  function addNode() {
    const index = nodeCounterRef.current
    nodeCounterRef.current += 1
    setNodes((prev) => [
      ...prev,
      {
        id: `node_${index}`,
        label: `New claim ${index}`,
        kind: 'claim',
      },
    ])
  }

  function updateNode(targetId: string, patch: Partial<ArgumentNode>) {
    setNodes((prevNodes) =>
      prevNodes.map((node) => (node.id === targetId ? { ...node, ...patch } : node)),
    )
  }

  function deleteNode(targetId: string) {
    setNodes((prevNodes) => prevNodes.filter((node) => node.id !== targetId))
    setEdges((prevEdges) =>
      prevEdges.filter((edge) => edge.from !== targetId && edge.to !== targetId),
    )
  }

  function addEdge(edge: ArgumentEdge) {
    setEdges((prevEdges) => [...prevEdges, edge])
  }

  function deleteEdge(edgeId: string) {
    setEdges((prevEdges) => prevEdges.filter((edge) => edge.id !== edgeId))
  }

  const canEvaluate = nodes.length > 0 && !isEvaluating

  async function handleEvaluate() {
    setRequestError(null)
    setIsEvaluating(true)
    const payload = {
      nodes,
      edges: edges.map((edge) => ({
        from: edge.from,
        to: edge.to,
        relation: edge.relation,
        weight: edge.weight,
      })),
      params,
    }

    try {
      const response = await evaluateGraph(payload)
      setResult(response)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unexpected evaluation error.'
      if (error instanceof EvaluateAPIError && error.diagnosticsErrors.length > 0) {
        setRequestError(`${message} — ${error.diagnosticsErrors.join('; ')}`)
      } else {
        setRequestError(message)
      }
      setResult(null)
    } finally {
      setIsEvaluating(false)
    }
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>Dialectic Lab</h1>
        <p>
          Build a weighted support/attack argument graph and evaluate which claims are
          resilient under criticism.
        </p>
      </header>

      <div className="toolbar-row">
        <SampleLoader
          onLoad={({ nodes: nextNodes, edges: nextEdges, params: nextParams }) => {
            setNodes(nextNodes)
            setEdges(
              nextEdges.map((edge) => ({
                ...edge,
                id: crypto.randomUUID(),
              })),
            )
            setParams(nextParams)
            setResult(null)
            setRequestError(null)
          }}
        />
      </div>

      <main className="content-grid">
        <div className="editor-grid">
          <NodeEditor
            nodes={nodes}
            onAdd={addNode}
            onUpdate={updateNode}
            onDelete={deleteNode}
          />
          <EdgeEditor
            nodes={nodes}
            edges={edges}
            onAdd={addEdge}
            onDelete={deleteEdge}
          />
          <ParameterPanel params={params} onChange={setParams} />
          <section className="panel">
            <h2>Run Evaluation</h2>
            <button
              className="button primary"
              type="button"
              onClick={handleEvaluate}
              disabled={!canEvaluate}
            >
              {isEvaluating ? 'Evaluating…' : 'Evaluate argument graph'}
            </button>
            {!nodes.length && (
              <p className="status neutral">Add at least one node to evaluate.</p>
            )}
          </section>
        </div>
        <section>
          <ResultPanel
            result={result}
            nodes={nodes}
            isLoading={isEvaluating}
            errorMessage={requestError}
          />
        </section>
      </main>
    </div>
  )
}

export default App
