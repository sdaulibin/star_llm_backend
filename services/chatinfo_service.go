package services

import (
	"star_llm_backend_n/logs"
	"star_llm_backend_n/models"

	"gorm.io/gorm"
)

// DeleteChatInfo 删除单个对话信息及其相关消息
func DeleteChatInfo(sessionID string) error {
	// 使用事务确保数据一致性
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除对话相关的所有消息
		if err := models.DeleteMessagesBySessionIDWithTx(tx, sessionID); err != nil {
			logs.Logger.Error("[错误] 删除会话ID=%s的消息失败: %v", sessionID, err)
			return err
		}

		// 2. 删除对话信息
		if err := models.DeleteChatInfoWithTx(tx, sessionID); err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteChatInfos 批量删除对话信息及其相关消息
func DeleteChatInfos(sessionIDs []string) ([]string, error) {
	// 记录失败的会话ID
	failedSessionIDs := []string{}

	// 循环删除每个会话及其相关消息
	for _, sessionID := range sessionIDs {
		// 使用事务确保数据一致性
		err := models.DB.Transaction(func(tx *gorm.DB) error {
			// 1. 删除对话相关的所有消息
			if err := models.DeleteMessagesBySessionIDWithTx(tx, sessionID); err != nil {
				logs.Logger.Error("[错误] 删除会话ID=%s的消息失败: %v", sessionID, err)
				return err
			}

			// 2. 删除对话信息
			if err := models.DeleteChatInfoWithTx(tx, sessionID); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			logs.Logger.Error("[错误] 删除会话ID=%s失败: %v", sessionID, err)
			failedSessionIDs = append(failedSessionIDs, sessionID)
		}
	}

	return failedSessionIDs, nil
}
