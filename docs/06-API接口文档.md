# API接口文档

## 基础信息

- 基础URL: `http://localhost:8080/api/v1`
- 认证方式: JWT Bearer Token
- 响应格式: JSON

## 通用响应格式

```json
{
    "code": 0,
    "message": "success",
    "data": {}
}
```

错误码说明：
- `0`: 成功
- `1`: 通用错误
- `401`: 未授权
- `403`: 禁止访问
- `404`: 资源不存在
- `500`: 服务器错误

---

## 简历服务 (Resume Service)

### 上传简历文件

```
POST /resumes/upload
Content-Type: multipart/form-data
```

说明：
- 需要携带 JWT Bearer Token
- 求职者门户上传时，后端会根据当前登录用户自动绑定 `talent_id`

参数：
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | File | 是 | 简历文件(PDF/DOC/DOCX) |
| talent_id | int | 否 | 关联人才ID |
| job_id | int | 否 | 关联职位ID |

响应：
```json
{
    "code": 0,
    "message": "简历上传成功",
    "data": {
        "id": 1,
        "file_name": "张三_简历.pdf",
        "file_url": "/api/v1/resumes/file/xxx.pdf",
        "status": "pending"
    }
}
```

### 解析简历

```
POST /resumes/parse
Content-Type: application/json
```

请求体：
```json
{
    "text": "简历文本内容"
}
```

响应：
```json
{
    "code": 0,
    "message": "解析成功",
    "data": {
        "name": "张三",
        "phone": "13800138000",
        "email": "zhangsan@example.com",
        "education": "本科",
        "experience": "5年",
        "skills": ["Go", "Python", "Docker"],
        "risk_items": [],
        "risk_score": 0
    }
}
```

### AI智能解析

```
POST /ai/parse
Content-Type: multipart/form-data
```

参数：
```json
{
    "file": "简历文件",
    "jd_text": "职位描述（可选）"
}
```

响应：
```json
{
    "code": 0,
    "message": "AI解析成功",
    "data": {
        "parsed_resume": {...},
        "match_score": 85,
        "ai_summary": "该候选人..."
    }
}
```

---

## 推荐服务 (Recommendation Service)

### 为人才推荐职位

```
POST /recommendations/jobs-for-talent
Content-Type: application/json
```

请求体：
```json
{
    "id": 1,
    "name": "张三",
    "skills": ["Go", "Docker", "Kubernetes"],
    "experience": 5,
    "education": "本科",
    "location": "北京",
    "salary": "30-40K"
}
```

响应：
```json
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "id": 1,
            "name": "高级Go开发工程师",
            "score": 85.5,
            "reason": "高度匹配",
            "match_level": "high",
            "match_details": [
                "匹配技能: Go, Docker, Kubernetes",
                "经验完全匹配",
                "地理位置匹配"
            ]
        }
    ]
}
```

### 为职位推荐人才

```
POST /recommendations/talents-for-job
Content-Type: application/json
```

请求体：
```json
{
    "id": 1,
    "title": "高级Go开发工程师",
    "skills": ["Go", "Docker", "Kubernetes"],
    "location": "北京",
    "level": "senior",
    "salary": "30-50K"
}
```

### 生成归因报告

```
POST /recommendations/attribution-report
Content-Type: application/json
```

请求体：
```json
{
    "talent_id": 1,
    "job_id": 1
}
```

响应：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "summary": "该候选人与职位高度匹配",
        "match_score": 85,
        "dimensions": [
            {
                "name": "技能匹配",
                "score": 90,
                "weight": 50,
                "details": "掌握核心技能"
            }
        ],
        "strengths": ["核心技术栈匹配"],
        "gaps": ["缺少K8s生产经验"],
        "recommendation": "建议优先面试"
    }
}
```

---

## 评估服务 (Evaluator Service)

### 执行AI评估

```
POST /evaluate
Content-Type: application/json
```

请求体：
```json
{
    "resume_id": 1,
    "jd_text": "职位描述...",
    "position_id": 1
}
```

响应：
```json
{
    "code": 0,
    "message": "评估完成",
    "data": {
        "score": 85,
        "evaluation": "综合评估报告...",
        "details": {...}
    }
}
```

---

## 用户服务 (User Service)

### 用户登录

```
POST /auth/login
Content-Type: application/json
```

请求体：
```json
{
    "username": "admin",
    "password": "admin123"
}
```

响应：
```json
{
    "code": 0,
    "message": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
        "user": {
            "id": 1,
            "username": "admin",
            "role": "admin"
        }
    }
}
```

### 用户注册

```
POST /auth/register
Content-Type: application/json
```

请求体：
```json
{
    "username": "newuser",
    "password": "password123",
    "email": "user@example.com",
    "role": "hr"
}
```

---

## 职位服务 (Job Service)

### 获取职位列表

```
GET /jobs?page=1&page_size=10&status=open
```

### 创建职位

```
POST /jobs
Content-Type: application/json
```

请求体：
```json
{
    "title": "高级Go开发工程师",
    "department": "技术部",
    "location": "北京",
    "salary": "30-50K",
    "level": "senior",
    "skills": ["Go", "Docker", "Kubernetes"],
    "description": "职位描述..."
}
```

---

## 人才服务 (Talent Service)

### 获取人才列表

```
GET /talents?page=1&page_size=10&status=active
```

### 创建人才档案

```
POST /talents
Content-Type: application/json
```

请求体：
```json
{
    "name": "张三",
    "phone": "13800138000",
    "email": "zhangsan@example.com",
    "skills": ["Go", "Python"],
    "experience": 5,
    "education": "本科",
    "location": "北京"
}
```

---

## 错误处理

### 常见错误响应

```json
{
    "code": 1,
    "message": "参数错误: file字段不能为空",
    "data": null
}
```

```json
{
    "code": 401,
    "message": "未授权，请先登录",
    "data": null
}
```

```json
{
    "code": 500,
    "message": "服务器内部错误",
    "data": null
}
```
