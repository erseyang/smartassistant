# 设备类插件开发指南

开发前先阅读插件设计概要：[插件系统设计技术概要](plugin-module.md)

使用 [plugin-sdk](../../pkg/plugin/sdk) 可以忽略不重要的逻辑，快速实现插件

## 插件实现

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

### 根据设备的通讯协议实现
以下仅列举可能的协议实现，实际的实现需要视具体的协议来定

`plugin/lib/protocol.go`

```go
package lib

import "github.com/zhiting-tech/smartassistant/pkg/thingmodel"

type Protocol struct {
}

func (p Protocol) Connect() error {
	return nil
}

func (p Protocol) Disconnect() error {
	return nil
}

func (p Protocol) Close() error {
	return nil
}

func (p Protocol) IsOnline() bool {
	return true
}

func (p Protocol) SetAttribute(attr string, val interface{}) error {
	return nil
}

type AttributeChange struct {
	Attr thingmodel.Attribute
	Val  interface{}
}

func (p Protocol) AttributeChange() chan AttributeChange {
	return nil
}


```

### 根据设备的发现协议实现

`plugin/lib/discover.go`

```go
package lib

func Discover() Protocol {
    return &Protocol{}
}

```

### 定义设备

实现前推荐先了解插件设备生命周期管理：[插件设备生命周期管理](../guide/plugin-device-lifecycle.md)

`plugin/device/light.go`
```go
package device

import (
    "plugin/lib"
    
    "github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2"
    "github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2/definer"
)

type Light struct {
    protocol lib.Protocol
}

func NewLight(protocol lib.Protocol) sdk.Device {
    return &Light{protocol}
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
    return l.protocol.Connect()
}
func (l Light) Disconnect(iid string) error {
    return l.protocol.Disconnect()
}

```
### 实现设备资源回收接口
`plugin/device/light.go`
```go
func (l Light) Close() error {
    return l.protocol.Close()
}
```
### 实现设备是否在线接口
该方法用来判断设备是否在线，会被频繁调用，尽量使用非阻塞的实现
`plugin/device/light.go`
```go
func (l Light) Online(iid string) bool {
    return l.protocol.IsOnline()
}
```

### 通过协议定义设备属性和实现控制

`plugin/device/attribute.go`
```go
package plugin

import (
    "plugin/lib"
)

// OnOff
// 通过实现 thingmodel.IAttribute 的接口，以便sdk调用
type OnOff struct {
	protocol *lib.Protocol
}

func (l OnOff) Set(val interface{}) error {
	return l.protocol.SetAttribute("on_off", val)
}

```

### 实现设备定义和生成物模型接口

sdk中提供了预定义的设备模型，使用模型可以方便SA有效进行管理和控制

- 物模型概念请参考[物模型](thing-model.md)
- 设备模型请参考[设备模型](device-thing-model.md)

### 注意！：
- definer.Definer重复New Service会导致物模型错误，在开发中要注意不要重复New Service
- 定义设备如果在SA上被使用了后重新修改有可能导致未知错误，比如属性功能被SA的权限或者场景使用后

`plugin/device/light.go`
```go

func (l Light) Define(df *definer.Definer) error {
    info := df.Instance(l.Info().IID).NewInfo()
    info.WithAttribute(thingmodel.Manufacturer).SetVal(l.Info().Manufacturer)
    info.WithAttribute(thingmodel.Model).SetVal(l.Info().Model)
    info.WithAttribute(thingmodel.Identify).SetVal(l.Info().IID)
    info.WithAttribute(thingmodel.Version).SetVal("1.0.0")
    
    light := df.Instance(l.Info().IID).NewLight()
    light.Enable(thingmodel.OnOff, OnOff{l.protocol}).SetVal("on")
    
    // 为了使得物模型属性状态和设备状态保持一致，需要开发者主动同步设备状态变化
    go func() {
        for change:=range l.protocol.AttributeChange(){
            // 此处同步设备属性状态
            light.Notify(change.Attr, change.Val)
        }
    }()
}


```

### 自定义服务或属性：

### 注意！：
- 通过自定义实现的功能无法用于云对云和homebridge等需要统一物模型的服务

#### 重写属性权限

对每个属性和配置都可以重置权限

`plugin/device/light.go`
```go
thingmodel.OnOff.SetPermissions(
    thingmodel.AttributePermissionWrite,
    thingmodel.AttributePermissionRead,
    thingmodel.AttributePermissionNotify,
)
```
#### 自定义服务
`plugin/device/light.go`
```go
package plugin

import (
    "github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2/definer"
    "github.com/zhiting-tech/smartassistant/pkg/thingmodel"
)

const newServiceType thingmodel.ServiceType = "new_service"

var newService = definer.NewService(newServiceType)

```

#### 自定义属性：

`plugin/device/light.go`
```go
package plugin

import (
    "github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2/definer"
    "github.com/zhiting-tech/smartassistant/pkg/thingmodel"
)

var newAttr = thingmodel.Attribute{
	Type:    "new_attr",
	ValType: String,
	Permission: thingmodel.SetPermissions(
		AttributePermissionRead,
	),
}

```

### 初始化和运行

定义好设备和实现方法后，运行插件服务（包括grpc和http服务）

`plugin/main.go`
```go
package main

import (
    "context"
    "log"
    
    sdk "github.com/zhiting-tech/smartassistant/pkg/plugin/sdk/v2"
    
    "plugin/device"
    "plugin/lib"
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
    
	protocols:=lib.Discover()
	
	for {
		select {
		case <-ctx.Done():
			return
		case protocol:=<-protocols:
            devices<-device.NewLight(protocol)
        }
    }
}
```

这样服务就会运行起来，并通过SA的etcd地址0.0.0.0:2379注册插件服务， SA会通过etcd发现插件服务并且建立通道开始通信并且转发请求和命令

### 镜像编译和部署

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


### Demo

[demo-plugin](../../examples/plugin-demo) :
通过上文的插件实现教程实现的示例插件；这是一个模拟设备写的一个简单插件服务，不依赖硬件，实现了核心插件的功能

## 快速开始

[快速开始](../tutorial/plugin-quickstart.md)