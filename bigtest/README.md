# bigtest

`bigtest` 是一个可复制到其他项目里的真实流量自动化测试框架。它的目标是：少写配置，自动登录、自动发现页面、自动抓取接口、自动生成用例、自动回放，并在安全边界内测试写操作。

## 能做什么

- 自动识别登录页输入框和登录按钮，减少 selector 配置
- 自动发现路由：支持 Vue Router、React Router 常见文件、`src/views` / `src/pages` 页面目录、页面 DOM 链接兜底
- 自动探索页面：访问页面、填写输入框、点击安全按钮
- 自动抓取真实浏览器发出的 API 请求
- 自动去重并生成接口回放用例
- 自动识别 API 资源覆盖情况
- 支持安全写操作：只允许创建、更新、删除本轮 `BIGTEST_<runId>` 数据
- 输出原始抓包、生成用例、安全写结果、发现结果、回放结果和 Markdown 报告

## 最小配置

复制 `config.example.json` 为 `config.local.json`，通常只需要改这几个值：

```json
{
  "env": {
    "frontendBaseUrl": "http://localhost:5173",
    "replayBaseUrl": "http://localhost:8080"
  },
  "auth": {
    "loginPath": "/login",
    "username": "admin",
    "password": "admin123"
  }
}
```

如果登录页结构特殊，再补：

```json
{
  "auth": {
    "usernameSelector": "input[name='username']",
    "passwordSelector": "input[type='password']",
    "submitSelector": "button[type='submit']"
  }
}
```

## 运行

```bash
cd bigtest
npm install
npm run run -- ./config.local.json
```

如果本机已安装 Chrome，可以在配置里加：

```json
{
  "browser": {
    "channel": "chrome"
  }
}
```

## 自动发现策略

`bigtest` 默认会自动开启页面探索。发现顺序是：

1. 读取常见路由文件，例如 `frontend/src/router/index.ts`、`src/router/index.ts`、`src/routes.ts`、`src/App.tsx`
2. 扫描页面目录，例如 `frontend/src/views`、`frontend/src/pages`、`src/views`、`src/pages`
3. 登录后从页面里的 `<a href>` 和链接元素继续扩展

默认不会访问：

- `/ai`
- `/evaluate`
- `/evaluation`
- `/upload`
- `/dev-logs`

## 页面元素策略

按钮会被分成 4 类：

- `safe`：搜索、查询、查看、详情、返回、重置
- `medium`：新增、编辑、保存、提交、发布
- `danger`：删除、移除、清空、停用、关闭、上传、评估
- `unknown`：无法判断

默认框架只点击 `safe`。如果希望做更完整的全站覆盖，可以在配置里开启：

```json
{
  "autoExplore": {
    "clickRiskLevels": ["safe", "medium"]
  }
}
```

这样会覆盖“新增/编辑/保存/提交/发布/确定”等功能按钮。删除、上传、AI 评估等仍建议保留在 `dangerClickTexts` 或排除路径中，不在页面探索阶段随机点击。

输入框会按字段信息自动生成测试值：

- 搜索框：`测试`
- 邮箱：`bigtest_<runId>@example.com`
- 手机号：`13800000000`
- 日期：未来日期
- 标题/名称：`BIGTEST_<runId>_名称`
- 内容/描述：`BIGTEST_<runId>_自动化测试内容`

## 安全写操作

`bigtest` 支持 `POST / PUT / DELETE`，但有硬性保护：

- 创建数据必须包含本轮 `BIGTEST_<runId>` 前缀
- 创建成功后必须能提取 ID
- 更新、删除前必须校验 ID 已登记
- 命中 AI、评估、上传等危险路径会直接拦截
- 结束后自动清理本轮创建的数据

你可以手动配置高价值 CRUD 场景，也可以让框架从真实流量中推断资源型接口。

```json
{
  "writeSafety": {
    "scenarios": [
      {
        "name": "notice-crud-owned-data",
        "resourceType": "notice",
        "create": {
          "method": "POST",
          "path": "/api/v1/notices",
          "body": {
            "title": "${prefix}_公告",
            "content": "${prefix}_自动创建的公告内容",
            "status": "draft"
          },
          "idPath": "data.id",
          "expectedStatuses": [200, 201]
        },
        "verify": {
          "method": "GET",
          "path": "/api/v1/notices/${id}",
          "expectedStatuses": [200]
        },
        "update": {
          "method": "PUT",
          "path": "/api/v1/notices/${id}",
          "body": {
            "title": "${prefix}_公告_已更新",
            "content": "${prefix}_自动更新的公告内容",
            "status": "draft"
          },
          "expectedStatuses": [200]
        },
        "delete": {
          "method": "DELETE",
          "path": "/api/v1/notices/${id}",
          "expectedStatuses": [200, 204]
        }
      }
    ]
  }
}
```

