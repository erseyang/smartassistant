package task

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/plugin"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2"
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2/definer"

	"github.com/jinzhu/now"
	"gorm.io/gorm"
)

var (
	manager     Manager
	managerOnce sync.Once
)

const (
	sceneDefaultDelaySecond = 1 // 默认自动循环场景延迟秒数
)

type Manager interface {
	AddSceneTask(entity.Scene)
	AddSceneTaskForceDelay(entity.Scene) // 场景任务强制延迟频率控制
	DeleteSceneTask(sceneID int)
	RestartSceneTask(sceneID int) error
	DeviceStateChange(d entity.Device, attr definer.AttributeEvent) error
	Run(ctx context.Context)
}

type sceneTasksManager struct {
	queue *queueServe
	tasks map[string]*Task // TaskID -> *Task
	mu    sync.Mutex
}

func newSceneTasksManager(queue *queueServe) *sceneTasksManager {
	return &sceneTasksManager{
		queue: queue,
		tasks: make(map[string]*Task),
	}
}

func (st *sceneTasksManager) AddTasks(tasks ...*Task) {
	st.mu.Lock()
	for _, task := range tasks {
		st.tasks[task.ID] = task
	}
	st.mu.Unlock()
}

func (st *sceneTasksManager) Executed(taskID string) {
	st.mu.Lock()
	delete(st.tasks, taskID)
	st.mu.Unlock()
}

func (st *sceneTasksManager) Remove(taskID string) {
	st.mu.Lock()
	task, ok := st.tasks[taskID]
	if ok {
		delete(st.tasks, taskID)
		st.queue._remove(task.index)
	}
	st.mu.Unlock()
}

func (st *sceneTasksManager) RemoveAll() {
	st.mu.Lock()
	for _, task := range st.tasks {
		st.queue._remove(task.index)
	}
	st.tasks = make(map[string]*Task)
	st.mu.Unlock()
}

// LocalManager Task 服务
type LocalManager struct {
	queue        *queueServe
	runningScene sync.Map     // 正在执行的场景的id -> queue index
	scenes       sync.Map     // 保存queue中记录所有与entity.Scene相关的未执行的场景 sceneID -> *SceneTasks
	sceneRunTime sceneRunInfo // 记录各场景执行时间
}

type sceneRunInfo struct {
	mutex   sync.Mutex
	runData map[int]runData
}

type runData struct {
	runningTime int64
	hasNext     bool
}

func NewLocalManager() *LocalManager {
	return &LocalManager{
		queue: newQueueServe(),
	}
}

func SetManager(m Manager) {
	managerOnce.Do(func() {
		manager = m
	})
}

func GetManager() Manager {
	managerOnce.Do(func() {
		manager = &LocalManager{
			queue: newQueueServe(),
		}
	})
	return manager
}

// Run 启动服务，扫描插件并且连接通讯
func (m *LocalManager) Run(ctx context.Context) {
	logger.Info("starting task manager")
	go m.queue.start(ctx)
	// 重启时编排任务
	m.addSceneTaskByTime(time.Now())
	// 每天 23:55:00 进行第二天任务编排
	m.addArrangeSceneTask(now.EndOfDay().Add(-5 * time.Minute))
	// TODO 扫描已安装的插件，并且启动，连接 state change...
	<-ctx.Done()
	// TODO 断开连接
	logger.Warning("task manager stopped")
}

// addSceneTaskByTime 编排场景任务
func (m *LocalManager) addSceneTaskByTime(t time.Time) {
	scenes, err := entity.GetPendingScenesByTime(t)
	if err != nil {
		logger.Errorf("get execute scenes err %v", err)
		return
	}
	for _, scene := range scenes {
		// 没有定时触发条件 不加入队列
		if !scene.HaveTimeCondition() {
			continue
		}
		m.AddSceneTaskWithTime(scene, t)
	}
}

// addArrangeSceneTask 每天定时编排场景任务
func (m *LocalManager) addArrangeSceneTask(executeTime time.Time) {
	var f TaskFunc
	f = func(task *Task) error {
		m.addSceneTaskByTime(executeTime.AddDate(0, 0, 1))

		// 将下一个定时编排任务排进队列
		m.addArrangeSceneTask(executeTime.AddDate(0, 0, 1))
		return nil
	}

	task := NewTaskAt(f, executeTime)
	m.pushTask(task, "daily arrange scene task")
}

