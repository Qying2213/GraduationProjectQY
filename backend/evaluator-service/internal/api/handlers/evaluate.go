package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"evaluator-service/internal/api/middleware"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/models"
	"evaluator-service/internal/repository"
	"evaluator-service/internal/script"
	"evaluator-service/internal/thirdparty/coze"
	"evaluator-service/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *Handlers) Evaluate(c *gin.Context) {
	file, err := c.FormFile("resume")
	if err != nil {
		bad(c, err)
		return
	}
	userID := middleware.GetUserID(c)

	// 获取 JD：优先使用 position_id，否则使用直接传入的 jd
	var jd string
	positionIDStr := c.PostForm("position_id")
	if positionIDStr != "" {
		// 从数据库获取岗位 JD
		positionSvc := h.getPositionService()
		var positionID uint
		if _, err := fmt.Sscanf(positionIDStr, "%d", &positionID); err != nil {
			bad(c, fmt.Errorf("无效的岗位ID"))
			return
		}
		position, err := positionSvc.GetByIDAndUser(positionID, userID)
		if err != nil {
			fail(c, err)
			return
		}
		if position == nil {
			c.JSON(400, gin.H{"error": "岗位不存在"})
			return
		}
		jd = position.GetJDText()
	} else {
		jd = c.PostForm("jd")
	}

	tmpPath, err := utils.SaveUploadedTemp(utils.Join(h.cfg.Storage.BaseDir, h.cfg.Storage.TempDir), file)
	if err != nil {
		fail(c, err)
		return
	}
	defer utils.RemoveQuiet(tmpPath)

	// Read PDF bytes and run Coze workflow (single upload also goes through Coze)
	pdfBytes, err := os.ReadFile(tmpPath)
	if err != nil {
		fail(c, err)
		return
	}
	name := utils.ExtractCandidateName(file.Filename)
	h.log.Info("Starting evaluation",
		logging.KV("candidate_name", name),
		logging.KV("filename", file.Filename),
		logging.KV("pdf_size", len(pdfBytes)),
		logging.KV("jd_len", len(jd)))

	cozeCtx, cancel := context.WithTimeout(c.Request.Context(), 300*time.Second)
	defer cancel()

	h.log.Info("Calling Coze workflow")
	cozeData, cozeErr := coze.RunWorkflow(cozeCtx, name, jd, pdfBytes)

	// 如果 Coze 调用失败，记录错误但继续使用空数据（会触发错误）
	if cozeErr != nil {
		h.log.Error("Coze workflow failed", logging.Err(cozeErr))
		cozeData = nil
	} else {
		h.log.Info("Coze workflow completed successfully",
			logging.KV("has_data", cozeData != nil),
			logging.KV("data_keys", getCozeDataKeys(cozeData)))
	}

	// Continue evaluation using in-memory bytes and Coze data
	h.log.Info("Starting resume evaluation")
	out, err := h.svc.EvaluateSingleBytesWithUser(pdfBytes, file.Filename, jd, "", cozeData, userID)
	if err != nil {
		h.log.Error("Resume evaluation failed", logging.Err(err))
		fail(c, err)
		return
	}

	h.log.Info("Evaluation completed successfully",
		logging.KV("candidate_id", out.Candidate.ID),
		logging.KV("total_score", out.Candidate.TotalScore))

	// 评估完成后触发自动推送（异步，不阻塞响应）
	cand := out.Candidate
	go h.triggerAutoPushIfEnabled(userID, []models.Candidate{*cand})

	ok(c, gin.H{
		"success":        true,
		"candidate_id":   cand.ID,
		"candidate_name": cand.Name,
		"score": gin.H{
			"total":      cand.TotalScore,
			"grade":      cand.Grade,
			"age":        cand.AgeScore,
			"experience": cand.ExperienceScore,
			"education":  cand.EducationScore,
			"company":    cand.CompanyScore,
			"tech":       cand.TechScore,
			"project":    cand.ProjectScore,
		},
		"jd_match":         cand.JDMatch,
		"recommendation":   cand.Recommendation,
		"report_markdown":  out.ReportMD,
		"report_html":      out.ReportHTML,
		"coze_report_json": cand.CozeReportJSON,
		"saved_path":       out.ReportMDPath,
	})
}

