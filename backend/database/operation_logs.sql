-- 操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    username VARCHAR(100),
    action VARCHAR(50) NOT NULL,
    module VARCHAR(50) NOT NULL,
    content TEXT,
    ip VARCHAR(50),
    user_agent TEXT,
    status VARCHAR(20) DEFAULT 'success',
    details JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_operation_logs_user_id ON operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_action ON operation_logs(action);
CREATE INDEX IF NOT EXISTS idx_operation_logs_module ON operation_logs(module);
CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs(created_at);

-- 插入模拟数据
INSERT INTO operation_logs (user_id, username, action, module, content, ip, user_agent, status, details) VALUES
(1, 'admin', 'login', 'user', '用户登录系统', '192.168.1.100', 'Chrome/120.0.0.0', 'success', '{}'),
(1, 'admin', 'create', 'job', '创建了新职位：高级Go开发工程师', '192.168.1.100', 'Chrome/120.0.0.0', 'success', '{"job_id": 1}'),
(2, 'hr_zhang', 'create', 'talent', '创建了新人才记录：张三', '192.168.1.101', 'Chrome/120.0.0.0', 'success', '{"talent_id": 1}'),
(2, 'hr_zhang', 'update', 'interview', '更新了面试状态为已完成', '192.168.1.101', 'Chrome/120.0.0.0', 'success', '{"interview_id": 1}'),
(1, 'admin', 'export', 'talent', '导出了人才列表数据', '192.168.1.100', 'Chrome/120.0.0.0', 'success', '{"count": 50}'),
(3, 'interviewer_chen', 'update', 'interview', '提交了面试反馈', '192.168.1.102', 'Chrome/120.0.0.0', 'success', '{"interview_id": 2}'),
(2, 'hr_zhang', 'create', 'interview', '安排了面试：李四 - 前端工程师', '192.168.1.101', 'Chrome/120.0.0.0', 'success', '{"interview_id": 3}'),
(1, 'admin', 'delete', 'resume', '删除了简历记录 #123', '192.168.1.100', 'Chrome/120.0.0.0', 'success', '{"resume_id": 123}'),
(1, 'admin', 'update', 'system', '更新了系统配置', '192.168.1.100', 'Chrome/120.0.0.0', 'success', '{}'),
(2, 'hr_zhang', 'login', 'user', '用户登录系统', '192.168.1.101', 'Chrome/120.0.0.0', 'success', '{}');
