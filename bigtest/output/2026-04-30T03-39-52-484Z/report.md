# bigtest 真实流量自动化测试报告

- Run ID: `2026-04-30T03-39-52-484Z`
- 配置文件: `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/config.local.json`
- 输出目录: `output/2026-04-30T03-39-52-484Z`
- 自动探索: 已开启
- 安全写操作: 已开启，测试数据前缀 `BIGTEST_2026-04-30T03-39-52-484Z`

## 总览

| Profile | 原始请求 | 生成用例 | 安全写通过 | 安全写失败 | 回放通过 | 回放失败 | 回放跳过 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| admin-smoke | 111 | 36 | 1 | 0 | 35 | 0 | 1 |

## admin-smoke

### 安全写操作结果

| 场景 | 状态 | 资源 | 步骤数 | 说明 |
| --- | --- | --- | ---: | --- |
| notice-crud-owned-data | PASS | notice | 4 | BIGTEST 数据创建、更新和清理完成 |

### 去重后生成的用例

| 用例ID | 方法 | 路径 | 归并次数 | 需鉴权 |
| --- | --- | --- | ---: | --- |
| admin-smoke-1 | POST | `/api/v1/login` | 1 | 否 |
| admin-smoke-2 | GET | `/api/v1/stats/dashboard` | 4 | 是 |
| admin-smoke-3 | GET | `/api/v1/messages/unread-count` | 29 | 是 |
| admin-smoke-4 | GET | `/api/v1/messages?page=1&page_size=5` | 4 | 是 |
| admin-smoke-5 | GET | `/api/v1/talents?page=1&page_size=20` | 4 | 是 |
| admin-smoke-6 | GET | `/api/v1/jobs?page=1&page_size=5` | 4 | 否 |
| admin-smoke-7 | GET | `/api/v1/jobs?page=1&page_size=10` | 6 | 是 |
| admin-smoke-8 | GET | `/api/v1/jobs?page=1&page_size=10&search=%E5%90%8E%E7%AB%AF` | 2 | 是 |
| admin-smoke-9 | GET | `/api/v1/jobs/1` | 3 | 是 |
| admin-smoke-10 | GET | `/api/v1/jobs/1/applications?page=1&page_size=10` | 2 | 是 |
| admin-smoke-11 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | 3 | 是 |
| admin-smoke-12 | GET | `/api/v1/notices?page=1&page_size=10` | 4 | 是 |
| admin-smoke-13 | GET | `/api/v1/notices?keyword=%E6%A0%A1%E5%9B%AD%E6%8B%9B%E8%81%98&page=1&page_size=10` | 2 | 是 |
| admin-smoke-14 | GET | `/api/v1/notices/1` | 3 | 是 |
| admin-smoke-15 | GET | `/api/v1/conversations/unread-count` | 7 | 是 |
| admin-smoke-16 | GET | `/api/v1/jobs?keyword=%E5%90%8E%E7%AB%AF&page=1&page_size=10&status=open` | 1 | 是 |
| admin-smoke-17 | GET | `/api/v1/applications?page=1&page_size=1000&talent_id=me` | 4 | 是 |
| admin-smoke-18 | GET | `/api/v1/jobs?page=1&page_size=10&status=open` | 2 | 是 |
| admin-smoke-19 | GET | `/api/v1/applications?page=1&page_size=20` | 1 | 是 |
| admin-smoke-20 | GET | `/api/v1/jobs/62` | 2 | 是 |
| admin-smoke-21 | GET | `/api/v1/resumes/online` | 1 | 是 |
| admin-smoke-22 | GET | `/api/v1/resumes?page_size=10` | 1 | 是 |
| admin-smoke-23 | GET | `/api/v1/conversations` | 4 | 是 |
| admin-smoke-24 | GET | `/api/v1/talents?page=1&page_size=10` | 2 | 是 |
| admin-smoke-25 | GET | `/api/v1/talents?page=1&page_size=10&search=%E6%B5%8B%E8%AF%95` | 1 | 是 |
| admin-smoke-26 | GET | `/api/v1/jobs?page_size=100&status=open` | 1 | 是 |
| admin-smoke-27 | GET | `/api/v1/resumes?page=1&page_size=10` | 2 | 是 |
| admin-smoke-28 | GET | `/api/v1/resumes?page=1&page_size=10&search=%E6%B5%8B%E8%AF%95` | 1 | 是 |
| admin-smoke-29 | GET | `/api/v1/jobs?page=1&page_size=20` | 1 | 是 |
| admin-smoke-30 | GET | `/api/v1/applications?page=1&page_size=50` | 1 | 是 |
| admin-smoke-31 | GET | `/api/v1/interviews/1/feedback` | 1 | 是 |
| admin-smoke-32 | GET | `/api/v1/interviews/1` | 1 | 是 |
| admin-smoke-33 | GET | `/api/v1/logs?page=1&page_size=20` | 2 | 是 |
| admin-smoke-34 | GET | `/api/v1/evaluations?latest_only=true&page=1&page_size=10&sort_by=created_at&sort_order=desc` | 2 | 是 |
| admin-smoke-35 | GET | `/api/v1/evaluations/stats?latest_only=true` | 1 | 是 |
| admin-smoke-36 | GET | `/api/v1/evaluations?latest_only=true&page=1&page_size=10&search=%E6%B5%8B%E8%AF%95&sort_by=created_at&sort_order=desc` | 1 | 是 |

