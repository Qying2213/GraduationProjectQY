<template>
  <div class="work-experience-form">
    <div class="section-header">
      <span class="section-title">工作经历</span>
      <el-button type="primary" link @click="addExperience">
        <el-icon><Plus /></el-icon> 添加工作经历
      </el-button>
    </div>

    <el-empty v-if="experiences.length === 0" description="暂无工作经历，点击上方按钮添加" :image-size="80" />

    <div v-else class="experience-list">
      <div
        v-for="(exp, index) in experiences"
        :key="index"
        class="experience-item"
      >
        <div class="item-header">
          <span class="item-index">工作经历 {{ index + 1 }}</span>
          <el-button type="danger" link size="small" @click="removeExperience(index)">
            <el-icon><Delete /></el-icon> 删除
          </el-button>
        </div>

        <el-form
          :ref="(el: any) => setFormRef(el, index)"
          :model="exp"
          :rules="rules"
          label-width="80px"
          label-position="top"
        >
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item label="公司名称" prop="company">
                <el-input
                  v-model="exp.company"
                  placeholder="请输入公司名称"
                  maxlength="100"
                  @input="emitUpdate"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item label="职位名称" prop="position">
                <el-input
                  v-model="exp.position"
                  placeholder="请输入职位名称"
                  maxlength="100"
                  @input="emitUpdate"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <el-form-item label="开始时间" prop="start_date">
                <el-date-picker
                  v-model="exp.start_date"
                  type="month"
                  placeholder="选择开始时间"
                  format="YYYY-MM"
                  value-format="YYYY-MM"
                  style="width: 100%"
                  @change="emitUpdate"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <el-form-item label="结束时间" prop="end_date">
                <el-date-picker
                  v-model="exp.end_date"
                  type="month"
                  placeholder="选择结束时间"
                  format="YYYY-MM"
                  value-format="YYYY-MM"
                  style="width: 100%"
                  :disabled="exp.is_current"
                  @change="emitUpdate"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <el-form-item label="工作地点" prop="location">
                <el-input
                  v-model="exp.location"
                  placeholder="请输入工作地点"
                  maxlength="100"
                  @input="emitUpdate"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <el-form-item label=" ">
                <el-checkbox
                  v-model="exp.is_current"
                  @change="handleCurrentChange(index)"
                >
                  至今在职
                </el-checkbox>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="工作描述" prop="description">
                <el-input
                  v-model="exp.description"
                  type="textarea"
                  :rows="3"
                  placeholder="请描述您的工作职责和成就"
                  maxlength="1000"
                  show-word-limit
                  @input="emitUpdate"
                />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { Plus, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

// 工作经历接口
export interface WorkExperience {
  company: string
  position: string
  start_date: string
  end_date?: string
  is_current: boolean
  description?: string
  location?: string
}

const props = defineProps<{
  modelValue: WorkExperience[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: WorkExperience[]): void
}>()

// 表单引用数组
const formRefs = ref<(FormInstance | null)[]>([])

const setFormRef = (el: FormInstance | null, index: number) => {
  formRefs.value[index] = el
}

// 工作经历列表
const experiences = ref<WorkExperience[]>([])

// 表单验证规则
const rules = reactive<FormRules<WorkExperience>>({
  company: [
    { required: true, message: '请输入公司名称', trigger: 'blur' },
    { max: 100, message: '公司名称不能超过100个字符', trigger: 'blur' }
  ],
  position: [
    { required: true, message: '请输入职位名称', trigger: 'blur' },
    { max: 100, message: '职位名称不能超过100个字符', trigger: 'blur' }
  ],
  start_date: [
    { required: true, message: '请选择开始时间', trigger: 'change' }
  ],
  description: [
    { max: 1000, message: '工作描述不能超过1000个字符', trigger: 'blur' }
  ],
  location: [
    { max: 100, message: '工作地点不能超过100个字符', trigger: 'blur' }
  ]
})

// 监听 props 变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal && Array.isArray(newVal)) {
      experiences.value = newVal.map(exp => ({ ...exp }))
    } else {
      experiences.value = []
    }
  },
  { immediate: true, deep: true }
)

// 添加工作经历
const addExperience = () => {
  experiences.value.push({
    company: '',
    position: '',
    start_date: '',
    end_date: '',
    is_current: false,
    description: '',
    location: ''
  })
  emitUpdate()
}

// 删除工作经历
const removeExperience = (index: number) => {
  experiences.value.splice(index, 1)
  formRefs.value.splice(index, 1)
  emitUpdate()
}

// 处理"至今在职"变化
const handleCurrentChange = (index: number) => {
  if (experiences.value[index].is_current) {
    experiences.value[index].end_date = ''
  }
  emitUpdate()
}

// 发送更新事件
const emitUpdate = () => {
  emit('update:modelValue', [...experiences.value])
}

// 验证所有表单
const validate = async (): Promise<boolean> => {
  if (experiences.value.length === 0) return true
  
  const results = await Promise.all(
    formRefs.value.map(async (formRef) => {
      if (!formRef) return true
      try {
        await formRef.validate()
        return true
      } catch {
        return false
      }
    })
  )
  
  return results.every(result => result)
}

// 重置表单
const resetForm = () => {
  experiences.value = []
  formRefs.value = []
  emitUpdate()
}

// 暴露方法给父组件
defineExpose({
  validate,
  resetForm
})
</script>

<style scoped lang="scss">
.work-experience-form {
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #f1f5f9;

    .section-title {
      font-size: 16px;
      font-weight: 600;
      color: #1e293b;
    }
  }

  .experience-list {
    .experience-item {
      background: #f8fafc;
      border-radius: 8px;
      padding: 16px;
      margin-bottom: 16px;

      &:last-child {
        margin-bottom: 0;
      }

      .item-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;

        .item-index {
          font-size: 14px;
          font-weight: 500;
          color: #0ea5e9;
        }
      }

      .el-form {
        .el-form-item {
          margin-bottom: 16px;
        }
      }
    }
  }
}
</style>
