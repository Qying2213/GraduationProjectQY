# 智能人才招聘管理平台 - 代码规范

> 📖 返回 [项目首页](../README.md) | 相关文档：[系统架构](ARCHITECTURE.md) | [快速启动](QUICKSTART.md)

---

## 1. 项目结构

### 1.1 前端结构

```
frontend/src/
├── api/                    # API 接口封装
│   ├── user.ts            # 用户相关接口
│   ├── job.ts             # 职位相关接口
│   ├── talent.ts          # 人才相关接口
│   ├── resume.ts          # 简历相关接口
│   ├── interview.ts       # 面试相关接口
│   ├── message.ts         # 消息相关接口
│   └── recommendation.ts  # 推荐相关接口
├── components/             # 公共组件
│   ├── layout/            # 布局组件
│   ├── common/            # 通用组件
│   └── charts/            # 图表组件
├── views/                  # 页面视图
│   ├── auth/              # 认证页面
│   ├── dashboard/         # 仪表板
│   ├── talents/           # 人才管理
│   ├── jobs/              # 职位管理
│   ├── resumes/           # 简历管理
│   ├── recommend/         # 智能推荐
│   ├── interviews/        # 面试管理
│   ├── calendar/          # 面试日历
│   ├── kanban/            # 招聘看板
│   ├── messages/          # 消息中心
│   ├── reports/           # 数据报表
│   ├── portal/            # 求职者门户
│   ├── profile/           # 个人中心
│   └── system/            # 系统设置
├── store/                  # Pinia 状态管理
├── router/                 # 路由配置
├── types/                  # TypeScript 类型定义
├── utils/                  # 工具函数
└── styles/                 # 全局样式
```

### 1.2 后端结构

```
backend/
├── gateway/                # API 网关
├── user-service/           # 用户服务
│   ├── handlers/          # 请求处理器
│   ├── models/            # 数据模型
│   └── main.go            # 入口文件
├── job-service/            # 职位服务
├── interview-service/      # 面试服务
├── resume-service/         # 简历服务
├── message-service/        # 消息服务
├── talent-service/         # 人才服务
├── recommendation-service/ # 推荐服务
├── log-service/            # 日志服务
├── evaluator-service/      # AI评估服务
│   ├── cmd/server/        # 入口
│   ├── internal/          # 内部模块
│   │   ├── api/          # API 路由
│   │   ├── config/       # 配置
│   │   ├── database/     # 数据库
│   │   ├── repository/   # 数据访问
│   │   ├── service/      # 业务逻辑
│   │   └── thirdparty/   # 第三方集成
│   └── pkg/               # 公共包
├── common/                 # 公共模块
│   ├── config/            # 配置管理
│   ├── elasticsearch/     # ES 客户端
│   ├── middleware/        # 中间件
│   └── response/          # 统一响应
└── database/               # 数据库脚本
```

---

## 2. 命名规范

### 2.1 前端

| 类型 | 规范 | 示例 |
|------|------|------|
| 文件名 | PascalCase | `TalentList.vue` |
| 组件名 | PascalCase | `TalentCard` |
| 变量 | camelCase | `talentList` |
| 常量 | UPPER_SNAKE_CASE | `API_BASE_URL` |
| 类型 | PascalCase | `TalentInfo` |
| 接口 | PascalCase + I前缀 | `ITalent` |

### 2.2 后端

| 类型 | 规范 | 示例 |
|------|------|------|
| 包名 | 小写 | `handlers` |
| 文件名 | snake_case | `talent_handler.go` |
| 结构体 | PascalCase | `TalentHandler` |
| 方法 | PascalCase | `CreateTalent` |
| 变量 | camelCase | `talentList` |
| 常量 | PascalCase | `DefaultPageSize` |

---

## 3. 代码风格

### 3.1 TypeScript/Vue

```typescript
// 组件定义
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { TalentInfo } from '@/types'

// Props
const props = defineProps<{
  talentId: number
}>()

// Emits
const emit = defineEmits<{
  (e: 'update', talent: TalentInfo): void
}>()

// 响应式数据
const loading = ref(false)
const talent = ref<TalentInfo | null>(null)

// 计算属性
const fullName = computed(() => talent.value?.name || '')

// 方法
const fetchTalent = async () => {
  loading.value = true
  try {
    // ...
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(() => {
  fetchTalent()
})
</script>
```

### 3.2 Go

```go
// Handler 定义
type TalentHandler struct {
    db *gorm.DB
}

func NewTalentHandler(db *gorm.DB) *TalentHandler {
    return &TalentHandler{db: db}
}

// 方法实现
func (h *TalentHandler) CreateTalent(c *gin.Context) {
    var talent models.Talent
    if err := c.ShouldBindJSON(&talent); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    1,
            "message": err.Error(),
        })
        return
    }

    if err := h.db.Create(&talent).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    1,
            "message": "创建失败",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    0,
        "message": "success",
        "data":    talent,
    })
}
```

---

## 4. API 规范

### 4.1 RESTful 设计

| 操作 | HTTP 方法 | 路径示例 |
|------|----------|---------|
| 列表 | GET | /api/v1/talents |
| 详情 | GET | /api/v1/talents/:id |
| 创建 | POST | /api/v1/talents |
| 更新 | PUT | /api/v1/talents/:id |
| 删除 | DELETE | /api/v1/talents/:id |

### 4.2 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

---

## 5. Git 规范

### 5.1 分支命名

| 类型 | 格式 | 示例 |
|------|------|------|
| 功能 | feature/xxx | feature/talent-search |
| 修复 | fix/xxx | fix/login-error |
| 优化 | refactor/xxx | refactor/api-structure |

### 5.2 提交信息

```
<type>(<scope>): <subject>

feat(talent): 添加人才搜索功能
fix(auth): 修复登录token过期问题
docs(readme): 更新部署文档
```

---

## 📚 相关文档

| 文档 | 说明 |
|------|------|
| [📖 项目首页](../README.md) | 项目概述 |
| [📐 系统架构](ARCHITECTURE.md) | 架构设计 |
| [🚀 快速启动](QUICKSTART.md) | 环境配置 |
| [🧪 测试指南](TEST_GUIDE.md) | 测试方法 |
