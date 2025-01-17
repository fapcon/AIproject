// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.0
// source: internal_events.proto

package intev

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	decimalpb "studentgit.kata.academy/quant/torque/pkg/decimalpb"
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
	Exchange_ExchangeGateIO  Exchange = 2
	Exchange_ExchangeBybit   Exchange = 3
)

// Enum value maps for Exchange.
var (
	Exchange_name = map[int32]string{
		0: "ExchangeUnknown",
		1: "ExchangeOKX",
		2: "ExchangeGateIO",
		3: "ExchangeBybit",
	}
	Exchange_value = map[string]int32{
		"ExchangeUnknown": 0,
		"ExchangeOKX":     1,
		"ExchangeGateIO":  2,
		"ExchangeBybit":   3,
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
	return file_internal_events_proto_enumTypes[0].Descriptor()
}

func (Exchange) Type() protoreflect.EnumType {
	return &file_internal_events_proto_enumTypes[0]
}

func (x Exchange) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Exchange.Descriptor instead.
func (Exchange) EnumDescriptor() ([]byte, []int) {
	return file_internal_events_proto_rawDescGZIP(), []int{0}
}

type OrderType int32

const (
	OrderType_ORDER_TYPE_UNSPECIFIED OrderType = 0
	OrderType_ORDER_TYPE_MARKET      OrderType = 1
	OrderType_ORDER_TYPE_LIMIT       OrderType = 2
)

// Enum value maps for OrderType.
var (
	OrderType_name = map[int32]string{
		0: "ORDER_TYPE_UNSPECIFIED",
		1: "ORDER_TYPE_MARKET",
		2: "ORDER_TYPE_LIMIT",
	}
	OrderType_value = map[string]int32{
		"ORDER_TYPE_UNSPECIFIED": 0,
		"ORDER_TYPE_MARKET":      1,
		"ORDER_TYPE_LIMIT":       2,
	}
)

func (x OrderType) Enum() *OrderType {
	p := new(OrderType)
	*p = x
	return p
}

func (x OrderType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderType) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_events_proto_enumTypes[1].Descriptor()
}

func (OrderType) Type() protoreflect.EnumType {
	return &file_internal_events_proto_enumTypes[1]
}

func (x OrderType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderType.Descriptor instead.
func (OrderType) EnumDescriptor() ([]byte, []int) {
	return file_internal_events_proto_rawDescGZIP(), []int{1}
}

type OrderSide int32

const (
	OrderSide_ORDER_SIDE_UNSPECIFIED OrderSide = 0
	OrderSide_ORDER_SIDE_BUY         OrderSide = 1
	OrderSide_ORDER_SIDE_SELL        OrderSide = 2
)

// Enum value maps for OrderSide.
var (
	OrderSide_name = map[int32]string{
		0: "ORDER_SIDE_UNSPECIFIED",
		1: "ORDER_SIDE_BUY",
		2: "ORDER_SIDE_SELL",
	}
	OrderSide_value = map[string]int32{
		"ORDER_SIDE_UNSPECIFIED": 0,
		"ORDER_SIDE_BUY":         1,
		"ORDER_SIDE_SELL":        2,
	}
)

func (x OrderSide) Enum() *OrderSide {
	p := new(OrderSide)
	*p = x
	return p
}

func (x OrderSide) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderSide) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_events_proto_enumTypes[2].Descriptor()
}

func (OrderSide) Type() protoreflect.EnumType {
	return &file_internal_events_proto_enumTypes[2]
}

func (x OrderSide) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderSide.Descriptor instead.
func (OrderSide) EnumDescriptor() ([]byte, []int) {
	return file_internal_events_proto_rawDescGZIP(), []int{2}
}

type TimeInForce int32

const (
	TimeInForce_TIME_IN_FORCE_UNSPECIFIED TimeInForce = 0
	TimeInForce_TIME_IN_FORCE_GTC         TimeInForce = 1
	TimeInForce_TIME_IN_FORCE_IOC         TimeInForce = 2
	TimeInForce_TIME_IN_FORCE_FOK         TimeInForce = 3
)

