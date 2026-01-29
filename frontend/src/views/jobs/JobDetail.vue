<template>
  <div class="job-detail-page" v-loading="loading">
    <!-- 返回按钮 -->
    <div class="page-header">
      <el-button text @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
    </div>

    <!-- 职位信息卡片 -->
    <div class="job-info-card card" v-if="job">
      <div class="job-header">
        <div class="job-main">
          <h1 class="job-title">{{ job.title }}</h1>
          <div class="job-salary">{{ job.salary }}</div>
        </div>
        <div class="job-actions">
          <el-tag :type="getStatusType(job.status)" size="large" effect="light">
            {{ getStatusText(job.status) }}
          </el-tag>
          <el-button type="primary" @click="openEditDialog" v-if="canEdit">
            <el-icon><Edit /></el-icon>
            编辑职位
          </el-button>
        </div>
      </div>

      <div class="job-meta">
        <span class="meta-item">
          <el-icon><Location /></el-icon>
          {{ job.location }}
        </span>
        <span class="meta-item">
          <el-icon><OfficeBuilding /></el-icon>
          {{ job.department }}
        </span>
        <span class="meta-item">
          <el-icon><Timer /></el-icon>
          {{ getJobType(job.type) }}
        </span>
        <span class="meta-item">
          <el-icon><User /></el-icon>
          {{ applicationsTotal }} 人申请
        </span>
      </div>
    </div>

    <!-- Tab 切换 -->
    <el-tabs v-model="activeTab" class="detail-tabs">
      <!-- 职位详情 Tab -->
      <el-tab-pane label="职位详情" name="info">
        <div class="tab-content" v-if="job">
          <div class="detail-section">
            <h3>技能要求</h3>
            <div class="skills-list">
              <el-tag v-for="skill in job.skills" :key="skill" type="info">
                {{ skill }}
              </el-tag>
              <span v-if="!job.skills?.length" class="empty-text">暂无技能要求</span>
            </div>
          </div>

          <div class="detail-section">
            <h3>职位描述</h3>
            <p class="description">{{ job.description || '暂无描述' }}</p>
          </div>

          <div class="detail-section" v-if="job.requirements?.length">
            <h3>任职要求</h3>
            <ul class="requirements-list">
              <li v-for="(req, index) in job.requirements" :key="index">{{ req }}</li>
            </ul>
          </div>

          <div class="detail-section" v-if="job.benefits?.length">
            <h3>福利待遇</h3>
            <div class="benefits-list">
              <el-tag v-for="benefit in job.benefits" :key="benefit" type="success" effect="light">
                {{ benefit }}
              </el-tag>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- 申请管理 Tab -->
      <el-tab-pane label="申请管理" name="applications">
        <div class="tab-content">
          <!-- 筛选栏 -->
          <div class="filter-bar">
            <el-select v-model="statusFilter" placeholder="筛选状态" clearable style="width: 150px" @change="fetchApplications">
              <el-option label="待处理" value="pending" />
              <el-option label="审核中" value="reviewing" />
              <el-option label="面试中" value="interview" />
              <el-option label="已录用" value="offer" />
              <el-option label="已拒绝" value="rejected" />
            </el-select>
            <span class="filter-info">共 {{ applicationsTotal }} 条申请</span>
          </div>

          <!-- 申请列表 -->
          <div class="applications-list" v-loading="applicationsLoading">
            <div 
              class="application-card" 
              v-for="app in applications" 
              :key="app.id"
            >
              <div class="applicant-info">
                <el-avatar :size="48" :icon="UserFilled" />
                <div class="applicant-details">
                  <h4 class="applicant-name">{{ app.candidate_name || '未知候选人' }}</h4>
                  <div class="applicant-contact">
                    <span v-if="app.candidate_email">
                      <el-icon><Message /></el-icon>
                      {{ app.candidate_email }}
                    </span>
                    <span v-if="app.candidate_phone">
                      <el-icon><Phone /></el-icon>
                      {{ app.candidate_phone }}
                    </span>
                  </div>
                  <p class="resume-summary" v-if="app.resume_summary">{{ app.resume_summary }}</p>
                </div>
              </div>

              <div class="application-meta">
                <el-tag :type="getApplicationStatusType(app.status)" effect="light">
                  {{ getApplicationStatusText(app.status) }}
                </el-tag>
                <span class="apply-time">
                  <el-icon><Clock /></el-icon>
                  {{ formatDate(app.created_at) }}
                </span>
              </div>

              <div class="application-actions">
                <el-dropdown trigger="click" @command="(cmd: string) => handleStatusChange(app, cmd)">
                  <el-button type="primary" size="small">
                    更新状态
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="pending" :disabled="app.status === 'pending'">
                        <el-icon><Clock /></el-icon> 待处理
                      </el-dropdown-item>
                      <el-dropdown-item command="reviewing" :disabled="app.status === 'reviewing'">
                        <el-icon><View /></el-icon> 审核中
                      </el-dropdown-item>
                      <el-dropdown-item command="interview" :disabled="app.status === 'interview'">
                        <el-icon><Calendar /></el-icon> 面试中
                      </el-dropdown-item>
                      <el-dropdown-item command="offer" :disabled="app.status === 'offer'">
                        <el-icon><CircleCheck /></el-icon> 已录用
                      </el-dropdown-item>
                      <el-dropdown-item command="rejected" :disabled="app.status === 'rejected'">
                        <el-icon><CircleClose /></el-icon> 已拒绝
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>

                <el-button size="small" @click="viewApplicantDetail(app)">
                  <el-icon><Document /></el-icon>
                  查看详情
                </el-button>
                <el-button size="small" type="success" @click="openScheduleDialog(app)" v-if="app.status !== 'rejected' && app.status !== 'offer'">
                  <el-icon><Calendar /></el-icon>
                  安排面试
                </el-button>
              </div>
            </div>

            <el-empty v-if="applications.length === 0 && !applicationsLoading" description="暂无申请记录" />
          </div>

          <!-- 分页 -->
          <div class="pagination" v-if="applicationsTotal > 0">
            <el-pagination
              v-model:current-page="applicationsPage"
              v-model:page-size="applicationsPageSize"
              :total="applicationsTotal"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next"
              @current-change="fetchApplications"
              @size-change="fetchApplications"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 申请人详情抽屉 -->
    <el-drawer
      v-model="applicantDrawerVisible"
      title="申请人详情"
      size="500px"
    >
      <div class="applicant-detail" v-if="currentApplicant">
        <div class="detail-header">
          <el-avatar :size="64" :icon="UserFilled" />
          <div class="header-info">
            <h2>{{ currentApplicant.candidate_name || '未知候选人' }}</h2>
            <el-tag :type="getApplicationStatusType(currentApplicant.status)" effect="light">
              {{ getApplicationStatusText(currentApplicant.status) }}
            </el-tag>
          </div>
        </div>

        <el-divider />

        <div class="info-section">
          <h4>联系方式</h4>
          <p v-if="currentApplicant.candidate_email">
            <el-icon><Message /></el-icon>
            {{ currentApplicant.candidate_email }}
          </p>
          <p v-if="currentApplicant.candidate_phone">
            <el-icon><Phone /></el-icon>
            {{ currentApplicant.candidate_phone }}
          </p>
        </div>

        <div class="info-section" v-if="currentApplicant.resume_summary">
          <h4>简历摘要</h4>
          <p>{{ currentApplicant.resume_summary }}</p>
        </div>

        <div class="info-section" v-if="currentApplicant.cover_letter">
          <h4>求职信</h4>
          <p>{{ currentApplicant.cover_letter }}</p>
        </div>

        <div class="info-section">
          <h4>申请时间</h4>
          <p>{{ formatDate(currentApplicant.created_at) }}</p>
        </div>

        <div class="info-section" v-if="currentApplicant.notes">
          <h4>备注</h4>
          <p>{{ currentApplicant.notes }}</p>
        </div>

        <div class="drawer-actions">
          <el-dropdown trigger="click" @command="(cmd: string) => currentApplicant && handleStatusChange(currentApplicant, cmd)">
            <el-button type="primary">
              更新状态
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="pending">待处理</el-dropdown-item>
                <el-dropdown-item command="reviewing">审核中</el-dropdown-item>
                <el-dropdown-item command="interview">面试中</el-dropdown-item>
                <el-dropdown-item command="offer">已录用</el-dropdown-item>
                <el-dropdown-item command="rejected">已拒绝</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button @click="currentApplicant && openScheduleDialog(currentApplicant)" v-if="currentApplicant.status !== 'rejected' && currentApplicant.status !== 'offer'">
            <el-icon><Calendar /></el-icon>
            安排面试
          </el-button>
        </div>
      </div>
    </el-drawer>

    <!-- 面试安排对话框 - Requirement 7.1, 7.2 -->
    <el-dialog
      v-model="scheduleDialogVisible"
      title="安排面试"
      width="600px"
      destroy-on-close
    >
      <el-form 
        ref="scheduleFormRef"
        :model="scheduleForm" 
        :rules="scheduleRules"
        label-width="100px"
      >
        <el-form-item label="候选人">
          <el-input :value="scheduleForm.candidate_name" disabled />
        </el-form-item>
        <el-form-item label="应聘职位">
          <el-input :value="job?.title" disabled />
        </el-form-item>
        <el-form-item label="面试类型" prop="type">
          <el-select v-model="scheduleForm.type" placeholder="选择面试类型" style="width: 100%">
            <el-option label="初试" value="initial" />
            <el-option label="复试" value="second" />
            <el-option label="终面" value="final" />
            <el-option label="HR面试" value="hr" />
          </el-select>
        </el-form-item>
        <el-form-item label="面试日期" prop="date">
          <el-date-picker
            v-model="scheduleForm.date"
            type="date"
            placeholder="选择日期"
            style="width: 100%"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            :disabled-date="disablePastDate"
          />
        </el-form-item>
        <el-form-item label="面试时间" prop="time">
          <el-time-picker
            v-model="scheduleForm.time"
            placeholder="选择时间"
            style="width: 100%"
            format="HH:mm"
            value-format="HH:mm"
          />
        </el-form-item>
        <el-form-item label="时长(分钟)" prop="duration">
          <el-input-number v-model="scheduleForm.duration" :min="30" :max="180" :step="30" />
        </el-form-item>
        <el-form-item label="面试官" prop="interviewer">
          <el-input v-model="scheduleForm.interviewer" placeholder="输入面试官姓名" />
        </el-form-item>
        <el-form-item label="面试方式" prop="method">
          <el-radio-group v-model="scheduleForm.method">
            <el-radio value="onsite">现场面试</el-radio>
            <el-radio value="video">视频面试</el-radio>
            <el-radio value="phone">电话面试</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="地点/链接" prop="location">
          <el-input 
            v-model="scheduleForm.location" 
            :placeholder="scheduleForm.method === 'onsite' ? '输入会议室地点' : scheduleForm.method === 'video' ? '输入视频会议链接' : '输入联系电话'" 
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="scheduleForm.notes" type="textarea" :rows="3" placeholder="面试注意事项..." />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="scheduleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitSchedule" :loading="scheduleSubmitting">
          确认安排
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>


