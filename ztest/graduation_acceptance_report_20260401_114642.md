# 毕业设计验收测试报告

生成时间：2026-04-01 11:46:42

## 汇总

- 总检查项：34
- 通过：31
- 失败：3
- 警告：0
- 跳过：0

## 结果明细

| 领域 | 检查项 | 结果 | 说明 |
|---|---|---|---|
| 健康检查 | gateway 服务可用 | PASS | 通过 |
| 健康检查 | user 服务可用 | PASS | 通过 |
| 健康检查 | job 服务可用 | PASS | 通过 |
| 健康检查 | interview 服务可用 | PASS | 通过 |
| 健康检查 | resume 服务可用 | PASS | 通过 |
| 健康检查 | message 服务可用 | PASS | 通过 |
| 健康检查 | talent 服务可用 | FAIL | HTTPConnectionPool(host='localhost', port=8086): Max retries exceeded with url: /health (Caused by NewConnectionError("HTTPConnection(host='localhost', port=8086): Failed to establish a new connection: [Errno 61] Connection refused")) |
| 健康检查 | recommendation 服务可用 | PASS | 通过 |
| 认证 | HR 注册 | PASS | 通过 |
| 认证 | HR 登录 | PASS | 通过 |
| 认证 | 候选人 注册 | PASS | 通过 |
| 认证 | 候选人 登录 | PASS | 通过 |
| RBAC | 候选人禁止访问人才后台列表 | FAIL | HTTP 502, data={'raw': ''} |
| RBAC | 候选人禁止访问操作日志 | PASS | 通过 |
| RBAC | HR 可访问操作日志 | PASS | 通过 |
| 目标一-简历解析 | 文本简历解析可提取核心字段 | PASS | 通过 |
| 目标一-风控 | 时间冲突风控可识别异常简历 | PASS | 通过 |
| 目标一-AI能力 | AI 配置可用 | PASS | 通过 |
| 目标二-可解释推荐 | 职位推荐结果包含 reason 和 match_details | PASS | 通过 |
| 目标二-语义匹配 | 人才推荐结果包含匹配解释 | PASS | 通过 |
| 目标二-统计 | 推荐统计接口可访问 | PASS | 通过 |
| 业务链路 | 候选人在线简历保存成功 | PASS | 通过 |
| 业务链路 | 候选人附件简历上传成功 | PASS | 通过 |
| 业务链路 | HR 创建职位成功 | PASS | 通过 |
| 业务链路 | 候选人成功投递职位 | PASS | 通过 |
| 业务链路 | 候选人可查看我的投递记录 | PASS | 通过 |
| 业务链路 | HR 可查看简历列表 | PASS | 通过 |
| 业务链路 | HR 可查看人才列表 | FAIL | HTTP 502, data={'raw': ''} |
| 业务链路 | HR 可成功安排面试 | PASS | 通过 |
| 业务链路 | 面试安排后候选人收到消息通知 | PASS | 通过 |
| 目标三-性能 | 职位列表 满足 QPS/延迟指标 | PASS | 通过 |
| 目标三-性能 | 人才列表 满足 QPS/延迟指标 | PASS | 通过 |
| 目标三-性能 | 推荐统计 满足 QPS/延迟指标 | PASS | 通过 |
| 目标三-性能 | 核心读接口整体满足开题性能目标 | PASS | 通过 |

## 结论

当前版本仍有未通过项，需要继续修复后再作为最终验收证据。