import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUserStore } from '../user'

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

// Mock API
vi.mock('@/api/auth', () => ({
  authApi: {
    login: vi.fn().mockResolvedValue({
      data: {
        code: 0,
        data: {
          token: 'test-token',
          user: {
            id: 1,
            username: 'testuser',
            email: 'test@example.com',
            role: 'admin'
          }
        }
      }
    }),
    register: vi.fn().mockResolvedValue({ data: { code: 0 } }),
    getProfile: vi.fn().mockResolvedValue({
      data: {
        code: 0,
        data: { id: 1, username: 'testuser', email: 'test@example.com', role: 'admin' }
      }
    }),
    updateProfile: vi.fn().mockResolvedValue({
      data: {
        code: 0,
        data: { id: 1, username: 'testuser', email: 'test@example.com', role: 'admin' }
      }
    })
  }
}))

describe('User Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorageMock.getItem.mockReturnValue(null)
  })

  describe('初始状态', () => {
    it('应该有正确的初始状态', () => {
      const store = useUserStore()
      expect(store.user).toBeNull()
      expect(store.token).toBe('')
      expect(store.isLoggedIn).toBe(false)
    })
  })

  describe('登录功能', () => {
    it('登录成功应该设置用户信息和token', async () => {
      const store = useUserStore()
      await store.login('testuser', 'password123')

      expect(store.token).toBe('test-token')
      expect(store.user).toEqual({
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
        role: 'admin'
      })
      expect(store.isLoggedIn).toBe(true)
      expect(localStorageMock.setItem).toHaveBeenCalledWith('token', 'test-token')
    })
  })

  describe('登出功能', () => {
    it('登出应该清除用户信息', async () => {
      const store = useUserStore()
      await store.login('testuser', 'password123')
      store.logout()

      expect(store.user).toBeNull()
      expect(store.token).toBe('')
      expect(store.isLoggedIn).toBe(false)
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('token')
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('user')
    })
  })

  describe('角色判断', () => {
    it('应该正确判断管理员角色', async () => {
      const store = useUserStore()
      await store.login('testuser', 'password123')

      expect(store.isAdmin).toBe(true)
      expect(store.role).toBe('admin')
    })
  })
})
