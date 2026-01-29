<template>
  <div class="skills-form">
    <div class="section-header">
      <span class="section-title">技能标签</span>
      <span class="skill-count">{{ skills.length }} 个技能</span>
    </div>

    <div class="skills-input-area">
      <el-autocomplete
        v-model="inputValue"
        :fetch-suggestions="querySearch"
        placeholder="输入技能名称，按回车添加"
        clearable
        class="skill-input"
        @keyup.enter="addSkill"
        @select="handleSelect"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
        <template #append>
          <el-button type="primary" @click="addSkill" :disabled="!inputValue.trim()">
            添加
          </el-button>
        </template>
      </el-autocomplete>
    </div>

    <div class="skills-tags" v-if="skills.length > 0">
      <el-tag
        v-for="(skill, index) in skills"
        :key="index"
        closable
        size="large"
        type="primary"
        effect="light"
        @close="removeSkill(index)"
      >
        {{ skill }}
      </el-tag>
    </div>

    <el-empty v-else description="暂无技能标签，请在上方输入框添加" :image-size="60" />

    <div class="skill-suggestions" v-if="showSuggestions">
      <div class="suggestions-header">
        <span>热门技能推荐</span>
      </div>
      <div class="suggestions-list">
        <el-tag
          v-for="suggestion in filteredSuggestions"
          :key="suggestion"
          size="default"
          effect="plain"
          class="suggestion-tag"
          @click="addSuggestion(suggestion)"
        >
          <el-icon><Plus /></el-icon> {{ suggestion }}
        </el-tag>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { Search, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  modelValue: string[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string[]): void
}>()

// 技能列表
const skills = ref<string[]>([])

// 输入值
const inputValue = ref('')

// 常用技能建议
const commonSkills = [
  // 编程语言
  'JavaScript', 'TypeScript', 'Python', 'Java', 'Go', 'C++', 'C#', 'PHP', 'Ruby', 'Swift', 'Kotlin', 'Rust',
  // 前端
  'Vue.js', 'React', 'Angular', 'HTML5', 'CSS3', 'Sass', 'Less', 'Webpack', 'Vite', 'Node.js',
  // 后端
  'Spring Boot', 'Django', 'Flask', 'Express.js', 'Gin', 'Laravel', 'ASP.NET',
  // 数据库
  'MySQL', 'PostgreSQL', 'MongoDB', 'Redis', 'Elasticsearch', 'Oracle', 'SQL Server',
  // 云服务/DevOps
  'Docker', 'Kubernetes', 'AWS', 'Azure', 'GCP', 'Linux', 'Git', 'CI/CD', 'Jenkins',
  // 数据/AI
  '机器学习', '深度学习', 'TensorFlow', 'PyTorch', '数据分析', 'Pandas', 'NumPy',
  // 设计
  'Figma', 'Sketch', 'Photoshop', 'UI设计', 'UX设计',
  // 产品/管理
  '项目管理', '敏捷开发', 'Scrum', '产品设计', '需求分析',
  // 软技能
  '团队协作', '沟通能力', '问题解决', '领导力', '英语'
]

// 是否显示建议
const showSuggestions = computed(() => {
  return filteredSuggestions.value.length > 0
})

// 过滤后的建议（排除已添加的技能）
const filteredSuggestions = computed(() => {
  return commonSkills
    .filter(skill => !skills.value.includes(skill))
    .slice(0, 20)
})

// 监听 props 变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal && Array.isArray(newVal)) {
      skills.value = [...newVal]
    } else {
      skills.value = []
    }
  },
  { immediate: true, deep: true }
)

// 自动完成搜索
const querySearch = (queryString: string, cb: (results: { value: string }[]) => void) => {
  const results = queryString
    ? commonSkills
        .filter(skill => 
          skill.toLowerCase().includes(queryString.toLowerCase()) &&
          !skills.value.includes(skill)
        )
        .map(skill => ({ value: skill }))
    : []
  cb(results)
}

// 处理选择建议
const handleSelect = (item: { value: string }) => {
  addSkillByName(item.value)
  inputValue.value = ''
}

// 添加技能
const addSkill = () => {
  const skillName = inputValue.value.trim()
  if (skillName) {
    addSkillByName(skillName)
    inputValue.value = ''
  }
}

// 通过名称添加技能
const addSkillByName = (skillName: string) => {
  if (!skillName) return
  
  // 检查是否已存在
  if (skills.value.some(s => s.toLowerCase() === skillName.toLowerCase())) {
    ElMessage.warning('该技能已添加')
    return
  }
  
  // 检查数量限制
  if (skills.value.length >= 30) {
    ElMessage.warning('最多添加30个技能')
    return
  }
  
  skills.value.push(skillName)
  emitUpdate()
}

// 添加建议的技能
const addSuggestion = (skill: string) => {
  addSkillByName(skill)
}

// 删除技能
const removeSkill = (index: number) => {
  skills.value.splice(index, 1)
  emitUpdate()
}

// 发送更新事件
const emitUpdate = () => {
  emit('update:modelValue', [...skills.value])
}

// 验证（技能是可选的，所以总是返回 true）
const validate = async (): Promise<boolean> => {
  return true
}

// 重置
const resetForm = () => {
  skills.value = []
  inputValue.value = ''
  emitUpdate()
}

// 暴露方法给父组件
defineExpose({
  validate,
  resetForm
})
</script>

<style scoped lang="scss">
.skills-form {
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

    .skill-count {
      font-size: 13px;
      color: #94a3b8;
    }
  }

  .skills-input-area {
    margin-bottom: 16px;

    .skill-input {
      width: 100%;

      :deep(.el-input-group__append) {
        padding: 0;
        
        .el-button {
          margin: 0;
          border-radius: 0 4px 4px 0;
        }
      }
    }
  }

  .skills-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 24px;

    .el-tag {
      font-size: 14px;
      padding: 8px 12px;
      height: auto;
    }
  }

  .skill-suggestions {
    background: #f8fafc;
    border-radius: 8px;
    padding: 16px;

    .suggestions-header {
      font-size: 14px;
      font-weight: 500;
      color: #64748b;
      margin-bottom: 12px;
    }

    .suggestions-list {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;

      .suggestion-tag {
        cursor: pointer;
        transition: all 0.2s;

        &:hover {
          background: #0ea5e9;
          color: white;
          border-color: #0ea5e9;
        }

        .el-icon {
          margin-right: 4px;
        }
      }
    }
  }
}
</style>