<script setup lang="ts">
import { ref, computed, onMounted, watch, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { jobApi } from '@/api/job'
import { applicationApi } from '@/api/application'
import { interviewApi } from '@/api/interview'
import type { Job } from '@/types'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft, Edit, Location, OfficeBuilding, Timer, User, UserFilled,
  Message, Phone, Clock, ArrowDown, View, Calendar, CircleCheck, CircleClose,
  Document
} from '@element-plus/icons-vue'
import { usePermissionStore } from '@/store/permission'

// Types
interface ApplicationWithCandidate {
  id: number
  job_id: number
  talent_id: number
  resume_id: number
  status: string
  cover_letter: string
  notes: string
  created_at: string
  updated_at: string
  candidate_name: string
  candidate_email: string
  candidate_phone: string
  resume_summary: string
}

interface ScheduleForm {
  application_id: number
  candidate_id: number
  candidate_name: string
  position_id: number
  position: string
  type: string
  date: string
  time: string
  duration: number
  interviewer_id: number
  interviewer: string
  method: string
  location: string
  notes: string
}

const route = useRoute()
const router = useRouter()
const permissionStore = usePermissionStore()

// 权限检查
const canEdit = computed(() => permissionStore.hasPermission('job:edit'))

// 状态
const loading = ref(false)
const job = ref<Job | null>(null)
const activeTab = ref('info')

