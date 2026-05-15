# bigtest 真实流量自动化测试报告

- Run ID: `2026-04-30_15-11-30`
- 配置文件: `/Users/qinyang/Desktop/GraduationProjectQY/bigtest/config.local.json`
- 输出目录: `output/2026-04-30_15-11-30`
- 自动探索: 已开启
- 安全写操作: 已开启，测试数据前缀 `BIGTEST_2026-04-30_15-11-30`

## 测试结论

本次测试共访问 28 个页面，捕获 105 条真实 API 请求，生成 31 条接口回放用例。
接口回放通过 30 条，失败 0 条，跳过 1 条；安全写操作通过 1 个场景，失败 0 个场景。

**结论：本次自动化测试未发现接口回放失败或安全写操作失败。**

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

## 总览

| Profile | 页面 | 输入 | 按钮 | 原始请求 | 生成用例 | 安全写通过 | 安全写失败 | 回放通过 | 回放失败 | 回放跳过 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| auto-default | 28 | 551/98 | 542/103 | 105 | 31 | 1 | 0 | 30 | 0 | 1 |

## auto-default

### 自动发现覆盖

- 发现路由: 36
- 实际访问页面: 28
- 输入框: 551 个，已自动填写 98 个
- 按钮/链接: 542 个，安全点击 103 个
- 自动取消确认弹窗: 0
- API 资源数: 10
- 安全拦截/跳过: 182

### API 资源覆盖

| 资源 | 方法 | 用例数 |
| --- | --- | ---: |
| stats | GET | 1 |
| messages | GET | 3 |
| talents | GET | 2 |
| jobs | GET | 8 |
| conversations | GET | 2 |
| applications | GET | 3 |
| resumes | GET | 4 |
| interviews | GET | 2 |
| logs | GET | 1 |
| notices | DELETE, GET, POST, PUT | 8 |

### 被跳过的页面操作

| 风险 | 文案 | 页面 |
| --- | --- | --- |
| danger | AI评估 | `http://localhost:5173/dashboard` |
| danger | 投递简历 | `http://localhost:5173/dashboard` |
| danger | AI评估 | `http://localhost:5173/dashboard` |
| empty | <empty> | `http://localhost:5173/dashboard` |
| danger | 投递简历 | `http://localhost:5173/talents` |
| danger | AI评估 | `http://localhost:5173/talents` |
| empty | <empty> | `http://localhost:5173/talents` |
| empty | <empty> | `http://localhost:5173/talents` |
| empty | <empty> | `http://localhost:5173/talents` |
| empty | <empty> | `http://localhost:5173/talents` |
| danger | 删除 | `http://localhost:5173/talents` |
| danger | 删除 | `http://localhost:5173/talents` |
| danger | 删除 | `http://localhost:5173/talents` |
| danger | 删除 | `http://localhost:5173/talents` |
| danger | 删除 | `http://localhost:5173/talents` |
| danger | 删除 | `http://localhost:5173/talents` |
| danger | 删除 | `http://localhost:5173/talents` |
| danger | 投递简历 | `http://localhost:5173/talents/1` |
| danger | AI评估 | `http://localhost:5173/talents/1` |
| empty | <empty> | `http://localhost:5173/talents/1` |
| danger | 投递简历 | `http://localhost:5173/jobs` |
| danger | AI评估 | `http://localhost:5173/jobs` |
| empty | <empty> | `http://localhost:5173/jobs` |
| empty | <empty> | `http://localhost:5173/jobs` |
| empty | <empty> | `http://localhost:5173/jobs` |
| empty | <empty> | `http://localhost:5173/jobs` |
| empty | <empty> | `http://localhost:5173/jobs` |
| empty | <empty> | `http://localhost:5173/jobs` |
| empty | <empty> | `http://localhost:5173/jobs` |
| empty | <empty> | `http://localhost:5173/jobs` |
### 安全写操作结果

