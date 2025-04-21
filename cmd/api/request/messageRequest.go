package request

// GetMessagesRequest 获取消息列表的请求参数
type GetMessagesRequest struct {
	UserID    string `json:"user_id" form:"user_id" binding:"required"`
	SessionID string `json:"session_id" form:"session_id"`
	Query     string `json:"query"`
	Page      int    `json:"page" form:"page" binding:"required,min=1"`
	PageSize  int    `json:"page_size" form:"page_size" binding:"required,min=1,max=1000"`
}

// UpdateCollectStatusRequest 更新消息收藏状态的请求参数
type UpdateCollectStatusRequest struct {
	MessageID string `json:"message_id" binding:"required"`
	SessionID string `json:"session_id" binding:"required"`
	IsCollect bool   `json:"is_collect""`
}

type DeleteMessageRequest struct {
	MessageID string `json:"message_id" binding:"required"`
	SessionID string `json:"session_id" binding:"required"`
}

// DeleteMessagesRequest 批量删除消息的请求参数
type DeleteMessagesRequest struct {
	MessageIDs []string `json:"message_ids" binding:"required"`
	SessionID  string   `json:"session_id" binding:"required"`
}
