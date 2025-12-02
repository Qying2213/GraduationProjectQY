<template>
  <div class="interview-calendar">
    <!-- Header -->
    <div class="page-header">
      <div class="header-left">
        <h1>面试日历</h1>
        <p class="subtitle">管理和查看所有面试安排</p>
      </div>
      <div class="header-right">
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
        <el-button type="primary" @click="showAddDialog = true">
          <el-icon><Plus /></el-icon>
          安排面试
        </el-button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon today">
          <el-icon><Calendar /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ todayInterviews.length }}</span>
          <span class="stat-label">今日面试</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon week">
          <el-icon><DataLine /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ weekInterviews.length }}</span>
          <span class="stat-label">本周面试</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon pending">
          <el-icon><Clock /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ pendingInterviews.length }}</span>
          <span class="stat-label">待进行</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon completed">
          <el-icon><CircleCheck /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ completedInterviews.length }}</span>
          <span class="stat-label">已完成</span>
        </div>
      </div>
    </div>

    <!-- Calendar Container -->
    <div class="calendar-container">
      <!-- Calendar Header -->
      <div class="calendar-header">
        <div class="calendar-nav">
          <el-button-group>
            <el-button @click="prevMonth">
              <el-icon><ArrowLeft /></el-icon>
            </el-button>
            <el-button @click="goToToday">今天</el-button>
            <el-button @click="nextMonth">
              <el-icon><ArrowRight /></el-icon>
            </el-button>
          </el-button-group>
          <h2 class="current-month">{{ currentMonthLabel }}</h2>
        </div>
        <div class="view-switch">
          <el-radio-group v-model="viewMode" size="default">
            <el-radio-button value="month">月视图</el-radio-button>
            <el-radio-button value="week">周视图</el-radio-button>
            <el-radio-button value="list">列表</el-radio-button>
          </el-radio-group>
        </div>
      </div>

      <!-- Month View -->
      <div v-if="viewMode === 'month'" class="month-view">
        <div class="weekday-header">
          <div v-for="day in weekdays" :key="day" class="weekday">{{ day }}</div>
        </div>
        <div class="calendar-grid">
          <div
            v-for="(day, index) in calendarDays"
            :key="index"
            class="calendar-day"
            :class="{
              'other-month': !day.isCurrentMonth,
              'today': day.isToday,
              'has-events': getDayInterviews(day.date).length > 0
            }"
            @click="selectDate(day.date)"
          >
            <span class="day-number">{{ day.day }}</span>
            <div class="day-events">
              <div
                v-for="interview in getDayInterviews(day.date).slice(0, 3)"
                :key="interview.id"
                class="event-dot"
                :class="interview.type"
                :title="`${interview.time} ${interview.candidateName}`"
                @click.stop="showInterviewDetail(interview)"
              >
                <span class="event-time">{{ interview.time }}</span>
                <span class="event-name">{{ interview.candidateName }}</span>
              </div>
              <div
                v-if="getDayInterviews(day.date).length > 3"
                class="more-events"
              >
                +{{ getDayInterviews(day.date).length - 3 }} 更多
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Week View -->
      <div v-if="viewMode === 'week'" class="week-view">
        <div class="week-header">
          <div class="time-column"></div>
          <div
            v-for="day in weekDays"
            :key="day.date"
            class="week-day-header"
            :class="{ 'today': day.isToday }"
          >
            <span class="week-day-name">{{ day.dayName }}</span>
            <span class="week-day-number">{{ day.day }}</span>
          </div>
        </div>
        <div class="week-body">
          <div class="time-slots">
            <div v-for="hour in timeSlots" :key="hour" class="time-slot">
              <span class="time-label">{{ hour }}:00</span>
            </div>
          </div>
          <div class="week-grid">
            <div
              v-for="day in weekDays"
              :key="day.date"
              class="week-day-column"
            >
              <div
                v-for="interview in getDayInterviews(day.date)"
                :key="interview.id"
                class="week-event"
                :class="interview.type"
                :style="getEventStyle(interview)"
                @click="showInterviewDetail(interview)"
              >
                <span class="event-time">{{ interview.time }}</span>
                <span class="event-title">{{ interview.candidateName }}</span>
                <span class="event-position">{{ interview.position }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- List View -->
      <div v-if="viewMode === 'list'" class="list-view">
        <el-table :data="sortedInterviews" stripe style="width: 100%">
          <el-table-column prop="date" label="日期" width="120">
            <template #default="{ row }">
              {{ formatDate(row.date) }}
            </template>
          </el-table-column>
          <el-table-column prop="time" label="时间" width="100" />
          <el-table-column prop="candidateName" label="候选人" width="120" />
          <el-table-column prop="position" label="应聘职位" />
          <el-table-column prop="type" label="面试类型" width="120">
            <template #default="{ row }">
              <el-tag :type="getTypeTagType(row.type)" size="small">
                {{ getTypeLabel(row.type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="interviewer" label="面试官" width="120" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusTagType(row.status)" size="small">
                {{ getStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small" @click="showInterviewDetail(row)">
                查看
              </el-button>
              <el-button link type="warning" size="small" @click="editInterview(row)">
                编辑
              </el-button>
              <el-button link type="danger" size="small" @click="cancelInterview(row)">
                取消
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- Today's Schedule Sidebar -->
    <div class="today-sidebar">
      <h3>今日安排</h3>
      <div v-if="todayInterviews.length === 0" class="empty-today">
        <el-icon><Calendar /></el-icon>
        <p>今天没有面试安排</p>
      </div>
      <div v-else class="today-list">
        <div
          v-for="interview in todayInterviews"
          :key="interview.id"
          class="today-item"
          :class="interview.type"
          @click="showInterviewDetail(interview)"
        >
          <div class="today-time">{{ interview.time }}</div>
          <div class="today-info">
            <span class="today-name">{{ interview.candidateName }}</span>
            <span class="today-position">{{ interview.position }}</span>
            <span class="today-interviewer">
              <el-icon><User /></el-icon>
              {{ interview.interviewer }}
            </span>
          </div>
          <el-tag :type="getTypeTagType(interview.type)" size="small">
            {{ getTypeLabel(interview.type) }}
          </el-tag>
        </div>
      </div>
    </div>

    <!-- Add Interview Dialog -->
    <el-dialog
      v-model="showAddDialog"
      :title="isEditing ? '编辑面试' : '安排面试'"
      width="600px"
      destroy-on-close
    >
      <el-form :model="interviewForm" label-width="100px">
        <el-form-item label="候选人">
          <el-select v-model="interviewForm.candidateId" placeholder="选择候选人" style="width: 100%">
            <el-option
              v-for="candidate in candidates"
              :key="candidate.id"
              :label="candidate.name"
              :value="candidate.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="应聘职位">
          <el-select v-model="interviewForm.positionId" placeholder="选择职位" style="width: 100%">
            <el-option
              v-for="position in positions"
              :key="position.id"
              :label="position.name"
              :value="position.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="面试类型">
          <el-select v-model="interviewForm.type" placeholder="选择面试类型" style="width: 100%">
            <el-option label="初试" value="initial" />
            <el-option label="复试" value="second" />
            <el-option label="终面" value="final" />
            <el-option label="HR面试" value="hr" />
          </el-select>
        </el-form-item>
        <el-form-item label="面试日期">
          <el-date-picker
            v-model="interviewForm.date"
            type="date"
            placeholder="选择日期"
            style="width: 100%"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="面试时间">
          <el-time-picker
            v-model="interviewForm.time"
            placeholder="选择时间"
            style="width: 100%"
            format="HH:mm"
            value-format="HH:mm"
          />
        </el-form-item>
        <el-form-item label="时长(分钟)">
          <el-input-number v-model="interviewForm.duration" :min="30" :max="180" :step="30" />
        </el-form-item>
        <el-form-item label="面试官">
          <el-select v-model="interviewForm.interviewerId" placeholder="选择面试官" style="width: 100%">
            <el-option
              v-for="interviewer in interviewers"
              :key="interviewer.id"
              :label="interviewer.name"
              :value="interviewer.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="面试方式">
          <el-radio-group v-model="interviewForm.method">
            <el-radio value="onsite">现场面试</el-radio>
            <el-radio value="video">视频面试</el-radio>
            <el-radio value="phone">电话面试</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="会议室/链接">
          <el-input v-model="interviewForm.location" :placeholder="interviewForm.method === 'onsite' ? '输入会议室' : '输入会议链接'" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="interviewForm.notes" type="textarea" :rows="3" placeholder="面试注意事项..." />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveInterview">
          {{ isEditing ? '保存' : '确认安排' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Interview Detail Drawer -->
    <el-drawer
      v-model="showDetailDrawer"
      :title="'面试详情'"
      size="480px"
    >
      <div v-if="selectedInterview" class="interview-detail">
        <div class="detail-header">
          <el-avatar :size="64" class="candidate-avatar">
            {{ selectedInterview.candidateName.charAt(0) }}
          </el-avatar>
          <div class="candidate-info">
            <h3>{{ selectedInterview.candidateName }}</h3>
            <p>{{ selectedInterview.position }}</p>
          </div>
        </div>

        <el-divider />

        <div class="detail-section">
          <h4>面试信息</h4>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">日期</span>
              <span class="value">{{ formatDate(selectedInterview.date) }}</span>
            </div>
            <div class="info-item">
              <span class="label">时间</span>
              <span class="value">{{ selectedInterview.time }}</span>
            </div>
            <div class="info-item">
              <span class="label">类型</span>
              <el-tag :type="getTypeTagType(selectedInterview.type)" size="small">
                {{ getTypeLabel(selectedInterview.type) }}
              </el-tag>
            </div>
            <div class="info-item">
              <span class="label">状态</span>
              <el-tag :type="getStatusTagType(selectedInterview.status)" size="small">
                {{ getStatusLabel(selectedInterview.status) }}
              </el-tag>
            </div>
            <div class="info-item">
              <span class="label">面试官</span>
              <span class="value">{{ selectedInterview.interviewer }}</span>
            </div>
            <div class="info-item">
              <span class="label">面试方式</span>
              <span class="value">{{ getMethodLabel(selectedInterview.method) }}</span>
            </div>
            <div class="info-item full-width">
              <span class="label">地点/链接</span>
              <span class="value">{{ selectedInterview.location }}</span>
            </div>
          </div>
        </div>

        <el-divider />

        <div class="detail-section">
          <h4>备注</h4>
          <p class="notes">{{ selectedInterview.notes || '暂无备注' }}</p>
        </div>

        <div class="detail-actions">
          <el-button @click="editInterview(selectedInterview)">
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          <el-button type="danger" @click="cancelInterview(selectedInterview)">
            <el-icon><Close /></el-icon>
            取消面试
          </el-button>
          <el-button type="success" @click="completeInterview(selectedInterview)">
            <el-icon><Check /></el-icon>
            完成面试
          </el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Calendar, DataLine, Clock, CircleCheck,
  ArrowLeft, ArrowRight, User, Edit, Close, Check,
  Download, ArrowDown, Document, DocumentCopy
} from '@element-plus/icons-vue'
import { exportToExcel, exportToCsv, interviewExportColumns } from '@/utils/export'

interface Interview {
  id: number
  candidateId: number
  candidateName: string
  positionId: number
  position: string
  type: 'initial' | 'second' | 'final' | 'hr'
  date: string
  time: string
  duration: number
  interviewerId: number
  interviewer: string
  method: 'onsite' | 'video' | 'phone'
  location: string
  status: 'scheduled' | 'completed' | 'cancelled'
  notes: string
}

interface InterviewForm {
  candidateId: number | null
  positionId: number | null
  type: string
  date: string
  time: string
  duration: number
  interviewerId: number | null
  method: string
  location: string
  notes: string
}

const viewMode = ref<'month' | 'week' | 'list'>('month')
const currentDate = ref(new Date())
const showAddDialog = ref(false)
const showDetailDrawer = ref(false)
const selectedInterview = ref<Interview | null>(null)
const isEditing = ref(false)

const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
const timeSlots = [8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18]

// Mock data
const candidates = ref([
  { id: 1, name: '张三' },
  { id: 2, name: '李四' },
  { id: 3, name: '王五' },
  { id: 4, name: '赵六' },
  { id: 5, name: '钱七' }
])

const positions = ref([
  { id: 1, name: '高级前端工程师' },
  { id: 2, name: '后端开发工程师' },
  { id: 3, name: '产品经理' },
  { id: 4, name: 'UI设计师' },
  { id: 5, name: '数据分析师' }
])

const interviewers = ref([
  { id: 1, name: '刘经理' },
  { id: 2, name: '陈总监' },
  { id: 3, name: '周主管' },
  { id: 4, name: 'HR小李' }
])

const interviews = ref<Interview[]>([
  {
    id: 1,
    candidateId: 1,
    candidateName: '张三',
    positionId: 1,
    position: '高级前端工程师',
    type: 'initial',
    date: getTodayString(),
    time: '09:30',
    duration: 60,
    interviewerId: 1,
    interviewer: '刘经理',
    method: 'onsite',
    location: '3楼会议室A',
    status: 'scheduled',
    notes: '请准备技术面试问题，重点考察Vue3和TypeScript'
  },
  {
    id: 2,
    candidateId: 2,
    candidateName: '李四',
    positionId: 2,
    position: '后端开发工程师',
    type: 'second',
    date: getTodayString(),
    time: '14:00',
    duration: 90,
    interviewerId: 2,
    interviewer: '陈总监',
    method: 'video',
    location: 'https://meeting.example.com/xyz',
    status: 'scheduled',
    notes: '复试，请评估系统设计能力'
  },
  {
    id: 3,
    candidateId: 3,
    candidateName: '王五',
    positionId: 3,
    position: '产品经理',
    type: 'final',
    date: getDateString(1),
    time: '10:00',
    duration: 60,
    interviewerId: 2,
    interviewer: '陈总监',
    method: 'onsite',
    location: '5楼总监办公室',
    status: 'scheduled',
    notes: '终面，请准备offer谈判'
  },
  {
    id: 4,
    candidateId: 4,
    candidateName: '赵六',
    positionId: 4,
    position: 'UI设计师',
    type: 'hr',
    date: getDateString(2),
    time: '15:30',
    duration: 45,
    interviewerId: 4,
    interviewer: 'HR小李',
    method: 'phone',
    location: '010-12345678',
    status: 'scheduled',
    notes: 'HR面试，了解薪资期望和入职时间'
  },
  {
    id: 5,
    candidateId: 5,
    candidateName: '钱七',
    positionId: 5,
    position: '数据分析师',
    type: 'initial',
    date: getDateString(-1),
    time: '11:00',
    duration: 60,
    interviewerId: 3,
    interviewer: '周主管',
    method: 'video',
    location: 'https://meeting.example.com/abc',
    status: 'completed',
    notes: '初试已完成，表现优秀'
  }
])

const interviewForm = ref<InterviewForm>({
  candidateId: null,
  positionId: null,
  type: 'initial',
  date: '',
  time: '',
  duration: 60,
  interviewerId: null,
  method: 'onsite',
  location: '',
  notes: ''
})

// Helper functions
function getTodayString(): string {
  const today = new Date()
  return today.toISOString().split('T')[0]
}

function getDateString(offset: number): string {
  const date = new Date()
  date.setDate(date.getDate() + offset)
  return date.toISOString().split('T')[0]
}

// Computed
const currentMonthLabel = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth() + 1
  return `${year}年${month}月`
})

const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const days = []

  // Previous month days
  const startDay = firstDay.getDay()
  const prevMonth = new Date(year, month, 0)
  for (let i = startDay - 1; i >= 0; i--) {
    const day = prevMonth.getDate() - i
    const date = new Date(year, month - 1, day)
    days.push({
      day,
      date: date.toISOString().split('T')[0],
      isCurrentMonth: false,
      isToday: false
    })
  }

  // Current month days
  const today = new Date()
  for (let i = 1; i <= lastDay.getDate(); i++) {
    const date = new Date(year, month, i)
    const isToday = date.toDateString() === today.toDateString()
    days.push({
      day: i,
      date: date.toISOString().split('T')[0],
      isCurrentMonth: true,
      isToday
    })
  }

  // Next month days
  const remaining = 42 - days.length
  for (let i = 1; i <= remaining; i++) {
    const date = new Date(year, month + 1, i)
    days.push({
      day: i,
      date: date.toISOString().split('T')[0],
      isCurrentMonth: false,
      isToday: false
    })
  }

  return days
})

const weekDays = computed(() => {
  const today = new Date()
  const currentDay = today.getDay()
  const weekStart = new Date(today)
  weekStart.setDate(today.getDate() - currentDay)

  const days = []
  for (let i = 0; i < 7; i++) {
    const date = new Date(weekStart)
    date.setDate(weekStart.getDate() + i)
    days.push({
      date: date.toISOString().split('T')[0],
      day: date.getDate(),
      dayName: weekdays[i],
      isToday: date.toDateString() === today.toDateString()
    })
  }
  return days
})

const todayInterviews = computed(() => {
  const today = getTodayString()
  return interviews.value
    .filter(i => i.date === today && i.status === 'scheduled')
    .sort((a, b) => a.time.localeCompare(b.time))
})

const weekInterviews = computed(() => {
  const today = new Date()
  const weekStart = new Date(today)
  weekStart.setDate(today.getDate() - today.getDay())
  const weekEnd = new Date(weekStart)
  weekEnd.setDate(weekStart.getDate() + 6)

  return interviews.value.filter(i => {
    const date = new Date(i.date)
    return date >= weekStart && date <= weekEnd
  })
})

const pendingInterviews = computed(() => {
  return interviews.value.filter(i => i.status === 'scheduled')
})

const completedInterviews = computed(() => {
  return interviews.value.filter(i => i.status === 'completed')
})

const sortedInterviews = computed(() => {
  return [...interviews.value].sort((a, b) => {
    const dateCompare = a.date.localeCompare(b.date)
    if (dateCompare !== 0) return dateCompare
    return a.time.localeCompare(b.time)
  })
})

// Methods
function getDayInterviews(date: string): Interview[] {
  return interviews.value
    .filter(i => i.date === date && i.status !== 'cancelled')
    .sort((a, b) => a.time.localeCompare(b.time))
}

function getEventStyle(interview: Interview) {
  const hour = parseInt(interview.time.split(':')[0])
  const minute = parseInt(interview.time.split(':')[1])
  const top = (hour - 8) * 60 + minute
  const height = interview.duration

  return {
    top: `${top}px`,
    height: `${height}px`
  }
}

function prevMonth() {
  const date = new Date(currentDate.value)
  date.setMonth(date.getMonth() - 1)
  currentDate.value = date
}

function nextMonth() {
  const date = new Date(currentDate.value)
  date.setMonth(date.getMonth() + 1)
  currentDate.value = date
}

function goToToday() {
  currentDate.value = new Date()
}

function selectDate(date: string) {
  interviewForm.value.date = date
  showAddDialog.value = true
}

function showInterviewDetail(interview: Interview) {
  selectedInterview.value = interview
  showDetailDrawer.value = true
}

function editInterview(interview: Interview) {
  isEditing.value = true
  interviewForm.value = {
    candidateId: interview.candidateId,
    positionId: interview.positionId,
    type: interview.type,
    date: interview.date,
    time: interview.time,
    duration: interview.duration,
    interviewerId: interview.interviewerId,
    method: interview.method,
    location: interview.location,
    notes: interview.notes
  }
  showDetailDrawer.value = false
  showAddDialog.value = true
}

function saveInterview() {
  const candidate = candidates.value.find(c => c.id === interviewForm.value.candidateId)
  const position = positions.value.find(p => p.id === interviewForm.value.positionId)
  const interviewer = interviewers.value.find(i => i.id === interviewForm.value.interviewerId)

  if (!candidate || !position || !interviewer) {
    ElMessage.warning('请填写完整信息')
    return
  }

  if (isEditing.value && selectedInterview.value) {
    // Update existing
    const index = interviews.value.findIndex(i => i.id === selectedInterview.value!.id)
    if (index !== -1) {
      interviews.value[index] = {
        ...interviews.value[index],
        candidateId: interviewForm.value.candidateId!,
        candidateName: candidate.name,
        positionId: interviewForm.value.positionId!,
        position: position.name,
        type: interviewForm.value.type as Interview['type'],
        date: interviewForm.value.date,
        time: interviewForm.value.time,
        duration: interviewForm.value.duration,
        interviewerId: interviewForm.value.interviewerId!,
        interviewer: interviewer.name,
        method: interviewForm.value.method as Interview['method'],
        location: interviewForm.value.location,
        notes: interviewForm.value.notes
      }
      ElMessage.success('面试已更新')
    }
  } else {
    // Add new
    const newInterview: Interview = {
      id: Date.now(),
      candidateId: interviewForm.value.candidateId!,
      candidateName: candidate.name,
      positionId: interviewForm.value.positionId!,
      position: position.name,
      type: interviewForm.value.type as Interview['type'],
      date: interviewForm.value.date,
      time: interviewForm.value.time,
      duration: interviewForm.value.duration,
      interviewerId: interviewForm.value.interviewerId!,
      interviewer: interviewer.name,
      method: interviewForm.value.method as Interview['method'],
      location: interviewForm.value.location,
      status: 'scheduled',
      notes: interviewForm.value.notes
    }
    interviews.value.push(newInterview)
    ElMessage.success('面试已安排')
  }

  resetForm()
  showAddDialog.value = false
}

function cancelInterview(interview: Interview) {
  ElMessageBox.confirm('确定要取消这场面试吗？', '取消面试', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    const index = interviews.value.findIndex(i => i.id === interview.id)
    if (index !== -1) {
      interviews.value[index].status = 'cancelled'
      ElMessage.success('面试已取消')
      showDetailDrawer.value = false
    }
  }).catch(() => {})
}

function completeInterview(interview: Interview) {
  const index = interviews.value.findIndex(i => i.id === interview.id)
  if (index !== -1) {
    interviews.value[index].status = 'completed'
    ElMessage.success('面试已标记为完成')
    showDetailDrawer.value = false
  }
}

function resetForm() {
  isEditing.value = false
  interviewForm.value = {
    candidateId: null,
    positionId: null,
    type: 'initial',
    date: '',
    time: '',
    duration: 60,
    interviewerId: null,
    method: 'onsite',
    location: '',
    notes: ''
  }
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}月${date.getDate()}日`
}

function getTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    initial: '初试',
    second: '复试',
    final: '终面',
    hr: 'HR面'
  }
  return labels[type] || type
}

function getTypeTagType(type: string): 'primary' | 'success' | 'warning' | 'info' {
  const types: Record<string, 'primary' | 'success' | 'warning' | 'info'> = {
    initial: 'primary',
    second: 'warning',
    final: 'success',
    hr: 'info'
  }
  return types[type] || 'info'
}

function getStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    scheduled: '待进行',
    completed: '已完成',
    cancelled: '已取消'
  }
  return labels[status] || status
}

function getStatusTagType(status: string): 'primary' | 'success' | 'danger' | 'info' {
  const types: Record<string, 'primary' | 'success' | 'danger' | 'info'> = {
    scheduled: 'primary',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

function getMethodLabel(method: string): string {
  const labels: Record<string, string> = {
    onsite: '现场面试',
    video: '视频面试',
    phone: '电话面试'
  }
  return labels[method] || method
}

// 导出数据
function handleExport(format: 'excel' | 'csv') {
  if (interviews.value.length === 0) {
    ElMessage.warning('暂无数据可导出')
    return
  }

  const options = {
    filename: '面试安排',
    sheetName: '面试数据',
    columns: interviewExportColumns,
    data: interviews.value
  }

  if (format === 'excel') {
    exportToExcel(options)
    ElMessage.success('Excel 导出成功')
  } else {
    exportToCsv(options)
    ElMessage.success('CSV 导出成功')
  }
}
</script>

<style scoped lang="scss">
.interview-calendar {
  display: grid;
  grid-template-columns: 1fr 320px;
  grid-template-rows: auto auto 1fr;
  gap: 24px;
  min-height: calc(100vh - 112px);
}

.page-header {
  grid-column: 1 / -1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--card-bg);
  padding: 24px;
  border-radius: 16px;
  box-shadow: var(--card-shadow);

  h1 {
    margin: 0;
    font-size: 24px;
    color: var(--text-primary);
  }

  .subtitle {
    margin: 4px 0 0;
    color: var(--text-secondary);
    font-size: 14px;
  }
}

.stats-row {
  grid-column: 1 / -1;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;

  .stat-card {
    background: var(--card-bg);
    padding: 20px;
    border-radius: 12px;
    box-shadow: var(--card-shadow);
    display: flex;
    align-items: center;
    gap: 16px;

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;

      .el-icon {
        font-size: 24px;
        color: white;
      }

      &.today { background: linear-gradient(135deg, #6366f1, #8b5cf6); }
      &.week { background: linear-gradient(135deg, #3b82f6, #06b6d4); }
      &.pending { background: linear-gradient(135deg, #f59e0b, #f97316); }
      &.completed { background: linear-gradient(135deg, #10b981, #059669); }
    }

    .stat-info {
      display: flex;
      flex-direction: column;

      .stat-value {
        font-size: 28px;
        font-weight: 700;
        color: var(--text-primary);
      }

      .stat-label {
        font-size: 13px;
        color: var(--text-secondary);
      }
    }
  }
}

.calendar-container {
  background: var(--card-bg);
  border-radius: 16px;
  box-shadow: var(--card-shadow);
  padding: 24px;
  overflow: hidden;

  .calendar-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    .calendar-nav {
      display: flex;
      align-items: center;
      gap: 16px;

      .current-month {
        margin: 0;
        font-size: 20px;
        font-weight: 600;
        color: var(--text-primary);
      }
    }
  }
}

.month-view {
  .weekday-header {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    margin-bottom: 8px;

    .weekday {
      text-align: center;
      padding: 12px;
      font-weight: 600;
      color: var(--text-secondary);
      font-size: 13px;
    }
  }

  .calendar-grid {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    overflow: hidden;

    .calendar-day {
      min-height: 100px;
      padding: 8px;
      border-right: 1px solid var(--border-color);
      border-bottom: 1px solid var(--border-color);
      background: var(--card-bg);
      cursor: pointer;
      transition: background-color 0.2s;

      &:nth-child(7n) { border-right: none; }
      &:nth-last-child(-n+7) { border-bottom: none; }

      &:hover {
        background: var(--bg-tertiary);
      }

      &.other-month {
        background: var(--bg-secondary);

        .day-number { color: var(--text-muted); }
      }

      &.today {
        .day-number {
          background: var(--primary-color);
          color: white;
          border-radius: 50%;
          width: 28px;
          height: 28px;
          display: flex;
          align-items: center;
          justify-content: center;
        }
      }

      .day-number {
        font-weight: 500;
        color: var(--text-primary);
        margin-bottom: 4px;
      }

      .day-events {
        display: flex;
        flex-direction: column;
        gap: 2px;

        .event-dot {
          padding: 2px 6px;
          border-radius: 4px;
          font-size: 11px;
          cursor: pointer;
          display: flex;
          gap: 4px;
          overflow: hidden;

          &.initial { background: rgba(99, 102, 241, 0.15); color: #6366f1; }
          &.second { background: rgba(245, 158, 11, 0.15); color: #f59e0b; }
          &.final { background: rgba(16, 185, 129, 0.15); color: #10b981; }
          &.hr { background: rgba(107, 114, 128, 0.15); color: #6b7280; }

          .event-time { font-weight: 500; }
          .event-name {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
        }

        .more-events {
          font-size: 11px;
          color: var(--text-secondary);
          padding: 2px 6px;
        }
      }
    }
  }
}

.week-view {
  .week-header {
    display: grid;
    grid-template-columns: 60px repeat(7, 1fr);
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 12px;

    .time-column { }

    .week-day-header {
      text-align: center;

      &.today {
        .week-day-number {
          background: var(--primary-color);
          color: white;
          border-radius: 50%;
          width: 32px;
          height: 32px;
          display: inline-flex;
          align-items: center;
          justify-content: center;
        }
      }

      .week-day-name {
        display: block;
        font-size: 12px;
        color: var(--text-secondary);
      }

      .week-day-number {
        font-size: 18px;
        font-weight: 600;
        color: var(--text-primary);
      }
    }
  }

  .week-body {
    display: grid;
    grid-template-columns: 60px repeat(7, 1fr);
    position: relative;
    height: 660px;
    overflow-y: auto;

    .time-slots {
      .time-slot {
        height: 60px;
        border-bottom: 1px solid var(--border-light);

        .time-label {
          font-size: 11px;
          color: var(--text-secondary);
          padding: 4px;
        }
      }
    }

    .week-grid {
      display: contents;

      .week-day-column {
        border-left: 1px solid var(--border-light);
        position: relative;
        background:
          repeating-linear-gradient(
            to bottom,
            transparent 0px,
            transparent 59px,
            var(--border-light) 59px,
            var(--border-light) 60px
          );

        .week-event {
          position: absolute;
          left: 4px;
          right: 4px;
          border-radius: 6px;
          padding: 4px 8px;
          cursor: pointer;
          overflow: hidden;
          z-index: 1;

          &.initial { background: rgba(99, 102, 241, 0.9); color: white; }
          &.second { background: rgba(245, 158, 11, 0.9); color: white; }
          &.final { background: rgba(16, 185, 129, 0.9); color: white; }
          &.hr { background: rgba(107, 114, 128, 0.9); color: white; }

          .event-time {
            font-size: 11px;
            font-weight: 500;
          }

          .event-title {
            display: block;
            font-size: 12px;
            font-weight: 600;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .event-position {
            font-size: 10px;
            opacity: 0.9;
          }
        }
      }
    }
  }
}

.list-view {
  :deep(.el-table) {
    background: transparent;

    th.el-table__cell {
      background: var(--bg-tertiary);
    }

    tr {
      background: var(--card-bg);

      &:hover > td.el-table__cell {
        background: var(--bg-tertiary);
      }
    }
  }
}

.today-sidebar {
  background: var(--card-bg);
  border-radius: 16px;
  box-shadow: var(--card-shadow);
  padding: 24px;
  height: fit-content;

  h3 {
    margin: 0 0 16px;
    font-size: 16px;
    color: var(--text-primary);
  }

  .empty-today {
    text-align: center;
    padding: 40px 0;
    color: var(--text-secondary);

    .el-icon {
      font-size: 48px;
      margin-bottom: 12px;
      color: var(--text-muted);
    }

    p { margin: 0; }
  }

  .today-list {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .today-item {
      padding: 16px;
      border-radius: 12px;
      cursor: pointer;
      transition: transform 0.2s;
      border-left: 4px solid;

      &:hover { transform: translateX(4px); }

      &.initial { background: rgba(99, 102, 241, 0.1); border-color: #6366f1; }
      &.second { background: rgba(245, 158, 11, 0.1); border-color: #f59e0b; }
      &.final { background: rgba(16, 185, 129, 0.1); border-color: #10b981; }
      &.hr { background: rgba(107, 114, 128, 0.1); border-color: #6b7280; }

      .today-time {
        font-size: 13px;
        font-weight: 600;
        color: var(--text-primary);
        margin-bottom: 8px;
      }

      .today-info {
        display: flex;
        flex-direction: column;
        gap: 4px;
        margin-bottom: 8px;

        .today-name {
          font-weight: 600;
          color: var(--text-primary);
        }

        .today-position {
          font-size: 13px;
          color: var(--text-secondary);
        }

        .today-interviewer {
          font-size: 12px;
          color: var(--text-muted);
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }
  }
}

.interview-detail {
  .detail-header {
    display: flex;
    align-items: center;
    gap: 16px;

    .candidate-avatar {
      background: var(--gradient-1);
      font-size: 24px;
    }

    .candidate-info {
      h3 {
        margin: 0;
        font-size: 20px;
        color: var(--text-primary);
      }

      p {
        margin: 4px 0 0;
        color: var(--text-secondary);
      }
    }
  }

  .detail-section {
    h4 {
      margin: 0 0 12px;
      font-size: 14px;
      color: var(--text-secondary);
    }

    .info-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 16px;

      .info-item {
        display: flex;
        flex-direction: column;
        gap: 4px;

        &.full-width { grid-column: 1 / -1; }

        .label {
          font-size: 12px;
          color: var(--text-muted);
        }

        .value {
          font-size: 14px;
          color: var(--text-primary);
        }
      }
    }

    .notes {
      color: var(--text-secondary);
      font-size: 14px;
      line-height: 1.6;
    }
  }

  .detail-actions {
    margin-top: 24px;
    display: flex;
    gap: 12px;
  }
}

@media (max-width: 1200px) {
  .interview-calendar {
    grid-template-columns: 1fr;
  }

  .today-sidebar {
    order: -1;
  }

  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-row {
    grid-template-columns: 1fr;
  }

  .page-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }

  .calendar-header {
    flex-direction: column;
    gap: 16px;
  }
}
</style>
