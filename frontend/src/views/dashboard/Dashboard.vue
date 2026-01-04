<template>
  <div class="dashboard">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">仪表板</h1>
        <p class="page-subtitle">欢迎回来，{{ userStore.user?.username || '用户' }}！这是您的数据概览</p>
      </div>
      <div class="header-right">
        <el-button type="success" @click="$router.push('/data-screen')">
          <el-icon><DataAnalysis /></el-icon>
          数据大屏
        </el-button>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          size="default"
          style="width: 260px"
        />
        <el-button type="primary" :icon="Refresh" @click="refreshData">刷新数据</el-button>
      </div>
    </div>

    <!-- Stats Cards -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :lg="6" v-for="(stat, index) in statsCards" :key="index">
        <div class="stat-card" :style="{ animationDelay: `${index * 0.1}s` }">
          <div class="stat-card-inner" :class="stat.colorClass">
            <div class="stat-icon-wrapper">
              <div class="stat-icon">
                <el-icon :size="28"><component :is="stat.icon" /></el-icon>
              </div>
              <div class="stat-trend" :class="stat.trend > 0 ? 'up' : 'down'" v-if="stat.trend !== 0">
                <el-icon><component :is="stat.trend > 0 ? ArrowUp : ArrowDown" /></el-icon>
                <span>{{ Math.abs(stat.trend) }}%</span>
              </div>
            </div>
            <div class="stat-content">
              <div class="stat-value">
                <span class="number" ref="statNumbers">{{ stat.displayValue }}</span>
                <span class="suffix" v-if="stat.suffix">{{ stat.suffix }}</span>
              </div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
            <div class="stat-bg-icon">
              <el-icon :size="100"><component :is="stat.icon" /></el-icon>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Charts Row -->
    <el-row :gutter="20" class="charts-row">
      <el-col :xs="24" :lg="16">
        <div class="card chart-card">
          <div class="card-header">
            <h3>数据趋势</h3>
            <el-radio-group v-model="trendType" size="small">
              <el-radio-button label="week">本周</el-radio-button>
              <el-radio-button label="month">本月</el-radio-button>
              <el-radio-button label="year">全年</el-radio-button>
            </el-radio-group>
          </div>
          <div ref="trendChartRef" style="width: 100%; height: 350px"></div>
        </div>
      </el-col>

      <el-col :xs="24" :lg="8">
        <div class="card chart-card">
          <div class="card-header">
            <h3>职位状态分布</h3>
            <el-dropdown trigger="click">
              <el-button text :icon="MoreFilled" />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item>导出图表</el-dropdown-item>
                  <el-dropdown-item>查看详情</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <div ref="jobChartRef" style="width: 100%; height: 350px"></div>
        </div>
      </el-col>
    </el-row>

    <!-- Quick Actions & Recent Activities -->
    <el-row :gutter="20" class="bottom-row">
      <el-col :xs="24" :lg="8">
        <div class="card quick-actions-card">
          <div class="card-header">
            <h3>快捷操作</h3>
          </div>
          <div class="quick-actions">
            <div class="action-item" v-for="(action, index) in quickActions" :key="index"
                 @click="handleAction(action.route)">
              <div class="action-icon" :style="{ background: action.gradient }">
                <el-icon :size="24"><component :is="action.icon" /></el-icon>
              </div>
              <div class="action-info">
                <span class="action-title">{{ action.title }}</span>
                <span class="action-desc">{{ action.desc }}</span>
              </div>
              <el-icon class="action-arrow"><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :lg="16">
        <div class="card recent-activities-card">
          <div class="card-header">
            <h3>最近活动</h3>
            <el-button type="primary" size="small" @click="$router.push('/messages')">查看全部</el-button>
          </div>
          <div class="activities-list">
            <div class="activity-item" v-for="(activity, index) in recentActivities" :key="index">
              <div class="activity-avatar" :style="{ background: activity.color }">
                <el-icon><component :is="activity.icon" /></el-icon>
              </div>
              <div class="activity-content">
                <div class="activity-title">{{ activity.title }}</div>
                <div class="activity-desc">{{ activity.description }}</div>
              </div>
              <div class="activity-time">{{ activity.time }}</div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Top Talents & Hot Jobs -->
    <el-row :gutter="20" class="data-row">
      <el-col :xs="24" :lg="12">
        <div class="card data-card">
          <div class="card-header">
            <h3>热门人才</h3>
            <el-button type="primary" size="small" @click="$router.push('/talents')">查看更多</el-button>
          </div>
          <div class="talent-list">
            <div class="talent-item" v-for="(talent, index) in topTalents" :key="index"
                 @click="goToTalent(talent.id)" style="cursor: pointer;">
              <div class="talent-rank" :class="{ top: index < 3 }">{{ index + 1 }}</div>
              <el-avatar :size="40" :style="{ background: getAvatarColor(index) }">
                {{ talent.name.charAt(0) }}
              </el-avatar>
              <div class="talent-info">
                <div class="talent-name">{{ talent.name }}</div>
                <div class="talent-skills">
                  <el-tag v-for="skill in talent.skills.slice(0, 2)" :key="skill" size="small" type="info">
                    {{ skill }}
                  </el-tag>
                </div>
              </div>
              <div class="talent-score">
                <span class="score">{{ talent.score }}</span>
                <span class="label">匹配度</span>
              </div>
            </div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :lg="12">
        <div class="card data-card">
          <div class="card-header">
            <h3>热门职位</h3>
            <el-button type="primary" size="small" @click="$router.push('/jobs')">查看更多</el-button>
          </div>
          <div class="job-list">
            <div class="job-item" v-for="(job, index) in hotJobs" :key="index"
                 @click="goToJob(job.id)" style="cursor: pointer;">
              <div class="job-icon" :style="{ background: job.color }">
                <el-icon :size="20"><Suitcase /></el-icon>
              </div>
              <div class="job-info">
                <div class="job-title">{{ job.title }}</div>
                <div class="job-meta">
                  <span><el-icon><Location /></el-icon>{{ job.location }}</span>
                  <span><el-icon><Money /></el-icon>{{ job.salary }}</span>
                </div>
              </div>
              <div class="job-applicants">
                <span class="count">{{ job.applicants }}</span>
                <span class="label">申请</span>
              </div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, markRaw } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { statsApi } from '@/api/stats'
