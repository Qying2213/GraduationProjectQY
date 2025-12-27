package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"evaluator-service/internal/ai"
	"evaluator-service/internal/config"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/models"
	"evaluator-service/internal/repository"
	"evaluator-service/internal/utils"
)

type ResumeService struct {
	cfg   *config.Config
	log   *logging.Logger
	repo  *repository.CandidateRepository
	aiFac *ai.Factory
}

func NewResumeService(cfg *config.Config, log *logging.Logger, repo *repository.CandidateRepository) *ResumeService {
	return &ResumeService{cfg: cfg, log: log, repo: repo, aiFac: ai.NewFactory(cfg)}
}

// ==================== PDF 存储相关 ====================

// SaveResumePDF 保存简历PDF文件到服务器
// 返回保存的文件路径，如果保存失败返回空字符串（不阻断流程）
func (s *ResumeService) SaveResumePDF(pdfBytes []byte, userID uint, candidateName string) string {
	if len(pdfBytes) == 0 {
		return ""
	}

	// 构建存储目录: data/resumes/{userID}/
	resumeDir := filepath.Join(s.cfg.Storage.BaseDir, s.cfg.Storage.Resumes, fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(resumeDir, 0755); err != nil {
		s.log.Error("Failed to create resume directory", logging.Err(err), logging.KV("dir", resumeDir))
		return ""
	}

	// 生成唯一文件名: {timestamp}_{candidateName}.pdf
	safeName := sanitizeFilename(candidateName)
	filename := fmt.Sprintf("%d_%s.pdf", time.Now().UnixNano(), safeName)
	pdfPath := filepath.Join(resumeDir, filename)

	// 写入文件
	if err := os.WriteFile(pdfPath, pdfBytes, 0644); err != nil {
		s.log.Error("Failed to save resume PDF", logging.Err(err), logging.KV("path", pdfPath))
		return ""
	}

	s.log.Info("Resume PDF saved", logging.KV("path", pdfPath), logging.KV("size", len(pdfBytes)))
	return pdfPath
}

// sanitizeFilename 清理文件名中的非法字符
func sanitizeFilename(name string) string {
	// 移除或替换非法字符
	name = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '_'
		}
		return r
	}, name)
	name = strings.TrimSpace(name)
	if name == "" {
		name = "unknown"
	}
	return name
}

// cleanJSONString 清理 JSON 字符串，去掉 markdown 代码块标记，并尝试修复截断的 JSON
func cleanJSONString(s string) string {
	s = strings.TrimSpace(s)
	// 去掉开头的 ```json 或 ```
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
	}
	// 去掉结尾的 ```
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
	}
	s = strings.TrimSpace(s)

	// 尝试修复截断的 JSON：补全缺失的括号
	if s != "" && !json.Valid([]byte(s)) {
		// 统计括号数量
		openBraces := strings.Count(s, "{")
		closeBraces := strings.Count(s, "}")
		openBrackets := strings.Count(s, "[")
		closeBrackets := strings.Count(s, "]")

		// 补全缺失的括号
		for i := 0; i < openBrackets-closeBrackets; i++ {
			s += "]"
		}
		for i := 0; i < openBraces-closeBraces; i++ {
			s += "}"
		}

		// 如果还是无效，尝试在末尾加上常见的结束符
		if !json.Valid([]byte(s)) {
			s = strings.TrimRight(s, ", \t\n\r")
			if !json.Valid([]byte(s)) {
				s += "\"}]}"
			}
		}
	}

	return s
}

// extractCozeReportJSON 从 Coze 返回数据中提取完整报告 JSON
func extractCozeReportJSON(cozeData map[string]any) string {
	if cozeData == nil {
		return ""
	}
	if outputStr, ok := cozeData["output"].(string); ok && outputStr != "" {
		return cleanJSONString(outputStr)
	}
	if resultStr, ok := cozeData["result"].(string); ok {
		return cleanJSONString(resultStr)
	}
	return ""
}

// CozeEvaluationResult 封装 coze 评估的结果
type CozeEvaluationResult struct {
	ResumeMD           string
	JDMatch            models.JDMatchResult
	Requirement        models.RequirementResult
	Score              models.ScoringResult
	InterviewQuestions []models.InterviewQuestion
	Recommendation     string
}