// 申请管理状态
const applicationsLoading = ref(false)
const applications = ref<ApplicationWithCandidate[]>([])
const applicationsPage = ref(1)
const applicationsPageSize = ref(10)
const applicationsTotal = ref(0)
const statusFilter = ref('')

// 申请人详情抽屉
const applicantDrawerVisible = ref(false)
const currentApplicant = ref<ApplicationWithCandidate | null>(null)

// 面试安排对话框状态 - Requirement 7.1, 7.2
const scheduleDialogVisible = ref(false)
const scheduleSubmitting = ref(false)
const scheduleFormRef = ref<FormInstance>()
const scheduleForm = reactive<ScheduleForm>({
  application_id: 0,
  candidate_id: 0,
  candidate_name: '',
  position_id: 0,
  position: '',
  type: 'initial',
  date: '',
  time: '',
  duration: 60,
  interviewer_id: 0,
  interviewer: '',
  method: 'onsite',
  location: '',
  notes: ''
})

// 面试安排表单验证规则 - Requirement 7.5: 日期必须是未来时间
const validateFutureDate = (_rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('请选择面试日期'))
    return
  }
  const selectedDate = new Date(value)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  if (selectedDate < today) {
    callback(new Error('面试日期必须是未来时间'))
  } else {
    callback()
  }
}

