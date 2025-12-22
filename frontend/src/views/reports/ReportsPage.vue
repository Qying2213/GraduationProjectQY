<template>
  <div class="reports-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>数据报表</h1>
        <p class="subtitle">全面的招聘数据分析与可视化</p>
      </div>
      <div class="header-right">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          :shortcuts="dateShortcuts"
          style="width: 280px"
          @change="fetchReportData"
        />
        <el-dropdown @command="handleExport" trigger="click">
          <el-button type="primary">
            <el-icon><Download /></el-icon>
            导出报表
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="pdf">导出 PDF</el-dropdown-item>
              <el-dropdown-item command="excel">导出 Excel</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 核心指标卡片 -->
    <div class="metrics-grid">
      <div class="metric-card" v-for="(metric, index) in coreMetrics" :key="index">
        <div class="metric-header">
          <div class="metric-icon" :style="{ background: metric.gradient }">
            <el-icon :size="24"><component :is="metric.icon" /></el-icon>
          </div>
          <div class="metric-trend" :class="metric.trend > 0 ? 'up' : 'down'">
            <el-icon><component :is="metric.trend > 0 ? ArrowUp : ArrowDown" /></el-icon>
            {{ Math.abs(metric.trend) }}%
          </div>
        </div>
        <div class="metric-value">{{ metric.value }}</div>
        <div class="metric-label">{{ metric.label }}</div>
        <div class="metric-compare">较上期 {{ metric.trend > 0 ? '增长' : '下降' }}</div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-section">
      <el-row :gutter="20">
        <!-- 招聘漏斗 -->
        <el-col :xs="24" :lg="12">
          <div class="chart-card">
            <div class="chart-header">
              <h3>招聘漏斗分析</h3>
              <el-radio-group v-model="funnelPeriod" size="small">
                <el-radio-button label="week">本周</el-radio-button>
                <el-radio-button label="month">本月</el-radio-button>
                <el-radio-button label="quarter">本季</el-radio-button>
              </el-radio-group>
            </div>
            <div ref="funnelChartRef" class="chart-container"></div>
            <div class="funnel-stats">
              <div class="funnel-stat" v-for="(stat, index) in funnelStats" :key="index">
                <span class="stat-label">{{ stat.label }}</span>
                <span class="stat-value">{{ stat.value }}</span>
                <span class="stat-rate">{{ stat.rate }}</span>
              </div>
            </div>
          </div>
        </el-col>

        <!-- 职位热度排行 -->
        <el-col :xs="24" :lg="12">
          <div class="chart-card">
            <div class="chart-header">
              <h3>职位热度排行</h3>
              <el-button type="primary" size="small" @click="$router.push('/jobs')">查看全部</el-button>
            </div>
            <div ref="jobRankChartRef" class="chart-container"></div>
          </div>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px">
        <!-- 招聘趋势 -->
        <el-col :xs="24" :lg="16">
          <div class="chart-card">
            <div class="chart-header">
              <h3>招聘趋势分析</h3>
              <div class="chart-legend">
                <span class="legend-item"><i class="dot" style="background: #667eea"></i>简历投递</span>
                <span class="legend-item"><i class="dot" style="background: #43e97b"></i>面试安排</span>
                <span class="legend-item"><i class="dot" style="background: #f093fb"></i>录用人数</span>
              </div>
            </div>
            <div ref="trendChartRef" class="chart-container large"></div>
          </div>
        </el-col>

        <!-- 渠道效果分析 -->
        <el-col :xs="24" :lg="8">
          <div class="chart-card">
            <div class="chart-header">
              <h3>渠道效果分析</h3>
            </div>
            <div ref="channelChartRef" class="chart-container"></div>
            <div class="channel-list">
              <div class="channel-item" v-for="(channel, index) in channelData" :key="index">
                <div class="channel-info">
                  <span class="channel-dot" :style="{ background: channel.color }"></span>
                  <span class="channel-name">{{ channel.name }}</span>
                </div>
                <div class="channel-stats">
                  <span class="channel-count">{{ channel.count }}人</span>
                  <span class="channel-rate">{{ channel.rate }}%</span>
                </div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 详细数据表格 -->
    <div class="data-section">
      <el-row :gutter="20">
        <!-- 部门招聘进度 -->
        <el-col :xs="24" :lg="12">
          <div class="data-card">
            <div class="card-header">
              <h3>部门招聘进度</h3>
            </div>
            <el-table :data="departmentProgress" stripe>
              <el-table-column prop="department" label="部门" width="120" />
              <el-table-column prop="target" label="目标" width="80" />
              <el-table-column prop="hired" label="已录用" width="80" />
              <el-table-column label="完成率">
                <template #default="{ row }">
                  <div class="progress-cell">
                    <el-progress
                      :percentage="row.progress"
                      :color="getProgressColor(row.progress)"
                      :stroke-width="8"
                    />
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-col>

        <!-- 面试官效率排行 -->
        <el-col :xs="24" :lg="12">
          <div class="data-card">
            <div class="card-header">
              <h3>面试官效率排行</h3>
            </div>
            <div class="interviewer-list">
              <div class="interviewer-item" v-for="(item, index) in interviewerRank" :key="index">
                <div class="rank" :class="{ top: index < 3 }">{{ index + 1 }}</div>
                <el-avatar :size="40" :style="{ background: getAvatarColor(index) }">
                  {{ item.name.charAt(0) }}
                </el-avatar>
                <div class="interviewer-info">
                  <span class="name">{{ item.name }}</span>
                  <span class="dept">{{ item.department }}</span>
                </div>
                <div class="interviewer-stats">
                  <div class="stat">
                    <span class="value">{{ item.interviews }}</span>
                    <span class="label">面试数</span>
                  </div>
                  <div class="stat">
                    <span class="value">{{ item.passRate }}%</span>
                    <span class="label">通过率</span>
                  </div>
                  <div class="stat">
                    <span class="value">{{ item.avgScore }}</span>
                    <span class="label">平均评分</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, markRaw } from 'vue'
