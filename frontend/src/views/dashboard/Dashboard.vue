<template>
  <div class="dashboard">
    <h1 class="page-title">仪表板</h1>
    
    <!-- Stats Cards -->
    <el-row :gutter="24" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card card" style="border-left: 4px solid #667eea">
          <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.totalTalents }}</div>
            <div class="stat-label">总人才数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card card" style="border-left: 4px solid #f093fb">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%)">
            <el-icon><Suitcase /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.openJobs }}</div>
            <div class="stat-label">在招职位</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card card" style="border-left: 4px solid #4facfe">
          <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.totalApplications }}</div>
            <div class="stat-label">申请总数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card card" style="border-left: 4px solid #43e97b">
          <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)">
            <el-icon><MagicStick /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.successfulMatches }}</div>
            <div class="stat-label">成功匹配</div>
          </div>
        </div>
      </el-col>
    </el-row>
    
    <!-- Charts Row -->
    <el-row :gutter="24" class="charts-row">
      <el-col :xs="24" :md="12">
        <div class="card chart-card">
          <h3>职位统计</h3>
          <div ref="jobChartRef" style="width: 100%; height: 300px"></div>
        </div>
      </el-col>
      
      <el-col :xs="24" :md="12">
        <div class="card chart-card">
          <h3>申请趋势</h3>
          <div ref="trendChartRef" style="width: 100%; height: 300px"></div>
        </div>
      </el-col>
    </el-row>
    
    <!-- Recent Activities -->
    <div class="card recent-activities">
      <h3>最近活动</h3>
      <el-timeline>
        <el-timeline-item
          v-for="(activity, index) in recentActivities"
          :key="index"
          :timestamp="activity.time"
          placement="top"
        >
          {{ activity.text }}
        </el-timeline-item>
      </el-timeline>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'
import { jobApi } from '@/api/job'
import { recommendationApi } from '@/api/recommendation'

const jobChartRef = ref<HTMLElement>()
const trendChartRef = ref<HTMLElement>()

const stats = ref({
  totalTalents: 156,
  openJobs: 23,
  totalApplications: 89,
  successfulMatches: 45
})

const recentActivities = ref([
  { time: '5分钟前', text: '新增人才：张三' },
  { time: '15分钟前', text: '职位发布：高级Go开发工程师' },
  { time: '1小时前', text: '简历审核：李四的简历已通过' },
  { time: '2小时前', text: '智能推荐：5个匹配候选人' },
  { time: '3小时前', text: '面试安排：王五 - 前端开发工程师' }
])

const initJobChart = () => {
  if (!jobChartRef.value) return
  
  const chart = echarts.init(jobChartRef.value)
  const option = {
    tooltip: {
      trigger: 'item'
    },
    legend: {
      bottom: '0%'
    },
    series: [
      {
        name: '职位状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 20,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 23, name: '在招', itemStyle: { color: '#667eea' } },
          { value: 12, name: '已关闭', itemStyle: { color: '#f093fb' } },
          { value: 8, name: '已填补', itemStyle: { color: '#43e97b' } }
        ]
      }
    ]
  }
  chart.setOption(option)
}

const initTrendChart = () => {
  if (!trendChartRef.value) return
  
  const chart = echarts.init(trendChartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '申请数',
        type: 'bar',
        data: [12, 15, 8, 20, 18, 10, 14],
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#667eea' },
            { offset: 1, color: '#764ba2' }
          ])
        }
      }
    ]
  }
  chart.setOption(option)
}

const fetchStats = async () => {
  try {
    const [jobRes, recommendRes] = await Promise.all([
      jobApi.getStats(),
      recommendationApi.getStats()
    ])
    
    if (jobRes.data.code === 0 && jobRes.data.data) {
      stats.value.openJobs = jobRes.data.data.open_jobs || 0
      stats.value.totalApplications = jobRes.data.data.total_jobs || 0
    }
    
    if (recommendRes.data.code === 0 && recommendRes.data.data) {
      stats.value.successfulMatches = recommendRes.data.data.successful_matches || 0
    }
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

onMounted(() => {
  initJobChart()
  initTrendChart()
  fetchStats()
})
</script>

<style scoped lang="scss">
.dashboard {
  .page-title {
    font-size: 28px;
    font-weight: 700;
    margin-bottom: 24px;
    color: var(--text-primary);
  }
}

.stats-row {
  margin-bottom: 24px;
  
  .stat-card {
    display: flex;
    align-items: center;
    gap: 20px;
    padding: 20px;
    margin-bottom: 16px;
    
    .stat-icon {
      width: 60px;
      height: 60px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 28px;
    }
    
    .stat-content {
      flex: 1;
      
      .stat-value {
        font-size: 32px;
        font-weight: 700;
        color: var(--text-primary);
        line-height: 1;
        margin-bottom: 8px;
      }
      
      .stat-label {
        font-size: 14px;
        color: var(--text-secondary);
      }
    }
  }
}

.charts-row {
  margin-bottom: 24px;
}

.chart-card {
  h3 {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 20px;
    color: var(--text-primary);
  }
}

.recent-activities {
  h3 {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 20px;
    color: var(--text-primary);
  }
  
  :deep(.el-timeline-item__timestamp) {
    color: var(--text-secondary);
    font-size: 12px;
  }
}
</style>
