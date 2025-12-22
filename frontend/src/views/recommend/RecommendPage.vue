<template>
  <div class="recommend-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1>智能推荐</h1>
        <p class="subtitle">基于 AI 算法为您精准匹配最佳人才与职位</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="refreshRecommendations">
          <el-icon><Refresh /></el-icon>
          刷新推荐
        </el-button>
      </div>
    </div>

    <!-- 推荐模式切换 -->
    <div class="mode-tabs">
      <div
        class="mode-tab"
        :class="{ active: activeMode === 'talent' }"
        @click="activeMode = 'talent'"
      >
        <div class="tab-icon">
          <el-icon><User /></el-icon>
        </div>
        <div class="tab-info">
          <span class="tab-title">人才推荐</span>
          <span class="tab-desc">为职位匹配合适人才</span>
        </div>
      </div>
      <div
        class="mode-tab"
        :class="{ active: activeMode === 'job' }"
        @click="activeMode = 'job'"
      >
        <div class="tab-icon">
          <el-icon><Suitcase /></el-icon>
        </div>
        <div class="tab-info">
          <span class="tab-title">职位推荐</span>
          <span class="tab-desc">为人才匹配合适职位</span>
        </div>
      </div>
    </div>

    <!-- 筛选条件 -->
    <div class="filter-section card">
      <div class="filter-row">
        <div class="filter-item">
          <label>{{ activeMode === 'talent' ? '选择职位' : '选择人才' }}</label>
          <el-select
            v-model="selectedTarget"
            :placeholder="activeMode === 'talent' ? '请选择要推荐人才的职位' : '请选择要推荐职位的人才'"
            style="width: 300px"
            filterable
          >
            <el-option
              v-for="item in targetOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </div>
        <div class="filter-item">
          <label>最低匹配度</label>
          <el-slider
            v-model="minMatchScore"
            :min="0"
            :max="100"
            :format-tooltip="(val: number) => `${val}%`"
            style="width: 200px"
          />
        </div>
        <div class="filter-item">
          <label>排序方式</label>
          <el-select v-model="sortBy" style="width: 150px">
            <el-option label="匹配度优先" value="match" />
            <el-option label="最新优先" value="latest" />
            <el-option label="经验优先" value="experience" />
          </el-select>
        </div>
      </div>
    </div>

    <!-- 推荐结果 -->
    <div class="recommendation-results" v-loading="loading">
      <div v-if="recommendations.length === 0 && !loading" class="empty-state">
        <el-empty description="请选择职位/人才以获取推荐结果">
          <el-button type="primary" @click="loadRecommendations">开始匹配</el-button>
        </el-empty>
      </div>

      <div v-else class="recommendation-grid">
        <div
          v-for="(item, index) in recommendations"
          :key="item.id"
          class="recommendation-card"
          :style="{ animationDelay: `${index * 0.1}s` }"
        >
          <!-- 匹配度圆环 -->
          <div class="match-score-ring" :class="getScoreClass(item.matchScore)">
            <svg viewBox="0 0 100 100">
              <circle
                class="bg-circle"
                cx="50" cy="50" r="45"
                fill="none"
                stroke-width="8"
              />
              <circle
                class="progress-circle"
                cx="50" cy="50" r="45"
                fill="none"
                stroke-width="8"
                :stroke-dasharray="`${item.matchScore * 2.83} 283`"
                stroke-linecap="round"
              />
            </svg>
            <div class="score-value">
              <span class="number">{{ item.matchScore }}</span>
              <span class="percent">%</span>
            </div>
            <div class="score-label">匹配度</div>
          </div>

          <!-- 信息卡片 -->
          <div class="card-content">
            <div class="card-header">
              <div class="avatar" v-if="activeMode === 'talent'">
                <img :src="item.avatar || getDefaultAvatar(item.name)" :alt="item.name" />
              </div>
              <div class="company-logo" v-else>
                <el-icon><OfficeBuilding /></el-icon>
              </div>
              <div class="basic-info">
                <h3 class="name">{{ item.name }}</h3>
                <p class="sub-info">
                  {{ activeMode === 'talent' ? item.position : item.company }}
                </p>
              </div>
              <el-tag
                :type="item.status === 'active' ? 'success' : 'info'"
                size="small"
              >
                {{ item.status === 'active' ? '活跃' : '一般' }}
              </el-tag>
            </div>

            <!-- 关键信息 -->
            <div class="key-info">
              <div class="info-item" v-if="activeMode === 'talent'">
                <el-icon><Calendar /></el-icon>
                <span>{{ item.experience }}年经验</span>
              </div>
              <div class="info-item" v-if="activeMode === 'talent'">
                <el-icon><Location /></el-icon>
                <span>{{ item.location }}</span>
              </div>
              <div class="info-item" v-if="activeMode === 'talent'">
                <el-icon><Money /></el-icon>
                <span>{{ item.salary }}</span>
              </div>
              <div class="info-item" v-if="activeMode === 'job'">
                <el-icon><Location /></el-icon>
                <span>{{ item.location }}</span>
              </div>
              <div class="info-item" v-if="activeMode === 'job'">
                <el-icon><Money /></el-icon>
                <span>{{ item.salaryRange }}</span>
              </div>
              <div class="info-item" v-if="activeMode === 'job'">
                <el-icon><Tickets /></el-icon>
                <span>{{ item.type }}</span>
              </div>
            </div>

            <!-- 技能标签 -->
            <div class="skills">
              <el-tag
                v-for="skill in item.skills?.slice(0, 4)"
                :key="skill"
                size="small"
                :type="item.matchedSkills?.includes(skill) ? 'success' : 'info'"
                effect="plain"
              >
                {{ skill }}
                <el-icon v-if="item.matchedSkills?.includes(skill)" class="match-icon">
                  <Check />
                </el-icon>
              </el-tag>
              <el-tag v-if="item.skills?.length > 4" size="small" type="info">
                +{{ item.skills.length - 4 }}
              </el-tag>
            </div>

            <!-- 匹配原因 -->
            <div class="match-reasons">
              <div class="reason-title">
                <el-icon><MagicStick /></el-icon>
                <span>匹配亮点</span>
              </div>
              <div class="reason-list">
                <div
                  v-for="(reason, idx) in item.matchReasons"
                  :key="idx"
                  class="reason-item"
                >
                  <el-icon><CircleCheck /></el-icon>
                  <span>{{ reason }}</span>
                </div>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="card-actions">
              <el-button type="primary" @click="viewDetail(item)">
                查看详情
              </el-button>
              <el-button @click="sendInvitation(item)">
                {{ activeMode === 'talent' ? '发送邀请' : '申请职位' }}
              </el-button>
              <el-button text type="info" @click="addToFavorites(item)">
                <el-icon><Star /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="recommendations.length > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[6, 12, 24, 48]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 详情抽屉 -->
    <el-drawer
      v-model="showDetailDrawer"
      :title="activeMode === 'talent' ? '人才详情' : '职位详情'"
      size="600px"
    >
      <div class="detail-drawer" v-if="currentDetail">
        <!-- 匹配度大圆环 -->
        <div class="detail-match-score" :class="getScoreClass(currentDetail.matchScore)">
          <div class="score-ring">
            <svg viewBox="0 0 120 120">
              <circle
                class="bg-circle"
                cx="60" cy="60" r="54"
                fill="none"
                stroke-width="10"
              />
              <circle
                class="progress-circle"
                cx="60" cy="60" r="54"
                fill="none"
                stroke-width="10"
                :stroke-dasharray="`${currentDetail.matchScore * 3.39} 339`"
                stroke-linecap="round"
              />
            </svg>
            <div class="score-content">
              <span class="score-number">{{ currentDetail.matchScore }}</span>
              <span class="score-percent">%</span>
              <span class="score-text">匹配度</span>
            </div>
          </div>
        </div>

        <!-- 匹配维度分析 -->
        <div class="match-dimensions">
          <h4>匹配维度分析</h4>
          <div class="dimension-list">
            <div
              v-for="dim in currentDetail.dimensions"
              :key="dim.name"
              class="dimension-item"
            >
              <div class="dim-header">
                <span class="dim-name">{{ dim.name }}</span>
                <span class="dim-score">{{ dim.score }}%</span>
              </div>
              <el-progress
                :percentage="dim.score"
                :stroke-width="8"
                :color="getDimensionColor(dim.score)"
              />
            </div>
          </div>
        </div>

        <!-- 详细信息 -->
        <div class="detail-info">
          <h4>{{ activeMode === 'talent' ? '人才信息' : '职位信息' }}</h4>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="姓名/职位">
              {{ currentDetail.name }}
            </el-descriptions-item>
            <el-descriptions-item :label="activeMode === 'talent' ? '当前职位' : '公司'">
              {{ activeMode === 'talent' ? currentDetail.position : currentDetail.company }}
            </el-descriptions-item>
            <el-descriptions-item label="工作地点">
              {{ currentDetail.location }}
            </el-descriptions-item>
            <el-descriptions-item label="薪资">
              {{ activeMode === 'talent' ? currentDetail.salary : currentDetail.salaryRange }}
            </el-descriptions-item>
            <el-descriptions-item v-if="activeMode === 'talent'" label="工作经验">
              {{ currentDetail.experience }}年
            </el-descriptions-item>
            <el-descriptions-item label="技能">
              <el-tag
                v-for="skill in currentDetail.skills"
                :key="skill"
                size="small"
                style="margin-right: 4px; margin-bottom: 4px;"
              >
                {{ skill }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 推荐理由 -->
        <div class="recommend-reasons">
          <h4>推荐理由</h4>
          <div class="reasons-content">
            <div
              v-for="(reason, idx) in currentDetail.matchReasons"
              :key="idx"
              class="reason-row"
            >
              <el-icon class="reason-icon"><CircleCheckFilled /></el-icon>
              <span>{{ reason }}</span>
            </div>
          </div>
        </div>

        <!-- 操作区 -->
        <div class="detail-actions">
          <el-button type="primary" size="large" @click="sendInvitation(currentDetail)">
            {{ activeMode === 'talent' ? '发送面试邀请' : '申请该职位' }}
          </el-button>
          <el-button size="large" @click="addToFavorites(currentDetail)">
            <el-icon><Star /></el-icon>
            收藏
          </el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh, User, Suitcase, OfficeBuilding, Calendar, Location,
  Money, Tickets, Check, MagicStick, CircleCheck, Star,
  CircleCheckFilled
} from '@element-plus/icons-vue'

