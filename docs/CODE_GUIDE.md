# 项目代码指南 - 快速上手

本文档帮助你快速理解项目代码结构，方便修改和扩展。

---

## 一、项目整体结构

```
├── backend/                 # 后端（Go 微服务）
│   ├── gateway/            # API 网关 - 所有请求入口
│   ├── user-service/       # 用户服务 - 登录注册
│   ├── talent-service/     # 人才服务 - 人才库管理
│   ├── job-service/        # 职位服务 - 职位管理
│   ├── resume-service/     # 简历服务 - 简历上传解析
│   ├── recommendation-service/  # 推荐服务 - AI匹配
│   ├── interview-service/  # 面试服务 - 面试安排
│   ├── message-service/    # 消息服务 - 通知
│   ├── common/             # 公共模块
│   └── database/           # 数据库脚本
│
├── frontend/               # 前端（Vue3）
│   └── src/
│       ├── api/           # 接口请求
│       ├── views/         # 页面
│       ├── components/    # 组件
│       ├── store/         # 状态管理
│       ├── router/        # 路由
│       └── utils/         # 工具函数
│
└── docs/                   # 文档
```

---

## 二、后端代码结构

### 每个微服务的结构（以 user-service 为例）

```
backend/user-service/
├── main.go              # 入口文件，启动服务
├── handlers/            # 处理 HTTP 请求
│   └── user_handler.go  # 用户相关接口
├── models/              # 数据模型
│   └── user.go          # User 结构体
├── go.mod               # Go 依赖
└── go.sum
```

### 核心文件说明

#### 1. main.go - 服务入口
```go
// 每个服务的 main.go 做3件事：
// 1. 连接数据库
// 2. 注册路由
// 3. 启动 HTTP 服务

func main() {
    // 连接数据库
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    
    // 创建 Gin 路由
    r := gin.Default()
    
    // 注册接口
    handler := handlers.NewUserHandler(db)
    r.POST("/api/v1/login", handler.Login)
    r.GET("/api/v1/users", handler.ListUsers)
    
    // 启动服务
    r.Run(":8081")
}
```

#### 2. handlers/ - 接口处理
```go
// handlers/user_handler.go
type UserHandler struct {
    DB *gorm.DB
}

// 登录接口
func (h *UserHandler) Login(c *gin.Context) {
    var req LoginRequest
    c.ShouldBindJSON(&req)  // 解析请求体
    
    // 查询数据库
    var user models.User
    h.DB.Where("username = ?", req.Username).First(&user)
    
    // 返回响应
    c.JSON(200, gin.H{"token": token, "user": user})
}
```

#### 3. models/ - 数据模型
```go
// models/user.go
type User struct {
    ID        uint   `json:"id" gorm:"primaryKey"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password  string `json:"-"`  // json:"-" 表示不返回给前端
    Role      string `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 各服务端口

| 服务 | 端口 | 主要接口 |
|------|------|----------|
| gateway | 8080 | 转发所有请求 |
| user-service | 8081 | /api/v1/login, /api/v1/users |
| talent-service | 8082 | /api/v1/talents |
| job-service | 8083 | /api/v1/jobs |
| resume-service | 8084 | /api/v1/resumes |
| recommendation-service | 8085 | /api/v1/recommend |
| message-service | 8086 | /api/v1/messages |
| interview-service | 8087 | /api/v1/interviews |

---

## 三、前端代码结构

```
frontend/src/
├── api/                    # API 接口封装
│   ├── auth.ts            # 登录注册接口
│   ├── talent.ts          # 人才接口
│   ├── job.ts             # 职位接口
│   └── ...
│
├── views/                  # 页面组件
│   ├── auth/              # 登录注册页
│   │   └── Login.vue
│   ├── dashboard/         # 仪表板
│   │   └── Dashboard.vue
│   ├── talents/           # 人才管理
│   │   ├── TalentList.vue
│   │   └── TalentDetail.vue
│   ├── jobs/              # 职位管理
│   ├── portal/            # 前台求职端
│   │   ├── PortalHome.vue
│   │   ├── PortalJobList.vue
│   │   └── ...
│   └── ...
│
├── components/             # 可复用组件
│   ├── layout/            # 布局组件
│   │   ├── MainLayout.vue    # 后台布局
│   │   └── PortalLayout.vue  # 前台布局
│   └── common/            # 通用组件
│
├── store/                  # Pinia 状态管理
│   ├── user.ts            # 用户状态
│   ├── permission.ts      # 权限状态
│   └── theme.ts           # 主题状态
│
├── router/                 # 路由配置
│   └── index.ts
│
├── types/                  # TypeScript 类型
│   └── index.ts
│
├── styles/                 # 全局样式
│   └── global.scss
│
└── utils/                  # 工具函数
    └── request.ts         # Axios 封装
```

### 核心文件说明

#### 1. api/ - 接口封装
```typescript
// api/talent.ts
import request from '@/utils/request'

