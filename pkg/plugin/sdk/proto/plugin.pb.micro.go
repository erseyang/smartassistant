// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: plugin.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Plugin service

func NewPluginEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Plugin service

type PluginService interface {
	Discover(ctx context.Context, in *Empty, opts ...client.CallOption) (Plugin_DiscoverService, error)
	StateChange(ctx context.Context, in *Empty, opts ...client.CallOption) (Plugin_StateChangeService, error)
	GetAttributes(ctx context.Context, in *GetAttributesReq, opts ...client.CallOption) (*GetAttributesResp, error)
	SetAttributes(ctx context.Context, in *SetAttributesReq, opts ...client.CallOption) (*SetAttributesResp, error)
}

type pluginService struct {
	c    client.Client
	name string
}

func NewPluginService(name string, c client.Client) PluginService {
	return &pluginService{
		c:    c,
		name: name,
	}
}

func (c *pluginService) Discover(ctx context.Context, in *Empty, opts ...client.CallOption) (Plugin_DiscoverService, error) {
	req := c.c.NewRequest(c.name, "Plugin.Discover", &Empty{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &pluginServiceDiscover{stream}, nil
}

type Plugin_DiscoverService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*Device, error)
}

type pluginServiceDiscover struct {
	stream client.Stream
}

func (x *pluginServiceDiscover) Close() error {
	return x.stream.Close()
}

func (x *pluginServiceDiscover) Context() context.Context {
	return x.stream.Context()
}

func (x *pluginServiceDiscover) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *pluginServiceDiscover) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *pluginServiceDiscover) Recv() (*Device, error) {
	m := new(Device)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pluginService) StateChange(ctx context.Context, in *Empty, opts ...client.CallOption) (Plugin_StateChangeService, error) {
	req := c.c.NewRequest(c.name, "Plugin.StateChange", &Empty{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &pluginServiceStateChange{stream}, nil
}

type Plugin_StateChangeService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*State, error)
}

type pluginServiceStateChange struct {
	stream client.Stream
}

func (x *pluginServiceStateChange) Close() error {
	return x.stream.Close()
}

func (x *pluginServiceStateChange) Context() context.Context {
	return x.stream.Context()
}

func (x *pluginServiceStateChange) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *pluginServiceStateChange) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *pluginServiceStateChange) Recv() (*State, error) {
	m := new(State)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pluginService) GetAttributes(ctx context.Context, in *GetAttributesReq, opts ...client.CallOption) (*GetAttributesResp, error) {
	req := c.c.NewRequest(c.name, "Plugin.GetAttributes", in)
	out := new(GetAttributesResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginService) SetAttributes(ctx context.Context, in *SetAttributesReq, opts ...client.CallOption) (*SetAttributesResp, error) {
	req := c.c.NewRequest(c.name, "Plugin.SetAttributes", in)
	out := new(SetAttributesResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Plugin service

type PluginHandler interface {
	Discover(context.Context, *Empty, Plugin_DiscoverStream) error
	StateChange(context.Context, *Empty, Plugin_StateChangeStream) error
	GetAttributes(context.Context, *GetAttributesReq, *GetAttributesResp) error
	SetAttributes(context.Context, *SetAttributesReq, *SetAttributesResp) error
}

func RegisterPluginHandler(s server.Server, hdlr PluginHandler, opts ...server.HandlerOption) error {
	type plugin interface {
		Discover(ctx context.Context, stream server.Stream) error
		StateChange(ctx context.Context, stream server.Stream) error
		GetAttributes(ctx context.Context, in *GetAttributesReq, out *GetAttributesResp) error
		SetAttributes(ctx context.Context, in *SetAttributesReq, out *SetAttributesResp) error
	}
	type Plugin struct {
		plugin
	}
	h := &pluginHandler{hdlr}
	return s.Handle(s.NewHandler(&Plugin{h}, opts...))
}

type pluginHandler struct {
	PluginHandler
}

func (h *pluginHandler) Discover(ctx context.Context, stream server.Stream) error {
	m := new(Empty)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.PluginHandler.Discover(ctx, m, &pluginDiscoverStream{stream})
}

type Plugin_DiscoverStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Device) error
}

type pluginDiscoverStream struct {
	stream server.Stream
}

func (x *pluginDiscoverStream) Close() error {
	return x.stream.Close()
}

func (x *pluginDiscoverStream) Context() context.Context {
	return x.stream.Context()
}

func (x *pluginDiscoverStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *pluginDiscoverStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *pluginDiscoverStream) Send(m *Device) error {
	return x.stream.Send(m)
}

func (h *pluginHandler) StateChange(ctx context.Context, stream server.Stream) error {
	m := new(Empty)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.PluginHandler.StateChange(ctx, m, &pluginStateChangeStream{stream})
}

type Plugin_StateChangeStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*State) error
}

type pluginStateChangeStream struct {
	stream server.Stream
}

func (x *pluginStateChangeStream) Close() error {
	return x.stream.Close()
}

func (x *pluginStateChangeStream) Context() context.Context {
	return x.stream.Context()
}

func (x *pluginStateChangeStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *pluginStateChangeStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *pluginStateChangeStream) Send(m *State) error {
	return x.stream.Send(m)
}

func (h *pluginHandler) GetAttributes(ctx context.Context, in *GetAttributesReq, out *GetAttributesResp) error {
	return h.PluginHandler.GetAttributes(ctx, in, out)
}

func (h *pluginHandler) SetAttributes(ctx context.Context, in *SetAttributesReq, out *SetAttributesResp) error {
	return h.PluginHandler.SetAttributes(ctx, in, out)
}
