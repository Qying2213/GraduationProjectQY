<template>
  <div class="my-resume">
    <div class="page-container">
      <div class="page-header">
        <div class="header-left">
          <h1>我的简历</h1>
          <!-- 简历完整度指示器 -->
          <div class="completeness-indicator" v-if="!isEditing">
            <el-progress
              :percentage="completenessPercentage"
              :status="completenessPercentage >= 80 ? 'success' : completenessPercentage >= 50 ? '' : 'warning'"
              :stroke-width="8"
              style="width: 120px"
            />
            <span class="completeness-text">完整度 {{ completenessPercentage }}%</span>
          </div>
        </div>
        <div class="header-actions">
          <template v-if="isEditing">
            <el-button @click="cancelEdit">取消</el-button>
            <el-button type="primary" :loading="saving" @click="saveResume">
              <el-icon><Check /></el-icon> 保存简历
            </el-button>
          </template>
          <template v-else>
            <el-button type="primary" @click="startEdit">
              <el-icon><Edit /></el-icon> 编辑简历
            </el-button>
            <el-button @click="showUploadDialog = true">
              <el-icon><Upload /></el-icon> 上传简历
            </el-button>
          </template>
        </div>
      </div>

      <el-row :gutter="24" v-loading="loading">
        <!-- 在线简历 -->
        <el-col :xs="24" :lg="16">
          <div class="resume-card">
            <div class="card-header">
              <h2>在线简历</h2>
              <span v-if="resumeData.updated_at" class="update-time">
                最后更新: {{ formatDate(resumeData.updated_at) }}
              </span>
            </div>

            <!-- 编辑模式 -->
            <template v-if="isEditing">
              <el-tabs v-model="activeTab" class="edit-tabs">
                <el-tab-pane label="基本信息" name="basic">
                  <BasicInfoForm
                    ref="basicInfoFormRef"
                    v-model="editData.basic_info"
                  />
                </el-tab-pane>
                <el-tab-pane label="工作经历" name="work">
                  <WorkExperienceForm
                    ref="workExperienceFormRef"
                    v-model="editData.work_experience"
                  />
                </el-tab-pane>
                <el-tab-pane label="教育经历" name="education">
                  <EducationForm
                    ref="educationFormRef"
                    v-model="editData.education"
                  />
                </el-tab-pane>
                <el-tab-pane label="技能特长" name="skills">
                  <SkillsForm
                    ref="skillsFormRef"
                    v-model="editData.skills"
                  />
                </el-tab-pane>
              </el-tabs>
            </template>

            <!-- 查看模式 -->
            <template v-else>
              <!-- 基本信息 -->
              <div class="resume-section">
                <div class="section-header">
                  <el-icon><User /></el-icon>
                  <h3>基本信息</h3>
                </div>
                <div class="info-grid">
                  <div class="info-item">
                    <label>姓名</label>
                    <span>{{ resumeData.basic_info.name || '-' }}</span>
                  </div>
                  <div class="info-item">
                    <label>性别</label>
                    <span>{{ resumeData.basic_info.gender || '-' }}</span>
                  </div>
                  <div class="info-item">
                    <label>年龄</label>
                    <span>{{ resumeData.basic_info.age ? resumeData.basic_info.age + '岁' : '-' }}</span>
                  </div>
                  <div class="info-item">
                    <label>手机</label>
                    <span>{{ resumeData.basic_info.phone || '-' }}</span>
                  </div>
                  <div class="info-item">
                    <label>邮箱</label>
                    <span>{{ resumeData.basic_info.email || '-' }}</span>
                  </div>
                  <div class="info-item">
                    <label>现居地</label>
                    <span>{{ resumeData.basic_info.location || '-' }}</span>
                  </div>
                </div>
                <div v-if="resumeData.basic_info.summary" class="summary-section">
                  <label>个人简介</label>
                  <p>{{ resumeData.basic_info.summary }}</p>
                </div>
              </div>

              <!-- 工作经历 -->
              <div class="resume-section">
                <div class="section-header">
                  <el-icon><Suitcase /></el-icon>
                  <h3>工作经历</h3>
                </div>
                <div class="experience-item" v-for="(exp, index) in resumeData.work_experience" :key="index">
                  <div class="exp-header">
                    <div class="exp-title">
                      <h4>{{ exp.company }}</h4>
                      <span class="exp-position">{{ exp.position }}</span>
                    </div>
                    <span class="exp-time">{{ exp.start_date }} - {{ exp.is_current ? '至今' : exp.end_date }}</span>
                  </div>
                  <p class="exp-desc" v-if="exp.description">{{ exp.description }}</p>
                </div>
                <el-empty v-if="resumeData.work_experience.length === 0" description="暂无工作经历" :image-size="60" />
              </div>

              <!-- 教育经历 -->
              <div class="resume-section">
                <div class="section-header">
                  <el-icon><School /></el-icon>
                  <h3>教育经历</h3>
                </div>
                <div class="experience-item" v-for="(edu, index) in resumeData.education" :key="index">
                  <div class="exp-header">
                    <div class="exp-title">
                      <h4>{{ edu.school }}</h4>
                      <span class="exp-position">{{ edu.major }} · {{ edu.degree }}</span>
                    </div>
                    <span class="exp-time">{{ edu.start_date }} - {{ edu.is_current ? '至今' : edu.end_date }}</span>
                  </div>
                  <p class="exp-desc" v-if="edu.activities">{{ edu.activities }}</p>
                </div>
                <el-empty v-if="resumeData.education.length === 0" description="暂无教育经历" :image-size="60" />
              </div>

              <!-- 技能特长 -->
              <div class="resume-section">
                <div class="section-header">
                  <el-icon><Medal /></el-icon>
                  <h3>技能特长</h3>
                </div>
                <div class="skills-list" v-if="resumeData.skills.length > 0">
                  <el-tag v-for="skill in resumeData.skills" :key="skill" size="large">{{ skill }}</el-tag>
                </div>
                <el-empty v-else description="暂无技能标签" :image-size="60" />
              </div>
            </template>
          </div>
        </el-col>

        <!-- 附件简历 -->
        <el-col :xs="24" :lg="8">
          <div class="attachment-card">
            <h3>附件简历</h3>
            <div class="attachment-list">
              <div class="attachment-item" v-for="file in attachments" :key="file.id">
                <div class="file-icon">
                  <el-icon :size="24"><Document /></el-icon>
                </div>
                <div class="file-info">
                  <span class="file-name">{{ file.name }}</span>
                  <span class="file-meta">{{ file.size }} · {{ file.uploadTime }}</span>
                </div>
                <div class="file-actions">
                  <el-button link type="primary" size="small" @click="previewAttachment(file.id)">预览</el-button>
                  <el-button link type="danger" size="small" @click="deleteAttachment(file.id)">删除</el-button>
                </div>
              </div>
              <el-empty v-if="attachments.length === 0" description="暂无附件简历" :image-size="80" />
            </div>
          </div>

          <div class="tips-card">
            <h3>简历优化建议</h3>
            <ul>
              <li :class="{ completed: hasBasicInfo }">
                <el-icon><component :is="hasBasicInfo ? CircleCheck : Warning" /></el-icon>
                完善基本信息
              </li>
              <li :class="{ completed: hasWorkExperience }">
                <el-icon><component :is="hasWorkExperience ? CircleCheck : Warning" /></el-icon>
                添加工作经历
              </li>
              <li :class="{ completed: hasEducation }">
                <el-icon><component :is="hasEducation ? CircleCheck : Warning" /></el-icon>
                添加教育经历
              </li>
              <li :class="{ completed: hasSkills }">
                <el-icon><component :is="hasSkills ? CircleCheck : Warning" /></el-icon>
                添加技能标签
              </li>
              <li :class="{ completed: attachments.length > 0 }">
                <el-icon><component :is="attachments.length > 0 ? CircleCheck : Warning" /></el-icon>
                上传附件简历
              </li>
            </ul>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 上传弹窗 -->
    <el-dialog v-model="showUploadDialog" title="上传简历" width="500px">
      <el-upload
        drag
        :auto-upload="false"
        accept=".pdf,.doc,.docx"
        :limit="1"
        v-model:file-list="fileList"
      >
        <el-icon class="el-icon--upload" :size="48"><UploadFilled /></el-icon>
        <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip">支持 PDF、DOC、DOCX 格式，文件大小不超过 10MB</div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button type="primary" :loading="uploading" :disabled="fileList.length === 0" @click="handleUpload">上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Upload, User, Suitcase, School, Medal, Document,
  CircleCheck, Warning, UploadFilled, Edit, Check
} from '@element-plus/icons-vue'
import { BasicInfoForm, WorkExperienceForm, EducationForm, SkillsForm } from '@/components/resume'
import type { BasicInfo, WorkExperience, Education } from '@/components/resume'
import { resumeApi } from '@/api/resume'
import type { OnlineResumeRequest } from '@/api/resume'
import { useUserStore } from '@/store/user'
import request from '@/utils/request'

