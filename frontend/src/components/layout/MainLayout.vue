<template>
  <div class="main-layout">
    <el-container>
      <!-- Sidebar -->
      <el-aside :width="isCollapse ? '64px' : '240px'" class="sidebar">
        <div class="logo">
          <div class="logo-icon">
            <el-icon><MagicStick /></el-icon>
          </div>
          <h2 v-show="!isCollapse" class="logo-text">人才运营平台</h2>
        </div>

        <el-menu
          :default-active="$route.path"
          router
          class="sidebar-menu"
          :collapse="isCollapse"
        >
          <el-menu-item index="/dashboard">
            <el-icon><DataAnalysis /></el-icon>
            <span>仪表板</span>
          </el-menu-item>

          <el-menu-item index="/talents">
            <el-icon><User /></el-icon>
            <span>人才管理</span>
          </el-menu-item>

          <el-menu-item index="/jobs">
            <el-icon><Suitcase /></el-icon>
            <span>职位管理</span>
          </el-menu-item>

          <el-menu-item index="/resumes">
            <el-icon><Document /></el-icon>
            <span>简历管理</span>
          </el-menu-item>

          <el-menu-item index="/recommend">
            <el-icon><TrendCharts /></el-icon>
            <span>智能推荐</span>
          </el-menu-item>

          <el-menu-item index="/kanban">
            <el-icon><Operation /></el-icon>
            <span>招聘看板</span>
          </el-menu-item>

          <el-menu-item index="/calendar">
            <el-icon><Calendar /></el-icon>
            <span>面试日历</span>
          </el-menu-item>

          <el-menu-item index="/messages">
            <el-icon><ChatDotRound /></el-icon>
            <span>消息中心</span>
            <el-badge v-if="unreadCount > 0" :value="unreadCount" class="badge" />
          </el-menu-item>

          <el-menu-item index="/reports">
            <el-icon><DataLine /></el-icon>
            <span>数据报表</span>
          </el-menu-item>

          <el-sub-menu index="system">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>系统管理</span>
            </template>
            <el-menu-item index="/roles">
              <el-icon><Key /></el-icon>
              <span>权限管理</span>
            </el-menu-item>
            <el-menu-item index="/logs">
              <el-icon><Document /></el-icon>
              <span>操作日志</span>
            </el-menu-item>
            <el-menu-item index="/settings">
              <el-icon><Setting /></el-icon>
              <span>系统设置</span>
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-aside>

      <!-- Main Content -->
      <el-container>
        <el-header class="header">
          <div class="header-left">
            <el-icon @click="toggleSidebar" class="collapse-icon">
              <Fold v-if="!isCollapse" />
              <Expand v-else />
            </el-icon>
          </div>

          <div class="header-right">
            <!-- 快捷入口按钮 -->
            <el-tooltip content="前台投递简历" placement="bottom">
              <el-button type="primary" plain size="small" @click="goToPortal" class="quick-btn">
                <el-icon><Upload /></el-icon>
                投递简历
              </el-button>
            </el-tooltip>
            
            <el-tooltip content="AI智能评估系统" placement="bottom">
              <el-button type="success" plain size="small" @click="goToEvaluator" class="quick-btn">
                <el-icon><MagicStick /></el-icon>
                AI评估
              </el-button>
            </el-tooltip>

            <!-- Theme Switcher -->
            <el-dropdown @command="handleThemeChange" trigger="click">
              <div class="theme-btn">
                <el-icon v-if="themeStore.actualTheme === 'light'"><Sunny /></el-icon>
                <el-icon v-else><Moon /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="light" :class="{ active: themeStore.mode === 'light' }">
                    <el-icon><Sunny /></el-icon>
                    浅色模式
                  </el-dropdown-item>
                  <el-dropdown-item command="dark" :class="{ active: themeStore.mode === 'dark' }">
                    <el-icon><Moon /></el-icon>
                    深色模式
                  </el-dropdown-item>
                  <el-dropdown-item command="system" :class="{ active: themeStore.mode === 'system' }">
                    <el-icon><Monitor /></el-icon>
                    跟随系统
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>

            <!-- Notifications -->
            <el-badge :value="unreadCount" :hidden="unreadCount === 0">
              <el-icon class="icon-btn" @click="goToMessages">
                <Bell />
              </el-icon>
            </el-badge>

            <!-- User Menu -->
            <el-dropdown @command="handleCommand">
              <div class="user-info">
                <el-avatar :src="userStore.user?.avatar" class="avatar">
                  {{ userStore.user?.username?.charAt(0).toUpperCase() }}
                </el-avatar>
                <span class="username">{{ userStore.user?.username }}</span>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">
                    <el-icon><User /></el-icon>
                    个人中心
                  </el-dropdown-item>
                  <el-dropdown-item divided command="logout">
                    <el-icon><SwitchButton /></el-icon>
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <el-main class="main-content">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useThemeStore, type ThemeMode } from '@/store/theme'
import { usePermissionStore } from '@/store/permission'
import { messageApi } from '@/api/message'
import { ElMessage } from 'element-plus'
import {
  DataAnalysis, User, Suitcase, Document, TrendCharts,
  ChatDotRound, Fold, Expand, Bell, SwitchButton, MagicStick,
  Sunny, Moon, Monitor, Operation, Calendar, Key, Setting, DataLine, Upload
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const themeStore = useThemeStore()
const permissionStore = usePermissionStore()
const isCollapse = ref(false)
const unreadCount = ref(0)

// 初始化主题和权限
onMounted(() => {
  themeStore.init()
  permissionStore.init()
})

const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

const goToMessages = () => {
  router.push('/messages')
}

const goToPortal = () => {
  // 跳转到前台投递简历页面
  window.open('/portal/jobs', '_blank')
}

const goToEvaluator = () => {
  // 跳转到 AI 评估系统（8090端口）
  window.open('http://localhost:8090', '_blank')
}

const handleCommand = (command: string) => {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    userStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  }
}

