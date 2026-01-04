-- 追加申请数据 (applications表)
-- 字段: talent_id, job_id, resume_id, stage, status, source, notes

DO $$
DECLARE
    v_job_id INTEGER;
    v_talent_id INTEGER;
BEGIN
    -- 张伟申请Go开发
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%Go%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '张伟' OR email LIKE 'zhangwei%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'interview', 'active', '招聘网站', '技术面试通过', NOW() - INTERVAL '10 days', NOW() - INTERVAL '3 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 李强申请Go开发
    SELECT id INTO v_talent_id FROM talents WHERE name = '李强' OR email LIKE 'liqiang%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'screening', 'active', '内推', '简历已审核', NOW() - INTERVAL '8 days', NOW() - INTERVAL '5 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 陈浩申请架构师
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%架构%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '陈浩' OR email LIKE 'chenhao%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'offer', 'active', '猎头', '已发offer', NOW() - INTERVAL '15 days', NOW() - INTERVAL '2 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 赵雪申请前端
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%前端%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '赵雪' OR email LIKE 'zhaoxue%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'hired', 'active', '招聘网站', '已入职', NOW() - INTERVAL '20 days', NOW() - INTERVAL '5 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 孙婷申请前端
    SELECT id INTO v_talent_id FROM talents WHERE name = '孙婷' OR email LIKE 'sunting%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'interview', 'active', '校招', '二面中', NOW() - INTERVAL '12 days', NOW() - INTERVAL '4 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 马超申请DevOps
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%DevOps%' OR title LIKE '%运维%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '马超' OR email LIKE 'machao%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'interview', 'active', '招聘网站', '技术面试中', NOW() - INTERVAL '7 days', NOW() - INTERVAL '3 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 林峰申请云原生/SRE
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%云原生%' OR title LIKE '%SRE%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '林峰' OR email LIKE 'linfeng%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'offer', 'active', '猎头', '已发offer', NOW() - INTERVAL '12 days', NOW() - INTERVAL '2 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 何伟申请AI算法
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%AI%' OR title LIKE '%算法%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '何伟' OR email LIKE 'hewei%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'interview', 'active', '招聘网站', '技术面试中', NOW() - INTERVAL '9 days', NOW() - INTERVAL '4 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 杨帆申请数据分析
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%数据分析%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '杨帆' OR email LIKE 'yangfan%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'hired', 'active', '内推', '已入职', NOW() - INTERVAL '18 days', NOW() - INTERVAL '5 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 谢敏申请产品经理
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%产品%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '谢敏' OR email LIKE 'xiemin%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'interview', 'active', '招聘网站', '二面中', NOW() - INTERVAL '8 days', NOW() - INTERVAL '3 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 王磊申请Java
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%Java%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '王磊' OR email LIKE 'wanglei%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'interview', 'active', '招聘网站', '技术面试中', NOW() - INTERVAL '6 days', NOW() - INTERVAL '2 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 冯刚申请Go开发
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%Go%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '冯刚' OR email LIKE 'fenggang%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'applied', 'active', '招聘网站', NULL, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

    -- 董丽申请前端
    SELECT id INTO v_job_id FROM jobs WHERE title LIKE '%前端%' LIMIT 1;
    SELECT id INTO v_talent_id FROM talents WHERE name = '董丽' OR email LIKE 'dongli%' LIMIT 1;
    IF v_job_id IS NOT NULL AND v_talent_id IS NOT NULL THEN
        INSERT INTO applications (job_id, talent_id, stage, status, source, notes, created_at, updated_at)
        VALUES (v_job_id, v_talent_id, 'offer', 'active', '猎头', '已发offer', NOW() - INTERVAL '14 days', NOW() - INTERVAL '3 days')
        ON CONFLICT (talent_id, job_id) DO NOTHING;
    END IF;

END $$;