// EvaluationResultType 评估结果类型
type EvaluationResultType string

const (
	ResultTypeNew         EvaluationResultType = "new"         // 新评估
	ResultTypeReevaluated EvaluationResultType = "reevaluated" // 重新评估
	ResultTypeSkipped     EvaluationResultType = "skipped"     // 跳过
)

// EvaluateBatch: 使用 SSE 流式响应，实时推送评估进度
func (h *Handlers) EvaluateBatch(c *gin.Context) {
	jdFromReq := c.PostForm("jd")
	criteria := c.PostForm("criteria")
	forceReevaluate := c.PostForm("force_reevaluate") == "true" // 是否强制重新评估
	userID := middleware.GetUserID(c)

	// 使用用户登录时存储的凭据
	repoCred := repository.NewCredentialRepository(h.repo.DB())
	cred, err := repoCred.GetByOrgAndUser("motern", userID)
	if err != nil {
		fail(c, err)
		return
	}
	if cred == nil {
		c.JSON(400, gin.H{"error": "凭据未找到，请重新登录", "code": "CREDENTIALS_NOT_FOUND"})
		return
	}
	// 解密密码
	encKey := getEncryptionKey(h.cfg.Credentials.EncKey)
	password, decErr := utils.DecryptAESGCM(encKey, cred.PasswordCipher)
	if decErr != nil {
		h.log.Error("decrypt credential failed", logging.KV("org", cred.Org), logging.Err(decErr))
		c.JSON(400, gin.H{"error": "凭据解密失败，请重新登录"})
		return
	}

	// 注入给 Python 的环境变量
	env := map[string]string{
		"WT_USERNAME": cred.Account,
		"WT_PASSWORD": password,
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// SSE 发送函数
	sendSSE := func(event string, data interface{}) {
		jsonData, _ := json.Marshal(data)
		fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, string(jsonData))
		c.Writer.Flush()
	}

	// 发送抓取开始事件
	sendSSE("status", gin.H{"stage": "fetching", "message": "正在从招聘系统获取简历..."})

	// 由 Python 脚本内部自动完成分页抓取，这里只需调用一次
	h.log.Info("Starting Wintalent fetch")
	fetchCtx, cancel := script.WithTimeout(c.Request.Context(), 300*time.Second)
	allItems, err := script.RunWintalentFetchWithEnv(fetchCtx, "internal/script/wintalent_fetch.py", env)
	cancel()
	if err != nil {
		h.log.Error("Wintalent fetch failed", logging.Err(err))
		sendSSE("error", gin.H{"error": err.Error()})
		return
	}

	total := len(allItems)
	h.log.Info("Wintalent fetched items", logging.KV("items_total", total))

	// 构建 ApplyID 列表和映射
	applyIDs := make([]string, 0, total)
	applyIDToName := make(map[string]string)
	applyIDToItem := make(map[string]script.WintalentItem)
	for _, item := range allItems {
		if item.ApplyID != "" {
			applyIDs = append(applyIDs, item.ApplyID)
			applyIDToName[item.ApplyID] = item.Name
			applyIDToItem[item.ApplyID] = item
		}
	}

	// 检测重复候选人
	sendSSE("status", gin.H{"stage": "checking", "message": "正在检测已评估的候选人..."})
	dupResult, _ := h.svc.CheckDuplicates(applyIDs, applyIDToName, userID)

	// 获取前端传递的跳过标志（用户已确认跳过重复候选人）
	skipDuplicates := c.PostForm("skip_duplicates") == "true"

	// 如果有重复且未强制重新评估且用户未确认跳过，发送重复检测事件
	if dupResult.DuplicateCount > 0 && !forceReevaluate && !skipDuplicates {
		// 构建重复候选人信息
		duplicatesInfo := make([]gin.H, 0, len(dupResult.Duplicates))
		for _, dup := range dupResult.Duplicates {
			duplicatesInfo = append(duplicatesInfo, gin.H{
				"name":         dup.Name,
				"apply_id":     dup.ApplyID,
				"evaluated_at": dup.EvaluatedAt.Format("2006-01-02 15:04"),
				"total_score":  dup.TotalScore,
				"grade":        dup.Grade,
			})
		}

		// 发送重复检测事件，等待前端确认
		sendSSE("duplicates_found", gin.H{
			"duplicates":      duplicatesInfo,
			"total_fetched":   dupResult.TotalFetched,
			"duplicate_count": dupResult.DuplicateCount,
			"new_count":       len(dupResult.NewApplyIDs),
		})
		return // 等待前端确认后重新调用
	}

	// 确定要评估的候选人列表
	var itemsToEvaluate []script.WintalentItem
	var skipApplyIDs map[string]bool

	if forceReevaluate {
		// 强制重新评估：评估所有候选人
		itemsToEvaluate = allItems
		skipApplyIDs = make(map[string]bool)
	} else {
		// 跳过重复：只评估新候选人
		skipApplyIDs = make(map[string]bool)
		for _, dup := range dupResult.Duplicates {
			skipApplyIDs[dup.ApplyID] = true
		}
		itemsToEvaluate = make([]script.WintalentItem, 0, len(dupResult.NewApplyIDs))
		for _, item := range allItems {
			if !skipApplyIDs[item.ApplyID] {
				itemsToEvaluate = append(itemsToEvaluate, item)
			}
		}
	}

	// 如果所有候选人都已评估且不强制重新评估
	if len(itemsToEvaluate) == 0 {
		sendSSE("complete", gin.H{
			"success":       true,
			"total":         0,
			"results":       []interface{}{},
			"skipped_count": dupResult.DuplicateCount,
			"message":       "所有候选人均已评估，无需重复操作",
		})
		return
	}

	evalTotal := len(itemsToEvaluate)
	h.log.Info("Starting evaluation", logging.KV("eval_total", evalTotal), logging.KV("skipped", dupResult.DuplicateCount))

	// 发送总数
	sendSSE("total", gin.H{"total": evalTotal, "skipped": dupResult.DuplicateCount})
	sendSSE("status", gin.H{"stage": "evaluating", "message": fmt.Sprintf("开始评估 %d 份简历...", evalTotal)})

	type res struct {
		Filename        string               `json:"filename"`
		CandidateID     uint                 `json:"candidate_id"`
		CandidateName   string               `json:"candidate_name"`
		ApplyID         string               `json:"apply_id"`
		TotalScore      float64              `json:"total_score"`
		Grade           string               `json:"grade"`
		JDMatch         int                  `json:"jd_match"`
		AgeScore        int                  `json:"age_score"`
		ExperienceScore int                  `json:"experience_score"`
		EducationScore  int                  `json:"education_score"`
		CompanyScore    int                  `json:"company_score"`
		TechScore       int                  `json:"tech_score"`
		ProjectScore    int                  `json:"project_score"`
		Recommendation  string               `json:"recommendation"`
		ReportHTML      string               `json:"report_html"`
		ReportMarkdown  string               `json:"report_markdown"`
		CozeReportJSON  string               `json:"coze_report_json"`
		SavedPath       string               `json:"saved_path"`
		Rank            int                  `json:"rank"`
		ResultType      EvaluationResultType `json:"result_type"`
		Error           string               `json:"error,omitempty"`
	}

	concurrency := h.cfg.Batch.Concurrency
	if concurrency <= 0 {
		concurrency = 5
	}
	sem := make(chan struct{}, concurrency)
	var g errgroup.Group
	results := make([]res, evalTotal)

	// 进度计数器
	var completed int32

	for i := range itemsToEvaluate {
		i := i
		sem <- struct{}{}
		g.Go(func() error {
			defer func() { <-sem }()
			item := itemsToEvaluate[i]

			// Build jd_text: prefer request JD, else from item.JD
			jdText := jdFromReq
			if strings.TrimSpace(jdText) == "" && item.JD != nil {
				if sc, ok := item.JD["serviceCondition"].(string); ok && sc != "" {
					jdText = sc
				}
				if wc, ok := item.JD["workContent"].(string); ok && wc != "" {
					if jdText != "" {
						jdText += "\n\n"
					}
					jdText += wc
				}
			}

			// Decode PDF bytes
			var pdfBytes []byte
			if item.ResumePDFB64 != "" {
				b, err := base64.StdEncoding.DecodeString(item.ResumePDFB64)
				if err == nil {
					pdfBytes = b
				}
			}
			name := strings.TrimSpace(item.Name)
			if name == "" {
				name = "unknown"
			}
			filename := name + ".pdf"

			// 检查是否是重新评估（已存在的候选人）
			var existingCandidate *models.Candidate
			var resultType EvaluationResultType = ResultTypeNew
			if forceReevaluate && item.ApplyID != "" {
				existingCandidate, _ = h.svc.GetExistingCandidateByApplyID(item.ApplyID, userID)
				if existingCandidate != nil {
					resultType = ResultTypeReevaluated
				}
			}

			// Call Coze with in-memory bytes (no temp file)
			cozeCtx, cancel := context.WithTimeout(c.Request.Context(), 300*time.Second)
			cozeData, cozeErr := coze.RunWorkflow(cozeCtx, name, jdText, pdfBytes)
			cancel()
			// Coze 错误不阻断后续评估流程，仅记录
			if cozeErr != nil {
				h.log.Error("Coze workflow failed", logging.Err(cozeErr), logging.KV("candidate", name))
				cozeData = nil
			}

			// Continue evaluation using in-memory bytes with ApplyID
			out, evalErr := h.svc.EvaluateSingleBytesWithApplyID(pdfBytes, filename, jdText, criteria, cozeData, userID, item.ApplyID, existingCandidate)
			if evalErr != nil {
				results[i] = res{
					Filename:      filename,
					CandidateName: utils.ExtractCandidateName(filename),
					ApplyID:       item.ApplyID,
					ResultType:    resultType,
					Error:         firstErrMsg(cozeErr, evalErr),
				}
			} else {
				cand := out.Candidate
				results[i] = res{
					Filename:        filename,
					CandidateID:     cand.ID,
					CandidateName:   cand.Name,
					ApplyID:         item.ApplyID,
					TotalScore:      cand.TotalScore,
					Grade:           cand.Grade,
					JDMatch:         cand.JDMatch,
					AgeScore:        cand.AgeScore,
					ExperienceScore: cand.ExperienceScore,
					EducationScore:  cand.EducationScore,
					CompanyScore:    cand.CompanyScore,
					TechScore:       cand.TechScore,
					ProjectScore:    cand.ProjectScore,
					Recommendation:  cand.Recommendation,
					ReportHTML:      out.ReportHTML,
					ReportMarkdown:  out.ReportMD,
					CozeReportJSON:  cand.CozeReportJSON,
					SavedPath:       out.ReportMDPath,
					ResultType:      resultType,
				}
			}

			// 更新进度并发送 SSE 事件
			done := atomic.AddInt32(&completed, 1)
			sendSSE("progress", gin.H{
				"completed":   done,
				"total":       evalTotal,
				"percent":     int(float64(done) / float64(evalTotal) * 100),
				"current":     name,
				"result_type": resultType,
			})

			return nil
		})
	}
	_ = g.Wait()

	// 排序并设置排名
	sort.Slice(results, func(i, j int) bool { return results[i].TotalScore > results[j].TotalScore })
	for i := range results {
		results[i].Rank = i + 1
	}

	// 统计结果类型
	var newCount, reevaluatedCount int
	for _, r := range results {
		switch r.ResultType {
		case ResultTypeNew:
			newCount++
		case ResultTypeReevaluated:
			reevaluatedCount++
		}
	}

	// 收集成功评估的候选人用于自动推送
	var successCandidates []models.Candidate
	for _, r := range results {
		if r.Error == "" && r.CandidateID > 0 {
			successCandidates = append(successCandidates, models.Candidate{
				ID:             r.CandidateID,
				Name:           r.CandidateName,
				TotalScore:     r.TotalScore,
				Grade:          r.Grade,
				JDMatch:        r.JDMatch,
				Recommendation: r.Recommendation,
			})
		}
	}

	// 批量评估完成后触发自动推送（异步，不阻塞 SSE 响应）
	if len(successCandidates) > 0 {
		go h.triggerAutoPushIfEnabled(userID, successCandidates)
	}

	// 发送最终结果
	sendSSE("complete", gin.H{
		"success":           true,
		"total":             evalTotal,
		"results":           results,
		"new_count":         newCount,
		"reevaluated_count": reevaluatedCount,
		"skipped_count":     dupResult.DuplicateCount,
	})
}