const handleThemeChange = (mode: ThemeMode) => {
  themeStore.setMode(mode)
}

const fetchUnreadCount = async () => {
  if (userStore.user?.id) {
    try {
      const res = await messageApi.getUnreadCount(userStore.user.id)
      if (res.data.code === 0 && res.data.data) {
        unreadCount.value = res.data.data.unread_count || 0
      }
    } catch (error) {
      console.error('Failed to fetch unread count:', error)
    }
  }
}

onMounted(() => {
  fetchUnreadCount()
  // 每分钟刷新一次未读消息数
  setInterval(fetchUnreadCount, 60000)
})
</script>

<style scoped lang="scss">
.main-layout {
  height: 100vh;
  overflow: hidden;
}

.el-container {
  height: 100%;
}

.sidebar {
  background: var(--sidebar-bg);
  box-shadow: 1px 0 3px rgba(0, 0, 0, 0.05);
  transition: width 0.3s ease;
  overflow: hidden;
  border-right: 1px solid var(--border-color);

  .logo {
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    border-bottom: 1px solid var(--border-color);
    padding: 0 16px;

    .logo-icon {
      width: 36px;
      height: 36px;
      background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%);
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;

      .el-icon {
        font-size: 20px;
        color: white;
      }
    }

    .logo-text {
      margin: 0;
      font-size: 18px;
      font-weight: 700;
      color: var(--text-primary);
      white-space: nowrap;
    }
  }

  .sidebar-menu {
    border-right: none;
    background: transparent;

    :deep(.el-menu-item) {
      color: var(--sidebar-text);
      margin: 4px 8px;
      border-radius: 8px;

      &:hover {
        background: var(--sidebar-item-hover);
        color: var(--sidebar-text-active);
      }

      &.is-active {
        background: var(--sidebar-item-active);
        color: var(--sidebar-text-active);
        font-weight: 600;
      }
    }

    :deep(.el-sub-menu__title) {
      color: var(--sidebar-text);
      margin: 4px 8px;
      border-radius: 8px;

      &:hover {
        background: var(--sidebar-item-hover);
        color: var(--sidebar-text-active);
      }
    }

    &.el-menu--collapse {
      :deep(.el-menu-item) {
        margin: 4px;
      }
    }
  }
}

.header {
  background: var(--header-bg);
  box-shadow: var(--header-shadow);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;

  .header-left {
    .collapse-icon {
      font-size: 20px;
      cursor: pointer;
      color: var(--text-secondary);
      transition: color 0.3s;

      &:hover {
        color: var(--primary-color);
      }
    }
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;

    .quick-btn {
      border-radius: 8px;
      font-weight: 500;
      
      .el-icon {
        margin-right: 4px;
      }
    }

    .theme-btn {
      width: 36px;
      height: 36px;
      border-radius: 10px;
      background: var(--bg-tertiary);
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.3s ease;

      .el-icon {
        font-size: 18px;
        color: var(--text-secondary);
      }

      &:hover {
        background: var(--primary-color);

        .el-icon {
          color: white;
        }
      }
    }

    .icon-btn {
      font-size: 20px;
      cursor: pointer;
      color: var(--text-secondary);
      transition: color 0.3s;

      &:hover {
        color: var(--primary-color);
      }
    }

    .user-info {
      display: flex;
      align-items: center;
      gap: 12px;
      cursor: pointer;
      padding: 6px 12px;
      border-radius: 10px;
      transition: background-color 0.3s ease;

      &:hover {
        background: var(--bg-tertiary);
      }

      .avatar {
        background: var(--gradient-1);
      }

      .username {
        font-weight: 500;
        color: var(--text-primary);
      }
    }
  }
}

.main-content {
  background: var(--bg-secondary);
  padding: 24px;
  overflow-y: auto;
}

.badge {
  position: absolute;
  top: -5px;
  right: -5px;
}

// Dropdown menu active state
:deep(.el-dropdown-menu__item) {
  &.active {
    color: var(--primary-color);
    background: rgba(14, 165, 233, 0.1);
  }

  .el-icon {
    margin-right: 8px;
  }
}
</style>
