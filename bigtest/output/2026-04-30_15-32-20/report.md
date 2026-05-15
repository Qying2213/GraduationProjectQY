# bigtest 真实流量自动化测试报告

- Run ID: `2026-04-30_15-32-20`
- 配置文件: `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/config.local.json`
- 输出目录: `output/2026-04-30_15-32-20`
- 自动探索: 已开启
- 写操作场景: 已开启，测试数据前缀 `BIGTEST_2026-04-30_15-32-20`

## 测试结论

本次测试共访问 28 个页面，捕获 212 条真实 API 请求，生成 27 条接口回放用例。
接口回放通过 26 条，失败 1 条，跳过 0 条；写操作场景通过 4 个，失败 1 个。

**结论：本次自动化测试存在失败项，请优先查看“回放结果”和“写操作场景步骤”。**

## 测试环境与策略

| 项目 | 值 |
| --- | --- |
| 前端地址 | `http://localhost:5173` |
| 后端回放地址 | `http://localhost:8080` |
| 捕获 URL 规则 | `/api/v1/` |
| 回放方法 | `GET, POST, PUT, DELETE` |
| 排除路径 | `/api/v1/ai, /api/v1/evaluations, /upload` |
| 页面探索最大路由数 | 120 |
| 每页最大点击数 | 15 |
| 截图策略 | 最多 5 张；按钮点击失败默认不截图，只记录 URL 和错误 |
| AI 评估策略 | 明确排除 `/api/v1/ai`、`/api/v1/evaluations`、`/ai/evaluate`、`/evaluate` 等路径，不触发真实 AI 评估、上传或外部模型调用 |

## 总览

| Profile | 页面 | 输入 | 按钮 | 原始请求 | 生成用例 | 安全写通过 | 安全写失败 | 回放通过 | 回放失败 | 回放跳过 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| auto-default | 28 | 437/133 | 542/160 | 212 | 27 | 4 | 1 | 26 | 1 | 0 |

## 模块覆盖统计

| 模块 | 页面访问 | API 用例 | GET | POST | PUT | PATCH | DELETE | 未覆盖接口 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 公告 | 3 | 6 | 3 | 1 | 1 | 0 | 1 | 0 |
| 简历 | 2 | 3 | 3 | 0 | 0 | 0 | 0 | 13 |
| 面试 | 1 | 8 | 3 | 3 | 1 | 0 | 1 | 7 |
| 其他 | 12 | 1 | 0 | 1 | 0 | 0 | 0 | 0 |
| 人才 | 2 | 6 | 3 | 1 | 1 | 0 | 1 | 2 |
| 认证 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 1 |
| 日志 | 1 | 1 | 1 | 0 | 0 | 0 | 0 | 4 |
| 统计 | 0 | 1 | 1 | 0 | 0 | 0 | 0 | 6 |
| 投递 | 1 | 3 | 3 | 0 | 0 | 0 | 0 | 3 |
| 推荐 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 12 |
| 消息 | 1 | 8 | 6 | 1 | 0 | 0 | 1 | 8 |
| 用户 | 1 | 0 | 0 | 0 | 0 | 0 | 0 | 3 |
| 职位 | 4 | 11 | 8 | 1 | 1 | 0 | 1 | 2 |
| ws | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 1 |

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

