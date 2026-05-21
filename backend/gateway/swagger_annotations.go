package main

// SwaggerAny 表示未展开的通用请求体。
type SwaggerAny map[string]interface{}

// SwaggerResponse 是平台接口的通用成功响应结构。
type SwaggerResponse struct {
	Code    int         `json:"code" example:"0"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data"`
}

// SwaggerErrorResponse 是平台接口的通用错误响应结构。
type SwaggerErrorResponse struct {
	Code    int    `json:"code" example:"1"`
	Message string `json:"message" example:"error message"`
	Error   string `json:"error,omitempty"`
}

// RegisterRequest 是用户注册请求。
type RegisterRequest struct {
	Username string `json:"username" example:"candidate01"`
	Email    string `json:"email" example:"candidate01@example.com"`
	Password string `json:"password" example:"123456"`
	Role     string `json:"role" example:"candidate" enums:"admin,hr,hr_manager,recruiter,candidate"`
	Name     string `json:"name" example:"张三"`
	Phone    string `json:"phone" example:"13800000000"`
}

// LoginRequest 是用户登录请求。
type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin123"`
}

// ProfileUpdateRequest 是当前用户资料更新请求。
type ProfileUpdateRequest struct {
	Name   string `json:"name" example:"管理员"`
	Phone  string `json:"phone" example:"13800000000"`
	Avatar string `json:"avatar" example:"https://example.com/avatar.png"`
	Email  string `json:"email" example:"admin@example.com"`
}

// TalentUpsertRequest 是人才库新增或更新请求。
type TalentUpsertRequest struct {
	Name        string   `json:"name" example:"李四"`
	Email       string   `json:"email" example:"lisi@example.com"`
	Phone       string   `json:"phone" example:"13900000000"`
	Position    string   `json:"position" example:"Java 后端工程师"`
	Skills      []string `json:"skills" example:"Go,PostgreSQL,Vue"`
	Experience  int      `json:"experience" example:"3"`
	Education   string   `json:"education" example:"本科"`
	Location    string   `json:"location" example:"成都"`
	Status      string   `json:"status" example:"active"`
	Source      string   `json:"source" example:"校园招聘"`
	Description string   `json:"description" example:"具备后端开发和微服务项目经验"`
}

// JobUpsertRequest 是职位新增或更新请求。
type JobUpsertRequest struct {
	Title        string   `json:"title" example:"后端开发工程师"`
	Department   string   `json:"department" example:"研发部"`
	Location     string   `json:"location" example:"成都"`
	Type         string   `json:"type" example:"full-time"`
	Status       string   `json:"status" example:"open"`
	SalaryMin    int      `json:"salary_min" example:"12000"`
	SalaryMax    int      `json:"salary_max" example:"20000"`
	Experience   string   `json:"experience" example:"3-5年"`
	Education    string   `json:"education" example:"本科"`
	Skills       []string `json:"skills" example:"Go,Redis,PostgreSQL"`
	Description  string   `json:"description" example:"负责招聘平台后端服务开发"`
	Requirements string   `json:"requirements" example:"熟悉微服务、数据库和缓存"`
}

// OnlineResumeRequest 是在线简历保存请求。
type OnlineResumeRequest struct {
	Name       string   `json:"name" example:"王五"`
	Phone      string   `json:"phone" example:"13700000000"`
	Email      string   `json:"email" example:"wangwu@example.com"`
	Education  string   `json:"education" example:"本科"`
	Experience string   `json:"experience" example:"2年 Java 开发经验"`
	Skills     []string `json:"skills" example:"Java,Spring Boot,MySQL"`
	Projects   string   `json:"projects" example:"智能招聘平台、在线简历系统"`
	Summary    string   `json:"summary" example:"具备后端开发和项目协作经验"`
}

// ResumeTextRequest 是简历文本解析请求。
type ResumeTextRequest struct {
	Text string `json:"text" example:"张三，本科，熟悉 Go、Vue、PostgreSQL，有微服务项目经验。"`
}

// ResumeMatchRequest 是简历岗位匹配请求。
type ResumeMatchRequest struct {
	ResumeID uint   `json:"resume_id" example:"1"`
	JobID    uint   `json:"job_id" example:"10"`
	Text     string `json:"text" example:"候选人简历文本，可选"`
}

// ResumeRiskCheckRequest 是简历风险检查请求。
type ResumeRiskCheckRequest struct {
	ResumeID uint   `json:"resume_id" example:"1"`
	Text     string `json:"text" example:"候选人简历文本"`
	JobID    uint   `json:"job_id" example:"10"`
}

// ResumeJobUpdateRequest 是简历关联职位请求。
type ResumeJobUpdateRequest struct {
	JobID uint `json:"job_id" example:"10"`
}

// StatusUpdateRequest 是通用状态更新请求。
type StatusUpdateRequest struct {
	Status string `json:"status" example:"screening"`
	Remark string `json:"remark" example:"进入初筛阶段"`
}

// ApplicationCreateRequest 是职位投递创建请求。
type ApplicationCreateRequest struct {
	JobID     uint   `json:"job_id" example:"10"`
	ResumeID  uint   `json:"resume_id" example:"1"`
	TalentID  uint   `json:"talent_id" example:"1"`
	CoverNote string `json:"cover_note" example:"我对该岗位非常感兴趣"`
}

// AIEvaluateRequest 是 AI 简历评估请求。
type AIEvaluateRequest struct {
	ResumeID uint `json:"resume_id" example:"1"`
	JobID    uint `json:"job_id" example:"10"`
}

// BatchEvaluateRequest 是 AI 批量评估请求。
type BatchEvaluateRequest struct {
	ResumeIDs []uint `json:"resume_ids" example:"1,2,3"`
	JobID     uint   `json:"job_id" example:"10"`
}