// DeleteSceneTask 删除场景任务
func (m *LocalManager) DeleteSceneTask(sceneID int) {
	// 现时需求如果场景对应的任务已运行，则不需要处理
	value, ok := m.scenes.LoadAndDelete(sceneID)
	if ok {
		sceneTasks := value.(*sceneTasksManager)
		sceneTasks.RemoveAll()
	}
}

// addSceneTaskByID 根据场景id执行场景（执行或者开启时调用）
// forceDelay 强制延迟执行控制频率
func (m *LocalManager) addSceneTaskByID(sceneID int, forceDelay bool) error {
	scene, err := entity.GetSceneInfoById(sceneID)
	if err != nil {
		if errors2.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(status.SceneNotExist)
		}
		return errors.Wrap(err, errors.InternalServerErr)
	}
	if forceDelay {
		m.AddSceneTaskForceDelay(scene)
		return nil
	}
	m.AddSceneTask(scene)
	return nil
}

// AddSceneTask 添加场景任务（执行或者开启时调用）
func (m *LocalManager) AddSceneTask(scene entity.Scene) {
	if scene.AutoRun {
		m.AddSceneTaskWithTime(scene, time.Now())
	} else {
		// 手动场景
		m.ManualScene(scene, false)
	}

}

// AddSceneTaskForceDelay 强制延迟场景任务
func (m *LocalManager) AddSceneTaskForceDelay(scene entity.Scene) {
	if scene.AutoRun {
		m.AddSceneTaskWithTime(scene, time.Now())
	} else {
		// 强制频率控制的场景
		if !m.checkScenesRunOnce(scene.ID){
			return
		}
		// 手动场景
		m.ManualScene(scene, true)
	}
}


func (m *LocalManager) AddSceneTaskWithTime(scene entity.Scene, t time.Time) {
	var task *Task
	date := now.New(t)
	// 开启自动场景
	// 找到定时条件的时间
	for _, c := range scene.SceneConditions {
		if c.ConditionType == entity.ConditionTypeTiming {

			// 获取任务今天的下次执行时间
			execTime := date.BeginningOfDay().Add(c.TimingAt.Sub(now.New(c.TimingAt).BeginningOfDay()))
			if execTime.Before(time.Now()) || execTime.After(date.EndOfDay()) {
				//logger.Debugf("now:%v,invalid next execute time:%v", time.Now(), execTime)
				continue
			}

			if !IsConditionsSatisfied(scene, true) {
				//logger.Debugf("auto scene:%d's conditions not satisfied", scene.ID)
				continue
			}
			task = NewTaskAt(m.wrapSceneFunc(scene, false), execTime)
			m.pushTask(task, scene)
			continue
		}
	}
}

// ManualScene 非自动场景
// forceDelay 默认频率控制
func (m *LocalManager) ManualScene(scene entity.Scene, forceDelay bool) {
	var task *Task
	task = NewTask(m.wrapSceneFunc(scene, forceDelay), 0)
	m.pushTask(task, scene)
}

func (m *LocalManager) pushTask(task *Task, target interface{}) {
	task.WithWrapper(m.sceneTaskManageWrapper(task, target), taskLogWrapper(target))
	m.queue.push(task)
}

// RestartSceneTask 重启场景对应的任务（就是删除然后重新添加任务）
func (m *LocalManager) RestartSceneTask(sceneID int) error {
	scene, err := entity.GetSceneInfoById(sceneID)
	if err != nil {
		if errors2.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(status.SceneNotExist)
		}
		return errors.Wrap(err, errors.InternalServerErr)
	}
	if !scene.AutoRun { // 手动执行的任务不需要重启
		return nil
	}
	m.DeleteSceneTask(sceneID)
	return m.addSceneTaskByID(sceneID, false)
}

func (m *LocalManager) addRunningScene(sceneID int, queueIndex int) {
	m.runningScene.Store(sceneID, queueIndex)
}

