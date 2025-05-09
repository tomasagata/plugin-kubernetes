// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: clusterRegistration/clusterRegistration.proto

package clusterRegistration

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

type CS1Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HelloServiceManager string `protobuf:"bytes,1,opt,name=hello_service_manager,json=helloServiceManager,proto3" json:"hello_service_manager,omitempty"`
}

func (x *CS1Message) Reset() {
	*x = CS1Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clusterRegistration_clusterRegistration_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CS1Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CS1Message) ProtoMessage() {}

func (x *CS1Message) ProtoReflect() protoreflect.Message {
	mi := &file_clusterRegistration_clusterRegistration_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CS1Message.ProtoReflect.Descriptor instead.
func (*CS1Message) Descriptor() ([]byte, []int) {
	return file_clusterRegistration_clusterRegistration_proto_rawDescGZIP(), []int{0}
}

func (x *CS1Message) GetHelloServiceManager() string {
	if x != nil {
		return x.HelloServiceManager
	}
	return ""
}

type SC1Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HelloClusterManager string `protobuf:"bytes,1,opt,name=hello_cluster_manager,json=helloClusterManager,proto3" json:"hello_cluster_manager,omitempty"`
}

func (x *SC1Message) Reset() {
	*x = SC1Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clusterRegistration_clusterRegistration_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SC1Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SC1Message) ProtoMessage() {}

func (x *SC1Message) ProtoReflect() protoreflect.Message {
	mi := &file_clusterRegistration_clusterRegistration_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SC1Message.ProtoReflect.Descriptor instead.
func (*SC1Message) Descriptor() ([]byte, []int) {
	return file_clusterRegistration_clusterRegistration_proto_rawDescGZIP(), []int{1}
}

func (x *SC1Message) GetHelloClusterManager() string {
	if x != nil {
		return x.HelloClusterManager
	}
	return ""
}

type CS2Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ManagerPort          int32       `protobuf:"varint,1,opt,name=manager_port,json=managerPort,proto3" json:"manager_port,omitempty"`
	NetworkComponentPort int32       `protobuf:"varint,2,opt,name=network_component_port,json=networkComponentPort,proto3" json:"network_component_port,omitempty"`
	ClusterName          string      `protobuf:"bytes,3,opt,name=cluster_name,json=clusterName,proto3" json:"cluster_name,omitempty"`
	ClusterInfo          []*KeyValue `protobuf:"bytes,4,rep,name=cluster_info,json=clusterInfo,proto3" json:"cluster_info,omitempty"`
	ClusterLocation      string      `protobuf:"bytes,5,opt,name=cluster_location,json=clusterLocation,proto3" json:"cluster_location,omitempty"`
}

func (x *CS2Message) Reset() {
	*x = CS2Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clusterRegistration_clusterRegistration_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CS2Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CS2Message) ProtoMessage() {}

func (x *CS2Message) ProtoReflect() protoreflect.Message {
	mi := &file_clusterRegistration_clusterRegistration_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CS2Message.ProtoReflect.Descriptor instead.
func (*CS2Message) Descriptor() ([]byte, []int) {
	return file_clusterRegistration_clusterRegistration_proto_rawDescGZIP(), []int{2}
}

func (x *CS2Message) GetManagerPort() int32 {
	if x != nil {
		return x.ManagerPort
	}
	return 0
}

func (x *CS2Message) GetNetworkComponentPort() int32 {
	if x != nil {
		return x.NetworkComponentPort
	}
	return 0
}

func (x *CS2Message) GetClusterName() string {
	if x != nil {
		return x.ClusterName
	}
	return ""
}

func (x *CS2Message) GetClusterInfo() []*KeyValue {
	if x != nil {
		return x.ClusterInfo
	}
	return nil
}

func (x *CS2Message) GetClusterLocation() string {
	if x != nil {
		return x.ClusterLocation
	}
	return ""
}

type KeyValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"` // Ist das ein String, oder kann es auch etwas anderes sein?
}

func (x *KeyValue) Reset() {
	*x = KeyValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clusterRegistration_clusterRegistration_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValue) ProtoMessage() {}

func (x *KeyValue) ProtoReflect() protoreflect.Message {
	mi := &file_clusterRegistration_clusterRegistration_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValue.ProtoReflect.Descriptor instead.
func (*KeyValue) Descriptor() ([]byte, []int) {
	return file_clusterRegistration_clusterRegistration_proto_rawDescGZIP(), []int{3}
}

func (x *KeyValue) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *KeyValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// SC2Message represents a message in Protocol Buffers.
type SC2Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *SC2Message) Reset() {
	*x = SC2Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clusterRegistration_clusterRegistration_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SC2Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SC2Message) ProtoMessage() {}

