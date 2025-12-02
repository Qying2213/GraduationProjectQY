-- =====================================================
-- 人才运营平台 - 数据库表结构
-- PostgreSQL 14+
-- =====================================================

-- 启用 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================
-- 1. 用户表
-- =====================================================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'viewer',
    avatar VARCHAR(500),
    phone VARCHAR(20),
    department VARCHAR(50),
    position VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);

COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.role IS '角色: admin, hr_manager, recruiter, interviewer, viewer';
COMMENT ON COLUMN users.status IS '状态: active, inactive, suspended';

-- =====================================================
-- 2. 职位表
-- =====================================================
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    requirements TEXT[],
    salary VARCHAR(50),
    location VARCHAR(50),
    type VARCHAR(20) NOT NULL DEFAULT 'full-time',
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    created_by INTEGER REFERENCES users(id),
    department VARCHAR(50),
    level VARCHAR(20),
    skills TEXT[],
    benefits TEXT[],
    headcount INTEGER DEFAULT 1,
    urgent BOOLEAN DEFAULT FALSE,
    deadline DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_type ON jobs(type);
CREATE INDEX idx_jobs_location ON jobs(location);
CREATE INDEX idx_jobs_created_by ON jobs(created_by);

COMMENT ON TABLE jobs IS '职位表';
COMMENT ON COLUMN jobs.type IS '类型: full-time, part-time, contract, internship';
COMMENT ON COLUMN jobs.status IS '状态: open, closed, filled, paused';
COMMENT ON COLUMN jobs.level IS '级别: junior, mid, senior, expert, management';

-- =====================================================
-- 3. 人才表
-- =====================================================
CREATE TABLE IF NOT EXISTS talents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    skills TEXT[],
    experience INTEGER DEFAULT 0,
    education VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    tags TEXT[],
    user_id INTEGER REFERENCES users(id),
    location VARCHAR(50),
    salary VARCHAR(50),
    summary TEXT,
    gender VARCHAR(10),
    age INTEGER,
    current_company VARCHAR(100),
    current_position VARCHAR(100),
    source VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_talents_email ON talents(email);
CREATE INDEX idx_talents_status ON talents(status);
CREATE INDEX idx_talents_skills ON talents USING GIN(skills);
CREATE INDEX idx_talents_location ON talents(location);

COMMENT ON TABLE talents IS '人才表';
COMMENT ON COLUMN talents.status IS '状态: active, hired, pending, rejected';
COMMENT ON COLUMN talents.education IS '学历: 高中, 大专, 本科, 硕士, 博士';

