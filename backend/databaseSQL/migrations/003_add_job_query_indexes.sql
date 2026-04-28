-- Migration: 003_add_job_query_indexes.sql
-- Description: Add high-value indexes for job listing and application count queries
-- Requirements: job-service list performance, applicant count aggregation
--
-- This migration is idempotent - safe to run multiple times
-- All CREATE INDEX statements use IF NOT EXISTS to prevent errors on re-run

-- ============================================================================
-- 职位列表高频查询索引 (Job List Query Indexes)
-- ============================================================================
-- Used when:
-- 1. 列表按 created_at 倒序分页
-- 2. status 过滤 + created_at 倒序分页
-- 3. type / level / education 精确过滤

CREATE INDEX IF NOT EXISTS idx_jobs_active_created_at
    ON jobs(created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_jobs_active_status_created_at
    ON jobs(status, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_jobs_type
    ON jobs(type);

CREATE INDEX IF NOT EXISTS idx_jobs_level
    ON jobs(level);

CREATE INDEX IF NOT EXISTS idx_jobs_education
    ON jobs(education);

-- ============================================================================
-- 申请人数统计索引 (Application Count Index)
-- ============================================================================
-- Used when:
-- 1. 按 job_id 聚合统计申请人数
-- 2. 查询职位申请列表时过滤已软删除数据

CREATE INDEX IF NOT EXISTS idx_applications_job_active
    ON applications(job_id)
    WHERE deleted_at IS NULL;

-- ============================================================================
-- Migration complete
-- ============================================================================
