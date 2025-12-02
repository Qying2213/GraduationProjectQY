<template>
  <div class="user-profile">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>个人中心</h1>
      <p class="subtitle">管理您的账户信息和偏好设置</p>
    </div>

    <div class="profile-content">
      <!-- 左侧导航 -->
      <div class="profile-sidebar">
        <div class="user-card">
          <div class="avatar-wrapper" @click="showAvatarUpload = true">
            <img :src="userInfo.avatar || defaultAvatar" :alt="userInfo.name" />
            <div class="avatar-overlay">
              <el-icon><Camera /></el-icon>
              <span>更换头像</span>
            </div>
          </div>
          <h3 class="user-name">{{ userInfo.name }}</h3>
          <p class="user-role">{{ userInfo.role }}</p>
          <div class="user-stats">
            <div class="stat-item">
              <span class="stat-value">{{ userInfo.jobCount }}</span>
              <span class="stat-label">发布职位</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ userInfo.talentCount }}</span>
              <span class="stat-label">收藏人才</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ userInfo.interviewCount }}</span>
              <span class="stat-label">面试安排</span>
            </div>
          </div>
        </div>

        <div class="nav-menu">
          <div
            v-for="item in menuItems"
            :key="item.key"
            class="menu-item"
            :class="{ active: activeTab === item.key }"
            @click="activeTab = item.key"
          >
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </div>
        </div>
      </div>

      <!-- 右侧内容 -->
      <div class="profile-main">
        <!-- 基本信息 -->
        <div v-show="activeTab === 'basic'" class="section-card">
          <div class="section-header">
            <h2>基本信息</h2>
            <el-button
              v-if="!isEditing"
              type="primary"
              @click="isEditing = true"
            >
              <el-icon><Edit /></el-icon>
              编辑信息
            </el-button>
            <div v-else class="edit-actions">
              <el-button @click="cancelEdit">取消</el-button>
              <el-button type="primary" @click="saveProfile">保存</el-button>
            </div>
          </div>

          <el-form
            ref="profileFormRef"
            :model="profileForm"
            :rules="profileRules"
            label-width="100px"
            :disabled="!isEditing"
          >
            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="姓名" prop="name">
                  <el-input v-model="profileForm.name" placeholder="请输入姓名" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="性别" prop="gender">
                  <el-radio-group v-model="profileForm.gender">
                    <el-radio label="male">男</el-radio>
                    <el-radio label="female">女</el-radio>
                  </el-radio-group>
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="手机号" prop="phone">
                  <el-input v-model="profileForm.phone" placeholder="请输入手机号" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="邮箱" prop="email">
                  <el-input v-model="profileForm.email" placeholder="请输入邮箱" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="公司" prop="company">
                  <el-input v-model="profileForm.company" placeholder="请输入公司名称" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="职位" prop="position">
                  <el-input v-model="profileForm.position" placeholder="请输入职位" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="所在城市" prop="city">
                  <el-select v-model="profileForm.city" placeholder="请选择城市" style="width: 100%">
                    <el-option label="北京" value="北京" />
                    <el-option label="上海" value="上海" />
                    <el-option label="广州" value="广州" />
                    <el-option label="深圳" value="深圳" />
                    <el-option label="杭州" value="杭州" />
                    <el-option label="成都" value="成都" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="入职日期" prop="joinDate">
                  <el-date-picker
                    v-model="profileForm.joinDate"
                    type="date"
                    placeholder="请选择入职日期"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="个人简介" prop="bio">
              <el-input
                v-model="profileForm.bio"
                type="textarea"
                :rows="4"
                placeholder="请输入个人简介"
              />
            </el-form-item>
          </el-form>
        </div>

        <!-- 账户安全 -->
        <div v-show="activeTab === 'security'" class="section-card">
          <div class="section-header">
            <h2>账户安全</h2>
          </div>

          <div class="security-items">
            <div class="security-item">
              <div class="security-info">
                <div class="security-icon">
                  <el-icon><Lock /></el-icon>
                </div>
                <div class="security-detail">
                  <h4>登录密码</h4>
                  <p>定期更换密码可以保护账户安全</p>
                </div>
              </div>
              <el-button @click="showPasswordDialog = true">修改密码</el-button>
            </div>

            <div class="security-item">
              <div class="security-info">
                <div class="security-icon phone">
                  <el-icon><Iphone /></el-icon>
                </div>
                <div class="security-detail">
                  <h4>手机绑定</h4>
                  <p>已绑定：{{ userInfo.phone }}</p>
                </div>
              </div>
              <el-button>更换手机</el-button>
            </div>

            <div class="security-item">
              <div class="security-info">
                <div class="security-icon email">
                  <el-icon><Message /></el-icon>
                </div>
                <div class="security-detail">
                  <h4>邮箱绑定</h4>
                  <p>已绑定：{{ userInfo.email }}</p>
                </div>
              </div>
              <el-button>更换邮箱</el-button>
            </div>

            <div class="security-item">
              <div class="security-info">
                <div class="security-icon wechat">
                  <el-icon><ChatDotRound /></el-icon>
                </div>
                <div class="security-detail">
                  <h4>微信绑定</h4>
                  <p>{{ userInfo.wechatBound ? '已绑定' : '未绑定' }}</p>
                </div>
              </div>
              <el-button type="primary">{{ userInfo.wechatBound ? '解除绑定' : '立即绑定' }}</el-button>
            </div>
          </div>

          <!-- 登录记录 -->
          <div class="login-history">
            <h3>最近登录记录</h3>
            <el-table :data="loginHistory" style="width: 100%">
              <el-table-column prop="time" label="登录时间" width="180" />
              <el-table-column prop="ip" label="IP地址" width="150" />
              <el-table-column prop="location" label="登录地点" width="150" />
              <el-table-column prop="device" label="设备" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === '成功' ? 'success' : 'danger'" size="small">
                    {{ row.status }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>

        <!-- 通知设置 -->
        <div v-show="activeTab === 'notification'" class="section-card">
          <div class="section-header">
            <h2>通知设置</h2>
          </div>

          <div class="notification-settings">
            <div class="notification-group">
              <h3>消息通知</h3>
              <div class="notification-item">
                <div class="notification-info">
                  <h4>系统消息</h4>
                  <p>接收系统公告、维护通知等</p>
                </div>
                <el-switch v-model="notificationSettings.system" />
              </div>
              <div class="notification-item">
                <div class="notification-info">
                  <h4>面试邀约</h4>
                  <p>接收新的面试邀请通知</p>
                </div>
                <el-switch v-model="notificationSettings.interview" />
              </div>
              <div class="notification-item">
                <div class="notification-info">
                  <h4>简历投递</h4>
                  <p>有求职者投递简历时通知</p>
                </div>
                <el-switch v-model="notificationSettings.resume" />
              </div>
              <div class="notification-item">
                <div class="notification-info">
                  <h4>人才推荐</h4>
                  <p>接收智能推荐的匹配人才</p>
                </div>
                <el-switch v-model="notificationSettings.recommend" />
              </div>
            </div>

            <div class="notification-group">
              <h3>通知方式</h3>
              <div class="notification-item">
                <div class="notification-info">
                  <h4>站内消息</h4>
                  <p>在消息中心接收通知</p>
                </div>
                <el-switch v-model="notificationSettings.inApp" />
              </div>
              <div class="notification-item">
                <div class="notification-info">
                  <h4>邮件通知</h4>
                  <p>通过邮件接收重要通知</p>
                </div>
                <el-switch v-model="notificationSettings.email" />
              </div>
              <div class="notification-item">
                <div class="notification-info">
                  <h4>短信通知</h4>
                  <p>通过短信接收紧急通知</p>
                </div>
                <el-switch v-model="notificationSettings.sms" />
              </div>
            </div>
          </div>

          <div class="save-settings">
            <el-button type="primary" @click="saveNotificationSettings">保存设置</el-button>
          </div>
        </div>

        <!-- 操作记录 -->
        <div v-show="activeTab === 'activity'" class="section-card">
          <div class="section-header">
            <h2>操作记录</h2>
          </div>

          <el-timeline>
            <el-timeline-item
              v-for="(activity, index) in activities"
              :key="index"
              :timestamp="activity.time"
              placement="top"
              :type="activity.type"
            >
              <div class="activity-card">
                <div class="activity-icon" :class="activity.category">
                  <el-icon><component :is="activity.icon" /></el-icon>
                </div>
                <div class="activity-content">
                  <h4>{{ activity.title }}</h4>
                  <p>{{ activity.description }}</p>
                </div>
              </div>
            </el-timeline-item>
          </el-timeline>

          <div class="load-more">
            <el-button text type="primary">加载更多</el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 头像上传对话框 -->
    <el-dialog v-model="showAvatarUpload" title="更换头像" width="500px">
      <div class="avatar-upload-content">
        <el-upload
          class="avatar-uploader"
          :show-file-list="false"
          :before-upload="beforeAvatarUpload"
          :http-request="handleAvatarUpload"
          accept="image/*"
        >
          <div v-if="avatarPreview" class="avatar-preview">
            <img :src="avatarPreview" alt="avatar" />
          </div>
          <div v-else class="upload-placeholder">
            <el-icon class="upload-icon"><Plus /></el-icon>
            <span>点击上传头像</span>
          </div>
        </el-upload>
        <div class="upload-tips">
          <p>支持 JPG、PNG 格式，大小不超过 2MB</p>
          <p>建议尺寸 200x200 像素</p>
        </div>
      </div>
      <template #footer>
        <el-button @click="showAvatarUpload = false">取消</el-button>
        <el-button type="primary" @click="confirmAvatarUpload">确认更换</el-button>
      </template>
    </el-dialog>

    <!-- 修改密码对话框 -->
    <el-dialog v-model="showPasswordDialog" title="修改密码" width="500px">
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
      >
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            show-password
            placeholder="请输入当前密码"
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            show-password
            placeholder="请输入新密码"
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            show-password
            placeholder="请再次输入新密码"
          />
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
import { ref, reactive, markRaw, type Component } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules, UploadRawFile } from 'element-plus'
import {
  Camera, Edit, Lock, Iphone, Message, ChatDotRound, Plus,
  User, Setting, Bell, Clock, Document, Suitcase, Star,
  Position
} from '@element-plus/icons-vue'

