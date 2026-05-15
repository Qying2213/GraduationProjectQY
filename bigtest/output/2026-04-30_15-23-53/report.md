# bigtest 真实流量自动化测试报告

- Run ID: `2026-04-30_15-23-53`
- 配置文件: `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/config.local.json`
- 输出目录: `output/2026-04-30_15-23-53`
- 自动探索: 已开启
- 安全写操作: 已开启，测试数据前缀 `BIGTEST_2026-04-30_15-23-53`

## 测试结论

本次测试共访问 28 个页面，捕获 104 条真实 API 请求，生成 30 条接口回放用例。
接口回放通过 29 条，失败 0 条，跳过 1 条；安全写操作通过 4 个场景，失败 1 个场景。

**结论：本次自动化测试存在失败项，请优先查看“回放结果”和“安全写操作步骤”。**

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
| AI 评估策略 | 明确排除 `/api/v1/ai`、`/api/v1/evaluations`、`/ai/evaluate`、`/evaluate` 等路径，不触发真实 AI 评估、上传或外部模型调用 |

## 总览

| Profile | 页面 | 输入 | 按钮 | 原始请求 | 生成用例 | 安全写通过 | 安全写失败 | 回放通过 | 回放失败 | 回放跳过 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| auto-default | 28 | 693/124 | 542/133 | 104 | 30 | 4 | 1 | 29 | 0 | 1 |

## 模块覆盖统计

| 模块 | 页面访问 | API 用例 | GET | POST | PUT | PATCH | DELETE | 未覆盖接口 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 公告 | 3 | 8 | 5 | 1 | 1 | 0 | 1 | 0 |
| 简历 | 2 | 4 | 4 | 0 | 0 | 0 | 0 | 13 |
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

### 被安全策略跳过的关键功能

