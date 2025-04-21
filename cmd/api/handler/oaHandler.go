package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"star_llm_backend_n/cmd/api/request"
	"star_llm_backend_n/cmd/api/response"
	"star_llm_backend_n/logs"

	"github.com/gin-gonic/gin"
)

// VerifyOAToken 验证OA系统的token
func VerifyOAToken(ctx *gin.Context) {
	var req request.OATokenVerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	request, err := json.Marshal(req)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	logs.Logger.Infof("OA系统的token认证请求: %s\n", string(request))
	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("POST", "http://oa.qdccb.cn:8080/peimc-customization/login/verifyIMGHToken", bytes.NewBuffer(request))
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		logs.Logger.Errorf("发送 HTTP 请求失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, "发送 HTTP 请求失败", nil)
		return
	}
	defer resp.Body.Close()
	// 解析响应
	var oaResp response.OATokenVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&oaResp); err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	logs.Logger.Infof("OA系统返回的token认证结果: %s<>%s\n", oaResp.Status, oaResp.Msg)
	if oaResp.Status == "0" {
		response.MkResponse(ctx, http.StatusOK, "success", nil)
	} else {
		response.MkResponse(ctx, http.StatusInternalServerError, oaResp.Msg, nil)
	}
}
