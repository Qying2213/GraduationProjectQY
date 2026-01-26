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
              <el-descriptions-item label="经验">{{ evaluateResult.parsed_experience || '-' }}</el-descriptions-item>
              <el-descriptions-item label="城市">{{ evaluateResult.parsed_location || '-' }}</el-descriptions-item>
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

// 开始评估
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
      evaluateResult.value = {
        // 解析信息
        parsed_name: data.candidate_name || data.parsed_name,
        parsed_phone: data.parsed_phone,
        parsed_email: data.parsed_email,
        parsed_education: data.parsed_education || getEducationFromGrade(data.education_score),
        parsed_experience: data.parsed_experience || `${data.experience_score || 0}年`,
        parsed_location: data.parsed_location,
        parsed_skills: data.matched_skills || data.parsed_skills || [],
        
        // 匹配分数
        match_score: data.total_score || data.match_score || 0,
        
        // 维度分析
        dimensions: [
          { name: 'JD匹配度', score: data.jd_match_score || 0 },
          { name: '工作经验', score: data.experience_score || 0 },
          { name: '学历背景', score: data.education_score || 0 },
          { name: '技术能力', score: data.tech_score || 0 },
          { name: '项目经验', score: data.project_score || 0 },
        ],
        
        // AI建议
        strengths: data.matched_skills || [],
        gaps: data.missing_skills || [],
        recommendation: data.recommendation || data.summary || ''
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
}
</style>
