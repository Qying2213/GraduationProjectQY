package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"evaluator-service/internal/api/middleware"
	"evaluator-service/internal/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) CandidatesStatsSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)
	total, err := h.repo.CountByUser(userID)
	if err != nil {
		fail(c, err)
		return
	}
	pending, _ := h.repo.CountByStatusAndUser("待面试", userID)
	interviewed, _ := h.repo.CountByStatusAndUser("已面试", userID)
	hired, _ := h.repo.CountByStatusAndUser("已录用", userID)
	rejected, _ := h.repo.CountByStatusAndUser("淘汰", userID)
	ok(c, gin.H{"total": total, "pending": pending, "interviewed": interviewed, "hired": hired, "rejected": rejected})
}

func (h *Handlers) GetCandidates(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var f repository.ListFilter
	f.Status = c.Query("status")
	f.Grade = c.Query("grade")
	if ms := c.Query("min_score"); ms != "" {
		if v, err := strconv.ParseFloat(ms, 64); err == nil {
			f.MinScore = &v
		}
	}
	f.Search = c.Query("search")
	list, err := h.repo.ListByUser(userID, f)
	if err != nil {
		fail(c, err)
		return
	}
	out := make([]gin.H, 0, len(list))
	for _, cnd := range list {
		out = append(out, gin.H{
			"id":               cnd.ID,
			"name":             cnd.Name,
			"filename":         cnd.Filename,
			"total_score":      cnd.TotalScore,
			"grade":            cnd.Grade,
			"jd_match":         cnd.JDMatch,
			"age_score":        cnd.AgeScore,
			"experience_score": cnd.ExperienceScore,
			"education_score":  cnd.EducationScore,
			"company_score":    cnd.CompanyScore,
			"tech_score":       cnd.TechScore,
			"project_score":    cnd.ProjectScore,
			"recommendation":   cnd.Recommendation,
			"status":           cnd.Status,
			"created_at":       cnd.CreatedAt.Format("2006-01-02 15:04"),
			"notes":            cnd.Notes,
		})
	}
	ok(c, out)
}

func (h *Handlers) GetCandidate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	cand, err := h.repo.GetByUser(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ok(c, gin.H{
		"id":                cand.ID,
		"name":              cand.Name,
		"filename":          cand.Filename,
		"pdf_path":          cand.PDFPath,
		"has_pdf":           cand.PDFPath != "" && fileExists(cand.PDFPath),
		"total_score":       cand.TotalScore,
		"grade":             cand.Grade,
		"jd_match":          cand.JDMatch,
		"age_score":         cand.AgeScore,
		"experience_score":  cand.ExperienceScore,
		"education_score":   cand.EducationScore,
		"company_score":     cand.CompanyScore,
		"tech_score":        cand.TechScore,
		"project_score":     cand.ProjectScore,
		"age_reason":        cand.AgeReason,
		"experience_reason": cand.ExperienceReason,
		"education_reason":  cand.EducationReason,
		"company_reason":    cand.CompanyReason,
		"tech_reason":       cand.TechReason,
		"project_reason":    cand.ProjectReason,
		"recommendation":    cand.Recommendation,
		"report_markdown":   cand.ReportMarkdown,
		"resume_markdown":   cand.ResumeMarkdown,
		"coze_report_json":  cand.CozeReportJSON,
		"status":            cand.Status,
		"created_at":        cand.CreatedAt.Format("2006-01-02 15:04"),
		"notes":             cand.Notes,
	})
}

func (h *Handlers) GetCandidateResume(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	cand, err := h.repo.GetByUser(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	if cand.PDFPath == "" || !fileExists(cand.PDFPath) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "简历文件不存在"})
		return
	}
	name := cand.Filename
	if name == "" {
		name = cand.Name + "_简历.pdf"
	}
	c.FileAttachment(cand.PDFPath, filepath.Base(name))
}

func (h *Handlers) UpdateCandidateStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	status := c.PostForm("status")
	cand, err := h.repo.GetByUser(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	cand.Status = status
	if err := h.repo.Update(cand); err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"success": true, "status": status})
}

func (h *Handlers) UpdateCandidateNotes(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	notes := c.PostForm("notes")
	cand, err := h.repo.GetByUser(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	cand.Notes = notes
	if err := h.repo.Update(cand); err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"success": true})
}

func (h *Handlers) DeleteCandidate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	cand, err := h.repo.GetByUser(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	if cand.PDFPath != "" {
		_ = os.Remove(cand.PDFPath)
	}
	if err := h.repo.DeleteByUser(uint(id), userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ok(c, gin.H{"success": true})
}

func (h *Handlers) DeleteAllCandidates(c *gin.Context) {
	userID := middleware.GetUserID(c)
	list, _ := h.repo.AllByUser(userID)
	for _, x := range list {
		if x.PDFPath != "" {
			_ = os.Remove(x.PDFPath)
		}
	}
	n, err := h.repo.DeleteAllByUser(userID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"success": true, "deleted": n})
}

