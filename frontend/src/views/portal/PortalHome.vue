<template>
  <div class="portal-home">
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-content">
        <h1>找到你的理想工作</h1>
        <p>智能匹配，精准推荐，让求职更简单</p>
        <div class="search-box">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索职位、公司或关键词..."
            size="large"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="searchCity" placeholder="城市" size="large" style="width: 140px">
            <el-option label="全国" value="" />
            <el-option label="北京" value="北京" />
            <el-option label="上海" value="上海" />
            <el-option label="深圳" value="深圳" />
            <el-option label="杭州" value="杭州" />
            <el-option label="广州" value="广州" />
          </el-select>
          <el-button type="primary" size="large" @click="handleSearch">
            搜索职位
          </el-button>
        </div>
        <div class="hot-keywords">
          <span>热门搜索：</span>
          <el-tag v-for="keyword in hotKeywords" :key="keyword" @click="quickSearch(keyword)">
            {{ keyword }}
          </el-tag>
        </div>
      </div>
      <div class="hero-stats">
        <div class="stat-item">
          <span class="stat-number">10,000+</span>
          <span class="stat-label">在线职位</span>
        </div>
        <div class="stat-item">
          <span class="stat-number">5,000+</span>
          <span class="stat-label">入驻企业</span>
        </div>
        <div class="stat-item">
          <span class="stat-number">98%</span>
          <span class="stat-label">匹配成功率</span>
        </div>
      </div>
    </section>

    <!-- 热门职位 -->
    <section class="section">
      <div class="section-container">
        <div class="section-header">
          <h2>热门职位</h2>
          <el-button text type="primary" @click="$router.push('/portal/jobs')">
            查看更多 <el-icon><ArrowRight /></el-icon>
          </el-button>
        </div>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :lg="6" v-for="job in hotJobs" :key="job.id">
            <div class="job-card" @click="goToJob(job.id)">
              <div class="job-header">
                <h3>{{ job.title }}</h3>
                <span class="salary">{{ job.salary }}</span>
              </div>
              <div class="job-company">{{ job.company }}</div>
              <div class="job-tags">
                <el-tag size="small" type="info">{{ job.location }}</el-tag>
                <el-tag size="small" type="info">{{ job.experience }}</el-tag>
                <el-tag size="small" type="info">{{ job.education }}</el-tag>
              </div>
              <div class="job-footer">
                <span class="post-time">{{ job.postTime }}</span>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </section>

    <!-- 热门分类 -->
    <section class="section section-gray">
      <div class="section-container">
        <div class="section-header">
          <h2>职位分类</h2>
        </div>
        <el-row :gutter="20">
          <el-col :xs="12" :sm="8" :lg="4" v-for="category in categories" :key="category.name">
            <div class="category-card" @click="searchByCategory(category.name)">
              <el-icon :size="32" :style="{ color: category.color }">
                <component :is="category.icon" />
              </el-icon>
              <span class="category-name">{{ category.name }}</span>
              <span class="category-count">{{ category.count }}个职位</span>
            </div>
          </el-col>
        </el-row>
      </div>
    </section>

    <!-- 平台优势 -->
    <section class="section">
      <div class="section-container">
        <div class="section-header center">
          <h2>为什么选择我们</h2>
          <p>智能化招聘平台，让求职更高效</p>
        </div>
        <el-row :gutter="40">
          <el-col :xs="24" :sm="8" v-for="feature in features" :key="feature.title">
            <div class="feature-card">
              <div class="feature-icon" :style="{ background: feature.color }">
                <el-icon :size="28"><component :is="feature.icon" /></el-icon>
              </div>
              <h3>{{ feature.title }}</h3>
              <p>{{ feature.desc }}</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, markRaw } from 'vue'
import { useRouter } from 'vue-router'
import {
  Search, ArrowRight, Monitor, DataLine, Cpu, Promotion, Service, Coin,
  MagicStick, Document, Timer
} from '@element-plus/icons-vue'

const router = useRouter()
const searchKeyword = ref('')
const searchCity = ref('')

const hotKeywords = ['前端开发', 'Java', '产品经理', 'UI设计', '数据分析', 'Go开发']

const hotJobs = ref([
  { id: 1, title: '高级前端工程师', salary: '25-45K', company: '科技有限公司', location: '北京', experience: '3-5年', education: '本科', postTime: '3天前' },
  { id: 2, title: '后端开发工程师', salary: '30-50K', company: '互联网公司', location: '上海', experience: '3-5年', education: '本科', postTime: '1天前' },
  { id: 3, title: '产品经理', salary: '20-35K', company: '创新科技', location: '深圳', experience: '2-4年', education: '本科', postTime: '2天前' },
  { id: 4, title: 'UI设计师', salary: '15-25K', company: '设计工作室', location: '杭州', experience: '1-3年', education: '本科', postTime: '今天' },
])

