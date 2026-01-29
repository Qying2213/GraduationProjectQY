<template>
  <div class="my-applications">
    <div class="page-container">
      <div class="page-header">
        <h1>我的投递</h1>
        <p>共投递 {{ totalCount }} 个职位</p>
      </div>

      <!-- 状态筛选 Tabs -->
      <div class="filter-tabs">
        <el-tabs v-model="statusFilter" @tab-change="handleStatusChange">
          <el-tab-pane label="全部" name="">
            <template #label>
              <span>全部 <el-badge :value="totalCount" :max="99" type="info" /></span>
            </template>
          </el-tab-pane>
          <el-tab-pane label="待处理" name="pending">
            <template #label>
              <span>待处理 <el-badge :value="pendingCount" :max="99" type="warning" v-if="pendingCount > 0" /></span>
            </template>
          </el-tab-pane>
          <el-tab-pane label="已查看" name="viewed">
            <template #label>
              <span>已查看 <el-badge :value="viewedCount" :max="99" type="primary" v-if="viewedCount > 0" /></span>
            </template>
          </el-tab-pane>
          <el-tab-pane label="面试中" name="interviewing">
            <template #label>
              <span>面试中 <el-badge :value="interviewingCount" :max="99" type="success" v-if="interviewingCount > 0" /></span>
            </template>
          </el-tab-pane>
          <el-tab-pane label="已录用" name="accepted">
            <template #label>
              <span>已录用 <el-badge :value="acceptedCount" :max="99" type="success" v-if="acceptedCount > 0" /></span>
            </template>
          </el-tab-pane>
          <el-tab-pane label="已拒绝" name="rejected">
            <template #label>
              <span>已拒绝 <el-badge :value="rejectedCount" :max="99" type="danger" v-if="rejectedCount > 0" /></span>
            </template>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 投递列表 -->
      <div class="application-list" v-loading="loading">
        <div class="application-item" v-for="app in filteredApplications" :key="app.id">
          <div class="app-main">
            <div class="job-info">
              <h3 @click="viewJob(app.job_id)" class="job-title-link">{{ app.job_title }}</h3>
              <p class="company">{{ app.company_name }} · {{ app.location }}</p>
              <p class="salary">{{ app.salary }}</p>
            </div>
            <div class="app-status">
              <el-tag :type="getStatusType(app.status)" size="large">{{ getStatusText(app.status) }}</el-tag>
              <span class="apply-time">投递于 {{ formatDate(app.created_at) }}</span>
            </div>
          </div>
          
          <!-- 状态时间线 -->
          <div class="app-timeline" v-if="app.timeline && app.timeline.length > 0">
            <el-timeline>
              <el-timeline-item
                v-for="(item, index) in app.timeline"
                :key="index"
                :type="getTimelineItemType(item.status)"
                :hollow="index !== 0"
                :timestamp="item.time"
                placement="top"
              >
                <div class="timeline-content">
                  <span class="timeline-title">{{ item.content }}</span>
                  <span class="timeline-desc" v-if="item.description">{{ item.description }}</span>
                </div>
              </el-timeline-item>
            </el-timeline>
          </div>
          
          <div class="app-actions">
            <el-button size="small" @click="viewJob(app.job_id)">查看职位</el-button>
            <el-button 
              v-if="canWithdraw(app.status)" 
              size="small" 
              type="danger" 
              plain 
              @click="handleWithdraw(app)"
            >
              撤回投递
            </el-button>
            <el-tag v-if="app.status === 'withdrawn'" type="info" size="small">已撤回</el-tag>
          </div>
        </div>

        <el-empty v-if="filteredApplications.length === 0 && !loading" description="暂无投递记录">
          <el-button type="primary" @click="goToJobs">去看看职位</el-button>
        </el-empty>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="total > pageSize">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next, jumper"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 撤回确认对话框 -->
    <el-dialog
      v-model="withdrawDialogVisible"
      title="撤回投递"
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="withdraw-dialog-content">
        <el-icon class="warning-icon"><WarningFilled /></el-icon>
        <p>确定要撤回对 <strong>{{ withdrawingApp?.job_title }}</strong> 的投递吗？</p>
        <p class="warning-text">撤回后，该投递记录将被删除，您可以重新投递该职位。</p>
      </div>
      <template #footer>
        <el-button @click="withdrawDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmWithdraw" :loading="withdrawing">确认撤回</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { WarningFilled } from '@element-plus/icons-vue'
import { applicationApi } from '@/api/application'

const router = useRouter()
const statusFilter = ref('')
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 撤回相关状态
const withdrawDialogVisible = ref(false)
const withdrawingApp = ref<ApplicationItem | null>(null)
const withdrawing = ref(false)

// 状态映射
const STATUS_MAP: Record<string, { text: string; type: string }> = {
  pending: { text: '待处理', type: 'warning' },
  viewed: { text: '已查看', type: 'primary' },
  interviewing: { text: '面试中', type: 'success' },
  accepted: { text: '已录用', type: 'success' },
  rejected: { text: '已拒绝', type: 'danger' },
  withdrawn: { text: '已撤回', type: 'info' }
}

