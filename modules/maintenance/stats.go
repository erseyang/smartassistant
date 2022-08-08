package maintenance

import (
	"github.com/zhiting-tech/smartassistant/modules/types"
	"go.uber.org/atomic"
	"time"
)

var maintenance stats

type stats struct {
	started          bool
	startTime        time.Time
	checkStop        chan struct{}
	connectUserId    string
	connected        *atomic.Bool
	resetProPassword bool
}

func CheckStatedAndConnected() (connected bool, started bool) {
	if maintenance.connected == nil {
		maintenance.connected = atomic.NewBool(false)
	}
	return maintenance.connected.Load(), maintenance.started
}

func GetMaintenanceStatus() int {
	if maintenance.started {
		if maintenance.connected.Load() {
			if maintenance.resetProPassword {
				return types.MaintenanceModeResetPassword
			} else {
				return types.MaintenanceModeConnected
			}
		} else {
			return types.MaintenanceModeDiscover
		}
	}
	return types.NormalMode
}

func ConnectMaintenanceMode() bool {
	if maintenance.started {
		if maintenance.connected == nil {
			maintenance.connected = atomic.NewBool(false)
		}
		return maintenance.connected.CAS(false, true)
	}
	return false
}

func SetConnectUserId(id string) {
	maintenance.connectUserId = id
}

func GetMaintenanceUserId() (userId string) {
	if maintenance.started {
		return maintenance.connectUserId
	}
	return ""
}
func RefreshMaintenanceStarTime() int64 {
	if maintenance.started {
		maintenance.startTime = time.Now()
	}
	return maintenance.startTime.Unix()
}

func ResetProPassword(v bool) bool {
	if maintenance.started {
		maintenance.resetProPassword = v
		return maintenance.resetProPassword
	}
	return false
}

func DirectResetProPassword() bool {
	if maintenance.started {
		return maintenance.resetProPassword
	}
	return false
}
