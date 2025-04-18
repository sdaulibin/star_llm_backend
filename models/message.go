package models

import (
	"log"
	"time"
)

// Message 表示数据库中的消息记录
type Message struct {
	ID             int       `gorm:"primaryKey"`
	UserID         string    `gorm:"size:10"`
	SessionID      string    `gorm:"size:32;not null"`
	MessageID      string    `gorm:"type:uuid"`
	ConversationID string    `gorm:"type:uuid"`
	Query          string    `gorm:"type:text;not null"`
	Answer         string    `gorm:"type:text;not null"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	IsSafe         bool      `gorm:"default:false;not null"`
	IsLike         bool      `gorm:"default:false;not null"`
	IsCollect      bool      `gorm:"default:false;not null"`
	IsDelete       bool      `gorm:"default:false;not null"`
	CurrentID      string    `gorm:"type:uuid"`
	IsStop         bool      `gorm:"default:false;not null"`
	FileID         string    `gorm:"type:uuid"`
	TaskID         string    `gorm:"type:uuid"`
}

// CreateMessage 创建新消息
func CreateMessage(message *Message) error {
	return DB.Create(message).Error
}

// GetMessageByID 通过ID获取消息（排除已删除的消息）
func GetMessageByID(id int) (*Message, error) {
	var message Message
	err := DB.Where("id = ? AND is_delete = ?", id, false).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetMessageByMessageIDAndSessionID 通过MessageID和SessionID获取消息（排除已删除的消息）
func GetMessageByMessageIDAndSessionID(messageID, sessionID string) (*Message, error) {
	var message Message
	err := DB.Where("message_id = ? AND session_id = ? AND is_delete = ?", messageID, sessionID, false).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// UpdateMessageLikeStatus 更新消息的点赞状态
func UpdateMessageLikeStatus(messageID, sessionID string, isLike bool) error {
	return DB.Model(&Message{}).Where("message_id = ? AND session_id = ?", messageID, sessionID).Update("is_like", isLike).Error
}

// UpdateMessageStopStatus 更新消息的停止状态
func UpdateMessageStopStatus(taskId string, isStop bool) error {
	log.Printf("[数据库] 更新消息: task_id=%s", taskId)
	return DB.Debug().Model(&Message{}).Where("task_id = ?", taskId).Update("is_stop", isStop).Error
}

// UpdateMessageByCurrentID 通过CurrentID更新消息
func UpdateMessageByCurrentID(message *Message) error {
	return DB.Model(&Message{}).Where("current_id = ?", message.CurrentID).Updates(map[string]interface{}{
		"answer":          message.Answer,
		"conversation_id": message.ConversationID,
		"message_id":      message.MessageID,
		"task_id":         message.TaskID,
		"updated_at":      message.UpdatedAt,
	}).Error
}

// DeleteMessage 逻辑删除消息（通过设置is_delete为true）
func DeleteMessage(messageID, sessionID string) error {
	return DB.Model(&Message{}).Where("message_id = ? AND session_id = ?", messageID, sessionID).Update("is_delete", true).Error
}

// GetMessagesByUserIDAndSessionID 获取用户特定会话的消息列表
func GetMessages(userID, sessionID string, page, pageSize int) ([]Message, int64, error) {
	var messages []Message
	var total int64

	offset := (page - 1) * pageSize

	db := DB.Model(&Message{}).Where("is_delete = ?", false)
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	if sessionID != "" {
		db = db.Where("session_id = ?", sessionID)
	}

	// 获取总记录数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// UpdateCollectStatus 更新消息的收藏状态
func UpdateCollectStatus(messageID, sessionID string, isCollect bool) error {
	return DB.Model(&Message{}).Where("message_id = ? AND session_id = ?", messageID, sessionID).Update("is_collect", isCollect).Error
}
