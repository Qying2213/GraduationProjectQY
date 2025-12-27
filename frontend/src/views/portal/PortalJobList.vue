<template>
  <div class="portal-job-list">
    <div class="page-container">
      <!-- 搜索筛选 -->
      <div class="search-section">
        <div class="search-bar">
          <el-input v-model="searchParams.keyword" placeholder="搜索职位、公司..." size="large" @keyup.enter="fetchJobs">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <el-select v-model="searchParams.city" placeholder="城市" size="large" clearable>
            <el-option label="全国" value="" />
            <el-option v-for="city in cities" :key="city" :label="city" :value="city" />
          </el-select>
          <el-button type="primary" size="large" @click="fetchJobs">搜索</el-button>
        </div>
        <div class="filter-tags">
          <div class="filter-group">
            <span class="filter-label">经验：</span>
            <el-tag v-for="exp in experienceOptions" :key="exp.value"
                    :type="searchParams.experience === exp.value ? '' : 'info'"
                    @click="filterByExperience(exp.value)">
              {{ exp.label }}
            </el-tag>
          </div>
          <div class="filter-group">
            <span class="filter-label">学历：</span>
            <el-tag v-for="edu in educationOptions" :key="edu.value"
                    :type="searchParams.education === edu.value ? '' : 'info'"
                    @click="filterByEducation(edu.value)">
              {{ edu.label }}
            </el-tag>
          </div>
        </div>
      </div>

      <div class="content-wrapper">
        <!-- 职位列表 -->
        <div class="job-list">
          <div class="list-header">
            <span>共找到 <strong>{{ total }}</strong> 个职位</span>
            <el-select v-model="sortBy" size="small" style="width: 120px">
              <el-option label="最新发布" value="latest" />
              <el-option label="薪资最高" value="salary" />
            </el-select>
          </div>

          <div v-loading="loading">
            <div class="job-item" v-for="job in jobs" :key="job.id" @click="goToDetail(job.id)">
              <div class="job-main">
                <div class="job-info">
                  <h3 class="job-title">{{ job.title }}</h3>
                  <div class="job-meta">
                    <span><el-icon><Location /></el-icon>{{ job.location }}</span>
                    <span><el-icon><Timer /></el-icon>{{ job.experience }}</span>
                    <span><el-icon><School /></el-icon>{{ job.education }}</span>
                  </div>
                  <div class="job-tags">
                    <el-tag v-for="skill in job.skills?.slice(0, 4)" :key="skill" size="small" type="info">
                      {{ skill }}
                    </el-tag>
                  </div>
                </div>
                <div class="job-salary">{{ job.salary }}</div>
              </div>
              <div class="job-company">
                <div class="company-logo">
                  <el-icon :size="24"><OfficeBuilding /></el-icon>
                </div>
                <div class="company-info">
                  <span class="company-name">{{ job.company }}</span>
                  <span class="company-meta">{{ job.companySize }} · {{ job.companyType }}</span>
                </div>
                <el-button type="primary" size="small" @click.stop="applyJob(job)">
                  投递简历
                </el-button>
              </div>
            </div>

            <el-empty v-if="!loading && jobs.length === 0" description="暂无匹配的职位" />
          </div>

          <div class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              :total="total"
              :page-size="pageSize"
              layout="prev, pager, next"
              @current-change="fetchJobs"
            />
          </div>
        </div>

        <!-- 侧边栏 -->
        <div class="sidebar">
          <div class="sidebar-card">
            <h4>热门搜索</h4>
            <div class="hot-tags">
              <el-tag v-for="tag in hotTags" :key="tag" @click="quickSearch(tag)">{{ tag }}</el-tag>
            </div>
          </div>
          <div class="sidebar-card">
            <h4>薪资分布</h4>
            <div class="salary-ranges">
              <div class="salary-item" v-for="range in salaryRanges" :key="range.label"
                   @click="filterBySalary(range.value)">
                <span>{{ range.label }}</span>
                <span class="count">{{ range.count }}个</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 投递弹窗 -->
    <el-dialog v-model="showApplyDialog" title="投递简历" width="500px">
      <div class="apply-content">
        <div class="apply-job-info">
          <h3>{{ currentJob?.title }}</h3>
          <p>{{ currentJob?.company }} · {{ currentJob?.location }}</p>
        </div>
        <el-form :model="applyForm" label-width="80px">
          <el-form-item label="简历">
            <el-select v-model="applyForm.resumeId" placeholder="选择简历" style="width: 100%">
              <el-option label="我的简历.pdf" :value="1" />
            </el-select>
            <div class="upload-tip">
              或 <el-button link type="primary">上传新简历</el-button>
            </div>
          </el-form-item>
          <el-form-item label="求职信">
            <el-input v-model="applyForm.coverLetter" type="textarea" :rows="4"
                      placeholder="简单介绍自己，为什么适合这个职位..." />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="showApplyDialog = false">取消</el-button>
        <el-button type="primary" @click="submitApplication" :loading="submitting">确认投递</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Location, Timer, School, OfficeBuilding } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const submitting = ref(false)
