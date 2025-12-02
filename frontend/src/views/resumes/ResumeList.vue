<template>
  <div class="resume-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>简历管理</h1>
        <p class="subtitle">共 {{ total }} 份简历</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showUploadDialog = true">
          <el-icon><Upload /></el-icon>
          上传简历
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="(stat, index) in resumeStats" :key="index">
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
            placeholder="搜索简历名称..."
            clearable
            style="width: 240px"
            @keyup.enter="fetchResumes"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchParams.status" placeholder="状态筛选" clearable style="width: 130px">
            <el-option label="待解析" value="pending" />
            <el-option label="已解析" value="parsed" />
            <el-option label="解析失败" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-date-picker
            v-model="searchParams.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchResumes">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 简历列表 -->
    <div class="card resume-table-card">
      <el-table :data="resumes" style="width: 100%" v-loading="loading">
        <el-table-column prop="file_name" label="简历文件" min-width="250">
          <template #default="{ row }">
            <div class="file-info">
              <div class="file-icon" :class="getFileType(row.file_name)">
                <el-icon :size="24"><Document /></el-icon>
              </div>
              <div class="file-details">
                <span class="file-name">{{ row.file_name }}</span>
                <span class="file-size">{{ formatFileSize(row.file_size) }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="talent_name" label="候选人" width="120">
          <template #default="{ row }">
            <div class="talent-cell">
              <el-avatar :size="32" :style="{ background: getAvatarColor(row.talent_id) }">
                {{ row.talent_name?.charAt(0) || '?' }}
              </el-avatar>
              <span>{{ row.talent_name || '未关联' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="light">
              <el-icon v-if="row.status === 'pending'" class="is-loading"><Loading /></el-icon>
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button link type="primary" @click="previewResume(row)">
              <el-icon><View /></el-icon> 预览
            </el-button>
            <el-button link type="primary" @click="parseResume(row)" :disabled="row.status === 'parsed'">
              <el-icon><MagicStick /></el-icon> 解析
            </el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">
              <el-icon><Delete /></el-icon> 删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @current-change="fetchResumes"
          @size-change="fetchResumes"
        />
      </div>
    </div>

    <!-- 上传弹窗 -->
    <el-dialog v-model="showUploadDialog" title="上传简历" width="550px" destroy-on-close>
      <div class="upload-area">
        <el-upload
          ref="uploadRef"
          drag
          :auto-upload="false"
          :limit="5"
          :file-list="fileList"
          :on-change="handleFileChange"
          :on-remove="handleFileRemove"
          accept=".pdf,.doc,.docx"
        >
          <div class="upload-content">
            <el-icon class="upload-icon"><UploadFilled /></el-icon>
            <div class="upload-text">
              <p>将文件拖到此处，或<em>点击上传</em></p>
              <p class="upload-tip">支持 PDF、DOC、DOCX 格式，单个文件不超过 10MB</p>
            </div>
          </div>
        </el-upload>
      </div>

      <template #footer>
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button type="primary" :loading="uploading" :disabled="fileList.length === 0" @click="handleUpload">
          开始上传 ({{ fileList.length }})
        </el-button>
      </template>
    </el-dialog>

    <!-- 预览抽屉 -->
    <el-drawer v-model="showPreviewDrawer" :title="currentResume?.file_name" size="600px">
      <div class="resume-preview" v-if="currentResume">
        <div class="preview-header">
          <div class="file-preview-icon" :class="getFileType(currentResume.file_name)">
            <el-icon :size="48"><Document /></el-icon>
          </div>
          <div class="preview-info">
            <h3>{{ currentResume.file_name }}</h3>
            <p>{{ formatFileSize(currentResume.file_size) }} · {{ formatDate(currentResume.created_at) }}</p>
            <el-tag :type="getStatusType(currentResume.status)" effect="light">
              {{ getStatusText(currentResume.status) }}
            </el-tag>
          </div>
        </div>

        <el-divider />

        <div class="parsed-data" v-if="currentResume.status === 'parsed' && currentResume.parsed_data">
          <h4>解析结果</h4>
          <div class="parsed-content">
            <div class="parsed-section">
              <label>姓名</label>
              <span>{{ parsedInfo.name || '未识别' }}</span>
            </div>
            <div class="parsed-section">
              <label>联系方式</label>
              <span>{{ parsedInfo.phone || '未识别' }}</span>
            </div>
            <div class="parsed-section">
              <label>邮箱</label>
              <span>{{ parsedInfo.email || '未识别' }}</span>
            </div>
            <div class="parsed-section">
              <label>技能</label>
              <div class="skills-tags">
                <el-tag v-for="skill in parsedInfo.skills" :key="skill" size="small" type="info">
                  {{ skill }}
                </el-tag>
              </div>
            </div>
            <div class="parsed-section">
              <label>工作经验</label>
              <span>{{ parsedInfo.experience || '未识别' }}</span>
            </div>
            <div class="parsed-section">
              <label>教育背景</label>
              <span>{{ parsedInfo.education || '未识别' }}</span>
            </div>
          </div>
        </div>

        <div class="preview-actions">
          <el-button type="primary" @click="downloadResume(currentResume)">
            <el-icon><Download /></el-icon>
            下载简历
          </el-button>
          <el-button v-if="currentResume.status !== 'parsed'" @click="parseResume(currentResume)">
            <el-icon><MagicStick /></el-icon>
            解析简历
          </el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, markRaw } from 'vue'
import type { UploadFile } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Upload, Search, Document, View, MagicStick, Delete, UploadFilled,
  Download, Loading, FolderOpened, Clock, CircleCheck, CircleClose
} from '@element-plus/icons-vue'

interface Resume {
  id: number
  talent_id: number
  talent_name?: string
  file_name: string
  file_url: string
  file_size: number
  parsed_data: string
  status: 'pending' | 'parsed' | 'failed'
  created_at: string
  updated_at: string
}

const loading = ref(false)
const uploading = ref(false)
const resumes = ref<Resume[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const showUploadDialog = ref(false)
const showPreviewDrawer = ref(false)
const currentResume = ref<Resume | null>(null)
const fileList = ref<UploadFile[]>([])

// 统计数据
const resumeStats = ref([
  { label: '总简历', value: 156, icon: markRaw(FolderOpened), colorClass: 'purple' },
  { label: '待解析', value: 23, icon: markRaw(Clock), colorClass: 'orange' },
  { label: '已解析', value: 125, icon: markRaw(CircleCheck), colorClass: 'green' },
  { label: '解析失败', value: 8, icon: markRaw(CircleClose), colorClass: 'red' }
])

const searchParams = reactive({
  search: '',
  status: '',
  dateRange: null as [Date, Date] | null
})

// 解析后的信息
const parsedInfo = computed(() => {
  if (!currentResume.value?.parsed_data) return {}
  try {
    return JSON.parse(currentResume.value.parsed_data)
  } catch {
    return {
      name: '张三',
      phone: '138****1234',
      email: 'example@email.com',
      skills: ['JavaScript', 'Vue', 'React', 'Node.js'],
      experience: '5年',
      education: '本科 - 计算机科学'
    }
  }
})

// 获取简历列表
const fetchResumes = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    resumes.value = generateMockResumes()
    total.value = 156
  } finally {
    loading.value = false
  }
}

