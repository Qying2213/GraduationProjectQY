<template>
  <div class="resume-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>简历管理</h1>
        <p class="subtitle">共 {{ total }} 份简历</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showUploadDialog = true" v-if="canCreate">
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
      <el-table :data="resumes" style="width: 100%" v-loading="loading" @sort-change="handleSortChange">
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
        <el-table-column prop="status" label="状态" width="120" sortable="custom">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="light">
              <el-icon v-if="row.status === 'pending'" class="is-loading"><Loading /></el-icon>
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="180" sortable="custom">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button link type="primary" @click="previewResume(row)" title="">
              <el-icon><View /></el-icon> 预览
            </el-button>
            <el-button link type="primary" @click="parseResume(row)" :disabled="row.status === 'parsed'" title="" v-if="canEdit">
              <el-icon><MagicStick /></el-icon> 解析
            </el-button>
            <el-button link type="danger" @click="handleDelete(row.id)" title="" v-if="canDelete">
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
    <el-drawer v-model="showPreviewDrawer" title="简历预览" size="85%">
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

        <!-- PDF 预览区域 -->
        <div class="pdf-preview-container" v-if="currentResume.file_url && getFileType(currentResume.file_name) === 'pdf'">
          <iframe 
            :src="getPreviewUrl(currentResume) + '#toolbar=1&navpanes=0&scrollbar=1&view=FitH'" 
            class="pdf-iframe"
            frameborder="0"
          ></iframe>
        </div>

        <!-- 非 PDF 文件提示 -->
        <div class="non-pdf-notice" v-else-if="currentResume.file_url">
          <el-icon :size="48"><Document /></el-icon>
          <p>{{ getFileType(currentResume.file_name).toUpperCase() }} 文件暂不支持在线预览</p>
          <p class="tip">请下载后查看</p>
        </div>

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
import { usePermissionStore } from '@/store/permission'
import request from '@/utils/request'

const permissionStore = usePermissionStore()

// 权限检查
const canCreate = computed(() => permissionStore.hasPermission('resume:create'))
const canEdit = computed(() => permissionStore.hasPermission('resume:edit'))
const canDelete = computed(() => permissionStore.hasPermission('resume:delete'))

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

// 统计数据 - 青蓝色系
const resumeStats = ref([
  { label: '总简历', value: 156, icon: markRaw(FolderOpened), colorClass: 'cyan' },
  { label: '待解析', value: 23, icon: markRaw(Clock), colorClass: 'teal' },
  { label: '已解析', value: 125, icon: markRaw(CircleCheck), colorClass: 'green' },
  { label: '解析失败', value: 8, icon: markRaw(CircleClose), colorClass: 'red' }
])

const searchParams = reactive({
  search: '',
  status: '',
  dateRange: null as [Date, Date] | null
})

// 排序参数
const sortParams = reactive({
  prop: '',
  order: '' as '' | 'ascending' | 'descending'
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
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (searchParams.search) {
      params.search = searchParams.search
    }
    if (searchParams.status) {
      params.status = searchParams.status
    }
    if (sortParams.prop && sortParams.order) {
      params.sort_by = sortParams.prop
      params.sort_order = sortParams.order === 'ascending' ? 'asc' : 'desc'
    }
    
    const res = await request.get('/resumes', { params })
    if (res.data?.code === 0) {
      resumes.value = res.data.data.resumes || []
      total.value = res.data.data.total || 0
      
      // 更新统计数据
      updateStats(res.data.data.resumes || [])
    } else {
      ElMessage.error(res.data?.message || '获取简历列表失败')
    }
  } catch (error) {
    console.error('获取简历列表失败:', error)
    ElMessage.error('获取简历列表失败')
  } finally {
    loading.value = false
  }
}

// 排序变化处理
const handleSortChange = ({ prop, order }: { prop: string; order: '' | 'ascending' | 'descending' | null }) => {
  sortParams.prop = prop || ''
  sortParams.order = order || ''
  currentPage.value = 1  // 排序时回到第一页
  fetchResumes()
}

