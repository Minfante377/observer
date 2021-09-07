// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.6.1
// source: src/api/api.proto

package api

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

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType int32  `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
	Date      string `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	User      string `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	Pwd       string `protobuf:"bytes,4,opt,name=pwd,proto3" json:"pwd,omitempty"`
	Cmd       string `protobuf:"bytes,5,opt,name=cmd,proto3" json:"cmd,omitempty"`
	Pid       string `protobuf:"bytes,6,opt,name=pid,proto3" json:"pid,omitempty"`
	Notes     string `protobuf:"bytes,7,opt,name=notes,proto3" json:"notes,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_api_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_src_api_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_src_api_api_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetEventType() int32 {
	if x != nil {
		return x.EventType
	}
	return 0
}

func (x *Event) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *Event) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *Event) GetPwd() string {
	if x != nil {
		return x.Pwd
	}
	return ""
}

func (x *Event) GetCmd() string {
	if x != nil {
		return x.Cmd
	}
	return ""
}

func (x *Event) GetPid() string {
	if x != nil {
		return x.Pid
	}
	return ""
}

func (x *Event) GetNotes() string {
	if x != nil {
		return x.Notes
	}
	return ""
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_api_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_src_api_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_src_api_api_proto_rawDescGZIP(), []int{1}
}

var File_src_api_api_proto protoreflect.FileDescriptor

var file_src_api_api_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x72, 0x63, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x22, 0x9a, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x77, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x70, 0x77, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x63,
	0x6d, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x70, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6e, 0x6f, 0x74, 0x65, 0x73, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x2c,
	0x0a, 0x06, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x22, 0x0a, 0x08, 0x4e, 0x65, 0x77, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x1a, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x09, 0x5a, 0x07,
	0x73, 0x72, 0x63, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_src_api_api_proto_rawDescOnce sync.Once
	file_src_api_api_proto_rawDescData = file_src_api_api_proto_rawDesc
)

func file_src_api_api_proto_rawDescGZIP() []byte {
	file_src_api_api_proto_rawDescOnce.Do(func() {
		file_src_api_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_src_api_api_proto_rawDescData)
	})
	return file_src_api_api_proto_rawDescData
}

var file_src_api_api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_src_api_api_proto_goTypes = []interface{}{
	(*Event)(nil), // 0: api.Event
	(*Empty)(nil), // 1: api.Empty
}
var file_src_api_api_proto_depIdxs = []int32{
	0, // 0: api.Events.NewEvent:input_type -> api.Event
	1, // 1: api.Events.NewEvent:output_type -> api.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_src_api_api_proto_init() }
func file_src_api_api_proto_init() {
	if File_src_api_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_src_api_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
		file_src_api_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_src_api_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_src_api_api_proto_goTypes,
		DependencyIndexes: file_src_api_api_proto_depIdxs,
		MessageInfos:      file_src_api_api_proto_msgTypes,
	}.Build()
	File_src_api_api_proto = out.File
	file_src_api_api_proto_rawDesc = nil
	file_src_api_api_proto_goTypes = nil
	file_src_api_api_proto_depIdxs = nil
}