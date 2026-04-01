# 毕业设计验收测试报告

生成时间：2026-04-01 11:37:18

## 汇总

- 总检查项：27
- 通过：24
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
| 健康检查 | talent 服务可用 | PASS | 通过 |
| 健康检查 | recommendation 服务可用 | PASS | 通过 |
| 认证 | HR 注册 | PASS | 通过 |
| 认证 | HR 登录 | PASS | 通过 |
| 认证 | 候选人 注册 | PASS | 通过 |
| 认证 | 候选人 登录 | PASS | 通过 |
| RBAC | 候选人禁止访问人才后台列表 | FAIL | HTTP 200, data={'code': 0, 'data': {'page': 1, 'page_size': 10, 'talents': [{'id': 26, 'created_at': '2026-03-31T11:34:16.633793+08:00', 'updated_at': '2026-03-31T11:34:16.633793+08:00', 'name': 'test', 'email': '2269022031@qq.com', 'phone': '15005089297', 'skills': None, 'experience': 0, 'education': '', 'status': 'active', 'tags': None, 'user_id': 37, 'location': '', 'salary': '', 'summary': '', 'gender': '', 'age': 0, 'current_company': '', 'current_position': '', 'source': '用户注册', 'match_score': 20}, {'id': 20, 'created_at': '2026-02-24T18:14:24.498678+08:00', 'updated_at': '2026-03-02T18:14:24.498678+08:00', 'name': '蔡静', 'email': 'caijing@gmail.com', 'phone': '13900000020', 'skills': ['合同审核', '法律咨询', '知识产权'], 'experience': 5, 'education': '硕士', 'status': 'active', 'tags': None, 'user_id': 25, 'location': '深圳', 'salary': '20K', 'summary': '', 'gender': '女', 'age': 29, 'current_company': '平安科技', 'current_position': '法务专员', 'source': '猎头', 'match_score': 55}, {'id': 19, 'created_at': '2026-02-23T18:14:24.498678+08:00', 'updated_at': '2026-03-02T18:14:24.498678+08:00', 'name': '冯宇', 'email': 'fengyu@gmail.com', 'phone': '13900000019', 'skills': ['客户服务', '团队管理', 'CRM'], 'experience': 4, 'education': '本科', 'status': 'active', 'tags': None, 'user_id': 24, 'location': '深圳', 'salary': '12K', 'summary': '', 'gender': '男', 'age': 27, 'current_company': '微众银行', 'current_position': '客服主管', 'source': '内推', 'match_score': 55}, {'id': 18, 'created_at': '2026-02-22T18:14:24.498678+08:00', 'updated_at': '2026-03-02T18:14:24.498678+08:00', 'name': '韩晓', 'email': 'hanxiao@gmail.com', 'phone': '13900000018', 'skills': ['行政管理', 'Office', '会议组织'], 'experience': 2, 'education': '大专', 'status': 'active', 'tags': None, 'user_id': 23, 'location': '深圳', 'salary': '8K', 'summary': '', 'gender': '女', 'age': 24, 'current_company': '顺丰科技', 'current_position': '行政专员', 'source': '招聘网站', 'match_score': 55}, {'id': 17, 'created_at': '2026-02-21T18:14:24.498678+08:00', 'updated_at': '2026-03-02T18:14:24.498678+08:00', 'name': '唐明', 'email': 'tangming@gmail.com', 'phone': '13900000017', 'skills': ['招聘', '员工关系', '培训', 'HRIS'], 'experience': 3, 'education': '本科', 'status': 'active', 'tags': None, 'user_id': 22, 'location': '北京', 'salary': '10K', 'summary': '', 'gender': '男', 'age': 26, 'current_company': '贝壳找房', 'current_position': '人事专员', 'source': '校招', 'match_score': 55}, {'id': 16, 'created_at': '2026-02-20T18:14:24.498678+08:00', 'updated_at': '2026-03-02T18:14:24.498678+08:00', 'name': '谢芳', 'email': 'xiefang@gmail.com', 'phone': '13900000016', 'skills': ['会计', '财务分析', 'Excel', 'SAP'], 'experience': 4, 'education': '本科', 'status': 'active', 'tags': None, 'user_id': 21, 'location': '杭州', 'salary': '15K', 'summary': '', 'gender': '女', 'age': 28, 'current_company': '蚂蚁金服', 'current_position': '财务专员', 'source': '招聘网站', 'match_score': 55}, {'id': 15, 'created_at': '2026-02-19T18:14:24.498678+08:00', 'updated_at': '2026-03-02T18:14:24.498678+08:00', 'name': '罗伟', 'email': 'luowei@gmail.com', 'phone': '13900000015', 'skills': ['市场推广', '品牌策划', 'SEM', 'SEO'], 'experience': 5, 'education': '本科', 'status': 'active', 'tags': None, 'user_id': 20, 'location': '上海', 'salary': '18K', 'summary': '', 'gender': '男', 'age': 29, 'current_company': '携程', 'current_position': '市场专员', 'source': '内推', 'match_score': 55}, {'id': 14, 'created_at': '2026-02-18T18:14:24.498678+08:00', 'updated_at': '2026-03-02T18:14:24.498678+08:00', 'name': '高丽', 'email': 'gaoli@gmail.com', 'phone': '13900000014', 'skills': ['数据分析', '用户运营', '活动策划'], 'experience': 2, 'education': '本科', 'status': 'active', 'tags': None, 'user_id': 19, 'location': '上海', 'salary': '12K', 'summary': '', 'gender': '女', 'age': 25, 'current_company': '拼多多', 'current_position': '运营专员', 'source': '招聘网站', 'match_score': 55}, {'id': 13, 'created_at': '2026-02-17T18:14:24.498678+08:00', 'updated_at': '2026-03-02T18:14:24.498678+08:00', 'name': '何军', 'email': 'hejun@gmail.com', 'phone': '13900000013', 'skills': ['Docker', 'Kubernetes', 'Jenkins', 'AWS', 'Terraform'], 'experience': 7, 'education': '本科', 'status': 'active', 'tags': None, 'user_id': 18, 'location': '北京', 'salary': '40K', 'summary': '', 'gender': '男', 'age': 31, 'current_company': '快手', 'current_position': 'DevOps工程师', 'source': '猎头', 'match_score': 55}, {'id': 12, 'created_at': '2026-02-16T18:14:24.498678+08:00', 'updated_at': '2026-03-02T18:14:24.498678+08:00', 'name': '林雪', 'email': 'linxue@gmail.com', 'phone': '13900000012', 'skills': ['Python', 'Django', 'TensorFlow', 'PyTorch'], 'experience': 4, 'education': '硕士', 'status': 'active', 'tags': None, 'user_id': 17, 'location': '北京', 'salary': '30K', 'summary': '', 'gender': '女', 'age': 27, 'current_company': '滴滴', 'current_position': 'Python开发', 'source': '招聘网站', 'match_score': 55}], 'total': 21}, 'message': 'success'} |
| RBAC | 候选人禁止访问操作日志 | FAIL | HTTP 200, data={'code': 0, 'data': {'logs': [], 'page': 1, 'page_size': 20, 'total': 0}, 'message': 'success', 'warning': 'Elasticsearch unavailable, returned empty logs'} |
| RBAC | HR 可访问操作日志 | PASS | 通过 |
| 目标一-简历解析 | 文本简历解析可提取核心字段 | PASS | 通过 |
| 目标一-风控 | 时间冲突风控可识别异常简历 | PASS | 通过 |
| 目标一-AI能力 | AI 配置可用 | PASS | 通过 |
| 目标二-可解释推荐 | 职位推荐结果包含 reason 和 match_details | PASS | 通过 |
| 目标二-语义匹配 | 人才推荐结果包含匹配解释 | PASS | 通过 |
| 目标二-统计 | 推荐统计接口可访问 | PASS | 通过 |
| 业务链路 | 候选人在线简历保存成功 | PASS | 通过 |
| 业务链路 | 候选人附件简历上传成功 | FAIL | HTTP 400, data={'code': 400, 'message': '当前用户还未创建人才档案'} |
| 目标三-性能 | 职位列表 满足 QPS/延迟指标 | PASS | 通过 |
| 目标三-性能 | 人才列表 满足 QPS/延迟指标 | PASS | 通过 |
| 目标三-性能 | 推荐统计 满足 QPS/延迟指标 | PASS | 通过 |
| 目标三-性能 | 核心读接口整体满足开题性能目标 | PASS | 通过 |

## 结论

当前版本仍有未通过项，需要继续修复后再作为最终验收证据。