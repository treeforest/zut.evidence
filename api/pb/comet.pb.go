// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: comet.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PushMsgReq_MsgType int32

const (
	PushMsgReq_Unknown        PushMsgReq_MsgType = 0
	PushMsgReq_Audit_DOING    PushMsgReq_MsgType = 1 // 待审核
	PushMsgReq_Apply_FAILED   PushMsgReq_MsgType = 2 // 审核失败
	PushMsgReq_Apply_DONE     PushMsgReq_MsgType = 3 // 审核成功
	PushMsgReq_Challenge_REQ  PushMsgReq_MsgType = 4 // 挑战请求
	PushMsgReq_Challenge_RESP PushMsgReq_MsgType = 5 // 挑战响应
)

// Enum value maps for PushMsgReq_MsgType.
var (
	PushMsgReq_MsgType_name = map[int32]string{
		0: "Unknown",
		1: "Audit_DOING",
		2: "Apply_FAILED",
		3: "Apply_DONE",
		4: "Challenge_REQ",
		5: "Challenge_RESP",
	}
	PushMsgReq_MsgType_value = map[string]int32{
		"Unknown":        0,
		"Audit_DOING":    1,
		"Apply_FAILED":   2,
		"Apply_DONE":     3,
		"Challenge_REQ":  4,
		"Challenge_RESP": 5,
	}
)

func (x PushMsgReq_MsgType) Enum() *PushMsgReq_MsgType {
	p := new(PushMsgReq_MsgType)
	*p = x
	return p
}

func (x PushMsgReq_MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PushMsgReq_MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_comet_proto_enumTypes[0].Descriptor()
}

func (PushMsgReq_MsgType) Type() protoreflect.EnumType {
	return &file_comet_proto_enumTypes[0]
}

func (x PushMsgReq_MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PushMsgReq_MsgType.Descriptor instead.
func (PushMsgReq_MsgType) EnumDescriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{0, 0}
}

type PushMsgReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid  int64              `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`                           // 推送的目标
	Type PushMsgReq_MsgType `protobuf:"varint,2,opt,name=type,proto3,enum=PushMsgReq_MsgType" json:"type,omitempty"` // 推送的消息类型
	Msg  string             `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`                            // 推送的消息
}

func (x *PushMsgReq) Reset() {
	*x = PushMsgReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMsgReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMsgReq) ProtoMessage() {}

func (x *PushMsgReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMsgReq.ProtoReflect.Descriptor instead.
func (*PushMsgReq) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{0}
}

func (x *PushMsgReq) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *PushMsgReq) GetType() PushMsgReq_MsgType {
	if x != nil {
		return x.Type
	}
	return PushMsgReq_Unknown
}

func (x *PushMsgReq) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type PushMsgReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushMsgReply) Reset() {
	*x = PushMsgReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMsgReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMsgReply) ProtoMessage() {}

func (x *PushMsgReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMsgReply.ProtoReflect.Descriptor instead.
func (*PushMsgReply) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{1}
}

var File_comet_proto protoreflect.FileDescriptor

var file_comet_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcb, 0x01, 0x0a, 0x0a, 0x50,
	0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x27, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x50, 0x75, 0x73, 0x68,
	0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x70, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x0f,
	0x0a, 0x0b, 0x41, 0x75, 0x64, 0x69, 0x74, 0x5f, 0x44, 0x4f, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12,
	0x10, 0x0a, 0x0c, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10,
	0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x5f, 0x44, 0x4f, 0x4e, 0x45, 0x10,
	0x03, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x52,
	0x45, 0x51, 0x10, 0x04, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x10, 0x05, 0x22, 0x0e, 0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68,
	0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x32, 0x66, 0x0a, 0x05, 0x43, 0x6f, 0x6d, 0x65,
	0x74, 0x12, 0x36, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x25, 0x0a, 0x07, 0x50, 0x75, 0x73,
	0x68, 0x4d, 0x73, 0x67, 0x12, 0x0b, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65,
	0x71, 0x1a, 0x0d, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x42, 0x06, 0x5a, 0x04, 0x2e, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_comet_proto_rawDescOnce sync.Once
	file_comet_proto_rawDescData = file_comet_proto_rawDesc
)

func file_comet_proto_rawDescGZIP() []byte {
	file_comet_proto_rawDescOnce.Do(func() {
		file_comet_proto_rawDescData = protoimpl.X.CompressGZIP(file_comet_proto_rawDescData)
	})
	return file_comet_proto_rawDescData
}

var file_comet_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_comet_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_comet_proto_goTypes = []interface{}{
	(PushMsgReq_MsgType)(0), // 0: PushMsgReq.MsgType
	(*PushMsgReq)(nil),      // 1: PushMsgReq
	(*PushMsgReply)(nil),    // 2: PushMsgReply
	(*emptypb.Empty)(nil),   // 3: google.protobuf.Empty
}
var file_comet_proto_depIdxs = []int32{
	0, // 0: PushMsgReq.type:type_name -> PushMsgReq.MsgType
	3, // 1: Comet.Ping:input_type -> google.protobuf.Empty
	1, // 2: Comet.PushMsg:input_type -> PushMsgReq
	3, // 3: Comet.Ping:output_type -> google.protobuf.Empty
	2, // 4: Comet.PushMsg:output_type -> PushMsgReply
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_comet_proto_init() }
func file_comet_proto_init() {
	if File_comet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_comet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMsgReq); i {
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
		file_comet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMsgReply); i {
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
			RawDescriptor: file_comet_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_comet_proto_goTypes,
		DependencyIndexes: file_comet_proto_depIdxs,
		EnumInfos:         file_comet_proto_enumTypes,
		MessageInfos:      file_comet_proto_msgTypes,
	}.Build()
	File_comet_proto = out.File
	file_comet_proto_rawDesc = nil
	file_comet_proto_goTypes = nil
	file_comet_proto_depIdxs = nil
}
