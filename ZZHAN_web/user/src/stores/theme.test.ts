import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useThemeStore } from './theme'

beforeEach(() => {
  setActivePinia(createPinia())
  localStorage.clear()
  document.documentElement.classList.remove('dark')
})

describe('theme', () => {
  it('defaults to light and toggles', () => {
    const s = useThemeStore()
    s.init()
    expect(s.dark).toBe(false)
    expect(document.documentElement.classList.contains('dark')).toBe(false)

    s.toggle()
    expect(s.dark).toBe(true)
    expect(localStorage.getItem('ct-theme')).toBe('dark')
    expect(document.documentElement.classList.contains('dark')).toBe(true)
  })

  it('restores saved light theme on init', () => {
    localStorage.setItem('ct-theme', 'light')
    const s = useThemeStore()
    s.init()
    expect(s.dark).toBe(false)
  })

  it('persists dark after toggle', () => {
    const s = useThemeStore()
    s.init()
    s.toggle()
    s.toggle()
    expect(s.dark).toBe(false)
    expect(localStorage.getItem('ct-theme')).toBe('light')
  })
})
