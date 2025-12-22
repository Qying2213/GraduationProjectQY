<template>
  <el-dialog
    v-model="visible"
    title="批量导入人才"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="import-container">
      <!-- 步骤条 -->
      <el-steps :active="currentStep" finish-status="success" align-center class="import-steps">
        <el-step title="上传文件" />
        <el-step title="数据预览" />
        <el-step title="导入结果" />
      </el-steps>

      <!-- 步骤1: 上传文件 -->
      <div v-show="currentStep === 0" class="step-content">
        <el-upload
          ref="uploadRef"
          class="upload-area"
          drag
          :auto-upload="false"
          :limit="1"
          accept=".xlsx,.xls,.csv"
          :on-change="handleFileChange"
          :on-exceed="handleExceed"
        >
          <el-icon class="upload-icon"><UploadFilled /></el-icon>
          <div class="upload-text">
            <p>将文件拖到此处，或<em>点击上传</em></p>
            <p class="upload-tip">支持 .xlsx, .xls, .csv 格式，单次最多导入 500 条</p>
          </div>
        </el-upload>

        <div class="template-download">
          <el-button link type="primary" @click="downloadTemplate">
            <el-icon><Download /></el-icon>
            下载导入模板
          </el-button>
        </div>

        <div class="import-tips">
          <h4>导入说明</h4>
          <ul>
            <li>请按照模板格式填写数据，必填字段：姓名、邮箱、电话</li>
            <li>技能字段请用逗号分隔，如：Vue,React,TypeScript</li>
            <li>工作经验请填写数字，单位为年</li>
            <li>重复的邮箱或电话将自动跳过</li>
          </ul>
        </div>
      </div>

      <!-- 步骤2: 数据预览 -->
      <div v-show="currentStep === 1" class="step-content">
        <div class="preview-header">
          <span>共解析 <strong>{{ previewData.length }}</strong> 条数据</span>
          <span class="error-count" v-if="errorCount > 0">
            <el-icon><Warning /></el-icon>
            {{ errorCount }} 条数据有问题
          </span>
        </div>

        <el-table :data="previewData" max-height="350" stripe size="small">
          <el-table-column type="index" width="50" label="#" />
          <el-table-column prop="name" label="姓名" width="100">
            <template #default="{ row }">
              <span :class="{ 'error-cell': !row.name }">{{ row.name || '缺失' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="email" label="邮箱" width="180">
            <template #default="{ row }">
              <span :class="{ 'error-cell': !row.email || row.emailError }">
                {{ row.email || '缺失' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="phone" label="电话" width="130" />
          <el-table-column prop="skills" label="技能">
            <template #default="{ row }">
              <div class="skills-preview">
                <el-tag v-for="skill in (row.skills || []).slice(0, 3)" :key="skill" size="small" type="info">
                  {{ skill }}
                </el-tag>
                <span v-if="(row.skills || []).length > 3">+{{ row.skills.length - 3 }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="experience" label="经验" width="80">
            <template #default="{ row }">
              {{ row.experience }}年
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.hasError ? 'danger' : 'success'" size="small">
                {{ row.hasError ? '异常' : '正常' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 步骤3: 导入结果 -->
      <div v-show="currentStep === 2" class="step-content">
        <div class="result-container">
          <div class="result-icon" :class="importResult.success ? 'success' : 'error'">
            <el-icon v-if="importResult.success"><CircleCheckFilled /></el-icon>
            <el-icon v-else><CircleCloseFilled /></el-icon>
          </div>
          <h3>{{ importResult.success ? '导入完成' : '导入失败' }}</h3>
          <p class="result-summary">
            成功导入 <strong>{{ importResult.successCount }}</strong> 条，
            失败 <strong>{{ importResult.failCount }}</strong> 条
          </p>

          <div class="result-details" v-if="importResult.failCount > 0">
            <h4>失败详情</h4>
            <el-table :data="importResult.failedItems" max-height="200" size="small">
              <el-table-column prop="row" label="行号" width="80" />
              <el-table-column prop="name" label="姓名" width="100" />
              <el-table-column prop="reason" label="失败原因" />
            </el-table>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button v-if="currentStep > 0 && currentStep < 2" @click="prevStep">上一步</el-button>
        <el-button v-if="currentStep === 0" type="primary" :disabled="!selectedFile" @click="parseFile">
          下一步
        </el-button>
        <el-button v-if="currentStep === 1" type="primary" :loading="importing" @click="startImport">
          开始导入
        </el-button>
        <el-button v-if="currentStep === 2" type="primary" @click="handleClose">
          完成
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, defineProps, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'
import * as XLSX from 'xlsx'
import {
  UploadFilled, Download, Warning, CircleCheckFilled, CircleCloseFilled
} from '@element-plus/icons-vue'

interface PreviewItem {
  name: string
  email: string
  phone: string
  skills: string[]
  experience: number
  education: string
  location: string
  salary: string
  hasError: boolean
  emailError?: boolean
}

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits(['update:modelValue', 'success'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const currentStep = ref(0)
const selectedFile = ref<File | null>(null)
const previewData = ref<PreviewItem[]>([])
const importing = ref(false)
const uploadRef = ref()

const importResult = ref({
  success: true,
  successCount: 0,
  failCount: 0,
  failedItems: [] as { row: number; name: string; reason: string }[]
})

const errorCount = computed(() => previewData.value.filter(item => item.hasError).length)

const handleFileChange = (file: any) => {
  selectedFile.value = file.raw
}

const handleExceed = () => {
  ElMessage.warning('只能上传一个文件')
}

const downloadTemplate = () => {
  const template = [
    ['姓名', '邮箱', '电话', '技能(逗号分隔)', '工作经验(年)', '学历', '所在地区', '期望薪资'],
    ['张三', 'zhangsan@example.com', '13800138000', 'Vue,React,TypeScript', '5', '本科', '北京', '20-30K'],
    ['李四', 'lisi@example.com', '13900139000', 'Go,Python,Docker', '3', '硕士', '上海', '25-35K']
  ]

  const ws = XLSX.utils.aoa_to_sheet(template)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '人才导入模板')
  XLSX.writeFile(wb, '人才导入模板.xlsx')
}

const parseFile = async () => {
  if (!selectedFile.value) return

  try {
    const data = await readExcelFile(selectedFile.value)
    previewData.value = data.map(row => ({
      ...row,
      hasError: !row.name || !row.email || !validateEmail(row.email),
      emailError: row.email && !validateEmail(row.email)
    }))
    currentStep.value = 1
  } catch (error) {
    ElMessage.error('文件解析失败，请检查文件格式')
  }
}

const readExcelFile = (file: File): Promise<PreviewItem[]> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      try {
        const data = new Uint8Array(e.target?.result as ArrayBuffer)
        const workbook = XLSX.read(data, { type: 'array' })
        const sheetName = workbook.SheetNames[0]
        const worksheet = workbook.Sheets[sheetName]
        const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 }) as any[][]

        // 跳过表头
        const rows = jsonData.slice(1).filter(row => row.length > 0)
        const result = rows.map(row => ({
          name: String(row[0] || ''),
          email: String(row[1] || ''),
          phone: String(row[2] || ''),
          skills: row[3] ? String(row[3]).split(/[,，]/).map(s => s.trim()).filter(Boolean) : [],
          experience: parseInt(row[4]) || 0,
          education: String(row[5] || ''),
          location: String(row[6] || ''),
          salary: String(row[7] || ''),
          hasError: false
        }))

        resolve(result)
      } catch (err) {
        reject(err)
      }
    }
    reader.onerror = reject
    reader.readAsArrayBuffer(file)
  })
}

const validateEmail = (email: string) => {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

const prevStep = () => {
  currentStep.value--
}

const startImport = async () => {
  importing.value = true

  // 模拟导入过程
  await new Promise(resolve => setTimeout(resolve, 2000))

  const validData = previewData.value.filter(item => !item.hasError)
  const invalidData = previewData.value.filter(item => item.hasError)

  importResult.value = {
    success: validData.length > 0,
    successCount: validData.length,
    failCount: invalidData.length,
    failedItems: invalidData.map((item, index) => ({
      row: index + 2,
      name: item.name || '未知',
      reason: !item.name ? '姓名缺失' : !item.email ? '邮箱缺失' : '邮箱格式错误'
    }))
  }

  importing.value = false
  currentStep.value = 2

  if (validData.length > 0) {
    emit('success', validData)
  }
}

const handleClose = () => {
  currentStep.value = 0
  selectedFile.value = null
  previewData.value = []
  uploadRef.value?.clearFiles()
  visible.value = false
}
</script>

<style scoped lang="scss">
.import-container {
  padding: 20px 0;
}

.import-steps {
  margin-bottom: 32px;
}

.step-content {
  min-height: 300px;
}

.upload-area {
  :deep(.el-upload-dragger) {
    padding: 40px;
    border-radius: 12px;
    border: 2px dashed var(--border-color);
    transition: all 0.3s ease;

    &:hover {
      border-color: var(--primary-color);
    }
  }

  .upload-icon {
    font-size: 48px;
    color: var(--text-tertiary);
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
      color: var(--text-tertiary);
      margin-top: 8px;
    }
  }
}

.template-download {
  text-align: center;
  margin-top: 16px;
}

.import-tips {
  margin-top: 24px;
  padding: 16px;
  background: var(--bg-tertiary);
  border-radius: 8px;

  h4 {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 12px 0;
  }

  ul {
    margin: 0;
    padding-left: 20px;

    li {
      font-size: 13px;
      color: var(--text-secondary);
      margin-bottom: 6px;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  strong {
    color: var(--primary-color);
  }

  .error-count {
    display: flex;
    align-items: center;
    gap: 4px;
    color: #ef4444;
    font-size: 14px;
  }
}

.error-cell {
  color: #ef4444;
}

.skills-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}

.result-container {
  text-align: center;
  padding: 40px 0;

  .result-icon {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 20px;
    font-size: 40px;

    &.success {
      background: rgba(67, 233, 123, 0.1);
      color: #43e97b;
    }

    &.error {
      background: rgba(239, 68, 68, 0.1);
      color: #ef4444;
    }
  }

  h3 {
    font-size: 20px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 8px 0;
  }

  .result-summary {
    color: var(--text-secondary);
    margin: 0 0 24px 0;

    strong {
      color: var(--primary-color);
    }
  }

  .result-details {
    text-align: left;
    margin-top: 24px;
    padding: 16px;
    background: var(--bg-tertiary);
    border-radius: 8px;

    h4 {
      font-size: 14px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 12px 0;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