// 类型定义
interface MenuItem {
  key: string
  label: string
  icon: Component
}

interface Activity {
  time: string
  title: string
  description: string
  category: string
  icon: Component
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
}

// 默认头像
const defaultAvatar = 'https://api.dicebear.com/7.x/avataaars/svg?seed=default'

// 状态
const activeTab = ref('basic')
const isEditing = ref(false)
const showAvatarUpload = ref(false)
const showPasswordDialog = ref(false)
const avatarPreview = ref('')

// 用户信息
const userInfo = reactive({
  name: '张三',
  role: 'HR经理',
  avatar: '',
  phone: '138****8888',
  email: 'zhang***@example.com',
  jobCount: 12,
  talentCount: 86,
  interviewCount: 24,
  wechatBound: false
})

// 表单
const profileFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()

const profileForm = reactive({
  name: '张三',
  gender: 'male',
  phone: '13812345678',
  email: 'zhangsan@example.com',
  company: '字节跳动',
  position: 'HR经理',
  city: '北京',
  joinDate: '2022-03-15',
  bio: '资深HR，专注于技术人才招聘，擅长人才识别与团队建设。'
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 验证规则
const profileRules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

const passwordRules: FormRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 菜单项
const menuItems = ref<MenuItem[]>([
  { key: 'basic', label: '基本信息', icon: markRaw(User) },
  { key: 'security', label: '账户安全', icon: markRaw(Lock) },
  { key: 'notification', label: '通知设置', icon: markRaw(Bell) },
  { key: 'activity', label: '操作记录', icon: markRaw(Clock) }
])

// 通知设置
const notificationSettings = reactive({
  system: true,
  interview: true,
  resume: true,
  recommend: true,
  inApp: true,
  email: true,
  sms: false
})

// 登录记录
const loginHistory = ref([
  { time: '2024-01-10 09:30:22', ip: '192.168.1.100', location: '北京', device: 'Chrome / Windows', status: '成功' },
  { time: '2024-01-09 18:45:11', ip: '192.168.1.100', location: '北京', device: 'Safari / macOS', status: '成功' },
  { time: '2024-01-09 08:20:33', ip: '10.0.0.55', location: '上海', device: 'Mobile / iOS', status: '成功' },
  { time: '2024-01-08 22:15:44', ip: '203.156.78.90', location: '深圳', device: 'Chrome / Windows', status: '失败' },
  { time: '2024-01-08 14:30:00', ip: '192.168.1.100', location: '北京', device: 'Chrome / Windows', status: '成功' }
])

// 操作记录
const activities = ref<Activity[]>([
  {
    time: '2024-01-10 10:30',
    title: '发布了新职位',
    description: '发布了"高级前端工程师"职位，薪资范围30-50K',
    category: 'job',
    icon: markRaw(Suitcase),
    type: 'primary'
  },
  {
    time: '2024-01-10 09:15',
    title: '收藏了人才',
    description: '将"李明 - 资深前端工程师"添加到收藏夹',
    category: 'talent',
    icon: markRaw(Star),
    type: 'success'
  },
  {
    time: '2024-01-09 16:45',
    title: '安排了面试',
    description: '为"王强 - 后端工程师"安排了技术面试，时间：1月15日 14:00',
    category: 'interview',
    icon: markRaw(Position),
    type: 'warning'
  },
  {
    time: '2024-01-09 11:20',
    title: '查看了简历',
    description: '查看了"张伟"的简历详情',
    category: 'resume',
    icon: markRaw(Document)
  },
  {
    time: '2024-01-08 15:30',
    title: '修改了职位信息',
    description: '更新了"产品经理"职位的薪资范围和工作要求',
    category: 'job',
    icon: markRaw(Edit)
  }
])

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
  // 重置表单
}

// 保存个人信息
const saveProfile = async () => {
  if (!profileFormRef.value) return
  await profileFormRef.value.validate()
  ElMessage.success('个人信息保存成功')
  isEditing.value = false
}

// 头像上传前检查
const beforeAvatarUpload = (file: UploadRawFile) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

// 处理头像上传
const handleAvatarUpload = (options: { file: File }) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    avatarPreview.value = e.target?.result as string
  }
  reader.readAsDataURL(options.file)
}

