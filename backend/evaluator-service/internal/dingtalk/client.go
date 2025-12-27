package dingtalk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"evaluator-service/internal/logging"
	"evaluator-service/internal/models"

	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	streamclient "github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
)

type DingTalkMessage struct {
	MsgType  string       `json:"msgtype"`
	Text     DingTalkText `json:"text,omitempty"`
	Markdown *Markdown    `json:"markdown,omitempty"`
	At       DingTalkAt   `json:"at"`
}

type DingTalkText struct {
	Content string `json:"content"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingTalkAt struct {
	AtUserIds []string `json:"atUserIds,omitempty"`
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingTalkResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type Client struct {
	config       *models.DingTalkConfig
	log          *logging.Logger
	streamClient *streamclient.StreamClient
	onMessage    func(ctx context.Context, content string, senderID string) error
}

func NewClient(config *models.DingTalkConfig, log *logging.Logger) *Client {
	return &Client{
		config: config,
		log:    log,
	}
}

// sign 生成钉钉签名（如果secret为空则返回空字符串）
func (c *Client) sign(secret string, timestamp int64) string {
	if secret == "" {
		return ""
	}
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// SendTextMessage 发送文本消息
func (c *Client) SendTextMessage(content string, atUserIds []string, isAtAll bool) error {
	timestamp := time.Now().UnixMilli()
	sign := c.sign(c.config.Secret, timestamp)

	webhookURL := c.config.Webhook
	// 只有在有签名时才添加签名参数
	if sign != "" {
		webhookURL = fmt.Sprintf("%s&timestamp=%d&sign=%s",
			c.config.Webhook, timestamp, url.QueryEscape(sign))
	}

	c.log.Info("sending text message to dingtalk",
		logging.KV("webhook", c.config.Webhook),
		logging.KV("has_sign", sign != ""),
		logging.KV("at_users", atUserIds),
		logging.KV("content_len", len(content)))

	message := DingTalkMessage{
		MsgType: "text",
		Text: DingTalkText{
			Content: content,
		},
		At: DingTalkAt{
			AtUserIds: atUserIds,
			IsAtAll:   isAtAll,
		},
	}

	return c.sendMessage(webhookURL, message)
}

// SendMarkdownMessage 发送Markdown消息
func (c *Client) SendMarkdownMessage(title, content string, atUserIds []string, isAtAll bool) error {
	timestamp := time.Now().UnixMilli()
	sign := c.sign(c.config.Secret, timestamp)

	webhookURL := c.config.Webhook
	// 只有在有签名时才添加签名参数
	if sign != "" {
		webhookURL = fmt.Sprintf("%s&timestamp=%d&sign=%s",
			c.config.Webhook, timestamp, url.QueryEscape(sign))
	}

	c.log.Info("sending markdown message to dingtalk",
		logging.KV("webhook", c.config.Webhook),
		logging.KV("has_sign", sign != ""),
		logging.KV("title", title),
		logging.KV("at_users", atUserIds),
		logging.KV("at_count", len(atUserIds)),
		logging.KV("content_len", len(content)))

	// 钉钉Markdown消息的@功能需要在文本末尾添加@手机号
	// 但atUserIds使用的是钉钉UserID，需要通过atMobiles或在文本中@
	// 这里我们使用atMobiles字段（如果提供的是手机号）或atUserIds
	message := DingTalkMessage{
		MsgType: "markdown",
		Markdown: &Markdown{
			Title: title,
			Text:  content,
		},
		At: DingTalkAt{
			AtUserIds: atUserIds,
			IsAtAll:   isAtAll,
		},
	}

	return c.sendMessage(webhookURL, message)
}

func (c *Client) sendMessage(webhookURL string, message DingTalkMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		c.log.Error("marshal message failed", logging.Err(err))
		return fmt.Errorf("marshal message: %w", err)
	}

	c.log.Info("sending request to dingtalk",
		logging.KV("url", webhookURL),
		logging.KV("payload_size", len(jsonData)),
		logging.KV("payload", string(jsonData))) // 打印完整payload用于调试

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.log.Error("create request failed", logging.Err(err))
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.log.Error("send request failed", logging.Err(err))
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	c.log.Info("received response from dingtalk",
		logging.KV("status_code", resp.StatusCode))

	var dingResp DingTalkResponse
	if err := json.NewDecoder(resp.Body).Decode(&dingResp); err != nil {
		c.log.Error("decode response failed", logging.Err(err))
		return fmt.Errorf("decode response: %w", err)
	}

	c.log.Info("dingtalk response",
		logging.KV("errcode", dingResp.ErrCode),
		logging.KV("errmsg", dingResp.ErrMsg))

	if dingResp.ErrCode != 0 {
		c.log.Error("dingtalk api error",
			logging.KV("errcode", dingResp.ErrCode),
			logging.KV("errmsg", dingResp.ErrMsg))
		return fmt.Errorf("dingtalk api error: code=%d, msg=%s", dingResp.ErrCode, dingResp.ErrMsg)
	}

	c.log.Info("message sent successfully to dingtalk")
	return nil
}

// StartStream 启动Stream模式监听消息
func (c *Client) StartStream(ctx context.Context, onMessage func(ctx context.Context, content string, senderID string) error) error {
	if c.config.ClientID == "" || c.config.ClientSecret == "" {
		return fmt.Errorf("clientID or clientSecret is empty")
	}

	c.onMessage = onMessage

	cli := streamclient.NewStreamClient(
		streamclient.WithAppCredential(
			streamclient.NewAppCredentialConfig(c.config.ClientID, c.config.ClientSecret),
		),
	)

	cli.RegisterChatBotCallbackRouter(c.handleChatBotMessage)

	c.streamClient = cli

	go func() {
		if err := cli.Start(ctx); err != nil {
			c.log.Error("dingtalk stream client error", logging.Err(err))
		}
	}()

	c.log.Info("dingtalk stream client started")
	return nil
}

func (c *Client) handleChatBotMessage(ctx context.Context, data *chatbot.BotCallbackDataModel) (result []byte, resultErr error) {
	// 添加panic恢复机制，防止钉钉SDK的bug导致服务崩溃
	defer func() {
		if r := recover(); r != nil {
			c.log.Error("recovered from panic in handleChatBotMessage",
				logging.KV("panic", r),
				logging.KV("sender", data.SenderStaffId),
				logging.KV("content", data.Text.Content))
			result = []byte("")
			resultErr = nil // 不返回error，避免SDK重试
		}
	}()

	content := strings.TrimSpace(data.Text.Content)
	senderID := data.SenderStaffId

	c.log.Info("received dingtalk message",
		logging.KV("content", content),
		logging.KV("sender", senderID))

	if c.onMessage != nil {
		if err := c.onMessage(ctx, content, senderID); err != nil {
			c.log.Error("handle message error", logging.Err(err))
			// 回复错误消息（只回复一次，不返回error避免SDK重试）
			// 注意：这里可能触发SDK的panic（send on closed channel）
			replier := chatbot.NewChatbotReplier()
			_ = replier.SimpleReplyText(ctx, data.SessionWebhook, []byte("处理失败: "+err.Error()))
			return []byte(""), nil // 返回nil而不是err，避免重复回复
		}
	}

	return []byte(""), nil
}

// Close 关闭Stream客户端
func (c *Client) Close() {
	if c.streamClient != nil {
		c.streamClient.Close()
	}
}
