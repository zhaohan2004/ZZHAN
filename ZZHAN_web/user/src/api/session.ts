export class AccountDisabledError extends Error {
  constructor() {
    super('用户已被禁用')
  }
}

const ACCESS_TOKEN_KEY = 'ct-access-token'
const REFRESH_TOKEN_KEY = 'ct-refresh-token'

interface SessionTokens {
  accessToken: string | null
  refreshToken: string | null
}

type SessionChange =
  | { accessToken: string }
  | { accessToken: null; error?: AccountDisabledError }

const listeners = new Set<(change: SessionChange) => void>()

export function onSessionChange(listener: (change: SessionChange) => void): () => void {
  listeners.add(listener)
  return () => { listeners.delete(listener) }
}

export function readSession(): SessionTokens {
  return {
    accessToken: localStorage.getItem(ACCESS_TOKEN_KEY),
    refreshToken: localStorage.getItem(REFRESH_TOKEN_KEY),
  }
}

export function isCurrentSession(session: SessionTokens): boolean {
  const current = readSession()
  // 刷新访问令牌时，会话仍由同一枚刷新令牌标识。
  return current.refreshToken === session.refreshToken &&
    (session.refreshToken !== null || current.accessToken === session.accessToken)
}

export function saveSession(accessToken: string, refreshToken: string): void {
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken)
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
}

export function updateAccessToken(session: SessionTokens, accessToken: string): boolean {
  if (!isCurrentSession(session)) return false

  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken)
  for (const listener of listeners) listener({ accessToken })
  return true
}

export function clearSession(session: SessionTokens, error?: AccountDisabledError): void {
  if (!isCurrentSession(session)) return

  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
  for (const listener of listeners) listener({ accessToken: null, error })
}
