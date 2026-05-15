# bigtest 真实流量自动化测试报告

- Run ID: `2026-04-30_16-17-30`
- 配置文件: `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/config.local.json`
- 输出目录: `output/2026-04-30_16-17-30`
- 自动探索: 已开启
- 写操作场景: 已开启，测试数据前缀 `BIGTEST_2026-04-30_16-17-30`

## 测试结论

本次测试共访问 28 个页面，捕获 206 条真实 API 请求，生成 145 条接口回放用例。
接口回放通过 112 条，失败 9 条，跳过 24 条；写操作场景通过 4 个，失败 1 个。
用例来源：真实流量 27 条，OpenAPI 自动生成 78 条，负向异常用例 40 条。
失败分层：真实失败 1 条，疑似依赖/环境失败 3 条，负向用例失败 4 条，自动用例数据不足跳过 24 条。

**结论：本次自动化测试存在失败项，请优先查看“回放结果”和“写操作场景步骤”。**

## 测试环境与策略

| 项目 | 值 |
| --- | --- |
| 前端地址 | `http://localhost:5173` |
| 后端回放地址 | `http://localhost:8080` |
| 捕获 URL 规则 | `/api/v1/` |
| 回放方法 | `GET, POST, PUT, PATCH, DELETE` |
| 排除路径 | `/api/v1/ai, /api/v1/evaluations, /upload` |
| 页面探索最大路由数 | 120 |
| 每页最大点击数 | 15 |
| 截图策略 | 最多 5 张；按钮点击失败默认不截图，只记录 URL 和错误 |
| AI 评估策略 | 明确排除 `/api/v1/ai`、`/api/v1/evaluations`、`/ai/evaluate`、`/evaluate` 等路径，不触发真实 AI 评估、上传或外部模型调用 |

## 总览

| Profile | 页面 | 输入 | 按钮 | 原始请求 | 生成用例 | 写操作通过 | 写操作失败 | 回放通过 | 回放失败 | 回放跳过 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| auto-default | 28 | 437/133 | 541/160 | 206 | 145 | 4 | 1 | 112 | 9 | 24 |

## 模块覆盖统计

| 模块 | 页面访问 | API 用例 | GET | POST | PUT | PATCH | DELETE | 未覆盖接口 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 公告 | 3 | 17 | 5 | 4 | 4 | 0 | 4 | 0 |
| 简历 | 2 | 16 | 6 | 6 | 3 | 0 | 1 | 0 |
| 面试 | 1 | 20 | 8 | 8 | 2 | 0 | 2 | 0 |
| 其他 | 12 | 2 | 0 | 2 | 0 | 0 | 0 | 0 |
| 人才 | 2 | 22 | 10 | 4 | 4 | 0 | 4 | 0 |
| 认证 | 0 | 1 | 0 | 1 | 0 | 0 | 0 | 0 |
| 日志 | 1 | 5 | 4 | 0 | 0 | 0 | 1 | 0 |
| 统计 | 0 | 7 | 7 | 0 | 0 | 0 | 0 | 0 |
| 投递 | 1 | 13 | 4 | 3 | 3 | 0 | 3 | 0 |
| 推荐 | 0 | 26 | 1 | 25 | 0 | 0 | 0 | 0 |
| 消息 | 1 | 18 | 10 | 4 | 2 | 0 | 2 | 0 |
| 用户 | 1 | 3 | 2 | 0 | 1 | 0 | 0 | 0 |
| 职位 | 4 | 16 | 10 | 2 | 2 | 0 | 2 | 0 |

## 未覆盖清单

### 未覆盖页面

| 页面 | 模块 | 来源 |
| --- | --- | --- |
| `/login` | 日志 | router |
| `/register` | 认证 | router |
| `/dev-logs` | 日志 | router |
| `/portal/login` | 日志 | router |
| `/portal/register` | 认证 | router |
| `/ai-evaluate` | 其他 | router |
| `/ai-process` | 其他 | router |
| `/evaluation-results` | 其他 | router |

### 未覆盖接口

- 无。

### 未覆盖写操作

- 无。

### 主动排除功能（AI/上传）

| 排除规则 | 原因 |
| --- | --- |
| `/api/v1/ai` | 不触发 AI 评估、上传或外部模型调用 |
| `/api/v1/evaluations` | 不触发 AI 评估、上传或外部模型调用 |
| `/upload` | 不触发 AI 评估、上传或外部模型调用 |
| `/api/v1/ai/` | 不触发 AI 评估、上传或外部模型调用 |
| `/ai/evaluate` | 不触发 AI 评估、上传或外部模型调用 |
| `/evaluate` | 不触发 AI 评估、上传或外部模型调用 |
| `/evaluation` | 不触发 AI 评估、上传或外部模型调用 |

## auto-default

### 自动发现覆盖

- 发现路由: 36
- 实际访问页面: 28
- 输入框: 437 个，已自动填写 133 个
- 按钮/链接: 541 个，已点击 160 个
- 自动取消确认弹窗: 0
- API 资源数: 15
- 跳过/排除: 43

### API 资源覆盖

| 资源 | 方法 | 用例数 |
| --- | --- | ---: |
| stats | GET | 7 |
| messages | DELETE, GET, POST, PUT | 10 |
| talents | DELETE, GET, POST, PUT | 22 |
| jobs | DELETE, GET, POST, PUT | 16 |
| conversations | GET, POST, PUT | 6 |
| applications | DELETE, GET, POST, PUT | 13 |
| resumes | DELETE, GET, POST, PUT | 16 |
| interviews | DELETE, GET, POST, PUT | 20 |
| logs | DELETE, GET | 5 |
| notices | DELETE, GET, POST, PUT | 17 |
| recommendations | GET, POST | 26 |
| online-status | GET | 2 |
| profile | GET, PUT | 2 |
| users | GET | 1 |
| register | POST | 1 |

### 被跳过的页面操作

| 风险 | 文案 | 页面 |
| --- | --- | --- |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
| disabled | el-button el-button--primary el-button--large is-disabled 开始AI智能评估 请上传简历并填写职位描述后开始评估 开始AI智能评估 el-button el-button--primary el-button--large is-disabled | `http://localhost:5173/ai-evaluate` |
### 写操作场景结果

本轮写操作场景已执行 5 次 `DELETE` 删除请求。

| 场景 | 来源 | 状态 | 资源 | 步骤数 | 说明 |
| --- | --- | --- | --- | ---: | --- |
| job-crud-owned-data | 配置 | PASS | jobs | 4 | BIGTEST 数据创建、更新、删除完成 |
| talent-crud-owned-data | 配置 | FAIL | talents | 4 | BIGTEST 数据创建、更新、删除完成 |
| message-create-query-delete-owned-data | 配置 | PASS | messages | 3 | BIGTEST 数据创建、更新、删除完成 |
| interview-crud-owned-data | 配置 | PASS | interviews | 6 | BIGTEST 数据创建、更新、删除完成 |
| notice-crud-owned-data | 配置 | PASS | notice | 4 | BIGTEST 数据创建、更新、删除完成 |