// OCRRequest 是 OCR 文本提取请求。
type OCRRequest struct {
	FileURL string `json:"file_url" example:"/uploads/resumes/resume.pdf"`
	Text    string `json:"text" example:"可选，直接传入待提取文本"`
}

// RecommendationRequest 是推荐相关请求。
type RecommendationRequest struct {
	TalentID uint   `json:"talent_id" example:"1"`
	JobID    uint   `json:"job_id" example:"10"`
	ResumeID uint   `json:"resume_id" example:"1"`
	Text     string `json:"text" example:"候选人或职位画像文本"`
	Limit    int    `json:"limit" example:"10"`
}

// RAGQueryRequest 是 RAG 查询请求。
type RAGQueryRequest struct {
	Query string `json:"query" example:"推荐适合 Go 后端工程师的岗位"`
	TopK  int    `json:"top_k" example:"5"`
	Type  string `json:"type" example:"job" enums:"job,talent,resume,all"`
}

// RAGIndexRequest 是 RAG 索引请求。
type RAGIndexRequest struct {
	ID      uint   `json:"id" example:"1"`
	Type    string `json:"type" example:"talent" enums:"talent,job,resume"`
	Content string `json:"content" example:"待索引的人才、职位或简历文本"`
}

// MessageCreateRequest 是站内信发送请求。
type MessageCreateRequest struct {
	ReceiverID uint   `json:"receiver_id" example:"2"`
	Title      string `json:"title" example:"面试通知"`
	Content    string `json:"content" example:"请明天下午参加一面"`
	Type       string `json:"type" example:"interview"`
}

// ConversationCreateRequest 是创建或获取会话请求。
type ConversationCreateRequest struct {
	TargetUserID uint `json:"target_user_id" example:"2"`
	JobID        uint `json:"job_id" example:"10"`
}

// ChatMessageRequest 是聊天消息发送请求。
type ChatMessageRequest struct {
	Content string `json:"content" example:"您好，我想了解一下这个岗位。"`
	Type    string `json:"type" example:"text"`
}

// NoticeUpsertRequest 是公告新增或更新请求。
type NoticeUpsertRequest struct {
	Title    string `json:"title" example:"校园招聘宣讲会通知"`
	Content  string `json:"content" example:"本周五下午 3 点举行线上宣讲会。"`
	Status   string `json:"status" example:"published" enums:"draft,published"`
	IsPinned bool   `json:"is_pinned" example:"false"`
	Priority string `json:"priority" example:"normal" enums:"normal,high,urgent"`
}

// InterviewUpsertRequest 是面试安排新增或更新请求。
type InterviewUpsertRequest struct {
	ApplicationID uint   `json:"application_id" example:"1"`
	TalentID      uint   `json:"talent_id" example:"1"`
	JobID         uint   `json:"job_id" example:"10"`
	InterviewerID uint   `json:"interviewer_id" example:"3"`
	Interviewer   string `json:"interviewer" example:"HR 张经理"`
	Type          string `json:"type" example:"technical"`
	ScheduledAt   string `json:"scheduled_at" example:"2026-05-21T15:00:00+08:00"`
	Location      string `json:"location" example:"腾讯会议"`
	Remark        string `json:"remark" example:"请提前准备项目介绍"`
}

// CompleteInterviewRequest 是面试完成请求。
type CompleteInterviewRequest struct {
	Result string `json:"result" example:"passed"`
	Remark string `json:"remark" example:"候选人基础扎实，建议进入下一轮"`
}

// InterviewFeedbackRequest 是面试反馈请求。
type InterviewFeedbackRequest struct {
	Score      int    `json:"score" example:"85"`
	Feedback   string `json:"feedback" example:"沟通能力良好，后端基础扎实"`
	Suggestion string `json:"suggestion" example:"建议进入二面"`
	Result     string `json:"result" example:"passed"`
}

// RescheduleInterviewRequest 是面试改期请求。
type RescheduleInterviewRequest struct {
	ScheduledAt string `json:"scheduled_at" example:"2026-05-22T10:00:00+08:00"`
	Reason      string `json:"reason" example:"候选人时间冲突，重新约定面试时间"`
}

// swaggerHealth godoc
// @Summary Gateway 健康检查
// @Description 返回 Gateway 状态和当前服务注册表。
// @Tags 网关
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /health [get]
func swaggerHealth() {}

// swaggerRegister godoc
// @Summary 用户注册
// @Description 创建求职者或 HR 用户账号。
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册参数：用户名、邮箱、密码、角色、姓名和手机号"
// @Success 201 {object} SwaggerResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/register [post]
func swaggerRegister() {}

// swaggerLogin godoc
// @Summary 用户登录
// @Description 使用用户名/邮箱和密码登录，返回 JWT Token。
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录参数：用户名/邮箱和密码"
// @Success 200 {object} SwaggerResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Router /api/v1/login [post]
func swaggerLogin() {}

// swaggerGetProfile godoc
// @Summary 获取当前用户资料
// @Description 根据 JWT 获取当前登录用户资料。
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Router /api/v1/profile [get]
func swaggerGetProfile() {}

// swaggerUpdateProfile godoc
// @Summary 更新当前用户资料
// @Description 更新当前登录用户的姓名、手机号、头像等资料。
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ProfileUpdateRequest true "用户资料：姓名、手机号、头像和邮箱"
// @Success 200 {object} SwaggerResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Router /api/v1/profile [put]
func swaggerUpdateProfile() {}

// swaggerListUsers godoc
// @Summary 用户列表
// @Description 查询系统用户列表，需要管理员或 HR 类角色。
// @Tags 用户
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} SwaggerResponse
// @Failure 403 {object} SwaggerErrorResponse
// @Router /api/v1/users [get]
func swaggerListUsers() {}