import * as echarts from 'echarts'
import {
  User, Suitcase, Document, MagicStick, ArrowUp, ArrowDown,
  Refresh, MoreFilled, ArrowRight, Plus, Search, Bell, Location, Money,
  EditPen, Upload, ChatDotRound, UserFilled, Briefcase, TrendCharts, DataAnalysis
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const jobChartRef = ref<HTMLElement>()
const trendChartRef = ref<HTMLElement>()
let jobChart: echarts.ECharts | null = null
let trendChart: echarts.ECharts | null = null

const dateRange = ref<[Date, Date] | null>(null)
const trendType = ref('week')

// 统计卡片数据
const statsCards = ref([
  {
    label: '总人才数',
    value: 0,
    displayValue: 0,
    suffix: '',
    icon: markRaw(User),
    colorClass: 'sky',
    trend: 0
  },
  {
    label: '在招职位',
    value: 0,
    displayValue: 0,
    suffix: '个',
    icon: markRaw(Suitcase),
    colorClass: 'cyan',
    trend: 0
  },
  {
    label: '本月申请',
    value: 0,
    displayValue: 0,
    suffix: '',
    icon: markRaw(Document),
    colorClass: 'blue',
    trend: 0
  },
  {
    label: '成功匹配',
    value: 0,
    displayValue: 0,
    suffix: '%',
    icon: markRaw(MagicStick),
    colorClass: 'teal',
    trend: 0
  }
])

// 从后端获取统计数据
const fetchDashboardStats = async () => {
  try {
    const res = await statsApi.getDashboardStats()
    if (res.data.code === 0) {
      const data = res.data.data
      statsCards.value[0].value = data.total_talents
      statsCards.value[0].trend = data.talent_trend
      statsCards.value[1].value = data.total_jobs
      statsCards.value[1].trend = data.job_trend
      statsCards.value[2].value = data.total_applications
      statsCards.value[2].trend = data.application_trend
      statsCards.value[3].value = data.match_rate
      animateNumbers()
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 快捷操作 - 天蓝色系
const quickActions = [
  {
    title: '发布新职位',
    desc: '创建并发布招聘职位',
    icon: markRaw(Plus),
    gradient: 'linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%)',
    route: '/jobs'
  },
  {
    title: '搜索人才',
    desc: '在人才库中搜索匹配候选人',
    icon: markRaw(Search),
    gradient: 'linear-gradient(135deg, #06b6d4 0%, #14b8a6 100%)',
    route: '/talents'
  },
  {
    title: '上传简历',
    desc: '批量导入候选人简历',
    icon: markRaw(Upload),
    gradient: 'linear-gradient(135deg, #38bdf8 0%, #22d3ee 100%)',
    route: '/resumes'
  },
  {
    title: '智能推荐',
    desc: '查看AI推荐的匹配结果',
    icon: markRaw(MagicStick),
    gradient: 'linear-gradient(135deg, #22c55e 0%, #10b981 100%)',
    route: '/recommend'
  }
]

// 最近活动 - 从后端获取
const recentActivities = ref<any[]>([])

// 从后端获取最近活动（消息/操作日志）
const fetchRecentActivities = async () => {
  const colors = ['#0ea5e9', '#06b6d4', '#38bdf8', '#22c55e', '#14b8a6']
  const icons = [markRaw(UserFilled), markRaw(Briefcase), markRaw(Document), markRaw(TrendCharts), markRaw(ChatDotRound)]
  try {
    const res = await fetch('/api/v1/messages?page=1&page_size=5')
    const data = await res.json()
    if (data.code === 0 && data.data?.list) {
      recentActivities.value = data.data.list.map((m: any, i: number) => ({
        title: m.title || '系统消息',
        description: m.content?.slice(0, 30) || '',
        time: formatTime(m.created_at),
        icon: icons[i % icons.length],
        color: colors[i % colors.length]
      }))
    }
  } catch (error) {
    console.error('获取最近活动失败:', error)
  }
}

// 格式化时间
const formatTime = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  return `${days}天前`
}

// 热门人才
const topTalents = ref<any[]>([])

// 热门职位
const hotJobs = ref<any[]>([])

// 从后端获取热门人才
const fetchTopTalents = async () => {
  try {
    const res = await fetch('/api/v1/talents?page=1&page_size=5')
    const data = await res.json()
    if (data.code === 0 && data.data?.list) {
      topTalents.value = data.data.list.map((t: any) => ({
        id: t.id,
        name: t.name,
        skills: t.skills || [],
        score: t.experience ? t.experience * 10 : 80
      }))
    }
  } catch (error) {
    console.error('获取热门人才失败:', error)
  }
}

// 从后端获取热门职位
const fetchHotJobs = async () => {
  const colors = ['#0ea5e9', '#06b6d4', '#38bdf8', '#22c55e', '#14b8a6']
  try {
    const res = await fetch('/api/v1/jobs?page=1&page_size=5')
    const data = await res.json()
    if (data.code === 0 && data.data?.jobs) {
      hotJobs.value = data.data.jobs.map((j: any, i: number) => ({
        id: j.id,
        title: j.title,
        location: j.location,
        salary: j.salary,
        applicants: j.applicants || 0,
        color: colors[i % colors.length]
      }))
    }
  } catch (error) {
    console.error('获取热门职位失败:', error)
  }
}

// 跳转到人才详情
const goToTalent = (id: number) => {
  router.push(`/talents/${id}`)
}

// 跳转到职位详情
const goToJob = (id: number) => {
  router.push(`/jobs/${id}`)
}

// 获取头像颜色 - 天蓝色系
const getAvatarColor = (index: number) => {
  const colors = ['#0ea5e9', '#06b6d4', '#38bdf8', '#22c55e', '#14b8a6']
  return colors[index % colors.length]
}

// 数字动画
const animateNumbers = () => {
  statsCards.value.forEach((stat) => {
    const duration = 1500
    const startTime = Date.now()
    const startValue = 0
    const endValue = stat.value

    const animate = () => {
      const elapsed = Date.now() - startTime
      const progress = Math.min(elapsed / duration, 1)
      // easeOutQuart
      const easeProgress = 1 - Math.pow(1 - progress, 4)
      stat.displayValue = Math.round(startValue + (endValue - startValue) * easeProgress)

      if (progress < 1) {
        requestAnimationFrame(animate)
      }
    }
    requestAnimationFrame(animate)
  })
}

// 初始化趋势图表 - 天蓝色系
const initTrendChart = () => {
  if (!trendChartRef.value) return

  trendChart = echarts.init(trendChartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: { backgroundColor: '#0ea5e9' }
      }
    },
    legend: {
      data: ['新增人才', '职位申请', '成功匹配'],
      bottom: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '12%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      axisLine: { lineStyle: { color: '#e2e8f0' } },
      axisLabel: { color: '#64748b' }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: '#f1f5f9', type: 'dashed' } },
      axisLabel: { color: '#64748b' }
    },
    series: [
      {
        name: '新增人才',
        type: 'line',
        smooth: true,
        lineStyle: { width: 3, color: '#0ea5e9' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(14, 165, 233, 0.3)' },
            { offset: 1, color: 'rgba(14, 165, 233, 0.05)' }
          ])
        },
        itemStyle: { color: '#0ea5e9' },
        label: {
          show: true,
          position: 'top',
          color: '#0ea5e9',
          fontSize: 11,
          fontWeight: 'bold'
        },
        data: [12, 19, 15, 25, 22, 18, 28]
      },
      {
        name: '职位申请',
        type: 'line',
        smooth: true,
        lineStyle: { width: 3, color: '#06b6d4' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(6, 182, 212, 0.3)' },
            { offset: 1, color: 'rgba(6, 182, 212, 0.05)' }
          ])
        },
        itemStyle: { color: '#06b6d4' },
        label: {
          show: true,
          position: 'top',
          color: '#06b6d4',
          fontSize: 11,
          fontWeight: 'bold'
        },
        data: [20, 32, 28, 45, 38, 30, 42]
      },
      {
        name: '成功匹配',
        type: 'line',
        smooth: true,
        lineStyle: { width: 3, color: '#22c55e' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(34, 197, 94, 0.3)' },
            { offset: 1, color: 'rgba(34, 197, 94, 0.05)' }
          ])
        },
        itemStyle: { color: '#22c55e' },
        label: {
          show: true,
          position: 'top',
          color: '#22c55e',
          fontSize: 11,
          fontWeight: 'bold'
        },
        data: [5, 8, 6, 12, 10, 8, 15]
      }
    ]
  }
  trendChart.setOption(option)
}

