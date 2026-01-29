<template>
  <div class="education-form">
    <div class="section-header">
      <span class="section-title">教育经历</span>
      <el-button type="primary" link @click="addEducation">
        <el-icon><Plus /></el-icon> 添加教育经历
      </el-button>
    </div>

    <el-empty v-if="educations.length === 0" description="暂无教育经历，点击上方按钮添加" :image-size="80" />

    <div v-else class="education-list">
      <div
        v-for="(edu, index) in educations"
        :key="index"
        class="education-item"
      >
        <div class="item-header">
          <span class="item-index">教育经历 {{ index + 1 }}</span>
          <el-button type="danger" link size="small" @click="removeEducation(index)">
            <el-icon><Delete /></el-icon> 删除
          </el-button>
        </div>

        <el-form
          :ref="(el: any) => setFormRef(el, index)"
          :model="edu"
          :rules="rules"
          label-width="80px"
          label-position="top"
        >
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item label="学校名称" prop="school">
                <el-input
                  v-model="edu.school"
                  placeholder="请输入学校名称"
                  maxlength="100"
                  @input="emitUpdate"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item label="专业" prop="major">
                <el-input
                  v-model="edu.major"
                  placeholder="请输入专业名称"
                  maxlength="100"
                  @input="emitUpdate"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <el-form-item label="学历" prop="degree">
                <el-select
                  v-model="edu.degree"
                  placeholder="请选择学历"
                  style="width: 100%"
                  @change="emitUpdate"
                >
                  <el-option
                    v-for="degree in degreeOptions"
                    :key="degree.value"
                    :label="degree.label"
                    :value="degree.value"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <el-form-item label="GPA" prop="gpa">
                <el-input
                  v-model="edu.gpa"
                  placeholder="如：3.8/4.0"
                  maxlength="20"
                  @input="emitUpdate"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <el-form-item label="开始时间" prop="start_date">
                <el-date-picker
                  v-model="edu.start_date"
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
                  v-model="edu.end_date"
                  type="month"
                  placeholder="选择结束时间"
                  format="YYYY-MM"
                  value-format="YYYY-MM"
                  style="width: 100%"
                  :disabled="edu.is_current"
                  @change="emitUpdate"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <el-form-item label=" ">
                <el-checkbox
                  v-model="edu.is_current"
                  @change="handleCurrentChange(index)"
                >
                  在读
                </el-checkbox>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="校园活动/荣誉" prop="activities">
                <el-input
                  v-model="edu.activities"
                  type="textarea"
                  :rows="2"
                  placeholder="请描述您的校园活动、获奖经历或荣誉"
                  maxlength="500"
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

// 教育经历接口
export interface Education {
  school: string
  degree: string
  major: string
  start_date: string
  end_date?: string
  is_current: boolean
  gpa?: string
  activities?: string
}

const props = defineProps<{
  modelValue: Education[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Education[]): void
}>()

// 学历选项
const degreeOptions = [
  { label: '高中', value: '高中' },
  { label: '中专', value: '中专' },
  { label: '大专', value: '大专' },
  { label: '本科', value: '本科' },
  { label: '硕士', value: '硕士' },
  { label: '博士', value: '博士' },
  { label: 'MBA', value: 'MBA' },
  { label: '其他', value: '其他' }
]

// 表单引用数组
const formRefs = ref<(FormInstance | null)[]>([])

const setFormRef = (el: FormInstance | null, index: number) => {
  formRefs.value[index] = el
}

// 教育经历列表
const educations = ref<Education[]>([])

// 表单验证规则
const rules = reactive<FormRules<Education>>({
  school: [
    { required: true, message: '请输入学校名称', trigger: 'blur' },
    { max: 100, message: '学校名称不能超过100个字符', trigger: 'blur' }
  ],
  degree: [
    { required: true, message: '请选择学历', trigger: 'change' }
  ],
  major: [
    { required: true, message: '请输入专业名称', trigger: 'blur' },
    { max: 100, message: '专业名称不能超过100个字符', trigger: 'blur' }
  ],
  start_date: [
    { required: true, message: '请选择开始时间', trigger: 'change' }
  ],
  gpa: [
    { max: 20, message: 'GPA不能超过20个字符', trigger: 'blur' }
  ],
  activities: [
    { max: 500, message: '校园活动描述不能超过500个字符', trigger: 'blur' }
  ]
})

// 监听 props 变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal && Array.isArray(newVal)) {
      educations.value = newVal.map(edu => ({ ...edu }))
    } else {
      educations.value = []
    }
  },
  { immediate: true, deep: true }
)

// 添加教育经历
const addEducation = () => {
  educations.value.push({
    school: '',
    degree: '',
    major: '',
    start_date: '',
    end_date: '',
    is_current: false,
    gpa: '',
    activities: ''
  })
  emitUpdate()
}

// 删除教育经历
const removeEducation = (index: number) => {
  educations.value.splice(index, 1)
  formRefs.value.splice(index, 1)
  emitUpdate()
}

// 处理"在读"变化
const handleCurrentChange = (index: number) => {
  if (educations.value[index].is_current) {
    educations.value[index].end_date = ''
  }
  emitUpdate()
}

// 发送更新事件
const emitUpdate = () => {
  emit('update:modelValue', [...educations.value])
}

// 验证所有表单
const validate = async (): Promise<boolean> => {
  if (educations.value.length === 0) return true
  
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
  educations.value = []
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
.education-form {
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

  .education-list {
    .education-item {
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
          color: #10b981;
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
