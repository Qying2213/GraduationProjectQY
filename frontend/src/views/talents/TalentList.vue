<template>
  <div class="talent-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>人才管理</h1>
        <p class="subtitle">共 {{ total }} 位人才</p>
      </div>
      <div class="header-actions">
        <el-button-group>
          <el-button :type="viewMode === 'table' ? 'primary' : 'default'" @click="viewMode = 'table'">
            <el-icon><List /></el-icon>
          </el-button>
          <el-button :type="viewMode === 'card' ? 'primary' : 'default'" @click="viewMode = 'card'">
            <el-icon><Grid /></el-icon>
          </el-button>
        </el-button-group>
        <el-dropdown @command="handleExport" trigger="click">
          <el-button>
            <el-icon><Download /></el-icon>
            导出
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="excel">
                <el-icon><Document /></el-icon>
                导出 Excel
              </el-dropdown-item>
              <el-dropdown-item command="csv">
                <el-icon><DocumentCopy /></el-icon>
                导出 CSV
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button type="primary" @click="openCreateDialog" v-if="canCreate">
          <el-icon><Plus /></el-icon>
          新增人才
        </el-button>
      </div>
    </div>

    <!-- 搜索筛选 -->
    <div class="card search-card">
      <el-form :inline="true" :model="searchParams" class="search-form">
        <el-form-item>
          <el-input
            v-model="searchParams.search"
            placeholder="搜索姓名、邮箱、技能..."
            clearable
            style="width: 280px"
            @clear="fetchTalents"
            @keyup.enter="fetchTalents"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchParams.status" placeholder="状态筛选" clearable style="width: 140px">
            <el-option label="活跃" value="active" />
            <el-option label="已雇佣" value="hired" />
            <el-option label="待处理" value="pending" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchParams.experience" placeholder="经验要求" clearable style="width: 140px">
            <el-option label="应届生" value="0" />
            <el-option label="1-3年" value="1-3" />
            <el-option label="3-5年" value="3-5" />
            <el-option label="5-10年" value="5-10" />
            <el-option label="10年以上" value="10+" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchTalents">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格视图 -->
    <div v-if="viewMode === 'table'" class="card table-card">
      <el-table :data="talents" style="width: 100%" v-loading="loading" row-key="id"
                @row-click="openDetailDrawer">
        <el-table-column prop="name" label="姓名" width="150">
          <template #default="{ row }">
            <div class="talent-name-cell">
              <el-avatar :size="36" :style="{ background: getAvatarColor(row.id) }">
                {{ row.name?.charAt(0) }}
              </el-avatar>
              <div class="name-info">
                <span class="name">{{ row.name }}</span>
                <span class="email">{{ row.email }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="电话" width="140" />
        <el-table-column prop="skills" label="技能" min-width="200">
          <template #default="{ row }">
            <div class="skills-cell">
              <el-tag v-for="skill in row.skills?.slice(0, 3)" :key="skill" size="small"
                      :type="getSkillType(skill)">
                {{ skill }}
              </el-tag>
              <el-tag v-if="row.skills?.length > 3" size="small" type="info">
                +{{ row.skills.length - 3 }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="experience" label="经验" width="100">
          <template #default="{ row }">
            {{ row.experience }}年
          </template>
        </el-table-column>
        <el-table-column prop="location" label="地区" width="100" />
        <el-table-column prop="salary" label="期望薪资" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="light">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="160">
          <template #default="{ row }">
            <el-button link type="primary" @click.stop="openEditDialog(row)" v-if="canEdit">编辑</el-button>
            <el-button link type="primary" @click.stop="openDetailDrawer(row)">详情</el-button>
            <el-button link type="danger" @click.stop="handleDelete(row.id)" v-if="canDelete">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchTalents"
          @size-change="fetchTalents"
        />
      </div>
    </div>

    <!-- 卡片视图 -->
    <div v-else class="card-view">
      <el-row :gutter="20" v-loading="loading">
        <el-col :xs="24" :sm="12" :lg="8" :xl="6" v-for="talent in talents" :key="talent.id">
          <div class="talent-card" @click="openDetailDrawer(talent)">
            <div class="card-header">
              <el-avatar :size="56" :style="{ background: getAvatarColor(talent.id) }">
                {{ talent.name?.charAt(0) }}
              </el-avatar>
              <el-dropdown trigger="click" @click.stop>
                <el-button text :icon="MoreFilled" />
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="openEditDialog(talent)" v-if="canEdit">
                      <el-icon><Edit /></el-icon> 编辑
                    </el-dropdown-item>
                    <el-dropdown-item divided @click="handleDelete(talent.id)" v-if="canDelete">
                      <el-icon><Delete /></el-icon> 删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
            <div class="card-body">
              <h3 class="talent-name">{{ talent.name }}</h3>
              <p class="talent-title">{{ talent.experience }}年经验 · {{ talent.location }}</p>
              <div class="talent-skills">
                <el-tag v-for="skill in talent.skills?.slice(0, 3)" :key="skill" size="small"
                        :type="getSkillType(skill)">
                  {{ skill }}
                </el-tag>
              </div>
            </div>
            <div class="card-footer">
              <div class="salary">
                <el-icon><Money /></el-icon>
                <span>{{ talent.salary || '面议' }}</span>
              </div>
              <el-tag :type="getStatusType(talent.status)" effect="light" size="small">
                {{ getStatusText(talent.status) }}
              </el-tag>
            </div>
          </div>
        </el-col>
      </el-row>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[12, 24, 48]"
          layout="total, sizes, prev, pager, next"
          @current-change="fetchTalents"
          @size-change="fetchTalents"
        />
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑人才' : '新增人才'"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="talentForm"
        :rules="formRules"
        label-width="100px"
        label-position="top"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="talentForm.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="talentForm.email" placeholder="请输入邮箱" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="电话" prop="phone">
              <el-input v-model="talentForm.phone" placeholder="请输入电话" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="工作经验(年)" prop="experience">
              <el-input-number v-model="talentForm.experience" :min="0" :max="50" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所在地区" prop="location">
              <el-input v-model="talentForm.location" placeholder="请输入地区" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="期望薪资" prop="salary">
              <el-input v-model="talentForm.salary" placeholder="如：20-30K" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="技能标签" prop="skills">
          <el-select
            v-model="talentForm.skills"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="选择或输入技能"
            style="width: 100%"
          >
            <el-option v-for="skill in commonSkills" :key="skill" :label="skill" :value="skill" />
          </el-select>
        </el-form-item>
        <el-form-item label="学历" prop="education">
          <el-select v-model="talentForm.education" placeholder="请选择学历" style="width: 100%">
            <el-option label="高中及以下" value="高中" />
            <el-option label="大专" value="大专" />
            <el-option label="本科" value="本科" />
            <el-option label="硕士" value="硕士" />
            <el-option label="博士" value="博士" />
          </el-select>
        </el-form-item>
        <el-form-item label="个人简介" prop="summary">
          <el-input
            v-model="talentForm.summary"
            type="textarea"
            :rows="4"
            placeholder="请输入个人简介"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="talentForm.status">
            <el-radio label="active">活跃</el-radio>
            <el-radio label="pending">待处理</el-radio>
            <el-radio label="hired">已雇佣</el-radio>
            <el-radio label="rejected">已拒绝</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ isEdit ? '保存修改' : '确认添加' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 详情抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      :title="currentTalent?.name || '人才详情'"
      size="450px"
    >
      <div class="talent-detail" v-if="currentTalent">
        <div class="detail-header">
          <el-avatar :size="80" :style="{ background: getAvatarColor(currentTalent.id) }">
            {{ currentTalent.name?.charAt(0) }}
          </el-avatar>
          <div class="detail-info">
            <h2>{{ currentTalent.name }}</h2>
            <p>{{ currentTalent.experience }}年经验 · {{ currentTalent.location }}</p>
            <el-tag :type="getStatusType(currentTalent.status)" effect="light">
              {{ getStatusText(currentTalent.status) }}
            </el-tag>
          </div>
        </div>

        <el-divider />

        <div class="detail-section">
          <h4>联系方式</h4>
          <div class="info-item">
            <el-icon><Message /></el-icon>
            <span>{{ currentTalent.email }}</span>
          </div>
          <div class="info-item">
            <el-icon><Phone /></el-icon>
            <span>{{ currentTalent.phone }}</span>
          </div>
        </div>

        <div class="detail-section">
          <h4>期望薪资</h4>
          <div class="salary-display">{{ currentTalent.salary || '面议' }}</div>
        </div>

        <div class="detail-section">
          <h4>技能标签</h4>
          <div class="skills-display">
            <el-tag v-for="skill in currentTalent.skills" :key="skill" :type="getSkillType(skill)">
              {{ skill }}
            </el-tag>
          </div>
        </div>

        <div class="detail-section">
          <h4>学历背景</h4>
          <p>{{ currentTalent.education || '未填写' }}</p>
        </div>

        <div class="detail-section">
          <h4>个人简介</h4>
          <p class="summary">{{ currentTalent.summary || '暂无简介' }}</p>
        </div>

        <div class="detail-actions">
          <el-button type="primary" @click="openEditDialog(currentTalent)">
            <el-icon><Edit /></el-icon>
            编辑信息
          </el-button>
          <el-button @click="$router.push('/recommend')">
            <el-icon><MagicStick /></el-icon>
            智能推荐
          </el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { talentApi } from '@/api/talent'
import type { Talent } from '@/types'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import {
  Search, Plus, List, Grid, MoreFilled, Edit, Delete, Money,
  Message, Phone, MagicStick, Download, ArrowDown, Document, DocumentCopy
} from '@element-plus/icons-vue'
import { exportToExcel, exportToCsv, talentExportColumns } from '@/utils/export'
import { usePermissionStore } from '@/store/permission'

const permissionStore = usePermissionStore()

// 权限检查
const canCreate = computed(() => permissionStore.hasPermission('talent:create'))
const canEdit = computed(() => permissionStore.hasPermission('talent:edit'))
const canDelete = computed(() => permissionStore.hasPermission('talent:delete'))
const canExport = computed(() => permissionStore.hasPermission('talent:export'))

const loading = ref(false)
const submitting = ref(false)
const talents = ref<Talent[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const viewMode = ref<'table' | 'card'>('table')
const dialogVisible = ref(false)
const drawerVisible = ref(false)
const isEdit = ref(false)
const currentTalent = ref<Talent | null>(null)
const formRef = ref<FormInstance>()

const searchParams = reactive({
  search: '',
  status: '',
  experience: ''
})

const talentForm = reactive<Partial<Talent>>({
  name: '',
  email: '',
  phone: '',
  skills: [],
  experience: 0,
  education: '',
  status: 'active',
  location: '',
  salary: '',
  summary: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱', trigger: 'blur' }
  ],
  phone: [{ required: true, message: '请输入电话', trigger: 'blur' }],
  skills: [{ required: true, message: '请选择至少一个技能', trigger: 'change' }]
}

const commonSkills = [
  'JavaScript', 'TypeScript', 'Vue', 'React', 'Angular', 'Node.js',
  'Go', 'Python', 'Java', 'C++', 'Rust', 'PHP',
  'MySQL', 'PostgreSQL', 'MongoDB', 'Redis',
  'Docker', 'Kubernetes', 'AWS', 'Linux',
  'Git', 'CI/CD', 'Agile', 'Scrum'
]

// 获取人才列表
const fetchTalents = async () => {
  loading.value = true
  try {
    const res = await talentApi.list({
      page: currentPage.value,
      page_size: pageSize.value,
      ...searchParams
    })

    if (res.data.code === 0 && res.data.data) {
      talents.value = res.data.data.talents || []
      total.value = res.data.data.total || 0
    } else {
      ElMessage.error(res.data?.message || '获取人才列表失败')
    }
  } catch (error) {
    console.error('获取人才列表失败:', error)
    ElMessage.error('获取人才列表失败，请检查后端服务是否启动')
    talents.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 生成模拟数据
const generateMockTalents = (): Talent[] => {
  const names = ['张三', '李四', '王五', '赵六', '钱七', '孙八', '周九', '吴十']
  const locations = ['北京', '上海', '深圳', '杭州', '广州', '成都', '南京', '武汉']
  const skills = ['Go', 'Python', 'Java', 'Vue', 'React', 'TypeScript', 'Node.js', 'MySQL', 'Redis', 'Docker']
  const statuses: Talent['status'][] = ['active', 'hired', 'pending', 'rejected']

  return Array.from({ length: pageSize.value }, (_, i) => ({
    id: (currentPage.value - 1) * pageSize.value + i + 1,
    name: names[i % names.length],
    email: `user${i + 1}@example.com`,
    phone: `138${String(Math.random()).slice(2, 10)}`,
    skills: skills.sort(() => Math.random() - 0.5).slice(0, Math.floor(Math.random() * 4) + 2),
    experience: Math.floor(Math.random() * 10) + 1,
    education: ['本科', '硕士', '博士'][Math.floor(Math.random() * 3)],
    status: statuses[Math.floor(Math.random() * statuses.length)],
    tags: [],
    location: locations[Math.floor(Math.random() * locations.length)],
    salary: `${Math.floor(Math.random() * 30) + 10}-${Math.floor(Math.random() * 30) + 40}K`,
    summary: '资深工程师，具有丰富的项目经验和团队协作能力。',
    user_id: i + 1,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString()
  }))
}

// 重置搜索
const resetSearch = () => {
  searchParams.search = ''
  searchParams.status = ''
  searchParams.experience = ''
  currentPage.value = 1
  fetchTalents()
}

// 打开新增弹窗
const openCreateDialog = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

// 打开编辑弹窗
const openEditDialog = (talent: Talent) => {
  isEdit.value = true
  Object.assign(talentForm, talent)
  dialogVisible.value = true
}

// 打开详情抽屉
const openDetailDrawer = (talent: Talent) => {
  currentTalent.value = talent
  drawerVisible.value = true
}

// 重置表单
const resetForm = () => {
  Object.assign(talentForm, {
    name: '',
    email: '',
    phone: '',
    skills: [],
    experience: 0,
    education: '',
    status: 'active',
    location: '',
    salary: '',
    summary: ''
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value && talentForm.id) {
          await talentApi.update(talentForm.id, talentForm)
          ElMessage.success('修改成功')
        } else {
          await talentApi.create(talentForm)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        fetchTalents()
      } catch (error) {
        ElMessage.error(isEdit.value ? '修改失败' : '添加失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

// 删除人才
const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这条记录吗？', '确认删除', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })

    await talentApi.delete(id)
    ElMessage.success('删除成功')
    fetchTalents()
  } catch {
    // 用户取消
  }
}

// 获取头像颜色 - 青蓝色系
const getAvatarColor = (id: number) => {
  const colors = ['#00b8d4', '#26c6da', '#4dd0e1', '#00c853', '#0097a7']
  return colors[id % colors.length]
}

// 获取技能标签类型
const getSkillType = (skill: string) => {
  const frontendSkills = ['Vue', 'React', 'Angular', 'JavaScript', 'TypeScript', 'HTML', 'CSS']
  const backendSkills = ['Go', 'Python', 'Java', 'Node.js', 'PHP', 'C++', 'Rust']
  const dbSkills = ['MySQL', 'PostgreSQL', 'MongoDB', 'Redis', 'Elasticsearch']

  if (frontendSkills.includes(skill)) return 'success'
  if (backendSkills.includes(skill)) return 'warning'
  if (dbSkills.includes(skill)) return 'danger'
  return 'info'
}

// 获取状态类型
const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    active: 'success',
    hired: 'info',
    rejected: 'danger',
    pending: 'warning'
  }
  return map[status] || ''
}

// 获取状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '活跃',
    hired: '已雇佣',
    rejected: '已拒绝',
    pending: '待处理'
  }
  return map[status] || status
}

// 导出数据
const handleExport = (format: 'excel' | 'csv') => {
  if (talents.value.length === 0) {
    ElMessage.warning('暂无数据可导出')
    return
  }

  const exportData = talents.value.map(t => ({
    ...t,
    currentCompany: t.summary?.split('，')[0] || '',
    currentPosition: '',
    expectedSalary: t.salary,
    createTime: t.created_at
  }))

  const options = {
    filename: '人才列表',
    sheetName: '人才数据',
    columns: talentExportColumns,
    data: exportData
  }

  if (format === 'excel') {
    exportToExcel(options)
    ElMessage.success('Excel 导出成功')
  } else {
    exportToCsv(options)
    ElMessage.success('CSV 导出成功')
  }
}

onMounted(() => {
  fetchTalents()
})
</script>

<style scoped lang="scss">
.talent-list {
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
  flex-wrap: wrap;
  gap: 16px;

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

  .header-actions {
    display: flex;
    gap: 12px;
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

.table-card {
  .talent-name-cell {
    display: flex;
    align-items: center;
    gap: 12px;

    .name-info {
      display: flex;
      flex-direction: column;

      .name {
        font-weight: 600;
        color: var(--text-primary);
      }

      .email {
        font-size: 12px;
        color: var(--text-muted);
      }
    }
  }

  .skills-cell {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

// 卡片视图
.card-view {
  .talent-card {
    background: var(--bg-primary);
    border-radius: 16px;
    padding: 20px;
    margin-bottom: 20px;
    box-shadow: var(--shadow-card);
    border: 1px solid var(--border-light);
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-4px);
      box-shadow: var(--shadow-lg);
      border-color: var(--primary-color);
    }

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 16px;
    }

    .card-body {
      .talent-name {
        font-size: 18px;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 4px 0;
      }

      .talent-title {
        font-size: 14px;
        color: var(--text-secondary);
        margin: 0 0 12px 0;
      }

      .talent-skills {
        display: flex;
        flex-wrap: wrap;
        gap: 6px;
      }
    }

    .card-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-top: 16px;
      padding-top: 16px;
      border-top: 1px solid var(--border-light);

      .salary {
        display: flex;
        align-items: center;
        gap: 6px;
        color: var(--primary-color);
        font-weight: 600;
      }
    }
  }
}

// 详情抽屉
.talent-detail {
  .detail-header {
    display: flex;
    align-items: center;
    gap: 20px;

    .detail-info {
      h2 {
        font-size: 22px;
        font-weight: 700;
        color: var(--text-primary);
        margin: 0 0 4px 0;
      }

      p {
        color: var(--text-secondary);
        font-size: 14px;
        margin: 0 0 8px 0;
      }
    }
  }

  .detail-section {
    margin-bottom: 24px;

    h4 {
      font-size: 14px;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 12px 0;
    }

    .info-item {
      display: flex;
      align-items: center;
      gap: 8px;
      color: var(--text-secondary);
      margin-bottom: 8px;
    }

    .salary-display {
      font-size: 20px;
      font-weight: 700;
      color: var(--primary-color);
    }

    .skills-display {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }

    .summary {
      color: var(--text-secondary);
      line-height: 1.6;
    }
  }

  .detail-actions {
    display: flex;
    gap: 12px;
    margin-top: 32px;
  }
}
</style>
