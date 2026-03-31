<template>
  <div class="notice-edit" v-loading="loading">
    <div class="page-header">
      <h1>编辑公告</h1>
      <el-button @click="goBack">返回</el-button>
    </div>

    <div class="form-card">
      <el-form :model="form" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="请输入标题" />
        </el-form-item>

        <el-form-item label="内容">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="8"
            placeholder="请输入公告内容"
          />
        </el-form-item>

        <el-form-item label="状态">
          <el-select v-model="form.status" style="width: 100%">
            <el-option label="草稿" value="draft" />
            <el-option label="已发布" value="published" />
          </el-select>
        </el-form-item>

        <el-form-item label="优先级">
          <el-select v-model="form.priority" style="width: 100%">
            <el-option label="普通" value="normal" />
            <el-option label="重要" value="high" />
            <el-option label="紧急" value="urgent" />
          </el-select>
        </el-form-item>

        <el-form-item label="置顶">
          <el-switch v-model="form.is_pinned" />
        </el-form-item>

        <el-form-item>
          <el-button @click="goBack">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting"
            >保存</el-button
          >
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { noticeApi, type NoticeFormData,NoticePriority} from '@/api/notice'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const submitting = ref(false)

const form = reactive<NoticeFormData>({
  title: '',
  content: '',
  status: 'draft',
  is_pinned: false,
  priority: 'normal'
})
const normalizePriority = (value?: string): NoticePriority => {
  if (value === 'urgent' || value === 'high') return value
  return 'normal'
}
const noticeId = Number(route.params.id)

const loadNotice = async () => {
  if (!noticeId) {
    ElMessage.error('公告ID无效')
    return
  }

  loading.value = true
  try {
    const res = await noticeApi.get(noticeId)
    if (res.data.code === 0) {
      form.title = res.data.data.title
      form.content = res.data.data.content
      form.status = res.data.data.status
      form.is_pinned = res.data.data.is_pinned
      form.priority = normalizePriority(res.data.data.priority)
    } else {
      ElMessage.error(res.data.message || '加载详情失败')
    }
  } catch (error) {
    ElMessage.error('加载公告失败')
  } finally {
    loading.value = false
  }
}

const submitForm = async () => {
  if (!form.title.trim()) {
    ElMessage.warning('请输入标题')
    return
  }

  submitting.value = true
  try {
    const res = await noticeApi.update(noticeId, {
      title: form.title,
      content: form.content,
      status: form.status,
      is_pinned: form.is_pinned,
      priority: form.priority
    })

    if (res.data.code === 0) {
      ElMessage.success('更新成功')
      router.push('/notices')
    } else {
      ElMessage.error(res.data.message || '更新失败')
    }
  } catch (error) {
    ElMessage.error('更新公告失败')
  } finally {
    submitting.value = false
  }
}

const goBack = () => {
  router.push('/notices')
}

onMounted(() => {
  loadNotice()
})
</script>

<style scoped lang="scss">
.notice-edit {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.form-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
}
</style>
