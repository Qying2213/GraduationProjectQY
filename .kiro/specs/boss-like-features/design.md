# Design Document: Boss-Like Features

## Overview

本设计文档描述在现有智能人才运营平台基础上扩展的核心招聘功能。项目已有完整的微服务架构和基础功能，本次设计聚焦于**增量开发**，在现有代码基础上添加新功能。

### 现有功能基础

项目已实现的功能：
- **用户系统**：登录注册、JWT认证、角色权限（admin/hr/candidate）
- **职位管理**：职位CRUD、列表查询、搜索筛选（job-service）
- **简历管理**：简历上传、解析、AI评估（resume-service）
- **人才管理**：人才档案、标签管理（talent-service）
- **面试管理**：面试安排、日历视图（interview-service）
- **消息通知**：系统消息、WebSocket推送（message-service）
- **前端页面**：求职者端（PortalJobList、MyApplications、MyResume）、HR后台

### 本次扩展目标

在现有基础上完善以下功能：
1. **完善申请流程**：增强投递功能、状态追踪、通知机制
2. **在线简历编辑**：支持在线编辑简历信息（目前只有上传功能）
3. **即时聊天功能**：求职者与HR之间的实时聊天（新增）
4. **HR候选人管理**：增强筛选、状态管理功能

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Frontend (Vue3)                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │
│  │  求职者端     │  │  企业端/HR   │  │  聊天模块     │               │
│  │  Portal      │  │  Admin       │  │  Chat        │               │
│  └──────────────┘  └──────────────┘  └──────────────┘               │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Gateway (Gin + JWT)                             │
│              路由分发 / 认证鉴权 / 请求转发                           │
└─────────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│ job-service  │      │message-service│     │resume-service│
│  职位管理     │      │  消息/聊天    │      │  简历管理     │
└──────────────┘      └──────────────┘      └──────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      PostgreSQL + Redis                              │
│              业务数据存储 / 会话缓存 / 在线状态                        │
└─────────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 1. 后端组件

#### 1.1 Job Service 扩展（现有服务增强）

job-service 已有职位CRUD功能，需要扩展以下接口：

**现有接口（无需修改）：**
- `GET /jobs` - 职位列表
- `GET /jobs/:id` - 职位详情
- `POST /jobs` - 创建职位
- `PUT /jobs/:id` - 更新职位
- `DELETE /jobs/:id` - 删除职位

**新增接口：**

```go
// 获取职位申请列表（HR视角）
GET /jobs/:id/applications
Response: {
    "code": 0,
    "data": {
        "applications": [...],
        "total": 100
    }
}

// 获取职位统计（申请人数等）
GET /jobs/:id/stats
Response: {
    "code": 0,
    "data": {
        "applicants": 50,
        "views": 200
    }
}
```

#### 1.2 Application Service（现有功能增强）

applications 表和基础CRUD已存在于 job-service 中，需要增强以下功能：

**现有功能（需增强）：**
- 申请创建时自动发送通知给HR
- 申请状态更新时通知求职者
- 防止重复申请检查

**现有接口（需增强逻辑）：**

```go
// 创建申请
POST /applications
Request: {
    "job_id": 1,
    "talent_id": 1,
    "resume_id": 1,
    "cover_letter": "..."
}

// 获取申请列表
GET /applications?talent_id=1&job_id=1&status=pending

// 更新申请状态
PUT /applications/:id
Request: {
    "status": "interview",
    "notes": "..."
}

// 删除/撤回申请
DELETE /applications/:id
```

#### 1.3 Message Service 扩展（新增聊天功能）

message-service 已有系统消息和 WebSocket Hub，需要扩展支持即时聊天：

**现有功能：**
- 系统消息发送和接收
- WebSocket 连接管理（Hub）
- 消息已读状态

**新增数据模型：**

