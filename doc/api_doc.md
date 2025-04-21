# Star LLM Backend API 文档

本文档提供了 Star LLM Backend 的 API 接口说明，包括请求方法、URL、参数和响应格式。

## 基础信息

- 基础路径: `/sllb/api/`
- 所有接口均使用 JSON 格式进行数据交换
- 所有请求需要遵循 CORS 策略
- 以下示例中使用 `http://localhost:8080` 作为基础URL，实际使用时请替换为真实服务器地址

## API 接口列表

### 1. 文件上传

**请求方法**: POST

**URL**: `/sllb/api/files/upload`

**处理函数**: `handler.FileUpload`

**请求参数**:

- `user` (表单字段): 用户标识，必填
- `file` (表单文件): 要上传的文件，支持的格式包括:
  - 文本和文档格式: .txt, .markdown, .mdx, .pdf, .html, .xlsx, .xls, .docx, .csv, .md, .htm 等

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/files/upload \
  -F "user=user123" \
  -F "file=@/path/to/document.pdf"
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": {
    "file_id": "文件ID",
    "file_name": "文件名称"
  }
}
```

### 2. 消息反馈

**请求方法**: POST

**URL**: `/sllb/api/messages/:message_id/feedbacks`

**处理函数**: `handler.FeedBack`

**URL 参数**:
- `message_id`: 消息ID

**请求体**:
```json
{
  "rating": "like或dislike",
  "user": "用户ID",
  "content": "反馈内容",
  "session_id": "会话ID"
}
```

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/messages/msg123/feedbacks \
  -H "Content-Type: application/json" \
  -d '{
    "rating": "like",
    "user": "user123",
    "content": "非常有帮助的回答",
    "session_id": "session456"
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": null
}
```

### 3. 停止聊天消息

**请求方法**: POST

**URL**: `/sllb/api/chat-messages/:task_id/stop`

**处理函数**: `handler.StopChatMessage`

**URL 参数**:
- `task_id`: 任务ID

**请求体**:
```json
{
  "user": "用户ID",
  "session_id": "会话ID"
}
```

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-messages/task123/stop \
  -H "Content-Type: application/json" \
  -d '{
    "user": "user123",
    "session_id": "session456"
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": null
}
```

### 4. 发送聊天消息

**请求方法**: POST

**URL**: `/sllb/api/chat-messages`

**处理函数**: `handler.ChatMessage`

**请求体**:
```json
{
  "inputs": {},
  "query": "用户问题",
  "response_mode": "blocking或streaming",
  "conversation_id": "对话ID",
  "user": "用户ID",
  "files": [
    {
      "type": "document",
      "transfer_method": "local_file",
      "upload_file_id": "上传文件ID"
    }
  ],
  "session_id": "会话ID"
}
```

**curl调用示例**:

阻塞模式:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-messages \
  -H "Content-Type: application/json" \
  -d '{
    "inputs": {},
    "query": "如何使用这个API？",
    "response_mode": "blocking",
    "conversation_id": "conv123",
    "user": "user123",
    "files": [
      {
        "type": "document",
        "transfer_method": "local_file",
        "upload_file_id": "file123"
      }
    ],
    "session_id": "session456"
  }'
```

流式模式:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-messages \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{
    "inputs": {},
    "query": "如何使用这个API？",
    "response_mode": "streaming",
    "conversation_id": "conv123",
    "user": "user123",
    "files": [],
    "session_id": "session456"
  }'
```

**响应格式**:

阻塞模式响应:
```json
{
  "code": 200,
  "result": "success",
  "data": {
    "task_id": "任务ID",
    "message_id": "消息ID",
    "conversation_id": "对话ID",
    "answer": "回答内容",
    "metadata": {}
  }
}
```

流式模式响应: 使用 Server-Sent Events (SSE) 格式返回数据流

### 5. 获取建议问题

**请求方法**: GET

**URL**: `/sllb/api/messages/:message_id/suggested`

**处理函数**: `handler.Suggested`

**URL 参数**:
- `message_id`: 消息ID

**查询参数**:
- `user`: 用户ID

**curl调用示例**:
```bash
curl -X GET "http://localhost:8080/sllb/api/messages/msg123/suggested?user=user123"
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": ["建议问题1", "建议问题2"]
}
```

### 6. 创建聊天信息

**请求方法**: POST

**URL**: `/sllb/api/chat-info/create`

**处理函数**: `handler.CreateChatInfo`

**请求体**:
```json
{
  "user_id": "用户ID",
  "session_id": "会话ID",
  "chat_name": "对话名称"
}
```

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-info/create \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "session_id": "session456",
    "chat_name": "API使用讨论"
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": {
    "chat_info": {
      "user_id": "用户ID",
      "session_id": "会话ID",
      "chat_name": "对话名称",
      "created_at": "创建时间",
      "updated_at": "更新时间"
    }
  }
}
```

### 7. 获取用户的所有对话信息

**请求方法**: POST

**URL**: `/sllb/api/chat-info/get`

**处理函数**: `handler.GetChatInfos`