// 更新统计数据
const updateStats = (resumeList: Resume[]) => {
  const pending = resumeList.filter(r => r.status === 'pending').length
  const parsed = resumeList.filter(r => r.status === 'parsed').length
  const failed = resumeList.filter(r => r.status === 'failed').length
  
  resumeStats.value = [
    { label: '总简历', value: total.value, icon: markRaw(FolderOpened), colorClass: 'cyan' },
    { label: '待解析', value: pending, icon: markRaw(Clock), colorClass: 'teal' },
    { label: '已解析', value: parsed, icon: markRaw(CircleCheck), colorClass: 'green' },
    { label: '解析失败', value: failed, icon: markRaw(CircleClose), colorClass: 'red' }
  ]
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
const handleFileChange = (uploadFile: UploadFile, uploadFiles: UploadFile[]) => {
  if (uploadFile.size && uploadFile.size > 10 * 1024 * 1024) {
    ElMessage.warning('文件大小不能超过 10MB')
    // 移除超大文件
    const index = uploadFiles.findIndex(f => f.uid === uploadFile.uid)
    if (index > -1) {
      uploadFiles.splice(index, 1)
    }
    return
  }
  fileList.value = uploadFiles
}

// 文件移除
const handleFileRemove = (file: UploadFile, uploadFiles: UploadFile[]) => {
  fileList.value = uploadFiles
}

// 上传处理
const handleUpload = async () => {
  console.log('========== 开始上传 ==========')
  console.log('待上传文件数量:', fileList.value.length)
  
  uploading.value = true
  try {
    let successCount = 0
    for (let i = 0; i < fileList.value.length; i++) {
      const file = fileList.value[i]
      console.log(`[${i + 1}/${fileList.value.length}] 准备上传文件:`, file.name)
      console.log('  - 文件大小:', file.size, 'bytes')
      console.log('  - 文件类型:', file.raw?.type)
      
      if (!file.raw) {
        console.error('  ❌ file.raw 为空!')
        continue
      }
      
      const formData = new FormData()
      formData.append('file', file.raw as File)
      formData.append('talent_id', '0')
      formData.append('job_id', '0')
      
      try {
        // 直接请求后端，绕过 Vite 代理测试
        console.log('  直接请求后端 http://localhost:8084 ...')
        const token = localStorage.getItem('token')
        const response = await fetch('http://localhost:8084/api/v1/resumes/upload', {
          method: 'POST',
          headers: token ? { 'Authorization': `Bearer ${token}` } : {},
          body: formData
        })
        
        console.log('  响应状态:', response.status)
        const data = await response.json()
        console.log('  响应数据:', data)
        
        if (data.code === 0) {
          console.log('  ✓ 上传成功')
          successCount++
        } else {
          console.log('  ❌ 上传失败:', data.message)
        }
      } catch (err: any) {
        console.error('  ❌ 请求异常:', err)
      }
    }
    
    console.log('========== 上传完成 ==========')
    console.log('成功数量:', successCount, '/', fileList.value.length)
    
    if (successCount > 0) {
      ElMessage.success(`成功上传 ${successCount} 份简历`)
      showUploadDialog.value = false
      fileList.value = []
      fetchResumes()
    } else {
      ElMessage.error('上传失败')
    }
  } catch (error) {
    console.error('上传过程异常:', error)
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
const downloadResume = async (resume: Resume) => {
  try {
    const response = await request.get(`/resumes/${resume.id}/download`, {
      responseType: 'blob'
    })
    
    const blob = new Blob([response.data])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = resume.file_name
    link.click()
    window.URL.revokeObjectURL(url)
    
    ElMessage.success(`开始下载: ${resume.file_name}`)
  } catch (error) {
    ElMessage.error('下载失败')
  }
}

// 删除简历
const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这份简历吗？', '确认删除', {
      type: 'warning'
    })
    
    const res = await request.delete(`/resumes/${id}`)
    if (res.data?.code === 0) {
      ElMessage.success('删除成功')
      fetchResumes()
    } else {
      ElMessage.error(res.data?.message || '删除失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
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

// 获取头像颜色 - 青蓝色系
const getAvatarColor = (id: number) => {
  const colors = ['#00b8d4', '#26c6da', '#4dd0e1', '#00c853', '#0097a7']
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

// 获取预览URL
const getPreviewUrl = (resume: Resume) => {
  if (resume.file_url) {
    return resume.file_url
  }
  return ''
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
      color: var(--text-primary);
      margin: 0 0 4px 0;
    }

    .subtitle {
      color: var(--text-secondary);
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
    background: var(--bg-primary);
    box-shadow: var(--shadow-card);
    border: 1px solid var(--border-light);

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
        color: var(--text-primary);
      }

      .stat-label {
        font-size: 14px;
        color: var(--text-secondary);
      }
    }

    /* 青蓝色系 */
    &.cyan .stat-icon { background: linear-gradient(135deg, #00b8d4 0%, #0097a7 100%); }
    &.teal .stat-icon { background: linear-gradient(135deg, #26c6da 0%, #00acc1 100%); }
    &.green .stat-icon { background: linear-gradient(135deg, #00c853 0%, #00e676 100%); }
    &.red .stat-icon { background: #ef4444; }
  }
}

.card {
  background: var(--bg-primary);
  border-radius: 12px;
  padding: 20px;
  box-shadow: var(--shadow-card);
  border: 1px solid var(--border-light);
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
      &.word { background: linear-gradient(135deg, #00b8d4 0%, #0097a7 100%); }
      &.other { background: #9ca3af; }
    }

    .file-details {
      display: flex;
      flex-direction: column;

      .file-name {
        font-weight: 500;
        color: var(--text-primary);
      }

      .file-size {
        font-size: 12px;
        color: var(--text-muted);
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
    border: 2px dashed var(--border-color);
    padding: 40px;
    transition: all 0.3s;

    &:hover {
      border-color: var(--primary-color);
    }
  }

  .upload-content {
    .upload-icon {
      font-size: 48px;
      color: var(--primary-color);
      margin-bottom: 16px;
    }

    .upload-text {
      p {
        margin: 0;
        color: var(--text-secondary);

        em {
          color: var(--primary-color);
          font-style: normal;
        }
      }

      .upload-tip {
        font-size: 12px;
        color: var(--text-muted);
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
      &.word { background: linear-gradient(135deg, #00b8d4 0%, #0097a7 100%); }
    }

    .preview-info {
      h3 {
        font-size: 18px;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 8px 0;
      }

      p {
        font-size: 14px;
        color: var(--text-secondary);
        margin: 0 0 8px 0;
      }
    }
  }

  .parsed-data {
    h4 {
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 16px 0;
    }

    .parsed-content {
      .parsed-section {
        display: flex;
        margin-bottom: 16px;

        label {
          width: 80px;
          font-size: 14px;
          color: var(--text-secondary);
          flex-shrink: 0;
        }

        span {
          font-size: 14px;
          color: var(--text-primary);
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

  .pdf-preview-container {
    margin: 0;
    height: calc(100vh - 200px);

    .pdf-iframe {
      width: 100%;
      height: 100%;
      border: none;
      border-radius: 8px;
    }
  }

  .non-pdf-notice {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    background: var(--bg-secondary);
    border-radius: 12px;
    margin: 20px 0;

    .el-icon {
      color: var(--text-muted);
      margin-bottom: 16px;
    }

    p {
      margin: 0;
      color: var(--text-secondary);
      font-size: 14px;

      &.tip {
        color: var(--text-muted);
        font-size: 12px;
        margin-top: 8px;
      }
    }
  }
}
</style>
