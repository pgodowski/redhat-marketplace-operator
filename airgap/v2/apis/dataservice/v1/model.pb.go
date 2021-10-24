// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.1
// source: dataservice/v1/model.proto

package dataservicev1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// represents the file info
type FileInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name       string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Size       uint32                 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	Source     string                 `protobuf:"bytes,4,opt,name=source,proto3" json:"source,omitempty"`
	SourceType string                 `protobuf:"bytes,5,opt,name=sourceType,proto3" json:"sourceType,omitempty"`
	Checksum   string                 `protobuf:"bytes,10,opt,name=checksum,proto3" json:"checksum,omitempty"`
	MimeType   string                 `protobuf:"bytes,11,opt,name=mimeType,proto3" json:"mimeType,omitempty"`
	CreatedAt  *timestamppb.Timestamp `protobuf:"bytes,15,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt  *timestamppb.Timestamp `protobuf:"bytes,16,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt  *timestamppb.Timestamp `protobuf:"bytes,17,opt,name=deleted_at,json=deletedAt,proto3,oneof" json:"deleted_at,omitempty"`
	Metadata   map[string]string      `protobuf:"bytes,20,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *FileInfo) Reset() {
	*x = FileInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dataservice_v1_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileInfo) ProtoMessage() {}

func (x *FileInfo) ProtoReflect() protoreflect.Message {
	mi := &file_dataservice_v1_model_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileInfo.ProtoReflect.Descriptor instead.
func (*FileInfo) Descriptor() ([]byte, []int) {
	return file_dataservice_v1_model_proto_rawDescGZIP(), []int{0}
}

func (x *FileInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *FileInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FileInfo) GetSize() uint32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *FileInfo) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *FileInfo) GetSourceType() string {
	if x != nil {
		return x.SourceType
	}
	return ""
}

func (x *FileInfo) GetChecksum() string {
	if x != nil {
		return x.Checksum
	}
	return ""
}

func (x *FileInfo) GetMimeType() string {
	if x != nil {
		return x.MimeType
	}
	return ""
}

func (x *FileInfo) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *FileInfo) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *FileInfo) GetDeletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *FileInfo) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type FileKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Source     string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	SourceType string `protobuf:"bytes,3,opt,name=sourceType,proto3" json:"sourceType,omitempty"`
}

func (x *FileKey) Reset() {
	*x = FileKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dataservice_v1_model_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileKey) ProtoMessage() {}

func (x *FileKey) ProtoReflect() protoreflect.Message {
	mi := &file_dataservice_v1_model_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileKey.ProtoReflect.Descriptor instead.
func (*FileKey) Descriptor() ([]byte, []int) {
	return file_dataservice_v1_model_proto_rawDescGZIP(), []int{1}
}

func (x *FileKey) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FileKey) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *FileKey) GetSourceType() string {
	if x != nil {
		return x.SourceType
	}
	return ""
}

var File_dataservice_v1_model_proto protoreflect.FileDescriptor

var file_dataservice_v1_model_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x64, 0x61,
	0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f,
	0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8a,
	0x04, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x73,
	0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x3f, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x42, 0x04, 0xe2, 0x41, 0x01, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x3f, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x10, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x42, 0x04, 0xe2, 0x41, 0x01, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x44, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x11, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x04, 0xe2, 0x41, 0x01, 0x03, 0x48, 0x00, 0x52, 0x09, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x14, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a,
	0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0d, 0x0a, 0x0b,
	0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x22, 0x55, 0x0a, 0x07, 0x46,
	0x69, 0x6c, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x42, 0xe0, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x4d, 0x6f, 0x64, 0x65, 0x6c,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x65, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x65, 0x64, 0x68, 0x61, 0x74, 0x2d, 0x6d, 0x61, 0x72, 0x6b, 0x65,
	0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x2f, 0x72, 0x65, 0x64, 0x68, 0x61, 0x74, 0x2d, 0x6d, 0x61,
	0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x2d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x2f, 0x61, 0x69, 0x72, 0x67, 0x61, 0x70, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69,
	0x73, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31,
	0x3b, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x44, 0x58, 0x58, 0xaa, 0x02, 0x0e, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0e, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1a, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x0f, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dataservice_v1_model_proto_rawDescOnce sync.Once
	file_dataservice_v1_model_proto_rawDescData = file_dataservice_v1_model_proto_rawDesc
)

func file_dataservice_v1_model_proto_rawDescGZIP() []byte {
	file_dataservice_v1_model_proto_rawDescOnce.Do(func() {
		file_dataservice_v1_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_dataservice_v1_model_proto_rawDescData)
	})
	return file_dataservice_v1_model_proto_rawDescData
}

var file_dataservice_v1_model_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_dataservice_v1_model_proto_goTypes = []interface{}{
	(*FileInfo)(nil),              // 0: dataservice.v1.FileInfo
	(*FileKey)(nil),               // 1: dataservice.v1.FileKey
	nil,                           // 2: dataservice.v1.FileInfo.MetadataEntry
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_dataservice_v1_model_proto_depIdxs = []int32{
	3, // 0: dataservice.v1.FileInfo.created_at:type_name -> google.protobuf.Timestamp
	3, // 1: dataservice.v1.FileInfo.updated_at:type_name -> google.protobuf.Timestamp
	3, // 2: dataservice.v1.FileInfo.deleted_at:type_name -> google.protobuf.Timestamp
	2, // 3: dataservice.v1.FileInfo.metadata:type_name -> dataservice.v1.FileInfo.MetadataEntry
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_dataservice_v1_model_proto_init() }
func file_dataservice_v1_model_proto_init() {
	if File_dataservice_v1_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dataservice_v1_model_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileInfo); i {
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
		file_dataservice_v1_model_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileKey); i {
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
	file_dataservice_v1_model_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_dataservice_v1_model_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_dataservice_v1_model_proto_goTypes,
		DependencyIndexes: file_dataservice_v1_model_proto_depIdxs,
		MessageInfos:      file_dataservice_v1_model_proto_msgTypes,
	}.Build()
	File_dataservice_v1_model_proto = out.File
	file_dataservice_v1_model_proto_rawDesc = nil
	file_dataservice_v1_model_proto_goTypes = nil
	file_dataservice_v1_model_proto_depIdxs = nil
}
