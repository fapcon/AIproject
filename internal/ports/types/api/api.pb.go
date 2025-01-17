// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.0
// source: api.proto

package api

import (
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

type Exchange int32

const (
	Exchange_ExchangeUnknown Exchange = 0
	Exchange_ExchangeOKX     Exchange = 1
)

// Enum value maps for Exchange.
var (
	Exchange_name = map[int32]string{
		0: "ExchangeUnknown",
		1: "ExchangeOKX",
	}
	Exchange_value = map[string]int32{
		"ExchangeUnknown": 0,
		"ExchangeOKX":     1,
	}
)

func (x Exchange) Enum() *Exchange {
	p := new(Exchange)
	*p = x
	return p
}

func (x Exchange) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Exchange) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_enumTypes[0].Descriptor()
}

func (Exchange) Type() protoreflect.EnumType {
	return &file_api_proto_enumTypes[0]
}

func (x Exchange) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Exchange.Descriptor instead.
func (Exchange) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{0}
}

type Instruments struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instruments []*Instrument `protobuf:"bytes,1,rep,name=instruments,proto3" json:"instruments,omitempty"`
}

func (x *Instruments) Reset() {
	*x = Instruments{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Instruments) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Instruments) ProtoMessage() {}

func (x *Instruments) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Instruments.ProtoReflect.Descriptor instead.
func (*Instruments) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{0}
}

func (x *Instruments) GetInstruments() []*Instrument {
	if x != nil {
		return x.Instruments
	}
	return nil
}

type InstrumentDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InstrumentDetails []*InstrumentDetail `protobuf:"bytes,1,rep,name=instrument_details,json=instrumentDetails,proto3" json:"instrument_details,omitempty"`
}

func (x *InstrumentDetails) Reset() {
	*x = InstrumentDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstrumentDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstrumentDetails) ProtoMessage() {}

func (x *InstrumentDetails) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstrumentDetails.ProtoReflect.Descriptor instead.
func (*InstrumentDetails) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{1}
}

func (x *InstrumentDetails) GetInstrumentDetails() []*InstrumentDetail {
	if x != nil {
		return x.InstrumentDetails
	}
	return nil
}

type Instrument struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exchange       Exchange `protobuf:"varint,1,opt,name=Exchange,proto3,enum=api.Exchange" json:"Exchange,omitempty"`
	LocalSymbol    string   `protobuf:"bytes,2,opt,name=LocalSymbol,proto3" json:"LocalSymbol,omitempty"`
	ExchangeSymbol string   `protobuf:"bytes,3,opt,name=ExchangeSymbol,proto3" json:"ExchangeSymbol,omitempty"`
}

func (x *Instrument) Reset() {
	*x = Instrument{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Instrument) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Instrument) ProtoMessage() {}

func (x *Instrument) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Instrument.ProtoReflect.Descriptor instead.
func (*Instrument) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{2}
}

func (x *Instrument) GetExchange() Exchange {
	if x != nil {
		return x.Exchange
	}
	return Exchange_ExchangeUnknown
}

func (x *Instrument) GetLocalSymbol() string {
	if x != nil {
		return x.LocalSymbol
	}
	return ""
}

func (x *Instrument) GetExchangeSymbol() string {
	if x != nil {
		return x.ExchangeSymbol
	}
	return ""
}

type InstrumentDetail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exchange       Exchange `protobuf:"varint,1,opt,name=Exchange,proto3,enum=api.Exchange" json:"Exchange,omitempty"`
	LocalSymbol    string   `protobuf:"bytes,2,opt,name=LocalSymbol,proto3" json:"LocalSymbol,omitempty"`
	ExchangeSymbol string   `protobuf:"bytes,3,opt,name=ExchangeSymbol,proto3" json:"ExchangeSymbol,omitempty"`
	PricePrecision int32    `protobuf:"varint,4,opt,name=PricePrecision,proto3" json:"PricePrecision,omitempty"`
	SizePrecision  int32    `protobuf:"varint,5,opt,name=SizePrecision,proto3" json:"SizePrecision,omitempty"`
	MinLot         float32  `protobuf:"fixed32,6,opt,name=MinLot,proto3" json:"MinLot,omitempty"`
}

func (x *InstrumentDetail) Reset() {
	*x = InstrumentDetail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstrumentDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstrumentDetail) ProtoMessage() {}

func (x *InstrumentDetail) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstrumentDetail.ProtoReflect.Descriptor instead.
func (*InstrumentDetail) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{3}
}

func (x *InstrumentDetail) GetExchange() Exchange {
	if x != nil {
		return x.Exchange
	}
	return Exchange_ExchangeUnknown
}

func (x *InstrumentDetail) GetLocalSymbol() string {
	if x != nil {
		return x.LocalSymbol
	}
	return ""
}

func (x *InstrumentDetail) GetExchangeSymbol() string {
	if x != nil {
		return x.ExchangeSymbol
	}
	return ""
}

