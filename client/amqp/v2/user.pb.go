// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

package v2

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = proto.Marshal
	_ = fmt.Errorf
	_ = math.Inf
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type User struct {
	Firstname            *string  `protobuf:"bytes,1,req,name=Firstname" json:"Firstname,omitempty"`
	Lastname             *string  `protobuf:"bytes,2,opt,name=Lastname" json:"Lastname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{0}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}

func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}

func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}

func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}

func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetFirstname() string {
	if m != nil && m.Firstname != nil {
		return *m.Firstname
	}
	return ""
}

func (m *User) GetLastname() string {
	if m != nil && m.Lastname != nil {
		return *m.Lastname
	}
	return ""
}

func init() {
	proto.RegisterType((*User)(nil), "amqp.User")
}

func init() { proto.RegisterFile("user.proto", fileDescriptor_116e343673f7ffaf) }

var fileDescriptor_116e343673f7ffaf = []byte{
	// 88 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x2d, 0x4e, 0x2d,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0xcc, 0x2d, 0x2c, 0x50, 0x72, 0xe0, 0x62,
	0x09, 0x2d, 0x4e, 0x2d, 0x12, 0x92, 0xe1, 0xe2, 0x74, 0xcb, 0x2c, 0x2a, 0x2e, 0xc9, 0x4b, 0xcc,
	0x4d, 0x95, 0x60, 0x54, 0x60, 0xd2, 0xe0, 0x0c, 0x42, 0x08, 0x08, 0x49, 0x71, 0x71, 0xf8, 0x24,
	0x42, 0x25, 0x99, 0x14, 0x18, 0x35, 0x38, 0x83, 0xe0, 0x7c, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x87, 0x2e, 0xfc, 0x85, 0x54, 0x00, 0x00, 0x00,
}
