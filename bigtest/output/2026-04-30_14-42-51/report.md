# bigtest 真实流量自动化测试报告

- Run ID: `2026-04-30_14-42-51`
- 配置文件: `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/config.local.json`
- 输出目录: `output/2026-04-30_14-42-51`
- 自动探索: 已开启
- 安全写操作: 已开启，测试数据前缀 `BIGTEST_2026-04-30_14-42-51`

## 总览

| Profile | 页面 | 输入 | 按钮 | 原始请求 | 生成用例 | 安全写通过 | 安全写失败 | 回放通过 | 回放失败 | 回放跳过 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| auto-default | 28 | 407/69 | 542/50 | 90 | 35 | 1 | 0 | 34 | 0 | 1 |

## auto-default

### 自动发现覆盖

- 发现路由: 36
- 实际访问页面: 28
- 输入框: 407 个，已自动填写 69 个
- 按钮/链接: 542 个，安全点击 50 个
- 自动取消确认弹窗: 0
- API 资源数: 10
- 安全拦截/跳过: 341

### API 资源覆盖

| 资源 | 方法 | 用例数 |
| --- | --- | ---: |
| stats | GET | 1 |
| messages | GET | 3 |
| talents | GET | 2 |
| jobs | GET | 10 |
| conversations | GET | 2 |
| applications | GET | 3 |
| resumes | GET | 4 |
| interviews | GET | 2 |
| logs | GET | 2 |
| notices | GET | 5 |

### 被跳过的页面操作

| 风险 | 文案 | 页面 |
| --- | --- | --- |
| unknown | 首页 | `http://localhost:5173/portal` |
| unknown | 职位列表 | `http://localhost:5173/portal` |
| unknown | 企业招聘 | `http://localhost:5173/portal` |
| unknown | 消息 | `http://localhost:5173/portal` |
| unknown | a admin | `http://localhost:5173/portal` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs?keyword=%E5%90%8E%E7%AB%AF&city=` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs?keyword=%E5%90%8E%E7%AB%AF&city=` |
| unknown | 首页 | `http://localhost:5173/portal/jobs` |
| unknown | 职位列表 | `http://localhost:5173/portal/jobs` |
| unknown | 企业招聘 | `http://localhost:5173/portal/jobs` |
| unknown | 消息 | `http://localhost:5173/portal/jobs` |
| unknown | a admin | `http://localhost:5173/portal/jobs` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs` |
| unknown | 投递简历 | `http://localhost:5173/portal/jobs` |
| unknown | <empty> | `http://localhost:5173/portal/jobs` |
| unknown | <empty> | `http://localhost:5173/portal/jobs` |
| unknown | 首页 | `http://localhost:5173/portal/jobs/1` |
| unknown | 职位列表 | `http://localhost:5173/portal/jobs/1` |
| unknown | 企业招聘 | `http://localhost:5173/portal/jobs/1` |
| unknown | 消息 | `http://localhost:5173/portal/jobs/1` |
| unknown | a admin | `http://localhost:5173/portal/jobs/1` |
| unknown | 立即投递 | `http://localhost:5173/portal/jobs/1` |
### 安全写操作结果

