# 智能人才招聘管理平台 - 测试指南

> 📖 返回 [项目首页](../README.md) | 相关文档：[系统架构](ARCHITECTURE.md) | [快速启动](QUICKSTART.md) | [代码规范](CODE_GUIDE.md)

---

## 1. 测试概述

本项目包含后端 API 测试和前端单元测试两部分。

---

## 2. 后端 API 测试

### 2.1 运行测试

```bash
cd backend
chmod +x test_api.sh
./test_api.sh
```

### 2.2 测试覆盖范围

| 服务 | 测试内容 |
|------|----------|
| user-service | 登录、注册、用户列表 |
| job-service | 职位CRUD、筛选、搜索、统计 |
| talent-service | 人才CRUD、搜索 |
| resume-service | 简历管理、状态更新、申请管理 |
| interview-service | 面试安排、反馈、统计 |
| message-service | 消息发送、已读标记、未读统计 |
| log-service | ES日志查询、统计 |
| recommendation-service | 智能推荐 |

### 2.3 手动测试示例

**用户登录**：
```bash
curl -X POST http://localhost:8081/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

**获取职位列表**：
```bash
curl http://localhost:8082/api/v1/jobs
```

**创建面试**：
```bash
curl -X POST http://localhost:8083/api/v1/interviews \
  -H "Content-Type: application/json" \
  -d '{
    "candidate_id": 1,
    "candidate_name": "张三",
    "position_id": 1,
    "position": "高级Go开发工程师",
    "type": "initial",
    "date": "2024-01-15",
    "time": "14:00",
    "interviewer": "陈总监",
    "method": "onsite"
  }'
```

**智能推荐**：
```bash
curl -X POST http://localhost:8087/api/v1/recommendations/talents-for-job \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "title": "高级Go开发工程师",
    "skills": ["Go", "Docker", "Kubernetes"],
    "location": "北京",
    "level": "senior"
  }'
```

---

## 3. 前端单元测试

### 3.1 运行测试

```bash
cd frontend
npm run test
```

### 3.2 测试覆盖范围

| 模块 | 测试内容 |
|------|----------|
| API | 接口调用、错误处理 |
| Store | Pinia 状态管理 |
| Utils | 工具函数 |
| Components | 组件渲染、交互 |

### 3.3 测试框架

- **Vitest**：测试运行器
- **Vue Test Utils**：Vue 组件测试工具

---

## 4. 功能测试清单

### 4.1 用户认证

- [ ] 用户注册
- [ ] 用户登录
- [ ] 登出
- [ ] Token 过期处理

### 4.2 人才管理

- [ ] 人才列表展示
- [ ] 人才搜索筛选
- [ ] 人才详情查看
- [ ] 人才创建/编辑/删除

### 4.3 职位管理

- [ ] 职位列表展示
- [ ] 职位搜索筛选
- [ ] 职位详情查看
- [ ] 职位创建/编辑/删除
- [ ] 职位状态切换

### 4.4 简历管理

- [ ] 简历列表展示
- [ ] 简历上传
- [ ] 简历状态更新
- [ ] AI评估功能

### 4.5 智能推荐

- [ ] 为职位推荐人才
- [ ] 为人才推荐职位
- [ ] 匹配度展示
- [ ] 推荐详情查看

### 4.6 面试管理

- [ ] 面试日历展示
- [ ] 面试安排
- [ ] 面试反馈提交
- [ ] 面试状态更新

### 4.7 消息中心

- [ ] 消息列表展示
- [ ] 消息分类筛选
- [ ] 标记已读
- [ ] 消息删除

### 4.8 数据报表

- [ ] 仪表板数据展示
- [ ] 招聘漏斗图
- [ ] 趋势图表
- [ ] 数据大屏

---

## 📚 相关文档

| 文档 | 说明 |
|------|------|
| [📖 项目首页](../README.md) | 项目概述 |
| [📐 系统架构](ARCHITECTURE.md) | 架构设计 |
| [🚀 快速启动](QUICKSTART.md) | 环境配置 |
| [📝 代码规范](CODE_GUIDE.md) | 开发指南 |
