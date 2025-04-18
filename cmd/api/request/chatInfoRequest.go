package request

// CreateChatInfoRequest 表示创建对话信息的请求
type CreateChatInfoRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	SessionID string `json:"session_id"`
	ChatName  string `json:"chat_name" binding:"required"`
}

// GetChatInfosByUserIDRequest 表示获取用户所有对话信息的请求
type GetChatInfosByUserIDRequest struct {
	UserID string `json:"user_id" form:"user_id" binding:"required"`
}

// UpdateChatInfoRequest 表示更新对话信息的请求
type UpdateChatInfoRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	ChatName  string `json:"chat_name" binding:"required"`
}

// DeleteChatInfoRequest 表示删除对话信息的请求
type DeleteChatInfoRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}
