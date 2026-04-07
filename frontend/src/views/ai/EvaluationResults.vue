<template>
  <div class="evaluation-results-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-left">
        <h1>评估结果管理</h1>
        <p class="subtitle">查看每份简历最新的AI评估结果</p>
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
            <div class="stat-label">已评估简历数</div>
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
        <el-table-column prop="id" label="评估ID" width="90" />
        <el-table-column prop="resume_id" label="简历ID" width="90" />
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
        <el-table-column label="最新评估时间" width="160" sortable="custom" prop="created_at">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="openResume(row)">
              <el-icon><Document /></el-icon>
              简历
            </el-button>
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
          <el-descriptions-item label="学校">{{ currentEval.parsed_school || '-' }}</el-descriptions-item>
          <el-descriptions-item label="经验">{{ currentEval.parsed_experience || '-' }}</el-descriptions-item>
          <el-descriptions-item label="城市">{{ currentEval.parsed_location || '-' }}</el-descriptions-item>
          <el-descriptions-item label="评级">{{ currentEval.grade || '-' }}</el-descriptions-item>
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

        <div v-if="currentEval.jd_summary || currentEval.jd_matched_skills?.length || currentEval.jd_missing_skills?.length" class="detail-section">
          <h4>JD 匹配详细分析</h4>
          <p v-if="currentEval.jd_summary" class="report-summary">{{ currentEval.jd_summary }}</p>
          <el-row :gutter="16">
            <el-col :span="12">
              <div class="report-box strengths">
                <h5><el-icon><CircleCheck /></el-icon> 匹配到的能力</h5>
                <ul>
                  <li v-for="(item, i) in currentEval.jd_matched_skills || []" :key="`jd-matched-${i}`">{{ item }}</li>
                </ul>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="report-box gaps">
                <h5><el-icon><Warning /></el-icon> 缺失或待确认</h5>
                <ul>
                  <li v-for="(item, i) in currentEval.jd_missing_skills || []" :key="`jd-missing-${i}`">{{ item }}</li>
                </ul>
              </div>
            </el-col>
          </el-row>
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

        <div v-if="currentEval.dimension_details?.length" class="detail-section">
          <h4>各维度详细说明</h4>
          <div class="dimension-detail-list">
            <div v-for="item in currentEval.dimension_details" :key="item.name" class="dimension-detail-item">
              <div class="dimension-detail-head">
                <span class="dimension-title">{{ item.name }}</span>
                <span class="dimension-score">{{ item.raw_score }}/{{ item.max_score }}</span>
              </div>
              <p class="dimension-description">{{ item.description || '暂无说明' }}</p>
            </div>
          </div>
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

        <div v-if="currentEval.recommendation_conclusion || currentEval.recommendation_reason || currentEval.salary_suggestion || currentEval.suitable_roles?.length || currentEval.interview_focus?.length" class="detail-section">
          <h4>录用建议</h4>
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="结论">{{ currentEval.recommendation_conclusion || '-' }}</el-descriptions-item>
            <el-descriptions-item label="结论理由">{{ currentEval.recommendation_reason || '-' }}</el-descriptions-item>
            <el-descriptions-item label="薪资建议">{{ currentEval.salary_suggestion || '-' }}</el-descriptions-item>
            <el-descriptions-item label="适合岗位">
              <span v-if="currentEval.suitable_roles?.length">{{ currentEval.suitable_roles.join('、') }}</span>
              <span v-else>-</span>
            </el-descriptions-item>
          </el-descriptions>
          <div v-if="currentEval.interview_focus?.length" class="interview-focus">
            <h5>面试重点</h5>
            <ul>
              <li v-for="(item, i) in currentEval.interview_focus || []" :key="`focus-${i}`">{{ item }}</li>
            </ul>
          </div>
        </div>

        <div v-if="currentEval.core_strengths?.length || currentEval.improvement_items?.length || currentEval.risk_tips?.length" class="detail-section">
          <h4>综合评价</h4>
          <el-row :gutter="16">
            <el-col :span="8">
              <div class="report-box strengths">
                <h5><el-icon><CircleCheck /></el-icon> 核心优势</h5>
                <ul>
                  <li v-for="(item, i) in currentEval.core_strengths || []" :key="`strength-${i}`">{{ item }}</li>
                </ul>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="report-box gaps">
                <h5><el-icon><Warning /></el-icon> 待提升项</h5>
                <ul>
                  <li v-for="(item, i) in currentEval.improvement_items || []" :key="`improve-${i}`">{{ item }}</li>
                </ul>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="report-box risks">
                <h5><el-icon><Warning /></el-icon> 风险提示</h5>
                <ul>
                  <li v-for="(item, i) in currentEval.risk_tips || []" :key="`risk-tip-${i}`">{{ item }}</li>
                </ul>
              </div>
            </el-col>
          </el-row>
        </div>

        <div v-if="currentEval.interview_questions?.length" class="detail-section">
          <h4>建议面试题</h4>
          <el-collapse>
            <el-collapse-item
              v-for="(question, i) in currentEval.interview_questions"
              :key="`question-${i}`"
              :title="question.title || `问题 ${i + 1}`"
              :name="String(i)"
            >
              <p><strong>考察点：</strong>{{ question.focus || '-' }}</p>
              <p><strong>参考答案要点：</strong></p>
              <ul>
                <li v-for="(point, idx) in question.answer_points || []" :key="`point-${i}-${idx}`">{{ point }}</li>
              </ul>
            </el-collapse-item>
          </el-collapse>
        </div>

        <div v-if="currentEval.parsed_report" class="detail-section">
          <h4>原始 Coze 结构化结果</h4>
          <pre class="raw-report">{{ formatReport(currentEval.parsed_report) }}</pre>
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
    const res = await request.get('/evaluations', {
      params: {
        latest_only: true,
        page: pagination.page,
        page_size: pagination.pageSize,
        search: filters.search,
        status: filters.status,
        eval_type: filters.evalType,
        sort_by: sortParams.sortBy,
        sort_order: sortParams.sortOrder
      }
    })

    if (res.data?.code === 0) {
      evaluations.value = (res.data.data?.evaluations || []).map(normalizeEvaluation)
      pagination.total = res.data.data?.total || 0
    }
  } catch (e) {
    console.error('加载数据失败', e)
  } finally {
    loading.value = false
  }
}