export const talentApi = {
  // 获取人才列表
  getList(params: any) {
    return request.get('/api/v1/talents', { params })
  },
  
  // 创建人才
  create(data: any) {
    return request.post('/api/v1/talents', data)
  },
  
  // 更新人才
  update(id: number, data: any) {
    return request.put(`/api/v1/talents/${id}`, data)
  },
  
  // 删除人才
  delete(id: number) {
    return request.delete(`/api/v1/talents/${id}`)
  }
}
```

#### 2. views/ - 页面组件
```vue
<!-- views/talents/TalentList.vue -->
<template>
  <div class="talent-list">
    <!-- 搜索栏 -->
    <el-form :model="searchForm" inline>
      <el-form-item label="姓名">
        <el-input v-model="searchForm.name" />
      </el-form-item>
      <el-button @click="handleSearch">搜索</el-button>
    </el-form>
    
    <!-- 表格 -->
    <el-table :data="tableData">
      <el-table-column prop="name" label="姓名" />
      <el-table-column prop="skills" label="技能" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button @click="handleEdit(row)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { talentApi } from '@/api/talent'

const tableData = ref([])
const searchForm = ref({ name: '' })

// 获取数据
const fetchData = async () => {
  const res = await talentApi.getList(searchForm.value)
  tableData.value = res.data.data.talents
}

onMounted(() => {
  fetchData()
})
</script>
```

#### 3. store/ - 状态管理
```typescript
// store/user.ts
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || ''
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin'
  },
  
  actions: {
    async login(username: string, password: string) {
      const res = await authApi.login({ username, password })
      this.token = res.data.data.token
      this.user = res.data.data.user
      localStorage.setItem('token', this.token)
    },
    
    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('token')
    }
  }
})
```

#### 4. router/ - 路由配置
```typescript
// router/index.ts
const routes = [
  {
    path: '/login',
    component: () => import('@/views/auth/Login.vue')
  },
  {
    path: '/',
    component: () => import('@/components/layout/MainLayout.vue'),
    children: [
      { path: 'dashboard', component: () => import('@/views/dashboard/Dashboard.vue') },
      { path: 'talents', component: () => import('@/views/talents/TalentList.vue') },
      // ...
    ]
  },
  {
    path: '/portal',
    component: () => import('@/components/layout/PortalLayout.vue'),
    children: [
      { path: '', component: () => import('@/views/portal/PortalHome.vue') },
      // ...
    ]
  }
]
```

---

## 四、常见修改场景

### 场景1：修改页面样式

找到对应的 `.vue` 文件，修改 `<style>` 部分：
```
frontend/src/views/xxx/XxxPage.vue
```

### 场景2：修改接口返回数据

找到对应服务的 handler：
```
backend/xxx-service/handlers/xxx_handler.go
```

### 场景3：添加新页面

1. 创建页面文件：`frontend/src/views/xxx/NewPage.vue`
2. 添加路由：`frontend/src/router/index.ts`
3. 如需新接口，添加 API：`frontend/src/api/xxx.ts`

### 场景4：添加新接口

1. 后端添加 handler 方法
2. 在 main.go 注册路由
3. 前端 api/ 目录添加调用方法

### 场景5：修改数据库表

1. 修改 `backend/database/schema.sql`
2. 修改对应服务的 `models/xxx.go`
3. 重新导入数据库

---

## 五、关键文件速查表

| 要改什么 | 文件位置 |
|----------|----------|
| 登录逻辑 | `backend/user-service/handlers/user_handler.go` |
| 登录页面 | `frontend/src/views/auth/Login.vue` |
| 权限控制 | `frontend/src/store/permission.ts` |
| 路由配置 | `frontend/src/router/index.ts` |
| 全局样式 | `frontend/src/styles/global.scss` |
| 主题颜色 | `frontend/src/store/theme.ts` |
| 数据库表 | `backend/database/schema.sql` |
| 测试数据 | `backend/database/mock_data.sql` |
| 简历解析 | `backend/resume-service/parser/resume_parser.go` |
| AI智能评估 | `backend/resume-service/evaluator/coze_evaluator.go` |
| 推荐算法 | `backend/recommendation-service/handlers/` |

---

## 六、调试技巧

### 前端调试
```bash
# 启动开发服务器（支持热更新）
cd frontend && npm run dev

# 浏览器打开 F12 查看 Console 和 Network
```

### 后端调试
```bash
# 查看某个服务的日志
cd backend/user-service && go run main.go

# 日志会打印在终端
```

### 数据库调试
```bash
# 连接数据库
psql -U postgres -d talent_platform

# 查看表
\dt

# 查询数据
SELECT * FROM users;
```

---

## 七、技术栈速查

| 技术 | 用途 | 文档 |
|------|------|------|
| Vue 3 | 前端框架 | https://vuejs.org |
| Element Plus | UI 组件 | https://element-plus.org |
| Pinia | 状态管理 | https://pinia.vuejs.org |
| TypeScript | 类型检查 | https://typescriptlang.org |
| Go | 后端语言 | https://go.dev |
| Gin | Web 框架 | https://gin-gonic.com |
| GORM | ORM | https://gorm.io |
| PostgreSQL | 数据库 | https://postgresql.org |

---

有问题随时问！
