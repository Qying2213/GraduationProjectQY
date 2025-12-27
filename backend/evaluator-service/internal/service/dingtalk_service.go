package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"evaluator-service/internal/dingtalk"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/models"
	"evaluator-service/internal/repository"

	"github.com/robfig/cron/v3"
)

// PushSession 推送会话，记录每次推送的候选人列表
type PushSession struct {
	SessionID    string    // 会话ID（时间戳）
	CandidateIDs []uint    // 候选人ID列表（按推送顺序）
	PushTime     time.Time // 推送时间
}

type DingTalkService struct {
	repo             *repository.CandidateRepository
	dtRepo           *repository.DingTalkRepository
	log              *logging.Logger
	client           *dingtalk.Client
	cron             *cron.Cron
	cancelFunc       context.CancelFunc
	sessionCache     map[string]*PushSession // 会话缓存：sessionID -> PushSession
	currentSessionID string                  // 当前最新的会话ID
	sessionMutex     sync.RWMutex            // 保护会话缓存的并发访问
}

func NewDingTalkService(
	repo *repository.CandidateRepository,
	dtRepo *repository.DingTalkRepository,
	log *logging.Logger,
) *DingTalkService {
	service := &DingTalkService{
		repo:         repo,
		dtRepo:       dtRepo,
		log:          log,
		cron:         cron.New(),
		sessionCache: make(map[string]*PushSession),
	}

	// 启动定期清理过期会话的goroutine
	go service.cleanExpiredSessions()

	return service
}

// Start 启动钉钉服务（定时任务 + Stream监听）
func (s *DingTalkService) Start(ctx context.Context) error {
	config, err := s.dtRepo.Get()
	if err != nil {
		return fmt.Errorf("get dingtalk config: %w", err)
	}

	if config == nil || !config.Enabled {
		s.log.Info("dingtalk service disabled")
		return nil
	}

	// 创建钉钉客户端
	s.client = dingtalk.NewClient(config, s.log)

	// 启动定时任务
	if err := s.startCronJob(config); err != nil {
		return fmt.Errorf("start cron job: %w", err)
	}

	// 启动Stream监听
	ctx, cancel := context.WithCancel(ctx)
	s.cancelFunc = cancel

	if err := s.client.StartStream(ctx, s.handleMessage); err != nil {
		return fmt.Errorf("start stream: %w", err)
	}

	s.log.Info("dingtalk service started")
	return nil
}

// Stop 停止钉钉服务
func (s *DingTalkService) Stop() {
	if s.cron != nil {
		s.cron.Stop()
	}
	if s.client != nil {
		s.client.Close()
	}
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	s.log.Info("dingtalk service stopped")
}

// Restart 重启钉钉服务
func (s *DingTalkService) Restart(ctx context.Context) error {
	s.log.Info("restarting dingtalk service")
	s.Stop()

	// 创建新的context，避免使用被取消的context
	newCtx := context.Background()
	if err := s.Start(newCtx); err != nil {
		s.log.Error("failed to restart dingtalk service", logging.Err(err))
		return err
	}

	s.log.Info("dingtalk service restarted successfully")
	return nil
}

// startCronJob 启动定时任务
func (s *DingTalkService) startCronJob(config *models.DingTalkConfig) error {
	// 解析推送时间 (格式: HH:MM)
	parts := strings.Split(config.PushTime, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid push_time format: %s", config.PushTime)
	}

	hour := parts[0]
	minute := parts[1]

	// Cron表达式: 分 时 日 月 周
	cronExpr := fmt.Sprintf("%s %s * * *", minute, hour)

	_, err := s.cron.AddFunc(cronExpr, func() {
		s.log.Info("dingtalk cron job triggered")
		if err := s.pushDailyCandidates(); err != nil {
			s.log.Error("push daily candidates failed", logging.Err(err))
		}
	})

	if err != nil {
		return fmt.Errorf("add cron job: %w", err)
	}

	s.cron.Start()
	s.log.Info("dingtalk cron job started", logging.KV("schedule", cronExpr))
	return nil
}

// pushDailyCandidates 推送每日候选人列表（定时任务调用）
func (s *DingTalkService) pushDailyCandidates() error {
	config, err := s.dtRepo.Get()
	if err != nil {
		return err
	}

	if config == nil || !config.Enabled {
		s.log.Info("dingtalk config not enabled, skip daily push")
		return nil
	}

	// 查询通知次数为0的候选人，按分数降序
	candidates, err := s.repo.FindUnnotified(config.PushLimit)
	if err != nil {
		return fmt.Errorf("find unnotified candidates: %w", err)
	}

	if len(candidates) == 0 {
		s.log.Info("no unnotified candidates to push")
		return nil
	}

	// 复用统一的推送方法
	return s.PushCandidates(candidates, config)
}

