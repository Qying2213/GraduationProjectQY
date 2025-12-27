package handlers

import (
	"net/http"
	"path/filepath"

	"evaluator-service/internal/config"
	"evaluator-service/internal/database"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/repository"
	"evaluator-service/internal/service"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	cfg       *config.Config
	log       *logging.Logger
	repo      *repository.CandidateRepository
	dtRepo    *repository.DingTalkRepository
	userRepo  *repository.UserRepository
	credRepo  *repository.CredentialRepository
	svc       *service.ResumeService
	exprt     *service.ExportService
	dtService *service.DingTalkService
	authSvc   *service.AuthService
}

func New(cfg *config.Config, log *logging.Logger, dtService *service.DingTalkService, authSvc *service.AuthService) *Handlers {
	repo := repository.NewCandidateRepository(database.DB)
	dtRepo := repository.NewDingTalkRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	credRepo := repository.NewCredentialRepository(database.DB)
	svc := service.NewResumeService(cfg, log, repo)
	ex := service.NewExportService(cfg, log)
	return &Handlers{
		cfg:       cfg,
		log:       log,
		repo:      repo,
		dtRepo:    dtRepo,
		userRepo:  userRepo,
		credRepo:  credRepo,
		svc:       svc,
		exprt:     ex,
		dtService: dtService,
		authSvc:   authSvc,
	}
}

func (h *Handlers) Index(c *gin.Context) {
	index := filepath.Join(h.cfg.Storage.StaticDir, "index.html")
	c.File(index)
}

func ok(c *gin.Context, v any)      { c.JSON(http.StatusOK, v) }
func bad(c *gin.Context, err error) { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) }
func fail(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
