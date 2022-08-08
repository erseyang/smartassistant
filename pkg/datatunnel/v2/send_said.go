package datatunnel

import (
	"time"

	"github.com/zhiting-tech/smartassistant/pkg/datatunnel/v2/control"
	"github.com/zhiting-tech/smartassistant/pkg/datatunnel/v2/proto"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
)

type TempConnCertInfo struct {
	SAID       string
	ExpireTime time.Duration
}

// AllowTempConnCert 临时连接凭证
func (c *ProxyControlClient) AllowTempConnCert(pcsc *control.ProxyControlStreamContext, req *proto.TempConnectionCertRequest) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn("AllowTempConnCert err:", err)
		}
	}()

	var (
		caller *control.RemoteCaller
	)

	caller, err = c.base.NewRemoteCaller(methodAllowTempConnCert, c.version)
	if err != nil {
		return
	}
	_, err = caller.Call(pcsc, req)

	return
}
