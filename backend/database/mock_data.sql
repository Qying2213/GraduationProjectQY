-- 智能人才运营平台测试数据
-- PostgreSQL

-- 清空现有数据（按依赖顺序）
TRUNCATE TABLE evaluation_results, interview_feedbacks, interviews, applications, messages, resumes, talents, jobs, users RESTART IDENTITY CASCADE;

-- 插入用户数据
INSERT INTO users (username, password, email, phone, role, department, status) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'admin@company.com', '13800000001', 'admin', '技术部', 'active'),
('hr_li', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'hr_li@company.com', '13800000002', 'hr', '人力资源部', 'active'),
('hr_wang', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'hr_wang@company.com', '13800000003', 'hr', '人力资源部', 'active'),
('interviewer_zhang', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'zhang@company.com', '13800000004', 'interviewer', '技术部', 'active'),
('interviewer_liu', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'liu@company.com', '13800000005', 'interviewer', '产品部', 'active');

-- 插入职位数据
INSERT INTO jobs (title, department, location, salary, level, type, description, requirements, skills, headcount, status, publisher_id) VALUES
('高级Go开发工程师', '技术部', '北京', '30-50K', 'senior', 'full-time', '负责后端微服务架构设计与开发', '5年以上Go开发经验，熟悉微服务架构', 'Go,Docker,Kubernetes,微服务,PostgreSQL,Redis', 3, 'open', 1),
('前端架构师', '技术部', '上海', '35-55K', 'senior', 'full-time', '负责前端技术架构设计与团队管理', '8年以上前端开发经验，精通Vue/React', 'Vue,React,TypeScript,Webpack,Node.js', 1, 'open', 1),
('全栈开发工程师', '产品部', '深圳', '25-40K', 'mid', 'full-time', '负责产品功能开发与维护', '3年以上全栈开发经验', 'Go,Vue,PostgreSQL,Redis,Docker', 2, 'open', 1),
('AI算法工程师', '技术部', '北京', '40-60K', 'senior', 'full-time', '负责AI模型研发与优化', '硕士及以上学历，3年以上AI开发经验', 'Python,PyTorch,TensorFlow,机器学习,深度学习,NLP', 2, 'open', 1),
('产品经理', '产品部', '杭州', '25-35K', 'mid', 'full-time', '负责产品规划与需求分析', '3年以上产品经验，有B端产品经验优先', '产品设计,需求分析,项目管理,数据分析', 1, 'open', 2),
('DevOps工程师', '技术部', '北京', '30-45K', 'mid', 'full-time', '负责CI/CD流程建设与运维', '3年以上DevOps经验', 'Docker,Kubernetes,Jenkins,Ansible,Linux,Shell', 1, 'open', 1);

-- 插入人才数据
INSERT INTO talents (name, phone, email, gender, age, education, school, major, experience, skills, location, salary, current_company, current_position, source, status) VALUES
('张三', '13900000001', 'zhangsan@email.com', '男', 28, '本科', '北京大学', '计算机科学', 5, 'Go,Docker,Kubernetes,Redis,微服务,PostgreSQL', '北京', '30-40K', '字节跳动', '高级开发工程师', '猎聘网', 'active'),
('李四', '13900000002', 'lisi@email.com', '男', 30, '硕士', '清华大学', '软件工程', 7, 'Vue,React,TypeScript,Node.js,Webpack,前端架构', '上海', '35-50K', '阿里巴巴', '前端技术专家', 'BOSS直聘', 'active'),
('王五', '13900000003', 'wangwu@email.com', '男', 26, '本科', '浙江大学', '计算机科学', 4, 'Go,Vue,PostgreSQL,Docker,Redis,全栈开发', '深圳', '25-35K', '腾讯', '全栈工程师', '官网投递', 'active'),
('赵六', '13900000004', 'zhaoliu@email.com', '女', 29, '硕士', '上海交通大学', '人工智能', 4, 'Python,PyTorch,TensorFlow,机器学习,深度学习,NLP,CV', '北京', '40-55K', '百度', 'AI算法工程师', '内部推荐', 'active'),
('钱七', '13900000005', 'qianqi@email.com', '女', 27, '本科', '复旦大学', '产品设计', 4, '产品设计,需求分析,项目管理,数据分析,用户研究', '杭州', '25-30K', '网易', '产品经理', '猎聘网', 'active'),
('孙八', '13900000006', 'sunba@email.com', '男', 31, '本科', '华中科技大学', '计算机科学', 8, 'Docker,Kubernetes,Jenkins,Ansible,Linux,Shell,AWS', '北京', '35-45K', '美团', 'DevOps专家', 'BOSS直聘', 'active'),
('周九', '13900000007', 'zhoujiu@email.com', '男', 25, '本科', '南京大学', '软件工程', 3, 'Go,Python,MySQL,Redis,Docker', '南京', '20-30K', '苏宁', '后端开发', '官网投递', 'active'),
('吴十', '13900000008', 'wushi@email.com', '女', 28, '硕士', '中国科学技术大学', '计算机科学', 5, 'Java,Spring,MySQL,Redis,微服务,分布式', '合肥', '25-35K', '科大讯飞', '高级Java开发', '内部推荐', 'active');

