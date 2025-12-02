<template>
  <div class="dashboard">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">仪表板</h1>
        <p class="page-subtitle">欢迎回来，{{ userStore.user?.username || '用户' }}！这是您的数据概览</p>
      </div>
      <div class="header-right">
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
            <el-button text type="primary" @click="$router.push('/messages')">查看全部</el-button>
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
            <el-button text type="primary" @click="$router.push('/talents')">查看更多</el-button>
          </div>
          <div class="talent-list">
            <div class="talent-item" v-for="(talent, index) in topTalents" :key="index">
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
            <el-button text type="primary" @click="$router.push('/jobs')">查看更多</el-button>
          </div>
          <div class="job-list">
            <div class="job-item" v-for="(job, index) in hotJobs" :key="index">
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
import * as echarts from 'echarts'
import {
  User, Suitcase, Document, MagicStick, ArrowUp, ArrowDown,
  Refresh, MoreFilled, ArrowRight, Plus, Search, Bell, Location, Money,
  EditPen, Upload, ChatDotRound, UserFilled, Briefcase, TrendCharts
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
    value: 1560,
    displayValue: 0,
    suffix: '',
    icon: markRaw(User),
    colorClass: 'purple',
    trend: 12.5
  },
  {
    label: '在招职位',
    value: 48,
    displayValue: 0,
    suffix: '个',
    icon: markRaw(Suitcase),
    colorClass: 'pink',
    trend: 8.2
  },
  {
    label: '本月申请',
    value: 326,
    displayValue: 0,
    suffix: '',
    icon: markRaw(Document),
    colorClass: 'blue',
    trend: -3.1
  },
  {
    label: '成功匹配',
    value: 89,
    displayValue: 0,
    suffix: '%',
    icon: markRaw(MagicStick),
    colorClass: 'green',
    trend: 5.8
  }
])

// 快捷操作
const quickActions = [
  {
    title: '发布新职位',
    desc: '创建并发布招聘职位',
    icon: markRaw(Plus),
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    route: '/jobs'
  },
  {
    title: '搜索人才',
    desc: '在人才库中搜索匹配候选人',
    icon: markRaw(Search),
    gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    route: '/talents'
  },
  {
    title: '上传简历',
    desc: '批量导入候选人简历',
    icon: markRaw(Upload),
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    route: '/resumes'
  },
  {
    title: '智能推荐',
    desc: '查看AI推荐的匹配结果',
    icon: markRaw(MagicStick),
    gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    route: '/recommend'
  }
]

// 最近活动
const recentActivities = ref([
  {
    title: '新增候选人',
    description: '张三已加入人才库',
    time: '5分钟前',
    icon: markRaw(UserFilled),
    color: '#667eea'
  },
  {
    title: '职位发布',
    description: '高级Go开发工程师职位已上线',
    time: '15分钟前',
    icon: markRaw(Briefcase),
    color: '#f093fb'
  },
  {
    title: '简历审核',
    description: '李四的简历已通过初审',
    time: '1小时前',
    icon: markRaw(Document),
    color: '#4facfe'
  },
  {
    title: '智能匹配',
    description: '为前端开发岗位推荐了5位候选人',
    time: '2小时前',
    icon: markRaw(TrendCharts),
    color: '#43e97b'
  },
  {
    title: '面试安排',
    description: '王五的面试已安排在明天下午3点',
    time: '3小时前',
    icon: markRaw(ChatDotRound),
    color: '#f5576c'
  }
])

// 热门人才
const topTalents = ref([
  { name: '张三', skills: ['Go', 'Python', 'Kubernetes'], score: 95 },
  { name: '李四', skills: ['React', 'TypeScript', 'Node.js'], score: 92 },
  { name: '王五', skills: ['Java', 'Spring', 'MySQL'], score: 88 },
  { name: '赵六', skills: ['Python', 'TensorFlow', 'PyTorch'], score: 85 },
  { name: '钱七', skills: ['Vue', 'Element Plus', 'Vite'], score: 82 }
])

// 热门职位
const hotJobs = ref([
  { title: '高级Go开发工程师', location: '北京', salary: '30-50K', applicants: 128, color: '#667eea' },
  { title: '前端架构师', location: '上海', salary: '40-60K', applicants: 96, color: '#f093fb' },
  { title: 'AI算法工程师', location: '深圳', salary: '50-80K', applicants: 87, color: '#4facfe' },
  { title: '产品经理', location: '杭州', salary: '25-40K', applicants: 156, color: '#43e97b' },
  { title: 'DevOps工程师', location: '广州', salary: '30-45K', applicants: 64, color: '#f5576c' }
])