```go
// 会话表
type Conversation struct {
    ID            uint      `gorm:"primarykey"`
    ParticipantA  uint      `gorm:"not null"` // 用户A ID
    ParticipantB  uint      `gorm:"not null"` // 用户B ID
    LastMessageID *uint     // 最后一条消息ID
    LastMessageAt *time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// 聊天消息表
type ChatMessage struct {
    ID             uint      `gorm:"primarykey"`
    ConversationID uint      `gorm:"not null"`
    SenderID       uint      `gorm:"not null"`
    Content        string    `gorm:"type:text"`
    MessageType    string    `gorm:"default:'text'"` // text, image, file
    IsRead         bool      `gorm:"default:false"`
    CreatedAt      time.Time
}
```

**新增接口：**

```go
// 获取会话列表
GET /conversations
Response: {
    "code": 0,
    "data": {
        "conversations": [
            {
                "id": 1,
                "participant": {...},
                "last_message": {...},
                "unread_count": 3
            }
        ]
    }
}

// 获取/创建会话
POST /conversations
Request: {
    "participant_id": 2
}

// 获取会话消息
GET /conversations/:id/messages?page=1&page_size=20

// 发送消息
POST /conversations/:id/messages
Request: {
    "content": "Hello",
    "message_type": "text"
}

// 标记会话已读
PUT /conversations/:id/read
```

#### 1.4 WebSocket 扩展

扩展现有 Hub 支持聊天消息：

```go
// 消息类型
const (
    MsgTypeChat           = "chat"           // 聊天消息
    MsgTypeChatRead       = "chat_read"      // 消息已读
    MsgTypeOnlineStatus   = "online_status"  // 在线状态
    MsgTypeTyping         = "typing"         // 正在输入
)

// 聊天消息结构
type ChatWebSocketMessage struct {
    Type           string `json:"type"`
    ConversationID uint   `json:"conversation_id"`
    Message        struct {
        ID        uint   `json:"id"`
        SenderID  uint   `json:"sender_id"`
        Content   string `json:"content"`
        CreatedAt string `json:"created_at"`
    } `json:"message"`
}
```

#### 1.5 Resume Service 扩展（在线编辑）

resume-service 已有文件上传和AI解析功能，需要新增在线编辑功能：

**现有功能：**
- 简历文件上传
- 简历解析（AI）
- 简历列表查询

**新增接口：**

```go
// 获取在线简历数据
GET /resumes/online
Response: {
    "code": 0,
    "data": {
        "basic_info": {...},
        "work_experience": [...],
        "education": [...],
        "skills": [...]
    }
}

// 保存在线简历
PUT /resumes/online
Request: {
    "basic_info": {
        "name": "张三",
        "phone": "13800138000",
        "email": "test@example.com",
        "location": "北京"
    },
    "work_experience": [...],
    "education": [...],
    "skills": ["Go", "Python"]
}
```

#### 1.6 Interview Service 扩展（现有服务增强）

interview-service 已有面试管理功能，需要增强与申请的关联：

```go
// 从申请创建面试
POST /interviews
Request: {
    "application_id": 1,
    "date": "2024-01-20",
    "time": "14:00",
    "duration": 60,
    "method": "online",
    "location": "腾讯会议链接"
}
```

### 2. 前端组件

#### 2.1 求职者端页面（现有页面增强）

| 页面 | 路由 | 现状 | 需要增强 |
|------|------|------|---------|
| PortalJobList | /portal/jobs | ✅ 已实现 | 完善投递弹窗逻辑 |
| PortalJobDetail | /portal/jobs/:id | ✅ 已实现 | 增强投递功能 |
| MyApplications | /portal/my-applications | ✅ 已实现 | 完善状态追踪、时间线 |
| MyResume | /portal/my-resume | ✅ 已实现 | 新增在线编辑功能 |
| PortalChat | /portal/chat | ❌ 新增 | 聊天页面 |

#### 2.2 企业端/HR页面（现有页面增强）