const jobs = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const sortBy = ref('latest')
const showApplyDialog = ref(false)
const currentJob = ref<any>(null)

const searchParams = reactive({
  keyword: '',
  city: '',
  experience: '',
  education: '',
  salary: ''
})

const applyForm = reactive({
  resumeId: null as number | null,
  coverLetter: ''
})

const cities = ['北京', '上海', '深圳', '杭州', '广州', '成都', '南京', '武汉']
const hotTags = ['前端开发', 'Java', 'Python', '产品经理', 'UI设计', '数据分析']

const experienceOptions = [
  { label: '不限', value: '' },
  { label: '应届生', value: '0' },
  { label: '1-3年', value: '1-3' },
  { label: '3-5年', value: '3-5' },
  { label: '5-10年', value: '5-10' },
]

const educationOptions = [
  { label: '不限', value: '' },
  { label: '大专', value: '大专' },
  { label: '本科', value: '本科' },
  { label: '硕士', value: '硕士' },
]

const salaryRanges = [
  { label: '10K以下', value: '0-10', count: 320 },
  { label: '10-20K', value: '10-20', count: 580 },
  { label: '20-30K', value: '20-30', count: 420 },
  { label: '30-50K', value: '30-50', count: 280 },
  { label: '50K以上', value: '50+', count: 120 },
]

