<template>
  <div class="ai-process-flow">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>AI处理流程监控</h1>
      <p class="subtitle">实时查看简历解析、Embedding向量化、RAG检索增强的处理状态</p>
    </div>

    <!-- 当前处理任务（支持多个） -->
    <el-card class="current-task-card" v-if="currentTasks.length > 0">
      <template #header>
        <div class="card-header">
          <span>当前处理任务 ({{ currentTasks.length }}个)</span>
          <el-tag type="warning" effect="dark">
            <el-icon class="is-loading"><Loading /></el-icon>
            处理中
          </el-tag>
        </div>
      </template>
      
      <!-- 多任务列表 -->
      <div v-for="(task, index) in currentTasks" :key="task.resumeId" class="task-item" :class="{ 'task-divider': index > 0 }">
        <div class="task-info">
          <div class="task-file">
            <el-icon><Document /></el-icon>
            <span>{{ task.fileName }}</span>
          </div>
          <div class="task-job" v-if="task.jobTitle">
            <el-icon><Suitcase /></el-icon>
            <span>{{ task.jobTitle }}</span>
          </div>
        </div>

        <!-- 处理流程步骤 -->
        <el-steps :active="task.currentStep" finish-status="success" align-center class="process-steps">
          <el-step title="文件上传" :icon="Upload">
            <template #description>
              <span v-if="task.currentStep > 0" class="step-done">✓ 完成</span>
              <span v-else-if="task.currentStep === 0" class="step-processing">上传中...</span>
            </template>
          </el-step>
          <el-step title="OCR识别" :icon="Document">
            <template #description>
              <span v-if="task.currentStep > 1" class="step-done">✓ 完成</span>
              <span v-else-if="task.currentStep === 1" class="step-processing">识别中...</span>
            </template>
          </el-step>
          <el-step title="Embedding向量化" :icon="Connection">
            <template #description>
              <span v-if="task.currentStep > 2" class="step-done">✓ 完成</span>
              <span v-else-if="task.currentStep === 2" class="step-processing">向量化中...</span>
            </template>
          </el-step>
          <el-step title="RAG检索匹配" :icon="Search">
            <template #description>
              <span v-if="task.currentStep > 3" class="step-done">✓ 完成</span>
              <span v-else-if="task.currentStep === 3" class="step-processing">检索中...</span>
            </template>
          </el-step>
          <el-step title="AI智能评估" :icon="MagicStick">
            <template #description>
              <span v-if="task.currentStep > 4" class="step-done">✓ 完成</span>
              <span v-else-if="task.currentStep === 4" class="step-processing">评估中...</span>
            </template>
          </el-step>
        </el-steps>

        <!-- 当前步骤详情 -->
        <div class="step-detail" v-if="task.stepDetail">
          <el-alert :title="task.stepDetail" type="info" :closable="false" show-icon />
        </div>
      </div>
    </el-card>

    <!-- 无任务提示 -->
    <el-card v-else class="no-task-card">
      <el-empty description="当前没有正在处理的任务">
        <el-button type="primary" @click="goToResumeList">
          <el-icon><Upload /></el-icon>
          去上传简历
        </el-button>
      </el-empty>
    </el-card>

    <!-- 处理历史记录 -->
    <el-card class="history-card">
      <template #header>
        <div class="card-header">
          <span>处理历史记录</span>
          <el-button text type="primary" @click="refreshHistory">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-table :data="processHistory" v-loading="loading" style="width: 100%">
        <el-table-column prop="file_name" label="简历文件" min-width="200">
          <template #default="{ row }">
            <div class="file-cell">
              <el-icon><Document /></el-icon>
              <span>{{ row.resume_name || row.file_name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="job_title" label="求职职位" width="150" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="match_score" label="匹配度" width="120">
          <template #default="{ row }">
            <el-progress 
              v-if="row.match_score" 
              :percentage="Math.round(row.match_score)" 
              :color="getScoreColor(row.match_score)"
              :stroke-width="8"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="处理时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row)" v-if="row.status === 'completed'">
              <el-icon><View /></el-icon> 查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @current-change="fetchHistory"
          @size-change="fetchHistory"
        />
      </div>
    </el-card>

    <!-- 详情抽屉 -->
    <el-drawer v-model="showDetailDrawer" title="处理详情" size="600px">
      <div class="detail-content" v-if="currentDetail">
        <!-- 基本信息 -->
        <div class="detail-section">
          <h4>基本信息</h4>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="简历文件">{{ currentDetail.resume_name }}</el-descriptions-item>
            <el-descriptions-item label="求职职位">{{ currentDetail.job_title || '-' }}</el-descriptions-item>
            <el-descriptions-item label="候选人">{{ currentDetail.parsed_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="处理时间">{{ formatDate(currentDetail.created_at) }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- Embedding结果 -->
        <div class="detail-section">
          <h4><el-icon><Connection /></el-icon> Embedding向量化</h4>
          <div class="embedding-info">
            <el-tag type="success">向量维度: 1024</el-tag>
            <el-tag type="info">模型: doubao-embedding-large</el-tag>
          </div>
        </div>

        <!-- RAG检索结果 -->
        <div class="detail-section">
          <h4><el-icon><Search /></el-icon> RAG检索匹配</h4>
          <div class="rag-info">
            <p>检索到 <strong>{{ currentDetail.rag_matches || 5 }}</strong> 个相似职位</p>
            <p>最高相似度: <strong>{{ (currentDetail.max_similarity || 0.85) * 100 }}%</strong></p>
          </div>
        </div>

        <!-- AI评估结果 -->
        <div class="detail-section">
          <h4><el-icon><MagicStick /></el-icon> AI智能评估</h4>
          <div class="score-display">
            <el-progress 
              type="dashboard" 
              :percentage="Math.round(currentDetail.match_score || 0)" 
              :color="getScoreColor(currentDetail.match_score)"
              :width="120"
            />
          </div>
          
          <div class="dimensions" v-if="currentDetail.dimensions">
            <div v-for="dim in parseDimensions(currentDetail.report_dimensions)" :key="dim.name" class="dimension-item">
              <span class="dim-name">{{ dim.name }}</span>
              <el-progress :percentage="dim.score" :stroke-width="8" :color="getScoreColor(dim.score)" />
            </div>
          </div>

          <div class="recommendation" v-if="currentDetail.report_recommendation">
            <h5>AI建议</h5>
            <p>{{ currentDetail.report_recommendation }}</p>
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Document, Suitcase, Upload, Connection, Search, MagicStick,
  Loading, Refresh, View
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()

// 当前处理任务（支持多个）
const currentTasks = ref<any[]>([])

// 历史记录
const processHistory = ref<any[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 详情
const showDetailDrawer = ref(false)
const currentDetail = ref<any>(null)

// 轮询定时器
let pollTimer: any = null

// 获取处理历史
const fetchHistory = async () => {
  loading.value = true
  try {
    const res = await request.get('/evaluations', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    if (res.data?.code === 0) {
      processHistory.value = res.data.data?.evaluations || []
      total.value = res.data.data?.total || 0
    }
  } catch (error) {
    console.error('获取历史记录失败:', error)
  } finally {
    loading.value = false
  }
}

// 刷新历史
const refreshHistory = () => {
  currentPage.value = 1
  fetchHistory()
}

// 查看详情
const viewDetail = (row: any) => {
  currentDetail.value = row
  showDetailDrawer.value = true
}

// 跳转到简历列表
const goToResumeList = () => {
  router.push('/resumes')
}

// 解析维度数据
const parseDimensions = (dimensionsStr: string) => {
  if (!dimensionsStr) return []
  try {
    return JSON.parse(dimensionsStr)
  } catch {
    return []
  }
}

// 状态相关
const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success',
    failed: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    processing: '处理中',
    completed: '已完成',
    failed: '失败'
  }
  return map[status] || status
}

const getScoreColor = (score: number) => {
  if (score >= 80) return '#67c23a'
  if (score >= 60) return '#e6a23c'
  return '#f56c6c'
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

// 检查当前任务状态
const checkCurrentTask = async () => {
  try {
    const res = await request.get('/ai/current-task')
    if (res.data?.code === 0 && res.data.data) {
      // 支持多任务
      if (res.data.data.tasks) {
        currentTasks.value = res.data.data.tasks
      } else if (res.data.data.resumeId) {
        // 兼容单任务格式
        currentTasks.value = [res.data.data]
      } else {
        currentTasks.value = []
      }
    } else {
      currentTasks.value = []
    }
  } catch {
    currentTasks.value = []
  }
}

// 开始轮询
const startPolling = () => {
  pollTimer = setInterval(() => {
    checkCurrentTask()
  }, 2000)
}

onMounted(() => {
  fetchHistory()
  checkCurrentTask()
  startPolling()
})

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
  }
})
</script>

<style scoped lang="scss">
.ai-process-flow {
  padding: 24px;
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.page-header {
  margin-bottom: 24px;

  h1 {
    font-size: 24px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0 0 8px 0;
  }

  .subtitle {
    color: var(--text-secondary);
    font-size: 14px;
    margin: 0;
  }
}

.current-task-card {
  margin-bottom: 24px;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .task-item {
    &.task-divider {
      margin-top: 24px;
      padding-top: 24px;
      border-top: 1px dashed var(--border-color);
    }
  }

  .task-info {
    display: flex;
    gap: 24px;
    margin-bottom: 24px;

    .task-file, .task-job {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 15px;
      color: var(--text-primary);
    }
  }

  .process-steps {
    margin: 24px 0;
  }

  .step-done {
    color: #67c23a;
    font-size: 12px;
  }

  .step-processing {
    color: #409eff;
    font-size: 12px;
  }

  .step-detail {
    margin-top: 16px;
  }
}

.no-task-card {
  margin-bottom: 24px;
}

.history-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .file-cell {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

// 详情抽屉
.detail-content {
  .detail-section {
    margin-bottom: 24px;

    h4 {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 16px 0;
      padding-bottom: 8px;
      border-bottom: 1px solid var(--border-color);
    }

    h5 {
      font-size: 14px;
      font-weight: 600;
      margin: 16px 0 8px 0;
    }
  }

  .embedding-info, .rag-info {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;

    p {
      margin: 4px 0;
      color: var(--text-secondary);
    }
  }

  .score-display {
    text-align: center;
    margin: 16px 0;
  }

  .dimensions {
    .dimension-item {
      margin-bottom: 12px;

      .dim-name {
        display: block;
        font-size: 13px;
        color: var(--text-secondary);
        margin-bottom: 4px;
      }
    }
  }

  .recommendation {
    background: var(--bg-secondary);
    padding: 16px;
    border-radius: 8px;

    p {
      margin: 0;
      color: var(--text-primary);
      line-height: 1.6;
    }
  }
}
</style>
