// Package rag 实现平台中的检索增强推荐能力。
// vector_store.go 负责向量存储和检索：把人才、职位、简历索引为 Embedding，
// 再为 AI 评估和推荐解释返回语义最相似的记录。
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

// VectorStore 是项目内的向量仓库。
//
// 早期设计中提到 pgvector，但当前实现把向量存为 JSONB。这样普通 PostgreSQL
// 环境也能直接部署，同时仍然可以在 Go 代码里计算余弦相似度。
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

// TalentEmbedding 存储单个人才的可检索语义画像。
// content 和向量一起保存，这样 RAG 命中后能返回可读上下文，而不只是一个 ID。
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

// IndexTalent 为人才画像创建或刷新向量索引。
// 使用 upsert 保证索引操作幂等，因此重建 RAG 索引是安全的。
func (vs *VectorStore) IndexTalent(ctx context.Context, talentID uint, content string) error {
	// 获取向量
	emb, err := vs.embeddingClient.GetEmbedding(ctx, content)
	if err != nil {
		return fmt.Errorf("获取向量失败: %w", err)
	}

	// 转换为 pgvector 格式
	embStr := vectorToString(emb)

	// 使用原生 SQL upsert；embedding 使用 JSON 文本，兼容未安装 pgvector 的环境。
	sql := `INSERT INTO talent_embeddings (talent_id, content, embedding) 
			VALUES (?, ?, ?) 
			ON CONFLICT (talent_id) DO UPDATE 
			SET content = EXCLUDED.content, embedding = EXCLUDED.embedding, updated_at = CURRENT_TIMESTAMP`
	return vs.db.WithContext(ctx).Exec(sql, talentID, content, embStr).Error
}

// IndexJob 为职位描述创建向量索引，用于相似职位推荐和人岗匹配。
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

// IndexResume 为已评估简历创建向量索引，让后续评估可以把历史简历作为 RAG 示例参考。
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

// SearchSimilarTalents 在人才向量索引中做语义检索。
// query 可以是 JD、简历文本，也可以是 HR 输入的自然语言问题。
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

// SearchSimilarJobs 在职位向量索引中做语义检索。
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

// SearchSimilarResumes 在简历向量索引中做语义检索。
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

// vectorToString 将向量序列化为 JSONB 可存储的字符串。
// 单独封装这一层，后续如果迁移到原生 pgvector 会更容易替换。
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
