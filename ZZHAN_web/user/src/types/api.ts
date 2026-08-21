/**
 * API 契约类型 — 统一响应包 / 分页 / 错误码。
 */

/** 统一响应包：{ code, message, data } */
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

/** 分页响应：data.list + data.total */
export interface Paged<T> {
  list: T[]
  total: number
}

/** 错误码（docs/api.md 通用约定） */
export const ERROR_CODES = {
  OK: 0,
  BAD_REQUEST: 40001,
  UNAUTHORIZED: 40100,
  FORBIDDEN: 40300,
  NOT_FOUND: 40400,
  SERVER_ERROR: 50000,
} as const

export type ErrorCode = (typeof ERROR_CODES)[keyof typeof ERROR_CODES]
