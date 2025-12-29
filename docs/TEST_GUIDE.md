# 智能人才运营平台 - 功能测试指南

## 一、快速API测试

### 运行自动化测试脚本
```bash
cd backend
chmod +x test_api.sh
./test_api.sh
```

**预期结果**：90个测试用例全部通过，通过率100%

测试覆盖：
- 用户服务（10个）：登录、注册、用户列表
- 职位服务（14个）：职位CRUD、筛选、搜索
- 人才服务（14个）：人才CRUD、搜索
- 简历服务（14个）：简历管理、状态更新
- 面试服务（16个）：面试安排、反馈
- 消息服务（10个）：消息发送、已读标记
- 日志服务（8个）：ES日志查询、统计
- 综合测试（4个）：ES日志验证

## 二、手动功能测试清单

### 1. 用户认证模块

| 功能 | 测试步骤 | 预期结果 |
|------|----------|----------|
| 登录 | 访问 http://localhost:3000/login，输入 admin/password123 | 登录成功，跳转到仪表板 |
| 登录失败 | 输入错误密码 | 显示错误提示 |
| 退出登录 | 点击右上角头像，选择退出 | 返回登录页 |

### 2. 职位管理模块

| 功能 | 测试步骤 | 预期结果 |
|------|----------|----------|
| 职位列表 | 点击左侧菜单"职位管理" | 显示职位列表，有12条数据 |
| 职位筛选 | 选择状态/地点筛选 | 列表按条件过滤 |
| 职位搜索 | 输入关键词搜索 | 显示匹配结果 |
| 职位详情 | 点击某个职位 | 显示职位详细信息 |
| 新建职位 | 点击"新建职位"按钮 | 弹出表单，填写后保存成功 |

### 3. 人才管理模块

| 功能 | 测试步骤 | 预期结果 |
|------|----------|----------|
| 人才列表 | 点击左侧菜单"人才库" | 显示人才列表，有20条数据 |
| 人才筛选 | 按技能/状态筛选 | 列表按条件过滤 |
| 人才详情 | 点击某个人才 | 显示人才详细信息 |

### 4. 简历管理模块

| 功能 | 测试步骤 | 预期结果 |
|------|----------|----------|
| 简历列表 | 点击左侧菜单"简历管理" | 显示简历列表，有20条数据 |
| 上传简历 | 点击"上传简历"，选择PDF文件 | 上传成功，列表刷新 |
| 简历排序 | 点击"按时间排序"/"按状态排序" | 列表重新排序 |
| 下载简历 | 点击下载按钮 | 下载简历文件 |
| 删除简历 | 点击删除按钮 | 确认后删除成功 |

### 5. 面试管理模块

| 功能 | 测试步骤 | 预期结果 |
|------|----------|----------|
| 面试列表 | 点击左侧菜单"面试管理" | 显示面试列表，有20条数据 |
| 面试详情 | 点击某场面试 | 显示面试详细信息 |
| 面试反馈 | 在面试详情页提交反馈 | 反馈保存成功 |

### 6. 消息中心

| 功能 | 测试步骤 | 预期结果 |
|------|----------|----------|
| 消息列表 | 点击左侧菜单"消息中心" | 显示消息列表 |
| 未读消息 | 查看顶部消息图标 | 显示未读消息数量 |
| 标记已读 | 点击消息 | 消息标记为已读 |

### 7. 前台求职门户

| 功能 | 测试步骤 | 预期结果 |
|------|----------|----------|
| 职位浏览 | 访问 http://localhost:3000/portal/jobs | 显示公开职位列表 |
| 职位筛选 | 按经验/学历筛选 | 列表按条件过滤 |
| 职位详情 | 点击某个职位 | 显示职位详情和投递按钮 |
| 投递简历 | 点击"投递简历" | 跳转到登录或上传简历 |

### 8. AI简历评估

| 功能 | 测试步骤 | 预期结果 |
|------|----------|----------|
| 访问评估系统 | 点击顶部"AI评估"按钮或访问 http://localhost:8090 | 显示AI评估登录页 |
| 登录 | 选择"graduate"模式登录 | 登录成功 |
| 批量评估 | 选择简历进行AI评估 | 显示评估结果 |

## 三、API接口测试 (curl命令)

### 用户登录
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

### 获取职位列表
```bash
curl http://localhost:8082/api/v1/jobs
```

### 获取职位详情
```bash
curl http://localhost:8082/api/v1/jobs/1
```

### 获取人才列表
```bash
curl http://localhost:8086/api/v1/talents
```

### 获取简历列表
```bash
curl http://localhost:8084/api/v1/resumes
```

### 获取面试列表
```bash
curl http://localhost:8083/api/v1/interviews
```

### 获取消息列表
```bash
curl "http://localhost:8085/api/v1/messages?user_id=1"
```

## 四、数据库验证

```bash
# 连接数据库
psql -U qinyang -d talent_platform

# 查看各表数据量
SELECT 'users' as table_name, count(*) FROM users
UNION ALL SELECT 'jobs', count(*) FROM jobs
UNION ALL SELECT 'talents', count(*) FROM talents
UNION ALL SELECT 'resumes', count(*) FROM resumes
UNION ALL SELECT 'interviews', count(*) FROM interviews
UNION ALL SELECT 'messages', count(*) FROM messages;

# 查看用户数据
SELECT id, username, email, role FROM users;

# 查看职位数据
SELECT id, title, status, location FROM jobs;

# 查看人才数据
SELECT id, name, email, status FROM talents;
```

## 五、测试账号

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| admin | password123 | 超级管理员 | 拥有所有权限 |
| hr_zhang | password123 | HR主管 | 招聘流程管理 |
| hr_li | password123 | 招聘专员 | 日常招聘工作 |
| tech_chen | password123 | 面试官 | 技术面试 |
| viewer_test | password123 | 只读用户 | 只能查看 |

## 六、服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| frontend | 5173 | 前端页面 (Vite) |
| user-service | 8081 | 用户认证服务 |
| job-service | 8082 | 职位管理服务 |
| interview-service | 8083 | 面试管理服务 |
| resume-service | 8084 | 简历管理服务 |
| message-service | 8085 | 消息服务 |
| talent-service | 8086 | 人才管理服务 |
| log-service | 8088 | 日志查询服务 |
| evaluator-service | 8090 | AI评估服务 |
| PostgreSQL | 5432 | 数据库 |
| Elasticsearch | 9200 | 日志存储 |
| Kibana | 5601 | 日志可视化 |

## 七、前端单元测试

```bash
cd frontend
npm run test
```

## 八、常见问题排查

### Q: API测试失败
1. 检查所有后端服务是否启动
2. 检查数据库连接是否正常
3. 检查模拟数据是否导入

### Q: 前端页面报错
1. 检查浏览器控制台错误信息
2. 确认后端服务正常运行
3. 清除浏览器缓存后重试

### Q: ES日志功能不可用
1. 确认Elasticsearch已启动：`curl http://localhost:9200`
2. 日志功能会自动降级，不影响其他功能