### 写操作场景步骤

| 场景 | 步骤 | 方法 | 路径 | 期望 | 实际 | 状态 |
| --- | --- | --- | --- | --- | ---: | --- |
| job-crud-owned-data | create 创建测试数据 | POST | `/api/v1/jobs` | 200, 201 | 201 | PASS |
| job-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/jobs/67` | 200 | 200 | PASS |
| job-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/jobs/67` | 200 | 200 | PASS |
| job-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/jobs/67` | 200, 204 | 200 | PASS |
| talent-crud-owned-data | create 创建测试数据 | POST | `/api/v1/talents` | 200, 201 | 201 | PASS |
| talent-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/talents/122` | 200 | 200 | PASS |
| talent-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/talents/122` | 200 | 500 | FAIL |
| talent-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/talents/122` | 200, 204 | 200 | PASS |
| message-create-query-delete-owned-data | create 创建测试数据 | POST | `/api/v1/messages` | 200, 201 | 201 | PASS |
| message-create-query-delete-owned-data | verify 验证测试数据 | GET | `/api/v1/messages` | 200 | 200 | PASS |
| message-create-query-delete-owned-data | delete 删除测试数据 | DELETE | `/api/v1/messages/55` | 200, 204 | 200 | PASS |
| interview-crud-owned-data | create 创建测试数据 | POST | `/api/v1/interviews` | 200, 201 | 201 | PASS |
| interview-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/interviews/117` | 200 | 200 | PASS |
| interview-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/interviews/117` | 200 | 200 | PASS |
| interview-crud-owned-data | reschedule 改期测试面试 | POST | `/api/v1/interviews/117/reschedule` | 200 | 200 | PASS |
| interview-crud-owned-data | complete 完成测试面试 | POST | `/api/v1/interviews/117/complete` | 200 | 200 | PASS |
| interview-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/interviews/117` | 200, 204 | 200 | PASS |
| notice-crud-owned-data | create 创建测试数据 | POST | `/api/v1/notices` | 200, 201 | 200 | PASS |
| notice-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/notices/36` | 200 | 200 | PASS |
| notice-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/notices/36` | 200 | 200 | PASS |
| notice-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/notices/36` | 200, 204 | 200 | PASS |

### 去重后生成的用例

