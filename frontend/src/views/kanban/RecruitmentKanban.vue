<template>
  <div class="recruitment-kanban">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1>招聘看板</h1>
        <p class="subtitle">可视化管理招聘流程，拖拽候选人卡片更新状态</p>
      </div>
      <div class="header-actions">
        <el-select v-model="selectedJob" placeholder="选择职位" style="width: 200px" clearable>
          <el-option
            v-for="job in jobs"
            :key="job.id"
            :label="job.title"
            :value="job.id"
          />
        </el-select>
        <el-button type="primary" @click="showAddCandidate = true">
          <el-icon><Plus /></el-icon>
          添加候选人
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <div
        v-for="stage in stages"
        :key="stage.id"
        class="stat-card"
        :style="{ borderTopColor: stage.color }"
      >
        <div class="stat-value">{{ getCandidateCount(stage.id) }}</div>
        <div class="stat-label">{{ stage.name }}</div>
      </div>
    </div>

    <!-- 看板区域 -->
    <div class="kanban-board">
      <div
        v-for="stage in stages"
        :key="stage.id"
        class="kanban-column"
        @dragover.prevent
        @drop="handleDrop($event, stage.id)"
      >
        <div class="column-header" :style="{ backgroundColor: stage.color + '15' }">
          <div class="column-title">
            <span class="color-dot" :style="{ backgroundColor: stage.color }"></span>
            <span>{{ stage.name }}</span>
          </div>
          <el-badge :value="getCandidateCount(stage.id)" type="info" />
        </div>

        <div class="column-content">
          <div
            v-for="candidate in getCandidatesByStage(stage.id)"
            :key="candidate.id"
            class="candidate-card"
            draggable="true"
            @dragstart="handleDragStart($event, candidate)"
            @dragend="handleDragEnd"
            :class="{ dragging: draggingId === candidate.id }"
          >
            <div class="card-header">
              <el-avatar :size="40" :src="candidate.avatar">
                {{ candidate.name.charAt(0) }}
              </el-avatar>
              <div class="candidate-info">
                <h4 class="candidate-name">{{ candidate.name }}</h4>
                <p class="candidate-position">{{ candidate.position }}</p>
              </div>
              <el-dropdown trigger="click" @command="(cmd: string) => handleCardAction(cmd, candidate)">
                <el-icon class="more-btn"><MoreFilled /></el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="view">查看详情</el-dropdown-item>
                    <el-dropdown-item command="schedule">安排面试</el-dropdown-item>
                    <el-dropdown-item command="notes">添加备注</el-dropdown-item>
                    <el-dropdown-item command="delete" divided>
                      <span style="color: #ef4444">移除候选人</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>

            <div class="card-body">
              <div class="card-tags">
                <el-tag
                  v-for="skill in candidate.skills.slice(0, 3)"
                  :key="skill"
                  size="small"
                  type="info"
                >
                  {{ skill }}
                </el-tag>
              </div>

              <div class="card-meta">
                <span class="meta-item">
                  <el-icon><Location /></el-icon>
                  {{ candidate.location }}
                </span>
                <span class="meta-item">
                  <el-icon><Suitcase /></el-icon>
                  {{ candidate.experience }}年
                </span>
              </div>

              <div class="card-footer">
                <span class="apply-time">
                  <el-icon><Clock /></el-icon>
                  {{ formatTime(candidate.applyTime) }}
                </span>
                <div class="match-score" :class="getScoreClass(candidate.matchScore)">
                  {{ candidate.matchScore }}%匹配
                </div>
              </div>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="getCandidateCount(stage.id) === 0" class="empty-column">
            <el-icon><Box /></el-icon>
            <span>暂无候选人</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 候选人详情抽屉 -->
    <el-drawer
      v-model="showCandidateDetail"
      :title="currentCandidate?.name + ' - 候选人详情'"
      size="500px"
    >
      <div class="candidate-detail" v-if="currentCandidate">
        <div class="detail-header">
          <el-avatar :size="80" :src="currentCandidate.avatar">
            {{ currentCandidate.name.charAt(0) }}
          </el-avatar>
          <div class="detail-info">
            <h2>{{ currentCandidate.name }}</h2>
            <p>{{ currentCandidate.position }}</p>
            <div class="detail-tags">
              <el-tag
                v-for="skill in currentCandidate.skills"
                :key="skill"
                size="small"
              >
                {{ skill }}
              </el-tag>
            </div>
          </div>
        </div>

        <el-divider />

        <el-descriptions :column="1" border>
          <el-descriptions-item label="当前阶段">
            <el-tag :color="getCurrentStage(currentCandidate.stage)?.color" effect="dark">
              {{ getCurrentStage(currentCandidate.stage)?.name }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="匹配度">
            <span :class="'score-' + getScoreClass(currentCandidate.matchScore)">
              {{ currentCandidate.matchScore }}%
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="工作经验">
            {{ currentCandidate.experience }}年
          </el-descriptions-item>
          <el-descriptions-item label="所在城市">
            {{ currentCandidate.location }}
          </el-descriptions-item>
          <el-descriptions-item label="期望薪资">
            {{ currentCandidate.salary }}
          </el-descriptions-item>
          <el-descriptions-item label="投递时间">
            {{ currentCandidate.applyTime }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="detail-actions">
          <el-button type="primary" @click="moveToNextStage(currentCandidate)">
            推进到下一阶段
          </el-button>
          <el-button @click="scheduleInterview(currentCandidate)">
            安排面试
          </el-button>
        </div>

        <!-- 操作历史 -->
        <div class="action-history">
          <h4>操作记录</h4>
          <el-timeline>
            <el-timeline-item
              v-for="(record, index) in currentCandidate.history"
              :key="index"
              :timestamp="record.time"
              placement="top"
            >
              {{ record.action }}
            </el-timeline-item>
          </el-timeline>
        </div>
      </div>
    </el-drawer>

    <!-- 添加候选人对话框 -->
    <el-dialog v-model="showAddCandidate" title="添加候选人" width="500px">
      <el-form :model="newCandidateForm" label-width="80px">
        <el-form-item label="姓名">
          <el-input v-model="newCandidateForm.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="职位">
          <el-input v-model="newCandidateForm.position" placeholder="请输入当前职位" />
        </el-form-item>
        <el-form-item label="城市">
          <el-input v-model="newCandidateForm.location" placeholder="请输入所在城市" />
        </el-form-item>
        <el-form-item label="经验">
          <el-input-number v-model="newCandidateForm.experience" :min="0" :max="30" />
          <span style="margin-left: 8px">年</span>
        </el-form-item>
        <el-form-item label="技能">
          <el-select
            v-model="newCandidateForm.skills"
            multiple
            filterable
            allow-create
            placeholder="请选择或输入技能"
            style="width: 100%"
          >
            <el-option label="Vue" value="Vue" />
            <el-option label="React" value="React" />
            <el-option label="TypeScript" value="TypeScript" />
            <el-option label="Node.js" value="Node.js" />
            <el-option label="Python" value="Python" />
            <el-option label="Java" value="Java" />
            <el-option label="Go" value="Go" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddCandidate = false">取消</el-button>
        <el-button type="primary" @click="addCandidate">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Plus, MoreFilled, Location, Suitcase, Clock, Box
} from '@element-plus/icons-vue'

// 类型定义
interface Candidate {
  id: number
  name: string
  position: string
  avatar?: string
  skills: string[]
  location: string
  experience: number
  salary: string
  matchScore: number
  stage: string
  applyTime: string
  history: { time: string; action: string }[]
}

interface Stage {
  id: string
  name: string
  color: string
}

// 招聘阶段
const stages = ref<Stage[]>([
  { id: 'applied', name: '简历投递', color: '#6366f1' },
  { id: 'screening', name: '简历筛选', color: '#8b5cf6' },
  { id: 'interview', name: '面试中', color: '#f59e0b' },
  { id: 'offer', name: '待发Offer', color: '#10b981' },
  { id: 'hired', name: '已录用', color: '#06b6d4' }
])

// 职位列表
const jobs = ref([
  { id: 1, title: '高级前端工程师' },
  { id: 2, title: '产品经理' },
  { id: 3, title: '后端工程师' },
  { id: 4, title: 'UI设计师' }
])

// 状态
const selectedJob = ref<number | null>(null)
const draggingId = ref<number | null>(null)
const showCandidateDetail = ref(false)
const showAddCandidate = ref(false)
const currentCandidate = ref<Candidate | null>(null)

// 新候选人表单
const newCandidateForm = reactive({
  name: '',
  position: '',
  location: '',
  experience: 3,
  skills: [] as string[]
})

// 候选人数据
const candidates = ref<Candidate[]>([
  {
    id: 1,
    name: '张伟',
    position: '高级前端工程师',
    skills: ['Vue', 'React', 'TypeScript'],
    location: '北京',
    experience: 5,
    salary: '30-40K',
    matchScore: 92,
    stage: 'interview',
    applyTime: '2024-01-10',
    history: [
      { time: '2024-01-10 10:00', action: '投递简历' },
      { time: '2024-01-11 14:00', action: '通过简历筛选' },
      { time: '2024-01-12 09:00', action: '安排技术面试' }
    ]
  },
  {
    id: 2,
    name: '李娜',
    position: '前端工程师',
    skills: ['Vue', 'JavaScript', 'CSS'],
    location: '上海',
    experience: 3,
    salary: '20-28K',
    matchScore: 85,
    stage: 'screening',
    applyTime: '2024-01-12',
    history: [
      { time: '2024-01-12 15:00', action: '投递简历' }
    ]
  },
  {
    id: 3,
    name: '王强',
    position: '全栈工程师',
    skills: ['Vue', 'Node.js', 'Python'],
    location: '深圳',
    experience: 4,
    salary: '25-35K',
    matchScore: 78,
    stage: 'applied',
    applyTime: '2024-01-13',
    history: [
      { time: '2024-01-13 09:30', action: '投递简历' }
    ]
  },
  {
    id: 4,
    name: '刘芳',
    position: '资深前端',
    skills: ['React', 'TypeScript', 'Webpack'],
    location: '杭州',
    experience: 6,
    salary: '35-45K',
    matchScore: 88,
    stage: 'offer',
    applyTime: '2024-01-08',
    history: [
      { time: '2024-01-08 11:00', action: '投递简历' },
      { time: '2024-01-09 10:00', action: '通过简历筛选' },
      { time: '2024-01-10 14:00', action: '完成技术面试' },
      { time: '2024-01-11 16:00', action: '完成HR面试' },
      { time: '2024-01-12 10:00', action: '进入Offer阶段' }
    ]
  },
  {
    id: 5,
    name: '陈明',
    position: '前端开发',
    skills: ['Vue', 'Element Plus', 'Git'],
    location: '成都',
    experience: 2,
    salary: '15-22K',
    matchScore: 72,
    stage: 'applied',
    applyTime: '2024-01-14',
    history: [
      { time: '2024-01-14 08:00', action: '投递简历' }
    ]
  },
  {
    id: 6,
    name: '赵雪',
    position: '中级前端',
    skills: ['React', 'Redux', 'Ant Design'],
    location: '广州',
    experience: 3,
    salary: '22-30K',
    matchScore: 80,
    stage: 'interview',
    applyTime: '2024-01-09',
    history: [
      { time: '2024-01-09 10:00', action: '投递简历' },
      { time: '2024-01-10 11:00', action: '通过简历筛选' },
      { time: '2024-01-11 15:00', action: '安排技术面试' }
    ]
  },
  {
    id: 7,
    name: '孙磊',
    position: '前端架构师',
    skills: ['Vue', 'React', '微前端', '性能优化'],
    location: '北京',
    experience: 8,
    salary: '45-60K',
    matchScore: 95,
    stage: 'hired',
    applyTime: '2024-01-05',
    history: [
      { time: '2024-01-05 09:00', action: '投递简历' },
      { time: '2024-01-05 14:00', action: '通过简历筛选' },
      { time: '2024-01-06 10:00', action: '完成技术面试' },
      { time: '2024-01-07 14:00', action: '完成HR面试' },
      { time: '2024-01-08 10:00', action: '发送Offer' },
      { time: '2024-01-09 10:00', action: '接受Offer，已录用' }
    ]
  }
])

// 获取阶段候选人
const getCandidatesByStage = (stageId: string) => {
  return candidates.value.filter(c => c.stage === stageId)
}

// 获取阶段候选人数量
const getCandidateCount = (stageId: string) => {
  return getCandidatesByStage(stageId).length
}

// 获取当前阶段信息
const getCurrentStage = (stageId: string) => {
  return stages.value.find(s => s.id === stageId)
}

// 获取匹配度样式类
const getScoreClass = (score: number) => {
  if (score >= 90) return 'excellent'
  if (score >= 75) return 'good'
  if (score >= 60) return 'normal'
  return 'low'
}

// 格式化时间
const formatTime = (time: string) => {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / 86400000)

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days}天前`
  return time
}

// 拖拽开始
const handleDragStart = (event: DragEvent, candidate: Candidate) => {
  draggingId.value = candidate.id
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('text/plain', String(candidate.id))
  }
}

// 拖拽结束
const handleDragEnd = () => {
  draggingId.value = null
}

// 放置
const handleDrop = (event: DragEvent, stageId: string) => {
  event.preventDefault()
  if (event.dataTransfer) {
    const candidateId = parseInt(event.dataTransfer.getData('text/plain'))
    const candidate = candidates.value.find(c => c.id === candidateId)
    if (candidate && candidate.stage !== stageId) {
      const oldStage = getCurrentStage(candidate.stage)?.name
      const newStage = getCurrentStage(stageId)?.name
      candidate.stage = stageId
      candidate.history.push({
        time: new Date().toLocaleString(),
        action: `从"${oldStage}"移动到"${newStage}"`
      })
      ElMessage.success(`已将 ${candidate.name} 移动到 ${newStage}`)
    }
  }
  draggingId.value = null
}

// 卡片操作
const handleCardAction = (command: string, candidate: Candidate) => {
  currentCandidate.value = candidate
  if (command === 'view') {
    showCandidateDetail.value = true
  } else if (command === 'schedule') {
    ElMessage.info('即将跳转到面试安排页面')
  } else if (command === 'notes') {
    ElMessage.info('备注功能开发中')
  } else if (command === 'delete') {
    candidates.value = candidates.value.filter(c => c.id !== candidate.id)
    ElMessage.success('已移除候选人')
  }
}

// 移动到下一阶段
const moveToNextStage = (candidate: Candidate) => {
  const currentIndex = stages.value.findIndex(s => s.id === candidate.stage)
  if (currentIndex < stages.value.length - 1) {
    const nextStage = stages.value[currentIndex + 1]
    candidate.stage = nextStage.id
    candidate.history.push({
      time: new Date().toLocaleString(),
      action: `推进到"${nextStage.name}"`
    })
    ElMessage.success(`已推进到 ${nextStage.name}`)
  } else {
    ElMessage.warning('已经是最后阶段')
  }
}

// 安排面试
const scheduleInterview = (candidate: Candidate) => {
  ElMessage.info(`为 ${candidate.name} 安排面试`)
}

// 添加候选人
const addCandidate = () => {
  if (!newCandidateForm.name) {
    ElMessage.warning('请输入候选人姓名')
    return
  }

  const newId = Math.max(...candidates.value.map(c => c.id)) + 1
  candidates.value.push({
    id: newId,
    name: newCandidateForm.name,
    position: newCandidateForm.position || '未知职位',
    skills: newCandidateForm.skills,
    location: newCandidateForm.location || '未知',
    experience: newCandidateForm.experience,
    salary: '面议',
    matchScore: Math.floor(Math.random() * 30) + 70,
    stage: 'applied',
    applyTime: new Date().toISOString().split('T')[0],
    history: [
      { time: new Date().toLocaleString(), action: '添加到候选人库' }
    ]
  })

  showAddCandidate.value = false
  newCandidateForm.name = ''
  newCandidateForm.position = ''
  newCandidateForm.location = ''
  newCandidateForm.experience = 3
  newCandidateForm.skills = []

  ElMessage.success('添加成功')
}
</script>

<style scoped lang="scss">
.recruitment-kanban {
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

  .header-actions {
    display: flex;
    gap: 12px;
  }
}

// 统计卡片
.stats-row {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  flex: 1;
  background: var(--bg-primary);
  border-radius: 12px;
  padding: 16px 20px;
  border-top: 3px solid;
  box-shadow: var(--shadow-card);

  .stat-value {
    font-size: 28px;
    font-weight: 700;
    color: var(--text-primary);
  }

  .stat-label {
    font-size: 14px;
    color: var(--text-secondary);
    margin-top: 4px;
  }
}

// 看板区域
.kanban-board {
  display: flex;
  gap: 16px;
  overflow-x: auto;
  padding-bottom: 16px;

  &::-webkit-scrollbar {
    height: 8px;
  }
}

.kanban-column {
  flex: 1;
  min-width: 280px;
  max-width: 320px;
  background: var(--bg-primary);
  border-radius: 16px;
  box-shadow: var(--shadow-card);
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 280px);

  .column-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    border-radius: 16px 16px 0 0;

    .column-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-weight: 600;
      color: var(--text-primary);

      .color-dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;
      }
    }
  }

  .column-content {
    flex: 1;
    padding: 12px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 12px;

    &::-webkit-scrollbar {
      width: 4px;
    }
  }
}

// 候选人卡片
.candidate-card {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 14px;
  cursor: grab;
  transition: all 0.3s ease;
  border: 2px solid transparent;

  &:hover {
    box-shadow: var(--shadow-md);
    border-color: var(--primary-color);
  }

  &.dragging {
    opacity: 0.5;
    transform: rotate(3deg);
  }

  &:active {
    cursor: grabbing;
  }

  .card-header {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    margin-bottom: 12px;

    .candidate-info {
      flex: 1;
      min-width: 0;

      .candidate-name {
        font-size: 15px;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 2px 0;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .candidate-position {
        font-size: 12px;
        color: var(--text-secondary);
        margin: 0;
      }
    }

    .more-btn {
      cursor: pointer;
      color: var(--text-tertiary);
      padding: 4px;

      &:hover {
        color: var(--primary-color);
      }
    }
  }

  .card-body {
    .card-tags {
      display: flex;
      flex-wrap: wrap;
      gap: 4px;
      margin-bottom: 10px;

      .el-tag {
        font-size: 11px;
        padding: 0 6px;
        height: 20px;
        line-height: 18px;
      }
    }

    .card-meta {
      display: flex;
      gap: 12px;
      font-size: 12px;
      color: var(--text-secondary);
      margin-bottom: 10px;

      .meta-item {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }

    .card-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .apply-time {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 11px;
        color: var(--text-tertiary);
      }

      .match-score {
        font-size: 11px;
        font-weight: 600;
        padding: 2px 8px;
        border-radius: 10px;

        &.excellent {
          background: rgba(16, 185, 129, 0.15);
          color: #10b981;
        }

        &.good {
          background: rgba(99, 102, 241, 0.15);
          color: #6366f1;
        }

        &.normal {
          background: rgba(245, 158, 11, 0.15);
          color: #f59e0b;
        }

        &.low {
          background: rgba(239, 68, 68, 0.15);
          color: #ef4444;
        }
      }
    }
  }
}

// 空状态
.empty-column {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: var(--text-tertiary);

  .el-icon {
    font-size: 32px;
    margin-bottom: 8px;
  }

  span {
    font-size: 13px;
  }
}

// 候选人详情
.candidate-detail {
  .detail-header {
    display: flex;
    gap: 20px;
    align-items: center;
    margin-bottom: 20px;

    .detail-info {
      h2 {
        font-size: 22px;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 4px 0;
      }

      p {
        font-size: 14px;
        color: var(--text-secondary);
        margin: 0 0 12px 0;
      }

      .detail-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 6px;
      }
    }
  }

  .detail-actions {
    display: flex;
    gap: 12px;
    margin: 24px 0;
  }

  .action-history {
    margin-top: 24px;

    h4 {
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 16px 0;
    }
  }
}

.score-excellent { color: #10b981; }
.score-good { color: #6366f1; }
.score-normal { color: #f59e0b; }
.score-low { color: #ef4444; }

// 响应式
@media (max-width: 768px) {
  .recruitment-kanban {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;

    .header-actions {
      width: 100%;
      flex-direction: column;
    }
  }

  .stats-row {
    flex-wrap: wrap;

    .stat-card {
      flex: 1 1 calc(50% - 8px);
      min-width: 140px;
    }
  }

  .kanban-column {
    min-width: 260px;
  }
}
</style>
