import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// 权限类型
export type Permission =
  | 'dashboard:view'
  | 'talent:view' | 'talent:create' | 'talent:edit' | 'talent:delete' | 'talent:export'
  | 'job:view' | 'job:create' | 'job:edit' | 'job:delete' | 'job:export'
  | 'resume:view' | 'resume:create' | 'resume:edit' | 'resume:delete' | 'resume:export'
  | 'kanban:view' | 'kanban:edit'
  | 'calendar:view' | 'calendar:create' | 'calendar:edit' | 'calendar:delete'
  | 'message:view' | 'message:send'
  | 'recommend:view' | 'recommend:use'
  | 'user:view' | 'user:create' | 'user:edit' | 'user:delete'
  | 'role:view' | 'role:create' | 'role:edit' | 'role:delete'
  | 'system:settings'

// 角色类型
export interface Role {
  id: number
  name: string
  code: string
  description: string
  permissions: Permission[]
  createdAt: string
  updatedAt: string
}

// 预定义角色
export const PREDEFINED_ROLES: Record<string, Role> = {
  admin: {
    id: 1,
    name: '超级管理员',
    code: 'admin',
    description: '拥有系统所有权限',
    permissions: [
      'dashboard:view',
      'talent:view', 'talent:create', 'talent:edit', 'talent:delete', 'talent:export',
      'job:view', 'job:create', 'job:edit', 'job:delete', 'job:export',
      'resume:view', 'resume:create', 'resume:edit', 'resume:delete', 'resume:export',
      'kanban:view', 'kanban:edit',
      'calendar:view', 'calendar:create', 'calendar:edit', 'calendar:delete',
      'message:view', 'message:send',
      'recommend:view', 'recommend:use',
      'user:view', 'user:create', 'user:edit', 'user:delete',
      'role:view', 'role:create', 'role:edit', 'role:delete',
      'system:settings'
    ],
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01'
  },
  hr_manager: {
    id: 2,
    name: 'HR主管',
    code: 'hr_manager',
    description: '负责招聘流程管理和人才库管理',
    permissions: [
      'dashboard:view',
      'talent:view', 'talent:create', 'talent:edit', 'talent:delete', 'talent:export',
      'job:view', 'job:create', 'job:edit', 'job:delete', 'job:export',
      'resume:view', 'resume:create', 'resume:edit', 'resume:delete', 'resume:export',
      'kanban:view', 'kanban:edit',
      'calendar:view', 'calendar:create', 'calendar:edit', 'calendar:delete',
      'message:view', 'message:send',
      'recommend:view', 'recommend:use',
      'user:view'
    ],
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01'
  },
  recruiter: {
    id: 3,
    name: '招聘专员',
    code: 'recruiter',
    description: '负责日常招聘工作',
    permissions: [
      'dashboard:view',
      'talent:view', 'talent:create', 'talent:edit',
      'job:view',
      'resume:view', 'resume:create', 'resume:edit',
      'kanban:view', 'kanban:edit',
      'calendar:view', 'calendar:create', 'calendar:edit',
      'message:view', 'message:send',
      'recommend:view', 'recommend:use'
    ],
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01'
  },
  interviewer: {
    id: 4,
    name: '面试官',
    code: 'interviewer',
    description: '参与面试评估',
    permissions: [
      'dashboard:view',
      'talent:view',
      'job:view',
      'resume:view',
      'kanban:view',
      'calendar:view',
      'message:view', 'message:send'
    ],
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01'
  },
  viewer: {
    id: 5,
    name: '只读用户',
    code: 'viewer',
    description: '只能查看数据',
    permissions: [
      'dashboard:view',
      'talent:view',
      'job:view',
      'resume:view',
      'kanban:view',
      'calendar:view',
      'message:view'
    ],
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01'
  }
}

