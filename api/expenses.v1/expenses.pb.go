// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: expenses.v1/expenses.proto

package expensesv1

import (
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

type CreateExpenseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount    uint64                 `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
	Date      *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	ReceiptId *uint64                `protobuf:"varint,3,opt,name=receipt_id,json=receiptId,proto3,oneof" json:"receipt_id,omitempty"`
}

func (x *CreateExpenseRequest) Reset() {
	*x = CreateExpenseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_expenses_v1_expenses_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateExpenseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateExpenseRequest) ProtoMessage() {}

func (x *CreateExpenseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_expenses_v1_expenses_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateExpenseRequest.ProtoReflect.Descriptor instead.
func (*CreateExpenseRequest) Descriptor() ([]byte, []int) {
	return file_expenses_v1_expenses_proto_rawDescGZIP(), []int{0}
}

func (x *CreateExpenseRequest) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *CreateExpenseRequest) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *CreateExpenseRequest) GetReceiptId() uint64 {
	if x != nil && x.ReceiptId != nil {
		return *x.ReceiptId
	}
	return 0
}

type CreateExpenseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Expense *Expense `protobuf:"bytes,1,opt,name=expense,proto3" json:"expense,omitempty"`
}

func (x *CreateExpenseResponse) Reset() {
	*x = CreateExpenseResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_expenses_v1_expenses_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateExpenseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateExpenseResponse) ProtoMessage() {}

func (x *CreateExpenseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_expenses_v1_expenses_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateExpenseResponse.ProtoReflect.Descriptor instead.
func (*CreateExpenseResponse) Descriptor() ([]byte, []int) {
	return file_expenses_v1_expenses_proto_rawDescGZIP(), []int{1}
}

func (x *CreateExpenseResponse) GetExpense() *Expense {
	if x != nil {
		return x.Expense
	}
	return nil
}

type UpdateExpenseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ReceiptId   *uint64                `protobuf:"varint,2,opt,name=receipt_id,json=receiptId,proto3,oneof" json:"receipt_id,omitempty"`
	Amount      *uint64                `protobuf:"varint,3,opt,name=amount,proto3,oneof" json:"amount,omitempty"`
	Date        *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=date,proto3,oneof" json:"date,omitempty"`
	Category    *string                `protobuf:"bytes,5,opt,name=category,proto3,oneof" json:"category,omitempty"`
	Subcategory *string                `protobuf:"bytes,6,opt,name=subcategory,proto3,oneof" json:"subcategory,omitempty"`
	Description *uint64                `protobuf:"varint,7,opt,name=description,proto3,oneof" json:"description,omitempty"`
}

func (x *UpdateExpenseRequest) Reset() {
	*x = UpdateExpenseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_expenses_v1_expenses_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateExpenseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateExpenseRequest) ProtoMessage() {}

func (x *UpdateExpenseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_expenses_v1_expenses_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateExpenseRequest.ProtoReflect.Descriptor instead.
func (*UpdateExpenseRequest) Descriptor() ([]byte, []int) {
	return file_expenses_v1_expenses_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateExpenseRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateExpenseRequest) GetReceiptId() uint64 {
	if x != nil && x.ReceiptId != nil {
		return *x.ReceiptId
	}
	return 0
}

func (x *UpdateExpenseRequest) GetAmount() uint64 {
	if x != nil && x.Amount != nil {
		return *x.Amount
	}
	return 0
}

func (x *UpdateExpenseRequest) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *UpdateExpenseRequest) GetCategory() string {
	if x != nil && x.Category != nil {
		return *x.Category
	}
	return ""
}

func (x *UpdateExpenseRequest) GetSubcategory() string {
	if x != nil && x.Subcategory != nil {
		return *x.Subcategory
	}
	return ""
}

func (x *UpdateExpenseRequest) GetDescription() uint64 {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return 0
}

type UpdateExpenseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Expense *Expense `protobuf:"bytes,1,opt,name=expense,proto3" json:"expense,omitempty"`
}

func (x *UpdateExpenseResponse) Reset() {
	*x = UpdateExpenseResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_expenses_v1_expenses_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateExpenseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateExpenseResponse) ProtoMessage() {}