func (s *ResumeService) evaluateWithCoze(cozeData map[string]any) (*CozeEvaluationResult, error) {
	client := ai.NewCozeClient(cozeData)

	resumeMD, err := client.Structure("")
	if err != nil {
		return nil, fmt.Errorf("failed to structure resume: %w", err)
	}

	jdMatch, err := client.EvaluateJD(resumeMD, "")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate JD match: %w", err)
	}

	reqRes, err := client.EvaluateRequirement(resumeMD)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate requirements: %w", err)
	}

	scoreRes, err := client.Score(resumeMD, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to score resume: %w", err)
	}

	interviewQuestions, _ := client.GenerateInterviewQuestions(resumeMD, models.EvaluationResult{
		JDMatch:        jdMatch,
		Requirement:    reqRes,
		Recommendation: getRecommendationFromCoze(cozeData, scoreRes, jdMatch, reqRes),
	})

	rec := getRecommendationFromCoze(cozeData, scoreRes, jdMatch, reqRes)

	return &CozeEvaluationResult{
		ResumeMD:           resumeMD,
		JDMatch:            jdMatch,
		Requirement:        reqRes,
		Score:              scoreRes,
		InterviewQuestions: interviewQuestions,
		Recommendation:     rec,
	}, nil
}

func getRecommendationFromCoze(cozeData map[string]any, score models.ScoringResult, jd models.JDMatchResult, req models.RequirementResult) string {
	if cozeData == nil {
		return recommendation(score, jd, req)
	}

	resultStr, ok := cozeData["output"].(string)
	if !ok {
		resultStr, ok = cozeData["result"].(string)
	}
	if !ok {
		return recommendation(score, jd, req)
	}

	// 清理JSON字符串（去掉markdown代码块标记等）
	resultStr = cleanJSONString(resultStr)

	var resultData map[string]any
	if err := json.Unmarshal([]byte(resultStr), &resultData); err != nil {
		return recommendation(score, jd, req)
	}

	if rec, ok := resultData["录用建议"].(map[string]any); ok {
		if conclusion, ok := rec["结论"].(string); ok && conclusion != "" {
			return conclusion
		}
	}

	return recommendation(score, jd, req)
}

type EvaluateOutput struct {
	Candidate    *models.Candidate
	ReportMD     string
	ReportHTML   string
	ReportMDPath string
}

func (s *ResumeService) EvaluateSingle(pdfPath, filename, jd, criteria string, cozeData map[string]any) (*EvaluateOutput, error) {
	return s.EvaluateSingleWithUser(pdfPath, filename, jd, criteria, cozeData, 0)
}

func (s *ResumeService) EvaluateSingleWithUser(pdfPath, filename, jd, criteria string, cozeData map[string]any, userID uint) (*EvaluateOutput, error) {
	cozeRes, err := s.evaluateWithCoze(cozeData)
	if err != nil {
		return nil, err
	}

	candName := utils.ExtractCandidateName(filename)
	eval := models.EvaluationResult{
		JDMatch:        cozeRes.JDMatch,
		Requirement:    cozeRes.Requirement,
		Recommendation: cozeRes.Recommendation,
	}

	reportMD := s.buildReport(cozeRes.ResumeMD, eval, cozeRes.Score, filename, candName)
	reportHTML, err := utils.MarkdownToHTML(reportMD)
	if err != nil {
		return nil, err
	}

	cozeReportJSON := extractCozeReportJSON(cozeData)

	cand := &models.Candidate{
		UserID:           userID,
		Name:             candName,
		Filename:         filename,
		TotalScore:       cozeRes.Score.TotalScore,
		Grade:            cozeRes.Score.Grade,
		JDMatch:          cozeRes.JDMatch.Score,
		AgeScore:         cozeRes.Score.AgeScore,
		ExperienceScore:  cozeRes.Score.ExperienceScore,
		EducationScore:   cozeRes.Score.EducationScore,
		CompanyScore:     cozeRes.Score.CompanyScore,
		TechScore:        cozeRes.Score.TechScore,
		ProjectScore:     cozeRes.Score.ProjectScore,
		AgeReason:        cozeRes.Score.AgeReason,
		ExperienceReason: cozeRes.Score.ExperienceReason,
		EducationReason:  cozeRes.Score.EducationReason,
		CompanyReason:    cozeRes.Score.CompanyReason,
		TechReason:       cozeRes.Score.TechReason,
		ProjectReason:    cozeRes.Score.ProjectReason,
		Recommendation:   cozeRes.Recommendation,
		ReportMarkdown:   reportMD,
		ResumeMarkdown:   cozeRes.ResumeMD,
		CozeReportJSON:   cozeReportJSON,
	}
	if err := s.repo.Create(cand); err != nil {
		return nil, err
	}

	return &EvaluateOutput{Candidate: cand, ReportMD: reportMD, ReportHTML: reportHTML}, nil
}