// 获取头像颜色
const getAvatarColor = (index: number) => {
  const colors = ['#667eea', '#f093fb', '#4facfe', '#43e97b', '#f5576c']
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

// 初始化趋势图表
const initTrendChart = () => {
  if (!trendChartRef.value) return

  trendChart = echarts.init(trendChartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: { backgroundColor: '#6a7985' }
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
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      axisLine: { lineStyle: { color: '#e5e7eb' } },
      axisLabel: { color: '#6b7280' }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: '#f3f4f6', type: 'dashed' } },
      axisLabel: { color: '#6b7280' }
    },
    series: [
      {
        name: '新增人才',
        type: 'line',
        smooth: true,
        lineStyle: { width: 3, color: '#667eea' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(102, 126, 234, 0.3)' },
            { offset: 1, color: 'rgba(102, 126, 234, 0.05)' }
          ])
        },
        itemStyle: { color: '#667eea' },
        data: [12, 19, 15, 25, 22, 18, 28]
      },
      {
        name: '职位申请',
        type: 'line',
        smooth: true,
        lineStyle: { width: 3, color: '#f093fb' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(240, 147, 251, 0.3)' },
            { offset: 1, color: 'rgba(240, 147, 251, 0.05)' }
          ])
        },
        itemStyle: { color: '#f093fb' },
        data: [20, 32, 28, 45, 38, 30, 42]
      },
      {
        name: '成功匹配',
        type: 'line',
        smooth: true,
        lineStyle: { width: 3, color: '#43e97b' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(67, 233, 123, 0.3)' },
            { offset: 1, color: 'rgba(67, 233, 123, 0.05)' }
          ])
        },
        itemStyle: { color: '#43e97b' },
        data: [5, 8, 6, 12, 10, 8, 15]
      }
    ]
  }
  trendChart.setOption(option)
}

// 初始化职位图表
const initJobChart = () => {
  if (!jobChartRef.value) return

  jobChart = echarts.init(jobChartRef.value)
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: '5%',
      top: 'center',
      itemGap: 16,
      textStyle: { color: '#6b7280' }
    },
    series: [
      {
        name: '职位状态',
        type: 'pie',
        radius: ['50%', '75%'],
        center: ['35%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 8,
          borderColor: '#fff',
          borderWidth: 3
        },
        label: {
          show: false
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          },
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.2)'
          }
        },
        labelLine: { show: false },
        data: [
          { value: 48, name: '招聘中', itemStyle: { color: '#667eea' } },
          { value: 23, name: '已暂停', itemStyle: { color: '#f093fb' } },
          { value: 15, name: '已关闭', itemStyle: { color: '#9ca3af' } },
          { value: 32, name: '已完成', itemStyle: { color: '#43e97b' } }
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
  animateNumbers()
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
  animateNumbers()
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
      color: #1a1a2e;
      margin: 0 0 4px 0;
    }

    .page-subtitle {
      color: #6b7280;
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
      background: white;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
      overflow: hidden;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
      }

      &.purple { border-left: 4px solid #667eea; }
      &.pink { border-left: 4px solid #f093fb; }
      &.blue { border-left: 4px solid #4facfe; }
      &.green { border-left: 4px solid #43e97b; }

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

      &.purple .stat-icon { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
      &.pink .stat-icon { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
      &.blue .stat-icon { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
      &.green .stat-icon { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }

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
            color: #1a1a2e;
            line-height: 1;
          }

          .suffix {
            font-size: 16px;
            color: #6b7280;
          }
        }

        .stat-label {
          font-size: 14px;
          color: #6b7280;
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
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      font-size: 18px;
      font-weight: 600;
      color: #1a1a2e;
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
        background: #f9fafb;
        cursor: pointer;
        transition: all 0.3s ease;

        &:hover {
          background: #f3f4f6;
          transform: translateX(4px);

          .action-arrow {
            opacity: 1;
            transform: translateX(0);
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
            color: #1a1a2e;
            margin-bottom: 2px;
          }

          .action-desc {
            display: block;
            font-size: 13px;
            color: #6b7280;
          }
        }

        .action-arrow {
          color: #9ca3af;
          opacity: 0;
          transform: translateX(-10px);
          transition: all 0.3s ease;
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
        border-bottom: 1px solid #f3f4f6;

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
            color: #1a1a2e;
            margin-bottom: 2px;
          }

          .activity-desc {
            font-size: 13px;
            color: #6b7280;
          }
        }

        .activity-time {
          font-size: 12px;
          color: #9ca3af;
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
      background: #f9fafb;
      transition: all 0.3s ease;

      &:hover {
        background: #f3f4f6;
      }

      .talent-rank {
        width: 28px;
        height: 28px;
        border-radius: 8px;
        background: #e5e7eb;
        color: #6b7280;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        font-weight: 600;

        &.top {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: white;
        }
      }

      .talent-info {
        flex: 1;

        .talent-name {
          font-size: 14px;
          font-weight: 600;
          color: #1a1a2e;
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
          color: #667eea;
        }

        .label {
          font-size: 12px;
          color: #9ca3af;
        }
      }
    }

    .job-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px;
      border-radius: 12px;
      background: #f9fafb;
      transition: all 0.3s ease;

      &:hover {
        background: #f3f4f6;
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
          color: #1a1a2e;
          margin-bottom: 4px;
        }

        .job-meta {
          display: flex;
          gap: 16px;
          font-size: 12px;
          color: #6b7280;

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
          color: #f093fb;
        }

        .label {
          font-size: 12px;
          color: #9ca3af;
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
