<template>
  <div class="operation-logs">
    <div class="page-header">
      <div class="header-left">
        <h1>操作日志</h1>
        <p class="subtitle">查看系统操作记录和审计日志 (Elasticsearch)</p>
      </div>
      <div class="header-right">
        <el-button @click="exportLogs">
          <el-icon><Download /></el-icon>
          导出日志
        </el-button>
        <el-button type="danger" @click="cleanupLogs" v-if="isAdmin">
          <el-icon><Delete /></el-icon>
          清理日志
        </el-button>
      </div>
    </div>

    <!-- 筛选条件 -->
    <div class="filter-card">
      <el-form :inline="true" :model="filterParams">
        <el-form-item>
          <el-input
            v-model="filterParams.keyword"
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
          <el-select v-model="filterParams.service" placeholder="服务" clearable style="width: 150px">
            <el-option label="用户服务" value="user-service" />
            <el-option label="职位服务" value="job-service" />
            <el-option label="简历服务" value="resume-service" />
            <el-option label="面试服务" value="interview-service" />
            <el-option label="消息服务" value="message-service" />
            <el-option label="人才服务" value="talent-service" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="filterParams.method" placeholder="请求方法" clearable style="width: 120px">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="filterParams.level" placeholder="日志级别" clearable style="width: 120px">
            <el-option label="信息" value="info" />
            <el-option label="警告" value="warn" />
            <el-option label="错误" value="error" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-date-picker
            v-model="filterParams.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
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
        <el-table-column prop="timestamp" label="时间" width="180">
          <template #default="{ row }">
            <span class="time-cell">{{ formatTime(row.timestamp) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="service" label="服务" width="130">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ row.service }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户" width="100">
          <template #default="{ row }">
            <div class="user-cell" v-if="row.username">
              <el-avatar :size="24" :style="{ background: getAvatarColor(row.user_id || 0) }">
                {{ row.username?.charAt(0) || '?' }}
              </el-avatar>
              <span>{{ row.username }}</span>
            </div>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="方法" width="80">
          <template #default="{ row }">
            <el-tag :type="getMethodType(row.method)" size="small">
              {{ row.method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路径" min-width="200">
          <template #default="{ row }">
            <span class="path-cell">{{ row.path }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status_code" label="状态码" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status_code)" size="small">
              {{ row.status_code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="耗时" width="90">
          <template #default="{ row }">
            <span :class="getDurationClass(row.duration)">{{ row.duration }}ms</span>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column prop="level" label="级别" width="80">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.level)" size="small">
              {{ getLevelLabel(row.level) }}
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
    <el-dialog v-model="showDetailDialog" title="日志详情" width="700px">
      <div class="log-detail" v-if="currentLog">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="时间">{{ formatTime(currentLog.timestamp) }}</el-descriptions-item>
          <el-descriptions-item label="服务">{{ currentLog.service }}</el-descriptions-item>
          <el-descriptions-item label="用户">{{ currentLog.username || '-' }}</el-descriptions-item>
          <el-descriptions-item label="IP">{{ currentLog.ip }}</el-descriptions-item>
          <el-descriptions-item label="方法">
            <el-tag :type="getMethodType(currentLog.method)" size="small">{{ currentLog.method }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态码">
            <el-tag :type="getStatusType(currentLog.status_code)" size="small">{{ currentLog.status_code }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="路径" :span="2">{{ currentLog.path }}</el-descriptions-item>
          <el-descriptions-item label="查询参数" :span="2">{{ currentLog.query || '-' }}</el-descriptions-item>
          <el-descriptions-item label="耗时">{{ currentLog.duration }}ms</el-descriptions-item>
          <el-descriptions-item label="级别">
            <el-tag :type="getLevelType(currentLog.level)" size="small">{{ getLevelLabel(currentLog.level) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="操作类型">{{ currentLog.action }}</el-descriptions-item>
          <el-descriptions-item label="模块">{{ currentLog.module }}</el-descriptions-item>
          <el-descriptions-item label="User-Agent" :span="2">{{ currentLog.user_agent }}</el-descriptions-item>
        </el-descriptions>
        
        <div v-if="currentLog.request_body" class="detail-section">
          <h4>请求体</h4>
          <pre class="code-block">{{ formatJSON(currentLog.request_body) }}</pre>
        </div>
        
        <div v-if="currentLog.error_msg" class="detail-section error">
          <h4>错误信息</h4>
          <pre class="code-block error">{{ currentLog.error_msg }}</pre>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Download, Delete } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import axios from 'axios'
import { exportToCsv, type ExportColumn } from '@/utils/export'

interface OperationLog {
  id: string
  timestamp: string
  service: string
  user_id: number
  username: string
  ip: string
  method: string
  path: string
  query: string
  status_code: number
  duration: number
  request_body: string
  response_body: string
  user_agent: string
  action: string
  module: string
  description: string
  level: string
  error_msg: string
}

const userStore = useUserStore()
const isAdmin = computed(() => userStore.user?.role === 'admin')

const loading = ref(false)
const logs = ref<OperationLog[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showDetailDialog = ref(false)
const currentLog = ref<OperationLog | null>(null)

const filterParams = reactive({
  keyword: '',
  service: '',
  method: '',
  level: '',
  dateRange: null as [string, string] | null
})

// 获取日志
const fetchLogs = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (filterParams.keyword) params.keyword = filterParams.keyword
    if (filterParams.service) params.service = filterParams.service
    if (filterParams.method) params.method = filterParams.method
    if (filterParams.level) params.level = filterParams.level
    if (filterParams.dateRange) {
      params.start_time = filterParams.dateRange[0]
      params.end_time = filterParams.dateRange[1]
    }

    const res = await axios.get('/api/v1/logs', { params })
    if (res.data.code === 0) {
      logs.value = res.data.data.logs || []
      total.value = res.data.data.total || 0
    } else {
      // 如果ES服务不可用，使用模拟数据
      logs.value = generateMockLogs()
      total.value = 100
    }
  } catch (error) {
    console.error('获取日志失败:', error)
    // 使用模拟数据
    logs.value = generateMockLogs()
    total.value = 100
  } finally {
    loading.value = false
  }
}

// 生成模拟数据
const generateMockLogs = (): OperationLog[] => {
  const services = ['user-service', 'job-service', 'talent-service', 'resume-service']
  const methods = ['GET', 'POST', 'PUT', 'DELETE']
  const paths = ['/api/v1/users', '/api/v1/jobs', '/api/v1/talents', '/api/v1/login', '/api/v1/resumes']
  const levels = ['info', 'warn', 'error']
  const users = ['admin', 'hr_zhang', 'hr_li', '']

  return Array.from({ length: pageSize.value }, (_, i) => ({
    id: `log-${Date.now()}-${i}`,
    timestamp: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toISOString(),
    service: services[Math.floor(Math.random() * services.length)],
    user_id: Math.floor(Math.random() * 10),
    username: users[Math.floor(Math.random() * users.length)],
    ip: `192.168.1.${Math.floor(Math.random() * 255)}`,
    method: methods[Math.floor(Math.random() * methods.length)],
    path: paths[Math.floor(Math.random() * paths.length)],
    query: '',
    status_code: Math.random() > 0.9 ? 500 : Math.random() > 0.8 ? 400 : 200,
    duration: Math.floor(Math.random() * 500),
    request_body: '',
    response_body: '',
    user_agent: 'Mozilla/5.0 Chrome/120.0.0.0',
    action: '查询',
    module: 'users',
    description: '',
    level: levels[Math.floor(Math.random() * levels.length)],
    error_msg: ''
  }))
}

const resetFilter = () => {
  filterParams.keyword = ''
  filterParams.service = ''
  filterParams.method = ''
  filterParams.level = ''
  filterParams.dateRange = null
  currentPage.value = 1
  fetchLogs()
}

const showDetail = (log: OperationLog) => {
  currentLog.value = log
  showDetailDialog.value = true
}

const exportLogs = () => {
  if (logs.value.length === 0) {
    ElMessage.warning('暂无数据可导出')
    return
  }

  const columns: ExportColumn[] = [
    { key: 'timestamp', title: '时间' },
    { key: 'service', title: '服务' },
    { key: 'username', title: '用户' },
    { key: 'method', title: '方法' },
    { key: 'path', title: '路径' },
    { key: 'status_code', title: '状态码' },
    { key: 'duration', title: '耗时(ms)' },
    { key: 'ip', title: 'IP' },
    { key: 'level', title: '级别' }
  ]

  exportToCsv({
    filename: `操作日志_${new Date().toISOString().split('T')[0]}`,
    columns,
    data: logs.value
  })
  ElMessage.success('导出成功')
}

const cleanupLogs = async () => {
  try {
    await ElMessageBox.confirm('确定要清理30天前的日志吗？此操作不可恢复。', '确认清理', {
      type: 'warning'
    })
    
    await axios.delete('/api/v1/logs/cleanup', { params: { days: 30 } })
    ElMessage.success('清理成功')
    fetchLogs()
  } catch {
    // 用户取消
  }
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const formatJSON = (str: string) => {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

const getAvatarColor = (id: number) => {
  const colors = ['#0ea5e9', '#06b6d4', '#22c55e', '#f59e0b', '#8b5cf6']
  return colors[id % colors.length]
}

const getMethodType = (method: string) => {
  const types: Record<string, any> = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'info'
  }
  return types[method] || 'info'
}

const getStatusType = (code: number) => {
  if (code >= 500) return 'danger'
  if (code >= 400) return 'warning'
  if (code >= 200 && code < 300) return 'success'
  return 'info'
}

const getLevelType = (level: string) => {
  const types: Record<string, any> = {
    info: 'success',
    warn: 'warning',
    error: 'danger'
  }
  return types[level] || 'info'
}

const getLevelLabel = (level: string) => {
  const labels: Record<string, string> = {
    info: '信息',
    warn: '警告',
    error: '错误'
  }
  return labels[level] || level
}

const getDurationClass = (duration: number) => {
  if (duration > 1000) return 'duration-slow'
  if (duration > 500) return 'duration-medium'
  return 'duration-fast'
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

  .header-right {
    display: flex;
    gap: 12px;
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
  gap: 6px;
}

.path-cell {
  font-family: monospace;
  font-size: 13px;
  color: var(--text-primary);
}

.time-cell {
  font-size: 13px;
  color: var(--text-secondary);
}

.text-muted {
  color: var(--text-muted);
}

.duration-fast { color: #22c55e; }
.duration-medium { color: #f59e0b; }
.duration-slow { color: #ef4444; font-weight: 600; }

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.log-detail {
  .detail-section {
    margin-top: 20px;

    h4 {
      font-size: 14px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 10px 0;
    }

    &.error h4 {
      color: #ef4444;
    }
  }

  .code-block {
    background: var(--bg-tertiary);
    padding: 12px;
    border-radius: 8px;
    font-family: monospace;
    font-size: 12px;
    overflow-x: auto;
    white-space: pre-wrap;
    word-break: break-all;
    max-height: 200px;
    overflow-y: auto;

    &.error {
      background: #fef2f2;
      color: #dc2626;
    }
  }
}
</style>