const userStore = useUserStore()

// 状态
const loading = ref(false)
const saving = ref(false)
const uploading = ref(false)
const isEditing = ref(false)
const showUploadDialog = ref(false)
const fileList = ref<any[]>([])
const activeTab = ref('basic')

// 表单组件引用
const basicInfoFormRef = ref<InstanceType<typeof BasicInfoForm> | null>(null)
const workExperienceFormRef = ref<InstanceType<typeof WorkExperienceForm> | null>(null)
const educationFormRef = ref<InstanceType<typeof EducationForm> | null>(null)
const skillsFormRef = ref<InstanceType<typeof SkillsForm> | null>(null)

// 简历数据（查看模式）
interface ResumeData {
  id: number
  basic_info: BasicInfo
  work_experience: WorkExperience[]
  education: Education[]
  skills: string[]
  is_complete: boolean
  updated_at: string
}

const resumeData = reactive<ResumeData>({
  id: 0,
  basic_info: {
    name: '',
    phone: '',
    email: '',
    location: '',
    avatar: '',
    gender: '',
    age: undefined,
    summary: ''
  },
  work_experience: [],
  education: [],
  skills: [],
  is_complete: false,
  updated_at: ''
})

// 编辑数据（编辑模式）
const editData = reactive<OnlineResumeRequest>({
  basic_info: {
    name: '',
    phone: '',
    email: '',
    location: '',
    avatar: '',
    gender: '',
    age: undefined,
    summary: ''
  },
  work_experience: [],
  education: [],
  skills: []
})

