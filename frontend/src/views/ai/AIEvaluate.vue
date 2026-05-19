<template>
  <div class="ai-evaluate-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-left">
        <h1>AI智能评估系统</h1>
        <p class="subtitle">基于大模型的简历智能解析与人岗匹配</p>
      </div>
      <div class="header-right">
        <el-tag type="success" v-if="aiConfigured">
          <el-icon><CircleCheck /></el-icon>
          AI服务已连接
        </el-tag>
        <el-tag type="warning" v-else>
          <el-icon><Warning /></el-icon>
          AI服务未配置
        </el-tag>
      </div>
    </div>

    <!-- 主要内容区：简历 + JD 并排 -->
    <el-row :gutter="20" class="input-section">
      <!-- 左侧：简历上传 -->
      <el-col :span="12">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><Document /></el-icon>
              <span>简历文件</span>
            </div>
          </template>
          
          <el-upload
            class="resume-upload"
            drag
            :action="uploadUrl"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
            :show-file-list="false"
            accept=".pdf,.doc,.docx"
          >
            <div v-if="uploadedFile" class="uploaded-file">
              <el-icon class="file-icon"><Document /></el-icon>
              <span class="file-name">{{ uploadedFile.name }}</span>
              <el-tag type="success" size="small">已上传</el-tag>
            </div>
            <div v-else>
              <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
              <div class="el-upload__text">
                拖拽简历到此处，或 <em>点击上传</em>
              </div>
              <div class="el-upload__tip">
                支持 PDF、DOC、DOCX 格式
              </div>
            </div>
          </el-upload>
        </el-card>
      </el-col>

      <!-- 右侧：JD输入 -->
      <el-col :span="12">
        <el-card class="input-card">
          <template #header>
            <div class="card-header">
              <el-icon><Suitcase /></el-icon>
              <span>职位描述 (JD)</span>
              <el-select 
                v-model="selectedJobId" 
                placeholder="或选择已有职位"
                clearable
                size="small"
                class="job-select"
                @change="onJobSelect"
              >
                <el-option
                  v-for="job in jobList"
                  :key="job.id"
                  :label="job.title"
                  :value="job.id"
                />
              </el-select>
            </div>
          </template>
          
          <el-input
            v-model="jdText"
            type="textarea"
            :rows="8"
            placeholder="请输入职位描述（JD），包括：
• 岗位职责
• 任职要求
• 技能要求
• 学历/经验要求等

