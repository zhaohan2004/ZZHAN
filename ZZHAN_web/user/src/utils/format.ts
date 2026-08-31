export function fmtNum(n: number): string {
  if (n >= 10000) return (n / 10000).toFixed(1).replace(/\.0$/, '') + 'w'
  if (n >= 1000) return (n / 1000).toFixed(1).replace(/\.0$/, '') + 'k'
  return String(n)
}
export function fmtDate(s: string): string {
  if (!s) return ''
  // 如果已经是 YYYY-MM-DD 格式，直接返回
  if (/^\d{4}-\d{2}-\d{2}$/.test(s)) return s
  // 否则取前10位（处理 ISO8601 格式）
  return s.slice(0, 10)
}
export function readTime(content: string): number { return Math.max(2, Math.round(content.length / 380)) }