const validateFutureTime = (_rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('请选择面试时间'))
    return
  }
  if (scheduleForm.date) {
    const now = new Date()
    const selectedDateTime = new Date(`${scheduleForm.date}T${value}`)
    if (selectedDateTime <= now) {
      callback(new Error('面试时间必须是未来时间'))
    } else {
      callback()
    }
  }
  callback()
}

const scheduleRules = reactive<FormRules>({
  type: [{ required: true, message: '请选择面试类型', trigger: 'change' }],
  date: [{ required: true, validator: validateFutureDate, trigger: 'change' }],
  time: [{ required: true, validator: validateFutureTime, trigger: 'change' }],
  duration: [{ required: true, message: '请设置面试时长', trigger: 'change' }],
  interviewer: [{ required: true, message: '请输入面试官姓名', trigger: 'blur' }],
  method: [{ required: true, message: '请选择面试方式', trigger: 'change' }],
  location: [{ required: true, message: '请输入面试地点或链接', trigger: 'blur' }]
})

// 禁用过去的日期 - Requirement 7.5
const disablePastDate = (date: Date) => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return date < today
}

// 获取职位ID
const jobId = computed(() => Number(route.params.id))

// 获取职位详情
const fetchJob = async () => {
  if (!jobId.value) return
  
  loading.value = true
  try {
    const res = await jobApi.get(jobId.value)
    if (res.data.code === 0 && res.data.data) {
      job.value = res.data.data
    } else {
      ElMessage.error(res.data?.message || '获取职位详情失败')
    }
  } catch (error) {
    console.error('获取职位详情失败:', error)
    ElMessage.error('获取职位详情失败')
  } finally {
    loading.value = false
  }
}

// 获取申请列表
const fetchApplications = async () => {
  if (!jobId.value) return
  
  applicationsLoading.value = true
  try {
    const res = await applicationApi.getJobApplications(jobId.value, {
      page: applicationsPage.value,
      page_size: applicationsPageSize.value,
      status: statusFilter.value || undefined
    })
    
    if (res.data.code === 0 && res.data.data) {
      applications.value = res.data.data.applications || []
      applicationsTotal.value = res.data.data.total || 0
    } else {
      ElMessage.error(res.data?.message || '获取申请列表失败')
    }
  } catch (error) {
    console.error('获取申请列表失败:', error)
    ElMessage.error('获取申请列表失败')
    applications.value = []
    applicationsTotal.value = 0
  } finally {
    applicationsLoading.value = false
  }
}

