/** API 契约类型 — 统一响应包 / 分页 / 错误码。 */

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface Paged<T> {
  list: T[]
  total: number
}

export const ERROR_CODES = {
  OK: 0,
  BAD_REQUEST: 40001,
  UNAUTHORIZED: 40100,
  FORBIDDEN: 40300,
  NOT_FOUND: 40400,
  SERVER_ERROR: 50000,
} as const

export type ErrorCode = (typeof ERROR_CODES)[keyof typeof ERROR_CODES]