// 初始化职位图表 - 天蓝色系
const initJobChart = () => {
  if (!jobChartRef.value) return

  jobChart = echarts.init(jobChartRef.value)
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c}个 ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: '2%',
      top: 'center',
      itemGap: 16,
      textStyle: { color: '#64748b' }
    },
    series: [
      {
        name: '职位状态',
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['35%', '50%'],
        avoidLabelOverlap: true,
        itemStyle: {
          borderRadius: 8,
          borderColor: '#fff',
          borderWidth: 3
        },
        label: {
          show: true,
          position: 'inside',
          formatter: '{c}个\n{d}%',
          fontSize: 11,
          lineHeight: 14,
          color: '#fff',
          fontWeight: 'bold'
        },
        labelLine: {
          show: false
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 12,
            fontWeight: 'bold'
          },
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.2)'
          }
        },
        data: [
          { value: 48, name: '招聘中', itemStyle: { color: '#0ea5e9' } },
          { value: 23, name: '已暂停', itemStyle: { color: '#06b6d4' } },
          { value: 15, name: '已关闭', itemStyle: { color: '#94a3b8' } },
          { value: 32, name: '已完成', itemStyle: { color: '#22c55e' } }
        ]
      }
    ]
  }
  jobChart.setOption(option)
}

// 处理操作
const handleAction = (route: string) => {
  router.push(route)
}

