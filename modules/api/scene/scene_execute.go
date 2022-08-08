package scene

import (
	"strconv"

	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/task"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"

	"github.com/gin-gonic/gin"

	"github.com/zhiting-tech/smartassistant/pkg/errors"
)

// ExecuteSceneReq 场景执行接口请求参数
type ExecuteSceneReq struct {
	IsExecute bool `json:"is_execute"`
}

// ExecuteScene 用于处理场景执行接口的请求
func ExecuteScene(c *gin.Context) {
	var (
		req ExecuteSceneReq
		err error
	)
	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	if err = c.BindJSON(&req); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	sceneId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = errors.New(errors.BadRequest)
		return
	}

	if err = entity.CheckSceneExitById(sceneId); err != nil {
		return
	}

	if err = req.executeScene(c, sceneId); err != nil {
		return
	}
}

func (req ExecuteSceneReq) executeScene(c *gin.Context, sceneId int) (err error) {

	controlPermission, err := CheckUserControlPermission(c, sceneId, session.Get(c).UserID)
	if !controlPermission {
		err = errors.New(status.DeviceOrSceneControlDeny)
		return
	}

	scene, err := entity.GetSceneInfoById(sceneId)
	if err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	if scene.AutoRun {
		// 自动场景先更新数据库
		if err = entity.SwitchAutoScene(&scene, req.IsExecute); err != nil {
			return
		}
	}

	if req.IsExecute {
		task.GetManager().AddSceneTask(scene)
	} else {
		task.GetManager().DeleteSceneTask(scene.ID)
	}

	return
}