import * as echarts from 'echarts'
import {
  Download, ArrowDown, ArrowUp,
  User, Suitcase, Document, TrendCharts
} from '@element-plus/icons-vue'

const dateRange = ref<[Date, Date] | null>(null)
const funnelPeriod = ref('month')

const funnelChartRef = ref<HTMLElement>()
const jobRankChartRef = ref<HTMLElement>()
const trendChartRef = ref<HTMLElement>()
const channelChartRef = ref<HTMLElement>()

let funnelChart: echarts.ECharts | null = null
let jobRankChart: echarts.ECharts | null = null
let trendChart: echarts.ECharts | null = null
let channelChart: echarts.ECharts | null = null

const dateShortcuts = [
  { text: '最近一周', value: () => {
    const end = new Date()
    const start = new Date()
    start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
    return [start, end]
  }},
  { text: '最近一月', value: () => {
    const end = new Date()
    const start = new Date()
    start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
    return [start, end]
  }},
  { text: '最近三月', value: () => {
    const end = new Date()
    const start = new Date()
    start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
    return [start, end]
  }}
]

// 核心指标
const coreMetrics = ref([
  {
    label: '简历投递量',
    value: '1,256',
    trend: 12.5,
    icon: markRaw(Document),
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  },
  {
    label: '面试安排数',
    value: '328',
    trend: 8.3,
    icon: markRaw(User),
    gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)'
  },
  {
    label: '录用人数',
    value: '45',
    trend: -3.2,
    icon: markRaw(Suitcase),
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)'
  },
  {
    label: '平均招聘周期',
    value: '18天',
    trend: -15.6,
    icon: markRaw(TrendCharts),
    gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)'
  }
])

// 漏斗统计
const funnelStats = ref([
  { label: '简历筛选率', value: '45%', rate: '↑ 5%' },
  { label: '面试通过率', value: '62%', rate: '↑ 3%' },
  { label: '录用转化率', value: '28%', rate: '↓ 2%' }
])

// 渠道数据
const channelData = ref([
  { name: '官网投递', count: 456, rate: 36, color: '#667eea' },
  { name: '猎聘网', count: 312, rate: 25, color: '#f093fb' },
  { name: 'BOSS直聘', count: 234, rate: 19, color: '#4facfe' },
  { name: '内部推荐', count: 156, rate: 12, color: '#43e97b' },
  { name: '其他渠道', count: 98, rate: 8, color: '#f5576c' }
])

// 部门招聘进度
const departmentProgress = ref([
  { department: '技术部', target: 20, hired: 15, progress: 75 },
  { department: '产品部', target: 8, hired: 6, progress: 75 },
  { department: '设计部', target: 5, hired: 5, progress: 100 },
  { department: '市场部', target: 10, hired: 4, progress: 40 },
  { department: '运营部', target: 6, hired: 3, progress: 50 }
])