| 页面 | 路由 | 现状 | 需要增强 |
|------|------|------|---------|
| JobList | /jobs | ✅ 已实现 | 显示申请人数 |
| JobDetail | /jobs/:id | ✅ 已实现 | 新增申请管理Tab |
| TalentList | /talents | ✅ 已实现 | 增强筛选功能 |
| TalentDetail | /talents/:id | ✅ 已实现 | 完善简历展示 |
| ChatCenter | /messages/chat | ❌ 新增 | HR聊天中心 |

#### 2.3 聊天组件

```typescript
// 聊天相关组件
components/
  chat/
    ConversationList.vue    // 会话列表
    ChatWindow.vue          // 聊天窗口
    MessageItem.vue         // 消息项
    ChatInput.vue           // 输入框
```

#### 2.4 API 模块

```typescript
// frontend/src/api/chat.ts
export const chatApi = {
    // 获取会话列表
    getConversations(params?: { page?: number; page_size?: number })
    
    // 创建/获取会话
    createConversation(participantId: number)
    
    // 获取消息列表
    getMessages(conversationId: number, params?: { page?: number; page_size?: number })
    
    // 发送消息
    sendMessage(conversationId: number, content: string, type?: string)
    
    // 标记已读
    markAsRead(conversationId: number)
    
    // 获取未读数
    getUnreadCount()
}
```

## Data Models

### 数据库表结构扩展

```sql
-- 会话表
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    participant_a INTEGER REFERENCES users(id) NOT NULL,
    participant_b INTEGER REFERENCES users(id) NOT NULL,
    last_message_id INTEGER,
    last_message_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(participant_a, participant_b)
);

-- 聊天消息表
CREATE TABLE IF NOT EXISTS chat_messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER REFERENCES conversations(id) NOT NULL,
    sender_id INTEGER REFERENCES users(id) NOT NULL,
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text',
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_conversations_participants ON conversations(participant_a, participant_b);
CREATE INDEX idx_conversations_last_message ON conversations(last_message_at DESC);
CREATE INDEX idx_chat_messages_conversation ON chat_messages(conversation_id, created_at DESC);
CREATE INDEX idx_chat_messages_unread ON chat_messages(conversation_id, is_read) WHERE is_read = FALSE;
```

### 前端类型定义扩展

```typescript
// frontend/src/types/chat.ts
export interface Conversation {
    id: number
    participant: User
    last_message?: ChatMessage
    unread_count: number
    updated_at: string
}

export interface ChatMessage {
    id: number
    conversation_id: number
    sender_id: number
    content: string
    message_type: 'text' | 'image' | 'file'
    is_read: boolean
    created_at: string
}

export interface ConversationListResponse {
    conversations: Conversation[]
    total: number
}
```

### Redis 缓存结构

```
# 用户在线状态
online:user:{user_id} -> "1" (TTL: 60s, 心跳续期)

# 会话未读数
unread:conversation:{conversation_id}:user:{user_id} -> count

# 用户总未读数
unread:total:user:{user_id} -> count
```



## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Job Search Filter Correctness

*For any* job search with keyword, location, or experience filters applied, all returned jobs SHALL contain the keyword in title or description, match the selected location, and match the experience level requirement respectively.

**Validates: Requirements 1.2, 1.3, 1.4**

### Property 2: Application Creation Status

*For any* valid application submission with job_id, talent_id, and resume_id, the created application record SHALL have status "pending".

**Validates: Requirements 2.2**

### Property 3: Duplicate Application Prevention

*For any* candidate-job pair where an application already exists, attempting to create another application SHALL be rejected and return an error.

**Validates: Requirements 2.4**

### Property 4: Resume Requirement for Application

*For any* candidate without at least one resume, attempting to submit an application SHALL be rejected with a validation error.

**Validates: Requirements 2.6**

### Property 5: Application Status Filter Correctness

*For any* application list query with a status filter, all returned applications SHALL have the specified status value.

**Validates: Requirements 3.4, 6.2**

### Property 6: Application Withdrawal Deletion