| 场景 | 来源 | 状态 | 资源 | 步骤数 | 说明 |
| --- | --- | --- | --- | ---: | --- |
| notice-crud-owned-data | 配置 | PASS | notice | 4 | BIGTEST 数据创建、更新和清理完成 |

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
| auto-default-8 | GET | `/api/v1/jobs?keyword=%E5%90%8E%E7%AB%AF&page=1&page_size=10&status=open` | 1 | 是 |
| auto-default-9 | GET | `/api/v1/applications?page=1&page_size=1000&talent_id=me` | 4 | 是 |
| auto-default-10 | GET | `/api/v1/jobs?page=1&page_size=10&status=open` | 2 | 是 |
| auto-default-11 | GET | `/api/v1/jobs/1` | 2 | 是 |
| auto-default-12 | GET | `/api/v1/applications?page=1&page_size=20` | 1 | 是 |
| auto-default-13 | GET | `/api/v1/jobs/62` | 2 | 是 |
| auto-default-14 | GET | `/api/v1/resumes/online` | 1 | 是 |
| auto-default-15 | GET | `/api/v1/resumes?page_size=10` | 1 | 是 |
| auto-default-16 | GET | `/api/v1/conversations` | 4 | 是 |
| auto-default-17 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | 2 | 是 |
| auto-default-18 | GET | `/api/v1/talents?page=1&page_size=10` | 1 | 是 |
| auto-default-19 | GET | `/api/v1/jobs?page=1&page_size=10` | 4 | 是 |
| auto-default-20 | GET | `/api/v1/jobs/1/applications?page=1&page_size=10` | 1 | 是 |
| auto-default-21 | GET | `/api/v1/jobs?page_size=100&status=open` | 1 | 是 |
| auto-default-22 | GET | `/api/v1/resumes?page=1&page_size=10` | 2 | 是 |
| auto-default-23 | GET | `/api/v1/resumes?page=1&page_size=10&search=%E6%B5%8B%E8%AF%95` | 1 | 是 |
| auto-default-24 | GET | `/api/v1/jobs?page=1&page_size=20` | 1 | 是 |
| auto-default-25 | GET | `/api/v1/applications?page=1&page_size=50` | 1 | 是 |
| auto-default-26 | GET | `/api/v1/interviews/1/feedback` | 1 | 是 |
| auto-default-27 | GET | `/api/v1/interviews/1` | 1 | 是 |
| auto-default-28 | GET | `/api/v1/jobs?page=1&page_size=10&search=%E5%90%8E%E7%AB%AF` | 1 | 是 |
| auto-default-29 | GET | `/api/v1/logs?page=1&page_size=20` | 2 | 是 |
| auto-default-30 | GET | `/api/v1/logs?end_time=2026-05-01&keyword=%E6%B5%8B%E8%AF%95&page=1&page_size=20&start_time=2026-05-01` | 1 | 是 |
| auto-default-31 | GET | `/api/v1/notices/1` | 2 | 是 |
| auto-default-32 | GET | `/api/v1/notices?page=1&page_size=10` | 3 | 是 |
| auto-default-33 | GET | `/api/v1/notices/20` | 1 | 是 |
| auto-default-34 | GET | `/api/v1/notices?keyword=%E6%A0%A1%E5%9B%AD%E6%8B%9B%E8%81%98&page=1&page_size=10` | 1 | 是 |
| auto-default-35 | GET | `/api/v1/notices/11` | 1 | 是 |

### 回放结果

| 用例ID | 状态 | 期望 | 实际 | 说明 |
| --- | --- | ---: | ---: | --- |
| auto-default-1 | SKIP | - | - | POST body is not BIGTEST-owned |
| auto-default-2 | PASS | 200 | 200 | ok |
| auto-default-3 | PASS | 200 | 200 | ok |
| auto-default-4 | PASS | 200 | 200 | ok |
| auto-default-5 | PASS | 200 | 200 | ok |
| auto-default-6 | PASS | 200 | 200 | ok |
| auto-default-7 | PASS | 200 | 200 | ok |
| auto-default-8 | PASS | 200 | 200 | ok |
| auto-default-9 | PASS | 200 | 200 | ok |
| auto-default-10 | PASS | 200 | 200 | ok |
| auto-default-11 | PASS | 200 | 200 | ok |
| auto-default-12 | PASS | 200 | 200 | ok |
| auto-default-13 | PASS | 200 | 200 | ok |
| auto-default-14 | PASS | 200 | 200 | ok |
| auto-default-15 | PASS | 200 | 200 | ok |
| auto-default-16 | PASS | 200 | 200 | ok |
| auto-default-17 | PASS | 200 | 200 | ok |
| auto-default-18 | PASS | 200 | 200 | ok |
| auto-default-19 | PASS | 200 | 200 | ok |
| auto-default-20 | PASS | 200 | 200 | ok |
| auto-default-21 | PASS | 200 | 200 | ok |
| auto-default-22 | PASS | 200 | 200 | ok |
| auto-default-23 | PASS | 200 | 200 | ok |
| auto-default-24 | PASS | 200 | 200 | ok |
| auto-default-25 | PASS | 200 | 200 | ok |
| auto-default-26 | PASS | 200 | 200 | ok |
| auto-default-27 | PASS | 200 | 200 | ok |
| auto-default-28 | PASS | 200 | 200 | ok |
| auto-default-29 | PASS | 200 | 200 | ok |
| auto-default-30 | PASS | 200 | 200 | ok |
| auto-default-31 | PASS | 200 | 200 | ok |
| auto-default-32 | PASS | 200 | 200 | ok |
| auto-default-33 | PASS | 200 | 200 | ok |
| auto-default-34 | PASS | 200 | 200 | ok |
| auto-default-35 | PASS | 200 | 200 | ok |