**请求体**:
```json
{
  "user_id": "用户ID"
}
```

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-info/get \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123"
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": {
    "chat_infos": [
      {
        "user_id": "用户ID",
        "session_id": "会话ID",
        "chat_name": "对话名称",
        "created_at": "创建时间",
        "updated_at": "更新时间"
      },
      // 更多对话信息...
    ]
  }
}
```

### 8. 更新对话信息

**请求方法**: POST

**URL**: `/sllb/api/chat-info/update`

**处理函数**: `handler.UpdateChatInfo`

**请求体**:
```json
{
  "session_id": "会话ID",
  "chat_name": "新的对话名称"
}
```

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-info/update \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session456",
    "chat_name": "API高级用法讨论"
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": {
    "chat_info": {
      "user_id": "用户ID",
      "session_id": "会话ID",
      "chat_name": "新的对话名称",
      "created_at": "创建时间",
      "updated_at": "更新时间"
    }
  }
}
```

### 9. 删除对话信息

**请求方法**: POST

**URL**: `/sllb/api/chat-info/delete`

**处理函数**: `handler.DeleteChatInfo`

**请求体**:
```json
{
  "session_id": "会话ID"
}
```

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-info/delete \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session456"
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": null
}
```

### 10. 获取消息列表

**请求方法**: POST

**URL**: `/sllb/api/chat-messages/get`

**处理函数**: `handler.GetMessages`

**请求体**:
```json
{
  "user_id": "用户ID",
  "session_id": "会话ID",
  "page": 1,
  "page_size": 10
}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| user_id | string | 否 | 用户ID，不传则不按用户ID过滤 |
| session_id | string | 是 | 会话ID |
| page | int | 是 | 页码，从1开始 |
| page_size | int | 是 | 每页记录数，最大100 |

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-messages/get \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "session_id": "session456",
    "page": 1,
    "page_size": 10
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": {
    "messages": [
      {
        "ID": 1,
        "UserID": "user123",
        "SessionID": "session456",
        "MessageID": "msg789",
        "ConversationID": "conv123",
        "Query": "用户问题",
        "Answer": "AI回答",
        "CreatedAt": "2023-01-01T12:00:00Z",
        "UpdatedAt": "2023-01-01T12:01:00Z",
        "IsSafe": true,
        "IsLike": false,
        "IsCollect": true,
        "CurrentID": "current123",
        "IsStop": false,
        "FileID": "file123",
        "TaskID": "task123"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 11. 更新消息收藏状态

**请求方法**: POST

**URL**: `/sllb/api/chat-messages/collect`

**处理函数**: `handler.UpdateCollectStatus`

**请求体**:
```json
{
  "message_id": "msg789",
  "session_id": "session456",
  "is_collect": true
}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| message_id | string | 是 | 消息ID |
| session_id | string | 是 | 会话ID |
| is_collect | boolean | 是 | 收藏状态，true为收藏，false为取消收藏 |

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-messages/collect \
  -H "Content-Type: application/json" \
  -d '{
    "message_id": "msg789",
    "session_id": "session456",
    "is_collect": true
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": null
}
```

### 12. 删除消息

**请求方法**: POST

**URL**: `/sllb/api/chat-messages/delete`

**处理函数**: `handler.DeleteMessage`

**请求体**:
```json
{
  "message_id": "msg789",
  "session_id": "session456"
}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| message_id | string | 是 | 消息ID |
| session_id | string | 是 | 会话ID |

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-messages/delete \
  -H "Content-Type: application/json" \
  -d '{
    "message_id": "msg789",
    "session_id": "session456"
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": null
}
```

### 13. 批量删除消息

**请求方法**: POST

**URL**: `/sllb/api/chat-messagess/delete`

**处理函数**: `handler.DeleteMessages`

**请求体**:
```json
{
  "message_ids": ["msg1", "msg2"],
  "session_id": "session456"
}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| message_ids | array | 是 | 消息ID数组 |
| session_id | string | 是 | 会话ID |

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/chat-messagess/delete \
  -H "Content-Type: application/json" \
  -d '{
    "message_ids": ["msg1", "msg2"],
    "session_id": "session456"
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": null
}


### 13. OA系统单点登录

**请求方法**: POST

**URL**: `/sllb/api/oa/login`

**处理函数**: `handler.VerifyOAToken`

**请求体**:
```json
{
  "token": "OA系统令牌"
}
```

**curl调用示例**:
```bash
curl -X POST http://localhost:8080/sllb/api/oa/login \
  -H "Content-Type: application/json" \
  -d '{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": {
    "user_info": {
      "user_id": "用户ID",
      "user_name": "用户名称"
    }
  }
}
```

## 错误码说明

- 200: 成功 (result: "success")
- 400: 请求参数错误 (result: "param_invalid")
- 500: 服务器内部错误 (result: 错误信息)

## 通用响应格式

所有API响应均遵循以下格式:

```json
{
  "code": 状态码,
  "result": "状态描述",
  "data": 响应数据对象或null
}
```