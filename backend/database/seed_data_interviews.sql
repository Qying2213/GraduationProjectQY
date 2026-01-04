-- 追加面试数据 (interviews表)
-- 字段: candidate_id, candidate_name, position_id, position, type, date, time, duration,
--       interviewer_id, interviewer, method, location, status, notes, feedback, rating, created_by

DO $$
DECLARE
    v_interviewer_id INTEGER;
    v_interviewer_name VARCHAR(100);
    v_talent_id INTEGER;
    v_job_id INTEGER;
    v_job_title VARCHAR(200);
BEGIN
    -- 获取面试官
    SELECT id, username INTO v_interviewer_id, v_interviewer_name 
    FROM users WHERE role = 'interviewer' LIMIT 1;
    IF v_interviewer_id IS NULL THEN
        SELECT id, username INTO v_interviewer_id, v_interviewer_name 
        FROM users WHERE role = 'admin' LIMIT 1;
    END IF;
    IF v_interviewer_name IS NULL THEN v_interviewer_name := '面试官'; END IF;

    -- 张伟的面试
    SELECT id INTO v_talent_id FROM talents WHERE name = '张伟' OR email LIKE 'zhangwei%' LIMIT 1;
    SELECT id, title INTO v_job_id, v_job_title FROM jobs WHERE title LIKE '%Go%' LIMIT 1;
    IF v_talent_id IS NOT NULL AND v_job_id IS NOT NULL THEN
        INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, duration, interviewer_id, interviewer, method, location, status, notes, feedback, rating, created_by, created_at)
        VALUES (v_talent_id, '张伟', v_job_id, COALESCE(v_job_title, 'Go开发工程师'), 'initial', TO_CHAR(NOW() - INTERVAL '5 days', 'YYYY-MM-DD'), '14:00', 60, v_interviewer_id, v_interviewer_name, 'onsite', '会议室A', 'completed', '技术一面', '技术扎实，微服务经验丰富', 4, v_interviewer_id, NOW() - INTERVAL '8 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 李强的面试
    SELECT id INTO v_talent_id FROM talents WHERE name = '李强' OR email LIKE 'liqiang%' LIMIT 1;
    IF v_talent_id IS NOT NULL AND v_job_id IS NOT NULL THEN
        INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, duration, interviewer_id, interviewer, method, location, status, notes, created_by, created_at)
        VALUES (v_talent_id, '李强', v_job_id, COALESCE(v_job_title, 'Go开发工程师'), 'initial', TO_CHAR(NOW() + INTERVAL '2 days', 'YYYY-MM-DD'), '10:00', 60, v_interviewer_id, v_interviewer_name, 'video', '线上腾讯会议', 'scheduled', '技术面试', v_interviewer_id, NOW() - INTERVAL '3 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 陈浩的面试
    SELECT id INTO v_talent_id FROM talents WHERE name = '陈浩' OR email LIKE 'chenhao%' LIMIT 1;
    SELECT id, title INTO v_job_id, v_job_title FROM jobs WHERE title LIKE '%架构%' LIMIT 1;
    IF v_talent_id IS NOT NULL AND v_job_id IS NOT NULL THEN
        INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, duration, interviewer_id, interviewer, method, location, status, notes, feedback, rating, created_by, created_at)
        VALUES (v_talent_id, '陈浩', v_job_id, COALESCE(v_job_title, '架构师'), 'final', TO_CHAR(NOW() - INTERVAL '10 days', 'YYYY-MM-DD'), '10:00', 90, v_interviewer_id, v_interviewer_name, 'onsite', '会议室A', 'completed', '架构面试', '系统设计能力出色', 5, v_interviewer_id, NOW() - INTERVAL '12 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 赵雪的面试
    SELECT id INTO v_talent_id FROM talents WHERE name = '赵雪' OR email LIKE 'zhaoxue%' LIMIT 1;
    SELECT id, title INTO v_job_id, v_job_title FROM jobs WHERE title LIKE '%前端%' LIMIT 1;
    IF v_talent_id IS NOT NULL AND v_job_id IS NOT NULL THEN
        INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, duration, interviewer_id, interviewer, method, location, status, notes, feedback, rating, created_by, created_at)
        VALUES (v_talent_id, '赵雪', v_job_id, COALESCE(v_job_title, '前端工程师'), 'second', TO_CHAR(NOW() - INTERVAL '15 days', 'YYYY-MM-DD'), '10:00', 60, v_interviewer_id, v_interviewer_name, 'onsite', '会议室A', 'completed', '技术面试', '前端技术全面，有架构能力', 5, v_interviewer_id, NOW() - INTERVAL '18 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 孙婷的面试
    SELECT id INTO v_talent_id FROM talents WHERE name = '孙婷' OR email LIKE 'sunting%' LIMIT 1;
    IF v_talent_id IS NOT NULL AND v_job_id IS NOT NULL THEN
        INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, duration, interviewer_id, interviewer, method, location, status, notes, feedback, rating, created_by, created_at)
        VALUES (v_talent_id, '孙婷', v_job_id, COALESCE(v_job_title, '前端工程师'), 'initial', TO_CHAR(NOW() - INTERVAL '6 days', 'YYYY-MM-DD'), '14:00', 60, v_interviewer_id, v_interviewer_name, 'onsite', '会议室A', 'completed', '技术一面', 'React技术扎实', 4, v_interviewer_id, NOW() - INTERVAL '9 days')
        ON CONFLICT DO NOTHING;
        
        INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, duration, interviewer_id, interviewer, method, location, status, notes, created_by, created_at)
        VALUES (v_talent_id, '孙婷', v_job_id, COALESCE(v_job_title, '前端工程师'), 'second', TO_CHAR(NOW() + INTERVAL '1 day', 'YYYY-MM-DD'), '10:00', 60, v_interviewer_id, v_interviewer_name, 'video', '线上', 'scheduled', '技术二面', v_interviewer_id, NOW() - INTERVAL '3 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 马超的面试
    SELECT id INTO v_talent_id FROM talents WHERE name = '马超' OR email LIKE 'machao%' LIMIT 1;
    SELECT id, title INTO v_job_id, v_job_title FROM jobs WHERE title LIKE '%DevOps%' OR title LIKE '%运维%' LIMIT 1;
    IF v_talent_id IS NOT NULL AND v_job_id IS NOT NULL THEN
        INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, duration, interviewer_id, interviewer, method, location, status, notes, feedback, rating, created_by, created_at)
        VALUES (v_talent_id, '马超', v_job_id, COALESCE(v_job_title, 'DevOps工程师'), 'initial', TO_CHAR(NOW() - INTERVAL '4 days', 'YYYY-MM-DD'), '10:00', 60, v_interviewer_id, v_interviewer_name, 'video', '线上', 'completed', '技术面试', 'K8s经验丰富', 4, v_interviewer_id, NOW() - INTERVAL '6 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 何伟的面试
    SELECT id INTO v_talent_id FROM talents WHERE name = '何伟' OR email LIKE 'hewei%' LIMIT 1;
    SELECT id, title INTO v_job_id, v_job_title FROM jobs WHERE title LIKE '%AI%' OR title LIKE '%算法%' LIMIT 1;
    IF v_talent_id IS NOT NULL AND v_job_id IS NOT NULL THEN
        INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, duration, interviewer_id, interviewer, method, location, status, notes, feedback, rating, created_by, created_at)
        VALUES (v_talent_id, '何伟', v_job_id, COALESCE(v_job_title, 'AI算法工程师'), 'initial', TO_CHAR(NOW() - INTERVAL '5 days', 'YYYY-MM-DD'), '10:00', 90, v_interviewer_id, v_interviewer_name, 'onsite', '会议室A', 'completed', '算法面试', 'NLP能力强，有论文发表', 5, v_interviewer_id, NOW() - INTERVAL '7 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 王磊的面试
    SELECT id INTO v_talent_id FROM talents WHERE name = '王磊' OR email LIKE 'wanglei%' LIMIT 1;
    SELECT id, title INTO v_job_id, v_job_title FROM jobs WHERE title LIKE '%Java%' LIMIT 1;
    IF v_talent_id IS NOT NULL AND v_job_id IS NOT NULL THEN
        INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, duration, interviewer_id, interviewer, method, location, status, notes, created_by, created_at)
        VALUES (v_talent_id, '王磊', v_job_id, COALESCE(v_job_title, 'Java开发工程师'), 'initial', TO_CHAR(NOW() + INTERVAL '3 days', 'YYYY-MM-DD'), '14:00', 60, v_interviewer_id, v_interviewer_name, 'video', '线上', 'scheduled', '技术面试', v_interviewer_id, NOW() - INTERVAL '2 days')
        ON CONFLICT DO NOTHING;
    END IF;

END $$;
