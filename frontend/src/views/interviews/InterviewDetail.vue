<template>
  <div class="interview-detail" v-loading="loading">
    <div class="page-header">
      <div class="header-left">
        <el-button :icon="ArrowLeft" @click="$router.back()">返回</el-button>
        <h1>面试详情</h1>
      </div>
      <div class="header-actions">
        <el-button v-if="interview?.status === 'scheduled'" type="warning" @click="showRescheduleDialog = true">
          <el-icon><Clock /></el-icon>改期
        </el-button>
        <el-button v-if="interview?.status === 'scheduled'" type="danger" @click="handleCancel">
          <el-icon><Close /></el-icon>取消
        </el-button>
        <el-button v-if="interview?.status === 'scheduled'" type="success" @click="showFeedbackDialog = true">
          <el-icon><Check /></el-icon>完成并反馈
        </el-button>
      </div>
    </div>

    <el-row :gutter="20" v-if="interview">
      <el-col :span="16">
        <!-- 基本信息 -->
        <div class="card info-card">
          <div class="card-header">
            <h3>面试信息</h3>
            <el-tag :type="statusType" size="large">{{ statusText }}</el-tag>
          </div>
          <div class="info-grid">
            <div class="info-item">
              <label>候选人</label>
              <span class="value highlight">{{ interview.candidate_name }}</span>
            </div>
            <div class="info-item">
              <label>应聘职位</label>
              <span class="value">{{ interview.position }}</span>
            </div>
            <div class="info-item">
              <label>面试类型</label>
              <span class="value">{{ typeText }}</span>
            </div>
            <div class="info-item">
              <label>面试方式</label>
              <span class="value">
                <el-icon v-if="interview.method === 'onsite'"><Location /></el-icon>
                <el-icon v-else-if="interview.method === 'video'"><VideoCamera /></el-icon>
                <el-icon v-else><Phone /></el-icon>
                {{ methodText }}
              </span>
            </div>
            <div class="info-item">
              <label>面试时间</label>
              <span class="value">{{ interview.date }} {{ interview.time }}</span>
            </div>
            <div class="info-item">
              <label>时长</label>
              <span class="value">{{ interview.duration }} 分钟</span>
            </div>
            <div class="info-item">
              <label>面试官</label>
              <span class="value">{{ interview.interviewer }}</span>
            </div>
            <div class="info-item full-width">
              <label>地点/链接</label>
              <span class="value">{{ interview.location || '待定' }}</span>
            </div>
            <div class="info-item full-width" v-if="interview.notes">
              <label>备注</label>
              <span class="value notes">{{ interview.notes }}</span>
            </div>
          </div>
        </div>

        <!-- 面试反馈 -->
        <div class="card feedback-card" v-if="feedbacks.length > 0 || interview.status === 'completed'">
          <div class="card-header">
            <h3>面试反馈</h3>
          </div>
          <div v-if="feedbacks.length > 0">
            <div class="feedback-item" v-for="fb in feedbacks" :key="fb.id">
              <div class="feedback-header">
                <el-rate v-model="fb.rating" disabled />
                <el-tag :type="fb.recommendation === 'pass' ? 'success' : fb.recommendation === 'fail' ? 'danger' : 'warning'">
                  {{ fb.recommendation === 'pass' ? '建议录用' : fb.recommendation === 'fail' ? '不建议录用' : '待定' }}
                </el-tag>
              </div>
              <div class="feedback-content">
                <div class="feedback-section" v-if="fb.strengths">
                  <label>优势</label>
                  <p>{{ fb.strengths }}</p>
                </div>
                <div class="feedback-section" v-if="fb.weaknesses">
                  <label>不足</label>
                  <p>{{ fb.weaknesses }}</p>
                </div>
                <div class="feedback-section" v-if="fb.comments">
                  <label>综合评价</label>
                  <p>{{ fb.comments }}</p>
                </div>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无反馈" />
        </div>
      </el-col>

      <el-col :span="8">
        <!-- 时间线 -->
        <div class="card timeline-card">
          <div class="card-header">
            <h3>面试流程</h3>
          </div>
          <el-timeline>
            <el-timeline-item color="#67c23a" :hollow="interview.status !== 'completed'">
              <p class="timeline-title">面试完成</p>
              <p class="timeline-time" v-if="interview.status === 'completed'">已完成</p>
            </el-timeline-item>
            <el-timeline-item color="#409eff">
              <p class="timeline-title">面试进行中</p>
              <p class="timeline-time">{{ interview.date }} {{ interview.time }}</p>
            </el-timeline-item>
            <el-timeline-item color="#909399">
              <p class="timeline-title">面试已安排</p>
              <p class="timeline-time">{{ formatDate(interview.created_at) }}</p>
            </el-timeline-item>
          </el-timeline>
        </div>

        <!-- 快捷操作 -->
        <div class="card actions-card">
          <div class="card-header">
            <h3>快捷操作</h3>
          </div>
          <div class="action-buttons">
            <el-button type="primary" plain @click="viewCandidate">
              <el-icon><User /></el-icon>查看候选人
            </el-button>
            <el-button type="primary" plain @click="viewJob">
              <el-icon><Suitcase /></el-icon>查看职位
            </el-button>
            <el-button type="primary" plain @click="sendReminder">
              <el-icon><Bell /></el-icon>发送提醒
            </el-button>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 反馈对话框 -->
    <el-dialog v-model="showFeedbackDialog" title="提交面试反馈" width="600px">
      <el-form :model="feedbackForm" label-width="100px">
        <el-form-item label="评分" required>
          <el-rate v-model="feedbackForm.rating" show-text :texts="['很差', '较差', '一般', '良好', '优秀']" />
        </el-form-item>
        <el-form-item label="候选人优势">
          <el-input v-model="feedbackForm.strengths" type="textarea" :rows="3" placeholder="请描述候选人的优势" />
        </el-form-item>
        <el-form-item label="候选人不足">
          <el-input v-model="feedbackForm.weaknesses" type="textarea" :rows="3" placeholder="请描述候选人的不足" />
        </el-form-item>
        <el-form-item label="综合评价">
          <el-input v-model="feedbackForm.comments" type="textarea" :rows="3" placeholder="请输入综合评价" />
        </el-form-item>
        <el-form-item label="录用建议" required>
          <el-radio-group v-model="feedbackForm.recommendation">
            <el-radio label="pass">建议录用</el-radio>
            <el-radio label="pending">待定</el-radio>
            <el-radio label="fail">不建议录用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFeedbackDialog = false">取消</el-button>
        <el-button type="primary" @click="submitFeedback" :loading="submitting">提交</el-button>
      </template>
    </el-dialog>

    <!-- 改期对话框 -->
    <el-dialog v-model="showRescheduleDialog" title="重新安排面试" width="500px">
      <el-form :model="rescheduleForm" label-width="100px">
        <el-form-item label="新日期" required>
          <el-date-picker v-model="rescheduleForm.date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="新时间" required>
          <el-time-picker v-model="rescheduleForm.time" placeholder="选择时间" format="HH:mm" value-format="HH:mm" style="width: 100%" />
        </el-form-item>
        <el-form-item label="改期原因">
          <el-input v-model="rescheduleForm.reason" type="textarea" :rows="2" placeholder="请输入改期原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRescheduleDialog = false">取消</el-button>
        <el-button type="primary" @click="handleReschedule" :loading="submitting">确认改期</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Clock, Close, Check, Location, VideoCamera, Phone, User, Suitcase, Bell } from '@element-plus/icons-vue'
