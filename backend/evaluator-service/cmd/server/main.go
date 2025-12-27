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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := logging.New()
	defer logger.Sync()

	if err := database.Init(cfg, logger); err != nil {
		logger.Fatal("db init failed", logging.Err(err))
	}

	// Initialize Coze client
	coze.Init(cfg)

	// Initialize Python path from config
	if cfg.Python.Path != "" {
		script.SetPythonPath(cfg.Python.Path)
	}

	// Initialize services
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := repository.NewCandidateRepository(database.DB)
	dtRepo := repository.NewDingTalkRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	dtService := service.NewDingTalkService(repo, dtRepo, logger)
	authSvc := service.NewAuthService(cfg, logger, userRepo, database.DB)

	// Start DingTalk service
	if err := dtService.Start(ctx); err != nil {
		logger.Error("dingtalk service start failed", logging.Err(err))
	}
	defer dtService.Stop()

	// Pass services to router/handlers
	r := api.NewRouter(cfg, logger, dtService, authSvc)

	// Setup graceful shutdown
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

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", logging.Err(err))
	}
	logger.Info("server stopped")
}
