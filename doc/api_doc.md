# Star LLM Backend API 文档

本文档提供了 Star LLM Backend 的 API 接口说明，包括请求方法、URL、参数和响应格式。

## 基础信息

- 基础路径: `/sllb/api/`
- 所有接口均使用 JSON 格式进行数据交换
- 所有请求需要遵循 CORS 策略

## API 接口列表

### 1. 文件上传

**请求方法**: POST

**URL**: `/sllb/api/files/upload`

**处理函数**: `handler.FileUpload`

**请求参数**:

- `user` (表单字段): 用户标识，必填
- `file` (表单文件): 要上传的文件，支持的格式包括:
  - 文本和文档格式: .txt, .markdown, .mdx, .pdf 等

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

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": {
    "questions": ["建议问题1", "建议问题2"]
  }
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
  "name": "对话名称"
}
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
}
```

```
# 创建聊天信息
curl -X POST http://localhost:8090/sllb/api/chat-info/create \
-H "Content-Type: application/json" \
-d '{
  "session_id": "",
  "chat_name": "新对话",
  "user_id": "QD24000010"
}'
```

### 7. 获取用户的所有对话信息

**请求方法**: GET

**URL**: `/sllb/api/chat-info/get`

**处理函数**: `handler.GetChatInfosByUserID`

**URL 参数**:
- `user_id`: 用户ID (必填)

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": {
    "chat_infos": [
      {
        "id": 1,
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
```
# 获取用户聊天信息列表
curl -X POST "http://localhost:8090/sllb/api/chat-info/get" \
-H "Content-Type: application/json" \
-d '{
    "user_id": "QD24000010"
}'
```

### 8. 更新对话信息

**请求方法**: POST

**URL**: `/sllb/api/chat-info/update`

**处理函数**: `handler.UpdateChatInfo`

**请求体**:
```json
{
  "id": 1,
  "name": "新的对话名称"
}
```

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": {
    "chat_info": {
      "id": 1,
      "user_id": "用户ID",
      "session_id": "会话ID",
      "chat_name": "新的对话名称",
      "created_at": "创建时间",
      "updated_at": "更新时间"
    }
  }
}
```
```
curl -X POST "http://localhost:8090/sllb/api/chat-info/update" \
-H "Content-Type: application/json" \
-d '{
    "id": 1,
    "chat_name": "新的聊天会话名称"
}'
```
### 9. 删除对话信息

**请求方法**: POST

**URL**: `/sllb/api/chat-info/delete`

**处理函数**: `handler.DeleteChatInfo`

**请求体**:
```json
{
  "id": 1
}
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
```
curl -X POST "http://localhost:8090/sllb/api/chat-info/delete" \
-H "Content-Type: application/json" \
-d '{
    "id": 2
}'
```

### 10. 获取消息列表

**请求方法**: GET

**URL**: `/sllb/api/messages`

**处理函数**: 未指定

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| user_id | string | 否 | 用户ID，不传则不按用户ID过滤 |
| session_id | string | 是 | 会话ID |
| page | int | 是 | 页码，从1开始 |
| page_size | int | 是 | 每页记录数，最大100 |

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

**URL**: `/sllb/api/messages/collect`

**处理函数**: 未指定

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

**处理函数**: 未指定

**请求体**:
```json
{
  "id": 1
}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| id | int | 是 | 消息ID |

**响应格式**:

成功响应:
```json
{
  "code": 200,
  "result": "success",
  "data": null
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