// 更新申请状态
const handleStatusChange = async (app: ApplicationWithCandidate, newStatus: string) => {
  try {
    await ElMessageBox.confirm(
      `确定要将申请状态更新为"${getApplicationStatusText(newStatus)}"吗？`,
      '确认更新',
      { type: 'warning' }
    )
    
    const res = await applicationApi.update(app.id, { status: newStatus })
    if (res.data.code === 0) {
      ElMessage.success('状态更新成功')
      // 更新本地状态
      app.status = newStatus
      if (currentApplicant.value?.id === app.id) {
        currentApplicant.value.status = newStatus
      }
    } else {
      ElMessage.error(res.data?.message || '状态更新失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('状态更新失败:', error)
      ElMessage.error('状态更新失败')
    }
  }
}

// 查看申请人详情
const viewApplicantDetail = (app: ApplicationWithCandidate) => {
  currentApplicant.value = app
  applicantDrawerVisible.value = true
}

// 打开面试安排对话框 - Requirement 7.1: HR可以从申请页面安排面试
const openScheduleDialog = (app: ApplicationWithCandidate) => {
  // 重置表单
  scheduleForm.application_id = app.id
  scheduleForm.candidate_id = app.talent_id
  scheduleForm.candidate_name = app.candidate_name || '未知候选人'
  scheduleForm.position_id = app.job_id
  scheduleForm.position = job.value?.title || ''
  scheduleForm.type = 'initial'
  scheduleForm.date = ''
  scheduleForm.time = ''
  scheduleForm.duration = 60
  scheduleForm.interviewer_id = 0
  scheduleForm.interviewer = ''
  scheduleForm.method = 'onsite'
  scheduleForm.location = ''
  scheduleForm.notes = ''
  
  // 关闭抽屉，打开对话框
  applicantDrawerVisible.value = false
  scheduleDialogVisible.value = true
}

// 提交面试安排 - Requirement 7.2, 7.3, 7.5
const submitSchedule = async () => {
  if (!scheduleFormRef.value) return
  
  try {
    await scheduleFormRef.value.validate()
  } catch {
    return
  }
  
  scheduleSubmitting.value = true
  try {
    const res = await interviewApi.create({
      candidate_id: scheduleForm.candidate_id,
      candidate_name: scheduleForm.candidate_name,
      position_id: scheduleForm.position_id,
      position: scheduleForm.position,
      type: scheduleForm.type,
      date: scheduleForm.date,
      time: scheduleForm.time,
      duration: scheduleForm.duration,
      interviewer_id: scheduleForm.interviewer_id || 1, // 默认面试官ID
      interviewer: scheduleForm.interviewer,
      method: scheduleForm.method,
      location: scheduleForm.location,
      notes: scheduleForm.notes,
      application_id: scheduleForm.application_id
    })
    
    if (res.data.code === 0) {
      ElMessage.success('面试安排成功，已发送通知给候选人')
      scheduleDialogVisible.value = false
      
      // 更新申请状态为"面试中"
      const app = applications.value.find(a => a.id === scheduleForm.application_id)
      if (app && app.status !== 'interview') {
        await applicationApi.update(app.id, { status: 'interview' })
        app.status = 'interview'
      }
    } else {
      // 处理特定错误码
      if (res.data.code === 1005) {
        ElMessage.error(res.data.message || '面试日期必须是未来时间')
      } else {
        ElMessage.error(res.data?.message || '面试安排失败')
      }
    }
  } catch (error: any) {
    console.error('面试安排失败:', error)
    const errorMsg = error.response?.data?.message || '面试安排失败'
    ElMessage.error(errorMsg)
  } finally {
    scheduleSubmitting.value = false
  }
}

// 打开编辑弹窗
const openEditDialog = () => {
  router.push(`/jobs/${jobId.value}/edit`)
}

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取职位状态类型
const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    open: 'success',
    closed: 'info',
    filled: 'warning'
  }
  return map[status] || ''
}

// 获取职位状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    open: '招聘中',
    closed: '已关闭',
    filled: '已满员'
  }
  return map[status] || status
}

// 获取职位类型
const getJobType = (type: string) => {
  const map: Record<string, string> = {
    'full-time': '全职',
    'part-time': '兼职',
    'contract': '合同',
    'internship': '实习'
  }
  return map[type] || type
}

// 获取申请状态类型
const getApplicationStatusType = (status: string) => {
  const map: Record<string, any> = {
    pending: 'info',
    reviewing: 'warning',
    interview: 'primary',
    offer: 'success',
    rejected: 'danger'
  }
  return map[status] || 'info'
}