// 附件简历
interface Attachment {
  id: number
  name: string
  size: string
  uploadTime: string
}
const attachments = ref<Attachment[]>([])

// 计算属性：简历完整度
const hasBasicInfo = computed(() => {
  const info = resumeData.basic_info
  return !!(info.name && info.phone && info.email)
})

const hasWorkExperience = computed(() => resumeData.work_experience.length > 0)
const hasEducation = computed(() => resumeData.education.length > 0)
const hasSkills = computed(() => resumeData.skills.length > 0)

const completenessPercentage = computed(() => {
  let score = 0
  // 基本信息（必填项）- 40分
  const info = resumeData.basic_info
  if (info.name) score += 10
  if (info.phone) score += 10
  if (info.email) score += 10
  if (info.location) score += 5
  if (info.summary) score += 5
  // 工作经历 - 25分
  if (resumeData.work_experience.length > 0) score += 25
  // 教育经历 - 20分
  if (resumeData.education.length > 0) score += 20
  // 技能 - 15分
  if (resumeData.skills.length > 0) score += 15
  return Math.min(score, 100)
})

// 加载在线简历
const loadOnlineResume = async () => {
  loading.value = true
  try {
    const res = await resumeApi.getOnlineResume()
    if (res.data?.code === 0 && res.data.data) {
      const data = res.data.data
      resumeData.id = data.id || 0
      resumeData.basic_info = data.basic_info || {
        name: '',
        phone: '',
        email: '',
        location: '',
        avatar: '',
        gender: '',
        age: undefined,
        summary: ''
      }
      resumeData.work_experience = data.work_experience || []
      resumeData.education = data.education || []
      resumeData.skills = data.skills || []
      resumeData.is_complete = data.is_complete || false
      resumeData.updated_at = data.updated_at || ''
    } else {
      // 如果没有在线简历，从用户信息初始化
      const user = userStore.userInfo
      if (user) {
        resumeData.basic_info.name = user.real_name || user.username || ''
        resumeData.basic_info.email = user.email || ''
        resumeData.basic_info.phone = user.phone || ''
      }
    }
  } catch (error: any) {
    console.error('加载在线简历失败:', error)
    // 如果接口不存在，尝试从旧接口获取数据
    await loadResumeFromLegacyApi()
  } finally {
    loading.value = false
  }
}

