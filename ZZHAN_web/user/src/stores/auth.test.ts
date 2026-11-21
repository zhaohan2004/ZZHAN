import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { setActivePinia, createPinia, disposePinia, type Pinia } from 'pinia'
import * as authApi from '@/api/auth'
import { request } from '@/api/http'
import { AccountDisabledError } from '@/api/session'
import { useAuthStore } from './auth'

const { axiosRequest, toast } = vi.hoisted(() => ({ axiosRequest: vi.fn(), toast: vi.fn() }))

vi.mock('axios', async (importOriginal) => {
  const { default: axios } = await importOriginal<typeof import('axios')>()
  return { default: { create: () => ({ request: axiosRequest }), isAxiosError: axios.isAxiosError } }
})

vi.mock('@/composables/useToast', () => ({ useToast: () => ({ toast }) }))

let pinia: Pinia

beforeEach(() => {
  pinia = createPinia()
  setActivePinia(pinia)
  localStorage.clear()
  axiosRequest.mockReset()
  toast.mockClear()
  vi.stubEnv('VITE_API_MODE', 'mock')
})

afterEach(() => {
  disposePinia(pinia)
  window.history.replaceState({}, '', '/')
  vi.unstubAllEnvs()
  vi.restoreAllMocks()
})

function rejectedResponse(status: number, code: number) {
  return { isAxiosError: true, response: { status, data: { code } } }
}

function expectSignedOut(store: ReturnType<typeof useAuthStore>): void {
  expect(store.isLoggedIn).toBe(false)
  expect(store.user).toBeNull()
  expect(store.token).toBeNull()
  expect(store.refreshToken).toBeNull()
  expect(localStorage.getItem('ct-access-token')).toBeNull()
  expect(localStorage.getItem('ct-refresh-token')).toBeNull()
}

describe('auth', () => {
  it('checks runtime configuration before navigating when the build has no client ID', async () => {
    vi.stubEnv('VITE_GITHUB_CLIENT_ID', '')
    const getConfig = vi.spyOn(authApi, 'getGitHubConfig').mockResolvedValue({ client_id: '' })

    await expect(useAuthStore().loginWith('github')).rejects.toThrow('GitHub 登录尚未配置')
    expect(getConfig).toHaveBeenCalledOnce()
    expect(window.location.pathname).toBe('/')
  })

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

  it('reports disabled GitHub login with the account status', async () => {
    vi.stubEnv('VITE_API_MODE', 'api')
    const store = useAuthStore()
    window.history.replaceState({}, '', '/auth/github/callback?code=mock-code')
    axiosRequest.mockRejectedValueOnce(rejectedResponse(403, 1004))

    await expect(store.loginWith('github')).rejects.toBeInstanceOf(AccountDisabledError)

    expectSignedOut(store)
    expect(axiosRequest).toHaveBeenCalledOnce()
    expect(toast).toHaveBeenCalledExactlyOnceWith('用户已被禁用', 'error')
  })

  it('clears a disabled session on a protected request without attempting to refresh', async () => {
    const store = useAuthStore()
    window.history.replaceState({}, '', '/auth/github/callback?code=mock-code')
    await store.loginWith('github')
    expect(store.user).not.toBeNull()
    vi.stubEnv('VITE_API_MODE', 'api')
    axiosRequest.mockRejectedValueOnce(rejectedResponse(403, 1004))

    await expect(request({ method: 'POST', url: '/articles/example/like' })).rejects.toBeInstanceOf(AccountDisabledError)

    expectSignedOut(store)
    expect(axiosRequest).toHaveBeenCalledOnce()
    expect(toast).toHaveBeenCalledExactlyOnceWith('用户已被禁用', 'error')
  })

  it('propagates a disabled auto-refresh result to all pending requests', async () => {
    vi.stubEnv('VITE_API_MODE', 'api')
    const store = useAuthStore()
    store.setTokens('old-access', 'old-refresh')
    axiosRequest
      .mockRejectedValueOnce(rejectedResponse(401, 40100))
      .mockRejectedValueOnce(rejectedResponse(401, 40100))
      .mockRejectedValueOnce(rejectedResponse(403, 1004))

    const requests = Promise.allSettled([authApi.getCurrentUser(), authApi.getCurrentUser()])
    expect(await requests).toEqual([
      { status: 'rejected', reason: expect.any(AccountDisabledError) },
      { status: 'rejected', reason: expect.any(AccountDisabledError) },
    ])

    expectSignedOut(store)
    expect(axiosRequest).toHaveBeenCalledTimes(3)
    expect(axiosRequest).toHaveBeenLastCalledWith(expect.objectContaining({ url: '/auth/refresh' }))
    expect(toast).toHaveBeenCalledExactlyOnceWith('用户已被禁用', 'error')
  })

  it('preserves a newer login when a disabled response arrives for the previous session', async () => {
    vi.stubEnv('VITE_API_MODE', 'api')
    const store = useAuthStore()
    store.setTokens('old-access', 'old-refresh')
    let rejectPreviousRequest!: (reason: unknown) => void
    axiosRequest.mockReturnValueOnce(new Promise((_resolve, reject) => {
      rejectPreviousRequest = reject
    }))

    const previousRequest = authApi.getCurrentUser()
    store.setTokens('current-access', 'current-refresh')
    rejectPreviousRequest(rejectedResponse(403, 1004))

    await expect(previousRequest).rejects.toBeInstanceOf(AccountDisabledError)
    expect(store.token).toBe('current-access')
    expect(store.refreshToken).toBe('current-refresh')
    expect(localStorage.getItem('ct-access-token')).toBe('current-access')
    expect(toast).not.toHaveBeenCalled()
  })

  it('keeps a normal expired session working after a successful refresh', async () => {
    vi.stubEnv('VITE_API_MODE', 'api')
    const store = useAuthStore()
    store.setTokens('old-access', 'current-refresh')
    axiosRequest
      .mockRejectedValueOnce(rejectedResponse(401, 40100))
      .mockResolvedValueOnce({ status: 200, data: { code: 0, data: { access_token: 'new-access' } } })
      .mockResolvedValueOnce({ status: 200, data: { code: 0, data: { id: 4 } } })

    await expect(authApi.getCurrentUser()).resolves.toEqual({ id: 4 })

    expect(axiosRequest).toHaveBeenCalledTimes(3)
    expect(axiosRequest).toHaveBeenLastCalledWith(expect.objectContaining({
      headers: { Authorization: 'Bearer new-access' },
    }))
    expect(store.token).toBe('new-access')
    expect(store.refreshToken).toBe('current-refresh')
    expect(localStorage.getItem('ct-access-token')).toBe('new-access')
    expect(toast).not.toHaveBeenCalled()
  })

  it('clears outdated tokens after an unauthorized refresh response', async () => {
    vi.stubEnv('VITE_API_MODE', 'api')
    const store = useAuthStore()
    store.setTokens('outdated-access', 'outdated-refresh')
    axiosRequest
      .mockRejectedValueOnce(rejectedResponse(401, 40100))
      .mockRejectedValueOnce(rejectedResponse(401, 40100))

    await expect(authApi.getCurrentUser()).rejects.toThrow('登录已过期，请重新登录')

    expectSignedOut(store)
    expect(axiosRequest).toHaveBeenCalledTimes(2)
    expect(toast).not.toHaveBeenCalled()
  })
})