// Enum value maps for TimeInForce.
var (
	TimeInForce_name = map[int32]string{
		0: "TIME_IN_FORCE_UNSPECIFIED",
		1: "TIME_IN_FORCE_GTC",
		2: "TIME_IN_FORCE_IOC",
		3: "TIME_IN_FORCE_FOK",
	}
	TimeInForce_value = map[string]int32{
		"TIME_IN_FORCE_UNSPECIFIED": 0,
		"TIME_IN_FORCE_GTC":         1,
		"TIME_IN_FORCE_IOC":         2,
		"TIME_IN_FORCE_FOK":         3,
	}
)

func (x TimeInForce) Enum() *TimeInForce {
	p := new(TimeInForce)
	*p = x
	return p
}

func (x TimeInForce) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TimeInForce) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_events_proto_enumTypes[3].Descriptor()
}

func (TimeInForce) Type() protoreflect.EnumType {
	return &file_internal_events_proto_enumTypes[3]
}

func (x TimeInForce) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TimeInForce.Descriptor instead.
func (TimeInForce) EnumDescriptor() ([]byte, []int) {
	return file_internal_events_proto_rawDescGZIP(), []int{3}
}

type InternalEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Event:
	//
	//	*InternalEvent_SubscribeInstruments
	//	*InternalEvent_NewOrderRequested
	//
	//dispatcher:generate
	Event isInternalEvent_Event `protobuf_oneof:"event"`
}

func (x *InternalEvent) Reset() {
	*x = InternalEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_events_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InternalEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InternalEvent) ProtoMessage() {}

func (x *InternalEvent) ProtoReflect() protoreflect.Message {
	mi := &file_internal_events_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InternalEvent.ProtoReflect.Descriptor instead.
func (*InternalEvent) Descriptor() ([]byte, []int) {
	return file_internal_events_proto_rawDescGZIP(), []int{0}
}

func (m *InternalEvent) GetEvent() isInternalEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *InternalEvent) GetSubscribeInstruments() *SubscribeInstruments {
	if x, ok := x.GetEvent().(*InternalEvent_SubscribeInstruments); ok {
		return x.SubscribeInstruments
	}
	return nil
}

func (x *InternalEvent) GetNewOrderRequested() *NewOrderRequested {
	if x, ok := x.GetEvent().(*InternalEvent_NewOrderRequested); ok {
		return x.NewOrderRequested
	}
	return nil
}

type isInternalEvent_Event interface {
	isInternalEvent_Event()
}

type InternalEvent_SubscribeInstruments struct {
	SubscribeInstruments *SubscribeInstruments `protobuf:"bytes,1,opt,name=subscribe_instruments,json=subscribeInstruments,proto3,oneof"`
}

type InternalEvent_NewOrderRequested struct {
	NewOrderRequested *NewOrderRequested `protobuf:"bytes,2,opt,name=new_order_requested,json=newOrderRequested,proto3,oneof"`
}

func (*InternalEvent_SubscribeInstruments) isInternalEvent_Event() {}

func (*InternalEvent_NewOrderRequested) isInternalEvent_Event() {}

type SubscribeInstruments struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instruments []*Instrument `protobuf:"bytes,1,rep,name=instruments,proto3" json:"instruments,omitempty"`
}

func (x *SubscribeInstruments) Reset() {
	*x = SubscribeInstruments{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_events_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeInstruments) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeInstruments) ProtoMessage() {}

func (x *SubscribeInstruments) ProtoReflect() protoreflect.Message {
	mi := &file_internal_events_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeInstruments.ProtoReflect.Descriptor instead.
func (*SubscribeInstruments) Descriptor() ([]byte, []int) {
	return file_internal_events_proto_rawDescGZIP(), []int{1}
}

func (x *SubscribeInstruments) GetInstruments() []*Instrument {
	if x != nil {
		return x.Instruments
	}
	return nil
}

type Instrument struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exchange       Exchange `protobuf:"varint,1,opt,name=Exchange,proto3,enum=intev.Exchange" json:"Exchange,omitempty"`
	LocalSymbol    string   `protobuf:"bytes,2,opt,name=LocalSymbol,proto3" json:"LocalSymbol,omitempty"`
	ExchangeSymbol string   `protobuf:"bytes,3,opt,name=ExchangeSymbol,proto3" json:"ExchangeSymbol,omitempty"`
}

