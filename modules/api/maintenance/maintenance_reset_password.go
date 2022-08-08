package maintenance

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/maintenance"
)

type ResetProPasswordReq struct {
	IsReset bool `json:"is_reset"`
}

type ResetProPasswordResp struct {
	IsReset           bool  `json:"is_reset"`
	MaintainStartTime int64 `json:"maintenance_start_time"`
}

func resetProPassword(c *gin.Context) {
	var (
		err  error
		req  ResetProPasswordReq
		resp ResetProPasswordResp
	)
	defer func() {
		t, exists := c.Get(TimeKey)
		v, ok := t.(int64)
		if exists && ok {
			v = maintenance.RefreshMaintenanceStarTime()
		}
		resp.MaintainStartTime = v
		response.HandleResponse(c, err, &resp)
	}()
	reset := maintenance.ResetProPassword(req.IsReset)
	resp.IsReset = reset

}
