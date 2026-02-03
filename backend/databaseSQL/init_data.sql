-- 智能人才运营平台 - 初始化测试数据
-- 执行方式: psql -d talent_platform -f backend/database/init_data.sql

-- 清空所有表（如果存在）
DO $$
BEGIN
    -- 只在表存在时清空
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users') THEN
        TRUNCATE TABLE chat_messages, conversations, applications, resumes, interviews, talents, jobs, users RESTART IDENTITY CASCADE;
    END IF;
END $$;

-- =====================================================
-- 1. 用户数据（明文密码）
-- =====================================================
INSERT INTO users (username, email, password, role, real_name, phone, department, position, status, created_at, updated_at) VALUES
-- 管理员
('admin', 'admin@company.com', 'admin123', 'admin', '系统管理员', '13800000001', '技术部', '系统管理员', 'active', NOW(), NOW()),
-- HR (id: 2-5)
('hr1', 'hr1@company.com', '123456', 'hr', '张丽华', '13800000002', '人力资源部', 'HR经理', 'active', NOW(), NOW()),
('hr2', 'hr2@company.com', '123456', 'hr', '李明', '13800000003', '人力资源部', 'HR专员', 'active', NOW(), NOW()),
('hr3', 'hr3@company.com', '123456', 'hr', '王芳', '13800000004', '人力资源部', 'HR专员', 'active', NOW(), NOW()),
('hr4', 'hr4@company.com', '123456', 'hr', '赵强', '13800000005', '人力资源部', '招聘主管', 'active', NOW(), NOW()),
-- 求职者 (id: 6-25)
('candidate01', 'zhangwei@gmail.com', '123456', 'candidate', '张伟', '13900000001', '', '', 'active', NOW(), NOW()),
('candidate02', 'lina@gmail.com', '123456', 'candidate', '李娜', '13900000002', '', '', 'active', NOW(), NOW()),
('candidate03', 'wangqiang@gmail.com', '123456', 'candidate', '王强', '13900000003', '', '', 'active', NOW(), NOW()),
('candidate04', 'liuxia@gmail.com', '123456', 'candidate', '刘霞', '13900000004', '', '', 'active', NOW(), NOW()),
('candidate05', 'chenlong@gmail.com', '123456', 'candidate', '陈龙', '13900000005', '', '', 'active', NOW(), NOW()),
('candidate06', 'zhaomin@gmail.com', '123456', 'candidate', '赵敏', '13900000006', '', '', 'active', NOW(), NOW()),
('candidate07', 'sunhao@gmail.com', '123456', 'candidate', '孙浩', '13900000007', '', '', 'active', NOW(), NOW()),
('candidate08', 'zhouting@gmail.com', '123456', 'candidate', '周婷', '13900000008', '', '', 'active', NOW(), NOW()),
('candidate09', 'wugang@gmail.com', '123456', 'candidate', '吴刚', '13900000009', '', '', 'active', NOW(), NOW()),
('candidate10', 'zhengyan@gmail.com', '123456', 'candidate', '郑燕', '13900000010', '', '', 'active', NOW(), NOW()),
('candidate11', 'huangpeng@gmail.com', '123456', 'candidate', '黄鹏', '13900000011', '', '', 'active', NOW(), NOW()),
('candidate12', 'linxue@gmail.com', '123456', 'candidate', '林雪', '13900000012', '', '', 'active', NOW(), NOW()),
('candidate13', 'hejun@gmail.com', '123456', 'candidate', '何军', '13900000013', '', '', 'active', NOW(), NOW()),
('candidate14', 'gaoli@gmail.com', '123456', 'candidate', '高丽', '13900000014', '', '', 'active', NOW(), NOW()),
('candidate15', 'luowei@gmail.com', '123456', 'candidate', '罗伟', '13900000015', '', '', 'active', NOW(), NOW()),
('candidate16', 'xiefang@gmail.com', '123456', 'candidate', '谢芳', '13900000016', '', '', 'active', NOW(), NOW()),
('candidate17', 'tangming@gmail.com', '123456', 'candidate', '唐明', '13900000017', '', '', 'active', NOW(), NOW()),
('candidate18', 'hanxiao@gmail.com', '123456', 'candidate', '韩晓', '13900000018', '', '', 'active', NOW(), NOW()),
('candidate19', 'fengyu@gmail.com', '123456', 'candidate', '冯宇', '13900000019', '', '', 'active', NOW(), NOW()),
('candidate20', 'caijing@gmail.com', '123456', 'candidate', '蔡静', '13900000020', '', '', 'active', NOW(), NOW());