// 从旧接口加载简历数据（兼容）
const loadResumeFromLegacyApi = async () => {
  try {
    const user = userStore.userInfo
    if (user) {
      resumeData.basic_info.name = user.real_name || user.username || ''
      resumeData.basic_info.email = user.email || ''
      resumeData.basic_info.phone = user.phone || ''
    }

    const res = await request.get('/resumes', { params: { page_size: 1 } })
    if (res.data?.code === 0 && res.data.data?.resumes?.length > 0) {
      const data = res.data.data.resumes[0]
      if (data.parsed_data) {
        try {
          const parsed = typeof data.parsed_data === 'string' ? JSON.parse(data.parsed_data) : data.parsed_data
          resumeData.basic_info = {
            name: parsed.name || resumeData.basic_info.name,
            phone: parsed.phone || resumeData.basic_info.phone,
            email: parsed.email || resumeData.basic_info.email,
            location: parsed.location || '',
            gender: parsed.gender || '',
            age: parsed.age || undefined,
            summary: parsed.summary || ''
          }
          // 转换工作经历格式
          if (parsed.workExperience && Array.isArray(parsed.workExperience)) {
            resumeData.work_experience = parsed.workExperience.map((exp: any) => ({
              company: exp.company || '',
              position: exp.position || '',
              start_date: exp.startTime || exp.start_date || '',
              end_date: exp.endTime || exp.end_date || '',
              is_current: false,
              description: exp.description || '',
              location: exp.location || ''
            }))
          }
          // 转换教育经历格式
          if (parsed.education && Array.isArray(parsed.education)) {
            resumeData.education = parsed.education.map((edu: any) => ({
              school: edu.school || '',
              degree: edu.degree || '',
              major: edu.major || '',
              start_date: edu.startTime || edu.start_date || '',
              end_date: edu.endTime || edu.end_date || '',
              is_current: false,
              gpa: edu.gpa || '',
              activities: edu.activities || ''
            }))
          }
          resumeData.skills = parsed.skills || []
        } catch (e) {
          console.error('解析简历数据失败', e)
        }
      }
    }
  } catch (error) {
    console.error('从旧接口加载简历失败:', error)
  }
}

// 加载附件简历列表
const loadAttachments = async () => {
  try {
    const res = await request.get('/resumes', { params: { page_size: 10 } })
    if (res.data?.code === 0) {
      const resumes = res.data.data?.resumes || []
      attachments.value = resumes.map((r: any) => ({
        id: r.id,
        name: r.file_name,
        size: formatFileSize(r.file_size),
        uploadTime: r.created_at?.split('T')[0] || ''
      }))
    }
  } catch (error) {
    console.error('加载附件列表失败:', error)
  }
}

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (!bytes) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 开始编辑
const startEdit = () => {
  // 复制当前数据到编辑数据
  editData.basic_info = { ...resumeData.basic_info }
  editData.work_experience = resumeData.work_experience.map(exp => ({ ...exp }))
  editData.education = resumeData.education.map(edu => ({ ...edu }))
  editData.skills = [...resumeData.skills]
  isEditing.value = true
  activeTab.value = 'basic'
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
}

