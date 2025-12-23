<template>
  <div class="portal-layout">
    <!-- 顶部导航 -->
    <header class="portal-header">
      <div class="header-container">
        <div class="logo" @click="$router.push('/portal')">
          <el-icon :size="28"><Connection /></el-icon>
          <span>智能招聘平台</span>
        </div>
        <nav class="nav-menu">
          <router-link to="/portal" class="nav-item">首页</router-link>
          <router-link to="/portal/jobs" class="nav-item">职位列表</router-link>
          <router-link to="/portal/companies" class="nav-item">企业招聘</router-link>
        </nav>
        <div class="header-right">
          <template v-if="userStore.isLoggedIn">
            <el-dropdown trigger="click">
              <div class="user-info">
                <el-avatar :size="36" :style="{ background: '#0ea5e9' }">
                  {{ userStore.user?.username?.charAt(0) }}
                </el-avatar>
                <span class="username">{{ userStore.user?.username }}</span>
                <el-icon><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="$router.push('/portal/my-applications')">
                    <el-icon><Document /></el-icon> 我的投递
                  </el-dropdown-item>
                  <el-dropdown-item @click="$router.push('/portal/my-resume')">
                    <el-icon><User /></el-icon> 我的简历
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleLogout">
                    <el-icon><SwitchButton /></el-icon> 退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button @click="$router.push('/portal/login')">登录</el-button>
            <el-button type="primary" @click="$router.push('/portal/register')">注册</el-button>
          </template>
        </div>
      </div>
    </header>

    <!-- 主内容区 -->
    <main class="portal-main">
      <router-view />
    </main>

    <!-- 底部 -->
    <footer class="portal-footer">
      <div class="footer-container">
        <div class="footer-info">
          <p>© 2024 智能人才招聘平台 - 毕业设计项目</p>
          <p>技术栈：Vue3 + TypeScript + Element Plus + Go + PostgreSQL</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/store/user'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Connection, ArrowDown, Document, User, SwitchButton } from '@element-plus/icons-vue'

const userStore = useUserStore()
const router = useRouter()

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/portal')
}
</script>

<style scoped lang="scss">
.portal-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f8fafc;
}

.portal-header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 100;

  .header-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 24px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 20px;
    font-weight: 700;
    color: #0ea5e9;
    cursor: pointer;
  }

  .nav-menu {
    display: flex;
    gap: 32px;

    .nav-item {
      color: #64748b;
      text-decoration: none;
      font-weight: 500;
      transition: color 0.3s;

      &:hover, &.router-link-active {
        color: #0ea5e9;
      }
    }
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 12px;

    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      padding: 4px 8px;
      border-radius: 8px;
      transition: background 0.3s;

      &:hover {
        background: #f1f5f9;
      }

      .username {
        font-weight: 500;
        color: #1e293b;
      }
    }
  }
}

.portal-main {
  flex: 1;
}

.portal-footer {
  background: #1e293b;
  color: #94a3b8;
  padding: 32px 0;

  .footer-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 24px;
    text-align: center;

    p {
      margin: 4px 0;
      font-size: 14px;
    }
  }
}
</style>