// swaggerListTalents godoc
// @Summary 获取人才列表
// @Description 按状态、关键词、技能、地点和经验筛选人才。
// @Tags 人才
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "人才状态"
// @Param search query string false "搜索关键词"
// @Param skills query []string false "技能"
// @Param location query string false "地点"
// @Param experience query int false "经验年限"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/talents [get]
func swaggerListTalents() {}

// swaggerCreateTalent godoc
// @Summary 创建人才
// @Description 创建人才库候选人记录。
// @Tags 人才
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TalentUpsertRequest true "人才信息：基础资料、技能、经验、来源和状态"
// @Success 201 {object} SwaggerResponse
// @Router /api/v1/talents [post]
func swaggerCreateTalent() {}

// swaggerTalentStats godoc
// @Summary 获取人才统计
// @Description 获取人才库总量、状态分布等统计信息。
// @Tags 人才
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/talents/stats [get]
func swaggerTalentStats() {}

// swaggerSearchTalents godoc
// @Summary 搜索人才
// @Description 按关键词、学历和地点搜索人才。
// @Tags 人才
// @Produce json
// @Security BearerAuth
// @Param keyword query string false "关键词"
// @Param education query string false "学历"
// @Param location query string false "地点"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/talents/search [get]
func swaggerSearchTalents() {}

// swaggerGetTalent godoc
// @Summary 获取人才详情
// @Description 根据人才 ID 获取详情。
// @Tags 人才
// @Produce json
// @Security BearerAuth
// @Param id path int true "人才 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/talents/{id} [get]
func swaggerGetTalent() {}

// swaggerUpdateTalent godoc
// @Summary 更新人才
// @Description 更新人才库候选人信息。
// @Tags 人才
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "人才 ID"
// @Param request body TalentUpsertRequest true "人才信息：基础资料、技能、经验、来源和状态"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/talents/{id} [put]
func swaggerUpdateTalent() {}

// swaggerDeleteTalent godoc
// @Summary 删除人才
// @Description 删除人才库候选人记录。
// @Tags 人才
// @Produce json
// @Security BearerAuth
// @Param id path int true "人才 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/talents/{id} [delete]
func swaggerDeleteTalent() {}

// swaggerListJobs godoc
// @Summary 获取职位列表
// @Description 公开查询职位列表，支持关键词、状态、地点、学历、级别等筛选。
// @Tags 职位
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "职位状态"
// @Param type query string false "职位类型"
// @Param location query string false "地点"
// @Param education query string false "学历"
// @Param search query string false "搜索关键词"
// @Param keyword query string false "关键词"
// @Param level query string false "级别"
// @Param experience query string false "经验要求"
// @Param sort_by query string false "排序字段"
// @Param sort_order query string false "排序方向"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/jobs [get]
func swaggerListJobs() {}

// swaggerCreateJob godoc
// @Summary 创建职位
// @Description 创建招聘职位，需要管理员或 HR 类角色。
// @Tags 职位
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body JobUpsertRequest true "职位信息：标题、部门、薪资、技能、描述和要求"
// @Success 201 {object} SwaggerResponse
// @Router /api/v1/jobs [post]
func swaggerCreateJob() {}

// swaggerJobStats godoc
// @Summary 获取职位统计
// @Description 获取职位数量、状态分布等统计信息。
// @Tags 职位
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/jobs/stats [get]
func swaggerJobStats() {}

// swaggerHotJobs godoc
// @Summary 获取热门职位
// @Description 获取热门职位列表。
// @Tags 职位
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/jobs/hot [get]
func swaggerHotJobs() {}

// swaggerGetJob godoc
// @Summary 获取职位详情
// @Description 根据职位 ID 获取职位详情。
// @Tags 职位
// @Produce json
// @Param id path int true "职位 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/jobs/{id} [get]
func swaggerGetJob() {}

// swaggerUpdateJob godoc
// @Summary 更新职位
// @Description 更新职位信息，需要管理员或 HR 类角色。
// @Tags 职位
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "职位 ID"
// @Param request body JobUpsertRequest true "职位信息：标题、部门、薪资、技能、描述和要求"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/jobs/{id} [put]
func swaggerUpdateJob() {}

// swaggerDeleteJob godoc
// @Summary 删除职位
// @Description 删除职位，需要管理员或 HR 类角色。
// @Tags 职位
// @Produce json
// @Security BearerAuth
// @Param id path int true "职位 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/jobs/{id} [delete]
func swaggerDeleteJob() {}

// swaggerJobApplications godoc
// @Summary 获取职位投递列表
// @Description 查询某个职位下的投递申请，需要管理员或 HR 类角色。
// @Tags 职位
// @Produce json
// @Security BearerAuth
// @Param id path int true "职位 ID"
// @Param status query string false "投递状态"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/jobs/{id}/applications [get]
func swaggerJobApplications() {}

// swaggerListResumes godoc
// @Summary 获取简历列表
// @Description 按人才、状态和关键词查询简历列表。
// @Tags 简历
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param talent_id query int false "人才 ID"
// @Param status query string false "简历状态"
// @Param search query string false "搜索关键词"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes [get]
func swaggerListResumes() {}

// swaggerCreateResume godoc
// @Summary 创建或上传简历记录
// @Description 创建简历记录，或通过 multipart 上传简历文件。
// @Tags 简历
// @Accept mpfd
// @Produce json
// @Param file formData file false "简历文件"
// @Param talent_id formData int false "人才 ID"
// @Param job_id formData int false "职位 ID"
// @Success 201 {object} SwaggerResponse
// @Router /api/v1/resumes [post]
func swaggerCreateResume() {}

// swaggerUploadResume godoc
// @Summary 登录用户上传简历文件
// @Description 求职者门户上传简历，后端根据 JWT 自动绑定人才档案。
// @Tags 简历
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param file formData file true "简历文件"
// @Param talent_id formData int false "人才 ID"
// @Param job_id formData int false "职位 ID"
// @Success 201 {object} SwaggerResponse
// @Router /api/v1/resumes/upload [post]
func swaggerUploadResume() {}