// 验证所有表单
const validateAllForms = async (): Promise<boolean> => {
  const results = await Promise.all([
    basicInfoFormRef.value?.validate() ?? true,
    workExperienceFormRef.value?.validate() ?? true,
    educationFormRef.value?.validate() ?? true,
    skillsFormRef.value?.validate() ?? true
  ])
  return results.every(result => result)
}

// 保存简历
const saveResume = async () => {
  // 验证所有表单
  const isValid = await validateAllForms()
  if (!isValid) {
    ElMessage.warning('请检查表单中的错误信息')
    return
  }

  // 验证必填字段
  if (!editData.basic_info.name || !editData.basic_info.phone || !editData.basic_info.email) {
    ElMessage.warning('请填写必填信息：姓名、手机号、邮箱')
    activeTab.value = 'basic'
    return
  }

  saving.value = true
  try {
    const res = await resumeApi.saveOnlineResume(editData)
    if (res.data?.code === 0) {
      ElMessage.success('简历保存成功')
      // 更新显示数据
      if (res.data.data) {
        const data = res.data.data
        resumeData.id = data.id || resumeData.id
        resumeData.basic_info = data.basic_info || editData.basic_info
        resumeData.work_experience = data.work_experience || editData.work_experience
        resumeData.education = data.education || editData.education
        resumeData.skills = data.skills || editData.skills
        resumeData.is_complete = data.is_complete || false
        resumeData.updated_at = data.updated_at || new Date().toISOString()
      } else {
        // 如果响应没有返回数据，使用编辑数据
        resumeData.basic_info = { ...editData.basic_info }
        resumeData.work_experience = [...editData.work_experience]
        resumeData.education = [...editData.education]
        resumeData.skills = [...editData.skills]
        resumeData.updated_at = new Date().toISOString()
      }
      isEditing.value = false
    } else {
      ElMessage.error(res.data?.message || '保存失败')
    }
  } catch (error: any) {
    console.error('保存简历失败:', error)
    ElMessage.error(error.response?.data?.message || '保存失败，请稍后重试')
  } finally {
    saving.value = false
  }
}