// 刷新数据
const refreshData = () => {
  fetchDashboardStats()
}

// 监听图表类型变化
watch(trendType, () => {
  // 根据不同时间范围更新图表数据
  if (trendChart) {
    const dataMap: Record<string, string[]> = {
      week: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      month: ['第1周', '第2周', '第3周', '第4周'],
      year: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
    }
    trendChart.setOption({
      xAxis: { data: dataMap[trendType.value] }
    })
  }
})

// 窗口大小变化时重绘图表
const handleResize = () => {
  jobChart?.resize()
  trendChart?.resize()
}

onMounted(() => {
  fetchDashboardStats()
  fetchTopTalents()
  fetchHotJobs()
  fetchRecentActivities()
  initTrendChart()
  initJobChart()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  jobChart?.dispose()
  trendChart?.dispose()
})
</script>

<style scoped lang="scss">
.dashboard {
  animation: fadeIn 0.5s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;

  .header-left {
    .page-title {
      font-size: 28px;
      font-weight: 700;
      color: var(--text-primary);
      margin: 0 0 4px 0;
    }

    .page-subtitle {
      color: var(--text-secondary);
      font-size: 14px;
      margin: 0;
    }
  }

  .header-right {
    display: flex;
    gap: 12px;
    align-items: center;
  }
}

