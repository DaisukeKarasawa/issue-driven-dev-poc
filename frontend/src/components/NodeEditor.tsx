import type { ArgumentNode } from '../types'

interface NodeEditorProps {
  nodes: ArgumentNode[]
  onAdd: () => void
  onUpdate: (id: string, patch: Partial<ArgumentNode>) => void
  onDelete: (id: string) => void
}

export function NodeEditor({ nodes, onAdd, onUpdate, onDelete }: NodeEditorProps) {
  return (
    <section className="panel">
      <div className="panel-header">
        <h2>Nodes</h2>
        <button type="button" onClick={onAdd}>
          Add Node
        </button>
      </div>
      {nodes.length === 0 ? (
        <p className="empty">No nodes yet. Add claims or evidence.</p>
      ) : (
        <ul className="list">
          {nodes.map((node) => (
            <li key={node.id} className="list-item">
              <div className="field-row">
                <label htmlFor={`node-id-${node.id}`}>ID</label>
                <input
                  id={`node-id-${node.id}`}
                  value={node.id}
                  readOnly
                  placeholder="unique-node-id"
                />
              </div>
              <div className="field-row">
                <label htmlFor={`node-label-${node.id}`}>Label</label>
                <input
                  id={`node-label-${node.id}`}
                  value={node.label}
                  onChange={(event) => onUpdate(node.id, { label: event.target.value })}
                  placeholder="Describe this claim"
                />
              </div>
              <div className="field-row">
                <label htmlFor={`node-kind-${node.id}`}>Kind</label>
                <select
                  id={`node-kind-${node.id}`}
                  value={node.kind}
                  onChange={(event) =>
                    onUpdate(node.id, {
                      kind: event.target.value as ArgumentNode['kind'],
                    })
                  }
                >
                  <option value="claim">claim</option>
                  <option value="evidence">evidence</option>
                </select>
              </div>
              <button type="button" className="danger" onClick={() => onDelete(node.id)}>
                Remove
              </button>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}