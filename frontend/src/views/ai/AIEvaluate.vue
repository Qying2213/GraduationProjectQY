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

    <!-- 功能卡片 -->
    <el-row :gutter="20" class="feature-cards">
      <el-col :span="8">
        <el-card class="feature-card" shadow="hover" @click="activeTab = 'parse'">
          <div class="card-icon parse">
            <el-icon><Document /></el-icon>
          </div>
          <h3>简历智能解析</h3>
          <p>OCR + 大模型提取简历关键信息</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="feature-card" shadow="hover" @click="activeTab = 'match'">
          <div class="card-icon match">
            <el-icon><Connection /></el-icon>
          </div>
          <h3>人岗智能匹配</h3>
          <p>语义向量化 + 多维度匹配算法</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="feature-card" shadow="hover" @click="activeTab = 'report'">
          <div class="card-icon report">
            <el-icon><DataAnalysis /></el-icon>
          </div>
          <h3>归因报告生成</h3>
          <p>可解释的推荐理由与建议</p>
        </el-card>
      </el-col>
    </el-row>

    <!-- 主要内容区 -->
    <el-card class="main-content">
      <el-tabs v-model="activeTab">
        <!-- 简历解析 -->
        <el-tab-pane label="简历智能解析" name="parse">
          <div class="tab-content">
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="upload-section">
                  <h4>上传简历文件</h4>
                  <el-upload
                    class="resume-upload"
                    drag
                    :action="uploadUrl"
                    :headers="uploadHeaders"
                    :on-success="handleUploadSuccess"
                    :on-error="handleUploadError"
                    :before-upload="beforeUpload"
                    accept=".pdf,.doc,.docx"
                  >
                    <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
                    <div class="el-upload__text">
                      拖拽文件到此处，或 <em>点击上传</em>
                    </div>
                    <template #tip>
                      <div class="el-upload__tip">
                        支持 PDF、DOC、DOCX 格式，文件大小不超过 10MB
                      </div>
                    </template>
                  </el-upload>

                  <div class="jd-input" v-if="uploadedResumeId">
                    <h4>职位描述（可选）</h4>
                    <el-input
                      v-model="jdText"
                      type="textarea"
                      :rows="4"
                      placeholder="输入职位描述，AI将根据JD进行针对性解析..."
                    />
                    <el-button 
                      type="primary" 
                      @click="startAIParse" 
                      :loading="parsing"
                      class="parse-btn"
                    >
                      <el-icon><MagicStick /></el-icon>
                      开始AI解析
                    </el-button>
                  </div>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="result-section">
                  <h4>解析结果</h4>
                  <div v-if="parseResult" class="parse-result">
                    <el-descriptions :column="1" border>
                      <el-descriptions-item label="姓名">
                        {{ parseResult.name || '-' }}
                      </el-descriptions-item>
                      <el-descriptions-item label="手机">
                        {{ parseResult.phone || '-' }}
                      </el-descriptions-item>
                      <el-descriptions-item label="邮箱">
                        {{ parseResult.email || '-' }}
                      </el-descriptions-item>
                      <el-descriptions-item label="学历">
                        {{ parseResult.education || '-' }}
                      </el-descriptions-item>
                      <el-descriptions-item label="工作经验">
                        {{ parseResult.experience || '-' }}
                      </el-descriptions-item>
                      <el-descriptions-item label="所在城市">
                        {{ parseResult.location || '-' }}
                      </el-descriptions-item>
                      <el-descriptions-item label="技能标签">
                        <el-tag 
                          v-for="skill in parseResult.skills" 
                          :key="skill" 
                          size="small"
                          class="skill-tag"
                        >
                          {{ skill }}
                        </el-tag>
                        <span v-if="!parseResult.skills?.length">-</span>
                      </el-descriptions-item>
                    </el-descriptions>

                    <!-- 风控提示 -->
                    <div v-if="parseResult.risk_items?.length" class="risk-section">
                      <h5><el-icon><Warning /></el-icon> 风控提示</h5>
                      <el-alert
                        v-for="(risk, index) in parseResult.risk_items"
                        :key="index"
                        :title="risk.message"
                        :type="risk.level === 'high' ? 'error' : risk.level === 'warning' ? 'warning' : 'info'"
                        :closable="false"
                        show-icon
                        class="risk-alert"
                      />
                    </div>
                  </div>
                  <el-empty v-else description="上传简历后查看解析结果" />
                </div>
              </el-col>
            </el-row>
          </div>
        </el-tab-pane>

        <!-- 人岗匹配 -->
        <el-tab-pane label="人岗智能匹配" name="match">
          <div class="tab-content">
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="match-form">
                  <h4>选择人才</h4>
                  <el-select 
                    v-model="selectedTalentId" 
                    placeholder="选择人才"
                    filterable
                    class="full-width"
                  >
                    <el-option
                      v-for="talent in talentList"
                      :key="talent.id"
                      :label="talent.name"
                      :value="talent.id"
                    />
                  </el-select>

                  <h4>选择职位</h4>
                  <el-select 
                    v-model="selectedJobId" 
                    placeholder="选择职位"
                    filterable
                    class="full-width"
                  >
                    <el-option
                      v-for="job in jobList"
                      :key="job.id"
                      :label="job.title"
                      :value="job.id"
                    />
                  </el-select>

                  <el-button 
                    type="primary" 
                    @click="startMatch" 
                    :loading="matching"
                    :disabled="!selectedTalentId || !selectedJobId"
                    class="match-btn"
                  >
                    <el-icon><Connection /></el-icon>
                    开始匹配分析
                  </el-button>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="match-result" v-if="matchResult">
                  <div class="score-display">
                    <el-progress 
                      type="dashboard" 
                      :percentage="matchResult.match_score" 
                      :color="getScoreColor(matchResult.match_score)"
                      :width="150"
                    >
                      <template #default="{ percentage }">
                        <span class="score-value">{{ percentage }}</span>
                        <span class="score-label">匹配度</span>
                      </template>
                    </el-progress>
                  </div>

                  <div class="match-details">
                    <h5>匹配维度分析</h5>
                    <div v-for="dim in matchResult.dimensions" :key="dim.name" class="dimension-item">
                      <span class="dim-name">{{ dim.name }}</span>
                      <el-progress 
                        :percentage="dim.score" 
                        :stroke-width="10"
                        :color="getScoreColor(dim.score)"
                      />
                    </div>
                  </div>
                </div>
                <el-empty v-else description="选择人才和职位后查看匹配结果" />
              </el-col>
            </el-row>
          </div>
        </el-tab-pane>

        <!-- 归因报告 -->
        <el-tab-pane label="归因报告" name="report">
          <div class="tab-content">
            <div v-if="attributionReport" class="attribution-report">
              <div class="report-header">
                <h4>{{ attributionReport.summary }}</h4>
                <el-tag :type="getMatchLevelType(attributionReport.match_score)" size="large">
                  综合得分: {{ attributionReport.match_score }}分
                </el-tag>
              </div>

              <el-row :gutter="20">
                <el-col :span="12">
                  <div class="report-section strengths">
                    <h5><el-icon><CircleCheck /></el-icon> 优势</h5>
                    <ul>
                      <li v-for="(item, index) in attributionReport.strengths" :key="index">
                        {{ item }}
                      </li>
                    </ul>
                  </div>
                </el-col>
                <el-col :span="12">
                  <div class="report-section gaps">
                    <h5><el-icon><Warning /></el-icon> 待提升</h5>
                    <ul>
                      <li v-for="(item, index) in attributionReport.gaps" :key="index">
                        {{ item }}
                      </li>
                    </ul>
                  </div>
                </el-col>
              </el-row>

              <div class="recommendation">
                <h5><el-icon><ChatLineSquare /></el-icon> AI建议</h5>
                <p>{{ attributionReport.recommendation }}</p>
              </div>
            </div>
            <el-empty v-else description="完成人岗匹配后生成归因报告" />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Document, Connection, DataAnalysis, UploadFilled, MagicStick,
  Warning, CircleCheck, ChatLineSquare
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const activeTab = ref('parse')
const aiConfigured = ref(false)

