package handler

import (
	"net/http"
	"star_llm_backend_n/cmd/api/request"
	"star_llm_backend_n/cmd/api/response"
	"star_llm_backend_n/dify"
	"star_llm_backend_n/logs"
	"star_llm_backend_n/models"

	"github.com/gin-gonic/gin"
)

func StopChatMessage(ctx *gin.Context) {
	// 从URL中提取task_id
	task_id := ctx.Param("task_id")
	if task_id == "" {
		logs.Logger.Error("[错误] 读取task_id失败")
		response.MkResponse(ctx, http.StatusBadRequest, "[错误] 读取task_id失败", nil)
		return
	}

	// 解析请求体
	stopRequest := request.StopRequest{}
	err := ctx.Bind(&stopRequest)
	if err != nil {
		logs.Logger.Error("[错误] 解析停止响应请求体失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	logs.Logger.Infof("[提取] 停止响应信息: user=%s, session_id=%s", stopRequest.User, stopRequest.SessionID)

	// 获取Dify客户端
	difyClient := dify.GetDifyClient(ctx)

	// 转发反馈到Dify
	err = difyClient.StopChatMessage(
		ctx.Request.Context(),
		stopRequest,
		task_id,
	)

	if err != nil {
		logs.Logger.Errorf("[错误] 转发停止响应到Dify失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, "转发停止响应到Dify失败: "+err.Error(), nil)
		return
	} else {
		logs.Logger.Info("[数据库] 更新消息: task_id=%s", task_id)
		err := models.UpdateMessageStopStatus(task_id, true)
		if err != nil {
			logs.Logger.Error("[错误]保存stop 状态到数据库失败: %v", err)
		}
	}
	response.MkResponse(ctx, http.StatusOK, "success", nil)
}
