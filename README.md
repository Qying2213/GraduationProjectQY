# 智能人才运营平台

一个基于微服务架构的智能人才运营平台，采用 Go 微服务后端 + Vue3 前端的现代化技术栈。

## 项目概述

本项目是一个完整的毕业设计项目，实现了人才招聘管理的核心功能，包括人才管理、职位管理、智能推荐、简历解析、消息通知等模块。

### 主要特性

- ✨ **微服务架构**: 7个独立微服务，易于扩展和维护
- 🤖 **智能推荐**: 基于技能匹配的人岗推荐算法
- 📊 **数据可视化**: ECharts 实现的数据分析大屏
- 🎨 **现代化 UI**: 渐变色设计、响应式布局
- 🔐 **安全认证**: JWT 身份验证和角色权限控制
- 📱 **多角色支持**: 管理员、HR、候选人三种角色

## 技术栈

### 后端
- **语言**: Golang 1.21+
- **框架**: Gin
- **数据库**: PostgreSQL + Redis
- **ORM**: GORM
- **认证**: JWT

### 前端
- **框架**: Vue 3 (Composition API)
- **语言**: TypeScript 5+
- **构建工具**: Vite
- **UI 组件**: Element Plus
- **状态管理**: Pinia
- **图表**: ECharts

## 项目结构

```
abc/
├── backend/                    # 后端微服务
│   ├── gateway/               # API 网关 (:8080)
│   ├── user-service/          # 用户服务 (:8081)
│   ├── talent-service/        # 人才服务 (:8082)
│   ├── job-service/           # 职位服务 (:8083)
│   ├── resume-service/        # 简历服务 (:8084)
│   ├── recommendation-service/ # 推荐服务 (:8085)
│   ├── message-service/       # 消息服务 (:8086)
│   └── common/                # 公共模块
└── frontend/                   # 前端应用 (:3000)
    ├── src/
    │   ├── api/               # API 接口
    │   ├── components/        # 组件
    │   ├── views/             # 页面
    │   ├── router/            # 路由
    │   ├── store/             # 状态管理
    │   └── styles/            # 样式
    └── public/
```

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 12+
- Redis 6+

### 1. 数据库配置

创建 PostgreSQL 数据库：

```sql
CREATE DATABASE talent_platform;
```

### 2. 启动后端服务

每个服务需要独立启动：

```bash
# 1. 用户服务
cd backend/user-service
go mod tidy
go run main.go

# 2. 人才服务
cd backend/talent-service
go mod tidy
go run main.go

# 3. 职位服务
cd backend/job-service
go mod tidy
go run main.go

# 4. 简历服务
cd backend/resume-service
go mod tidy
go run main.go

# 5. 推荐服务
cd backend/recommendation-service
go mod tidy
go run main.go

# 6. 消息服务
cd backend/message-service
go mod tidy
go run main.go

# 7. API 网关
cd backend/gateway
go mod tidy
go run main.go
```

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

访问 http://localhost:3000

## 核心功能

### 1. 用户管理
- 用户注册 / 登录
- 多角色权限（管理员、HR、候选人）
- 个人信息管理

### 2. 人才管理
- 人才信息 CRUD
- 高级搜索（技能、经验、地区）
- 标签分类

### 3. 职位管理
- 职位发布 / 编辑
- 职位状态管理（开放、关闭、已填补）
- 职位统计

### 4. 简历管理
- 简历上传
- 简历解析
- 申请状态跟踪

### 5. 智能推荐
- 为人才推荐合适职位
- 为职位推荐匹配候选人
- 基于技能、经验、地理位置的匹配算法

### 6. 消息通知
- 站内消息
- 未读消息提醒
- 消息已读标记

### 7. 数据可视化
- 仪表板统计
- ECharts 图表展示
- 最近活动时间线

## API 文档

所有 API 通过网关统一访问：`http://localhost:8080/api/v1`

详细 API 文档请查看：
- [后端 API 文档](backend/README.md)
- [前端 README](frontend/README.md)

## 开发指南

### 后端开发

1. 每个微服务独立开发和部署
2. 使用 GORM 进行数据库操作
3. 遵循 RESTful API 规范
4. 统一的错误处理和响应格式

### 前端开发

1. 使用 Vue 3 Composition API
2. TypeScript 严格模式
3. 组件化开发
4. 响应式设计

## 推荐算法

智能推荐算法基于以下维度：

- **技能匹配** (60%权重): 候选人技能与职位要求的匹配度
- **经验匹配** (20%权重): 工作年限与职位级别的适配性
- **地理位置** (20%权重): 候选人位置与职位地点的匹配

## 部署说明

### 开发环境
- 所有服务在本地独立运行
- 使用 API 网关统一转发

### 生产环境
- 建议使用 Docker 容器化部署
- 使用 Nginx 作为反向代理
- PostgreSQL 和 Redis 使用独立服务器

## 常见问题

**Q: 后端服务启动失败？**
A: 检查 PostgreSQL 数据库是否正常运行，确保数据库连接配置正确

**Q: 前端无法连接后端？**
A: 确保 API 网关正在运行，检查 Vite 代理配置

**Q: 推荐功能不工作？**
A: 推荐服务使用模拟数据，需要先有人才和职位数据

## 项目截图

（这里可以添加项目运行后的截图）

## License

MIT

## 作者

qinyang - 毕业设计项目 2024