// 确认更换头像
const confirmAvatarUpload = () => {
  if (avatarPreview.value) {
    userInfo.avatar = avatarPreview.value
    ElMessage.success('头像更换成功')
    showAvatarUpload.value = false
    avatarPreview.value = ''
  }
}

// 修改密码
const changePassword = async () => {
  if (!passwordFormRef.value) return
  await passwordFormRef.value.validate()
  ElMessage.success('密码修改成功')
  showPasswordDialog.value = false
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}

// 保存通知设置
const saveNotificationSettings = () => {
  ElMessage.success('通知设置已保存')
}
</script>

<style scoped lang="scss">
.user-profile {
  padding: 24px;
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.page-header {
  margin-bottom: 24px;

  h1 {
    font-size: 28px;
    font-weight: 700;
    color: #1a1a2e;
    margin: 0 0 8px 0;
  }

  .subtitle {
    font-size: 14px;
    color: #6b7280;
    margin: 0;
  }
}

.profile-content {
  display: flex;
  gap: 24px;
}

// 左侧边栏
.profile-sidebar {
  width: 280px;
  flex-shrink: 0;
}

.user-card {
  background: white;
  border-radius: 20px;
  padding: 32px 24px;
  text-align: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  margin-bottom: 20px;

  .avatar-wrapper {
    width: 100px;
    height: 100px;
    margin: 0 auto 16px;
    border-radius: 50%;
    position: relative;
    cursor: pointer;
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .avatar-overlay {
      position: absolute;
      inset: 0;
      background: rgba(0, 0, 0, 0.5);
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      color: white;
      opacity: 0;
      transition: opacity 0.3s ease;

      .el-icon {
        font-size: 24px;
        margin-bottom: 4px;
      }

      span {
        font-size: 12px;
      }
    }

    &:hover .avatar-overlay {
      opacity: 1;
    }
  }

  .user-name {
    font-size: 20px;
    font-weight: 600;
    color: #1a1a2e;
    margin: 0 0 8px 0;
  }

  .user-role {
    font-size: 14px;
    color: #6b7280;
    margin: 0 0 20px 0;
  }

  .user-stats {
    display: flex;
    justify-content: center;
    gap: 24px;
    padding-top: 20px;
    border-top: 1px solid #f0f2f5;

    .stat-item {
      text-align: center;

      .stat-value {
        display: block;
        font-size: 20px;
        font-weight: 700;
        color: #667eea;
      }

      .stat-label {
        font-size: 12px;
        color: #9ca3af;
      }
    }
  }
}

.nav-menu {
  background: white;
  border-radius: 16px;
  padding: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

  .menu-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 16px;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
    color: #4b5563;

    &:hover {
      background: #f5f7fa;
    }

    &.active {
      background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
      color: #667eea;
      font-weight: 600;
    }

    .el-icon {
      font-size: 20px;
    }
  }
}

