// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: payment-service/v1/payment-service.proto

package payment_servicev1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EchoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EchoRequest) Reset() {
	*x = EchoRequest{}
	mi := &file_payment_service_v1_payment_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EchoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoRequest) ProtoMessage() {}

func (x *EchoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_payment_service_v1_payment_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoRequest.ProtoReflect.Descriptor instead.
func (*EchoRequest) Descriptor() ([]byte, []int) {
	return file_payment_service_v1_payment_service_proto_rawDescGZIP(), []int{0}
}

func (x *EchoRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type EchoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EchoResponse) Reset() {
	*x = EchoResponse{}
	mi := &file_payment_service_v1_payment_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EchoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoResponse) ProtoMessage() {}

func (x *EchoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_payment_service_v1_payment_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoResponse.ProtoReflect.Descriptor instead.
func (*EchoResponse) Descriptor() ([]byte, []int) {
	return file_payment_service_v1_payment_service_proto_rawDescGZIP(), []int{1}
}

func (x *EchoResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_payment_service_v1_payment_service_proto protoreflect.FileDescriptor

const file_payment_service_v1_payment_service_proto_rawDesc = "" +
	"\n" +
	"(payment-service/v1/payment-service.proto\x12\x12payment_service.v1\x1a\x1cgoogle/api/annotations.proto\"'\n" +
	"\vEchoRequest\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"(\n" +
	"\fEchoResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage2[\n" +
	"\x0ePaymentService\x12I\n" +
	"\x04Echo\x12\x1f.payment_service.v1.EchoRequest\x1a .payment_service.v1.EchoResponseB\xe0\x01\n" +
	"\x16com.payment_service.v1B\x13PaymentServiceProtoP\x01ZLgithub.com/MaxFando/lms/payment-service/payment-service/v1;payment_servicev1\xa2\x02\x03PXX\xaa\x02\x11PaymentService.V1\xca\x02\x11PaymentService\\V1\xe2\x02\x1dPaymentService\\V1\\GPBMetadata\xea\x02\x12PaymentService::V1b\x06proto3"

var (
	file_payment_service_v1_payment_service_proto_rawDescOnce sync.Once
	file_payment_service_v1_payment_service_proto_rawDescData []byte
)

func file_payment_service_v1_payment_service_proto_rawDescGZIP() []byte {
	file_payment_service_v1_payment_service_proto_rawDescOnce.Do(func() {
		file_payment_service_v1_payment_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_payment_service_v1_payment_service_proto_rawDesc), len(file_payment_service_v1_payment_service_proto_rawDesc)))
	})
	return file_payment_service_v1_payment_service_proto_rawDescData
}

var file_payment_service_v1_payment_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_payment_service_v1_payment_service_proto_goTypes = []any{
	(*EchoRequest)(nil),  // 0: payment_service.v1.EchoRequest
	(*EchoResponse)(nil), // 1: payment_service.v1.EchoResponse
}
var file_payment_service_v1_payment_service_proto_depIdxs = []int32{
	0, // 0: payment_service.v1.PaymentService.Echo:input_type -> payment_service.v1.EchoRequest
	1, // 1: payment_service.v1.PaymentService.Echo:output_type -> payment_service.v1.EchoResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_payment_service_v1_payment_service_proto_init() }
func file_payment_service_v1_payment_service_proto_init() {
	if File_payment_service_v1_payment_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_payment_service_v1_payment_service_proto_rawDesc), len(file_payment_service_v1_payment_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_payment_service_v1_payment_service_proto_goTypes,
		DependencyIndexes: file_payment_service_v1_payment_service_proto_depIdxs,
		MessageInfos:      file_payment_service_v1_payment_service_proto_msgTypes,
	}.Build()
	File_payment_service_v1_payment_service_proto = out.File
	file_payment_service_v1_payment_service_proto_goTypes = nil
	file_payment_service_v1_payment_service_proto_depIdxs = nil
}
