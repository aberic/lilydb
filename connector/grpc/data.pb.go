// Code generated by protoc-gen-go. DO NOT EDIT.
// source: connector/grpc/data.proto

// protoc --java_out=plugins=grpc:. api/*.proto
// protoc --go_out=plugins=grpc:. api/*.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// FormType 表类型
type FormType int32

const (
	// SQL 关联数据库类型
	FormType_SQL FormType = 0
	// Doc 文档数据库类型
	FormType_Doc FormType = 1
)

var FormType_name = map[int32]string{
	0: "SQL",
	1: "Doc",
}

var FormType_value = map[string]int32{
	"SQL": 0,
	"Doc": 1,
}

func (x FormType) String() string {
	return proto.EnumName(FormType_name, int32(x))
}

func (FormType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_43e42cbf821258b1, []int{0}
}

// Lily 数据库引擎对象
type Lily struct {
	// databases 数据库集合
	Databases            map[string]*Database `protobuf:"bytes,1,rep,name=databases,proto3" json:"databases,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Lily) Reset()         { *m = Lily{} }
func (m *Lily) String() string { return proto.CompactTextString(m) }
func (*Lily) ProtoMessage()    {}
func (*Lily) Descriptor() ([]byte, []int) {
	return fileDescriptor_43e42cbf821258b1, []int{0}
}

func (m *Lily) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Lily.Unmarshal(m, b)
}
func (m *Lily) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Lily.Marshal(b, m, deterministic)
}
func (m *Lily) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Lily.Merge(m, src)
}
func (m *Lily) XXX_Size() int {
	return xxx_messageInfo_Lily.Size(m)
}
func (m *Lily) XXX_DiscardUnknown() {
	xxx_messageInfo_Lily.DiscardUnknown(m)
}

var xxx_messageInfo_Lily proto.InternalMessageInfo

func (m *Lily) GetDatabases() map[string]*Database {
	if m != nil {
		return m.Databases
	}
	return nil
}

// Database 数据库对象
type Database struct {
	// ID 数据库唯一ID，不能改变
	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	// Name 数据库名称，根据需求可以随时变化
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	// Comment 数据库描述
	Comment string `protobuf:"bytes,3,opt,name=Comment,proto3" json:"Comment,omitempty"`
	// Forms 数据库表集合
	Forms                map[string]*Form `protobuf:"bytes,4,rep,name=Forms,proto3" json:"Forms,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Database) Reset()         { *m = Database{} }
func (m *Database) String() string { return proto.CompactTextString(m) }
func (*Database) ProtoMessage()    {}
func (*Database) Descriptor() ([]byte, []int) {
	return fileDescriptor_43e42cbf821258b1, []int{1}
}

func (m *Database) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Database.Unmarshal(m, b)
}
func (m *Database) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Database.Marshal(b, m, deterministic)
}
func (m *Database) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Database.Merge(m, src)
}
func (m *Database) XXX_Size() int {
	return xxx_messageInfo_Database.Size(m)
}
func (m *Database) XXX_DiscardUnknown() {
	xxx_messageInfo_Database.DiscardUnknown(m)
}

var xxx_messageInfo_Database proto.InternalMessageInfo

func (m *Database) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Database) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Database) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *Database) GetForms() map[string]*Form {
	if m != nil {
		return m.Forms
	}
	return nil
}