| 用例ID | 来源 | 方法 | 路径 | 归并次数 | 需鉴权 |
| --- | --- | --- | --- | ---: | --- |
| auto-default-1 | 真实流量 | POST | `/api/v1/login` | 1 | 否 |
| auto-default-2 | 真实流量 | GET | `/api/v1/stats/dashboard` | 3 | 是 |
| auto-default-3 | 真实流量 | GET | `/api/v1/messages/unread-count` | 22 | 是 |
| auto-default-4 | 真实流量 | GET | `/api/v1/messages?page=1&page_size=5` | 3 | 是 |
| auto-default-5 | 真实流量 | GET | `/api/v1/talents?page=1&page_size=20` | 3 | 是 |
| auto-default-6 | 真实流量 | GET | `/api/v1/jobs?page=1&page_size=5` | 3 | 否 |
| auto-default-7 | 真实流量 | GET | `/api/v1/conversations/unread-count` | 57 | 是 |
| auto-default-8 | 真实流量 | GET | `/api/v1/jobs?page=1&page_size=10&status=open` | 28 | 是 |
| auto-default-9 | 真实流量 | GET | `/api/v1/applications?page=1&page_size=1000&talent_id=me` | 29 | 是 |
| auto-default-10 | 真实流量 | GET | `/api/v1/conversations` | 17 | 是 |
| auto-default-11 | 真实流量 | GET | `/api/v1/jobs/1` | 2 | 是 |
| auto-default-12 | 真实流量 | GET | `/api/v1/applications?page=1&page_size=20` | 1 | 是 |
| auto-default-13 | 真实流量 | GET | `/api/v1/resumes/online` | 1 | 是 |
| auto-default-14 | 真实流量 | GET | `/api/v1/resumes?page_size=10` | 1 | 是 |
| auto-default-15 | 真实流量 | GET | `/api/v1/jobs?page_size=100&status=open` | 22 | 是 |
| auto-default-16 | 真实流量 | GET | `/api/v1/talents?page=1&page_size=10` | 1 | 是 |
| auto-default-17 | 真实流量 | GET | `/api/v1/jobs?page=1&page_size=10` | 1 | 是 |
| auto-default-18 | 真实流量 | GET | `/api/v1/jobs/1/applications?page=1&page_size=10` | 1 | 是 |
| auto-default-19 | 真实流量 | GET | `/api/v1/resumes?page=1&page_size=10` | 1 | 是 |
| auto-default-20 | 真实流量 | GET | `/api/v1/jobs?page=1&page_size=20` | 1 | 是 |
| auto-default-21 | 真实流量 | GET | `/api/v1/applications?page=1&page_size=50` | 1 | 是 |
| auto-default-22 | 真实流量 | GET | `/api/v1/interviews/1` | 1 | 是 |
| auto-default-23 | 真实流量 | GET | `/api/v1/interviews/1/feedback` | 1 | 是 |
| auto-default-24 | 真实流量 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | 1 | 是 |
| auto-default-25 | 真实流量 | GET | `/api/v1/logs?page=1&page_size=20` | 1 | 是 |
| auto-default-26 | 真实流量 | GET | `/api/v1/notices/1` | 2 | 是 |
| auto-default-27 | 真实流量 | GET | `/api/v1/notices?page=1&page_size=10` | 1 | 是 |
| auto-default-openapi-28 | OpenAPI | POST | `/api/v1/talents` | 1 | 是 |
| auto-default-openapi-29 | OpenAPI | GET | `/api/v1/talents/122` | 1 | 是 |
| auto-default-openapi-30 | OpenAPI | PUT | `/api/v1/talents/122` | 1 | 是 |
| auto-default-openapi-31 | OpenAPI | DELETE | `/api/v1/talents/122` | 1 | 是 |
| auto-default-openapi-32 | OpenAPI | GET | `/api/v1/talents/search?keyword=%E6%B5%8B%E8%AF%95` | 1 | 是 |
| auto-default-openapi-33 | OpenAPI | GET | `/api/v1/talents/stats` | 1 | 是 |
| auto-default-openapi-34 | OpenAPI | POST | `/api/v1/notices` | 1 | 是 |
| auto-default-openapi-35 | OpenAPI | PUT | `/api/v1/notices/1` | 1 | 是 |
| auto-default-openapi-36 | OpenAPI | DELETE | `/api/v1/notices/1` | 1 | 是 |
| auto-default-openapi-37 | OpenAPI | POST | `/api/v1/applications` | 1 | 是 |
| auto-default-openapi-38 | OpenAPI | PUT | `/api/v1/applications/1` | 1 | 是 |
| auto-default-openapi-39 | OpenAPI | DELETE | `/api/v1/applications/1` | 1 | 是 |
| auto-default-openapi-40 | OpenAPI | POST | `/api/v1/recommendations/attribution-report` | 1 | 是 |
| auto-default-openapi-41 | OpenAPI | POST | `/api/v1/recommendations/batch` | 1 | 是 |
| auto-default-openapi-42 | OpenAPI | POST | `/api/v1/recommendations/jobs-for-talent` | 1 | 是 |
| auto-default-openapi-43 | OpenAPI | POST | `/api/v1/recommendations/rag/index-all` | 1 | 是 |
| auto-default-openapi-44 | OpenAPI | POST | `/api/v1/recommendations/rag/index-job` | 1 | 是 |
| auto-default-openapi-45 | OpenAPI | POST | `/api/v1/recommendations/rag/index-resume` | 1 | 是 |
| auto-default-openapi-46 | OpenAPI | POST | `/api/v1/recommendations/rag/index-talent` | 1 | 是 |
| auto-default-openapi-47 | OpenAPI | POST | `/api/v1/recommendations/rag/match` | 1 | 是 |
| auto-default-openapi-48 | OpenAPI | POST | `/api/v1/recommendations/rag/query` | 1 | 是 |
| auto-default-openapi-49 | OpenAPI | POST | `/api/v1/recommendations/semantic-match` | 1 | 是 |
| auto-default-openapi-50 | OpenAPI | GET | `/api/v1/recommendations/stats` | 1 | 是 |
| auto-default-openapi-51 | OpenAPI | POST | `/api/v1/recommendations/talents-for-job` | 1 | 是 |
| auto-default-openapi-52 | OpenAPI | GET | `/api/v1/logs/actions` | 1 | 是 |
| auto-default-openapi-53 | OpenAPI | DELETE | `/api/v1/logs/cleanup` | 1 | 是 |
| auto-default-openapi-54 | OpenAPI | GET | `/api/v1/logs/services` | 1 | 是 |
| auto-default-openapi-55 | OpenAPI | GET | `/api/v1/logs/stats` | 1 | 是 |
| auto-default-openapi-56 | OpenAPI | POST | `/api/v1/conversations` | 1 | 是 |
| auto-default-openapi-57 | OpenAPI | GET | `/api/v1/conversations/1/messages?page=1&page_size=10` | 1 | 是 |
| auto-default-openapi-58 | OpenAPI | POST | `/api/v1/conversations/1/messages` | 1 | 是 |
| auto-default-openapi-59 | OpenAPI | PUT | `/api/v1/conversations/1/read` | 1 | 是 |
| auto-default-openapi-60 | OpenAPI | POST | `/api/v1/messages` | 1 | 是 |
| auto-default-openapi-61 | OpenAPI | DELETE | `/api/v1/messages/55` | 1 | 是 |
| auto-default-openapi-62 | OpenAPI | PUT | `/api/v1/messages/55/read` | 1 | 是 |
| auto-default-openapi-63 | OpenAPI | GET | `/api/v1/messages/stats` | 1 | 是 |
| auto-default-openapi-64 | OpenAPI | GET | `/api/v1/online-status` | 1 | 是 |
| auto-default-openapi-65 | OpenAPI | GET | `/api/v1/online-status/1` | 1 | 是 |
| auto-default-openapi-66 | OpenAPI | GET | `/api/v1/profile` | 1 | 是 |
| auto-default-openapi-67 | OpenAPI | PUT | `/api/v1/profile` | 1 | 是 |
| auto-default-openapi-68 | OpenAPI | GET | `/api/v1/users?page=1&page_size=10` | 1 | 是 |
| auto-default-openapi-69 | OpenAPI | POST | `/api/v1/resumes` | 1 | 否 |
| auto-default-openapi-70 | OpenAPI | DELETE | `/api/v1/resumes/1` | 1 | 否 |
| auto-default-openapi-71 | OpenAPI | GET | `/api/v1/resumes/1/download` | 1 | 否 |
| auto-default-openapi-72 | OpenAPI | PUT | `/api/v1/resumes/1/job` | 1 | 是 |
| auto-default-openapi-73 | OpenAPI | PUT | `/api/v1/resumes/1/status` | 1 | 是 |
| auto-default-openapi-74 | OpenAPI | GET | `/api/v1/resumes/evaluation?page=1&page_size=10&status=open` | 1 | 否 |
| auto-default-openapi-75 | OpenAPI | GET | `/api/v1/resumes/file/test.pdf` | 1 | 否 |
| auto-default-openapi-76 | OpenAPI | POST | `/api/v1/resumes/match` | 1 | 否 |
| auto-default-openapi-77 | OpenAPI | PUT | `/api/v1/resumes/online` | 1 | 是 |
| auto-default-openapi-78 | OpenAPI | POST | `/api/v1/resumes/parse` | 1 | 否 |
| auto-default-openapi-79 | OpenAPI | POST | `/api/v1/resumes/risk-check` | 1 | 否 |
| auto-default-openapi-80 | OpenAPI | POST | `/api/v1/resumes/risk-check/education` | 1 | 否 |
| auto-default-openapi-81 | OpenAPI | POST | `/api/v1/resumes/risk-check/time-conflict` | 1 | 否 |
| auto-default-openapi-82 | OpenAPI | GET | `/api/v1/stats/channels` | 1 | 否 |
| auto-default-openapi-83 | OpenAPI | GET | `/api/v1/stats/department-progress` | 1 | 否 |
| auto-default-openapi-84 | OpenAPI | GET | `/api/v1/stats/funnel` | 1 | 否 |
| auto-default-openapi-85 | OpenAPI | GET | `/api/v1/stats/interviewer-rank` | 1 | 否 |
| auto-default-openapi-86 | OpenAPI | GET | `/api/v1/stats/job-rank` | 1 | 否 |
| auto-default-openapi-87 | OpenAPI | GET | `/api/v1/stats/trend` | 1 | 否 |
| auto-default-openapi-88 | OpenAPI | POST | `/api/v1/jobs` | 1 | 是 |
| auto-default-openapi-89 | OpenAPI | PUT | `/api/v1/jobs/1` | 1 | 是 |
| auto-default-openapi-90 | OpenAPI | DELETE | `/api/v1/jobs/1` | 1 | 是 |
| auto-default-openapi-91 | OpenAPI | GET | `/api/v1/jobs/hot` | 1 | 否 |
| auto-default-openapi-92 | OpenAPI | GET | `/api/v1/jobs/stats` | 1 | 否 |
| auto-default-openapi-93 | OpenAPI | POST | `/api/v1/register` | 1 | 否 |
| auto-default-openapi-94 | OpenAPI | GET | `/api/v1/interviews?page=1&page_size=10&status=open` | 1 | 是 |
| auto-default-openapi-95 | OpenAPI | POST | `/api/v1/interviews` | 1 | 是 |
| auto-default-openapi-96 | OpenAPI | PUT | `/api/v1/interviews/1` | 1 | 是 |
| auto-default-openapi-97 | OpenAPI | DELETE | `/api/v1/interviews/1` | 1 | 是 |
| auto-default-openapi-98 | OpenAPI | POST | `/api/v1/interviews/1/cancel` | 1 | 是 |
| auto-default-openapi-99 | OpenAPI | POST | `/api/v1/interviews/1/complete` | 1 | 是 |
| auto-default-openapi-100 | OpenAPI | POST | `/api/v1/interviews/1/feedback` | 1 | 是 |
| auto-default-openapi-101 | OpenAPI | POST | `/api/v1/interviews/1/reschedule` | 1 | 是 |
| auto-default-openapi-102 | OpenAPI | GET | `/api/v1/interviews/candidate/122` | 1 | 是 |
| auto-default-openapi-103 | OpenAPI | GET | `/api/v1/interviews/interviewer/1` | 1 | 是 |
| auto-default-openapi-104 | OpenAPI | GET | `/api/v1/interviews/stats` | 1 | 是 |
| auto-default-openapi-105 | OpenAPI | GET | `/api/v1/interviews/today` | 1 | 是 |
| auto-default-negative-106 | 负向 | GET | `/api/v1/talents?page=1&page_size=10&status=open&search=%E6%B5%8B%E8%AF%95` | 1 | 是 |
| auto-default-negative-107 | 负向 | POST | `/api/v1/talents` | 1 | 是 |
| auto-default-negative-108 | 负向 | POST | `/api/v1/talents` | 1 | 是 |
| auto-default-negative-109 | 负向 | GET | `/api/v1/talents/122` | 1 | 是 |
| auto-default-negative-110 | 负向 | PUT | `/api/v1/talents/122` | 1 | 是 |
| auto-default-negative-111 | 负向 | PUT | `/api/v1/talents/99999999` | 1 | 是 |
| auto-default-negative-112 | 负向 | DELETE | `/api/v1/talents/122` | 1 | 是 |
| auto-default-negative-113 | 负向 | DELETE | `/api/v1/talents/99999999` | 1 | 是 |
| auto-default-negative-114 | 负向 | GET | `/api/v1/talents/search?keyword=%E6%B5%8B%E8%AF%95` | 1 | 是 |
| auto-default-negative-115 | 负向 | GET | `/api/v1/talents/stats` | 1 | 是 |
| auto-default-negative-116 | 负向 | GET | `/api/v1/notices?page=1&page_size=10&keyword=%E6%B5%8B%E8%AF%95&status=open` | 1 | 是 |
| auto-default-negative-117 | 负向 | POST | `/api/v1/notices` | 1 | 是 |
| auto-default-negative-118 | 负向 | POST | `/api/v1/notices` | 1 | 是 |
| auto-default-negative-119 | 负向 | GET | `/api/v1/notices/1` | 1 | 是 |
| auto-default-negative-120 | 负向 | PUT | `/api/v1/notices/1` | 1 | 是 |
| auto-default-negative-121 | 负向 | PUT | `/api/v1/notices/99999999` | 1 | 是 |
| auto-default-negative-122 | 负向 | DELETE | `/api/v1/notices/1` | 1 | 是 |
| auto-default-negative-123 | 负向 | DELETE | `/api/v1/notices/99999999` | 1 | 是 |
| auto-default-negative-124 | 负向 | POST | `/api/v1/login` | 1 | 否 |
| auto-default-negative-125 | 负向 | GET | `/api/v1/applications?page=1&page_size=10&status=open` | 1 | 是 |
| auto-default-negative-126 | 负向 | POST | `/api/v1/applications` | 1 | 是 |
| auto-default-negative-127 | 负向 | POST | `/api/v1/applications` | 1 | 是 |
| auto-default-negative-128 | 负向 | PUT | `/api/v1/applications/1` | 1 | 是 |
| auto-default-negative-129 | 负向 | PUT | `/api/v1/applications/99999999` | 1 | 是 |
| auto-default-negative-130 | 负向 | DELETE | `/api/v1/applications/1` | 1 | 是 |
| auto-default-negative-131 | 负向 | DELETE | `/api/v1/applications/99999999` | 1 | 是 |
| auto-default-negative-132 | 负向 | POST | `/api/v1/recommendations/attribution-report` | 1 | 是 |
| auto-default-negative-133 | 负向 | POST | `/api/v1/recommendations/attribution-report` | 1 | 是 |
| auto-default-negative-134 | 负向 | POST | `/api/v1/recommendations/batch` | 1 | 是 |
| auto-default-negative-135 | 负向 | POST | `/api/v1/recommendations/batch` | 1 | 是 |
| auto-default-negative-136 | 负向 | POST | `/api/v1/recommendations/jobs-for-talent` | 1 | 是 |
| auto-default-negative-137 | 负向 | POST | `/api/v1/recommendations/jobs-for-talent` | 1 | 是 |
| auto-default-negative-138 | 负向 | POST | `/api/v1/recommendations/rag/index-all` | 1 | 是 |
| auto-default-negative-139 | 负向 | POST | `/api/v1/recommendations/rag/index-all` | 1 | 是 |
| auto-default-negative-140 | 负向 | POST | `/api/v1/recommendations/rag/index-job` | 1 | 是 |
| auto-default-negative-141 | 负向 | POST | `/api/v1/recommendations/rag/index-job` | 1 | 是 |
| auto-default-negative-142 | 负向 | POST | `/api/v1/recommendations/rag/index-resume` | 1 | 是 |
| auto-default-negative-143 | 负向 | POST | `/api/v1/recommendations/rag/index-resume` | 1 | 是 |
| auto-default-negative-144 | 负向 | POST | `/api/v1/recommendations/rag/index-talent` | 1 | 是 |
| auto-default-negative-145 | 负向 | POST | `/api/v1/recommendations/rag/index-talent` | 1 | 是 |