| Profile | 功能 | 页面 | 原因 |
| --- | --- | --- | --- |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/dashboard` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/dashboard` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/dashboard` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/talents` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/talents` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/talents/1` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/talents/1` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin人才详情人才详情页面开发中...后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/talents/1` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/jobs` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/jobs` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/jobs/1` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/jobs/1` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary 上传简历 上传简历 el-button el-button--primary | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | Sort by 上传时间 caret-wrapper 上传时间 Sort by 上传时间 caret-wrapper | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/recommend` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/recommend` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/kanban` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/kanban` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/calendar` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/calendar` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin面试日历管理和查看所有面试安排 导出 安排面试 2今日面试5本周面试4待进行1已完成今天2026年4月月视图周视图列表周日周一周二周三周四周五周六29303112345678910111213141516171819202122232425262728293011:00钱七109:30张三14:00李四210:00王五315:30赵六456789今日安排09:30张三高级前端工程师 刘经理初试14:00李四后端开发工程师 陈总监复试面试详情后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/calendar` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/interviews/1` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/interviews/1` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin面试日历管理和查看所有面试安排 导出 安排面试 2今日面试5本周面试4待进行1已完成今天2026年4月月视图周视图列表周日周一周二周三周四周五周六29303112345678910111213141516171819202122232425262728293011:00钱七109:30张三14:00李四210:00王五315:30赵六456789今日安排09:30张三高级前端工程师 刘经理初试14:00李四后端开发工程师 陈总监复试面试详情后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/calendar` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/messages` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/messages` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger el-button--small is-text 删除 删除 el-button el-button--danger el-button--small is-text | `http://localhost:5173/messages` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger el-button--small is-text 删除 删除 el-button el-button--danger el-button--small is-text | `http://localhost:5173/messages` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger el-button--small is-text 删除 删除 el-button el-button--danger el-button--small is-text | `http://localhost:5173/messages` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin消息中心及时掌握系统通知、面试邀约等重要信息 全部已读 10全部消息0未读消息3面试邀约0系统通知全部消息未读消息面试邀约简历投递系统通知系统公告 全选 38面试安排提醒您为候选人李娜创建的面试安排将在今天15:00开始，请提前准备。 系统 2026/4/2 删除 24面试邀请 - Go开发工程师您好，张三！ 您已被邀请参加【Go开发工程师】职位的初试。 面试详情： - 面试类型：初试 - 面试方式：现场面试 - 面试时间：2026-04-01 14:00 - 时长：60分钟 - 面试官：李四 - 地点/链接：会议室A 请准时参加面试，祝您面试顺利！ 系统 2026/3/31 删除 25面试改期通知 - Go开发您好，张三！ 您的【Go开发】职位面试时间已调整。 原时间：2026-12-25 14:00 新时间：2026-04-02 15:30 改期原因：面试官出差 请按新时间准时参加面试，祝您面试顺利！ 系统 2026/3/31 删除 Total 1210/page 1 2消息详情后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/messages` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/chat` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/chat` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin即时聊天与候选人实时沟通，高效完成招聘流程 刷新 0会话列表1 个会话张张伟03/27你好选择一个会话开始聊天从左侧列表中选择候选人，开始实时沟通后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/chat` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/profile` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/profile` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin个人中心管理您的账户信息和偏好设置更换头像admin管理员12发布职位86收藏人才24面试安排基本信息账户安全通知设置操作记录基本信息 编辑信息 姓名性别男女手机号邮箱公司职位所在城市北京入职日期个人简介账户安全登录密码定期更换密码可以保护账户安全修改密码手机绑定已绑定：13800000001更换手机邮箱绑定已绑定：admin@company.com更换邮箱微信绑定未绑定立即绑定最近登录记录登录时间IP地址登录地点设备状态2024-01-10 09:30:22192.168.1.100北京Chrome / Windows成功2024-01-09 18:45:11192.168.1.100北京Safari / macOS成功2024-01-09 08:20:3310.0.0.55上海Mobile / iOS成功2024-01-08 22:15:44203.156.78.90深圳Chrome / Windows失败2024-01-08 14:30:00192.168.1.100北京Chrome / Windows成功通知设置消息通知系统消息接收系统公告、维护通知等面试邀约接收新的面试邀请通知简历投递有求职者投递简历时通知人才推荐接收智能推荐的匹配人才通知方式站内消息在消息中心接收通知邮件通知通过邮件接收重要通知短信通知通过短信接收紧急通知保存设置操作记录2024-01-10 10:30发布了新职位发布了"高级前端工程师"职位，薪资范围30-50K2024-01-10 09:15收藏了人才将"李明 - 资深前端工程师"添加到收藏夹2024-01-09 16:45安排了面试为"王强 - 后端工程师"安排了技术面试，时间：1月15日 14:002024-01-09 11:20查看了简历查看了"张伟"的简历详情2024-01-08 15:30修改了职位信息更新了"产品经理"职位的薪资范围和工作要求加载更多后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/profile` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/roles` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/roles` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin权限管理管理系统角色和权限配置 新建角色 当前角色超级管理员拥有系统所有权限超级管理员角色列表共 5 个角色系统角色超级管理员拥有系统所有权限 35 项权限 系统角色HR主管负责招聘流程管理和人才库管理 27 项权限 系统角色招聘专员负责日常招聘工作 17 项权限 系统角色面试官参与面试评估 8 项权限 系统角色只读用户只能查看数据 7 项权限 权限总览超级管理员仪表板1/1查看仪表板人才管理5/5查看人才创建人才编辑人才删除人才导出人才职位管理5/5查看职位发布职位编辑职位删除职位导出职位简历管理5/5查看简历上传简历编辑简历删除简历导出简历招聘看板2/2查看看板编辑看板面试日历4/4查看日历创建面试编辑面试删除面试消息中心2/2查看消息发送消息智能推荐2/2查看推荐使用推荐用户管理4/4查看用户创建用户编辑用户删除用户角色管理4/4查看角色创建角色编辑角色删除角色系统设置1/1系统设置后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/roles` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/reports` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/reports` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin数据报表全面的招聘数据分析与可视化至 导出报表 12.5% 1,256简历投递量较上期 增长 8.3% 328面试安排数较上期 增长 3.2% 45录用人数较上期 下降 15.6% 18天平均招聘周期较上期 下降招聘漏斗分析本周本月本季简历筛选率45%↑ 5%面试通过率62%↑ 3%录用转化率28%↓ 2%职位热度排行查看全部招聘趋势分析简历投递面试安排录用人数渠道效果分析官网投递456人36%猎聘网312人25%BOSS直聘234人19%内部推荐156人12%其他渠道98人8%部门招聘进度部门目标已录用完成率技术部201575%产品部8675%设计部55100%市场部10440%运营部6350%面试官效率排行1陈陈总监技术部45面试数68%通过率4.5平均评分2刘刘经理产品部38面试数72%通过率4.3平均评分3周周主管技术部32面试数65%通过率4.2平均评分4王王总监设计部28面试数78%通过率4.6平均评分5HHR小李人力资源56面试数82%通过率4.4平均评分后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/reports` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/settings` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/settings` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin个人中心管理您的账户信息和偏好设置更换头像admin管理员12发布职位86收藏人才24面试安排基本信息账户安全通知设置操作记录基本信息 编辑信息 姓名性别男女手机号邮箱公司职位所在城市北京入职日期个人简介账户安全登录密码定期更换密码可以保护账户安全修改密码手机绑定已绑定：13800000001更换手机邮箱绑定已绑定：admin@company.com更换邮箱微信绑定未绑定立即绑定最近登录记录登录时间IP地址登录地点设备状态2024-01-10 09:30:22192.168.1.100北京Chrome / Windows成功2024-01-09 18:45:11192.168.1.100北京Safari / macOS成功2024-01-09 08:20:3310.0.0.55上海Mobile / iOS成功2024-01-08 22:15:44203.156.78.90深圳Chrome / Windows失败2024-01-08 14:30:00192.168.1.100北京Chrome / Windows成功通知设置消息通知系统消息接收系统公告、维护通知等面试邀约接收新的面试邀请通知简历投递有求职者投递简历时通知人才推荐接收智能推荐的匹配人才通知方式站内消息在消息中心接收通知邮件通知通过邮件接收重要通知短信通知通过短信接收紧急通知保存设置操作记录2024-01-10 10:30发布了新职位发布了"高级前端工程师"职位，薪资范围30-50K2024-01-10 09:15收藏了人才将"李明 - 资深前端工程师"添加到收藏夹2024-01-09 16:45安排了面试为"王强 - 后端工程师"安排了技术面试，时间：1月15日 14:002024-01-09 11:20查看了简历查看了"张伟"的简历详情2024-01-08 15:30修改了职位信息更新了"产品经理"职位的薪资范围和工作要求加载更多后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/profile` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/logs` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/logs` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/notices/1/edit` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/notices/1/edit` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/notices` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/notices` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 编辑删除 编辑 el-button el-button--primary is-link | `http://localhost:5173/notices` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger is-link 编辑删除 删除 el-button el-button--danger is-link | `http://localhost:5173/notices` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary is-link 编辑删除 编辑 el-button el-button--primary is-link | `http://localhost:5173/notices` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--danger is-link 编辑删除 删除 el-button el-button--danger is-link | `http://localhost:5173/notices` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin公告管理练习一个完整的 Vue3 页面状态置顶优先级搜索新建公告标题状态优先级置顶创建时间操作校园招聘宣讲会物料准备已发布普通-2026-03-25T14:52:11.414267+08:00编辑删除2026校园招聘正式启动已发布普通-2026-03-13T11:38:29.571787+08:00编辑删除Total 210/page 1 Go to新建公告标题内容状态草稿优先级普通置顶取消保存后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/notices` | danger 策略拦截，避免误操作真实数据或外部依赖 |
| auto-default | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/notices/1` | danger 策略拦截，避免误操作真实数据或外部依赖 |

