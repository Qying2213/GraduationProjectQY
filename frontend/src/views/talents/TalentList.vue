<template>
  <div class="talent-list">
    <div class="page-header">
      <h1>人才管理</h1>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        新增人才
      </el-button>
    </div>
    
    <!-- Search Bar -->
    <div class="card search-bar">
      <el-form :inline="true">
        <el-form-item>
          <el-input
            v-model="searchParams.search"
            placeholder="搜索姓名、邮箱"
            clearable
            @clear="fetchTalents"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchParams.status" placeholder="状态" clearable>
            <el-option label="在职" value="active" />
            <el-option label="已雇佣" value="hired" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchTalents">搜索</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- Table -->
    <div class="card">
      <el-table :data="talents" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="phone" label="电话" width="150" />
        <el-table-column prop="skills" label="技能" width="200">
          <template #default="{ row }">
            <el-tag v-for="skill in row.skills?.slice(0, 3)" :key="skill" size="small" style="margin-right: 4px">
              {{ skill }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="experience" label="经验" width="100" />
        <el-table-column prop="location" label="地区" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchTalents"
          @size-change="fetchTalents"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { talentApi } from '@/api/talent'
import type { Talent } from '@/types'
import { ElMessage, ElMessageBox } from 'element-plus'

const talents = ref<Talent[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const showCreateDialog = ref(false)

const searchParams = ref({
  search: '',
  status: ''
})

const fetchTalents = async () => {
  loading.value = true
  try {
    const res = await talentApi.list({
      page: currentPage.value,
      page_size: pageSize.value,
      ...searchParams.value
    })
    
    if (res.data.code === 0 && res.data.data) {
      talents.value = res.data.data.talents || []
      total.value = res.data.data.total || 0
    }
  } catch (error) {
    ElMessage.error('获取人才列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (talent: Talent) => {
  ElMessage.info('编辑功能开发中')
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个人才吗？', '提示', {
      type: 'warning'
    })
    
    await talentApi.delete(id)
    ElMessage.success('删除成功')
    fetchTalents()
  } catch (error) {
    // 用户取消
  }
}

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    active: 'success',
    hired: 'info',
    rejected: 'danger',
    pending: 'warning'
  }
  return map[status] || ''
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '活跃',
    hired: '已雇佣',
    rejected: '已拒绝',
    pending: '待处理'
  }
  return map[status] || status
}

onMounted(() => {
  fetchTalents()
})
</script>

<style scoped lang="scss">
.talent-list {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    
    h1 {
      font-size: 24px;
      font-weight: 700;
      color: var(--text-primary);
      margin: 0;
    }
  }
  
  .search-bar {
    margin-bottom: 16px;
  }
  
  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
