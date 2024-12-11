// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: webhooks/inboundwebhooksapi/v1/api.proto

package inboundwebhooksapiv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	v1 "libs/backend/proto-gen/go/common/v1"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// UserRegisteredRequest defines the incoming data for the user created event from
// Auth0 Actions after the user has registered for the first time
type UserRegisteredRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FirstName            string                `protobuf:"bytes,1,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName             string                `protobuf:"bytes,2,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Nickname             string                `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Username             string                `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	EmailAddress         string                `protobuf:"bytes,5,opt,name=email_address,json=emailAddress,proto3" json:"email_address,omitempty"`
	EmailAddressVerified bool                  `protobuf:"varint,6,opt,name=email_address_verified,json=emailAddressVerified,proto3" json:"email_address_verified,omitempty"`
	PhoneNumber          string                `protobuf:"bytes,7,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	PhoneNumberVerified  bool                  `protobuf:"varint,8,opt,name=phone_number_verified,json=phoneNumberVerified,proto3" json:"phone_number_verified,omitempty"`
	Strategy             string                `protobuf:"bytes,9,opt,name=strategy,proto3" json:"strategy,omitempty"`
	UserMetadata         map[string]*anypb.Any `protobuf:"bytes,10,rep,name=user_metadata,json=userMetadata,proto3" json:"user_metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *UserRegisteredRequest) Reset() {
	*x = UserRegisteredRequest{}
	mi := &file_webhooks_inboundwebhooksapi_v1_api_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserRegisteredRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserRegisteredRequest) ProtoMessage() {}

func (x *UserRegisteredRequest) ProtoReflect() protoreflect.Message {
	mi := &file_webhooks_inboundwebhooksapi_v1_api_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserRegisteredRequest.ProtoReflect.Descriptor instead.
func (*UserRegisteredRequest) Descriptor() ([]byte, []int) {
	return file_webhooks_inboundwebhooksapi_v1_api_proto_rawDescGZIP(), []int{0}
}

func (x *UserRegisteredRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *UserRegisteredRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *UserRegisteredRequest) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *UserRegisteredRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UserRegisteredRequest) GetEmailAddress() string {
	if x != nil {
		return x.EmailAddress
	}
	return ""
}

func (x *UserRegisteredRequest) GetEmailAddressVerified() bool {
	if x != nil {
		return x.EmailAddressVerified
	}
	return false
}

func (x *UserRegisteredRequest) GetPhoneNumber() string {
	if x != nil {
		return x.PhoneNumber
	}
	return ""
}

func (x *UserRegisteredRequest) GetPhoneNumberVerified() bool {
	if x != nil {
		return x.PhoneNumberVerified
	}
	return false
}

func (x *UserRegisteredRequest) GetStrategy() string {
	if x != nil {
		return x.Strategy
	}
	return ""
}

func (x *UserRegisteredRequest) GetUserMetadata() map[string]*anypb.Any {
	if x != nil {
		return x.UserMetadata
	}
	return nil
}

var File_webhooks_inboundwebhooksapi_v1_api_proto protoreflect.FileDescriptor

