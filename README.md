# 智能人才运营平台

一个基于微服务架构的智能人才运营平台，采用 Go 微服务后端 + Vue3 前端的现代化技术栈。

## 项目概述

本项目是一个完整的毕业设计项目，实现了人才招聘管理的核心功能，包括人才管理、职位管理、智能推荐、简历解析、消息通知、面试管理等模块。

### 主要特性

- ✨ **微服务架构**: 8个独立微服务，易于扩展和维护
- 🤖 **智能推荐**: 基于多维度匹配的人岗推荐算法（技能、经验、学历、地理位置、薪资）
- 📊 **数据可视化**: ECharts 实现的数据分析大屏
- 🎨 **现代化 UI**: 渐变色设计、响应式布局
- 🔐 **安全认证**: JWT 身份验证和角色权限控制
- 📱 **多角色支持**: 管理员、HR、面试官、候选人等角色
- 📅 **面试管理**: 完整的面试安排、反馈、改期流程
- 🔔 **实时通知**: WebSocket 实时消息推送
- 📤 **数据导出**: 支持 Excel/CSV 格式导出
- 📥 **批量导入**: 支持批量导入人才数据
- 🐳 **容器化部署**: Docker Compose 一键部署
- 📝 **API 文档**: OpenAPI 3.0 规范文档
- ✅ **测试覆盖**: 单元测试和集成测试

## 技术栈

### 后端
- **语言**: Golang 1.21+
- **框架**: Gin
- **数据库**: PostgreSQL + Redis
- **ORM**: GORM
- **认证**: JWT
- **日志**: Zap
- **WebSocket**: Gorilla WebSocket
- **测试**: Go Testing + Testify

### 前端
- **框架**: Vue 3 (Composition API)
- **语言**: TypeScript 5+
- **构建工具**: Vite
- **UI 组件**: Element Plus
- **状态管理**: Pinia
- **图表**: ECharts
- **测试**: Vitest
- **导出**: XLSX

### DevOps
- **容器化**: Docker + Docker Compose
- **CI/CD**: GitHub Actions (可选)
- **API 文档**: OpenAPI 3.0 / Swagger

## 项目结构

```
talent-platform/
├── backend/                    # 后端微服务
│   ├── gateway/               # API 网关 (:8080)
│   ├── user-service/          # 用户服务 (:8081)
│   ├── talent-service/        # 人才服务 (:8082)
│   ├── job-service/           # 职位服务 (:8083)
│   ├── resume-service/        # 简历服务 (:8084)
│   ├── recommendation-service/ # 推荐服务 (:8085)
│   ├── message-service/       # 消息服务 (:8086)
│   ├── interview-service/     # 面试服务 (:8087)
│   ├── common/                # 公共模块
│   ├── database/              # 数据库脚本
│   └── docs/                  # API 文档
├── frontend/                   # 前端应用 (:3000)
├── docker-compose.yml         # Docker 编排
├── Makefile                   # 开发命令
└── README.md
```

## 快速开始

### 方式一：Docker 部署（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd talent-platform

# 一键启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

访问地址：
- 前端应用: http://localhost:3000
- API 网关: http://localhost:8080
- API 文档: http://localhost:8080/swagger

### 方式二：本地开发

#### 前置要求
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Redis 6+ (可选)

#### 1. 初始化数据库

```bash
# 创建数据库
psql -U postgres -c "CREATE DATABASE talent_platform;"

# 导入表结构和测试数据
psql -U postgres -d talent_platform -f backend/database/schema.sql
psql -U postgres -d talent_platform -f backend/database/mock_data.sql
```

#### 2. 启动后端服务

```bash
# 使用 Makefile
make dev-backend

# 或使用启动脚本
./start-backend.sh
```

#### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

### 使用 Makefile

```bash
make help          # 查看所有命令
make dev           # 启动开发环境
make test          # 运行测试
make docker-up     # Docker 启动
make db-init       # 初始化数据库
```

