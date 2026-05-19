package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"evaluator-service/internal/api"
	"evaluator-service/internal/config"
	"evaluator-service/internal/database"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/repository"
	"evaluator-service/internal/script"
	"evaluator-service/internal/service"
	"evaluator-service/internal/thirdparty/coze"
)

func main() {
	// evaluator-service 是独立的评估后台服务，主要服务于外部招聘系统抓取、
	// 批量评估、钉钉推送和报告导出。它不是当前毕业设计主链路的入口；
	// 毕设主 AI 评估入口在 resume-service 中。
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := logging.New()
	defer logger.Sync()

	if err := database.Init(cfg, logger); err != nil {
		logger.Fatal("db init failed", logging.Err(err))
	}

	// 初始化 Coze 客户端，供独立评估服务调用外部工作流。
	coze.Init(cfg)

	// 从配置初始化 Python 路径，批量抓取外部招聘系统简历时会用到脚本执行器。
	if cfg.Python.Path != "" {
		script.SetPythonPath(cfg.Python.Path)
	}

	// 初始化仓储和服务。这里的 DingTalkService/AuthService 只属于 evaluator-service 自己的后台。
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := repository.NewCandidateRepository(database.DB)
	dtRepo := repository.NewDingTalkRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	dtService := service.NewDingTalkService(repo, dtRepo, logger)
	authSvc := service.NewAuthService(cfg, logger, userRepo, database.DB)

	// 启动钉钉推送服务，用于把评估后的候选人结果同步到外部通知渠道。
	if err := dtService.Start(ctx); err != nil {
		logger.Error("dingtalk service start failed", logging.Err(err))
	}
	defer dtService.Stop()

	// 把服务注入路由和 handler，形成独立 evaluator 后台 API。
	r := api.NewRouter(cfg, logger, dtService, authSvc)

	// 设置优雅关闭，收到中断信号后先停止 HTTP 服务和后台任务。
	srv := &http.Server{
		Addr:    cfg.Server.Address(),
		Handler: r,
	}

	go func() {
		logger.Info("server start", logging.KV("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", logging.Err(err))
		}
	}()

	// 等待中断信号。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", logging.Err(err))
	}
	logger.Info("server stopped")
}