### 回放结果

| 用例ID | 方法 | 路径 | 状态 | 期望 | 实际 | 耗时(ms) | 说明 |
| --- | --- | --- | --- | ---: | ---: | ---: | --- |
| auto-default-1 | POST | `/api/v1/login` | PASS | 200 | 200 | 59 | ok |
| auto-default-2 | GET | `/api/v1/stats/dashboard` | PASS | 200 | 200 | 53 | ok |
| auto-default-3 | GET | `/api/v1/messages/unread-count` | PASS | 200 | 200 | 53 | ok |
| auto-default-4 | GET | `/api/v1/messages?page=1&page_size=5` | PASS | 200 | 200 | 54 | ok |
| auto-default-5 | GET | `/api/v1/talents?page=1&page_size=20` | PASS | 200 | 200 | 58 | ok |
| auto-default-6 | GET | `/api/v1/jobs?page=1&page_size=5` | PASS | 200 | 200 | 54 | ok |
| auto-default-7 | GET | `/api/v1/conversations/unread-count` | PASS | 200 | 200 | 55 | ok |
| auto-default-8 | GET | `/api/v1/jobs?page=1&page_size=10&status=open` | PASS | 200 | 200 | 55 | ok |
| auto-default-9 | GET | `/api/v1/applications?page=1&page_size=1000&talent_id=me` | PASS | 200 | 200 | 56 | ok |
| auto-default-10 | GET | `/api/v1/conversations` | PASS | 200 | 200 | 55 | ok |
| auto-default-11 | GET | `/api/v1/jobs/1` | PASS | 200 | 200 | 55 | ok |
| auto-default-12 | GET | `/api/v1/applications?page=1&page_size=20` | PASS | 200 | 200 | 60 | ok |
| auto-default-13 | GET | `/api/v1/resumes/online` | PASS | 200 | 200 | 54 | ok |
| auto-default-14 | GET | `/api/v1/resumes?page_size=10` | PASS | 200 | 200 | 58 | ok |
| auto-default-15 | GET | `/api/v1/jobs?page_size=100&status=open` | PASS | 200 | 200 | 54 | ok |
| auto-default-16 | GET | `/api/v1/talents?page=1&page_size=10` | PASS | 200 | 200 | 55 | ok |
| auto-default-17 | GET | `/api/v1/jobs?page=1&page_size=10` | PASS | 200 | 200 | 54 | ok |
| auto-default-18 | GET | `/api/v1/jobs/1/applications?page=1&page_size=10` | PASS | 200 | 200 | 55 | ok |
| auto-default-19 | GET | `/api/v1/resumes?page=1&page_size=10` | PASS | 200 | 200 | 62 | ok |
| auto-default-20 | GET | `/api/v1/jobs?page=1&page_size=20` | PASS | 200 | 200 | 54 | ok |
| auto-default-21 | GET | `/api/v1/applications?page=1&page_size=50` | PASS | 200 | 200 | 66 | ok |
| auto-default-22 | GET | `/api/v1/interviews/1` | PASS | 200 | 200 | 56 | ok |
| auto-default-23 | GET | `/api/v1/interviews/1/feedback` | PASS | 200 | 200 | 55 | ok |
| auto-default-24 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | PASS | 200 | 200 | 56 | ok |
| auto-default-25 | GET | `/api/v1/logs?page=1&page_size=20` | PASS | 200 | 200 | 54 | ok |
| auto-default-26 | GET | `/api/v1/notices/1` | PASS | 404 | 404 | 55 | ok |
| auto-default-27 | GET | `/api/v1/notices?page=1&page_size=10` | PASS | 200 | 200 | 56 | ok |
| auto-default-openapi-28 | POST | `/api/v1/talents` | PASS | 201 | 201 | 57 | ok |
| auto-default-openapi-29 | GET | `/api/v1/talents/123` | PASS | 200 | 200 | 52 | ok |
| auto-default-openapi-30 | PUT | `/api/v1/talents/123` | FAIL | 200 | 500 | 53 | status mismatch |
| auto-default-openapi-31 | DELETE | `/api/v1/talents/123` | PASS | 200 | 200 | 54 | ok |
| auto-default-openapi-32 | GET | `/api/v1/talents/search?keyword=%E6%B5%8B%E8%AF%95` | PASS | 200 | 200 | 56 | ok |
| auto-default-openapi-33 | GET | `/api/v1/talents/stats` | PASS | 200 | 200 | 56 | ok |
| auto-default-openapi-34 | POST | `/api/v1/notices` | FAIL | 201 | 200 | 56 | status mismatch |
| auto-default-openapi-35 | PUT | `/api/v1/notices/1` | SKIP | 200 | 404 | 55 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-36 | DELETE | `/api/v1/notices/1` | SKIP | 200 | 404 | 54 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-37 | POST | `/api/v1/applications` | SKIP | 201 | 400 | 55 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-38 | PUT | `/api/v1/applications/1` | PASS | 200 | 200 | 56 | ok |
| auto-default-openapi-39 | DELETE | `/api/v1/applications/1` | SKIP | 200 | 403 | 54 | OpenAPI 自动用例缺少对应角色或资源权限 |
| auto-default-openapi-40 | POST | `/api/v1/recommendations/attribution-report` | PASS | 200 | 200 | 57 | ok |
| auto-default-openapi-41 | POST | `/api/v1/recommendations/batch` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-42 | POST | `/api/v1/recommendations/jobs-for-talent` | PASS | 200 | 200 | 54 | ok |
| auto-default-openapi-43 | POST | `/api/v1/recommendations/rag/index-all` | PASS | 200 | 200 | 106 | ok |
| auto-default-openapi-44 | POST | `/api/v1/recommendations/rag/index-job` | FAIL | 200 | 500 | 53 | status mismatch |
| auto-default-openapi-45 | POST | `/api/v1/recommendations/rag/index-resume` | SKIP | 200 | 400 | 52 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-46 | POST | `/api/v1/recommendations/rag/index-talent` | FAIL | 200 | 500 | 53 | status mismatch |
| auto-default-openapi-47 | POST | `/api/v1/recommendations/rag/match` | PASS | 200 | 200 | 69 | ok |
| auto-default-openapi-48 | POST | `/api/v1/recommendations/rag/query` | PASS | 200 | 200 | 60 | ok |
| auto-default-openapi-49 | POST | `/api/v1/recommendations/semantic-match` | SKIP | 200 | 400 | 54 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-50 | GET | `/api/v1/recommendations/stats` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-51 | POST | `/api/v1/recommendations/talents-for-job` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-52 | GET | `/api/v1/logs/actions` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-53 | DELETE | `/api/v1/logs/cleanup` | FAIL | 200 | 500 | 54 | status mismatch |
| auto-default-openapi-54 | GET | `/api/v1/logs/services` | PASS | 200 | 200 | 53 | ok |
| auto-default-openapi-55 | GET | `/api/v1/logs/stats` | PASS | 200 | 200 | 54 | ok |
| auto-default-openapi-56 | POST | `/api/v1/conversations` | SKIP | 200 | 400 | 54 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-57 | GET | `/api/v1/conversations/1/messages?page=1&page_size=10` | SKIP | 200 | 403 | 55 | OpenAPI 自动用例缺少对应角色或资源权限 |
| auto-default-openapi-58 | POST | `/api/v1/conversations/1/messages` | SKIP | 201 | 403 | 55 | OpenAPI 自动用例缺少对应角色或资源权限 |
| auto-default-openapi-59 | PUT | `/api/v1/conversations/1/read` | SKIP | 200 | 403 | 54 | OpenAPI 自动用例缺少对应角色或资源权限 |
| auto-default-openapi-60 | POST | `/api/v1/messages` | PASS | 201 | 201 | 53 | ok |
| auto-default-openapi-61 | DELETE | `/api/v1/messages/58` | PASS | 200 | 200 | 53 | ok |
| auto-default-openapi-62 | PUT | `/api/v1/messages/58/read` | SKIP | 200 | 404 | 54 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-63 | GET | `/api/v1/messages/stats` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-64 | GET | `/api/v1/online-status` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-65 | GET | `/api/v1/online-status/1` | SKIP | 200 | 400 | 55 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-66 | GET | `/api/v1/profile` | PASS | 200 | 200 | 62 | ok |
| auto-default-openapi-67 | PUT | `/api/v1/profile` | PASS | 200 | 200 | 87 | ok |
| auto-default-openapi-68 | GET | `/api/v1/users?page=1&page_size=10` | PASS | 200 | 200 | 56 | ok |
| auto-default-openapi-69 | POST | `/api/v1/resumes` | PASS | 201 | 201 | 58 | ok |
| auto-default-openapi-70 | DELETE | `/api/v1/resumes/168` | PASS | 200 | 200 | 56 | ok |
| auto-default-openapi-71 | GET | `/api/v1/resumes/168/download` | SKIP | 200 | 404 | 55 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-72 | PUT | `/api/v1/resumes/168/job` | SKIP | 200 | 404 | 55 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-73 | PUT | `/api/v1/resumes/168/status` | SKIP | 200 | 404 | 53 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-74 | GET | `/api/v1/resumes/evaluation?page=1&page_size=10&status=open` | PASS | 200 | 200 | 54 | ok |
| auto-default-openapi-75 | GET | `/api/v1/resumes/file/test.pdf` | SKIP | 200 | 404 | 52 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-76 | POST | `/api/v1/resumes/match` | SKIP | 200 | 400 | 52 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-77 | PUT | `/api/v1/resumes/online` | SKIP | 200 | 400 | 53 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-78 | POST | `/api/v1/resumes/parse` | SKIP | 200 | 400 | 53 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-79 | POST | `/api/v1/resumes/risk-check` | PASS | 200 | 200 | 51 | ok |
| auto-default-openapi-80 | POST | `/api/v1/resumes/risk-check/education` | PASS | 200 | 200 | 53 | ok |
| auto-default-openapi-81 | POST | `/api/v1/resumes/risk-check/time-conflict` | PASS | 200 | 200 | 54 | ok |
| auto-default-openapi-82 | GET | `/api/v1/stats/channels` | PASS | 200 | 200 | 54 | ok |
| auto-default-openapi-83 | GET | `/api/v1/stats/department-progress` | PASS | 200 | 200 | 54 | ok |
| auto-default-openapi-84 | GET | `/api/v1/stats/funnel` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-85 | GET | `/api/v1/stats/interviewer-rank` | PASS | 200 | 200 | 54 | ok |
| auto-default-openapi-86 | GET | `/api/v1/stats/job-rank` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-87 | GET | `/api/v1/stats/trend` | PASS | 200 | 200 | 65 | ok |
| auto-default-openapi-88 | POST | `/api/v1/jobs` | PASS | 201 | 201 | 57 | ok |
| auto-default-openapi-89 | PUT | `/api/v1/jobs/68` | PASS | 200 | 200 | 54 | ok |
| auto-default-openapi-90 | DELETE | `/api/v1/jobs/68` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-91 | GET | `/api/v1/jobs/hot` | PASS | 200 | 200 | 53 | ok |
| auto-default-openapi-92 | GET | `/api/v1/jobs/stats` | PASS | 200 | 200 | 54 | ok |
| auto-default-openapi-93 | POST | `/api/v1/register` | SKIP | 201 | 400 | 52 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-94 | GET | `/api/v1/interviews?page=1&page_size=10&status=open` | PASS | 200 | 200 | 56 | ok |
| auto-default-openapi-95 | POST | `/api/v1/interviews` | PASS | 201 | 201 | 56 | ok |
| auto-default-openapi-96 | PUT | `/api/v1/interviews/118` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-97 | DELETE | `/api/v1/interviews/118` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-98 | POST | `/api/v1/interviews/118/cancel` | SKIP | 200 | 404 | 56 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-99 | POST | `/api/v1/interviews/118/complete` | SKIP | 200 | 404 | 54 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-100 | POST | `/api/v1/interviews/118/feedback` | SKIP | 200 | 404 | 52 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-101 | POST | `/api/v1/interviews/118/reschedule` | SKIP | 200 | 404 | 54 | OpenAPI 自动用例缺少业务前置数据或字段模板 |
| auto-default-openapi-102 | GET | `/api/v1/interviews/candidate/123` | PASS | 200 | 200 | 53 | ok |
| auto-default-openapi-103 | GET | `/api/v1/interviews/interviewer/118` | PASS | 200 | 200 | 55 | ok |
| auto-default-openapi-104 | GET | `/api/v1/interviews/stats` | PASS | 200 | 200 | 56 | ok |
| auto-default-openapi-105 | GET | `/api/v1/interviews/today` | PASS | 200 | 200 | 54 | ok |
| auto-default-negative-106 | GET | `/api/v1/talents?page=1&page_size=10&status=open&search=%E6%B5%8B%E8%AF%95` | PASS | 401, 403 | 401 | 52 | ok |
| auto-default-negative-107 | POST | `/api/v1/talents` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-108 | POST | `/api/v1/talents` | PASS | 400, 401, 403, 404, 422 | 400 | 52 | ok |
| auto-default-negative-109 | GET | `/api/v1/talents/99999999` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-110 | PUT | `/api/v1/talents/99999999` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-111 | PUT | `/api/v1/talents/99999999` | PASS | 400, 401, 403, 404, 422 | 404 | 57 | ok |
| auto-default-negative-112 | DELETE | `/api/v1/talents/99999999` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-113 | DELETE | `/api/v1/talents/99999999` | FAIL | 400, 401, 403, 404, 422 | 200 | 55 | status mismatch |
| auto-default-negative-114 | GET | `/api/v1/talents/search?keyword=%E6%B5%8B%E8%AF%95` | PASS | 401, 403 | 401 | 56 | ok |
| auto-default-negative-115 | GET | `/api/v1/talents/stats` | PASS | 401, 403 | 401 | 54 | ok |
| auto-default-negative-116 | GET | `/api/v1/notices?page=1&page_size=10&keyword=%E6%B5%8B%E8%AF%95&status=open` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-117 | POST | `/api/v1/notices` | PASS | 401, 403 | 401 | 54 | ok |
| auto-default-negative-118 | POST | `/api/v1/notices` | PASS | 400, 401, 403, 404, 422 | 400 | 54 | ok |
| auto-default-negative-119 | GET | `/api/v1/notices/99999999` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-120 | PUT | `/api/v1/notices/99999999` | PASS | 401, 403 | 401 | 54 | ok |
| auto-default-negative-121 | PUT | `/api/v1/notices/99999999` | PASS | 400, 401, 403, 404, 422 | 400 | 53 | ok |
| auto-default-negative-122 | DELETE | `/api/v1/notices/99999999` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-123 | DELETE | `/api/v1/notices/99999999` | PASS | 400, 401, 403, 404, 422 | 404 | 54 | ok |
| auto-default-negative-124 | POST | `/api/v1/login` | FAIL | 400, 404, 422 | 200 | 58 | status mismatch |
| auto-default-negative-125 | GET | `/api/v1/applications?page=1&page_size=10&status=open` | PASS | 401, 403 | 401 | 51 | ok |
| auto-default-negative-126 | POST | `/api/v1/applications` | PASS | 401, 403 | 401 | 54 | ok |
| auto-default-negative-127 | POST | `/api/v1/applications` | PASS | 400, 401, 403, 404, 422 | 400 | 54 | ok |
| auto-default-negative-128 | PUT | `/api/v1/applications/99999999` | PASS | 401, 403 | 401 | 54 | ok |
| auto-default-negative-129 | PUT | `/api/v1/applications/99999999` | PASS | 400, 401, 403, 404, 422 | 404 | 56 | ok |
| auto-default-negative-130 | DELETE | `/api/v1/applications/99999999` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-131 | DELETE | `/api/v1/applications/99999999` | PASS | 400, 401, 403, 404, 422 | 404 | 55 | ok |
| auto-default-negative-132 | POST | `/api/v1/recommendations/attribution-report` | PASS | 401, 403 | 401 | 54 | ok |
| auto-default-negative-133 | POST | `/api/v1/recommendations/attribution-report` | PASS | 400, 401, 403, 404, 422 | 400 | 53 | ok |
| auto-default-negative-134 | POST | `/api/v1/recommendations/batch` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-135 | POST | `/api/v1/recommendations/batch` | FAIL | 400, 401, 403, 404, 422 | 200 | 55 | status mismatch |
| auto-default-negative-136 | POST | `/api/v1/recommendations/jobs-for-talent` | PASS | 401, 403 | 401 | 54 | ok |
| auto-default-negative-137 | POST | `/api/v1/recommendations/jobs-for-talent` | PASS | 400, 401, 403, 404, 422 | 400 | 55 | ok |
| auto-default-negative-138 | POST | `/api/v1/recommendations/rag/index-all` | PASS | 401, 403 | 401 | 52 | ok |
| auto-default-negative-139 | POST | `/api/v1/recommendations/rag/index-all` | FAIL | 400, 401, 403, 404, 422 | 200 | 89 | status mismatch |
| auto-default-negative-140 | POST | `/api/v1/recommendations/rag/index-job` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-141 | POST | `/api/v1/recommendations/rag/index-job` | PASS | 400, 401, 403, 404, 422 | 400 | 53 | ok |
| auto-default-negative-142 | POST | `/api/v1/recommendations/rag/index-resume` | PASS | 401, 403 | 401 | 51 | ok |
| auto-default-negative-143 | POST | `/api/v1/recommendations/rag/index-resume` | PASS | 400, 401, 403, 404, 422 | 400 | 53 | ok |
| auto-default-negative-144 | POST | `/api/v1/recommendations/rag/index-talent` | PASS | 401, 403 | 401 | 53 | ok |
| auto-default-negative-145 | POST | `/api/v1/recommendations/rag/index-talent` | PASS | 400, 401, 403, 404, 422 | 400 | 54 | ok |