-- =====================================================
-- 2. 职位数据 (24个职位: 18个招聘中, 3个已关闭, 3个待审核)
-- =====================================================
INSERT INTO jobs (title, department, location, salary, type, education, description, requirements, benefits, skills, status, headcount, created_by, created_at, updated_at) VALUES
-- 招聘中的职位 (18个)
('高级Go后端工程师', '技术部', '深圳', '25K-40K', 'full-time', '本科', '负责公司核心业务系统的后端开发，参与系统架构设计和技术选型。', ARRAY['精通Go语言，3年以上开发经验', '熟悉MySQL、Redis、Kafka等中间件', '有微服务架构经验优先'], ARRAY['五险一金', '带薪年假', '弹性工作', '免费三餐'], ARRAY['Go', 'MySQL', 'Redis', 'Kafka', 'Docker'], 'open', 3, 2, NOW() - INTERVAL '30 days', NOW()),
('前端开发工程师', '技术部', '深圳', '18K-30K', 'full-time', '本科', '负责公司Web前端开发，与后端工程师协作完成产品功能。', ARRAY['精通Vue3/React，熟悉TypeScript', '熟悉前端工程化和性能优化'], ARRAY['五险一金', '带薪年假', '弹性工作'], ARRAY['Vue', 'React', 'TypeScript', 'Webpack'], 'open', 2, 2, NOW() - INTERVAL '25 days', NOW()),
('产品经理', '产品部', '北京', '20K-35K', 'full-time', '本科', '负责公司核心产品的规划和设计，推动产品迭代优化。', ARRAY['3年以上B端产品经验', '熟悉敏捷开发流程'], ARRAY['五险一金', '带薪年假', '股票期权'], ARRAY['产品设计', 'Axure', '数据分析'], 'open', 1, 2, NOW() - INTERVAL '20 days', NOW()),
('UI设计师', '设计部', '上海', '15K-25K', 'full-time', '本科', '负责公司产品的UI设计，输出高质量的设计稿。', ARRAY['精通Figma/Sketch', '有完整的B端设计经验'], ARRAY['五险一金', '带薪年假', '设计培训'], ARRAY['Figma', 'Sketch', 'Photoshop'], 'open', 1, 3, NOW() - INTERVAL '18 days', NOW()),
('数据分析师', '数据部', '深圳', '18K-28K', 'full-time', '本科', '负责公司业务数据分析，输出数据报告和洞察。', ARRAY['熟练使用SQL和Python', '熟悉数据可视化工具'], ARRAY['五险一金', '带薪年假', '数据培训'], ARRAY['Python', 'SQL', 'Tableau', '数据分析'], 'open', 2, 3, NOW() - INTERVAL '15 days', NOW()),
('Java开发工程师', '技术部', '广州', '20K-35K', 'full-time', '本科', '负责公司Java后端系统开发和维护。', ARRAY['精通Java，熟悉Spring Boot', '熟悉MySQL、Redis'], ARRAY['五险一金', '带薪年假', '技术培训'], ARRAY['Java', 'Spring', 'MySQL', 'Redis'], 'open', 2, 2, NOW() - INTERVAL '12 days', NOW()),
('Python开发工程师', '技术部', '深圳', '22K-38K', 'full-time', '本科', '负责AI相关后端服务开发。', ARRAY['精通Python', '熟悉机器学习框架'], ARRAY['五险一金', '带薪年假'], ARRAY['Python', 'TensorFlow', 'PyTorch'], 'open', 2, 2, NOW() - INTERVAL '10 days', NOW()),
('DevOps工程师', '技术部', '深圳', '25K-45K', 'full-time', '本科', '负责公司CI/CD流程和基础设施建设。', ARRAY['熟悉Docker、Kubernetes', '有云平台经验'], ARRAY['五险一金', '带薪年假'], ARRAY['Docker', 'Kubernetes', 'Jenkins', 'AWS'], 'open', 1, 2, NOW() - INTERVAL '8 days', NOW()),
('测试工程师', '技术部', '深圳', '15K-25K', 'full-time', '本科', '负责产品质量保障和自动化测试。', ARRAY['熟悉自动化测试框架', '有性能测试经验'], ARRAY['五险一金', '带薪年假'], ARRAY['Selenium', 'JMeter', 'Python'], 'open', 2, 3, NOW() - INTERVAL '7 days', NOW()),
('Android开发工程师', '移动端', '深圳', '20K-35K', 'full-time', '本科', '负责Android客户端开发。', ARRAY['精通Kotlin/Java', '熟悉Android SDK'], ARRAY['五险一金', '带薪年假'], ARRAY['Kotlin', 'Java', 'Android'], 'open', 1, 2, NOW() - INTERVAL '6 days', NOW()),
('iOS开发工程师', '移动端', '深圳', '22K-38K', 'full-time', '本科', '负责iOS客户端开发。', ARRAY['精通Swift/OC', '熟悉iOS开发'], ARRAY['五险一金', '带薪年假'], ARRAY['Swift', 'Objective-C', 'iOS'], 'open', 1, 2, NOW() - INTERVAL '5 days', NOW()),
('运营专员', '运营部', '北京', '10K-18K', 'full-time', '本科', '负责产品运营和用户增长。', ARRAY['有互联网运营经验', '数据敏感度高'], ARRAY['五险一金', '带薪年假'], ARRAY['数据分析', '用户运营', '活动策划'], 'open', 2, 3, NOW() - INTERVAL '4 days', NOW()),
('市场专员', '市场部', '上海', '12K-20K', 'full-time', '本科', '负责市场推广和品牌建设。', ARRAY['有市场推广经验', '良好的沟通能力'], ARRAY['五险一金', '带薪年假'], ARRAY['市场推广', '品牌策划', '活动执行'], 'open', 1, 3, NOW() - INTERVAL '3 days', NOW()),
('财务专员', '财务部', '深圳', '10K-15K', 'full-time', '本科', '负责日常财务工作。', ARRAY['会计相关专业', '熟悉财务软件'], ARRAY['五险一金', '带薪年假'], ARRAY['会计', '财务分析', 'Excel'], 'open', 1, 3, NOW() - INTERVAL '2 days', NOW()),
('人事专员', '人力资源部', '深圳', '8K-12K', 'full-time', '本科', '负责招聘和员工关系管理。', ARRAY['人力资源相关专业', '有招聘经验'], ARRAY['五险一金', '带薪年假'], ARRAY['招聘', '员工关系', '培训'], 'open', 1, 2, NOW() - INTERVAL '1 day', NOW()),
('行政专员', '行政部', '深圳', '6K-10K', 'full-time', '大专', '负责日常行政事务。', ARRAY['有行政工作经验', '细心负责'], ARRAY['五险一金', '带薪年假'], ARRAY['行政管理', 'Office'], 'open', 1, 3, NOW(), NOW()),
('客服专员', '客服部', '深圳', '6K-10K', 'full-time', '大专', '负责客户服务和问题处理。', ARRAY['有客服经验', '沟通能力强'], ARRAY['五险一金', '带薪年假'], ARRAY['客户服务', '沟通协调'], 'open', 2, 3, NOW(), NOW()),
('法务专员', '法务部', '深圳', '15K-25K', 'full-time', '本科', '负责公司法律事务。', ARRAY['法学专业', '有企业法务经验'], ARRAY['五险一金', '带薪年假'], ARRAY['合同审核', '法律咨询'], 'open', 1, 2, NOW(), NOW()),
-- 已关闭的职位 (3个)
('高级架构师', '技术部', '深圳', '50K-80K', 'full-time', '本科', '负责系统架构设计。', ARRAY['10年以上开发经验', '有大型系统架构经验'], ARRAY['五险一金', '股票期权'], ARRAY['架构设计', '分布式系统'], 'closed', 1, 2, NOW() - INTERVAL '60 days', NOW()),
('CTO助理', '技术部', '深圳', '30K-50K', 'full-time', '硕士', '协助CTO处理技术管理工作。', ARRAY['5年以上技术管理经验'], ARRAY['五险一金', '股票期权'], ARRAY['技术管理', '项目管理'], 'closed', 1, 2, NOW() - INTERVAL '45 days', NOW()),
('技术总监', '技术部', '北京', '60K-100K', 'full-time', '本科', '负责技术团队管理。', ARRAY['15年以上技术经验', '有团队管理经验'], ARRAY['五险一金', '股票期权'], ARRAY['团队管理', '技术规划'], 'closed', 1, 2, NOW() - INTERVAL '40 days', NOW()),
-- 待审核的职位 (3个)
('实习生-前端', '技术部', '深圳', '3K-5K', 'intern', '本科在读', '前端开发实习。', ARRAY['熟悉HTML/CSS/JS', '有学习热情'], ARRAY['实习补贴', '转正机会'], ARRAY['HTML', 'CSS', 'JavaScript'], 'pending', 3, 3, NOW(), NOW()),
('实习生-后端', '技术部', '深圳', '3K-5K', 'intern', '本科在读', '后端开发实习。', ARRAY['熟悉一门编程语言', '有学习热情'], ARRAY['实习补贴', '转正机会'], ARRAY['Java', 'Python', 'Go'], 'pending', 3, 3, NOW(), NOW()),
('实习生-产品', '产品部', '深圳', '3K-5K', 'intern', '本科在读', '产品实习。', ARRAY['对产品有热情', '逻辑思维强'], ARRAY['实习补贴', '转正机会'], ARRAY['产品分析', 'Axure'], 'pending', 2, 3, NOW(), NOW());