## 输出目录

每次运行都会生成：

```text
bigtest/output/<run-id>/
  report.md
  <profile-name>/
    raw-captures.json
    generated-cases.json
    discovery-results.json
    write-safety-results.json
    replay-results.json
```

## 报告会展示

- 自动发现了多少路由
- 实际访问了多少页面
- 发现并填写了多少输入框
- 发现并安全点击了多少按钮
- 哪些按钮因为风险被跳过
- 识别了哪些 API 资源
- 哪些接口回放通过、失败或跳过
- 哪些安全写操作成功或失败

## 迁移到别的项目

1. 把 `bigtest/` 复制到目标项目根目录。
2. 执行 `cd bigtest && npm install`。
3. 复制 `config.example.json` 为 `config.local.json`。
4. 修改前端地址、后端地址、登录账号密码。
5. 运行 `npm run run -- ./config.local.json`。
6. 如果自动发现不准，再少量补 selector、排除路径或写操作场景。

## 注意事项

- 第一次跑建议在测试环境或本地环境运行。
- 不要把“删除/保存/提交/发布”加入 `safeClickTexts`。
- 写操作只会操作 `BIGTEST_` 数据，但仍建议用测试数据库。
- AI、评估、上传类接口默认跳过，避免外部依赖和文件副作用。
- 删除操作应通过 `writeSafety.scenarios` 执行：先创建本轮 `BIGTEST_<runId>` 测试数据，再更新/删除该数据，不删除真实业务数据。
# bigtest

`bigtest` 是一个基于**真实浏览器网络流量**的自动化测试框架。

它的目标不是从 Swagger 猜接口，而是：

1. 用真实浏览器登录系统
2. 按你配置的页面路径和操作步骤跑一遍真实业务流
3. 抓取前端真正发出的 `/api/v1/*` 请求
4. 对请求做归一化、去重、脱敏，生成标准化测试用例
5. 用新登录拿到的 token 自动回放这些用例
6. 自动读取前端路由做安全页面探索
7. 用 `BIGTEST_` 前缀创建测试数据，只改/删自己创建的数据
8. 校验状态码、响应结构，并输出报告

## 框架特点

- 基于 Playwright 抓真实浏览器网络事件，不解析 F12 面板 UI
- 支持把账号密码和页面步骤都写在配置文件里
- 支持 `goto / click / fill / press / select / wait / check / uncheck`
- 支持按请求指纹去重，避免轮询和重复请求把测试污染
- 支持自动读取 `frontend/src/router/index.ts` 并探索页面
- 支持安全写操作：`POST / PUT / DELETE` 只允许操作本轮登记的 `BIGTEST_` 数据
- 自动重建 `Authorization`，不会直接复用旧 token
- 生成：
  - 原始抓包
  - 去重后的标准化测试用例
  - 安全写操作结果
  - 回放结果
  - Markdown 报告

## 目录结构

```text
bigtest/
  config.example.json
  package.json
  README.md
  src/
    index.mjs
    framework.mjs
    utils.mjs
  output/
    <run-id>/
      <profile-name>/
        raw-captures.json
        generated-cases.json
        write-safety-results.json
        replay-results.json
      report.md
```

## 安装

在项目根目录执行：

```bash
cd bigtest
npm install
```

如果本机没有 Playwright 可用浏览器，可以再执行：

```bash
npx playwright install chromium
```

如果你本机装了 Chrome，配置里也可以直接用：

```json
"browser": {
  "channel": "chrome"
}
```

## 配置文件

先复制一份本地配置：

```bash
cp bigtest/config.example.json bigtest/config.local.json
```

然后修改 `bigtest/config.local.json`。

### 顶层配置

- `env.frontendBaseUrl`: 前端地址
- `env.replayBaseUrl`: 回放时实际请求的后端地址
- `browser`: 浏览器启动配置
- `capture`: 抓包过滤规则
- `generation`: 是否写出原始抓包和生成用例
- `replay`: 回放规则
- `autoExplore`: 自动探索页面规则
- `writeSafety`: 安全写操作规则
- `profiles`: 多角色/多场景配置

### profile 配置

每个 profile 代表一类角色或一条业务流，比如：

- `admin-smoke`
- `hr-notice-flow`
- `candidate-job-search`

每个 profile 里有两块最重要：

#### `auth`

- 登录页路径
- 用户名密码
- 用户名输入框 selector
- 密码输入框 selector
- 登录按钮 selector
- 登录成功 URL 关键字

#### `flow.steps`

支持的步骤：

- `goto`
- `click`
- `fill`
- `press`
- `select`
- `check`
- `uncheck`
- `wait`

### 一个更偏业务的步骤例子