### 失败分层明细

| 分类 | 用例 | 方法 | 路径 | 实际状态 | 原因 |
| --- | --- | --- | --- | ---: | --- |
| 真实失败 | auto-default-openapi-30 | PUT | `/api/v1/talents/123` | 500 | 接口返回 5xx 或请求执行失败 |
| 契约不一致 | auto-default-openapi-34 | POST | `/api/v1/notices` | 200 | OpenAPI 期望状态和实际业务返回不一致 |
| 依赖/环境 | auto-default-openapi-44 | POST | `/api/v1/recommendations/rag/index-job` | 500 | 依赖服务、数据库结构或外部组件未就绪 |
| 依赖/环境 | auto-default-openapi-46 | POST | `/api/v1/recommendations/rag/index-talent` | 500 | 依赖服务、数据库结构或外部组件未就绪 |
| 依赖/环境 | auto-default-openapi-53 | DELETE | `/api/v1/logs/cleanup` | 500 | 依赖服务、数据库结构或外部组件未就绪 |
| 负向失败 | auto-default-negative-113 | DELETE | `/api/v1/talents/99999999` | 200 | 负向用例没有被明确拒绝，或触发了 5xx |
| 负向失败 | auto-default-negative-124 | POST | `/api/v1/login` | 200 | 负向用例没有被明确拒绝，或触发了 5xx |
| 负向失败 | auto-default-negative-135 | POST | `/api/v1/recommendations/batch` | 200 | 负向用例没有被明确拒绝，或触发了 5xx |
| 负向失败 | auto-default-negative-139 | POST | `/api/v1/recommendations/rag/index-all` | 200 | 负向用例没有被明确拒绝，或触发了 5xx |