func (s *ResumeService) EvaluateSingleBytes(pdfBytes []byte, filename, jd, criteria string, cozeData map[string]any) (*EvaluateOutput, error) {
	return s.EvaluateSingleBytesWithUser(pdfBytes, filename, jd, criteria, cozeData, 0)
}

func (s *ResumeService) EvaluateSingleBytesWithUser(pdfBytes []byte, filename, jd, criteria string, cozeData map[string]any, userID uint) (*EvaluateOutput, error) {
	cozeRes, err := s.evaluateWithCoze(cozeData)
	if err != nil {
		return nil, err
	}

	candName := utils.ExtractCandidateName(filename)
	eval := models.EvaluationResult{
		JDMatch:        cozeRes.JDMatch,
		Requirement:    cozeRes.Requirement,
		Recommendation: cozeRes.Recommendation,
	}

	reportMD := s.buildReport(cozeRes.ResumeMD, eval, cozeRes.Score, filename, candName)
	reportHTML, err := utils.MarkdownToHTML(reportMD)
	if err != nil {
		return nil, err
	}

	cozeReportJSON := extractCozeReportJSON(cozeData)

	// 保存PDF文件到服务器
	pdfPath := s.SaveResumePDF(pdfBytes, userID, candName)

	// 为单个评估生成唯一的ApplyID，避免唯一约束冲突
	applyID := fmt.Sprintf("manual_%d", time.Now().UnixNano())

	cand := &models.Candidate{
		UserID:           userID,
		ApplyID:          applyID,
		Name:             candName,
		Filename:         filename,
		PDFPath:          pdfPath, // 设置PDF文件路径
		TotalScore:       cozeRes.Score.TotalScore,
		Grade:            cozeRes.Score.Grade,
		JDMatch:          cozeRes.JDMatch.Score,
		AgeScore:         cozeRes.Score.AgeScore,
		ExperienceScore:  cozeRes.Score.ExperienceScore,
		EducationScore:   cozeRes.Score.EducationScore,
		CompanyScore:     cozeRes.Score.CompanyScore,
		TechScore:        cozeRes.Score.TechScore,
		ProjectScore:     cozeRes.Score.ProjectScore,
		AgeReason:        cozeRes.Score.AgeReason,
		ExperienceReason: cozeRes.Score.ExperienceReason,
		EducationReason:  cozeRes.Score.EducationReason,
		CompanyReason:    cozeRes.Score.CompanyReason,
		TechReason:       cozeRes.Score.TechReason,
		ProjectReason:    cozeRes.Score.ProjectReason,
		Recommendation:   cozeRes.Recommendation,
		ReportMarkdown:   reportMD,
		ResumeMarkdown:   cozeRes.ResumeMD,
		CozeReportJSON:   cozeReportJSON,
	}
	if err := s.repo.Create(cand); err != nil {
		return nil, err
	}

	return &EvaluateOutput{Candidate: cand, ReportMD: reportMD, ReportHTML: reportHTML}, nil
}

func recommendation(score models.ScoringResult, jd models.JDMatchResult, req models.RequirementResult) string {
	if score.AgeScore == 0 {
		return "不推荐（年龄黑名单）"
	}
	if jd.Score >= 70 && req.OverallPass {
		return "推荐"
	}
	if jd.Score >= 50 || req.OverallPass {
		return "待定"
	}
	return "不推荐"
}

