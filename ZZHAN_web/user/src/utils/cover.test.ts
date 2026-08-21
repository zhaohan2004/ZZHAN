import { describe, it, expect } from 'vitest'
import { coverArt } from './cover'

describe('coverArt', () => {
  it('returns a data:image/svg+xml data URI', () => {
    const uri = coverArt('Go 并发模型深度解析', 'Go', 1)
    expect(uri.startsWith('data:image/svg+xml')).toBe(true)
  })

  it('is deterministic for the same inputs', () => {
    const a = coverArt('Go 并发模型深度解析', 'Go', 1)
    const b = coverArt('Go 并发模型深度解析', 'Go', 1)
    expect(a).toBe(b)
  })

  it('URL-encodes the SVG (# -> %23, < -> %3C)', () => {
    const uri = coverArt('Go 并发模型深度解析', 'Go', 1)
    expect(uri).toContain('%23')
    expect(uri).toContain('%3C')
    expect(uri).not.toContain('#')
  })
})
