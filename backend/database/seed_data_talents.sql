-- 追加人才数据 (talents表)
-- 使用 WHERE NOT EXISTS 避免重复

-- 后端开发人才
INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '张伟', 'zhangwei_dev@example.com', '13900139001', '{"Go","Docker","Kubernetes","Redis","MySQL","微服务","gRPC"}', 6, '本科', 'active', '北京', '35-45K', '6年后端开发经验，精通Go语言，有丰富的微服务架构经验。', NOW() - INTERVAL '15 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'zhangwei_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '李强', 'liqiang_dev@example.com', '13900139002', '{"Go","Python","PostgreSQL","Redis","Docker","Gin"}', 4, '硕士', 'active', '上海', '28-38K', '4年后端开发经验，Go和Python双语言开发者。', NOW() - INTERVAL '12 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'liqiang_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '王磊', 'wanglei_dev@example.com', '13900139003', '{"Java","Spring Boot","Spring Cloud","MySQL","Redis","Kafka"}', 5, '本科', 'active', '杭州', '30-40K', '5年Java开发经验，精通Spring全家桶。', NOW() - INTERVAL '10 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'wanglei_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '陈浩', 'chenhao_dev@example.com', '13900139004', '{"Go","Rust","C++","Linux","网络编程","高性能"}', 8, '硕士', 'active', '北京', '50-70K', '8年系统开发经验，精通Go和Rust，有基础架构开发经验。', NOW() - INTERVAL '20 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'chenhao_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '刘洋', 'liuyang_dev@example.com', '13900139005', '{"Python","Django","Flask","PostgreSQL","Celery","Redis"}', 3, '本科', 'active', '深圳', '22-30K', '3年Python后端开发经验。', NOW() - INTERVAL '8 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'liuyang_dev@example.com');

-- 前端开发人才
INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '赵雪', 'zhaoxue_dev@example.com', '13900139006', '{"Vue","TypeScript","React","Webpack","Node.js","Element Plus"}', 5, '本科', 'active', '北京', '30-40K', '5年前端开发经验，精通Vue和React。', NOW() - INTERVAL '18 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'zhaoxue_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '孙婷', 'sunting_dev@example.com', '13900139007', '{"React","TypeScript","Redux","Ant Design","Webpack"}', 4, '本科', 'active', '上海', '25-35K', '4年前端开发经验，专注React技术栈。', NOW() - INTERVAL '14 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'sunting_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '周明', 'zhouming_dev@example.com', '13900139008', '{"Vue","JavaScript","CSS","Element Plus","Vite","Git"}', 2, '本科', 'active', '广州', '15-22K', '2年前端开发经验，熟悉Vue框架。', NOW() - INTERVAL '6 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'zhouming_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '吴芳', 'wufang_dev@example.com', '13900139009', '{"Vue","React","TypeScript","小程序","Taro","移动端"}', 3, '硕士', 'active', '深圳', '22-30K', '3年前端开发经验，熟悉多端开发。', NOW() - INTERVAL '9 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'wufang_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '董丽', 'dongli_dev@example.com', '13900139010', '{"Vue","React","微前端","Webpack","性能优化"}', 6, '本科', 'active', '上海', '35-50K', '6年前端开发经验，有微前端架构经验。', NOW() - INTERVAL '16 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'dongli_dev@example.com');

-- DevOps人才
INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '郑凯', 'zhengkai_dev@example.com', '13900139011', '{"Vue","Go","PostgreSQL","Docker","Redis","Nginx"}', 4, '本科', 'active', '深圳', '28-38K', '4年全栈开发经验。', NOW() - INTERVAL '11 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'zhengkai_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '马超', 'machao_dev@example.com', '13900139012', '{"Docker","Kubernetes","Jenkins","Ansible","Linux","Shell","Prometheus"}', 5, '本科', 'active', '北京', '30-45K', '5年运维开发经验，精通容器化。', NOW() - INTERVAL '13 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'machao_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '林峰', 'linfeng_dev@example.com', '13900139013', '{"Kubernetes","Terraform","AWS","GCP","CI/CD","Go"}', 4, '硕士', 'active', '上海', '35-50K', '4年SRE经验，熟悉云原生技术栈。', NOW() - INTERVAL '17 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'linfeng_dev@example.com');

-- 数据/AI人才
INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '杨帆', 'yangfan_dev@example.com', '13900139015', '{"Python","SQL","Tableau","Excel","数据分析","统计学"}', 3, '硕士', 'active', '杭州', '20-30K', '3年数据分析经验。', NOW() - INTERVAL '7 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'yangfan_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '何伟', 'hewei_dev@example.com', '13900139017', '{"Python","PyTorch","TensorFlow","机器学习","深度学习","NLP","CV"}', 4, '博士', 'active', '北京', '50-70K', '4年算法研发经验，专注NLP方向。', NOW() - INTERVAL '22 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'hewei_dev@example.com');

-- 产品/设计人才
INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '谢敏', 'xiemin_dev@example.com', '13900139019', '{"产品设计","需求分析","Axure","数据分析","用户研究"}', 4, '本科', 'active', '北京', '28-40K', '4年产品经理经验。', NOW() - INTERVAL '10 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'xiemin_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '罗琳', 'luolin_dev@example.com', '13900139021', '{"Figma","Sketch","Photoshop","设计系统","UI设计"}', 3, '本科', 'active', '深圳', '18-28K', '3年UI设计经验。', NOW() - INTERVAL '8 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'luolin_dev@example.com');

-- 测试人才
INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '韩冰', 'hanbing_dev@example.com', '13900139024', '{"Selenium","Python","接口测试","性能测试","自动化测试","JMeter"}', 4, '本科', 'active', '上海', '22-32K', '4年测试开发经验。', NOW() - INTERVAL '15 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'hanbing_dev@example.com');

-- 更多人才
INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '冯刚', 'fenggang_dev@example.com', '13900139025', '{"Go","gRPC","etcd","分布式系统"}', 7, '硕士', 'active', '北京', '45-60K', '7年后端开发经验，专注分布式系统。', NOW() - INTERVAL '20 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'fenggang_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '曹阳', 'caoyang_dev@example.com', '13900139028', '{"Java","Spring","MySQL"}', 2, '本科', 'pending', '成都', '15-22K', '2年Java开发经验。', NOW() - INTERVAL '3 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'caoyang_dev@example.com');

INSERT INTO talents (name, email, phone, skills, experience, education, status, location, salary, summary, created_at, updated_at)
SELECT '邓超', 'dengchao_dev@example.com', '13900139030', '{"Python","数据分析"}', 1, '本科', 'active', '武汉', '10-15K', '1年数据分析经验。', NOW() - INTERVAL '4 days', NOW()
WHERE NOT EXISTS (SELECT 1 FROM talents WHERE email = 'dengchao_dev@example.com');
