-- 智能人才运营平台数据库Schema
-- PostgreSQL

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    role VARCHAR(20) DEFAULT 'user',
    avatar VARCHAR(255),
    department VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 职位表
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    department VARCHAR(100),
    location VARCHAR(100),
    salary VARCHAR(50),
    level VARCHAR(20),
    type VARCHAR(20) DEFAULT 'full-time',
    description TEXT,
    requirements TEXT,
    skills TEXT,
    headcount INTEGER DEFAULT 1,
    status VARCHAR(20) DEFAULT 'open',
    publisher_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 人才表
CREATE TABLE IF NOT EXISTS talents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(100),
    gender VARCHAR(10),
    age INTEGER,
    education VARCHAR(20),
    school VARCHAR(100),
    major VARCHAR(100),
    experience INTEGER DEFAULT 0,
    skills TEXT,
    location VARCHAR(100),
    salary VARCHAR(50),
    current_company VARCHAR(100),
    current_position VARCHAR(100),
    source VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    resume_url VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 简历表
CREATE TABLE IF NOT EXISTS resumes (
    id SERIAL PRIMARY KEY,
    talent_id INTEGER REFERENCES talents(id),
    file_name VARCHAR(255),
    file_path VARCHAR(500),
    file_size INTEGER,
    file_type VARCHAR(50),
    parsed_data TEXT,
    match_score INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 申请表
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    talent_id INTEGER REFERENCES talents(id),
    job_id INTEGER REFERENCES jobs(id),
    resume_id INTEGER REFERENCES resumes(id),
    status VARCHAR(20) DEFAULT 'pending',
    stage VARCHAR(50) DEFAULT 'resume_screening',
    match_score INTEGER,
    notes TEXT,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 面试表
CREATE TABLE IF NOT EXISTS interviews (
    id SERIAL PRIMARY KEY,
    application_id INTEGER REFERENCES applications(id),
    talent_id INTEGER REFERENCES talents(id),
    job_id INTEGER REFERENCES jobs(id),
    interviewer VARCHAR(100),
    interviewer_id INTEGER REFERENCES users(id),
    interview_type VARCHAR(50),
    scheduled_time TIMESTAMP,
    duration INTEGER DEFAULT 60,
    location VARCHAR(255),
    meeting_link VARCHAR(255),
    status VARCHAR(20) DEFAULT 'scheduled',
    result VARCHAR(20),
    feedback TEXT,
    score INTEGER,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 面试反馈表
CREATE TABLE IF NOT EXISTS interview_feedbacks (
    id SERIAL PRIMARY KEY,
    interview_id INTEGER REFERENCES interviews(id),
    interviewer_id INTEGER REFERENCES users(id),
    overall_score INTEGER,
    technical_score INTEGER,
    communication_score INTEGER,
    culture_fit_score INTEGER,
    strengths TEXT,
    weaknesses TEXT,
    recommendation VARCHAR(50),
    comments TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
CREATE INDEX IF NOT EXISTS idx_interviews_scheduled ON interviews(scheduled_time);
CREATE INDEX IF NOT EXISTS idx_messages_receiver ON messages(receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_read ON messages(is_read);
CREATE INDEX IF NOT EXISTS idx_evaluation_results_status ON evaluation_results(status);
CREATE INDEX IF NOT EXISTS idx_evaluation_results_match_level ON evaluation_results(match_level);