.stats-row {
  margin-bottom: 20px;

  .stat-card {
    margin-bottom: 20px;
    animation: slideUp 0.5s ease forwards;
    opacity: 0;

    @keyframes slideUp {
      from { opacity: 0; transform: translateY(20px); }
      to { opacity: 1; transform: translateY(0); }
    }

    .stat-card-inner {
      position: relative;
      padding: 24px;
      border-radius: 16px;
      background: var(--bg-primary);
      box-shadow: var(--shadow-card);
      border: 1px solid var(--border-light);
      overflow: hidden;
      transition: box-shadow 0.3s ease;

      &:hover {
        box-shadow: var(--shadow-md);
      }

      /* 天蓝色系 */
      &.sky { border-left: 4px solid #0ea5e9; }
      &.cyan { border-left: 4px solid #06b6d4; }
      &.blue { border-left: 4px solid #38bdf8; }
      &.teal { border-left: 4px solid #14b8a6; }

      .stat-icon-wrapper {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 16px;

        .stat-icon {
          width: 56px;
          height: 56px;
          border-radius: 14px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
        }

        .stat-trend {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 13px;
          font-weight: 600;
          padding: 4px 8px;
          border-radius: 20px;

          &.up {
            color: #10b981;
            background: rgba(16, 185, 129, 0.1);
          }

          &.down {
            color: #ef4444;
            background: rgba(239, 68, 68, 0.1);
          }
        }
      }

      /* 天蓝色系图标背景 */
      &.sky .stat-icon { background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%); }
      &.cyan .stat-icon { background: linear-gradient(135deg, #06b6d4 0%, #14b8a6 100%); }
      &.blue .stat-icon { background: linear-gradient(135deg, #38bdf8 0%, #22d3ee 100%); }
      &.teal .stat-icon { background: linear-gradient(135deg, #14b8a6 0%, #22c55e 100%); }

      .stat-content {
        position: relative;
        z-index: 1;

        .stat-value {
          display: flex;
          align-items: baseline;
          gap: 4px;

          .number {
            font-size: 32px;
            font-weight: 700;
            color: var(--text-primary);
            line-height: 1;
          }

          .suffix {
            font-size: 16px;
            color: var(--text-secondary);
          }
        }

        .stat-label {
          font-size: 14px;
          color: var(--text-secondary);
          margin-top: 8px;
        }
      }

      .stat-bg-icon {
        position: absolute;
        right: -20px;
        bottom: -20px;
        opacity: 0.05;
        color: #1a1a2e;
      }
    }
  }
}

.card {
  background: var(--bg-primary);
  border-radius: 16px;
  padding: 24px;
  box-shadow: var(--shadow-card);
  border: 1px solid var(--border-light);

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      font-size: 18px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0;
    }
  }
}

.charts-row {
  margin-bottom: 20px;

  .chart-card {
    height: 100%;
    min-height: 430px;
  }
}

.bottom-row {
  margin-bottom: 20px;

  .quick-actions-card {
    height: 100%;

    .quick-actions {
      display: flex;
      flex-direction: column;
      gap: 12px;

      .action-item {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 16px;
        border-radius: 12px;
        background: var(--bg-secondary);
        cursor: pointer;
        transition: background-color 0.3s ease;

        &:hover {
          background: var(--bg-tertiary);

          .action-arrow {
            opacity: 1;
          }
        }

        .action-icon {
          width: 48px;
          height: 48px;
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          flex-shrink: 0;
        }

        .action-info {
          flex: 1;

          .action-title {
            display: block;
            font-size: 15px;
            font-weight: 600;
            color: var(--text-primary);
            margin-bottom: 2px;
          }

          .action-desc {
            display: block;
            font-size: 13px;
            color: var(--text-secondary);
          }
        }

        .action-arrow {
          color: var(--text-muted);
          opacity: 0;
          transition: opacity 0.3s ease;
        }
      }
    }
  }

  .recent-activities-card {
    height: 100%;

    .activities-list {
      display: flex;
      flex-direction: column;
      gap: 16px;

      .activity-item {
        display: flex;
        align-items: center;
        gap: 16px;
        padding-bottom: 16px;
        border-bottom: 1px solid var(--border-light);

        &:last-child {
          padding-bottom: 0;
          border-bottom: none;
        }

        .activity-avatar {
          width: 44px;
          height: 44px;
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          flex-shrink: 0;
        }

        .activity-content {
          flex: 1;

          .activity-title {
            font-size: 14px;
            font-weight: 600;
            color: var(--text-primary);
            margin-bottom: 2px;
          }

          .activity-desc {
            font-size: 13px;
            color: var(--text-secondary);
          }
        }

        .activity-time {
          font-size: 12px;
          color: var(--text-muted);
          white-space: nowrap;
        }
      }
    }
  }
}

.data-row {
  .data-card {
    height: 100%;
    margin-bottom: 20px;

    .talent-list, .job-list {
      display: flex;
      flex-direction: column;
      gap: 12px;
    }

    .talent-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px;
      border-radius: 12px;
      background: var(--bg-secondary);
      transition: all 0.3s ease;

      &:hover {
        background: var(--bg-tertiary);
      }

      .talent-rank {
        width: 28px;
        height: 28px;
        border-radius: 8px;
        background: var(--bg-tertiary);
        color: var(--text-secondary);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        font-weight: 600;

        &.top {
          background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%);
          color: white;
        }
      }

      .talent-info {
        flex: 1;

        .talent-name {
          font-size: 14px;
          font-weight: 600;
          color: var(--text-primary);
          margin-bottom: 4px;
        }

        .talent-skills {
          display: flex;
          gap: 4px;
        }
      }

      .talent-score {
        text-align: right;

        .score {
          display: block;
          font-size: 18px;
          font-weight: 700;
          color: #0ea5e9;
        }

        .label {
          font-size: 12px;
          color: var(--text-muted);
        }
      }
    }

    .job-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px;
      border-radius: 12px;
      background: var(--bg-secondary);
      transition: all 0.3s ease;

      &:hover {
        background: var(--bg-tertiary);
      }

      .job-icon {
        width: 44px;
        height: 44px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        flex-shrink: 0;
      }

      .job-info {
        flex: 1;

        .job-title {
          font-size: 14px;
          font-weight: 600;
          color: var(--text-primary);
          margin-bottom: 4px;
        }

        .job-meta {
          display: flex;
          gap: 16px;
          font-size: 12px;
          color: var(--text-secondary);

          span {
            display: flex;
            align-items: center;
            gap: 4px;
          }
        }
      }

      .job-applicants {
        text-align: right;

        .count {
          display: block;
          font-size: 18px;
          font-weight: 700;
          color: #06b6d4;
        }

        .label {
          font-size: 12px;
          color: var(--text-muted);
        }
      }
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .page-header {
    .header-right {
      width: 100%;
      flex-direction: column;

      .el-date-picker {
        width: 100% !important;
      }

      .el-button {
        width: 100%;
      }
    }
  }
}
</style>
