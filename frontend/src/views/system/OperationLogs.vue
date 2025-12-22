<template>
  <div class="operation-logs">
    <div class="page-header">
      <div class="header-left">
        <h1>操作日志</h1>
        <p class="subtitle">查看系统操作记录和审计日志</p>
      </div>
      <div class="header-right">
        <el-button @click="exportLogs">
          <el-icon><Download /></el-icon>
          导出日志
        </el-button>
      </div>
    </div>

    <!-- 筛选条件 -->
    <div class="filter-card">
      <el-form :inline="true" :model="filterParams">
        <el-form-item>
          <el-input
            v-model="filterParams.search"
            placeholder="搜索操作内容..."
            clearable
            style="width: 220px"
            @keyup.enter="fetchLogs"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="filterParams.action" placeholder="操作类型" clearable style="width: 140px">
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
            <el-option label="登录" value="login" />
            <el-option label="登出" value="logout" />
            <el-option label="导出" value="export" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="filterParams.module" placeholder="功能模块" clearable style="width: 140px">
            <el-option label="用户管理" value="user" />
            <el-option label="人才管理" value="talent" />
            <el-option label="职位管理" value="job" />
            <el-option label="面试管理" value="interview" />
            <el-option label="简历管理" value="resume" />
            <el-option label="系统设置" value="system" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-date-picker
            v-model="filterParams.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 260px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchLogs">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 日志列表 -->
    <div class="logs-card">
      <el-table :data="logs" v-loading="loading" stripe>
        <el-table-column prop="time" label="操作时间" width="180">
          <template #default="{ row }">
            <span class="time-cell">{{ formatTime(row.time) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="user" label="操作用户" width="120">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="28" :style="{ background: getAvatarColor(row.userId) }">
                {{ row.user.charAt(0) }}
              </el-avatar>
              <span>{{ row.user }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)" size="small">
              {{ getActionLabel(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="module" label="功能模块" width="120">
          <template #default="{ row }">
            <span class="module-cell">
              <el-icon><component :is="getModuleIcon(row.module)" /></el-icon>
              {{ getModuleLabel(row.module) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="操作内容" min-width="250">
          <template #default="{ row }">
            <span class="content-cell">{{ row.content }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="showDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchLogs"
          @size-change="fetchLogs"
        />
      </div>
    </div>

    <!-- 详情弹窗 -->
    <el-dialog v-model="showDetailDialog" title="操作详情" width="600px">
      <div class="log-detail" v-if="currentLog">
        <div class="detail-row">
          <span class="label">操作时间</span>
          <span class="value">{{ formatTime(currentLog.time) }}</span>
        </div>
        <div class="detail-row">
          <span class="label">操作用户</span>
          <span class="value">{{ currentLog.user }}</span>
        </div>
        <div class="detail-row">
          <span class="label">操作类型</span>
          <el-tag :type="getActionType(currentLog.action)" size="small">
            {{ getActionLabel(currentLog.action) }}
          </el-tag>
        </div>
        <div class="detail-row">
          <span class="label">功能模块</span>
          <span class="value">{{ getModuleLabel(currentLog.module) }}</span>
        </div>
        <div class="detail-row">
          <span class="label">操作内容</span>
          <span class="value">{{ currentLog.content }}</span>
        </div>
        <div class="detail-row">
          <span class="label">IP地址</span>
          <span class="value">{{ currentLog.ip }}</span>
        </div>
        <div class="detail-row">
          <span class="label">浏览器</span>
          <span class="value">{{ currentLog.userAgent }}</span>
        </div>
        <div class="detail-row" v-if="currentLog.details">
          <span class="label">详细数据</span>
          <pre class="value code">{{ JSON.stringify(currentLog.details, null, 2) }}</pre>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, markRaw } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Search, Download, User, Suitcase, Document, Calendar, Setting, DataAnalysis
} from '@element-plus/icons-vue'

interface OperationLog {
  id: number
  time: string
  user: string
  userId: number
  action: string
  module: string
  content: string
  ip: string
  userAgent: string
  status: 'success' | 'failed'
  details?: any
}

const loading = ref(false)
const logs = ref<OperationLog[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showDetailDialog = ref(false)
const currentLog = ref<OperationLog | null>(null)

const filterParams = reactive({
  search: '',
  action: '',
  module: '',
  dateRange: null as [Date, Date] | null
})

const moduleIcons: Record<string, any> = {
  user: markRaw(User),
  talent: markRaw(User),
  job: markRaw(Suitcase),
  interview: markRaw(Calendar),
  resume: markRaw(Document),
  system: markRaw(Setting)
}

const fetchLogs = () => {
  loading.value = true
  // 模拟数据
  setTimeout(() => {
    logs.value = generateMockLogs()
    total.value = 256
    loading.value = false
  }, 500)
}

const generateMockLogs = (): OperationLog[] => {
  const users = ['admin', 'hr_zhang', 'hr_li', 'tech_chen']
  const actions = ['create', 'update', 'delete', 'login', 'export']
  const modules = ['user', 'talent', 'job', 'interview', 'resume']
  const contents = [
    '创建了新人才记录：张三',
    '更新了职位信息：高级Go开发工程师',
    '删除了简历记录 #123',
    '用户登录系统',
    '导出了人才列表数据',
    '安排了面试：李四 - 前端工程师',
    '更新了面试状态为已完成',
    '创建了新职位：产品经理'
  ]

  return Array.from({ length: pageSize.value }, (_, i) => ({
    id: (currentPage.value - 1) * pageSize.value + i + 1,
    time: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toISOString(),
    user: users[Math.floor(Math.random() * users.length)],
    userId: Math.floor(Math.random() * 10) + 1,
    action: actions[Math.floor(Math.random() * actions.length)],
    module: modules[Math.floor(Math.random() * modules.length)],
    content: contents[Math.floor(Math.random() * contents.length)],
    ip: `192.168.1.${Math.floor(Math.random() * 255)}`,
    userAgent: 'Chrome/120.0.0.0 Safari/537.36',
    status: Math.random() > 0.1 ? 'success' : 'failed',
    details: { oldValue: '旧值', newValue: '新值' }
  }))
}

const resetFilter = () => {
  filterParams.search = ''
  filterParams.action = ''
  filterParams.module = ''
  filterParams.dateRange = null
  fetchLogs()
}

const showDetail = (log: OperationLog) => {
  currentLog.value = log
  showDetailDialog.value = true
}

const exportLogs = () => {
  ElMessage.success('日志导出已开始')
}

const formatTime = (time: string) => {
  return new Date(time).toLocaleString('zh-CN')
}

const getAvatarColor = (id: number) => {
  const colors = ['#667eea', '#f093fb', '#4facfe', '#43e97b', '#f5576c']
  return colors[id % colors.length]
}

const getActionType = (action: string) => {
  const types: Record<string, any> = {
    create: 'success',
    update: 'warning',
    delete: 'danger',
    login: 'primary',
    logout: 'info',
    export: 'info'
  }
  return types[action] || 'info'
}

const getActionLabel = (action: string) => {
  const labels: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除',
    login: '登录',
    logout: '登出',
    export: '导出'
  }
  return labels[action] || action
}

const getModuleIcon = (module: string) => {
  return moduleIcons[module] || DataAnalysis
}

const getModuleLabel = (module: string) => {
  const labels: Record<string, string> = {
    user: '用户管理',
    talent: '人才管理',
    job: '职位管理',
    interview: '面试管理',
    resume: '简历管理',
    system: '系统设置'
  }
  return labels[module] || module
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped lang="scss">
.operation-logs {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  .header-left {
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
}

.filter-card, .logs-card {
  background: var(--bg-primary);
  border-radius: 12px;
  padding: 20px;
  box-shadow: var(--shadow-card);
  margin-bottom: 20px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.module-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-secondary);
}

.time-cell {
  font-size: 13px;
  color: var(--text-secondary);
}

.content-cell {
  color: var(--text-primary);
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.log-detail {
  .detail-row {
    display: flex;
    padding: 12px 0;
    border-bottom: 1px solid var(--border-light);

    &:last-child {
      border-bottom: none;
    }

    .label {
      width: 100px;
      color: var(--text-secondary);
      flex-shrink: 0;
    }

    .value {
      flex: 1;
      color: var(--text-primary);

      &.code {
        background: var(--bg-tertiary);
        padding: 12px;
        border-radius: 8px;
        font-family: monospace;
        font-size: 12px;
        overflow-x: auto;
      }
    }
  }
}
</style>
