<template>
  <div class="evaluation-results-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-left">
        <h1>评估结果管理</h1>
        <p class="subtitle">查看和管理所有AI评估记录</p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon total">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.total }}</div>
            <div class="stat-label">总评估数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon high">
            <el-icon><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.high_match }}</div>
            <div class="stat-label">高匹配</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon medium">
            <el-icon><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.medium_match }}</div>
            <div class="stat-label">中匹配</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon score">
            <el-icon><TrendCharts /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.avg_match_score?.toFixed(1) || 0 }}</div>
            <div class="stat-label">平均分</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 筛选和搜索 -->
    <el-card class="filter-card">
      <el-row :gutter="16">
        <el-col :span="6">
          <el-input
            v-model="filters.search"
            placeholder="搜索姓名/简历名称"
            clearable
            @clear="loadData"
            @keyup.enter="loadData"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.status" placeholder="状态" clearable @change="loadData">
            <el-option label="已完成" value="completed" />
            <el-option label="处理中" value="pending" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.evalType" placeholder="评估类型" clearable @change="loadData">
            <el-option label="简历解析" value="parse" />
            <el-option label="人岗匹配" value="match" />
            <el-option label="完整评估" value="full" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="loadData">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card">
      <el-table
        :data="evaluations"
        v-loading="loading"
        stripe
        @sort-change="handleSortChange"
      >
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="resume_name" label="简历名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="parsed_name" label="候选人" width="100" />
        <el-table-column label="匹配度" width="120" sortable="custom" prop="match_score">
          <template #default="{ row }">
            <div class="score-cell">
              <el-progress
                :percentage="row.match_score"
                :color="getScoreColor(row.match_score)"
                :stroke-width="8"
                :show-text="false"
                style="width: 60px"
              />
              <span class="score-text">{{ row.match_score?.toFixed(1) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="匹配等级" width="100">
          <template #default="{ row }">
            <el-tag :type="getMatchLevelType(row.match_level)" size="small">
              {{ getMatchLevelText(row.match_level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="风险分" width="90" sortable="custom" prop="risk_score">
          <template #default="{ row }">
            <el-tag :type="getRiskType(row.risk_score)" size="small">
              {{ row.risk_score || 0 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="parsed_education" label="学历" width="80" />
        <el-table-column prop="parsed_experience" label="经验" width="80" />
        <el-table-column label="技能" min-width="150">
          <template #default="{ row }">
            <div class="skills-cell">
              <el-tag
                v-for="skill in (row.parsed_skills || []).slice(0, 3)"
                :key="skill"
                size="small"
                class="skill-tag"
              >
                {{ skill }}
              </el-tag>
              <span v-if="(row.parsed_skills || []).length > 3" class="more-skills">
                +{{ row.parsed_skills.length - 3 }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="评估时间" width="160" sortable="custom" prop="created_at">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row)">
              <el-icon><View /></el-icon>
              详情
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- 详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      title="评估详情"
      width="800px"
      destroy-on-close
    >
      <div v-if="currentEval" class="eval-detail">
        <!-- 基本信息 -->
        <el-descriptions title="候选人信息" :column="3" border>
          <el-descriptions-item label="姓名">{{ currentEval.parsed_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="手机">{{ currentEval.parsed_phone || '-' }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ currentEval.parsed_email || '-' }}</el-descriptions-item>
          <el-descriptions-item label="学历">{{ currentEval.parsed_education || '-' }}</el-descriptions-item>
          <el-descriptions-item label="经验">{{ currentEval.parsed_experience || '-' }}</el-descriptions-item>
          <el-descriptions-item label="城市">{{ currentEval.parsed_location || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- 技能标签 -->
        <div class="detail-section">
          <h4>技能标签</h4>
          <div class="skills-wrapper">
            <el-tag
              v-for="skill in currentEval.parsed_skills || []"
              :key="skill"
              class="skill-tag"
            >
              {{ skill }}
            </el-tag>
            <span v-if="!currentEval.parsed_skills?.length">暂无</span>
          </div>
        </div>

        <!-- 匹配结果 -->
        <div class="detail-section">
          <h4>匹配分析</h4>
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="score-display">
                <el-progress
                  type="dashboard"
                  :percentage="currentEval.match_score || 0"
                  :color="getScoreColor(currentEval.match_score)"
                  :width="120"
                />
                <div class="score-label">匹配度</div>
              </div>
            </el-col>
            <el-col :span="16">
              <div v-if="currentEval.report_dimensions" class="dimensions">
                <div
                  v-for="dim in currentEval.report_dimensions"
                  :key="dim.name"
                  class="dimension-item"
                >
                  <span class="dim-name">{{ dim.name }}</span>
                  <el-progress
                    :percentage="dim.score"
                    :stroke-width="8"
                    :color="getScoreColor(dim.score)"
                  />
                </div>
              </div>
            </el-col>
          </el-row>
        </div>

        <!-- 归因报告 -->
        <div v-if="currentEval.report_summary" class="detail-section">
          <h4>AI评估报告</h4>
          <p class="report-summary">{{ currentEval.report_summary }}</p>

          <el-row :gutter="16">
            <el-col :span="12">
              <div class="report-box strengths">
                <h5><el-icon><CircleCheck /></el-icon> 优势</h5>
                <ul>
                  <li v-for="(item, i) in currentEval.report_strengths || []" :key="i">{{ item }}</li>
                </ul>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="report-box gaps">
                <h5><el-icon><Warning /></el-icon> 待提升</h5>
                <ul>
                  <li v-for="(item, i) in currentEval.report_gaps || []" :key="i">{{ item }}</li>
                </ul>
              </div>
            </el-col>
          </el-row>

          <div v-if="currentEval.report_recommendation" class="recommendation-box">
            <h5><el-icon><ChatLineSquare /></el-icon> AI建议</h5>
            <p>{{ currentEval.report_recommendation }}</p>
          </div>
        </div>

        <!-- 风控提示 -->
        <div v-if="currentEval.risk_items?.length" class="detail-section">
          <h4>风控提示</h4>
          <el-alert
            v-for="(risk, i) in currentEval.risk_items"
            :key="i"
            :title="risk.message"
            :type="risk.level === 'high' ? 'error' : risk.level === 'warning' ? 'warning' : 'info'"
            :closable="false"
            show-icon
            class="risk-alert"
          />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Document, CircleCheck, Warning, TrendCharts, Search,
  View, Delete, ChatLineSquare
} from '@element-plus/icons-vue'
import request from '@/utils/request'

// 数据
const loading = ref(false)
const evaluations = ref<any[]>([])
const stats = ref<any>({})
const detailVisible = ref(false)
const currentEval = ref<any>(null)

// 筛选条件
const filters = reactive({
  search: '',
  status: '',
  evalType: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 排序
const sortParams = reactive({
  sortBy: 'created_at',
  sortOrder: 'desc'
})

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/evaluations', {
      params: {
        page: pagination.page,
        page_size: pagination.pageSize,
        search: filters.search,
        status: filters.status,
        eval_type: filters.evalType,
        sort_by: sortParams.sortBy,
        sort_order: sortParams.sortOrder
      }
    })

    if (res.code === 0) {
      evaluations.value = res.data.evaluations || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载数据失败', e)
  } finally {
    loading.value = false
  }
}

// 加载统计
const loadStats = async () => {
  try {
    const res = await request.get('/api/v1/evaluations/stats')
    if (res.code === 0) {
      stats.value = res.data
    }
  } catch (e) {
    console.error('加载统计失败', e)
  }
}

// 重置筛选
const resetFilters = () => {
  filters.search = ''
  filters.status = ''
  filters.evalType = ''
  pagination.page = 1
  loadData()
}

// 排序变化
const handleSortChange = ({ prop, order }: any) => {
  sortParams.sortBy = prop || 'created_at'
  sortParams.sortOrder = order === 'ascending' ? 'asc' : 'desc'
  loadData()
}

// 查看详情
const viewDetail = (row: any) => {
  currentEval.value = row
  detailVisible.value = true
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评估记录吗？', '提示', {
      type: 'warning'
    })

    const res = await request.delete(`/api/v1/evaluations/${row.id}`)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      loadData()
      loadStats()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (e) {
    // 取消删除
  }
}

// 工具函数
const getScoreColor = (score: number) => {
  if (score >= 80) return '#67c23a'
  if (score >= 60) return '#e6a23c'
  return '#f56c6c'
}

const getMatchLevelType = (level: string) => {
  const map: Record<string, string> = {
    high: 'success',
    medium: 'warning',
    low: 'danger'
  }
  return map[level] || 'info'
}

const getMatchLevelText = (level: string) => {
  const map: Record<string, string> = {
    high: '高匹配',
    medium: '中匹配',
    low: '低匹配'
  }
  return map[level] || '-'
}

const getRiskType = (score: number) => {
  if (score >= 50) return 'danger'
  if (score >= 20) return 'warning'
  return 'success'
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    completed: 'success',
    pending: 'warning',
    failed: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    completed: '已完成',
    pending: '处理中',
    failed: '失败'
  }
  return map[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadData()
  loadStats()
})
</script>

<style scoped lang="scss">
.evaluation-results-page {
  .page-header {
    margin-bottom: 20px;

    h1 {
      margin: 0 0 8px 0;
      font-size: 24px;
      color: var(--text-primary);
    }

    .subtitle {
      margin: 0;
      color: var(--text-secondary);
      font-size: 14px;
    }
  }

  .stats-row {
    margin-bottom: 20px;

    .stat-card {
      display: flex;
      align-items: center;
      padding: 16px;

      .stat-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 16px;

        .el-icon {
          font-size: 24px;
          color: white;
        }

        &.total { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
        &.high { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
        &.medium { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
        &.score { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
      }

      .stat-info {
        .stat-value {
          font-size: 24px;
          font-weight: bold;
          color: var(--text-primary);
        }

        .stat-label {
          font-size: 13px;
          color: var(--text-secondary);
        }
      }
    }
  }

  .filter-card {
    margin-bottom: 20px;

    .el-select {
      width: 100%;
    }
  }

  .table-card {
    .score-cell {
      display: flex;
      align-items: center;
      gap: 8px;

      .score-text {
        font-weight: 500;
        min-width: 30px;
      }
    }

    .skills-cell {
      display: flex;
      flex-wrap: wrap;
      gap: 4px;
      align-items: center;

      .skill-tag {
        margin: 0;
      }

      .more-skills {
        font-size: 12px;
        color: var(--text-secondary);
      }
    }

    .pagination-wrapper {
      margin-top: 20px;
      display: flex;
      justify-content: flex-end;
    }
  }

  .eval-detail {
    .detail-section {
      margin-top: 24px;
      padding-top: 16px;
      border-top: 1px solid var(--border-color);

      h4 {
        margin: 0 0 16px 0;
        font-size: 15px;
        color: var(--text-primary);
      }

      h5 {
        margin: 0 0 12px 0;
        font-size: 14px;
        display: flex;
        align-items: center;
        gap: 6px;
      }
    }

    .skills-wrapper {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;

      .skill-tag {
        margin: 0;
      }
    }

    .score-display {
      text-align: center;

      .score-label {
        margin-top: 8px;
        font-size: 14px;
        color: var(--text-secondary);
      }
    }

    .dimensions {
      .dimension-item {
        margin-bottom: 12px;

        .dim-name {
          display: block;
          margin-bottom: 4px;
          font-size: 13px;
          color: var(--text-secondary);
        }
      }
    }

    .report-summary {
      margin: 0 0 16px 0;
      padding: 12px;
      background: var(--bg-secondary);
      border-radius: 8px;
      color: var(--text-primary);
      line-height: 1.6;
    }

    .report-box {
      background: var(--bg-secondary);
      border-radius: 8px;
      padding: 12px;

      &.strengths h5 { color: #67c23a; }
      &.gaps h5 { color: #e6a23c; }

      ul {
        margin: 0;
        padding-left: 20px;

        li {
          margin-bottom: 6px;
          color: var(--text-primary);
        }
      }
    }

    .recommendation-box {
      margin-top: 16px;
      padding: 12px;
      background: linear-gradient(135deg, #667eea10 0%, #764ba210 100%);
      border-radius: 8px;

      h5 { color: #667eea; }

      p {
        margin: 0;
        color: var(--text-primary);
        line-height: 1.6;
      }
    }

    .risk-alert {
      margin-bottom: 8px;
    }
  }
}
</style>
