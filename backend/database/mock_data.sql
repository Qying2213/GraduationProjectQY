-- 智能人才运营平台 - 模拟数据
-- 执行前请确保已创建数据库结构 (schema.sql)

-- 清空现有数据（按外键依赖顺序）
TRUNCATE TABLE interview_feedbacks, interviews, applications, resumes, messages, operation_logs, talents, jobs, users, roles RESTART IDENTITY CASCADE;

-- =====================================================
-- 1. 角色数据
-- =====================================================
INSERT INTO roles (name, code, description, permissions) VALUES
('超级管理员', 'admin', '拥有系统所有权限', ARRAY['*']),
('HR主管', 'hr_manager', '负责招聘流程管理', ARRAY['talent:*', 'job:*', 'resume:*', 'interview:*', 'message:*']),
('招聘专员', 'recruiter', '负责日常招聘工作', ARRAY['talent:view', 'talent:create', 'talent:edit', 'job:view', 'resume:*', 'interview:*']),
('面试官', 'interviewer', '参与面试评估', ARRAY['talent:view', 'job:view', 'interview:view', 'interview:feedback']),
('只读用户', 'viewer', '只能查看数据', ARRAY['talent:view', 'job:view', 'resume:view', 'interview:view']);

-- =====================================================
-- 2. 用户数据 (密码都是 password123，使用 bcrypt 加密)
-- =====================================================
INSERT INTO users (username, email, password, role, avatar, phone, department, position, status, real_name) VALUES
('admin', 'admin@company.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGqFKOJGJPHCYWJqFh0Iq2', 'admin', 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin', '13800000001', '技术部', '系统管理员', 'active', '管理员'),
('hr_zhang', 'zhang.hr@company.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGqFKOJGJPHCYWJqFh0Iq2', 'hr_manager', 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang', '13800000002', '人力资源部', 'HR总监', 'active', '张华'),
('hr_li', 'li.hr@company.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGqFKOJGJPHCYWJqFh0Iq2', 'recruiter', 'https://api.dicebear.com/7.x/avataaars/svg?seed=li', '13800000003', '人力资源部', '招聘专员', 'active', '李明'),
('hr_wang', 'wang.hr@company.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGqFKOJGJPHCYWJqFh0Iq2', 'recruiter', 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang', '13800000004', '人力资源部', '招聘专员', 'active', '王芳'),
('tech_chen', 'chen.tech@company.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGqFKOJGJPHCYWJqFh0Iq2', 'interviewer', 'https://api.dicebear.com/7.x/avataaars/svg?seed=chen', '13800000005', '技术部', '技术总监', 'active', '陈强'),
('tech_liu', 'liu.tech@company.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGqFKOJGJPHCYWJqFh0Iq2', 'interviewer', 'https://api.dicebear.com/7.x/avataaars/svg?seed=liu', '13800000006', '技术部', '前端负责人', 'active', '刘洋'),
('tech_zhao', 'zhao.tech@company.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGqFKOJGJPHCYWJqFh0Iq2', 'interviewer', 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhao', '13800000007', '技术部', '后端负责人', 'active', '赵磊'),
('product_sun', 'sun.product@company.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGqFKOJGJPHCYWJqFh0Iq2', 'interviewer', 'https://api.dicebear.com/7.x/avataaars/svg?seed=sun', '13800000008', '产品部', '产品总监', 'active', '孙婷'),
('viewer_test', 'viewer@company.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQlLBgXGqFKOJGJPHCYWJqFh0Iq2', 'viewer', 'https://api.dicebear.com/7.x/avataaars/svg?seed=viewer', '13800000009', '运营部', '运营专员', 'active', '周杰');

-- =====================================================
-- 3. 职位数据
-- =====================================================
INSERT INTO jobs (title, description, requirements, salary, location, type, status, created_by, department, level, skills, benefits, headcount, urgent, deadline) VALUES
('高级前端工程师', '负责公司核心产品的前端架构设计和开发工作，主导前端技术选型，优化前端性能，提升用户体验。', 
 ARRAY['5年以上前端开发经验', '精通Vue3/React等主流框架', '熟悉TypeScript', '有大型项目架构经验'], 
 '35-50K·15薪', '北京', 'full-time', 'open', 5, '技术部', 'senior', 
 ARRAY['Vue3', 'TypeScript', 'React', 'Webpack', 'Node.js'], 
 ARRAY['五险一金', '年终奖', '股票期权', '带薪年假', '弹性工作'], 2, false, '2025-03-01'),

('后端开发工程师（Go）', '参与公司微服务架构的设计和开发，负责高并发、高可用系统的建设。', 
 ARRAY['3年以上后端开发经验', '精通Go语言', '熟悉微服务架构', '了解Kubernetes'], 
 '30-45K·14薪', '北京', 'full-time', 'open', 5, '技术部', 'mid', 
 ARRAY['Go', 'gRPC', 'MySQL', 'Redis', 'Kubernetes', 'Docker'], 
 ARRAY['五险一金', '年终奖', '带薪年假', '免费三餐'], 3, true, '2025-02-15'),

('产品经理', '负责公司B端产品的规划和设计，深入理解用户需求，输出高质量的产品方案。', 
 ARRAY['3年以上B端产品经验', '熟悉互联网产品开发流程', '优秀的文档能力', '数据分析能力'], 
 '25-40K·13薪', '上海', 'full-time', 'open', 8, '产品部', 'mid', 
 ARRAY['Axure', 'Figma', 'SQL', '数据分析', '用户研究'], 
 ARRAY['五险一金', '年终奖', '带薪年假', '弹性工作'], 1, false, NULL),

('UI设计师', '负责公司产品的视觉设计和交互设计，打造优质的用户体验。', 
 ARRAY['3年以上UI设计经验', '精通Figma/Sketch', '有B端产品设计经验', '了解前端开发基础'], 
 '20-35K·13薪', '深圳', 'full-time', 'open', 8, '设计部', 'mid', 
 ARRAY['Figma', 'Sketch', 'Adobe XD', 'Photoshop', 'Illustrator'], 
 ARRAY['五险一金', '年终奖', '带薪年假', '设计培训'], 1, false, NULL),

('数据分析师', '负责业务数据的分析和挖掘，为产品和运营决策提供数据支持。', 
 ARRAY['2年以上数据分析经验', '精通SQL和Python', '熟悉常用BI工具', '良好的数据敏感度'], 
 '20-30K·13薪', '杭州', 'full-time', 'open', 8, '数据部', 'junior', 
 ARRAY['Python', 'SQL', 'Tableau', 'Excel', '统计学'], 
 ARRAY['五险一金', '年终奖', '带薪年假', '培训机会'], 2, false, NULL),

('DevOps工程师', '负责公司CI/CD流程建设和运维自动化，保障系统稳定运行。', 
 ARRAY['3年以上运维经验', '精通Linux系统', '熟悉Docker/K8s', '有CI/CD搭建经验'], 
 '25-40K·14薪', '北京', 'full-time', 'open', 5, '技术部', 'mid', 
 ARRAY['Linux', 'Docker', 'Kubernetes', 'Jenkins', 'Ansible', 'Python'], 
 ARRAY['五险一金', '年终奖', '带薪年假', '技术培训'], 1, false, NULL),

('测试工程师', '负责产品质量保障，制定测试计划和用例，执行功能和性能测试。', 
 ARRAY['2年以上测试经验', '熟悉测试流程和方法', '会自动化测试', '了解性能测试'], 
 '15-25K·13薪', '广州', 'full-time', 'open', 5, '技术部', 'junior', 
 ARRAY['Selenium', 'JMeter', 'Postman', 'Python', 'SQL'], 
 ARRAY['五险一金', '年终奖', '带薪年假'], 2, false, NULL),

('Java开发工程师', '参与公司核心业务系统开发，负责模块设计和编码实现。', 
 ARRAY['3年以上Java开发经验', '熟悉Spring Boot', '了解微服务架构', '熟悉MySQL'], 
 '25-40K·14薪', '成都', 'full-time', 'open', 7, '技术部', 'mid', 
 ARRAY['Java', 'Spring Boot', 'MyBatis', 'MySQL', 'Redis', 'RabbitMQ'], 
 ARRAY['五险一金', '年终奖', '带薪年假', '餐补'], 2, false, NULL),

('前端开发实习生', '参与前端项目开发，学习前端技术栈。', 
 ARRAY['在校本科及以上学生', '了解HTML/CSS/JavaScript', '有Vue/React基础优先', '每周至少4天实习'], 
 '150-200元/天', '北京', 'internship', 'open', 6, '技术部', 'junior', 
 ARRAY['HTML', 'CSS', 'JavaScript', 'Vue'], 
 ARRAY['实习证明', '转正机会', '导师指导'], 3, false, NULL),

('人力资源专员', '负责招聘渠道维护、简历筛选、面试安排等工作。', 
 ARRAY['1年以上HR相关经验', '熟悉招聘流程', '良好的沟通能力', '熟练使用办公软件'], 
 '10-15K·13薪', '北京', 'full-time', 'open', 2, '人力资源部', 'junior', 
 ARRAY['招聘', '面试', '人力资源', 'Excel'], 
 ARRAY['五险一金', '年终奖', '带薪年假', '节日福利'], 1, false, NULL),

('Python开发工程师', '负责数据平台开发和维护，参与AI模型部署。', 
 ARRAY['3年以上Python经验', '熟悉Django/Flask', '了解机器学习基础'], 
 '25-35K·13薪', '上海', 'full-time', 'closed', 5, '技术部', 'mid', 
 ARRAY['Python', 'Django', 'Flask', 'MySQL', 'TensorFlow'], 
 ARRAY['五险一金', '年终奖'], 1, false, NULL),

('运营专员', '负责平台日常运营工作，用户增长和活跃度提升。', 
 ARRAY['1年以上运营经验', '熟悉数据分析', '有互联网运营经验优先'], 
 '8-12K·13薪', '杭州', 'full-time', 'filled', 8, '运营部', 'junior', 
 ARRAY['运营', '数据分析', '文案', '活动策划'], 
 ARRAY['五险一金', '年终奖'], 1, false, NULL);