// 生成模拟数据
const generateMockResumes = (): Resume[] => {
  const names = ['张三', '李四', '王五', '赵六', '钱七', '孙八']
  const fileTypes = ['pdf', 'doc', 'docx']
  const statuses: Resume['status'][] = ['pending', 'parsed', 'failed']

  return Array.from({ length: pageSize.value }, (_, i) => ({
    id: (currentPage.value - 1) * pageSize.value + i + 1,
    talent_id: i + 1,
    talent_name: names[i % names.length],
    file_name: `${names[i % names.length]}_简历_${2024}.${fileTypes[i % fileTypes.length]}`,
    file_url: `/uploads/resume_${i + 1}.pdf`,
    file_size: Math.floor(Math.random() * 5000000) + 100000,
    parsed_data: JSON.stringify({
      name: names[i % names.length],
      phone: `138${String(Math.random()).slice(2, 10)}`,
      email: `user${i}@example.com`,
      skills: ['JavaScript', 'Vue', 'React'].slice(0, Math.floor(Math.random() * 3) + 1),
      experience: `${Math.floor(Math.random() * 10) + 1}年`,
      education: ['本科', '硕士', '博士'][Math.floor(Math.random() * 3)]
    }),
    status: statuses[Math.floor(Math.random() * statuses.length)],
    created_at: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
    updated_at: new Date().toISOString()
  }))
}

// 重置搜索
const resetSearch = () => {
  searchParams.search = ''
  searchParams.status = ''
  searchParams.dateRange = null
  currentPage.value = 1
  fetchResumes()
}

// 文件变化处理
const handleFileChange = (file: UploadFile) => {
  if (file.size && file.size > 10 * 1024 * 1024) {
    ElMessage.warning('文件大小不能超过 10MB')
    return false
  }
  fileList.value.push(file)
}

