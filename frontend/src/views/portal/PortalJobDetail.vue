<template>
  <div class="portal-job-detail">
    <div class="page-container">
      <div class="main-content">
        <!-- 职位信息 -->
        <div class="job-header-card">
          <div class="job-header">
            <div class="job-info">
              <h1>{{ job.title }}</h1>
              <div class="job-salary">{{ job.salary }}</div>
            </div>
            <el-button type="primary" size="large" @click="showApplyDialog = true">
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
          <div class="description" v-html="job.description"></div>
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
          <el-form-item label="简历">
            <el-select v-model="applyForm.resumeId" placeholder="选择简历" style="width: 100%">
              <el-option label="我的简历.pdf" :value="1" />
            </el-select>
          </el-form-item>
          <el-form-item label="求职信">
            <el-input v-model="applyForm.coverLetter" type="textarea" :rows="4" placeholder="简单介绍自己..." />
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
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Location, Timer, School, Calendar, OfficeBuilding, CircleCheck } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const showApplyDialog = ref(false)
const submitting = ref(false)

const job = ref({
  id: 1,
  title: '高级前端工程师',
  salary: '25-45K·15薪',
  location: '北京·朝阳区',
  experience: '3-5年',
  education: '本科',
  postTime: '3天前发布',
  tags: ['Vue3', 'TypeScript', '大厂背景优先'],
  description: `
    <p>我们正在寻找一位经验丰富的高级前端工程师，加入我们的核心产品团队。</p>
    <p><strong>工作内容：</strong></p>
    <p>1. 负责公司核心产品的前端架构设计和开发工作</p>
    <p>2. 主导前端技术选型，制定开发规范</p>
    <p>3. 优化前端性能，提升用户体验</p>
    <p>4. 指导初级工程师，进行代码审查</p>
  `,
  requirements: [
    '5年以上前端开发经验',
    '精通Vue3/React等主流框架',
    '熟悉TypeScript，有大型项目经验',
    '了解Node.js，有全栈开发经验优先',
    '良好的沟通能力和团队协作精神'
  ],
  skills: ['Vue3', 'TypeScript', 'React', 'Webpack', 'Node.js', 'Git'],
  benefits: ['五险一金', '年终奖', '股票期权', '带薪年假', '弹性工作', '免费三餐', '健身房'],
  company: '科技有限公司',
  companyType: '互联网',
  companySize: '500-1000人',
  companyDesc: '我们是一家专注于企业服务的科技公司，致力于用技术提升企业效率。',
  address: '北京市朝阳区望京SOHO T1 20层'
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
  } finally {
    submitting.value = false
  }
}

const goToJob = (id: number) => {
  router.push(`/portal/jobs/${id}`)
}

onMounted(() => {
  // 根据 route.params.id 获取职位详情
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
}

@media (max-width: 900px) {
  .portal-job-detail .page-container {
    grid-template-columns: 1fr;
  }
}
</style>