// swaggerResumesForEvaluation godoc
// @Summary 获取待评估简历列表
// @Description 获取自动评估系统使用的简历列表。
// @Tags 简历
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "简历状态"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/evaluation [get]
func swaggerResumesForEvaluation() {}

// swaggerResumeFile godoc
// @Summary 获取简历文件
// @Description 根据文件名返回简历文件流。
// @Tags 简历
// @Produce octet-stream
// @Param filename path string true "文件名"
// @Success 200 {file} file
// @Router /api/v1/resumes/file/{filename} [get]
func swaggerResumeFile() {}

// swaggerGetOnlineResume godoc
// @Summary 获取在线简历
// @Description 获取当前登录用户的在线简历。
// @Tags 简历
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/online [get]
func swaggerGetOnlineResume() {}

// swaggerSaveOnlineResume godoc
// @Summary 保存在线简历
// @Description 保存当前登录用户在线编辑的简历。
// @Tags 简历
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body OnlineResumeRequest true "在线简历：基础信息、教育经历、技能、项目经历和个人总结"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/online [put]
func swaggerSaveOnlineResume() {}

// swaggerParseResume godoc
// @Summary 解析简历文本
// @Description 根据传入的纯文本解析简历结构化信息。
// @Tags 简历
// @Accept json
// @Produce json
// @Param request body ResumeTextRequest true "简历文本：用于解析结构化字段"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/parse [post]
func swaggerParseResume() {}

// swaggerMatchResume godoc
// @Summary 简历职位匹配
// @Description 根据简历文本和职位要求计算匹配度。
// @Tags 简历
// @Accept json
// @Produce json
// @Param request body ResumeMatchRequest true "匹配参数：简历 ID、职位 ID 或待匹配文本"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/match [post]
func swaggerMatchResume() {}

// swaggerResumeRiskCheck godoc
// @Summary 简历风险综合检查
// @Description 检查简历内容中的潜在风险。
// @Tags 简历
// @Accept json
// @Produce json
// @Param request body ResumeRiskCheckRequest true "风险检查参数：简历 ID、职位 ID 或简历文本"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/risk-check [post]
func swaggerResumeRiskCheck() {}

// swaggerResumeTimeConflict godoc
// @Summary 工作经历时间冲突检查
// @Description 检查简历工作经历是否存在时间冲突。
// @Tags 简历
// @Accept json
// @Produce json
// @Param request body ResumeRiskCheckRequest true "时间冲突检查参数：简历 ID、职位 ID 或简历文本"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/risk-check/time-conflict [post]
func swaggerResumeTimeConflict() {}

// swaggerResumeEducationRisk godoc
// @Summary 教育经历风险检查
// @Description 检查简历教育经历真实性风险。
// @Tags 简历
// @Accept json
// @Produce json
// @Param request body ResumeRiskCheckRequest true "学历风险检查参数：简历 ID、职位 ID 或简历文本"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/risk-check/education [post]
func swaggerResumeEducationRisk() {}

// swaggerGetResume godoc
// @Summary 获取简历详情
// @Description 根据简历 ID 获取详情。
// @Tags 简历
// @Produce json
// @Param id path int true "简历 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/{id} [get]
func swaggerGetResume() {}

// swaggerDeleteResume godoc
// @Summary 删除简历
// @Description 删除指定简历。
// @Tags 简历
// @Produce json
// @Param id path int true "简历 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/{id} [delete]
func swaggerDeleteResume() {}

// swaggerDownloadResume godoc
// @Summary 下载简历文件
// @Description 根据简历 ID 下载文件。
// @Tags 简历
// @Produce octet-stream
// @Param id path int true "简历 ID"
// @Success 200 {file} file
// @Router /api/v1/resumes/{id}/download [get]
func swaggerDownloadResume() {}

// swaggerUpdateResumeJob godoc
// @Summary 为简历关联职位
// @Description 为已有简历补充或修改关联职位。
// @Tags 简历
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "简历 ID"
// @Param request body ResumeJobUpdateRequest true "职位关联参数：目标职位 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/{id}/job [put]
func swaggerUpdateResumeJob() {}

// swaggerUpdateResumeStatus godoc
// @Summary 更新简历状态
// @Description 更新简历处理状态。
// @Tags 简历
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "简历 ID"
// @Param request body StatusUpdateRequest true "状态参数：目标状态和备注"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/resumes/{id}/status [put]
func swaggerUpdateResumeStatus() {}

// swaggerListApplications godoc
// @Summary 获取投递列表
// @Description 查询职位投递列表，talent_id 可传 me 表示当前登录用户。
// @Tags 投递
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param job_id query int false "职位 ID"
// @Param talent_id query string false "人才 ID 或 me"
// @Param status query string false "投递状态"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/applications [get]
func swaggerListApplications() {}

// swaggerCreateApplication godoc
// @Summary 创建职位投递
// @Description 当前登录求职者投递职位。
// @Tags 投递
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ApplicationCreateRequest true "投递参数：职位、简历、人才和附言"
// @Success 201 {object} SwaggerResponse
// @Router /api/v1/applications [post]
func swaggerCreateApplication() {}

// swaggerUpdateApplication godoc
// @Summary 更新投递状态
// @Description 更新投递申请状态和备注。
// @Tags 投递
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "投递 ID"
// @Param request body StatusUpdateRequest true "投递状态参数：目标状态和备注"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/applications/{id} [put]
func swaggerUpdateApplication() {}

// swaggerDeleteApplication godoc
// @Summary 撤回或删除投递
// @Description 删除投递申请或求职者撤回申请。
// @Tags 投递
// @Produce json
// @Security BearerAuth
// @Param id path int true "投递 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/applications/{id} [delete]
func swaggerDeleteApplication() {}

