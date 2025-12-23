import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import PortalLayout from '../layout/PortalLayout.vue'
import ElementPlus from 'element-plus'

// Mock router
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/portal', component: { template: '<div>Home</div>' } },
    { path: '/portal/jobs', component: { template: '<div>Jobs</div>' } },
    { path: '/portal/login', component: { template: '<div>Login</div>' } },
  ]
})

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

describe('PortalLayout', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorageMock.getItem.mockReturnValue(null)
  })

  it('应该正确渲染导航栏', () => {
    const wrapper = mount(PortalLayout, {
      global: {
        plugins: [router, ElementPlus],
        stubs: ['router-view', 'router-link']
      }
    })

    expect(wrapper.find('.portal-header').exists()).toBe(true)
    expect(wrapper.find('.logo').exists()).toBe(true)
    expect(wrapper.find('.nav-menu').exists()).toBe(true)
  })

  it('未登录时应该显示登录注册按钮', () => {
    const wrapper = mount(PortalLayout, {
      global: {
        plugins: [router, ElementPlus],
        stubs: ['router-view', 'router-link']
      }
    })

    expect(wrapper.text()).toContain('登录')
    expect(wrapper.text()).toContain('注册')
  })

  it('应该包含页脚', () => {
    const wrapper = mount(PortalLayout, {
      global: {
        plugins: [router, ElementPlus],
        stubs: ['router-view', 'router-link']
      }
    })

    expect(wrapper.find('.portal-footer').exists()).toBe(true)
  })
})