### 回放结果

| 用例ID | 状态 | 期望 | 实际 | 说明 |
| --- | --- | ---: | ---: | --- |
| admin-smoke-1 | SKIP | - | - | POST body is not BIGTEST-owned |
| admin-smoke-2 | PASS | 200 | 200 | ok |
| admin-smoke-3 | PASS | 200 | 200 | ok |
| admin-smoke-4 | PASS | 200 | 200 | ok |
| admin-smoke-5 | PASS | 200 | 200 | ok |
| admin-smoke-6 | PASS | 200 | 200 | ok |
| admin-smoke-7 | PASS | 200 | 200 | ok |
| admin-smoke-8 | PASS | 200 | 200 | ok |
| admin-smoke-9 | PASS | 200 | 200 | ok |
| admin-smoke-10 | PASS | 200 | 200 | ok |
| admin-smoke-11 | PASS | 200 | 200 | ok |
| admin-smoke-12 | PASS | 200 | 200 | ok |
| admin-smoke-13 | PASS | 200 | 200 | ok |
| admin-smoke-14 | PASS | 200 | 200 | ok |
| admin-smoke-15 | PASS | 200 | 200 | ok |
| admin-smoke-16 | PASS | 200 | 200 | ok |
| admin-smoke-17 | PASS | 200 | 200 | ok |
| admin-smoke-18 | PASS | 200 | 200 | ok |
| admin-smoke-19 | PASS | 200 | 200 | ok |
| admin-smoke-20 | PASS | 200 | 200 | ok |
| admin-smoke-21 | PASS | 200 | 200 | ok |
| admin-smoke-22 | PASS | 200 | 200 | ok |
| admin-smoke-23 | PASS | 200 | 200 | ok |
| admin-smoke-24 | PASS | 200 | 200 | ok |
| admin-smoke-25 | PASS | 200 | 200 | ok |
| admin-smoke-26 | PASS | 200 | 200 | ok |
| admin-smoke-27 | PASS | 200 | 200 | ok |
| admin-smoke-28 | PASS | 200 | 200 | ok |
| admin-smoke-29 | PASS | 200 | 200 | ok |
| admin-smoke-30 | PASS | 200 | 200 | ok |
| admin-smoke-31 | PASS | 200 | 200 | ok |
| admin-smoke-32 | PASS | 200 | 200 | ok |
| admin-smoke-33 | PASS | 200 | 200 | ok |
| admin-smoke-34 | PASS | 200 | 200 | ok |
| admin-smoke-35 | PASS | 200 | 200 | ok |
| admin-smoke-36 | PASS | 200 | 200 | ok |
