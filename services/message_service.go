package services

import (
	"log"
	"time"

	"star_llm_backend_n/cmd/api/request"
	"star_llm_backend_n/cmd/api/response"
	"star_llm_backend_n/logs"
	"star_llm_backend_n/models"

	"github.com/google/uuid"
)

// SaveMessageToDB 保存消息到数据库
func SaveMessageToDB(currentID, sessionID, query, answer, userID, conversationID, fileId, messageID string) error {
	// 创建消息对象
	message := &models.Message{
		CurrentID:      currentID,
		UserID:         userID,
		SessionID:      sessionID,
		MessageID:      messageID,
		ConversationID: conversationID,
		Query:          query,
		Answer:         answer,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		IsSafe:         false,
		IsLike:         false,
		FileID:         fileId,
		TaskID:         uuid.Nil.String(),
	}

	// 如果提供了messageID，则使用它，否则生成一个新的UUID
	if len(messageID) > 0 && messageID != uuid.Nil.String() {
		message.MessageID = messageID
	} else {
		message.MessageID = uuid.New().String()
	}

	log.Printf("[服务] 保存消息到数据库: message_id=%s\n >>>>content:%s", message.MessageID, message.Answer)
	return models.CreateMessage(message)
}

// UpdateMessageToDB 根据CurrentID更新消息到数据库
func UpdateMessageToDB(currentID, answer, conversationID, messageID, taskID string) error {
	// 创建消息对象
	message := &models.Message{
		CurrentID:      currentID,
		Answer:         answer,
		ConversationID: conversationID,
		MessageID:      messageID,
		UpdatedAt:      time.Now(),
		TaskID:         taskID,
	}

	log.Printf("[服务] 更新消息到数据库: CurrentID=%s\n >>>>content:%s", message.CurrentID, message.Answer)
	return models.UpdateMessageByCurrentID(message)
}

// ConvertMessageToResponse 将 models.Message 转换为 response.MessageResponse
func ConvertMessageToResponse(message *models.Message) *response.MessageResponse {
	return &response.MessageResponse{
		ID:             message.ID,
		UserID:         message.UserID,
		SessionID:      message.SessionID,
		MessageID:      message.MessageID,
		ConversationID: message.ConversationID,
		Query:          message.Query,
		Answer:         message.Answer,
		CreatedAt:      message.CreatedAt,
		UpdatedAt:      message.UpdatedAt,
		IsSafe:         message.IsSafe,
		IsLike:         message.IsLike,
		IsCollect:      message.IsCollect,
		IsDelete:       message.IsDelete,
		CurrentID:      message.CurrentID,
		IsStop:         message.IsStop,
		FileID:         message.FileID,
		TaskID:         message.TaskID,
		File:           []response.FileJSON{},
	}
}

// GetMessagesWithResponse 获取消息列表并转换为 MessageResponse 格式
func GetMessages(getMessagesRequest request.GetMessagesRequest, page, pageSize int) ([]response.MessageResponse, int64, error) {
	// 调用 models 包中的 GetMessages 方法获取消息列表
	messages, total, err := models.GetMessages(getMessagesRequest, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	logs.Logger.Infof("[服务] 获取消息列表messages: %v", messages)

	// 转换为 MessageResponse 格式
	responseMessages := make([]response.MessageResponse, 0, len(messages))
	for _, msg := range messages {
		msgResponse := ConvertMessageToResponse(&msg)

		// 如果消息有关联的文件，获取文件信息
		if msg.FileID != "" && msg.FileID != uuid.Nil.String() {
			file, err := models.GetFileByFileID(msg.FileID)
			if err == nil {
				// 将文件信息添加到响应中
				fileJSON := response.FileJSON{
					ID:               file.ID,
					UserID:           file.UserID,
					FileID:           file.FileID,
					OriginalFilename: file.OriginalFilename,
					LocalFilename:    file.LocalFilename,
					FilePath:         file.FilePath,
					FileSize:         file.FileSize,
					FileType:         file.FileType,
					CreatedAt:        file.CreatedAt,
					UpdatedAt:        file.UpdatedAt,
				}
				msgResponse.File = append(msgResponse.File, fileJSON)
			}
			logs.Logger.Infof("[服务] 获取文件信息: %s,%v\n", file.FileID, file)
		}

		responseMessages = append(responseMessages, *msgResponse)
	}
	logs.Logger.Infof("[服务] 获取消息列表responseMessages: %v", responseMessages)
	return responseMessages, total, nil
}
