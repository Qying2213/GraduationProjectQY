<template>
  <div class="message-center">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1>消息中心</h1>
        <p class="subtitle">及时掌握系统通知、面试邀约等重要信息</p>
      </div>
      <div class="header-actions">
        <el-button @click="markAllAsRead" :disabled="unreadCount === 0">
          <el-icon><Check /></el-icon>
          全部已读
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon all">
          <el-icon><Message /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ totalCount }}</span>
          <span class="stat-label">全部消息</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon unread">
          <el-icon><Bell /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ unreadCount }}</span>
          <span class="stat-label">未读消息</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon interview">
          <el-icon><Calendar /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ interviewCount }}</span>
          <span class="stat-label">面试邀约</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon system">
          <el-icon><Setting /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ systemCount }}</span>
          <span class="stat-label">系统通知</span>
        </div>
      </div>
    </div>

    <div class="content-wrapper">
      <!-- 左侧分类 -->
      <div class="category-sidebar">
        <div class="category-list">
          <div
            v-for="cat in categories"
            :key="cat.key"
            class="category-item"
            :class="{ active: activeCategory === cat.key }"
            @click="activeCategory = cat.key"
          >
            <div class="category-icon" :class="cat.key">
              <el-icon><component :is="cat.icon" /></el-icon>
            </div>
            <span class="category-name">{{ cat.name }}</span>
            <el-badge
              v-if="getCategoryUnread(cat.key) > 0"
              :value="getCategoryUnread(cat.key)"
              class="category-badge"
            />
          </div>
        </div>
      </div>

      <!-- 右侧消息列表 -->
      <div class="message-list-wrapper">
        <!-- 工具栏 -->
        <div class="toolbar">
          <div class="toolbar-left">
            <el-checkbox
              v-model="selectAll"
              :indeterminate="isIndeterminate"
              @change="handleSelectAll"
            >
              全选
            </el-checkbox>
            <el-button
              v-if="selectedIds.length > 0"
              type="primary"
              text
              @click="batchMarkAsRead"
            >
              标记已读
            </el-button>
            <el-button
              v-if="selectedIds.length > 0"
              type="danger"
              text
              @click="batchDelete"
            >
              删除
            </el-button>
          </div>
          <div class="toolbar-right">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索消息内容"
              prefix-icon="Search"
              clearable
              style="width: 240px"
            />
          </div>
        </div>

        <!-- 消息列表 -->
        <div class="message-list" v-loading="loading">
          <div
            v-for="msg in filteredMessages"
            :key="msg.id"
            class="message-item"
            :class="{ unread: !msg.isRead, selected: selectedIds.includes(msg.id) }"
            @click="viewMessage(msg)"
          >
            <el-checkbox
              v-model="selectedIds"
              :label="msg.id"
              @click.stop
            />
            <div class="message-icon" :class="msg.type">
              <el-icon><component :is="getMessageIcon(msg.type)" /></el-icon>
            </div>
            <div class="message-content">
              <div class="message-header">
                <span class="message-title">{{ msg.title }}</span>
                <el-tag
                  v-if="!msg.isRead"
                  type="danger"
                  size="small"
                  effect="dark"
                >
                  新
                </el-tag>
              </div>
              <p class="message-summary">{{ msg.content }}</p>
              <div class="message-meta">
                <span class="message-sender">
                  <el-icon><User /></el-icon>
                  {{ msg.sender }}
                </span>
                <span class="message-time">
                  <el-icon><Clock /></el-icon>
                  {{ formatTime(msg.createdAt) }}
                </span>
              </div>
            </div>
            <div class="message-actions" @click.stop>
              <el-button
                v-if="!msg.isRead"
                type="primary"
                text
                size="small"
                @click="markAsRead(msg)"
              >
                标记已读
              </el-button>
              <el-button
                type="danger"
                text
                size="small"
                @click="deleteMessage(msg)"
              >
                删除
              </el-button>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="filteredMessages.length === 0 && !loading" class="empty-state">
            <el-empty description="暂无消息" />
          </div>
        </div>

        <!-- 分页 -->
        <div class="pagination-wrapper" v-if="filteredMessages.length > 0">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50]"
            :total="total"
            layout="total, sizes, prev, pager, next"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </div>
    </div>

    <!-- 消息详情抽屉 -->
    <el-drawer
      v-model="showDetailDrawer"
      title="消息详情"
      size="500px"
    >
      <div class="message-detail" v-if="currentMessage">
        <div class="detail-header">
          <div class="detail-icon" :class="currentMessage.type">
            <el-icon><component :is="getMessageIcon(currentMessage.type)" /></el-icon>
          </div>
          <div class="detail-info">
            <h3 class="detail-title">{{ currentMessage.title }}</h3>
            <div class="detail-meta">
              <span>
                <el-icon><User /></el-icon>
                {{ currentMessage.sender }}
              </span>
              <span>
                <el-icon><Clock /></el-icon>
                {{ formatTime(currentMessage.createdAt) }}
              </span>
            </div>
          </div>
        </div>

        <el-divider />

        <div class="detail-content">
          {{ currentMessage.content }}
        </div>

        <!-- 附加信息 -->
        <div class="detail-extra" v-if="currentMessage.extra">
          <h4>相关信息</h4>
          <el-descriptions :column="1" border>
            <el-descriptions-item
              v-for="(value, key) in currentMessage.extra"
              :key="key"
              :label="getExtraLabel(key as string)"
            >
              {{ value }}
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 操作按钮 -->
        <div class="detail-actions">
          <el-button
            v-if="currentMessage.type === 'interview'"
            type="primary"
            size="large"
          >
            <el-icon><Check /></el-icon>
            接受邀请
          </el-button>
          <el-button
            v-if="currentMessage.type === 'interview'"
            size="large"
          >
            <el-icon><Close /></el-icon>
            婉拒邀请
          </el-button>
          <el-button
            v-if="currentMessage.type === 'application'"
            type="primary"
            size="large"
          >
            <el-icon><View /></el-icon>
            查看简历
          </el-button>
          <el-button type="danger" size="large" @click="deleteMessage(currentMessage)">
            <el-icon><Delete /></el-icon>
            删除消息
          </el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, markRaw, type Component } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Message, Bell, Calendar, Setting, Check, User, Clock,
  Promotion, Document, Close, View, Delete, InfoFilled
} from '@element-plus/icons-vue'
import { messageApi } from '@/api/message'

