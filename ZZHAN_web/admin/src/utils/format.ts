/** 数字 / 日期 / 阅读时长格式化（复刻前台行为）。 */
export function fmtNum(n: number): string {
  if (n >= 10000) return (n / 10000).toFixed(1).replace(/\.0$/, '') + 'w'
  if (n >= 1000) return (n / 1000).toFixed(1).replace(/\.0$/, '') + 'k'
  return String(n)
}
export function fmtDate(s: string): string {
  return s
}
export function readTime(content: string): number {
  return Math.max(2, Math.round(content.length / 380))
}
