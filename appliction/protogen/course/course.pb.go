// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: course/course.proto

package course

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type Course struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CourseId      string                 `protobuf:"bytes,1,opt,name=course_id,json=courseId,proto3" json:"course_id,omitempty"`
	CourseTitle   string                 `protobuf:"bytes,2,opt,name=course_title,json=courseTitle,proto3" json:"course_title,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Course) Reset() {
	*x = Course{}
	mi := &file_course_course_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Course) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Course) ProtoMessage() {}

func (x *Course) ProtoReflect() protoreflect.Message {
	mi := &file_course_course_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Course.ProtoReflect.Descriptor instead.
func (*Course) Descriptor() ([]byte, []int) {
	return file_course_course_proto_rawDescGZIP(), []int{0}
}

func (x *Course) GetCourseId() string {
	if x != nil {
		return x.CourseId
	}
	return ""
}

func (x *Course) GetCourseTitle() string {
	if x != nil {
		return x.CourseTitle
	}
	return ""
}

func (x *Course) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Course) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_course_course_proto protoreflect.FileDescriptor

const file_course_course_proto_rawDesc = "" +
	"\n" +
	"\x13course/course.proto\x12\n" +
	"lms.course\x1a\x1fgoogle/protobuf/timestamp.proto\"\xbe\x01\n" +
	"\x06Course\x12\x1b\n" +
	"\tcourse_id\x18\x01 \x01(\tR\bcourseId\x12!\n" +
	"\fcourse_title\x18\x02 \x01(\tR\vcourseTitle\x129\n" +
	"\n" +
	"created_at\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAtB<Z:github.com/multi-tenant-cms-golang/lms-sys/protogen/courseb\x06proto3"

var (
	file_course_course_proto_rawDescOnce sync.Once
	file_course_course_proto_rawDescData []byte
)

func file_course_course_proto_rawDescGZIP() []byte {
	file_course_course_proto_rawDescOnce.Do(func() {
		file_course_course_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_course_course_proto_rawDesc), len(file_course_course_proto_rawDesc)))
	})
	return file_course_course_proto_rawDescData
}

var file_course_course_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_course_course_proto_goTypes = []any{
	(*Course)(nil),                // 0: lms.course.Course
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_course_course_proto_depIdxs = []int32{
	1, // 0: lms.course.Course.created_at:type_name -> google.protobuf.Timestamp
	1, // 1: lms.course.Course.updated_at:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_course_course_proto_init() }
func file_course_course_proto_init() {
	if File_course_course_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_course_course_proto_rawDesc), len(file_course_course_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_course_course_proto_goTypes,
		DependencyIndexes: file_course_course_proto_depIdxs,
		MessageInfos:      file_course_course_proto_msgTypes,
	}.Build()
	File_course_course_proto = out.File
	file_course_course_proto_goTypes = nil
	file_course_course_proto_depIdxs = nil
}
