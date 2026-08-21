/** Mock 错误与上下文（后台 mock 使用）。 */

export class MockError extends Error {
  code: number
  constructor(code: number, message: string) {
    super(message)
    this.name = 'MockError'
    this.code = code
  }
}

export function notFound(): never {
  throw new MockError(40400, '资源不存在')
}

export interface MockContext {
  path: string[]
  params: Record<string, unknown>
  data: unknown
}