// buildCandidateListMarkdown 构建候选人列表Markdown
func (s *DingTalkService) buildCandidateListMarkdown(candidates []models.Candidate, atUserIds []string) string {
	var sb strings.Builder

	sb.WriteString("## 📋 待面试候选人列表\n\n")
	sb.WriteString(fmt.Sprintf("> 共 **%d** 位候选人待处理\n\n", len(candidates)))

	for i, c := range candidates {
		gradeEmoji := s.getGradeEmoji(c.Grade)
		sb.WriteString(fmt.Sprintf("**%d. %s** %s\n", i+1, c.Name, gradeEmoji))
		sb.WriteString(fmt.Sprintf("- 评分: **%.1f** 分 | 评级: **%s**\n", c.TotalScore, c.Grade))
		sb.WriteString(fmt.Sprintf("- JD匹配: %d%% | 状态: %s\n", c.JDMatch, c.Status))
		sb.WriteString(fmt.Sprintf("- 建议: %s\n", c.Recommendation))
		sb.WriteString("\n")
	}

	sb.WriteString("---\n")
	sb.WriteString("💡 回复序号查看详情，如：`1` 或 `1,2,3`\n\n")

	// 在消息末尾添加@用户（钉钉Markdown消息需要这样才能@生效）
	if len(atUserIds) > 0 {
		sb.WriteString("\n")
		for _, userId := range atUserIds {
			if userId != "" {
				sb.WriteString(fmt.Sprintf("@%s ", userId))
			}
		}
	}

	return sb.String()
}

func (s *DingTalkService) getGradeEmoji(grade string) string {
	switch grade {
	case "A":
		return "🌟"
	case "B":
		return "✨"
	case "C":
		return "⭐"
	case "D":
		return "💫"
	default:
		return "📄"
	}
}

// handleMessage 处理钉钉消息回调
func (s *DingTalkService) handleMessage(ctx context.Context, content string, senderID string) error {
	content = strings.TrimSpace(content)

	s.log.Info("processing message",
		logging.KV("content", content),
		logging.KV("sender", senderID))

	// 解析候选人序号（支持多个序号）
	nums := s.parseNumbers(content)
	if len(nums) == 0 {
		// 不包含有效数字，忽略
		s.log.Info("message contains no valid numbers, ignoring", logging.KV("content", content))
		return nil
	}

	// 从最新的推送会话中获取候选人列表
	session := s.getCurrentSession()
	if session == nil {
		s.log.Warn("no active push session found")
		return fmt.Errorf("暂无候选人列表，请先推送候选人")
	}

	s.log.Info("using push session",
		logging.KV("session_id", session.SessionID),
		logging.KV("candidate_count", len(session.CandidateIDs)),
		logging.KV("push_time", session.PushTime.Format("2006-01-02 15:04:05")),
		logging.KV("requested_nums", nums))

	// 验证所有序号是否在范围内
	var invalidNums []int
	for _, num := range nums {
		if num < 1 || num > len(session.CandidateIDs) {
			invalidNums = append(invalidNums, num)
		}
	}

	if len(invalidNums) > 0 {
		return fmt.Errorf("序号 %v 超出范围，请输入 1-%d", invalidNums, len(session.CandidateIDs))
	}

	// 查询所有候选人详情
	var candidates []*models.Candidate
	for _, num := range nums {
		candidateID := session.CandidateIDs[num-1]

		candidate, err := s.repo.FindByID(candidateID)
		if err != nil {
			s.log.Error("find candidate by id failed",
				logging.KV("id", candidateID),
				logging.KV("num", num),
				logging.Err(err))
			continue
		}

		if candidate == nil {
			s.log.Error("candidate not found",
				logging.KV("id", candidateID),
				logging.KV("num", num))
			continue
		}

		candidates = append(candidates, candidate)
	}

	if len(candidates) == 0 {
		return fmt.Errorf("未找到候选人信息")
	}

	s.log.Info("found candidates",
		logging.KV("count", len(candidates)))

	// 构建并发送详细信息
	if len(candidates) == 1 {
		// 单个候选人：发送详细信息
		detail := s.buildCandidateDetail(candidates[0])
		if err := s.client.SendMarkdownMessage(
			fmt.Sprintf("候选人详情 - %s", candidates[0].Name),
			detail,
			[]string{senderID},
			false,
		); err != nil {
			return err
		}
	} else {
		// 多个候选人：发送汇总信息
		summary := s.buildCandidatesSummary(candidates)
		if err := s.client.SendMarkdownMessage(
			fmt.Sprintf("候选人汇总 - 共%d人", len(candidates)),
			summary,
			[]string{senderID},
			false,
		); err != nil {
			return err
		}
	}

	s.log.Info("sent candidate details",
		logging.KV("count", len(candidates)),
		logging.KV("sender", senderID))

	return nil
}