// 解析技能字符串
const parseSkills = (skills: any) => {
  if (!skills) return []
  if (Array.isArray(skills)) return skills
  try {
    return JSON.parse(skills)
  } catch {
    return []
  }
}

// 解析维度数据
const parseDimensions = (dimensions: any) => {
  if (!dimensions) return []
  if (Array.isArray(dimensions)) return dimensions
  try {
    return JSON.parse(dimensions)
  } catch {
    return []
  }
}

// 加载统计
const loadStats = async () => {
  try {
    const res = await request.get('/evaluations/stats', {
      params: {
        latest_only: true
      }
    })
    if (res.data?.code === 0) {
      stats.value = res.data.data || {}
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
  currentEval.value = normalizeEvaluation(row)
  detailVisible.value = true
}

const openResume = (row: any) => {
  if (!row?.resume_id) {
    ElMessage.warning('该记录未关联简历')
    return
  }
  window.open(`/api/v1/resumes/${row.resume_id}/download`, '_blank')
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除这条最新评估结果吗？删除后该简历可能回退显示更早的一条评估记录。', '提示', {
      type: 'warning'
    })

    const res = await request.delete(`/evaluations/${row.id}`)
    if (res.data?.code === 0) {
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

const asRecord = (value: any): Record<string, any> => {
  if (value && typeof value === 'object' && !Array.isArray(value)) {
    return value
  }
  return {}
}

const getString = (value: any) => {
  return typeof value === 'string' ? value : ''
}

const getFirstString = (record: Record<string, any>, keys: string[]) => {
  for (const key of keys) {
    const value = getString(record[key])
    if (value) return value
  }
  return ''
}

const getNumber = (value: any) => {
  if (typeof value === 'number') return value
  if (typeof value === 'string') {
    const n = Number(value)
    return Number.isFinite(n) ? n : 0
  }
  return 0
}

const normalizeList = (value: any): string[] => {
  if (Array.isArray(value)) {
    return value.map((item) => String(item)).filter(Boolean)
  }
  if (typeof value === 'string' && value.trim()) {
    return value.split(/[、,，\n]/).map((item) => item.trim()).filter(Boolean)
  }
  return []
}

const normalizeInterviewQuestions = (value: any) => {
  if (!Array.isArray(value)) return []
  return value.map((item) => {
    const obj = asRecord(item)
    return {
      title: getString(obj['题目']),
      focus: getString(obj['考察点']),
      answer_points: normalizeList(obj['参考答案要点'])
    }
  }).filter((item) => item.title || item.focus || item.answer_points.length)
}

const buildDimensionDetails = (dimensionMap: Record<string, any>) => {
  return Object.entries(dimensionMap).map(([name, raw]) => {
    const detail = asRecord(raw)
    return {
      name,
      raw_score: getNumber(detail['得分']),
      max_score: getNumber(detail['满分']) || 10,
      description: getString(detail['说明'])
    }
  })
}

const buildDimensionProgress = (dimensionMap: Record<string, any>, fallback: any[]) => {
  const mapped = Object.entries(dimensionMap).map(([name, raw]) => {
    const detail = asRecord(raw)
    const score = getNumber(detail['得分'])
    const maxScore = getNumber(detail['满分']) || 10
    return {
      name,
      score: maxScore > 0 ? Math.round((score / maxScore) * 100) : 0
    }
  })
  return mapped.length ? mapped : fallback
}

const normalizeEvaluation = (e: any) => {
  const parsedReport = asRecord(e.parsed_report)
  const basicInfo = asRecord(parsedReport['基本信息'])
  const jdMatch = asRecord(parsedReport['JD匹配度'])
  const dimensionMap = asRecord(parsedReport['各维度得分'])
  const recommendation = asRecord(parsedReport['录用建议'])
  const summary = asRecord(parsedReport['综合评价'])

  const fallbackDimensions = parseDimensions(e.report_dimensions)

  return {
    ...e,
    parsed_name: e.parsed_name || getString(basicInfo['姓名']),
    parsed_phone: e.parsed_phone || getFirstString(basicInfo, ['手机', '手机号', '联系电话', '电话', 'phone']),
    parsed_email: e.parsed_email || getFirstString(basicInfo, ['邮箱', '邮箱地址', '电子邮箱', 'email']),
    parsed_education: e.parsed_education || getString(basicInfo['学历']),
    parsed_school: e.parsed_school || getString(basicInfo['学校']),
    parsed_experience: e.parsed_experience || getString(basicInfo['工作经验']),
    parsed_location: e.parsed_location || getString(basicInfo['城市']) || getString(basicInfo['地点']),
    grade: getString(basicInfo['评级']),
    parsed_skills: parseSkills(e.parsed_skills),
    report_dimensions: buildDimensionProgress(dimensionMap, fallbackDimensions),
    dimension_details: buildDimensionDetails(dimensionMap),
    report_strengths: parseSkills(e.report_strengths),
    report_gaps: parseSkills(e.report_gaps),
    jd_summary: getString(jdMatch['匹配总结']) || e.report_summary,
    jd_matched_skills: normalizeList(jdMatch['匹配的技能']),
    jd_missing_skills: normalizeList(jdMatch['缺失的技能']) || parseSkills(e.match_details),
    recommendation_conclusion: getString(recommendation['结论']) || e.report_recommendation,
    recommendation_reason: getString(recommendation['结论理由']),
    salary_suggestion: getString(recommendation['薪资建议']),
    suitable_roles: normalizeList(recommendation['适合岗位']),
    interview_focus: normalizeList(recommendation['面试重点']),
    interview_questions: normalizeInterviewQuestions(recommendation['面试题目']),
    core_strengths: normalizeList(summary['核心优势']) || parseSkills(e.report_strengths),
    improvement_items: normalizeList(summary['待提升项']) || parseSkills(e.report_gaps),
    risk_tips: normalizeList(summary['风险提示']),
  }
}

const formatReport = (report: any) => {
  try {
    return JSON.stringify(report, null, 2)
  } catch {
    return String(report || '')
  }
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

    .dimension-detail-list {
      display: grid;
      gap: 12px;
    }

    .dimension-detail-item {
      padding: 12px 14px;
      border-radius: 10px;
      background: var(--bg-secondary);

      .dimension-detail-head {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 8px;
      }

      .dimension-title {
        font-weight: 600;
      }

      .dimension-score {
        font-size: 12px;
        color: var(--text-secondary);
      }

      .dimension-description {
        margin: 0;
        white-space: pre-wrap;
        line-height: 1.7;
        color: var(--text-secondary);
      }
    }

    .report-box {
      background: var(--bg-secondary);
      border-radius: 8px;
      padding: 12px;

      &.strengths h5 { color: #67c23a; }
      &.gaps h5 { color: #e6a23c; }
      &.risks h5 { color: #ea580c; }

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

    .interview-focus {
      margin-top: 16px;

      h5 {
        margin: 0 0 8px;
      }
    }

    .raw-report {
      margin: 0;
      padding: 16px;
      border-radius: 10px;
      background: #0f172a;
      color: #e2e8f0;
      font-size: 12px;
      line-height: 1.6;
      white-space: pre-wrap;
      word-break: break-word;
      max-height: 420px;
      overflow: auto;
    }

    .risk-alert {
      margin-bottom: 8px;
    }
  }
}
</style>