// 简历解析相关
const uploadedResumeId = ref<number | null>(null)
const jdText = ref('')
const parsing = ref(false)
const parseResult = ref<any>(null)

// 人岗匹配相关
const selectedTalentId = ref<number | null>(null)
const selectedJobId = ref<number | null>(null)
const matching = ref(false)
const matchResult = ref<any>(null)
const attributionReport = ref<any>(null)

// 数据列表
const talentList = ref<any[]>([])
const jobList = ref<any[]>([])

// 上传配置
const uploadUrl = '/api/v1/resumes/upload'
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 检查AI配置
const checkAIConfig = async () => {
  try {
    const res = await request.get('/api/v1/ai/config')
    aiConfigured.value = res.data?.configured || false
  } catch {
    aiConfigured.value = false
  }
}

// 加载人才列表
const loadTalents = async () => {
  try {
    const res = await request.get('/api/v1/talents', { params: { page_size: 100 } })
    talentList.value = res.data?.talents || []
  } catch (e) {
    console.error('加载人才列表失败', e)
  }
}

// 加载职位列表
const loadJobs = async () => {
  try {
    const res = await request.get('/api/v1/jobs', { params: { page_size: 100 } })
    jobList.value = res.data?.jobs || []
  } catch (e) {
    console.error('加载职位列表失败', e)
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
const handleUploadSuccess = (response: any) => {
  if (response.code === 0) {
    uploadedResumeId.value = response.data.id
    ElMessage.success('简历上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

// 上传失败
const handleUploadError = () => {
  ElMessage.error('上传失败，请重试')
}

// 开始AI解析
const startAIParse = async () => {
  if (!uploadedResumeId.value) return

  parsing.value = true
  try {
    const res = await request.post('/api/v1/ai/parse', {
      resume_id: uploadedResumeId.value,
      jd_text: jdText.value
    })
    
    if (res.code === 0) {
      parseResult.value = res.data?.parsed_result || res.data
      ElMessage.success('解析完成')
    } else {
      ElMessage.error(res.message || '解析失败')
    }
  } catch (e: any) {
    ElMessage.error(e.message || '解析失败')
  } finally {
    parsing.value = false
  }
}

// 开始匹配
const startMatch = async () => {
  if (!selectedTalentId.value || !selectedJobId.value) return

  matching.value = true
  try {
    const res = await request.post('/api/v1/recommendations/attribution-report', {
      talent_id: selectedTalentId.value,
      job_id: selectedJobId.value
    })
    
    if (res.code === 0) {
      matchResult.value = {
        match_score: res.data.match_score,
        dimensions: res.data.dimensions
      }
      attributionReport.value = res.data
      activeTab.value = 'report'
      ElMessage.success('匹配分析完成')
    } else {
      ElMessage.error(res.message || '匹配失败')
    }
  } catch (e: any) {
    ElMessage.error(e.message || '匹配失败')
  } finally {
    matching.value = false
  }
}

// 获取分数颜色
const getScoreColor = (score: number) => {
  if (score >= 80) return '#67c23a'
  if (score >= 60) return '#e6a23c'
  return '#f56c6c'
}

// 获取匹配等级类型
const getMatchLevelType = (score: number) => {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'danger'
}

onMounted(() => {
  checkAIConfig()
  loadTalents()
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

  .feature-cards {
    margin-bottom: 24px;

    .feature-card {
      cursor: pointer;
      text-align: center;
      transition: all 0.3s;

      &:hover {
        transform: translateY(-4px);
      }

      .card-icon {
        width: 60px;
        height: 60px;
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 16px;

        .el-icon {
          font-size: 28px;
          color: white;
        }

        &.parse {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }

        &.match {
          background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
        }

        &.report {
          background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
        }
      }

      h3 {
        margin: 0 0 8px 0;
        font-size: 16px;
        color: var(--text-primary);
      }

      p {
        margin: 0;
        font-size: 13px;
        color: var(--text-secondary);
      }
    }
  }

  .main-content {
    .tab-content {
      padding: 20px 0;
    }

    h4 {
      margin: 0 0 16px 0;
      font-size: 15px;
      color: var(--text-primary);
    }

    h5 {
      margin: 16px 0 12px 0;
      font-size: 14px;
      color: var(--text-primary);
      display: flex;
      align-items: center;
      gap: 6px;
    }

    .full-width {
      width: 100%;
      margin-bottom: 16px;
    }

    .parse-btn, .match-btn {
      margin-top: 16px;
      width: 100%;
    }

    .skill-tag {
      margin: 2px 4px 2px 0;
    }

    .risk-section {
      margin-top: 20px;
      padding-top: 16px;
      border-top: 1px solid var(--border-color);

      .risk-alert {
        margin-bottom: 8px;
      }
    }

    .score-display {
      text-align: center;
      margin-bottom: 24px;

      .score-value {
        display: block;
        font-size: 32px;
        font-weight: bold;
        color: var(--text-primary);
      }

      .score-label {
        font-size: 14px;
        color: var(--text-secondary);
      }
    }

    .match-details {
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

    .attribution-report {
      .report-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 24px;
        padding-bottom: 16px;
        border-bottom: 1px solid var(--border-color);

        h4 {
          margin: 0;
          font-size: 18px;
        }
      }

      .report-section {
        background: var(--bg-secondary);
        border-radius: 12px;
        padding: 16px;
        margin-bottom: 16px;

        &.strengths h5 {
          color: #67c23a;
        }

        &.gaps h5 {
          color: #e6a23c;
        }

        ul {
          margin: 0;
          padding-left: 20px;

          li {
            margin-bottom: 8px;
            color: var(--text-primary);
          }
        }
      }

      .recommendation {
        background: linear-gradient(135deg, #667eea20 0%, #764ba220 100%);
        border-radius: 12px;
        padding: 16px;

        h5 {
          color: #667eea;
        }

        p {
          margin: 0;
          color: var(--text-primary);
          line-height: 1.6;
        }
      }
    }
  }

  .resume-upload {
    :deep(.el-upload-dragger) {
      padding: 40px 20px;
    }
  }

  .jd-input {
    margin-top: 24px;
  }
}
</style>