// 类型定义
interface MessageItem {
  id: number
  type: 'system' | 'interview' | 'application' | 'notification'
  title: string
  content: string
  sender: string
  isRead: boolean
  createdAt: string
  extra?: Record<string, string>
}

interface Category {
  key: string
  name: string
  icon: Component
}

// 分类
const categories = ref<Category[]>([
  { key: 'all', name: '全部消息', icon: markRaw(Message) },
  { key: 'unread', name: '未读消息', icon: markRaw(Bell) },
  { key: 'interview', name: '面试邀约', icon: markRaw(Calendar) },
  { key: 'application', name: '简历投递', icon: markRaw(Document) },
  { key: 'notification', name: '系统通知', icon: markRaw(Promotion) },
  { key: 'system', name: '系统公告', icon: markRaw(Setting) }
])

// 状态
const activeCategory = ref('all')
const messages = ref<MessageItem[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')
const selectedIds = ref<number[]>([])
const showDetailDrawer = ref(false)
const currentMessage = ref<MessageItem | null>(null)

// 计算属性
const totalCount = computed(() => messages.value.length)
const unreadCount = computed(() => messages.value.filter(m => !m.isRead).length)
const interviewCount = computed(() => messages.value.filter(m => m.type === 'interview').length)
const systemCount = computed(() => messages.value.filter(m => m.type === 'system').length)

const selectAll = computed({
  get: () => selectedIds.value.length === filteredMessages.value.length && filteredMessages.value.length > 0,
  set: (val) => {
    selectedIds.value = val ? filteredMessages.value.map(m => m.id) : []
  }
})

const isIndeterminate = computed(() => {
  return selectedIds.value.length > 0 && selectedIds.value.length < filteredMessages.value.length
})

const filteredMessages = computed(() => {
  let result = messages.value

  // 按分类筛选
  if (activeCategory.value === 'unread') {
    result = result.filter(m => !m.isRead)
  } else if (activeCategory.value !== 'all') {
    result = result.filter(m => m.type === activeCategory.value)
  }

  // 按关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(m =>
      m.title.toLowerCase().includes(keyword) ||
      m.content.toLowerCase().includes(keyword)
    )
  }

  return result
})

// 获取分类未读数
const getCategoryUnread = (key: string) => {
  if (key === 'all') return unreadCount.value
  if (key === 'unread') return unreadCount.value
  return messages.value.filter(m => m.type === key && !m.isRead).length
}