func (x *Instrument) Reset() {
	*x = Instrument{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_events_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Instrument) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Instrument) ProtoMessage() {}

func (x *Instrument) ProtoReflect() protoreflect.Message {
	mi := &file_internal_events_proto_msgTypes[2]
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
	return file_internal_events_proto_rawDescGZIP(), []int{2}
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

type NewOrderRequested struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientOrderId   string             `protobuf:"bytes,1,opt,name=client_order_id,json=clientOrderId,proto3" json:"client_order_id,omitempty"`
	LocalInstrument *LocalInstrument   `protobuf:"bytes,2,opt,name=local_instrument,json=localInstrument,proto3" json:"local_instrument,omitempty"`
	Price           *decimalpb.Decimal `protobuf:"bytes,3,opt,name=price,proto3" json:"price,omitempty"`
	Qty             *decimalpb.Decimal `protobuf:"bytes,4,opt,name=qty,proto3" json:"qty,omitempty"`
	OrderSide       OrderSide          `protobuf:"varint,5,opt,name=order_side,json=orderSide,proto3,enum=intev.OrderSide" json:"order_side,omitempty"`
	Type            OrderType          `protobuf:"varint,6,opt,name=type,proto3,enum=intev.OrderType" json:"type,omitempty"`
	TimeInForce     TimeInForce        `protobuf:"varint,7,opt,name=time_in_force,json=timeInForce,proto3,enum=intev.TimeInForce" json:"time_in_force,omitempty"`
}

func (x *NewOrderRequested) Reset() {
	*x = NewOrderRequested{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_events_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewOrderRequested) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewOrderRequested) ProtoMessage() {}

func (x *NewOrderRequested) ProtoReflect() protoreflect.Message {
	mi := &file_internal_events_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewOrderRequested.ProtoReflect.Descriptor instead.
func (*NewOrderRequested) Descriptor() ([]byte, []int) {
	return file_internal_events_proto_rawDescGZIP(), []int{3}
}

func (x *NewOrderRequested) GetClientOrderId() string {
	if x != nil {
		return x.ClientOrderId
	}
	return ""
}

func (x *NewOrderRequested) GetLocalInstrument() *LocalInstrument {
	if x != nil {
		return x.LocalInstrument
	}
	return nil
}

func (x *NewOrderRequested) GetPrice() *decimalpb.Decimal {
	if x != nil {
		return x.Price
	}
	return nil
}

func (x *NewOrderRequested) GetQty() *decimalpb.Decimal {
	if x != nil {
		return x.Qty
	}
	return nil
}

func (x *NewOrderRequested) GetOrderSide() OrderSide {
	if x != nil {
		return x.OrderSide
	}
	return OrderSide_ORDER_SIDE_UNSPECIFIED
}

func (x *NewOrderRequested) GetType() OrderType {
	if x != nil {
		return x.Type
	}
	return OrderType_ORDER_TYPE_UNSPECIFIED
}

func (x *NewOrderRequested) GetTimeInForce() TimeInForce {
	if x != nil {
		return x.TimeInForce
	}
	return TimeInForce_TIME_IN_FORCE_UNSPECIFIED
}

type LocalInstrument struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exchange Exchange `protobuf:"varint,1,opt,name=exchange,proto3,enum=intev.Exchange" json:"exchange,omitempty"`
	Symbol   string   `protobuf:"bytes,2,opt,name=Symbol,proto3" json:"Symbol,omitempty"`
}

func (x *LocalInstrument) Reset() {
	*x = LocalInstrument{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_events_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocalInstrument) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocalInstrument) ProtoMessage() {}

