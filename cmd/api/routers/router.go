package routers

import (
	"net/http"
	"star_llm_backend_n/cmd/api/handler"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	var router *gin.Engine
	// if config.GlobalConfig.DebugMode {
	// 	gin.SetMode(gin.DebugMode)
	// 	router = gin.Default()
	// 	//visit http://0.0.0.0:9090/debug/pprof/
	// 	pprof.Register(router)
	// } else {
	// 	gin.SetMode(gin.ReleaseMode)
	// 	router = gin.Default()
	// }

	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()

	router.GET("/", func(context *gin.Context) {
		//context.String(http.StatusOK, "hello world!")
		context.Redirect(http.StatusMovedPermanently, "https://aisit.qdccb.cn:9900/chat/static")
	})

	// 使用CORS中间件
	router.Use(handler.CORSMiddleware())

	difyApi := router.Group("/sllb/api/")
	{
		difyApi.POST("files/upload", handler.FileUpload)
		difyApi.POST("messages/:message_id/feedbacks", handler.FeedBack)
		difyApi.POST("chat-messages/:task_id/stop", handler.StopChatMessage)
		difyApi.POST("chat-messages", handler.ChatMessage)
		difyApi.GET("messages/:message_id/suggested", handler.Suggested)
		difyApi.POST("chat-info/create", handler.CreateChatInfo)
		difyApi.POST("chat-info/get", handler.GetChatInfos)
		difyApi.POST("chat-info/update", handler.UpdateChatInfo)
		difyApi.POST("chat-info/delete", handler.DeleteChatInfo)
		// 消息管理相关接口
		difyApi.POST("chat-messages/get", handler.GetMessages)
		difyApi.POST("chat-messages/collect", handler.UpdateCollectStatus)
		difyApi.POST("chat-messages/delete", handler.DeleteMessage)
		// OA系统单点登录接口
		difyApi.POST("oa/login", handler.VerifyOAToken)
	}
	return router
}
