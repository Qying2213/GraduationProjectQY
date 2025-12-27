-- =====================================================
-- 4. 人才数据
-- =====================================================
INSERT INTO talents (name, email, phone, skills, experience, education, status, tags, location, salary, summary, gender, age, current_company, current_position, source) VALUES
('张伟', 'zhangwei@email.com', '13900000001', ARRAY['Vue3', 'TypeScript', 'React', 'Node.js'], 6, '本科', 'active', ARRAY['高级前端', '架构能力'], '北京', '40-50K', '6年前端开发经验，擅长Vue和React', '男', 28, '字节跳动', '高级前端工程师', '猎聘'),
('李娜', 'lina@email.com', '13900000002', ARRAY['Go', 'gRPC', 'Kubernetes', 'Docker'], 4, '硕士', 'active', ARRAY['后端开发', '微服务'], '北京', '35-45K', '4年Go开发经验，熟悉云原生', '女', 27, '阿里巴巴', '后端工程师', 'Boss直聘'),
('王强', 'wangqiang@email.com', '13900000003', ARRAY['Java', 'Spring Boot', 'MySQL', 'Redis'], 5, '本科', 'active', ARRAY['Java开发', '中级'], '成都', '30-40K', '5年Java开发，有电商项目经验', '男', 29, '京东', 'Java开发工程师', '拉勾'),
('刘芳', 'liufang@email.com', '13900000004', ARRAY['产品设计', 'Axure', 'SQL', '数据分析'], 4, '本科', 'active', ARRAY['产品经理', 'B端'], '上海', '30-40K', '4年B端产品经验，擅长数据产品', '女', 28, '美团', '产品经理', '猎聘'),
('陈明', 'chenming@email.com', '13900000005', ARRAY['Figma', 'Sketch', 'UI设计', '交互设计'], 3, '本科', 'active', ARRAY['UI设计', '视觉设计'], '深圳', '25-35K', '3年UI设计经验，有B端设计作品', '男', 26, '腾讯', 'UI设计师', 'Boss直聘'),
('赵丽', 'zhaoli@email.com', '13900000006', ARRAY['Python', 'SQL', 'Tableau', '数据分析'], 2, '硕士', 'active', ARRAY['数据分析', '统计'], '杭州', '20-30K', '2年数据分析经验，统计学硕士', '女', 25, '网易', '数据分析师', '智联'),
('孙浩', 'sunhao@email.com', '13900000007', ARRAY['Linux', 'Docker', 'Kubernetes', 'Jenkins'], 4, '本科', 'active', ARRAY['DevOps', '运维'], '北京', '30-40K', '4年运维经验，熟悉K8s', '男', 28, '百度', 'DevOps工程师', '拉勾'),
('周婷', 'zhouting@email.com', '13900000008', ARRAY['Selenium', 'JMeter', 'Python', '自动化测试'], 3, '本科', 'active', ARRAY['测试', '自动化'], '广州', '18-25K', '3年测试经验，会自动化测试', '女', 27, '华为', '测试工程师', 'Boss直聘'),
('吴杰', 'wujie@email.com', '13900000009', ARRAY['Vue', 'JavaScript', 'HTML', 'CSS'], 1, '本科', 'pending', ARRAY['前端', '应届'], '北京', '8-12K', '应届毕业生，有Vue项目经验', '男', 23, NULL, NULL, '校招'),
('郑雪', 'zhengxue@email.com', '13900000010', ARRAY['招聘', '人力资源', 'Excel', '沟通'], 2, '本科', 'active', ARRAY['HR', '招聘'], '北京', '12-18K', '2年HR经验，熟悉招聘流程', '女', 26, '滴滴', 'HR专员', '智联'),
('黄磊', 'huanglei@email.com', '13900000011', ARRAY['Go', 'Python', 'MySQL', 'Redis'], 3, '本科', 'hired', ARRAY['后端', '全栈'], '上海', '25-35K', '3年后端开发，Go和Python都熟', '男', 27, '拼多多', '后端工程师', '猎聘'),
('杨梅', 'yangmei@email.com', '13900000012', ARRAY['React', 'TypeScript', 'Node.js', 'GraphQL'], 5, '硕士', 'active', ARRAY['前端', '全栈'], '北京', '40-55K', '5年前端经验，有全栈能力', '女', 29, '小米', '高级前端工程师', 'Boss直聘'),
('徐涛', 'xutao@email.com', '13900000013', ARRAY['Java', 'Spring Cloud', '微服务', 'Kafka'], 7, '本科', 'active', ARRAY['架构师', '后端'], '杭州', '50-70K', '7年Java经验，有架构设计能力', '男', 32, '蚂蚁金服', '技术专家', '猎聘'),
('朱琳', 'zhulin@email.com', '13900000014', ARRAY['产品设计', '用户研究', 'Figma', '数据分析'], 6, '硕士', 'active', ARRAY['产品总监', '高级'], '北京', '45-60K', '6年产品经验，带过10人团队', '女', 31, '快手', '产品总监', '猎头'),
('马超', 'machao@email.com', '13900000015', ARRAY['Python', 'TensorFlow', 'PyTorch', '机器学习'], 4, '博士', 'active', ARRAY['AI', '算法'], '北京', '50-80K', 'AI算法博士，有NLP项目经验', '男', 30, '商汤科技', '算法工程师', '猎头'),
('何静', 'hejing@email.com', '13900000016', ARRAY['Vue', 'React', '小程序', 'Flutter'], 4, '本科', 'rejected', ARRAY['前端', '移动端'], '深圳', '30-40K', '4年前端经验，会小程序和Flutter', '女', 28, '大疆', '前端工程师', 'Boss直聘'),
('林峰', 'linfeng@email.com', '13900000017', ARRAY['Java', 'Go', 'MySQL', 'MongoDB'], 5, '本科', 'active', ARRAY['后端', '数据库'], '成都', '35-45K', '5年后端经验，Java和Go双修', '男', 29, '字节跳动', '后端工程师', '拉勾'),
('谢敏', 'xiemin@email.com', '13900000018', ARRAY['UI设计', 'UX设计', 'Figma', 'After Effects'], 5, '本科', 'active', ARRAY['设计', '动效'], '上海', '30-40K', '5年设计经验，擅长动效设计', '女', 28, '哔哩哔哩', '高级设计师', '猎聘'),
('罗刚', 'luogang@email.com', '13900000019', ARRAY['测试', '性能测试', 'LoadRunner', 'Python'], 6, '本科', 'active', ARRAY['测试', '性能'], '广州', '25-35K', '6年测试经验，专注性能测试', '男', 30, '唯品会', '测试专家', '智联'),
('唐欣', 'tangxin@email.com', '13900000020', ARRAY['运营', '数据分析', '用户增长', '活动策划'], 3, '本科', 'hired', ARRAY['运营', '增长'], '杭州', '15-22K', '3年运营经验，擅长用户增长', '女', 26, '网易', '运营专员', 'Boss直聘');