## auto-default

### 自动发现覆盖

- 发现路由: 36
- 实际访问页面: 28
- 输入框: 693 个，已自动填写 124 个
- 按钮/链接: 542 个，安全点击 133 个
- 自动取消确认弹窗: 0
- API 资源数: 10
- 安全拦截/跳过: 82

### API 资源覆盖

| 资源 | 方法 | 用例数 |
| --- | --- | ---: |
| stats | GET | 1 |
| messages | DELETE, GET, POST | 6 |
| talents | DELETE, GET, POST, PUT | 6 |
| jobs | DELETE, GET, POST, PUT | 11 |
| conversations | GET | 2 |
| applications | GET | 3 |
| resumes | GET | 4 |
| interviews | DELETE, GET, POST, PUT | 8 |
| logs | GET | 1 |
| notices | DELETE, GET, POST, PUT | 8 |

### 被跳过的页面操作

| 风险 | 文案 | 页面 |
| --- | --- | --- |
| danger | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/dashboard` |
| danger | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/dashboard` |
| danger | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/dashboard` |
| danger | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/talents` |
| danger | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/talents` |
| danger | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/talents/1` |
| danger | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/talents/1` |
| danger | dev-logs-entry 人才运营平台仪表板人才管理职位管理简历管理智能推荐AI智能评估AI处理流程评估结果招聘看板面试日历消息中心即时聊天公告管理数据报表系统管理权限管理操作日志系统设置 投递简历 AI评估 0admin人才详情人才详情页面开发中...后台运行终端日志无需登录即可查看 后台运行终端日志 无需登录即可查看 dev-logs-entry | `http://localhost:5173/talents/1` |
| danger | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/jobs` |
| danger | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/jobs` |
| danger | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/jobs/1` |
| danger | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/jobs/1` |
| danger | el-button el-button--primary el-button--small is-plain quick-btn 投递简历 AI评估 0admin 投递简历 el-button el-button--primary el-button--small is-plain quick-btn | `http://localhost:5173/resumes` |
| danger | el-button el-button--success el-button--small is-plain quick-btn 投递简历 AI评估 0admin AI评估 el-button el-button--success el-button--small is-plain quick-btn | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary 上传简历 上传简历 el-button el-button--primary | `http://localhost:5173/resumes` |
| danger | Sort by 上传时间 caret-wrapper 上传时间 Sort by 上传时间 caret-wrapper | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--danger is-link 预览 解析 删除 删除 el-button el-button--danger is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary is-link 预览 解析 删除 预览 el-button el-button--primary is-link | `http://localhost:5173/resumes` |
| danger | el-button el-button--primary is-link 预览 解析 删除 解析 el-button el-button--primary is-link | `http://localhost:5173/resumes` |
### 安全写操作结果

