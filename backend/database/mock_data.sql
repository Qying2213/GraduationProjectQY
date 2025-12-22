-- =====================================================
-- 人才运营平台 - Mock 数据
-- =====================================================

-- 清空现有数据（按外键依赖顺序）
TRUNCATE TABLE interview_feedbacks, interviews, applications, resumes, messages, talents, jobs, users RESTART IDENTITY CASCADE;

-- =====================================================
-- 1. 用户表 (users)
-- =====================================================
INSERT INTO users (username, email, password, role, avatar, phone, department, position, status, created_at, updated_at) VALUES
-- 密码都是 'password123' 的 bcrypt 加密
('admin', 'admin@company.com', '$2a$10$yGMtPeEFpcTrNiGsIEWXueRkbPEp5zJQZNX4kHuuh9W9sqKHxG2qa', 'admin', 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin', '13800000001', '技术部', '系统管理员', 'active', NOW() - INTERVAL '180 days', NOW()),
('hr_zhang', 'zhang.hr@company.com', '$2a$10$yGMtPeEFpcTrNiGsIEWXueRkbPEp5zJQZNX4kHuuh9W9sqKHxG2qa', 'hr_manager', 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang', '13800000002', '人力资源部', 'HR总监', 'active', NOW() - INTERVAL '150 days', NOW()),
('hr_li', 'li.hr@company.com', '$2a$10$yGMtPeEFpcTrNiGsIEWXueRkbPEp5zJQZNX4kHuuh9W9sqKHxG2qa', 'recruiter', 'https://api.dicebear.com/7.x/avataaars/svg?seed=li', '13800000003', '人力资源部', '招聘专员', 'active', NOW() - INTERVAL '120 days', NOW()),
('hr_wang', 'wang.hr@company.com', '$2a$10$yGMtPeEFpcTrNiGsIEWXueRkbPEp5zJQZNX4kHuuh9W9sqKHxG2qa', 'recruiter', 'https://api.dicebear.com/7.x/avataaars/svg?seed=wang', '13800000004', '人力资源部', '招聘专员', 'active', NOW() - INTERVAL '90 days', NOW()),
('tech_chen', 'chen.tech@company.com', '$2a$10$yGMtPeEFpcTrNiGsIEWXueRkbPEp5zJQZNX4kHuuh9W9sqKHxG2qa', 'interviewer', 'https://api.dicebear.com/7.x/avataaars/svg?seed=chen', '13800000005', '技术部', '技术总监', 'active', NOW() - INTERVAL '200 days', NOW()),
('tech_liu', 'liu.tech@company.com', '$2a$10$yGMtPeEFpcTrNiGsIEWXueRkbPEp5zJQZNX4kHuuh9W9sqKHxG2qa', 'interviewer', 'https://api.dicebear.com/7.x/avataaars/svg?seed=liu', '13800000006', '技术部', '前端负责人', 'active', NOW() - INTERVAL '180 days', NOW()),
('tech_zhao', 'zhao.tech@company.com', '$2a$10$yGMtPeEFpcTrNiGsIEWXueRkbPEp5zJQZNX4kHuuh9W9sqKHxG2qa', 'interviewer', 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhao', '13800000007', '技术部', '后端负责人', 'active', NOW() - INTERVAL '160 days', NOW()),
('product_sun', 'sun.product@company.com', '$2a$10$yGMtPeEFpcTrNiGsIEWXueRkbPEp5zJQZNX4kHuuh9W9sqKHxG2qa', 'interviewer', 'https://api.dicebear.com/7.x/avataaars/svg?seed=sun', '13800000008', '产品部', '产品总监', 'active', NOW() - INTERVAL '140 days', NOW()),
('viewer_test', 'viewer@company.com', '$2a$10$yGMtPeEFpcTrNiGsIEWXueRkbPEp5zJQZNX4kHuuh9W9sqKHxG2qa', 'viewer', 'https://api.dicebear.com/7.x/avataaars/svg?seed=viewer', '13800000009', '运营部', '运营专员', 'active', NOW() - INTERVAL '30 days', NOW());

-- =====================================================
-- 2. 职位表 (jobs)
-- =====================================================
INSERT INTO jobs (title, description, requirements, salary, location, type, status, created_by, department, level, skills, benefits, created_at, updated_at) VALUES
('高级前端工程师',
'负责公司核心产品的前端架构设计和开发工作，主导前端技术选型，优化前端性能，提升用户体验。需要有较强的技术视野和团队协作能力。',
ARRAY['5年以上前端开发经验', '精通Vue3/React等主流框架', '熟悉TypeScript', '有大型项目架构经验', '良好的沟通能力'],
'35-50K·15薪', '北京', 'full-time', 'open', 5, '技术部', 'senior',
ARRAY['Vue3', 'TypeScript', 'React', 'Webpack', 'Node.js'],
ARRAY['五险一金', '年终奖', '股票期权', '带薪年假', '弹性工作'],
NOW() - INTERVAL '30 days', NOW()),

('后端开发工程师（Go）',
'参与公司微服务架构的设计和开发，负责高并发、高可用系统的建设。需要扎实的编程功底和分布式系统经验。',
ARRAY['3年以上后端开发经验', '精通Go语言', '熟悉微服务架构', '了解Kubernetes', '有高并发系统经验优先'],
'30-45K·14薪', '北京', 'full-time', 'open', 5, '技术部', 'mid',
ARRAY['Go', 'gRPC', 'MySQL', 'Redis', 'Kubernetes', 'Docker'],
ARRAY['五险一金', '年终奖', '带薪年假', '免费三餐', '健身房'],
NOW() - INTERVAL '25 days', NOW()),

('产品经理',
'负责公司B端产品的规划和设计，深入理解用户需求，输出高质量的产品方案。需要有较强的逻辑思维和数据分析能力。',
ARRAY['3年以上B端产品经验', '熟悉互联网产品开发流程', '优秀的文档能力', '数据分析能力', '有HR SaaS经验优先'],
'25-40K·13薪', '上海', 'full-time', 'open', 8, '产品部', 'mid',
ARRAY['Axure', 'Figma', 'SQL', '数据分析', '用户研究'],
ARRAY['五险一金', '年终奖', '带薪年假', '弹性工作'],
NOW() - INTERVAL '20 days', NOW()),

('UI设计师',
'负责公司产品的视觉设计和交互设计，打造优质的用户体验。需要有出色的审美能力和设计功底。',
ARRAY['3年以上UI设计经验', '精通Figma/Sketch', '有B端产品设计经验', '了解前端开发基础', '作品集优秀'],
'20-35K·13薪', '深圳', 'full-time', 'open', 8, '设计部', 'mid',
ARRAY['Figma', 'Sketch', 'Adobe XD', 'Photoshop', 'Illustrator'],
ARRAY['五险一金', '年终奖', '带薪年假', '设计培训'],
NOW() - INTERVAL '15 days', NOW()),

('数据分析师',
'负责业务数据的分析和挖掘，为产品和运营决策提供数据支持。需要扎实的数据分析能力和业务理解能力。',
ARRAY['2年以上数据分析经验', '精通SQL和Python', '熟悉常用BI工具', '良好的数据敏感度', '有互联网行业经验优先'],
'20-30K·13薪', '杭州', 'full-time', 'open', 8, '数据部', 'junior',
ARRAY['Python', 'SQL', 'Tableau', 'Excel', '统计学'],
ARRAY['五险一金', '年终奖', '带薪年假', '培训机会'],
NOW() - INTERVAL '10 days', NOW()),

('DevOps工程师',
'负责公司CI/CD流程建设和运维自动化，保障系统稳定运行。需要全面的运维技能和自动化开发能力。',
ARRAY['3年以上运维经验', '精通Linux系统', '熟悉Docker/K8s', '有CI/CD搭建经验', '会脚本编程'],
'25-40K·14薪', '北京', 'full-time', 'open', 5, '技术部', 'mid',
ARRAY['Linux', 'Docker', 'Kubernetes', 'Jenkins', 'Ansible', 'Python'],
ARRAY['五险一金', '年终奖', '带薪年假', '技术培训'],
NOW() - INTERVAL '8 days', NOW()),

('测试工程师',
'负责产品质量保障，制定测试计划和用例，执行功能和性能测试。需要细致的工作态度和测试思维。',
ARRAY['2年以上测试经验', '熟悉测试流程和方法', '会自动化测试', '了解性能测试', '有接口测试经验'],
'15-25K·13薪', '广州', 'full-time', 'open', 5, '技术部', 'junior',
ARRAY['Selenium', 'JMeter', 'Postman', 'Python', 'SQL'],
ARRAY['五险一金', '年终奖', '带薪年假'],
NOW() - INTERVAL '5 days', NOW()),

('Java开发工程师',
'参与公司核心业务系统开发，负责模块设计和编码实现。需要扎实的Java基础和良好的编码习惯。',
ARRAY['3年以上Java开发经验', '熟悉Spring Boot', '了解微服务架构', '熟悉MySQL', '有分布式系统经验优先'],
'25-40K·14薪', '成都', 'full-time', 'open', 7, '技术部', 'mid',
ARRAY['Java', 'Spring Boot', 'MyBatis', 'MySQL', 'Redis', 'RabbitMQ'],
ARRAY['五险一金', '年终奖', '带薪年假', '餐补'],
NOW() - INTERVAL '3 days', NOW()),

('实习生-前端开发',
'参与前端项目开发，学习前端技术栈。需要有一定的前端基础和学习热情。',
ARRAY['在校本科及以上学生', '了解HTML/CSS/JavaScript', '有Vue/React基础优先', '每周至少4天实习'],
'150-200元/天', '北京', 'internship', 'open', 6, '技术部', 'junior',
ARRAY['HTML', 'CSS', 'JavaScript', 'Vue'],
ARRAY['实习证明', '转正机会', '导师指导'],
NOW() - INTERVAL '2 days', NOW()),

('人力资源专员',
'负责招聘渠道维护、简历筛选、面试安排等工作。需要良好的沟通能力和服务意识。',
ARRAY['1年以上HR相关经验', '熟悉招聘流程', '良好的沟通能力', '熟练使用办公软件', '有互联网公司经验优先'],
'10-15K·13薪', '北京', 'full-time', 'open', 2, '人力资源部', 'junior',
ARRAY['招聘', '面试', '人力资源', 'Excel'],
ARRAY['五险一金', '年终奖', '带薪年假', '节日福利'],
NOW() - INTERVAL '1 day', NOW()),

-- 已关闭的职位
('Python开发工程师',
'负责数据平台开发和维护。',
ARRAY['3年以上Python经验', '熟悉Django/Flask'],
'25-35K·13薪', '上海', 'full-time', 'closed', 5, '技术部', 'mid',
ARRAY['Python', 'Django', 'Flask', 'MySQL'],
ARRAY['五险一金', '年终奖'],
NOW() - INTERVAL '60 days', NOW() - INTERVAL '10 days'),

('运营专员',
'负责平台日常运营工作。',
ARRAY['1年以上运营经验', '熟悉数据分析'],
'8-12K·13薪', '杭州', 'full-time', 'filled', 8, '运营部', 'junior',
ARRAY['运营', '数据分析', '文案'],
ARRAY['五险一金', '年终奖'],
NOW() - INTERVAL '45 days', NOW() - INTERVAL '5 days');

-- =====================================================
-- 3. 人才表 (talents)
-- =====================================================
INSERT INTO talents (name, email, phone, skills, experience, education, status, tags, user_id, location, salary, summary, gender, age, current_company, current_position, created_at, updated_at) VALUES
('张伟', 'zhangwei@email.com', '13912345001', ARRAY['Vue3', 'TypeScript', 'React', 'Node.js', 'Webpack'], 6, '本科', 'active', ARRAY['高潜力', '技术专家'], 3, '北京', '35-45K', '6年前端开发经验，精通Vue全家桶，有大型项目架构经验。曾主导多个百万级用户产品的前端开发。', '男', 29, '字节跳动', '高级前端工程师', NOW() - INTERVAL '45 days', NOW()),

('李娜', 'lina@email.com', '13912345002', ARRAY['Go', 'Python', 'MySQL', 'Redis', 'Kubernetes'], 5, '硕士', 'active', ARRAY['技术专家'], 3, '北京', '40-55K', '5年后端开发经验，精通Go和分布式系统。曾负责千万级DAU系统的架构设计。', '女', 28, '美团', '后端技术专家', NOW() - INTERVAL '40 days', NOW()),

('王磊', 'wanglei@email.com', '13912345003', ARRAY['Java', 'Spring Boot', 'MySQL', 'Redis', 'RabbitMQ'], 4, '本科', 'active', ARRAY['稳定'], 4, '上海', '25-35K', '4年Java开发经验，熟悉微服务架构，有电商系统开发经验。', '男', 27, '拼多多', 'Java开发工程师', NOW() - INTERVAL '35 days', NOW()),

('刘芳', 'liufang@email.com', '13912345004', ARRAY['产品设计', 'Axure', 'SQL', '数据分析', '用户研究'], 4, '本科', 'active', ARRAY['沟通能力强'], 3, '上海', '25-40K', '4年B端产品经验，擅长需求分析和产品设计。曾负责多个SaaS产品从0到1的设计。', '女', 28, '有赞', '高级产品经理', NOW() - INTERVAL '30 days', NOW()),

('陈强', 'chenqiang@email.com', '13912345005', ARRAY['Figma', 'Sketch', 'Photoshop', 'Illustrator', '动效设计'], 3, '本科', 'active', ARRAY['创意'], 4, '深圳', '20-30K', '3年UI设计经验，有优秀的审美能力，擅长B端产品设计。', '男', 26, '腾讯', 'UI设计师', NOW() - INTERVAL '28 days', NOW()),

('赵敏', 'zhaomin@email.com', '13912345006', ARRAY['Python', 'SQL', 'Tableau', 'Excel', '机器学习'], 3, '硕士', 'active', ARRAY['数据敏感'], 3, '杭州', '22-32K', '3年数据分析经验，擅长用户行为分析和业务建模。', '女', 27, '阿里巴巴', '数据分析师', NOW() - INTERVAL '25 days', NOW()),

('孙涛', 'suntao@email.com', '13912345007', ARRAY['Linux', 'Docker', 'Kubernetes', 'Jenkins', 'Prometheus'], 5, '本科', 'active', ARRAY['技术专家'], 4, '北京', '30-45K', '5年运维经验，精通容器化和自动化运维，有大规模集群管理经验。', '男', 30, '快手', 'SRE工程师', NOW() - INTERVAL '22 days', NOW()),

('周婷', 'zhouting@email.com', '13912345008', ARRAY['Selenium', 'JMeter', 'Python', 'Postman', '性能测试'], 3, '本科', 'active', ARRAY['细致'], 4, '广州', '15-22K', '3年测试经验，熟悉自动化测试和性能测试。', '女', 26, '网易', '测试工程师', NOW() - INTERVAL '20 days', NOW()),

('吴杰', 'wujie@email.com', '13912345009', ARRAY['Vue3', 'React', 'TypeScript', 'Webpack', 'Vite'], 3, '本科', 'active', ARRAY['成长快'], 3, '成都', '20-30K', '3年前端经验，技术栈全面，学习能力强。', '男', 26, '字节跳动', '前端工程师', NOW() - INTERVAL '18 days', NOW()),

('郑雪', 'zhengxue@email.com', '13912345010', ARRAY['Go', 'gRPC', 'MySQL', 'Redis', 'Kafka'], 2, '硕士', 'active', ARRAY['高潜力'], 4, '北京', '25-35K', '2年Go开发经验，有微服务和高并发系统经验。', '女', 25, '小米', '后端开发工程师', NOW() - INTERVAL '15 days', NOW()),

('黄浩', 'huanghao@email.com', '13912345011', ARRAY['Java', 'Spring Cloud', 'MySQL', 'MongoDB', 'Elasticsearch'], 7, '本科', 'hired', ARRAY['技术专家', '架构师'], 3, '北京', '50-70K', '7年Java开发经验，有丰富的架构设计经验。', '男', 32, '京东', '架构师', NOW() - INTERVAL '50 days', NOW() - INTERVAL '5 days'),

('林小红', 'linxiaohong@email.com', '13912345012', ARRAY['招聘', 'HRBP', '员工关系', '培训'], 4, '本科', 'active', ARRAY['沟通能力强'], 3, '北京', '15-25K', '4年HR经验，擅长招聘和员工关系管理。', '女', 28, '滴滴', 'HRBP', NOW() - INTERVAL '12 days', NOW()),

('马云飞', 'mayunfei@email.com', '13912345013', ARRAY['React', 'TypeScript', 'Next.js', 'GraphQL'], 4, '本科', 'pending', ARRAY[]::text[], 4, '深圳', '28-38K', '4年前端经验，专注React技术栈。', '男', 27, '腾讯', '前端开发工程师', NOW() - INTERVAL '10 days', NOW()),

('杨丽', 'yangli@email.com', '13912345014', ARRAY['Python', 'Django', 'FastAPI', 'PostgreSQL'], 3, '本科', 'active', ARRAY[]::text[], 3, '上海', '22-32K', '3年Python开发经验，有Web开发和数据处理经验。', '女', 26, '携程', 'Python开发工程师', NOW() - INTERVAL '8 days', NOW()),

('徐明', 'xuming@email.com', '13912345015', ARRAY['iOS', 'Swift', 'Objective-C', 'Flutter'], 5, '本科', 'active', ARRAY['移动端专家'], 4, '杭州', '35-50K', '5年iOS开发经验，有多款App Store上架产品。', '男', 29, '网易', 'iOS开发工程师', NOW() - INTERVAL '6 days', NOW()),

('何芳', 'hefang@email.com', '13912345016', ARRAY['Android', 'Kotlin', 'Java', 'Flutter'], 4, '本科', 'active', ARRAY[]::text[], 3, '广州', '25-35K', '4年Android开发经验。', '女', 27, 'OPPO', 'Android开发工程师', NOW() - INTERVAL '5 days', NOW()),

('谢伟', 'xiewei@email.com', '13912345017', ARRAY['Vue', 'JavaScript', 'CSS', 'HTML'], 1, '本科', 'active', ARRAY['应届生'], 4, '北京', '12-18K', '1年前端经验，基础扎实，学习能力强。', '男', 23, '', '前端开发实习生', NOW() - INTERVAL '3 days', NOW()),

('罗琳', 'luolin@email.com', '13912345018', ARRAY['产品运营', '数据分析', '用户增长'], 2, '本科', 'rejected', ARRAY[]::text[], 3, '上海', '12-18K', '2年运营经验。', '女', 25, '拼多多', '运营专员', NOW() - INTERVAL '30 days', NOW() - INTERVAL '20 days'),

('高强', 'gaoqiang@email.com', '13912345019', ARRAY['Node.js', 'Express', 'MongoDB', 'AWS'], 3, '硕士', 'active', ARRAY['全栈'], 4, '深圳', '25-38K', '3年全栈开发经验。', '男', 27, '华为', '全栈工程师', NOW() - INTERVAL '2 days', NOW()),

('唐雨', 'tangyu@email.com', '13912345020', ARRAY['机器学习', 'Python', 'TensorFlow', 'PyTorch'], 2, '博士', 'active', ARRAY['AI专家', '高学历'], 3, '北京', '40-60K', '2年AI研发经验，专注NLP方向。', '女', 28, '百度', 'AI算法工程师', NOW() - INTERVAL '1 day', NOW());

-- =====================================================
-- 4. 简历表 (resumes)
-- =====================================================
INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, created_at, updated_at) VALUES
(1, 1, '/uploads/resumes/zhangwei_resume.pdf', '张伟_高级前端工程师_简历.pdf', 'interviewed', 92, NOW() - INTERVAL '40 days', NOW() - INTERVAL '5 days'),
(2, 2, '/uploads/resumes/lina_resume.pdf', '李娜_后端开发工程师_简历.pdf', 'offered', 95, NOW() - INTERVAL '35 days', NOW() - INTERVAL '3 days'),
(3, 8, '/uploads/resumes/wanglei_resume.pdf', '王磊_Java开发工程师_简历.pdf', 'reviewing', 78, NOW() - INTERVAL '30 days', NOW()),
(4, 3, '/uploads/resumes/liufang_resume.pdf', '刘芳_产品经理_简历.pdf', 'interviewed', 88, NOW() - INTERVAL '28 days', NOW() - INTERVAL '10 days'),
(5, 4, '/uploads/resumes/chenqiang_resume.pdf', '陈强_UI设计师_简历.pdf', 'pending', 85, NOW() - INTERVAL '25 days', NOW()),
(6, 5, '/uploads/resumes/zhaomin_resume.pdf', '赵敏_数据分析师_简历.pdf', 'interviewed', 90, NOW() - INTERVAL '22 days', NOW() - INTERVAL '8 days'),
(7, 6, '/uploads/resumes/suntao_resume.pdf', '孙涛_DevOps工程师_简历.pdf', 'offered', 93, NOW() - INTERVAL '20 days', NOW() - INTERVAL '2 days'),
(8, 7, '/uploads/resumes/zhouting_resume.pdf', '周婷_测试工程师_简历.pdf', 'reviewing', 82, NOW() - INTERVAL '18 days', NOW()),
(9, 1, '/uploads/resumes/wujie_resume.pdf', '吴杰_前端工程师_简历.pdf', 'pending', 75, NOW() - INTERVAL '15 days', NOW()),
(10, 2, '/uploads/resumes/zhengxue_resume.pdf', '郑雪_后端开发工程师_简历.pdf', 'interviewed', 87, NOW() - INTERVAL '12 days', NOW() - INTERVAL '5 days'),
(11, 8, '/uploads/resumes/huanghao_resume.pdf', '黄浩_架构师_简历.pdf', 'hired', 96, NOW() - INTERVAL '45 days', NOW() - INTERVAL '5 days'),
(12, 10, '/uploads/resumes/linxiaohong_resume.pdf', '林小红_HR专员_简历.pdf', 'reviewing', 80, NOW() - INTERVAL '10 days', NOW()),
(13, 1, '/uploads/resumes/mayunfei_resume.pdf', '马云飞_前端工程师_简历.pdf', 'pending', 72, NOW() - INTERVAL '8 days', NOW()),
(14, 11, '/uploads/resumes/yangli_resume.pdf', '杨丽_Python工程师_简历.pdf', 'rejected', 68, NOW() - INTERVAL '55 days', NOW() - INTERVAL '40 days'),
(17, 9, '/uploads/resumes/xiewei_resume.pdf', '谢伟_前端实习生_简历.pdf', 'pending', 70, NOW() - INTERVAL '2 days', NOW()),
(19, 2, '/uploads/resumes/gaoqiang_resume.pdf', '高强_全栈工程师_简历.pdf', 'pending', 76, NOW() - INTERVAL '1 day', NOW()),
(20, 2, '/uploads/resumes/tangyu_resume.pdf', '唐雨_AI算法工程师_简历.pdf', 'reviewing', 65, NOW() - INTERVAL '1 day', NOW());

-- =====================================================
-- 5. 面试表 (interviews)
-- =====================================================
INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, duration, interviewer_id, interviewer, method, location, status, notes, feedback, rating, created_by, created_at, updated_at) VALUES
-- 今天的面试
(1, '张伟', 1, '高级前端工程师', 'final', CURRENT_DATE::text, '10:00', 60, 5, '陈总监', 'onsite', '5楼会议室A', 'scheduled', '终面，评估整体能力和团队匹配度', NULL, 0, 3, NOW() - INTERVAL '2 days', NOW()),
(6, '赵敏', 5, '数据分析师', 'second', CURRENT_DATE::text, '14:30', 60, 8, '孙总监', 'video', 'https://meeting.company.com/data123', 'scheduled', '复试，重点考察SQL和业务分析能力', NULL, 0, 4, NOW() - INTERVAL '3 days', NOW()),
(10, '郑雪', 2, '后端开发工程师（Go）', 'initial', CURRENT_DATE::text, '16:00', 90, 7, '赵主管', 'onsite', '3楼会议室B', 'scheduled', '技术初试，考察Go语言和系统设计', NULL, 0, 3, NOW() - INTERVAL '1 day', NOW()),

-- 明天的面试
(4, '刘芳', 3, '产品经理', 'hr', (CURRENT_DATE + INTERVAL '1 day')::date::text, '10:30', 45, 2, '张HR总监', 'onsite', '2楼HR办公室', 'scheduled', 'HR面试，谈薪资和入职时间', NULL, 0, 3, NOW() - INTERVAL '5 days', NOW()),
(5, '陈强', 4, 'UI设计师', 'initial', (CURRENT_DATE + INTERVAL '1 day')::date::text, '14:00', 60, 8, '孙总监', 'video', 'https://meeting.company.com/design456', 'scheduled', '初试，需要展示作品集', NULL, 0, 4, NOW() - INTERVAL '2 days', NOW()),

-- 后天的面试
(3, '王磊', 8, 'Java开发工程师', 'initial', (CURRENT_DATE + INTERVAL '2 day')::date::text, '09:30', 90, 7, '赵主管', 'onsite', '3楼会议室A', 'scheduled', '技术初试', NULL, 0, 4, NOW() - INTERVAL '1 day', NOW()),
(8, '周婷', 7, '测试工程师', 'initial', (CURRENT_DATE + INTERVAL '2 day')::date::text, '15:00', 60, 6, '刘经理', 'phone', '13800000006', 'scheduled', '电话初试', NULL, 0, 4, NOW() - INTERVAL '1 day', NOW()),

-- 已完成的面试
(1, '张伟', 1, '高级前端工程师', 'initial', (CURRENT_DATE - INTERVAL '15 day')::date::text, '10:00', 90, 6, '刘经理', 'onsite', '3楼会议室A', 'completed', '技术初试', '技术基础扎实，Vue3掌握深入，有架构思维。建议进入复试。', 4, 3, NOW() - INTERVAL '20 days', NOW() - INTERVAL '15 days'),
(1, '张伟', 1, '高级前端工程师', 'second', (CURRENT_DATE - INTERVAL '10 day')::date::text, '14:00', 60, 7, '赵主管', 'onsite', '5楼会议室B', 'completed', '技术复试', '系统设计能力强，代码质量高。建议终面。', 5, 3, NOW() - INTERVAL '15 days', NOW() - INTERVAL '10 days'),
(2, '李娜', 2, '后端开发工程师（Go）', 'initial', (CURRENT_DATE - INTERVAL '20 day')::date::text, '09:30', 90, 7, '赵主管', 'onsite', '3楼会议室B', 'completed', '技术初试', 'Go语言精通，有高并发系统经验。强烈推荐。', 5, 3, NOW() - INTERVAL '25 days', NOW() - INTERVAL '20 days'),
(2, '李娜', 2, '后端开发工程师（Go）', 'second', (CURRENT_DATE - INTERVAL '15 day')::date::text, '14:30', 60, 5, '陈总监', 'onsite', '5楼会议室A', 'completed', '技术复试+终面', '技术能力优秀，团队协作能力强。建议发offer。', 5, 3, NOW() - INTERVAL '20 days', NOW() - INTERVAL '15 days'),
(4, '刘芳', 3, '产品经理', 'initial', (CURRENT_DATE - INTERVAL '18 day')::date::text, '10:00', 60, 8, '孙总监', 'video', 'https://meeting.company.com/pm789', 'completed', '初试', '产品思维清晰，有B端产品经验。进入复试。', 4, 3, NOW() - INTERVAL '22 days', NOW() - INTERVAL '18 days'),
(4, '刘芳', 3, '产品经理', 'second', (CURRENT_DATE - INTERVAL '12 day')::date::text, '14:00', 60, 8, '孙总监', 'onsite', '5楼产品部', 'completed', '复试', '需求分析能力强，文档能力优秀。建议HR面。', 4, 3, NOW() - INTERVAL '16 days', NOW() - INTERVAL '12 days'),
(6, '赵敏', 5, '数据分析师', 'initial', (CURRENT_DATE - INTERVAL '10 day')::date::text, '11:00', 60, 8, '孙总监', 'video', 'https://meeting.company.com/data456', 'completed', '初试', 'SQL能力强，有业务sense。进入复试。', 4, 4, NOW() - INTERVAL '15 days', NOW() - INTERVAL '10 days'),
(7, '孙涛', 6, 'DevOps工程师', 'initial', (CURRENT_DATE - INTERVAL '12 day')::date::text, '09:30', 90, 5, '陈总监', 'onsite', '3楼会议室A', 'completed', '技术初试', 'K8s和CI/CD经验丰富，动手能力强。', 5, 4, NOW() - INTERVAL '18 days', NOW() - INTERVAL '12 days'),
(7, '孙涛', 6, 'DevOps工程师', 'second', (CURRENT_DATE - INTERVAL '8 day')::date::text, '14:00', 60, 5, '陈总监', 'onsite', '5楼会议室A', 'completed', '复试', '综合能力优秀，建议发offer。', 5, 4, NOW() - INTERVAL '12 days', NOW() - INTERVAL '8 days'),
(11, '黄浩', 8, 'Java开发工程师', 'initial', (CURRENT_DATE - INTERVAL '35 day')::date::text, '10:00', 90, 7, '赵主管', 'onsite', '3楼会议室B', 'completed', '技术初试', '架构能力强，经验丰富。', 5, 3, NOW() - INTERVAL '40 days', NOW() - INTERVAL '35 days'),
(11, '黄浩', 8, 'Java开发工程师', 'final', (CURRENT_DATE - INTERVAL '30 day')::date::text, '14:00', 60, 5, '陈总监', 'onsite', '5楼会议室A', 'completed', '终面', '各方面都很优秀，强烈建议录用。', 5, 3, NOW() - INTERVAL '35 days', NOW() - INTERVAL '30 days'),

-- 已取消的面试
(13, '马云飞', 1, '高级前端工程师', 'initial', (CURRENT_DATE - INTERVAL '5 day')::date::text, '10:00', 60, 6, '刘经理', 'onsite', '3楼会议室A', 'cancelled', '候选人临时有事', NULL, 0, 4, NOW() - INTERVAL '10 days', NOW() - INTERVAL '5 days');

-- =====================================================
-- 6. 消息表 (messages)
-- =====================================================
INSERT INTO messages (sender_id, receiver_id, type, title, content, is_read, created_at) VALUES
-- 系统通知
(1, 3, 'system', '新简历提醒', '您收到了一份新简历：张伟 应聘 高级前端工程师 职位', true, NOW() - INTERVAL '40 days'),
(1, 3, 'system', '新简历提醒', '您收到了一份新简历：李娜 应聘 后端开发工程师（Go） 职位', true, NOW() - INTERVAL '35 days'),
(1, 4, 'system', '新简历提醒', '您收到了一份新简历：王磊 应聘 Java开发工程师 职位', true, NOW() - INTERVAL '30 days'),
(1, 3, 'system', '面试安排提醒', '面试已安排：张伟 - 高级前端工程师 初试，时间：明天 10:00', true, NOW() - INTERVAL '16 days'),
(1, 6, 'interview', '面试邀请', '您有一场新的面试需要参加：张伟 - 高级前端工程师 初试', true, NOW() - INTERVAL '16 days'),

-- 面试反馈
(6, 3, 'feedback', '面试反馈已提交', '刘经理 已提交 张伟 的面试反馈，评分：4星', true, NOW() - INTERVAL '15 days'),
(7, 3, 'feedback', '面试反馈已提交', '赵主管 已提交 张伟 的面试反馈，评分：5星', true, NOW() - INTERVAL '10 days'),
(7, 3, 'feedback', '面试反馈已提交', '赵主管 已提交 李娜 的面试反馈，评分：5星', true, NOW() - INTERVAL '20 days'),

-- Offer相关
(2, 3, 'offer', 'Offer审批通过', '李娜 的 Offer 已通过审批，请尽快发放', true, NOW() - INTERVAL '10 days'),
(2, 3, 'offer', 'Offer审批通过', '孙涛 的 Offer 已通过审批，请尽快发放', false, NOW() - INTERVAL '5 days'),

-- 未读消息
(1, 3, 'system', '新简历提醒', '您收到了一份新简历：高强 应聘 后端开发工程师（Go） 职位', false, NOW() - INTERVAL '1 day'),
(1, 3, 'system', '新简历提醒', '您收到了一份新简历：唐雨 应聘 后端开发工程师（Go） 职位', false, NOW() - INTERVAL '1 day'),
(1, 4, 'system', '面试安排提醒', '面试已安排：王磊 - Java开发工程师 初试，时间：后天 09:30', false, NOW() - INTERVAL '1 day'),
(5, 3, 'feedback', '面试反馈已提交', '陈总监 已提交 张伟 终面评估，请查看', false, NOW() - INTERVAL '2 days'),
(1, 3, 'reminder', '面试提醒', '您今天有 3 场面试待进行，请做好准备', false, NOW()),

-- 沟通消息
(3, 6, 'chat', '关于张伟的面试', '刘经理，张伟的技术初试您觉得怎么样？', true, NOW() - INTERVAL '14 days'),
(6, 3, 'chat', '回复：关于张伟的面试', '他Vue3掌握得很好，架构思维也不错，建议进入复试。', true, NOW() - INTERVAL '14 days'),
(3, 5, 'chat', '李娜的offer', '陈总监，李娜的offer谈判进展如何？', true, NOW() - INTERVAL '8 days'),
(5, 3, 'chat', '回复：李娜的offer', '已经谈好了，她接受我们的offer，下周一入职。', true, NOW() - INTERVAL '8 days');

-- =====================================================
-- 7. 申请/应聘记录表 (applications) - 如果有的话
-- =====================================================
-- 这个表用于记录候选人的完整应聘历程

-- =====================================================
-- 统计数据验证
-- =====================================================
-- SELECT '用户数' as metric, COUNT(*) as count FROM users
-- UNION ALL SELECT '职位数', COUNT(*) FROM jobs
-- UNION ALL SELECT '人才数', COUNT(*) FROM talents
-- UNION ALL SELECT '简历数', COUNT(*) FROM resumes
-- UNION ALL SELECT '面试数', COUNT(*) FROM interviews
-- UNION ALL SELECT '消息数', COUNT(*) FROM messages;