// 面试官排行
const interviewerRank = ref([
  { name: '陈总监', department: '技术部', interviews: 45, passRate: 68, avgScore: 4.5 },
  { name: '刘经理', department: '产品部', interviews: 38, passRate: 72, avgScore: 4.3 },
  { name: '周主管', department: '技术部', interviews: 32, passRate: 65, avgScore: 4.2 },
  { name: '王总监', department: '设计部', interviews: 28, passRate: 78, avgScore: 4.6 },
  { name: 'HR小李', department: '人力资源', interviews: 56, passRate: 82, avgScore: 4.4 }
])

const getProgressColor = (progress: number) => {
  if (progress >= 80) return '#43e97b'
  if (progress >= 50) return '#667eea'
  return '#f5576c'
}

const getAvatarColor = (index: number) => {
  const colors = ['#667eea', '#f093fb', '#4facfe', '#43e97b', '#f5576c']
  return colors[index % colors.length]
}

const fetchReportData = () => {
  // 获取报表数据
}

const handleExport = (format: string) => {
  console.log('Export as', format)
}

// 初始化图表
const initCharts = () => {
  // 漏斗图
  if (funnelChartRef.value) {
    funnelChart = echarts.init(funnelChartRef.value)
    funnelChart.setOption({
      tooltip: { trigger: 'item', formatter: '{b}: {c}' },
      series: [{
        type: 'funnel',
        left: '10%',
        width: '80%',
        label: { position: 'inside', formatter: '{b}\n{c}' },
        itemStyle: { borderWidth: 0 },
        data: [
          { value: 1256, name: '简历投递', itemStyle: { color: '#667eea' } },
          { value: 565, name: '简历筛选', itemStyle: { color: '#818cf8' } },
          { value: 328, name: '面试安排', itemStyle: { color: '#f093fb' } },
          { value: 156, name: '面试通过', itemStyle: { color: '#4facfe' } },
          { value: 45, name: '录用入职', itemStyle: { color: '#43e97b' } }
        ]
      }]
    })
  }

  // 职位热度排行
  if (jobRankChartRef.value) {
    jobRankChart = echarts.init(jobRankChartRef.value)
    jobRankChart.setOption({
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      grid: { left: '3%', right: '10%', bottom: '3%', top: '3%', containLabel: true },
      xAxis: { type: 'value', show: false },
      yAxis: {
        type: 'category',
        data: ['数据分析师', 'UI设计师', '产品经理', '前端工程师', 'Go开发工程师'],
        axisLine: { show: false },
        axisTick: { show: false }
      },
      series: [{
        type: 'bar',
        data: [
          { value: 89, itemStyle: { color: '#667eea', borderRadius: [0, 4, 4, 0] } },
          { value: 112, itemStyle: { color: '#818cf8', borderRadius: [0, 4, 4, 0] } },
          { value: 134, itemStyle: { color: '#f093fb', borderRadius: [0, 4, 4, 0] } },
          { value: 156, itemStyle: { color: '#4facfe', borderRadius: [0, 4, 4, 0] } },
          { value: 198, itemStyle: { color: '#43e97b', borderRadius: [0, 4, 4, 0] } }
        ],
        label: { show: true, position: 'right', formatter: '{c}人' }
      }]
    })
  }

  // 趋势图
  if (trendChartRef.value) {
    trendChart = echarts.init(trendChartRef.value)
    trendChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', top: '10%', containLabel: true },
      xAxis: {
        type: 'category',
        data: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月'],
        axisLine: { lineStyle: { color: '#e5e7eb' } },
        axisLabel: { color: '#6b7280' }
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        splitLine: { lineStyle: { color: '#f3f4f6', type: 'dashed' } }
      },
      series: [
        {
          name: '简历投递',
          type: 'line',
          smooth: true,
          data: [120, 132, 101, 134, 90, 230, 210, 182, 191, 234, 290, 330],
          lineStyle: { color: '#667eea', width: 3 },
          itemStyle: { color: '#667eea' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(102, 126, 234, 0.3)' },
              { offset: 1, color: 'rgba(102, 126, 234, 0.05)' }
            ])
          }
        },
        {
          name: '面试安排',
          type: 'line',
          smooth: true,
          data: [45, 52, 38, 48, 35, 85, 78, 68, 72, 88, 108, 125],
          lineStyle: { color: '#43e97b', width: 3 },
          itemStyle: { color: '#43e97b' }
        },
        {
          name: '录用人数',
          type: 'line',
          smooth: true,
          data: [8, 12, 6, 10, 5, 18, 15, 12, 14, 16, 22, 28],
          lineStyle: { color: '#f093fb', width: 3 },
          itemStyle: { color: '#f093fb' }
        }
      ]
    })
  }

  // 渠道饼图
  if (channelChartRef.value) {
    channelChart = echarts.init(channelChartRef.value)
    channelChart.setOption({
      tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
      series: [{
        type: 'pie',
        radius: ['50%', '70%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: false,
        label: { show: false },
        data: channelData.value.map(c => ({
          value: c.count,
          name: c.name,
          itemStyle: { color: c.color }
        }))
      }]
    })
  }
}

