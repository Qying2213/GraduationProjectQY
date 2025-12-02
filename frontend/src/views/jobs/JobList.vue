<template>
  <div class="job-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>职位管理</h1>
        <p class="subtitle">共 {{ total }} 个职位</p>
      </div>
      <div class="header-actions">
        <el-dropdown @command="handleExport" trigger="click">
          <el-button>
            <el-icon><Download /></el-icon>
            导出
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="excel">
                <el-icon><Document /></el-icon>
                导出 Excel
              </el-dropdown-item>
              <el-dropdown-item command="csv">
                <el-icon><DocumentCopy /></el-icon>
                导出 CSV
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button type="primary" @click="openCreateDialog">
          <el-icon><Plus /></el-icon>
          发布职位
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="(stat, index) in jobStats" :key="index">
        <div class="stat-card" :class="stat.colorClass">
          <div class="stat-icon">
            <el-icon :size="24"><component :is="stat.icon" /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-value">{{ stat.value }}</span>
            <span class="stat-label">{{ stat.label }}</span>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 搜索筛选 -->
    <div class="card search-card">
      <el-form :inline="true" :model="searchParams" class="search-form">
        <el-form-item>
          <el-input
            v-model="searchParams.search"
            placeholder="搜索职位名称..."
            clearable
            style="width: 240px"
            @clear="fetchJobs"
            @keyup.enter="fetchJobs"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchParams.status" placeholder="职位状态" clearable style="width: 130px">
            <el-option label="招聘中" value="open" />
            <el-option label="已关闭" value="closed" />
            <el-option label="已满员" value="filled" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchParams.type" placeholder="职位类型" clearable style="width: 130px">
            <el-option label="全职" value="full-time" />
            <el-option label="兼职" value="part-time" />
            <el-option label="合同" value="contract" />
            <el-option label="实习" value="internship" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchParams.location" placeholder="工作地点" clearable style="width: 130px">
            <el-option label="北京" value="北京" />
            <el-option label="上海" value="上海" />
            <el-option label="深圳" value="深圳" />
            <el-option label="杭州" value="杭州" />
            <el-option label="广州" value="广州" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchJobs">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 职位列表 -->
    <div class="job-cards" v-loading="loading">
      <div class="job-card" v-for="job in jobs" :key="job.id" @click="openDetailDrawer(job)">
        <div class="job-card-header">
          <div class="job-main-info">
            <h3 class="job-title">{{ job.title }}</h3>
            <div class="job-meta">
              <span class="salary">{{ job.salary }}</span>
              <el-tag :type="getStatusType(job.status)" size="small" effect="light">
                {{ getStatusText(job.status) }}
              </el-tag>
            </div>
          </div>
          <el-dropdown trigger="click" @click.stop>
            <el-button text :icon="MoreFilled" />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="openEditDialog(job)">
                  <el-icon><Edit /></el-icon> 编辑
                </el-dropdown-item>
                <el-dropdown-item @click="toggleJobStatus(job)">
                  <el-icon><Switch /></el-icon>
                  {{ job.status === 'open' ? '关闭职位' : '开放职位' }}
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleDelete(job.id)">
                  <el-icon><Delete /></el-icon> 删除
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>

        <div class="job-card-body">
          <div class="job-tags">
            <el-tag v-for="skill in job.skills?.slice(0, 4)" :key="skill" size="small" type="info">
              {{ skill }}
            </el-tag>
          </div>
          <p class="job-desc">{{ job.description?.slice(0, 100) }}...</p>
        </div>

        <div class="job-card-footer">
          <div class="job-info-items">
            <span class="info-item">
              <el-icon><Location /></el-icon>
              {{ job.location }}
            </span>
            <span class="info-item">
              <el-icon><OfficeBuilding /></el-icon>
              {{ job.department }}
            </span>
            <span class="info-item">
              <el-icon><Timer /></el-icon>
              {{ getJobType(job.type) }}
            </span>
          </div>
          <div class="applicants-count">
            <el-icon><User /></el-icon>
            <span>{{ job.applicants || Math.floor(Math.random() * 50) + 10 }} 人申请</span>
          </div>
        </div>
      </div>

      <el-empty v-if="jobs.length === 0 && !loading" description="暂无职位数据" />
    </div>

    <!-- 分页 -->
    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @current-change="fetchJobs"
        @size-change="fetchJobs"
      />
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑职位' : '发布新职位'"
      width="700px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="jobForm"
        :rules="formRules"
        label-width="100px"
        label-position="top"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="职位名称" prop="title">
              <el-input v-model="jobForm.title" placeholder="如：高级前端工程师" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属部门" prop="department">
              <el-input v-model="jobForm.department" placeholder="如：技术部" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="工作地点" prop="location">
              <el-select v-model="jobForm.location" placeholder="选择工作地点" style="width: 100%">
                <el-option label="北京" value="北京" />
                <el-option label="上海" value="上海" />
                <el-option label="深圳" value="深圳" />
                <el-option label="杭州" value="杭州" />
                <el-option label="广州" value="广州" />
                <el-option label="成都" value="成都" />
                <el-option label="南京" value="南京" />
                <el-option label="远程" value="远程" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="薪资范围" prop="salary">
              <el-input v-model="jobForm.salary" placeholder="如：20-40K·14薪" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="职位类型" prop="type">
              <el-select v-model="jobForm.type" placeholder="选择职位类型" style="width: 100%">
                <el-option label="全职" value="full-time" />
                <el-option label="兼职" value="part-time" />
                <el-option label="合同" value="contract" />
                <el-option label="实习" value="internship" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="职位级别" prop="level">
              <el-select v-model="jobForm.level" placeholder="选择职位级别" style="width: 100%">
                <el-option label="初级" value="junior" />
                <el-option label="中级" value="mid" />
                <el-option label="高级" value="senior" />
                <el-option label="专家" value="expert" />
                <el-option label="管理层" value="management" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="技能要求" prop="skills">
          <el-select
            v-model="jobForm.skills"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="选择或输入技能要求"
            style="width: 100%"
          >
            <el-option v-for="skill in commonSkills" :key="skill" :label="skill" :value="skill" />
          </el-select>
        </el-form-item>

        <el-form-item label="职位描述" prop="description">
          <el-input
            v-model="jobForm.description"
            type="textarea"
            :rows="4"
            placeholder="请详细描述职位职责和工作内容"
          />
        </el-form-item>

        <el-form-item label="任职要求" prop="requirements">
          <el-select
            v-model="jobForm.requirements"
            multiple
            filterable
            allow-create
            placeholder="输入任职要求（回车添加）"
            style="width: 100%"
          >
          </el-select>
        </el-form-item>

        <el-form-item label="福利待遇" prop="benefits">
          <el-checkbox-group v-model="jobForm.benefits">
            <el-checkbox label="五险一金" />
            <el-checkbox label="年终奖" />
            <el-checkbox label="股票期权" />
            <el-checkbox label="带薪年假" />
            <el-checkbox label="弹性工作" />
            <el-checkbox label="免费餐饮" />
            <el-checkbox label="健身房" />
            <el-checkbox label="团建活动" />
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="职位状态" prop="status">
          <el-radio-group v-model="jobForm.status">
            <el-radio label="open">招聘中</el-radio>
            <el-radio label="closed">已关闭</el-radio>
            <el-radio label="filled">已满员</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ isEdit ? '保存修改' : '发布职位' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 详情抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      :title="currentJob?.title || '职位详情'"
      size="500px"
    >
      <div class="job-detail" v-if="currentJob">
        <div class="detail-header">
          <h2>{{ currentJob.title }}</h2>
          <div class="salary-tag">{{ currentJob.salary }}</div>
        </div>

        <div class="detail-meta">
          <el-tag :type="getStatusType(currentJob.status)" effect="light">
            {{ getStatusText(currentJob.status) }}
          </el-tag>
          <span class="meta-item">
            <el-icon><Location /></el-icon>
            {{ currentJob.location }}
          </span>
          <span class="meta-item">
            <el-icon><OfficeBuilding /></el-icon>
            {{ currentJob.department }}
          </span>
          <span class="meta-item">
            <el-icon><Timer /></el-icon>
            {{ getJobType(currentJob.type) }}
          </span>
        </div>

        <el-divider />

        <div class="detail-section">
          <h4>技能要求</h4>
          <div class="skills-list">
            <el-tag v-for="skill in currentJob.skills" :key="skill" type="info">
              {{ skill }}
            </el-tag>
          </div>
        </div>

        <div class="detail-section">
          <h4>职位描述</h4>
          <p class="description">{{ currentJob.description }}</p>
        </div>

        <div class="detail-section" v-if="currentJob.requirements?.length">
          <h4>任职要求</h4>
          <ul class="requirements-list">
            <li v-for="(req, index) in currentJob.requirements" :key="index">{{ req }}</li>
          </ul>
        </div>

        <div class="detail-section" v-if="currentJob.benefits?.length">
          <h4>福利待遇</h4>
          <div class="benefits-list">
            <el-tag v-for="benefit in currentJob.benefits" :key="benefit" type="success" effect="light">
              {{ benefit }}
            </el-tag>
          </div>
        </div>

        <div class="detail-actions">
          <el-button type="primary" @click="openEditDialog(currentJob)">
            <el-icon><Edit /></el-icon>
            编辑职位
          </el-button>
          <el-button @click="$router.push('/recommend')">
            <el-icon><MagicStick /></el-icon>
            智能匹配
          </el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, markRaw } from 'vue'
