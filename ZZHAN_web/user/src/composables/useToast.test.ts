import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { defineComponent, h } from 'vue'
import { mount } from '@vue/test-utils'
import { useToast } from './useToast'

const Host = defineComponent({
  setup() {
    const { toast } = useToast()
    return { toast }
  },
  render() {
    return h('div', 'host')
  },
})

describe('useToast', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    document.body.innerHTML = ''
  })
  afterEach(() => {
    vi.useRealTimers()
    document.body.innerHTML = ''
  })

  it('renders .toast after toast("hi") and removes it ~3s later', () => {
    const wrapper = mount(Host)
    wrapper.vm.toast('hi')

    const toastEl = document.querySelector('.toast')
    expect(toastEl).toBeTruthy()
    expect(toastEl?.textContent).toContain('hi')
    expect(document.querySelector('.toast-wrap')).toBeTruthy()

    vi.advanceTimersByTime(2600)
    expect(toastEl?.classList.contains('out')).toBe(true)

    vi.advanceTimersByTime(320)
    expect(document.querySelector('.toast')).toBeNull()
  })

  it('supports success/error/info types with icon', () => {
    const wrapper = mount(Host)
    wrapper.vm.toast('ok', 'success')
    const el = document.querySelector('.toast')
    expect(el?.classList.contains('success')).toBe(true)
    expect(el?.querySelector('.t-ico')).toBeTruthy()
  })

  it('reuses the same .toast-wrap for multiple toasts', () => {
    const wrapper = mount(Host)
    wrapper.vm.toast('a')
    wrapper.vm.toast('b')
    expect(document.querySelectorAll('.toast').length).toBe(2)
    expect(document.querySelectorAll('.toast-wrap').length).toBe(1)
  })
})
