package job

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
	"sync"
)

var once sync.Once
var jobServer *JobServer

type JobServer struct {
	Cron *cron.Cron
}

func GetJobServer() *JobServer {
	once.Do(func() {
		jobServer = &JobServer{
			Cron: cron.New(),
		}
	})
	return jobServer
}

func (s *JobServer) Run(ctx context.Context) {
	if _, err := s.Cron.AddFunc("59 23 * * *", RemoveLogTask); err != nil {
		logger.Info(err)
	}

	s.Cron.Start()
	<-ctx.Done()
	logger.Warning("job server stopped")
}