const handleResize = () => {
  funnelChart?.resize()
  jobRankChart?.resize()
  trendChart?.resize()
  channelChart?.resize()
}

onMounted(() => {
  initCharts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  funnelChart?.dispose()
  jobRankChart?.dispose()
  trendChart?.dispose()
  channelChart?.dispose()
})
</script>

<style scoped lang="scss">
.reports-page {
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
  margin-bottom: 24px;

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

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;

  @media (max-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }

  .metric-card {
    background: var(--bg-primary);
    border-radius: 16px;
    padding: 24px;
    box-shadow: var(--shadow-card);

    .metric-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 16px;

      .metric-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
      }

      .metric-trend {
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

    .metric-value {
      font-size: 32px;
      font-weight: 700;
      color: var(--text-primary);
      margin-bottom: 4px;
    }

    .metric-label {
      font-size: 14px;
      color: var(--text-secondary);
      margin-bottom: 8px;
    }

    .metric-compare {
      font-size: 12px;
      color: var(--text-tertiary);
    }
  }
}

.chart-card, .data-card {
  background: var(--bg-primary);
  border-radius: 16px;
  padding: 24px;
  box-shadow: var(--shadow-card);
  height: 100%;

  .chart-header, .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0;
    }

    .chart-legend {
      display: flex;
      gap: 16px;

      .legend-item {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
        color: var(--text-secondary);

        .dot {
          width: 8px;
          height: 8px;
          border-radius: 50%;
        }
      }
    }
  }

  .chart-container {
    height: 280px;

    &.large {
      height: 350px;
    }
  }
}

.funnel-stats {
  display: flex;
  justify-content: space-around;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--border-light);

  .funnel-stat {
    text-align: center;

    .stat-label {
      display: block;
      font-size: 12px;
      color: var(--text-secondary);
      margin-bottom: 4px;
    }
    .stat-value {
      display: block;
      font-size: 20px;
      font-weight: 700;
      color: var(--text-primary);
    }
    .stat-rate {
      font-size: 12px;
      color: #10b981;
    }
  }
}

.channel-list {
  margin-top: 16px;

  .channel-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid var(--border-light);

    &:last-child {
      border-bottom: none;
    }

    .channel-info {
      display: flex;
      align-items: center;
      gap: 8px;

      .channel-dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;
      }
      .channel-name {
        font-size: 14px;
        color: var(--text-primary);
      }
    }

    .channel-stats {
      display: flex;
      gap: 16px;

      .channel-count {
        font-size: 14px;
        font-weight: 600;
        color: var(--text-primary);
      }
      .channel-rate {
        font-size: 14px;
        color: var(--text-secondary);
      }
    }
  }
}

.data-section {
  margin-top: 20px;

  .progress-cell {
    width: 100%;
  }
}

.interviewer-list {
  .interviewer-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px 0;
    border-bottom: 1px solid var(--border-light);

    &:last-child {
      border-bottom: none;
    }

    .rank {
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
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
      }
    }

    .interviewer-info {
      flex: 1;

      .name {
        display: block;
        font-size: 14px;
        font-weight: 600;
        color: var(--text-primary);
      }
      .dept {
        font-size: 12px;
        color: var(--text-secondary);
      }
    }

    .interviewer-stats {
      display: flex;
      gap: 24px;

      .stat {
        text-align: center;

        .value {
          display: block;
          font-size: 16px;
          font-weight: 700;
          color: var(--primary-color);
        }
        .label {
          font-size: 11px;
          color: var(--text-tertiary);
        }
      }
    }
  }
}
</style>
