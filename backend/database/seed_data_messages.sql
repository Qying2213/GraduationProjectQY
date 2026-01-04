-- 追加消息数据 (messages表)
-- 字段: sender_id, receiver_id, type, title, content, is_read, read_at, created_at

DO $$
DECLARE
    v_admin_id INTEGER;
    v_hr_id INTEGER;
    v_interviewer_id INTEGER;
BEGIN
    -- 获取用户ID
    SELECT id INTO v_admin_id FROM users WHERE role = 'admin' LIMIT 1;
    SELECT id INTO v_hr_id FROM users WHERE role = 'hr' OR role = 'hr_manager' OR role = 'recruiter' LIMIT 1;
    SELECT id INTO v_interviewer_id FROM users WHERE role = 'interviewer' LIMIT 1;

    IF v_admin_id IS NULL THEN v_admin_id := 1; END IF;
    IF v_hr_id IS NULL THEN v_hr_id := v_admin_id; END IF;
    IF v_interviewer_id IS NULL THEN v_interviewer_id := v_admin_id; END IF;

    -- 系统通知
    INSERT INTO messages (sender_id, receiver_id, type, title, content, is_read, created_at) VALUES
    (v_admin_id, v_hr_id, 'system', '系统升级通知', '系统将于本周六凌晨2点-6点进行升级维护，届时系统将暂停服务。', FALSE, NOW() - INTERVAL '2 days'),
    (v_admin_id, v_hr_id, 'system', '新功能上线', '智能推荐功能已上线，系统将根据职位要求自动推荐匹配的候选人。', FALSE, NOW() - INTERVAL '1 day')
    ON CONFLICT DO NOTHING;

    -- 面试通知
    INSERT INTO messages (sender_id, receiver_id, type, title, content, is_read, created_at) VALUES
    (v_admin_id, v_hr_id, 'interview', '面试安排提醒', '您有一场面试即将开始：候选人张伟，职位高级Go开发工程师，时间：明天下午2点。', FALSE, NOW() - INTERVAL '1 day'),
    (v_admin_id, v_hr_id, 'interview', '面试安排提醒', '您有一场面试即将开始：候选人李强，职位Go开发工程师，时间：后天上午10点。', FALSE, NOW() - INTERVAL '12 hours'),
    (v_admin_id, v_interviewer_id, 'interview', '面试任务通知', '您需要进行技术面试：候选人孙婷，职位高级前端工程师，请准备面试问题。', FALSE, NOW() - INTERVAL '2 days')
    ON CONFLICT DO NOTHING;

    -- 反馈通知
    INSERT INTO messages (sender_id, receiver_id, type, title, content, is_read, created_at) VALUES
    (v_admin_id, v_hr_id, 'feedback', '新简历投递', '有新的简历投递：冯刚 申请了 高级Go开发工程师 职位，请及时查看。', FALSE, NOW() - INTERVAL '2 days'),
    (v_admin_id, v_hr_id, 'feedback', '新简历投递', '有新的简历投递：曹阳 申请了 Java开发工程师 职位，请及时查看。', FALSE, NOW() - INTERVAL '3 days'),
    (v_admin_id, v_hr_id, 'feedback', '新简历投递', '有新的简历投递：邓超 申请了 数据分析师 职位，请及时查看。', FALSE, NOW() - INTERVAL '4 days')
    ON CONFLICT DO NOTHING;

    -- Offer通知
    INSERT INTO messages (sender_id, receiver_id, type, title, content, is_read, created_at) VALUES
    (v_admin_id, v_hr_id, 'offer', 'Offer发放通知', '候选人陈浩的Offer已发放，职位：资深Go架构师，请跟进候选人回复。', FALSE, NOW() - INTERVAL '2 days'),
    (v_admin_id, v_hr_id, 'offer', 'Offer发放通知', '候选人董丽的Offer已发放，职位：前端技术专家，请跟进候选人回复。', FALSE, NOW() - INTERVAL '3 days'),
    (v_admin_id, v_hr_id, 'offer', 'Offer接受通知', '恭喜！候选人赵雪已接受Offer，职位：高级前端工程师，预计入职时间：下周一。', TRUE, NOW() - INTERVAL '5 days')
    ON CONFLICT DO NOTHING;

    -- 提醒通知
    INSERT INTO messages (sender_id, receiver_id, type, title, content, is_read, created_at) VALUES
    (v_admin_id, v_hr_id, 'reminder', '入职提醒', '候选人赵雪将于明天入职，职位：高级前端工程师，请做好入职准备工作。', TRUE, NOW() - INTERVAL '6 days'),
    (v_admin_id, v_hr_id, 'reminder', '入职提醒', '候选人杨帆将于下周一入职，职位：数据分析师，请做好入职准备工作。', TRUE, NOW() - INTERVAL '8 days'),
    (v_admin_id, v_hr_id, 'reminder', '待办事项提醒', '您有5份简历待审核，3场面试待安排，请及时处理。', FALSE, NOW() - INTERVAL '6 hours'),
    (v_admin_id, v_interviewer_id, 'reminder', '面试反馈提醒', '请及时提交候选人张伟的面试反馈，以便HR跟进后续流程。', FALSE, NOW() - INTERVAL '4 days')
    ON CONFLICT DO NOTHING;

END $$;
