/** 图形验证码生成 — 字符集 + canvas 绘制（复刻后台 genCaptcha）。 */

const CHARS = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
const COLORS = ['#a5b4fc', '#67e8f9', '#f0abfc', '#86efac']

export interface CaptchaResult {
  text: string
  draw(ctx: CanvasRenderingContext2D, w: number, h: number): void
}

export function genCaptcha(): CaptchaResult {
  let text = ''
  for (let i = 0; i < 4; i++) text += CHARS[Math.floor(Math.random() * CHARS.length)]
  return {
    text,
    draw(ctx, w, h) {
      const g = ctx.createLinearGradient(0, 0, w, h)
      g.addColorStop(0, '#1e1b4b')
      g.addColorStop(1, '#164e63')
      ctx.fillStyle = g
      ctx.fillRect(0, 0, w, h)
      for (let i = 0; i < 4; i++) {
        ctx.strokeStyle = 'rgba(255,255,255,' + (0.12 + Math.random() * 0.15) + ')'
        ctx.beginPath()
        ctx.moveTo(Math.random() * w, Math.random() * h)
        ctx.lineTo(Math.random() * w, Math.random() * h)
        ctx.stroke()
      }
      for (let i = 0; i < 4; i++) {
        const ch = text[i]
        ctx.save()
        ctx.translate(16 + i * 26 + Math.random() * 6, 28)
        ctx.rotate((Math.random() - 0.5) * 0.5)
        ctx.font = '700 20px "JetBrains Mono", monospace'
        ctx.fillStyle = COLORS[i]
        ctx.fillText(ch, 0, 0)
        ctx.restore()
      }
      for (let i = 0; i < 30; i++) {
        ctx.fillStyle = 'rgba(255,255,255,' + Math.random() * 0.4 + ')'
        ctx.fillRect(Math.random() * w, Math.random() * h, 1.4, 1.4)
      }
    },
  }
}