interface TimelineItem {
  content: string
  time: string
  status: string
  description?: string
}

interface ApplicationItem {
  id: number
  job_id: number
  job_title: string
  company_name: string
  location: string
  salary: string
  status: string
  created_at: string
  updated_at: string
  viewed_at?: string
  interview_time?: string
  timeline: TimelineItem[]
}

const applications = ref<ApplicationItem[]>([])

// 加载我的申请列表
const loadApplications = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    // 如果有状态筛选，添加到参数中
    if (statusFilter.value) {
      params.status = statusFilter.value
    }
    
    const res = await applicationApi.list(params)
    
    if (res.data?.code === 0) {
      const data = res.data.data
      const appList = data?.applications || data || []
      total.value = data?.total || appList.length
      
      applications.value = appList.map((app: any) => ({
        id: app.id,
        job_id: app.job_id,
        job_title: app.job_title || app.job?.title || '未知职位',
        company_name: app.company_name || app.job?.company_name || '未知公司',
        location: app.location || app.job?.location || '',
        salary: app.salary || app.job?.salary || '',
        status: app.status || 'pending',
        created_at: app.created_at,
        updated_at: app.updated_at,
        viewed_at: app.viewed_at,
        interview_time: app.interview_time,
        timeline: buildTimeline(app)
      }))
    }
  } catch (error) {
    console.error('加载申请列表失败:', error)
    ElMessage.error('加载申请列表失败')
  } finally {
    loading.value = false
  }
}

// 构建状态时间线
const buildTimeline = (app: any): TimelineItem[] => {
  const timeline: TimelineItem[] = []
  const status = app.status || 'pending'
  
  // 根据当前状态和历史记录构建时间线
  // 最新的状态在最上面
  
  if (status === 'accepted') {
    timeline.push({
      content: '恭喜！您已被录用',
      time: formatDateTime(app.accepted_at || app.updated_at),
      status: 'accepted',
      description: '请等待HR联系您办理入职手续'
    })
  }
  
  if (status === 'rejected') {
    timeline.push({
      content: '很遗憾，未能通过筛选',
      time: formatDateTime(app.rejected_at || app.updated_at),
      status: 'rejected',
      description: '感谢您的投递，祝您求职顺利'
    })
  }
  
  if (status === 'withdrawn') {
    timeline.push({
      content: '您已撤回投递',
      time: formatDateTime(app.withdrawn_at || app.updated_at),
      status: 'withdrawn'
    })
  }
  
  if (status === 'interviewing' || ['accepted', 'rejected'].includes(status)) {
    timeline.push({
      content: '收到面试邀请',
      time: formatDateTime(app.interview_time || app.interviewing_at || app.updated_at),
      status: 'interviewing',
      description: app.interview_location ? `面试地点：${app.interview_location}` : undefined
    })
  }
  
  if (['viewed', 'interviewing', 'accepted', 'rejected'].includes(status)) {
    timeline.push({
      content: 'HR已查看您的简历',
      time: formatDateTime(app.viewed_at || app.updated_at),
      status: 'viewed'
    })
  }
  
  // 投递成功始终显示在最下面
  timeline.push({
    content: '投递成功',
    time: formatDateTime(app.created_at),
    status: 'pending',
    description: '等待HR查看'
  })
  
  return timeline
}

// 格式化日期时间
const formatDateTime = (dateStr: string): string => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 格式化日期（仅日期部分）
const formatDate = (dateStr: string): string => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

onMounted(() => {
  loadApplications()
})

// 计算属性：根据状态筛选后的申请列表
const filteredApplications = computed(() => {
  // 如果已经在API层面筛选了，直接返回
  if (statusFilter.value) {
    return applications.value
  }
  return applications.value
})

// 各状态数量统计
const totalCount = computed(() => total.value)
const pendingCount = computed(() => applications.value.filter(a => a.status === 'pending').length)
const viewedCount = computed(() => applications.value.filter(a => a.status === 'viewed').length)
const interviewingCount = computed(() => applications.value.filter(a => a.status === 'interviewing').length)
const acceptedCount = computed(() => applications.value.filter(a => a.status === 'accepted').length)
const rejectedCount = computed(() => applications.value.filter(a => a.status === 'rejected').length)

// 获取状态标签类型
const getStatusType = (status: string): string => {
  return STATUS_MAP[status]?.type || 'info'
}

// 获取状态文本
const getStatusText = (status: string): string => {
  return STATUS_MAP[status]?.text || status
}

// 获取时间线项目类型
const getTimelineItemType = (status: string): string => {
  const typeMap: Record<string, string> = {
    pending: 'primary',
    viewed: 'primary',
    interviewing: 'success',
    accepted: 'success',
    rejected: 'danger',
    withdrawn: 'info'
  }
  return typeMap[status] || 'primary'
}

// 判断是否可以撤回
const canWithdraw = (status: string): boolean => {
  // 只有待处理和已查看状态可以撤回
  return ['pending', 'viewed'].includes(status)
}

// 处理状态筛选变化
const handleStatusChange = () => {
  currentPage.value = 1
  loadApplications()
}