const fetchJobs = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: pageSize.value,
      status: 'open'  // 只显示开放的职位
    }
    
    // 添加筛选参数
    if (searchParams.keyword) {
      params.keyword = searchParams.keyword
    }
    if (searchParams.city) {
      params.location = searchParams.city
    }
    if (searchParams.experience) {
      params.experience = searchParams.experience
    }
    if (searchParams.education) {
      params.education = searchParams.education
    }
    if (sortBy.value === 'salary') {
      params.sort_by = 'salary'
      params.sort_order = 'desc'
    }
    
    const res = await request.get('/jobs', { params })
    
    if (res.data?.code === 0 && res.data.data) {
      // 转换后端数据格式为前端显示格式
      jobs.value = (res.data.data.jobs || []).map((job: any) => ({
        id: job.id,
        title: job.title,
        salary: job.salary || '面议',
        company: job.department || '公司',
        companySize: '100-500人',
        companyType: '互联网',
        location: job.location || '全国',
        experience: formatExperience(job.level),
        education: '本科',
        skills: job.skills || []
      }))
      total.value = res.data.data.total || 0
    } else {
      ElMessage.error(res.data?.message || '获取职位列表失败')
    }
  } catch (error) {
    console.error('获取职位列表失败:', error)
    ElMessage.error('获取职位列表失败')
    jobs.value = []
    total.value = 0
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

const goToDetail = (id: number) => {
  router.push(`/portal/jobs/${id}`)
}

const applyJob = (job: any) => {
  currentJob.value = job
  showApplyDialog.value = true
}

const submitApplication = async () => {
  if (!applyForm.resumeId) {
    ElMessage.warning('请选择简历')
    return
  }
  submitting.value = true
  try {
    await new Promise(r => setTimeout(r, 1000))
    ElMessage.success('投递成功！')
    showApplyDialog.value = false
    applyForm.resumeId = null
    applyForm.coverLetter = ''
  } finally {
    submitting.value = false
  }
}

const quickSearch = (keyword: string) => {
  searchParams.keyword = keyword
  fetchJobs()
}

const filterByExperience = (value: string) => {
  searchParams.experience = value
  fetchJobs()
}

const filterByEducation = (value: string) => {
  searchParams.education = value
  fetchJobs()
}

const filterBySalary = (value: string) => {
  searchParams.salary = value
  fetchJobs()
}

onMounted(() => {
  if (route.query.keyword) searchParams.keyword = route.query.keyword as string
  if (route.query.city) searchParams.city = route.query.city as string
  fetchJobs()
})
</script>

<style scoped lang="scss">
.portal-job-list {
  padding: 24px;
  background: #f8fafc;
  min-height: calc(100vh - 160px);

  .page-container {
    max-width: 1200px;
    margin: 0 auto;
  }

  .search-section {
    background: white;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 24px;

    .search-bar {
      display: flex;
      gap: 12px;
      margin-bottom: 16px;

      .el-input { flex: 1; }
    }

    .filter-tags {
      display: flex;
      flex-direction: column;
      gap: 12px;

      .filter-group {
        display: flex;
        align-items: center;
        gap: 8px;
        flex-wrap: wrap;

        .filter-label {
          color: #64748b;
          font-size: 14px;
        }

        .el-tag {
          cursor: pointer;
        }
      }
    }
  }

  .content-wrapper {
    display: grid;
    grid-template-columns: 1fr 280px;
    gap: 24px;
  }

  .job-list {
    .list-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      color: #64748b;

      strong { color: #0ea5e9; }
    }

    .job-item {
      background: white;
      border-radius: 12px;
      padding: 20px;
      margin-bottom: 16px;
      cursor: pointer;
      transition: all 0.3s;
      border: 1px solid #e2e8f0;

      &:hover {
        box-shadow: 0 4px 12px rgba(0,0,0,0.08);
        border-color: #0ea5e9;
      }

      .job-main {
        display: flex;
        justify-content: space-between;
        margin-bottom: 16px;

        .job-title {
          font-size: 18px;
          font-weight: 600;
          color: #1e293b;
          margin-bottom: 8px;
        }

        .job-meta {
          display: flex;
          gap: 16px;
          color: #64748b;
          font-size: 14px;
          margin-bottom: 12px;

          span {
            display: flex;
            align-items: center;
            gap: 4px;
          }
        }

        .job-tags {
          display: flex;
          gap: 6px;
        }

        .job-salary {
          font-size: 20px;
          font-weight: 600;
          color: #0ea5e9;
        }
      }

      .job-company {
        display: flex;
        align-items: center;
        gap: 12px;
        padding-top: 16px;
        border-top: 1px solid #f1f5f9;

        .company-logo {
          width: 48px;
          height: 48px;
          background: #f1f5f9;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #94a3b8;
        }

        .company-info {
          flex: 1;

          .company-name {
            display: block;
            font-weight: 500;
            color: #1e293b;
          }

          .company-meta {
            font-size: 12px;
            color: #94a3b8;
          }
        }
      }
    }

    .pagination {
      display: flex;
      justify-content: center;
      margin-top: 24px;
    }
  }

  .sidebar {
    .sidebar-card {
      background: white;
      border-radius: 12px;
      padding: 20px;
      margin-bottom: 16px;

      h4 {
        font-size: 16px;
        font-weight: 600;
        color: #1e293b;
        margin-bottom: 16px;
      }

      .hot-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;

        .el-tag { cursor: pointer; }
      }

      .salary-ranges {
        .salary-item {
          display: flex;
          justify-content: space-between;
          padding: 8px 0;
          cursor: pointer;
          color: #64748b;
          border-bottom: 1px solid #f1f5f9;

          &:hover { color: #0ea5e9; }
          &:last-child { border-bottom: none; }

          .count { font-size: 12px; }
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

    h3 { font-size: 16px; margin-bottom: 4px; }
    p { color: #64748b; font-size: 14px; margin: 0; }
  }

  .upload-tip {
    margin-top: 8px;
    font-size: 12px;
    color: #94a3b8;
  }
}

@media (max-width: 900px) {
  .portal-job-list .content-wrapper {
    grid-template-columns: 1fr;
  }

  .sidebar { display: none; }
}
</style>
