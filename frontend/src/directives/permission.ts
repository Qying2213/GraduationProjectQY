import type { Directive, DirectiveBinding } from 'vue'
import { usePermissionStore, type Permission } from '@/store/permission'

/**
 * 权限指令
 * 用法:
 * v-permission="'talent:create'" - 单个权限
 * v-permission="['talent:create', 'talent:edit']" - 多个权限（满足任一即可）
 * v-permission:all="['talent:create', 'talent:edit']" - 多个权限（需要全部满足）
 */
export const permission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    checkPermission(el, binding)
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    checkPermission(el, binding)
  }
}

function checkPermission(el: HTMLElement, binding: DirectiveBinding) {
  const permissionStore = usePermissionStore()
  const value = binding.value as Permission | Permission[]
  const modifier = binding.arg // 'all' 表示需要所有权限

  if (!value) {
    return
  }

  let hasPermission = false

  if (Array.isArray(value)) {
    if (modifier === 'all') {
      hasPermission = permissionStore.hasAllPermissions(value)
    } else {
      hasPermission = permissionStore.hasAnyPermission(value)
    }
  } else {
    hasPermission = permissionStore.hasPermission(value)
  }

  if (!hasPermission) {
    // 移除元素
    el.parentNode?.removeChild(el)
  }
}

/**
 * 权限检查函数（用于在 setup 中使用）
 */
export function usePermission() {
  const permissionStore = usePermissionStore()

  const checkPermission = (permission: Permission | Permission[], requireAll = false): boolean => {
    if (Array.isArray(permission)) {
      return requireAll
        ? permissionStore.hasAllPermissions(permission)
        : permissionStore.hasAnyPermission(permission)
    }
    return permissionStore.hasPermission(permission)
  }

  return {
    hasPermission: permissionStore.hasPermission,
    hasAnyPermission: permissionStore.hasAnyPermission,
    hasAllPermissions: permissionStore.hasAllPermissions,
    checkPermission
  }
}

export default permission