-- =====================================================
-- 4. 简历表
-- =====================================================
CREATE TABLE IF NOT EXISTS resumes (
    id SERIAL PRIMARY KEY,
    talent_id INTEGER REFERENCES talents(id) ON DELETE CASCADE,
    job_id INTEGER REFERENCES jobs(id) ON DELETE SET NULL,
    file_path VARCHAR(500),
    file_name VARCHAR(200),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    match_score INTEGER DEFAULT 0,
    parse_result JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_resumes_talent_id ON resumes(talent_id);
CREATE INDEX idx_resumes_job_id ON resumes(job_id);
CREATE INDEX idx_resumes_status ON resumes(status);

COMMENT ON TABLE resumes IS '简历表';
COMMENT ON COLUMN resumes.status IS '状态: pending, reviewing, interviewed, offered, hired, rejected';
COMMENT ON COLUMN resumes.match_score IS 'AI匹配度分数 0-100';

-- =====================================================
-- 5. 面试表
-- =====================================================
CREATE TABLE IF NOT EXISTS interviews (
    id SERIAL PRIMARY KEY,
    candidate_id INTEGER NOT NULL,
    candidate_name VARCHAR(100) NOT NULL,
    position_id INTEGER NOT NULL,
    position VARCHAR(200) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'initial',
    date VARCHAR(20) NOT NULL,
    time VARCHAR(10) NOT NULL,
    duration INTEGER DEFAULT 60,
    interviewer_id INTEGER REFERENCES users(id),
    interviewer VARCHAR(100) NOT NULL,
    method VARCHAR(20) NOT NULL DEFAULT 'onsite',
    location VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    notes TEXT,
    feedback TEXT,
    rating INTEGER DEFAULT 0,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_interviews_candidate_id ON interviews(candidate_id);
CREATE INDEX idx_interviews_interviewer_id ON interviews(interviewer_id);
CREATE INDEX idx_interviews_date ON interviews(date);
CREATE INDEX idx_interviews_status ON interviews(status);

COMMENT ON TABLE interviews IS '面试表';
COMMENT ON COLUMN interviews.type IS '类型: initial(初试), second(复试), final(终面), hr(HR面)';
COMMENT ON COLUMN interviews.method IS '方式: onsite(现场), video(视频), phone(电话)';
COMMENT ON COLUMN interviews.status IS '状态: scheduled, completed, cancelled, no_show';
COMMENT ON COLUMN interviews.rating IS '评分: 1-5';

-- =====================================================
-- 6. 面试反馈表
-- =====================================================
CREATE TABLE IF NOT EXISTS interview_feedbacks (
    id SERIAL PRIMARY KEY,
    interview_id INTEGER REFERENCES interviews(id) ON DELETE CASCADE,
    interviewer_id INTEGER REFERENCES users(id),
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    strengths TEXT,
    weaknesses TEXT,
    comments TEXT,
    recommendation VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_interview_feedbacks_interview_id ON interview_feedbacks(interview_id);

COMMENT ON TABLE interview_feedbacks IS '面试反馈表';
COMMENT ON COLUMN interview_feedbacks.recommendation IS '建议: pass, fail, pending';

-- =====================================================
-- 7. 消息表
-- =====================================================
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER REFERENCES users(id),
    receiver_id INTEGER REFERENCES users(id) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'system',
    title VARCHAR(200) NOT NULL,
    content TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_is_read ON messages(is_read);
CREATE INDEX idx_messages_type ON messages(type);

COMMENT ON TABLE messages IS '消息表';
COMMENT ON COLUMN messages.type IS '类型: system, interview, feedback, offer, reminder, chat';

-- =====================================================
-- 8. 应聘记录表
-- =====================================================
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    talent_id INTEGER REFERENCES talents(id) ON DELETE CASCADE,
    job_id INTEGER REFERENCES jobs(id) ON DELETE CASCADE,
    resume_id INTEGER REFERENCES resumes(id),
    stage VARCHAR(50) NOT NULL DEFAULT 'applied',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    source VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(talent_id, job_id)
);

CREATE INDEX idx_applications_talent_id ON applications(talent_id);
CREATE INDEX idx_applications_job_id ON applications(job_id);
CREATE INDEX idx_applications_stage ON applications(stage);

COMMENT ON TABLE applications IS '应聘记录表';
COMMENT ON COLUMN applications.stage IS '阶段: applied, screening, interview, offer, hired, rejected';

-- =====================================================
-- 9. 角色权限表（可选，用于RBAC）
-- =====================================================
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    permissions TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE roles IS '角色表';

-- 插入预设角色
INSERT INTO roles (name, code, description, permissions) VALUES
('超级管理员', 'admin', '拥有系统所有权限', ARRAY['*']),
('HR主管', 'hr_manager', '负责招聘流程管理', ARRAY['talent:*', 'job:*', 'resume:*', 'interview:*', 'message:*']),
('招聘专员', 'recruiter', '负责日常招聘工作', ARRAY['talent:view', 'talent:create', 'talent:edit', 'job:view', 'resume:*', 'interview:*']),
('面试官', 'interviewer', '参与面试评估', ARRAY['talent:view', 'job:view', 'interview:view', 'interview:feedback']),
('只读用户', 'viewer', '只能查看数据', ARRAY['talent:view', 'job:view', 'resume:view', 'interview:view'])
ON CONFLICT (code) DO NOTHING;

-- =====================================================
-- 10. 操作日志表（可选）
-- =====================================================
CREATE TABLE IF NOT EXISTS operation_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id INTEGER,
    details JSONB,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_operation_logs_user_id ON operation_logs(user_id);
CREATE INDEX idx_operation_logs_action ON operation_logs(action);
CREATE INDEX idx_operation_logs_created_at ON operation_logs(created_at);

COMMENT ON TABLE operation_logs IS '操作日志表';

-- =====================================================
-- 触发器：自动更新 updated_at
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为各表添加触发器
DO $$
DECLARE
    t text;
BEGIN
    FOR t IN
        SELECT table_name
        FROM information_schema.columns
        WHERE column_name = 'updated_at'
        AND table_schema = 'public'
    LOOP
        EXECUTE format('
            DROP TRIGGER IF EXISTS update_%I_updated_at ON %I;
            CREATE TRIGGER update_%I_updated_at
                BEFORE UPDATE ON %I
                FOR EACH ROW
                EXECUTE FUNCTION update_updated_at_column();
        ', t, t, t, t);
    END LOOP;
END $$;

-- =====================================================
-- 完成
-- =====================================================
SELECT 'Database schema created successfully!' as status;
