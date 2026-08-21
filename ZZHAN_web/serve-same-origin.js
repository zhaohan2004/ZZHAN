/**
 * 同源预览服务器：同一个端口同时 serve 前台(user)与后台(admin)两个 dist。
 *  - /admin/*   -> ZZHAN_web/admin/dist  (后台，子路径部署)
 *  - 其他       -> ZZHAN_web/user/dist   (前台，根路径)
 * 两端都在 localhost:8080，共享同一 origin 的 localStorage，
 * 从而验证「后台改设置 -> 前台 /site 生效」。
 * 纯 Node 原生实现，无第三方依赖；history 模式路由有 fallback。
 */
const http = require('http')
const fs = require('fs')
const path = require('path')

const ROOT = __dirname
const ADMIN_DIR = path.join(ROOT, 'admin', 'dist')
const USER_DIR = path.join(ROOT, 'user', 'dist')
const PORT = Number(process.env.PORT || 8080)

const MIME = {
  '.html': 'text/html; charset=utf-8',
  '.js': 'text/javascript; charset=utf-8',
  '.mjs': 'text/javascript; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.svg': 'image/svg+xml',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.gif': 'image/gif',
  '.webp': 'image/webp',
  '.ico': 'image/x-icon',
  '.woff': 'font/woff',
  '.woff2': 'font/woff2',
  '.ttf': 'font/ttf',
  '.txt': 'text/plain; charset=utf-8',
  '.map': 'application/json',
}

function send(res, code, filePath) {
  const ext = path.extname(filePath).toLowerCase()
  res.writeHead(code, {
    'Content-Type': MIME[ext] || 'application/octet-stream',
    'Cache-Control': 'no-cache',
  })
  fs.createReadStream(filePath).pipe(res)
}

function resolveFile(baseDir, urlPath) {
  // 防止路径穿越
  const safe = path.normalize(urlPath).replace(/^(\.\.[/\\])+/, '')
  let fp = path.join(baseDir, safe)
  if (fs.existsSync(fp) && fs.statSync(fp).isDirectory()) {
    fp = path.join(fp, 'index.html')
  }
  if (!fs.existsSync(fp) || fs.statSync(fp).isDirectory()) {
    // history 模式 fallback：命中路由但无实体文件 -> index.html
    fp = path.join(baseDir, 'index.html')
  }
  return fp
}

http
  .createServer((req, res) => {
    const raw = (req.url || '/').split('?')[0]
    let urlPath = decodeURIComponent(raw)
    let baseDir = USER_DIR

    if (urlPath === '/admin' || urlPath.startsWith('/admin/')) {
      baseDir = ADMIN_DIR
      urlPath = urlPath.replace(/^\/admin/, '') || '/'
    }

    // 资源请求（带后缀）直接按文件处理，失败即 404
    const isAsset = /\.(js|mjs|css|json|svg|png|jpg|jpeg|gif|webp|ico|woff2?|ttf|map)$/i.test(urlPath)
    if (isAsset) {
      const fp = path.join(baseDir, urlPath)
      if (fs.existsSync(fp)) return send(res, 200, fp)
      return send(res, 404, path.join(ROOT, '404.txt'))
    }

    send(res, 200, resolveFile(baseDir, urlPath))
  })
  .listen(PORT, () => {
    console.log(`\n  同源预览已启动:  http://localhost:${PORT}/       (前台 user)`)
    console.log(`                   http://localhost:${PORT}/admin/  (后台 admin)`)
    console.log('  两端共享 localStorage(ct-site-settings)，后台保存设置后前台刷新即生效。\n')
  })
