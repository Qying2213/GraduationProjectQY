package embedding

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Client 火山引擎豆包Embedding客户端
type Client struct {
	endpoint  string
	accessKey string
	secretKey string
	model     string
	client    *http.Client
}

// Config 配置
type Config struct {
	Endpoint  string // https://ark.cn-beijing.volces.com/api/v3
	AccessKey string
	SecretKey string
	Model     string // doubao-embedding
}

// NewClient 创建客户端
func NewClient(cfg *Config) *Client {
	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://ark.cn-beijing.volces.com/api/v3"
	}
	if cfg.Model == "" {
		cfg.Model = "doubao-embedding"
	}

	return &Client{
		endpoint:  cfg.Endpoint,
		accessKey: cfg.AccessKey,
		secretKey: cfg.SecretKey,
		model:     cfg.Model,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

// EmbeddingRequest 请求结构
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// EmbeddingResponse 响应结构
type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// GetEmbeddings 获取文本向量
func (c *Client) GetEmbeddings(ctx context.Context, texts []string) ([][]float64, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	reqBody := EmbeddingRequest{
		Model: c.model,
		Input: texts,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	url := c.endpoint + "/embeddings"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 火山引擎签名认证
	c.signRequest(req, jsonBody)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API错误 %d: %s", resp.StatusCode, string(body))
	}

	var result EmbeddingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 按index排序并提取embedding
	embeddings := make([][]float64, len(result.Data))
	for _, item := range result.Data {
		if item.Index < len(embeddings) {
			embeddings[item.Index] = item.Embedding
		}
	}

	return embeddings, nil
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

// signRequest 火山引擎V4签名
func (c *Client) signRequest(req *http.Request, body []byte) {
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	date := timestamp[:8]

	// 设置必要的头
	req.Header.Set("X-Date", timestamp)
	req.Header.Set("Host", req.URL.Host)

	// 计算payload hash
	hashedPayload := sha256Hash(body)

	// 构建规范请求
	signedHeaders := "content-type;host;x-date"
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\nx-date:%s\n",
		req.Header.Get("Content-Type"),
		req.URL.Host,
		timestamp)

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		req.Method,
		req.URL.Path,
		req.URL.RawQuery,
		canonicalHeaders,
		signedHeaders,
		hashedPayload)

	// 构建待签名字符串
	credentialScope := fmt.Sprintf("%s/cn-beijing/ark/request", date)
	stringToSign := fmt.Sprintf("HMAC-SHA256\n%s\n%s\n%s",
		timestamp,
		credentialScope,
		sha256Hash([]byte(canonicalRequest)))

	// 计算签名
	kDate := hmacSHA256([]byte("VOLC"+c.secretKey), date)
	kRegion := hmacSHA256(kDate, "cn-beijing")
	kService := hmacSHA256(kRegion, "ark")
	kSigning := hmacSHA256(kService, "request")
	signature := hex.EncodeToString(hmacSHA256(kSigning, stringToSign))

	// 设置Authorization头
	authHeader := fmt.Sprintf("HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		c.accessKey, credentialScope, signedHeaders, signature)
	req.Header.Set("Authorization", authHeader)
}

func sha256Hash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
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

// EuclideanDistance 计算欧氏距离
func EuclideanDistance(a, b []float64) float64 {
	if len(a) != len(b) {
		return math.MaxFloat64
	}

	var sum float64
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

// FindMostSimilar 找到最相似的文本
func FindMostSimilar(target []float64, candidates [][]float64) (int, float64) {
	if len(candidates) == 0 {
		return -1, 0
	}

	maxSim := -1.0
	maxIdx := 0

	for i, candidate := range candidates {
		sim := CosineSimilarity(target, candidate)
		if sim > maxSim {
			maxSim = sim
			maxIdx = i
		}
	}

	return maxIdx, maxSim
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

	// 按相似度降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	return results
}

// BatchSimilarity 批量计算相似度矩阵
func BatchSimilarity(embeddings [][]float64) [][]float64 {
	n := len(embeddings)
	matrix := make([][]float64, n)

	for i := 0; i < n; i++ {
		matrix[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if i == j {
				matrix[i][j] = 1.0
			} else if j < i {
				matrix[i][j] = matrix[j][i]
			} else {
				matrix[i][j] = CosineSimilarity(embeddings[i], embeddings[j])
			}
		}
	}

	return matrix
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

// CombineEmbeddings 组合多个向量（平均）
func CombineEmbeddings(embeddings [][]float64) []float64 {
	if len(embeddings) == 0 {
		return nil
	}

	dim := len(embeddings[0])
	combined := make([]float64, dim)

	for _, emb := range embeddings {
		for i, v := range emb {
			combined[i] += v
		}
	}

	n := float64(len(embeddings))
	for i := range combined {
		combined[i] /= n
	}

	return NormalizeEmbedding(combined)
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

// SkillMatcher 技能匹配器
type SkillMatcher struct {
	client *Client
	cache  map[string][]float64
}

// NewSkillMatcher 创建技能匹配器
func NewSkillMatcher(client *Client) *SkillMatcher {
	return &SkillMatcher{
		client: client,
		cache:  make(map[string][]float64),
	}
}

// MatchSkills 匹配技能
func (m *SkillMatcher) MatchSkills(ctx context.Context, talentSkills, jobSkills []string) (float64, []string, error) {
	if len(jobSkills) == 0 {
		return 0.5, []string{"职位未指定技能要求"}, nil
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
	totalScore := 0.0

	for i, jobEmb := range jobEmbeddings {
		maxSim := 0.0
		for _, talentEmb := range talentEmbeddings {
			sim := CosineSimilarity(jobEmb, talentEmb)
			if sim > maxSim {
				maxSim = sim
			}
		}

		if maxSim > 0.8 { // 相似度阈值
			matchedSkills = append(matchedSkills, jobSkills[i])
			totalScore += 1.0
		} else if maxSim > 0.6 {
			totalScore += 0.5
		}
	}

	score := totalScore / float64(len(jobSkills))
	details := []string{}
	if len(matchedSkills) > 0 {
		details = append(details, "匹配技能: "+strings.Join(matchedSkills, ", "))
	}

	return score, details, nil
}

// exactMatch 精确匹配（降级方案）
func (m *SkillMatcher) exactMatch(talentSkills, jobSkills []string) (float64, []string, error) {
	matchedSkills := []string{}

	talentSet := make(map[string]bool)
	for _, s := range talentSkills {
		talentSet[strings.ToLower(s)] = true
	}

	for _, js := range jobSkills {
		if talentSet[strings.ToLower(js)] {
			matchedSkills = append(matchedSkills, js)
		}
	}

	score := float64(len(matchedSkills)) / float64(len(jobSkills))
	details := []string{}
	if len(matchedSkills) > 0 {
		details = append(details, "匹配技能: "+strings.Join(matchedSkills, ", "))
	}

	return score, details, nil
}
