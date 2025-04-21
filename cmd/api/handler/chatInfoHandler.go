package handler

import (
	"net/http"
	"star_llm_backend_n/cmd/api/request"
	"star_llm_backend_n/cmd/api/response"
	"star_llm_backend_n/logs"
	"star_llm_backend_n/models"
	"star_llm_backend_n/services"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateChatInfo 处理创建对话信息的请求
func CreateChatInfo(ctx *gin.Context) {
	// 解析请求体
	createRequest := request.CreateChatInfoRequest{}
	err := ctx.Bind(&createRequest)
	if err != nil {
		logs.Logger.Error("[错误] 解析创建对话信息请求体失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	logs.Logger.Infof("[提取] 创建对话信息: userID=%s, sessionID=%s, name=%s",
		createRequest.UserID, createRequest.SessionID, createRequest.ChatName)

	sessionId := createRequest.SessionID
	if sessionId == "" {
		sessionId = strings.ReplaceAll(uuid.NewString(), "-", "")
	}
	// 创建对话信息对象
	chatInfo := &models.ChatInfo{
		UserID:    createRequest.UserID,
		SessionID: sessionId,
		ChatName:  createRequest.ChatName,
	}

	// 调用模型层创建对话信息
	err = models.CreateChatInfo(chatInfo)
	if err != nil {
		logs.Logger.Error("[错误] 创建对话信息失败: %s", err.Error())
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 返回成功响应
	response.MkResponse(ctx, http.StatusOK, response.Success, gin.H{
		"chat_info": chatInfo,
	})
}

// GetChatInfosByUserID 处理获取用户所有对话信息的请求
func GetChatInfos(ctx *gin.Context) {
	// 解析请求参数
	getRequest := request.GetChatInfosByUserIDRequest{}
	err := ctx.Bind(&getRequest)
	if err != nil {
		logs.Logger.Error("[错误] 解析获取用户对话信息请求参数失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	logs.Logger.Infof("[提取] 获取用户对话信息: userID=%s", getRequest.UserID)

	// 调用模型层获取用户所有对话信息
	chatInfos, err := models.GetChatInfos(getRequest.UserID, getRequest.ChatName)
	if err != nil {
		logs.Logger.Error("[错误] 获取用户对话信息失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 返回成功响应
	response.MkResponse(ctx, http.StatusOK, response.Success, gin.H{
		"chat_infos": chatInfos,
	})
}

// UpdateChatInfo 处理更新对话信息的请求
func UpdateChatInfo(ctx *gin.Context) {
	// 解析请求体
	updateRequest := request.UpdateChatInfoRequest{}
	err := ctx.Bind(&updateRequest)
	if err != nil {
		logs.Logger.Error("[错误] 解析更新对话信息请求体失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	logs.Logger.Infof("[提取] 更新对话信息: ID=%d, name=%s", updateRequest.SessionID, updateRequest.ChatName)

	// 获取现有对话信息
	chatInfo, err := models.GetChatInfoBySessionId(updateRequest.SessionID)
	if err != nil {
		logs.Logger.Error("[错误] 获取对话信息失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 更新对话信息
	chatInfo.ChatName = updateRequest.ChatName
	chatInfo.UpdatedAt = time.Now()

	// 调用模型层更新对话信息
	err = models.UpdateChatInfo(chatInfo)
	if err != nil {
		logs.Logger.Error("[错误] 更新对话信息失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 返回成功响应
	response.MkResponse(ctx, http.StatusOK, response.Success, gin.H{
		"chat_info": chatInfo,
	})
}

// DeleteChatInfo 处理删除对话信息的请求
func DeleteChatInfo(ctx *gin.Context) {
	// 解析请求体
	deleteRequest := request.DeleteChatInfoRequest{}
	err := ctx.Bind(&deleteRequest)
	if err != nil {
		logs.Logger.Error("[错误] 解析删除对话信息请求体失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	logs.Logger.Infof("[提取] 删除对话信息: ID=%s", deleteRequest.SessionID)

	// 调用服务层删除对话信息及相关消息
	err = services.DeleteChatInfo(deleteRequest.SessionID)
	if err != nil {
		logs.Logger.Error("[错误] 删除对话信息失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 返回成功响应
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

// DeleteChatInfos 批量删除对话信息的请求
func DeleteChatInfos(ctx *gin.Context) {
	// 解析请求体
	deleteRequest := request.DeleteChatInfosRequest{}
	err := ctx.Bind(&deleteRequest)
	if err != nil {
		logs.Logger.Error("[错误] 解析批量删除对话信息请求体失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	logs.Logger.Infof("[提取] 批量删除对话信息: 会话数量=%d", len(deleteRequest.SessionIDs))

	// 调用服务层批量删除对话信息及相关消息
	failedSessionIDs, err := services.DeleteChatInfos(deleteRequest.SessionIDs)

	// 检查是否有删除失败的情况
	if len(failedSessionIDs) > 0 {
		response.MkResponse(ctx, http.StatusOK, "部分会话删除失败", gin.H{
			"failed_session_ids": failedSessionIDs,
		})
		return
	}

	// 返回成功响应
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}
