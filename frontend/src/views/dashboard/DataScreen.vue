<template>
  <div class="data-screen">
    <!-- 头部 -->
    <header class="screen-header">
      <div class="header-left">
        <span class="time">{{ currentTime }}</span>
      </div>
      <h1 class="title">智能人才招聘管理平台 - 数据大屏</h1>
      <div class="header-right">
        <span class="date">{{ currentDate }}</span>
        <el-button type="primary" size="small" @click="$router.push('/dashboard')">返回</el-button>
      </div>
    </header>

    <!-- 主体内容 -->
    <main class="screen-main">
      <!-- 左侧 -->
      <div class="screen-left">
        <!-- 招聘概览 -->
        <div class="panel">
          <div class="panel-title">招聘概览</div>
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-value">{{ stats.totalTalents }}</div>
              <div class="stat-label">人才总数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ stats.totalJobs }}</div>
              <div class="stat-label">在招职位</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ stats.totalInterviews }}</div>
              <div class="stat-label">面试安排</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ stats.totalApplications }}</div>
              <div class="stat-label">简历投递</div>
            </div>
          </div>
        </div>

        <!-- 招聘漏斗 -->
        <div class="panel">
          <div class="panel-title">招聘漏斗</div>
          <div ref="funnelChart" class="chart-container"></div>
        </div>

        <!-- 职位类型分布 -->
        <div class="panel">
          <div class="panel-title">职位类型分布</div>
          <div ref="jobTypeChart" class="chart-container"></div>
        </div>
      </div>

      <!-- 中间 -->
      <div class="screen-center">
        <!-- 地图 -->
        <div class="panel map-panel">
          <div class="panel-title">人才地域分布</div>
          <div ref="mapChart" class="map-container"></div>
        </div>

        <!-- 实时动态 -->
        <div class="panel">
          <div class="panel-title">实时动态</div>
          <div class="activity-list">
            <div v-for="(item, index) in activities" :key="index" class="activity-item">
              <span class="activity-icon" :style="{ background: item.color }">
                {{ item.icon }}
              </span>
              <span class="activity-text">{{ item.text }}</span>
              <span class="activity-time">{{ item.time }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧 -->
      <div class="screen-right">
        <!-- 招聘趋势 -->
        <div class="panel">
          <div class="panel-title">招聘趋势（近7天）</div>
          <div ref="trendChart" class="chart-container"></div>
        </div>

        <!-- 热门技能 -->
        <div class="panel">
          <div class="panel-title">热门技能需求</div>
          <div ref="skillChart" class="chart-container"></div>
        </div>

        <!-- 面试官排行 -->
        <div class="panel">
          <div class="panel-title">面试官排行</div>
          <div class="rank-list">
            <div v-for="(item, index) in interviewerRank" :key="index" class="rank-item">
              <span class="rank-num" :class="{ top: index < 3 }">{{ index + 1 }}</span>
              <span class="rank-name">{{ item.name }}</span>
              <span class="rank-count">{{ item.count }}场</span>
              <div class="rank-bar">
                <div class="rank-bar-inner" :style="{ width: item.percent + '%' }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'

const currentTime = ref('')
const currentDate = ref('')
let timer: number | null = null

// 统计数据
const stats = ref({
  totalTalents: 156,
  totalJobs: 48,
  totalInterviews: 32,
  totalApplications: 289
})

// 实时动态
const activities = ref([
  { icon: '📄', text: '张三 投递了 高级Go开发工程师', time: '刚刚', color: '#0ea5e9' },
  { icon: '📅', text: '李四 的面试已安排', time: '2分钟前', color: '#22c55e' },
  { icon: '✅', text: '王五 通过了技术面试', time: '5分钟前', color: '#8b5cf6' },
  { icon: '🎉', text: '赵六 已发送Offer', time: '10分钟前', color: '#f59e0b' },
  { icon: '📄', text: '钱七 投递了 前端架构师', time: '15分钟前', color: '#0ea5e9' },
])

// 面试官排行
const interviewerRank = ref([
  { name: '陈强', count: 28, percent: 100 },
  { name: '王芳', count: 24, percent: 86 },
  { name: '李明', count: 20, percent: 71 },
  { name: '张伟', count: 16, percent: 57 },
  { name: '刘洋', count: 12, percent: 43 },
])

// 图表引用
const funnelChart = ref<HTMLElement>()
const jobTypeChart = ref<HTMLElement>()
const mapChart = ref<HTMLElement>()
const trendChart = ref<HTMLElement>()
const skillChart = ref<HTMLElement>()

let charts: echarts.ECharts[] = []

// 更新时间
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour12: false })
  currentDate.value = now.toLocaleDateString('zh-CN', { 
    year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' 
  })
}

