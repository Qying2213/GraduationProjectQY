<template>
  <div class="portal-companies">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1>企业招聘</h1>
        <p>发现优质企业，开启职业新篇章</p>
      </div>
    </div>

    <!-- 搜索区域 -->
    <div class="search-section">
      <div class="search-container">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索企业名称、行业..."
          size="large"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select v-model="selectedIndustry" placeholder="行业" size="large" clearable>
          <el-option label="全部行业" value="" />
          <el-option label="互联网/IT" value="互联网/IT" />
          <el-option label="金融" value="金融" />
          <el-option label="教育" value="教育" />
          <el-option label="医疗健康" value="医疗健康" />
          <el-option label="电商" value="电商" />
        </el-select>
        <el-select v-model="selectedScale" placeholder="规模" size="large" clearable>
          <el-option label="全部规模" value="" />
          <el-option label="50人以下" value="small" />
          <el-option label="50-200人" value="medium" />
          <el-option label="200-1000人" value="large" />
          <el-option label="1000人以上" value="enterprise" />
        </el-select>
        <el-button type="primary" size="large" @click="handleSearch">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
      </div>
    </div>

    <!-- 企业列表 -->
    <div class="companies-container">
      <div class="companies-grid">
        <div 
          v-for="company in companies" 
          :key="company.id" 
          class="company-card"
          @click="viewCompany(company)"
        >
          <div class="company-header">
            <div class="company-logo" :style="{ background: company.color }">
              {{ company.name.charAt(0) }}
            </div>
            <div class="company-info">
              <h3>{{ company.name }}</h3>
              <div class="company-tags">
                <el-tag size="small" type="info">{{ company.industry }}</el-tag>
                <el-tag size="small" type="info">{{ company.scale }}</el-tag>
              </div>
            </div>
          </div>
          <p class="company-desc">{{ company.description }}</p>
          <div class="company-footer">
            <div class="job-count">
              <el-icon><Briefcase /></el-icon>
              <span>{{ company.jobCount }} 个在招职位</span>
            </div>
            <div class="company-location">
              <el-icon><Location /></el-icon>
              <span>{{ company.location }}</span>
            </div>
          </div>
          <div class="company-benefits">
            <el-tag v-for="benefit in company.benefits.slice(0, 4)" :key="benefit" size="small" effect="plain">
              {{ benefit }}
            </el-tag>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Search, Briefcase, Location } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const searchKeyword = ref('')
const selectedIndustry = ref('')
const selectedScale = ref('')
const currentPage = ref(1)
const pageSize = ref(9)
const total = ref(12)

