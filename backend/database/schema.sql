-- 智能人才运营平台数据库Schema
-- PostgreSQL

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'candidate',
    avatar VARCHAR(255),
    phone VARCHAR(20),
    department VARCHAR(50),
    position VARCHAR(50),
    real_name VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 职位表
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    requirements TEXT[],
    salary VARCHAR(100),
    location VARCHAR(100),
    type VARCHAR(20) DEFAULT 'full-time',
    status VARCHAR(20) DEFAULT 'open',
    created_by INTEGER REFERENCES users(id),
    department VARCHAR(100),
    level VARCHAR(50),
    skills TEXT[],
    benefits TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 人才表
CREATE TABLE IF NOT EXISTS talents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    skills TEXT[],
    experience INTEGER DEFAULT 0,
    education VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    tags TEXT[],
    user_id INTEGER REFERENCES users(id),
    location VARCHAR(100),
    salary VARCHAR(50),
    summary TEXT,
    gender VARCHAR(10),
    age INTEGER,
    current_company VARCHAR(100),
    current_position VARCHAR(100),
    source VARCHAR(50),
    resume_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 简历表
CREATE TABLE IF NOT EXISTS resumes (
    id SERIAL PRIMARY KEY,
    talent_id INTEGER REFERENCES talents(id),
    job_id INTEGER REFERENCES jobs(id),
    file_name VARCHAR(255),
    file_path VARCHAR(500),
    file_url VARCHAR(500),
    file_size BIGINT,
    file_type VARCHAR(50),
    parsed_data TEXT,
    extracted_text TEXT,
    match_score INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 申请表
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    talent_id INTEGER REFERENCES talents(id),
    job_id INTEGER REFERENCES jobs(id),
    resume_id INTEGER REFERENCES resumes(id),
    status VARCHAR(20) DEFAULT 'pending',
    cover_letter TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 面试表
CREATE TABLE IF NOT EXISTS interviews (
    id SERIAL PRIMARY KEY,
    candidate_id INTEGER REFERENCES talents(id),
    candidate_name VARCHAR(100) NOT NULL,
    position_id INTEGER REFERENCES jobs(id),
    position VARCHAR(200) NOT NULL,
    type VARCHAR(20) DEFAULT 'initial',
    date VARCHAR(20) NOT NULL,
    time VARCHAR(10) NOT NULL,
    duration INTEGER DEFAULT 60,
    interviewer_id INTEGER REFERENCES users(id),
    interviewer VARCHAR(100) NOT NULL,
    method VARCHAR(20) DEFAULT 'onsite',
    location VARCHAR(500),
    status VARCHAR(20) DEFAULT 'scheduled',
    notes TEXT,
    feedback TEXT,
    rating INTEGER DEFAULT 0,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 面试反馈表
CREATE TABLE IF NOT EXISTS interview_feedbacks (
    id SERIAL PRIMARY KEY,
    interview_id INTEGER REFERENCES interviews(id),
    interviewer_id INTEGER REFERENCES users(id),
    rating INTEGER NOT NULL,
    strengths TEXT,
    weaknesses TEXT,
    comments TEXT,
    recommendation VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 消息表
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER REFERENCES users(id),
    receiver_id INTEGER REFERENCES users(id),
    title VARCHAR(200),
    content TEXT,
    type VARCHAR(50) DEFAULT 'notification',
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- AI评估结果表
CREATE TABLE IF NOT EXISTS evaluation_results (
    id SERIAL PRIMARY KEY,
    resume_id INTEGER REFERENCES resumes(id),
    talent_id INTEGER REFERENCES talents(id),
    job_id INTEGER REFERENCES jobs(id),
    resume_name VARCHAR(255),
    resume_file VARCHAR(500),
    parsed_name VARCHAR(100),
    parsed_phone VARCHAR(20),
    parsed_email VARCHAR(100),
    parsed_education VARCHAR(50),
    parsed_experience VARCHAR(50),
    parsed_location VARCHAR(100),
    parsed_skills TEXT,
    match_score DECIMAL(5,2),
    match_level VARCHAR(20),
    match_details TEXT,
    risk_score DECIMAL(5,2) DEFAULT 0,
    risk_items TEXT,
    status VARCHAR(20) DEFAULT 'completed',
    eval_type VARCHAR(50),
    report_summary TEXT,
    report_recommendation TEXT,
    report_strengths TEXT,
    report_gaps TEXT,
    report_dimensions TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_department ON jobs(department);
CREATE INDEX IF NOT EXISTS idx_talents_status ON talents(status);
CREATE INDEX IF NOT EXISTS idx_talents_location ON talents(location);
CREATE INDEX IF NOT EXISTS idx_applications_status ON applications(status);
CREATE INDEX IF NOT EXISTS idx_applications_talent ON applications(talent_id);
CREATE INDEX IF NOT EXISTS idx_applications_job ON applications(job_id);
CREATE INDEX IF NOT EXISTS idx_interviews_status ON interviews(status);
CREATE INDEX IF NOT EXISTS idx_interviews_date ON interviews(date);
CREATE INDEX IF NOT EXISTS idx_interviews_candidate ON interviews(candidate_id);
CREATE INDEX IF NOT EXISTS idx_interviews_interviewer ON interviews(interviewer_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver ON messages(receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_read ON messages(is_read);
CREATE INDEX IF NOT EXISTS idx_evaluation_results_status ON evaluation_results(status);
CREATE INDEX IF NOT EXISTS idx_evaluation_results_match_level ON evaluation_results(match_level);
