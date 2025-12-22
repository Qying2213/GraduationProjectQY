<template>
  <el-dialog v-model="visible" :title="title" width="600px" @close="handleClose">
    <div class="batch-import">
      <!-- 上传区域 -->
      <el-upload
        ref="uploadRef"
        class="upload-area"
        drag
        :auto-upload="false"
        :limit="1"
        :accept="acceptTypes"
        :on-change="handleFileChange"
        :on-exceed="handleExceed"
      >
        <el-icon class="upload-icon"><Upload /></el-icon>
        <div class="upload-text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="upload-tip">
            支持 {{ acceptTypes }} 格式，单次最多导入 {{ maxRows }} 条数据
          </div>
        </template>
      </el-upload>

      <!-- 模板下载 -->
      <div class="template-section">
        <el-button type="primary" link @click="downloadTemplate">
          <el-icon><Download /></el-icon>
          下载导入模板
        </el-button>
      </div>

      <!-- 预览数据 -->
      <div v-if="previewData.length > 0" class="preview-section">
        <div class="preview-header">
          <span>数据预览（共 {{ previewData.length }} 条）</span>
          <el-button type="danger" link @click="clearFile">清除</el-button>
        </div>
        <el-table :data="previewData.slice(0, 5)" max-height="300" size="small">
          <el-table-column
            v-for="col in columns"
            :key="col.key"
            :prop="col.key"
            :label="col.title"
            :width="col.width"
          />
        </el-table>
        <div v-if="previewData.length > 5" class="preview-more">
          ... 还有 {{ previewData.length - 5 }} 条数据
        </div>
      </div>

      <!-- 验证结果 -->
      <div v-if="validationErrors.length > 0" class="validation-section">
        <el-alert type="warning" :closable="false">
          <template #title>
            发现 {{ validationErrors.length }} 条数据存在问题
          </template>
          <div class="error-list">
            <div v-for="(error, index) in validationErrors.slice(0, 5)" :key="index" class="error-item">
              第 {{ error.row }} 行: {{ error.message }}
            </div>
            <div v-if="validationErrors.length > 5">
              ... 还有 {{ validationErrors.length - 5 }} 条错误
            </div>
          </div>
        </el-alert>
      </div>

      <!-- 导入进度 -->
      <div v-if="importing" class="progress-section">
        <el-progress :percentage="importProgress" :status="importStatus" />
        <div class="progress-text">
          已导入 {{ importedCount }} / {{ previewData.length }} 条
        </div>
      </div>
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button
        type="primary"
        :loading="importing"
        :disabled="previewData.length === 0 || validationErrors.length > 0"
        @click="handleImport"
      >
        {{ importing ? '导入中...' : '开始导入' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload, Download } from '@element-plus/icons-vue'
import * as XLSX from 'xlsx'
import type { UploadFile, UploadInstance } from 'element-plus'

interface Column {
  key: string
  title: string
  width?: number
  required?: boolean
  validator?: (value: any) => string | null
}

interface ValidationError {
  row: number
  column: string
  message: string
}

const props = withDefaults(defineProps<{
  modelValue: boolean
  title?: string
  columns: Column[]
  acceptTypes?: string
  maxRows?: number
  templateData?: any[]
  importFn: (data: any[]) => Promise<{ success: number; failed: number }>
}>(), {
  title: '批量导入',
  acceptTypes: '.xlsx,.xls,.csv',
  maxRows: 1000,
  templateData: () => [],
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success', result: { success: number; failed: number }): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const uploadRef = ref<UploadInstance>()
const previewData = ref<any[]>([])
const validationErrors = ref<ValidationError[]>([])
const importing = ref(false)
const importProgress = ref(0)
const importedCount = ref(0)
const importStatus = ref<'' | 'success' | 'exception'>('')

// 处理文件选择
const handleFileChange = (file: UploadFile) => {
  if (!file.raw) return

  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const data = e.target?.result
      const workbook = XLSX.read(data, { type: 'binary' })
      const sheetName = workbook.SheetNames[0]
      const worksheet = workbook.Sheets[sheetName]
      const jsonData = XLSX.utils.sheet_to_json(worksheet)

      if (jsonData.length > props.maxRows) {
        ElMessage.warning(`数据量超过限制，最多支持 ${props.maxRows} 条`)
        clearFile()
        return
      }

      // 映射列名
      previewData.value = jsonData.map((row: any) => {
        const mapped: any = {}
        props.columns.forEach(col => {
          mapped[col.key] = row[col.title] ?? ''
        })
        return mapped
      })

      // 验证数据
      validateData()
    } catch (error) {
      ElMessage.error('文件解析失败，请检查文件格式')
      clearFile()
    }
  }
  reader.readAsBinaryString(file.raw)
}

// 验证数据
const validateData = () => {
  validationErrors.value = []

  previewData.value.forEach((row, index) => {
    props.columns.forEach(col => {
      // 必填验证
      if (col.required && !row[col.key]) {
        validationErrors.value.push({
          row: index + 2, // Excel行号从2开始（1是表头）
          column: col.key,
          message: `${col.title} 不能为空`,
        })
      }

      // 自定义验证
      if (col.validator && row[col.key]) {
        const error = col.validator(row[col.key])
        if (error) {
          validationErrors.value.push({
            row: index + 2,
            column: col.key,
            message: error,
          })
        }
      }
    })
  })
}

// 下载模板
const downloadTemplate = () => {
  const headers = props.columns.map(col => col.title)
  const templateRows = props.templateData.length > 0
    ? props.templateData.map(row => props.columns.map(col => row[col.key] ?? ''))
    : [props.columns.map(() => '')]

  const ws = XLSX.utils.aoa_to_sheet([headers, ...templateRows])
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '导入模板')
  XLSX.writeFile(wb, `${props.title}_模板.xlsx`)
}