| 模块 | 方法 | 路径 |
| --- | --- | --- |
| ws | GET | `/api/v1/ws` |
| 人才 | GET | `/api/v1/talents/search` |
| 人才 | GET | `/api/v1/talents/stats` |
| 投递 | DELETE | `/api/v1/applications/{id}` |
| 投递 | PUT | `/api/v1/applications/{id}` |
| 投递 | POST | `/api/v1/applications` |
| 推荐 | POST | `/api/v1/recommendations/attribution-report` |
| 推荐 | POST | `/api/v1/recommendations/batch` |
| 推荐 | POST | `/api/v1/recommendations/jobs-for-talent` |
| 推荐 | POST | `/api/v1/recommendations/rag/index-all` |
| 推荐 | POST | `/api/v1/recommendations/rag/index-job` |
| 推荐 | POST | `/api/v1/recommendations/rag/index-resume` |
| 推荐 | POST | `/api/v1/recommendations/rag/index-talent` |
| 推荐 | POST | `/api/v1/recommendations/rag/match` |
| 推荐 | POST | `/api/v1/recommendations/rag/query` |
| 推荐 | POST | `/api/v1/recommendations/semantic-match` |
| 推荐 | GET | `/api/v1/recommendations/stats` |
| 推荐 | POST | `/api/v1/recommendations/talents-for-job` |
| 日志 | GET | `/api/v1/logs/actions` |
| 日志 | DELETE | `/api/v1/logs/cleanup` |
| 日志 | GET | `/api/v1/logs/services` |
| 日志 | GET | `/api/v1/logs/stats` |
| 消息 | GET | `/api/v1/conversations/{id}/messages` |
| 消息 | POST | `/api/v1/conversations/{id}/messages` |
| 消息 | PUT | `/api/v1/conversations/{id}/read` |
| 消息 | POST | `/api/v1/conversations` |
| 消息 | PUT | `/api/v1/messages/{id}/read` |
| 消息 | GET | `/api/v1/messages/stats` |
| 消息 | GET | `/api/v1/online-status/{user_id}` |
| 消息 | GET | `/api/v1/online-status` |
| 用户 | GET | `/api/v1/profile` |
| 用户 | PUT | `/api/v1/profile` |
| 用户 | GET | `/api/v1/users` |
| 简历 | GET | `/api/v1/resumes/{id}/download` |
| 简历 | PUT | `/api/v1/resumes/{id}/job` |
| 简历 | PUT | `/api/v1/resumes/{id}/status` |
| 简历 | DELETE | `/api/v1/resumes/{id}` |
| 简历 | GET | `/api/v1/resumes/evaluation` |
| 简历 | GET | `/api/v1/resumes/file/{filename}` |
| 简历 | POST | `/api/v1/resumes/match` |
| 简历 | PUT | `/api/v1/resumes/online` |
| 简历 | POST | `/api/v1/resumes/parse` |
| 简历 | POST | `/api/v1/resumes/risk-check/education` |
| 简历 | POST | `/api/v1/resumes/risk-check/time-conflict` |
| 简历 | POST | `/api/v1/resumes/risk-check` |
| 简历 | POST | `/api/v1/resumes` |
| 统计 | GET | `/api/v1/stats/channels` |
| 统计 | GET | `/api/v1/stats/department-progress` |
| 统计 | GET | `/api/v1/stats/funnel` |
| 统计 | GET | `/api/v1/stats/interviewer-rank` |
| 统计 | GET | `/api/v1/stats/job-rank` |
| 统计 | GET | `/api/v1/stats/trend` |
| 职位 | GET | `/api/v1/jobs/hot` |
| 职位 | GET | `/api/v1/jobs/stats` |
| 认证 | POST | `/api/v1/register` |
| 面试 | POST | `/api/v1/interviews/{id}/cancel` |
| 面试 | POST | `/api/v1/interviews/{id}/feedback` |
| 面试 | GET | `/api/v1/interviews/candidate/{candidate_id}` |
| 面试 | GET | `/api/v1/interviews/interviewer/{interviewer_id}` |
| 面试 | GET | `/api/v1/interviews/stats` |
| 面试 | GET | `/api/v1/interviews/today` |
| 面试 | GET | `/api/v1/interviews` |

### 未覆盖写操作

