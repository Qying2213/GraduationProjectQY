<template>
  <div class="main-layout">
    <el-container>
      <!-- Sidebar -->
      <el-aside width="240px" class="sidebar">
        <div class="logo">
          <h2 class="gradient-text">人才运营平台</h2>
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
            <el-icon><MagicStick /></el-icon>
            <span>智能推荐</span>
          </el-menu-item>
          
          <el-menu-item index="/messages">
            <el-icon><ChatDotRound /></el-icon>
            <span>消息中心</span>
            <el-badge v-if="unreadCount > 0" :value="unreadCount" class="badge" />
          </el-menu-item>
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
            <el-badge :value="unreadCount" :hidden="unreadCount === 0">
              <el-icon class="icon-btn" @click="goToMessages">
                <Bell />
              </el-icon>
            </el-badge>
            
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
import { messageApi } from '@/api/message'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const isCollapse = ref(false)
const unreadCount = ref(0)

const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

const goToMessages = () => {
  router.push('/messages')
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
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  
  .logo {
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    
    h2 {
      margin: 0;
      font-size: 20px;
      font-weight: 700;
      color: white;
      background: none;
      -webkit-background-clip: unset;
      -webkit-text-fill-color: white;
    }
  }
  
  .sidebar-menu {
    border-right: none;
    background: transparent;
    
    :deep(.el-menu-item) {
      color: rgba(255, 255, 255, 0.8);
      
      &:hover {
        background: rgba(255, 255, 255, 0.1);
        color: white;
      }
      
      &.is-active {
        background: rgba(255, 255, 255, 0.2);
        color: white;
      }
    }
  }
}

.header {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
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
    gap: 24px;
    
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
</style>