// parseNumbers 解析消息中的数字序号（支持多种分隔符）
func (s *DingTalkService) parseNumbers(content string) []int {
	// 支持的分隔符：逗号(中英文)、顿号、空格、分号、斜杠
	// 替换所有分隔符为空格
	replacer := strings.NewReplacer(
		",", " ", // 英文逗号
		"，", " ", // 中文逗号
		"、", " ", // 顿号
		"；", " ", // 中文分号
		";", " ", // 英文分号
		"/", " ", // 斜杠
		"|", " ", // 竖线
		"\t", " ", // 制表符
	)
	normalized := replacer.Replace(content)

	// 按空格分割
	parts := strings.Fields(normalized)

	// 解析每个部分为数字
	var nums []int
	seen := make(map[int]bool) // 去重

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		// 移除 num > 0 的限制，允许解析0和负数，让后续验证逻辑处理
		if err == nil && !seen[num] {
			nums = append(nums, num)
			seen[num] = true
		}
	}

	return nums
}

// buildCandidatesSummary 构建多个候选人的汇总信息
func (s *DingTalkService) buildCandidatesSummary(candidates []*models.Candidate) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## 📋 候选人汇总（共%d人）\n\n", len(candidates)))

	for i, c := range candidates {
		gradeEmoji := s.getGradeEmoji(c.Grade)
		sb.WriteString(fmt.Sprintf("### %d. %s %s\n\n", i+1, c.Name, gradeEmoji))

		// 基本信息（简化版）
		sb.WriteString(fmt.Sprintf("- **总分**: %.1f 分 | **评级**: %s\n", c.TotalScore, c.Grade))
		sb.WriteString(fmt.Sprintf("- **JD匹配**: %d%% | **状态**: %s\n", c.JDMatch, c.Status))
		sb.WriteString(fmt.Sprintf("- **建议**: %s\n", c.Recommendation))

		// 各维度得分（紧凑版）
		sb.WriteString(fmt.Sprintf("- **得分明细**: 年龄%d | 经验%d | 学历%d | 公司%d | 技术%d | 项目%d\n",
			c.AgeScore, c.ExperienceScore, c.EducationScore,
			c.CompanyScore, c.TechScore, c.ProjectScore))

		if i < len(candidates)-1 {
			sb.WriteString("\n---\n\n")
		}
	}

	sb.WriteString("\n\n💡 **提示**: 回复单个序号可查看详细信息")

	return sb.String()
}

// buildCandidateDetail 构建候选人详细信息
func (s *DingTalkService) buildCandidateDetail(c *models.Candidate) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## 👤 %s\n\n", c.Name))

	// 基本信息
	sb.WriteString("### 📊 综合评分\n")
	sb.WriteString(fmt.Sprintf("- **总分**: %.1f 分\n", c.TotalScore))
	sb.WriteString(fmt.Sprintf("- **评级**: %s %s\n", c.Grade, s.getGradeEmoji(c.Grade)))
	sb.WriteString(fmt.Sprintf("- **JD匹配度**: %d%%\n", c.JDMatch))
	sb.WriteString(fmt.Sprintf("- **录用建议**: %s\n\n", c.Recommendation))

	// 各维度得分
	sb.WriteString("### 📈 各维度得分\n")
	sb.WriteString(fmt.Sprintf("- **年龄**: %d/10 - %s\n", c.AgeScore, c.AgeReason))
	sb.WriteString(fmt.Sprintf("- **工作经验**: %d/25 - %s\n", c.ExperienceScore, c.ExperienceReason))
	sb.WriteString(fmt.Sprintf("- **学历背景**: %d/20 - %s\n", c.EducationScore, c.EducationReason))
	sb.WriteString(fmt.Sprintf("- **公司背景**: %d/15 - %s\n", c.CompanyScore, c.CompanyReason))
	sb.WriteString(fmt.Sprintf("- **技术能力**: %d/25 - %s\n", c.TechScore, c.TechReason))
	sb.WriteString(fmt.Sprintf("- **项目经验**: %d/15 - %s\n\n", c.ProjectScore, c.ProjectReason))

	// 状态信息
	sb.WriteString("### 📝 状态信息\n")
	sb.WriteString(fmt.Sprintf("- **当前状态**: %s\n", c.Status))
	if c.Notes != "" {
		sb.WriteString(fmt.Sprintf("- **备注**: %s\n", c.Notes))
	}
	sb.WriteString(fmt.Sprintf("- **创建时间**: %s\n", c.CreatedAt.Format("2006-01-02 15:04")))

	return sb.String()
}

