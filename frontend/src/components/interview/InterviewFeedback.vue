<template>
  <div class="interview-feedback">
    <div class="feedback-header">
      <h3>面试评价</h3>
      <el-tag :type="getStatusType(interview.status)" size="small">
        {{ getStatusLabel(interview.status) }}
      </el-tag>
    </div>

    <!-- 候选人信息 -->
    <div class="candidate-info">
      <el-avatar :size="56" :style="{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }">
        {{ interview.candidateName?.charAt(0) }}
      </el-avatar>
      <div class="info">
        <h4>{{ interview.candidateName }}</h4>
        <p>{{ interview.position }}</p>
        <span class="meta">{{ interview.date }} {{ interview.time }} · {{ getTypeLabel(interview.type) }}</span>
      </div>
    </div>

    <el-divider />

    <!-- 评分区域 -->
    <div class="rating-section">
      <h4>综合评分</h4>
      <div class="overall-rating">
        <el-rate
          v-model="feedback.overallRating"
          :colors="ratingColors"
          :texts="ratingTexts"
          show-text
          size="large"
        />
      </div>

      <div class="dimension-ratings">
        <div class="dimension-item" v-for="dim in dimensions" :key="dim.key">
          <span class="dim-label">{{ dim.label }}</span>
          <el-rate
            v-model="feedback.dimensions[dim.key]"
            :colors="ratingColors"
            size="default"
          />
          <span class="dim-score">{{ feedback.dimensions[dim.key] || 0 }}分</span>
        </div>
      </div>
    </div>

    <el-divider />

    <!-- 评价内容 -->
    <div class="feedback-content">
      <div class="content-item">
        <h4>
          <el-icon><CircleCheck /></el-icon>
          优势亮点
        </h4>
        <el-input
          v-model="feedback.strengths"
          type="textarea"
          :rows="3"
          placeholder="请描述候选人的优势和亮点..."
        />
      </div>

      <div class="content-item">
        <h4>
          <el-icon><Warning /></el-icon>
          待改进项
        </h4>
        <el-input
          v-model="feedback.weaknesses"
          type="textarea"
          :rows="3"
          placeholder="请描述候选人需要改进的地方..."
        />
      </div>

      <div class="content-item">
        <h4>
          <el-icon><ChatDotRound /></el-icon>
          面试记录
        </h4>
        <el-input
          v-model="feedback.notes"
          type="textarea"
          :rows="4"
          placeholder="记录面试过程中的关键问答和表现..."
        />
      </div>
    </div>

    <el-divider />

    <!-- 录用建议 -->
    <div class="recommendation-section">
      <h4>录用建议</h4>
      <el-radio-group v-model="feedback.recommendation" class="recommendation-options">
        <el-radio-button label="strong_hire">
          <el-icon><CircleCheckFilled /></el-icon>
          强烈推荐
        </el-radio-button>
        <el-radio-button label="hire">
          <el-icon><Select /></el-icon>
          建议录用
        </el-radio-button>
        <el-radio-button label="pending">
          <el-icon><Clock /></el-icon>
          待定
        </el-radio-button>
        <el-radio-button label="no_hire">
          <el-icon><CloseBold /></el-icon>
          不建议录用
        </el-radio-button>
      </el-radio-group>

      <div class="recommendation-reason" v-if="feedback.recommendation">
        <el-input
          v-model="feedback.recommendationReason"
          type="textarea"
          :rows="2"
          :placeholder="getRecommendationPlaceholder()"
        />
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="feedback-actions">
      <el-button @click="saveDraft">保存草稿</el-button>
      <el-button type="primary" @click="submitFeedback" :loading="submitting">
        提交评价
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, defineProps, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'
import {
  CircleCheck, Warning, ChatDotRound, CircleCheckFilled,
  Select, Clock, CloseBold
} from '@element-plus/icons-vue'

interface Interview {
  id: number
  candidateName: string
  position: string
  date: string
  time: string
  type: string
  status: string
}