-- =====================================================
-- 3. 人才数据 (20个人才，对应20个求职者用户)
-- =====================================================
INSERT INTO talents (name, email, phone, gender, age, education, experience, current_company, current_position, salary, location, skills, status, source, user_id, created_at, updated_at) VALUES
('张伟', 'zhangwei@gmail.com', '13900000001', '男', 28, '本科', 5, '腾讯科技', 'Go开发工程师', '30K', '深圳', ARRAY['Go', 'MySQL', 'Redis', 'Docker', 'Kubernetes'], 'active', '招聘网站', 6, NOW() - INTERVAL '25 days', NOW()),
('李娜', 'lina@gmail.com', '13900000002', '女', 26, '硕士', 3, '阿里巴巴', '前端工程师', '25K', '杭州', ARRAY['Vue', 'React', 'TypeScript', 'Node.js'], 'active', '内推', 7, NOW() - INTERVAL '24 days', NOW()),
('王强', 'wangqiang@gmail.com', '13900000003', '男', 32, '本科', 8, '字节跳动', '高级后端工程师', '45K', '北京', ARRAY['Java', 'Spring', 'MySQL', 'Kafka', '微服务'], 'active', '猎头', 8, NOW() - INTERVAL '23 days', NOW()),
('刘霞', 'liuxia@gmail.com', '13900000004', '女', 27, '本科', 4, '美团', '产品经理', '28K', '北京', ARRAY['产品设计', 'Axure', '数据分析', 'SQL'], 'active', '招聘网站', 9, NOW() - INTERVAL '22 days', NOW()),
('陈龙', 'chenlong@gmail.com', '13900000005', '男', 30, '本科', 6, '华为', 'Java开发工程师', '35K', '深圳', ARRAY['Java', 'Spring Boot', 'MySQL', 'Redis'], 'active', '校招', 10, NOW() - INTERVAL '21 days', NOW()),
('赵敏', 'zhaomin@gmail.com', '13900000006', '女', 25, '本科', 2, '网易', 'UI设计师', '18K', '广州', ARRAY['Figma', 'Sketch', 'Photoshop', 'Illustrator'], 'active', '招聘网站', 11, NOW() - INTERVAL '20 days', NOW()),
('孙浩', 'sunhao@gmail.com', '13900000007', '男', 29, '硕士', 5, '百度', '数据分析师', '32K', '北京', ARRAY['Python', 'SQL', 'Tableau', 'Spark'], 'active', '内推', 12, NOW() - INTERVAL '19 days', NOW()),
('周婷', 'zhouting@gmail.com', '13900000008', '女', 24, '本科', 2, '小米', '测试工程师', '15K', '北京', ARRAY['Selenium', 'JMeter', 'Python', '自动化测试'], 'active', '招聘网站', 13, NOW() - INTERVAL '18 days', NOW()),
('吴刚', 'wugang@gmail.com', '13900000009', '男', 35, '本科', 12, '京东', '技术总监', '60K', '北京', ARRAY['架构设计', '团队管理', 'Java', 'Go'], 'active', '猎头', 14, NOW() - INTERVAL '17 days', NOW()),
('郑燕', 'zhengyan@gmail.com', '13900000010', '女', 26, '本科', 3, 'OPPO', 'Android开发', '22K', '深圳', ARRAY['Kotlin', 'Java', 'Android SDK'], 'active', '招聘网站', 15, NOW() - INTERVAL '16 days', NOW()),
('黄鹏', 'huangpeng@gmail.com', '13900000011', '男', 28, '本科', 4, 'vivo', 'iOS开发', '25K', '深圳', ARRAY['Swift', 'Objective-C', 'iOS'], 'active', '内推', 16, NOW() - INTERVAL '15 days', NOW()),
('林雪', 'linxue@gmail.com', '13900000012', '女', 27, '硕士', 4, '滴滴', 'Python开发', '30K', '北京', ARRAY['Python', 'Django', 'TensorFlow', 'PyTorch'], 'active', '招聘网站', 17, NOW() - INTERVAL '14 days', NOW()),
('何军', 'hejun@gmail.com', '13900000013', '男', 31, '本科', 7, '快手', 'DevOps工程师', '40K', '北京', ARRAY['Docker', 'Kubernetes', 'Jenkins', 'AWS', 'Terraform'], 'active', '猎头', 18, NOW() - INTERVAL '13 days', NOW()),
('高丽', 'gaoli@gmail.com', '13900000014', '女', 25, '本科', 2, '拼多多', '运营专员', '12K', '上海', ARRAY['数据分析', '用户运营', '活动策划'], 'active', '招聘网站', 19, NOW() - INTERVAL '12 days', NOW()),
('罗伟', 'luowei@gmail.com', '13900000015', '男', 29, '本科', 5, '携程', '市场专员', '18K', '上海', ARRAY['市场推广', '品牌策划', 'SEM', 'SEO'], 'active', '内推', 20, NOW() - INTERVAL '11 days', NOW()),
('谢芳', 'xiefang@gmail.com', '13900000016', '女', 28, '本科', 4, '蚂蚁金服', '财务专员', '15K', '杭州', ARRAY['会计', '财务分析', 'Excel', 'SAP'], 'active', '招聘网站', 21, NOW() - INTERVAL '10 days', NOW()),
('唐明', 'tangming@gmail.com', '13900000017', '男', 26, '本科', 3, '贝壳找房', '人事专员', '10K', '北京', ARRAY['招聘', '员工关系', '培训', 'HRIS'], 'active', '校招', 22, NOW() - INTERVAL '9 days', NOW()),
('韩晓', 'hanxiao@gmail.com', '13900000018', '女', 24, '大专', 2, '顺丰科技', '行政专员', '8K', '深圳', ARRAY['行政管理', 'Office', '会议组织'], 'active', '招聘网站', 23, NOW() - INTERVAL '8 days', NOW()),
('冯宇', 'fengyu@gmail.com', '13900000019', '男', 27, '本科', 4, '微众银行', '客服主管', '12K', '深圳', ARRAY['客户服务', '团队管理', 'CRM'], 'active', '内推', 24, NOW() - INTERVAL '7 days', NOW()),
('蔡静', 'caijing@gmail.com', '13900000020', '女', 29, '硕士', 5, '平安科技', '法务专员', '20K', '深圳', ARRAY['合同审核', '法律咨询', '知识产权'], 'active', '猎头', 25, NOW() - INTERVAL '6 days', NOW());


