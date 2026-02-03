# ztest - 测试脚本目录

本目录包含项目的全面测试脚本。

## 测试脚本

| 脚本               | 说明            | 用法                                              |
| ------------------ | --------------- | ------------------------------------------------- |
| `run_all_tests.sh` | Bash 全面测试   | `chmod +x run_all_tests.sh && ./run_all_tests.sh` |
| `test_api.py`      | Python API 测试 | `python3 test_api.py`                             |

## 前置条件

1. 确保所有后端服务已启动:

   ```bash
   cd /Users/qinyang/Desktop/GraduationProjectQY
   ./start-all.sh
   ```

2. Python 测试需要 requests 库:
   ```bash
   pip install requests
   ```

## 测试覆盖

- ✅ 服务健康检查 (8 个微服务)
- ✅ 用户认证 (注册/登录/获取信息)
- ✅ 职位管理 (列表/创建)
- ✅ 简历管理 (列表)
- ✅ 人才管理 (列表)
- ✅ 消息服务 (列表/未读数)
- ✅ 推荐服务 (列表)
- ✅ 面试服务 (列表)
- ✅ 统计接口 (7 个端点)
- ✅ AI 评估接口
- ✅ 潜在问题检测 (Redis/ES)

## 测试报告

运行后会生成 `test_report_*.md` 文件，包含详细结果。