// sceneTaskManageWrapper 记录当前任务队列中与entity.Scene相关的未执行的场景任务
func (m *LocalManager) sceneTaskManageWrapper(task *Task, target interface{}) WrapperFunc {
	var (
		sceneTasks  *sceneTasksManager
		wrapperFunc WrapperFunc
	)

	// 判断是否处理entity.Scene
	scene, ok := target.(entity.Scene)
	if !ok {
		wrapperFunc = func(f TaskFunc) TaskFunc {
			return f
		}
	} else {

		// 统计当前所有entity.Scene相关的场景任务
		value, ok := m.scenes.Load(scene.ID)
		if !ok {
			sceneTasks = newSceneTasksManager(m.queue)
			value, ok = m.scenes.LoadOrStore(scene.ID, sceneTasks)
		}
		sceneTasks = value.(*sceneTasksManager)
		sceneTasks.AddTasks(task)

		wrapperFunc = func(f TaskFunc) TaskFunc {
			return func(task *Task) error {
				// 不存在则说明当前场景已经移除
				value, ok := m.scenes.Load(scene.ID)
				if !ok {
					return nil
				}
				// 移除要执行的任务
				sceneTasks := value.(*sceneTasksManager)
				sceneTasks.Executed(task.ID)

				return f(task)
			}
		}
	}

	return wrapperFunc
}

// wrapSceneFunc  包装场景为 TaskFunc
// forceDelayAble 是否需要强制执行默认延迟
func (m *LocalManager) wrapSceneFunc(sc entity.Scene, forceDelayAble bool) (f TaskFunc) {
	var (
		delayIc int
	)
	return func(t *Task) error {
		scene, err := entity.GetSceneInfoById(sc.ID)
		if err != nil {
			if errors2.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(status.SceneNotExist)
			}
			return errors.Wrap(err, errors.InternalServerErr)
		}
		// TODO 过滤旧版本场景, sc.version < scene.version

		if scene.Deleted.Valid { // 已删除的场景不执行
			return errors.New(status.SceneNotExist)
		}
		// TODO 此代码达到其功能，需清理
		m.addRunningScene(scene.ID, t.index)
		for _, sceneTask := range scene.SceneTasks {
			if forceDelayAble && sceneTask.DelaySeconds == 0 {
				delayIc = sceneDefaultDelaySecond
			}
			delay := time.Duration(sceneTask.DelaySeconds+delayIc) * time.Second
			task := NewTask(m.wrapTaskToFunc(sceneTask), delay).WithParent(t)

			if sceneTask.Type == entity.TaskTypeSmartDevice { // 控制设备
				if len(sceneTask.Attributes) == 0 {
					continue
				}
				deviceID := sceneTask.DeviceID
				var device entity.Device
				device, err := entity.GetDeviceByIDWithUnscoped(deviceID)
				if err == nil {
					m.pushTask(task, device)
				}
			} else {
				controlScene, err := entity.GetSceneByIDWithUnscoped(sceneTask.ControlSceneID)
				if err == nil {
					m.pushTask(task, controlScene)
				}
			}
		}
		// 场景执行结束解除占用状态
		if forceDelayAble {
			// 延迟指定时间
			time.Sleep(time.Second * sceneDefaultDelaySecond)
			m.relieveRunSceneStatus(scene.ID)
		}
		return nil
	}
}

// wrapTaskToFunc 包装场景任务为 TaskFunc
func (m *LocalManager) wrapTaskToFunc(task entity.SceneTask) (f TaskFunc) {
	return func(t *Task) error {
		// TODO 判断权限、判断场景是否有修改
		//logger.Debugf("execute task:%d,type:%d\n", task.ID, task.Type)
		switch task.Type {
		case entity.TaskTypeSmartDevice: // 控制设备
			return m.executeDevice(task)
		case entity.TaskTypeManualRun: // 执行场景
			return m.addSceneTaskByID(task.ControlSceneID, true)
		case entity.TaskTypeEnableAutoRun: // 开启场景
			return m.setSceneOn(task.ControlSceneID)
		case entity.TaskTypeDisableAutoRun: // 关闭场景
			return m.setSceneOff(task.ControlSceneID)
		}
		return nil
	}
}

// executeDevice 控制设备执行
func (m *LocalManager) executeDevice(task entity.SceneTask) (err error) {

	var ds []entity.Attribute
	if err := json.Unmarshal(task.Attributes, &ds); err != nil {
		logger.Error(err)
		return err
	}
	for _, d := range ds {
		var device entity.Device
		device, err = entity.GetDeviceByID(task.DeviceID)
		if err != nil {
			if errors2.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(status.DeviceNotExist)
			}
			return errors.Wrap(err, http.StatusInternalServerError)
		}
		//logger.Debugf("execute device command device id:%d instance id:%s attr:%s val:%v",
		//	device.ID, device.IID, "d.Attribute.Attribute", d.Attribute.Val)

		setReq := sdk.SetRequest{
			Attributes: []sdk.SetAttribute{
				{
					IID: device.IID,
					AID: d.Attribute.AID,
					Val: d.Attribute.Val,
				},
			}}
		err = controlDeviceFailRetry(device, setReq)
		if err != nil {
			identify := plugin.Identify{
				PluginID: device.PluginID,
				IID:      device.IID,
				AreaID:   device.AreaID,
			}
			return errors.Wrapf(err, status.DeviceOffline, identify.ID())
		}
	}
	return
}

