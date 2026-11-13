import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { mount, type VueWrapper } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { nextTick } from 'vue'
import { siteData } from '@/api/mock/data'
import { useSiteStore } from '@/stores/site'
import HeroTerminal from './HeroTerminal.vue'

describe('HeroTerminal progressive reveal', () => {
  let wrapper: VueWrapper | undefined
  let motionQuery: EventTarget & { matches: boolean; media: string }

  beforeEach(() => {
    vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval'] })
    setActivePinia(createPinia())
    motionQuery = Object.assign(new EventTarget(), {
      matches: false,
      media: '(prefers-reduced-motion: reduce)',
    })
    vi.stubGlobal('matchMedia', vi.fn(() => motionQuery))
  })

  afterEach(() => {
    wrapper?.unmount()
    wrapper = undefined
    vi.useRealTimers()
    vi.unstubAllGlobals()
  })

  function mountNote(raw?: string): void {
    if (raw !== undefined) useSiteStore().site = { ...siteData, hero_terminal: raw }
    wrapper = mount(HeroTerminal)
  }

  function visibleLines(): string[] {
    return wrapper!.findAll('.note-lines > div')
      .filter((line) => !line.classes('note-line-pending'))
      .map((line) => line.text())
  }

  it('reserves every line and reveals them in order at 430 ms intervals', async () => {
    mountNote('tk|go run server.go\ncm|# ready\nfn|done')

    expect(wrapper!.findAll('.note-lines > div').map((line) => line.text()))
      .toEqual(['go run server.go', '# ready', 'done'])
    expect(visibleLines()).toEqual([])

    await vi.advanceTimersByTimeAsync(429)
    expect(visibleLines()).toEqual([])
    await vi.advanceTimersByTimeAsync(1)
    expect(visibleLines()).toEqual(['go run server.go'])
    await vi.advanceTimersByTimeAsync(430)
    expect(visibleLines()).toEqual(['go run server.go', '# ready'])
    await vi.advanceTimersByTimeAsync(430)
    expect(visibleLines()).toEqual(['go run server.go', '# ready', 'done'])
    expect(vi.getTimerCount()).toBe(0)
  })

  it('restarts from the first line when delayed site settings replace the fallback', async () => {
    mountNote()
    await vi.advanceTimersByTimeAsync(860)
    expect(visibleLines()).toHaveLength(2)

    useSiteStore().site = { ...siteData, hero_terminal: 'tk|new command\nfn|new result' }
    await nextTick()
    expect(visibleLines()).toEqual([])
    expect(wrapper!.findAll('.note-lines > div').map((line) => line.text()))
      .toEqual(['new command', 'new result'])
    expect(vi.getTimerCount()).toBe(1)

    await vi.advanceTimersByTimeAsync(430)
    expect(visibleLines()).toEqual(['new command'])
    await vi.advanceTimersByTimeAsync(430)
    expect(visibleLines()).toEqual(['new command', 'new result'])
    expect(vi.getTimerCount()).toBe(0)
  })

  it('plays newly added lines after the previous animation has finished', async () => {
    mountNote('tk|first')
    await vi.advanceTimersByTimeAsync(430)
    expect(vi.getTimerCount()).toBe(0)

    useSiteStore().site!.hero_terminal = 'tk|first\nfn|second\ncm|third'
    await nextTick()
    expect(visibleLines()).toEqual([])
    await vi.advanceTimersByTimeAsync(1290)
    expect(visibleLines()).toEqual(['first', 'second', 'third'])
    expect(vi.getTimerCount()).toBe(0)
  })

  it('keeps configured spacer rows and text styles while ignoring raw blank lines', async () => {
    mountNote('\n tk|command\n tk|\n cm|# comment\nplain text\n \n')
    await vi.advanceTimersByTimeAsync(1720)

    expect(visibleLines()).toEqual(['command', '', '# comment', 'plain text'])
    const rows = wrapper!.findAll('.note-lines > div')
    expect(rows[0]!.classes()).toContain('tk')
    expect(rows[2]!.classes()).toContain('cm')
    expect(rows[3]!.classes()).toContain('tk')
  })

  it.each([undefined, ''])('plays fallback content when settings are missing or empty: %s', async (raw) => {
    mountNote(raw)
    await vi.runAllTimersAsync()

    expect(visibleLines()[0]).toBe('$ whoami')
    expect(visibleLines().length).toBeGreaterThan(1)
    expect(vi.getTimerCount()).toBe(0)
  })

  it('does not schedule a reveal when configured content has no readable lines', () => {
    mountNote(' \n \n')

    expect(wrapper!.findAll('.note-lines > div')).toHaveLength(0)
    expect(vi.getTimerCount()).toBe(0)
  })

  it('shows all content immediately when reduced motion is enabled', async () => {
    motionQuery.matches = true
    mountNote('tk|command\nfn|result')
    await nextTick()

    expect(visibleLines()).toEqual(['command', 'result'])
    expect(vi.getTimerCount()).toBe(0)
  })

  it('finishes immediately when reduced motion is enabled during playback', async () => {
    mountNote('tk|command\nfn|result')
    await vi.advanceTimersByTimeAsync(430)
    expect(visibleLines()).toEqual(['command'])

    motionQuery.matches = true
    motionQuery.dispatchEvent(new Event('change'))
    await nextTick()

    expect(visibleLines()).toEqual(['command', 'result'])
    expect(vi.getTimerCount()).toBe(0)
  })

  it('cleans up the active timer and motion listener when leaving the homepage', () => {
    const removeListener = vi.spyOn(motionQuery, 'removeEventListener')
    mountNote('tk|command\nfn|result')
    expect(vi.getTimerCount()).toBe(1)

    wrapper!.unmount()
    wrapper = undefined

    expect(vi.getTimerCount()).toBe(0)
    expect(removeListener).toHaveBeenCalledWith('change', expect.any(Function))
    motionQuery.matches = true
    motionQuery.dispatchEvent(new Event('change'))
    expect(vi.getTimerCount()).toBe(0)
  })
})