func (s *ResumeService) buildReport(resumeMD string, eval models.EvaluationResult, score models.ScoringResult, pdfName, candidate string) string {
	b, _ := json.Marshal(eval.JDMatch.MatchedSkills)
	ms := string(b)
	b, _ = json.Marshal(eval.JDMatch.MissingSkills)
	mis := string(b)
	return strings.Join([]string{
		"# " + candidate + " - 简历评估报告",
		"> 生成时间：" + time.Now().Format("2006-01-02 15:04:05"),
		"> 候选人：" + candidate,
		"> 简历文件：" + pdfName,
		"\n## 📊 综合评分",
		fmt.Sprintf("年龄: %d/10\n", score.AgeScore),
		fmt.Sprintf("经验: %d/25\n", score.ExperienceScore),
		fmt.Sprintf("学历: %d/20\n", score.EducationScore),
		fmt.Sprintf("公司: %d/15\n", score.CompanyScore),
		fmt.Sprintf("技术: %d/25\n", score.TechScore),
		fmt.Sprintf("项目: %d/15\n", score.ProjectScore),
		fmt.Sprintf("总分: %.1f  评级: %s\n", score.TotalScore, score.Grade),
		"\n## 🎯 JD匹配度分析\n",
		fmt.Sprintf("匹配分数：%d/100\n", eval.JDMatch.Score),
		"匹配技能: " + ms + "\n",
		"缺失技能: " + mis + "\n",
		"匹配总结: " + eval.JDMatch.Summary + "\n",
		"\n## ✅ 用人标准评估\n",
		fmt.Sprintf("学历: %v  详情: %s\n", eval.Requirement.EducationPass, eval.Requirement.EducationDetail),
		fmt.Sprintf("经验: %v  年限: %.2f\n", eval.Requirement.ExperiencePass, valueOrZero(eval.Requirement.ExperienceYears)),
		fmt.Sprintf("黑名单: %v  命中: %v\n", eval.Requirement.BlacklistPass, eval.Requirement.BlacklistHits),
		fmt.Sprintf("综合结果: %v\n", eval.Requirement.OverallPass),
		"\n## 📄 简历结构化内容\n\n" + resumeMD + "\n",
		"\n## 📌 最终建议\n\n推荐结果：" + eval.Recommendation + "\n",
	}, "\n")
}

