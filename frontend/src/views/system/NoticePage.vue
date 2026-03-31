<template>
  <div class="notice-page">
    <div class="page-header">
      <h1>公告管理</h1>
      <p>练习一个完整的 Vue3 页面</p>
    </div>

    <div class="toolbar">
      <el-input
        v-model="keyword"
        placeholder="搜索公告标题"
        clearable
        style="width: 260px"
        @keyup.enter="handleSearch"
      />
      <el-select v-model="status" placeholder="状态" clearable style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="已发布" value="published" />
        <el-option label="草稿" value="draft" />
      </el-select>
      <el-select v-model="pinned" placeholder="置顶" clearable style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="仅置顶" value="true" />
        <el-option label="仅非置顶" value="false" />
      </el-select>
      <el-select v-model="priority" placeholder="优先级" clearable style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="紧急" value="urgent" />
        <el-option label="重要" value="high" />
        <el-option label="普通" value="normal" />
      </el-select>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
      <el-button type="success" @click="openCreateDialog">新建公告</el-button>
    </div>

    <el-table :data="notices" v-loading="loading" style="width: 100%">
      <el-table-column prop="title" label="标题">
        <template #default="{ row }">
          <el-button link type="primary" @click="goToDetail(row.id)">
            {{ row.title }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态">
        <template #default="{ row }">
          <el-tag :type="row.status === 'published' ? 'success' : 'info'">
            {{ row.status === "published" ? "已发布" : "草稿" }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="优先级" width="100">
        <template #default="{ row }">
          <el-tag class="priority-tag" :class="priorityTagClass(row.priority)">
            {{ priorityLabel(row.priority) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="置顶" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.is_pinned" type="warning">置顶</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" />
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button link type="primary" @click="goToEdit(row.id)">编辑</el-button>
          <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>

    <el-dialog v-model="dialogVisible" title="新建公告" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="内容">
          <el-input v-model="form.content" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width: 100%">
            <el-option label="草稿" value="draft" />
            <el-option label="发布" value="published" />
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
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting"
          >保存</el-button
        >
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { noticeApi, type Notice, type NoticeFormData,NoticePriority } from '@/api/notice'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const keyword = ref('')
const notices = ref<Notice[]>([])
const status=ref('')
const router=useRouter()
const priority = ref<NoticePriority | ''>('')
const form = reactive<NoticeFormData>({
  title: '',
  content: '',
  status: 'draft',
  is_pinned: false,
  priority: 'normal'
})
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const pinned=ref('')
const priorityLabel = (priority: string) => {
  switch (priority) {
    case 'urgent':
      return '紧急'
    case 'high':
      return '重要'
    case 'normal':
      return '普通'
    default:
      return priority
  }
}

const normalizePriority = (value?: string): NoticePriority => {
  if (value === 'urgent' || value === 'high') return value
  return 'normal'
}

const priorityTagClass = (value?: string) => {
  const p = normalizePriority(value);
  return `priority-tag--${p}`;
};

const loadNotices = async () => {
  loading.value = true
  try {
    const res = await noticeApi.list({
  keyword: keyword.value,
  status: status.value,
  is_pinned: pinned.value === '' ? undefined : pinned.value === 'true',
  priority: priority.value || undefined, 
  page: currentPage.value,
  page_size: pageSize.value,
})
    if (res.data.code === 0) {
     notices.value = res.data.data.notices || []
    total.value = res.data.data.total || 0
    } else {
      ElMessage.error(res.data.message || '加载失败222')

    }
  } catch (error) {
    ElMessage.error('加载公告失败')
  } finally {
    loading.value = false
  }
}
const handleSearch = () => {
  currentPage.value = 1
  loadNotices()
}

const handleCurrentChange = () => {
  loadNotices()
}

const handleSizeChange = () => {
  currentPage.value = 1
  loadNotices()
}

const openCreateDialog = () => {
  form.title = ''
  form.content = ''
  form.status = 'draft'
  form.is_pinned = false
  form.priority = 'normal'
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!form.title.trim()) {
    ElMessage.warning('请输入标题')
    return
  }

  submitting.value = true
  try {
    const res = await noticeApi.create({
      title: form.title,
      content: form.content,
      status: form.status,
      is_pinned: form.is_pinned,
      priority: form.priority
    })

    if (res.data.code === 0) {
      ElMessage.success('创建成功')
      dialogVisible.value = false
      loadNotices()
    } else {
      ElMessage.error(res.data.message || '创建失败')
    }
  } catch (error) {
    ElMessage.error('创建公告失败')
  } finally {
    submitting.value = false
  }
}

const goToEdit = (id: number) => {
  router.push(`/notices/${id}/edit`)
}
const goToDetail = (id: number) => {
  router.push(`/notices/${id}`)
}
const handleDelete = async (row: Notice) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除公告《${row.title}》吗？`,
      '删除确认',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    const res = await noticeApi.remove(row.id)

    if (res.data.code === 0) {
      ElMessage.success('删除成功')
      await loadNotices()
    } else {
      ElMessage.error(res.data.message || '删除失败')
    }
  } catch (error: any) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    ElMessage.error('删除公告失败')
  }
}

onMounted(() => {
  loadNotices()
})
</script>

<style scoped lang="scss">
.notice-page {
  padding: 24px;
}

.page-header {
  margin-bottom: 20px;
}

.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
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