// Form 数据库表对象
type Form struct {
	// ID 表唯一ID，不能改变
	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	// Name 表名，根据需求可以随时变化
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	// Comment 表描述
	Comment string `protobuf:"bytes,3,opt,name=Comment,proto3" json:"Comment,omitempty"`
	// FormType 表类型 SQL/Doc
	FormType FormType `protobuf:"varint,4,opt,name=FormType,proto3,enum=api.FormType" json:"FormType,omitempty"`
	// Indexes 索引ID集合
	Indexes              map[string]*Index `protobuf:"bytes,5,rep,name=Indexes,proto3" json:"Indexes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Form) Reset()         { *m = Form{} }
func (m *Form) String() string { return proto.CompactTextString(m) }
func (*Form) ProtoMessage()    {}
func (*Form) Descriptor() ([]byte, []int) {
	return fileDescriptor_43e42cbf821258b1, []int{2}
}

func (m *Form) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Form.Unmarshal(m, b)
}
func (m *Form) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Form.Marshal(b, m, deterministic)
}
func (m *Form) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Form.Merge(m, src)
}
func (m *Form) XXX_Size() int {
	return xxx_messageInfo_Form.Size(m)
}
func (m *Form) XXX_DiscardUnknown() {
	xxx_messageInfo_Form.DiscardUnknown(m)
}

var xxx_messageInfo_Form proto.InternalMessageInfo

func (m *Form) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Form) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Form) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *Form) GetFormType() FormType {
	if m != nil {
		return m.FormType
	}
	return FormType_SQL
}

func (m *Form) GetIndexes() map[string]*Index {
	if m != nil {
		return m.Indexes
	}
	return nil
}

// Index 索引对象
type Index struct {
	// ID 索引唯一ID
	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	// Primary 是否主键
	Primary bool `protobuf:"varint,2,opt,name=Primary,proto3" json:"Primary,omitempty"`
	// KeyStructure 按照规范结构组成的索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
	KeyStructure         string   `protobuf:"bytes,3,opt,name=KeyStructure,proto3" json:"KeyStructure,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Index) Reset()         { *m = Index{} }
func (m *Index) String() string { return proto.CompactTextString(m) }
func (*Index) ProtoMessage()    {}
func (*Index) Descriptor() ([]byte, []int) {
	return fileDescriptor_43e42cbf821258b1, []int{3}
}

func (m *Index) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Index.Unmarshal(m, b)
}
func (m *Index) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Index.Marshal(b, m, deterministic)
}
func (m *Index) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Index.Merge(m, src)
}
func (m *Index) XXX_Size() int {
	return xxx_messageInfo_Index.Size(m)
}
func (m *Index) XXX_DiscardUnknown() {
	xxx_messageInfo_Index.DiscardUnknown(m)
}

var xxx_messageInfo_Index proto.InternalMessageInfo

func (m *Index) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Index) GetPrimary() bool {
	if m != nil {
		return m.Primary
	}
	return false
}

func (m *Index) GetKeyStructure() string {
	if m != nil {
		return m.KeyStructure
	}
	return ""
}