func (x *LocalInstrument) ProtoReflect() protoreflect.Message {
	mi := &file_internal_events_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocalInstrument.ProtoReflect.Descriptor instead.
func (*LocalInstrument) Descriptor() ([]byte, []int) {
	return file_internal_events_proto_rawDescGZIP(), []int{4}
}

func (x *LocalInstrument) GetExchange() Exchange {
	if x != nil {
		return x.Exchange
	}
	return Exchange_ExchangeUnknown
}

func (x *LocalInstrument) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

var File_internal_events_proto protoreflect.FileDescriptor

var file_internal_events_proto_rawDesc = []byte{
	0x0a, 0x15, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x69, 0x6e, 0x74, 0x65, 0x76, 0x1a, 0x1b,
	0x70, 0x6b, 0x67, 0x2f, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x70, 0x62, 0x2f, 0x64, 0x65,
	0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb8, 0x01, 0x0a, 0x0d,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x52, 0x0a,
	0x15, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x5f, 0x69, 0x6e, 0x73, 0x74, 0x72,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x76, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x49, 0x6e,
	0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x48, 0x00, 0x52, 0x14, 0x73, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0x4a, 0x0a, 0x13, 0x6e, 0x65, 0x77, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x69, 0x6e, 0x74, 0x65, 0x76, 0x2e, 0x4e, 0x65, 0x77, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x48, 0x00, 0x52, 0x11, 0x6e, 0x65, 0x77, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x42, 0x07, 0x0a,
	0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x4b, 0x0a, 0x14, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x62, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x33,
	0x0a, 0x0b, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x76, 0x2e, 0x49, 0x6e, 0x73, 0x74,
	0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0b, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x22, 0x83, 0x01, 0x0a, 0x0a, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x2b, 0x0a, 0x08, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x76, 0x2e, 0x45, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x08, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x53, 0x79, 0x6d, 0x62, 0x6f,
	0x6c, 0x12, 0x26, 0x0a, 0x0e, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x53, 0x79, 0x6d,
	0x62, 0x6f, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x45, 0x78, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x22, 0xd9, 0x02, 0x0a, 0x11, 0x4e, 0x65,
	0x77, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x12,
	0x26, 0x0a, 0x0f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x41, 0x0a, 0x10, 0x6c, 0x6f, 0x63, 0x61, 0x6c,
	0x5f, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x76, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x49,
	0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0f, 0x6c, 0x6f, 0x63, 0x61, 0x6c,
	0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x26, 0x0a, 0x05, 0x70, 0x72,
	0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x64, 0x65, 0x63, 0x69,
	0x6d, 0x61, 0x6c, 0x2e, 0x44, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x52, 0x05, 0x70, 0x72, 0x69,
	0x63, 0x65, 0x12, 0x22, 0x0a, 0x03, 0x71, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x44, 0x65, 0x63, 0x69, 0x6d, 0x61,
	0x6c, 0x52, 0x03, 0x71, 0x74, 0x79, 0x12, 0x2f, 0x0a, 0x0a, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f,
	0x73, 0x69, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x69, 0x6e, 0x74,
	0x65, 0x76, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x69, 0x64, 0x65, 0x52, 0x09, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x53, 0x69, 0x64, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x76, 0x2e, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x36, 0x0a,
	0x0d, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x69, 0x6e, 0x5f, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x76, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x49, 0x6e, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x52, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x49, 0x6e,
	0x46, 0x6f, 0x72, 0x63, 0x65, 0x22, 0x56, 0x0a, 0x0f, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x49, 0x6e,
	0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x2b, 0x0a, 0x08, 0x65, 0x78, 0x63, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x69, 0x6e, 0x74,
	0x65, 0x76, 0x2e, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x08, 0x65, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x2a, 0x57, 0x0a,
	0x08, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x13, 0x0a, 0x0f, 0x45, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x0f,
	0x0a, 0x0b, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4f, 0x4b, 0x58, 0x10, 0x01, 0x12,
	0x12, 0x0a, 0x0e, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x47, 0x61, 0x74, 0x65, 0x49,
	0x4f, 0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x42,
	0x79, 0x62, 0x69, 0x74, 0x10, 0x03, 0x2a, 0x54, 0x0a, 0x09, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x15, 0x0a, 0x11, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x41,
	0x52, 0x4b, 0x45, 0x54, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x49, 0x4d, 0x49, 0x54, 0x10, 0x02, 0x2a, 0x50, 0x0a, 0x09,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x69, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x4f, 0x52, 0x44,
	0x45, 0x52, 0x5f, 0x53, 0x49, 0x44, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53,
	0x49, 0x44, 0x45, 0x5f, 0x42, 0x55, 0x59, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x4f, 0x52, 0x44,
	0x45, 0x52, 0x5f, 0x53, 0x49, 0x44, 0x45, 0x5f, 0x53, 0x45, 0x4c, 0x4c, 0x10, 0x02, 0x2a, 0x71,
	0x0a, 0x0b, 0x54, 0x69, 0x6d, 0x65, 0x49, 0x6e, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x12, 0x1d, 0x0a,
	0x19, 0x54, 0x49, 0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x43, 0x45, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11,
	0x54, 0x49, 0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x43, 0x45, 0x5f, 0x47, 0x54,
	0x43, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x49, 0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46,
	0x4f, 0x52, 0x43, 0x45, 0x5f, 0x49, 0x4f, 0x43, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x49,
	0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x43, 0x45, 0x5f, 0x46, 0x4f, 0x4b, 0x10,
	0x03, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x69, 0x6e, 0x74, 0x65, 0x76, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_events_proto_rawDescOnce sync.Once
	file_internal_events_proto_rawDescData = file_internal_events_proto_rawDesc
)

func file_internal_events_proto_rawDescGZIP() []byte {
	file_internal_events_proto_rawDescOnce.Do(func() {
		file_internal_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_events_proto_rawDescData)
	})
	return file_internal_events_proto_rawDescData
}

