-- 添加消息数据
-- 执行方式: psql -d talent_platform -f backend/database/add_messages.sql

-- =====================================================
-- 1. 系统消息/通知数据 (messages 表)
-- =====================================================
INSERT INTO messages (sender_id, receiver_id, type, title, content, is_read, created_at) VALUES
-- 给 admin 的消息
(2, 1, 'system', '新简历投递通知', '张伟投递了高级Go后端工程师职位，请及时查看。', false, NOW() - INTERVAL '2 hours'),
(3, 1, 'system', '面试安排提醒', '明天上午10:00有一场技术面试，候选人：李娜，职位：前端开发工程师。', false, NOW() - INTERVAL '1 hour'),
(4, 1, 'interview', '面试反馈已提交', 'HR王芳已提交对候选人王强的面试反馈，评分：5星。', true, NOW() - INTERVAL '3 hours'),
-- 给 HR 的消息
(1, 2, 'system', '系统更新通知', '系统已升级到最新版本，新增AI智能评估功能。', true, NOW() - INTERVAL '1 day'),
(1, 2, 'resume', '简历解析完成', '候选人陈龙的简历已完成AI解析，匹配度85%。', false, NOW() - INTERVAL '30 minutes'),
(1, 3, 'interview', '面试时间变更', '候选人赵敏的面试时间已调整为2月1日下午2点。', false, NOW() - INTERVAL '4 hours'),
(1, 4, 'system', '招聘数据周报', '本周新增简历28份，安排面试15场，发放Offer 6个。', true, NOW() - INTERVAL '2 days'),
-- 给求职者的消息
(2, 6, 'interview', '面试邀请', '恭喜您！您投递的高级Go后端工程师职位已通过简历筛选，诚邀您参加面试。', false, NOW() - INTERVAL '1 day'),
(2, 7, 'resume', '简历已查看', '您投递的前端开发工程师职位，HR已查看您的简历。', true, NOW() - INTERVAL '2 days'),
(3, 8, 'offer', 'Offer通知', '恭喜您！您已通过所有面试环节，我们向您发出正式Offer。', false, NOW() - INTERVAL '12 hours'),
(2, 9, 'interview', '面试结果通知', '您的产品经理职位面试已通过，请等待下一轮面试安排。', true, NOW() - INTERVAL '3 days'),
(4, 10, 'system', '简历投递成功', '您的简历已成功投递Java开发工程师职位，请耐心等待。', true, NOW() - INTERVAL '5 days'),
-- 更多系统通知
(1, 2, 'announcement', '公司公告', '公司将于下周一举办技术分享会，欢迎参加。', false, NOW() - INTERVAL '6 hours'),
(1, 3, 'announcement', '节假日通知', '春节假期安排：1月28日至2月4日放假，2月5日正常上班。', true, NOW() - INTERVAL '1 week'),
(1, 4, 'system', '权限变更通知', '您的账号已被授予招聘主管权限。', true, NOW() - INTERVAL '10 days');

-- =====================================================
-- 完成
-- =====================================================
SELECT '消息数据添加完成！' as message;
SELECT '消息数量: ' || COUNT(*) FROM messages;
