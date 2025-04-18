package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// ChatInfo 表示数据库中的对话信息记录
type ChatInfo struct {
	ID        int       `gorm:"primaryKey"`
	UserID    string    `gorm:"size:10"`
	SessionID string    `gorm:"type:uuid;not null"`
	ChatName  string    `gorm:"size:255;not null"`
	IsDelete  bool      `gorm:"default:false;not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (ChatInfo) TableName() string {
	return "chat_info" // 指定为单数形式
}

// CreateChatInfo 创建新的对话信息
func CreateChatInfo(chatInfo *ChatInfo) error {
	return DB.Create(chatInfo).Error
}

// GetChatInfoByID 通过ID获取对话信息
func GetChatInfoBySessionId(session_id string) (*ChatInfo, error) {
	var chatInfo ChatInfo
	err := DB.Where("session_id = ? AND is_delete = ?", session_id, false).First(&chatInfo).Error
	if err != nil {
		// 可以在这里判断是否是记录未找到的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &chatInfo, nil
}

// GetChatInfoByConversationID 通过对话ID获取对话信息
func GetChatInfoByConversationID(conversationID string) (*ChatInfo, error) {
	var chatInfo ChatInfo
	err := DB.Where("conversation_id = ?", conversationID).First(&chatInfo).Error
	if err != nil {
		return nil, err
	}
	return &chatInfo, nil
}

// GetChatInfosByUserID 获取用户的所有对话信息
func GetChatInfos(userID string) ([]ChatInfo, error) {
	var chatInfos []ChatInfo
	err := DB.Where("user_id = ? AND is_delete = ?", userID, false).Order("updated_at DESC").Find(&chatInfos).Error
	if err != nil {
		return nil, err
	}
	return chatInfos, nil
}

// UpdateChatInfo 更新对话信息
func UpdateChatInfo(chatInfo *ChatInfo) error {
	return DB.Model(&ChatInfo{}).Where("session_id = ?", chatInfo.SessionID).Updates(map[string]interface{}{
		"chat_name":  chatInfo.ChatName,
		"updated_at": chatInfo.UpdatedAt,
	}).Error
}

// DeleteChatInfo 删除对话信息（逻辑删除）
func DeleteChatInfo(session_id string) error {
	return DB.Model(&ChatInfo{}).Where("session_id = ?", session_id).Updates(map[string]interface{}{
		"is_delete":  true,
		"updated_at": time.Now(),
	}).Error
}