import { interviewApi, type Interview, type InterviewFeedback } from '@/api/interview'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const interview = ref<Interview | null>(null)
const feedbacks = ref<InterviewFeedback[]>([])
const showFeedbackDialog = ref(false)
const showRescheduleDialog = ref(false)

const feedbackForm = ref({
  rating: 3,
  strengths: '',
  weaknesses: '',
  comments: '',
  recommendation: 'pending' as 'pass' | 'fail' | 'pending'
})

const rescheduleForm = ref({
  date: '',
  time: '',
  reason: ''
})

const statusType = computed(() => {
  const map: Record<string, string> = {
    scheduled: 'primary',
    completed: 'success',
    cancelled: 'danger',
    no_show: 'warning'
  }
  return map[interview.value?.status || ''] || 'info'
})

const statusText = computed(() => {
  const map: Record<string, string> = {
    scheduled: '已安排',
    completed: '已完成',
    cancelled: '已取消',
    no_show: '爽约'
  }
  return map[interview.value?.status || ''] || '未知'
})

const typeText = computed(() => {
  const map: Record<string, string> = {
    initial: '初试',
    second: '复试',
    final: '终面',
    hr: 'HR面试'
  }
  return map[interview.value?.type || ''] || '未知'
})

const methodText = computed(() => {
  const map: Record<string, string> = {
    onsite: '现场面试',
    video: '视频面试',
    phone: '电话面试'
  }
  return map[interview.value?.method || ''] || '未知'
})

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const fetchInterview = async () => {
  loading.value = true
  try {
    const id = Number(route.params.id)
    const [interviewRes, feedbackRes] = await Promise.all([
      interviewApi.get(id),
      interviewApi.getFeedback(id)
    ])
    if (interviewRes.data.code === 0) {
      interview.value = interviewRes.data.data
    }
    if (feedbackRes.data.code === 0) {
      feedbacks.value = feedbackRes.data.data || []
    }
  } catch (error) {
    ElMessage.error('获取面试详情失败')
  } finally {
    loading.value = false
  }
}

