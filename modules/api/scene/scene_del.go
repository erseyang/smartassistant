package scene

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/task"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"strconv"
)

type DeleteSceneReq struct {
	SceneIds []int `json:"scene_ids"`
}

// BatchDeleteScene 用于处理删除场景接口的请求
func BatchDeleteScene(c *gin.Context) {
	var err error
	var req DeleteSceneReq
	defer func() {
		response.HandleResponse(c, err, nil)

	}()

	if err = c.BindJSON(&req); err != nil {
		err = errors.New(errors.BadRequest)
		return
	}

	u := session.Get(c)
	if !entity.JudgePermit(u.UserID, types.SceneDel) {
		err = errors.New(status.SceneDeleteDeny)
		return
	}

	// 删除场景需要与用户属于同一个家庭
	scenes, err := entity.GetAreaScenesByIDs(u.AreaID, req.SceneIds)
	if err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	if len(scenes) != len(req.SceneIds) {
		err = errors.New(status.SceneParamIncorrectErr)
		return
	}

	for _, sceneID := range req.SceneIds {
		if err = entity.DeleteScene(sceneID); err != nil {
			return
		}

		task.GetManager().DeleteSceneTask(sceneID)
	}

}

// DeleteScene 用于处理删除单个场景接口的请求
func DeleteScene(c *gin.Context) {
	var err error
	defer func() {
		response.HandleResponse(c, err, nil)

	}()

	sceneId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = errors.New(errors.BadRequest)
		return
	}

	if !entity.JudgePermit(session.Get(c).UserID, types.SceneDel) {
		err = errors.New(status.SceneDeleteDeny)
		return
	}

	if err = entity.DeleteScene(sceneId); err != nil {
		return
	}

	task.GetManager().DeleteSceneTask(sceneId)

}
