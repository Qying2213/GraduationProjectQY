package handlers

import (
	"net/http"
	"resume-service/riskcheck"

	"github.com/gin-gonic/gin"
)

// RiskCheckHandler 风控检查处理器
type RiskCheckHandler struct {
	checker *riskcheck.RiskChecker
}

// NewRiskCheckHandler 创建风控检查处理器
func NewRiskCheckHandler() *RiskCheckHandler {
	return &RiskCheckHandler{
		checker: riskcheck.NewRiskChecker(),
	}
}

// CheckResumeRisk 检查简历风险
// POST /api/v1/resumes/risk-check
func (h *RiskCheckHandler) CheckResumeRisk(c *gin.Context) {
	var req riskcheck.ResumeData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	result := h.checker.Check(&req)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    result,
	})
}

// CheckTimeConflict 仅检查时间冲突
// POST /api/v1/resumes/risk-check/time-conflict
func (h *RiskCheckHandler) CheckTimeConflict(c *gin.Context) {
	var req riskcheck.ResumeData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	risks := h.checker.CheckTimeConflict(&req)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"has_risk":   len(risks) > 0,
			"risk_items": risks,
		},
	})
}

// CheckEducationFraud 仅检查学历造假
// POST /api/v1/resumes/risk-check/education
func (h *RiskCheckHandler) CheckEducationFraud(c *gin.Context) {
	var req riskcheck.ResumeData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	risks := h.checker.CheckEducationFraud(&req)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"has_risk":   len(risks) > 0,
			"risk_items": risks,
		},
	})
}