// swaggerAIConfig godoc
// @Summary 检查 AI 配置
// @Description 返回 Coze、OCR、RAG 等 AI 能力配置状态。
// @Tags AI
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/ai/config [get]
func swaggerAIConfig() {}

// swaggerAICurrentTask godoc
// @Summary 获取当前 AI 任务
// @Description 获取当前正在处理的 AI 评估任务。
// @Tags AI
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/ai/current-task [get]
func swaggerAICurrentTask() {}

// swaggerAIEvaluate godoc
// @Summary 根据简历 ID 发起 AI 评估
// @Description 基于简历 ID 和 JD 文本或职位 ID 发起 AI 评估。
// @Tags AI
// @Accept json
// @Produce json
// @Param request body AIEvaluateRequest true "评估参数：简历 ID 和目标职位 ID"
// @Success 200 {object} SwaggerResponse
// @Failure 503 {object} SwaggerErrorResponse
// @Router /api/v1/ai/evaluate [post]
func swaggerAIEvaluate() {}

// swaggerAIEvaluateUpload godoc
// @Summary 上传 PDF 并发起 AI 评估
// @Description 上传 PDF 简历文件并立即进行 AI 评估。
// @Tags AI
// @Accept mpfd
// @Produce json
// @Param file formData file true "PDF 简历"
// @Param jd_text formData string false "JD 文本"
// @Param candidate_name formData string false "候选人姓名"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/ai/evaluate/upload [post]
func swaggerAIEvaluateUpload() {}

// swaggerAIEvaluateBatch godoc
// @Summary 批量发起 AI 评估
// @Description 按简历 ID 列表批量发起 AI 评估。
// @Tags AI
// @Accept json
// @Produce json
// @Param request body BatchEvaluateRequest true "批量评估参数：简历 ID 列表和目标职位 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/ai/evaluate/batch [post]
func swaggerAIEvaluateBatch() {}

// swaggerAIEvaluateResult godoc
// @Summary 获取 AI 评估结果
// @Description 根据评估 ID 获取 AI 评估结果。
// @Tags AI
// @Produce json
// @Param id path int true "评估 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/ai/evaluate/{id}/result [get]
func swaggerAIEvaluateResult() {}

// swaggerAIParse godoc
// @Summary AI 智能解析简历
// @Description 调用 AI 能力解析简历。
// @Tags AI
// @Accept json
// @Produce json
// @Param request body ResumeTextRequest true "AI 解析参数：简历文本"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/ai/parse [post]
func swaggerAIParse() {}

// swaggerAIOCR godoc
// @Summary OCR 文本提取
// @Description 从简历文件中提取文本。
// @Tags AI
// @Accept json
// @Produce json
// @Param request body OCRRequest true "OCR 参数：文件地址或待提取文本"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/ai/ocr [post]
func swaggerAIOCR() {}

// swaggerListEvaluations godoc
// @Summary 获取评估记录列表
// @Description 查询 AI 评估记录。
// @Tags 评估
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "状态"
// @Param eval_type query string false "评估类型"
// @Param search query string false "关键词"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/evaluations [get]
func swaggerListEvaluations() {}

// swaggerEvaluationStats godoc
// @Summary 获取评估统计
// @Description 获取 AI 评估记录统计。
// @Tags 评估
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/evaluations/stats [get]
func swaggerEvaluationStats() {}

// swaggerGetEvaluation godoc
// @Summary 获取评估详情
// @Description 根据 ID 获取评估记录详情。
// @Tags 评估
// @Produce json
// @Param id path int true "评估 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/evaluations/{id} [get]
func swaggerGetEvaluation() {}

// swaggerDeleteEvaluation godoc
// @Summary 删除评估记录
// @Description 删除指定 AI 评估记录。
// @Tags 评估
// @Produce json
// @Param id path int true "评估 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/evaluations/{id} [delete]
func swaggerDeleteEvaluation() {}

// swaggerEvaluationProcess godoc
// @Summary 获取评估流程详情
// @Description 获取 AI 评估流程、OCR、Embedding、RAG、LLM 步骤信息。
// @Tags 评估
// @Produce json
// @Param id path int true "评估 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/evaluations/{id}/process [get]
func swaggerEvaluationProcess() {}

// swaggerRecommendJobsForTalent godoc
// @Summary 为人才推荐职位
// @Description 根据人才画像推荐匹配职位。
// @Tags 推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RecommendationRequest true "人才画像：人才 ID、简历 ID、画像文本和返回数量"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/jobs-for-talent [post]
func swaggerRecommendJobsForTalent() {}

// swaggerRecommendTalentsForJob godoc
// @Summary 为职位推荐人才
// @Description 根据职位画像推荐匹配人才。
// @Tags 推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RecommendationRequest true "职位画像：职位 ID、岗位文本和返回数量"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/talents-for-job [post]
func swaggerRecommendTalentsForJob() {}

// swaggerRecommendationStats godoc
// @Summary 获取推荐统计
// @Description 获取推荐和匹配相关统计。
// @Tags 推荐
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/stats [get]
func swaggerRecommendationStats() {}

// swaggerRecommendationBatch godoc
// @Summary 批量推荐
// @Description 批量计算人才与职位推荐。
// @Tags 推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RecommendationRequest true "批量推荐参数：人才、职位或简历范围"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/batch [post]
func swaggerRecommendationBatch() {}

// swaggerAttributionReport godoc
// @Summary 生成推荐归因报告
// @Description 生成指定人才与职位的匹配归因报告。
// @Tags 推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RecommendationRequest true "归因参数：人才 ID、职位 ID 和匹配文本"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/attribution-report [post]
func swaggerAttributionReport() {}