// 右侧主内容
.profile-main {
  flex: 1;
}

.section-card {
  background: white;
  border-radius: 20px;
  padding: 28px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    h2 {
      font-size: 20px;
      font-weight: 600;
      color: #1a1a2e;
      margin: 0;
    }

    .edit-actions {
      display: flex;
      gap: 12px;
    }
  }
}

// 安全设置
.security-items {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 32px;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: #f8f9fc;
  border-radius: 14px;

  .security-info {
    display: flex;
    align-items: center;
    gap: 16px;

    .security-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      display: flex;
      align-items: center;
      justify-content: center;

      .el-icon {
        font-size: 22px;
        color: white;
      }

      &.phone { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
      &.email { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
      &.wechat { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }
    }

    .security-detail {
      h4 {
        font-size: 15px;
        font-weight: 600;
        color: #1a1a2e;
        margin: 0 0 4px 0;
      }

      p {
        font-size: 13px;
        color: #6b7280;
        margin: 0;
      }
    }
  }
}

.login-history {
  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #1a1a2e;
    margin: 0 0 16px 0;
  }
}

// 通知设置
.notification-settings {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.notification-group {
  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #1a1a2e;
    margin: 0 0 16px 0;
    padding-bottom: 12px;
    border-bottom: 1px solid #f0f2f5;
  }
}

