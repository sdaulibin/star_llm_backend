package request

// OATokenVerifyRequest 定义验证token的请求结构
type OATokenVerifyRequest struct {
	Username string `json:"username" binding:"required"` // 域账号，QD开头的用户号通过MD5加密后的字符串
	Token    string `json:"token" binding:"required"`    // 令牌，通过username和time进行MD5加密的32位字符串
	Time     string `json:"time" binding:"required"`     // 16位时间戳，生成token时返回的16位系统时间
}
