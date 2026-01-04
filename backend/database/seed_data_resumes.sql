-- 追加简历数据 (resumes表)
-- 字段: talent_id, job_id, file_path, file_name, status, match_score, parse_result

DO $$
DECLARE
    v_talent_id INTEGER;
    v_job_id INTEGER;
BEGIN
    -- 张伟的简历
    SELECT id INTO v_talent_id FROM talents WHERE email = 'zhangwei_dev@example.com' LIMIT 1;
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%Go%' LIMIT 1;
    IF v_talent_id IS NOT NULL THEN
        INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result, created_at, updated_at)
        VALUES (v_talent_id, v_job_id, '/uploads/resumes/zhangwei_resume.pdf', '张伟_简历.pdf', 'reviewing', 92,
        '{"name":"张伟","experience":6,"education":"本科","skills":["Go","Docker","Kubernetes","Redis","MySQL"]}',
        NOW() - INTERVAL '15 days', NOW() - INTERVAL '10 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 李强的简历
    SELECT id INTO v_talent_id FROM talents WHERE email = 'liqiang_dev@example.com' LIMIT 1;
    IF v_talent_id IS NOT NULL THEN
        INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result, created_at, updated_at)
        VALUES (v_talent_id, v_job_id, '/uploads/resumes/liqiang_resume.pdf', '李强_简历.pdf', 'pending', 85,
        '{"name":"李强","experience":4,"education":"硕士","skills":["Go","Python","PostgreSQL","Redis"]}',
        NOW() - INTERVAL '12 days', NOW() - INTERVAL '8 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 陈浩的简历
    SELECT id INTO v_talent_id FROM talents WHERE email = 'chenhao_dev@example.com' LIMIT 1;
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%架构%' LIMIT 1;
    IF v_talent_id IS NOT NULL THEN
        INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result, created_at, updated_at)
        VALUES (v_talent_id, v_job_id, '/uploads/resumes/chenhao_resume.pdf', '陈浩_简历.pdf', 'offered', 96,
        '{"name":"陈浩","experience":8,"education":"硕士","skills":["Go","Rust","C++","分布式系统"]}',
        NOW() - INTERVAL '20 days', NOW() - INTERVAL '15 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 赵雪的简历
    SELECT id INTO v_talent_id FROM talents WHERE email = 'zhaoxue_dev@example.com' LIMIT 1;
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%前端%' LIMIT 1;
    IF v_talent_id IS NOT NULL THEN
        INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result, created_at, updated_at)
        VALUES (v_talent_id, v_job_id, '/uploads/resumes/zhaoxue_resume.pdf', '赵雪_简历.pdf', 'hired', 94,
        '{"name":"赵雪","experience":5,"education":"本科","skills":["Vue","TypeScript","React","Webpack"]}',
        NOW() - INTERVAL '18 days', NOW() - INTERVAL '12 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 孙婷的简历
    SELECT id INTO v_talent_id FROM talents WHERE email = 'sunting_dev@example.com' LIMIT 1;
    IF v_talent_id IS NOT NULL THEN
        INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result, created_at, updated_at)
        VALUES (v_talent_id, v_job_id, '/uploads/resumes/sunting_resume.pdf', '孙婷_简历.pdf', 'interviewed', 88,
        '{"name":"孙婷","experience":4,"education":"本科","skills":["React","TypeScript","Redux"]}',
        NOW() - INTERVAL '14 days', NOW() - INTERVAL '10 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 马超的简历
    SELECT id INTO v_talent_id FROM talents WHERE email = 'machao_dev@example.com' LIMIT 1;
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%DevOps%' OR title LIKE '%运维%' LIMIT 1;
    IF v_talent_id IS NOT NULL THEN
        INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result, created_at, updated_at)
        VALUES (v_talent_id, v_job_id, '/uploads/resumes/machao_resume.pdf', '马超_简历.pdf', 'interviewed', 90,
        '{"name":"马超","experience":5,"education":"本科","skills":["Docker","Kubernetes","Jenkins"]}',
        NOW() - INTERVAL '13 days', NOW() - INTERVAL '10 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 何伟的简历
    SELECT id INTO v_talent_id FROM talents WHERE email = 'hewei_dev@example.com' LIMIT 1;
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%AI%' OR title LIKE '%算法%' LIMIT 1;
    IF v_talent_id IS NOT NULL THEN
        INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result, created_at, updated_at)
        VALUES (v_talent_id, v_job_id, '/uploads/resumes/hewei_resume.pdf', '何伟_简历.pdf', 'interviewed', 95,
        '{"name":"何伟","experience":4,"education":"博士","skills":["Python","PyTorch","NLP"]}',
        NOW() - INTERVAL '22 days', NOW() - INTERVAL '18 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 杨帆的简历
    SELECT id INTO v_talent_id FROM talents WHERE email = 'yangfan_dev@example.com' LIMIT 1;
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%数据分析%' LIMIT 1;
    IF v_talent_id IS NOT NULL THEN
        INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result, created_at, updated_at)
        VALUES (v_talent_id, v_job_id, '/uploads/resumes/yangfan_resume.pdf', '杨帆_简历.pdf', 'hired', 87,
        '{"name":"杨帆","experience":3,"education":"硕士","skills":["Python","SQL","Tableau"]}',
        NOW() - INTERVAL '7 days', NOW() - INTERVAL '5 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 韩冰的简历
    SELECT id INTO v_talent_id FROM talents WHERE email = 'hanbing_dev@example.com' LIMIT 1;
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%测试%' LIMIT 1;
    IF v_talent_id IS NOT NULL THEN
        INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result, created_at, updated_at)
        VALUES (v_talent_id, v_job_id, '/uploads/resumes/hanbing_resume.pdf', '韩冰_简历.pdf', 'offered', 89,
        '{"name":"韩冰","experience":4,"education":"本科","skills":["Selenium","Python","自动化测试"]}',
        NOW() - INTERVAL '15 days', NOW() - INTERVAL '12 days')
        ON CONFLICT DO NOTHING;
    END IF;

    -- 冯刚的简历
    SELECT id INTO v_talent_id FROM talents WHERE email = 'fenggang_dev@example.com' LIMIT 1;
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%Go%' LIMIT 1;
    IF v_talent_id IS NOT NULL THEN
        INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result, created_at, updated_at)
        VALUES (v_talent_id, v_job_id, '/uploads/resumes/fenggang_resume.pdf', '冯刚_简历.pdf', 'pending', 93,
        '{"name":"冯刚","experience":7,"education":"硕士","skills":["Go","gRPC","etcd","分布式系统"]}',
        NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days')
        ON CONFLICT DO NOTHING;
    END IF;

END $$;