// 处理分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadApplications()
}

// 查看职位详情
const viewJob = (jobId: number) => {
  router.push(`/portal/jobs/${jobId}`)
}

// 跳转到职位列表
const goToJobs = () => {
  router.push('/portal/jobs')
}

// 处理撤回操作
const handleWithdraw = (app: ApplicationItem) => {
  withdrawingApp.value = app
  withdrawDialogVisible.value = true
}

// 确认撤回
const confirmWithdraw = async () => {
  if (!withdrawingApp.value) return
  
  withdrawing.value = true
  try {
    const res = await applicationApi.delete(withdrawingApp.value.id)
    
    if (res.data?.code === 0) {
      ElMessage.success('投递已撤回')
      withdrawDialogVisible.value = false
      
      // 从列表中移除该申请
      applications.value = applications.value.filter(a => a.id !== withdrawingApp.value?.id)
      total.value = Math.max(0, total.value - 1)
      
      // 如果当前页没有数据了，回到上一页
      if (applications.value.length === 0 && currentPage.value > 1) {
        currentPage.value--
        loadApplications()
      }
    } else {
      ElMessage.error(res.data?.message || '撤回失败')
    }
  } catch (error: any) {
    console.error('撤回失败:', error)
    ElMessage.error(error.response?.data?.message || '撤回失败，请稍后重试')
  } finally {
    withdrawing.value = false
    withdrawingApp.value = null
  }
}
</script>

<style scoped lang="scss">
.my-applications {
  padding: 24px;
  background: #f8fafc;
  min-height: calc(100vh - 160px);

  .page-container {
    max-width: 900px;
    margin: 0 auto;
  }

  .page-header {
    margin-bottom: 24px;

    h1 {
      font-size: 24px;
      font-weight: 700;
      color: #1e293b;
      margin: 0 0 4px 0;
    }

    p {
      color: #64748b;
      margin: 0;
    }
  }

  .filter-tabs {
    margin-bottom: 24px;
    background: white;
    border-radius: 12px;
    padding: 8px 16px;
    border: 1px solid #e2e8f0;

    :deep(.el-tabs__header) {
      margin: 0;
    }

    :deep(.el-tabs__nav-wrap::after) {
      display: none;
    }

    :deep(.el-tabs__item) {
      padding: 0 20px;
      height: 44px;
      line-height: 44px;
    }

    :deep(.el-badge__content) {
      margin-left: 6px;
    }
  }

  .application-list {
    .application-item {
      background: white;
      border-radius: 12px;
      padding: 20px;
      margin-bottom: 16px;
      border: 1px solid #e2e8f0;
      transition: box-shadow 0.2s ease;

      &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
      }

      .app-main {
        display: flex;
        justify-content: space-between;
        margin-bottom: 16px;

        .job-info {
          h3.job-title-link {
            font-size: 18px;
            font-weight: 600;
            color: #1e293b;
            margin: 0 0 8px 0;
            cursor: pointer;
            transition: color 0.2s ease;

            &:hover {
              color: #0ea5e9;
            }
          }

          .company {
            color: #64748b;
            margin: 0 0 4px 0;
            font-size: 14px;
          }

          .salary {
            color: #0ea5e9;
            font-weight: 600;
            margin: 0;
            font-size: 16px;
          }
        }

        .app-status {
          text-align: right;
          display: flex;
          flex-direction: column;
          align-items: flex-end;
          gap: 8px;

          .apply-time {
            font-size: 12px;
            color: #94a3b8;
          }
        }
      }

      .app-timeline {
        background: #f8fafc;
        border-radius: 8px;
        padding: 16px 20px;
        margin-bottom: 16px;

        :deep(.el-timeline) {
          padding-left: 0;
        }

        :deep(.el-timeline-item) {
          padding-bottom: 12px;

          &:last-child {
            padding-bottom: 0;
          }
        }

        :deep(.el-timeline-item__wrapper) {
          padding-left: 20px;
        }

        :deep(.el-timeline-item__timestamp) {
          color: #94a3b8;
          font-size: 12px;
        }

        .timeline-content {
          display: flex;
          flex-direction: column;
          gap: 4px;

          .timeline-title {
            color: #475569;
            font-size: 14px;
            font-weight: 500;
          }

          .timeline-desc {
            color: #94a3b8;
            font-size: 12px;
          }
        }
      }

      .app-actions {
        display: flex;
        gap: 12px;
        padding-top: 12px;
        border-top: 1px solid #f1f5f9;
      }
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 24px;
    padding: 16px;
    background: white;
    border-radius: 12px;
    border: 1px solid #e2e8f0;
  }
}

// 撤回对话框样式
.withdraw-dialog-content {
  text-align: center;
  padding: 20px 0;

  .warning-icon {
    font-size: 48px;
    color: #f59e0b;
    margin-bottom: 16px;
  }

  p {
    margin: 0 0 8px 0;
    color: #1e293b;
    font-size: 16px;

    strong {
      color: #0ea5e9;
    }
  }

  .warning-text {
    color: #94a3b8;
    font-size: 14px;
  }
}
</style>
