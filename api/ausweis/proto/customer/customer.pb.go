// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: customer/customer.proto

package customer

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Customer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId                  int64    `protobuf:"varint,1,opt,name=customerId,proto3" json:"customerId,omitempty"`
	Email                       string   `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Phone                       *string  `protobuf:"bytes,3,opt,name=phone,proto3,oneof" json:"phone,omitempty"`
	PhoneSms                    *string  `protobuf:"bytes,4,opt,name=phoneSms,proto3,oneof" json:"phoneSms,omitempty"`
	Firstname                   string   `protobuf:"bytes,5,opt,name=firstname,proto3" json:"firstname,omitempty"`
	Lastname                    string   `protobuf:"bytes,6,opt,name=lastname,proto3" json:"lastname,omitempty"`
	Fullname                    string   `protobuf:"bytes,7,opt,name=fullname,proto3" json:"fullname,omitempty"`
	Company                     *string  `protobuf:"bytes,8,opt,name=company,proto3,oneof" json:"company,omitempty"`
	Blocked                     bool     `protobuf:"varint,9,opt,name=blocked,proto3" json:"blocked,omitempty"`
	Landed                      bool     `protobuf:"varint,10,opt,name=landed,proto3" json:"landed,omitempty"`
	CreatedAt                   string   `protobuf:"bytes,11,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	CheckoutPaymentRestrictions []string `protobuf:"bytes,12,rep,name=checkoutPaymentRestrictions,proto3" json:"checkoutPaymentRestrictions,omitempty"`
}

func (x *Customer) Reset() {
	*x = Customer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_customer_customer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Customer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Customer) ProtoMessage() {}

func (x *Customer) ProtoReflect() protoreflect.Message {
	mi := &file_customer_customer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Customer.ProtoReflect.Descriptor instead.
func (*Customer) Descriptor() ([]byte, []int) {
	return file_customer_customer_proto_rawDescGZIP(), []int{0}
}

func (x *Customer) GetCustomerId() int64 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

func (x *Customer) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Customer) GetPhone() string {
	if x != nil && x.Phone != nil {
		return *x.Phone
	}
	return ""
}

func (x *Customer) GetPhoneSms() string {
	if x != nil && x.PhoneSms != nil {
		return *x.PhoneSms
	}
	return ""
}

func (x *Customer) GetFirstname() string {
	if x != nil {
		return x.Firstname
	}
	return ""
}

func (x *Customer) GetLastname() string {
	if x != nil {
		return x.Lastname
	}
	return ""
}

func (x *Customer) GetFullname() string {
	if x != nil {
		return x.Fullname
	}
	return ""
}

func (x *Customer) GetCompany() string {
	if x != nil && x.Company != nil {
		return *x.Company
	}
	return ""
}

func (x *Customer) GetBlocked() bool {
	if x != nil {
		return x.Blocked
	}
	return false
}

func (x *Customer) GetLanded() bool {
	if x != nil {
		return x.Landed
	}
	return false
}

func (x *Customer) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *Customer) GetCheckoutPaymentRestrictions() []string {
	if x != nil {
		return x.CheckoutPaymentRestrictions
	}
	return nil
}

type CustomersList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Customer `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Total int64       `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *CustomersList) Reset() {
	*x = CustomersList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_customer_customer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CustomersList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CustomersList) ProtoMessage() {}

func (x *CustomersList) ProtoReflect() protoreflect.Message {
	mi := &file_customer_customer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CustomersList.ProtoReflect.Descriptor instead.
func (*CustomersList) Descriptor() ([]byte, []int) {
	return file_customer_customer_proto_rawDescGZIP(), []int{1}
}

func (x *CustomersList) GetItems() []*Customer {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *CustomersList) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_customer_customer_proto protoreflect.FileDescriptor

var file_customer_customer_proto_rawDesc = []byte{
	0x0a, 0x17, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x2f, 0x63, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x63, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x65, 0x72, 0x22, 0xa6, 0x03, 0x0a, 0x08, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72,
	0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x19, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x1f, 0x0a, 0x08, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x53, 0x6d, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x53, 0x6d, 0x73, 0x88,
	0x01, 0x01, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x66, 0x75, 0x6c, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x66, 0x75, 0x6c, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x70,
	0x61, 0x6e, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x07, 0x63, 0x6f, 0x6d,
	0x70, 0x61, 0x6e, 0x79, 0x88, 0x01, 0x01, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x65,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x6e, 0x64, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x6c, 0x61, 0x6e, 0x64, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x40, 0x0a, 0x1b, 0x63, 0x68, 0x65, 0x63, 0x6b,
	0x6f, 0x75, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x74, 0x72, 0x69,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x09, 0x52, 0x1b, 0x63, 0x68,
	0x65, 0x63, 0x6b, 0x6f, 0x75, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x74, 0x72, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x70, 0x68,
	0x6f, 0x6e, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x53, 0x6d, 0x73,
	0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x22, 0x4f, 0x0a, 0x0d,
	0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x28, 0x0a,
	0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63,
	0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x2e, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72,
	0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x42, 0x13, 0x5a,
	0x11, 0x2e, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_customer_customer_proto_rawDescOnce sync.Once
	file_customer_customer_proto_rawDescData = file_customer_customer_proto_rawDesc
)

func file_customer_customer_proto_rawDescGZIP() []byte {
	file_customer_customer_proto_rawDescOnce.Do(func() {
		file_customer_customer_proto_rawDescData = protoimpl.X.CompressGZIP(file_customer_customer_proto_rawDescData)
	})
	return file_customer_customer_proto_rawDescData
}

var file_customer_customer_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_customer_customer_proto_goTypes = []any{
	(*Customer)(nil),      // 0: customer.Customer
	(*CustomersList)(nil), // 1: customer.CustomersList
}
var file_customer_customer_proto_depIdxs = []int32{
	0, // 0: customer.CustomersList.items:type_name -> customer.Customer
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_customer_customer_proto_init() }
func file_customer_customer_proto_init() {
	if File_customer_customer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_customer_customer_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Customer); i {
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
		file_customer_customer_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CustomersList); i {
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
	file_customer_customer_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_customer_customer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_customer_customer_proto_goTypes,
		DependencyIndexes: file_customer_customer_proto_depIdxs,
		MessageInfos:      file_customer_customer_proto_msgTypes,
	}.Build()
	File_customer_customer_proto = out.File
	file_customer_customer_proto_rawDesc = nil
	file_customer_customer_proto_goTypes = nil
	file_customer_customer_proto_depIdxs = nil
}