// 类型定义
interface Recommendation {
  id: number
  name: string
  avatar?: string
  position?: string
  company?: string
  location: string
  salary?: string
  salaryRange?: string
  experience?: number
  type?: string
  skills: string[]
  matchedSkills?: string[]
  matchScore: number
  matchReasons: string[]
  status: string
  dimensions?: MatchDimension[]
}

interface MatchDimension {
  name: string
  score: number
}

interface TargetOption {
  id: number
  name: string
}

// 状态
const activeMode = ref<'talent' | 'job'>('talent')
const selectedTarget = ref<number | null>(null)
const minMatchScore = ref(60)
const sortBy = ref('match')
const loading = ref(false)
const recommendations = ref<Recommendation[]>([])
const currentPage = ref(1)
const pageSize = ref(6)
const total = ref(0)
const showDetailDrawer = ref(false)
const currentDetail = ref<Recommendation | null>(null)

// 计算目标选项
const targetOptions = computed<TargetOption[]>(() => {
  if (activeMode.value === 'talent') {
    return [
      { id: 1, name: '高级前端工程师 - 技术部' },
      { id: 2, name: '产品经理 - 产品部' },
      { id: 3, name: 'Java后端开发 - 技术部' },
      { id: 4, name: 'UI设计师 - 设计部' },
      { id: 5, name: '数据分析师 - 数据部' }
    ]
  } else {
    return [
      { id: 1, name: '张伟 - 前端工程师' },
      { id: 2, name: '李娜 - 产品经理' },
      { id: 3, name: '王强 - 后端工程师' },
      { id: 4, name: '刘芳 - UI设计师' },
      { id: 5, name: '陈明 - 数据分析师' }
    ]
  }
})

