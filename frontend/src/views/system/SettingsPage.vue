<template>
  <div class="settings-page">
    <div class="page-header">
      <h1>系统设置</h1>
      <p class="subtitle">管理系统配置和个人偏好</p>
    </div>

    <div class="settings-container">
      <!-- 侧边导航 -->
      <div class="settings-nav">
        <div
          v-for="section in settingsSections"
          :key="section.key"
          class="nav-item"
          :class="{ active: activeSection === section.key }"
          @click="activeSection = section.key"
        >
          <el-icon><component :is="section.icon" /></el-icon>
          <span>{{ section.label }}</span>
        </div>
      </div>

      <!-- 设置内容 -->
      <div class="settings-content">
        <!-- 个人信息 -->
        <div v-show="activeSection === 'profile'" class="settings-section">
          <h2>个人信息</h2>
          <p class="section-desc">管理您的个人资料和账户信息</p>

          <div class="profile-card">
            <div class="avatar-section">
              <el-avatar :size="100" :style="{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }">
                {{ userStore.user?.username?.charAt(0).toUpperCase() }}
              </el-avatar>
              <el-button type="primary" plain size="small">更换头像</el-button>
            </div>

            <el-form :model="profileForm" label-width="100px" class="profile-form">
              <el-form-item label="用户名">
                <el-input v-model="profileForm.username" disabled />
              </el-form-item>
              <el-form-item label="邮箱">
                <el-input v-model="profileForm.email" />
              </el-form-item>
              <el-form-item label="手机号">
                <el-input v-model="profileForm.phone" />
              </el-form-item>
              <el-form-item label="部门">
                <el-input v-model="profileForm.department" />
              </el-form-item>
              <el-form-item label="职位">
                <el-input v-model="profileForm.position" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="saveProfile">保存修改</el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>

        <!-- 安全设置 -->
        <div v-show="activeSection === 'security'" class="settings-section">
          <h2>安全设置</h2>
          <p class="section-desc">管理您的密码和安全选项</p>

          <div class="security-card">
            <div class="security-item">
              <div class="security-info">
                <h4>登录密码</h4>
                <p>定期更换密码可以提高账户安全性</p>
              </div>
              <el-button @click="showPasswordDialog = true">修改密码</el-button>
            </div>

            <div class="security-item">
              <div class="security-info">
                <h4>两步验证</h4>
                <p>开启后登录时需要输入验证码</p>
              </div>
              <el-switch v-model="securitySettings.twoFactor" />
            </div>

            <div class="security-item">
              <div class="security-info">
                <h4>登录通知</h4>
                <p>当账户在新设备登录时发送通知</p>
              </div>
              <el-switch v-model="securitySettings.loginNotify" />
            </div>

            <div class="security-item">
              <div class="security-info">
                <h4>会话管理</h4>
                <p>查看和管理您的登录会话</p>
              </div>
              <el-button type="primary" size="small">查看会话</el-button>
            </div>
          </div>
        </div>

        <!-- 通知设置 -->
        <div v-show="activeSection === 'notification'" class="settings-section">
          <h2>通知设置</h2>
          <p class="section-desc">管理您接收通知的方式</p>

          <div class="notification-card">
            <div class="notification-group">
              <h4>邮件通知</h4>
              <div class="notification-item">
                <span>新简历投递</span>
                <el-switch v-model="notificationSettings.emailResume" />
              </div>
              <div class="notification-item">
                <span>面试提醒</span>
                <el-switch v-model="notificationSettings.emailInterview" />
              </div>
              <div class="notification-item">
                <span>系统公告</span>
                <el-switch v-model="notificationSettings.emailSystem" />
              </div>
            </div>

            <div class="notification-group">
              <h4>站内通知</h4>
              <div class="notification-item">
                <span>新简历投递</span>
                <el-switch v-model="notificationSettings.siteResume" />
              </div>
              <div class="notification-item">
                <span>面试提醒</span>
                <el-switch v-model="notificationSettings.siteInterview" />
              </div>
              <div class="notification-item">
                <span>消息回复</span>
                <el-switch v-model="notificationSettings.siteMessage" />
              </div>
            </div>

            <div class="notification-group">
              <h4>提醒时间</h4>
              <div class="notification-item">
                <span>面试前提醒</span>
                <el-select v-model="notificationSettings.interviewReminder" style="width: 150px">
                  <el-option label="15分钟前" value="15" />
                  <el-option label="30分钟前" value="30" />
                  <el-option label="1小时前" value="60" />
                  <el-option label="1天前" value="1440" />
                </el-select>
              </div>
            </div>
          </div>
        </div>

        <!-- 外观设置 -->
        <div v-show="activeSection === 'appearance'" class="settings-section">
          <h2>外观设置</h2>
          <p class="section-desc">自定义界面外观和显示偏好</p>

          <div class="appearance-card">
            <div class="appearance-item">
              <div class="appearance-info">
                <h4>主题模式</h4>
                <p>选择您喜欢的界面主题</p>
              </div>
              <el-radio-group v-model="appearanceSettings.theme" @change="handleThemeChange">
                <el-radio-button label="light">
                  <el-icon><Sunny /></el-icon> 浅色
                </el-radio-button>
                <el-radio-button label="dark">
                  <el-icon><Moon /></el-icon> 深色
                </el-radio-button>
                <el-radio-button label="system">
                  <el-icon><Monitor /></el-icon> 跟随系统
                </el-radio-button>
              </el-radio-group>
            </div>

            <div class="appearance-item">
              <div class="appearance-info">
                <h4>主题色</h4>
                <p>选择系统主题色</p>
              </div>
              <div class="color-picker">
                <div
                  v-for="color in themeColors"
                  :key="color.value"
                  class="color-option"
                  :class="{ active: appearanceSettings.primaryColor === color.value }"
                  :style="{ background: color.value }"
                  @click="appearanceSettings.primaryColor = color.value"
                />
              </div>
            </div>

            <div class="appearance-item">
              <div class="appearance-info">
                <h4>侧边栏</h4>
                <p>默认展开或收起侧边栏</p>
              </div>
              <el-switch v-model="appearanceSettings.sidebarExpanded" active-text="展开" inactive-text="收起" />
            </div>

            <div class="appearance-item">
              <div class="appearance-info">
                <h4>紧凑模式</h4>
                <p>减少界面间距，显示更多内容</p>
              </div>
              <el-switch v-model="appearanceSettings.compactMode" />
            </div>
          </div>
        </div>

        <!-- 数据管理 -->
        <div v-show="activeSection === 'data'" class="settings-section">
          <h2>数据管理</h2>
          <p class="section-desc">管理您的数据导出和清理</p>

          <div class="data-card">
            <div class="data-item">
              <div class="data-info">
                <h4>导出数据</h4>
                <p>导出您的所有数据，包括人才、职位、面试记录等</p>
              </div>
              <el-button type="primary" plain @click="exportAllData">
                <el-icon><Download /></el-icon>
                导出全部数据
              </el-button>
            </div>

            <div class="data-item">
              <div class="data-info">
                <h4>清除缓存</h4>
                <p>清除本地缓存数据，释放存储空间</p>
              </div>
              <el-button @click="clearCache">清除缓存</el-button>
            </div>

            <div class="data-item danger">
              <div class="data-info">
                <h4>删除账户</h4>
                <p>永久删除您的账户和所有数据，此操作不可恢复</p>
              </div>
              <el-button type="danger" plain @click="showDeleteAccountDialog = true">删除账户</el-button>
            </div>
          </div>
        </div>

        <!-- 关于 -->
        <div v-show="activeSection === 'about'" class="settings-section">
          <h2>关于系统</h2>
          <p class="section-desc">系统版本和相关信息</p>

          <div class="about-card">
            <div class="about-logo">
              <div class="logo-icon">
                <el-icon :size="48"><MagicStick /></el-icon>
              </div>
              <h3>智能人才运营平台</h3>
              <p>版本 1.0.0</p>
            </div>

            <div class="about-info">
              <div class="info-item">
                <span class="label">前端框架</span>
                <span class="value">Vue 3 + TypeScript</span>
              </div>
              <div class="info-item">
                <span class="label">UI 组件</span>
                <span class="value">Element Plus</span>
              </div>
              <div class="info-item">
                <span class="label">后端框架</span>
                <span class="value">Go + Gin</span>
              </div>
              <div class="info-item">
                <span class="label">数据库</span>
                <span class="value">PostgreSQL</span>
              </div>
            </div>

            <div class="about-links">
              <el-button type="primary" plain size="small">
                <el-icon><Document /></el-icon>
                使用文档
              </el-button>
              <el-button type="primary" plain size="small">
                <el-icon><ChatDotRound /></el-icon>
                反馈建议
              </el-button>
              <el-button type="primary" plain size="small">
                <el-icon><InfoFilled /></el-icon>
                更新日志
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 修改密码弹窗 -->
    <el-dialog v-model="showPasswordDialog" title="修改密码" width="450px">
      <el-form :model="passwordForm" label-width="100px">
        <el-form-item label="当前密码">
          <el-input v-model="passwordForm.oldPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="passwordForm.newPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" @click="changePassword">确认修改</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, markRaw } from 'vue'