// 获取消息图标
const getMessageIcon = (type: string): Component => {
  const icons: Record<string, Component> = {
    system: markRaw(Setting),
    interview: markRaw(Calendar),
    application: markRaw(Document),
    notification: markRaw(Promotion)
  }
  return icons[type] || markRaw(InfoFilled)
}

// 获取额外信息标签
const getExtraLabel = (key: string) => {
  const labels: Record<string, string> = {
    position: '应聘职位',
    company: '公司名称',
    interviewTime: '面试时间',
    interviewAddress: '面试地址',
    contact: '联系人',
    phone: '联系电话'
  }
  return labels[key] || key
}

// 格式化时间
const formatTime = (time: string) => {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString()
}

// 生成模拟数据
const generateMockMessages = (): MessageItem[] => {
  return [
    {
      id: 1,
      type: 'interview',
      title: '面试邀请 - 高级前端工程师',
      content: '尊敬的张先生，您好！感谢您对我司职位的关注。我们诚挚邀请您参加高级前端工程师岗位的面试。面试时间为2024年1月15日下午14:00，地点为北京市朝阳区望京SOHO T1栋15层。请携带个人简历及身份证件。如有疑问，请联系HR李小姐。',
      sender: '字节跳动 HR',
      isRead: false,
      createdAt: new Date(Date.now() - 1800000).toISOString(),
      extra: {
        position: '高级前端工程师',
        company: '字节跳动',
        interviewTime: '2024年1月15日 14:00',
        interviewAddress: '北京市朝阳区望京SOHO T1栋15层',
        contact: '李小姐',
        phone: '138****8888'
      }
    },
    {
      id: 2,
      type: 'application',
      title: '简历投递提醒 - 王伟投递了您的职位',
      content: '王伟向您发布的"资深后端工程师"职位投递了简历，请及时查看处理。该候选人具有5年Java开发经验，曾就职于阿里巴巴、腾讯等知名企业。',
      sender: '系统通知',
      isRead: false,
      createdAt: new Date(Date.now() - 3600000).toISOString(),
      extra: {
        position: '资深后端工程师'
      }
    },
    {
      id: 3,
      type: 'notification',
      title: '您的职位即将过期',
      content: '您发布的"产品经理"职位将于3天后过期，如需继续招聘，请及时刷新职位或重新发布。过期后职位将自动下线，求职者将无法查看和投递。',
      sender: '系统通知',
      isRead: false,
      createdAt: new Date(Date.now() - 7200000).toISOString()
    },
    {
      id: 4,
      type: 'interview',
      title: '面试结果通知 - 技术面试通过',
      content: '恭喜您！您已通过我司高级前端工程师岗位的技术面试。接下来将安排HR面试，具体时间我们将另行通知。请保持电话畅通。',
      sender: '阿里巴巴 技术部',
      isRead: true,
      createdAt: new Date(Date.now() - 86400000).toISOString()
    },
    {
      id: 5,
      type: 'system',
      title: '系统升级公告',
      content: '为了给您提供更好的服务体验，我们将于2024年1月20日凌晨2:00-6:00进行系统升级维护。升级期间部分功能可能无法正常使用，请提前做好安排。给您带来的不便，敬请谅解。',
      sender: '系统管理员',
      isRead: true,
      createdAt: new Date(Date.now() - 172800000).toISOString()
    },
    {
      id: 6,
      type: 'application',
      title: '简历投递提醒 - 李娜投递了您的职位',
      content: '李娜向您发布的"UI设计师"职位投递了简历。该候选人有4年UI设计经验，精通Figma、Sketch等设计工具，有丰富的移动端和Web端设计经验。',
      sender: '系统通知',
      isRead: true,
      createdAt: new Date(Date.now() - 259200000).toISOString(),
      extra: {
        position: 'UI设计师'
      }
    },
    {
      id: 7,
      type: 'notification',
      title: '人才推荐 - 发现5位匹配人才',
      content: '根据您的职位需求，系统为您智能推荐了5位高度匹配的人才。其中3位匹配度超过90%，建议您及时查看并主动沟通。',
      sender: '智能推荐系统',
      isRead: false,
      createdAt: new Date(Date.now() - 14400000).toISOString()
    },
    {
      id: 8,
      type: 'system',
      title: '账户安全提醒',
      content: '检测到您的账户于今日在新设备上登录，登录地点：北京市。如非本人操作，请立即修改密码并联系客服。',
      sender: '安全中心',
      isRead: true,
      createdAt: new Date(Date.now() - 432000000).toISOString()
    }
  ]
}

