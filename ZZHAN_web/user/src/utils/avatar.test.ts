import { describe, it, expect } from 'vitest'
import { initialsAvatar } from './avatar'

describe('initialsAvatar', () => {
  it('returns a data:image/svg+xml data URI', () => {
    const uri = initialsAvatar('阿轩', '#3b82f6', '#38bdf8')
    expect(uri.startsWith('data:image/svg+xml')).toBe(true)
  })

  it('is deterministic for the same inputs', () => {
    const a = initialsAvatar('阿轩', '#3b82f6', '#38bdf8')
    const b = initialsAvatar('阿轩', '#3b82f6', '#38bdf8')
    expect(a).toBe(b)
  })

  it('URL-encodes the SVG (# -> %23)', () => {
    const uri = initialsAvatar('阿轩', '#3b82f6', '#38bdf8')
    expect(uri).toContain('%23')
    expect(uri).not.toContain('#')
  })
})