| 模块 | 方法 | 路径 | 建议 |
| --- | --- | --- | --- |
| 投递 | DELETE | `/api/v1/applications/{id}` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 投递 | PUT | `/api/v1/applications/{id}` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 投递 | POST | `/api/v1/applications` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/attribution-report` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/batch` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/jobs-for-talent` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/rag/index-all` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/rag/index-job` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/rag/index-resume` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/rag/index-talent` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/rag/match` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/rag/query` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/semantic-match` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 推荐 | POST | `/api/v1/recommendations/talents-for-job` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 日志 | DELETE | `/api/v1/logs/cleanup` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 消息 | POST | `/api/v1/conversations/{id}/messages` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 消息 | PUT | `/api/v1/conversations/{id}/read` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 消息 | POST | `/api/v1/conversations` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 消息 | PUT | `/api/v1/messages/{id}/read` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 用户 | PUT | `/api/v1/profile` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 简历 | PUT | `/api/v1/resumes/{id}/job` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 简历 | PUT | `/api/v1/resumes/{id}/status` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 简历 | DELETE | `/api/v1/resumes/{id}` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 简历 | POST | `/api/v1/resumes/match` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 简历 | PUT | `/api/v1/resumes/online` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 简历 | POST | `/api/v1/resumes/parse` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 简历 | POST | `/api/v1/resumes/risk-check/education` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 简历 | POST | `/api/v1/resumes/risk-check/time-conflict` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 简历 | POST | `/api/v1/resumes/risk-check` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 简历 | POST | `/api/v1/resumes` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 认证 | POST | `/api/v1/register` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 面试 | POST | `/api/v1/interviews/{id}/cancel` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |
| 面试 | POST | `/api/v1/interviews/{id}/feedback` | 建议补充 writeSafety.scenarios，只操作 BIGTEST 数据 |

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
- 按钮/链接: 542 个，安全点击 160 个
- 自动取消确认弹窗: 0
- API 资源数: 10
- 跳过/排除: 0

### API 资源覆盖

| 资源 | 方法 | 用例数 |
| --- | --- | ---: |
| stats | GET | 1 |
| messages | DELETE, GET, POST | 6 |
| talents | DELETE, GET, POST, PUT | 6 |
| jobs | DELETE, GET, POST, PUT | 11 |
| conversations | GET | 2 |
| applications | GET | 3 |
| resumes | GET | 3 |
| interviews | DELETE, GET, POST, PUT | 8 |
| logs | GET | 1 |
| notices | DELETE, GET, POST, PUT | 6 |
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
| job-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/jobs/64` | 200 | 200 | PASS |
| job-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/jobs/64` | 200 | 200 | PASS |
| job-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/jobs/64` | 200, 204 | 200 | PASS |
| talent-crud-owned-data | create 创建测试数据 | POST | `/api/v1/talents` | 200, 201 | 201 | PASS |
| talent-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/talents/117` | 200 | 200 | PASS |
| talent-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/talents/117` | 200 | 500 | FAIL |
| talent-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/talents/117` | 200, 204 | 200 | PASS |
| message-create-query-delete-owned-data | create 创建测试数据 | POST | `/api/v1/messages` | 200, 201 | 201 | PASS |
| message-create-query-delete-owned-data | verify 验证测试数据 | GET | `/api/v1/messages` | 200 | 200 | PASS |
| message-create-query-delete-owned-data | delete 删除测试数据 | DELETE | `/api/v1/messages/43` | 200, 204 | 200 | PASS |
| interview-crud-owned-data | create 创建测试数据 | POST | `/api/v1/interviews` | 200, 201 | 201 | PASS |
| interview-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/interviews/114` | 200 | 200 | PASS |
| interview-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/interviews/114` | 200 | 200 | PASS |
| interview-crud-owned-data | reschedule 改期测试面试 | POST | `/api/v1/interviews/114/reschedule` | 200 | 200 | PASS |
| interview-crud-owned-data | complete 完成测试面试 | POST | `/api/v1/interviews/114/complete` | 200 | 200 | PASS |
| interview-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/interviews/114` | 200, 204 | 200 | PASS |
| notice-crud-owned-data | create 创建测试数据 | POST | `/api/v1/notices` | 200, 201 | 200 | PASS |
| notice-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/notices/31` | 200 | 200 | PASS |
| notice-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/notices/31` | 200 | 200 | PASS |
| notice-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/notices/31` | 200, 204 | 200 | PASS |

### 去重后生成的用例

