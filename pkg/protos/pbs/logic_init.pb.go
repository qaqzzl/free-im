// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: logic_init.proto

package pbs

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type TokenAuthReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *AuthMessage `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TokenAuthReq) Reset() {
	*x = TokenAuthReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_init_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenAuthReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenAuthReq) ProtoMessage() {}

func (x *TokenAuthReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_init_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenAuthReq.ProtoReflect.Descriptor instead.
func (*TokenAuthReq) Descriptor() ([]byte, []int) {
	return file_logic_init_proto_rawDescGZIP(), []int{0}
}

func (x *TokenAuthReq) GetMessage() *AuthMessage {
	if x != nil {
		return x.Message
	}
	return nil
}

type TokenAuthResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Statu bool `protobuf:"varint,1,opt,name=statu,proto3" json:"statu,omitempty"`
}

func (x *TokenAuthResp) Reset() {
	*x = TokenAuthResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_init_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenAuthResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenAuthResp) ProtoMessage() {}

func (x *TokenAuthResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_init_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenAuthResp.ProtoReflect.Descriptor instead.
func (*TokenAuthResp) Descriptor() ([]byte, []int) {
	return file_logic_init_proto_rawDescGZIP(), []int{1}
}

func (x *TokenAuthResp) GetStatu() bool {
	if x != nil {
		return x.Statu
	}
	return false
}

type MessageReceiveReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *MessageItem `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *MessageReceiveReq) Reset() {
	*x = MessageReceiveReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_init_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageReceiveReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageReceiveReq) ProtoMessage() {}

func (x *MessageReceiveReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_init_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageReceiveReq.ProtoReflect.Descriptor instead.
func (*MessageReceiveReq) Descriptor() ([]byte, []int) {
	return file_logic_init_proto_rawDescGZIP(), []int{2}
}

func (x *MessageReceiveReq) GetMessage() *MessageItem {
	if x != nil {
		return x.Message
	}
	return nil
}

type MessageReceiveResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MessageReceiveResp) Reset() {
	*x = MessageReceiveResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_init_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageReceiveResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageReceiveResp) ProtoMessage() {}

func (x *MessageReceiveResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_init_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageReceiveResp.ProtoReflect.Descriptor instead.
func (*MessageReceiveResp) Descriptor() ([]byte, []int) {
	return file_logic_init_proto_rawDescGZIP(), []int{3}
}

type MessageACKReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId string `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
}

func (x *MessageACKReq) Reset() {
	*x = MessageACKReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_init_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageACKReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageACKReq) ProtoMessage() {}

func (x *MessageACKReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_init_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageACKReq.ProtoReflect.Descriptor instead.
func (*MessageACKReq) Descriptor() ([]byte, []int) {
	return file_logic_init_proto_rawDescGZIP(), []int{4}
}

func (x *MessageACKReq) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

type MessageACKResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MessageACKResp) Reset() {
	*x = MessageACKResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_init_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageACKResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageACKResp) ProtoMessage() {}

func (x *MessageACKResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_init_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageACKResp.ProtoReflect.Descriptor instead.
func (*MessageACKResp) Descriptor() ([]byte, []int) {
	return file_logic_init_proto_rawDescGZIP(), []int{5}
}

type MessageSyncReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId    string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	MessageId string `protobuf:"bytes,2,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
}

func (x *MessageSyncReq) Reset() {
	*x = MessageSyncReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_init_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageSyncReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageSyncReq) ProtoMessage() {}

func (x *MessageSyncReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_init_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageSyncReq.ProtoReflect.Descriptor instead.
func (*MessageSyncReq) Descriptor() ([]byte, []int) {
	return file_logic_init_proto_rawDescGZIP(), []int{6}
}

func (x *MessageSyncReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MessageSyncReq) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

type MessageSyncResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MessageSyncResp) Reset() {
	*x = MessageSyncResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_init_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageSyncResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageSyncResp) ProtoMessage() {}

