package handler

import (
	"net/http"
	"star_llm_backend_n/cmd/api/request"
	"star_llm_backend_n/cmd/api/response"
	"star_llm_backend_n/logs"
	"star_llm_backend_n/models"

	"github.com/gin-gonic/gin"
)

// GetMessages 获取消息列表
func GetMessages(ctx *gin.Context) {
	// 解析请求参数
	// 解析请求参数
	getRequest := request.GetMessagesRequest{}
	err := ctx.Bind(&getRequest)
	if err != nil {
		logs.Logger.Error("[错误] 解析获取用户对话信息请求参数失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	logs.Logger.Infof("[请求] 获取消息列表: user_id=%s, session_id=%s, page=%d, page_size=%d",
		getRequest.UserID, getRequest.SessionID, getRequest.Page, getRequest.PageSize)

	// 调用模型层获取消息列表
	messages, total, err := models.GetMessages(getRequest.UserID, getRequest.SessionID, getRequest.Query, getRequest.Page, getRequest.PageSize)
	if err != nil {
		logs.Logger.Errorf("[错误] 获取消息列表失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, "获取消息列表失败", nil)
		return
	}

	// 构造响应数据
	data := gin.H{
		"messages":  messages,
		"total":     total,
		"page":      getRequest.Page,
		"page_size": getRequest.PageSize,
	}

	logs.Logger.Infof("[成功] 获取消息列表成功: 共%d条记录", total)
	response.MkResponse(ctx, http.StatusOK, "获取消息列表成功", data)
}

// UpdateCollectStatus 更新消息收藏状态
func UpdateCollectStatus(ctx *gin.Context) {
	// 解析请求体
	var req request.UpdateCollectStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logs.Logger.Errorf("[错误] 解析更新收藏状态请求参数失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	logs.Logger.Infof("[请求] 更新消息收藏状态: message_id=%s, session_id=%s, is_collect=%v",
		req.MessageID, req.SessionID, req.IsCollect)

	// 调用模型层更新收藏状态
	err := models.UpdateCollectStatus(req.MessageID, req.SessionID, req.IsCollect)
	if err != nil {
		logs.Logger.Errorf("[错误] 更新消息收藏状态失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, "更新消息收藏状态失败", nil)
		return
	}

	logs.Logger.Infof("[成功] 更新消息收藏状态成功: message_id=%s, session_id=%s, is_collect=%v",
		req.MessageID, req.SessionID, req.IsCollect)
	response.MkResponse(ctx, http.StatusOK, "更新消息收藏状态成功", nil)
}

// DeleteMessage 删除消息
func DeleteMessage(ctx *gin.Context) {
	// 解析请求体
	var req request.DeleteMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logs.Logger.Errorf("[错误] 解析删除消息请求参数失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	logs.Logger.Infof("[请求] 删除消息: message_id=%s, session_id=%s", req.MessageID, req.SessionID)
	// 调用模型层删除消息
	err := models.DeleteMessage(req.MessageID, req.SessionID)
	if err != nil {
		logs.Logger.Errorf("[错误] 删除消息失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, "删除消息失败", nil)
		return
	}

	logs.Logger.Infof("[成功] 删除消息成功: message_id=%s", req.MessageID)
	response.MkResponse(ctx, http.StatusOK, "删除消息成功", nil)
}

// DeleteMessages 批量删除消息
func DeleteMessages(ctx *gin.Context) {
	// 解析请求体
	deleteRequest := request.DeleteMessagesRequest{}
	err := ctx.Bind(&deleteRequest)
	if err != nil {
		logs.Logger.Error("[错误] 解析批量删除消息请求体失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	logs.Logger.Infof("[请求] 批量删除消息: 消息数量=%d, session_id=%s", len(deleteRequest.MessageIDs), deleteRequest.SessionID)

	// 记录失败的消息ID
	failedMessageIDs := []string{}

	// 循环删除每个消息
	for _, messageID := range deleteRequest.MessageIDs {
		err := models.DeleteMessage(messageID, deleteRequest.SessionID)
		if err != nil {
			logs.Logger.Error("[错误] 删除消息ID=%s失败: %v", messageID, err)
			failedMessageIDs = append(failedMessageIDs, messageID)
		}
	}

	// 检查是否有删除失败的情况
	if len(failedMessageIDs) > 0 {
		response.MkResponse(ctx, http.StatusOK, "部分消息删除失败", gin.H{
			"failed_message_ids": failedMessageIDs,
		})
		return
	}

	// 返回成功响应
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}
