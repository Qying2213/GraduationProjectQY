-- 为 admin 用户添加聊天会话
-- 执行方式: psql -d talent_platform -f backend/database/add_admin_chats.sql

-- =====================================================
-- 1. 为 admin (id=1) 添加会话
-- =====================================================
INSERT INTO conversations (participant_a, participant_b, last_message_at, created_at, updated_at) VALUES
(1, 6, NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '2 days', NOW()),
(1, 7, NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 day', NOW()),
(1, 8, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '3 days', NOW());

-- =====================================================
-- 2. 添加聊天消息
-- =====================================================
-- admin 和 张伟(id=6) 的对话 (conversation_id=5)
INSERT INTO chat_messages (conversation_id, sender_id, content, message_type, is_read, created_at) VALUES
(5, 1, '张伟你好，我是系统管理员，欢迎加入我们的人才平台。', 'text', true, NOW() - INTERVAL '2 days'),
(5, 6, '您好管理员，谢谢！请问有什么需要注意的吗？', 'text', true, NOW() - INTERVAL '2 days' + INTERVAL '5 minutes'),
(5, 1, '请确保您的简历信息完整，这样可以获得更好的职位匹配。', 'text', true, NOW() - INTERVAL '1 day'),
(5, 6, '好的，我会尽快完善简历。', 'text', false, NOW() - INTERVAL '30 minutes');

-- admin 和 李娜(id=7) 的对话 (conversation_id=6)
INSERT INTO chat_messages (conversation_id, sender_id, content, message_type, is_read, created_at) VALUES
(6, 1, '李娜你好，看到你投递了前端工程师职位，有什么问题可以随时问我。', 'text', true, NOW() - INTERVAL '1 day'),
(6, 7, '谢谢！请问面试一般是什么流程？', 'text', true, NOW() - INTERVAL '1 day' + INTERVAL '10 minutes'),
(6, 1, '一般是简历筛选 -> 技术面试 -> HR面试 -> 发放Offer。', 'text', false, NOW() - INTERVAL '1 hour');

-- admin 和 王强(id=8) 的对话 (conversation_id=7)
INSERT INTO chat_messages (conversation_id, sender_id, content, message_type, is_read, created_at) VALUES
(7, 8, '管理员您好，我想咨询一下高级后端工程师的薪资范围。', 'text', true, NOW() - INTERVAL '3 days'),
(7, 1, '您好，高级后端工程师薪资范围是25K-40K，具体根据面试表现和经验确定。', 'text', true, NOW() - INTERVAL '3 days' + INTERVAL '15 minutes'),
(7, 8, '明白了，谢谢您的解答！', 'text', false, NOW() - INTERVAL '2 hours');

-- 更新会话的 last_message_id
UPDATE conversations SET last_message_id = 14 WHERE id = 5;
UPDATE conversations SET last_message_id = 17 WHERE id = 6;
UPDATE conversations SET last_message_id = 20 WHERE id = 7;

-- =====================================================
-- 完成
-- =====================================================
SELECT 'Admin聊天数据添加完成！' as message;
SELECT '会话数量: ' || COUNT(*) FROM conversations;
SELECT '聊天消息数量: ' || COUNT(*) FROM chat_messages;