var file_internal_events_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_internal_events_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_internal_events_proto_goTypes = []interface{}{
	(Exchange)(0),                // 0: intev.Exchange
	(OrderType)(0),               // 1: intev.OrderType
	(OrderSide)(0),               // 2: intev.OrderSide
	(TimeInForce)(0),             // 3: intev.TimeInForce
	(*InternalEvent)(nil),        // 4: intev.InternalEvent
	(*SubscribeInstruments)(nil), // 5: intev.SubscribeInstruments
	(*Instrument)(nil),           // 6: intev.Instrument
	(*NewOrderRequested)(nil),    // 7: intev.NewOrderRequested
	(*LocalInstrument)(nil),      // 8: intev.LocalInstrument
	(*decimalpb.Decimal)(nil),    // 9: decimal.Decimal
}
var file_internal_events_proto_depIdxs = []int32{
	5,  // 0: intev.InternalEvent.subscribe_instruments:type_name -> intev.SubscribeInstruments
	7,  // 1: intev.InternalEvent.new_order_requested:type_name -> intev.NewOrderRequested
	6,  // 2: intev.SubscribeInstruments.instruments:type_name -> intev.Instrument
	0,  // 3: intev.Instrument.Exchange:type_name -> intev.Exchange
	8,  // 4: intev.NewOrderRequested.local_instrument:type_name -> intev.LocalInstrument
	9,  // 5: intev.NewOrderRequested.price:type_name -> decimal.Decimal
	9,  // 6: intev.NewOrderRequested.qty:type_name -> decimal.Decimal
	2,  // 7: intev.NewOrderRequested.order_side:type_name -> intev.OrderSide
	1,  // 8: intev.NewOrderRequested.type:type_name -> intev.OrderType
	3,  // 9: intev.NewOrderRequested.time_in_force:type_name -> intev.TimeInForce
	0,  // 10: intev.LocalInstrument.exchange:type_name -> intev.Exchange
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_internal_events_proto_init() }
func file_internal_events_proto_init() {
	if File_internal_events_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_events_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InternalEvent); i {
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
		file_internal_events_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeInstruments); i {
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
		file_internal_events_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_events_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewOrderRequested); i {
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
		file_internal_events_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LocalInstrument); i {
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
	file_internal_events_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*InternalEvent_SubscribeInstruments)(nil),
		(*InternalEvent_NewOrderRequested)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_events_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_events_proto_goTypes,
		DependencyIndexes: file_internal_events_proto_depIdxs,
		EnumInfos:         file_internal_events_proto_enumTypes,
		MessageInfos:      file_internal_events_proto_msgTypes,
	}.Build()
	File_internal_events_proto = out.File
	file_internal_events_proto_rawDesc = nil
	file_internal_events_proto_goTypes = nil
	file_internal_events_proto_depIdxs = nil
}