// 初始化漏斗图
const initFunnelChart = () => {
  if (!funnelChart.value) return
  const chart = echarts.init(funnelChart.value)
  charts.push(chart)
  
  chart.setOption({
    color: ['#0ea5e9', '#06b6d4', '#22c55e', '#f59e0b', '#8b5cf6'],
    tooltip: { trigger: 'item', formatter: '{b}: {c}' },
    series: [{
      type: 'funnel',
      left: '5%',
      width: '50%',
      minSize: '20%',
      label: { 
        show: true, 
        position: 'inside', 
        color: '#fff',
        fontSize: 11,
        formatter: '{b}'
      },
      labelLine: { show: false },
      data: [
        { value: 289, name: '简历投递' },
        { value: 180, name: '简历筛选' },
        { value: 98, name: '面试邀约' },
        { value: 52, name: '面试通过' },
        { value: 28, name: '发送Offer' }
      ]
    }],
    graphic: [
      { type: 'text', left: '62%', top: '12%', style: { text: '289', fill: '#0ea5e9', fontSize: 14, fontWeight: 'bold' } },
      { type: 'text', left: '62%', top: '30%', style: { text: '180', fill: '#06b6d4', fontSize: 14, fontWeight: 'bold' } },
      { type: 'text', left: '62%', top: '48%', style: { text: '98', fill: '#22c55e', fontSize: 14, fontWeight: 'bold' } },
      { type: 'text', left: '62%', top: '66%', style: { text: '52', fill: '#f59e0b', fontSize: 14, fontWeight: 'bold' } },
      { type: 'text', left: '62%', top: '84%', style: { text: '28', fill: '#8b5cf6', fontSize: 14, fontWeight: 'bold' } }
    ]
  })
}

// 初始化职位类型图
const initJobTypeChart = () => {
  if (!jobTypeChart.value) return
  const chart = echarts.init(jobTypeChart.value)
  charts.push(chart)
  
  const data = [
    { value: 18, name: '技术研发' },
    { value: 12, name: '产品设计' },
    { value: 8, name: '市场运营' },
    { value: 6, name: '人力行政' },
    { value: 4, name: '财务法务' }
  ]
  
  chart.setOption({
    color: ['#0ea5e9', '#06b6d4', '#22c55e', '#f59e0b', '#8b5cf6'],
    tooltip: { trigger: 'item', formatter: '{b}: {c}个 ({d}%)' },
    legend: {
      orient: 'vertical',
      right: '2%',
      top: 'center',
      textStyle: { color: '#94a3b8', fontSize: 11 },
      itemWidth: 10,
      itemHeight: 10
    },
    series: [{
      type: 'pie',
      radius: ['30%', '55%'],
      center: ['35%', '50%'],
      label: { 
        show: true,
        position: 'inside',
        color: '#fff',
        fontSize: 11,
        fontWeight: 'bold',
        formatter: '{c}个'
      },
      labelLine: { show: false },
      data: data
    }]
  })
}

// 初始化地图
const initMapChart = () => {
  if (!mapChart.value) return
  const chart = echarts.init(mapChart.value)
  charts.push(chart)
  
  // 简化版：用散点图模拟地图
  chart.setOption({
    backgroundColor: 'transparent',
    tooltip: { trigger: 'item' },
    geo: {
      map: 'china',
      roam: false,
      silent: true,
      itemStyle: { areaColor: '#1a3a5c', borderColor: '#0ea5e9' }
    },
    series: [{
      type: 'scatter',
      coordinateSystem: 'geo',
      symbolSize: (val: number[]) => val[2] / 2,
      itemStyle: { color: '#0ea5e9', shadowBlur: 10, shadowColor: '#0ea5e9' },
      data: [
        { name: '北京', value: [116.46, 39.92, 85] },
        { name: '上海', value: [121.48, 31.22, 72] },
        { name: '深圳', value: [114.07, 22.62, 58] },
        { name: '杭州', value: [120.19, 30.26, 45] },
        { name: '广州', value: [113.23, 23.16, 38] },
        { name: '成都', value: [104.06, 30.67, 32] }
      ]
    }]
  })
}

// 初始化趋势图
const initTrendChart = () => {
  if (!trendChart.value) return
  const chart = echarts.init(trendChart.value)
  charts.push(chart)
  
  chart.setOption({
    tooltip: { trigger: 'axis' },
    legend: {
      data: ['简历投递', '面试安排'],
      bottom: 0,
      textStyle: { color: '#94a3b8', fontSize: 11 }
    },
    grid: { left: '3%', right: '4%', bottom: '15%', top: '10%', containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      axisLine: { lineStyle: { color: '#0ea5e9' } },
      axisLabel: { color: '#94a3b8' }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      splitLine: { lineStyle: { color: '#1e3a5f' } },
      axisLabel: { color: '#94a3b8' }
    },
    series: [
      {
        name: '简历投递',
        type: 'line',
        smooth: true,
        areaStyle: { color: 'rgba(14, 165, 233, 0.3)' },
        lineStyle: { color: '#0ea5e9', width: 2 },
        itemStyle: { color: '#0ea5e9' },
        label: {
          show: true,
          position: 'top',
          color: '#0ea5e9',
          fontSize: 10,
          fontWeight: 'bold'
        },
        data: [32, 45, 38, 52, 48, 28, 35]
      },
      {
        name: '面试安排',
        type: 'line',
        smooth: true,
        areaStyle: { color: 'rgba(34, 197, 94, 0.3)' },
        lineStyle: { color: '#22c55e', width: 2 },
        itemStyle: { color: '#22c55e' },
        label: {
          show: true,
          position: 'top',
          color: '#22c55e',
          fontSize: 10,
          fontWeight: 'bold'
        },
        data: [12, 18, 15, 22, 20, 10, 14]
      }
    ]
  })
}