### 结果分析

| 指标 | 结果 | 说明 |
| --- | ---: | --- |
| 页面访问数 | 28 | 自动探索实际打开的前端页面数量 |
| 输入框填充 | 133/437 | 仅填写可见、可编辑且能识别字段含义的输入框 |
| 按钮点击 | 160/541 | 按 clickRiskLevels 配置执行，测试环境可允许 danger/empty 按钮 |
| API 资源数 | 15 | 按 /api/v1/{resource} 归类统计 |
| 用例来源 | 27/78/40 | 真实流量/OpenAPI 自动生成/负向异常用例 |
| 失败分层 | 1/3/4 | 真实失败/依赖环境失败/负向用例失败 |
| 回放失败 | 9 | 状态码、响应结构或请求错误不符合预期 |
| 回放跳过 | 24 | 多为非 BIGTEST 自有写请求，避免污染真实数据 |

### 跳过原因统计

| 风险类型 | 数量 | 说明 |
| --- | ---: | --- |
| disabled | 19 | 未分类操作 |

### 失败接口复现命令

#### auto-default-openapi-30 PUT /api/v1/talents/123

```bash
curl -i -X 'PUT' 'http://localhost:8080/api/v1/talents/123' -H 'accept: application/json' -H 'content-type: application/json' -H 'Authorization: <redacted>' --data '{
  "name": "BIGTEST_2026-04-30_16-17-30_自动人才",
  "email": "bigtest_2026-04-30_16-17-30@example.com",
  "phone": "13800000000",
  "skills": [
    "Go",
    "Vue"
  ],
  "experience": 3,
  "education": "本科",
  "status": "active",
  "tags": [
    "BIGTEST"
  ],
  "location": "测试城市",
  "salary": "15k-25k",
  "summary": "BIGTEST_2026-04-30_16-17-30_自动人才摘要",
  "gender": "未知",
  "age": 28,
  "current_company": "BIGTEST 公司",
  "current_position": "后端工程师",
  "source": "bigtest"
}'
```

