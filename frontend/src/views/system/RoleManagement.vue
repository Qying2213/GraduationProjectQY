<template>
  <div class="role-management">
    <!-- Header -->
    <div class="page-header">
      <div class="header-left">
        <h1>权限管理</h1>
        <p class="subtitle">管理系统角色和权限配置</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="openCreateDialog">
          <el-icon><Plus /></el-icon>
          新建角色
        </el-button>
      </div>
    </div>

    <!-- Current Role Info -->
    <div class="current-role-card">
      <div class="role-info">
        <el-avatar :size="48" class="role-avatar">
          <el-icon><UserFilled /></el-icon>
        </el-avatar>
        <div class="role-details">
          <span class="role-label">当前角色</span>
          <h3>{{ permissionStore.currentRole?.name || '未设置' }}</h3>
          <p>{{ permissionStore.currentRole?.description }}</p>
        </div>
      </div>
      <el-select
        v-model="selectedRoleCode"
        placeholder="切换角色"
        @change="handleRoleChange"
        style="width: 200px"
      >
        <el-option
          v-for="role in permissionStore.roles"
          :key="role.code"
          :label="role.name"
          :value="role.code"
        />
      </el-select>
    </div>

    <!-- Role List -->
    <div class="role-list-card">
      <div class="card-header">
        <h3>角色列表</h3>
        <span class="count">共 {{ permissionStore.roles.length }} 个角色</span>
      </div>

      <div class="role-grid">
        <div
          v-for="role in permissionStore.roles"
          :key="role.id"
          class="role-card"
          :class="{ active: permissionStore.currentRole?.code === role.code }"
          @click="viewRoleDetail(role)"
        >
          <div class="role-card-header">
            <div class="role-icon" :class="getRoleIconClass(role.code)">
              <el-icon><component :is="getRoleIcon(role.code)" /></el-icon>
            </div>
            <el-dropdown v-if="role.id > 5" trigger="click" @click.stop>
              <el-button text :icon="MoreFilled" />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="openEditDialog(role)">
                    <el-icon><Edit /></el-icon> 编辑
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleDeleteRole(role.id)">
                    <el-icon><Delete /></el-icon> 删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-tag v-else size="small" type="info">系统角色</el-tag>
          </div>
          <div class="role-card-body">
            <h4>{{ role.name }}</h4>
            <p>{{ role.description }}</p>
          </div>
          <div class="role-card-footer">
            <span class="permission-count">
              <el-icon><Key /></el-icon>
              {{ role.permissions.length }} 项权限
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Permission Overview -->
    <div class="permission-overview-card">
      <div class="card-header">
        <h3>权限总览</h3>
        <el-tag>{{ selectedRole?.name || '请选择角色' }}</el-tag>
      </div>

      <div v-if="selectedRole" class="permission-groups">
        <div v-for="group in PERMISSION_GROUPS" :key="group.name" class="permission-group">
          <div class="group-header">
            <span class="group-name">{{ group.name }}</span>
            <span class="group-count">
              {{ getGroupPermissionCount(group) }}/{{ group.permissions.length }}
            </span>
          </div>
          <div class="group-permissions">
            <div
              v-for="perm in group.permissions"
              :key="perm.code"
              class="permission-item"
              :class="{ active: selectedRole.permissions.includes(perm.code as Permission) }"
            >
              <el-icon v-if="selectedRole.permissions.includes(perm.code as Permission)">
                <CircleCheck />
              </el-icon>
              <el-icon v-else><CircleClose /></el-icon>
              <span>{{ perm.name }}</span>
            </div>
          </div>
        </div>
      </div>

      <el-empty v-else description="请选择一个角色查看权限详情" />
    </div>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑角色' : '新建角色'"
      width="700px"
      destroy-on-close
    >
      <el-form :model="roleForm" label-width="100px" label-position="top">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="角色名称" required>
              <el-input v-model="roleForm.name" placeholder="如：HR主管" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色编码" required>
              <el-input v-model="roleForm.code" placeholder="如：hr_manager" :disabled="isEdit" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="角色描述">
          <el-input v-model="roleForm.description" type="textarea" :rows="2" placeholder="描述该角色的职责..." />
        </el-form-item>

        <el-form-item label="权限配置">
          <div class="permission-config">
            <div v-for="group in PERMISSION_GROUPS" :key="group.name" class="config-group">
              <div class="config-group-header">
                <el-checkbox
                  :model-value="isGroupAllChecked(group)"
                  :indeterminate="isGroupIndeterminate(group)"
                  @change="toggleGroupPermissions(group, $event as boolean)"
                >
                  {{ group.name }}
                </el-checkbox>
              </div>
              <div class="config-group-body">
                <el-checkbox
                  v-for="perm in group.permissions"
                  :key="perm.code"
                  :model-value="roleForm.permissions.includes(perm.code as Permission)"
                  @change="togglePermission(perm.code as Permission, $event as boolean)"
                >
                  {{ perm.name }}
                </el-checkbox>
              </div>
            </div>
          </div>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, UserFilled, MoreFilled, Edit, Delete, Key,
  CircleCheck, CircleClose, User, Setting, View, Management
} from '@element-plus/icons-vue'
import {
  usePermissionStore,
  PERMISSION_GROUPS,
  type Role,
  type Permission
} from '@/store/permission'

