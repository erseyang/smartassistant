package homebridge

import (
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"net/http"
)

type GetAuthorizationResp struct {
	IsAuthorized bool `json:"is_authorized"`
}

func GetAuthorization(c *gin.Context) {
	var (
		err  error
		resp GetAuthorizationResp
	)
	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	// 请求homebridge获取授权状态
	b, err := reqToHomeBridge(c, http.MethodGet, getHomeBridgeApi("authorization"), nil)
	if err != nil {
		err = errors.Wrap(err, status.HomeBridgeServiceError)
		return
	}

	resp.IsAuthorized = gjson.GetBytes(b, "data").Get("is_authorized").Bool()
	return

}
