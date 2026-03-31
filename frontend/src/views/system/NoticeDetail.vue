<template>
  <div class="notice-detail" v-loading="loading">
    <div class="page-header">
      <el-button @click="goBack">返回</el-button>
    </div>

    <div v-if="notice" class="detail-card">
      <h1>{{ notice.title }}</h1>
      <div class="meta">
        <el-tag :type="notice.status === 'published' ? 'success' : 'info'">
          {{ notice.status === 'published' ? '已发布' : '草稿' }}
        </el-tag>


        <el-tag class="priority-tag" :class="priorityTagClass(notice.priority)">
          {{ priorityLabel(notice.priority) }}
        </el-tag>

        <el-tag v-if="notice.is_pinned" type="warning">置顶</el-tag>
        <span>创建时间: {{ notice.created_at }}</span>
        <span>创建人ID: {{ notice.created_by ?? '_' }}</span>
      </div>

      <div class="content">{{ notice.content }}</div>
    </div>

    <el-empty v-else description="公告不存在" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { noticeApi, type Notice, type NoticePriority } from '@/api/notice'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const notice = ref<Notice | null>(null)

const normalizePriority = (value?: string): NoticePriority => {
  if (value === 'urgent' || value === 'high') return value
  return 'normal'
}

const priorityLabel = (value?: string) => {
  const p = normalizePriority(value)
  if (p === 'urgent') return '紧急'
  if (p === 'high') return '重要'
  return '普通'
}

const priorityTagClass = (value?: string) => {
  const p = normalizePriority(value)
  return `priority-tag--${p}`
}

const loadNotice = async () => {
  const id = Number(route.params.id)
  if (!id) {
    ElMessage.error('公告ID无效')
    return
  }

  loading.value = true
  try {
    const res = await noticeApi.get(id)
    if (res.data.code === 0) {
      notice.value = res.data.data
    } else {
      ElMessage.error(res.data.message || '加载失败')
    }
  } catch (error: any) {
    if (error.response?.status === 404) {
      ElMessage.error('公告不存在')
    } else {
      ElMessage.error('加载公告详情失败')
    }
  } finally {
    loading.value = false
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
.notice-detail {
  padding: 24px;
}
.page-header {
  margin-bottom: 20px;
}
.detail-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
}
.meta {
  display: flex;
  gap: 16px;
  margin: 16px 0 24px;
  color: #666;
  font-size: 14px;
  flex-wrap: wrap;
}
.content {
  line-height: 1.8;
  color: #333;
  white-space: pre-wrap;
}


.priority-tag {
  border: none;
  font-weight: 600;
}
.priority-tag--normal {
  color: #595959;
  background: #f0f1f3;
}
.priority-tag--high {
  color: #ad6800;
  background: #fff7e6;
}
.priority-tag--urgent {
  color: #cf1322;
  background: #fff1f0;
}
</style>