const props = defineProps<{
  interview: Interview
}>()

const emit = defineEmits(['submit', 'save-draft'])

const submitting = ref(false)

const ratingColors = ['#f5576c', '#f59e0b', '#f59e0b', '#43e97b', '#43e97b']
const ratingTexts = ['不合格', '待提升', '一般', '良好', '优秀']

const dimensions = [
  { key: 'technical', label: '专业技能' },
  { key: 'communication', label: '沟通表达' },
  { key: 'logic', label: '逻辑思维' },
  { key: 'teamwork', label: '团队协作' },
  { key: 'potential', label: '发展潜力' }
]

const feedback = reactive({
  overallRating: 0,
  dimensions: {
    technical: 0,
    communication: 0,
    logic: 0,
    teamwork: 0,
    potential: 0
  } as Record<string, number>,
  strengths: '',
  weaknesses: '',
  notes: '',
  recommendation: '',
  recommendationReason: ''
})

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    scheduled: 'primary',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    scheduled: '待进行',
    completed: '已完成',
    cancelled: '已取消'
  }
  return labels[status] || status
}

const getTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    initial: '初试',
    second: '复试',
    final: '终面',
    hr: 'HR面'
  }
  return labels[type] || type
}

const getRecommendationPlaceholder = () => {
  const placeholders: Record<string, string> = {
    strong_hire: '请说明强烈推荐的理由...',
    hire: '请说明建议录用的理由...',
    pending: '请说明需要进一步考察的方面...',
    no_hire: '请说明不建议录用的原因...'
  }
  return placeholders[feedback.recommendation] || '请说明理由...'
}

const saveDraft = () => {
  emit('save-draft', { ...feedback, interviewId: props.interview.id })
  ElMessage.success('草稿已保存')
}

const submitFeedback = async () => {
  if (feedback.overallRating === 0) {
    ElMessage.warning('请填写综合评分')
    return
  }
  if (!feedback.recommendation) {
    ElMessage.warning('请选择录用建议')
    return
  }

  submitting.value = true
  try {
    emit('submit', { ...feedback, interviewId: props.interview.id })
    ElMessage.success('评价已提交')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped lang="scss">
.interview-feedback {
  padding: 24px;
  background: var(--bg-primary);
  border-radius: 16px;
}

.feedback-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }
}

.candidate-info {
  display: flex;
  align-items: center;
  gap: 16px;

  .info {
    h4 {
      font-size: 18px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 4px 0;
    }

    p {
      font-size: 14px;
      color: var(--text-secondary);
      margin: 0 0 4px 0;
    }

    .meta {
      font-size: 12px;
      color: var(--text-tertiary);
    }
  }
}

.rating-section {
  h4 {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 16px 0;
  }

  .overall-rating {
    text-align: center;
    padding: 20px;
    background: var(--bg-tertiary);
    border-radius: 12px;
    margin-bottom: 20px;

    :deep(.el-rate) {
      height: auto;
    }

    :deep(.el-rate__icon) {
      font-size: 32px;
    }
  }

  .dimension-ratings {
    display: grid;
    gap: 12px;

    .dimension-item {
      display: flex;
      align-items: center;
      gap: 12px;

      .dim-label {
        width: 80px;
        font-size: 14px;
        color: var(--text-secondary);
      }

      .dim-score {
        width: 40px;
        font-size: 14px;
        font-weight: 600;
        color: var(--primary-color);
      }
    }
  }
}

.feedback-content {
  .content-item {
    margin-bottom: 20px;

    &:last-child {
      margin-bottom: 0;
    }

    h4 {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 12px 0;
    }
  }
}

.recommendation-section {
  h4 {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 16px 0;
  }

  .recommendation-options {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;

    :deep(.el-radio-button__inner) {
      display: flex;
      align-items: center;
      gap: 6px;
    }
  }

  .recommendation-reason {
    margin-top: 16px;
  }
}

.feedback-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--border-light);
}
</style>