import { useUserStore } from '@/store/user'
import { useThemeStore, type ThemeMode } from '@/store/theme'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  User, Lock, Bell, Brush, FolderOpened, InfoFilled,
  Sunny, Moon, Monitor, Download, Document, ChatDotRound, MagicStick
} from '@element-plus/icons-vue'

const userStore = useUserStore()
const themeStore = useThemeStore()

const activeSection = ref('profile')
const showPasswordDialog = ref(false)
const showDeleteAccountDialog = ref(false)

const settingsSections = [
  { key: 'profile', label: '个人信息', icon: markRaw(User) },
  { key: 'security', label: '安全设置', icon: markRaw(Lock) },
  { key: 'notification', label: '通知设置', icon: markRaw(Bell) },
  { key: 'appearance', label: '外观设置', icon: markRaw(Brush) },
  { key: 'data', label: '数据管理', icon: markRaw(FolderOpened) },
  { key: 'about', label: '关于系统', icon: markRaw(InfoFilled) }
]

const themeColors = [
  { name: '默认紫', value: '#667eea' },
  { name: '活力蓝', value: '#3b82f6' },
  { name: '清新绿', value: '#10b981' },
  { name: '热情橙', value: '#f59e0b' },
  { name: '优雅粉', value: '#ec4899' },
  { name: '沉稳灰', value: '#6b7280' }
]

