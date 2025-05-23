// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v4.22.2
// source: proto/api.proto

package calculator_proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

// Ответ с задачей и флагом о готовности задачи
type TaskResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Флаг, говорящий о готовности задачи
	Enabled bool `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	// Задача, ели enabled == true
	Task          *Task `protobuf:"bytes,2,opt,name=task,proto3" json:"task,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TaskResponse) Reset() {
	*x = TaskResponse{}
	mi := &file_proto_api_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskResponse) ProtoMessage() {}

func (x *TaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskResponse.ProtoReflect.Descriptor instead.
func (*TaskResponse) Descriptor() ([]byte, []int) {
	return file_proto_api_proto_rawDescGZIP(), []int{0}
}

func (x *TaskResponse) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *TaskResponse) GetTask() *Task {
	if x != nil {
		return x.Task
	}
	return nil
}

// Задача для агента
type Task struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ExpId         string                 `protobuf:"bytes,2,opt,name=expId,proto3" json:"expId,omitempty"`
	Arg1          float64                `protobuf:"fixed64,3,opt,name=arg1,proto3" json:"arg1,omitempty"`
	Arg2          float64                `protobuf:"fixed64,4,opt,name=arg2,proto3" json:"arg2,omitempty"`
	Operation     string                 `protobuf:"bytes,5,opt,name=operation,proto3" json:"operation,omitempty"`
	OperationTime *durationpb.Duration   `protobuf:"bytes,6,opt,name=operationTime,proto3" json:"operationTime,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Task) Reset() {
	*x = Task{}
	mi := &file_proto_api_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_proto_api_proto_rawDescGZIP(), []int{1}
}

func (x *Task) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Task) GetExpId() string {
	if x != nil {
		return x.ExpId
	}
	return ""
}

func (x *Task) GetArg1() float64 {
	if x != nil {
		return x.Arg1
	}
	return 0
}

func (x *Task) GetArg2() float64 {
	if x != nil {
		return x.Arg2
	}
	return 0
}

func (x *Task) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

func (x *Task) GetOperationTime() *durationpb.Duration {
	if x != nil {
		return x.OperationTime
	}
	return nil
}

type TaskResult struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Result        float64                `protobuf:"fixed64,2,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TaskResult) Reset() {
	*x = TaskResult{}
	mi := &file_proto_api_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskResult) ProtoMessage() {}

func (x *TaskResult) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskResult.ProtoReflect.Descriptor instead.
func (*TaskResult) Descriptor() ([]byte, []int) {
	return file_proto_api_proto_rawDescGZIP(), []int{2}
}

func (x *TaskResult) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TaskResult) GetResult() float64 {
	if x != nil {
		return x.Result
	}
	return 0
}

var File_proto_api_proto protoreflect.FileDescriptor

const file_proto_api_proto_rawDesc = "" +
	"\n" +
	"\x0fproto/api.proto\x12\n" +
	"calculator\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1egoogle/protobuf/duration.proto\"N\n" +
	"\fTaskResponse\x12\x18\n" +
	"\aenabled\x18\x01 \x01(\bR\aenabled\x12$\n" +
	"\x04task\x18\x02 \x01(\v2\x10.calculator.TaskR\x04task\"\xb3\x01\n" +
	"\x04Task\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x14\n" +
	"\x05expId\x18\x02 \x01(\tR\x05expId\x12\x12\n" +
	"\x04arg1\x18\x03 \x01(\x01R\x04arg1\x12\x12\n" +
	"\x04arg2\x18\x04 \x01(\x01R\x04arg2\x12\x1c\n" +
	"\toperation\x18\x05 \x01(\tR\toperation\x12?\n" +
	"\roperationTime\x18\x06 \x01(\v2\x19.google.protobuf.DurationR\roperationTime\"4\n" +
	"\n" +
	"TaskResult\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x16\n" +
	"\x06result\x18\x02 \x01(\x01R\x06result2\x84\x01\n" +
	"\n" +
	"APIService\x12;\n" +
	"\aGetTask\x12\x16.google.protobuf.Empty\x1a\x18.calculator.TaskResponse\x129\n" +
	"\aSetTask\x12\x16.calculator.TaskResult\x1a\x16.google.protobuf.EmptyB\x12Z\x10calculator.protob\x06proto3"

var (
	file_proto_api_proto_rawDescOnce sync.Once
	file_proto_api_proto_rawDescData []byte
)

func file_proto_api_proto_rawDescGZIP() []byte {
	file_proto_api_proto_rawDescOnce.Do(func() {
		file_proto_api_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_api_proto_rawDesc), len(file_proto_api_proto_rawDesc)))
	})
	return file_proto_api_proto_rawDescData
}

var file_proto_api_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_api_proto_goTypes = []any{
	(*TaskResponse)(nil),        // 0: calculator.TaskResponse
	(*Task)(nil),                // 1: calculator.Task
	(*TaskResult)(nil),          // 2: calculator.TaskResult
	(*durationpb.Duration)(nil), // 3: google.protobuf.Duration
	(*emptypb.Empty)(nil),       // 4: google.protobuf.Empty
}
var file_proto_api_proto_depIdxs = []int32{
	1, // 0: calculator.TaskResponse.task:type_name -> calculator.Task
	3, // 1: calculator.Task.operationTime:type_name -> google.protobuf.Duration
	4, // 2: calculator.APIService.GetTask:input_type -> google.protobuf.Empty
	2, // 3: calculator.APIService.SetTask:input_type -> calculator.TaskResult
	0, // 4: calculator.APIService.GetTask:output_type -> calculator.TaskResponse
	4, // 5: calculator.APIService.SetTask:output_type -> google.protobuf.Empty
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_api_proto_init() }
func file_proto_api_proto_init() {
	if File_proto_api_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_api_proto_rawDesc), len(file_proto_api_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_api_proto_goTypes,
		DependencyIndexes: file_proto_api_proto_depIdxs,
		MessageInfos:      file_proto_api_proto_msgTypes,
	}.Build()
	File_proto_api_proto = out.File
	file_proto_api_proto_goTypes = nil
	file_proto_api_proto_depIdxs = nil
}