本轮安全写操作已执行 5 次 `DELETE` 删除请求，删除目标均为本轮创建并登记的 `BIGTEST_` 测试数据。

| 场景 | 来源 | 状态 | 资源 | 步骤数 | 说明 |
| --- | --- | --- | --- | ---: | --- |
| job-crud-owned-data | 配置 | PASS | jobs | 4 | BIGTEST 数据创建、更新、删除完成 |
| talent-crud-owned-data | 配置 | FAIL | talents | 4 | BIGTEST 数据创建、更新、删除完成 |
| message-create-query-delete-owned-data | 配置 | PASS | messages | 3 | BIGTEST 数据创建、更新、删除完成 |
| interview-crud-owned-data | 配置 | PASS | interviews | 6 | BIGTEST 数据创建、更新、删除完成 |
| notice-crud-owned-data | 配置 | PASS | notice | 4 | BIGTEST 数据创建、更新、删除完成 |

### 安全写操作步骤

| 场景 | 步骤 | 方法 | 路径 | 期望 | 实际 | 状态 |
| --- | --- | --- | --- | --- | ---: | --- |
| job-crud-owned-data | create 创建测试数据 | POST | `/api/v1/jobs` | 200, 201 | 201 | PASS |
| job-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/jobs/63` | 200 | 200 | PASS |
| job-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/jobs/63` | 200 | 200 | PASS |
| job-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/jobs/63` | 200, 204 | 200 | PASS |
| talent-crud-owned-data | create 创建测试数据 | POST | `/api/v1/talents` | 200, 201 | 201 | PASS |
| talent-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/talents/116` | 200 | 200 | PASS |
| talent-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/talents/116` | 200 | 500 | FAIL |
| talent-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/talents/116` | 200, 204 | 200 | PASS |
| message-create-query-delete-owned-data | create 创建测试数据 | POST | `/api/v1/messages` | 200, 201 | 201 | PASS |
| message-create-query-delete-owned-data | verify 验证测试数据 | GET | `/api/v1/messages` | 200 | 200 | PASS |
| message-create-query-delete-owned-data | delete 删除测试数据 | DELETE | `/api/v1/messages/40` | 200, 204 | 200 | PASS |
| interview-crud-owned-data | create 创建测试数据 | POST | `/api/v1/interviews` | 200, 201 | 201 | PASS |
| interview-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/interviews/113` | 200 | 200 | PASS |
| interview-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/interviews/113` | 200 | 200 | PASS |
| interview-crud-owned-data | reschedule 改期测试面试 | POST | `/api/v1/interviews/113/reschedule` | 200 | 200 | PASS |
| interview-crud-owned-data | complete 完成测试面试 | POST | `/api/v1/interviews/113/complete` | 200 | 200 | PASS |
| interview-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/interviews/113` | 200, 204 | 200 | PASS |
| notice-crud-owned-data | create 创建测试数据 | POST | `/api/v1/notices` | 200, 201 | 200 | PASS |
| notice-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/notices/30` | 200 | 200 | PASS |
| notice-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/notices/30` | 200 | 200 | PASS |
| notice-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/notices/30` | 200, 204 | 200 | PASS |