func (x *SC2Message) ProtoReflect() protoreflect.Message {
	mi := &file_clusterRegistration_clusterRegistration_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SC2Message.ProtoReflect.Descriptor instead.
func (*SC2Message) Descriptor() ([]byte, []int) {
	return file_clusterRegistration_clusterRegistration_proto_rawDescGZIP(), []int{4}
}

func (x *SC2Message) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_clusterRegistration_clusterRegistration_proto protoreflect.FileDescriptor

var file_clusterRegistration_clusterRegistration_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x13, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x40, 0x0a, 0x0a, 0x43, 0x53, 0x31, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x32, 0x0a, 0x15, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x13, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x22, 0x40, 0x0a, 0x0a, 0x53, 0x43, 0x31, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x32, 0x0a, 0x15, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x5f, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x13, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x22, 0xf5, 0x01, 0x0a, 0x0a, 0x43, 0x53, 0x32,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x34, 0x0a, 0x16, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x5f,
	0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x14, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x50, 0x6f, 0x72, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x40, 0x0a, 0x0c, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x69,
	0x6e, 0x66, 0x6f, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0b, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x29, 0x0a, 0x10, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x32, 0x0a, 0x08, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x1c, 0x0a, 0x0a, 0x53, 0x43, 0x32, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x32, 0xc7, 0x01, 0x0a, 0x10, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x5a, 0x0a, 0x14, 0x68, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x5f, 0x69, 0x6e, 0x69, 0x74, 0x5f, 0x67, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x12,
	0x1f, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x53, 0x31, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x1a, 0x1f, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x43, 0x31, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x22, 0x00, 0x12, 0x57, 0x0a, 0x11, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x5f, 0x69, 0x6e,
	0x69, 0x74, 0x5f, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x12, 0x1f, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43,
	0x53, 0x32, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x1f, 0x2e, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x53, 0x43, 0x32, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x42, 0x3b, 0x5a, 0x39,
	0x6f, 0x61, 0x6b, 0x65, 0x73, 0x74, 0x72, 0x61, 0x2f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x72,
	0x61, 0x74, 0x65, 0x2d, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_clusterRegistration_clusterRegistration_proto_rawDescOnce sync.Once
	file_clusterRegistration_clusterRegistration_proto_rawDescData = file_clusterRegistration_clusterRegistration_proto_rawDesc
)

func file_clusterRegistration_clusterRegistration_proto_rawDescGZIP() []byte {
	file_clusterRegistration_clusterRegistration_proto_rawDescOnce.Do(func() {
		file_clusterRegistration_clusterRegistration_proto_rawDescData = protoimpl.X.CompressGZIP(file_clusterRegistration_clusterRegistration_proto_rawDescData)
	})
	return file_clusterRegistration_clusterRegistration_proto_rawDescData
}

var file_clusterRegistration_clusterRegistration_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_clusterRegistration_clusterRegistration_proto_goTypes = []interface{}{
	(*CS1Message)(nil), // 0: clusterRegistration.CS1Message
	(*SC1Message)(nil), // 1: clusterRegistration.SC1Message
	(*CS2Message)(nil), // 2: clusterRegistration.CS2Message
	(*KeyValue)(nil),   // 3: clusterRegistration.KeyValue
	(*SC2Message)(nil), // 4: clusterRegistration.SC2Message
}
var file_clusterRegistration_clusterRegistration_proto_depIdxs = []int32{
	3, // 0: clusterRegistration.CS2Message.cluster_info:type_name -> clusterRegistration.KeyValue
	0, // 1: clusterRegistration.register_cluster.handle_init_greeting:input_type -> clusterRegistration.CS1Message
	2, // 2: clusterRegistration.register_cluster.handle_init_final:input_type -> clusterRegistration.CS2Message
	1, // 3: clusterRegistration.register_cluster.handle_init_greeting:output_type -> clusterRegistration.SC1Message
	4, // 4: clusterRegistration.register_cluster.handle_init_final:output_type -> clusterRegistration.SC2Message
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_clusterRegistration_clusterRegistration_proto_init() }
func file_clusterRegistration_clusterRegistration_proto_init() {
	if File_clusterRegistration_clusterRegistration_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_clusterRegistration_clusterRegistration_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CS1Message); i {
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
		file_clusterRegistration_clusterRegistration_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SC1Message); i {
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
		file_clusterRegistration_clusterRegistration_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CS2Message); i {
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
		file_clusterRegistration_clusterRegistration_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValue); i {
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
		file_clusterRegistration_clusterRegistration_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SC2Message); i {
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
			RawDescriptor: file_clusterRegistration_clusterRegistration_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_clusterRegistration_clusterRegistration_proto_goTypes,
		DependencyIndexes: file_clusterRegistration_clusterRegistration_proto_depIdxs,
		MessageInfos:      file_clusterRegistration_clusterRegistration_proto_msgTypes,
	}.Build()
	File_clusterRegistration_clusterRegistration_proto = out.File
	file_clusterRegistration_clusterRegistration_proto_rawDesc = nil
	file_clusterRegistration_clusterRegistration_proto_goTypes = nil
	file_clusterRegistration_clusterRegistration_proto_depIdxs = nil
}
