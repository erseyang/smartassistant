package maintenance

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/maintenance"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
)

func exitMaintainMode(c *gin.Context) {
	var (
		err error
	)
	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	user := session.Get(c)
	if user == nil {
		err = errors.Wrap(err, status.InvalidUserCredentials)
		return
	}

	logger.Infof("user %d exit maintenance model", user.UserID)
	maintenance.ExitMaintenance()
}
