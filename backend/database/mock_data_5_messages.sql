-- =====================================================
-- 9. 消息数据
-- =====================================================
INSERT INTO messages (sender_id, receiver_id, type, title, content, is_read, read_at) VALUES
-- 系统消息
(NULL, 1, 'system', '系统升级通知', '系统将于今晚22:00-24:00进行升级维护，届时部分功能可能无法使用。', true, '2025-12-20 09:00:00'),
(NULL, 2, 'system', '新功能上线', 'AI简历评估功能已上线，可在简历管理页面使用。', true, '2025-12-21 10:00:00'),
(NULL, 3, 'system', '新功能上线', 'AI简历评估功能已上线，可在简历管理页面使用。', false, NULL),
(NULL, 4, 'system', '新功能上线', 'AI简历评估功能已上线，可在简历管理页面使用。', false, NULL),

-- 面试相关消息
(3, 5, 'interview', '面试安排通知', '您有一场新的面试安排：张伟 - 高级前端工程师，时间：2025-12-26 10:00', true, '2025-12-25 09:00:00'),
(3, 6, 'interview', '面试安排通知', '您有一场新的面试安排：杨梅 - 高级前端工程师，时间：2025-12-30 10:00', false, NULL),
(4, 7, 'interview', '面试安排通知', '您有一场新的面试安排：王强 - Java开发工程师，时间：2025-12-31 14:00', false, NULL),
(3, 5, 'interview', '面试安排通知', '您有一场新的面试安排：孙浩 - DevOps工程师，时间：2025-12-27 10:00', false, NULL),

-- 反馈消息
(6, 3, 'feedback', '面试反馈已提交', '刘洋已提交张伟的面试反馈，评分：4分，建议：通过', true, '2025-12-20 12:00:00'),
(5, 3, 'feedback', '面试反馈已提交', '陈强已提交张伟的复试反馈，评分：5分，建议：通过', true, '2025-12-23 16:00:00'),
(7, 3, 'feedback', '面试反馈已提交', '赵磊已提交李娜的面试反馈，评分：5分，建议：通过', true, '2025-12-18 16:00:00'),
(8, 4, 'feedback', '面试反馈已提交', '孙婷已提交刘芳的面试反馈，评分：4分，建议：通过', true, '2025-12-19 12:00:00'),

-- Offer消息
(2, 3, 'offer', 'Offer审批通过', '李娜的Offer已审批通过，请尽快联系候选人。', true, '2025-12-24 17:00:00'),
(2, 4, 'offer', 'Offer审批通过', '朱琳的Offer已审批通过，请尽快联系候选人。', true, '2025-12-21 17:00:00'),
(2, 4, 'offer', 'Offer审批通过', '罗刚的Offer已审批通过，请尽快联系候选人。', false, NULL),

-- 提醒消息
(NULL, 3, 'reminder', '待处理简历提醒', '您有5份简历待筛选，请及时处理。', false, NULL),
(NULL, 4, 'reminder', '待处理简历提醒', '您有3份简历待筛选，请及时处理。', false, NULL),
(NULL, 6, 'reminder', '面试提醒', '您明天有1场面试安排，请做好准备。', false, NULL),
(NULL, 7, 'reminder', '面试提醒', '您后天有1场面试安排，请做好准备。', false, NULL),

-- 聊天消息
(3, 4, 'chat', '关于产品经理候选人', '刘芳的面试反馈很好，你觉得可以推进到HR面吗？', true, '2025-12-22 15:00:00'),
(4, 3, 'chat', '回复：关于产品经理候选人', '可以的，我来安排HR面试时间。', true, '2025-12-22 15:30:00'),
(2, 5, 'chat', '技术面试官安排', '下周有几个后端候选人需要面试，麻烦安排一下时间。', false, NULL);

-- =====================================================
-- 10. 操作日志数据
-- =====================================================
INSERT INTO operation_logs (user_id, action, resource_type, resource_id, details, ip_address, user_agent) VALUES
(1, 'login', 'user', 1, '{"method": "password"}', '192.168.1.100', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(2, 'login', 'user', 2, '{"method": "password"}', '192.168.1.101', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'),
(3, 'create', 'job', 1, '{"title": "高级前端工程师"}', '192.168.1.102', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(3, 'create', 'job', 2, '{"title": "后端开发工程师（Go）"}', '192.168.1.102', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(4, 'create', 'job', 3, '{"title": "产品经理"}', '192.168.1.103', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'),
(3, 'upload', 'resume', 1, '{"file_name": "张伟_高级前端工程师简历.pdf"}', '192.168.1.102', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(3, 'upload', 'resume', 2, '{"file_name": "李娜_Go开发工程师简历.pdf"}', '192.168.1.102', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(4, 'upload', 'resume', 4, '{"file_name": "刘芳_产品经理简历.pdf"}', '192.168.1.103', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'),
(3, 'create', 'interview', 1, '{"candidate": "张伟", "type": "initial"}', '192.168.1.102', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(3, 'create', 'interview', 4, '{"candidate": "李娜", "type": "initial"}', '192.168.1.102', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(6, 'submit', 'feedback', 1, '{"interview_id": 1, "rating": 4}', '192.168.1.104', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(5, 'submit', 'feedback', 2, '{"interview_id": 2, "rating": 5}', '192.168.1.105', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(2, 'approve', 'offer', 2, '{"candidate": "李娜", "position": "后端开发工程师（Go）"}', '192.168.1.101', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'),
(2, 'approve', 'offer', 14, '{"candidate": "朱琳", "position": "产品经理"}', '192.168.1.101', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'),
(3, 'update', 'application', 11, '{"stage": "hired", "candidate": "黄磊"}', '192.168.1.102', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(4, 'update', 'application', 20, '{"stage": "hired", "candidate": "唐欣"}', '192.168.1.103', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'),
(1, 'export', 'report', NULL, '{"type": "monthly_recruitment"}', '192.168.1.100', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'),
(2, 'view', 'dashboard', NULL, '{"page": "recruitment_overview"}', '192.168.1.101', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)');

-- =====================================================
-- 完成提示
-- =====================================================
SELECT '模拟数据插入完成！' as status;
SELECT '用户数: ' || COUNT(*) FROM users;
SELECT '职位数: ' || COUNT(*) FROM jobs;
SELECT '人才数: ' || COUNT(*) FROM talents;
SELECT '简历数: ' || COUNT(*) FROM resumes;
SELECT '应聘记录数: ' || COUNT(*) FROM applications;
SELECT '面试数: ' || COUNT(*) FROM interviews;
SELECT '面试反馈数: ' || COUNT(*) FROM interview_feedbacks;
SELECT '消息数: ' || COUNT(*) FROM messages;
SELECT '操作日志数: ' || COUNT(*) FROM operation_logs;