本轮安全写操作已执行 1 次 `DELETE` 删除请求，删除目标均为本轮创建并登记的 `BIGTEST_` 测试数据。

| 场景 | 来源 | 状态 | 资源 | 步骤数 | 说明 |
| --- | --- | --- | --- | ---: | --- |
| notice-crud-owned-data | 配置 | PASS | notice | 4 | BIGTEST 数据创建、更新、删除完成 |

### 安全写操作步骤

| 场景 | 步骤 | 方法 | 路径 | 期望 | 实际 | 状态 |
| --- | --- | --- | --- | --- | ---: | --- |
| notice-crud-owned-data | create 创建测试数据 | POST | `/api/v1/notices` | 200, 201 | 200 | PASS |
| notice-crud-owned-data | verify 验证测试数据 | GET | `/api/v1/notices/29` | 200 | 200 | PASS |
| notice-crud-owned-data | update 更新测试数据 | PUT | `/api/v1/notices/29` | 200 | 200 | PASS |
| notice-crud-owned-data | delete 删除测试数据 | DELETE | `/api/v1/notices/29` | 200, 204 | 200 | PASS |

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
| auto-default-21 | GET | `/api/v1/jobs/62` | 1 | 是 |
| auto-default-22 | GET | `/api/v1/jobs?page=1&page_size=20` | 1 | 是 |
| auto-default-23 | GET | `/api/v1/applications?page=1&page_size=50` | 1 | 是 |
| auto-default-24 | GET | `/api/v1/interviews/1` | 1 | 是 |
| auto-default-25 | GET | `/api/v1/interviews/1/feedback` | 1 | 是 |
| auto-default-26 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | 1 | 是 |
| auto-default-27 | GET | `/api/v1/logs?page=1&page_size=20` | 1 | 是 |
| auto-default-28 | GET | `/api/v1/notices/1` | 2 | 是 |
| auto-default-29 | GET | `/api/v1/notices?page=1&page_size=10` | 3 | 是 |
| auto-default-30 | GET | `/api/v1/notices/20` | 1 | 是 |
| auto-default-31 | GET | `/api/v1/notices?keyword=%E6%A0%A1%E5%9B%AD%E6%8B%9B%E8%81%98&page=1&page_size=10` | 1 | 是 |

### 回放结果