| 用例ID | 方法 | 路径 | 归并次数 | 需鉴权 |
| --- | --- | --- | ---: | --- |
| auto-default-1 | POST | `/api/v1/login` | 1 | 否 |
| auto-default-2 | GET | `/api/v1/stats/dashboard` | 3 | 是 |
| auto-default-3 | GET | `/api/v1/messages/unread-count` | 22 | 是 |
| auto-default-4 | GET | `/api/v1/messages?page=1&page_size=5` | 3 | 是 |
| auto-default-5 | GET | `/api/v1/talents?page=1&page_size=20` | 3 | 是 |
| auto-default-6 | GET | `/api/v1/jobs?page=1&page_size=5` | 3 | 否 |
| auto-default-7 | GET | `/api/v1/conversations/unread-count` | 63 | 是 |
| auto-default-8 | GET | `/api/v1/jobs?page=1&page_size=10&status=open` | 28 | 是 |
| auto-default-9 | GET | `/api/v1/applications?page=1&page_size=1000&talent_id=me` | 29 | 是 |
| auto-default-10 | GET | `/api/v1/conversations` | 17 | 是 |
| auto-default-11 | GET | `/api/v1/jobs/1` | 2 | 是 |
| auto-default-12 | GET | `/api/v1/applications?page=1&page_size=20` | 1 | 是 |
| auto-default-13 | GET | `/api/v1/resumes/online` | 1 | 是 |
| auto-default-14 | GET | `/api/v1/resumes?page_size=10` | 1 | 是 |
| auto-default-15 | GET | `/api/v1/jobs?page_size=100&status=open` | 22 | 是 |
| auto-default-16 | GET | `/api/v1/talents?page=1&page_size=10` | 1 | 是 |
| auto-default-17 | GET | `/api/v1/jobs?page=1&page_size=10` | 1 | 是 |
| auto-default-18 | GET | `/api/v1/jobs/1/applications?page=1&page_size=10` | 1 | 是 |
| auto-default-19 | GET | `/api/v1/resumes?page=1&page_size=10` | 1 | 是 |
| auto-default-20 | GET | `/api/v1/jobs?page=1&page_size=20` | 1 | 是 |
| auto-default-21 | GET | `/api/v1/applications?page=1&page_size=50` | 1 | 是 |
| auto-default-22 | GET | `/api/v1/interviews/1` | 1 | 是 |
| auto-default-23 | GET | `/api/v1/interviews/1/feedback` | 1 | 是 |
| auto-default-24 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | 1 | 是 |
| auto-default-25 | GET | `/api/v1/logs?page=1&page_size=20` | 1 | 是 |
| auto-default-26 | GET | `/api/v1/notices/1` | 2 | 是 |
| auto-default-27 | GET | `/api/v1/notices?page=1&page_size=10` | 1 | 是 |

### 回放结果

