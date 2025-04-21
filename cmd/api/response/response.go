package response

import "time"

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

type MessageResponse struct {
	ID             int
	UserID         string
	SessionID      string
	MessageID      string
	ConversationID string
	Query          string
	Answer         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IsSafe         bool
	IsLike         bool
	IsCollect      bool
	IsDelete       bool
	CurrentID      string
	IsStop         bool
	FileID         string
	TaskID         string
	File           []FileJSON
}

// FileJSON 表示用于 JSON 序列化的文件记录
type FileJSON struct {
	ID               int
	UserID           string
	FileID           string
	OriginalFilename string
	LocalFilename    string
	FilePath         string
	FileSize         int64
	FileType         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