var file_webhooks_inboundwebhooksapi_v1_api_proto_rawDesc = []byte{
	0x0a, 0x28, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2f, 0x69, 0x6e, 0x62, 0x6f, 0x75,
	0x6e, 0x64, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x69, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x15, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdd, 0x04, 0x0a, 0x15, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x26, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x09,
	0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x09, 0x6c, 0x61, 0x73,
	0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48,
	0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x23, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2c, 0x0a, 0x0d, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x60, 0x01, 0x52, 0x0c, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x34, 0x0a, 0x16, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x14, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x2a, 0x0a,
	0x0c, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x0b, 0x70, 0x68,
	0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x32, 0x0a, 0x15, 0x70, 0x68, 0x6f,
	0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69,
	0x65, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x13, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x23, 0x0a,
	0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65,
	0x67, 0x79, 0x12, 0x6c, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x47, 0x2e, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x69, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x1a, 0x55, 0x0a, 0x11, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0xaa, 0x01, 0x0a, 0x16, 0x49, 0x6e, 0x62, 0x6f,
	0x75, 0x6e, 0x64, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x41, 0x75, 0x74, 0x68, 0x41,
	0x50, 0x49, 0x12, 0x8f, 0x01, 0x0a, 0x0e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x65, 0x64, 0x12, 0x35, 0x2e, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73,
	0x2e, 0x69, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x34,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2e, 0x3a, 0x01, 0x2a, 0x22, 0x29, 0x2f, 0x76, 0x31, 0x2f, 0x77,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2f, 0x69, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x2f,
	0x61, 0x75, 0x74, 0x68, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x65, 0x64, 0x42, 0x97, 0x02, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x65, 0x62,
	0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x69, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x77, 0x65, 0x62,
	0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x42, 0x08, 0x41, 0x70, 0x69,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4d, 0x6c, 0x69, 0x62, 0x73, 0x2f, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2d, 0x67, 0x65, 0x6e, 0x2f,
	0x67, 0x6f, 0x2f, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2f, 0x69, 0x6e, 0x62, 0x6f,
	0x75, 0x6e, 0x64, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x3b, 0x69, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b,
	0x73, 0x61, 0x70, 0x69, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x57, 0x49, 0x58, 0xaa, 0x02, 0x1e, 0x57,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x77,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1e,
	0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x5c, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64,
	0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x61, 0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x2a, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x5c, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e,
	0x64, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x61, 0x70, 0x69, 0x5c, 0x56, 0x31, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x20, 0x57, 0x65,
	0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x3a, 0x3a, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x77,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x61, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_webhooks_inboundwebhooksapi_v1_api_proto_rawDescOnce sync.Once
	file_webhooks_inboundwebhooksapi_v1_api_proto_rawDescData = file_webhooks_inboundwebhooksapi_v1_api_proto_rawDesc
)

func file_webhooks_inboundwebhooksapi_v1_api_proto_rawDescGZIP() []byte {
	file_webhooks_inboundwebhooksapi_v1_api_proto_rawDescOnce.Do(func() {
		file_webhooks_inboundwebhooksapi_v1_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_webhooks_inboundwebhooksapi_v1_api_proto_rawDescData)
	})
	return file_webhooks_inboundwebhooksapi_v1_api_proto_rawDescData
}

var file_webhooks_inboundwebhooksapi_v1_api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_webhooks_inboundwebhooksapi_v1_api_proto_goTypes = []any{
	(*UserRegisteredRequest)(nil), // 0: webhooks.inboundwebhooksapi.v1.UserRegisteredRequest
	nil,                           // 1: webhooks.inboundwebhooksapi.v1.UserRegisteredRequest.UserMetadataEntry
	(*anypb.Any)(nil),             // 2: google.protobuf.Any
	(*v1.Empty)(nil),              // 3: common.v1.Empty
}
var file_webhooks_inboundwebhooksapi_v1_api_proto_depIdxs = []int32{
	1, // 0: webhooks.inboundwebhooksapi.v1.UserRegisteredRequest.user_metadata:type_name -> webhooks.inboundwebhooksapi.v1.UserRegisteredRequest.UserMetadataEntry
	2, // 1: webhooks.inboundwebhooksapi.v1.UserRegisteredRequest.UserMetadataEntry.value:type_name -> google.protobuf.Any
	0, // 2: webhooks.inboundwebhooksapi.v1.InboundWebhooksAuthAPI.UserRegistered:input_type -> webhooks.inboundwebhooksapi.v1.UserRegisteredRequest
	3, // 3: webhooks.inboundwebhooksapi.v1.InboundWebhooksAuthAPI.UserRegistered:output_type -> common.v1.Empty
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_webhooks_inboundwebhooksapi_v1_api_proto_init() }
func file_webhooks_inboundwebhooksapi_v1_api_proto_init() {
	if File_webhooks_inboundwebhooksapi_v1_api_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_webhooks_inboundwebhooksapi_v1_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_webhooks_inboundwebhooksapi_v1_api_proto_goTypes,
		DependencyIndexes: file_webhooks_inboundwebhooksapi_v1_api_proto_depIdxs,
		MessageInfos:      file_webhooks_inboundwebhooksapi_v1_api_proto_msgTypes,
	}.Build()
	File_webhooks_inboundwebhooksapi_v1_api_proto = out.File
	file_webhooks_inboundwebhooksapi_v1_api_proto_rawDesc = nil
	file_webhooks_inboundwebhooksapi_v1_api_proto_goTypes = nil
	file_webhooks_inboundwebhooksapi_v1_api_proto_depIdxs = nil
}
