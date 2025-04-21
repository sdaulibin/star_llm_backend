package models

import (
	"errors"
	"star_llm_backend_n/logs"
	"time"

	"gorm.io/gorm"
)

// ChatInfo 表示数据库中的对话信息记录
type ChatInfo struct {
	ID             int       `gorm:"primaryKey"`
	UserID         string    `gorm:"size:10"`
	SessionID      string    `gorm:"type:uuid;not null"`
	ConversationID string    `gorm:"type:uuid"`
	ChatName       string    `gorm:"size:255;not null"`
	IsDelete       bool      `gorm:"default:false;not null"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (ChatInfo) TableName() string {
	return "chat_info" // 指定为单数形式
}

// CreateChatInfo 创建新的对话信息
func CreateChatInfo(chatInfo *ChatInfo) error {
	logs.Logger.Infof("[创建] 创建对话信息: session_id=%s, chat_info=%v\n", chatInfo.SessionID, chatInfo)
	// 确保ChatName字段不为空
	if chatInfo.ChatName == "" {
		logs.Logger.Warn("[警告] 创建对话信息时ChatName为空")
		chatInfo.ChatName = "新对话"
	}

	// 创建记录
	err := DB.Create(chatInfo).Error

	// 验证创建后的记录
	if err == nil {
		var savedChat ChatInfo
		DB.Where("session_id = ?", chatInfo.SessionID).First(&savedChat)
		logs.Logger.Infof("[创建] 保存后的对话信息: session_id=%s, chat_name=%s\n", savedChat.SessionID, savedChat.ChatName)
	}

	return err
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
func GetChatInfos(userID, chatName string) ([]ChatInfo, error) {
	var chatInfos []ChatInfo
	db := DB.Where("user_id = ? AND is_delete = ?", userID, false)

	if chatName != "" {
		db = db.Where("chat_name like ?", "%"+chatName+"%")
	}
	if err := db.Order("updated_at DESC").Find(&chatInfos).Error; err != nil {
		return nil, err
	}
	return chatInfos, nil
}

// UpdateChatInfo 更新对话信息
func UpdateChatInfo(chatInfo *ChatInfo) error {
	logs.Logger.Infof("[更新] 更新对话信息: session_id=%s, chat_info=%v\n", chatInfo.SessionID, chatInfo)
	return DB.Model(&ChatInfo{}).Where("session_id = ?", chatInfo.SessionID).Updates(map[string]interface{}{
		"chat_name":       chatInfo.ChatName,
		"updated_at":      chatInfo.UpdatedAt,
		"conversation_id": chatInfo.ConversationID,
	}).Error
}

// DeleteChatInfo 删除对话信息（逻辑删除）
func DeleteChatInfo(session_id string) error {
	return DeleteChatInfoWithTx(DB, session_id)
}

// DeleteChatInfoWithTx 支持事务的删除聊天信息方法
func DeleteChatInfoWithTx(tx *gorm.DB, session_id string) error {
	return tx.Model(&ChatInfo{}).Where("session_id = ?", session_id).Updates(map[string]interface{}{
		"is_delete":  true,
		"updated_at": time.Now(),
	}).Error
}