### 去重后生成的用例

| 用例ID | 方法 | 路径 | 归并次数 | 需鉴权 |
| --- | --- | --- | ---: | --- |
| auto-default-1 | POST | `/api/v1/login` | 1 | 否 |
| auto-default-2 | GET | `/api/v1/stats/dashboard` | 3 | 是 |
| auto-default-3 | GET | `/api/v1/messages/unread-count` | 24 | 是 |
| auto-default-4 | GET | `/api/v1/messages?page=1&page_size=5` | 3 | 是 |
| auto-default-5 | GET | `/api/v1/talents?page=1&page_size=20` | 3 | 是 |
| auto-default-6 | GET | `/api/v1/jobs?page=1&page_size=5` | 3 | 否 |
| auto-default-7 | GET | `/api/v1/conversations/unread-count` | 7 | 是 |
| auto-default-8 | GET | `/api/v1/jobs?page=1&page_size=10&status=open` | 8 | 是 |
| auto-default-9 | GET | `/api/v1/applications?page=1&page_size=1000&talent_id=me` | 9 | 是 |
| auto-default-10 | GET | `/api/v1/conversations` | 17 | 是 |
| auto-default-11 | GET | `/api/v1/jobs/1` | 2 | 是 |
| auto-default-12 | GET | `/api/v1/applications?page=1&page_size=20` | 1 | 是 |
| auto-default-13 | GET | `/api/v1/resumes/online` | 1 | 是 |
| auto-default-14 | GET | `/api/v1/resumes?page_size=10` | 1 | 是 |
| auto-default-15 | GET | `/api/v1/talents?page=1&page_size=10` | 1 | 是 |
| auto-default-16 | GET | `/api/v1/jobs?page=1&page_size=10` | 2 | 是 |
| auto-default-17 | GET | `/api/v1/jobs/1/applications?page=1&page_size=10` | 1 | 是 |
| auto-default-18 | GET | `/api/v1/jobs?page_size=100&status=open` | 1 | 是 |
| auto-default-19 | GET | `/api/v1/resumes?page=1&page_size=10` | 2 | 是 |
| auto-default-20 | GET | `/api/v1/resumes?page=1&page_size=10&search=%E6%B5%8B%E8%AF%95` | 1 | 是 |
| auto-default-21 | GET | `/api/v1/jobs?page=1&page_size=20` | 1 | 是 |
| auto-default-22 | GET | `/api/v1/applications?page=1&page_size=50` | 1 | 是 |
| auto-default-23 | GET | `/api/v1/interviews/1` | 1 | 是 |
| auto-default-24 | GET | `/api/v1/interviews/1/feedback` | 1 | 是 |
| auto-default-25 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | 1 | 是 |
| auto-default-26 | GET | `/api/v1/logs?page=1&page_size=20` | 1 | 是 |
| auto-default-27 | GET | `/api/v1/notices/1` | 2 | 是 |
| auto-default-28 | GET | `/api/v1/notices?page=1&page_size=10` | 3 | 是 |
| auto-default-29 | GET | `/api/v1/notices/20` | 1 | 是 |
| auto-default-30 | GET | `/api/v1/notices?keyword=%E6%A0%A1%E5%9B%AD%E6%8B%9B%E8%81%98&page=1&page_size=10` | 1 | 是 |

### 回放结果

