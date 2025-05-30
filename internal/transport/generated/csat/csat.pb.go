// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: csat.proto

package csat

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

// ############### Survey ###############
type GetSurveyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSurveyRequest) Reset() {
	*x = GetSurveyRequest{}
	mi := &file_csat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSurveyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSurveyRequest) ProtoMessage() {}

func (x *GetSurveyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_csat_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSurveyRequest.ProtoReflect.Descriptor instead.
func (*GetSurveyRequest) Descriptor() ([]byte, []int) {
	return file_csat_proto_rawDescGZIP(), []int{0}
}

func (x *GetSurveyRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type SurveyWithQuestionsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SurveyId      string                 `protobuf:"bytes,1,opt,name=surveyId,proto3" json:"surveyId,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Questions     []*QuestionResponseDTO `protobuf:"bytes,4,rep,name=questions,proto3" json:"questions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SurveyWithQuestionsResponse) Reset() {
	*x = SurveyWithQuestionsResponse{}
	mi := &file_csat_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SurveyWithQuestionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SurveyWithQuestionsResponse) ProtoMessage() {}

func (x *SurveyWithQuestionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_csat_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SurveyWithQuestionsResponse.ProtoReflect.Descriptor instead.
func (*SurveyWithQuestionsResponse) Descriptor() ([]byte, []int) {
	return file_csat_proto_rawDescGZIP(), []int{1}
}

func (x *SurveyWithQuestionsResponse) GetSurveyId() string {
	if x != nil {
		return x.SurveyId
	}
	return ""
}

func (x *SurveyWithQuestionsResponse) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *SurveyWithQuestionsResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *SurveyWithQuestionsResponse) GetQuestions() []*QuestionResponseDTO {
	if x != nil {
		return x.Questions
	}
	return nil
}

type QuestionResponseDTO struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QuestionId    string                 `protobuf:"bytes,1,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"`
	Text          string                 `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QuestionResponseDTO) Reset() {
	*x = QuestionResponseDTO{}
	mi := &file_csat_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QuestionResponseDTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuestionResponseDTO) ProtoMessage() {}

func (x *QuestionResponseDTO) ProtoReflect() protoreflect.Message {
	mi := &file_csat_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuestionResponseDTO.ProtoReflect.Descriptor instead.
func (*QuestionResponseDTO) Descriptor() ([]byte, []int) {
	return file_csat_proto_rawDescGZIP(), []int{2}
}

func (x *QuestionResponseDTO) GetQuestionId() string {
	if x != nil {
		return x.QuestionId
	}
	return ""
}

func (x *QuestionResponseDTO) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type BriefSurvey struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BriefSurvey) Reset() {
	*x = BriefSurvey{}
	mi := &file_csat_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BriefSurvey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BriefSurvey) ProtoMessage() {}

func (x *BriefSurvey) ProtoReflect() protoreflect.Message {
	mi := &file_csat_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BriefSurvey.ProtoReflect.Descriptor instead.
func (*BriefSurvey) Descriptor() ([]byte, []int) {
	return file_csat_proto_rawDescGZIP(), []int{3}
}

func (x *BriefSurvey) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BriefSurvey) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

type SurveysList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Surveys       []*BriefSurvey         `protobuf:"bytes,1,rep,name=surveys,proto3" json:"surveys,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SurveysList) Reset() {
	*x = SurveysList{}
	mi := &file_csat_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SurveysList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SurveysList) ProtoMessage() {}

func (x *SurveysList) ProtoReflect() protoreflect.Message {
	mi := &file_csat_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SurveysList.ProtoReflect.Descriptor instead.
func (*SurveysList) Descriptor() ([]byte, []int) {
	return file_csat_proto_rawDescGZIP(), []int{4}
}

func (x *SurveysList) GetSurveys() []*BriefSurvey {
	if x != nil {
		return x.Surveys
	}
	return nil
}

// ############### Answers ###############
type SubmitAnswerRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SurveyId      string                 `protobuf:"bytes,1,opt,name=surveyId,proto3" json:"surveyId,omitempty"`
	Answers       []*AnswerRequestDTO    `protobuf:"bytes,2,rep,name=answers,proto3" json:"answers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitAnswerRequest) Reset() {
	*x = SubmitAnswerRequest{}
	mi := &file_csat_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitAnswerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitAnswerRequest) ProtoMessage() {}

func (x *SubmitAnswerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_csat_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitAnswerRequest.ProtoReflect.Descriptor instead.
func (*SubmitAnswerRequest) Descriptor() ([]byte, []int) {
	return file_csat_proto_rawDescGZIP(), []int{5}
}

func (x *SubmitAnswerRequest) GetSurveyId() string {
	if x != nil {
		return x.SurveyId
	}
	return ""
}

func (x *SubmitAnswerRequest) GetAnswers() []*AnswerRequestDTO {
	if x != nil {
		return x.Answers
	}
	return nil
}

type AnswerRequestDTO struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QuestionId    string                 `protobuf:"bytes,1,opt,name=questionId,proto3" json:"questionId,omitempty"`
	Value         uint32                 `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AnswerRequestDTO) Reset() {
	*x = AnswerRequestDTO{}
	mi := &file_csat_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AnswerRequestDTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnswerRequestDTO) ProtoMessage() {}

