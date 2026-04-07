import type {
  APIErrorResponse,
  EvaluateRequest,
  EvaluationResponse,
} from './types'

export class EvaluateAPIError extends Error {
  diagnosticsErrors: string[]

  constructor(message: string, diagnosticsErrors: string[] = []) {
    super(message)
    this.name = 'EvaluateAPIError'
    this.diagnosticsErrors = diagnosticsErrors
  }
}

export async function evaluateGraph(
  payload: EvaluateRequest,
): Promise<EvaluationResponse> {
  const response = await fetch('/api/evaluate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`
    let diagnosticsErrors: string[] = []

    try {
      const body = (await response.json()) as APIErrorResponse
      if (body.message) {
        message = body.message
      }
      diagnosticsErrors = body.diagnostics?.errors ?? []
    } catch {
      // ignore parse error and keep fallback message
    }

    throw new EvaluateAPIError(message, diagnosticsErrors)
  }

  return (await response.json()) as EvaluationResponse
}