AI将根据JD对简历进行针对性评估..."
          />
        </el-card>
      </el-col>
    </el-row>

    <!-- 开始评估按钮 -->
    <div class="action-section">
      <el-button 
        type="primary" 
        size="large"
        @click="startEvaluate" 
        :loading="evaluating"
        :disabled="!canEvaluate"
      >
        <el-icon><MagicStick /></el-icon>
        开始AI智能评估
      </el-button>
      <p class="action-tip" v-if="!canEvaluate">
        请上传简历并填写职位描述后开始评估
      </p>
    </div>

    <!-- 评估结果区域 -->
    <el-card v-if="evaluateResult" class="result-card">
      <template #header>
        <div class="result-header">
          <span>评估结果</span>
          <el-tag :type="getMatchLevelType(evaluateResult.match_score)" size="large">
            匹配度: {{ evaluateResult.match_score?.toFixed(1) }}分
          </el-tag>
        </div>
      </template>

      <el-row :gutter="24">
        <!-- 左侧：解析信息 -->
        <el-col :span="8">
          <div class="result-section">
            <h4><el-icon><User /></el-icon> 候选人信息</h4>
            <el-descriptions :column="1" border size="small">
              <el-descriptions-item label="姓名">{{ evaluateResult.parsed_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="手机">{{ evaluateResult.parsed_phone || '-' }}</el-descriptions-item>
              <el-descriptions-item label="邮箱">{{ evaluateResult.parsed_email || '-' }}</el-descriptions-item>
              <el-descriptions-item label="学历">{{ evaluateResult.parsed_education || '-' }}</el-descriptions-item>
              <el-descriptions-item label="学校">{{ evaluateResult.parsed_school || '-' }}</el-descriptions-item>
              <el-descriptions-item label="经验">{{ evaluateResult.parsed_experience || '-' }}</el-descriptions-item>
              <el-descriptions-item label="城市">{{ evaluateResult.parsed_location || '-' }}</el-descriptions-item>
              <el-descriptions-item label="评级">{{ evaluateResult.grade || '-' }}</el-descriptions-item>
            </el-descriptions>
            
            <div class="skills-section" v-if="evaluateResult.parsed_skills?.length">
              <h5>技能标签</h5>
              <div class="skills-wrapper">
                <el-tag 
                  v-for="skill in evaluateResult.parsed_skills" 
                  :key="skill" 
                  size="small"
                  class="skill-tag"
                >
                  {{ skill }}
                </el-tag>
              </div>
            </div>
          </div>
        </el-col>

        <!-- 中间：匹配分析 -->
        <el-col :span="8">
          <div class="result-section">
            <h4><el-icon><DataAnalysis /></el-icon> 匹配分析</h4>
            
            <div class="score-display">
              <el-progress 
                type="dashboard" 
                :percentage="Math.round(evaluateResult.match_score || 0)" 
                :color="getScoreColor(evaluateResult.match_score)"
                :width="120"
              >
                <template #default="{ percentage }">
                  <span class="score-value">{{ percentage }}</span>
                  <span class="score-label">匹配度</span>
                </template>
              </el-progress>
            </div>

            <div v-if="evaluateResult.dimensions" class="dimensions">
              <div v-for="dim in evaluateResult.dimensions" :key="dim.name" class="dimension-item">
                <span class="dim-name">{{ dim.name }}</span>
                <el-progress 
                  :percentage="Math.round(dim.score)" 
                  :stroke-width="8"
                  :color="getScoreColor(dim.score)"
                />
              </div>
            </div>
          </div>
        </el-col>

        <!-- 右侧：AI建议 -->
        <el-col :span="8">
          <div class="result-section">
            <h4><el-icon><ChatLineSquare /></el-icon> AI评估建议</h4>
            
            <div v-if="evaluateResult.strengths?.length" class="advice-box strengths">
              <h5><el-icon><CircleCheck /></el-icon> 优势</h5>
              <ul>
                <li v-for="(item, i) in evaluateResult.strengths" :key="i">{{ item }}</li>
              </ul>
            </div>

            <div v-if="evaluateResult.gaps?.length" class="advice-box gaps">
              <h5><el-icon><Warning /></el-icon> 待提升</h5>
              <ul>
                <li v-for="(item, i) in evaluateResult.gaps" :key="i">{{ item }}</li>
              </ul>
            </div>

            <div v-if="evaluateResult.recommendation" class="recommendation-box">
              <h5>综合建议</h5>
              <p>{{ evaluateResult.recommendation }}</p>
            </div>
          </div>
        </el-col>
      </el-row>

      <div v-if="evaluateResult.jd_summary || evaluateResult.jd_matched_skills?.length || evaluateResult.jd_missing_skills?.length" class="detail-block">
        <h4>JD 匹配详细分析</h4>
        <p v-if="evaluateResult.jd_summary" class="detail-summary">{{ evaluateResult.jd_summary }}</p>

        <el-row :gutter="16">
          <el-col :span="12">
            <div class="report-box strengths">
              <h5><el-icon><CircleCheck /></el-icon> 匹配到的能力</h5>
              <ul>
                <li v-for="(item, i) in evaluateResult.jd_matched_skills || []" :key="`matched-${i}`">{{ item }}</li>
              </ul>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="report-box gaps">
              <h5><el-icon><Warning /></el-icon> 缺失或待确认</h5>
              <ul>
                <li v-for="(item, i) in evaluateResult.jd_missing_skills || []" :key="`missing-${i}`">{{ item }}</li>
              </ul>
            </div>
          </el-col>
        </el-row>
      </div>

      <div v-if="evaluateResult.dimension_details?.length" class="detail-block">
        <h4>各维度详细说明</h4>
        <div class="dimension-detail-list">
          <div v-for="item in evaluateResult.dimension_details" :key="item.name" class="dimension-detail-item">
            <div class="dimension-detail-head">
              <span class="dimension-title">{{ item.name }}</span>
              <span class="dimension-score">{{ item.raw_score }}/{{ item.max_score }}</span>
            </div>
            <p class="dimension-description">{{ item.description || '暂无说明' }}</p>
          </div>
        </div>
      </div>

      <div v-if="evaluateResult.recommendation_conclusion || evaluateResult.recommendation_reason || evaluateResult.salary_suggestion || evaluateResult.suitable_roles?.length || evaluateResult.interview_focus?.length" class="detail-block">
        <h4>录用建议</h4>
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="结论">{{ evaluateResult.recommendation_conclusion || '-' }}</el-descriptions-item>
          <el-descriptions-item label="结论理由">{{ evaluateResult.recommendation_reason || '-' }}</el-descriptions-item>
          <el-descriptions-item label="薪资建议">{{ evaluateResult.salary_suggestion || '-' }}</el-descriptions-item>
          <el-descriptions-item label="适合岗位">
            <span v-if="evaluateResult.suitable_roles?.length">{{ evaluateResult.suitable_roles.join('、') }}</span>
            <span v-else>-</span>
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="evaluateResult.interview_focus?.length" class="interview-focus">
          <h5>面试重点</h5>
          <ul>
            <li v-for="(item, i) in evaluateResult.interview_focus" :key="`focus-${i}`">{{ item }}</li>
          </ul>
        </div>
      </div>

      <div v-if="evaluateResult.core_strengths?.length || evaluateResult.improvement_items?.length || evaluateResult.risk_tips?.length" class="detail-block">
        <h4>综合评价</h4>
        <el-row :gutter="16">
          <el-col :span="8">
            <div class="report-box strengths">
              <h5><el-icon><CircleCheck /></el-icon> 核心优势</h5>
              <ul>
                <li v-for="(item, i) in evaluateResult.core_strengths || []" :key="`strength-${i}`">{{ item }}</li>
              </ul>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="report-box gaps">
              <h5><el-icon><Warning /></el-icon> 待提升项</h5>
              <ul>
                <li v-for="(item, i) in evaluateResult.improvement_items || []" :key="`improve-${i}`">{{ item }}</li>
              </ul>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="report-box risks">
              <h5><el-icon><Warning /></el-icon> 风险提示</h5>
              <ul>
                <li v-for="(item, i) in evaluateResult.risk_tips || []" :key="`risk-${i}`">{{ item }}</li>
              </ul>
            </div>
          </el-col>
        </el-row>
      </div>

      <div v-if="evaluateResult.interview_questions?.length" class="detail-block">
        <h4>建议面试题</h4>
        <el-collapse>
          <el-collapse-item
            v-for="(question, i) in evaluateResult.interview_questions"
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

      <div v-if="evaluateResult.parsed_report" class="detail-block">
        <h4>原始 Coze 结构化结果</h4>
        <pre class="raw-report">{{ formatReport(evaluateResult.parsed_report) }}</pre>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Document, Suitcase, UploadFilled, MagicStick, User,
  Warning, CircleCheck, ChatLineSquare, DataAnalysis
} from '@element-plus/icons-vue'
import request from '@/utils/request'

// AIEvaluate 是毕业设计核心演示页。
// 页面负责把“上传简历 + 选择/输入 JD”组合成一次 AI 评估请求，并展示后端返回的解析、
// 匹配分、风险项和推荐理由。
const aiConfigured = ref(false)

// 简历相关
const uploadedResumeId = ref<number | null>(null)
const uploadedFile = ref<{ name: string } | null>(null)

// JD相关
const jdText = ref('')
const selectedJobId = ref<number | null>(null)
const jobList = ref<any[]>([])

// 评估状态
const evaluating = ref(false)
const evaluateResult = ref<any>(null)

// 上传配置
const uploadUrl = '/api/v1/resumes/upload'
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 是否可以开始评估
const canEvaluate = computed(() => {
  return uploadedResumeId.value && jdText.value.trim().length > 0
})

// 检查AI配置
const checkAIConfig = async () => {
  try {
    const res = await request.get('/ai/config')
    aiConfigured.value = res.data?.configured || false
  } catch {
    aiConfigured.value = false
  }
}

// 加载职位列表
const loadJobs = async () => {
  try {
    const res = await request.get('/jobs', { params: { page_size: 100, status: 'open' } })
    console.log('[AIEvaluate] 职位列表响应:', res.data)
    // 后端返回格式: { code: 0, data: { jobs: [...], total: n } }
    if (res.data?.code === 0 && res.data?.data?.jobs) {
      jobList.value = res.data.data.jobs
    } else if (res.data?.data) {
      // 兼容其他格式
      jobList.value = Array.isArray(res.data.data) ? res.data.data : []
    } else {
      jobList.value = []
    }
    console.log('[AIEvaluate] 加载到职位数量:', jobList.value.length)
  } catch (e) {
    console.error('[AIEvaluate] 加载职位列表失败', e)
    jobList.value = []
  }
}

// 选择职位时填充JD
const onJobSelect = (jobId: number) => {
  if (!jobId) return
  const job = jobList.value.find(j => j.id === jobId)
  if (job) {
    // 组合职位信息生成JD
    const parts = []
    if (job.title) parts.push(`职位名称：${job.title}`)
    if (job.department) parts.push(`所属部门：${job.department}`)
    if (job.location) parts.push(`工作地点：${job.location}`)
    if (job.salary) parts.push(`薪资范围：${job.salary}`)
    if (job.description) parts.push(`\n岗位职责：\n${job.description}`)
    if (job.requirements) parts.push(`\n任职要求：\n${job.requirements}`)
    if (job.skills) parts.push(`\n技能要求：${job.skills}`)
    
    jdText.value = parts.join('\n')
  }
}

// 上传前检查
const beforeUpload = (file: File) => {
  const isValidType = ['application/pdf', 'application/msword', 
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document'].includes(file.type)
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isValidType) {
    ElMessage.error('只支持 PDF、DOC、DOCX 格式!')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB!')
    return false
  }
  return true
}

// 上传成功
const handleUploadSuccess = (response: any, file: any) => {
  if (response.code === 0) {
    uploadedResumeId.value = response.data.id
    uploadedFile.value = { name: file.name }
    ElMessage.success('简历上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

// 上传失败
const handleUploadError = () => {
  ElMessage.error('上传失败，请重试')
}

// 开始评估。
// 前端只负责提交 resume_id、JD 和可选 job_id；OCR、Embedding、RAG 和 Coze 评分都在后端完成。
const startEvaluate = async () => {
  if (!uploadedResumeId.value || !jdText.value.trim()) {
    ElMessage.warning('请上传简历并填写职位描述')
    return
  }

  evaluating.value = true
  evaluateResult.value = null

  try {
    // 调用AI评估接口
    const response = await request.post('/ai/evaluate', {
      resume_id: uploadedResumeId.value,
      jd_text: jdText.value,
      job_id: selectedJobId.value
    })
    
    const res = response.data || response
    
    if (res.code === 0) {
      // 处理评估结果
      const data = res.data
      const parsedReport = normalizeParsedReport(data.parsed_report, data.raw_result)
      const basicInfo = asRecord(parsedReport['基本信息'])
      const jdMatch = asRecord(parsedReport['JD匹配度'])
      const dimensionMap = asRecord(parsedReport['各维度得分'])
      const recommendation = asRecord(parsedReport['录用建议'])
      const summary = asRecord(parsedReport['综合评价'])

      evaluateResult.value = {
        // 解析信息
        parsed_name: data.candidate_name || data.parsed_name || getString(basicInfo['姓名']),
        parsed_phone: data.parsed_phone || getFirstString(basicInfo, ['手机', '手机号', '联系电话', '电话', 'phone']),
        parsed_email: data.parsed_email || getFirstString(basicInfo, ['邮箱', '邮箱地址', '电子邮箱', 'email']),
        parsed_education: data.parsed_education || getString(basicInfo['学历']) || getEducationFromGrade(data.education_score),
        parsed_school: getString(basicInfo['学校']),
        parsed_experience: data.parsed_experience || getString(basicInfo['工作经验']) || `${data.experience_score || 0}年`,
        parsed_location: data.parsed_location || getString(basicInfo['城市']) || getString(basicInfo['地点']),
        parsed_skills: data.matched_skills || data.parsed_skills || toStringArray(jdMatch['匹配的技能']),
        grade: data.grade || getString(basicInfo['评级']),
        
        // 匹配分数
        match_score: data.total_score || data.match_score || 0,
        
        // 维度分析
        dimensions: buildDimensionProgress(data, dimensionMap),
        dimension_details: buildDimensionDetails(dimensionMap),
        
        // AI建议
        strengths: data.matched_skills || [],
        gaps: data.missing_skills || [],
        recommendation: data.recommendation || data.summary || '',

        jd_summary: getString(jdMatch['匹配总结']) || data.summary || '',
        jd_matched_skills: toStringArray(jdMatch['匹配的技能']) || data.matched_skills || [],
        jd_missing_skills: toStringArray(jdMatch['缺失的技能']) || data.missing_skills || [],

        recommendation_conclusion: getString(recommendation['结论']) || data.recommendation || '',
        recommendation_reason: getString(recommendation['结论理由']),
        salary_suggestion: getString(recommendation['薪资建议']),
        suitable_roles: normalizeList(recommendation['适合岗位']),
        interview_focus: normalizeList(recommendation['面试重点']),
        interview_questions: normalizeInterviewQuestions(recommendation['面试题目']),

        core_strengths: normalizeList(summary['核心优势']),
        improvement_items: normalizeList(summary['待提升项']),
        risk_tips: normalizeList(summary['风险提示']),
        parsed_report: parsedReport
      }
      
      ElMessage.success('评估完成')
    } else {
      ElMessage.error(res.message || '评估失败')
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || e.message || '评估失败')
  } finally {
    evaluating.value = false
  }
}

// 辅助函数
const getEducationFromGrade = (score: number) => {
  if (score >= 90) return '博士'
  if (score >= 80) return '硕士'
  if (score >= 70) return '本科'
  if (score >= 60) return '大专'
  return '-'
}

const getScoreColor = (score: number) => {
  if (score >= 80) return '#67c23a'
  if (score >= 60) return '#e6a23c'
  return '#f56c6c'
}

const getMatchLevelType = (score: number) => {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'danger'
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

const toStringArray = (value: any): string[] => {
  if (!Array.isArray(value)) return []
  return value.map((item) => String(item)).filter(Boolean)
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

const normalizeParsedReport = (parsedReport: any, rawResult: any) => {
  if (parsedReport && typeof parsedReport === 'object') {
    return parsedReport
  }

  const raw = asRecord(rawResult)
  const data = raw.data
  if (typeof data === 'string') {
    try {
      const parsed = JSON.parse(data)
      if (parsed?.result) {
        return JSON.parse(parsed.result)
      }
      if (parsed?.output) {
        return JSON.parse(parsed.output)
      }
    } catch {
      return {}
    }
  }

  const dataObj = asRecord(data)
  if (typeof dataObj.result === 'string') {
    try {
      return JSON.parse(dataObj.result)
    } catch {
      return {}
    }
  }
  if (typeof dataObj.output === 'string') {
    try {
      return JSON.parse(dataObj.output)
    } catch {
      return {}
    }
  }
  return {}
}

const buildDimensionProgress = (data: any, dimensionMap: Record<string, any>) => {
  const defaultDimensions = [
    { name: 'JD匹配度', score: data.jd_match_score || 0 },
    { name: '年龄适配', score: (data.age_score || 0) * 10 },
    { name: '工作经验', score: (data.experience_score || 0) * 10 },
    { name: '学历背景', score: (data.education_score || 0) * 10 },
    { name: '公司背景', score: (data.company_score || 0) * 10 },
    { name: '技术能力', score: (data.tech_score || 0) * 10 },
    { name: '项目经验', score: (data.project_score || 0) * 10 },
  ]

  const mapped = Object.entries(dimensionMap).map(([name, raw]) => {
    const detail = asRecord(raw)
    const score = getNumber(detail['得分'])
    const maxScore = getNumber(detail['满分']) || 10
    const percent = maxScore > 0 ? Math.round((score / maxScore) * 100) : 0
    return {
      name,
      score: percent,
    }
  })

  return mapped.length ? mapped : defaultDimensions
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

const formatReport = (report: any) => {
  try {
    return JSON.stringify(report, null, 2)
  } catch {
    return String(report || '')
  }
}

onMounted(() => {
  checkAIConfig()
  loadJobs()
})
</script>

<style scoped lang="scss">
.ai-evaluate-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;

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

  .input-section {
    margin-bottom: 24px;
  }

  .input-card {
    height: 100%;
    
    .card-header {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .el-icon {
        font-size: 18px;
        color: var(--primary-color);
      }
      
      .job-select {
        margin-left: auto;
        width: 180px;
      }
    }

    .resume-upload {
      :deep(.el-upload-dragger) {
        padding: 40px 20px;
        min-height: 180px;
        display: flex;
        flex-direction: column;
        justify-content: center;
      }
      
      .uploaded-file {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 12px;
        
        .file-icon {
          font-size: 48px;
          color: var(--primary-color);
        }
        
        .file-name {
          font-size: 14px;
          color: var(--text-primary);
          max-width: 200px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }

    :deep(.el-textarea__inner) {
      min-height: 180px !important;
    }
  }

  .action-section {
    text-align: center;
    margin-bottom: 24px;
    
    .el-button {
      padding: 12px 48px;
      font-size: 16px;
      
      .el-icon {
        margin-right: 8px;
      }
    }
    
    .action-tip {
      margin-top: 12px;
      font-size: 13px;
      color: var(--text-secondary);
    }
  }

  .result-card {
    .result-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      span {
        font-size: 16px;
        font-weight: 600;
      }
    }

    .result-section {
      h4 {
        display: flex;
        align-items: center;
        gap: 8px;
        margin: 0 0 16px 0;
        font-size: 15px;
        color: var(--text-primary);
        padding-bottom: 12px;
        border-bottom: 1px solid var(--border-color);
        
        .el-icon {
          color: var(--primary-color);
        }
      }
      
      h5 {
        display: flex;
        align-items: center;
        gap: 6px;
        margin: 12px 0 8px 0;
        font-size: 13px;
        color: var(--text-secondary);
      }
    }

    .skills-section {
      margin-top: 16px;
      
      .skills-wrapper {
        display: flex;
        flex-wrap: wrap;
        gap: 6px;
      }
      
      .skill-tag {
        margin: 0;
      }
    }

    .score-display {
      text-align: center;
      margin-bottom: 20px;

      .score-value {
        display: block;
        font-size: 28px;
        font-weight: bold;
        color: var(--text-primary);
      }

      .score-label {
        font-size: 12px;
        color: var(--text-secondary);
      }
    }

    .dimensions {
      .dimension-item {
        margin-bottom: 10px;

        .dim-name {
          display: block;
          margin-bottom: 4px;
          font-size: 12px;
          color: var(--text-secondary);
        }
      }
    }

    .advice-box {
      background: var(--bg-secondary);
      border-radius: 8px;
      padding: 12px;
      margin-bottom: 12px;

      &.strengths h5 { color: #67c23a; }
      &.gaps h5 { color: #e6a23c; }

      ul {
        margin: 0;
        padding-left: 18px;
        font-size: 13px;

        li {
          margin-bottom: 4px;
          color: var(--text-primary);
        }
      }
    }

    .recommendation-box {
      background: linear-gradient(135deg, #667eea10 0%, #764ba210 100%);
      border-radius: 8px;
      padding: 12px;

      h5 { 
        color: #667eea; 
        margin-top: 0;
      }

      p {
        margin: 0;
        font-size: 13px;
        color: var(--text-primary);
        line-height: 1.6;
      }
    }
  }

  .detail-block {
    margin-top: 24px;
    padding-top: 20px;
    border-top: 1px solid var(--el-border-color-lighter);

    h4 {
      margin: 0 0 16px;
      font-size: 16px;
      color: var(--text-primary);
    }

    .detail-summary {
      margin: 0 0 16px;
      padding: 14px 16px;
      border-radius: 12px;
      background: var(--el-fill-color-light);
      color: var(--text-primary);
      line-height: 1.7;
    }
  }

  .dimension-detail-list {
    display: grid;
    gap: 12px;
  }

  .dimension-detail-item {
    padding: 14px 16px;
    border-radius: 12px;
    background: var(--el-fill-color-light);

    .dimension-detail-head {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;
    }

    .dimension-title {
      font-weight: 600;
      color: var(--text-primary);
    }

    .dimension-score {
      font-size: 13px;
      color: var(--text-secondary);
    }

    .dimension-description {
      margin: 0;
      line-height: 1.7;
      color: var(--text-secondary);
      white-space: pre-wrap;
    }
  }

  .interview-focus {
    margin-top: 16px;

    h5 {
      margin: 0 0 10px;
    }
  }

  .raw-report {
    margin: 0;
    padding: 16px;
    border-radius: 12px;
    background: #0f172a;
    color: #e2e8f0;
    font-size: 12px;
    line-height: 1.6;
    white-space: pre-wrap;
    word-break: break-word;
    max-height: 420px;
    overflow: auto;
  }

  .report-box.risks {
    background: #fff7ed;
    h5 { color: #ea580c; }
  }
}
</style>
