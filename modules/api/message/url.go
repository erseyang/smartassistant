package message

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/middleware"
)

// RegisterMessagesRouter 注册与消息中心相关的路由及其处理函数
func RegisterMessagesRouter(r gin.IRouter) {
	msgGroup := r.Group("/messages")
	msgGroup.Use(middleware.RequireToken)
	msgGroup.GET("", GetMessages)
	msgGroup.GET("/silence", SilenceSettingInfo)
	msgGroup.PUT("/silence", UpdateSilenceSetting)
	msgGroup.GET("/silence/devices", SilenceDevList)
	msgGroup.PUT("/silence/devices/:id", UpdateSilenceDev)
	msgGroup.PUT("/notification/setting", UpdateNotificationSetting)
	msgGroup.PUT("/notification/setting/:id", UpdateNotificationDev)
	msgGroup.GET("/notification/setting", NotificationSettingInfo)
	msgGroup.GET("/status", GetReadStatus)

	r.POST("/messages/notification", Notification)

}