// createPushSession 创建新的推送会话
func (s *DingTalkService) createPushSession(candidates []models.Candidate) {
	// 生成会话ID（使用时间戳）
	sessionID := fmt.Sprintf("push_%d", time.Now().Unix())

	// 提取候选人ID列表（保持推送顺序）
	candidateIDs := make([]uint, len(candidates))
	for i, c := range candidates {
		candidateIDs[i] = c.ID
	}

	session := &PushSession{
		SessionID:    sessionID,
		CandidateIDs: candidateIDs,
		PushTime:     time.Now(),
	}

	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()

	// 保存到缓存
	s.sessionCache[sessionID] = session
	s.currentSessionID = sessionID

	s.log.Info("created push session",
		logging.KV("session_id", sessionID),
		logging.KV("candidate_count", len(candidateIDs)),
		logging.KV("candidate_ids", candidateIDs))
}

// getCurrentSession 获取当前最新的推送会话
func (s *DingTalkService) getCurrentSession() *PushSession {
	s.sessionMutex.RLock()
	defer s.sessionMutex.RUnlock()

	if s.currentSessionID == "" {
		return nil
	}

	return s.sessionCache[s.currentSessionID]
}

// cleanExpiredSessions 定期清理过期的会话（24小时前的）
func (s *DingTalkService) cleanExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour) // 每小时清理一次
	defer ticker.Stop()

	for range ticker.C {
		s.sessionMutex.Lock()

		expireTime := time.Now().Add(-24 * time.Hour)
		deletedCount := 0

		for sessionID, session := range s.sessionCache {
			if session.PushTime.Before(expireTime) {
				delete(s.sessionCache, sessionID)
				deletedCount++
			}
		}

		if deletedCount > 0 {
			s.log.Info("cleaned expired push sessions",
				logging.KV("deleted_count", deletedCount),
				logging.KV("remaining_count", len(s.sessionCache)))
		}

		s.sessionMutex.Unlock()
	}
}

// PushCandidates 推送候选人列表（统一的推送方法，供定时任务和API调用）
func (s *DingTalkService) PushCandidates(candidates []models.Candidate, config *models.DingTalkConfig) error {
	if len(candidates) == 0 {
		s.log.Info("no candidates to push")
		return nil
	}

	// 解析@的用户ID
	atUserIDs := s.parseAtUserIDs(config.AtUserIDs)

	// 构建Markdown消息
	content := s.buildCandidateListMarkdown(candidates, atUserIDs)

	// 发送消息（使用服务的client，确保一致性）
	if err := s.client.SendMarkdownMessage("📋 待面试候选人列表", content, atUserIDs, false); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	// 更新候选人通知状态
	s.updateNotifyStatus(candidates)

	// 创建新的推送会话
	s.createPushSession(candidates)

	s.log.Info("pushed candidates successfully",
		logging.KV("count", len(candidates)))

	return nil
}