const categories = ref([
  { name: '技术开发', count: 3200, icon: markRaw(Monitor), color: '#0ea5e9' },
  { name: '产品运营', count: 1800, icon: markRaw(DataLine), color: '#8b5cf6' },
  { name: '人工智能', count: 960, icon: markRaw(Cpu), color: '#10b981' },
  { name: '市场营销', count: 1200, icon: markRaw(Promotion), color: '#f59e0b' },
  { name: '客户服务', count: 800, icon: markRaw(Service), color: '#ec4899' },
  { name: '金融财务', count: 650, icon: markRaw(Coin), color: '#06b6d4' },
])

const features = ref([
  { title: '智能匹配', desc: '基于AI算法，精准匹配职位与人才，提高求职效率', icon: markRaw(MagicStick), color: 'linear-gradient(135deg, #0ea5e9, #06b6d4)' },
  { title: '简历解析', desc: '一键上传简历，AI自动解析提取关键信息', icon: markRaw(Document), color: 'linear-gradient(135deg, #8b5cf6, #a855f7)' },
  { title: '快速响应', desc: '企业快速反馈，缩短求职周期', icon: markRaw(Timer), color: 'linear-gradient(135deg, #10b981, #14b8a6)' },
])

const handleSearch = () => {
  router.push({ path: '/portal/jobs', query: { keyword: searchKeyword.value, city: searchCity.value } })
}

const quickSearch = (keyword: string) => {
  searchKeyword.value = keyword
  handleSearch()
}

const searchByCategory = (category: string) => {
  router.push({ path: '/portal/jobs', query: { category } })
}

const goToJob = (id: number) => {
  router.push(`/portal/jobs/${id}`)
}
</script>

<style scoped lang="scss">
.portal-home {
  .hero-section {
    background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%);
    padding: 80px 24px;
    text-align: center;
    color: white;

    .hero-content {
      max-width: 800px;
      margin: 0 auto;

      h1 {
        font-size: 48px;
        font-weight: 700;
        margin-bottom: 16px;
      }

      p {
        font-size: 20px;
        opacity: 0.9;
        margin-bottom: 40px;
      }

      .search-box {
        display: flex;
        gap: 12px;
        justify-content: center;
        margin-bottom: 24px;

        :deep(.el-input) {
          width: 400px;
        }
      }

      .hot-keywords {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        flex-wrap: wrap;

        span { opacity: 0.8; }

        .el-tag {
          cursor: pointer;
          background: rgba(255,255,255,0.2);
          border: none;
          color: white;

          &:hover { background: rgba(255,255,255,0.3); }
        }
      }
    }

    .hero-stats {
      display: flex;
      justify-content: center;
      gap: 80px;
      margin-top: 60px;

      .stat-item {
        display: flex;
        flex-direction: column;

        .stat-number {
          font-size: 36px;
          font-weight: 700;
        }

        .stat-label {
          font-size: 14px;
          opacity: 0.8;
        }
      }
    }
  }

  .section {
    padding: 60px 24px;

    &.section-gray { background: #f1f5f9; }

    .section-container {
      max-width: 1200px;
      margin: 0 auto;
    }

    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 32px;

      &.center {
        flex-direction: column;
        text-align: center;

        p { color: #64748b; margin-top: 8px; }
      }

      h2 {
        font-size: 28px;
        font-weight: 700;
        color: #1e293b;
      }
    }
  }

  .job-card {
    background: white;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 20px;
    cursor: pointer;
    transition: all 0.3s;
    border: 1px solid #e2e8f0;

    &:hover {
      box-shadow: 0 8px 24px rgba(0,0,0,0.1);
      transform: translateY(-4px);
    }

    .job-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 8px;

      h3 {
        font-size: 16px;
        font-weight: 600;
        color: #1e293b;
      }

      .salary {
        color: #0ea5e9;
        font-weight: 600;
      }
    }

    .job-company {
      color: #64748b;
      font-size: 14px;
      margin-bottom: 12px;
    }

    .job-tags {
      display: flex;
      gap: 6px;
      margin-bottom: 12px;
    }

    .job-footer {
      color: #94a3b8;
      font-size: 12px;
    }
  }

  .category-card {
    background: white;
    border-radius: 12px;
    padding: 24px;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s;
    margin-bottom: 20px;

    &:hover {
      box-shadow: 0 8px 24px rgba(0,0,0,0.1);
      transform: translateY(-4px);
    }

    .category-name {
      display: block;
      font-weight: 600;
      color: #1e293b;
      margin: 12px 0 4px;
    }

    .category-count {
      font-size: 12px;
      color: #94a3b8;
    }
  }

  .feature-card {
    text-align: center;
    padding: 32px;

    .feature-icon {
      width: 64px;
      height: 64px;
      border-radius: 16px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0 auto 20px;
      color: white;
    }

    h3 {
      font-size: 20px;
      font-weight: 600;
      color: #1e293b;
      margin-bottom: 12px;
    }

    p {
      color: #64748b;
      line-height: 1.6;
    }
  }
}

@media (max-width: 768px) {
  .portal-home {
    .hero-section {
      padding: 40px 16px;

      .hero-content h1 { font-size: 28px; }

      .search-box {
        flex-direction: column;
        :deep(.el-input) { width: 100%; }
      }
    }

    .hero-stats { gap: 24px; }
  }
}
</style>
