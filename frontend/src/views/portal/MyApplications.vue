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
      <div class="application-list">
        <div class="application-item" v-for="app in filteredApplications" :key="app.id">
          <div class="app-main">
            <div class="job-info">
              <h3>{{ app.jobTitle }}</h3>
              <p class="company">{{ app.company }} · {{ app.location }}</p>
              <p class="salary">{{ app.salary }}</p>
            </div>
            <div class="app-status">
              <el-tag :type="getStatusType(app.status)">{{ getStatusText(app.status) }}</el-tag>
              <span class="apply-time">{{ app.applyTime }}</span>
            </div>
          </div>
          <div class="app-timeline" v-if="app.timeline.length > 0">
            <div class="timeline-item" v-for="(item, index) in app.timeline" :key="index">
              <span class="timeline-dot" :class="{ active: index === 0 }"></span>
              <span class="timeline-content">{{ item.content }}</span>
              <span class="timeline-time">{{ item.time }}</span>
            </div>
          </div>
          <div class="app-actions">
            <el-button size="small" @click="viewJob(app.jobId)">查看职位</el-button>
            <el-button size="small" type="danger" plain @click="withdrawApplication(app.id)">撤回投递</el-button>
          </div>
        </div>

        <el-empty v-if="filteredApplications.length === 0" description="暂无投递记录" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const statusFilter = ref('')

const applications = ref([
  {
    id: 1,
    jobId: 1,
    jobTitle: '高级前端工程师',
    company: '科技有限公司',
    location: '北京',
    salary: '25-45K',
    status: 'interview',
    applyTime: '2024-01-15',
    timeline: [
      { content: '收到面试邀请', time: '01-18 10:30' },
      { content: 'HR已查看简历', time: '01-16 14:20' },
      { content: '投递成功', time: '01-15 09:00' }
    ]
  },
  {
    id: 2,
    jobId: 2,
    jobTitle: '后端开发工程师',
    company: '互联网公司',
    location: '上海',
    salary: '30-50K',
    status: 'viewed',
    applyTime: '2024-01-14',
    timeline: [
      { content: 'HR已查看简历', time: '01-15 11:00' },
      { content: '投递成功', time: '01-14 16:30' }
    ]
  },
  {
    id: 3,
    jobId: 3,
    jobTitle: '产品经理',
    company: '创新科技',
    location: '深圳',
    salary: '20-35K',
    status: 'pending',
    applyTime: '2024-01-13',
    timeline: [
      { content: '投递成功', time: '01-13 10:00' }
    ]
  },
  {
    id: 4,
    jobId: 4,
    jobTitle: 'UI设计师',
    company: '设计工作室',
    location: '杭州',
    salary: '15-25K',
    status: 'rejected',
    applyTime: '2024-01-10',
    timeline: [
      { content: '不合适', time: '01-12 09:00' },
      { content: 'HR已查看简历', time: '01-11 14:00' },
      { content: '投递成功', time: '01-10 11:30' }
    ]
  }
])

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

const filterApplications = () => {}

const viewJob = (jobId: number) => {
  router.push(`/portal/jobs/${jobId}`)
}

const withdrawApplication = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要撤回这个投递吗？', '撤回投递', { type: 'warning' })
    applications.value = applications.value.filter(a => a.id !== id)
    ElMessage.success('已撤回')
  } catch {}
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
