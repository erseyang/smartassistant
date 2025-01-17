// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.2
// source: proxy_control.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ProxyControlStreamMsgType int32

const (
	// 请求
	ProxyControlStreamMsgType_REQUEST ProxyControlStreamMsgType = 0
	// 响应
	ProxyControlStreamMsgType_RESPONSE ProxyControlStreamMsgType = 1
	// 通知,不需要响应
	ProxyControlStreamMsgType_NOTIFY ProxyControlStreamMsgType = 2
)

// Enum value maps for ProxyControlStreamMsgType.
var (
	ProxyControlStreamMsgType_name = map[int32]string{
		0: "REQUEST",
		1: "RESPONSE",
		2: "NOTIFY",
	}
	ProxyControlStreamMsgType_value = map[string]int32{
		"REQUEST":  0,
		"RESPONSE": 1,
		"NOTIFY":   2,
	}
)

func (x ProxyControlStreamMsgType) Enum() *ProxyControlStreamMsgType {
	p := new(ProxyControlStreamMsgType)
	*p = x
	return p
}

func (x ProxyControlStreamMsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProxyControlStreamMsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_proxy_control_proto_enumTypes[0].Descriptor()
}

func (ProxyControlStreamMsgType) Type() protoreflect.EnumType {
	return &file_proxy_control_proto_enumTypes[0]
}

func (x ProxyControlStreamMsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProxyControlStreamMsgType.Descriptor instead.
func (ProxyControlStreamMsgType) EnumDescriptor() ([]byte, []int) {
	return file_proxy_control_proto_rawDescGZIP(), []int{0}
}

type ProxyControlStreamMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hdr  *ProxyControlStreamMsgHdr  `protobuf:"bytes,1,opt,name=Hdr,proto3" json:"Hdr,omitempty"`
	Body *ProxyControlStreamMsgBody `protobuf:"bytes,2,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (x *ProxyControlStreamMsg) Reset() {
	*x = ProxyControlStreamMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_control_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProxyControlStreamMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProxyControlStreamMsg) ProtoMessage() {}

func (x *ProxyControlStreamMsg) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_control_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProxyControlStreamMsg.ProtoReflect.Descriptor instead.
func (*ProxyControlStreamMsg) Descriptor() ([]byte, []int) {
	return file_proxy_control_proto_rawDescGZIP(), []int{0}
}

func (x *ProxyControlStreamMsg) GetHdr() *ProxyControlStreamMsgHdr {
	if x != nil {
		return x.Hdr
	}
	return nil
}

func (x *ProxyControlStreamMsg) GetBody() *ProxyControlStreamMsgBody {
	if x != nil {
		return x.Body
	}
	return nil
}

type ProxyControlStreamMsgHdr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId   int64                     `protobuf:"varint,1,opt,name=MessageId,proto3" json:"MessageId,omitempty"`
	MessageType ProxyControlStreamMsgType `protobuf:"varint,2,opt,name=MessageType,proto3,enum=proxy.control.proto.ProxyControlStreamMsgType" json:"MessageType,omitempty"`
}

func (x *ProxyControlStreamMsgHdr) Reset() {
	*x = ProxyControlStreamMsgHdr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_control_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProxyControlStreamMsgHdr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProxyControlStreamMsgHdr) ProtoMessage() {}

func (x *ProxyControlStreamMsgHdr) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_control_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProxyControlStreamMsgHdr.ProtoReflect.Descriptor instead.
func (*ProxyControlStreamMsgHdr) Descriptor() ([]byte, []int) {
	return file_proxy_control_proto_rawDescGZIP(), []int{1}
}

func (x *ProxyControlStreamMsgHdr) GetMessageId() int64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

func (x *ProxyControlStreamMsgHdr) GetMessageType() ProxyControlStreamMsgType {
	if x != nil {
		return x.MessageType
	}
	return ProxyControlStreamMsgType_REQUEST
}

type ProxyControlStreamMsgBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 方法版本
	Version int32 `protobuf:"varint,1,opt,name=Version,proto3" json:"Version,omitempty"`
	// 方法名
	Method string `protobuf:"bytes,2,opt,name=Method,proto3" json:"Method,omitempty"`
	// 参数
	Values [][]byte `protobuf:"bytes,3,rep,name=Values,proto3" json:"Values,omitempty"`
	// 状态码
	StatusCode int32 `protobuf:"varint,4,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
	// 自定义报错信息
	Reason string `protobuf:"bytes,5,opt,name=Reason,proto3" json:"Reason,omitempty"`
}

