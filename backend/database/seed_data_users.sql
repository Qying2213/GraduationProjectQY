-- 追加用户数据 (users表)
-- 密码都是 123456 的 bcrypt 哈希
-- 使用 ON CONFLICT DO NOTHING 避免与现有数据冲突

INSERT INTO users (username, email, password, role, status, created_at, updated_at) VALUES
-- 新增HR用户
('hr_wang', 'wang.hr@talentplatform.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGlmPzPdZGNqQxLUotkvJO', 'hr', 'active', NOW(), NOW()),
('hr_chen', 'chen.hr@talentplatform.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGlmPzPdZGNqQxLUotkvJO', 'hr', 'active', NOW(), NOW()),
('hr_liu', 'liu.hr@talentplatform.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGlmPzPdZGNqQxLUotkvJO', 'hr', 'active', NOW(), NOW()),

-- 新增面试官
('interviewer_zhao', 'zhao.tech@talentplatform.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGlmPzPdZGNqQxLUotkvJO', 'interviewer', 'active', NOW(), NOW()),
('interviewer_sun', 'sun.tech@talentplatform.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGlmPzPdZGNqQxLUotkvJO', 'interviewer', 'active', NOW(), NOW()),
('interviewer_qian', 'qian.tech@talentplatform.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGlmPzPdZGNqQxLUotkvJO', 'interviewer', 'active', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;
