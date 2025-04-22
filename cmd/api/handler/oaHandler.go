package handler

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"star_llm_backend_n/cmd/api/request"
	"star_llm_backend_n/cmd/api/response"
	"star_llm_backend_n/config"
	"star_llm_backend_n/logs"
	"strings"

	"github.com/gin-gonic/gin"
)

// VerifyOAToken 验证OA系统的token
func VerifyOAToken(ctx *gin.Context) {
	var req request.OATokenVerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 对username进行MD5加密并转换为大写
	hash := md5.Sum([]byte(req.Username))
	req.Username = hex.EncodeToString(hash[:])
	req.Username = strings.ToUpper(req.Username)

	request, err := json.Marshal(req)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("POST", config.GlobalConfig.OA.Url, bytes.NewBuffer(request))
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	logs.Logger.Infof("发送OA系统的token认证请求: %s\n", string(request))
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
