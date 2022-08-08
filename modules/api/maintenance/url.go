package maintenance

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/middleware"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/maintenance"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"strconv"
)

const TimeKey = "startTime"

func RegisterMaintenanceRouter(r gin.IRouter) {
	maintenanceGroup := r.Group("maintenance")
	maintenanceGroup.GET("/:id/token", requireMaintainModeStarted, connectAndRefreshToken)
	maintenanceGroup.DELETE("/exit", requireAccessToken, exitMaintainMode)
	maintenanceGroup.PUT("/:id/owner", requireAccessToken, requireSameArea, transferOwner)
	maintenanceGroup.POST("/password/reset", requireAccessToken, resetProPassword)
}

//requireAccessToken 维护模式需要校验accessToken、连接用户userId是否为连接用户Id
func requireAccessToken(c *gin.Context) {
	_, started := maintenance.CheckStatedAndConnected()
	if !started {
		err := errors.New(status.MaintenanceNotRun)
		response.HandleResponse(c, err, nil)
		c.Abort()
		return
	}
	token := c.Request.Header.Get(types.SATokenKey)
	queryToken := c.Query("token")
	if token == "" && queryToken != "" {
		c.Request.Header.Add(types.SATokenKey, queryToken)
	}
	ti, err := middleware.VerifyAccessToken(c)
	if err != nil {
		response.HandleResponse(c, err, nil)
		c.Abort()
		return
	}
	id := maintenance.GetMaintenanceUserId()
	if id != ti.GetUserID() {
		err = errors.New(status.GetUserTokenAuthDeny)
		response.HandleResponse(c, err, nil)
		c.Abort()
		return
	}
	c.Set(TimeKey, maintenance.RefreshMaintenanceStarTime())
}

//requireMaintainModeStarted 检查是否在运行维护模式
func requireMaintainModeStarted(c *gin.Context) {
	_, started := maintenance.CheckStatedAndConnected()
	if !started {
		err := errors.New(status.MaintenanceNotRun)
		response.HandleResponse(c, err, nil)
		c.Abort()
		return
	}
}

// requireSameArea  请求api需要在同一个家庭下
func requireSameArea(c *gin.Context) {

	u := session.Get(c)

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.HandleResponse(c, errors.Wrap(err, errors.BadRequest), nil)
		c.Abort()
		return
	}

	user, err := entity.GetUserByID(userID)
	if err != nil {
		response.HandleResponse(c, errors.Wrap(err, errors.InternalServerErr), nil)
		c.Abort()
		return
	}

	if u.AreaID != user.AreaID {
		response.HandleResponse(c, errors.New(status.Deny), nil)
		c.Abort()
	} else {
		c.Next()
	}

}