func (x *MessageSyncResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_init_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageSyncResp.ProtoReflect.Descriptor instead.
func (*MessageSyncResp) Descriptor() ([]byte, []int) {
	return file_logic_init_proto_rawDescGZIP(), []int{7}
}

var File_logic_init_proto protoreflect.FileDescriptor

var file_logic_init_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x5f, 0x69, 0x6e, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x70, 0x62, 0x73, 0x1a, 0x08, 0x69, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x3a, 0x0a, 0x0c, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65,
	0x71, 0x12, 0x2a, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x25, 0x0a,
	0x0d, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x75, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x22, 0x3f, 0x0a, 0x11, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x12, 0x2a, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x62, 0x73,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x2e, 0x0a, 0x0d, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x43, 0x4b, 0x52, 0x65, 0x71, 0x12, 0x1d, 0x0a, 0x0a,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0x10, 0x0a, 0x0e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x43, 0x4b, 0x52, 0x65, 0x73, 0x70, 0x22, 0x48, 0x0a,
	0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0x11, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x70, 0x32, 0xf3, 0x01, 0x0a, 0x09, 0x4c,
	0x6f, 0x67, 0x69, 0x63, 0x49, 0x6e, 0x69, 0x74, 0x12, 0x32, 0x0a, 0x09, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x41, 0x75, 0x74, 0x68, 0x12, 0x11, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x12, 0x41, 0x0a, 0x0e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x12, 0x16,
	0x2e, 0x70, 0x62, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x63, 0x65,
	0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12,
	0x35, 0x0a, 0x0a, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x43, 0x4b, 0x12, 0x12, 0x2e,
	0x70, 0x62, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x43, 0x4b, 0x52, 0x65,
	0x71, 0x1a, 0x13, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41,
	0x43, 0x4b, 0x52, 0x65, 0x73, 0x70, 0x12, 0x38, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x53, 0x79, 0x6e, 0x63, 0x12, 0x13, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x70, 0x62, 0x73,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x70,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_logic_init_proto_rawDescOnce sync.Once
	file_logic_init_proto_rawDescData = file_logic_init_proto_rawDesc
)

func file_logic_init_proto_rawDescGZIP() []byte {
	file_logic_init_proto_rawDescOnce.Do(func() {
		file_logic_init_proto_rawDescData = protoimpl.X.CompressGZIP(file_logic_init_proto_rawDescData)
	})
	return file_logic_init_proto_rawDescData
}

var file_logic_init_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_logic_init_proto_goTypes = []interface{}{
	(*TokenAuthReq)(nil),       // 0: pbs.TokenAuthReq
	(*TokenAuthResp)(nil),      // 1: pbs.TokenAuthResp
	(*MessageReceiveReq)(nil),  // 2: pbs.MessageReceiveReq
	(*MessageReceiveResp)(nil), // 3: pbs.MessageReceiveResp
	(*MessageACKReq)(nil),      // 4: pbs.MessageACKReq
	(*MessageACKResp)(nil),     // 5: pbs.MessageACKResp
	(*MessageSyncReq)(nil),     // 6: pbs.MessageSyncReq
	(*MessageSyncResp)(nil),    // 7: pbs.MessageSyncResp
	(*AuthMessage)(nil),        // 8: pbs.AuthMessage
	(*MessageItem)(nil),        // 9: pbs.MessageItem
}
var file_logic_init_proto_depIdxs = []int32{
	8, // 0: pbs.TokenAuthReq.message:type_name -> pbs.AuthMessage
	9, // 1: pbs.MessageReceiveReq.message:type_name -> pbs.MessageItem
	0, // 2: pbs.LogicInit.TokenAuth:input_type -> pbs.TokenAuthReq
	2, // 3: pbs.LogicInit.MessageReceive:input_type -> pbs.MessageReceiveReq
	4, // 4: pbs.LogicInit.MessageACK:input_type -> pbs.MessageACKReq
	6, // 5: pbs.LogicInit.MessageSync:input_type -> pbs.MessageSyncReq
	1, // 6: pbs.LogicInit.TokenAuth:output_type -> pbs.TokenAuthResp
	3, // 7: pbs.LogicInit.MessageReceive:output_type -> pbs.MessageReceiveResp
	5, // 8: pbs.LogicInit.MessageACK:output_type -> pbs.MessageACKResp
	7, // 9: pbs.LogicInit.MessageSync:output_type -> pbs.MessageSyncResp
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_logic_init_proto_init() }
func file_logic_init_proto_init() {
	if File_logic_init_proto != nil {
		return
	}
	file_im_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_logic_init_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenAuthReq); i {
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
		file_logic_init_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenAuthResp); i {
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
		file_logic_init_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageReceiveReq); i {
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
		file_logic_init_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageReceiveResp); i {
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
		file_logic_init_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageACKReq); i {
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
		file_logic_init_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageACKResp); i {
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
		file_logic_init_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageSyncReq); i {
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
		file_logic_init_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageSyncResp); i {
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
			RawDescriptor: file_logic_init_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_logic_init_proto_goTypes,
		DependencyIndexes: file_logic_init_proto_depIdxs,
		MessageInfos:      file_logic_init_proto_msgTypes,
	}.Build()
	File_logic_init_proto = out.File
	file_logic_init_proto_rawDesc = nil
	file_logic_init_proto_goTypes = nil
	file_logic_init_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// LogicInitClient is the client API for LogicInit service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LogicInitClient interface {
	// token 连接认证
	TokenAuth(ctx context.Context, in *TokenAuthReq, opts ...grpc.CallOption) (*TokenAuthResp, error)
	// 消息接收
	MessageReceive(ctx context.Context, in *MessageReceiveReq, opts ...grpc.CallOption) (*MessageReceiveResp, error)
	// 消息回执
	MessageACK(ctx context.Context, in *MessageACKReq, opts ...grpc.CallOption) (*MessageACKResp, error)
	// 消息同步
	MessageSync(ctx context.Context, in *MessageSyncReq, opts ...grpc.CallOption) (*MessageSyncResp, error)
}

