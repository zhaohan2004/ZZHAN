/** 未上传封面时，用技术分类和文章编号生成一张索引卡。 */
export function coverArt(
  title: string,
  catName: string,
  id: number,
  w = 800,
  h = 420,
  noText = false,
): string {
  const width = w || 800
  const height = h || 420
  const unit = Math.min(width, height)
  const colors = ['#dceff4', '#e9e8f4', '#f5efd9']
  const index = Math.abs(id + catName.length) % colors.length
  const categoryName = (catName || title || '?').slice(0, 6).toUpperCase()
  const category = categoryName.replace(
    /[&<>]/g,
    (char) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;' })[char] ?? char,
  )
  const label = noText
    ? ''
    : `<g text-anchor="middle" font-family="monospace" fill="#294653">
        <text y="${unit * 0.015}" font-size="${Math.min(unit * 0.17, unit * 0.74 / categoryName.length)}" font-weight="700">${category}</text>
        <text y="${unit * 0.19}" font-size="${unit * 0.07}" letter-spacing="${unit * 0.008}"># ${String(id).padStart(2, '0')}</text>
      </g>`
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
    <rect width="${width}" height="${height}" fill="${colors[index]}"/>
    <g transform="translate(${width / 2} ${height / 2})">
      <rect x="${-unit * 0.43}" y="${-unit * 0.33}" width="${unit * 0.86}" height="${unit * 0.66}"
        rx="${unit * 0.02}" fill="none" stroke="#294653" stroke-opacity=".24"/>
      <path d="M ${-unit * 0.1} ${-unit * 0.22} h ${unit * 0.2}" stroke="#206079" stroke-width="${unit * 0.009}"/>
      ${label}
    </g>
  </svg>`
  return 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(svg)
}
