package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"star_llm_backend_n/cmd/api/request"
	"star_llm_backend_n/cmd/api/response"
	"star_llm_backend_n/dify"
	"star_llm_backend_n/logs"
	"star_llm_backend_n/models"
	"star_llm_backend_n/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ChatMessage(ctx *gin.Context) {
	// 解析请求体
	chatMessageRequest := request.ChatMessageRequest{}
	err := ctx.Bind(&chatMessageRequest)
	if err != nil {
		logs.Logger.Error("[错误] 解析发送对话请求体失败: %v", err)
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	logs.Logger.Infof("[提取] 发送对话信息: query=%s", chatMessageRequest.Query)

	current_id := uuid.NewString()
	user_id := chatMessageRequest.User
	query := chatMessageRequest.Query
	session_id := chatMessageRequest.SessionID
	conversation_id := uuid.Nil.String()
	if len(chatMessageRequest.ConversationID) > 0 {
		conversation_id = chatMessageRequest.ConversationID
	}
	log.Printf("[提取] 用户ID: %s, 输入: %s, 会话ID: %s, 对话ID: %s", user_id, query, session_id, conversation_id)

	chatInfo, err := models.GetChatInfoBySessionId(session_id)
	if err != nil {
		logs.Logger.Error("[错误] 获取对话信息失败: %v", err)
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	} else {
		if chatInfo == nil {
			chatName := query
			if len(chatName) > 10 {
				chatName = chatName[:10]
			}
			models.CreateChatInfo(&models.ChatInfo{
				UserID:    user_id,
				SessionID: session_id,
				ChatName:  chatName,
			})
		}
	}

	// 提取文件信息
	file_id := uuid.Nil.String()
	if len(chatMessageRequest.Files) > 0 {
		for _, fileInfo := range chatMessageRequest.Files {
			if fileInfo.Type == "document" && fileInfo.TransferMethod == "local_file" {
				log.Printf("[提取] 文件上传ID: %s, 类型: %s, 传输方式: %s",
					fileInfo.UploadFileID, fileInfo.Type, fileInfo.TransferMethod)
				file_id = fileInfo.UploadFileID
			}
		}
	}
	// 检查是否需要保存消息到数据库
	if query != "" {
		err := services.SaveMessageToDB(current_id, session_id, query, "", user_id, conversation_id, uuid.Nil.String(), file_id)
		if err != nil {
			logs.Logger.Error("[错误] 保存对话消息到数据库失败: %v", err)
		}
	}

	// 获取Dify客户端
	difyClient := dify.GetDifyClient(ctx)

	// 设置响应头，支持流式传输
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
	ctx.Writer.Flush()

	// 创建一个字符串构建器来累积响应
	var fullAnswer strings.Builder
	var res_message_id, res_conversation_id, res_task_id string

	// 发送流式请求
	difyClient.SendChatMessageStream(context.Background(), chatMessageRequest, func(chunk response.StreamChunk) error {
		// 处理每个数据块
		//logs.Logger.Infof("事件: %s\n", chunk.Event)

		if chunk.Event == "message" {
			// 累积答案
			fullAnswer.WriteString(chunk.Answer)

			// 保存消息ID和对话ID
			if res_message_id == "" && chunk.MessageID != "" {
				res_message_id = chunk.MessageID
			}
			if res_conversation_id == "" && chunk.ConversationID != "" {
				res_conversation_id = chunk.ConversationID
			}
			if res_task_id == "" && chunk.TaskID != "" {
				res_task_id = chunk.TaskID
			}

			data, err := json.Marshal(chunk)
			if err != nil {
				logs.Logger.Errorf("序列化chunk失败: %v", err)
			} else {
				// 发送数据到前端
				ctx.Writer.Write([]byte("data: " + string(data) + "\n\n"))
				ctx.Writer.Flush()
				logs.Logger.Infof("收到部分回答: %s\n", string(data))
			}
		} else if chunk.Event == "message_end" {
			logs.Logger.Infof("消息结束，完整回答:\n%s\n", fullAnswer.String())
			logs.Logger.Infof("消息ID: %s\n", res_message_id)
			logs.Logger.Infof("对话ID: %s\n", res_conversation_id)
			logs.Logger.Infof("任务ID: %s\n", res_task_id)

			// 发送结束事件到前端
			endData, _ := json.Marshal(chunk)
			ctx.Writer.Write([]byte("data: " + string(endData) + "\n\n"))
			ctx.Writer.Write([]byte("event: end\ndata: \n\n"))
			ctx.Writer.Flush()

			// 在第一条消息时就更新数据库
			if res_conversation_id != "" && res_message_id != "" && res_task_id != "" {
				err := services.UpdateMessageToDB(current_id, fullAnswer.String(), res_conversation_id, res_message_id, res_task_id)
				if err != nil {
					log.Printf("[错误] 保存消息到数据库失败: %v", err)
				}
			}
		}

		return nil
	})
	//response.MkResponse(ctx, http.StatusOK, "success", nil)
}