-- =====================================================
-- 4. 简历数据 (20份简历，对应20个人才)
-- =====================================================
INSERT INTO resumes (talent_id, file_path, file_name, file_type, parsed_data, status, created_at, updated_at) VALUES
(1, '/uploads/resumes/zhangwei_resume.pdf', '张伟_Go开发工程师.pdf', 'application/pdf', '{"name":"张伟","education":"本科","experience":"5年","skills":["Go","MySQL","Redis"]}', 'parsed', NOW() - INTERVAL '25 days', NOW()),
(2, '/uploads/resumes/lina_resume.pdf', '李娜_前端工程师.pdf', 'application/pdf', '{"name":"李娜","education":"硕士","experience":"3年","skills":["Vue","React","TypeScript"]}', 'parsed', NOW() - INTERVAL '24 days', NOW()),
(3, '/uploads/resumes/wangqiang_resume.pdf', '王强_高级后端工程师.pdf', 'application/pdf', '{"name":"王强","education":"本科","experience":"8年","skills":["Java","Spring","微服务"]}', 'parsed', NOW() - INTERVAL '23 days', NOW()),
(4, '/uploads/resumes/liuxia_resume.pdf', '刘霞_产品经理.pdf', 'application/pdf', '{"name":"刘霞","education":"本科","experience":"4年","skills":["产品设计","Axure"]}', 'parsed', NOW() - INTERVAL '22 days', NOW()),
(5, '/uploads/resumes/chenlong_resume.pdf', '陈龙_Java开发工程师.pdf', 'application/pdf', '{"name":"陈龙","education":"本科","experience":"6年","skills":["Java","Spring Boot"]}', 'parsed', NOW() - INTERVAL '21 days', NOW()),
(6, '/uploads/resumes/zhaomin_resume.pdf', '赵敏_UI设计师.pdf', 'application/pdf', '{"name":"赵敏","education":"本科","experience":"2年","skills":["Figma","Sketch"]}', 'parsed', NOW() - INTERVAL '20 days', NOW()),
(7, '/uploads/resumes/sunhao_resume.pdf', '孙浩_数据分析师.pdf', 'application/pdf', '{"name":"孙浩","education":"硕士","experience":"5年","skills":["Python","SQL","Tableau"]}', 'parsed', NOW() - INTERVAL '19 days', NOW()),
(8, '/uploads/resumes/zhouting_resume.pdf', '周婷_测试工程师.pdf', 'application/pdf', '{"name":"周婷","education":"本科","experience":"2年","skills":["Selenium","JMeter"]}', 'parsed', NOW() - INTERVAL '18 days', NOW()),
(9, '/uploads/resumes/wugang_resume.pdf', '吴刚_技术总监.pdf', 'application/pdf', '{"name":"吴刚","education":"本科","experience":"12年","skills":["架构设计","团队管理"]}', 'parsed', NOW() - INTERVAL '17 days', NOW()),
(10, '/uploads/resumes/zhengyan_resume.pdf', '郑燕_Android开发.pdf', 'application/pdf', '{"name":"郑燕","education":"本科","experience":"3年","skills":["Kotlin","Android"]}', 'parsed', NOW() - INTERVAL '16 days', NOW()),
(11, '/uploads/resumes/huangpeng_resume.pdf', '黄鹏_iOS开发.pdf', 'application/pdf', '{"name":"黄鹏","education":"本科","experience":"4年","skills":["Swift","iOS"]}', 'parsed', NOW() - INTERVAL '15 days', NOW()),
(12, '/uploads/resumes/linxue_resume.pdf', '林雪_Python开发.pdf', 'application/pdf', '{"name":"林雪","education":"硕士","experience":"4年","skills":["Python","TensorFlow"]}', 'parsed', NOW() - INTERVAL '14 days', NOW()),
(13, '/uploads/resumes/hejun_resume.pdf', '何军_DevOps工程师.pdf', 'application/pdf', '{"name":"何军","education":"本科","experience":"7年","skills":["Docker","Kubernetes"]}', 'parsed', NOW() - INTERVAL '13 days', NOW()),
(14, '/uploads/resumes/gaoli_resume.pdf', '高丽_运营专员.pdf', 'application/pdf', '{"name":"高丽","education":"本科","experience":"2年","skills":["数据分析","用户运营"]}', 'parsed', NOW() - INTERVAL '12 days', NOW()),
(15, '/uploads/resumes/luowei_resume.pdf', '罗伟_市场专员.pdf', 'application/pdf', '{"name":"罗伟","education":"本科","experience":"5年","skills":["市场推广","SEM"]}', 'parsed', NOW() - INTERVAL '11 days', NOW()),
(16, '/uploads/resumes/xiefang_resume.pdf', '谢芳_财务专员.pdf', 'application/pdf', '{"name":"谢芳","education":"本科","experience":"4年","skills":["会计","财务分析"]}', 'parsed', NOW() - INTERVAL '10 days', NOW()),
(17, '/uploads/resumes/tangming_resume.pdf', '唐明_人事专员.pdf', 'application/pdf', '{"name":"唐明","education":"本科","experience":"3年","skills":["招聘","员工关系"]}', 'parsed', NOW() - INTERVAL '9 days', NOW()),
(18, '/uploads/resumes/hanxiao_resume.pdf', '韩晓_行政专员.pdf', 'application/pdf', '{"name":"韩晓","education":"大专","experience":"2年","skills":["行政管理","Office"]}', 'parsed', NOW() - INTERVAL '8 days', NOW()),
(19, '/uploads/resumes/fengyu_resume.pdf', '冯宇_客服主管.pdf', 'application/pdf', '{"name":"冯宇","education":"本科","experience":"4年","skills":["客户服务","团队管理"]}', 'parsed', NOW() - INTERVAL '7 days', NOW()),
(20, '/uploads/resumes/caijing_resume.pdf', '蔡静_法务专员.pdf', 'application/pdf', '{"name":"蔡静","education":"硕士","experience":"5年","skills":["合同审核","法律咨询"]}', 'parsed', NOW() - INTERVAL '6 days', NOW());


