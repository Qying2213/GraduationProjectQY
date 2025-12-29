package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	client *elasticsearch.Client
	once   sync.Once
)

// Config ES配置
type Config struct {
	Addresses []string
	Username  string
	Password  string
}

// GetClient 获取ES客户端单例
func GetClient() *elasticsearch.Client {
	once.Do(func() {
		cfg := elasticsearch.Config{
			Addresses: []string{getEnv("ES_URL", "http://localhost:9200")},
		}

		username := getEnv("ES_USERNAME", "")
		password := getEnv("ES_PASSWORD", "")
		if username != "" && password != "" {
			cfg.Username = username
			cfg.Password = password
		}

		var err error
		client, err = elasticsearch.NewClient(cfg)
		if err != nil {
			log.Printf("Error creating ES client: %s", err)
			return
		}

		// 测试连接
		res, err := client.Info()
		if err != nil {
			log.Printf("Error connecting to ES: %s", err)
			return
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("ES connection error: %s", res.String())
		} else {
			log.Println("Successfully connected to Elasticsearch")
		}
	})
	return client
}

// InitClient 初始化ES客户端
func InitClient(config *Config) error {
	cfg := elasticsearch.Config{
		Addresses: config.Addresses,
	}

	if config.Username != "" && config.Password != "" {
		cfg.Username = config.Username
		cfg.Password = config.Password
	}

	var err error
	client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("error creating ES client: %w", err)
	}

	// 测试连接
	res, err := client.Info()
	if err != nil {
		return fmt.Errorf("error connecting to ES: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ES connection error: %s", res.String())
	}

	log.Println("Successfully connected to Elasticsearch")
	return nil
}

// CreateIndex 创建索引
func CreateIndex(indexName string, mapping map[string]interface{}) error {
	es := GetClient()
	if es == nil {
		return fmt.Errorf("ES client not initialized")
	}

	// 检查索引是否存在
	res, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil // 索引已存在
	}

	// 创建索引
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		return err
	}

	res, err = es.Indices.Create(
		indexName,
		es.Indices.Create.WithBody(&buf),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	return nil
}

// IndexDocument 索引文档
func IndexDocument(indexName string, docID string, doc interface{}) error {
	es := GetClient()
	if es == nil {
		return fmt.Errorf("ES client not initialized")
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

// Search 搜索文档
func Search(indexName string, query map[string]interface{}) ([]map[string]interface{}, int64, error) {
	es := GetClient()
	if es == nil {
		return nil, 0, fmt.Errorf("ES client not initialized")
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(indexName),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	hits := result["hits"].(map[string]interface{})
	total := int64(hits["total"].(map[string]interface{})["value"].(float64))

	var docs []map[string]interface{}
	for _, hit := range hits["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		source["_id"] = hit.(map[string]interface{})["_id"]
		docs = append(docs, source)
	}

	return docs, total, nil
}

// DeleteByQuery 按条件删除
func DeleteByQuery(indexName string, query map[string]interface{}) error {
	es := GetClient()
	if es == nil {
		return fmt.Errorf("ES client not initialized")
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return err
	}

	res, err := es.DeleteByQuery(
		[]string{indexName},
		&buf,
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete error: %s", res.String())
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// BulkIndex 批量索引
func BulkIndex(indexName string, docs []interface{}) error {
	es := GetClient()
	if es == nil {
		return fmt.Errorf("ES client not initialized")
	}

	var buf bytes.Buffer
	for i, doc := range docs {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
				"_id":    fmt.Sprintf("%d-%d", time.Now().UnixNano(), i),
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return err
		}
	}

	res, err := es.Bulk(&buf, es.Bulk.WithRefresh("true"))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk index error: %s", res.String())
	}

	return nil
}
