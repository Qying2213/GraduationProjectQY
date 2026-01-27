<template>
  <div class="my-applications">
    <div class="page-container">
      <div class="page-header">
        <h1>我的投递</h1>
        <p>共投递 {{ applications.length }} 个职位</p>
      </div>

      <!-- 状态筛选 -->
      <div class="filter-tabs">
        <el-radio-group v-model="statusFilter" @change="filterApplications">
          <el-radio-button value="">全部 ({{ applications.length }})</el-radio-button>
          <el-radio-button value="pending">待查看 ({{ pendingCount }})</el-radio-button>
          <el-radio-button value="viewed">已查看 ({{ viewedCount }})</el-radio-button>
          <el-radio-button value="interview">面试邀请 ({{ interviewCount }})</el-radio-button>
          <el-radio-button value="rejected">不合适 ({{ rejectedCount }})</el-radio-button>
        </el-radio-group>
      </div>

      <!-- 投递列表 -->
      <div class="application-list" v-loading="loading">
        <div class="application-item" v-for="app in filteredApplications" :key="app.id">
          <div class="app-main">
            <div class="job-info">
              <h3>{{ app.job_title }}</h3>
              <p class="company">{{ app.company_name }} · {{ app.location }}</p>
              <p class="salary">{{ app.salary }}</p>
            </div>
            <div class="app-status">
              <el-tag :type="getStatusType(app.status)">{{ getStatusText(app.status) }}</el-tag>
              <span class="apply-time">{{ app.created_at?.split('T')[0] }}</span>
            </div>
          </div>
          <div class="app-timeline" v-if="app.timeline && app.timeline.length > 0">
            <div class="timeline-item" v-for="(item, index) in app.timeline" :key="index">
              <span class="timeline-dot" :class="{ active: index === 0 }"></span>
              <span class="timeline-content">{{ item.content }}</span>
              <span class="timeline-time">{{ item.time }}</span>
            </div>
          </div>
          <div class="app-actions">
            <el-button size="small" @click="viewJob(app.job_id)">查看职位</el-button>
            <el-button size="small" type="danger" plain @click="withdrawApplication(app.id)">撤回投递</el-button>
          </div>
        </div>

        <el-empty v-if="filteredApplications.length === 0 && !loading" description="暂无投递记录" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const statusFilter = ref('')
const loading = ref(false)

interface Application {
  id: number
  job_id: number
  job_title?: string
  company_name?: string
  location?: string
  salary?: string
  status: string
  created_at: string
  timeline?: { content: string; time: string }[]
}

const applications = ref<Application[]>([])

// 加载我的申请列表
const loadApplications = async () => {
  loading.value = true
  try {
    // 使用 /applications 接口，后端会根据 JWT 中的用户信息返回对应的申请
    const res = await request.get('/applications', { params: { page_size: 100 } })
    if (res.data?.code === 0) {
      applications.value = (res.data.data?.applications || res.data.data || []).map((app: any) => ({
        id: app.id,
        job_id: app.job_id,
        job_title: app.job_title || app.job?.title || '未知职位',
        company_name: app.company_name || app.job?.company_name || '未知公司',
        location: app.location || app.job?.location || '',
        salary: app.salary || app.job?.salary || '',
        status: app.status || 'pending',
        created_at: app.created_at,
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

// 构建时间线
const buildTimeline = (app: any) => {
  const timeline = []
  if (app.status === 'interview') {
    timeline.push({ content: '收到面试邀请', time: formatTime(app.interview_time || app.updated_at) })
  }
  if (app.status === 'rejected') {
    timeline.push({ content: '不合适', time: formatTime(app.updated_at) })
  }
  if (app.viewed_at) {
    timeline.push({ content: 'HR已查看简历', time: formatTime(app.viewed_at) })
  }
  timeline.push({ content: '投递成功', time: formatTime(app.created_at) })
  return timeline
}

const formatTime = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

onMounted(() => {
  loadApplications()
})

const filteredApplications = computed(() => {
  if (!statusFilter.value) return applications.value
  return applications.value.filter(app => app.status === statusFilter.value)
})

const pendingCount = computed(() => applications.value.filter(a => a.status === 'pending').length)
const viewedCount = computed(() => applications.value.filter(a => a.status === 'viewed').length)
const interviewCount = computed(() => applications.value.filter(a => a.status === 'interview').length)
const rejectedCount = computed(() => applications.value.filter(a => a.status === 'rejected').length)

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'info',
    viewed: 'warning',
    interview: 'success',
    rejected: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待查看',
    viewed: '已查看',
    interview: '面试邀请',
    rejected: '不合适'
  }
  return map[status] || status
}

const filterApplications = () => {
  loadApplications()
}

const viewJob = (jobId: number) => {
  router.push(`/portal/jobs/${jobId}`)
}

const withdrawApplication = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要撤回这个投递吗？', '撤回投递', { type: 'warning' })
    const res = await request.delete(`/applications/${id}`)
    if (res.data?.code === 0) {
      applications.value = applications.value.filter(a => a.id !== id)
      ElMessage.success('已撤回')
    } else {
      ElMessage.error(res.data?.message || '撤回失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('撤回失败')
    }
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
  }

  .application-list {
    .application-item {
      background: white;
      border-radius: 12px;
      padding: 20px;
      margin-bottom: 16px;
      border: 1px solid #e2e8f0;

      .app-main {
        display: flex;
        justify-content: space-between;
        margin-bottom: 16px;

        .job-info {
          h3 {
            font-size: 18px;
            font-weight: 600;
            color: #1e293b;
            margin: 0 0 8px 0;
          }

          .company {
            color: #64748b;
            margin: 0 0 4px 0;
          }

          .salary {
            color: #0ea5e9;
            font-weight: 600;
            margin: 0;
          }
        }

        .app-status {
          text-align: right;

          .apply-time {
            display: block;
            margin-top: 8px;
            font-size: 12px;
            color: #94a3b8;
          }
        }
      }

      .app-timeline {
        background: #f8fafc;
        border-radius: 8px;
        padding: 16px;
        margin-bottom: 16px;

        .timeline-item {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 8px 0;

          &:not(:last-child) {
            border-bottom: 1px dashed #e2e8f0;
          }

          .timeline-dot {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background: #cbd5e1;

            &.active { background: #0ea5e9; }
          }

          .timeline-content {
            flex: 1;
            color: #475569;
          }

          .timeline-time {
            font-size: 12px;
            color: #94a3b8;
          }
        }
      }

      .app-actions {
        display: flex;
        gap: 12px;
      }
    }
  }
}
</style>