| 用例ID | 方法 | 路径 | 状态 | 期望 | 实际 | 耗时(ms) | 说明 |
| --- | --- | --- | --- | ---: | ---: | ---: | --- |
| auto-default-1 | POST | `/api/v1/login` | SKIP | - | - | - | POST body is not BIGTEST-owned |
| auto-default-2 | GET | `/api/v1/stats/dashboard` | PASS | 200 | 200 | 0 | ok |
| auto-default-3 | GET | `/api/v1/messages/unread-count` | PASS | 200 | 200 | 1 | ok |
| auto-default-4 | GET | `/api/v1/messages?page=1&page_size=5` | PASS | 200 | 200 | 1 | ok |
| auto-default-5 | GET | `/api/v1/talents?page=1&page_size=20` | PASS | 200 | 200 | 4 | ok |
| auto-default-6 | GET | `/api/v1/jobs?page=1&page_size=5` | PASS | 200 | 200 | 1 | ok |
| auto-default-7 | GET | `/api/v1/conversations/unread-count` | PASS | 200 | 200 | 1 | ok |
| auto-default-8 | GET | `/api/v1/jobs?page=1&page_size=10&status=open` | PASS | 200 | 200 | 1 | ok |
| auto-default-9 | GET | `/api/v1/applications?page=1&page_size=1000&talent_id=me` | PASS | 200 | 200 | 3 | ok |
| auto-default-10 | GET | `/api/v1/conversations` | PASS | 200 | 200 | 1 | ok |
| auto-default-11 | GET | `/api/v1/jobs/1` | PASS | 200 | 200 | 0 | ok |
| auto-default-12 | GET | `/api/v1/applications?page=1&page_size=20` | PASS | 200 | 200 | 5 | ok |
| auto-default-13 | GET | `/api/v1/resumes/online` | PASS | 200 | 200 | 1 | ok |
| auto-default-14 | GET | `/api/v1/resumes?page_size=10` | PASS | 200 | 200 | 5 | ok |
| auto-default-15 | GET | `/api/v1/talents?page=1&page_size=10` | PASS | 200 | 200 | 2 | ok |
| auto-default-16 | GET | `/api/v1/jobs?page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |
| auto-default-17 | GET | `/api/v1/jobs/1/applications?page=1&page_size=10` | PASS | 200 | 200 | 2 | ok |
| auto-default-18 | GET | `/api/v1/jobs?page_size=100&status=open` | PASS | 200 | 200 | 1 | ok |
| auto-default-19 | GET | `/api/v1/resumes?page=1&page_size=10` | PASS | 200 | 200 | 5 | ok |
| auto-default-20 | GET | `/api/v1/resumes?page=1&page_size=10&search=%E6%B5%8B%E8%AF%95` | PASS | 200 | 200 | 2 | ok |
| auto-default-21 | GET | `/api/v1/jobs?page=1&page_size=20` | PASS | 200 | 200 | 1 | ok |
| auto-default-22 | GET | `/api/v1/applications?page=1&page_size=50` | PASS | 200 | 200 | 9 | ok |
| auto-default-23 | GET | `/api/v1/interviews/1` | PASS | 200 | 200 | 1 | ok |
| auto-default-24 | GET | `/api/v1/interviews/1/feedback` | PASS | 200 | 200 | 1 | ok |
| auto-default-25 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | PASS | 200 | 200 | 0 | ok |
| auto-default-26 | GET | `/api/v1/logs?page=1&page_size=20` | PASS | 200 | 200 | 1 | ok |
| auto-default-27 | GET | `/api/v1/notices/1` | PASS | 200 | 200 | 1 | ok |
| auto-default-28 | GET | `/api/v1/notices?page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |
| auto-default-29 | GET | `/api/v1/notices/20` | PASS | 200 | 200 | 0 | ok |
| auto-default-30 | GET | `/api/v1/notices?keyword=%E6%A0%A1%E5%9B%AD%E6%8B%9B%E8%81%98&page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |

### 结果分析

| 指标 | 结果 | 说明 |
| --- | ---: | --- |
| 页面访问数 | 28 | 自动探索实际打开的前端页面数量 |
| 输入框填充 | 124/693 | 仅填写可见、可编辑且能识别字段含义的输入框 |
| 按钮点击 | 133/542 | 已允许 safe、medium、unknown 按钮，danger/empty 仍按规则跳过 |
| API 资源数 | 10 | 按 /api/v1/{resource} 归类统计 |
| 回放失败 | 0 | 状态码、响应结构或请求错误不符合预期 |
| 回放跳过 | 1 | 多为非 BIGTEST 自有写请求，避免污染真实数据 |

### 跳过原因统计

| 风险类型 | 数量 | 说明 |
| --- | ---: | --- |
| danger | 81 | 高风险页面操作，如删除、上传、AI评估、真实投递等 |

### 错误现场

| 类型 | 数量 | 说明 |
| --- | ---: | --- |
| 页面失败 | 0 | 自动探索打开页面失败时保存截图和错误 |
| 按钮失败 | 90 | 页面按钮点击失败时保存截图和错误 |
| 控制台错误/警告 | 146 | 来源于 Playwright console 事件 |
| 网络失败请求 | 0 | 来源于 Playwright requestfailed 事件 |

| 类型 | URL/页面 | 错误 | 截图 |
| --- | --- | --- | --- |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(8)[22m [2m - locator resolved to <button type="button" data-v-075e9b23="" aria-disabled="false" class="el-button el-button--primary">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-form-item is-required asterisk-left el-form-item--label-top">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-form-item is-required asterisk-left el-form-item--label-top">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-form-item is-required asterisk-left el-form-item--label-top">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-1.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(9)[22m [2m - locator resolved to <button type="button" class="el-button" data-v-075e9b23="" aria-disabled="false">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-form-item is-required asterisk-left el-form-item--label-top">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-form-item is-required asterisk-left el-form-item--label-top">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-form-item is-required asterisk-left el-form-item--label-top">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-2.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(10)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-075e9b23="" id="el-id-3341-75" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-3341-76" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <input type="text" tabindex="0" autocomplete="off" placeholder="请输入地区" id="el-id-3341-169" class="el-input__inner"/> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <input type="text" tabindex="0" autocomplete="off" placeholder="请输入地区" id="el-id-3341-169" class="el-input__inner"/> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <input type="text" tabindex="0" autocomplete="off" placeholder="请输入地区" id="el-id-3341-169" class="el-input__inner"/> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-3.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(11)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-075e9b23="" id="el-id-3341-79" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-3341-80" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <input type="text" tabindex="0" autocomplete="off" id="el-id-3341-170" placeholder="如：20-30K" class="el-input__inner"/> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <input type="text" tabindex="0" autocomplete="off" id="el-id-3341-170" placeholder="如：20-30K" class="el-input__inner"/> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <input type="text" tabindex="0" autocomplete="off" id="el-id-3341-170" placeholder="如：20-30K" class="el-input__inner"/> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-4.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(12)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-075e9b23="" id="el-id-3341-83" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-3341-84" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-5.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(13)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-075e9b23="" id="el-id-3341-87" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-3341-88" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <textarea rows="4" tabindex="0" autocomplete="off" id="el-id-3341-173" placeholder="请输入个人简介" class="el-textarea__inner"></textarea> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <textarea rows="4" tabindex="0" autocomplete="off" id="el-id-3341-173" placeholder="请输入个人简介" class="el-textarea__inner"></textarea> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-6.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(14)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-075e9b23="" id="el-id-3341-91" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-3341-92" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <textarea rows="4" tabindex="0" autocomplete="off" id="el-id-3341-173" placeholder="请输入个人简介" class="el-textarea__inner"></textarea> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <textarea rows="4" tabindex="0" autocomplete="off" id="el-id-3341-173" placeholder="请输入个人简介" class="el-textarea__inner"></textarea> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-7.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(15)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-075e9b23="" id="el-id-3341-95" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-3341-96" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <strong>后台运行终端日志</strong> from <a href="/dev-logs" class="dev-logs-entry">…</a> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <strong>后台运行终端日志</strong> from <a href="/dev-logs" class="dev-logs-entry">…</a> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-8.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(16)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-075e9b23="" id="el-id-3341-99" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-3341-100" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <form data-v-075e9b23="" class="el-form el-form--default el-form--label-top">…</form> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div tabindex="-1" class="el-dialog">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <form data-v-075e9b23="" class="el-form el-form--default el-form--label-top">…</form> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-9.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(17)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-075e9b23="" id="el-id-3341-103" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-3341-104" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <form data-v-075e9b23="" class="el-form el-form--default el-form--label-top">…</form> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div tabindex="-1" class="el-dialog">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <form data-v-075e9b23="" class="el-form el-form--default el-form--label-top">…</form> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-10.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(18)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-075e9b23="" id="el-id-3341-107" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-3341-108" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-11.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(19)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-075e9b23="" id="el-id-3341-111" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-3341-112" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <form data-v-075e9b23="" class="el-form el-form--default el-form--label-top">…</form> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-075e9b23="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <form data-v-075e9b23="" class="el-form el-form--default el-form--label-top">…</form> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-12.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(20)[22m [2m - locator resolved to <button disabled type="button" class="btn-prev" aria-disabled="true" aria-label="Go to previous page">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is not enabled[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-13.png` |
| 按钮失败 | `http://localhost:5173/talents` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(21)[22m [2m - locator resolved to <button type="button" aria-disabled="false" class="btn-next is-last" aria-label="Go to next page">…</button>[22m [2m - attempting click action[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <a href="/dev-logs" class="dev-logs-entry">…</a> intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <a href="/dev-logs" class="dev-logs-entry">…</a> intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <a href="/dev-logs" class="dev-logs-entry">…</a> intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="新增人才" class="el-overlay-dialog" aria-describedby="el-id-3341-65">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-14.png` |
| 按钮失败 | `http://localhost:5173/jobs` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(6)[22m [2m - locator resolved to <button type="button" data-v-081df423="" aria-disabled="false" class="el-button el-button--primary">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="发布新职位" class="el-overlay-dialog" aria-describedby="el-id-4377-38">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-15.png` |
| 按钮失败 | `http://localhost:5173/jobs` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(7)[22m [2m - locator resolved to <button type="button" class="el-button" data-v-081df423="" aria-disabled="false">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="发布新职位" class="el-overlay-dialog" aria-describedby="el-id-4377-38">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-16.png` |
| 按钮失败 | `http://localhost:5173/jobs` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(8)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-081df423="" id="el-id-4377-45" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-4377-46" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-form-item is-required asterisk-left el-form-item--label-top">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="发布新职位" class="el-overlay-dialog" aria-describedby="el-id-4377-38">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-form-item is-required asterisk-left el-form-item--label-top">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-17.png` |
| 按钮失败 | `http://localhost:5173/jobs` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(9)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-081df423="" id="el-id-4377-50" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-4377-51" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="发布新职位" class="el-overlay-dialog" aria-describedby="el-id-4377-38">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="发布新职位" class="el-overlay-dialog" aria-describedby="el-id-4377-38">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="发布新职位" class="el-overlay-dialog" aria-describedby="el-id-4377-38">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-18.png` |
| 按钮失败 | `http://localhost:5173/jobs` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(10)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-081df423="" id="el-id-4377-55" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-4377-56" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <form data-v-081df423="" class="el-form el-form--default el-form--label-top">…</form> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div data-v-081df423="" class="el-col el-col-12 is-guttered">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m 2 × retrying click action[22m [2m - waiting 100ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="发布新职位" class="el-overlay-dialog" aria-describedby="el-id-4377-38">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m [2m - waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <form data-v-081df423="" class="el-form el-form--default el-form--label-top">…</form> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-19.png` |
| 按钮失败 | `http://localhost:5173/jobs` | locator.click: Timeout 1000ms exceeded. Call log: [2m - waiting for locator('button, a, .el-button, [role=\'button\']').nth(11)[22m [2m - locator resolved to <button tabindex="0" type="button" role="button" data-v-081df423="" id="el-id-4377-60" aria-haspopup="menu" aria-disabled="false" aria-expanded="false" aria-controls="el-id-4377-61" class="el-button is-text el-tooltip__trigger">…</button>[22m [2m - attempting click action[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="发布新职位" class="el-overlay-dialog" aria-describedby="el-id-4377-38">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 20ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="发布新职位" class="el-overlay-dialog" aria-describedby="el-id-4377-38">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 100ms[22m [2m 2 × waiting for element to be visible, enabled and stable[22m [2m - element is visible, enabled and stable[22m [2m - scrolling into view if needed[22m [2m - done scrolling[22m [2m - <div role="dialog" aria-modal="true" aria-label="发布新职位" class="el-overlay-dialog" aria-describedby="el-id-4377-38">…</div> from <div class="el-overlay el-modal-dialog">…</div> subtree intercepts pointer events[22m [2m - retrying click action[22m [2m - waiting 500ms[22m | `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/output/2026-04-30_15-23-53/auto-default/debug/button-failed-20.png` |

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
- 页面上的危险按钮没有直接点击；如需测试删除/投递等动作，请继续通过 `writeSafety.scenarios` 创建并操作 `BIGTEST_` 测试数据。