func (x *ProxyControlStreamMsgBody) Reset() {
	*x = ProxyControlStreamMsgBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_control_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProxyControlStreamMsgBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProxyControlStreamMsgBody) ProtoMessage() {}

func (x *ProxyControlStreamMsgBody) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_control_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProxyControlStreamMsgBody.ProtoReflect.Descriptor instead.
func (*ProxyControlStreamMsgBody) Descriptor() ([]byte, []int) {
	return file_proxy_control_proto_rawDescGZIP(), []int{2}
}

func (x *ProxyControlStreamMsgBody) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *ProxyControlStreamMsgBody) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *ProxyControlStreamMsgBody) GetValues() [][]byte {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *ProxyControlStreamMsgBody) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *ProxyControlStreamMsgBody) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type AuthenticateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SAID  string `protobuf:"bytes,1,opt,name=SAID,proto3" json:"SAID,omitempty"`
	SAKey string `protobuf:"bytes,2,opt,name=SAKey,proto3" json:"SAKey,omitempty"`
}

func (x *AuthenticateRequest) Reset() {
	*x = AuthenticateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_control_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateRequest) ProtoMessage() {}

func (x *AuthenticateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_control_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateRequest.ProtoReflect.Descriptor instead.
func (*AuthenticateRequest) Descriptor() ([]byte, []int) {
	return file_proxy_control_proto_rawDescGZIP(), []int{3}
}

func (x *AuthenticateRequest) GetSAID() string {
	if x != nil {
		return x.SAID
	}
	return ""
}

func (x *AuthenticateRequest) GetSAKey() string {
	if x != nil {
		return x.SAKey
	}
	return ""
}

type RegisterServiceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Services []*RegisterServiceValue `protobuf:"bytes,1,rep,name=Services,proto3" json:"Services,omitempty"`
}

func (x *RegisterServiceRequest) Reset() {
	*x = RegisterServiceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_control_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterServiceRequest) ProtoMessage() {}

func (x *RegisterServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_control_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterServiceRequest.ProtoReflect.Descriptor instead.
func (*RegisterServiceRequest) Descriptor() ([]byte, []int) {
	return file_proxy_control_proto_rawDescGZIP(), []int{4}
}

func (x *RegisterServiceRequest) GetServices() []*RegisterServiceValue {
	if x != nil {
		return x.Services
	}
	return nil
}

type RegisterServiceValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName string `protobuf:"bytes,1,opt,name=ServiceName,proto3" json:"ServiceName,omitempty"`
	ServicePort int32  `protobuf:"varint,2,opt,name=ServicePort,proto3" json:"ServicePort,omitempty"`
}

func (x *RegisterServiceValue) Reset() {
	*x = RegisterServiceValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_control_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterServiceValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterServiceValue) ProtoMessage() {}

func (x *RegisterServiceValue) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_control_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterServiceValue.ProtoReflect.Descriptor instead.
func (*RegisterServiceValue) Descriptor() ([]byte, []int) {
	return file_proxy_control_proto_rawDescGZIP(), []int{5}
}

func (x *RegisterServiceValue) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *RegisterServiceValue) GetServicePort() int32 {
	if x != nil {
		return x.ServicePort
	}
	return 0
}

type NewConnectionEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName string `protobuf:"bytes,1,opt,name=ServiceName,proto3" json:"ServiceName,omitempty"`
	Key         []byte `protobuf:"bytes,2,opt,name=Key,proto3" json:"Key,omitempty"`
	RemoteAddr  string `protobuf:"bytes,3,opt,name=RemoteAddr,proto3" json:"RemoteAddr,omitempty"`
}

func (x *NewConnectionEvent) Reset() {
	*x = NewConnectionEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_control_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewConnectionEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewConnectionEvent) ProtoMessage() {}

func (x *NewConnectionEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_control_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewConnectionEvent.ProtoReflect.Descriptor instead.
func (*NewConnectionEvent) Descriptor() ([]byte, []int) {
	return file_proxy_control_proto_rawDescGZIP(), []int{6}
}

func (x *NewConnectionEvent) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *NewConnectionEvent) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *NewConnectionEvent) GetRemoteAddr() string {
	if x != nil {
		return x.RemoteAddr
	}
	return ""
}

// TempConnectionCertRequest 临时连接凭据
type TempConnectionCertRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// SAID
	SaId string `protobuf:"bytes,1,opt,name=saId,proto3" json:"saId,omitempty"`
	// 二维码过期时间
	ExpireTime float64 `protobuf:"fixed64,2,opt,name=expire_time,json=expireTime,proto3" json:"expire_time,omitempty"`
}

func (x *TempConnectionCertRequest) Reset() {
	*x = TempConnectionCertRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_control_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TempConnectionCertRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TempConnectionCertRequest) ProtoMessage() {}

func (x *TempConnectionCertRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_control_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TempConnectionCertRequest.ProtoReflect.Descriptor instead.
func (*TempConnectionCertRequest) Descriptor() ([]byte, []int) {
	return file_proxy_control_proto_rawDescGZIP(), []int{7}
}

func (x *TempConnectionCertRequest) GetSaId() string {
	if x != nil {
		return x.SaId
	}
	return ""
}

func (x *TempConnectionCertRequest) GetExpireTime() float64 {
	if x != nil {
		return x.ExpireTime
	}
	return 0
}

var File_proxy_control_proto protoreflect.FileDescriptor

var file_proxy_control_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x01, 0x0a, 0x15, 0x50,
	0x72, 0x6f, 0x78, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x4d, 0x73, 0x67, 0x12, 0x3f, 0x0a, 0x03, 0x48, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2d, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x73, 0x67, 0x48, 0x64, 0x72,
	0x52, 0x03, 0x48, 0x64, 0x72, 0x12, 0x42, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x43,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x73, 0x67, 0x42,
	0x6f, 0x64, 0x79, 0x52, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x22, 0x8a, 0x01, 0x0a, 0x18, 0x50, 0x72,
	0x6f, 0x78, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x4d, 0x73, 0x67, 0x48, 0x64, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x49, 0x64, 0x12, 0x50, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x78,
	0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x72, 0x6f, 0x78, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0x9d, 0x01, 0x0a, 0x19, 0x50, 0x72, 0x6f, 0x78, 0x79,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x73, 0x67,
	0x42, 0x6f, 0x64, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x1e,
	0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x3f, 0x0a, 0x13, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e,
	0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x53, 0x41, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x53, 0x41, 0x49,
	0x44, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x41, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x53, 0x41, 0x4b, 0x65, 0x79, 0x22, 0x5f, 0x0a, 0x16, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x45, 0x0a, 0x08, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x22, 0x5a, 0x0a, 0x14, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x6f, 0x72,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x50, 0x6f, 0x72, 0x74, 0x22, 0x68, 0x0a, 0x12, 0x4e, 0x65, 0x77, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x1e,
	0x0a, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x41, 0x64, 0x64, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x41, 0x64, 0x64, 0x72, 0x22, 0x50,
	0x0a, 0x19, 0x54, 0x65, 0x6d, 0x70, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x43, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73,
	0x61, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x61, 0x49, 0x64, 0x12,
	0x1f, 0x0a, 0x0b, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x2a, 0x42, 0x0a, 0x19, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a,
	0x07, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45,
	0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x4f, 0x54, 0x49,
	0x46, 0x59, 0x10, 0x02, 0x32, 0x7e, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x12, 0x6b, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50,
	0x72, 0x6f, 0x78, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x4d, 0x73, 0x67, 0x1a, 0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x78, 0x79,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x73, 0x67,
	0x28, 0x01, 0x30, 0x01, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proxy_control_proto_rawDescOnce sync.Once
	file_proxy_control_proto_rawDescData = file_proxy_control_proto_rawDesc
)