const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定要取消这场面试吗？', '确认取消')
    await interviewApi.cancel(interview.value!.id)
    ElMessage.success('面试已取消')
    fetchInterview()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('取消失败')
    }
  }
}

const submitFeedback = async () => {
  if (!feedbackForm.value.rating) {
    ElMessage.warning('请选择评分')
    return
  }
  submitting.value = true
  try {
    await interviewApi.submitFeedback(interview.value!.id, feedbackForm.value)
    ElMessage.success('反馈提交成功')
    showFeedbackDialog.value = false
    fetchInterview()
  } catch (error) {
    ElMessage.error('提交失败')
  } finally {
    submitting.value = false
  }
}

const handleReschedule = async () => {
  if (!rescheduleForm.value.date || !rescheduleForm.value.time) {
    ElMessage.warning('请选择新的日期和时间')
    return
  }
  submitting.value = true
  try {
    await interviewApi.reschedule(interview.value!.id, rescheduleForm.value)
    ElMessage.success('面试已改期')
    showRescheduleDialog.value = false
    fetchInterview()
  } catch (error) {
    ElMessage.error('改期失败')
  } finally {
    submitting.value = false
  }
}

const viewCandidate = () => {
  router.push(`/talents/${interview.value?.candidate_id}`)
}

const viewJob = () => {
  router.push(`/jobs/${interview.value?.position_id}`)
}

const sendReminder = () => {
  ElMessage.success('提醒已发送')
}

onMounted(() => {
  fetchInterview()
})
</script>

<style scoped lang="scss">
.interview-detail {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    .header-left {
      display: flex;
      align-items: center;
      gap: 16px;

      h1 {
        margin: 0;
        font-size: 24px;
      }
    }

    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .card {
    background: white;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 20px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;

      h3 {
        margin: 0;
        font-size: 18px;
        color: #1a1a2e;
      }
    }
  }

  .info-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;

    .info-item {
      &.full-width {
        grid-column: span 2;
      }

      label {
        display: block;
        font-size: 13px;
        color: #909399;
        margin-bottom: 6px;
      }

      .value {
        font-size: 15px;
        color: #303133;
        display: flex;
        align-items: center;
        gap: 6px;

        &.highlight {
          font-weight: 600;
          color: #409eff;
        }

        &.notes {
          white-space: pre-wrap;
          line-height: 1.6;
        }
      }
    }
  }

  .feedback-item {
    padding: 16px;
    background: #f9fafb;
    border-radius: 8px;
    margin-bottom: 12px;

    .feedback-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
    }

    .feedback-section {
      margin-bottom: 12px;

      label {
        display: block;
        font-size: 13px;
        color: #909399;
        margin-bottom: 4px;
      }

      p {
        margin: 0;
        color: #303133;
        line-height: 1.6;
      }
    }
  }

  .timeline-card {
    .timeline-title {
      font-weight: 500;
      margin: 0 0 4px 0;
    }

    .timeline-time {
      font-size: 12px;
      color: #909399;
      margin: 0;
    }
  }

  .actions-card {
    .action-buttons {
      display: flex;
      flex-direction: column;
      gap: 12px;

      .el-button {
        justify-content: flex-start;
      }
    }
  }
}
</style>
