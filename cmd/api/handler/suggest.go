package handler

import (
	"net/http"
	"star_llm_backend_n/cmd/api/response"
	"star_llm_backend_n/dify"
	"star_llm_backend_n/logs"

	"github.com/gin-gonic/gin"
)

func Suggested(ctx *gin.Context) {
	// 从URL中提取message_id
	message_id := ctx.Param("message_id")
	if message_id == "" {
		logs.Logger.Error("[错误] 读取message_id失败")
		response.MkResponse(ctx, http.StatusBadRequest, "[错误] 读取message_id失败", nil)
		return
	}

	user_id := ctx.Query("user")
	if user_id == "" {
		logs.Logger.Error("[错误] 读取user_id失败")
		response.MkResponse(ctx, http.StatusBadRequest, "[错误] 读取user_id失败", nil)
		return
	}

	// 获取Dify客户端
	difyClient := dify.GetDifyClient(ctx)

	// 转发反馈到Dify
	suggestions, err := difyClient.Suggested(
		ctx.Request.Context(),
		message_id,
		user_id,
	)
	if err != nil {
		logs.Logger.Errorf("[错误] 转发获取下一轮列表到Dify失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, "转发获取下一轮列表到Dify失败: "+err.Error(), nil)
		return
	}

	// 打印建议问题
	logs.Logger.Info("建议的后续问题:\n")
	for i, suggestion := range suggestions {
		logs.Logger.Info("%d. %s\n", i+1, suggestion)
	}
	response.MkResponse(ctx, http.StatusOK, "success", suggestions)
}