func (h *Handlers) CandidatesStatsCharts(c *gin.Context) {
	userID := middleware.GetUserID(c)
	list, err := h.repo.AllByUser(userID)
	if err != nil {
		fail(c, err)
		return
	}
	gradeDist := map[string]int{"A": 0, "B": 0, "C": 0, "D": 0, "E": 0}
	statusDist := map[string]int{"待面试": 0, "已面试": 0, "已录用": 0, "淘汰": 0}
	recommendDist := map[string]int{"推荐": 0, "待定": 0, "不推荐": 0}
	scoreDist := map[string]int{"0-20": 0, "21-40": 0, "41-60": 0, "61-80": 0, "81-100": 0}
	dims := map[string]float64{"age": 0, "experience": 0, "education": 0, "company": 0, "tech": 0, "project": 0}
	dateStats := map[string]int{}
	total := len(list)
	for _, cnd := range list {
		if _, ok := gradeDist[cnd.Grade]; ok {
			gradeDist[cnd.Grade]++
		}
		if _, ok := statusDist[cnd.Status]; ok {
			statusDist[cnd.Status]++
		}
		rec := cnd.Recommendation
		if strings.Contains(rec, "推荐录用") || strings.Contains(rec, "建议录用") {
			recommendDist["推荐"]++
		} else if strings.Contains(rec, "谨慎考虑") || strings.Contains(rec, "待定") {
			recommendDist["待定"]++
		} else {
			recommendDist["不推荐"]++
		}
		s := cnd.TotalScore
		switch {
		case s <= 20:
			scoreDist["0-20"]++
		case s <= 40:
			scoreDist["21-40"]++
		case s <= 60:
			scoreDist["41-60"]++
		case s <= 80:
			scoreDist["61-80"]++
		default:
			scoreDist["81-100"]++
		}
		dims["age"] += float64(cnd.AgeScore)
		dims["experience"] += float64(cnd.ExperienceScore)
		dims["education"] += float64(cnd.EducationScore)
		dims["company"] += float64(cnd.CompanyScore)
		dims["tech"] += float64(cnd.TechScore)
		dims["project"] += float64(cnd.ProjectScore)
		key := cnd.CreatedAt.Format("2006-01-02")
		dateStats[key] = dateStats[key] + 1
	}
	if total > 0 {
		for k, v := range dims {
			dims[k] = round1(v / float64(total))
		}
	}
	keys := make([]string, 0, len(dateStats))
	for k := range dateStats {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(keys) > 7 {
		keys = keys[len(keys)-7:]
	}
	ordered := gin.H{}
	for _, k := range keys {
		ordered[k] = dateStats[k]
	}
	ok(c, gin.H{
		"total":                  total,
		"grade_distribution":     gradeDist,
		"status_distribution":    statusDist,
		"recommend_distribution": recommendDist,
		"score_distribution":     scoreDist,
		"dimension_average":      dims,
		"date_stats":             ordered,
	})
}

func (h *Handlers) CompareCandidates(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idsStr := c.PostForm("ids")
	parts := strings.Split(idsStr, ",")
	var ids []uint
	for _, p := range parts {
		if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			ids = append(ids, uint(n))
		}
	}
	list, err := h.repo.GetInByUser(ids, userID)
	if err != nil {
		fail(c, err)
		return
	}
	sort.Slice(list, func(i, j int) bool { return list[i].TotalScore > list[j].TotalScore })
	out := make([]gin.H, 0, len(list))
	for _, cnd := range list {
		out = append(out, gin.H{
			"id":                cnd.ID,
			"name":              cnd.Name,
			"total_score":       cnd.TotalScore,
			"grade":             cnd.Grade,
			"jd_match":          cnd.JDMatch,
			"age_score":         cnd.AgeScore,
			"experience_score":  cnd.ExperienceScore,
			"education_score":   cnd.EducationScore,
			"company_score":     cnd.CompanyScore,
			"tech_score":        cnd.TechScore,
			"project_score":     cnd.ProjectScore,
			"age_reason":        cnd.AgeReason,
			"experience_reason": cnd.ExperienceReason,
			"education_reason":  cnd.EducationReason,
			"company_reason":    cnd.CompanyReason,
			"tech_reason":       cnd.TechReason,
			"project_reason":    cnd.ProjectReason,
			"recommendation":    cnd.Recommendation,
			"status":            cnd.Status,
		})
	}
	ok(c, gin.H{"success": true, "candidates": out})
}

func (h *Handlers) ExportCompareReport(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idsStr := c.PostForm("ids")
	parts := strings.Split(idsStr, ",")
	var ids []uint
	for _, p := range parts {
		if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			ids = append(ids, uint(n))
		}
	}
	list, err := h.repo.GetInByUser(ids, userID)
	if err != nil {
		fail(c, err)
		return
	}
	sort.Slice(list, func(i, j int) bool { return list[i].TotalScore > list[j].TotalScore })
	b, err := h.exprt.ExcelCompare(list)
	if err != nil {
		fail(c, err)
		return
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=candidate_compare.xlsx")
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(b)
}

func fileExists(p string) bool { _, err := os.Stat(p); return err == nil }

func round1(f float64) float64 { return float64(int(f*10+0.5)) / 10 }