*For any* application withdrawal action, the application record SHALL be removed from the database and no longer appear in subsequent queries.

**Validates: Requirements 3.5**

### Property 7: Resume Data Persistence Round-Trip

*For any* resume update with valid data, saving and then retrieving the resume SHALL return data equivalent to what was submitted.

**Validates: Requirements 4.3**

### Property 8: File Upload Format Validation

*For any* resume file upload, files with formats other than PDF, DOC, DOCX or files larger than 10MB SHALL be rejected with a validation error.

**Validates: Requirements 4.7**

### Property 9: Job Creation Default Status

*For any* valid job creation submission, the created job record SHALL have status "open".

**Validates: Requirements 5.2**

### Property 10: Closed Job Search Exclusion

*For any* job with status "closed", it SHALL NOT appear in candidate-facing job search results.

**Validates: Requirements 5.4**

### Property 11: Job Soft Delete Behavior

*For any* job deletion action, the job record SHALL be soft-deleted (deleted_at timestamp set) rather than hard-deleted, and the record SHALL still exist in the database.

**Validates: Requirements 5.5**

### Property 12: Required Field Validation

*For any* job or resume submission missing required fields (job: title, description, location; resume: name, phone, email), the save operation SHALL be rejected with a validation error.

**Validates: Requirements 4.6, 5.6**

### Property 13: Status Update Notification

*For any* application or interview status update by HR, a notification message SHALL be created for the candidate with the status change information.

**Validates: Requirements 6.4, 7.4**

### Property 14: Interview Date Future Validation

*For any* interview scheduling attempt with a date in the past, the scheduling SHALL be rejected with a validation error.

**Validates: Requirements 7.5**

### Property 15: Chat Message Persistence Round-Trip

*For any* chat message sent in a conversation, the message SHALL be persisted to the database and retrievable with content matching what was sent.

**Validates: Requirements 8.5**

### Property 16: Offline Message Delivery

*For any* message sent to an offline user, the message SHALL be stored in the database and delivered via WebSocket when the user reconnects.

**Validates: Requirements 8.4**

### Property 17: Conversation Sorting Order

*For any* conversation list query, the conversations SHALL be sorted by last_message_at in descending order (most recent first).

**Validates: Requirements 9.1**

### Property 18: Unread Count Accuracy

*For any* conversation, the displayed unread_count SHALL equal the actual count of messages in that conversation where is_read is false and sender_id is not the current user.

**Validates: Requirements 9.3**

### Property 19: Mark As Read Behavior

*For any* conversation opened by a user, all messages in that conversation where sender_id is not the current user SHALL be marked as read (is_read = true).

**Validates: Requirements 9.4**

### Property 20: Unread Count Decrement on Mark Read

*For any* mark-as-read action on a conversation with N unread messages, the user's total unread count SHALL decrease by N.

**Validates: Requirements 10.6**

## Error Handling

### 后端错误处理

| 错误场景 | HTTP状态码 | 错误码 | 错误消息 |
|---------|-----------|--------|---------|
| 未授权访问 | 401 | 401 | "未授权，请先登录" |
| 权限不足 | 403 | 403 | "权限不足" |
| 资源不存在 | 404 | 404 | "资源不存在" |
| 参数验证失败 | 400 | 400 | "参数错误: {具体字段}" |
| 重复申请 | 400 | 1001 | "您已投递过该职位" |
| 简历不存在 | 400 | 1002 | "请先上传简历" |
| 文件格式错误 | 400 | 1003 | "不支持的文件格式" |
| 文件过大 | 400 | 1004 | "文件大小超过限制" |
| 面试日期无效 | 400 | 1005 | "面试日期必须是未来时间" |
| 服务器错误 | 500 | 500 | "服务器内部错误" |

### 前端错误处理

