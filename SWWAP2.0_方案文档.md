# SWWAP 2.0 方案文档

## 1. 文档定位

本文档定义 `SWWAP 2.0` 的目标形态、核心设计、配置模型、执行流程与扩展策略。

`SWWAP 2.0` 不是某个单一业务页面的测试脚本，而是一个可跨项目复用的、配置驱动的 Web 测试平台。  
它的目标不是“替代所有测试手段”，而是把以下三类能力统一起来：

- 代码级确定性校验
- AI 辅助语义/页面校验
- 人工审核留痕与统一汇总

---

## 2. 核心目标

SWWAP 2.0 的核心目标如下：

1. 同一套测试框架可在多个 Web 项目中复用。
2. 测试内容由配置驱动，而不是为每个项目重写脚本。
3. 对“稳定、明确、可断言”的问题优先使用代码校验。
4. 对“页面语义、视觉表现、可用性”的问题引入 AI 辅助判断。
5. 对“极难自动化或强依赖人工经验”的问题保留人工审核入口。
6. 所有执行结果统一沉淀为结构化结果和 Markdown 报告。

---

## 3. 设计原则

### 3.1 代码优先

凡是可以通过接口状态码、返回字段、分页参数、筛选结果、边界值行为明确判断的内容，必须优先使用代码断言。

### 3.2 AI 补位

AI 不负责替代核心事实判断，只负责补足以下两类能力：

- 语义理解
- 页面/文案/布局合理性判断

### 3.3 人工兜底

对于非常难自动化、自动化成本过高、或需要人工主观确认的问题，必须作为正式测试类型纳入系统，而不是游离在系统之外。

### 3.4 配置驱动

测试项目、测试套件、环境参数、执行器选择、输入与预期，都应通过配置描述，而不是写死在代码里。

### 3.5 跨项目复用

框架层不得写死“教练列表”“用户列表”“订单页”等业务含义。  
业务字段映射和项目差异，应放在项目级配置或项目适配层中。

---

## 4. 总体架构

SWWAP 2.0 推荐采用如下分层：

### 4.1 配置层

负责描述：

- 项目基础配置
- 测试套件
- 测试用例
- 报告输出策略

### 4.2 执行器层

负责执行不同类型的测试：

- `code`
- `browser`
- `ai`
- `manual`

### 4.3 断言层

负责处理：

- HTTP 断言
- JSON 断言
- 字段断言
- 列表结果断言
- 边界值断言
- AI 输出格式校验

### 4.4 报告层

负责输出：

- 结构化 JSON 结果
- Markdown 报告
- 原始响应归档
- 截图归档
- AI 判定归档

---

## 5. 执行器分类

SWWAP 2.0 建议内置 4 类执行器。

### 5.1 `code`

适用场景：

- 接口测试
- 搜索筛选
- 分页
- 边界值
- 极端值
- 返回字段校验
- 错误码和错误信息校验
- JSON 合法性校验

典型问题：

- `page=0` 是否返回合理错误
- `limit=0` 是否错误返回全部数据
- 搜索 `id=abc` 是否空列表且不报错
- 多条件搜索后每条记录是否满足筛选条件

### 5.2 `browser`

适用场景：

- 打开页面
- 输入筛选条件
- 点击搜索 / 重置
- 翻页
- 点击详情
- 打开抽屉 / 弹窗
- 截图采集
- DOM 采集

典型问题：

- 点击重置后表单是否恢复默认状态
- 横向滚动是否可用
- 点击关联课程是否打开正确弹窗

### 5.3 `ai`

适用场景：

- 文案理解
- 页面布局合理性
- 可用性问题
- 语义是否清晰
- 页面视觉问题辅助判断

典型问题：

- `name` 和 `show_name` 是否容易混淆
- 状态与展示状态是否容易理解
- 长文本是否影响可读性
- 空状态和错误提示是否足够清楚

### 5.4 `manual`

适用场景：

- 极难自动化
- 需要人工主观确认
- 多设备、多浏览器、多分辨率体验确认
- 复杂交互的最终验收

典型问题：

- 横向滚动在真实分辨率下是否顺手
- 固定列是否存在遮挡
- 页面整体交互体验是否顺畅

---

## 6. 配置模型

SWWAP 2.0 推荐拆为两类 TOML 文件：

### 6.1 基础配置文件

建议命名：

- `project.toml`
- 或 `swwap.base.toml`

只放稳定信息：

- 项目名
- 基础地址
- 登录方式
- 公共请求头
- 默认超时
- AI 配置
- 报告输出目录

示例：

```toml
[project]
name = "coach-admin"

[site]
base_url = "https://example.com/api/admin/v1"
login_path = "/auth/login"

[auth]
type = "token_login"
username = "admin"
password = "123456"
token_field = "data.token"
header_name = "Authorization"
header_prefix = "Bearer "

[headers]
apifoxToken = "${APIFOX_TOKEN}"

[defaults]
timeout_seconds = 20
page = 1
limit = 10

[ai]
enabled = true
provider = "openai"
endpoint = "https://example-openai-proxy.com/v1/responses"
model = "gpt-5.4"
api_key_env = "OPENAI_API_KEY"

[report]
output_dir = "./reports"
save_raw_response = true
save_screenshot = true
```