// parseAtUserIDs 解析@用户ID列表
func (s *DingTalkService) parseAtUserIDs(atUserIDsStr string) []string {
	if atUserIDsStr == "" {
		return []string{}
	}

	atUserIDs := strings.Split(atUserIDsStr, ",")
	result := make([]string, 0, len(atUserIDs))
	for _, id := range atUserIDs {
		trimmed := strings.TrimSpace(id)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// updateNotifyStatus 更新候选人通知状态
func (s *DingTalkService) updateNotifyStatus(candidates []models.Candidate) {
	now := time.Now()
	for i := range candidates {
		candidates[i].NotifyCount++
		if candidates[i].FirstNotifyAt == nil {
			candidates[i].FirstNotifyAt = &now
		}
		candidates[i].LastNotifyAt = &now
		if err := s.repo.Update(&candidates[i]); err != nil {
			s.log.Error("update candidate notify status failed",
				logging.KV("id", candidates[i].ID),
				logging.KV("name", candidates[i].Name),
				logging.Err(err))
		}
	}
}

// PushEvaluationResult 推送评估结果（评估完成后自动调用）
// candidates: 本次评估的候选人列表
// config: 钉钉配置
// isAutoPush: 是否为自动推送（影响消息标题）
func (s *DingTalkService) PushEvaluationResult(candidates []models.Candidate, config *models.DingTalkConfig, isAutoPush bool) error {
	if len(candidates) == 0 {
		s.log.Info("no candidates to push for evaluation result")
		return nil
	}

	// 解析@的用户ID
	atUserIDs := s.parseAtUserIDs(config.AtUserIDs)

	// 构建评估结果Markdown消息
	content := s.buildEvaluationResultMarkdown(candidates, atUserIDs, isAutoPush)

	// 确定消息标题
	title := "📋 评估结果通知"
	if isAutoPush {
		title = "🔔 评估完成自动通知"
	}

	// 创建临时客户端发送消息（因为可能是不同用户的配置）
	client := dingtalk.NewClient(config, s.log)

	// 发送消息
	if err := client.SendMarkdownMessage(title, content, atUserIDs, false); err != nil {
		return fmt.Errorf("send evaluation result message: %w", err)
	}

	// 更新候选人通知状态
	s.updateNotifyStatus(candidates)

	// 创建新的推送会话
	s.createPushSession(candidates)

	s.log.Info("pushed evaluation result successfully",
		logging.KV("count", len(candidates)),
		logging.KV("is_auto_push", isAutoPush))

	return nil
}

// buildEvaluationResultMarkdown 构建评估结果Markdown消息
func (s *DingTalkService) buildEvaluationResultMarkdown(candidates []models.Candidate, atUserIds []string, isAutoPush bool) string {
	var sb strings.Builder

	// 标题区分自动推送和手动推送
	if isAutoPush {
		sb.WriteString("## 🔔 评估完成自动通知\n\n")
		sb.WriteString(fmt.Sprintf("> 本次评估完成 **%d** 位候选人，以下是评估结果：\n\n", len(candidates)))
	} else {
		sb.WriteString("## 📋 评估结果通知\n\n")
		sb.WriteString(fmt.Sprintf("> 共 **%d** 位候选人评估完成\n\n", len(candidates)))
	}

	// 候选人列表（按分数降序排列）
	for i, c := range candidates {
		gradeEmoji := s.getGradeEmoji(c.Grade)
		sb.WriteString(fmt.Sprintf("**%d. %s** %s\n", i+1, c.Name, gradeEmoji))
		sb.WriteString(fmt.Sprintf("- 评分: **%.1f** 分 | 评级: **%s**\n", c.TotalScore, c.Grade))
		sb.WriteString(fmt.Sprintf("- JD匹配: %d%%\n", c.JDMatch))
		sb.WriteString(fmt.Sprintf("- 建议: %s\n", c.Recommendation))
		sb.WriteString("\n")
	}

	sb.WriteString("---\n")
	sb.WriteString("💡 回复序号查看详情，如：`1` 或 `1,2,3`\n\n")

	// 在消息末尾添加@用户
	if len(atUserIds) > 0 {
		sb.WriteString("\n")
		for _, userId := range atUserIds {
			if userId != "" {
				sb.WriteString(fmt.Sprintf("@%s ", userId))
			}
		}
	}

	return sb.String()
}

// PushEvaluationResultByUser 根据用户ID获取配置并推送评估结果
// 返回值: pushed (是否推送成功), error (错误信息)
func (s *DingTalkService) PushEvaluationResultByUser(candidates []models.Candidate, userID uint) (bool, error) {
	if len(candidates) == 0 {
		return false, nil
	}

	// 获取用户的钉钉配置
	config, err := s.dtRepo.GetByUser(userID)
	if err != nil {
		s.log.Error("get dingtalk config failed",
			logging.KV("user_id", userID),
			logging.Err(err))
		return false, err
	}

	// 检查配置是否满足自动推送条件
	if config == nil {
		s.log.Debug("dingtalk config not found, skip auto push",
			logging.KV("user_id", userID))
		return false, nil
	}

	if !config.Enabled {
		s.log.Debug("dingtalk not enabled, skip auto push",
			logging.KV("user_id", userID))
		return false, nil
	}

	if !config.AutoPushOnComplete {
		s.log.Debug("auto push on complete not enabled, skip",
			logging.KV("user_id", userID))
		return false, nil
	}

	if config.Webhook == "" {
		s.log.Debug("webhook not configured, skip auto push",
			logging.KV("user_id", userID))
		return false, nil
	}

	// 执行推送
	if err := s.PushEvaluationResult(candidates, config, true); err != nil {
		s.log.Error("auto push evaluation result failed",
			logging.KV("user_id", userID),
			logging.KV("candidate_count", len(candidates)),
			logging.Err(err))
		return false, err
	}

	s.log.Info("auto push evaluation result succeeded",
		logging.KV("user_id", userID),
		logging.KV("candidate_count", len(candidates)))

	return true, nil
}
