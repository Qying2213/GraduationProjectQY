# 智能招聘系统 - 测试套件

## 测试脚本说明

### 1. test_all.sh - 完整功能测试
```bash
chmod +x ztest/test_all.sh
./ztest/test_all.sh
```
123456789
覆盖所有服务的基本功能测试，包括：
- 服务健康检查
- 用户服务 (注册/登录/JWT)
- 职位服务 (CRUD)
- 人才服务 (CRUD)
- 简历服务 (上传/解析)
- AI评估服务 (OCR/Embedding/RAG/Coze)
- 推荐服务 (语义匹配/RAG)
- 面试服务
- 消息服务
- 数据库连接

### 2. test_api_detailed.sh - API详细测试
```bash
chmod +x ztest/test_api_detailed.sh
./ztest/test_api_detailed.sh
```
详细的API测试，显示完整的请求和响应内容。

### 3. test_ai_flow.sh - AI流程专项测试
```bash
chmod +x ztest/test_ai_flow.sh
./ztest/test_ai_flow.sh
```
专门测试AI评估的完整流程：
- PDF → OCR文本提取 → Embedding向量化 → RAG检索 → Coze AI评估

### 4. test_database.sh - 数据库测试
```bash
chmod +x ztest/test_database.sh
./ztest/test_database.sh
```
数据库相关测试：
- PostgreSQL连接
- pgvector扩展
- 表结构检查
- 数据统计
- 向量相似度查询

## 运行所有测试
```bash
# 赋予执行权限
chmod +x ztest/*.sh

# 运行完整测试
./ztest/test_all.sh

# 运行详细API测试
./ztest/test_api_detailed.sh

# 运行AI流程测试
./ztest/test_ai_flow.sh

# 运行数据库测试
./ztest/test_database.sh
```

## 测试结果
测试结果会保存在 `ztest/test_results_*.log` 文件中。

## 服务端口
| 服务 | 端口 |
|------|------|
| Gateway | 8080 |
| User Service | 8081 |
| Job Service | 8082 |
| Interview Service | 8083 |
| Resume Service | 8084 |
| Message Service | 8085 |
| Talent Service | 8086 |
| Recommendation Service | 8087 |
| Frontend | 5173 |

## 前置条件
1. 所有后端服务已启动
2. PostgreSQL数据库已启动
3. pgvector扩展已安装
4. 环境变量已配置 (COZE_TOKEN, ARK_API_KEY等)