const permissionStore = usePermissionStore()

const dialogVisible = ref(false)
const isEdit = ref(false)
const selectedRole = ref<Role | null>(null)
const selectedRoleCode = ref('')

const roleForm = ref({
  name: '',
  code: '',
  description: '',
  permissions: [] as Permission[]
})

// Get role icon based on role code
const getRoleIcon = (code: string) => {
  const iconMap: Record<string, any> = {
    admin: Setting,
    hr_manager: Management,
    recruiter: User,
    interviewer: View,
    viewer: View
  }
  return iconMap[code] || User
}

// Get role icon class for styling
const getRoleIconClass = (code: string) => {
  const classMap: Record<string, string> = {
    admin: 'admin',
    hr_manager: 'manager',
    recruiter: 'recruiter',
    interviewer: 'interviewer',
    viewer: 'viewer'
  }
  return classMap[code] || 'default'
}

// Get permission count for a group
const getGroupPermissionCount = (group: { permissions: { code: string }[] }) => {
  if (!selectedRole.value) return 0
  return group.permissions.filter(p =>
    selectedRole.value!.permissions.includes(p.code as Permission)
  ).length
}

// Check if all permissions in a group are checked
const isGroupAllChecked = (group: { permissions: { code: string }[] }) => {
  return group.permissions.every(p =>
    roleForm.value.permissions.includes(p.code as Permission)
  )
}

// Check if group is indeterminate
const isGroupIndeterminate = (group: { permissions: { code: string }[] }) => {
  const checkedCount = group.permissions.filter(p =>
    roleForm.value.permissions.includes(p.code as Permission)
  ).length
  return checkedCount > 0 && checkedCount < group.permissions.length
}

// Toggle all permissions in a group
const toggleGroupPermissions = (group: { permissions: { code: string }[] }, checked: boolean) => {
  group.permissions.forEach(p => {
    const perm = p.code as Permission
    if (checked) {
      if (!roleForm.value.permissions.includes(perm)) {
        roleForm.value.permissions.push(perm)
      }
    } else {
      const index = roleForm.value.permissions.indexOf(perm)
      if (index > -1) {
        roleForm.value.permissions.splice(index, 1)
      }
    }
  })
}

// Toggle single permission
const togglePermission = (permission: Permission, checked: boolean) => {
  if (checked) {
    if (!roleForm.value.permissions.includes(permission)) {
      roleForm.value.permissions.push(permission)
    }
  } else {
    const index = roleForm.value.permissions.indexOf(permission)
    if (index > -1) {
      roleForm.value.permissions.splice(index, 1)
    }
  }
}

// View role detail
const viewRoleDetail = (role: Role) => {
  selectedRole.value = role
}

// Handle role change
const handleRoleChange = (code: string) => {
  permissionStore.setRole(code)
  ElMessage.success(`已切换到 ${permissionStore.currentRole?.name} 角色`)
}

// Open create dialog
const openCreateDialog = () => {
  isEdit.value = false
  roleForm.value = {
    name: '',
    code: '',
    description: '',
    permissions: ['dashboard:view']
  }
  dialogVisible.value = true
}

// Open edit dialog
const openEditDialog = (role: Role) => {
  isEdit.value = true
  roleForm.value = {
    name: role.name,
    code: role.code,
    description: role.description,
    permissions: [...role.permissions]
  }
  dialogVisible.value = true
}

// Handle save
const handleSave = () => {
  if (!roleForm.value.name || !roleForm.value.code) {
    ElMessage.warning('请填写角色名称和编码')
    return
  }

  if (isEdit.value) {
    const role = permissionStore.roles.find(r => r.code === roleForm.value.code)
    if (role) {
      permissionStore.updateRole(role.id, {
        name: roleForm.value.name,
        description: roleForm.value.description,
        permissions: roleForm.value.permissions
      })
      ElMessage.success('角色更新成功')
    }
  } else {
    // Check if code already exists
    if (permissionStore.roles.some(r => r.code === roleForm.value.code)) {
      ElMessage.error('角色编码已存在')
      return
    }
    permissionStore.addRole({
      name: roleForm.value.name,
      code: roleForm.value.code,
      description: roleForm.value.description,
      permissions: roleForm.value.permissions
    })
    ElMessage.success('角色创建成功')
  }

  dialogVisible.value = false
}

// Handle delete role
const handleDeleteRole = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个角色吗？', '确认删除', {
      type: 'warning'
    })
    const success = permissionStore.deleteRole(id)
    if (success) {
      ElMessage.success('删除成功')
    } else {
      ElMessage.error('无法删除系统预设角色')
    }
  } catch {
    // User cancelled
  }
}

onMounted(() => {
  permissionStore.init()
  selectedRoleCode.value = permissionStore.currentRole?.code || 'admin'
  if (permissionStore.currentRole) {
    selectedRole.value = permissionStore.currentRole
  }
})
</script>

