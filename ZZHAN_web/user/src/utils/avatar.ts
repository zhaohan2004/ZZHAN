function svgURI(svg: string): string {
  return 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(svg)
}

/** SVG 头像占位图：线性渐变 + 首字符，data URI（移植自 assets/js/data.js initialsAvatar） */
export function initialsAvatar(name: string, c1: string, c2: string, size?: number): string {
  const s = size || 160
  const ch = String(name || '?').slice(0, 1).toUpperCase()
  const svg =
    '<svg xmlns="http://www.w3.org/2000/svg" width="' + s + '" height="' + s + '">' +
    '<defs><linearGradient id="g" x1="0" y1="0" x2="1" y2="1">' +
    '<stop offset="0" stop-color="' + (c1 || '#6b7280') + '"/><stop offset="1" stop-color="' + (c2 || '#9ca3af') + '"/>' +
    '</linearGradient></defs>' +
    '<rect width="' + s + '" height="' + s + '" rx="' + Math.round(s * 0.24) + '" fill="url(#g)"/>' +
    '<text x="50%" y="50%" dy=".35em" font-family="Arial,sans-serif" font-size="' + Math.round(s * 0.42) +
    '" font-weight="700" fill="rgba(255,255,255,.92)" text-anchor="middle">' + ch + '</text></svg>'
  return svgURI(svg)
}