// 获取申请状态文本
const getApplicationStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    reviewing: '审核中',
    interview: '面试中',
    offer: '已录用',
    rejected: '已拒绝'
  }
  return map[status] || status
}

// 监听 Tab 切换，切换到申请管理时加载数据
watch(activeTab, (newTab) => {
  if (newTab === 'applications' && applications.value.length === 0) {
    fetchApplications()
  }
})

// 初始化
onMounted(() => {
  fetchJob()
  // 预加载申请数量
  fetchApplications()
})
</script>


<style scoped lang="scss">
.job-detail-page {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.page-header {
  margin-bottom: 20px;
}

.card {
  background: var(--bg-primary);
  border-radius: 12px;
  padding: 24px;
  box-shadow: var(--shadow-card);
  border: 1px solid var(--border-light);
}

.job-info-card {
  margin-bottom: 20px;

  .job-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 20px;

    .job-main {
      .job-title {
        font-size: 28px;
        font-weight: 700;
        color: var(--text-primary);
        margin: 0 0 8px 0;
      }

      .job-salary {
        font-size: 24px;
        font-weight: 700;
        color: var(--primary-color);
      }
    }

    .job-actions {
      display: flex;
      align-items: center;
      gap: 12px;
    }
  }

  .job-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 24px;

    .meta-item {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 14px;
      color: var(--text-secondary);
    }
  }
}

.detail-tabs {
  background: var(--bg-primary);
  border-radius: 12px;
  padding: 20px;
  box-shadow: var(--shadow-card);
  border: 1px solid var(--border-light);
}

.tab-content {
  padding: 20px 0;
}

.detail-section {
  margin-bottom: 28px;

  h3 {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 16px 0;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--border-light);
  }

  .skills-list, .benefits-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .description {
    font-size: 14px;
    color: var(--text-secondary);
    line-height: 1.8;
    margin: 0;
    white-space: pre-wrap;
  }

  .requirements-list {
    margin: 0;
    padding-left: 20px;

    li {
      font-size: 14px;
      color: var(--text-secondary);
      line-height: 2;
    }
  }

  .empty-text {
    color: var(--text-secondary);
    font-size: 14px;
  }
}

// 申请管理样式
.filter-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;

  .filter-info {
    color: var(--text-secondary);
    font-size: 14px;
  }
}

.applications-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.application-card {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 20px;
  border: 1px solid var(--border-light);
  transition: all 0.3s ease;

  &:hover {
    border-color: var(--primary-color);
    box-shadow: var(--shadow-md);
  }

  .applicant-info {
    display: flex;
    gap: 16px;
    margin-bottom: 16px;

    .applicant-details {
      flex: 1;

      .applicant-name {
        font-size: 16px;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 8px 0;
      }

      .applicant-contact {
        display: flex;
        flex-wrap: wrap;
        gap: 16px;
        margin-bottom: 8px;

        span {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 13px;
          color: var(--text-secondary);
        }
      }

      .resume-summary {
        font-size: 13px;
        color: var(--text-secondary);
        margin: 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }
  }

  .application-meta {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 16px;

    .apply-time {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 13px;
      color: var(--text-secondary);
    }
  }

  .application-actions {
    display: flex;
    gap: 8px;
    padding-top: 16px;
    border-top: 1px solid var(--border-light);
  }
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

// 申请人详情抽屉
.applicant-detail {
  .detail-header {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;

    .header-info {
      h2 {
        font-size: 20px;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 8px 0;
      }
    }
  }

  .info-section {
    margin-bottom: 20px;

    h4 {
      font-size: 14px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 8px 0;
    }

    p {
      font-size: 14px;
      color: var(--text-secondary);
      margin: 0;
      line-height: 1.6;
      display: flex;
      align-items: center;
      gap: 6px;
    }
  }

  .drawer-actions {
    display: flex;
    gap: 12px;
    margin-top: 32px;
    padding-top: 20px;
    border-top: 1px solid var(--border-light);
  }
}

@media (max-width: 768px) {
  .job-info-card {
    .job-header {
      flex-direction: column;
      gap: 16px;

      .job-actions {
        width: 100%;
        justify-content: flex-start;
      }
    }
  }

  .application-card {
    .applicant-info {
      flex-direction: column;
    }

    .application-actions {
      flex-wrap: wrap;
    }
  }
}
</style>
