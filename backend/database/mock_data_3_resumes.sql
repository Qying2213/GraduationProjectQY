-- =====================================================
-- 5. 简历数据 (关联人才和职位)
-- =====================================================
INSERT INTO resumes (talent_id, job_id, file_path, file_name, status, match_score, parse_result) VALUES
(1, 1, 'uploads/resumes/zhangwei_resume.pdf', '张伟_高级前端工程师简历.pdf', 'interviewed', 92, '{"skills": ["Vue3", "TypeScript", "React"], "experience": 6}'),
(2, 2, 'uploads/resumes/lina_resume.pdf', '李娜_Go开发工程师简历.pdf', 'offered', 88, '{"skills": ["Go", "gRPC", "Kubernetes"], "experience": 4}'),
(3, 8, 'uploads/resumes/wangqiang_resume.pdf', '王强_Java开发工程师简历.pdf', 'reviewing', 85, '{"skills": ["Java", "Spring Boot"], "experience": 5}'),
(4, 3, 'uploads/resumes/liufang_resume.pdf', '刘芳_产品经理简历.pdf', 'interviewed', 90, '{"skills": ["产品设计", "数据分析"], "experience": 4}'),
(5, 4, 'uploads/resumes/chenming_resume.pdf', '陈明_UI设计师简历.pdf', 'pending', 82, '{"skills": ["Figma", "Sketch"], "experience": 3}'),
(6, 5, 'uploads/resumes/zhaoli_resume.pdf', '赵丽_数据分析师简历.pdf', 'reviewing', 87, '{"skills": ["Python", "SQL", "Tableau"], "experience": 2}'),
(7, 6, 'uploads/resumes/sunhao_resume.pdf', '孙浩_DevOps工程师简历.pdf', 'interviewed', 91, '{"skills": ["Docker", "Kubernetes"], "experience": 4}'),
(8, 7, 'uploads/resumes/zhouting_resume.pdf', '周婷_测试工程师简历.pdf', 'pending', 78, '{"skills": ["Selenium", "JMeter"], "experience": 3}'),
(9, 9, 'uploads/resumes/wujie_resume.pdf', '吴杰_前端实习生简历.pdf', 'pending', 65, '{"skills": ["Vue", "JavaScript"], "experience": 0}'),
(10, 10, 'uploads/resumes/zhengxue_resume.pdf', '郑雪_HR专员简历.pdf', 'reviewing', 80, '{"skills": ["招聘", "人力资源"], "experience": 2}'),
(11, 2, 'uploads/resumes/huanglei_resume.pdf', '黄磊_后端工程师简历.pdf', 'hired', 86, '{"skills": ["Go", "Python"], "experience": 3}'),
(12, 1, 'uploads/resumes/yangmei_resume.pdf', '杨梅_高级前端简历.pdf', 'reviewing', 94, '{"skills": ["React", "TypeScript"], "experience": 5}'),
(13, 2, 'uploads/resumes/xutao_resume.pdf', '徐涛_技术专家简历.pdf', 'interviewed', 95, '{"skills": ["Java", "微服务"], "experience": 7}'),
(14, 3, 'uploads/resumes/zhulin_resume.pdf', '朱琳_产品总监简历.pdf', 'offered', 93, '{"skills": ["产品设计", "团队管理"], "experience": 6}'),
(15, 11, 'uploads/resumes/machao_resume.pdf', '马超_算法工程师简历.pdf', 'pending', 88, '{"skills": ["Python", "机器学习"], "experience": 4}'),
(16, 1, 'uploads/resumes/hejing_resume.pdf', '何静_前端工程师简历.pdf', 'rejected', 72, '{"skills": ["Vue", "React"], "experience": 4}'),
(17, 8, 'uploads/resumes/linfeng_resume.pdf', '林峰_后端工程师简历.pdf', 'reviewing', 84, '{"skills": ["Java", "Go"], "experience": 5}'),
(18, 4, 'uploads/resumes/xiemin_resume.pdf', '谢敏_高级设计师简历.pdf', 'interviewed', 89, '{"skills": ["UI设计", "动效"], "experience": 5}'),
(19, 7, 'uploads/resumes/luogang_resume.pdf', '罗刚_测试专家简历.pdf', 'offered', 90, '{"skills": ["性能测试", "自动化"], "experience": 6}'),
(20, 12, 'uploads/resumes/tangxin_resume.pdf', '唐欣_运营专员简历.pdf', 'hired', 85, '{"skills": ["运营", "数据分析"], "experience": 3}');

-- =====================================================
-- 6. 应聘记录数据
-- =====================================================
INSERT INTO applications (talent_id, job_id, resume_id, stage, status, source, notes) VALUES
(1, 1, 1, 'interview', 'active', '猎聘', '技术面试通过，等待HR面'),
(2, 2, 2, 'offer', 'active', 'Boss直聘', '已发offer，等待回复'),
(3, 8, 3, 'screening', 'active', '拉勾', '简历筛选中'),
(4, 3, 4, 'interview', 'active', '猎聘', '产品总监面试中'),
(5, 4, 5, 'applied', 'active', 'Boss直聘', '刚投递'),
(6, 5, 6, 'screening', 'active', '智联', '简历评估中'),
(7, 6, 7, 'interview', 'active', '拉勾', '技术面试安排中'),
(8, 7, 8, 'applied', 'active', 'Boss直聘', '待筛选'),
(9, 9, 9, 'applied', 'active', '校招', '实习生候选'),
(10, 10, 10, 'screening', 'active', '智联', 'HR初筛中'),
(11, 2, 11, 'hired', 'closed', '猎聘', '已入职'),
(12, 1, 12, 'screening', 'active', 'Boss直聘', '高匹配度候选人'),
(13, 2, 13, 'interview', 'active', '猎聘', '架构师候选人'),
(14, 3, 14, 'offer', 'active', '猎头', '产品总监offer'),
(15, 11, 15, 'applied', 'active', '猎头', 'AI人才'),
(16, 1, 16, 'rejected', 'closed', 'Boss直聘', '技术面试未通过'),
(17, 8, 17, 'screening', 'active', '拉勾', '双语言候选人'),
(18, 4, 18, 'interview', 'active', '猎聘', '设计作品优秀'),
(19, 7, 19, 'offer', 'active', '智联', '测试专家offer'),
(20, 12, 20, 'hired', 'closed', 'Boss直聘', '已入职');
