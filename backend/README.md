# 智能人才运营平台 - 后端微服务

## 项目结构

```
backend/
├── gateway/                    # API网关 (:8080)
├── user-service/              # 用户服务 (:8081)
├── talent-service/            # 人才服务 (:8082)
├── job-service/               # 职位服务 (:8083)
├── resume-service/            # 简历服务 (:8084)
├── recommendation-service/    # 推荐服务 (:8085)
├── message-service/           # 消息服务 (:8086)
└── common/                    # 公共模块
    ├── database/              # 数据库连接
    ├── middleware/            # 中间件
    └── response/              # 响应工具
```

## 快速开始

### 前置要求

- Go 1.21+
- PostgreSQL 12+
- Redis 6+

### 数据库配置

1. 创建数据库：
```sql
CREATE DATABASE talent_platform;
```

2. 配置连接（在各服务的main.go中）：
```go
dsn := "host=localhost user=postgres password=postgres dbname=talent_platform port=5432 sslmode=disable TimeZone=Asia/Shanghai"
```

### 启动服务

每个服务独立启动，按以下顺序：

```bash
# 1. 启动用户服务
cd user-service
go mod tidy
go run main.go

# 2. 启动人才服务
cd talent-service
go mod tidy
go run main.go

# 3. 启动职位服务
cd job-service
go mod tidy
go run main.go

# 4. 启动简历服务
cd resume-service
go mod tidy
go run main.go

# 5. 启动推荐服务
cd recommendation-service
go mod tidy
go run main.go

# 6. 启动消息服务
cd message-service
go mod tidy
go run main.go

# 7. 启动API网关
cd gateway
go mod tidy
go run main.go
```

## API文档

### 统一访问地址
所有API通过网关访问：`http://localhost:8080/api/v1`

### 用户服务 API

#### 注册
```
POST /api/v1/register
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "role": "hr",
  "real_name": "张三",
  "phone": "13800138000"
}
```

#### 登录
```
POST /api/v1/login
{
  "username": "testuser",
  "password": "password123"
}
```

### 人才服务 API

#### 创建人才
```
POST /api/v1/talents
{
  "name": "李四",
  "email": "lisi@example.com",
  "phone": "13900139000",
  "skills": ["Go", "Vue", "Docker"],
  "experience": 5,
  "education": "本科",
  "location": "北京",
  "salary": "20-30K"
}
```

#### 获取人才列表
```
GET /api/v1/talents?page=1&page_size=10&status=active&search=李四
```

#### 高级搜索
```
GET /api/v1/talents/search?skills=Go&skills=Vue&min_experience=3&max_experience=8&location=北京
```

### 职位服务 API

#### 创建职位
```
POST /api/v1/jobs
{
  "title": "Go开发工程师",
  "description": "负责后端开发",
  "requirements": ["3年以上Go经验", "熟悉微服务"],
  "salary": "20-30K",
  "location": "北京",
  "type": "full-time",
  "skills": ["Go", "Docker", "Kubernetes"],
  "department": "技术部",
  "level": "mid"
}
```

#### 获取职位列表
```
GET /api/v1/jobs?page=1&page_size=10&status=open&type=full-time&location=北京
```

#### 职位统计
```
GET /api/v1/jobs/stats
```

### 简历服务 API

#### 上传简历
```
POST /api/v1/resumes
{
  "talent_id": 1,
  "file_name": "resume.pdf",
  "file_url": "/uploads/resume.pdf",
  "file_size": 102400
}
```

#### 创建申请
```
POST /api/v1/applications
{
  "job_id": 1,
  "talent_id": 1,
  "resume_id": 1,
  "cover_letter": "我对这个职位很感兴趣..."
}
```

### 推荐服务 API

#### 为人才推荐职位
```
POST /api/v1/recommendations/jobs-for-talent
{
  "id": 1,
  "name": "张三",
  "skills": ["Go", "Docker"],
  "experience": 5,
  "location": "北京"
}
```

#### 为职位推荐人才
```
POST /api/v1/recommendations/talents-for-job
{
  "id": 1,
  "title": "Go开发工程师",
  "skills": ["Go", "Docker"],
  "location": "北京"
}
```

### 消息服务 API

#### 发送消息
```
POST /api/v1/messages
{
  "from_id": 1,
  "to_id": 2,
  "title": "面试邀请",
  "content": "您好，我们邀请您参加面试...",
  "type": "notification"
}
```

#### 获取消息列表
```
GET /api/v1/messages?user_id=1&page=1&page_size=20&is_read=false
```

#### 未读消息数
```
GET /api/v1/messages/unread-count?user_id=1
```

## 服务端口

- API Gateway: 8080
- User Service: 8081
- Talent Service: 8082
- Job Service: 8083
- Resume Service: 8084
- Recommendation Service: 8085
- Message Service: 8086

## 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis (推荐服务)
- **认证**: JWT
