package request

type FeedBackRequest struct {
	Rating    string `json:"rating" binding:"required"`
	User      string `json:"user" binding:"required"`
	Content   string `json:"content"`
	SessionID string `json:"session_id"`
}

type StopRequest struct {
	User      string `json:"user" binding:"required"`
	SessionID string `json:"session_id"`
}

// File 表示上传到 Dify 的文件
type File struct {
	Type           string `json:"type"`                     // 文件类型，如 "image", "document" 等
	TransferMethod string `json:"transfer_method"`          // 传输方法，如 "remote_url", "local_file" 等
	URL            string `json:"url,omitempty"`            // 当 TransferMethod 为 "remote_url" 时使用
	UploadFileID   string `json:"upload_file_id,omitempty"` // 当 TransferMethod 为 "local_file" 时使用
}

// ChatRequest 表示发送到 Dify 的聊天请求
type ChatMessageRequest struct {
	Inputs         map[string]interface{} `json:"inputs"`
	Query          string                 `json:"query" binding:"required"`
	ResponseMode   string                 `json:"response_mode"` // "blocking" 或 "streaming"
	ConversationID string                 `json:"conversation_id,omitempty"`
	User           string                 `json:"user"`
	Files          []File                 `json:"files,omitempty"`
	SessionID      string                 `json:"session_id"`
}
