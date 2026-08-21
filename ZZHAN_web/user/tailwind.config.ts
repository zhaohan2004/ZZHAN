import type { Config } from 'tailwindcss'
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: { extend: {
    fontFamily: {
      sans: ['Inter', 'PingFang SC', 'Microsoft YaHei', 'system-ui', 'sans-serif'],
      mono: ['JetBrains Mono', 'monospace'],
    },
    colors: {
      base: 'var(--bg)', 'base-soft': 'var(--bg-soft)',
      card: 'var(--card)', 'card-2': 'var(--card-2)',
      text: { DEFAULT: 'var(--text)', 2: 'var(--text-2)', 3: 'var(--text-3)' },
      accent: { DEFAULT: 'var(--accent)', 2: 'var(--accent-2)', 3: 'var(--accent-3)' },
      danger: 'var(--danger)', success: 'var(--success)', warning: 'var(--warning)',
    },
    boxShadow: {
      card: 'var(--shadow)', lg: 'var(--shadow-lg)', glow: 'var(--shadow-glow)',
    },
  }},
} satisfies Config
