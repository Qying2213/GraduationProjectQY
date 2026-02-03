-- Migration: 002_create_online_resumes_table.sql
-- Description: Create online_resumes table for online resume editing functionality
-- Requirements: 4.1 (Display resume information), 4.3 (Persist updated information)
-- 
-- This migration is idempotent - safe to run multiple times
-- All CREATE statements use IF NOT EXISTS to prevent errors on re-run

-- ============================================================================
-- 在线简历表 (Online Resumes Table)
-- ============================================================================
-- Stores structured online resume data for candidates
-- Each user has one online resume (1:1 relationship)
-- Uses JSONB for flexible storage of work experience, education, and skills

CREATE TABLE IF NOT EXISTS online_resumes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL UNIQUE,  -- 一个用户只有一份在线简历
    talent_id INTEGER REFERENCES talents(id),              -- 可选关联人才档案
    
    -- 基本信息 (Basic Info)
    name VARCHAR(100),
    phone VARCHAR(20),
    email VARCHAR(100),
    location VARCHAR(100),
    avatar VARCHAR(500),
    gender VARCHAR(10),
    age INTEGER,
    summary TEXT,
    
    -- 结构化数据 (Structured Data - stored as JSONB)
    work_experience JSONB DEFAULT '[]'::jsonb,  -- 工作经历数组
    education JSONB DEFAULT '[]'::jsonb,        -- 教育经历数组
    skills JSONB DEFAULT '[]'::jsonb,           -- 技能数组
    
    -- 状态
    is_complete BOOLEAN DEFAULT FALSE,          -- 简历是否完整
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ============================================================================
-- 表注释 (Table Comments)
-- ============================================================================

COMMENT ON TABLE online_resumes IS '在线简历表 - 存储用户在线编辑的简历数据';
COMMENT ON COLUMN online_resumes.user_id IS '关联用户ID，一个用户只有一份在线简历';
COMMENT ON COLUMN online_resumes.talent_id IS '可选关联人才档案ID';
COMMENT ON COLUMN online_resumes.name IS '姓名（必填）';
COMMENT ON COLUMN online_resumes.phone IS '手机号（必填）';
COMMENT ON COLUMN online_resumes.email IS '邮箱（必填）';
COMMENT ON COLUMN online_resumes.location IS '所在地';
COMMENT ON COLUMN online_resumes.avatar IS '头像URL';
COMMENT ON COLUMN online_resumes.gender IS '性别';
COMMENT ON COLUMN online_resumes.age IS '年龄';
COMMENT ON COLUMN online_resumes.summary IS '个人简介';
COMMENT ON COLUMN online_resumes.work_experience IS '工作经历JSON数组';
COMMENT ON COLUMN online_resumes.education IS '教育经历JSON数组';
COMMENT ON COLUMN online_resumes.skills IS '技能JSON数组';
COMMENT ON COLUMN online_resumes.is_complete IS '简历是否完整（基本信息+经历+技能）';

-- ============================================================================
-- 索引 (Indexes)
-- ============================================================================

-- Index for looking up resume by user_id (already unique, but explicit index)
CREATE INDEX IF NOT EXISTS idx_online_resumes_user_id 
    ON online_resumes(user_id);

-- Index for looking up resume by talent_id
CREATE INDEX IF NOT EXISTS idx_online_resumes_talent_id 
    ON online_resumes(talent_id);

-- Index for soft delete queries
CREATE INDEX IF NOT EXISTS idx_online_resumes_deleted_at 
    ON online_resumes(deleted_at);

-- ============================================================================
-- JSONB 数据结构说明 (JSONB Data Structure Documentation)
-- ============================================================================
-- 
-- work_experience 结构:
-- [
--   {
--     "company": "公司名称",
--     "position": "职位",
--     "start_date": "2020-01",
--     "end_date": "2023-06",  // 空表示至今
--     "is_current": false,
--     "description": "工作描述",
--     "location": "工作地点"
--   }
-- ]
--
-- education 结构:
-- [
--   {
--     "school": "学校名称",
--     "degree": "本科/硕士/博士",
--     "major": "专业",
--     "start_date": "2016-09",
--     "end_date": "2020-06",
--     "is_current": false,
--     "gpa": "3.8",
--     "activities": "校园活动"
--   }
-- ]
--
-- skills 结构:
-- ["Go", "Python", "JavaScript", "PostgreSQL"]

-- ============================================================================
-- Migration complete
-- ============================================================================