func (x *InstrumentDetail) GetPricePrecision() int32 {
	if x != nil {
		return x.PricePrecision
	}
	return 0
}

func (x *InstrumentDetail) GetSizePrecision() int32 {
	if x != nil {
		return x.SizePrecision
	}
	return 0
}

func (x *InstrumentDetail) GetMinLot() float32 {
	if x != nil {
		return x.MinLot
	}
	return 0
}

var File_api_proto protoreflect.FileDescriptor

var file_api_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x40, 0x0a,
	0x0b, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x31, 0x0a, 0x0b,
	0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x0b, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22,
	0x59, 0x0a, 0x11, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x12, 0x44, 0x0a, 0x12, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x11, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d,
	0x65, 0x6e, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x81, 0x01, 0x0a, 0x0a, 0x49,
	0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x29, 0x0a, 0x08, 0x45, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x08, 0x45, 0x78, 0x63, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x53, 0x79, 0x6d,
	0x62, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4c, 0x6f, 0x63, 0x61, 0x6c,
	0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x26, 0x0a, 0x0e, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x22, 0xed,
	0x01, 0x0a, 0x10, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x12, 0x29, 0x0a, 0x08, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x63, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x52, 0x08, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c,
	0x12, 0x26, 0x0a, 0x0e, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x53, 0x79, 0x6d, 0x62,
	0x6f, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x26, 0x0a, 0x0e, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x50, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x50, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x24, 0x0a, 0x0d, 0x53, 0x69, 0x7a, 0x65, 0x50, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x53, 0x69, 0x7a, 0x65, 0x50, 0x72, 0x65,
	0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x69, 0x6e, 0x4c, 0x6f, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x4d, 0x69, 0x6e, 0x4c, 0x6f, 0x74, 0x2a, 0x30,
	0x0a, 0x08, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x13, 0x0a, 0x0f, 0x45, 0x78,
	0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12,
	0x0f, 0x0a, 0x0b, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4f, 0x4b, 0x58, 0x10, 0x01,
	0x32, 0x9e, 0x02, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x40, 0x0a, 0x14, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x49, 0x6e, 0x73, 0x74,
	0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x10, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e,
	0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x12, 0x42, 0x0a, 0x16, 0x55, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x10, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3f, 0x0a, 0x13, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x4e,
	0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x73, 0x12, 0x10, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x49, 0x0a, 0x17, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x79,
	0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x73, 0x12, 0x16, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_api_proto_rawDescOnce sync.Once
	file_api_proto_rawDescData = file_api_proto_rawDesc
)

func file_api_proto_rawDescGZIP() []byte {
	file_api_proto_rawDescOnce.Do(func() {
		file_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_rawDescData)
	})
	return file_api_proto_rawDescData
}

var file_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_proto_goTypes = []interface{}{
	(Exchange)(0),             // 0: api.Exchange
	(*Instruments)(nil),       // 1: api.Instruments
	(*InstrumentDetails)(nil), // 2: api.InstrumentDetails
	(*Instrument)(nil),        // 3: api.Instrument
	(*InstrumentDetail)(nil),  // 4: api.InstrumentDetail
	(*emptypb.Empty)(nil),     // 5: google.protobuf.Empty
}
var file_api_proto_depIdxs = []int32{
	3, // 0: api.Instruments.instruments:type_name -> api.Instrument
	4, // 1: api.InstrumentDetails.instrument_details:type_name -> api.InstrumentDetail
	0, // 2: api.Instrument.Exchange:type_name -> api.Exchange
	0, // 3: api.InstrumentDetail.Exchange:type_name -> api.Exchange
	1, // 4: api.APIService.SubscribeInstruments:input_type -> api.Instruments
	1, // 5: api.APIService.UnsubscribeInstruments:input_type -> api.Instruments
	1, // 6: api.APIService.ModifyNameResolvers:input_type -> api.Instruments
	2, // 7: api.APIService.ModifyInstrumentDetails:input_type -> api.InstrumentDetails
	5, // 8: api.APIService.SubscribeInstruments:output_type -> google.protobuf.Empty
	5, // 9: api.APIService.UnsubscribeInstruments:output_type -> google.protobuf.Empty
	5, // 10: api.APIService.ModifyNameResolvers:output_type -> google.protobuf.Empty
	5, // 11: api.APIService.ModifyInstrumentDetails:output_type -> google.protobuf.Empty
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_proto_init() }
func file_api_proto_init() {
	if File_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Instruments); i {
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
		file_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstrumentDetails); i {
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
		file_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Instrument); i {
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
		file_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstrumentDetail); i {
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
			RawDescriptor: file_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_goTypes,
		DependencyIndexes: file_api_proto_depIdxs,
		EnumInfos:         file_api_proto_enumTypes,
		MessageInfos:      file_api_proto_msgTypes,
	}.Build()
	File_api_proto = out.File
	file_api_proto_rawDesc = nil
	file_api_proto_goTypes = nil
	file_api_proto_depIdxs = nil
}