// 权限分组（用于权限配置界面）
export const PERMISSION_GROUPS = [
  {
    name: '仪表板',
    permissions: [
      { code: 'dashboard:view', name: '查看仪表板' }
    ]
  },
  {
    name: '人才管理',
    permissions: [
      { code: 'talent:view', name: '查看人才' },
      { code: 'talent:create', name: '创建人才' },
      { code: 'talent:edit', name: '编辑人才' },
      { code: 'talent:delete', name: '删除人才' },
      { code: 'talent:export', name: '导出人才' }
    ]
  },
  {
    name: '职位管理',
    permissions: [
      { code: 'job:view', name: '查看职位' },
      { code: 'job:create', name: '发布职位' },
      { code: 'job:edit', name: '编辑职位' },
      { code: 'job:delete', name: '删除职位' },
      { code: 'job:export', name: '导出职位' }
    ]
  },
  {
    name: '简历管理',
    permissions: [
      { code: 'resume:view', name: '查看简历' },
      { code: 'resume:create', name: '上传简历' },
      { code: 'resume:edit', name: '编辑简历' },
      { code: 'resume:delete', name: '删除简历' },
      { code: 'resume:export', name: '导出简历' }
    ]
  },
  {
    name: '招聘看板',
    permissions: [
      { code: 'kanban:view', name: '查看看板' },
      { code: 'kanban:edit', name: '编辑看板' }
    ]
  },
  {
    name: '面试日历',
    permissions: [
      { code: 'calendar:view', name: '查看日历' },
      { code: 'calendar:create', name: '创建面试' },
      { code: 'calendar:edit', name: '编辑面试' },
      { code: 'calendar:delete', name: '删除面试' }
    ]
  },
  {
    name: '消息中心',
    permissions: [
      { code: 'message:view', name: '查看消息' },
      { code: 'message:send', name: '发送消息' }
    ]
  },
  {
    name: '智能推荐',
    permissions: [
      { code: 'recommend:view', name: '查看推荐' },
      { code: 'recommend:use', name: '使用推荐' }
    ]
  },
  {
    name: '用户管理',
    permissions: [
      { code: 'user:view', name: '查看用户' },
      { code: 'user:create', name: '创建用户' },
      { code: 'user:edit', name: '编辑用户' },
      { code: 'user:delete', name: '删除用户' }
    ]
  },
  {
    name: '角色管理',
    permissions: [
      { code: 'role:view', name: '查看角色' },
      { code: 'role:create', name: '创建角色' },
      { code: 'role:edit', name: '编辑角色' },
      { code: 'role:delete', name: '删除角色' }
    ]
  },
  {
    name: '系统设置',
    permissions: [
      { code: 'system:settings', name: '系统设置' }
    ]
  }
]

export const usePermissionStore = defineStore('permission', () => {
  // 当前用户角色
  const currentRole = ref<Role | null>(null)

  // 所有角色列表
  const roles = ref<Role[]>(Object.values(PREDEFINED_ROLES))

  // 当前用户权限列表
  const permissions = computed<Permission[]>(() => {
    return currentRole.value?.permissions || []
  })

  // 检查是否有某个权限
  const hasPermission = (permission: Permission): boolean => {
    if (!currentRole.value) return false
    // 超级管理员拥有所有权限
    if (currentRole.value.code === 'admin') return true
    return permissions.value.includes(permission)
  }

  // 检查是否有多个权限中的任意一个
  const hasAnyPermission = (perms: Permission[]): boolean => {
    return perms.some(p => hasPermission(p))
  }

  // 检查是否拥有所有指定权限
  const hasAllPermissions = (perms: Permission[]): boolean => {
    return perms.every(p => hasPermission(p))
  }

  // 设置当前角色
  const setRole = (roleCode: string) => {
    const role = roles.value.find(r => r.code === roleCode)
    if (role) {
      currentRole.value = role
      localStorage.setItem('user-role', roleCode)
    }
  }

  // 添加角色
  const addRole = (role: Omit<Role, 'id' | 'createdAt' | 'updatedAt'>) => {
    const newRole: Role = {
      ...role,
      id: Date.now(),
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }
    roles.value.push(newRole)
    return newRole
  }

  // 更新角色
  const updateRole = (id: number, updates: Partial<Role>) => {
    const index = roles.value.findIndex(r => r.id === id)
    if (index !== -1) {
      roles.value[index] = {
        ...roles.value[index],
        ...updates,
        updatedAt: new Date().toISOString()
      }
      // 如果更新的是当前角色，同步更新
      if (currentRole.value?.id === id) {
        currentRole.value = roles.value[index]
      }
    }
  }

  // 删除角色
  const deleteRole = (id: number) => {
    const index = roles.value.findIndex(r => r.id === id)
    if (index !== -1) {
      // 不允许删除预定义角色
      if (roles.value[index].id <= 5) {
        return false
      }
      roles.value.splice(index, 1)
      return true
    }
    return false
  }

  // 初始化
  const init = () => {
    const savedRole = localStorage.getItem('user-role')
    if (savedRole) {
      setRole(savedRole)
    } else {
      // 默认设置为管理员角色（演示用）
      setRole('admin')
    }
  }

  return {
    currentRole,
    roles,
    permissions,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    setRole,
    addRole,
    updateRole,
    deleteRole,
    init
  }
})
