import { useState } from 'react'
import type { ArgumentEdge, ArgumentNode, Relation } from '../types'

type EdgeEditorProps = {
  edges: ArgumentEdge[]
  nodes: ArgumentNode[]
  onAdd: (edge: ArgumentEdge) => void
  onDelete: (id: string) => void
}

export function EdgeEditor({ edges, nodes, onAdd, onDelete }: EdgeEditorProps) {
  const [from, setFrom] = useState('')
  const [to, setTo] = useState('')
  const [relation, setRelation] = useState<Relation>('support')
  const [weight, setWeight] = useState('0.5')
  const [error, setError] = useState('')

  const addEdge = () => {
    const parsedWeight = Number(weight)
    if (!from || !to) {
      setError('Please select both from and to nodes.')
      return
    }
    if (!Number.isFinite(parsedWeight) || parsedWeight < 0 || parsedWeight > 1) {
      setError('Weight must be a number between 0 and 1.')
      return
    }
    if (from === to) {
      setError('Self-loops are not allowed.')
      return
    }

    onAdd({
      id: crypto.randomUUID(),
      from,
      to,
      relation,
      weight: parsedWeight,
    })

    setError('')
    setWeight('0.5')
  }

  return (
    <section className="panel">
      <h2>Edges</h2>
      <div className="row">
        <label>
          from
          <select value={from} onChange={(e) => setFrom(e.target.value)}>
            <option value="">select</option>
            {nodes.map((node) => (
              <option key={node.id} value={node.id}>
                {node.label} ({node.id.slice(0, 6)})
              </option>
            ))}
          </select>
        </label>

        <label>
          relation
          <select
            value={relation}
            onChange={(e) => setRelation(e.target.value as Relation)}
          >
            <option value="support">support</option>
            <option value="attack">attack</option>
          </select>
        </label>

        <label>
          to
          <select value={to} onChange={(e) => setTo(e.target.value)}>
            <option value="">select</option>
            {nodes.map((node) => (
              <option key={node.id} value={node.id}>
                {node.label} ({node.id.slice(0, 6)})
              </option>
            ))}
          </select>
        </label>

        <label>
          weight
          <input
            type="number"
            min={0}
            max={1}
            step={0.05}
            value={weight}
            onChange={(e) => setWeight(e.target.value)}
          />
        </label>
      </div>

      <button className="primary" onClick={addEdge} type="button">
        Add edge
      </button>
      {error ? <p className="error">{error}</p> : null}

      <table>
        <thead>
          <tr>
            <th>from</th>
            <th>relation</th>
            <th>to</th>
            <th>weight</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {edges.map((edge) => (
            <tr key={edge.id}>
              <td>{edge.from.slice(0, 8)}</td>
              <td>{edge.relation}</td>
              <td>{edge.to.slice(0, 8)}</td>
              <td>{edge.weight.toFixed(2)}</td>
              <td>
                <button
                  className="danger"
                  type="button"
                  onClick={() => onDelete(edge.id)}
                >
                  delete
                </button>
              </td>
            </tr>
          ))}
          {edges.length === 0 ? (
            <tr>
              <td colSpan={5}>No edges yet.</td>
            </tr>
          ) : null}
        </tbody>
      </table>
    </section>
  )
}