// 上传附件简历
const handleUpload = async () => {
  if (fileList.value.length === 0) {
    ElMessage.warning('请选择文件')
    return
  }

  uploading.value = true
  try {
    const file = fileList.value[0]
    const formData = new FormData()
    formData.append('file', file.raw)

    const res = await request.post('/resumes/upload', formData)
    if (res.data?.code === 0) {
      ElMessage.success('上传成功')
      showUploadDialog.value = false
      fileList.value = []
      loadAttachments()
    } else {
      ElMessage.error(res.data?.message || '上传失败')
    }
  } catch (error) {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

// 删除附件
const deleteAttachment = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这份简历吗？', '删除确认', { type: 'warning' })
    const res = await request.delete(`/resumes/${id}`)
    if (res.data?.code === 0) {
      attachments.value = attachments.value.filter(a => a.id !== id)
      ElMessage.success('已删除')
    } else {
      ElMessage.error(res.data?.message || '删除失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 预览附件
const previewAttachment = async (id: number) => {
  try {
    const res = await request.get(`/resumes/${id}`)
    if (res.data?.code === 0 && res.data.data?.file_url) {
      window.open(res.data.data.file_url, '_blank')
    } else {
      ElMessage.info('暂无预览')
    }
  } catch {
    ElMessage.error('获取预览失败')
  }
}

onMounted(() => {
  loadOnlineResume()
  loadAttachments()
})
</script>


<style scoped lang="scss">
.my-resume {
  padding: 24px;
  background: #f8fafc;
  min-height: calc(100vh - 160px);

  .page-container {
    max-width: 1200px;
    margin: 0 auto;
  }

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    flex-wrap: wrap;
    gap: 16px;

    .header-left {
      display: flex;
      align-items: center;
      gap: 24px;

      h1 {
        font-size: 24px;
        font-weight: 700;
        color: #1e293b;
        margin: 0;
      }

      .completeness-indicator {
        display: flex;
        align-items: center;
        gap: 12px;

        .completeness-text {
          font-size: 13px;
          color: #64748b;
          white-space: nowrap;
        }
      }
    }

    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .resume-card {
    background: white;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 24px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 24px;
      padding-bottom: 16px;
      border-bottom: 1px solid #f1f5f9;

      h2 {
        font-size: 20px;
        font-weight: 600;
        margin: 0;
      }

      .update-time {
        font-size: 13px;
        color: #94a3b8;
      }
    }

    .edit-tabs {
      :deep(.el-tabs__header) {
        margin-bottom: 24px;
      }

      :deep(.el-tabs__item) {
        font-size: 15px;
      }
    }

    .resume-section {
      margin-bottom: 32px;

      &:last-child { margin-bottom: 0; }

      .section-header {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 16px;

        .el-icon { color: #0ea5e9; }

        h3 {
          font-size: 16px;
          font-weight: 600;
          margin: 0;
        }
      }

      .info-grid {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 16px;

        .info-item {
          label {
            display: block;
            font-size: 12px;
            color: #94a3b8;
            margin-bottom: 4px;
          }

          span {
            color: #1e293b;
          }
        }
      }

      .summary-section {
        margin-top: 16px;
        padding-top: 16px;
        border-top: 1px dashed #e2e8f0;

        label {
          display: block;
          font-size: 12px;
          color: #94a3b8;
          margin-bottom: 8px;
        }

        p {
          color: #475569;
          line-height: 1.6;
          margin: 0;
        }
      }

      .experience-item {
        padding: 16px;
        background: #f8fafc;
        border-radius: 8px;
        margin-bottom: 12px;

        &:last-child { margin-bottom: 0; }

        .exp-header {
          display: flex;
          justify-content: space-between;
          margin-bottom: 8px;

          h4 {
            font-size: 15px;
            font-weight: 600;
            margin: 0 0 4px 0;
          }

          .exp-position {
            color: #64748b;
            font-size: 14px;
          }

          .exp-time {
            color: #94a3b8;
            font-size: 13px;
            white-space: nowrap;
          }
        }

        .exp-desc {
          color: #475569;
          font-size: 14px;
          line-height: 1.6;
          margin: 0;
        }
      }

      .skills-list {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
      }
    }
  }

  .attachment-card, .tips-card {
    background: white;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 16px;

    h3 {
      font-size: 16px;
      font-weight: 600;
      margin: 0 0 16px 0;
    }
  }

  .attachment-list {
    .attachment-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px;
      background: #f8fafc;
      border-radius: 8px;
      margin-bottom: 8px;

      .file-icon {
        width: 40px;
        height: 40px;
        background: #fee2e2;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #ef4444;
      }

      .file-info {
        flex: 1;

        .file-name {
          display: block;
          font-weight: 500;
          color: #1e293b;
        }

        .file-meta {
          font-size: 12px;
          color: #94a3b8;
        }
      }
    }
  }

  .tips-card {
    ul {
      list-style: none;
      padding: 0;
      margin: 0;

      li {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px 0;
        color: #475569;

        .el-icon {
          color: #f59e0b;
        }

        &.completed .el-icon {
          color: #10b981;
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .my-resume {
    .page-header {
      .header-left {
        flex-direction: column;
        align-items: flex-start;
        gap: 12px;
      }
    }

    .resume-card .resume-section .info-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }
}
</style>
