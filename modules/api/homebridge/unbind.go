package homebridge

import (
	errors2 "errors"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/zhiting-tech/smartassistant/modules/api/extension"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"net/http"
)

func Unbind(c *gin.Context) {
	var err error
	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	if err = UnbindHomeBridge(c); err != nil {
		return
	}
}

func UnbindHomeBridge(c *gin.Context) (err error) {
	if !extension.HasExtensionWithContext(c.Request.Context(), types.HomeBridge) {
		return
	}

	b, err := reqToHomeBridge(c, http.MethodDelete, getHomeBridgeApi("unbind"), nil)
	if err != nil {
		err = errors.Wrap(err, status.UnbindError)
		return
	}

	code := gjson.GetBytes(b, "status").Int()
	if code != 0 {
		e := errors2.New(gjson.GetBytes(b, "reason").String())
		err = errors.Wrap(e, status.UnbindError)
		return
	}
	return
}