```typescript
// 统一错误处理
const handleApiError = (error: any) => {
    const code = error.response?.data?.code
    const message = error.response?.data?.message
    
    switch (code) {
        case 401:
            // 跳转登录
            router.push('/portal/login')
            break
        case 1001:
            ElMessage.warning('您已投递过该职位')
            break
        case 1002:
            ElMessage.warning('请先上传简历')
            break
        default:
            ElMessage.error(message || '操作失败，请稍后重试')
    }
}
```

### WebSocket 错误处理

```typescript
// 连接断开重连
ws.onclose = () => {
    if (reconnectAttempts < maxReconnectAttempts) {
        setTimeout(() => reconnect(), reconnectDelay)
    }
}

// 消息发送失败重试
const sendWithRetry = async (message: any, retries = 3) => {
    for (let i = 0; i < retries; i++) {
        try {
            await send(message)
            return
        } catch (e) {
            if (i === retries - 1) throw e
            await delay(1000 * (i + 1))
        }
    }
}
```

## Testing Strategy

### 测试方法

本项目采用双重测试策略：

1. **单元测试**：验证具体示例、边界情况和错误条件
2. **属性测试**：验证跨所有输入的通用属性

两种测试方法互补，共同提供全面的测试覆盖。

### 属性测试配置

- **测试框架**：后端使用 `testing/quick` (Go)，前端使用 `fast-check` (TypeScript)
- **最小迭代次数**：每个属性测试至少运行 100 次
- **标签格式**：`Feature: boss-like-features, Property {number}: {property_text}`

### 后端测试示例

```go
// Property 2: Application Creation Status
func TestApplicationCreationStatus(t *testing.T) {
    // Feature: boss-like-features, Property 2: Application Creation Status
    f := func(jobID, talentID, resumeID uint) bool {
        if jobID == 0 || talentID == 0 || resumeID == 0 {
            return true // Skip invalid inputs
        }
        
        app := &Application{
            JobID:    jobID,
            TalentID: talentID,
            ResumeID: resumeID,
        }
        
        err := db.Create(app).Error
        if err != nil {
            return true // Skip if creation fails for other reasons
        }
        
        return app.Status == "pending"
    }
    
    if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
        t.Error(err)
    }
}
```

### 前端测试示例

```typescript
// Property 5: Application Status Filter Correctness
import fc from 'fast-check'

describe('Application Status Filter', () => {
    // Feature: boss-like-features, Property 5: Application Status Filter Correctness
    it('should return only applications with matching status', () => {
        fc.assert(
            fc.property(
                fc.array(fc.record({
                    id: fc.nat(),
                    status: fc.constantFrom('pending', 'viewed', 'interview', 'rejected')
                })),
                fc.constantFrom('pending', 'viewed', 'interview', 'rejected'),
                (applications, filterStatus) => {
                    const filtered = filterApplicationsByStatus(applications, filterStatus)
                    return filtered.every(app => app.status === filterStatus)
                }
            ),
            { numRuns: 100 }
        )
    })
})
```

### 单元测试覆盖

| 模块 | 测试重点 |
|------|---------|
| Job Service | 职位CRUD、搜索筛选、状态变更 |
| Application Service | 申请创建、状态更新、重复检测 |
| Chat Service | 消息发送、会话管理、已读状态 |
| WebSocket | 连接管理、消息广播、断线重连 |
| 前端组件 | 表单验证、状态管理、API调用 |

### 集成测试

```go
// 端到端申请流程测试
func TestApplicationFlow(t *testing.T) {
    // 1. 创建职位
    job := createTestJob(t)
    
    // 2. 创建求职者和简历
    candidate := createTestCandidate(t)
    resume := createTestResume(t, candidate.ID)
    
    // 3. 提交申请
    app := submitApplication(t, job.ID, candidate.ID, resume.ID)
    assert.Equal(t, "pending", app.Status)
    
    // 4. HR查看申请
    updateApplicationStatus(t, app.ID, "viewed")
    
    // 5. 验证通知
    notifications := getNotifications(t, candidate.UserID)
    assert.True(t, len(notifications) > 0)
}
```
