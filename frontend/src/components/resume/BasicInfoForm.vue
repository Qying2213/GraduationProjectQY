<template>
  <el-form
    ref="formRef"
    :model="formData"
    :rules="rules"
    label-width="80px"
    label-position="top"
  >
    <el-row :gutter="20">
      <el-col :xs="24" :sm="12" :md="8">
        <el-form-item label="姓名" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入姓名"
            maxlength="50"
            @input="emitUpdate"
          />
        </el-form-item>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8">
        <el-form-item label="手机号" prop="phone">
          <el-input
            v-model="formData.phone"
            placeholder="请输入手机号"
            maxlength="20"
            @input="emitUpdate"
          />
        </el-form-item>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8">
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="formData.email"
            placeholder="请输入邮箱"
            maxlength="100"
            @input="emitUpdate"
          />
        </el-form-item>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8">
        <el-form-item label="现居地" prop="location">
          <el-input
            v-model="formData.location"
            placeholder="请输入现居城市"
            maxlength="100"
            @input="emitUpdate"
          />
        </el-form-item>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8">
        <el-form-item label="性别" prop="gender">
          <el-select
            v-model="formData.gender"
            placeholder="请选择性别"
            clearable
            style="width: 100%"
            @change="emitUpdate"
          >
            <el-option label="男" value="男" />
            <el-option label="女" value="女" />
          </el-select>
        </el-form-item>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8">
        <el-form-item label="年龄" prop="age">
          <el-input-number
            v-model="formData.age"
            :min="16"
            :max="100"
            placeholder="请输入年龄"
            style="width: 100%"
            @change="emitUpdate"
          />
        </el-form-item>
      </el-col>
      <el-col :span="24">
        <el-form-item label="个人简介" prop="summary">
          <el-input
            v-model="formData.summary"
            type="textarea"
            :rows="4"
            placeholder="请简要介绍自己的职业背景、专业技能和求职意向"
            maxlength="500"
            show-word-limit
            @input="emitUpdate"
          />
        </el-form-item>
      </el-col>
    </el-row>
  </el-form>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'

// 基本信息接口
export interface BasicInfo {
  name: string
  phone: string
  email: string
  location: string
  avatar?: string
  gender?: string
  age?: number
  summary?: string
}

const props = defineProps<{
  modelValue: BasicInfo
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: BasicInfo): void
}>()

const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive<BasicInfo>({
  name: '',
  phone: '',
  email: '',
  location: '',
  avatar: '',
  gender: '',
  age: undefined,
  summary: ''
})

// 表单验证规则
const rules = reactive<FormRules<BasicInfo>>({
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 50, message: '姓名长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  location: [
    { max: 100, message: '现居地不能超过100个字符', trigger: 'blur' }
  ],
  summary: [
    { max: 500, message: '个人简介不能超过500个字符', trigger: 'blur' }
  ]
})

// 监听 props 变化，同步到表单数据
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal) {
      formData.name = newVal.name || ''
      formData.phone = newVal.phone || ''
      formData.email = newVal.email || ''
      formData.location = newVal.location || ''
      formData.avatar = newVal.avatar || ''
      formData.gender = newVal.gender || ''
      formData.age = newVal.age || undefined
      formData.summary = newVal.summary || ''
    }
  },
  { immediate: true, deep: true }
)

// 发送更新事件
const emitUpdate = () => {
  emit('update:modelValue', { ...formData })
}

// 验证表单
const validate = async (): Promise<boolean> => {
  if (!formRef.value) return false
  try {
    await formRef.value.validate()
    return true
  } catch {
    return false
  }
}

// 重置表单
const resetForm = () => {
  formRef.value?.resetFields()
}

// 暴露方法给父组件
defineExpose({
  validate,
  resetForm
})
</script>

<style scoped lang="scss">
.el-form {
  .el-form-item {
    margin-bottom: 18px;
  }
}
</style>
