package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"recommendation-service/embedding"
	"sort"
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
	Embedding string `gorm:"type:jsonb"`
}

func (TalentEmbedding) TableName() string { return "talent_embeddings" }

// JobEmbedding 职位向量记录
type JobEmbedding struct {
	ID        uint   `gorm:"primaryKey"`
	JobID     uint   `gorm:"uniqueIndex"`
	Content   string `gorm:"type:text"`
	Embedding string `gorm:"type:jsonb"`
}

func (JobEmbedding) TableName() string { return "job_embeddings" }

// ResumeEmbedding 简历向量记录
type ResumeEmbedding struct {
	ID        uint   `gorm:"primaryKey"`
	ResumeID  uint   `gorm:"uniqueIndex"`
	Content   string `gorm:"type:text"`
	Embedding string `gorm:"type:jsonb"`
}

func (ResumeEmbedding) TableName() string { return "resume_embeddings" }

// IndexTalent 索引人才信息
func (vs *VectorStore) IndexTalent(ctx context.Context, talentID uint, content string) error {
	// 获取向量
	emb, err := vs.embeddingClient.GetEmbedding(ctx, content)
	if err != nil {
		return fmt.Errorf("获取向量失败: %w", err)
	}

	// 转换为 pgvector 格式
	embStr := vectorToString(emb)

	// 使用原生SQL upsert；embedding 使用 JSON 文本，兼容未安装 pgvector 的环境
	sql := `INSERT INTO talent_embeddings (talent_id, content, embedding) 
			VALUES (?, ?, ?) 
			ON CONFLICT (talent_id) DO UPDATE 
			SET content = EXCLUDED.content, embedding = EXCLUDED.embedding, updated_at = CURRENT_TIMESTAMP`
	return vs.db.WithContext(ctx).Exec(sql, talentID, content, embStr).Error
}

// IndexJob 索引职位信息
func (vs *VectorStore) IndexJob(ctx context.Context, jobID uint, content string) error {
	emb, err := vs.embeddingClient.GetEmbedding(ctx, content)
	if err != nil {
		return fmt.Errorf("获取向量失败: %w", err)
	}

	embStr := vectorToString(emb)

	sql := `INSERT INTO job_embeddings (job_id, content, embedding) 
			VALUES (?, ?, ?) 
			ON CONFLICT (job_id) DO UPDATE 
			SET content = EXCLUDED.content, embedding = EXCLUDED.embedding, updated_at = CURRENT_TIMESTAMP`
	return vs.db.WithContext(ctx).Exec(sql, jobID, content, embStr).Error
}

// IndexResume 索引简历信息
func (vs *VectorStore) IndexResume(ctx context.Context, resumeID uint, content string) error {
	emb, err := vs.embeddingClient.GetEmbedding(ctx, content)
	if err != nil {
		return fmt.Errorf("获取向量失败: %w", err)
	}

	embStr := vectorToString(emb)

	sql := `INSERT INTO resume_embeddings (resume_id, content, embedding)
			VALUES (?, ?, ?)
			ON CONFLICT (resume_id) DO UPDATE
			SET content = EXCLUDED.content, embedding = EXCLUDED.embedding, updated_at = CURRENT_TIMESTAMP`
	return vs.db.WithContext(ctx).Exec(sql, resumeID, content, embStr).Error
}

// SearchResult 搜索结果
type SearchResult struct {
	ID         uint    `json:"id"`
	SourceType string  `json:"source_type"`
	Content    string  `json:"content"`
	Similarity float64 `json:"similarity"`
}

// SearchSimilarTalents 搜索相似人才
func (vs *VectorStore) SearchSimilarTalents(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	queryEmbedding, err := vs.embeddingClient.GetEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	var rows []TalentEmbedding
	if err := vs.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(rows))
	for _, row := range rows {
		rowEmbedding, err := parseStoredEmbedding(row.Embedding)
		if err != nil {
			continue
		}

		results = append(results, SearchResult{
			ID:         row.TalentID,
			SourceType: "talent",
			Content:    row.Content,
			Similarity: embedding.CosineSimilarity(queryEmbedding, rowEmbedding),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// SearchSimilarJobs 搜索相似职位
func (vs *VectorStore) SearchSimilarJobs(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	queryEmbedding, err := vs.embeddingClient.GetEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	var rows []JobEmbedding
	if err := vs.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(rows))
	for _, row := range rows {
		rowEmbedding, err := parseStoredEmbedding(row.Embedding)
		if err != nil {
			continue
		}

		results = append(results, SearchResult{
			ID:         row.JobID,
			SourceType: "job",
			Content:    row.Content,
			Similarity: embedding.CosineSimilarity(queryEmbedding, rowEmbedding),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// SearchSimilarResumes 搜索相似简历
func (vs *VectorStore) SearchSimilarResumes(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	queryEmbedding, err := vs.embeddingClient.GetEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	var rows []ResumeEmbedding
	if err := vs.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(rows))
	for _, row := range rows {
		rowEmbedding, err := parseStoredEmbedding(row.Embedding)
		if err != nil {
			continue
		}

		results = append(results, SearchResult{
			ID:         row.ResumeID,
			SourceType: "resume",
			Content:    row.Content,
			Similarity: embedding.CosineSimilarity(queryEmbedding, rowEmbedding),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// vectorToString 将向量转换为 pgvector 格式字符串
func vectorToString(vec []float64) string {
	data, err := json.Marshal(vec)
	if err != nil {
		strs := make([]string, len(vec))
		for i, v := range vec {
			strs[i] = fmt.Sprintf("%f", v)
		}
		return "[" + strings.Join(strs, ",") + "]"
	}
	return string(data)
}

func parseStoredEmbedding(raw string) ([]float64, error) {
	if raw == "" {
		return nil, fmt.Errorf("empty embedding")
	}

	var vec []float64
	if err := json.Unmarshal([]byte(raw), &vec); err != nil {
		return nil, err
	}
	return vec, nil
}
