<template>
  <div class="portal-job-detail">
    <div class="page-container" v-loading="loading">
      <div class="main-content">
        <!-- 职位信息 -->
        <div class="job-header-card">
          <div class="job-header">
            <div class="job-info">
              <h1>{{ job.title }}</h1>
              <div class="job-salary">{{ job.salary }}</div>
            </div>
            <el-button 
              v-if="hasApplied" 
              type="info" 
              size="large" 
              disabled
            >
              已投递
            </el-button>
            <el-button 
              v-else 
              type="primary" 
              size="large" 
              @click="handleApplyClick"
            >
              立即投递
            </el-button>
          </div>
          <div class="job-meta">
            <span><el-icon><Location /></el-icon>{{ job.location }}</span>
            <span><el-icon><Timer /></el-icon>{{ job.experience }}</span>
            <span><el-icon><School /></el-icon>{{ job.education }}</span>
            <span><el-icon><Calendar /></el-icon>{{ job.postTime }}</span>
          </div>
          <div class="job-tags">
            <el-tag v-for="tag in job.tags" :key="tag" type="info">{{ tag }}</el-tag>
          </div>
        </div>

        <!-- 职位描述 -->
        <div class="detail-card">
          <h2>职位描述</h2>
          <div class="description">{{ job.description }}</div>
        </div>

        <!-- 任职要求 -->
        <div class="detail-card">
          <h2>任职要求</h2>
          <ul class="requirements">
            <li v-for="(req, index) in job.requirements" :key="index">{{ req }}</li>
          </ul>
        </div>

        <!-- 技能要求 -->
        <div class="detail-card">
          <h2>技能要求</h2>
          <div class="skills">
            <el-tag v-for="skill in job.skills" :key="skill" size="large">{{ skill }}</el-tag>
          </div>
        </div>

        <!-- 福利待遇 -->
        <div class="detail-card">
          <h2>福利待遇</h2>
          <div class="benefits">
            <span v-for="benefit in job.benefits" :key="benefit" class="benefit-item">
              <el-icon><CircleCheck /></el-icon>{{ benefit }}
            </span>
          </div>
        </div>
      </div>

      <!-- 侧边栏 -->
      <div class="sidebar">
        <!-- 公司信息 -->
        <div class="company-card">
          <div class="company-header">
            <div class="company-logo">
              <el-icon :size="32"><OfficeBuilding /></el-icon>
            </div>
            <div class="company-info">
              <h3>{{ job.company }}</h3>
              <p>{{ job.companyType }} · {{ job.companySize }}</p>
            </div>
          </div>
          <div class="company-desc">
            {{ job.companyDesc }}
          </div>
          <el-button style="width: 100%">查看公司详情</el-button>
        </div>

        <!-- 工作地点 -->
        <div class="location-card">
          <h4>工作地点</h4>
          <p><el-icon><Location /></el-icon>{{ job.address }}</p>
        </div>

        <!-- 相似职位 -->
        <div class="similar-card">
          <h4>相似职位</h4>
          <div class="similar-job" v-for="sJob in similarJobs" :key="sJob.id" @click="goToJob(sJob.id)">
            <span class="similar-title">{{ sJob.title }}</span>
            <span class="similar-salary">{{ sJob.salary }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 投递弹窗 -->
    <el-dialog v-model="showApplyDialog" title="投递简历" width="500px">
      <div class="apply-content">
        <div class="apply-job-info">
          <h3>{{ job.title }}</h3>
          <p>{{ job.company }} · {{ job.location }}</p>
        </div>
        <el-form :model="applyForm" label-width="80px">
          <el-form-item label="简历" required>
            <el-select 
              v-model="applyForm.resumeId" 
              placeholder="选择简历" 
              style="width: 100%"
              :loading="loadingResumes"
            >
              <el-option 
                v-for="resume in userResumes" 
                :key="resume.id" 
                :label="resume.file_name || `简历 ${resume.id}`" 
                :value="resume.id" 
              />
            </el-select>
            <div class="upload-tip">
              <span v-if="userResumes.length === 0 && !loadingResumes" class="no-resume-tip">
                暂无简历，请先
                <el-button link type="primary" @click="goToUploadResume">上传简历</el-button>
              </span>
              <span v-else>
                或 <el-button link type="primary" @click="goToUploadResume">上传新简历</el-button>
              </span>
            </div>
          </el-form-item>
          <el-form-item label="求职信">
            <el-input 
              v-model="applyForm.coverLetter" 
              type="textarea" 
              :rows="4" 
              placeholder="简单介绍自己，为什么适合这个职位..." 
            />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="showApplyDialog = false">取消</el-button>
        <el-button 
          type="primary" 
          @click="submitApplication" 
          :loading="submitting"
          :disabled="!applyForm.resumeId"
        >
          确认投递
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Location, Timer, School, Calendar, OfficeBuilding, CircleCheck } from '@element-plus/icons-vue'
import request from '@/utils/request'
import { resumeApi } from '@/api/resume'
import { applicationApi } from '@/api/application'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const showApplyDialog = ref(false)
const submitting = ref(false)