// swaggerSemanticMatch godoc
// @Summary 语义相似度匹配
// @Description 计算两段文本的语义相似度。
// @Tags 推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RecommendationRequest true "语义匹配参数：人才/职位文本和返回数量"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/semantic-match [post]
func swaggerSemanticMatch() {}

// swaggerRAGQuery godoc
// @Summary RAG 查询
// @Description 查询向量索引中的人才、职位或简历内容。
// @Tags 推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RAGQueryRequest true "RAG 查询参数：问题、TopK 和检索类型"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/rag/query [post]
func swaggerRAGQuery() {}

// swaggerRAGIndexTalent godoc
// @Summary 索引单个人才
// @Description 将指定人才写入 RAG 索引。
// @Tags 推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RAGIndexRequest true "人才索引参数：人才 ID 和文本内容"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/rag/index-talent [post]
func swaggerRAGIndexTalent() {}

// swaggerRAGIndexJob godoc
// @Summary 索引单个职位
// @Description 将指定职位写入 RAG 索引。
// @Tags 推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RAGIndexRequest true "职位索引参数：职位 ID 和文本内容"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/rag/index-job [post]
func swaggerRAGIndexJob() {}

// swaggerRAGIndexResume godoc
// @Summary 索引单份简历
// @Description 将指定简历写入 RAG 索引。
// @Tags 推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RAGIndexRequest true "简历索引参数：简历 ID 和文本内容"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/rag/index-resume [post]
func swaggerRAGIndexResume() {}

// swaggerRAGIndexAll godoc
// @Summary 全量索引人才与职位
// @Description 重建人才与职位 RAG 索引。
// @Tags 推荐
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/rag/index-all [post]
func swaggerRAGIndexAll() {}

// swaggerRAGMatch godoc
// @Summary RAG 人岗匹配
// @Description 使用 RAG 能力对指定人才和职位进行匹配。
// @Tags 推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RecommendationRequest true "RAG 匹配参数：人才、职位、简历或文本"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/recommendations/rag/match [post]
func swaggerRAGMatch() {}

// swaggerListMessages godoc
// @Summary 获取系统消息列表
// @Description 获取当前用户的系统消息。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param type query string false "消息类型"
// @Param is_read query bool false "是否已读"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/messages [get]
func swaggerListMessages() {}

// swaggerSendMessage godoc
// @Summary 发送系统消息
// @Description 发送站内系统消息。
// @Tags 消息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body MessageCreateRequest true "消息内容：接收人、标题、内容和消息类型"
// @Success 201 {object} SwaggerResponse
// @Router /api/v1/messages [post]
func swaggerSendMessage() {}

// swaggerMessageStats godoc
// @Summary 获取消息统计
// @Description 获取当前用户消息统计。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/messages/stats [get]
func swaggerMessageStats() {}

// swaggerUnreadMessages godoc
// @Summary 获取未读消息数
// @Description 获取当前用户未读消息数量。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/messages/unread-count [get]
func swaggerUnreadMessages() {}

// swaggerMarkMessageRead godoc
// @Summary 标记消息已读
// @Description 将指定消息标记为已读。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Param id path int true "消息 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/messages/{id}/read [put]
func swaggerMarkMessageRead() {}

// swaggerDeleteMessage godoc
// @Summary 删除消息
// @Description 删除指定系统消息。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Param id path int true "消息 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/messages/{id} [delete]
func swaggerDeleteMessage() {}

// swaggerConversations godoc
// @Summary 获取会话列表
// @Description 获取当前用户聊天会话列表。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/conversations [get]
func swaggerConversations() {}

// swaggerCreateConversation godoc
// @Summary 创建或获取会话
// @Description 与指定用户创建或获取已有会话。
// @Tags 消息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ConversationCreateRequest true "会话参数：目标用户和关联职位"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/conversations [post]
func swaggerCreateConversation() {}

// swaggerConversationUnread godoc
// @Summary 获取会话未读总数
// @Description 获取当前用户所有会话未读消息数。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/conversations/unread-count [get]
func swaggerConversationUnread() {}

// swaggerConversationMessages godoc
// @Summary 获取会话消息
// @Description 获取指定会话中的聊天消息。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Param id path int true "会话 ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/conversations/{id}/messages [get]
func swaggerConversationMessages() {}

// swaggerSendConversationMessage godoc
// @Summary 发送会话消息
// @Description 向指定会话发送聊天消息。
// @Tags 消息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会话 ID"
// @Param request body ChatMessageRequest true "聊天消息：文本内容和消息类型"
// @Success 201 {object} SwaggerResponse
// @Router /api/v1/conversations/{id}/messages [post]
func swaggerSendConversationMessage() {}

// swaggerMarkConversationRead godoc
// @Summary 标记会话已读
// @Description 将指定会话中的消息标记为已读。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Param id path int true "会话 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/conversations/{id}/read [put]
func swaggerMarkConversationRead() {}

// swaggerOnlineStatus godoc
// @Summary 获取在线用户统计
// @Description 获取在线用户总数和列表。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/online-status [get]
func swaggerOnlineStatus() {}

// swaggerUserOnlineStatus godoc
// @Summary 获取指定用户在线状态
// @Description 查询指定用户是否在线。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "用户 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/online-status/{user_id} [get]
func swaggerUserOnlineStatus() {}

// swaggerWebSocket godoc
// @Summary WebSocket 实时消息连接
// @Description 建立实时消息 WebSocket 连接。
// @Tags 消息
// @Security BearerAuth
// @Success 101 {string} string "Switching Protocols"
// @Router /api/v1/ws [get]
func swaggerWebSocket() {}

// swaggerListNotices godoc
// @Summary 获取公告列表
// @Description 查询系统公告。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param keyword query string false "关键词"
// @Param status query string false "状态"
// @Param priority query string false "优先级"
// @Param is_pinned query bool false "是否置顶"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/notices [get]
func swaggerListNotices() {}

