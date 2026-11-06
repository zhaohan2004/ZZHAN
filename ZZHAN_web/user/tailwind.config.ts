import type { Config } from 'tailwindcss'
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: { extend: {
    fontFamily: {
      sans: ['var(--font-body)'],
      mono: ['var(--font-mono)'],
    },
    colors: {
      base: 'var(--bg)', 'base-soft': 'var(--bg-soft)',
      card: 'var(--card)', 'card-2': 'var(--card-2)',
      text: { DEFAULT: 'var(--text)', 2: 'var(--text-2)', 3: 'var(--text-3)' },
      accent: { DEFAULT: 'var(--accent)' },
      danger: 'var(--danger)', success: 'var(--success)', warning: 'var(--warning)',
    },
    boxShadow: {
      card: 'var(--shadow)', lg: 'var(--shadow-lg)', glow: 'var(--shadow-glow)',
    },
  }},
} satisfies Config
