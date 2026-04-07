package handlers

import (
	"resume-service/evaluator"
	"resume-service/models"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type aiSyncTalent struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Name            string
	Email           string
	Phone           string
	Skills          pq.StringArray `gorm:"type:text[]"`
	Experience      int
	Education       string
	Status          string
	Location        string
	Salary          string
	Summary         string
	CurrentCompany  string
	CurrentPosition string
	Source          string
	ResumeID        *uint
}

func (aiSyncTalent) TableName() string {
	return "talents"
}

func TestEnsureTalentForResumeCreatesTalentRecord(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&aiSyncTalent{}))

	handler := &AIEvaluateHandler{DB: db}

	resume := models.Resume{
		ID:       30,
		FileName: "candidate.pdf",
	}
	result := &evaluator.EvaluationResult{
		Name:          "张三",
		Summary:       "5年后端开发经验，熟悉Go和数据库",
		MatchedSkills: []string{"Go", "PostgreSQL"},
		ParsedReport: map[string]interface{}{
			"基本信息": map[string]interface{}{
				"姓名": "张三",
				"邮箱": "zhangsan@example.com",
				"手机": "13800000000",
				"学历": "本科",
				"城市": "上海",
			},
		},
	}

	talentID := handler.ensureTalentForResume(&resume, result)
	require.NotNil(t, talentID)

	var talent aiSyncTalent
	require.NoError(t, db.First(&talent, *talentID).Error)
	require.Equal(t, "张三", talent.Name)
	require.Equal(t, "zhangsan@example.com", talent.Email)
	require.Equal(t, "active", talent.Status)
	require.Equal(t, "AI评估导入", talent.Source)
	require.NotNil(t, talent.ResumeID)
	require.Equal(t, resume.ID, *talent.ResumeID)
}

func TestSaveEvaluationResultFallsBackToOCRContactInfo(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.Resume{}, &models.EvaluationResult{}))

	handler := &AIEvaluateHandler{DB: db}

	resume := models.Resume{
		FileName:      "candidate.pdf",
		FilePath:      "/tmp/candidate.pdf",
		ExtractedText: "张三\n手机号：13800000000\n邮箱：zhangsan@example.com\n本科\n上海\n5年后端开发经验",
		Status:        "parsed",
	}
	require.NoError(t, db.Create(&resume).Error)

	result := &evaluator.EvaluationResult{
		Name: "张三",
		ParsedReport: map[string]interface{}{
			"基本信息": map[string]interface{}{
				"姓名": "张三",
				"学历": "本科",
			},
		},
	}

	evalResult := handler.saveEvaluationResult(&resume, result, "张三", "ai_evaluate")
	require.Equal(t, "13800000000", evalResult.ParsedPhone)
	require.Equal(t, "zhangsan@example.com", evalResult.ParsedEmail)
	require.Equal(t, "本科", evalResult.ParsedEducation)
	require.Equal(t, "上海", evalResult.ParsedLocation)
}
