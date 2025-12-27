package api

import (
	"net/http"

	"evaluator-service/internal/api/handlers"
	"evaluator-service/internal/api/middleware"
	"evaluator-service/internal/config"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, log *logging.Logger, dtService *service.DingTalkService, authSvc *service.AuthService) http.Handler {
	g := gin.New()
	g.Use(middleware.RequestID())
	g.Use(middleware.Recovery(log))
	g.Use(middleware.CORS())
	g.Use(middleware.Logger(log))
	g.MaxMultipartMemory = int64(cfg.Batch.MaxUploadMB) * 1024 * 1024

	h := handlers.New(cfg, log, dtService, authSvc)

	// 公开路由（无需认证）
	g.GET("/", h.Index)
	g.Static("/static", cfg.Storage.StaticDir)

	// 认证路由（无需登录）
	auth := g.Group("/api/auth")
	{
		auth.POST("/login", h.Login)
	}

	// 受保护的 API 路由（需要认证）
	api := g.Group("/api")
	api.Use(middleware.Auth(authSvc))
	{
		// 认证相关
		api.GET("/auth/profile", h.GetProfile)
		api.POST("/auth/logout", h.Logout)

		tmpl := api.Group("/templates")
		tmpl.GET("/jd", h.GetJDTemplate)

		api.POST("/evaluate", h.Evaluate)
		api.POST("/evaluate/batch", h.EvaluateBatch)
		api.POST("/evaluate/batch/graduate", h.EvaluateBatchGraduate) // 从毕业设计后台获取简历评估

		// Position endpoints
		positions := api.Group("/positions")
		positions.GET("", h.GetPositions)
		positions.POST("/sync", h.SyncPositions)

		// Credentials endpoints
		creds := api.Group("/credentials")
		creds.GET("/status", h.CredentialStatus)
		creds.POST("", h.CredentialUpsert)

		export := api.Group("/export")
		export.POST("/pdf", h.ExportPDF)
		export.POST("/excel", h.ExportExcel)

		cands := api.Group("/candidates")
		cands.GET("/stats/summary", h.CandidatesStatsSummary)
		cands.GET("/stats/charts", h.CandidatesStatsCharts)
		cands.GET("", h.GetCandidates)
		cands.GET("/:id", h.GetCandidate)
		cands.GET("/:id/resume", h.GetCandidateResume)
		cands.PUT("/:id/status", h.UpdateCandidateStatus)
		cands.PUT("/:id/notes", h.UpdateCandidateNotes)
		cands.DELETE(":id", h.DeleteCandidate)
		cands.DELETE("", h.DeleteAllCandidates)
		cands.POST("/compare", h.CompareCandidates)
		cands.POST("/compare/export", h.ExportCompareReport)

		// DingTalk endpoints
		dt := api.Group("/dingtalk")
		dt.GET("/config", h.GetDingTalkConfig)            // 获取默认配置（兼容）
		dt.GET("/configs", h.ListDingTalkConfigs)         // 获取所有配置
		dt.GET("/configs/:id", h.GetDingTalkConfigByID)   // 获取指定配置
		dt.POST("/config", h.UpsertDingTalkConfig)        // 创建/更新配置
		dt.DELETE("/configs/:id", h.DeleteDingTalkConfig) // 删除配置
		dt.POST("/test", h.TestDingTalkPush)              // 测试推送
		dt.POST("/push", h.PushNow)                       // 立即推送候选人
	}

	return g
}
