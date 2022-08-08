package maintenance

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/oauth"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/maintenance"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"strconv"
)

type MaintainTokenReq struct {
	UserId int `uri:"id"  binding:"required"`
}

type MaintainTokenRes struct {
	MaintainStartTime int64  `json:"maintenance_start_time"`
	AccessToken       string `json:"access_token"`
}

func connectAndRefreshToken(c *gin.Context) {
	var (
		req  MaintainTokenReq
		resp MaintainTokenRes
		err  error
	)
	defer func() {
		response.HandleResponse(c, err, &resp)
	}()

	if err = c.ShouldBindUri(&req); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}
	token := c.Request.Header.Get(types.SATokenKey)
	queryToken := c.Query("token")
	if token == "" && queryToken != "" {
		token = queryToken
	}

	ti, err := oauth.GetOauthServer().Manager.LoadAccessToken(token)

	if err != nil || ti.GetUserID() != strconv.Itoa(req.UserId) {
		user, err2 := entity.GetUserByID(req.UserId)
		if err2 != nil {
			err = errors.Wrap(err2, errors.BadRequest)
			return
		}
		token, err2 := oauth.GetSAUserToken(user, c.Request)
		if err2 != nil {
			err = errors.Wrap(err2, errors.BadRequest)
			return
		}
		resp.AccessToken = token
	}

	connected := maintenance.ConnectMaintenanceMode()
	if !connected {
		id := maintenance.GetMaintenanceUserId()
		if id != strconv.Itoa(req.UserId) {
			resp.AccessToken = ""
			err = errors.New(status.MaintenanceConnected)
			return
		}
	}

	maintenance.SetConnectUserId(strconv.Itoa(req.UserId))
	resp.MaintainStartTime = maintenance.RefreshMaintenanceStarTime()

}