## 测试账号

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| admin | password123 | 超级管理员 | 拥有所有权限 |
| hr_zhang | password123 | HR主管 | HR总监 |
| hr_li | password123 | 招聘专员 | 招聘专员 |
| tech_chen | password123 | 面试官 | 技术总监 |
| tech_liu | password123 | 面试官 | 前端负责人 |
| viewer_test | password123 | 只读用户 | 仅查看权限 |

## 核心功能

### 1. 用户管理
- 用户注册 / 登录
- 多角色权限（管理员、HR、面试官、候选人）
- 个人信息管理
- 操作日志记录

### 2. 人才管理
- 人才信息 CRUD
- 高级搜索（技能、经验、地区）
- 标签分类
- 批量导入/导出

### 3. 职位管理
- 职位发布 / 编辑
- 职位状态管理（开放、关闭、已填补）
- 职位统计
- 紧急招聘标记

### 4. 简历管理
- 简历上传
- 简历解析
- 申请状态跟踪
- 匹配度评分

### 5. 面试管理
- 面试安排（初试、复试、终面、HR面）
- 面试日历视图
- 面试反馈录入
- 面试改期/取消
- 面试官日程管理

### 6. 智能推荐
- 为人才推荐合适职位
- 为职位推荐匹配候选人
- 多维度匹配算法：
  - 技能匹配 (50%)
  - 经验匹配 (20%)
  - 地理位置 (15%)
  - 学历匹配 (10%)
  - 薪资匹配 (5%)

### 7. 消息通知
- 站内消息
- 未读消息提醒
- WebSocket 实时推送
- 面试提醒通知

### 8. 数据可视化
- 仪表板统计
- 招聘漏斗分析
- 渠道效果分析
- 部门招聘进度
- 趋势图表

### 9. 数据导出
- 人才列表导出 Excel/CSV
- 职位列表导出
- 面试记录导出
- 自定义导出字段

## API 文档

API 文档遵循 OpenAPI 3.0 规范，位于 `backend/docs/swagger.yaml`。

所有 API 通过网关统一访问：`http://localhost:8080/api/v1`

主要接口：
- `POST /login` - 用户登录
- `POST /register` - 用户注册
- `GET /talents` - 获取人才列表
- `GET /jobs` - 获取职位列表
- `GET /interviews` - 获取面试列表
- `POST /recommendations/jobs-for-talent` - 为人才推荐职位
- `POST /recommendations/talents-for-job` - 为职位推荐人才

详细文档请查看 [API 文档](backend/docs/swagger.yaml)

## 测试

### 后端测试

```bash
# 运行所有后端测试
make test-backend

# 运行单个服务测试
cd backend/interview-service
go test -v ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 前端测试

```bash
cd frontend

# 运行测试
npm run test

# 监听模式
npm run test:watch

# 生成覆盖率报告
npm run test:coverage
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 生产环境建议

1. 使用 Nginx 作为反向代理
2. 配置 HTTPS 证书
3. 使用独立的 PostgreSQL 和 Redis 服务器
4. 配置日志收集（如 ELK）
5. 配置监控告警（如 Prometheus + Grafana）
6. 定期备份数据库

## 服务端口

| 服务名称 | 端口 | 说明 |
|----------|------|------|
| frontend | 3000 | 前端应用 |
| gateway | 8080 | API网关 |
| user-service | 8081 | 用户服务 |
| talent-service | 8082 | 人才服务 |
| job-service | 8083 | 职位服务 |
| resume-service | 8084 | 简历服务 |
| recommendation-service | 8085 | 推荐服务 |
| message-service | 8086 | 消息服务 |
| interview-service | 8087 | 面试服务 |

## 常见问题

**Q: 后端服务启动失败？**
A: 检查 PostgreSQL 数据库是否正常运行，确保数据库连接配置正确

**Q: 前端无法连接后端？**
A: 确保 API 网关正在运行，检查 Vite 代理配置

**Q: Docker 启动失败？**
A: 检查端口是否被占用，确保 Docker 和 Docker Compose 已正确安装

**Q: 测试运行失败？**
A: 确保已安装测试依赖，后端需要 testify，前端需要 vitest

## License

MIT

## 作者

qinyang - 毕业设计项目 2024