// 初始化技能图
const initSkillChart = () => {
  if (!skillChart.value) return
  const chart = echarts.init(skillChart.value)
  charts.push(chart)
  
  chart.setOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: '3%', right: '12%', bottom: '3%', top: '3%', containLabel: true },
    xAxis: {
      type: 'value',
      axisLine: { show: false },
      splitLine: { lineStyle: { color: '#1e3a5f' } },
      axisLabel: { color: '#94a3b8' }
    },
    yAxis: {
      type: 'category',
      data: ['Python', 'Java', 'Go', 'Vue', 'React'],
      axisLine: { lineStyle: { color: '#0ea5e9' } },
      axisLabel: { color: '#94a3b8' }
    },
    series: [{
      type: 'bar',
      barWidth: 16,
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
          { offset: 0, color: '#0ea5e9' },
          { offset: 1, color: '#22c55e' }
        ]),
        borderRadius: [0, 8, 8, 0]
      },
      label: {
        show: true,
        position: 'right',
        color: '#0ea5e9',
        fontSize: 11,
        fontWeight: 'bold',
        formatter: '{c}人'
      },
      data: [45, 42, 38, 35, 32]
    }]
  })
}

// 窗口大小变化
const handleResize = () => {
  charts.forEach(chart => chart.resize())
}

onMounted(async () => {
  updateTime()
  timer = window.setInterval(updateTime, 1000)
  
  // 动态加载中国地图
  try {
    const response = await fetch('https://geo.datav.aliyun.com/areas_v3/bound/100000_full.json')
    const chinaJson = await response.json()
    echarts.registerMap('china', chinaJson)
  } catch (e) {
    console.log('地图加载失败，使用简化版')
  }
  
  setTimeout(() => {
    initFunnelChart()
    initJobTypeChart()
    initMapChart()
    initTrendChart()
    initSkillChart()
  }, 100)
  
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
  charts.forEach(chart => chart.dispose())
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped lang="scss">
.data-screen {
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #0a1628 0%, #1a2a4a 100%);
  color: #fff;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.screen-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: rgba(14, 165, 233, 0.1);
  border-bottom: 1px solid rgba(14, 165, 233, 0.3);

  .title {
    font-size: 24px;
    font-weight: 600;
    background: linear-gradient(90deg, #0ea5e9, #22c55e);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .time, .date {
    color: #94a3b8;
    font-size: 14px;
  }
}

.screen-main {
  flex: 1;
  display: flex;
  padding: 16px;
  gap: 16px;
  overflow: hidden;
}

.screen-left, .screen-right {
  width: 25%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.screen-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.panel {
  background: rgba(14, 165, 233, 0.05);
  border: 1px solid rgba(14, 165, 233, 0.2);
  border-radius: 8px;
  padding: 16px;
  flex: 1;
  display: flex;
  flex-direction: column;

  .panel-title {
    font-size: 14px;
    font-weight: 600;
    color: #0ea5e9;
    margin-bottom: 12px;
    padding-left: 8px;
    border-left: 3px solid #0ea5e9;
  }
}

.map-panel {
  flex: 2;
}

.chart-container {
  flex: 1;
  min-height: 0;
}

.map-container {
  flex: 1;
  min-height: 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;

  .stat-item {
    text-align: center;
    padding: 12px;
    background: rgba(14, 165, 233, 0.1);
    border-radius: 8px;

    .stat-value {
      font-size: 28px;
      font-weight: 700;
      color: #0ea5e9;
    }

    .stat-label {
      font-size: 12px;
      color: #94a3b8;
      margin-top: 4px;
    }
  }
}

.activity-list {
  flex: 1;
  overflow-y: auto;

  .activity-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 0;
    border-bottom: 1px solid rgba(14, 165, 233, 0.1);

    .activity-icon {
      width: 28px;
      height: 28px;
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
    }

    .activity-text {
      flex: 1;
      font-size: 13px;
      color: #e2e8f0;
    }

    .activity-time {
      font-size: 11px;
      color: #64748b;
    }
  }
}

.rank-list {
  flex: 1;

  .rank-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 0;

    .rank-num {
      width: 20px;
      height: 20px;
      border-radius: 4px;
      background: #1e3a5f;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      font-weight: 600;

      &.top {
        background: linear-gradient(135deg, #f59e0b, #ef4444);
      }
    }

    .rank-name {
      width: 50px;
      font-size: 13px;
    }

    .rank-count {
      width: 40px;
      font-size: 12px;
      color: #0ea5e9;
    }

    .rank-bar {
      flex: 1;
      height: 6px;
      background: #1e3a5f;
      border-radius: 3px;
      overflow: hidden;

      .rank-bar-inner {
        height: 100%;
        background: linear-gradient(90deg, #0ea5e9, #22c55e);
        border-radius: 3px;
      }
    }
  }
}
</style>
