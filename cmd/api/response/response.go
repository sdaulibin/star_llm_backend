package response

// ChatResponse 表示从 Dify 接收的聊天响应
type ChatResponse struct {
	TaskID         string                 `json:"task_id"`
	MessageID      string                 `json:"message_id"`
	ConversationID string                 `json:"conversation_id"`
	Answer         string                 `json:"answer"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// StreamChunk 表示流式响应中的一个数据块
type StreamChunk struct {
	Event          string                 `json:"event"`
	TaskID         string                 `json:"task_id"`
	MessageID      string                 `json:"message_id"`
	ConversationID string                 `json:"conversation_id"`
	Answer         string                 `json:"answer"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// OATokenVerifyResponse 定义验证token的响应结构
type OATokenVerifyResponse struct {
	Status string `json:"status"` // 响应代号：0 认证通过，1 认证失败
	Msg    string `json:"msg"`    // 一级响应消息
}
