import { describe, it, expect, beforeAll } from 'vitest'
import { getSite, getArticles, getArticle, getCategories, getTags, getAbout } from './articles'
import type { ArticleDetail } from '../types/models'
beforeAll(() => { import.meta.env.VITE_API_MODE = 'mock' })
describe('mock adapter', () => {
  it('returns site info', async () => {
    const s = await getSite(); expect(s.name).toContain('博客'); expect(s.logo_text).toBeTruthy()
    const t = await getTags(); expect(t.length).toBeGreaterThan(0)
    const ab = await getAbout(); expect(ab.skills.length).toBeGreaterThan(0)
  })
  it('returns paged articles with category/tags', async () => {
    const { list, total } = await getArticles({ page: 1, size: 6 })
    expect(total).toBeGreaterThan(0); expect(list[0]).toHaveProperty('slug'); expect(list[0]).toHaveProperty('category_name')
  })
  it('returns article detail with content', async () => {
    const { list } = await getArticles({ page: 1, size: 1 })
    const a = await getArticle(list[0].slug)
    expect((a as ArticleDetail).content).toContain('#')
  })
  it('filters by category id', async () => {
    const cats = await getCategories(); const { list } = await getArticles({ category_id: cats[0].id })
    expect(list.every(x => x.category_id === cats[0].id)).toBe(true)
  })
})
