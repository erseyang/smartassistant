package sasignal

import (
	"context"
	"github.com/zhiting-tech/smartassistant/modules/maintenance"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type userSignalCounter struct {
	lastTime time.Time
	count    int
}

func HandleUserSignal(ctx context.Context) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer close(sig)
	// 用于记录 SIGTERM 接收情况
	uSigTimer := time.NewTimer(time.Second * 3)
	uSigCounter := 0
	for {
		select {
		case <-uSigTimer.C:
			if uSigCounter == 0 {
				continue
			}
			logger.Infof("sig count %v", uSigCounter)
			// 按键 N 次响应
			maintenance.HandleMaintenanceChange(uSigCounter)
			uSigCounter = 0
		case s := <-sig:
			if s == syscall.SIGTERM {
				uSigCounter++
				uSigTimer.Reset(time.Second * 3)
			}
		case <-ctx.Done():
			return
		}
	}
}