func file_proxy_control_proto_rawDescGZIP() []byte {
	file_proxy_control_proto_rawDescOnce.Do(func() {
		file_proxy_control_proto_rawDescData = protoimpl.X.CompressGZIP(file_proxy_control_proto_rawDescData)
	})
	return file_proxy_control_proto_rawDescData
}

var file_proxy_control_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proxy_control_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proxy_control_proto_goTypes = []interface{}{
	(ProxyControlStreamMsgType)(0),    // 0: proxy.control.proto.ProxyControlStreamMsgType
	(*ProxyControlStreamMsg)(nil),     // 1: proxy.control.proto.ProxyControlStreamMsg
	(*ProxyControlStreamMsgHdr)(nil),  // 2: proxy.control.proto.ProxyControlStreamMsgHdr
	(*ProxyControlStreamMsgBody)(nil), // 3: proxy.control.proto.ProxyControlStreamMsgBody
	(*AuthenticateRequest)(nil),       // 4: proxy.control.proto.AuthenticateRequest
	(*RegisterServiceRequest)(nil),    // 5: proxy.control.proto.RegisterServiceRequest
	(*RegisterServiceValue)(nil),      // 6: proxy.control.proto.RegisterServiceValue
	(*NewConnectionEvent)(nil),        // 7: proxy.control.proto.NewConnectionEvent
	(*TempConnectionCertRequest)(nil), // 8: proxy.control.proto.TempConnectionCertRequest
}
var file_proxy_control_proto_depIdxs = []int32{
	2, // 0: proxy.control.proto.ProxyControlStreamMsg.Hdr:type_name -> proxy.control.proto.ProxyControlStreamMsgHdr
	3, // 1: proxy.control.proto.ProxyControlStreamMsg.Body:type_name -> proxy.control.proto.ProxyControlStreamMsgBody
	0, // 2: proxy.control.proto.ProxyControlStreamMsgHdr.MessageType:type_name -> proxy.control.proto.ProxyControlStreamMsgType
	6, // 3: proxy.control.proto.RegisterServiceRequest.Services:type_name -> proxy.control.proto.RegisterServiceValue
	1, // 4: proxy.control.proto.ProxyController.ControlStream:input_type -> proxy.control.proto.ProxyControlStreamMsg
	1, // 5: proxy.control.proto.ProxyController.ControlStream:output_type -> proxy.control.proto.ProxyControlStreamMsg
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proxy_control_proto_init() }
func file_proxy_control_proto_init() {
	if File_proxy_control_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proxy_control_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProxyControlStreamMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proxy_control_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProxyControlStreamMsgHdr); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proxy_control_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProxyControlStreamMsgBody); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proxy_control_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proxy_control_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterServiceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proxy_control_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterServiceValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proxy_control_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewConnectionEvent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proxy_control_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TempConnectionCertRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proxy_control_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proxy_control_proto_goTypes,
		DependencyIndexes: file_proxy_control_proto_depIdxs,
		EnumInfos:         file_proxy_control_proto_enumTypes,
		MessageInfos:      file_proxy_control_proto_msgTypes,
	}.Build()
	File_proxy_control_proto = out.File
	file_proxy_control_proto_rawDesc = nil
	file_proxy_control_proto_goTypes = nil
	file_proxy_control_proto_depIdxs = nil
}