.notification-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #f0f2f5;

  &:last-child {
    border-bottom: none;
  }

  .notification-info {
    h4 {
      font-size: 14px;
      font-weight: 500;
      color: #1a1a2e;
      margin: 0 0 4px 0;
    }

    p {
      font-size: 13px;
      color: #6b7280;
      margin: 0;
    }
  }
}

.save-settings {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #f0f2f5;
  text-align: right;
}

// 操作记录
.activity-card {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: #f8f9fc;
  border-radius: 12px;

  .activity-icon {
    width: 40px;
    height: 40px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;

    .el-icon {
      font-size: 18px;
      color: white;
    }

    &.job { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
    &.talent { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
    &.interview { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
    &.resume { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }
  }

  .activity-content {
    h4 {
      font-size: 14px;
      font-weight: 600;
      color: #1a1a2e;
      margin: 0 0 4px 0;
    }

    p {
      font-size: 13px;
      color: #6b7280;
      margin: 0;
    }
  }
}

.load-more {
  text-align: center;
  padding-top: 20px;
}

// 头像上传
.avatar-upload-content {
  text-align: center;

  .avatar-uploader {
    .avatar-preview {
      width: 200px;
      height: 200px;
      border-radius: 50%;
      overflow: hidden;
      margin: 0 auto;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .upload-placeholder {
      width: 200px;
      height: 200px;
      border: 2px dashed #dcdfe6;
      border-radius: 50%;
      margin: 0 auto;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: border-color 0.3s ease;

      &:hover {
        border-color: #667eea;
      }

      .upload-icon {
        font-size: 40px;
        color: #c0c4cc;
        margin-bottom: 8px;
      }

      span {
        font-size: 14px;
        color: #606266;
      }
    }
  }

  .upload-tips {
    margin-top: 20px;

    p {
      font-size: 13px;
      color: #909399;
      margin: 4px 0;
    }
  }
}

// 响应式
@media (max-width: 1024px) {
  .profile-content {
    flex-direction: column;
  }

  .profile-sidebar {
    width: 100%;
  }

  .user-card {
    .user-stats {
      justify-content: space-around;
    }
  }

  .nav-menu {
    display: flex;
    overflow-x: auto;
    gap: 8px;

    &::-webkit-scrollbar {
      display: none;
    }

    .menu-item {
      flex-shrink: 0;
      padding: 12px 16px;
    }
  }
}

@media (max-width: 768px) {
  .user-profile {
    padding: 16px;
  }

  .section-card {
    padding: 20px;
  }

  .security-item {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;

    .el-button {
      width: 100%;
    }
  }
}
</style>