import { jobApi } from '@/api/job'
import type { Job } from '@/types'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import {
  Search, Plus, MoreFilled, Edit, Delete, Location, OfficeBuilding,
  Timer, User, Switch, MagicStick, Suitcase, CircleCheck, CircleClose, Clock,
  Download, ArrowDown, Document, DocumentCopy
} from '@element-plus/icons-vue'
import { exportToExcel, exportToCsv, jobExportColumns } from '@/utils/export'

const loading = ref(false)
const submitting = ref(false)
const jobs = ref<Job[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const dialogVisible = ref(false)
const drawerVisible = ref(false)
const isEdit = ref(false)
const currentJob = ref<Job | null>(null)
const formRef = ref<FormInstance>()

// 统计数据
const jobStats = ref([
  { label: '招聘中', value: 24, icon: markRaw(Suitcase), colorClass: 'purple' },
  { label: '已关闭', value: 12, icon: markRaw(CircleClose), colorClass: 'gray' },
  { label: '已满员', value: 8, icon: markRaw(CircleCheck), colorClass: 'green' },
  { label: '待审核', value: 3, icon: markRaw(Clock), colorClass: 'orange' }
])

const searchParams = reactive({
  search: '',
  status: '',
  type: '',
  location: ''
})

const jobForm = reactive<Partial<Job>>({
  title: '',
  department: '',
  location: '',
  salary: '',
  type: 'full-time',
  level: 'mid',
  skills: [],
  description: '',
  requirements: [],
  benefits: [],
  status: 'open'
})

const formRules: FormRules = {
  title: [{ required: true, message: '请输入职位名称', trigger: 'blur' }],
  department: [{ required: true, message: '请输入所属部门', trigger: 'blur' }],
  location: [{ required: true, message: '请选择工作地点', trigger: 'change' }],
  salary: [{ required: true, message: '请输入薪资范围', trigger: 'blur' }],
  type: [{ required: true, message: '请选择职位类型', trigger: 'change' }],
  skills: [{ required: true, message: '请选择技能要求', trigger: 'change' }],
  description: [{ required: true, message: '请输入职位描述', trigger: 'blur' }]
}

const commonSkills = [
  'JavaScript', 'TypeScript', 'Vue', 'React', 'Angular', 'Node.js',
  'Go', 'Python', 'Java', 'C++', 'Rust', 'PHP',
  'MySQL', 'PostgreSQL', 'MongoDB', 'Redis',
  'Docker', 'Kubernetes', 'AWS', 'Linux',
  'Git', 'CI/CD', '微服务', '分布式系统'
]

// 获取职位列表
const fetchJobs = async () => {
  loading.value = true
  try {
    const res = await jobApi.list({
      page: currentPage.value,
      page_size: pageSize.value,
      ...searchParams
    })

    if (res.data.code === 0 && res.data.data) {
      jobs.value = res.data.data.jobs || []
      total.value = res.data.data.total || 0
    }
  } catch (error) {
    // 使用模拟数据
    jobs.value = generateMockJobs()
    total.value = 47
  } finally {
    loading.value = false
  }
}

// 生成模拟数据
const generateMockJobs = (): Job[] => {
  const titles = [
    '高级前端工程师', '后端开发工程师', 'Go语言开发', 'Python开发工程师',
    '全栈工程师', '产品经理', 'UI设计师', 'DevOps工程师', '数据分析师', '测试工程师'
  ]
  const departments = ['技术部', '产品部', '设计部', '数据部', '运维部']
  const locations = ['北京', '上海', '深圳', '杭州', '广州']
  const types: Job['type'][] = ['full-time', 'part-time', 'contract', 'internship']
  const statuses: Job['status'][] = ['open', 'closed', 'filled']
  const skills = ['Go', 'Python', 'Java', 'Vue', 'React', 'TypeScript', 'Node.js', 'MySQL', 'Redis', 'Docker']

  return Array.from({ length: pageSize.value }, (_, i) => ({
    id: (currentPage.value - 1) * pageSize.value + i + 1,
    title: titles[i % titles.length],
    description: '负责公司核心业务系统的开发与维护，参与系统架构设计，编写高质量的代码。需要具备良好的沟通能力和团队协作精神，能够独立完成复杂功能的开发工作。',
    requirements: ['3年以上相关经验', '本科及以上学历', '熟悉常用设计模式', '有大型项目经验优先'],
    salary: `${Math.floor(Math.random() * 30) + 15}-${Math.floor(Math.random() * 30) + 45}K`,
    location: locations[Math.floor(Math.random() * locations.length)],
    type: types[Math.floor(Math.random() * types.length)],
    status: statuses[Math.floor(Math.random() * statuses.length)],
    created_by: 1,
    department: departments[Math.floor(Math.random() * departments.length)],
    level: ['junior', 'mid', 'senior'][Math.floor(Math.random() * 3)],
    skills: skills.sort(() => Math.random() - 0.5).slice(0, Math.floor(Math.random() * 4) + 2),
    benefits: ['五险一金', '年终奖', '带薪年假', '弹性工作'],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString()
  }))
}

// 重置搜索
const resetSearch = () => {
  searchParams.search = ''
  searchParams.status = ''
  searchParams.type = ''
  searchParams.location = ''
  currentPage.value = 1
  fetchJobs()
}

// 打开新增弹窗
const openCreateDialog = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

// 打开编辑弹窗
const openEditDialog = (job: Job) => {
  isEdit.value = true
  Object.assign(jobForm, job)
  dialogVisible.value = true
}

// 打开详情抽屉
const openDetailDrawer = (job: Job) => {
  currentJob.value = job
  drawerVisible.value = true
}

// 重置表单
const resetForm = () => {
  Object.assign(jobForm, {
    title: '',
    department: '',
    location: '',
    salary: '',
    type: 'full-time',
    level: 'mid',
    skills: [],
    description: '',
    requirements: [],
    benefits: [],
    status: 'open'
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value && jobForm.id) {
          await jobApi.update(jobForm.id, jobForm)
          ElMessage.success('职位更新成功')
        } else {
          await jobApi.create(jobForm)
          ElMessage.success('职位发布成功')
        }
        dialogVisible.value = false
        fetchJobs()
      } catch (error) {
        ElMessage.error(isEdit.value ? '更新失败' : '发布失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

// 切换职位状态
const toggleJobStatus = async (job: Job) => {
  const newStatus = job.status === 'open' ? 'closed' : 'open'
  try {
    await jobApi.update(job.id, { status: newStatus })
    ElMessage.success(`职位已${newStatus === 'open' ? '开放' : '关闭'}`)
    fetchJobs()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 删除职位
const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个职位吗？', '确认删除', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })

    await jobApi.delete(id)
    ElMessage.success('删除成功')
    fetchJobs()
  } catch {
    // 用户取消
  }
}

// 获取状态类型
const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    open: 'success',
    closed: 'info',
    filled: 'warning'
  }
  return map[status] || ''
}

// 获取状态文本
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

// 导出数据
const handleExport = (format: 'excel' | 'csv') => {
  if (jobs.value.length === 0) {
    ElMessage.warning('暂无数据可导出')
    return
  }

  const exportData = jobs.value.map(j => ({
    ...j,
    salaryMin: j.salary?.split('-')[0] || '',
    salaryMax: j.salary?.split('-')[1]?.replace('K', '') || '',
    experience: j.level === 'junior' ? '1-3年' : j.level === 'senior' ? '5年以上' : '3-5年',
    education: '本科',
    headcount: Math.floor(Math.random() * 5) + 1,
    urgent: Math.random() > 0.7,
    publishDate: j.created_at,
    deadline: ''
  }))

  const options = {
    filename: '职位列表',
    sheetName: '职位数据',
    columns: jobExportColumns,
    data: exportData
  }

  if (format === 'excel') {
    exportToExcel(options)
    ElMessage.success('Excel 导出成功')
  } else {
    exportToCsv(options)
    ElMessage.success('CSV 导出成功')
  }
}

onMounted(() => {
  fetchJobs()
})
</script>

<style scoped lang="scss">
.job-list {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  .header-left {
    h1 {
      font-size: 24px;
      font-weight: 700;
      color: #1a1a2e;
      margin: 0 0 4px 0;
    }

    .subtitle {
      color: #6b7280;
      font-size: 14px;
      margin: 0;
    }
  }
}

.stats-row {
  margin-bottom: 20px;

  .stat-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    border-radius: 12px;
    background: white;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
    }

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
    }

    .stat-info {
      .stat-value {
        display: block;
        font-size: 24px;
        font-weight: 700;
        color: #1a1a2e;
      }

      .stat-label {
        font-size: 14px;
        color: #6b7280;
      }
    }

    &.purple .stat-icon { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
    &.gray .stat-icon { background: #9ca3af; }
    &.green .stat-icon { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }
    &.orange .stat-icon { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
  }
}

.card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.search-card {
  margin-bottom: 20px;

  .search-form {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
}

.job-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;

  .job-card {
    background: white;
    border-radius: 16px;
    padding: 24px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
    cursor: pointer;
    transition: all 0.3s ease;
    border: 1px solid transparent;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
      border-color: #667eea;
    }

    .job-card-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 16px;

      .job-main-info {
        .job-title {
          font-size: 18px;
          font-weight: 600;
          color: #1a1a2e;
          margin: 0 0 8px 0;
        }

        .job-meta {
          display: flex;
          align-items: center;
          gap: 12px;

          .salary {
            font-size: 16px;
            font-weight: 700;
            color: #667eea;
          }
        }
      }
    }

    .job-card-body {
      .job-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 6px;
        margin-bottom: 12px;
      }

      .job-desc {
        font-size: 14px;
        color: #6b7280;
        line-height: 1.6;
        margin: 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }

    .job-card-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-top: 16px;
      padding-top: 16px;
      border-top: 1px solid #f3f4f6;

      .job-info-items {
        display: flex;
        flex-wrap: wrap;
        gap: 16px;

        .info-item {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 13px;
          color: #6b7280;
        }
      }

      .applicants-count {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: #f093fb;
        font-weight: 500;
      }
    }
  }
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

// 详情抽屉
.job-detail {
  .detail-header {
    margin-bottom: 16px;

    h2 {
      font-size: 24px;
      font-weight: 700;
      color: #1a1a2e;
      margin: 0 0 8px 0;
    }

    .salary-tag {
      font-size: 20px;
      font-weight: 700;
      color: #667eea;
    }
  }

  .detail-meta {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;

    .meta-item {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 14px;
      color: #6b7280;
    }
  }

  .detail-section {
    margin-bottom: 24px;

    h4 {
      font-size: 15px;
      font-weight: 600;
      color: #374151;
      margin: 0 0 12px 0;
    }

    .skills-list, .benefits-list {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }

    .description {
      font-size: 14px;
      color: #6b7280;
      line-height: 1.8;
      margin: 0;
    }

    .requirements-list {
      margin: 0;
      padding-left: 20px;

      li {
        font-size: 14px;
        color: #6b7280;
        line-height: 1.8;
      }
    }
  }

  .detail-actions {
    display: flex;
    gap: 12px;
    margin-top: 32px;
  }
}

@media (max-width: 768px) {
  .job-cards {
    grid-template-columns: 1fr;
  }

  .stats-row {
    .el-col {
      margin-bottom: 12px;
    }
  }
}
</style>