// 控制设备失败重试2S内最多重试三次
func controlDeviceFailRetry(device entity.Device, setReq sdk.SetRequest) (err error) {
	err = plugin.SetAttributes(context.Background(), device.PluginID, device.AreaID, setReq)
	if err == nil {
		return
	}
	timer := time.NewTimer(time.Second * 2)
	defer timer.Stop()
	for i := 0; i < 3; i++ {
		select {
		case <-timer.C:
			return
		default:
			err = plugin.SetAttributes(context.Background(), device.PluginID, device.AreaID, setReq)
			if err == nil {
				return
			}
		}
	}
	return
}

// SetSceneOn 开启场景
func (m *LocalManager) setSceneOn(sceneID int) (err error) {
	if err = entity.SwitchAutoSceneByID(sceneID, true); err != nil {
		return
	}
	if err := m.addSceneTaskByID(sceneID, true); err != nil {
		logger.Error(err)
	}
	return
}

// SetSceneOff 关闭场景
func (m *LocalManager) setSceneOff(sceneID int) (err error) {
	if err = entity.SwitchAutoSceneByID(sceneID, false); err != nil {
		return
	}
	m.DeleteSceneTask(sceneID)
	return
}

// DeviceStateChange 设备状态变化触发场景
func (m *LocalManager) DeviceStateChange(d entity.Device, ac definer.AttributeEvent) (err error) {

	deviceID := d.ID
	scenes, err := entity.GetScenesByCondition(deviceID, ac)
	if err != nil {
		return fmt.Errorf("can't find scenes with device %d %s %d change",
			deviceID, ac.IID, ac.AID)
	}

	// 遍历并包装场景为任务
	for _, scene := range scenes {
		// 判断该执行是否需要限制
		if !m.checkScenesRunOnce(scene.ID) {
			continue // 无需执行
		}
		scene, _ = entity.GetSceneInfoById(scene.ID)
		// 全部满足且有定时条件则不执行
		if scene.IsMatchAllCondition() && scene.HaveTimeCondition() {
			//logger.Debugf("device %d state %d changed but scenes %d not match time conditoin,ignore\n",
			//	deviceID, ac.AID, scene.ID)
			continue
		}

		if !IsConditionsSatisfied(scene, false) {
			//logger.Debugf("auto scene:%d's conditions not satisfied", scene.ID)
			continue
		}
		t := NewTask(m.wrapSceneFunc(scene, true), 0)
		m.pushTask(t, scene)
	}
	return
}

// 执行时间限制单位时间s内同一场景只允许执行一次
func (m *LocalManager) checkScenesRunOnce(sceneId int)  bool {
	m.sceneRunTime.mutex.Lock()
	defer m.sceneRunTime.mutex.Unlock()
	// 初始化
	if m.sceneRunTime.runData == nil {
		m.sceneRunTime.runData = make(map[int]runData)
	}
	// 场景上次执行时间间隔判断
	nowTime := time.Now().Unix()
	sceneRunData := m.sceneRunTime.runData[sceneId]
	// 没有时间限制则或有scene正在执行
	if sceneRunData.runningTime+sceneDefaultDelaySecond > nowTime && sceneRunData.hasNext {
		return false
	} else {
		// 不存在或者可以执行标记当前场景为执行状态
		m.sceneRunTime.runData[sceneId] = runData{
			runningTime: nowTime,
			hasNext:     true,
		}
		return true
	}
}

func (m *LocalManager) relieveRunSceneStatus(sceneId int) {
	m.sceneRunTime.mutex.Lock()
	defer m.sceneRunTime.mutex.Unlock()
	sceneRunData, ok := m.sceneRunTime.runData[sceneId]
	if !ok {
		m.sceneRunTime.runData[sceneId] = runData{
			hasNext: false,
		}
	}
	sceneRunData.hasNext = false
	m.sceneRunTime.runData[sceneId] = sceneRunData
}