<style scoped lang="scss">
.role-management {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--card-bg);
  padding: 24px;
  border-radius: 16px;
  box-shadow: var(--card-shadow);

  h1 {
    margin: 0;
    font-size: 24px;
    color: var(--text-primary);
  }

  .subtitle {
    margin: 4px 0 0;
    color: var(--text-secondary);
    font-size: 14px;
  }
}

.current-role-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
  border-radius: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;

  .role-info {
    display: flex;
    align-items: center;
    gap: 16px;

    .role-avatar {
      background: rgba(255, 255, 255, 0.2);

      .el-icon {
        font-size: 24px;
        color: white;
      }
    }

    .role-details {
      color: white;

      .role-label {
        font-size: 12px;
        opacity: 0.8;
      }

      h3 {
        margin: 4px 0;
        font-size: 20px;
      }

      p {
        margin: 0;
        font-size: 14px;
        opacity: 0.9;
      }
    }
  }

  :deep(.el-select) {
    .el-input__wrapper {
      background: rgba(255, 255, 255, 0.15);
      border: 1px solid rgba(255, 255, 255, 0.3);
      box-shadow: none;

      .el-input__inner {
        color: white;

        &::placeholder {
          color: rgba(255, 255, 255, 0.7);
        }
      }

      .el-input__suffix {
        color: white;
      }
    }
  }
}

.role-list-card,
.permission-overview-card {
  background: var(--card-bg);
  padding: 24px;
  border-radius: 16px;
  box-shadow: var(--card-shadow);

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      margin: 0;
      font-size: 18px;
      color: var(--text-primary);
    }

    .count {
      font-size: 14px;
      color: var(--text-secondary);
    }
  }
}

.role-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;

  .role-card {
    background: var(--bg-secondary);
    border: 2px solid transparent;
    border-radius: 12px;
    padding: 20px;
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      border-color: var(--primary-color);
      transform: translateY(-2px);
    }

    &.active {
      border-color: var(--primary-color);
      background: rgba(99, 102, 241, 0.05);
    }

    .role-card-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 12px;

      .role-icon {
        width: 40px;
        height: 40px;
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;

        .el-icon {
          font-size: 20px;
          color: white;
        }

        &.admin { background: linear-gradient(135deg, #667eea, #764ba2); }
        &.manager { background: linear-gradient(135deg, #f093fb, #f5576c); }
        &.recruiter { background: linear-gradient(135deg, #4facfe, #00f2fe); }
        &.interviewer { background: linear-gradient(135deg, #43e97b, #38f9d7); }
        &.viewer { background: #9ca3af; }
        &.default { background: #6b7280; }
      }
    }

    .role-card-body {
      h4 {
        margin: 0 0 4px;
        font-size: 16px;
        color: var(--text-primary);
      }

      p {
        margin: 0;
        font-size: 13px;
        color: var(--text-secondary);
        line-height: 1.5;
      }
    }

    .role-card-footer {
      margin-top: 12px;
      padding-top: 12px;
      border-top: 1px solid var(--border-color);

      .permission-count {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        color: var(--text-secondary);
      }
    }
  }
}

.permission-groups {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;

  .permission-group {
    background: var(--bg-secondary);
    border-radius: 12px;
    padding: 16px;

    .group-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      padding-bottom: 12px;
      border-bottom: 1px solid var(--border-color);

      .group-name {
        font-weight: 600;
        color: var(--text-primary);
      }

      .group-count {
        font-size: 12px;
        color: var(--text-secondary);
      }
    }

    .group-permissions {
      display: flex;
      flex-direction: column;
      gap: 8px;

      .permission-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px 12px;
        border-radius: 8px;
        font-size: 14px;
        color: var(--text-secondary);
        background: var(--card-bg);

        &.active {
          color: var(--success-color);
          background: rgba(16, 185, 129, 0.1);

          .el-icon {
            color: var(--success-color);
          }
        }

        .el-icon {
          font-size: 16px;
        }
      }
    }
  }
}

.permission-config {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  max-height: 400px;
  overflow-y: auto;
  padding: 4px;

  .config-group {
    background: var(--bg-secondary);
    border-radius: 12px;
    padding: 16px;

    .config-group-header {
      margin-bottom: 12px;
      padding-bottom: 12px;
      border-bottom: 1px solid var(--border-color);

      :deep(.el-checkbox__label) {
        font-weight: 600;
        color: var(--text-primary);
      }
    }

    .config-group-body {
      display: flex;
      flex-direction: column;
      gap: 8px;

      :deep(.el-checkbox) {
        margin-right: 0;
        height: auto;
        padding: 8px;
        border-radius: 8px;

        &:hover {
          background: var(--bg-tertiary);
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .current-role-card {
    flex-direction: column;
    gap: 16px;
    text-align: center;

    .role-info {
      flex-direction: column;
    }
  }

  .role-grid {
    grid-template-columns: 1fr;
  }

  .permission-groups {
    grid-template-columns: 1fr;
  }
}
</style>
