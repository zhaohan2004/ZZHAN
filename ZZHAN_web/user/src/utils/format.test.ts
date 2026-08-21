import { describe, it, expect } from 'vitest'
import { fmtNum, fmtDate, readTime } from './format'
describe('fmtNum', () => {
  it('formats large numbers to w/k', () => {
    expect(fmtNum(1286420)).toBe('128.6w')
    expect(fmtNum(12840)).toBe('1.3w')
    expect(fmtNum(862)).toBe('862')
  })
})
describe('readTime', () => {
  it('min 2 minutes, ~380 chars/min', () => {
    expect(readTime('')).toBe(2)
    expect(readTime('a'.repeat(760))).toBe(2)
    expect(readTime('a'.repeat(1140))).toBe(3)
  })
})
describe('fmtDate', () => {
  it('returns the input string unchanged', () => {
    expect(fmtDate('2026-08-18')).toBe('2026-08-18')
  })
})