// 用户简历列表
const userResumes = ref<any[]>([])
const loadingResumes = ref(false)

// 是否已投递
const hasApplied = ref(false)

const job = ref({
  id: 0,
  title: '',
  salary: '面议',
  location: '',
  experience: '',
  education: '本科',
  postTime: '',
  tags: [] as string[],
  description: '',
  requirements: [] as string[],
  skills: [] as string[],
  benefits: [] as string[],
  company: '',
  companyType: '互联网',
  companySize: '100-500人',
  companyDesc: '',
  address: ''
})

const similarJobs = ref([
  { id: 2, title: '前端开发工程师', salary: '20-35K' },
  { id: 3, title: '资深前端工程师', salary: '30-50K' },
  { id: 4, title: '前端架构师', salary: '40-60K' },
])

const applyForm = reactive({
  resumeId: null as number | null,
  coverLetter: ''
})

// 获取职位详情
const fetchJobDetail = async () => {
  const jobId = route.params.id
  if (!jobId) return
  
  loading.value = true
  try {
    const res = await request.get(`/jobs/${jobId}`)
    if (res.data?.code === 0 && res.data.data) {
      const jobData = res.data.data
      job.value = {
        id: jobData.id,
        title: jobData.title || '',
        salary: jobData.salary || '面议',
        location: jobData.location || '',
        experience: formatExperience(jobData.level),
        education: jobData.education || '不限',
        postTime: formatPostTime(jobData.created_at),
        tags: jobData.skills?.slice(0, 3) || [],
        description: jobData.description || '',
        requirements: parseRequirements(jobData.requirements),
        skills: jobData.skills || [],
        benefits: ['五险一金', '年终奖', '带薪年假', '弹性工作'],
        company: jobData.department || '公司',
        companyType: '互联网',
        companySize: '100-500人',
        companyDesc: '我们是一家专注于企业服务的科技公司，致力于用技术提升企业效率。',
        address: jobData.location || ''
      }
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

// 格式化经验要求
const formatExperience = (level: string) => {
  const map: Record<string, string> = {
    'junior': '1-3年',
    'mid': '3-5年',
    'senior': '5-10年',
    'expert': '10年以上',
    'management': '5年以上'
  }
  return map[level] || '不限'
}

// 格式化发布时间
const formatPostTime = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diffDays = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60 * 24))
  if (diffDays === 0) return '今天发布'
  if (diffDays === 1) return '昨天发布'
  if (diffDays < 7) return `${diffDays}天前发布`
  if (diffDays < 30) return `${Math.floor(diffDays / 7)}周前发布`
  return `${Math.floor(diffDays / 30)}月前发布`
}

// 解析职位要求
const parseRequirements = (requirements: string | string[] | undefined): string[] => {
  if (!requirements) return []
  if (Array.isArray(requirements)) return requirements
  // 如果是字符串，尝试按换行符分割
  return requirements.split('\n').filter(r => r.trim())
}

// 获取用户简历列表
const fetchUserResumes = async () => {
  loadingResumes.value = true
  try {
    const res = await resumeApi.list({ page: 1, page_size: 100 })
    if (res.data?.code === 0 && res.data.data) {
      userResumes.value = res.data.data.resumes || res.data.data || []
    }
  } catch (error) {
    console.error('获取简历列表失败:', error)
    userResumes.value = []
  } finally {
    loadingResumes.value = false
  }
}

// 检查是否已投递该职位
const checkIfApplied = async () => {
  const jobId = route.params.id
  if (!jobId) return
  
  try {
    const res = await applicationApi.getMyApplications({ page: 1, page_size: 1000 })
    if (res.data?.code === 0 && res.data.data) {
      const applications = res.data.data.applications || res.data.data || []
      hasApplied.value = applications.some((app: any) => app.job_id === Number(jobId))
    }
  } catch (error) {
    console.error('检查投递状态失败:', error)
  }
}

// 跳转到上传简历页面
const goToUploadResume = () => {
  showApplyDialog.value = false
  router.push('/portal/my-resume')
}

// 处理投递按钮点击
const handleApplyClick = async () => {
  // 检查是否已投递
  if (hasApplied.value) {
    ElMessage.warning('您已投递过该职位')
    return
  }
  
  // 获取用户简历列表
  await fetchUserResumes()
  
  // 检查是否有简历
  if (userResumes.value.length === 0) {
    ElMessage.warning('请先上传简历')
    showApplyDialog.value = true
    return
  }
  
  // 默认选择第一份简历
  if (userResumes.value.length > 0 && !applyForm.resumeId) {
    applyForm.resumeId = userResumes.value[0].id
  }
  
  showApplyDialog.value = true
}

