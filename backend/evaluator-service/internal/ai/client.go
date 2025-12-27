package ai

import (
	"evaluator-service/internal/config"
	"evaluator-service/internal/models"
)

type Client interface {
	OCR(images [][]byte) (string, error)
	Structure(text string) (string, error)
	EvaluateJD(resumeMD, jd string) (models.JDMatchResult, error)
	EvaluateRequirement(resumeMD string) (models.RequirementResult, error)
	Score(resumeMD, jd, criteria string) (models.ScoringResult, error)
	GenerateInterviewQuestions(resumeMD string, eval models.EvaluationResult) ([]models.InterviewQuestion, error)
}

type Factory struct{ cfg *config.Config }

func NewFactory(cfg *config.Config) *Factory { return &Factory{cfg: cfg} }

// New 方法已不再使用，保留用于兼容性
// 实际使用 NewCozeClient 直接创建客户端
func (f *Factory) New() (Client, error) {
	// 这个方法不再被使用，所有评估都通过 CozeClient 进行
	// 保留此方法以避免破坏现有接口
	return nil, nil
}
