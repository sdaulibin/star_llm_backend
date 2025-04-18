package dify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"star_llm_backend_n/cmd/api/request"
	"star_llm_backend_n/cmd/api/response"
	"star_llm_backend_n/config"
	"star_llm_backend_n/logs"
	"time"

	"github.com/gin-gonic/gin"
)

// Client 是 Dify API 的客户端
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient 创建一个新的 Dify API 客户端
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 180,
		},
	}
}

// SendChatMessage 发送聊天消息到 Dify API
func (client *Client) SendChatMessage(ctx context.Context, req request.ChatMessageRequest) (*response.ChatResponse, error) {
	if req.ResponseMode == "streaming" {
		return nil, fmt.Errorf("使用 SendChatMessageStream 方法处理流式响应")
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", client.BaseURL+"/v1/chat-messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+client.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 返回错误状态码 %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp response.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &chatResp, nil
}

// SendChatMessageStream 发送聊天消息并处理流式响应
func (client *Client) SendChatMessageStream(ctx context.Context, chatMessageRequest request.ChatMessageRequest, chunkHandler func(response.StreamChunk) error) error {
	chatMessageRequest.ResponseMode = "streaming"

	reqBody, err := json.Marshal(chatMessageRequest)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", client.BaseURL+"/v1/chat-messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+client.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.HTTPClient.Do(httpReq)
	if err != nil {
		// 检查是否是超时错误
		if err, ok := err.(net.Error); ok && err.Timeout() {
			logs.Logger.Error("请求超时: 超过180秒未收到完整响应")
			return fmt.Errorf("请求超时: 超过180秒未收到完整响应，请稍后重试")
		}
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API 返回错误状态码 %d: %s", resp.StatusCode, string(bodyBytes))
	}

	reader := NewSSEReader(resp.Body)
	for {
		event, err := reader.ReadEvent()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("读取事件流失败: %w", err)
		}

		if event.Data == "" {
			continue
		}

		var chunk response.StreamChunk
		if err := json.Unmarshal([]byte(event.Data), &chunk); err != nil {
			logs.Logger.Errorf("解析事件数据失败: %v, 数据: %s", err, event.Data)
			continue
		}

		if err := chunkHandler(chunk); err != nil {
			logs.Logger.Errorf("处理chunk数据块失败: %s", err)
			return fmt.Errorf("处理chunk数据块失败: %w", err)
		}

		// 如果收到消息结束事件，结束流式处理
		if chunk.Event == "message_end" {
			break
		}
	}

	return nil
}

// StopChatMessage 停止正在进行的聊天消息生成
func (client *Client) StopChatMessage(ctx context.Context, task_id string) error {
	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/chat-messages/%s/stop", client.BaseURL, task_id), nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+client.APIKey)

	resp, err := client.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API 返回错误状态码 %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// Suggested 获取下一轮建议问题列表
func (client *Client) Suggested(ctx context.Context, message_id, user_id string) ([]string, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v1/messages/%s/suggested?user=%s", client.BaseURL, message_id, user_id), nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+client.APIKey)

	resp, err := client.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 返回错误状态码 %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应
	var suggestedResp struct {
		Result string   `json:"result"`
		Data   []string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&suggestedResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if suggestedResp.Result != "success" {
		return nil, fmt.Errorf("API 返回失败结果: %s", suggestedResp.Result)
	}

	return suggestedResp.Data, nil
}

// FeedbackMessage 对消息进行反馈（点赞/点踩）
func (client *Client) FeedbackMessage(ctx context.Context, messageId string, feedbackRequest request.FeedBackRequest) error {
	logs.Logger.Infof("feedback请求体内容>>>> %s\n", feedbackRequest)
	reqBody, err := json.Marshal(feedbackRequest)
	logs.Logger.Infof("feedback请求体内容>>>> %s\n", reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/messages/%s/feedbacks", client.BaseURL, messageId), bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+client.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	logs.Logger.Infof("feedback请求地址>>>> %s\n", httpReq.URL.String())
	logs.Logger.Infof("feedback请求体内容>>>> %s\n", httpReq.Body)

	resp, err := client.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API 返回错误状态码 %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// getDifyClient 获取Dify客户端实例
func GetDifyClient(ctx *gin.Context) *Client {
	// 从应用配置或上下文中获取Dify API配置
	// 这里假设您已经在某处初始化了这些配置
	baseURL := config.GlobalConfig.API.BaseURL // 从配置中获取
	apiKey := config.GlobalConfig.API.Key      // 从配置中获取

	if apiKey == "" {
		// 如果上下文中没有，则从配置中获取
		apiKey = "your_default_api_key" // 从配置中获取默认值
	}

	return NewClient(baseURL, apiKey)
}
