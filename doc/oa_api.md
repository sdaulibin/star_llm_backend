# 青岛银行接入办公系统单点登录接口文档

## 目录
| 接口名称 | Url | 开发方 | 负责人 | 备注 |
| --- | --- | --- | --- | --- |
| 验证token接口 | /peimc-customization/login/verifyIMGHToken |  |  |  |

## 验证token接口
### 接口基本信息
| 项目 | 详情 |
| --- | --- |
| 接口名称 | 验证token接口 |
| 接口uri | /peimc - customization/login/verifyIMGHToken <br> 测试：https://oa.qdccb.cn:8080/ <br> 生产：https://oa.qdboffice.com:8080/ <br> 【说明】接口只能通过域名访问 <br> 测试环境域名host映射到专线：172.31.249.8 <br> 生产环境域名host映射到专线：172.31.253.10 |
| 请求方式 | POST |
| 报文传输格式 | application/json |
| 接口授权类型 |  |
| 是否需要登录 |  |
| 接口是否入参需要加签 |  |
| 接口是否入参需要加密 |  |
| 接口是否需要防重 |  |
| 接口是否出参需要加签 |  |
| 接口是否出参需要加密 |  |
| 接口版本 |  |
| 接口备注 |  |

### 请求参数
| 业务字段名称 | 类型 | 必须 | 说明 | 备注 |
| --- | --- | --- | --- | --- |
| username | String | 是 | 域账号 | QD开头的用户号通过MD5加密后的字符串 |
| token | String | 是 | 令牌 | 通过username和time进行MD5加密的32位字符串 |
| time | String | 是 | 16位时间戳 | 生成token时返回的16位系统时间 |

### 响应参数
| 业务字段名称 | 类型 | 必须 | 说明 | 备注 |
| --- | --- | --- | --- | --- |
| status | String | 是 | 响应代号 | 0 认证通过 <br> 1 认证失败 |
| msg | String | 是 | 一级响应消息 | 返回消息 |

```
curl -X POST "http://localhost:8090/sllb/api/oa/login" \
-H "Content-Type: application/json" \
-d '{
    "username": "QD24000010",
    "token": "md5_encrypted_token", 
    "time": "16_digit_timestamp"
}'
```