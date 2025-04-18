package handler

import (
	"log"
	"net/http"
	"star_llm_backend_n/cmd/api/request"
	"star_llm_backend_n/cmd/api/response"
	"star_llm_backend_n/dify"
	"star_llm_backend_n/logs"
	"star_llm_backend_n/models"

	"github.com/gin-gonic/gin"
)

func FeedBack(ctx *gin.Context) {
	// 从URL中提取message_id
	message_id := ctx.Param("message_id")
	if message_id == "" {
		logs.Logger.Error("[错误] 读取message_id失败")
		response.MkResponse(ctx, http.StatusBadRequest, "[错误] 读取message_id失败", nil)
		return
	}

	// 解析请求体
	feedbackRequest := request.FeedBackRequest{}
	err := ctx.Bind(&feedbackRequest)
	if err != nil {
		logs.Logger.Error("[错误] 解析反馈请求体失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	logs.Logger.Infof("[提取] 反馈信息: rating=%s, user=%s, session_id=%s",
		feedbackRequest.Rating, feedbackRequest.User, feedbackRequest.SessionID)

	// 根据rating更新is_like字段
	// 当rating为"like"时，is_like为true；否则为false
	isLike := feedbackRequest.Rating == "like"
	err = models.UpdateMessageLikeStatus(message_id, feedbackRequest.SessionID, isLike)
	if err != nil {
		log.Printf("[错误] 更新消息点赞状态失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	logs.Logger.Infof("[成功] 已更新消息点赞状态: message_id=%s, session_id=%s, is_like=%v",
		message_id, feedbackRequest.SessionID, isLike)

	// 获取Dify客户端
	difyClient := dify.GetDifyClient(ctx)

	// 转发反馈到Dify
	err = difyClient.FeedbackMessage(ctx.Request.Context(), message_id, feedbackRequest)

	if err != nil {
		logs.Logger.Errorf("[错误] 转发反馈到Dify失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, "转发反馈到Dify失败: "+err.Error(), nil)
		return
	}

	response.MkResponse(ctx, http.StatusOK, "success", nil)
}
