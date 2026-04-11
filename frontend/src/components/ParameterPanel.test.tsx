import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ParameterPanel } from './ParameterPanel'
import type { EvaluationParams } from '../types'

const baseParams: EvaluationParams = {
  damping: 0.7,
  baseline: 0.5,
  epsilon: 0.0001,
  maxIterations: 200,
}

describe('ParameterPanel', () => {
  afterEach(() => {
    cleanup()
  })

  it('rounds maxIterations decimal input to integer', () => {
    const onChange = vi.fn()
    render(<ParameterPanel params={baseParams} onChange={onChange} />)

    fireEvent.change(screen.getByLabelText(/maxIterations/i), {
      target: { value: '200.5' },
    })

    expect(onChange).toHaveBeenCalledWith({
      ...baseParams,
      maxIterations: 201,
    })
  })

  it('clamps maxIterations input to backend bounds', () => {
    const onChange = vi.fn()
    render(<ParameterPanel params={baseParams} onChange={onChange} />)
    const input = screen.getByLabelText(/maxIterations/i)

    fireEvent.change(input, {
      target: { value: '0' },
    })
    fireEvent.change(input, {
      target: { value: '999999' },
    })

    expect(onChange).toHaveBeenNthCalledWith(1, {
      ...baseParams,
      maxIterations: 1,
    })
    expect(onChange).toHaveBeenNthCalledWith(2, {
      ...baseParams,
      maxIterations: 2000,
    })
  })
})
