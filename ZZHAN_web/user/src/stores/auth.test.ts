import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'

beforeEach(() => {
  setActivePinia(createPinia())
  localStorage.clear()
  import.meta.env.VITE_API_MODE = 'mock'
})

describe('auth', () => {
  it('loginWith("github","code") fills token/user and logout clears', async () => {
    const s = useAuthStore()
    expect(s.isLoggedIn).toBe(false)

    const ok = await s.loginWith('github', 'code')
    expect(ok).toBe(true)
    expect(s.token).toBe('mock-access-token-github')
    expect(s.user?.provider).toBe('github')
    expect(s.user?.nickname).toBe('GitHub User')
    expect(s.need_profile).toBe(false)
    expect(localStorage.getItem('ct-access-token')).toBe('mock-access-token-github')
    expect(s.isLoggedIn).toBe(true)

    await s.logout()
    expect(s.token).toBeNull()
    expect(s.user).toBeNull()
    expect(s.need_profile).toBe(false)
    expect(localStorage.getItem('ct-access-token')).toBeNull()
    expect(s.isLoggedIn).toBe(false)
  })

  it('loginWith("wechat","code") works too', async () => {
    const s = useAuthStore()
    const ok = await s.loginWith('wechat', 'code')
    expect(ok).toBe(true)
    expect(s.token).toBe('mock-access-token-wechat')
    expect(s.user?.provider).toBe('wechat')
  })

  it('ensureAuth returns true when logged in, else opens login modal and returns false', async () => {
    const s = useAuthStore()
    const denied = await s.ensureAuth()
    expect(denied).toBe(false)
    expect(s.loginModalOpen).toBe(true)

    await s.loginWith('github', 'code')
    expect(s.loginModalOpen).toBe(false)
    const allowed = await s.ensureAuth()
    expect(allowed).toBe(true)
    expect(s.loginModalOpen).toBe(false)
  })
})