// 监听模式切换
watch(activeMode, () => {
  selectedTarget.value = null
  recommendations.value = []
})

// 监听筛选条件变化
watch([selectedTarget, minMatchScore, sortBy], () => {
  if (selectedTarget.value) {
    loadRecommendations()
  }
})

// 获取匹配度等级样式
const getScoreClass = (score: number) => {
  if (score >= 90) return 'excellent'
  if (score >= 75) return 'good'
  if (score >= 60) return 'normal'
  return 'low'
}

// 获取维度进度条颜色
const getDimensionColor = (score: number) => {
  if (score >= 90) return '#67c23a'
  if (score >= 75) return '#409eff'
  if (score >= 60) return '#e6a23c'
  return '#f56c6c'
}

// 获取默认头像
const getDefaultAvatar = (name: string) => {
  const colors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399']
  const index = name.charCodeAt(0) % colors.length
  return `https://api.dicebear.com/7.x/initials/svg?seed=${encodeURIComponent(name)}&backgroundColor=${colors[index].slice(1)}`
}

// 生成模拟推荐数据
const generateMockRecommendations = (): Recommendation[] => {
  if (activeMode.value === 'talent') {
    return [
      {
        id: 1,
        name: '张伟',
        position: '高级前端工程师',
        location: '北京',
        salary: '30-40K',
        experience: 6,
        skills: ['Vue', 'React', 'TypeScript', 'Node.js', 'Webpack'],
        matchedSkills: ['Vue', 'TypeScript'],
        matchScore: 95,
        matchReasons: ['技能高度匹配，熟练掌握Vue和TypeScript', '6年工作经验符合要求', '期望薪资在预算范围内', '地点匹配，可立即到岗'],
        status: 'active',
        dimensions: [
          { name: '技能匹配', score: 98 },
          { name: '经验匹配', score: 95 },
          { name: '薪资匹配', score: 90 },
          { name: '地点匹配', score: 100 },
          { name: '稳定性', score: 85 }
        ]
      },
      {
        id: 2,
        name: '李娜',
        position: '资深前端开发',
        location: '上海',
        salary: '25-35K',
        experience: 5,
        skills: ['Vue', 'JavaScript', 'CSS', 'Element Plus', 'Git'],
        matchedSkills: ['Vue', 'Element Plus'],
        matchScore: 88,
        matchReasons: ['Vue技术栈匹配', '有大型项目经验', '沟通能力强', '学习能力突出'],
        status: 'active',
        dimensions: [
          { name: '技能匹配', score: 90 },
          { name: '经验匹配', score: 88 },
          { name: '薪资匹配', score: 95 },
          { name: '地点匹配', score: 75 },
          { name: '稳定性', score: 90 }
        ]
      },
      {
        id: 3,
        name: '王强',
        position: '全栈工程师',
        location: '深圳',
        salary: '28-38K',
        experience: 4,
        skills: ['Vue', 'React', 'Node.js', 'Python', 'MySQL'],
        matchedSkills: ['Vue', 'Node.js'],
        matchScore: 82,
        matchReasons: ['全栈能力可独立承担项目', '有创业公司经验', '适应能力强'],
        status: 'active',
        dimensions: [
          { name: '技能匹配', score: 85 },
          { name: '经验匹配', score: 78 },
          { name: '薪资匹配', score: 88 },
          { name: '地点匹配', score: 70 },
          { name: '稳定性', score: 80 }
        ]
      },
      {
        id: 4,
        name: '刘芳',
        position: '前端工程师',
        location: '杭州',
        salary: '20-28K',
        experience: 3,
        skills: ['Vue', 'React', 'TypeScript', 'Sass', 'Vite'],
        matchedSkills: ['Vue', 'TypeScript', 'Vite'],
        matchScore: 78,
        matchReasons: ['技术栈契合度高', '有良好的代码规范', '积极主动'],
        status: 'normal',
        dimensions: [
          { name: '技能匹配', score: 88 },
          { name: '经验匹配', score: 70 },
          { name: '薪资匹配', score: 92 },
          { name: '地点匹配', score: 65 },
          { name: '稳定性', score: 75 }
        ]
      },
      {
        id: 5,
        name: '陈明',
        position: '前端开发工程师',
        location: '成都',
        salary: '18-25K',
        experience: 2,
        skills: ['Vue', 'JavaScript', 'HTML', 'CSS', 'Webpack'],
        matchedSkills: ['Vue'],
        matchScore: 72,
        matchReasons: ['基础扎实', '学习热情高', '性价比高'],
        status: 'normal',
        dimensions: [
          { name: '技能匹配', score: 75 },
          { name: '经验匹配', score: 60 },
          { name: '薪资匹配', score: 98 },
          { name: '地点匹配', score: 60 },
          { name: '稳定性', score: 70 }
        ]
      },
      {
        id: 6,
        name: '赵雪',
        position: '中级前端工程师',
        location: '广州',
        salary: '22-30K',
        experience: 4,
        skills: ['React', 'Vue', 'TypeScript', 'Ant Design', 'Redux'],
        matchedSkills: ['Vue', 'TypeScript'],
        matchScore: 68,
        matchReasons: ['React背景可快速上手Vue', '有移动端开发经验'],
        status: 'normal',
        dimensions: [
          { name: '技能匹配', score: 70 },
          { name: '经验匹配', score: 80 },
          { name: '薪资匹配', score: 85 },
          { name: '地点匹配', score: 55 },
          { name: '稳定性', score: 65 }
        ]
      }
    ]
  } else {
    return [
      {
        id: 1,
        name: '高级前端工程师',
        company: '字节跳动',
        location: '北京',
        salaryRange: '35-50K',
        type: '全职',
        skills: ['Vue', 'React', 'TypeScript', '微前端', '性能优化'],
        matchedSkills: ['Vue', 'React', 'TypeScript'],
        matchScore: 92,
        matchReasons: ['技术栈完全匹配', '大厂背景有利于职业发展', '薪资涨幅50%以上', '团队氛围好'],
        status: 'active',
        dimensions: [
          { name: '技能匹配', score: 95 },
          { name: '薪资提升', score: 90 },
          { name: '发展空间', score: 95 },
          { name: '公司规模', score: 98 },
          { name: '工作环境', score: 88 }
        ]
      },
      {
        id: 2,
        name: '前端技术专家',
        company: '阿里巴巴',
        location: '杭州',
        salaryRange: '40-60K',
        type: '全职',
        skills: ['Vue', 'Node.js', '跨端开发', '架构设计', '团队管理'],
        matchedSkills: ['Vue', 'Node.js'],
        matchScore: 85,
        matchReasons: ['晋升为技术专家的好机会', '阿里系资源丰富', '技术氛围浓厚'],
        status: 'active',
        dimensions: [
          { name: '技能匹配', score: 82 },
          { name: '薪资提升', score: 95 },
          { name: '发展空间', score: 92 },
          { name: '公司规模', score: 98 },
          { name: '工作环境', score: 85 }
        ]
      },
      {
        id: 3,
        name: '前端负责人',
        company: '美团',
        location: '北京',
        salaryRange: '45-55K',
        type: '全职',
        skills: ['Vue', 'React', '团队管理', '项目规划', '技术选型'],
        matchedSkills: ['Vue', 'React'],
        matchScore: 80,
        matchReasons: ['管理岗位，职业转型机会', '团队规模10人以上', '业务稳定'],
        status: 'active',
        dimensions: [
          { name: '技能匹配', score: 78 },
          { name: '薪资提升', score: 88 },
          { name: '发展空间', score: 90 },
          { name: '公司规模', score: 95 },
          { name: '工作环境', score: 82 }
        ]
      },
      {
        id: 4,
        name: '资深前端工程师',
        company: '腾讯',
        location: '深圳',
        salaryRange: '38-48K',
        type: '全职',
        skills: ['Vue', 'TypeScript', '小程序', '性能优化', 'WebGL'],
        matchedSkills: ['Vue', 'TypeScript'],
        matchScore: 76,
        matchReasons: ['腾讯系福利完善', '技术积累深厚', '可接触前沿技术'],
        status: 'normal',
        dimensions: [
          { name: '技能匹配', score: 80 },
          { name: '薪资提升', score: 75 },
          { name: '发展空间', score: 85 },
          { name: '公司规模', score: 98 },
          { name: '工作环境', score: 78 }
        ]
      },
      {
        id: 5,
        name: '全栈工程师',
        company: 'PingCAP',
        location: '北京',
        salaryRange: '35-45K',
        type: '全职',
        skills: ['Vue', 'Go', 'TypeScript', '分布式系统', 'Kubernetes'],
        matchedSkills: ['Vue', 'TypeScript'],
        matchScore: 70,
        matchReasons: ['技术驱动型公司', '开源文化浓厚', '学习成长空间大'],
        status: 'normal',
        dimensions: [
          { name: '技能匹配', score: 72 },
          { name: '薪资提升', score: 70 },
          { name: '发展空间', score: 88 },
          { name: '公司规模', score: 75 },
          { name: '工作环境', score: 90 }
        ]
      },
      {
        id: 6,
        name: '前端工程师',
        company: '蚂蚁金服',
        location: '杭州',
        salaryRange: '32-42K',
        type: '全职',
        skills: ['React', 'TypeScript', 'Ant Design', '数据可视化', '低代码'],
        matchedSkills: ['TypeScript'],
        matchScore: 65,
        matchReasons: ['金融科技领域机会', '蚂蚁系股票期权', '工作生活平衡'],
        status: 'normal',
        dimensions: [
          { name: '技能匹配', score: 65 },
          { name: '薪资提升', score: 72 },
          { name: '发展空间', score: 80 },
          { name: '公司规模', score: 95 },
          { name: '工作环境', score: 75 }
        ]
      }
    ]
  }
}