func firstErrMsg(errs ...error) string {
	for _, e := range errs {
		if e != nil {
			return e.Error()
		}
	}
	return ""
}

func getCozeDataKeys(data map[string]interface{}) []string {
	if data == nil {
		return []string{}
	}
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

// getEncryptionKey 获取加密密钥，确保 32 字节长度
func getEncryptionKey(configKey string) []byte {
	key := configKey
	if key == "" {
		key = "resume-evaluator-default-enc-key"
	}
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded, keyBytes)
		return padded
	}
	return keyBytes[:32]
}

// triggerAutoPushIfEnabled 检查配置并触发自动推送
// 该方法设计为异步调用，推送失败不影响主流程
func (h *Handlers) triggerAutoPushIfEnabled(userID uint, candidates []models.Candidate) {
	if len(candidates) == 0 {
		return
	}

	// 检查 dtService 是否可用
	if h.dtService == nil {
		h.log.Warn("dingtalk service not available, skip auto push")
		return
	}

	// 调用 DingTalkService 的推送方法
	pushed, err := h.dtService.PushEvaluationResultByUser(candidates, userID)
	if err != nil {
		h.log.Error("auto push failed",
			logging.KV("user_id", userID),
			logging.KV("candidate_count", len(candidates)),
			logging.Err(err))
		return
	}

	if pushed {
		h.log.Info("auto push triggered successfully",
			logging.KV("user_id", userID),
			logging.KV("candidate_count", len(candidates)))
	}
}