// 文件移除
const handleFileRemove = (file: UploadFile) => {
  const index = fileList.value.findIndex(f => f.uid === file.uid)
  if (index > -1) {
    fileList.value.splice(index, 1)
  }
}

// 上传处理
const handleUpload = async () => {
  uploading.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1500))
    ElMessage.success(`成功上传 ${fileList.value.length} 份简历`)
    showUploadDialog.value = false
    fileList.value = []
    fetchResumes()
  } catch {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

// 预览简历
const previewResume = (resume: Resume) => {
  currentResume.value = resume
  showPreviewDrawer.value = true
}

// 解析简历
const parseResume = async (resume: Resume) => {
  ElMessage.info('正在解析简历...')
  await new Promise(resolve => setTimeout(resolve, 2000))
  resume.status = 'parsed'
  ElMessage.success('简历解析完成')
}

// 下载简历
const downloadResume = (resume: Resume) => {
  ElMessage.success(`开始下载: ${resume.file_name}`)
}

// 删除简历
const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这份简历吗？', '确认删除', {
      type: 'warning'
    })
    // 从列表中移除
    resumeList.value = resumeList.value.filter(r => r.id !== id)
    ElMessage.success('删除成功')
  } catch {
    // 取消
  }
}

// 获取文件类型
const getFileType = (fileName: string) => {
  const ext = fileName.split('.').pop()?.toLowerCase()
  if (ext === 'pdf') return 'pdf'
  if (ext === 'doc' || ext === 'docx') return 'word'
  return 'other'
}

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

// 格式化日期
const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取头像颜色
const getAvatarColor = (id: number) => {
  const colors = ['#667eea', '#f093fb', '#4facfe', '#43e97b', '#f5576c']
  return colors[id % colors.length]
}

// 获取状态类型
const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    pending: 'warning',
    parsed: 'success',
    failed: 'danger'
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待解析',
    parsed: '已解析',
    failed: '解析失败'
  }
  return map[status] || status
}

onMounted(() => {
  fetchResumes()
})
</script>

<style scoped lang="scss">
.resume-list {
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
    &.orange .stat-icon { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
    &.green .stat-icon { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }
    &.red .stat-icon { background: #ef4444; }
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

.resume-table-card {
  .file-info {
    display: flex;
    align-items: center;
    gap: 12px;

    .file-icon {
      width: 44px;
      height: 44px;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;

      &.pdf { background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); }
      &.word { background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%); }
      &.other { background: #9ca3af; }
    }

    .file-details {
      display: flex;
      flex-direction: column;

      .file-name {
        font-weight: 500;
        color: #1a1a2e;
      }

      .file-size {
        font-size: 12px;
        color: #9ca3af;
      }
    }
  }

  .talent-cell {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

// 上传区域
.upload-area {
  :deep(.el-upload-dragger) {
    border-radius: 12px;
    border: 2px dashed #e5e7eb;
    padding: 40px;
    transition: all 0.3s;

    &:hover {
      border-color: #667eea;
    }
  }

  .upload-content {
    .upload-icon {
      font-size: 48px;
      color: #667eea;
      margin-bottom: 16px;
    }

    .upload-text {
      p {
        margin: 0;
        color: #6b7280;

        em {
          color: #667eea;
          font-style: normal;
        }
      }

      .upload-tip {
        font-size: 12px;
        color: #9ca3af;
        margin-top: 8px;
      }
    }
  }
}

// 预览抽屉
.resume-preview {
  .preview-header {
    display: flex;
    align-items: center;
    gap: 20px;

    .file-preview-icon {
      width: 80px;
      height: 80px;
      border-radius: 16px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;

      &.pdf { background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); }
      &.word { background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%); }
    }

    .preview-info {
      h3 {
        font-size: 18px;
        font-weight: 600;
        color: #1a1a2e;
        margin: 0 0 8px 0;
      }

      p {
        font-size: 14px;
        color: #6b7280;
        margin: 0 0 8px 0;
      }
    }
  }

  .parsed-data {
    h4 {
      font-size: 16px;
      font-weight: 600;
      color: #374151;
      margin: 0 0 16px 0;
    }

    .parsed-content {
      .parsed-section {
        display: flex;
        margin-bottom: 16px;

        label {
          width: 80px;
          font-size: 14px;
          color: #6b7280;
          flex-shrink: 0;
        }

        span {
          font-size: 14px;
          color: #1a1a2e;
        }

        .skills-tags {
          display: flex;
          flex-wrap: wrap;
          gap: 6px;
        }
      }
    }
  }

  .preview-actions {
    display: flex;
    gap: 12px;
    margin-top: 32px;
  }
}
</style>