// 加载推荐数据
const loadRecommendations = () => {
  loading.value = true

  setTimeout(() => {
    let data = generateMockRecommendations()

    // 筛选最低匹配度
    data = data.filter(item => item.matchScore >= minMatchScore.value)

    // 排序
    if (sortBy.value === 'match') {
      data.sort((a, b) => b.matchScore - a.matchScore)
    } else if (sortBy.value === 'experience' && activeMode.value === 'talent') {
      data.sort((a, b) => (b.experience || 0) - (a.experience || 0))
    }

    total.value = data.length
    recommendations.value = data
    loading.value = false
  }, 800)
}

// 刷新推荐
const refreshRecommendations = () => {
  if (selectedTarget.value) {
    loadRecommendations()
    ElMessage.success('推荐结果已刷新')
  } else {
    ElMessage.warning('请先选择职位或人才')
  }
}

// 查看详情
const viewDetail = (item: Recommendation) => {
  currentDetail.value = item
  showDetailDrawer.value = true
}

// 发送邀请
const sendInvitation = (item: Recommendation) => {
  ElMessage.success(`已向 ${item.name} 发送${activeMode.value === 'talent' ? '面试邀请' : '申请'}`)
}

// 添加收藏
const addToFavorites = (item: Recommendation) => {
  ElMessage.success(`已将 ${item.name} 添加到收藏`)
}