func (x *UpdateExpenseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_expenses_v1_expenses_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateExpenseResponse.ProtoReflect.Descriptor instead.
func (*UpdateExpenseResponse) Descriptor() ([]byte, []int) {
	return file_expenses_v1_expenses_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateExpenseResponse) GetExpense() *Expense {
	if x != nil {
		return x.Expense
	}
	return nil
}

type DeleteExpenseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteExpenseRequest) Reset() {
	*x = DeleteExpenseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_expenses_v1_expenses_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteExpenseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteExpenseRequest) ProtoMessage() {}

func (x *DeleteExpenseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_expenses_v1_expenses_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteExpenseRequest.ProtoReflect.Descriptor instead.
func (*DeleteExpenseRequest) Descriptor() ([]byte, []int) {
	return file_expenses_v1_expenses_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteExpenseRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DeleteExpenseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteExpenseResponse) Reset() {
	*x = DeleteExpenseResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_expenses_v1_expenses_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteExpenseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteExpenseResponse) ProtoMessage() {}

func (x *DeleteExpenseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_expenses_v1_expenses_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteExpenseResponse.ProtoReflect.Descriptor instead.
func (*DeleteExpenseResponse) Descriptor() ([]byte, []int) {
	return file_expenses_v1_expenses_proto_rawDescGZIP(), []int{5}
}

type ListExpensesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserEmail *string `protobuf:"bytes,1,opt,name=user_email,json=userEmail,proto3,oneof" json:"user_email,omitempty"`
	ReceiptId *string `protobuf:"bytes,2,opt,name=receipt_id,json=receiptId,proto3,oneof" json:"receipt_id,omitempty"`
}

func (x *ListExpensesRequest) Reset() {
	*x = ListExpensesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_expenses_v1_expenses_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListExpensesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListExpensesRequest) ProtoMessage() {}

func (x *ListExpensesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_expenses_v1_expenses_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListExpensesRequest.ProtoReflect.Descriptor instead.
func (*ListExpensesRequest) Descriptor() ([]byte, []int) {
	return file_expenses_v1_expenses_proto_rawDescGZIP(), []int{6}
}

func (x *ListExpensesRequest) GetUserEmail() string {
	if x != nil && x.UserEmail != nil {
		return *x.UserEmail
	}
	return ""
}

func (x *ListExpensesRequest) GetReceiptId() string {
	if x != nil && x.ReceiptId != nil {
		return *x.ReceiptId
	}
	return ""
}

type ListExpensesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Expenses []*Expense `protobuf:"bytes,1,rep,name=expenses,proto3" json:"expenses,omitempty"`
}

func (x *ListExpensesResponse) Reset() {
	*x = ListExpensesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_expenses_v1_expenses_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListExpensesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListExpensesResponse) ProtoMessage() {}

func (x *ListExpensesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_expenses_v1_expenses_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListExpensesResponse.ProtoReflect.Descriptor instead.
func (*ListExpensesResponse) Descriptor() ([]byte, []int) {
	return file_expenses_v1_expenses_proto_rawDescGZIP(), []int{7}
}

func (x *ListExpensesResponse) GetExpenses() []*Expense {
	if x != nil {
		return x.Expenses
	}
	return nil
}

type Expense struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ReceiptId   *uint64                `protobuf:"varint,2,opt,name=receipt_id,json=receiptId,proto3,oneof" json:"receipt_id,omitempty"`
	Amount      uint64                 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Date        *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=date,proto3" json:"date,omitempty"`
	Category    string                 `protobuf:"bytes,5,opt,name=category,proto3" json:"category,omitempty"`
	Subcategory string                 `protobuf:"bytes,6,opt,name=subcategory,proto3" json:"subcategory,omitempty"`
	Description uint64                 `protobuf:"varint,7,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *Expense) Reset() {
	*x = Expense{}
	if protoimpl.UnsafeEnabled {
		mi := &file_expenses_v1_expenses_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Expense) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Expense) ProtoMessage() {}

func (x *Expense) ProtoReflect() protoreflect.Message {
	mi := &file_expenses_v1_expenses_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Expense.ProtoReflect.Descriptor instead.
func (*Expense) Descriptor() ([]byte, []int) {
	return file_expenses_v1_expenses_proto_rawDescGZIP(), []int{8}
}

func (x *Expense) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Expense) GetReceiptId() uint64 {
	if x != nil && x.ReceiptId != nil {
		return *x.ReceiptId
	}
	return 0
}

func (x *Expense) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *Expense) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *Expense) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Expense) GetSubcategory() string {
	if x != nil {
		return x.Subcategory
	}
	return ""
}

func (x *Expense) GetDescription() uint64 {
	if x != nil {
		return x.Description
	}
	return 0
}

var File_expenses_v1_expenses_proto protoreflect.FileDescriptor

var file_expenses_v1_expenses_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2f, 0x65, 0x78,
	0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x65, 0x78,
	0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x91, 0x01, 0x0a, 0x14, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x22, 0x0a, 0x0a, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x48,
	0x00, 0x52, 0x09, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42,
	0x0d, 0x0a, 0x0b, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x5f, 0x69, 0x64, 0x22, 0x47,
	0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x65, 0x78, 0x70, 0x65, 0x6e,
	0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x6e,
	0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x07,
	0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x22, 0xdb, 0x02, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x22, 0x0a, 0x0a, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x09, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x49,
	0x64, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x48, 0x01, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x88, 0x01,
	0x01, 0x12, 0x33, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x02, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x79, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x0b, 0x73, 0x75, 0x62, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x0b,
	0x73, 0x75, 0x62, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x88, 0x01, 0x01, 0x12, 0x25,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x04, 0x48, 0x05, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70,
	0x74, 0x5f, 0x69, 0x64, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x42,
	0x07, 0x0a, 0x05, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x63, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x73, 0x75, 0x62, 0x63, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x47, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45,
	0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e,
	0x0a, 0x07, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78,
	0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x07, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x22, 0x26,
	0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0x17, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x7b, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x09, 0x75, 0x73,
	0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0a, 0x72, 0x65,
	0x63, 0x65, 0x69, 0x70, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01,
	0x52, 0x09, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0d,
	0x0a, 0x0b, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x0d, 0x0a,
	0x0b, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x5f, 0x69, 0x64, 0x22, 0x48, 0x0a, 0x14,
	0x4c, 0x69, 0x73, 0x74, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x08, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x65, 0x78,
	0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x22, 0xf4, 0x01, 0x0a, 0x07, 0x45, 0x78, 0x70, 0x65, 0x6e,
	0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x22, 0x0a, 0x0a, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x09, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70,
	0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2e,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x75,
	0x62, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x73, 0x75, 0x62, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0d,
	0x0a, 0x0b, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x5f, 0x69, 0x64, 0x32, 0xf6, 0x02,
	0x0a, 0x0f, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x58, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x78, 0x70, 0x65, 0x6e,
	0x73, 0x65, 0x12, 0x21, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x58, 0x0a, 0x0d, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x2e, 0x65,
	0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x22, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x58, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45,
	0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x78, 0x70, 0x65, 0x6e,
	0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x65, 0x78, 0x70, 0x65,
	0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x78,
	0x70, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x55, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x12,
	0x20, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x21, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0xa5, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x65,
	0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x45, 0x78, 0x70, 0x65,
	0x6e, 0x73, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x36, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x6e, 0x7a, 0x61, 0x6e, 0x69, 0x74,
	0x30, 0x2f, 0x6d, 0x63, 0x64, 0x75, 0x63, 0x6b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x78, 0x70,
	0x65, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x3b, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65,
	0x73, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x45, 0x58, 0x58, 0xaa, 0x02, 0x0b, 0x45, 0x78, 0x70, 0x65,
	0x6e, 0x73, 0x65, 0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73,
	0x65, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x17, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73,
	0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x0c, 0x45, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_expenses_v1_expenses_proto_rawDescOnce sync.Once
	file_expenses_v1_expenses_proto_rawDescData = file_expenses_v1_expenses_proto_rawDesc
)

func file_expenses_v1_expenses_proto_rawDescGZIP() []byte {
	file_expenses_v1_expenses_proto_rawDescOnce.Do(func() {
		file_expenses_v1_expenses_proto_rawDescData = protoimpl.X.CompressGZIP(file_expenses_v1_expenses_proto_rawDescData)
	})
	return file_expenses_v1_expenses_proto_rawDescData
}

var file_expenses_v1_expenses_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_expenses_v1_expenses_proto_goTypes = []any{
	(*CreateExpenseRequest)(nil),  // 0: expenses.v1.CreateExpenseRequest
	(*CreateExpenseResponse)(nil), // 1: expenses.v1.CreateExpenseResponse
	(*UpdateExpenseRequest)(nil),  // 2: expenses.v1.UpdateExpenseRequest
	(*UpdateExpenseResponse)(nil), // 3: expenses.v1.UpdateExpenseResponse
	(*DeleteExpenseRequest)(nil),  // 4: expenses.v1.DeleteExpenseRequest
	(*DeleteExpenseResponse)(nil), // 5: expenses.v1.DeleteExpenseResponse
	(*ListExpensesRequest)(nil),   // 6: expenses.v1.ListExpensesRequest
	(*ListExpensesResponse)(nil),  // 7: expenses.v1.ListExpensesResponse
	(*Expense)(nil),               // 8: expenses.v1.Expense
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
}
var file_expenses_v1_expenses_proto_depIdxs = []int32{
	9,  // 0: expenses.v1.CreateExpenseRequest.date:type_name -> google.protobuf.Timestamp
	8,  // 1: expenses.v1.CreateExpenseResponse.expense:type_name -> expenses.v1.Expense
	9,  // 2: expenses.v1.UpdateExpenseRequest.date:type_name -> google.protobuf.Timestamp
	8,  // 3: expenses.v1.UpdateExpenseResponse.expense:type_name -> expenses.v1.Expense
	8,  // 4: expenses.v1.ListExpensesResponse.expenses:type_name -> expenses.v1.Expense
	9,  // 5: expenses.v1.Expense.date:type_name -> google.protobuf.Timestamp
	0,  // 6: expenses.v1.ExpensesService.CreateExpense:input_type -> expenses.v1.CreateExpenseRequest
	2,  // 7: expenses.v1.ExpensesService.UpdateExpense:input_type -> expenses.v1.UpdateExpenseRequest
	4,  // 8: expenses.v1.ExpensesService.DeleteExpense:input_type -> expenses.v1.DeleteExpenseRequest
	6,  // 9: expenses.v1.ExpensesService.ListExpenses:input_type -> expenses.v1.ListExpensesRequest
	1,  // 10: expenses.v1.ExpensesService.CreateExpense:output_type -> expenses.v1.CreateExpenseResponse
	3,  // 11: expenses.v1.ExpensesService.UpdateExpense:output_type -> expenses.v1.UpdateExpenseResponse
	5,  // 12: expenses.v1.ExpensesService.DeleteExpense:output_type -> expenses.v1.DeleteExpenseResponse
	7,  // 13: expenses.v1.ExpensesService.ListExpenses:output_type -> expenses.v1.ListExpensesResponse
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_expenses_v1_expenses_proto_init() }
func file_expenses_v1_expenses_proto_init() {
	if File_expenses_v1_expenses_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_expenses_v1_expenses_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateExpenseRequest); i {
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
		file_expenses_v1_expenses_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateExpenseResponse); i {
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
		file_expenses_v1_expenses_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateExpenseRequest); i {
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
		file_expenses_v1_expenses_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateExpenseResponse); i {
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
		file_expenses_v1_expenses_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteExpenseRequest); i {
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
		file_expenses_v1_expenses_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteExpenseResponse); i {
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
		file_expenses_v1_expenses_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ListExpensesRequest); i {
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
		file_expenses_v1_expenses_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*ListExpensesResponse); i {
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
		file_expenses_v1_expenses_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*Expense); i {
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
	file_expenses_v1_expenses_proto_msgTypes[0].OneofWrappers = []any{}
	file_expenses_v1_expenses_proto_msgTypes[2].OneofWrappers = []any{}
	file_expenses_v1_expenses_proto_msgTypes[6].OneofWrappers = []any{}
	file_expenses_v1_expenses_proto_msgTypes[8].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_expenses_v1_expenses_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_expenses_v1_expenses_proto_goTypes,
		DependencyIndexes: file_expenses_v1_expenses_proto_depIdxs,
		MessageInfos:      file_expenses_v1_expenses_proto_msgTypes,
	}.Build()
	File_expenses_v1_expenses_proto = out.File
	file_expenses_v1_expenses_proto_rawDesc = nil
	file_expenses_v1_expenses_proto_goTypes = nil
	file_expenses_v1_expenses_proto_depIdxs = nil
}