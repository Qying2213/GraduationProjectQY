package embedding

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// Client 火山引擎豆包Embedding客户端
type Client struct {
	endpoint string
	apiKey   string // 火山引擎API Key (ARK_API_KEY)
	modelID  string // 模型ID
	client   *http.Client
	cache    sync.Map // 简单的内存缓存
}

var defaultClient *Client
var once sync.Once

// GetClient 获取默认客户端（单例）
func GetClient() *Client {
	once.Do(func() {
		defaultClient = NewClient(&Config{
			Endpoint: os.Getenv("VOLC_ENDPOINT"),
			APIKey:   os.Getenv("ARK_API_KEY"), // 使用 ARK_API_KEY
			ModelID:  os.Getenv("VOLC_MODEL_ID"),
		})
	})
	return defaultClient
}

// Config 配置
type Config struct {
	Endpoint string
	APIKey   string // 火山引擎API Key
	ModelID  string // 模型ID，如 doubao-embedding-large-text-250515
}

// NewClient 创建客户端
func NewClient(cfg *Config) *Client {
	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://ark.cn-beijing.volces.com/api/v3"
	}
	if cfg.ModelID == "" {
		cfg.ModelID = "doubao-embedding-large-text-250515" // 最新模型
	}

	return &Client{
		endpoint: cfg.Endpoint,
		apiKey:   cfg.APIKey,
		modelID:  cfg.ModelID,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

// IsConfigured 检查是否已配置
func (c *Client) IsConfigured() bool {
	return c.apiKey != ""
}

// MultimodalInput 多模态输入项
type MultimodalInput struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// EmbeddingRequest 请求结构（火山引擎多模态格式）
type EmbeddingRequest struct {
	Model          string            `json:"model"`
	Input          []MultimodalInput `json:"input"`
	Dimensions     int               `json:"dimensions,omitempty"`
	EncodingFormat string            `json:"encoding_format,omitempty"`
}

// EmbeddingResponse 响应结构（多模态格式 - 单个文本）
type EmbeddingResponse struct {
	Created int64 `json:"created"`
	Data    struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// getSingleEmbedding 获取单个文本的向量（多模态API每次只能处理一个文本）
func (c *Client) getSingleEmbedding(ctx context.Context, text string) ([]float64, error) {
	// 构建多模态格式的请求
	reqBody := EmbeddingRequest{
		Model: c.modelID,
		Input: []MultimodalInput{
			{Type: "text", Text: text},
		},
		Dimensions:     1024,
		EncodingFormat: "float",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	url := c.endpoint
	fmt.Printf("[Embedding] Calling API: %s\n", url)
	fmt.Printf("[Embedding] Model: %s, Text length: %d\n", c.modelID, len(text))
	fmt.Printf("[Embedding] Request body: %s\n", string(jsonBody))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	fmt.Printf("[Embedding] HTTP Status: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[Embedding] Error response: %s\n", string(body))
		return nil, fmt.Errorf("API返回错误: %d - %s", resp.StatusCode, string(body))
	}

	fmt.Printf("[Embedding] Success! Response: %s\n", string(body))

	var apiResp EmbeddingResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(apiResp.Data.Embedding) == 0 {
		return nil, fmt.Errorf("API返回空向量")
	}

	fmt.Printf("[Embedding] Got embedding with %d dimensions\n", len(apiResp.Data.Embedding))
	return apiResp.Data.Embedding, nil
}

// GetEmbeddings 获取文本向量（批量）
func (c *Client) GetEmbeddings(ctx context.Context, texts []string) ([][]float64, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	// 检查缓存
	results := make([][]float64, len(texts))
	uncachedTexts := []string{}
	uncachedIndices := []int{}

	for i, text := range texts {
		cacheKey := c.getCacheKey(text)
		if cached, ok := c.cache.Load(cacheKey); ok {
			results[i] = cached.([]float64)
		} else {
			uncachedTexts = append(uncachedTexts, text)
			uncachedIndices = append(uncachedIndices, i)
		}
	}

	// 如果所有都命中缓存
	if len(uncachedTexts) == 0 {
		fmt.Printf("[Embedding] All %d texts hit cache\n", len(texts))
		return results, nil
	}

	// 如果未配置，返回模拟向量
	if !c.IsConfigured() {
		fmt.Println("[Embedding] WARNING: API Key not configured, using mock embeddings")
		fmt.Printf("[Embedding] ARK_API_KEY value: '%s'\n", c.apiKey)
		for i, idx := range uncachedIndices {
			results[idx] = c.mockEmbedding(uncachedTexts[i])
		}
		return results, nil
	}

	fmt.Printf("[Embedding] Processing %d uncached texts (total: %d)\n", len(uncachedTexts), len(texts))

	// 多模态API每次只能处理一个文本，需要逐个调用
	for i, text := range uncachedTexts {
		idx := uncachedIndices[i]

		embedding, err := c.getSingleEmbedding(ctx, text)
		if err != nil {
			fmt.Printf("[Embedding] Error for text %d: %v, using mock\n", i, err)
			results[idx] = c.mockEmbedding(text)
		} else {
			results[idx] = embedding
			// 缓存结果
			cacheKey := c.getCacheKey(text)
			c.cache.Store(cacheKey, embedding)
		}
	}

	return results, nil
}

// GetEmbedding 获取单个文本的向量
func (c *Client) GetEmbedding(ctx context.Context, text string) ([]float64, error) {
	embeddings, err := c.GetEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("未获取到向量")
	}
	return embeddings[0], nil
}

// getCacheKey 生成缓存键
func (c *Client) getCacheKey(text string) string {
	h := sha256.Sum256([]byte(text))
	return hex.EncodeToString(h[:])
}

// mockEmbedding 生成模拟向量（用于测试或降级）
func (c *Client) mockEmbedding(text string) []float64 {
	// 基于文本内容生成确定性的模拟向量
	dim := 1024 // Doubao-embedding-large 维度
	embedding := make([]float64, dim)

	// 使用文本的hash来生成确定性的向量
	h := sha256.Sum256([]byte(text))
	for i := 0; i < dim; i++ {
		// 使用hash的不同部分生成向量值
		idx := i % 32
		embedding[i] = float64(h[idx])/255.0*2 - 1 // 归一化到[-1, 1]
	}

	// 归一化
	return NormalizeEmbedding(embedding)
}

// CosineSimilarity 计算余弦相似度
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// NormalizeEmbedding 归一化向量
func NormalizeEmbedding(embedding []float64) []float64 {
	var norm float64
	for _, v := range embedding {
		norm += v * v
	}
	norm = math.Sqrt(norm)

	if norm == 0 {
		return embedding
	}

	normalized := make([]float64, len(embedding))
	for i, v := range embedding {
		normalized[i] = v / norm
	}

	return normalized
}

// SimilarityResult 相似度结果
type SimilarityResult struct {
	Index      int     `json:"index"`
	Similarity float64 `json:"similarity"`
}

// RankBySimilarity 按相似度排序
func RankBySimilarity(target []float64, candidates [][]float64) []SimilarityResult {
	results := make([]SimilarityResult, len(candidates))

	for i, candidate := range candidates {
		results[i] = SimilarityResult{
			Index:      i,
			Similarity: CosineSimilarity(target, candidate),
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	return results
}

// SkillMatcher 技能匹配器
type SkillMatcher struct {
	client *Client
}

// NewSkillMatcher 创建技能匹配器
func NewSkillMatcher(client *Client) *SkillMatcher {
	return &SkillMatcher{client: client}
}

// MatchSkills 匹配技能（使用语义相似度）
func (m *SkillMatcher) MatchSkills(ctx context.Context, talentSkills, jobSkills []string) (float64, []string, error) {
	if len(jobSkills) == 0 {
		return 0.5, []string{"职位未指定技能要求"}, nil
	}

	if len(talentSkills) == 0 {
		return 0, []string{"候选人未提供技能信息"}, nil
	}

	// 获取所有技能的向量
	allSkills := append(talentSkills, jobSkills...)
	embeddings, err := m.client.GetEmbeddings(ctx, allSkills)
	if err != nil {
		// 降级到精确匹配
		return m.exactMatch(talentSkills, jobSkills)
	}

	talentEmbeddings := embeddings[:len(talentSkills)]
	jobEmbeddings := embeddings[len(talentSkills):]

	// 计算匹配分数
	matchedSkills := []string{}
	partialMatchSkills := []string{}
	totalScore := 0.0

	for i, jobEmb := range jobEmbeddings {
		maxSim := 0.0
		for _, talentEmb := range talentEmbeddings {
			sim := CosineSimilarity(jobEmb, talentEmb)
			if sim > maxSim {
				maxSim = sim
			}
		}

		if maxSim > 0.85 { // 高度匹配
			matchedSkills = append(matchedSkills, jobSkills[i])
			totalScore += 1.0
		} else if maxSim > 0.7 { // 部分匹配
			partialMatchSkills = append(partialMatchSkills, jobSkills[i])
			totalScore += 0.6
		} else if maxSim > 0.5 { // 相关
			totalScore += 0.3
		}
	}

	score := totalScore / float64(len(jobSkills))
	details := []string{}

	if len(matchedSkills) > 0 {
		details = append(details, "完全匹配: "+strings.Join(matchedSkills, ", "))
	}
	if len(partialMatchSkills) > 0 {
		details = append(details, "部分匹配: "+strings.Join(partialMatchSkills, ", "))
	}

	missingCount := len(jobSkills) - len(matchedSkills) - len(partialMatchSkills)
	if missingCount > 0 {
		details = append(details, fmt.Sprintf("缺少 %d 项核心技能", missingCount))
	}

	return score, details, nil
}

// exactMatch 精确匹配（降级方案）
func (m *SkillMatcher) exactMatch(talentSkills, jobSkills []string) (float64, []string, error) {
	matchedSkills := []string{}

	talentSet := make(map[string]bool)
	for _, s := range talentSkills {
		talentSet[strings.ToLower(strings.TrimSpace(s))] = true
	}

	for _, js := range jobSkills {
		jsLower := strings.ToLower(strings.TrimSpace(js))
		if talentSet[jsLower] {
			matchedSkills = append(matchedSkills, js)
		} else {
			// 检查部分匹配
			for ts := range talentSet {
				if strings.Contains(ts, jsLower) || strings.Contains(jsLower, ts) {
					matchedSkills = append(matchedSkills, js)
					break
				}
			}
		}
	}

	score := float64(len(matchedSkills)) / float64(len(jobSkills))
	details := []string{}
	if len(matchedSkills) > 0 {
		details = append(details, "匹配技能: "+strings.Join(matchedSkills, ", "))
	}

	return score, details, nil
}

// TextSimilarity 计算两段文本的语义相似度
func (c *Client) TextSimilarity(ctx context.Context, text1, text2 string) (float64, error) {
	embeddings, err := c.GetEmbeddings(ctx, []string{text1, text2})
	if err != nil {
		return 0, err
	}

	if len(embeddings) < 2 {
		return 0, fmt.Errorf("获取向量失败")
	}

	return CosineSimilarity(embeddings[0], embeddings[1]), nil
}