// 分页处理
const handleSizeChange = () => {
  currentPage.value = 1
  loadRecommendations()
}

const handlePageChange = () => {
  loadRecommendations()
}

// 初始化
onMounted(() => {
  // 默认选择第一个选项
  if (targetOptions.value.length > 0) {
    selectedTarget.value = targetOptions.value[0].id
  }
})
</script>

<style scoped lang="scss">
.recommend-page {
  padding: 24px;
  background: var(--bg-secondary);
  min-height: calc(100vh - 60px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  .header-content {
    h1 {
      font-size: 28px;
      font-weight: 700;
      color: var(--text-primary);
      margin: 0 0 8px 0;
    }

    .subtitle {
      font-size: 14px;
      color: var(--text-secondary);
      margin: 0;
    }
  }
}

// 模式切换标签
.mode-tabs {
  display: flex;
  gap: 20px;
  margin-bottom: 24px;

  .mode-tab {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px 24px;
    background: var(--bg-primary);
    border-radius: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
    border: 2px solid transparent;
    box-shadow: var(--shadow-card);

    &:hover {
      transform: translateY(-2px);
      box-shadow: var(--shadow-lg);
    }

    &.active {
      border-color: var(--primary-color);
      background: linear-gradient(135deg, rgba(0, 184, 212, 0.05) 0%, rgba(0, 151, 167, 0.05) 100%);

      .tab-icon {
        background: linear-gradient(135deg, #00b8d4 0%, #0097a7 100%);
        color: white;
      }

      .tab-title {
        color: var(--primary-color);
      }
    }

    .tab-icon {
      width: 56px;
      height: 56px;
      border-radius: 14px;
      background: var(--bg-tertiary);
      display: flex;
      align-items: center;
      justify-content: center;
      transition: all 0.3s ease;

      .el-icon {
        font-size: 24px;
      }
    }

    .tab-info {
      display: flex;
      flex-direction: column;
      gap: 4px;

      .tab-title {
        font-size: 18px;
        font-weight: 600;
        color: var(--text-primary);
      }

      .tab-desc {
        font-size: 14px;
        color: var(--text-secondary);
      }
    }
  }
}

// 筛选区域
.filter-section {
  padding: 20px 24px;
  margin-bottom: 24px;

  .filter-row {
    display: flex;
    align-items: center;
    gap: 32px;
    flex-wrap: wrap;
  }

  .filter-item {
    display: flex;
    align-items: center;
    gap: 12px;

    label {
      font-size: 14px;
      color: #4b5563;
      white-space: nowrap;
    }
  }
}

.card {
  background: var(--bg-primary);
  border-radius: 16px;
  box-shadow: var(--shadow-card);
  border: 1px solid var(--border-light);
}

// 推荐结果
.recommendation-results {
  min-height: 400px;
}

.empty-state {
  background: white;
  border-radius: 16px;
  padding: 80px 20px;
}

.recommendation-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 24px;
}

.recommendation-card {
  background: var(--bg-primary);
  border-radius: 20px;
  padding: 24px;
  position: relative;
  box-shadow: var(--shadow-card);
  border: 1px solid var(--border-light);
  transition: all 0.3s ease;
  animation: cardFadeIn 0.5s ease forwards;
  opacity: 0;

  &:hover {
    transform: translateY(-4px);
    box-shadow: var(--shadow-lg);
    border-color: var(--primary-color);
  }

  // 匹配度圆环
  .match-score-ring {
    position: absolute;
    top: -15px;
    right: 20px;
    width: 80px;
    height: 80px;

    svg {
      transform: rotate(-90deg);
    }

    .bg-circle {
      stroke: var(--bg-tertiary);
    }

    .progress-circle {
      stroke: var(--primary-color);
      transition: stroke-dasharray 1s ease;
    }

    .score-value {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      text-align: center;

      .number {
        font-size: 24px;
        font-weight: 700;
        color: var(--text-primary);
      }

      .percent {
        font-size: 12px;
        color: var(--text-secondary);
      }
    }

    .score-label {
      position: absolute;
      bottom: -5px;
      left: 50%;
      transform: translateX(-50%);
      font-size: 11px;
      color: var(--text-muted);
      white-space: nowrap;
    }

    &.excellent .progress-circle { stroke: #00c853; }
    &.good .progress-circle { stroke: #00b8d4; }
    &.normal .progress-circle { stroke: #f59e0b; }
    &.low .progress-circle { stroke: #ef4444; }
  }
}

@keyframes cardFadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.card-content {
  .card-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
    padding-right: 70px;

    .avatar {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      overflow: hidden;
      flex-shrink: 0;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .company-logo {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      background: linear-gradient(135deg, #00b8d4 0%, #0097a7 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      flex-shrink: 0;

      .el-icon {
        font-size: 24px;
      }
    }

    .basic-info {
      flex: 1;
      min-width: 0;

      .name {
        font-size: 18px;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 4px 0;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .sub-info {
        font-size: 14px;
        color: var(--text-secondary);
        margin: 0;
      }
    }
  }

  .key-info {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin-bottom: 16px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--border-light);

    .info-item {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 13px;
      color: var(--text-secondary);

      .el-icon {
        color: var(--text-muted);
      }
    }
  }

  .skills {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;

    .el-tag {
      border-radius: 6px;

      .match-icon {
        margin-left: 4px;
        font-size: 12px;
      }
    }
  }

  .match-reasons {
    background: linear-gradient(135deg, rgba(0, 184, 212, 0.05) 0%, rgba(0, 151, 167, 0.05) 100%);
    border-radius: 12px;
    padding: 12px 16px;
    margin-bottom: 16px;

    .reason-title {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 13px;
      font-weight: 600;
      color: var(--primary-color);
      margin-bottom: 8px;
    }

    .reason-list {
      display: flex;
      flex-direction: column;
      gap: 6px;

      .reason-item {
        display: flex;
        align-items: flex-start;
        gap: 6px;
        font-size: 12px;
        color: #4b5563;

        .el-icon {
          color: #10b981;
          margin-top: 2px;
          flex-shrink: 0;
        }
      }
    }
  }

  .card-actions {
    display: flex;
    gap: 12px;

    .el-button {
      flex: 1;
      border-radius: 10px;

      &:last-child {
        flex: 0;
      }
    }
  }
}

// 分页
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  padding: 20px;
  background: var(--bg-primary);
  border-radius: 16px;
}

// 详情抽屉
.detail-drawer {
  padding: 0 20px;

  .detail-match-score {
    display: flex;
    justify-content: center;
    margin-bottom: 32px;

    .score-ring {
      position: relative;
      width: 160px;
      height: 160px;

      svg {
        transform: rotate(-90deg);
      }

      .bg-circle {
        stroke: var(--bg-tertiary);
      }

      .progress-circle {
        stroke: var(--primary-color);
        transition: stroke-dasharray 1s ease;
      }

      .score-content {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        text-align: center;

        .score-number {
          font-size: 48px;
          font-weight: 700;
          color: var(--text-primary);
        }

        .score-percent {
          font-size: 18px;
          color: var(--text-secondary);
        }

        .score-text {
          display: block;
          font-size: 14px;
          color: var(--text-muted);
          margin-top: 4px;
        }
      }
    }

    &.excellent .progress-circle { stroke: #00c853; }
    &.good .progress-circle { stroke: #00b8d4; }
    &.normal .progress-circle { stroke: #f59e0b; }
    &.low .progress-circle { stroke: #ef4444; }
  }

  .match-dimensions {
    margin-bottom: 32px;

    h4 {
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 16px 0;
    }

    .dimension-list {
      display: flex;
      flex-direction: column;
      gap: 16px;
    }

    .dimension-item {
      .dim-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 8px;

        .dim-name {
          font-size: 14px;
          color: var(--text-secondary);
        }

        .dim-score {
          font-size: 14px;
          font-weight: 600;
          color: var(--text-primary);
        }
      }
    }
  }

  .detail-info {
    margin-bottom: 32px;

    h4 {
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 16px 0;
    }
  }

  .recommend-reasons {
    margin-bottom: 32px;

    h4 {
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 16px 0;
    }

    .reasons-content {
      background: linear-gradient(135deg, rgba(0, 184, 212, 0.05) 0%, rgba(0, 151, 167, 0.05) 100%);
      border-radius: 12px;
      padding: 16px;
    }

    .reason-row {
      display: flex;
      align-items: flex-start;
      gap: 10px;
      padding: 8px 0;
      font-size: 14px;
      color: var(--text-secondary);

      .reason-icon {
        color: #00c853;
        margin-top: 2px;
        flex-shrink: 0;
      }

      &:not(:last-child) {
        border-bottom: 1px dashed var(--border-light);
      }
    }
  }

  .detail-actions {
    display: flex;
    gap: 16px;
    padding-top: 16px;
    border-top: 1px solid var(--border-light);

    .el-button {
      flex: 1;
      border-radius: 12px;
      height: 48px;
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .recommend-page {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .mode-tabs {
    flex-direction: column;
  }

  .filter-section {
    .filter-row {
      flex-direction: column;
      align-items: stretch;
      gap: 16px;
    }

    .filter-item {
      flex-direction: column;
      align-items: flex-start;

      .el-select, .el-slider {
        width: 100% !important;
      }
    }
  }

  .recommendation-grid {
    grid-template-columns: 1fr;
  }

  .recommendation-card {
    .match-score-ring {
      width: 70px;
      height: 70px;
      top: -10px;
      right: 15px;

      .score-value .number {
        font-size: 20px;
      }
    }
  }
}
</style>