#### auto-default-openapi-34 POST /api/v1/notices

```bash
curl -i -X 'POST' 'http://localhost:8080/api/v1/notices' -H 'accept: application/json' -H 'content-type: application/json' -H 'Authorization: <redacted>' --data '{
  "title": "BIGTEST_2026-04-30_16-17-30_公告",
  "content": "BIGTEST_2026-04-30_16-17-30_自动公告内容",
  "status": "draft",
  "is_pinned": false,
  "priority": "normal"
}'
```

#### auto-default-openapi-44 POST /api/v1/recommendations/rag/index-job

```bash
curl -i -X 'POST' 'http://localhost:8080/api/v1/recommendations/rag/index-job' -H 'accept: application/json' -H 'content-type: application/json' -H 'Authorization: <redacted>' --data '{
  "job_id": 1,
  "talent_id": 1,
  "limit": 10,
  "query": "Go 后端"
}'
```

#### auto-default-openapi-46 POST /api/v1/recommendations/rag/index-talent

```bash
curl -i -X 'POST' 'http://localhost:8080/api/v1/recommendations/rag/index-talent' -H 'accept: application/json' -H 'content-type: application/json' -H 'Authorization: <redacted>' --data '{
  "job_id": 1,
  "talent_id": 1,
  "limit": 10,
  "query": "Go 后端"
}'
```

