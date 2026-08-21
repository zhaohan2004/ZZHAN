function svgURI(svg: string): string {
  return 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(svg)
}

/** SVG 封面占位图：确定性 hue 渐变 + grid 网格 + 径向光晕 + 等宽首字母，data URI（移植自 assets/js/data.js coverArt） */
export function coverArt(
  title: string,
  catName: string,
  id: number,
  w?: number,
  h?: number,
  noText?: boolean,
): string {
  const width = w || 800
  const height = h || 420
  const hue = ((id * 47 + catName.length * 13) % 360 + 360) % 360
  const c1 = 'hsl(' + hue + ',72%,58%)'
  const c2 = 'hsl(' + ((hue + 52) % 360) + ',68%,44%)'
  const letter = (catName || '?').slice(0, 2).toUpperCase()
  const textLayer = noText
    ? ''
    : '<text x="' + Math.round(width * 0.06) + '" y="' + Math.round(height * 0.52) +
      '" font-family="monospace" font-size="' + Math.round(height * 0.17) +
      '" font-weight="800" fill="rgba(255,255,255,.95)">' + letter + '</text>' +
      '<text x="' + Math.round(width * 0.06) + '" y="' + Math.round(height * 0.7) +
      '" font-family="monospace" font-size="' + Math.round(height * 0.085) +
      '" fill="rgba(255,255,255,.72)"># ' + String(id).padStart(2, '0') + ' · ' + (title || '').slice(0, 10) + '</text>'
  const svg =
    '<svg xmlns="http://www.w3.org/2000/svg" width="' + width + '" height="' + height + '" viewBox="0 0 ' + width + ' ' + height + '">' +
    '<defs>' +
    '<linearGradient id="cg' + id + '" x1="0" y1="0" x2="1" y2="1">' +
    '<stop offset="0" stop-color="' + c1 + '"/><stop offset="1" stop-color="' + c2 + '"/></linearGradient>' +
    '<pattern id="gp' + id + '" width="42" height="42" patternUnits="userSpaceOnUse">' +
    '<path d="M42 0H0V42" fill="none" stroke="rgba(255,255,255,.09)" stroke-width="1"/></pattern>' +
    '<radialGradient id="rg' + id + '" cx=".28" cy=".18" r=".95">' +
    '<stop offset="0" stop-color="rgba(255,255,255,.28)"/><stop offset="1" stop-color="rgba(255,255,255,0)"/></radialGradient>' +
    '</defs>' +
    '<rect width="' + width + '" height="' + height + '" fill="url(#cg' + id + ')"/>' +
    '<rect width="' + width + '" height="' + height + '" fill="url(#gp' + id + ')"/>' +
    '<circle cx="' + Math.round(width * 0.84) + '" cy="' + Math.round(height * 0.22) + '" r="' + Math.round(height * 0.4) + '" fill="url(#rg' + id + ')"/>' +
    '<circle cx="' + Math.round(width * 0.1) + '" cy="' + Math.round(height * 0.88) + '" r="' + Math.round(height * 0.3) + '" fill="none" stroke="rgba(255,255,255,.2)" stroke-width="2"/>' +
    textLayer +
    '</svg>'
  return svgURI(svg)
}
