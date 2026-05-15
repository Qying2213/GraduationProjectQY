# bigtest 真实流量自动化测试报告

- Run ID: `2026-04-30T03-09-26-200Z`
- 配置文件: `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/config.example.json`
- 输出目录: `output/2026-04-30T03-09-26-200Z`

## 总览

| Profile | 原始请求 | 生成用例 | 回放通过 | 回放失败 | 回放跳过 |
| --- | ---: | ---: | ---: | ---: | ---: |
| admin-smoke | 17 | 9 | 8 | 0 | 1 |

## admin-smoke

### 去重后生成的用例

| 用例ID | 方法 | 路径 | 归并次数 | 需鉴权 |
| --- | --- | --- | ---: | --- |
| admin-smoke-1 | POST | `/api/v1/login` | 1 | 否 |
| admin-smoke-2 | GET | `/api/v1/jobs?page=1&page_size=5` | 2 | 否 |
| admin-smoke-3 | GET | `/api/v1/stats/dashboard` | 2 | 是 |
| admin-smoke-4 | GET | `/api/v1/messages/unread-count` | 5 | 是 |
| admin-smoke-5 | GET | `/api/v1/messages?page=1&page_size=5` | 2 | 是 |
| admin-smoke-6 | GET | `/api/v1/talents?page=1&page_size=20` | 2 | 是 |
| admin-smoke-7 | GET | `/api/v1/jobs?page=1&page_size=10` | 1 | 是 |
| admin-smoke-8 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | 1 | 是 |
| admin-smoke-9 | GET | `/api/v1/notices?page=1&page_size=10` | 1 | 是 |

### 回放结果

| 用例ID | 状态 | 期望 | 实际 | 说明 |
| --- | --- | ---: | ---: | --- |
| admin-smoke-1 | SKIP | - | - | method POST not enabled |
| admin-smoke-2 | PASS | 200 | 200 | ok |
| admin-smoke-3 | PASS | 200 | 200 | ok |
| admin-smoke-4 | PASS | 200 | 200 | ok |
| admin-smoke-5 | PASS | 200 | 200 | ok |
| admin-smoke-6 | PASS | 200 | 200 | ok |
| admin-smoke-7 | PASS | 200 | 200 | ok |
| admin-smoke-8 | PASS | 200 | 200 | ok |
| admin-smoke-9 | PASS | 200 | 200 | ok |
