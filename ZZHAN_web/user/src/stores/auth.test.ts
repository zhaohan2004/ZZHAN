import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'

beforeEach(() => {
  setActivePinia(createPinia())
  localStorage.clear()
  vi.stubEnv('VITE_API_MODE', 'mock')
})

afterEach(() => {
  window.history.replaceState({}, '', '/')
  vi.unstubAllEnvs()
})

describe('auth', () => {
  it('loginWith fills token/user and logout clears', async () => {
    const s = useAuthStore()
    expect(s.isLoggedIn).toBe(false)

    window.history.replaceState({}, '', '/auth/github/callback?code=mock-code')
    const ok = await s.loginWith('github')
    expect(ok).toBe(true)
    expect(s.token).toBe('mock-access-token-github')
    expect(s.user?.provider).toBe('github')
    expect(s.user?.nickname).toBe('GitHub User')
    expect(localStorage.getItem('ct-access-token')).toBe('mock-access-token-github')
    expect(s.isLoggedIn).toBe(true)
    expect(window.location.search).toBe('')

    await s.logout()
    expect(s.token).toBeNull()
    expect(s.user).toBeNull()
    expect(localStorage.getItem('ct-access-token')).toBeNull()
    expect(s.isLoggedIn).toBe(false)
  })

  it('ensureAuth returns true when logged in, else opens login modal and returns false', async () => {
    const s = useAuthStore()
    const denied = await s.ensureAuth()
    expect(denied).toBe(false)
    expect(s.loginModalOpen).toBe(true)

    window.history.replaceState({}, '', '/auth/github/callback?code=mock-code')
    await s.loginWith('github')
    expect(s.loginModalOpen).toBe(false)
    const allowed = await s.ensureAuth()
    expect(allowed).toBe(true)
    expect(s.loginModalOpen).toBe(false)
  })
})