#### auto-default-openapi-53 DELETE /api/v1/logs/cleanup

```bash
curl -i -X 'DELETE' 'http://localhost:8080/api/v1/logs/cleanup' -H 'accept: application/json' -H 'content-type: application/json' -H 'Authorization: <redacted>'
```

#### auto-default-negative-113 DELETE /api/v1/talents/99999999

```bash
curl -i -X 'DELETE' 'http://localhost:8080/api/v1/talents/99999999' -H 'accept: application/json' -H 'content-type: application/json' -H 'Authorization: <redacted>'
```

#### auto-default-negative-124 POST /api/v1/login

```bash
curl -i -X 'POST' 'http://localhost:8080/api/v1/login' -H 'accept: application/json' -H 'content-type: application/json' --data '{
  "username": "admin",
  "password": "<redacted>"
}'
```

#### auto-default-negative-135 POST /api/v1/recommendations/batch

```bash
curl -i -X 'POST' 'http://localhost:8080/api/v1/recommendations/batch' -H 'accept: application/json' -H 'content-type: application/json' -H 'Authorization: <redacted>' --data '{
  "invalid": true,
  "id": "abc",
  "status": "__invalid__"
}'
```

#### auto-default-negative-139 POST /api/v1/recommendations/rag/index-all

```bash
curl -i -X 'POST' 'http://localhost:8080/api/v1/recommendations/rag/index-all' -H 'accept: application/json' -H 'content-type: application/json' -H 'Authorization: <redacted>' --data '{
  "invalid": true,
  "id": "abc",
  "status": "__invalid__"
}'
```


### 错误现场

| 类型 | 数量 | 说明 |
| --- | ---: | --- |
| 页面失败 | 0 | 自动探索打开页面失败时保存截图和错误 |
| 按钮失败 | 0 | 页面按钮点击失败时记录 URL、错误，按截图策略可选保存图片 |
| 控制台错误/警告 | 167 | 来源于 Playwright console 事件 |
| 网络失败请求 | 0 | 来源于 Playwright requestfailed 事件 |

| 控制台类型 | 页面 | 内容 |
| --- | --- | --- |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". null at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at（重复 16 次） |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". null at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object]（重复 16 次） |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <Router（重复 48 次） |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object]（重复 48 次） |
| error | `http://localhost:5173/dev-logs` | Failed to load resource: the server responded with a status of 500 (Internal Server Error)（重复 23 次） |
| warning | `http://localhost:5173/messages` | [Vue warn]: Invalid prop: type check failed for prop "modelValue". Expected Number \| String \| Boolean, got Array at <ElCheckbox modelValue= [] onUpdate:modelValue=fn label=50 ... > at <MessageCenter onVnodeUnmounted=fn<o |
| warning | `http://localhost:5173/messages` | [Vue warn]: Invalid prop: type check failed for prop "modelValue". Expected Number \| String \| Boolean, got Array at <ElCheckbox modelValue= [] onUpdate:modelValue=fn label=35 ... > at <MessageCenter onVnodeUnmounted=fn<o |
| warning | `http://localhost:5173/messages` | [Vue warn]: Invalid prop: type check failed for prop "modelValue". Expected Number \| String \| Boolean, got Array at <ElCheckbox modelValue= [] onUpdate:modelValue=fn label=36 ... > at <MessageCenter onVnodeUnmounted=fn<o |
| warning | `http://localhost:5173/messages` | [Vue warn]: Invalid prop: type check failed for prop "modelValue". Expected Number \| String \| Boolean, got Array at <ElCheckbox modelValue= [] onUpdate:modelValue=fn label=37 ... > at <MessageCenter onVnodeUnmounted=fn<o |
| warning | `http://localhost:5173/messages` | [Vue warn]: Invalid prop: type check failed for prop "modelValue". Expected Number \| String \| Boolean, got Array at <ElCheckbox modelValue= [] onUpdate:modelValue=fn label=38 ... > at <MessageCenter onVnodeUnmounted=fn<o |
| warning | `http://localhost:5173/messages` | [Vue warn]: Invalid prop: type check failed for prop "modelValue". Expected Number \| String \| Boolean, got Array at <ElCheckbox modelValue= [] onUpdate:modelValue=fn label=39 ... > at <MessageCenter onVnodeUnmounted=fn<o |
| warning | `http://localhost:5173/messages` | [Vue warn]: Invalid prop: type check failed for prop "modelValue". Expected Number \| String \| Boolean, got Array at <ElCheckbox modelValue= [] onUpdate:modelValue=fn label=24 ... > at <MessageCenter onVnodeUnmounted=fn<o |
| warning | `http://localhost:5173/messages` | [Vue warn]: Invalid prop: type check failed for prop "modelValue". Expected Number \| String \| Boolean, got Array at <ElCheckbox modelValue= [] onUpdate:modelValue=fn label=25 ... > at <MessageCenter onVnodeUnmounted=fn<o |
| warning | `http://localhost:5173/messages` | [Vue warn]: Invalid prop: type check failed for prop "modelValue". Expected Number \| String \| Boolean, got Array at <ElCheckbox modelValue= [] onUpdate:modelValue=fn label=23 ... > at <MessageCenter onVnodeUnmounted=fn<o |
| warning | `http://localhost:5173/messages` | [Vue warn]: Invalid prop: type check failed for prop "modelValue". Expected Number \| String \| Boolean, got Array at <ElCheckbox modelValue= [] onUpdate:modelValue=fn label=15 ... > at <MessageCenter onVnodeUnmounted=fn<o |
| error | `http://localhost:5173/notices/1/edit` | Failed to load resource: the server responded with a status of 404 (Not Found) |
| error | `http://localhost:5173/notices/1/edit` | [Response Error] /notices/1 Request failed with status code 404 |
| error | `http://localhost:5173/notices/1/edit` | [Response Error] Status: 404 Data: {code: 404, message: Notice not found} |
| error | `http://localhost:5173/notices/1` | Failed to load resource: the server responded with a status of 404 (Not Found) |
| error | `http://localhost:5173/notices/1` | [Response Error] /notices/1 Request failed with status code 404 |

### 后续建议

- 优先处理回放失败和写操作失败项，查看对应 JSON 明细定位请求与响应。