| 用例ID | 方法 | 路径 | 状态 | 期望 | 实际 | 耗时(ms) | 说明 |
| --- | --- | --- | --- | ---: | ---: | ---: | --- |
| auto-default-1 | POST | `/api/v1/login` | SKIP | - | - | - | POST body is not BIGTEST-owned |
| auto-default-2 | GET | `/api/v1/stats/dashboard` | PASS | 200 | 200 | 2 | ok |
| auto-default-3 | GET | `/api/v1/messages/unread-count` | PASS | 200 | 200 | 1 | ok |
| auto-default-4 | GET | `/api/v1/messages?page=1&page_size=5` | PASS | 200 | 200 | 2 | ok |
| auto-default-5 | GET | `/api/v1/talents?page=1&page_size=20` | PASS | 200 | 200 | 4 | ok |
| auto-default-6 | GET | `/api/v1/jobs?page=1&page_size=5` | PASS | 200 | 200 | 2 | ok |
| auto-default-7 | GET | `/api/v1/conversations/unread-count` | PASS | 200 | 200 | 1 | ok |
| auto-default-8 | GET | `/api/v1/jobs?page=1&page_size=10&status=open` | PASS | 200 | 200 | 2 | ok |
| auto-default-9 | GET | `/api/v1/applications?page=1&page_size=1000&talent_id=me` | PASS | 200 | 200 | 2 | ok |
| auto-default-10 | GET | `/api/v1/conversations` | PASS | 200 | 200 | 1 | ok |
| auto-default-11 | GET | `/api/v1/jobs/1` | PASS | 200 | 200 | 1 | ok |
| auto-default-12 | GET | `/api/v1/applications?page=1&page_size=20` | PASS | 200 | 200 | 7 | ok |
| auto-default-13 | GET | `/api/v1/resumes/online` | PASS | 200 | 200 | 1 | ok |
| auto-default-14 | GET | `/api/v1/resumes?page_size=10` | PASS | 200 | 200 | 6 | ok |
| auto-default-15 | GET | `/api/v1/talents?page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |
| auto-default-16 | GET | `/api/v1/jobs?page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |
| auto-default-17 | GET | `/api/v1/jobs/1/applications?page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |
| auto-default-18 | GET | `/api/v1/jobs?page_size=100&status=open` | PASS | 200 | 200 | 2 | ok |
| auto-default-19 | GET | `/api/v1/resumes?page=1&page_size=10` | PASS | 200 | 200 | 5 | ok |
| auto-default-20 | GET | `/api/v1/resumes?page=1&page_size=10&search=%E6%B5%8B%E8%AF%95` | PASS | 200 | 200 | 3 | ok |
| auto-default-21 | GET | `/api/v1/jobs/62` | PASS | 200 | 200 | 0 | ok |
| auto-default-22 | GET | `/api/v1/jobs?page=1&page_size=20` | PASS | 200 | 200 | 1 | ok |
| auto-default-23 | GET | `/api/v1/applications?page=1&page_size=50` | PASS | 200 | 200 | 10 | ok |
| auto-default-24 | GET | `/api/v1/interviews/1` | PASS | 200 | 200 | 2 | ok |
| auto-default-25 | GET | `/api/v1/interviews/1/feedback` | PASS | 200 | 200 | 1 | ok |
| auto-default-26 | GET | `/api/v1/messages?page=1&page_size=10&user_id=1` | PASS | 200 | 200 | 1 | ok |
| auto-default-27 | GET | `/api/v1/logs?page=1&page_size=20` | PASS | 200 | 200 | 2 | ok |
| auto-default-28 | GET | `/api/v1/notices/1` | PASS | 200 | 200 | 0 | ok |
| auto-default-29 | GET | `/api/v1/notices?page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |
| auto-default-30 | GET | `/api/v1/notices/20` | PASS | 200 | 200 | 1 | ok |
| auto-default-31 | GET | `/api/v1/notices?keyword=%E6%A0%A1%E5%9B%AD%E6%8B%9B%E8%81%98&page=1&page_size=10` | PASS | 200 | 200 | 1 | ok |

### 结果分析

| 指标 | 结果 | 说明 |
| --- | ---: | --- |
| 页面访问数 | 28 | 自动探索实际打开的前端页面数量 |
| 输入框填充 | 98/551 | 仅填写可见、可编辑且能识别字段含义的输入框 |
| 按钮点击 | 103/542 | 已允许 safe、medium、unknown 按钮，danger/empty 仍按规则跳过 |
| API 资源数 | 10 | 按 /api/v1/{resource} 归类统计 |
| 回放失败 | 0 | 状态码、响应结构或请求错误不符合预期 |
| 回放跳过 | 1 | 多为非 BIGTEST 自有写请求，避免污染真实数据 |

### 跳过原因统计

| 风险类型 | 数量 | 说明 |
| --- | ---: | --- |
| danger | 58 | 高风险页面操作，如删除、上传、AI评估、真实投递等 |
| empty | 123 | 按钮没有可识别文本，通常是图标按钮 |

### 后续建议

- 当前报告未发现失败接口，可作为本轮回归测试通过依据。
- 页面上的危险按钮没有直接点击；如需测试删除/投递等动作，请继续通过 `writeSafety.scenarios` 创建并操作 `BIGTEST_` 测试数据。
- 存在空文本图标按钮被跳过，若这些按钮很重要，建议给按钮补充 `aria-label` 或在测试框架中增加图标按钮识别规则。