// 加载消息
const loadMessages = async () => {
  loading.value = true
  try {
    const res = await messageApi.list({
      page: currentPage.value,
      page_size: pageSize.value,
      type: activeTab.value !== 'all' ? activeTab.value : undefined
    })
    
    if (res.data?.code === 0 && res.data.data) {
      messages.value = (res.data.data.messages || []).map((m: any) => ({
        id: m.id,
        type: m.type || 'system',
        title: m.title,
        content: m.content,
        sender: m.sender_name || '系统',
        isRead: m.is_read,
        createdAt: m.created_at
      }))
      total.value = res.data.data.total || 0
    } else {
      ElMessage.error(res.data?.message || '获取消息失败')
    }
  } catch (error) {
    console.error('获取消息失败:', error)
    ElMessage.error('获取消息失败')
    messages.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 查看消息
const viewMessage = (msg: MessageItem) => {
  currentMessage.value = msg
  showDetailDrawer.value = true
  if (!msg.isRead) {
    msg.isRead = true
  }
}

// 标记已读
const markAsRead = (msg: MessageItem) => {
  msg.isRead = true
  ElMessage.success('已标记为已读')
}

// 全部标记已读
const markAllAsRead = () => {
  messages.value.forEach(m => m.isRead = true)
  ElMessage.success('已全部标记为已读')
}

// 批量标记已读
const batchMarkAsRead = () => {
  messages.value.forEach(m => {
    if (selectedIds.value.includes(m.id)) {
      m.isRead = true
    }
  })
  selectedIds.value = []
  ElMessage.success('已标记为已读')
}

// 删除消息
const deleteMessage = async (msg: MessageItem) => {
  try {
    await ElMessageBox.confirm('确定要删除这条消息吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    messages.value = messages.value.filter(m => m.id !== msg.id)
    showDetailDrawer.value = false
    ElMessage.success('删除成功')
  } catch {
    // 用户取消
  }
}

// 批量删除
const batchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.value.length} 条消息吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    messages.value = messages.value.filter(m => !selectedIds.value.includes(m.id))
    selectedIds.value = []
    ElMessage.success('删除成功')
  } catch {
    // 用户取消
  }
}

// 全选处理
const handleSelectAll = (val: boolean) => {
  selectedIds.value = val ? filteredMessages.value.map(m => m.id) : []
}

// 分页处理
const handleSizeChange = () => {
  currentPage.value = 1
}

const handlePageChange = () => {
  // 分页逻辑
}

// 初始化
onMounted(() => {
  loadMessages()
})
</script>

