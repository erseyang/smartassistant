package maintenance

import (
	"context"
	"github.com/zhiting-tech/smartassistant/modules/api/extension"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/clouddisk"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/oauth"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
	"go.uber.org/atomic"
	"gopkg.in/oauth2.v3"
	"net/http"
	"time"
)

const START = 3
const STOP = 3
const RESTORE = 10

func HandleMaintenanceChange(event int) {
	if maintenance.connected == nil {
		maintenance.connected = atomic.NewBool(false)
	}
	if !maintenance.started {
		if event == START {
			logger.Info("enter maintenance model")
			maintenance.started = true
			maintenance.startTime = time.Now()
			if maintenance.checkStop != nil {
				close(maintenance.checkStop)
			} else {
				maintenance.checkStop = make(chan struct{})
			}
			go CheckMaintenanceDuration(maintenance.checkStop)
		}
	} else {
		if event >= STOP && event < RESTORE {
			logger.Info("exit maintenance model")
			ExitMaintenance()
		} else if event == RESTORE {
			logger.Info("restore Factory")
			ExitMaintenance()
			areas, err := entity.GetAreas()
			if err != nil {
			}
			for _, area := range areas {
				owner, err := entity.GetAreaOwner(area.ID)
				if err != nil {
					continue
				}
				req, err := http.NewRequest(http.MethodGet, "", nil)
				if err != nil {
					continue
				}
				req.Header.Add(types.GrantType, string(oauth2.ClientCredentials))
				token, err := oauth.GetSAUserToken(owner, req)
				if err != nil {
					continue
				}
				err = ProcessDelArea(area.ID, true, token)
				if err != nil {
					continue
				}
			}
		}
	}

}

func ExitMaintenance() {
	if maintenance.checkStop != nil {
		close(maintenance.checkStop)
		maintenance.checkStop = nil
	}
	maintenance = stats{}
}

func CheckMaintenanceDuration(checkStop chan struct{}) {
	tm := time.NewTicker(time.Minute)
	defer tm.Stop()
	for {
		if time.Now().Sub(maintenance.startTime) >= time.Duration(15)*time.Minute {
			ExitMaintenance()
			logger.Info("timeout auto exit maintenance model")
		} else if !maintenance.started {
			ExitMaintenance()
			logger.Info("maintenance model exited！")
			return
		}
		select {
		case <-tm.C:
			continue
		case <-checkStop:
			logger.Info("stop check maintenance duration")
			return

		}
	}
}

//ProcessDelArea 删除家庭
func ProcessDelArea(areaID uint64, isDelCloudDiskFile bool, accessToken string) (err error) {
	ctx := context.Background()
	if !extension.HasExtensionWithContext(ctx, types.CloudDisk) {
		err = clouddisk.DelAreaWithContext(ctx, areaID)
		return
	}
	_, err = clouddisk.DelCloudDiskWithContext(ctx, accessToken, isDelCloudDiskFile, areaID)

	return nil
}