func valueOrZero(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

// ==================== 重复检测相关 ====================

// DuplicateCandidate 重复候选人信息
type DuplicateCandidate struct {
	Name        string    `json:"name"`
	ApplyID     string    `json:"apply_id"`
	EvaluatedAt time.Time `json:"evaluated_at"`
	TotalScore  float64   `json:"total_score"`
	Grade       string    `json:"grade"`
}

// DuplicateCheckResult 重复检测结果
type DuplicateCheckResult struct {
	Duplicates     []DuplicateCandidate // 已存在的候选人
	NewApplyIDs    []string             // 新候选人的 ApplyID 列表
	DuplicateCount int                  // 重复数量
	TotalFetched   int                  // 总拉取数量
}

// CheckDuplicates 检测候选人列表中哪些已经评估过
// applyIDs: 从招聘系统拉取的候选人 ApplyID 列表
// applyIDToName: ApplyID 到姓名的映射
// userID: 当前用户ID
func (s *ResumeService) CheckDuplicates(applyIDs []string, applyIDToName map[string]string, userID uint) (*DuplicateCheckResult, error) {
	result := &DuplicateCheckResult{
		Duplicates:   make([]DuplicateCandidate, 0),
		NewApplyIDs:  make([]string, 0),
		TotalFetched: len(applyIDs),
	}

	if len(applyIDs) == 0 {
		return result, nil
	}

	// 查询数据库中已存在的候选人
	existingCandidates, err := s.repo.FindByApplyIDsAndUser(applyIDs, userID)
	if err != nil {
		s.log.Error("Failed to check duplicates", logging.Err(err))
		// 查询失败时，降级为不检测重复
		result.NewApplyIDs = applyIDs
		return result, nil
	}

	// 构建已存在的 ApplyID 集合
	existingMap := make(map[string]*models.Candidate)
	for i := range existingCandidates {
		c := &existingCandidates[i]
		existingMap[c.ApplyID] = c
	}

	// 分类：重复 vs 新增
	for _, applyID := range applyIDs {
		if applyID == "" {
			continue
		}
		if existing, ok := existingMap[applyID]; ok {
			result.Duplicates = append(result.Duplicates, DuplicateCandidate{
				Name:        existing.Name,
				ApplyID:     existing.ApplyID,
				EvaluatedAt: existing.CreatedAt,
				TotalScore:  existing.TotalScore,
				Grade:       existing.Grade,
			})
		} else {
			result.NewApplyIDs = append(result.NewApplyIDs, applyID)
		}
	}

	result.DuplicateCount = len(result.Duplicates)
	return result, nil
}

// EvaluateSingleBytesWithApplyID 评估单个候选人（带 ApplyID）
// 如果 existingCandidate 不为 nil，则更新已有记录；否则创建新记录
func (s *ResumeService) EvaluateSingleBytesWithApplyID(pdfBytes []byte, filename, jd, criteria string, cozeData map[string]any, userID uint, applyID string, existingCandidate *models.Candidate) (*EvaluateOutput, error) {
	cozeRes, err := s.evaluateWithCoze(cozeData)
	if err != nil {
		return nil, err
	}

	candName := utils.ExtractCandidateName(filename)
	eval := models.EvaluationResult{
		JDMatch:        cozeRes.JDMatch,
		Requirement:    cozeRes.Requirement,
		Recommendation: cozeRes.Recommendation,
	}

	reportMD := s.buildReport(cozeRes.ResumeMD, eval, cozeRes.Score, filename, candName)
	reportHTML, err := utils.MarkdownToHTML(reportMD)
	if err != nil {
		return nil, err
	}

	cozeReportJSON := extractCozeReportJSON(cozeData)

	// 保存PDF文件到服务器
	pdfPath := s.SaveResumePDF(pdfBytes, userID, candName)

	var cand *models.Candidate
	if existingCandidate != nil {
		// 重新评估：删除旧的PDF文件
		if existingCandidate.PDFPath != "" {
			if err := os.Remove(existingCandidate.PDFPath); err != nil {
				s.log.Warn("Failed to remove old PDF file", logging.Err(err), logging.KV("path", existingCandidate.PDFPath))
			}
		}

		// 更新已有记录（保留 ID 和 CreatedAt）
		cand = existingCandidate
		cand.Name = candName
		cand.Filename = filename
		cand.PDFPath = pdfPath // 更新PDF文件路径
		cand.TotalScore = cozeRes.Score.TotalScore
		cand.Grade = cozeRes.Score.Grade
		cand.JDMatch = cozeRes.JDMatch.Score
		cand.AgeScore = cozeRes.Score.AgeScore
		cand.ExperienceScore = cozeRes.Score.ExperienceScore
		cand.EducationScore = cozeRes.Score.EducationScore
		cand.CompanyScore = cozeRes.Score.CompanyScore
		cand.TechScore = cozeRes.Score.TechScore
		cand.ProjectScore = cozeRes.Score.ProjectScore
		cand.AgeReason = cozeRes.Score.AgeReason
		cand.ExperienceReason = cozeRes.Score.ExperienceReason
		cand.EducationReason = cozeRes.Score.EducationReason
		cand.CompanyReason = cozeRes.Score.CompanyReason
		cand.TechReason = cozeRes.Score.TechReason
		cand.ProjectReason = cozeRes.Score.ProjectReason
		cand.Recommendation = cozeRes.Recommendation
		cand.ReportMarkdown = reportMD
		cand.ResumeMarkdown = cozeRes.ResumeMD
		cand.CozeReportJSON = cozeReportJSON
		// UpdatedAt 会被 GORM 自动更新
		if err := s.repo.Update(cand); err != nil {
			return nil, err
		}
	} else {
		// 创建新记录
		cand = &models.Candidate{
			UserID:           userID,
			ApplyID:          applyID,
			Name:             candName,
			Filename:         filename,
			PDFPath:          pdfPath, // 设置PDF文件路径
			TotalScore:       cozeRes.Score.TotalScore,
			Grade:            cozeRes.Score.Grade,
			JDMatch:          cozeRes.JDMatch.Score,
			AgeScore:         cozeRes.Score.AgeScore,
			ExperienceScore:  cozeRes.Score.ExperienceScore,
			EducationScore:   cozeRes.Score.EducationScore,
			CompanyScore:     cozeRes.Score.CompanyScore,
			TechScore:        cozeRes.Score.TechScore,
			ProjectScore:     cozeRes.Score.ProjectScore,
			AgeReason:        cozeRes.Score.AgeReason,
			ExperienceReason: cozeRes.Score.ExperienceReason,
			EducationReason:  cozeRes.Score.EducationReason,
			CompanyReason:    cozeRes.Score.CompanyReason,
			TechReason:       cozeRes.Score.TechReason,
			ProjectReason:    cozeRes.Score.ProjectReason,
			Recommendation:   cozeRes.Recommendation,
			ReportMarkdown:   reportMD,
			ResumeMarkdown:   cozeRes.ResumeMD,
			CozeReportJSON:   cozeReportJSON,
		}
		if err := s.repo.Create(cand); err != nil {
			return nil, err
		}
	}

	return &EvaluateOutput{Candidate: cand, ReportMD: reportMD, ReportHTML: reportHTML}, nil
}

// GetExistingCandidateByApplyID 根据 ApplyID 获取已存在的候选人
func (s *ResumeService) GetExistingCandidateByApplyID(applyID string, userID uint) (*models.Candidate, error) {
	return s.repo.FindByApplyIDAndUser(applyID, userID)
}