-- =====================================================
-- 5. 申请数据 (30个申请，分布在不同职位和状态)
-- =====================================================
INSERT INTO applications (job_id, talent_id, resume_id, status, created_at, updated_at) VALUES
-- 高级Go后端工程师 (job_id=1) - 5个申请
(1, 1, 1, 'interview', NOW() - INTERVAL '20 days', NOW()),
(1, 3, 3, 'offer', NOW() - INTERVAL '18 days', NOW()),
(1, 5, 5, 'pending', NOW() - INTERVAL '5 days', NOW()),
-- 前端开发工程师 (job_id=2) - 4个申请
(2, 2, 2, 'interview', NOW() - INTERVAL '19 days', NOW()),
(2, 10, 10, 'rejected', NOW() - INTERVAL '15 days', NOW()),
(2, 11, 11, 'pending', NOW() - INTERVAL '3 days', NOW()),
(2, 12, 12, 'pending', NOW() - INTERVAL '2 days', NOW()),
-- 产品经理 (job_id=3) - 3个申请
(3, 4, 4, 'offer', NOW() - INTERVAL '15 days', NOW()),
(3, 7, 7, 'interview', NOW() - INTERVAL '10 days', NOW()),
(3, 14, 14, 'pending', NOW() - INTERVAL '4 days', NOW()),
-- UI设计师 (job_id=4) - 2个申请
(4, 6, 6, 'interview', NOW() - INTERVAL '12 days', NOW()),
(4, 2, 2, 'rejected', NOW() - INTERVAL '10 days', NOW()),
-- 数据分析师 (job_id=5) - 3个申请
(5, 7, 7, 'offer', NOW() - INTERVAL '10 days', NOW()),
(5, 12, 12, 'interview', NOW() - INTERVAL '8 days', NOW()),
(5, 4, 4, 'pending', NOW() - INTERVAL '3 days', NOW()),
-- Java开发工程师 (job_id=6) - 4个申请
(6, 5, 5, 'interview', NOW() - INTERVAL '8 days', NOW()),
(6, 3, 3, 'offer', NOW() - INTERVAL '6 days', NOW()),
(6, 1, 1, 'rejected', NOW() - INTERVAL '5 days', NOW()),
(6, 9, 9, 'pending', NOW() - INTERVAL '2 days', NOW()),
-- Python开发工程师 (job_id=7) - 3个申请
(7, 12, 12, 'interview', NOW() - INTERVAL '7 days', NOW()),
(7, 7, 7, 'pending', NOW() - INTERVAL '4 days', NOW()),
(7, 3, 3, 'pending', NOW() - INTERVAL '2 days', NOW()),
-- DevOps工程师 (job_id=8) - 2个申请
(8, 13, 13, 'offer', NOW() - INTERVAL '5 days', NOW()),
(8, 1, 1, 'interview', NOW() - INTERVAL '3 days', NOW()),
-- 测试工程师 (job_id=9) - 2个申请
(9, 8, 8, 'interview', NOW() - INTERVAL '4 days', NOW()),
(9, 2, 2, 'pending', NOW() - INTERVAL '2 days', NOW()),
-- Android开发工程师 (job_id=10) - 1个申请
(10, 10, 10, 'interview', NOW() - INTERVAL '3 days', NOW()),
-- iOS开发工程师 (job_id=11) - 1个申请
(11, 11, 11, 'offer', NOW() - INTERVAL '2 days', NOW());


