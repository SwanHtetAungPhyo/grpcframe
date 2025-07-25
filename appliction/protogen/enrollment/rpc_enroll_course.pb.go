// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: enrollment/rpc_enroll_course.proto

package enrollment

import (
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

type EnrollCourseRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StudentId     string                 `protobuf:"bytes,1,opt,name=student_id,json=studentId,proto3" json:"student_id,omitempty"`
	CourseId      string                 `protobuf:"bytes,2,opt,name=course_id,json=courseId,proto3" json:"course_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EnrollCourseRequest) Reset() {
	*x = EnrollCourseRequest{}
	mi := &file_enrollment_rpc_enroll_course_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnrollCourseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnrollCourseRequest) ProtoMessage() {}

func (x *EnrollCourseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_enrollment_rpc_enroll_course_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnrollCourseRequest.ProtoReflect.Descriptor instead.
func (*EnrollCourseRequest) Descriptor() ([]byte, []int) {
	return file_enrollment_rpc_enroll_course_proto_rawDescGZIP(), []int{0}
}

func (x *EnrollCourseRequest) GetStudentId() string {
	if x != nil {
		return x.StudentId
	}
	return ""
}

func (x *EnrollCourseRequest) GetCourseId() string {
	if x != nil {
		return x.CourseId
	}
	return ""
}

type EnrollCourseResponse struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	EnrollmentId   string                 `protobuf:"bytes,1,opt,name=enrollment_id,json=enrollmentId,proto3" json:"enrollment_id,omitempty"`
	StudentId      string                 `protobuf:"bytes,2,opt,name=student_id,json=studentId,proto3" json:"student_id,omitempty"`
	CourseId       string                 `protobuf:"bytes,3,opt,name=course_id,json=courseId,proto3" json:"course_id,omitempty"`
	EnrollmentDate string                 `protobuf:"bytes,4,opt,name=enrollment_date,json=enrollmentDate,proto3" json:"enrollment_date,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *EnrollCourseResponse) Reset() {
	*x = EnrollCourseResponse{}
	mi := &file_enrollment_rpc_enroll_course_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnrollCourseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnrollCourseResponse) ProtoMessage() {}

func (x *EnrollCourseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_enrollment_rpc_enroll_course_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnrollCourseResponse.ProtoReflect.Descriptor instead.
func (*EnrollCourseResponse) Descriptor() ([]byte, []int) {
	return file_enrollment_rpc_enroll_course_proto_rawDescGZIP(), []int{1}
}

func (x *EnrollCourseResponse) GetEnrollmentId() string {
	if x != nil {
		return x.EnrollmentId
	}
	return ""
}

func (x *EnrollCourseResponse) GetStudentId() string {
	if x != nil {
		return x.StudentId
	}
	return ""
}

func (x *EnrollCourseResponse) GetCourseId() string {
	if x != nil {
		return x.CourseId
	}
	return ""
}

func (x *EnrollCourseResponse) GetEnrollmentDate() string {
	if x != nil {
		return x.EnrollmentDate
	}
	return ""
}

var File_enrollment_rpc_enroll_course_proto protoreflect.FileDescriptor

const file_enrollment_rpc_enroll_course_proto_rawDesc = "" +
	"\n" +
	"\"enrollment/rpc_enroll_course.proto\x12\x0elms.enrollment\x1a\x1benrollment/enrollment.proto\"Q\n" +
	"\x13EnrollCourseRequest\x12\x1d\n" +
	"\n" +
	"student_id\x18\x01 \x01(\tR\tstudentId\x12\x1b\n" +
	"\tcourse_id\x18\x02 \x01(\tR\bcourseId\"\xa0\x01\n" +
	"\x14EnrollCourseResponse\x12#\n" +
	"\renrollment_id\x18\x01 \x01(\tR\fenrollmentId\x12\x1d\n" +
	"\n" +
	"student_id\x18\x02 \x01(\tR\tstudentId\x12\x1b\n" +
	"\tcourse_id\x18\x03 \x01(\tR\bcourseId\x12'\n" +
	"\x0fenrollment_date\x18\x04 \x01(\tR\x0eenrollmentDateB@Z>github.com/multi-tenant-cms-golang/lms-sys/protogen/enrollmentb\x06proto3"

var (
	file_enrollment_rpc_enroll_course_proto_rawDescOnce sync.Once
	file_enrollment_rpc_enroll_course_proto_rawDescData []byte
)

func file_enrollment_rpc_enroll_course_proto_rawDescGZIP() []byte {
	file_enrollment_rpc_enroll_course_proto_rawDescOnce.Do(func() {
		file_enrollment_rpc_enroll_course_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_enrollment_rpc_enroll_course_proto_rawDesc), len(file_enrollment_rpc_enroll_course_proto_rawDesc)))
	})
	return file_enrollment_rpc_enroll_course_proto_rawDescData
}

var file_enrollment_rpc_enroll_course_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_enrollment_rpc_enroll_course_proto_goTypes = []any{
	(*EnrollCourseRequest)(nil),  // 0: lms.enrollment.EnrollCourseRequest
	(*EnrollCourseResponse)(nil), // 1: lms.enrollment.EnrollCourseResponse
}
var file_enrollment_rpc_enroll_course_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_enrollment_rpc_enroll_course_proto_init() }
func file_enrollment_rpc_enroll_course_proto_init() {
	if File_enrollment_rpc_enroll_course_proto != nil {
		return
	}
	file_enrollment_enrollment_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_enrollment_rpc_enroll_course_proto_rawDesc), len(file_enrollment_rpc_enroll_course_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_enrollment_rpc_enroll_course_proto_goTypes,
		DependencyIndexes: file_enrollment_rpc_enroll_course_proto_depIdxs,
		MessageInfos:      file_enrollment_rpc_enroll_course_proto_msgTypes,
	}.Build()
	File_enrollment_rpc_enroll_course_proto = out.File
	file_enrollment_rpc_enroll_course_proto_goTypes = nil
	file_enrollment_rpc_enroll_course_proto_depIdxs = nil
}