const submitApplication = async () => {
  if (!applyForm.resumeId) {
    ElMessage.warning('请选择简历')
    return
  }
  
  if (!job.value.id) {
    ElMessage.error('职位信息错误')
    return
  }
  
  submitting.value = true
  try {
    const res = await applicationApi.create({
      job_id: job.value.id,
      talent_id: 0, // 后端会根据当前用户自动设置
      resume_id: applyForm.resumeId,
      cover_letter: applyForm.coverLetter || ''
    })
    
    if (res.data?.code === 0) {
      ElMessage.success('投递成功！')
      // 更新已投递状态
      hasApplied.value = true
      showApplyDialog.value = false
      // 重置表单
      applyForm.resumeId = null
      applyForm.coverLetter = ''
    } else {
      // 处理特定错误码
      const code = res.data?.code
      const message = res.data?.message
      
      if (code === 1001) {
        ElMessage.warning('您已投递过该职位')
        hasApplied.value = true
        showApplyDialog.value = false
      } else if (code === 1002) {
        ElMessage.warning('请先上传简历')
      } else {
        ElMessage.error(message || '投递失败，请稍后重试')
      }
    }
  } catch (error: any) {
    console.error('投递失败:', error)
    // 处理HTTP错误响应中的错误码
    const code = error.response?.data?.code
    const message = error.response?.data?.message
    
    if (code === 1001) {
      ElMessage.warning('您已投递过该职位')
      hasApplied.value = true
      showApplyDialog.value = false
    } else if (code === 1002) {
      ElMessage.warning('请先上传简历')
    } else {
      ElMessage.error(message || '投递失败，请稍后重试')
    }
  } finally {
    submitting.value = false
  }
}

const goToJob = (id: number) => {
  router.push(`/portal/jobs/${id}`)
}

onMounted(() => {
  fetchJobDetail()
  checkIfApplied()
})
</script>

<style scoped lang="scss">
.portal-job-detail {
  padding: 24px;
  background: #f8fafc;
  min-height: calc(100vh - 160px);

  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    display: grid;
    grid-template-columns: 1fr 320px;
    gap: 24px;
  }

  .job-header-card {
    background: white;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 20px;

    .job-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 16px;

      h1 {
        font-size: 24px;
        font-weight: 700;
        color: #1e293b;
        margin: 0 0 8px 0;
      }

      .job-salary {
        font-size: 28px;
        font-weight: 700;
        color: #0ea5e9;
      }
    }

    .job-meta {
      display: flex;
      gap: 24px;
      color: #64748b;
      margin-bottom: 16px;

      span {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }

    .job-tags {
      display: flex;
      gap: 8px;
    }
  }

  .detail-card {
    background: white;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 20px;

    h2 {
      font-size: 18px;
      font-weight: 600;
      color: #1e293b;
      margin: 0 0 16px 0;
      padding-bottom: 12px;
      border-bottom: 1px solid #f1f5f9;
    }

    .description {
      color: #475569;
      line-height: 1.8;

      :deep(p) { margin-bottom: 12px; }
    }

    .requirements {
      padding-left: 20px;
      color: #475569;

      li {
        margin-bottom: 8px;
        line-height: 1.6;
      }
    }

    .skills {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }

    .benefits {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;

      .benefit-item {
        display: flex;
        align-items: center;
        gap: 4px;
        color: #10b981;
      }
    }
  }

  .sidebar {
    .company-card, .location-card, .similar-card {
      background: white;
      border-radius: 12px;
      padding: 20px;
      margin-bottom: 16px;
    }

    .company-card {
      .company-header {
        display: flex;
        gap: 12px;
        margin-bottom: 16px;

        .company-logo {
          width: 56px;
          height: 56px;
          background: #f1f5f9;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #94a3b8;
        }

        .company-info {
          h3 {
            font-size: 16px;
            font-weight: 600;
            margin: 0 0 4px 0;
          }

          p {
            font-size: 13px;
            color: #64748b;
            margin: 0;
          }
        }
      }

      .company-desc {
        color: #64748b;
        font-size: 14px;
        line-height: 1.6;
        margin-bottom: 16px;
      }
    }

    .location-card {
      h4 {
        font-size: 14px;
        font-weight: 600;
        margin: 0 0 12px 0;
      }

      p {
        display: flex;
        align-items: center;
        gap: 4px;
        color: #64748b;
        margin: 0;
      }
    }

    .similar-card {
      h4 {
        font-size: 14px;
        font-weight: 600;
        margin: 0 0 12px 0;
      }

      .similar-job {
        display: flex;
        justify-content: space-between;
        padding: 10px 0;
        border-bottom: 1px solid #f1f5f9;
        cursor: pointer;

        &:hover .similar-title { color: #0ea5e9; }
        &:last-child { border-bottom: none; }

        .similar-title {
          color: #1e293b;
          font-size: 14px;
        }

        .similar-salary {
          color: #0ea5e9;
          font-size: 14px;
        }
      }
    }
  }
}

.apply-content {
  .apply-job-info {
    background: #f8fafc;
    padding: 16px;
    border-radius: 8px;
    margin-bottom: 20px;

    h3 { font-size: 16px; margin: 0 0 4px 0; }
    p { color: #64748b; font-size: 14px; margin: 0; }
  }

  .upload-tip {
    margin-top: 8px;
    font-size: 12px;
    color: #94a3b8;
    
    .no-resume-tip {
      color: #f59e0b;
    }
  }
}

@media (max-width: 900px) {
  .portal-job-detail .page-container {
    grid-template-columns: 1fr;
  }
}
</style>