<style scoped lang="scss">
.message-center {
  padding: 24px;
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  .header-content {
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
}

// 统计卡片
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  }

  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 14px;
    display: flex;
    align-items: center;
    justify-content: center;

    .el-icon {
      font-size: 24px;
      color: white;
    }

    &.all { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
    &.unread { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
    &.interview { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
    &.system { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }
  }

  .stat-info {
    display: flex;
    flex-direction: column;
    gap: 4px;

    .stat-value {
      font-size: 28px;
      font-weight: 700;
      color: #1a1a2e;
    }

    .stat-label {
      font-size: 14px;
      color: #6b7280;
    }
  }
}

// 内容区域
.content-wrapper {
  display: flex;
  gap: 24px;
}

// 左侧分类
.category-sidebar {
  width: 220px;
  flex-shrink: 0;

  .category-list {
    background: white;
    border-radius: 16px;
    padding: 12px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  }

  .category-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 16px;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
    position: relative;

    &:hover {
      background: #f5f7fa;
    }

    &.active {
      background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);

      .category-icon {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
      }

      .category-name {
        color: #667eea;
        font-weight: 600;
      }
    }

    .category-icon {
      width: 36px;
      height: 36px;
      border-radius: 10px;
      background: #f0f2f5;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: all 0.3s ease;

      .el-icon {
        font-size: 18px;
      }

      &.interview { background: rgba(79, 172, 254, 0.15); color: #4facfe; }
      &.application { background: rgba(240, 147, 251, 0.15); color: #f093fb; }
      &.notification { background: rgba(245, 87, 108, 0.15); color: #f5576c; }
      &.system { background: rgba(67, 233, 123, 0.15); color: #43e97b; }
    }

    .category-name {
      flex: 1;
      font-size: 14px;
      color: #4b5563;
    }

    .category-badge {
      :deep(.el-badge__content) {
        height: 18px;
        line-height: 18px;
        padding: 0 6px;
      }
    }
  }
}

// 右侧消息列表
.message-list-wrapper {
  flex: 1;
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f2f5;

  .toolbar-left {
    display: flex;
    align-items: center;
    gap: 16px;
  }
}

.message-list {
  min-height: 400px;

  .message-item {
    display: flex;
    align-items: flex-start;
    gap: 16px;
    padding: 20px;
    border-bottom: 1px solid #f0f2f5;
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      background: #fafbfc;
    }

    &.unread {
      background: linear-gradient(135deg, rgba(102, 126, 234, 0.03) 0%, rgba(118, 75, 162, 0.03) 100%);

      .message-title {
        font-weight: 600;
      }
    }

    &.selected {
      background: rgba(102, 126, 234, 0.08);
    }

    .el-checkbox {
      margin-top: 4px;
    }

    .message-icon {
      width: 44px;
      height: 44px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;

      .el-icon {
        font-size: 20px;
        color: white;
      }

      &.interview { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
      &.application { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
      &.notification { background: linear-gradient(135deg, #fa709a 0%, #fee140 100%); }
      &.system { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }
    }

    .message-content {
      flex: 1;
      min-width: 0;

      .message-header {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 8px;

        .message-title {
          font-size: 15px;
          color: #1a1a2e;
        }
      }

      .message-summary {
        font-size: 13px;
        color: #6b7280;
        margin: 0 0 12px 0;
        line-height: 1.6;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }

      .message-meta {
        display: flex;
        align-items: center;
        gap: 20px;
        font-size: 12px;
        color: #9ca3af;

        span {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }

    .message-actions {
      display: flex;
      flex-direction: column;
      gap: 4px;
      opacity: 0;
      transition: opacity 0.3s ease;
    }

    &:hover .message-actions {
      opacity: 1;
    }
  }
}

.empty-state {
  padding: 80px 20px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 20px;
  border-top: 1px solid #f0f2f5;
}

// 消息详情
.message-detail {
  .detail-header {
    display: flex;
    gap: 16px;
    margin-bottom: 24px;

    .detail-icon {
      width: 56px;
      height: 56px;
      border-radius: 14px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;

      .el-icon {
        font-size: 24px;
        color: white;
      }

      &.interview { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
      &.application { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
      &.notification { background: linear-gradient(135deg, #fa709a 0%, #fee140 100%); }
      &.system { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }
    }

    .detail-info {
      flex: 1;

      .detail-title {
        font-size: 18px;
        font-weight: 600;
        color: #1a1a2e;
        margin: 0 0 8px 0;
      }

      .detail-meta {
        display: flex;
        gap: 16px;
        font-size: 13px;
        color: #6b7280;

        span {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }
  }

  .detail-content {
    font-size: 14px;
    color: #4b5563;
    line-height: 1.8;
    margin-bottom: 24px;
  }

  .detail-extra {
    margin-bottom: 24px;

    h4 {
      font-size: 16px;
      font-weight: 600;
      color: #1a1a2e;
      margin: 0 0 16px 0;
    }
  }

  .detail-actions {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
    padding-top: 16px;
    border-top: 1px solid #f0f2f5;

    .el-button {
      border-radius: 10px;
    }
  }
}

// 响应式
@media (max-width: 1024px) {
  .content-wrapper {
    flex-direction: column;
  }

  .category-sidebar {
    width: 100%;

    .category-list {
      display: flex;
      overflow-x: auto;
      gap: 8px;
      padding: 12px;

      &::-webkit-scrollbar {
        display: none;
      }
    }

    .category-item {
      flex-shrink: 0;
      padding: 10px 14px;
    }
  }
}

@media (max-width: 768px) {
  .message-center {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .stats-row {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .stat-card {
    padding: 16px;

    .stat-icon {
      width: 48px;
      height: 48px;
    }

    .stat-info .stat-value {
      font-size: 24px;
    }
  }

  .toolbar {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;

    .toolbar-right {
      .el-input {
        width: 100% !important;
      }
    }
  }

  .message-list .message-item {
    padding: 16px;

    .message-actions {
      opacity: 1;
      flex-direction: row;
    }
  }
}
</style>
