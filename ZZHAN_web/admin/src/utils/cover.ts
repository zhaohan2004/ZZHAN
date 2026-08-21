function svgURI(svg: string): string {
  return 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(svg)
}

/** 封面占位图（复刻前台 coverArt）。 */
export function coverArt(title: string, catName: string, id: number, w?: number, h?: number): string {
  const width = w || 160
  const height = h || 96
  const hue = ((id * 47 + catName.length * 13) % 360 + 360) % 360
  const c1 = 'hsl(' + hue + ',72%,58%)'
  const c2 = 'hsl(' + ((hue + 52) % 360) + ',68%,44%)'
  const letter = (catName || '?').slice(0, 2).toUpperCase()
  const svg =
    '<svg xmlns="http://www.w3.org/2000/svg" width="' + width + '" height="' + height + '" viewBox="0 0 ' + width + ' ' + height + '">' +
    '<defs><linearGradient id="g' + id + '" x1="0" y1="0" x2="1" y2="1">' +
    '<stop offset="0" stop-color="' + c1 + '"/><stop offset="1" stop-color="' + c2 + '"/></linearGradient></defs>' +
    '<rect width="' + width + '" height="' + height + '" fill="url(#g' + id + ')"/>' +
    '<text x="50%" y="50%" dy=".35em" font-family="monospace" font-size="' + Math.round(height * 0.4) +
    '" font-weight="800" fill="rgba(255,255,255,.9)" text-anchor="middle">' + letter + '</text></svg>'
  return svgURI(svg)
}

/** 头像占位图（复刻前台 initialsAvatar）。 */
export function initialsAvatar(name: string, c1: string, c2: string, size?: number): string {
  const s = size || 64
  const ch = String(name || '?').slice(0, 1).toUpperCase()
  const svg =
    '<svg xmlns="http://www.w3.org/2000/svg" width="' + s + '" height="' + s + '">' +
    '<defs><linearGradient id="g" x1="0" y1="0" x2="1" y2="1">' +
    '<stop offset="0" stop-color="' + (c1 || '#3b82f6') + '"/><stop offset="1" stop-color="' + (c2 || '#3b82f6') + '"/>' +
    '</linearGradient></defs>' +
    '<rect width="' + s + '" height="' + s + '" rx="' + Math.round(s * 0.24) + '" fill="url(#g)"/>' +
    '<text x="50%" y="50%" dy=".35em" font-family="Arial,sans-serif" font-size="' + Math.round(s * 0.42) +
    '" font-weight="700" fill="rgba(255,255,255,.92)" text-anchor="middle">' + ch + '</text></svg>'
  return svgURI(svg)
}
