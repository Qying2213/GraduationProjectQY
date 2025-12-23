import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usePermissionStore, PREDEFINED_ROLES } from '../permission'

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

describe('Permission Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorageMock.getItem.mockReturnValue(null)
  })

  describe('预定义角色', () => {
    it('应该包含5个预定义角色', () => {
      expect(Object.keys(PREDEFINED_ROLES)).toHaveLength(5)
      expect(PREDEFINED_ROLES).toHaveProperty('admin')
      expect(PREDEFINED_ROLES).toHaveProperty('hr_manager')
      expect(PREDEFINED_ROLES).toHaveProperty('recruiter')
      expect(PREDEFINED_ROLES).toHaveProperty('interviewer')
      expect(PREDEFINED_ROLES).toHaveProperty('viewer')
    })

    it('管理员应该拥有所有权限', () => {
      const adminRole = PREDEFINED_ROLES.admin
      expect(adminRole.permissions.length).toBeGreaterThan(20)
      expect(adminRole.permissions).toContain('dashboard:view')
      expect(adminRole.permissions).toContain('talent:delete')
      expect(adminRole.permissions).toContain('system:settings')
    })

    it('只读用户应该只有查看权限', () => {
      const viewerRole = PREDEFINED_ROLES.viewer
      expect(viewerRole.permissions.every(p => p.includes(':view'))).toBe(true)
      expect(viewerRole.permissions).not.toContain('talent:create')
      expect(viewerRole.permissions).not.toContain('job:delete')
    })
  })

  describe('权限检查', () => {
    it('设置角色后应该能正确检查权限', () => {
      const store = usePermissionStore()
      store.setRole('viewer')

      expect(store.hasPermission('dashboard:view')).toBe(true)
      expect(store.hasPermission('talent:view')).toBe(true)
      expect(store.hasPermission('talent:create')).toBe(false)
      expect(store.hasPermission('talent:delete')).toBe(false)
    })

    it('管理员应该拥有所有权限', () => {
      const store = usePermissionStore()
      store.setRole('admin')

      expect(store.hasPermission('dashboard:view')).toBe(true)
      expect(store.hasPermission('talent:delete')).toBe(true)
      expect(store.hasPermission('system:settings')).toBe(true)
    })

    it('hasAnyPermission 应该正确工作', () => {
      const store = usePermissionStore()
      store.setRole('viewer')

      expect(store.hasAnyPermission(['talent:view', 'talent:create'])).toBe(true)
      expect(store.hasAnyPermission(['talent:create', 'talent:delete'])).toBe(false)
    })

    it('hasAllPermissions 应该正确工作', () => {
      const store = usePermissionStore()
      store.setRole('viewer')

      expect(store.hasAllPermissions(['dashboard:view', 'talent:view'])).toBe(true)
      expect(store.hasAllPermissions(['talent:view', 'talent:create'])).toBe(false)
    })
  })

  describe('角色管理', () => {
    it('应该能添加新角色', () => {
      const store = usePermissionStore()
      const initialCount = store.roles.length

      const newRole = store.addRole({
        name: '测试角色',
        code: 'test_role',
        description: '测试用角色',
        permissions: ['dashboard:view']
      })

      expect(store.roles.length).toBe(initialCount + 1)
      expect(newRole.name).toBe('测试角色')
      expect(newRole.code).toBe('test_role')
    })

    it('应该能更新角色', () => {
      const store = usePermissionStore()
      const newRole = store.addRole({
        name: '测试角色',
        code: 'test_role',
        description: '测试用角色',
        permissions: ['dashboard:view']
      })

      store.updateRole(newRole.id, { name: '更新后的角色' })

      const updatedRole = store.roles.find(r => r.id === newRole.id)
      expect(updatedRole?.name).toBe('更新后的角色')
    })

    it('不应该删除预定义角色', () => {
      const store = usePermissionStore()
      const result = store.deleteRole(1) // admin 角色 id

      expect(result).toBe(false)
      expect(store.roles.find(r => r.id === 1)).toBeDefined()
    })

    it('应该能删除自定义角色', () => {
      const store = usePermissionStore()
      const newRole = store.addRole({
        name: '测试角色',
        code: 'test_role',
        description: '测试用角色',
        permissions: ['dashboard:view']
      })

      const result = store.deleteRole(newRole.id)

      expect(result).toBe(true)
      expect(store.roles.find(r => r.id === newRole.id)).toBeUndefined()
    })
  })
})