-- =====================================================
-- 6. 面试数据 (15个面试，不同状态)
-- =====================================================
INSERT INTO interviews (candidate_id, candidate_name, position_id, position, type, date, time, interviewer, method, location, status, feedback, rating, created_at, updated_at) VALUES
-- 已完成的面试
(1, '张伟', 1, '高级Go后端工程师', 'technical', '2026-01-15', '10:00', '张丽华', 'onsite', '深圳总部3楼会议室A', 'completed', '技术能力扎实，Go语言掌握熟练，有微服务经验，沟通能力良好。建议进入下一轮。', 4, NOW() - INTERVAL '15 days', NOW()),
(3, '王强', 1, '高级Go后端工程师', 'final', '2026-01-18', '14:00', '李明', 'onsite', '深圳总部5楼会议室B', 'completed', '架构设计能力强，有大型项目经验，领导力突出。强烈推荐录用。', 5, NOW() - INTERVAL '12 days', NOW()),
(2, '李娜', 2, '前端开发工程师', 'technical', '2026-01-16', '09:30', '王芳', 'video', '腾讯会议', 'completed', 'Vue3掌握熟练，有TypeScript经验，代码规范。建议录用。', 4, NOW() - INTERVAL '14 days', NOW()),
(4, '刘霞', 3, '产品经理', 'final', '2026-01-20', '15:00', '赵强', 'onsite', '北京分部2楼会议室', 'completed', '产品思维清晰，有B端经验，数据分析能力强。推荐录用。', 5, NOW() - INTERVAL '10 days', NOW()),
(7, '孙浩', 5, '数据分析师', 'technical', '2026-01-22', '10:30', '张丽华', 'video', 'Zoom', 'completed', 'SQL能力强，Python熟练，有Spark经验。推荐进入终面。', 4, NOW() - INTERVAL '8 days', NOW()),
-- 已安排的面试
(6, '赵敏', 4, 'UI设计师', 'technical', '2026-01-31', '14:00', '李明', 'onsite', '上海分部设计中心', 'scheduled', NULL, 0, NOW() - INTERVAL '5 days', NOW()),
(5, '陈龙', 6, 'Java开发工程师', 'technical', '2026-02-01', '10:00', '王芳', 'video', '腾讯会议', 'scheduled', NULL, 0, NOW() - INTERVAL '4 days', NOW()),
(12, '林雪', 7, 'Python开发工程师', 'technical', '2026-02-02', '09:00', '赵强', 'onsite', '深圳总部3楼会议室C', 'scheduled', NULL, 0, NOW() - INTERVAL '3 days', NOW()),
(8, '周婷', 9, '测试工程师', 'technical', '2026-02-03', '14:30', '张丽华', 'video', 'Zoom', 'scheduled', NULL, 0, NOW() - INTERVAL '2 days', NOW()),
(10, '郑燕', 10, 'Android开发工程师', 'technical', '2026-02-04', '10:00', '李明', 'onsite', '深圳总部4楼会议室', 'scheduled', NULL, 0, NOW() - INTERVAL '1 day', NOW()),
(13, '何军', 8, 'DevOps工程师', 'final', '2026-02-05', '15:00', '王芳', 'onsite', '深圳总部5楼会议室A', 'scheduled', NULL, 0, NOW(), NOW()),
(1, '张伟', 8, 'DevOps工程师', 'technical', '2026-02-06', '09:30', '赵强', 'video', '腾讯会议', 'scheduled', NULL, 0, NOW(), NOW()),
-- 已取消的面试
(9, '吴刚', 19, '高级架构师', 'final', '2026-01-10', '14:00', '张丽华', 'onsite', '深圳总部6楼', 'cancelled', '职位已关闭', 0, NOW() - INTERVAL '25 days', NOW()),
-- 待确认的面试
(7, '孙浩', 3, '产品经理', 'technical', '2026-02-07', '10:00', '李明', 'video', 'Zoom', 'pending', NULL, 0, NOW(), NOW()),
(14, '高丽', 12, '运营专员', 'initial', '2026-02-08', '14:00', '王芳', 'phone', '电话面试', 'pending', NULL, 0, NOW(), NOW());