// swaggerCreateNotice godoc
// @Summary 创建公告
// @Description 创建系统公告。
// @Tags 消息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body NoticeUpsertRequest true "公告内容：标题、正文、状态、置顶和优先级"
// @Success 201 {object} SwaggerResponse
// @Router /api/v1/notices [post]
func swaggerCreateNotice() {}

// swaggerGetNotice godoc
// @Summary 获取公告详情
// @Description 根据 ID 获取公告详情。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Param id path int true "公告 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/notices/{id} [get]
func swaggerGetNotice() {}

// swaggerUpdateNotice godoc
// @Summary 更新公告
// @Description 更新公告内容。
// @Tags 消息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "公告 ID"
// @Param request body NoticeUpsertRequest true "公告内容：标题、正文、状态、置顶和优先级"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/notices/{id} [put]
func swaggerUpdateNotice() {}

// swaggerDeleteNotice godoc
// @Summary 删除公告
// @Description 删除指定公告。
// @Tags 消息
// @Produce json
// @Security BearerAuth
// @Param id path int true "公告 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/notices/{id} [delete]
func swaggerDeleteNotice() {}

// swaggerListInterviews godoc
// @Summary 获取面试列表
// @Description 查询面试安排列表。
// @Tags 面试
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "状态"
// @Param date query string false "日期"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param interviewer_id query int false "面试官 ID"
// @Param candidate_id query int false "候选人 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews [get]
func swaggerListInterviews() {}

// swaggerCreateInterview godoc
// @Summary 创建面试安排
// @Description 创建候选人面试安排。
// @Tags 面试
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body InterviewUpsertRequest true "面试安排：投递、人才、职位、面试官、时间和地点"
// @Success 201 {object} SwaggerResponse
// @Router /api/v1/interviews [post]
func swaggerCreateInterview() {}

// swaggerInterviewStats godoc
// @Summary 获取面试统计
// @Description 获取面试数量和状态统计。
// @Tags 面试
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/stats [get]
func swaggerInterviewStats() {}

// swaggerTodayInterviews godoc
// @Summary 获取今日面试
// @Description 获取当天的面试安排。
// @Tags 面试
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/today [get]
func swaggerTodayInterviews() {}

// swaggerInterviewerSchedule godoc
// @Summary 获取面试官日程
// @Description 查询指定面试官的面试日程。
// @Tags 面试
// @Produce json
// @Security BearerAuth
// @Param interviewer_id path int true "面试官 ID"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/interviewer/{interviewer_id} [get]
func swaggerInterviewerSchedule() {}

// swaggerCandidateInterviews godoc
// @Summary 获取候选人面试列表
// @Description 查询指定候选人的面试安排。
// @Tags 面试
// @Produce json
// @Security BearerAuth
// @Param candidate_id path int true "候选人 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/candidate/{candidate_id} [get]
func swaggerCandidateInterviews() {}

// swaggerGetInterview godoc
// @Summary 获取面试详情
// @Description 根据 ID 获取面试详情。
// @Tags 面试
// @Produce json
// @Security BearerAuth
// @Param id path int true "面试 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/{id} [get]
func swaggerGetInterview() {}

// swaggerUpdateInterview godoc
// @Summary 更新面试
// @Description 更新面试安排。
// @Tags 面试
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "面试 ID"
// @Param request body InterviewUpsertRequest true "面试信息：面试官、时间、地点和备注"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/{id} [put]
func swaggerUpdateInterview() {}

// swaggerDeleteInterview godoc
// @Summary 删除面试
// @Description 删除指定面试安排。
// @Tags 面试
// @Produce json
// @Security BearerAuth
// @Param id path int true "面试 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/{id} [delete]
func swaggerDeleteInterview() {}

// swaggerCancelInterview godoc
// @Summary 取消面试
// @Description 取消指定面试。
// @Tags 面试
// @Produce json
// @Security BearerAuth
// @Param id path int true "面试 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/{id}/cancel [post]
func swaggerCancelInterview() {}

// swaggerCompleteInterview godoc
// @Summary 完成面试
// @Description 标记面试完成，并可提交简要反馈与评分。
// @Tags 面试
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "面试 ID"
// @Param request body CompleteInterviewRequest false "完成参数：结果和备注"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/{id}/complete [post]
func swaggerCompleteInterview() {}

// swaggerGetInterviewFeedback godoc
// @Summary 获取面试反馈
// @Description 获取指定面试的反馈详情。
// @Tags 面试
// @Produce json
// @Security BearerAuth
// @Param id path int true "面试 ID"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/{id}/feedback [get]
func swaggerGetInterviewFeedback() {}

// swaggerSubmitInterviewFeedback godoc
// @Summary 提交面试反馈
// @Description 为指定面试提交评分、优劣势和录用建议。
// @Tags 面试
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "面试 ID"
// @Param request body InterviewFeedbackRequest true "面试反馈：评分、反馈、建议和结果"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/{id}/feedback [post]
func swaggerSubmitInterviewFeedback() {}

// swaggerRescheduleInterview godoc
// @Summary 面试改期
// @Description 修改面试日期和时间。
// @Tags 面试
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "面试 ID"
// @Param request body RescheduleInterviewRequest true "改期参数：新时间和改期原因"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/interviews/{id}/reschedule [post]
func swaggerRescheduleInterview() {}

// swaggerQueryLogs godoc
// @Summary 查询操作日志
// @Description 按服务、用户、路径、动作、时间范围等条件查询操作日志。
// @Tags 日志
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param service query string false "服务名"
// @Param username query string false "用户名"
// @Param method query string false "HTTP 方法"
// @Param path query string false "路径"
// @Param action query string false "动作"
// @Param module query string false "模块"
// @Param level query string false "级别"
// @Param keyword query string false "关键词"
// @Param user_id query int false "用户 ID"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/logs [get]
func swaggerQueryLogs() {}