### 6.2 用例配置文件

建议命名：

- `coach_list_cases.toml`
- `user_list_cases.toml`
- `order_detail_cases.toml`

只放：

- 测试内容
- 输入参数
- 预期结果
- 执行器类型
- AI 校验规则
- 人工审核规则

---

## 7. 推荐的用例结构

每条用例建议包含以下字段：

- `name`
- `module`
- `executor`
- `type`
- `tags`
- `input`
- `expect_http`
- `expect_api`
- `ai_check`
- `manual_check`

示例：

```toml
[[cases]]
name = "教练列表_按ID精确搜索_唯一命中"
module = "coach_list"
executor = "code"
type = "api_get"
tags = ["search", "stable"]
path = "/coaches/list"

[cases.input]
id = "2041821272066371584"
page = 1
limit = 10

[cases.expect_http]
status_code = 200
json_valid = true

[cases.expect_api]
total = 1
returned_count = 1
record_index = 0

[cases.expect_api.fields]
id = "2041821272066371584"
name = "test"
city = "北京市"
```

---

## 8. 教练列表测试点分类示例

以下以“教练管理 -> 教练列表”为例说明分类策略。

### 8.1 适合代码直接测试

- 默认第一页加载
- 搜索 `id`
- 搜索非法 `id`
- 搜索姓名
- 搜索不存在姓名
- 多条件组合搜索
- `page=0`
- `page=999`
- `limit=0`
- `limit` 超大
- 返回 JSON 是否合法
- 每条记录是否满足筛选条件

### 8.2 适合 AI 辅助测试

- `name` 和 `show_name` 是否容易混淆
- 状态与展示状态是否清晰
- 长常驻区域是否影响可读性
- 空图占位是否自然
- 空状态文案是否清晰
- 错误提示是否友好

### 8.3 适合混合模式

- 点击重置
- 封面图为空占位
- 状态与展示区分
- 加载失败提示

### 8.4 适合人工审核

- 横向滚动可用性
- 固定列遮挡
- 多分辨率可读性
- 复杂交互流畅度

---

## 9. 结果模型

每条用例最终应输出统一状态：

- `PASS`
- `FAIL`
- `WARN`
- `MANUAL_PENDING`
- `MANUAL_CONFIRMED`
- `SKIPPED`

每条结果应包含：

- 用例名称
- 执行器类型
- 输入参数
- 断言结果
- AI 判定结果
- 人工审核状态
- 原始响应引用
- 截图引用
- 最终结论

---

## 10. 报告输出

推荐每次执行后输出两类结果：

### 10.1 结构化结果

- `run.json`
- 便于机器读取
- 便于后续统计、聚合、回归对比

### 10.2 Markdown 报告

报告中应包含：

- 基本信息
- 执行范围
- 通过/失败统计
- 用例结果表
- 失败缺陷摘要
- 未自动化覆盖点
- 人工审核待确认项

---

## 11. 扩展性要求

为了后续跨项目复用，SWWAP 2.0 需要满足以下要求：

1. 核心代码中不写死具体项目字段。
2. 项目差异通过项目配置或项目适配层解决。
3. 用例文件可按模块拆分。
4. 新项目接入时，只需要新增：
   - 一个基础配置文件
   - 若干个用例配置文件
   - 必要时少量项目适配器

---

## 12. 推荐目录结构

```text
webtest/
  cmd/webtest/main.go
  internal/config/
  internal/executor/
  internal/assert/
  internal/report/
  internal/project/

projects/
  coach/
    project.toml
    suites/
      coach_list.toml
      coach_detail.toml
  another_project/
    project.toml
    suites/
      user_list.toml
```

---

## 13. 演进路线

### 阶段 1

先完成：

- 基础配置与用例配置分离
- `code` 执行器
- Markdown 报告
- JSON 结果输出

### 阶段 2

补充：

- `ai` 执行器
- 统一 AI 输出结构
- AI 与代码混合判定

### 阶段 3

补充：

- `browser` 执行器
- 截图与 DOM 采集
- 页面操作链

### 阶段 4

补充：

- `manual` 审核流
- 人工审核结果回填
- 统一汇总报告

---

## 14. 最终结论

SWWAP 2.0 最佳实践不是“把 AI 当成万能测试器”，而是建立一个分层明确、职责清晰的通用测试平台：

- 代码负责确定性事实
- AI 负责语义与页面辅助判断
- 人工负责复杂场景兜底
- TOML 负责配置驱动

如果按本方案落地，SWWAP 2.0 的价值不只是“测试教练列表”，而是可以演进为一个真正跨项目复用的 Web 自动化测试框架。
