# 开发您的第一个插件

此文档描述如何开发一个简单插件，面向插件开发者。

开发前先阅读插件设计概要：[插件系统设计技术概要](../guide/plugin-module.md)

## 插件实现
下面展示如何快速实现一个可控的（虚拟）灯设备的插件开发

### 新建项目，获取sdk

```
plugin
├── device
│   ├── attribute.go
│   ├── switch.go
│   └── light.go
├── lib
│   ├── protocol.go
│   └── discover.go
├── config.json
├── html
│   ├── index.html
│   └── images
...
├── Readme.md
├── .dockerignore
├── Dockerfile
└── main.go
```

```shell
go get github.com/zhiting-tech/smartassistant
```

### 定义设备

实现前推荐先了解插件设备生命周期管理：[插件设备生命周期管理](../guide/plugin-device-lifecycle.md)

`plugin/device/light.go`
```go
package device

import (
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2"
	"github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2/definer"
)

type Light struct {
}

func NewLight() sdk.Device {
	return &Light{}
}

func (l Light) Address() string {
	// TODO implement me
	panic("implement me")
}

func (l Light) Connect() error {
	// TODO implement me
	panic("implement me")
}

func (l Light) Disconnect(iid string) error {
	// TODO implement me
	panic("implement me")
}

func (l Light) Define(df *definer.Definer) error {
	// TODO implement me
	panic("implement me")
}

func (l Light) Info() sdk.DeviceInfo {
	// TODO implement me
	panic("implement me")
}

func (l Light) Online(iid string) bool {
	// TODO implement me
	panic("implement me")
}


```

### 实现设备地址接口
`plugin/device/light.go`
```go
func (l Light) Address() string {
	// 该地址返回改变会导致设备断开重连
	return "192.168.0.123:12345"
}
```
### 实现设备信息接口
`plugin/device/light.go`
```go
func (l Light) Info() sdk.DeviceInfo {
    return sdk.DeviceInfo{
        IID:          "iid",
        Model:        "model",
        Manufacturer: "manufacturer",
        AuthRequired: false,
    }
}
```
### 实现设备连接和断开接口
`plugin/device/light.go`
```go
func (l Light) Connect() error {
    // 此处由于是虚拟设备，所以没有连接过程
    // return l.protocol.Connect()
    return nil
}
func (l Light) Disconnect(iid string) error {
    // 此处由于是虚拟设备，所以没有断开连接过程
    // return l.protocol.Disconnect()
    return nil
}

```
### 实现设备资源回收接口
`plugin/device/light.go`
```go
func (l Light) Close() error {
    // 此处由于是虚拟设备，所以没有回收资源的过程
    // return l.protocol.Close()
    return nil
}
```
### 实现设备是否在线接口
`plugin/device/light.go`
```go
func (l Light) Online(iid string) bool {
    // 此处由于是虚拟设备，所以没判断实际是否在线
    // return l.protocol.IsOnline()
    return true
}
```

### 通过协议定义设备属性和实现控制

`plugin/device/attribute.go`

```go
package device

import "github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2/definer"

// OnOff 开关
// 通过实现 thingmodel.IAttribute 的接口，以便sdk调用
type OnOff struct {
	service *definer.BaseService
}

func (l *OnOff) Set(val interface{}) error {

	// 如果有设备，则跟据已实现的通讯协议控制设备
	// return l.protocol.SetSwitch(val)

	// 此处仅将设置的状态同步回SA
	return l.service.Notify(thingmodel.OnOff, val)
}


```

### 实现设备定义和生成物模型接口

`plugin/device/light.go`
```go

func (l Light) Define(df *definer.Definer) error {
    info := df.Instance(l.Info().IID).NewInfo()
    info.WithAttribute(thingmodel.Manufacturer).SetVal(l.Info().Manufacturer)
    info.WithAttribute(thingmodel.Model).SetVal(l.Info().Model)
    info.WithAttribute(thingmodel.Identify).SetVal(l.Info().IID)
    info.WithAttribute(thingmodel.Version).SetVal("1.0.0")
    
    light := df.Instance(l.Info().IID).NewLight()
    light.Enable(thingmodel.OnOff, OnOff{light}).SetVal("on")
	
	// 为了使得物模型属性状态和设备状态保持一致，需要开发者主动同步设备状态变化
    // go func() {
    //     for {
    //         // 此处同步设备属性状态
    //         light.Notify(thingmodel.OnOff, "on")
    //     }
    // }()
}


```

### 初始化和运行

`plugin/main.go`
```go
package main

import (
	"context"
	"log"
	
	sdk "github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2"

	"plugin/device"
)

func main() {
    p := sdk.NewPluginServer(Discover)
    err := sdk.Run(p)
    if err != nil {
        log.Panicln(err)
    }
}

// Discover 这里需要实现一个发现设备的方法,给sdk调用
func Discover(ctx context.Context, devices chan<- sdk.Device) {
    l := getDevice()
    devices <- l
    <-ctx.Done()
}

// getDevice 发现设备，由于没有设备，所以没有实现发现设备的逻辑
func getDevice() sdk.Device {
    return device.NewLight()
}
```

### Dockerfile编写和调试

暂时仅支持以镜像方式安装插件，调试正常后，编译成镜像提供给SA

1) Dockerfile示例参考

`plugin/Dockerfile`
```dockerfile
FROM golang:1.16-alpine as builder
RUN apk add build-base
COPY . /app
WORKDIR /app
RUN go env -w GOPROXY="goproxy.cn,direct"
RUN go build -ldflags="-w -s" -o demo-plugin

FROM alpine
WORKDIR /app
COPY --from=builder /app/demo-plugin /app/demo-plugin

# static file
COPY ./html ./html
ENTRYPOINT ["/app/demo-plugin"]

```

2) 编译镜像

```shell
docker build -f your_plugin_Dockerfile -t your_plugin_name
```

3) 运行插件

SA上通过**我的-支持品牌-我的插件-创作-添加插件**可以将插件安装到SA上调试


4) 更多：[设备类插件开发指南](../guide/device-plugin.md)

### Demo

[demo-plugin](../../examples/plugin-demo) :
通过上文的插件实现教程实现的示例插件；这是一个模拟设备写的一个简单插件服务，不依赖硬件，实现了核心插件的功能