// swaggerLogStats godoc
// @Summary 获取日志统计
// @Description 获取操作日志聚合统计。
// @Tags 日志
// @Produce json
// @Security BearerAuth
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/logs/stats [get]
func swaggerLogStats() {}

// swaggerLogServices godoc
// @Summary 获取日志服务列表
// @Description 获取日志中出现过的服务名列表。
// @Tags 日志
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/logs/services [get]
func swaggerLogServices() {}

// swaggerLogActions godoc
// @Summary 获取日志操作类型
// @Description 获取日志中出现过的操作类型列表。
// @Tags 日志
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/logs/actions [get]
func swaggerLogActions() {}

// swaggerLogCleanup godoc
// @Summary 清理旧日志
// @Description 按保留天数清理旧日志。
// @Tags 日志
// @Produce json
// @Security BearerAuth
// @Param days query int false "保留天数"
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/logs/cleanup [delete]
func swaggerLogCleanup() {}

// swaggerStatsDashboard godoc
// @Summary 获取仪表盘统计
// @Description 获取仪表盘核心统计指标。
// @Tags 统计
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/stats/dashboard [get]
func swaggerStatsDashboard() {}

// swaggerStatsFunnel godoc
// @Summary 获取招聘漏斗统计
// @Description 获取招聘流程各阶段漏斗数据。
// @Tags 统计
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/stats/funnel [get]
func swaggerStatsFunnel() {}

// swaggerStatsChannels godoc
// @Summary 获取渠道统计
// @Description 获取招聘渠道数据统计。
// @Tags 统计
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/stats/channels [get]
func swaggerStatsChannels() {}

// swaggerStatsDepartmentProgress godoc
// @Summary 获取部门招聘进度
// @Description 获取各部门招聘进度统计。
// @Tags 统计
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/stats/department-progress [get]
func swaggerStatsDepartmentProgress() {}

// swaggerStatsInterviewerRank godoc
// @Summary 获取面试官排名
// @Description 获取面试官维度的统计排名。
// @Tags 统计
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/stats/interviewer-rank [get]
func swaggerStatsInterviewerRank() {}

// swaggerStatsTrend godoc
// @Summary 获取招聘趋势
// @Description 获取招聘数据趋势图。
// @Tags 统计
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/stats/trend [get]
func swaggerStatsTrend() {}

// swaggerStatsJobRank godoc
// @Summary 获取职位排名
// @Description 获取职位维度的招聘排名。
// @Tags 统计
// @Produce json
// @Success 200 {object} SwaggerResponse
// @Router /api/v1/stats/job-rank [get]
func swaggerStatsJobRank() {}

// swaggerInternalMessage godoc
// @Summary 内部发送消息
// @Description 服务间调用的消息发送接口，可使用 X-Internal-Token。
// @Tags 内部接口
// @Accept json
// @Produce json
// @Security InternalToken
// @Param request body MessageCreateRequest true "内部消息内容：接收人、标题、内容和消息类型"
// @Success 201 {object} SwaggerResponse
// @Router /internal/messages [post]
func swaggerInternalMessage() {}

// swaggerInternalRAGQuery godoc
// @Summary 内部 RAG 查询
// @Description 服务间调用的 RAG 查询接口。
// @Tags 内部接口
// @Accept json
// @Produce json
// @Security InternalToken
// @Param request body RAGQueryRequest true "内部 RAG 查询参数：问题、TopK 和检索类型"
// @Success 200 {object} SwaggerResponse
// @Router /internal/recommendations/rag/query [post]
func swaggerInternalRAGQuery() {}

// swaggerInternalRAGIndexTalent godoc
// @Summary 内部索引人才
// @Description 服务间调用的人才索引接口。
// @Tags 内部接口
// @Accept json
// @Produce json
// @Security InternalToken
// @Param request body RAGIndexRequest true "内部人才索引参数：人才 ID 和文本内容"
// @Success 200 {object} SwaggerResponse
// @Router /internal/recommendations/rag/index-talent [post]
func swaggerInternalRAGIndexTalent() {}

// swaggerInternalRAGIndexJob godoc
// @Summary 内部索引职位
// @Description 服务间调用的职位索引接口。
// @Tags 内部接口
// @Accept json
// @Produce json
// @Security InternalToken
// @Param request body RAGIndexRequest true "内部职位索引参数：职位 ID 和文本内容"
// @Success 200 {object} SwaggerResponse
// @Router /internal/recommendations/rag/index-job [post]
func swaggerInternalRAGIndexJob() {}

// swaggerInternalRAGIndexResume godoc
// @Summary 内部索引简历
// @Description 服务间调用的简历索引接口。
// @Tags 内部接口
// @Accept json
// @Produce json
// @Security InternalToken
// @Param request body RAGIndexRequest true "内部简历索引参数：简历 ID 和文本内容"
// @Success 200 {object} SwaggerResponse
// @Router /internal/recommendations/rag/index-resume [post]
func swaggerInternalRAGIndexResume() {}

// swaggerInternalRAGIndexAll godoc
// @Summary 内部全量索引
// @Description 服务间调用的全量索引接口。
// @Tags 内部接口
// @Produce json
// @Security InternalToken
// @Success 200 {object} SwaggerResponse
// @Router /internal/recommendations/rag/index-all [post]
func swaggerInternalRAGIndexAll() {}

// swaggerInternalRAGMatch godoc
// @Summary 内部 RAG 人岗匹配
// @Description 服务间调用的人岗匹配接口。
// @Tags 内部接口
// @Accept json
// @Produce json
// @Security InternalToken
// @Param request body RecommendationRequest true "内部 RAG 匹配参数：人才、职位、简历或文本"
// @Success 200 {object} SwaggerResponse
// @Router /internal/recommendations/rag/match [post]
func swaggerInternalRAGMatch() {}