// EvaluateBatchGraduate: 从毕业设计后台获取简历并评估，使用 SSE 流式响应
func (h *Handlers) EvaluateBatchGraduate(c *gin.Context) {
	jdFromReq := c.PostForm("jd")
	criteria := c.PostForm("criteria")
	forceReevaluate := c.PostForm("force_reevaluate") == "true"
	userID := middleware.GetUserID(c)

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// SSE 发送函数
	sendSSE := func(event string, data interface{}) {
		jsonData, _ := json.Marshal(data)
		fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, string(jsonData))
		c.Writer.Flush()
	}

	// 发送抓取开始事件
	sendSSE("status", gin.H{"stage": "fetching", "message": "正在从毕业设计后台获取简历..."})

	// 设置环境变量（毕业设计后台地址）
	env := map[string]string{
		"GRADUATE_API_URL": h.cfg.Graduate.APIUrl,
	}

	// 调用 graduate_fetch.py 获取简历
	h.log.Info("Starting Graduate fetch", logging.KV("api_url", h.cfg.Graduate.APIUrl))
	fetchCtx, cancel := script.WithTimeout(c.Request.Context(), 300*time.Second)
	allItems, err := script.RunGraduateFetchWithEnv(fetchCtx, "internal/script/graduate_fetch.py", env)
	cancel()
	if err != nil {
		h.log.Error("Graduate fetch failed", logging.Err(err))
		sendSSE("error", gin.H{"error": err.Error()})
		return
	}

	total := len(allItems)
	h.log.Info("Graduate fetched items", logging.KV("items_total", total))

	// 构建 ApplyID 列表和映射
	applyIDs := make([]string, 0, total)
	applyIDToName := make(map[string]string)
	applyIDToItem := make(map[string]script.WintalentItem)
	for _, item := range allItems {
		if item.ApplyID != "" {
			applyIDs = append(applyIDs, item.ApplyID)
			applyIDToName[item.ApplyID] = item.Name
			applyIDToItem[item.ApplyID] = item
		}
	}

	// 检测重复候选人
	sendSSE("status", gin.H{"stage": "checking", "message": "正在检测已评估的候选人..."})
	dupResult, _ := h.svc.CheckDuplicates(applyIDs, applyIDToName, userID)

	// 获取前端传递的跳过标志
	skipDuplicates := c.PostForm("skip_duplicates") == "true"

	// 如果有重复且未强制重新评估且用户未确认跳过
	if dupResult.DuplicateCount > 0 && !forceReevaluate && !skipDuplicates {
		duplicatesInfo := make([]gin.H, 0, len(dupResult.Duplicates))
		for _, dup := range dupResult.Duplicates {
			duplicatesInfo = append(duplicatesInfo, gin.H{
				"name":         dup.Name,
				"apply_id":     dup.ApplyID,
				"evaluated_at": dup.EvaluatedAt.Format("2006-01-02 15:04"),
				"total_score":  dup.TotalScore,
				"grade":        dup.Grade,
			})
		}

		sendSSE("duplicates_found", gin.H{
			"duplicates":      duplicatesInfo,
			"total_fetched":   dupResult.TotalFetched,
			"duplicate_count": dupResult.DuplicateCount,
			"new_count":       len(dupResult.NewApplyIDs),
		})
		return
	}

	// 确定要评估的候选人列表
	var itemsToEvaluate []script.WintalentItem
	var skipApplyIDs map[string]bool

	if forceReevaluate {
		itemsToEvaluate = allItems
		skipApplyIDs = make(map[string]bool)
	} else {
		skipApplyIDs = make(map[string]bool)
		for _, dup := range dupResult.Duplicates {
			skipApplyIDs[dup.ApplyID] = true
		}
		itemsToEvaluate = make([]script.WintalentItem, 0, len(dupResult.NewApplyIDs))
		for _, item := range allItems {
			if !skipApplyIDs[item.ApplyID] {
				itemsToEvaluate = append(itemsToEvaluate, item)
			}
		}
	}

	if len(itemsToEvaluate) == 0 {
		sendSSE("complete", gin.H{
			"success":       true,
			"total":         0,
			"results":       []interface{}{},
			"skipped_count": dupResult.DuplicateCount,
			"message":       "所有候选人均已评估，无需重复操作",
		})
		return
	}

	evalTotal := len(itemsToEvaluate)
	h.log.Info("Starting evaluation", logging.KV("eval_total", evalTotal), logging.KV("skipped", dupResult.DuplicateCount))

	sendSSE("total", gin.H{"total": evalTotal, "skipped": dupResult.DuplicateCount})
	sendSSE("status", gin.H{"stage": "evaluating", "message": fmt.Sprintf("开始评估 %d 份简历...", evalTotal)})

	type res struct {
		Filename        string               `json:"filename"`
		CandidateID     uint                 `json:"candidate_id"`
		CandidateName   string               `json:"candidate_name"`
		ApplyID         string               `json:"apply_id"`
		TotalScore      float64              `json:"total_score"`
		Grade           string               `json:"grade"`
		JDMatch         int                  `json:"jd_match"`
		AgeScore        int                  `json:"age_score"`
		ExperienceScore int                  `json:"experience_score"`
		EducationScore  int                  `json:"education_score"`
		CompanyScore    int                  `json:"company_score"`
		TechScore       int                  `json:"tech_score"`
		ProjectScore    int                  `json:"project_score"`
		Recommendation  string               `json:"recommendation"`
		ReportHTML      string               `json:"report_html"`
		ReportMarkdown  string               `json:"report_markdown"`
		CozeReportJSON  string               `json:"coze_report_json"`
		SavedPath       string               `json:"saved_path"`
		Rank            int                  `json:"rank"`
		ResultType      EvaluationResultType `json:"result_type"`
		Error           string               `json:"error,omitempty"`
	}

	concurrency := h.cfg.Batch.Concurrency
	if concurrency <= 0 {
		concurrency = 5
	}
	sem := make(chan struct{}, concurrency)
	var g errgroup.Group
	results := make([]res, evalTotal)

	var completed int32

	for i := range itemsToEvaluate {
		i := i
		sem <- struct{}{}
		g.Go(func() error {
			defer func() { <-sem }()
			item := itemsToEvaluate[i]

			// 构建 JD 文本
			jdText := jdFromReq
			if strings.TrimSpace(jdText) == "" && item.JD != nil {
				// 从毕业设计后台的职位信息构建 JD
				if title, ok := item.JD["title"].(string); ok && title != "" {
					jdText = "职位：" + title + "\n"
				}
				if desc, ok := item.JD["description"].(string); ok && desc != "" {
					jdText += "描述：" + desc + "\n"
				}
				if req, ok := item.JD["requirements"].(string); ok && req != "" {
					jdText += "要求：" + req + "\n"
				}
			}

			// 解码 PDF
			var pdfBytes []byte
			if item.ResumePDFB64 != "" {
				b, err := base64.StdEncoding.DecodeString(item.ResumePDFB64)
				if err == nil {
					pdfBytes = b
				}
			}
			name := strings.TrimSpace(item.Name)
			if name == "" {
				name = "unknown"
			}
			filename := name + ".pdf"

			// 检查是否是重新评估
			var existingCandidate *models.Candidate
			var resultType EvaluationResultType = ResultTypeNew
			if forceReevaluate && item.ApplyID != "" {
				existingCandidate, _ = h.svc.GetExistingCandidateByApplyID(item.ApplyID, userID)
				if existingCandidate != nil {
					resultType = ResultTypeReevaluated
				}
			}

			// 调用 Coze 工作流
			cozeCtx, cancel := context.WithTimeout(c.Request.Context(), 300*time.Second)
			cozeData, cozeErr := coze.RunWorkflow(cozeCtx, name, jdText, pdfBytes)
			cancel()
			if cozeErr != nil {
				h.log.Error("Coze workflow failed", logging.Err(cozeErr), logging.KV("candidate", name))
				cozeData = nil
			}

			// 评估
			out, evalErr := h.svc.EvaluateSingleBytesWithApplyID(pdfBytes, filename, jdText, criteria, cozeData, userID, item.ApplyID, existingCandidate)
			if evalErr != nil {
				results[i] = res{
					Filename:      filename,
					CandidateName: utils.ExtractCandidateName(filename),
					ApplyID:       item.ApplyID,
					ResultType:    resultType,
					Error:         firstErrMsg(cozeErr, evalErr),
				}
			} else {
				cand := out.Candidate
				results[i] = res{
					Filename:        filename,
					CandidateID:     cand.ID,
					CandidateName:   cand.Name,
					ApplyID:         item.ApplyID,
					TotalScore:      cand.TotalScore,
					Grade:           cand.Grade,
					JDMatch:         cand.JDMatch,
					AgeScore:        cand.AgeScore,
					ExperienceScore: cand.ExperienceScore,
					EducationScore:  cand.EducationScore,
					CompanyScore:    cand.CompanyScore,
					TechScore:       cand.TechScore,
					ProjectScore:    cand.ProjectScore,
					Recommendation:  cand.Recommendation,
					ReportHTML:      out.ReportHTML,
					ReportMarkdown:  out.ReportMD,
					CozeReportJSON:  cand.CozeReportJSON,
					SavedPath:       out.ReportMDPath,
					ResultType:      resultType,
				}
			}

			done := atomic.AddInt32(&completed, 1)
			sendSSE("progress", gin.H{
				"completed":   done,
				"total":       evalTotal,
				"percent":     int(float64(done) / float64(evalTotal) * 100),
				"current":     name,
				"result_type": resultType,
			})

			return nil
		})
	}
	_ = g.Wait()

	// 排序并设置排名
	sort.Slice(results, func(i, j int) bool { return results[i].TotalScore > results[j].TotalScore })
	for i := range results {
		results[i].Rank = i + 1
	}

	// 统计结果类型
	var newCount, reevaluatedCount int
	for _, r := range results {
		switch r.ResultType {
		case ResultTypeNew:
			newCount++
		case ResultTypeReevaluated:
			reevaluatedCount++
		}
	}

	// 收集成功评估的候选人用于自动推送
	var successCandidates []models.Candidate
	for _, r := range results {
		if r.Error == "" && r.CandidateID > 0 {
			successCandidates = append(successCandidates, models.Candidate{
				ID:             r.CandidateID,
				Name:           r.CandidateName,
				TotalScore:     r.TotalScore,
				Grade:          r.Grade,
				JDMatch:        r.JDMatch,
				Recommendation: r.Recommendation,
			})
		}
	}

	if len(successCandidates) > 0 {
		go h.triggerAutoPushIfEnabled(userID, successCandidates)
	}

	sendSSE("complete", gin.H{
		"success":           true,
		"total":             evalTotal,
		"results":           results,
		"new_count":         newCount,
		"reevaluated_count": reevaluatedCount,
		"skipped_count":     dupResult.DuplicateCount,
	})
}