-- =====================================================
-- 7. 在线简历数据 (为部分求职者创建在线简历)
-- =====================================================
INSERT INTO online_resumes (user_id, talent_id, name, phone, email, location, gender, age, summary, work_experience, education, skills, is_complete, created_at, updated_at) VALUES
(6, 1, '张伟', '13900000001', 'zhangwei@gmail.com', '深圳', '男', 28, '5年Go开发经验，熟悉微服务架构，有大型分布式系统开发经验。',
 '[{"company":"腾讯科技","position":"Go开发工程师","startDate":"2021-03","endDate":"至今","description":"负责核心业务系统开发"},{"company":"华为","position":"后端开发","startDate":"2019-07","endDate":"2021-02","description":"参与云服务平台开发"}]',
 '[{"school":"华南理工大学","major":"计算机科学与技术","degree":"本科","startDate":"2015-09","endDate":"2019-06"}]',
 '[{"name":"Go","level":5},{"name":"MySQL","level":4},{"name":"Redis","level":4},{"name":"Docker","level":4}]',
 true, NOW() - INTERVAL '20 days', NOW()),
(7, 2, '李娜', '13900000002', 'lina@gmail.com', '杭州', '女', 26, '3年前端开发经验，精通Vue3和React，有大型项目经验。',
 '[{"company":"阿里巴巴","position":"前端工程师","startDate":"2022-06","endDate":"至今","description":"负责电商平台前端开发"}]',
 '[{"school":"浙江大学","major":"软件工程","degree":"硕士","startDate":"2020-09","endDate":"2022-06"}]',
 '[{"name":"Vue","level":5},{"name":"React","level":4},{"name":"TypeScript","level":4}]',
 true, NOW() - INTERVAL '18 days', NOW()),