// Selector 检索选择器
type Selector struct {
	// Conditions 条件查询
	Conditions []*Condition `protobuf:"bytes,1,rep,name=Conditions,proto3" json:"Conditions,omitempty"`
	// Skip 结果集跳过数量
	Skip uint32 `protobuf:"varint,2,opt,name=Skip,proto3" json:"Skip,omitempty"`
	// Sort 排序方式
	Sort *Sort `protobuf:"bytes,3,opt,name=Sort,proto3" json:"Sort,omitempty"`
	// Limit 结果集顺序数量
	Limit                uint32   `protobuf:"varint,4,opt,name=Limit,proto3" json:"Limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Selector) Reset()         { *m = Selector{} }
func (m *Selector) String() string { return proto.CompactTextString(m) }
func (*Selector) ProtoMessage()    {}
func (*Selector) Descriptor() ([]byte, []int) {
	return fileDescriptor_43e42cbf821258b1, []int{4}
}

func (m *Selector) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Selector.Unmarshal(m, b)
}
func (m *Selector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Selector.Marshal(b, m, deterministic)
}
func (m *Selector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Selector.Merge(m, src)
}
func (m *Selector) XXX_Size() int {
	return xxx_messageInfo_Selector.Size(m)
}
func (m *Selector) XXX_DiscardUnknown() {
	xxx_messageInfo_Selector.DiscardUnknown(m)
}

var xxx_messageInfo_Selector proto.InternalMessageInfo

func (m *Selector) GetConditions() []*Condition {
	if m != nil {
		return m.Conditions
	}
	return nil
}

func (m *Selector) GetSkip() uint32 {
	if m != nil {
		return m.Skip
	}
	return 0
}

func (m *Selector) GetSort() *Sort {
	if m != nil {
		return m.Sort
	}
	return nil
}

func (m *Selector) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

// Condition 条件查询
type Condition struct {
	// Param 参数名，由对象结构层级字段通过'.'组成，如
	//
	// ref := &ref{
	//		i: 1,
	//		s: "2",
	//		in: refIn{
	//			i: 3,
	//			s: "4",
	//		},
	//	}
	//
	// key可取'i','in.s'
	Param string `protobuf:"bytes,1,opt,name=Param,proto3" json:"Param,omitempty"`
	// Cond 条件 gt/lt/eq/dif 大于/小于/等于/不等
	Cond string `protobuf:"bytes,2,opt,name=Cond,proto3" json:"Cond,omitempty"`
	// Value 比较对象，支持int、string、float和bool
	Value                []byte   `protobuf:"bytes,3,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Condition) Reset()         { *m = Condition{} }
func (m *Condition) String() string { return proto.CompactTextString(m) }
func (*Condition) ProtoMessage()    {}
func (*Condition) Descriptor() ([]byte, []int) {
	return fileDescriptor_43e42cbf821258b1, []int{5}
}

func (m *Condition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Condition.Unmarshal(m, b)
}
func (m *Condition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Condition.Marshal(b, m, deterministic)
}
func (m *Condition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Condition.Merge(m, src)
}
func (m *Condition) XXX_Size() int {
	return xxx_messageInfo_Condition.Size(m)
}
func (m *Condition) XXX_DiscardUnknown() {
	xxx_messageInfo_Condition.DiscardUnknown(m)
}

var xxx_messageInfo_Condition proto.InternalMessageInfo

func (m *Condition) GetParam() string {
	if m != nil {
		return m.Param
	}
	return ""
}

func (m *Condition) GetCond() string {
	if m != nil {
		return m.Cond
	}
	return ""
}

func (m *Condition) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

// Sort 排序方式
type Sort struct {
	// Param 参数名，由对象结构层级字段通过'.'组成，如
	//
	// ref := &ref{
	//		i: 1,
	//		s: "2",
	//		in: refIn{
	//			i: 3,
	//			s: "4",
	//		},
	//	}
	//
	// key可取'i','in.s'
	Param string `protobuf:"bytes,1,opt,name=Param,proto3" json:"Param,omitempty"`
	// ASC 是否升序
	ASC                  bool     `protobuf:"varint,2,opt,name=ASC,proto3" json:"ASC,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Sort) Reset()         { *m = Sort{} }
func (m *Sort) String() string { return proto.CompactTextString(m) }
func (*Sort) ProtoMessage()    {}
func (*Sort) Descriptor() ([]byte, []int) {
	return fileDescriptor_43e42cbf821258b1, []int{6}
}

func (m *Sort) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Sort.Unmarshal(m, b)
}
func (m *Sort) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Sort.Marshal(b, m, deterministic)
}
func (m *Sort) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Sort.Merge(m, src)
}
func (m *Sort) XXX_Size() int {
	return xxx_messageInfo_Sort.Size(m)
}
func (m *Sort) XXX_DiscardUnknown() {
	xxx_messageInfo_Sort.DiscardUnknown(m)
}

var xxx_messageInfo_Sort proto.InternalMessageInfo

func (m *Sort) GetParam() string {
	if m != nil {
		return m.Param
	}
	return ""
}

func (m *Sort) GetASC() bool {
	if m != nil {
		return m.ASC
	}
	return false
}

func init() {
	proto.RegisterEnum("api.FormType", FormType_name, FormType_value)
	proto.RegisterType((*Lily)(nil), "api.Lily")
	proto.RegisterMapType((map[string]*Database)(nil), "api.Lily.DatabasesEntry")
	proto.RegisterType((*Database)(nil), "api.Database")
	proto.RegisterMapType((map[string]*Form)(nil), "api.Database.FormsEntry")
	proto.RegisterType((*Form)(nil), "api.Form")
	proto.RegisterMapType((map[string]*Index)(nil), "api.Form.IndexesEntry")
	proto.RegisterType((*Index)(nil), "api.Index")
	proto.RegisterType((*Selector)(nil), "api.Selector")
	proto.RegisterType((*Condition)(nil), "api.Condition")
	proto.RegisterType((*Sort)(nil), "api.Sort")
}

func init() { proto.RegisterFile("connector/grpc/data.proto", fileDescriptor_43e42cbf821258b1) }

var fileDescriptor_43e42cbf821258b1 = []byte{
	// 507 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0x51, 0xab, 0xd3, 0x30,
	0x18, 0x35, 0x6b, 0xe7, 0xd6, 0xef, 0x6e, 0x63, 0x04, 0x91, 0x38, 0x94, 0x3b, 0xea, 0xcb, 0x14,
	0xc9, 0x95, 0x2b, 0x88, 0xf8, 0xe6, 0xdd, 0xbc, 0x30, 0x36, 0x64, 0xa6, 0xea, 0x7b, 0xd6, 0x05,
	0x09, 0xb7, 0x6d, 0x4a, 0xd6, 0x89, 0x7d, 0xf6, 0xd1, 0xdf, 0xe4, 0x0f, 0xf2, 0x5f, 0x48, 0xd2,
	0xa6, 0xeb, 0x70, 0x6f, 0xf7, 0xed, 0xfb, 0xce, 0x39, 0x39, 0xc9, 0xf9, 0x92, 0xc0, 0x93, 0x58,
	0x65, 0x99, 0x88, 0x0b, 0xa5, 0xaf, 0xbe, 0xeb, 0x3c, 0xbe, 0xda, 0xf1, 0x82, 0xd3, 0x5c, 0xab,
	0x42, 0x61, 0x8f, 0xe7, 0x32, 0xfc, 0x8d, 0xc0, 0x5f, 0xcb, 0xa4, 0xc4, 0x6f, 0x21, 0x30, 0xdc,
	0x96, 0xef, 0xc5, 0x9e, 0xa0, 0xa9, 0x37, 0xbb, 0xb8, 0x26, 0x94, 0xe7, 0x92, 0x1a, 0x96, 0x2e,
	0x1c, 0xf5, 0x31, 0x2b, 0x74, 0xc9, 0x8e, 0xd2, 0xc9, 0x0a, 0x46, 0xa7, 0x24, 0x1e, 0x83, 0x77,
	0x27, 0x4a, 0x82, 0xa6, 0x68, 0x16, 0x30, 0x53, 0xe2, 0xe7, 0xd0, 0xfd, 0xc1, 0x93, 0x83, 0x20,
	0x9d, 0x29, 0x9a, 0x5d, 0x5c, 0x0f, 0xad, 0xaf, 0x5b, 0xc5, 0x2a, 0xee, 0x7d, 0xe7, 0x1d, 0x0a,
	0xff, 0x20, 0xe8, 0x3b, 0x1c, 0x8f, 0xa0, 0xb3, 0x5c, 0xd4, 0x36, 0x9d, 0xe5, 0x02, 0x63, 0xf0,
	0x3f, 0xf1, 0xb4, 0x32, 0x09, 0x98, 0xad, 0x31, 0x81, 0xde, 0x5c, 0xa5, 0xa9, 0xc8, 0x0a, 0xe2,
	0x59, 0xd8, 0xb5, 0x98, 0x42, 0xf7, 0x56, 0xe9, 0x74, 0x4f, 0xfc, 0x56, 0x16, 0xe7, 0x4d, 0x2d,
	0x55, 0x65, 0xa9, 0x64, 0x93, 0x39, 0xc0, 0x11, 0x3c, 0x93, 0xe1, 0xf2, 0x34, 0x43, 0x60, 0xfd,
	0xcc, 0x8a, 0xf6, 0xf9, 0xff, 0x22, 0xf0, 0x0d, 0x76, 0xcf, 0xb3, 0xbf, 0x80, 0xbe, 0x71, 0xf9,
	0x52, 0xe6, 0x82, 0xf8, 0x53, 0x34, 0x1b, 0xd5, 0x23, 0x73, 0x20, 0x6b, 0x68, 0xfc, 0x1a, 0x7a,
	0xcb, 0x6c, 0x27, 0x7e, 0x8a, 0x3d, 0xe9, 0xda, 0xa0, 0x8f, 0x1b, 0x25, 0xad, 0x89, 0x2a, 0xa6,
	0x93, 0x4d, 0x6e, 0x61, 0xd0, 0x26, 0xce, 0x44, 0x9d, 0x9e, 0x46, 0x05, 0xeb, 0x68, 0xd7, 0xb4,
	0xb3, 0x7e, 0x85, 0xae, 0xc5, 0xfe, 0xcb, 0x4a, 0xa0, 0xb7, 0xd1, 0x32, 0xe5, 0xba, 0xb4, 0x06,
	0x7d, 0xe6, 0x5a, 0x1c, 0xc2, 0x60, 0x25, 0xca, 0xa8, 0xd0, 0x87, 0xb8, 0x38, 0x68, 0x51, 0xc7,
	0x3e, 0xc1, 0xc2, 0x5f, 0x08, 0xfa, 0x91, 0x48, 0xec, 0x93, 0xc5, 0x14, 0x60, 0xae, 0xb2, 0x9d,
	0x2c, 0xa4, 0xca, 0xdc, 0xab, 0x1c, 0xd9, 0xe3, 0x34, 0x30, 0x6b, 0x29, 0xcc, 0x98, 0xa3, 0x3b,
	0x99, 0xdb, 0x7d, 0x87, 0xcc, 0xd6, 0xf8, 0x19, 0xf8, 0x91, 0xd2, 0xd5, 0x8c, 0xdd, 0xbd, 0x19,
	0x80, 0x59, 0x18, 0x3f, 0x82, 0xee, 0x5a, 0xa6, 0xb2, 0xb0, 0x83, 0x1e, 0xb2, 0xaa, 0x09, 0x57,
	0x10, 0x34, 0xb6, 0x46, 0xb2, 0xe1, 0x9a, 0xa7, 0x75, 0xc6, 0xaa, 0x31, 0x7b, 0x19, 0x89, 0xbb,
	0x52, 0x53, 0x1b, 0xe5, 0x37, 0x3b, 0x39, 0xb3, 0xd9, 0x80, 0x55, 0x4d, 0x48, 0xa1, 0xd9, 0xea,
	0x8c, 0xcf, 0x18, 0xbc, 0x0f, 0xd1, 0xbc, 0x1e, 0x95, 0x29, 0x5f, 0x3e, 0x3d, 0x5e, 0x3f, 0xee,
	0x81, 0x17, 0x7d, 0x5e, 0x8f, 0x1f, 0x98, 0x62, 0xa1, 0xe2, 0x31, 0xba, 0x79, 0x05, 0x97, 0x71,
	0x46, 0xf9, 0x56, 0x68, 0x19, 0xd3, 0x44, 0x26, 0xe5, 0x6e, 0x4b, 0x9b, 0x4f, 0x4e, 0xcd, 0x27,
	0xbf, 0x09, 0xcc, 0x3b, 0xdf, 0x98, 0x4f, 0xbe, 0x7d, 0x68, 0xff, 0xfa, 0x9b, 0x7f, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x3e, 0xc5, 0xe0, 0x1e, 0x08, 0x04, 0x00, 0x00,
}
