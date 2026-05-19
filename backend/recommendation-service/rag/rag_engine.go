// rag_engine.go 是 VectorStore 之上的编排层。
// 它把原始检索结果组织成 prompt 和匹配分析，让 RAG 能同时服务推荐页面和
// resume-service 的 AI 评估链路。
package rag

import (
	"context"
	"fmt"
	"recommendation-service/ai"
	"strings"

	"gorm.io/gorm"
)

// RAGEngine 组合向量检索和可选的 LLM 生成。
// VectorStore 回答“哪些记录相似”，Coze 负责把检索上下文转成可读解释。
type RAGEngine struct {
	vectorStore *VectorStore
	cozeClient  *ai.CozeClient
	db          *gorm.DB
}

// NewRAGEngine 创建RAG引擎
func NewRAGEngine(db *gorm.DB) *RAGEngine {
	return &RAGEngine{
		vectorStore: NewVectorStore(db),
		cozeClient:  ai.NewCozeClient(),
		db:          db,
	}
}

// GetVectorStore 获取向量存储
func (r *RAGEngine) GetVectorStore() *VectorStore {
	return r.vectorStore
}

// RAGRequest RAG请求
type RAGRequest struct {
	Query     string `json:"query"`      // 用户查询
	TopK      int    `json:"top_k"`      // 检索数量
	QueryType string `json:"query_type"` // talent 或 job
}

// RAGResponse RAG响应
type RAGResponse struct {
	Answer           string         `json:"answer"`            // 生成的回答
	RetrievedContext []SearchResult `json:"retrieved_context"` // 检索到的上下文
	Query            string         `json:"query"`             // 原始查询
}

// Query 对人才、职位或简历执行通用 RAG 查询。
// 它先检索上下文，再让 Coze 基于上下文回答；如果 Coze 不可用，也会降级返回
// 已检索到的记录。
func (r *RAGEngine) Query(ctx context.Context, req *RAGRequest) (*RAGResponse, error) {
	if req.TopK <= 0 {
		req.TopK = 5
	}

	// 1. 检索相似文档。QueryType 决定检索哪个向量表。
	var results []SearchResult
	var err error

	if req.QueryType == "job" {
		results, err = r.vectorStore.SearchSimilarJobs(ctx, req.Query, req.TopK)
	} else if req.QueryType == "resume" {
		results, err = r.vectorStore.SearchSimilarResumes(ctx, req.Query, req.TopK)
	} else {
		results, err = r.vectorStore.SearchSimilarTalents(ctx, req.Query, req.TopK)
	}

	if err != nil {
		return nil, fmt.Errorf("检索失败: %w", err)
	}

	// 2. 构建上下文。这里保留相似度，便于模型和前端解释检索依据。
	contextParts := make([]string, 0, len(results))
	for i, r := range results {
		contextParts = append(contextParts, fmt.Sprintf("[%d] (相似度: %.2f%%)\n%s", i+1, r.Similarity*100, r.Content))
	}
	context := strings.Join(contextParts, "\n\n")

	// 3. 调用 LLM 生成回答。RAG 的关键是限制模型只基于检索到的平台数据回答。
	prompt := fmt.Sprintf(`基于以下检索到的信息，回答用户的问题。

检索到的相关信息：
%s

用户问题：%s

请根据上述信息给出专业、准确的回答。如果信息不足，请说明。`, context, req.Query)

	answer, err := r.cozeClient.Chat(ctx, prompt)
	if err != nil {
		// 降级：直接返回检索结果
		answer = fmt.Sprintf("检索到 %d 条相关记录，请查看详情。", len(results))
	}

	return &RAGResponse{
		Answer:           answer,
		RetrievedContext: results,
		Query:            req.Query,
	}, nil
}

// MatchTalentToJob 使用 RAG 支持可解释的人岗匹配。
// 它会根据职位检索相似人才、根据人才检索相似职位，再让 LLM 总结优势、差距和录用建议。
func (r *RAGEngine) MatchTalentToJob(ctx context.Context, talentContent, jobContent string) (*MatchResult, error) {
	// 1. 用职位描述检索相似人才，作为“历史候选人画像”参考。
	similarTalents, err := r.vectorStore.SearchSimilarTalents(ctx, jobContent, 10)
	if err != nil {
		return nil, err
	}

	// 2. 用人才信息检索相似职位，帮助判断候选人更适合哪些岗位方向。
	similarJobs, err := r.vectorStore.SearchSimilarJobs(ctx, talentContent, 10)
	if err != nil {
		return nil, err
	}

	// 3. 构建匹配上下文，把目标对象和检索证据放在同一个 prompt 中。
	contextBuilder := strings.Builder{}
	contextBuilder.WriteString("【目标人才信息】\n")
	contextBuilder.WriteString(talentContent)
	contextBuilder.WriteString("\n\n【目标职位要求】\n")
	contextBuilder.WriteString(jobContent)

	if len(similarTalents) > 0 {
		contextBuilder.WriteString("\n\n【相似人才参考】\n")
		for i, t := range similarTalents[:min(3, len(similarTalents))] {
			contextBuilder.WriteString(fmt.Sprintf("%d. (相似度%.0f%%) %s\n", i+1, t.Similarity*100, truncate(t.Content, 200)))
		}
	}

	if len(similarJobs) > 0 {
		contextBuilder.WriteString("\n\n【相似职位参考】\n")
		for i, j := range similarJobs[:min(3, len(similarJobs))] {
			contextBuilder.WriteString(fmt.Sprintf("%d. (相似度%.0f%%) %s\n", i+1, j.Similarity*100, truncate(j.Content, 200)))
		}
	}

	// 4. 调用 LLM 生成匹配分析，输出可给 HR 直接阅读的建议。
	prompt := fmt.Sprintf(`作为专业的HR顾问，请分析以下人才与职位的匹配情况：

%s

请输出JSON格式的分析结果：
{
  "match_score": 0-100的匹配分数,
  "summary": "一句话总结",
  "strengths": ["优势1", "优势2"],
  "gaps": ["不足1", "不足2"],
  "recommendation": "建议"
}`, contextBuilder.String())

	analysis, err := r.cozeClient.Chat(ctx, prompt)
	if err != nil {
		// 降级处理
		analysis = "基于语义分析，该候选人与职位有一定匹配度。"
	}

	return &MatchResult{
		Analysis:       analysis,
		SimilarTalents: similarTalents,
		SimilarJobs:    similarJobs,
	}, nil
}

// MatchResult 匹配结果
type MatchResult struct {
	Analysis       string         `json:"analysis"`
	SimilarTalents []SearchResult `json:"similar_talents"`
	SimilarJobs    []SearchResult `json:"similar_jobs"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