func (x *AnswerRequestDTO) ProtoReflect() protoreflect.Message {
	mi := &file_csat_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnswerRequestDTO.ProtoReflect.Descriptor instead.
func (*AnswerRequestDTO) Descriptor() ([]byte, []int) {
	return file_csat_proto_rawDescGZIP(), []int{6}
}

func (x *AnswerRequestDTO) GetQuestionId() string {
	if x != nil {
		return x.QuestionId
	}
	return ""
}

func (x *AnswerRequestDTO) GetValue() uint32 {
	if x != nil {
		return x.Value
	}
	return 0
}

// ############### Statistics ###############
type GetStatisticsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SurveyId      string                 `protobuf:"bytes,1,opt,name=surveyId,proto3" json:"surveyId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetStatisticsRequest) Reset() {
	*x = GetStatisticsRequest{}
	mi := &file_csat_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStatisticsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStatisticsRequest) ProtoMessage() {}

func (x *GetStatisticsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_csat_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStatisticsRequest.ProtoReflect.Descriptor instead.
func (*GetStatisticsRequest) Descriptor() ([]byte, []int) {
	return file_csat_proto_rawDescGZIP(), []int{7}
}

func (x *GetStatisticsRequest) GetSurveyId() string {
	if x != nil {
		return x.SurveyId
	}
	return ""
}

type QuestionStatisticsDTO struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QuestionId    string                 `protobuf:"bytes,1,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"`
	Text          string                 `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	Stats         []uint32               `protobuf:"varint,3,rep,packed,name=stats,proto3" json:"stats,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QuestionStatisticsDTO) Reset() {
	*x = QuestionStatisticsDTO{}
	mi := &file_csat_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QuestionStatisticsDTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuestionStatisticsDTO) ProtoMessage() {}

func (x *QuestionStatisticsDTO) ProtoReflect() protoreflect.Message {
	mi := &file_csat_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuestionStatisticsDTO.ProtoReflect.Descriptor instead.
func (*QuestionStatisticsDTO) Descriptor() ([]byte, []int) {
	return file_csat_proto_rawDescGZIP(), []int{8}
}

func (x *QuestionStatisticsDTO) GetQuestionId() string {
	if x != nil {
		return x.QuestionId
	}
	return ""
}

func (x *QuestionStatisticsDTO) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *QuestionStatisticsDTO) GetStats() []uint32 {
	if x != nil {
		return x.Stats
	}
	return nil
}

type SurveyStatisticsResponse struct {
	state         protoimpl.MessageState   `protogen:"open.v1"`
	Description   string                   `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	Questions     []*QuestionStatisticsDTO `protobuf:"bytes,2,rep,name=questions,proto3" json:"questions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SurveyStatisticsResponse) Reset() {
	*x = SurveyStatisticsResponse{}
	mi := &file_csat_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SurveyStatisticsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SurveyStatisticsResponse) ProtoMessage() {}

func (x *SurveyStatisticsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_csat_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SurveyStatisticsResponse.ProtoReflect.Descriptor instead.
func (*SurveyStatisticsResponse) Descriptor() ([]byte, []int) {
	return file_csat_proto_rawDescGZIP(), []int{9}
}

func (x *SurveyStatisticsResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *SurveyStatisticsResponse) GetQuestions() []*QuestionStatisticsDTO {
	if x != nil {
		return x.Questions
	}
	return nil
}

var File_csat_proto protoreflect.FileDescriptor

const file_csat_proto_rawDesc = "" +
	"\n" +
	"\n" +
	"csat.proto\x12\x04csat\x1a\x1bgoogle/protobuf/empty.proto\"&\n" +
	"\x10GetSurveyRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\"\xaa\x01\n" +
	"\x1bSurveyWithQuestionsResponse\x12\x1a\n" +
	"\bsurveyId\x18\x01 \x01(\tR\bsurveyId\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x127\n" +
	"\tquestions\x18\x04 \x03(\v2\x19.csat.QuestionResponseDTOR\tquestions\"J\n" +
	"\x13QuestionResponseDTO\x12\x1f\n" +
	"\vquestion_id\x18\x01 \x01(\tR\n" +
	"questionId\x12\x12\n" +
	"\x04text\x18\x02 \x01(\tR\x04text\"3\n" +
	"\vBriefSurvey\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\":\n" +
	"\vSurveysList\x12+\n" +
	"\asurveys\x18\x01 \x03(\v2\x11.csat.BriefSurveyR\asurveys\"c\n" +
	"\x13SubmitAnswerRequest\x12\x1a\n" +
	"\bsurveyId\x18\x01 \x01(\tR\bsurveyId\x120\n" +
	"\aanswers\x18\x02 \x03(\v2\x16.csat.AnswerRequestDTOR\aanswers\"H\n" +
	"\x10AnswerRequestDTO\x12\x1e\n" +
	"\n" +
	"questionId\x18\x01 \x01(\tR\n" +
	"questionId\x12\x14\n" +
	"\x05value\x18\x02 \x01(\rR\x05value\"2\n" +
	"\x14GetStatisticsRequest\x12\x1a\n" +
	"\bsurveyId\x18\x01 \x01(\tR\bsurveyId\"b\n" +
	"\x15QuestionStatisticsDTO\x12\x1f\n" +
	"\vquestion_id\x18\x01 \x01(\tR\n" +
	"questionId\x12\x12\n" +
	"\x04text\x18\x02 \x01(\tR\x04text\x12\x14\n" +
	"\x05stats\x18\x03 \x03(\rR\x05stats\"w\n" +
	"\x18SurveyStatisticsResponse\x12 \n" +
	"\vdescription\x18\x01 \x01(\tR\vdescription\x129\n" +
	"\tquestions\x18\x02 \x03(\v2\x1b.csat.QuestionStatisticsDTOR\tquestions2\xb6\x02\n" +
	"\rSurveyService\x12S\n" +
	"\x16GetSurveyWithQuestions\x12\x16.csat.GetSurveyRequest\x1a!.csat.SurveyWithQuestionsResponse\x12A\n" +
	"\fSubmitAnswer\x12\x19.csat.SubmitAnswerRequest\x1a\x16.google.protobuf.Empty\x12Q\n" +
	"\x13GetSurveyStatistics\x12\x1a.csat.GetStatisticsRequest\x1a\x1e.csat.SurveyStatisticsResponse\x12:\n" +
	"\rGetAllSurveys\x12\x16.google.protobuf.Empty\x1a\x11.csat.SurveysListB4Z22025_1_ChillGuys/internal/transport/generated/csatb\x06proto3"

var (
	file_csat_proto_rawDescOnce sync.Once
	file_csat_proto_rawDescData []byte
)

func file_csat_proto_rawDescGZIP() []byte {
	file_csat_proto_rawDescOnce.Do(func() {
		file_csat_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_csat_proto_rawDesc), len(file_csat_proto_rawDesc)))
	})
	return file_csat_proto_rawDescData
}

var file_csat_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_csat_proto_goTypes = []any{
	(*GetSurveyRequest)(nil),            // 0: csat.GetSurveyRequest
	(*SurveyWithQuestionsResponse)(nil), // 1: csat.SurveyWithQuestionsResponse
	(*QuestionResponseDTO)(nil),         // 2: csat.QuestionResponseDTO
	(*BriefSurvey)(nil),                 // 3: csat.BriefSurvey
	(*SurveysList)(nil),                 // 4: csat.SurveysList
	(*SubmitAnswerRequest)(nil),         // 5: csat.SubmitAnswerRequest
	(*AnswerRequestDTO)(nil),            // 6: csat.AnswerRequestDTO
	(*GetStatisticsRequest)(nil),        // 7: csat.GetStatisticsRequest
	(*QuestionStatisticsDTO)(nil),       // 8: csat.QuestionStatisticsDTO
	(*SurveyStatisticsResponse)(nil),    // 9: csat.SurveyStatisticsResponse
	(*emptypb.Empty)(nil),               // 10: google.protobuf.Empty
}
var file_csat_proto_depIdxs = []int32{
	2,  // 0: csat.SurveyWithQuestionsResponse.questions:type_name -> csat.QuestionResponseDTO
	3,  // 1: csat.SurveysList.surveys:type_name -> csat.BriefSurvey
	6,  // 2: csat.SubmitAnswerRequest.answers:type_name -> csat.AnswerRequestDTO
	8,  // 3: csat.SurveyStatisticsResponse.questions:type_name -> csat.QuestionStatisticsDTO
	0,  // 4: csat.SurveyService.GetSurveyWithQuestions:input_type -> csat.GetSurveyRequest
	5,  // 5: csat.SurveyService.SubmitAnswer:input_type -> csat.SubmitAnswerRequest
	7,  // 6: csat.SurveyService.GetSurveyStatistics:input_type -> csat.GetStatisticsRequest
	10, // 7: csat.SurveyService.GetAllSurveys:input_type -> google.protobuf.Empty
	1,  // 8: csat.SurveyService.GetSurveyWithQuestions:output_type -> csat.SurveyWithQuestionsResponse
	10, // 9: csat.SurveyService.SubmitAnswer:output_type -> google.protobuf.Empty
	9,  // 10: csat.SurveyService.GetSurveyStatistics:output_type -> csat.SurveyStatisticsResponse
	4,  // 11: csat.SurveyService.GetAllSurveys:output_type -> csat.SurveysList
	8,  // [8:12] is the sub-list for method output_type
	4,  // [4:8] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_csat_proto_init() }
func file_csat_proto_init() {
	if File_csat_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_csat_proto_rawDesc), len(file_csat_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_csat_proto_goTypes,
		DependencyIndexes: file_csat_proto_depIdxs,
		MessageInfos:      file_csat_proto_msgTypes,
	}.Build()
	File_csat_proto = out.File
	file_csat_proto_goTypes = nil
	file_csat_proto_depIdxs = nil
}