(8, 3, '王强', '13900000003', 'wangqiang@gmail.com', '北京', '男', 32, '8年Java开发经验，有架构设计能力，带过10人团队。',
 '[{"company":"字节跳动","position":"高级后端工程师","startDate":"2020-01","endDate":"至今","description":"负责推荐系统后端架构"}]',
 '[{"school":"北京大学","major":"计算机科学","degree":"本科","startDate":"2012-09","endDate":"2016-06"}]',
 '[{"name":"Java","level":5},{"name":"Spring","level":5},{"name":"微服务","level":4}]',
 true, NOW() - INTERVAL '15 days', NOW()),
(9, 4, '刘霞', '13900000004', 'liuxia@gmail.com', '北京', '女', 27, '4年产品经理经验，专注B端产品，有完整产品生命周期管理经验。',
 '[{"company":"美团","position":"产品经理","startDate":"2021-07","endDate":"至今","description":"负责商家端产品规划"}]',
 '[{"school":"清华大学","major":"工商管理","degree":"本科","startDate":"2016-09","endDate":"2020-06"}]',
 '[{"name":"产品设计","level":5},{"name":"数据分析","level":4},{"name":"Axure","level":4}]',
 true, NOW() - INTERVAL '12 days', NOW()),
(12, 7, '孙浩', '13900000007', 'sunhao@gmail.com', '北京', '男', 29, '5年数据分析经验，精通Python和SQL，有大数据处理经验。',
 '[{"company":"百度","position":"数据分析师","startDate":"2021-03","endDate":"至今","description":"负责搜索业务数据分析"}]',
 '[{"school":"中国人民大学","major":"统计学","degree":"硕士","startDate":"2019-09","endDate":"2021-06"}]',
 '[{"name":"Python","level":5},{"name":"SQL","level":5},{"name":"Tableau","level":4}]',
 true, NOW() - INTERVAL '10 days', NOW());

-- =====================================================
-- 8. 会话和消息数据 (创建一些示例对话)
-- =====================================================
INSERT INTO conversations (participant_a, participant_b, last_message_at, created_at, updated_at) VALUES
(2, 6, NOW() - INTERVAL '1 hour', NOW() - INTERVAL '5 days', NOW()),
(2, 7, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '4 days', NOW()),
(3, 8, NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 days', NOW()),
(4, 9, NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', NOW());

INSERT INTO chat_messages (conversation_id, sender_id, content, message_type, is_read, created_at) VALUES
-- HR2 和 张伟 的对话
(1, 2, '您好，我是HR张丽华，看到您投递了高级Go后端工程师职位，想和您聊聊。', 'text', true, NOW() - INTERVAL '5 days'),
(1, 6, '您好，很高兴收到您的消息，我对这个职位很感兴趣。', 'text', true, NOW() - INTERVAL '5 days' + INTERVAL '10 minutes'),
(1, 2, '好的，请问您方便明天上午10点来公司面试吗？', 'text', true, NOW() - INTERVAL '4 days'),
(1, 6, '可以的，请问公司地址在哪里？', 'text', true, NOW() - INTERVAL '4 days' + INTERVAL '5 minutes'),
(1, 2, '深圳市南山区科技园3楼，到时候联系前台即可。', 'text', false, NOW() - INTERVAL '1 hour'),
-- HR2 和 李娜 的对话
(2, 2, '您好，我们正在招聘前端工程师，您的简历很匹配，方便聊聊吗？', 'text', true, NOW() - INTERVAL '4 days'),
(2, 7, '您好，可以的，请问有什么想了解的？', 'text', true, NOW() - INTERVAL '4 days' + INTERVAL '15 minutes'),
(2, 2, '主要想了解一下您的Vue3项目经验。', 'text', false, NOW() - INTERVAL '2 hours'),
-- HR3 和 王强 的对话
(3, 3, '王先生您好，我们有一个高级后端职位，薪资可以谈到45K，您有兴趣吗？', 'text', true, NOW() - INTERVAL '3 days'),
(3, 8, '您好，可以详细介绍一下职位要求吗？', 'text', false, NOW() - INTERVAL '3 hours');

-- 更新会话的最后消息ID
UPDATE conversations SET last_message_id = 5 WHERE id = 1;
UPDATE conversations SET last_message_id = 8 WHERE id = 2;
UPDATE conversations SET last_message_id = 10 WHERE id = 3;

-- =====================================================
-- 完成
-- =====================================================
SELECT '数据初始化完成！' as message;
SELECT '用户数量: ' || COUNT(*) FROM users;
SELECT '职位数量: ' || COUNT(*) FROM jobs;
SELECT '人才数量: ' || COUNT(*) FROM talents;
SELECT '简历数量: ' || COUNT(*) FROM resumes;
SELECT '申请数量: ' || COUNT(*) FROM applications;
SELECT '面试数量: ' || COUNT(*) FROM interviews;
