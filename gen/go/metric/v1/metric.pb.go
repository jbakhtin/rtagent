// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: metric/v1/metric.proto

package metricv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

type Type int32

const (
	Type_TYPE_UNSPECIFIED Type = 0
	Type_TYPE_GAUGE       Type = 1
	Type_TYPE_COUNTER     Type = 2
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "TYPE_UNSPECIFIED",
		1: "TYPE_GAUGE",
		2: "TYPE_COUNTER",
	}
	Type_value = map[string]int32{
		"TYPE_UNSPECIFIED": 0,
		"TYPE_GAUGE":       1,
		"TYPE_COUNTER":     2,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_metric_v1_metric_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_metric_v1_metric_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_metric_v1_metric_proto_rawDescGZIP(), []int{0}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  Type    `protobuf:"varint,1,opt,name=type,proto3,enum=metric.v1.Type" json:"type,omitempty"`
	Key   string  `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Delta uint64  `protobuf:"varint,3,opt,name=delta,proto3" json:"delta,omitempty"`
	Value float32 `protobuf:"fixed32,4,opt,name=value,proto3" json:"value,omitempty"`
	Hash  string  `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_v1_metric_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_metric_v1_metric_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_metric_v1_metric_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_TYPE_UNSPECIFIED
}

func (x *Metric) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Metric) GetDelta() uint64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Metric) GetValue() float32 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Metric) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type UpdateMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *UpdateMetricRequest) Reset() {
	*x = UpdateMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_v1_metric_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricRequest) ProtoMessage() {}

func (x *UpdateMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metric_v1_metric_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricRequest) Descriptor() ([]byte, []int) {
	return file_metric_v1_metric_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateMetricRequest) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

var File_metric_v1_metric_proto protoreflect.FileDescriptor

var file_metric_v1_metric_proto_rawDesc = []byte{
	0x0a, 0x16, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe0, 0x01,
	0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x6e, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x42, 0x49, 0xba, 0x48, 0x46, 0xba, 0x01, 0x3a, 0x0a,
	0x0e, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12,
	0x1d, 0x74, 0x79, 0x70, 0x65, 0x20, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x20, 0x6d, 0x75, 0x73,
	0x74, 0x20, 0x62, 0x65, 0x20, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x64, 0x1a, 0x09,
	0x74, 0x68, 0x69, 0x73, 0x20, 0x21, 0x3d, 0x20, 0x30, 0x82, 0x01, 0x06, 0x10, 0x01, 0x1a, 0x02,
	0x01, 0x02, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18, 0x19,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x1d, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18, 0x64, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x22, 0x40, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x2a, 0x3e, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x47, 0x41, 0x55, 0x47, 0x45, 0x10, 0x01,
	0x12, 0x10, 0x0a, 0x0c, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x45, 0x52,
	0x10, 0x02, 0x32, 0x58, 0x0a, 0x0e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x12, 0x1e, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x91, 0x01, 0x0a,
	0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x0b,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x62, 0x61, 0x6b, 0x68, 0x74,
	0x69, 0x6e, 0x2f, 0x72, 0x74, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x2f, 0x76, 0x31, 0x3b, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x4d, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x09, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_metric_v1_metric_proto_rawDescOnce sync.Once
	file_metric_v1_metric_proto_rawDescData = file_metric_v1_metric_proto_rawDesc
)

func file_metric_v1_metric_proto_rawDescGZIP() []byte {
	file_metric_v1_metric_proto_rawDescOnce.Do(func() {
		file_metric_v1_metric_proto_rawDescData = protoimpl.X.CompressGZIP(file_metric_v1_metric_proto_rawDescData)
	})
	return file_metric_v1_metric_proto_rawDescData
}

var file_metric_v1_metric_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_metric_v1_metric_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_metric_v1_metric_proto_goTypes = []interface{}{
	(Type)(0),                   // 0: metric.v1.Type
	(*Metric)(nil),              // 1: metric.v1.Metric
	(*UpdateMetricRequest)(nil), // 2: metric.v1.UpdateMetricRequest
	(*emptypb.Empty)(nil),       // 3: google.protobuf.Empty
}
var file_metric_v1_metric_proto_depIdxs = []int32{
	0, // 0: metric.v1.Metric.type:type_name -> metric.v1.Type
	1, // 1: metric.v1.UpdateMetricRequest.metric:type_name -> metric.v1.Metric
	2, // 2: metric.v1.MetricsService.UpdateMetric:input_type -> metric.v1.UpdateMetricRequest
	3, // 3: metric.v1.MetricsService.UpdateMetric:output_type -> google.protobuf.Empty
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_metric_v1_metric_proto_init() }
func file_metric_v1_metric_proto_init() {
	if File_metric_v1_metric_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_metric_v1_metric_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
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
		file_metric_v1_metric_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricRequest); i {
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
			RawDescriptor: file_metric_v1_metric_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metric_v1_metric_proto_goTypes,
		DependencyIndexes: file_metric_v1_metric_proto_depIdxs,
		EnumInfos:         file_metric_v1_metric_proto_enumTypes,
		MessageInfos:      file_metric_v1_metric_proto_msgTypes,
	}.Build()
	File_metric_v1_metric_proto = out.File
	file_metric_v1_metric_proto_rawDesc = nil
	file_metric_v1_metric_proto_goTypes = nil
	file_metric_v1_metric_proto_depIdxs = nil
}