-- 插入简历数据
INSERT INTO resumes (talent_id, file_name, file_path, file_size, file_type, status, match_score) VALUES
(1, '张三_简历.pdf', '/uploads/resumes/zhangsan.pdf', 1024000, 'application/pdf', 'parsed', 85),
(2, '李四_简历.pdf', '/uploads/resumes/lisi.pdf', 1536000, 'application/pdf', 'parsed', 90),
(3, '王五_简历.pdf', '/uploads/resumes/wangwu.pdf', 1280000, 'application/pdf', 'parsed', 78),
(4, '赵六_简历.pdf', '/uploads/resumes/zhaoliu.pdf', 2048000, 'application/pdf', 'parsed', 88),
(5, '钱七_简历.pdf', '/uploads/resumes/qianqi.pdf', 1024000, 'application/pdf', 'parsed', 82),
(6, '孙八_简历.pdf', '/uploads/resumes/sunba.pdf', 1536000, 'application/pdf', 'parsed', 86);

-- 插入申请数据
INSERT INTO applications (talent_id, job_id, resume_id, status, stage, match_score) VALUES
(1, 1, 1, 'interviewing', 'technical_interview', 85),
(2, 2, 2, 'interviewing', 'hr_interview', 90),
(3, 3, 3, 'pending', 'resume_screening', 78),
(4, 4, 4, 'offered', 'offer', 88),
(5, 5, 5, 'interviewing', 'technical_interview', 82),
(6, 6, 6, 'hired', 'onboarding', 86),
(1, 3, 1, 'rejected', 'resume_screening', 65),
(7, 1, NULL, 'pending', 'resume_screening', 72);

-- 插入面试数据
INSERT INTO interviews (application_id, talent_id, job_id, interviewer, interviewer_id, interview_type, scheduled_time, duration, location, status, result, score) VALUES
(1, 1, 1, '张工', 4, 'technical', CURRENT_TIMESTAMP + INTERVAL '1 day', 60, '会议室A', 'scheduled', NULL, NULL),
(2, 2, 2, '刘工', 5, 'hr', CURRENT_TIMESTAMP + INTERVAL '2 days', 45, '会议室B', 'scheduled', NULL, NULL),
(4, 4, 4, '张工', 4, 'technical', CURRENT_TIMESTAMP - INTERVAL '3 days', 60, '线上', 'completed', 'pass', 92),
(5, 5, 5, '刘工', 5, 'technical', CURRENT_TIMESTAMP + INTERVAL '3 days', 60, '会议室C', 'scheduled', NULL, NULL),
(6, 6, 6, '张工', 4, 'final', CURRENT_TIMESTAMP - INTERVAL '7 days', 60, '会议室A', 'completed', 'pass', 88);

-- 插入面试反馈
INSERT INTO interview_feedbacks (interview_id, interviewer_id, overall_score, technical_score, communication_score, culture_fit_score, strengths, weaknesses, recommendation, comments) VALUES
(3, 4, 92, 95, 88, 90, '技术能力强，AI算法经验丰富', '沟通表达可以更简洁', 'strong_hire', '非常优秀的候选人，建议尽快发offer'),
(5, 4, 88, 90, 85, 88, 'DevOps经验丰富，对K8s理解深入', '对新技术关注度可以提高', 'hire', '符合岗位要求，建议录用');

-- 插入消息数据
INSERT INTO messages (sender_id, receiver_id, title, content, type, is_read) VALUES
(1, 2, '新简历待审核', '有3份新简历等待您审核', 'notification', FALSE),
(1, 3, '面试安排提醒', '明天有2场面试需要安排', 'reminder', FALSE),
(4, 2, '面试反馈已提交', '张三的技术面试反馈已提交', 'notification', TRUE),
(2, 4, '面试安排通知', '请您明天下午2点面试候选人李四', 'notification', FALSE);

-- 插入评估结果数据
INSERT INTO evaluation_results (resume_id, talent_id, job_id, resume_name, parsed_name, parsed_education, parsed_experience, parsed_skills, match_score, match_level, status, eval_type, report_summary, report_recommendation) VALUES
(1, 1, 1, '张三_简历.pdf', '张三', '本科', '5年', '["Go","Docker","Kubernetes","Redis"]', 85.5, 'high', 'completed', 'ai_evaluate', '候选人与职位高度匹配，技术能力突出', '建议优先安排面试'),
(2, 2, 2, '李四_简历.pdf', '李四', '硕士', '7年', '["Vue","React","TypeScript"]', 90.2, 'high', 'completed', 'ai_evaluate', '候选人前端技术能力出色，架构经验丰富', '强烈推荐，建议尽快面试'),
(4, 4, 4, '赵六_简历.pdf', '赵六', '硕士', '4年', '["Python","PyTorch","机器学习"]', 88.0, 'high', 'completed', 'ai_evaluate', '候选人AI算法能力强，项目经验丰富', '建议录用');
