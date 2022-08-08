package homebridge

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/middleware"
)

func InitHomeBridgeRouter(r gin.IRouter) {
	hbGroup := r.Group("homebridge", middleware.RequireAccount, middleware.RequireOwner)

	hbGroup.POST("pin", GetPin)
	hbGroup.GET("authorization", GetAuthorization)
	hbGroup.DELETE("unbind", Unbind)

}
