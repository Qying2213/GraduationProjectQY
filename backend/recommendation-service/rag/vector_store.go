package rag

import (
	"context"
	"fmt"
	"recommendation-service/embedding"
	"strings"

	"gorm.io/gorm"
)

// VectorStore pgvector向量存储
type VectorStore struct {
	db              *gorm.DB
	embeddingClient *embedding.Client
}

// NewVectorStore 创建向量存储
func NewVectorStore(db *gorm.DB) *VectorStore {
	return &VectorStore{
		db:              db,
		embeddingClient: embedding.GetClient(),
	}
}

// TalentEmbedding 人才向量记录
type TalentEmbedding struct {
	ID        uint   `gorm:"primaryKey"`
	TalentID  uint   `gorm:"uniqueIndex"`
	Content   string `gorm:"type:text"`
	Embedding string `gorm:"type:vector(1024)"`
}

func (TalentEmbedding) TableName() string { return "talent_embeddings" }

// JobEmbedding 职位向量记录
type JobEmbedding struct {
	ID        uint   `gorm:"primaryKey"`
	JobID     uint   `gorm:"uniqueIndex"`
	Content   string `gorm:"type:text"`
	Embedding string `gorm:"type:vector(1024)"`
}

func (JobEmbedding) TableName() string { return "job_embeddings" }

// IndexTalent 索引人才信息
func (vs *VectorStore) IndexTalent(ctx context.Context, talentID uint, content string) error {
	// 获取向量
	emb, err := vs.embeddingClient.GetEmbedding(ctx, content)
	if err != nil {
		return fmt.Errorf("获取向量失败: %w", err)
	}

	// 转换为 pgvector 格式
	embStr := vectorToString(emb)

	// 使用原生SQL upsert
	sql := `INSERT INTO talent_embeddings (talent_id, content, embedding) 
			VALUES (?, ?, ?::vector) 
			ON CONFLICT (talent_id) DO UPDATE SET content = ?, embedding = ?::vector`
	return vs.db.Exec(sql, talentID, content, embStr, content, embStr).Error
}

// IndexJob 索引职位信息
func (vs *VectorStore) IndexJob(ctx context.Context, jobID uint, content string) error {
	emb, err := vs.embeddingClient.GetEmbedding(ctx, content)
	if err != nil {
		return fmt.Errorf("获取向量失败: %w", err)
	}

	embStr := vectorToString(emb)

	sql := `INSERT INTO job_embeddings (job_id, content, embedding) 
			VALUES (?, ?, ?::vector) 
			ON CONFLICT (job_id) DO UPDATE SET content = ?, embedding = ?::vector`
	return vs.db.Exec(sql, jobID, content, embStr, content, embStr).Error
}

// SearchResult 搜索结果
type SearchResult struct {
	ID         uint    `json:"id"`
	Content    string  `json:"content"`
	Similarity float64 `json:"similarity"`
}

// SearchSimilarTalents 搜索相似人才
func (vs *VectorStore) SearchSimilarTalents(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	emb, err := vs.embeddingClient.GetEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	embStr := vectorToString(emb)

	var results []SearchResult
	sql := `SELECT talent_id as id, content, 1 - (embedding <=> ?::vector) as similarity 
			FROM talent_embeddings 
			ORDER BY embedding <=> ?::vector 
			LIMIT ?`
	if err := vs.db.Raw(sql, embStr, embStr, limit).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// SearchSimilarJobs 搜索相似职位
func (vs *VectorStore) SearchSimilarJobs(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	emb, err := vs.embeddingClient.GetEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	embStr := vectorToString(emb)

	var results []SearchResult
	sql := `SELECT job_id as id, content, 1 - (embedding <=> ?::vector) as similarity 
			FROM job_embeddings 
			ORDER BY embedding <=> ?::vector 
			LIMIT ?`
	if err := vs.db.Raw(sql, embStr, embStr, limit).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// vectorToString 将向量转换为 pgvector 格式字符串
func vectorToString(vec []float64) string {
	strs := make([]string, len(vec))
	for i, v := range vec {
		strs[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(strs, ",") + "]"
}