```json
{
  "name": "notice-admin-flow",
  "auth": {
    "loginPath": "/login",
    "username": "admin",
    "password": "admin123",
    "usernameSelector": "input[placeholder='用户名 / 邮箱']",
    "passwordSelector": "input[placeholder='密码']",
    "submitSelector": "button.login-btn",
    "successUrlContains": "/dashboard",
    "postLoginWaitMs": 1500
  },
  "flow": {
    "steps": [
      { "type": "goto", "path": "/notices", "waitUntil": "networkidle", "waitForMs": 1200 },
      { "type": "click", "selector": "button:has-text('新建公告')", "waitForMs": 600 },
      { "type": "fill", "selector": "input[placeholder='请输入标题']", "value": "bigtest 公告" },
      { "type": "fill", "selector": "textarea", "value": "这是自动化测试框架生成的内容" },
      { "type": "wait", "ms": 1000 }
    ]
  }
}
```

## 运行

在项目根目录执行：

```bash
cd bigtest
npm run run -- ./config.local.json
```

如果不传配置路径，默认读取：

```text
bigtest/config.local.json
```

## 输出说明

每次运行都会生成一个新的 `run-id` 目录。

### `raw-captures.json`

浏览器抓到的原始请求，已经做了脱敏：

- token、cookie、password 不会明文写出

### `generated-cases.json`

这是去重、归一化之后的标准化测试用例。

这一步很重要，因为框架不是直接把零散抓包拿来重放，而是先整理成：

- 请求方法
- 规范化路径
- 是否需要鉴权
- 期望状态码
- 期望响应结构

### `write-safety-results.json`

安全写操作结果，包含本轮测试数据前缀、登记过的资源 ID、每个 CRUD 步骤的状态。

框架默认使用这样的保护链：

- 创建数据时，正文里必须包含本轮 `BIGTEST_<runId>` 前缀
- 创建成功后，从响应里提取 ID 并登记
- `PUT / DELETE` 前先检查目标 ID 是否已经登记
- 命中 `/ai/`、`/evaluate`、`/evaluation`、`/upload` 等危险路径会直接拦截
- 结束时自动执行清理步骤，删除自己创建的测试数据

### `replay-results.json`

实际自动回放的结果，包含：

- 期望状态码
- 实际状态码
- 是否通过
- 耗时
- 失败原因

### `report.md`

Markdown 报告，适合直接看整体结论。

## 自动探索策略

`autoExplore.enabled=true` 后，框架会读取前端路由文件，自动访问页面并尝试发现更多真实请求。

为了避免变成乱点机器人，它默认：

- 不访问 `/ai`、`/evaluate`、`/evaluation`、`/upload`、`/dev-logs`
- 每个页面最多点击 `maxClicksPerPage` 次
- 每个页面最多停留 `maxSecondsPerPage` 秒
- 只点击 `safeClickTexts` 白名单按钮
- 不点击 `删除/保存/提交/发布/上传/评估` 等危险按钮
- 出现确认弹窗自动取消
- 只自动填写搜索类输入框

## 安全写操作策略

`writeSafety.enabled=true` 后，框架可以执行配置里的 CRUD 场景。比如当前配置已经提供了 `notice-crud-owned-data`：

1. `POST /api/v1/notices` 创建标题带 `BIGTEST_<runId>` 前缀的公告
2. 从响应 `data.id` 读取 ID 并登记
3. `GET /api/v1/notices/${id}` 验证刚创建的数据
4. `PUT /api/v1/notices/${id}` 只更新登记过的 ID
5. 结束时 `DELETE /api/v1/notices/${id}` 自动清理

如果你以后要测职位、角色、消息，也是在 `writeSafety.scenarios` 里补新的资源场景，不需要改框架核心。

## 推荐用法

### 1. 先做一个 admin 冒烟流

- 登录
- 仪表盘
- 职位列表
- 消息列表
- 公告列表

### 2. 再补一个 candidate 流

- 求职端登录
- 职位列表
- 职位详情
- 我的投递

### 3. 再补一个 notice / job 的细场景

让配置里的 `flow.steps` 更接近真实操作。

## 后续怎么扩

这套框架下一步最值得扩的方向有 3 个：

1. 给职位、角色、消息等模块继续补 `writeSafety.scenarios`
2. 从列表响应里自动提取真实 ID，减少 `routeParamValues.id` 这种固定占位
3. 接入 CI，做定时回归

## 注意事项

- 写操作只建议在测试环境开启
- 不要把 `safeClickTexts` 配成“保存/提交/删除”这类危险按钮
- 如果页面请求有轮询接口，建议通过 `capture.ignoreUrlPatterns` 过滤
- 如果某些 query 参数是时间戳或随机数，建议放进 `capture.ignoreQueryKeys`
- 如果登录页 DOM 改了，优先更新 `auth.*Selector`