| 用例ID | 方法 | 路径 | 状态 | 期望 | 实际 | 耗时(ms) | 说明 |
| --- | --- | --- | --- | ---: | ---: | ---: | --- |
| auto-default-1 | POST | `/api/v1/login` | FAIL | 200 | 401 | 6 | status mismatch |
| auto-default-2 | GET | `/api/v1/stats/dashboard` | PASS | 200 | 200 | 2 | ok |
| auto-default-3 | GET | `/api/v1/messages/unread-count` | PASS | 200 | 200 | 2 | ok |
| auto-default-4 | GET | `/api/v1/messages?page=1&page_size=5` | PASS | 200 | 200 | 1 | ok |
| auto-default-5 | GET | `/api/v1/talents?page=1&page_size=20` | PASS | 200 | 200 | 4 | ok |
| auto-default-6 | GET | `/api/v1/jobs?page=1&page_size=5` | PASS | 200 | 200 | 1 | ok |
| auto-default-7 | GET | `/api/v1/conversations/unread-count` | PASS | 200 | 200 | 1 | ok |
| auto-default-8 | GET | `/api/v1/jobs?page=1&page_size=10&status=open` | PASS | 200 | 200 | 1 | ok |
| auto-default-9 | GET | `/api/v1/applications?page=1&page_size=1000&talent_id=me` | PASS | 200 | 200 | 2 | ok |
| auto-default-10 | GET | `/api/v1/conversations` | PASS | 200 | 200 | 1 | ok |
| auto-default-11 | GET | `/api/v1/jobs/1` | PASS | 200 | 200 | 1 | ok |
| auto-default-12 | GET | `/api/v1/applications?page=1&page_size=20` | PASS | 200 | 200 | 7 | ok |
| auto-default-13 | GET | `/api/v1/resumes/online` | PASS | 200 | 200 | 0 | ok |
| auto-default-14 | GET | `/api/v1/resumes?page_size=10` | PASS | 200 | 200 | 5 | ok |
| auto-default-15 | GET | `/api/v1/jobs?page_size=100&status=open` | PASS | 200 | 200 | 1 | ok |
| auto-default-16 | GET | `/api/v1/talents?page=1&page_size=10` | PASS | 200 | 200 | 2 | ok |
| auto-default-17 | GET | `/api/v1/jobs?page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |
| auto-default-18 | GET | `/api/v1/jobs/1/applications?page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |
| auto-default-19 | GET | `/api/v1/resumes?page=1&page_size=10` | PASS | 200 | 200 | 6 | ok |
| auto-default-20 | GET | `/api/v1/jobs?page=1&page_size=20` | PASS | 200 | 200 | 1 | ok |
| auto-default-21 | GET | `/api/v1/applications?page=1&page_size=50` | PASS | 200 | 200 | 9 | ok |
| auto-default-22 | GET | `/api/v1/interviews/1` | PASS | 200 | 200 | 1 | ok |
| auto-default-23 | GET | `/api/v1/interviews/1/feedback` | PASS | 200 | 200 | 0 | ok |
| auto-default-24 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | PASS | 200 | 200 | 1 | ok |
| auto-default-25 | GET | `/api/v1/logs?page=1&page_size=20` | PASS | 200 | 200 | 1 | ok |
| auto-default-26 | GET | `/api/v1/notices/1` | PASS | 200 | 200 | 1 | ok |
| auto-default-27 | GET | `/api/v1/notices?page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |

### 结果分析

| 指标 | 结果 | 说明 |
| --- | ---: | --- |
| 页面访问数 | 28 | 自动探索实际打开的前端页面数量 |
| 输入框填充 | 133/437 | 仅填写可见、可编辑且能识别字段含义的输入框 |
| 按钮点击 | 160/542 | 按 clickRiskLevels 配置执行，测试环境可允许 danger/empty 按钮 |
| API 资源数 | 10 | 按 /api/v1/{resource} 归类统计 |
| 回放失败 | 1 | 状态码、响应结构或请求错误不符合预期 |
| 回放跳过 | 0 | 多为非 BIGTEST 自有写请求，避免污染真实数据 |

### 失败接口复现命令

#### auto-default-1 POST /api/v1/login

```bash
curl -i -X 'POST' 'http://localhost:8080/api/v1/login' -H 'accept: application/json, text/plain, */*' -H 'accept-language: zh-CN,zh;q=0.9' -H 'content-type: application/json' -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/147.0.0.0 Safari/537.36' --data '{
  "username": "admin",
  "password": "<redacted>"
}'
```


### 错误现场

| 类型 | 数量 | 说明 |
| --- | ---: | --- |
| 页面失败 | 0 | 自动探索打开页面失败时保存截图和错误 |
| 按钮失败 | 19 | 页面按钮点击失败时保存截图和错误 |
| 控制台错误/警告 | 161 | 来源于 Playwright console 事件 |
| 网络失败请求 | 0 | 来源于 Playwright requestfailed 事件 |

| 类型 | URL/页面 | 错误 | 截图 |
| --- | --- | --- | --- |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |
| 按钮失败 | `http://localhost:5173/ai-evaluate` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(5)[22m [2m - locator resolved to <button disabled type="button" data-v-cf3bdcb2="" aria-disabled="true" class="el-button el-button--primary el-button--large is-disabled">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | - |

| 控制台类型 | 页面 | 内容 |
| --- | --- | --- |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". null at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". null at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". null at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". null at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". Proxy(Object) at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| error | `http://localhost:5173/dev-logs` | Failed to load resource: the server responded with a status of 500 (Internal Server Error) |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". null at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: custom validator check failed for prop "type". null at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |
| warning | `http://localhost:5173/portal/jobs` | Invalid prop: validation failed for prop "type". Expected one of ["primary", "success", "info", "warning", "danger"], got value "". null at <ElTag> at <PortalJobList> at <RouterView> at <PortalLayout> at <RouterView> at <App> [Object, Object, Object, Object, Object, Object] |

### 后续建议

- 优先处理回放失败和安全写失败项，查看对应 JSON 明细定位请求与响应。