// 处理超出限制
const handleExceed = () => {
  ElMessage.warning('只能上传一个文件')
}

// 清除文件
const clearFile = () => {
  uploadRef.value?.clearFiles()
  previewData.value = []
  validationErrors.value = []
}

// 执行导入
const handleImport = async () => {
  if (previewData.value.length === 0) return

  importing.value = true
  importProgress.value = 0
  importedCount.value = 0
  importStatus.value = ''

  try {
    // 分批导入
    const batchSize = 50
    const batches = []
    for (let i = 0; i < previewData.value.length; i += batchSize) {
      batches.push(previewData.value.slice(i, i + batchSize))
    }

    let totalSuccess = 0
    let totalFailed = 0

    for (const batch of batches) {
      const result = await props.importFn(batch)
      totalSuccess += result.success
      totalFailed += result.failed
      importedCount.value += batch.length
      importProgress.value = Math.round((importedCount.value / previewData.value.length) * 100)
    }

    importStatus.value = totalFailed === 0 ? 'success' : 'exception'

    ElMessage.success(`导入完成：成功 ${totalSuccess} 条，失败 ${totalFailed} 条`)
    emit('success', { success: totalSuccess, failed: totalFailed })

    if (totalFailed === 0) {
      setTimeout(() => {
        handleClose()
      }, 1500)
    }
  } catch (error) {
    importStatus.value = 'exception'
    ElMessage.error('导入失败，请重试')
  } finally {
    importing.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  clearFile()
  importing.value = false
  importProgress.value = 0
  importedCount.value = 0
  importStatus.value = ''
  visible.value = false
}

// 监听对话框关闭
watch(visible, (val) => {
  if (!val) {
    handleClose()
  }
})
</script>

<style scoped lang="scss">
.batch-import {
  .upload-area {
    :deep(.el-upload-dragger) {
      padding: 40px 20px;
    }

    .upload-icon {
      font-size: 48px;
      color: #c0c4cc;
      margin-bottom: 16px;
    }

    .upload-text {
      color: #606266;

      em {
        color: #409eff;
        font-style: normal;
      }
    }

    .upload-tip {
      font-size: 12px;
      color: #909399;
      margin-top: 8px;
    }
  }

  .template-section {
    text-align: center;
    margin: 16px 0;
  }

  .preview-section {
    margin-top: 20px;

    .preview-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      font-weight: 500;
    }

    .preview-more {
      text-align: center;
      color: #909399;
      font-size: 13px;
      margin-top: 8px;
    }
  }

  .validation-section {
    margin-top: 16px;

    .error-list {
      margin-top: 8px;
      font-size: 13px;

      .error-item {
        padding: 4px 0;
        color: #e6a23c;
      }
    }
  }

  .progress-section {
    margin-top: 20px;

    .progress-text {
      text-align: center;
      font-size: 13px;
      color: #606266;
      margin-top: 8px;
    }
  }
}
</style>