// 模拟企业数据
const companies = ref([
  {
    id: 1,
    name: '字节跳动',
    industry: '互联网/IT',
    scale: '10000人以上',
    location: '北京',
    description: '字节跳动是一家全球化的科技公司，致力于成为最懂用户的科技公司。',
    jobCount: 128,
    color: '#3370ff',
    benefits: ['五险一金', '带薪年假', '免费三餐', '健身房', '弹性工作']
  },
  {
    id: 2,
    name: '阿里巴巴',
    industry: '电商',
    scale: '10000人以上',
    location: '杭州',
    description: '阿里巴巴集团是全球领先的电子商务和云计算公司。',
    jobCount: 96,
    color: '#ff6a00',
    benefits: ['股票期权', '带薪年假', '节日福利', '培训发展']
  },
  {
    id: 3,
    name: '腾讯',
    industry: '互联网/IT',
    scale: '10000人以上',
    location: '深圳',
    description: '腾讯是中国领先的互联网增值服务提供商。',
    jobCount: 85,
    color: '#07c160',
    benefits: ['五险一金', '免费班车', '年度体检', '股票期权']
  },
  {
    id: 4,
    name: '美团',
    industry: '互联网/IT',
    scale: '10000人以上',
    location: '北京',
    description: '美团是中国领先的生活服务电子商务平台。',
    jobCount: 72,
    color: '#ffc300',
    benefits: ['餐补', '交通补贴', '带薪年假', '团建活动']
  },
  {
    id: 5,
    name: '京东',
    industry: '电商',
    scale: '10000人以上',
    location: '北京',
    description: '京东是中国领先的自营式电商企业。',
    jobCount: 68,
    color: '#e4393c',
    benefits: ['五险一金', '员工折扣', '年终奖', '培训机会']
  },
  {
    id: 6,
    name: '网易',
    industry: '互联网/IT',
    scale: '5000-10000人',
    location: '杭州',
    description: '网易是中国领先的互联网技术公司。',
    jobCount: 54,
    color: '#d43c33',
    benefits: ['免费三餐', '健身房', '带薪年假', '节日福利']
  },
  {
    id: 7,
    name: '百度',
    industry: '互联网/IT',
    scale: '10000人以上',
    location: '北京',
    description: '百度是全球最大的中文搜索引擎。',
    jobCount: 62,
    color: '#2932e1',
    benefits: ['五险一金', '股票期权', '免费班车', '年度体检']
  },
  {
    id: 8,
    name: '小米',
    industry: '互联网/IT',
    scale: '5000-10000人',
    location: '北京',
    description: '小米是一家专注于智能硬件和电子产品研发的公司。',
    jobCount: 45,
    color: '#ff6700',
    benefits: ['员工折扣', '带薪年假', '弹性工作', '培训发展']
  },
  {
    id: 9,
    name: '华为',
    industry: '互联网/IT',
    scale: '10000人以上',
    location: '深圳',
    description: '华为是全球领先的ICT基础设施和智能终端提供商。',
    jobCount: 156,
    color: '#cf0a2c',
    benefits: ['高薪资', '股票分红', '海外机会', '培训体系']
  }
])

const handleSearch = () => {
  ElMessage.info('搜索功能演示')
}

const viewCompany = (company: any) => {
  ElMessage.info(`查看 ${company.name} 详情`)
}

const handlePageChange = (page: number) => {
  currentPage.value = page
}
</script>

<style scoped lang="scss">
.portal-companies {
  min-height: 100vh;
  background: #f5f7fa;
}

.page-header {
  background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%);
  padding: 60px 20px;
  text-align: center;
  color: white;

  h1 {
    font-size: 36px;
    margin: 0 0 12px 0;
  }

  p {
    font-size: 16px;
    opacity: 0.9;
    margin: 0;
  }
}

.search-section {
  background: white;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);

  .search-container {
    max-width: 1000px;
    margin: 0 auto;
    display: flex;
    gap: 12px;

    .el-input {
      flex: 1;
    }

    .el-select {
      width: 140px;
    }
  }
}

.companies-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 20px;
}

.companies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 24px;
}

.company-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #eee;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
  }

  .company-header {
    display: flex;
    gap: 16px;
    margin-bottom: 16px;

    .company-logo {
      width: 56px;
      height: 56px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 24px;
      font-weight: bold;
      flex-shrink: 0;
    }

    .company-info {
      h3 {
        margin: 0 0 8px 0;
        font-size: 18px;
        color: #1a1a2e;
      }

      .company-tags {
        display: flex;
        gap: 8px;
      }
    }
  }

  .company-desc {
    color: #666;
    font-size: 14px;
    line-height: 1.6;
    margin: 0 0 16px 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .company-footer {
    display: flex;
    justify-content: space-between;
    margin-bottom: 12px;
    font-size: 13px;
    color: #666;

    .job-count, .company-location {
      display: flex;
      align-items: center;
      gap: 4px;
    }

    .job-count {
      color: #0ea5e9;
      font-weight: 500;
    }
  }

  .company-benefits {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 32px;
}

@media (max-width: 768px) {
  .search-container {
    flex-direction: column;

    .el-select {
      width: 100% !important;
    }
  }

  .companies-grid {
    grid-template-columns: 1fr;
  }
}
</style>