const profileForm = reactive({
  username: userStore.user?.username || '',
  email: userStore.user?.email || '',
  phone: userStore.user?.phone || '',
  department: '',
  position: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const securitySettings = reactive({
  twoFactor: false,
  loginNotify: true
})

const notificationSettings = reactive({
  emailResume: true,
  emailInterview: true,
  emailSystem: false,
  siteResume: true,
  siteInterview: true,
  siteMessage: true,
  interviewReminder: '30'
})

const appearanceSettings = reactive({
  theme: themeStore.mode,
  primaryColor: '#667eea',
  sidebarExpanded: true,
  compactMode: false
})

const saveProfile = () => {
  ElMessage.success('个人信息已保存')
}

const changePassword = () => {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  ElMessage.success('密码修改成功')
  showPasswordDialog.value = false
}

const handleThemeChange = (value: ThemeMode) => {
  themeStore.setMode(value)
}

const exportAllData = () => {
  ElMessage.success('数据导出已开始，请稍候...')
}

const clearCache = () => {
  localStorage.clear()
  ElMessage.success('缓存已清除')
}
</script>

<style scoped lang="scss">
.settings-page {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.page-header {
  margin-bottom: 24px;

  h1 {
    font-size: 24px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0 0 4px 0;
  }

  .subtitle {
    color: var(--text-secondary);
    font-size: 14px;
    margin: 0;
  }
}

.settings-container {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 24px;
  min-height: calc(100vh - 200px);

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.settings-nav {
  background: var(--bg-primary);
  border-radius: 16px;
  padding: 16px;
  box-shadow: var(--shadow-card);
  height: fit-content;

  .nav-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 16px;
    border-radius: 10px;
    cursor: pointer;
    color: var(--text-secondary);
    transition: all 0.3s ease;

    &:hover {
      background: var(--bg-tertiary);
      color: var(--text-primary);
    }

    &.active {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
    }

    .el-icon {
      font-size: 18px;
    }

    span {
      font-size: 14px;
      font-weight: 500;
    }
  }
}

.settings-content {
  background: var(--bg-primary);
  border-radius: 16px;
  padding: 32px;
  box-shadow: var(--shadow-card);
}

.settings-section {
  h2 {
    font-size: 20px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 8px 0;
  }

  .section-desc {
    color: var(--text-secondary);
    font-size: 14px;
    margin: 0 0 24px 0;
  }
}

.profile-card {
  .avatar-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
    margin-bottom: 32px;
    padding-bottom: 32px;
    border-bottom: 1px solid var(--border-light);
  }

  .profile-form {
    max-width: 500px;
  }
}

.security-card, .notification-card, .appearance-card, .data-card {
  .security-item, .notification-item, .appearance-item, .data-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 0;
    border-bottom: 1px solid var(--border-light);

    &:last-child {
      border-bottom: none;
    }

    &.danger {
      .data-info h4 {
        color: #ef4444;
      }
    }

    .security-info, .appearance-info, .data-info {
      h4 {
        font-size: 15px;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 4px 0;
      }

      p {
        font-size: 13px;
        color: var(--text-secondary);
        margin: 0;
      }
    }
  }
}

.notification-card {
  .notification-group {
    margin-bottom: 24px;

    &:last-child {
      margin-bottom: 0;
    }

    h4 {
      font-size: 14px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 12px 0;
      padding-bottom: 12px;
      border-bottom: 1px solid var(--border-light);
    }

    .notification-item {
      padding: 12px 0;
      border-bottom: none;

      span {
        font-size: 14px;
        color: var(--text-secondary);
      }
    }
  }
}

.color-picker {
  display: flex;
  gap: 12px;

  .color-option {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    cursor: pointer;
    transition: all 0.3s ease;
    border: 3px solid transparent;

    &:hover {
      transform: scale(1.1);
    }

    &.active {
      border-color: var(--text-primary);
      box-shadow: 0 0 0 2px var(--bg-primary);
    }
  }
}

.about-card {
  text-align: center;

  .about-logo {
    margin-bottom: 32px;

    .logo-icon {
      width: 100px;
      height: 100px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border-radius: 24px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      margin: 0 auto 16px;
    }

    h3 {
      font-size: 20px;
      font-weight: 700;
      color: var(--text-primary);
      margin: 0 0 4px 0;
    }

    p {
      color: var(--text-secondary);
      font-size: 14px;
      margin: 0;
    }
  }

  .about-info {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    max-width: 400px;
    margin: 0 auto 32px;
    text-align: left;

    .info-item {
      padding: 12px 16px;
      background: var(--bg-tertiary);
      border-radius: 10px;

      .label {
        display: block;
        font-size: 12px;
        color: var(--text-secondary);
        margin-bottom: 4px;
      }

      .value {
        font-size: 14px;
        font-weight: 600;
        color: var(--text-primary);
      }
    }
  }

  .about-links {
    display: flex;
    justify-content: center;
    gap: 24px;
  }
}
</style>