type logicInitClient struct {
	cc grpc.ClientConnInterface
}

func NewLogicInitClient(cc grpc.ClientConnInterface) LogicInitClient {
	return &logicInitClient{cc}
}

func (c *logicInitClient) TokenAuth(ctx context.Context, in *TokenAuthReq, opts ...grpc.CallOption) (*TokenAuthResp, error) {
	out := new(TokenAuthResp)
	err := c.cc.Invoke(ctx, "/pbs.LogicInit/TokenAuth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicInitClient) MessageReceive(ctx context.Context, in *MessageReceiveReq, opts ...grpc.CallOption) (*MessageReceiveResp, error) {
	out := new(MessageReceiveResp)
	err := c.cc.Invoke(ctx, "/pbs.LogicInit/MessageReceive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicInitClient) MessageACK(ctx context.Context, in *MessageACKReq, opts ...grpc.CallOption) (*MessageACKResp, error) {
	out := new(MessageACKResp)
	err := c.cc.Invoke(ctx, "/pbs.LogicInit/MessageACK", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicInitClient) MessageSync(ctx context.Context, in *MessageSyncReq, opts ...grpc.CallOption) (*MessageSyncResp, error) {
	out := new(MessageSyncResp)
	err := c.cc.Invoke(ctx, "/pbs.LogicInit/MessageSync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogicInitServer is the server API for LogicInit service.
type LogicInitServer interface {
	// token 连接认证
	TokenAuth(context.Context, *TokenAuthReq) (*TokenAuthResp, error)
	// 消息接收
	MessageReceive(context.Context, *MessageReceiveReq) (*MessageReceiveResp, error)
	// 消息回执
	MessageACK(context.Context, *MessageACKReq) (*MessageACKResp, error)
	// 消息同步
	MessageSync(context.Context, *MessageSyncReq) (*MessageSyncResp, error)
}

// UnimplementedLogicInitServer can be embedded to have forward compatible implementations.
type UnimplementedLogicInitServer struct {
}

func (*UnimplementedLogicInitServer) TokenAuth(context.Context, *TokenAuthReq) (*TokenAuthResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TokenAuth not implemented")
}
func (*UnimplementedLogicInitServer) MessageReceive(context.Context, *MessageReceiveReq) (*MessageReceiveResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MessageReceive not implemented")
}
func (*UnimplementedLogicInitServer) MessageACK(context.Context, *MessageACKReq) (*MessageACKResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MessageACK not implemented")
}
func (*UnimplementedLogicInitServer) MessageSync(context.Context, *MessageSyncReq) (*MessageSyncResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MessageSync not implemented")
}

func RegisterLogicInitServer(s *grpc.Server, srv LogicInitServer) {
	s.RegisterService(&_LogicInit_serviceDesc, srv)
}

func _LogicInit_TokenAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenAuthReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicInitServer).TokenAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pbs.LogicInit/TokenAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicInitServer).TokenAuth(ctx, req.(*TokenAuthReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogicInit_MessageReceive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageReceiveReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicInitServer).MessageReceive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pbs.LogicInit/MessageReceive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicInitServer).MessageReceive(ctx, req.(*MessageReceiveReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogicInit_MessageACK_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageACKReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicInitServer).MessageACK(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pbs.LogicInit/MessageACK",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicInitServer).MessageACK(ctx, req.(*MessageACKReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogicInit_MessageSync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageSyncReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicInitServer).MessageSync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pbs.LogicInit/MessageSync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicInitServer).MessageSync(ctx, req.(*MessageSyncReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _LogicInit_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pbs.LogicInit",
	HandlerType: (*LogicInitServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TokenAuth",
			Handler:    _LogicInit_TokenAuth_Handler,
		},
		{
			MethodName: "MessageReceive",
			Handler:    _LogicInit_MessageReceive_Handler,
		},
		{
			MethodName: "MessageACK",
			Handler:    _LogicInit_MessageACK_Handler,
		},
		{
			MethodName: "MessageSync",
			Handler:    _LogicInit_MessageSync_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logic_init.proto",
}
