// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sync.proto

package protos

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
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

type Test struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Test) Reset()         { *m = Test{} }
func (m *Test) String() string { return proto.CompactTextString(m) }
func (*Test) ProtoMessage()    {}
func (*Test) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{0}
}

func (m *Test) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Test.Unmarshal(m, b)
}
func (m *Test) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Test.Marshal(b, m, deterministic)
}
func (m *Test) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Test.Merge(m, src)
}
func (m *Test) XXX_Size() int {
	return xxx_messageInfo_Test.Size(m)
}
func (m *Test) XXX_DiscardUnknown() {
	xxx_messageInfo_Test.DiscardUnknown(m)
}

var xxx_messageInfo_Test proto.InternalMessageInfo

func (m *Test) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Test)(nil), "protos.Test")
}

func init() { proto.RegisterFile("sync.proto", fileDescriptor_5273b98214de8075) }

var fileDescriptor_5273b98214de8075 = []byte{
	// 137 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0xae, 0xcc, 0x4b,
	0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x03, 0x53, 0xc5, 0x52, 0x32, 0xe9, 0xf9, 0xf9,
	0xe9, 0x39, 0xa9, 0xfa, 0x89, 0x05, 0x99, 0xfa, 0x89, 0x79, 0x79, 0xf9, 0x25, 0x89, 0x25, 0x99,
	0xf9, 0x79, 0xc5, 0x10, 0x55, 0x4a, 0x52, 0x5c, 0x2c, 0x21, 0xa9, 0xc5, 0x25, 0x42, 0x42, 0x5c,
	0x2c, 0x79, 0x89, 0xb9, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x91, 0x2d,
	0x17, 0x2b, 0x48, 0xae, 0x58, 0xc8, 0x84, 0x8b, 0xdd, 0x3d, 0xb5, 0x04, 0xac, 0x8e, 0x07, 0xa2,
	0xaf, 0x58, 0x0f, 0xc4, 0x93, 0x42, 0xe1, 0x29, 0xf1, 0x36, 0x5d, 0x7e, 0x32, 0x99, 0x89, 0x5d,
	0x88, 0x55, 0xbf, 0x24, 0xb5, 0xb8, 0x24, 0x09, 0xe2, 0x00, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x2e, 0xb4, 0x8a, 0xd0, 0x95, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TestsClient is the client API for Tests service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TestsClient interface {
	// Stub test
	GetTest(ctx context.Context, in *Test, opts ...grpc.CallOption) (*Test, error)
}

type testsClient struct {
	cc *grpc.ClientConn
}

func NewTestsClient(cc *grpc.ClientConn) TestsClient {
	return &testsClient{cc}
}

func (c *testsClient) GetTest(ctx context.Context, in *Test, opts ...grpc.CallOption) (*Test, error) {
	out := new(Test)
	err := c.cc.Invoke(ctx, "/protos.Tests/GetTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TestsServer is the server API for Tests service.
type TestsServer interface {
	// Stub test
	GetTest(context.Context, *Test) (*Test, error)
}

func RegisterTestsServer(s *grpc.Server, srv TestsServer) {
	s.RegisterService(&_Tests_serviceDesc, srv)
}

func _Tests_GetTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Test)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestsServer).GetTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Tests/GetTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestsServer).GetTest(ctx, req.(*Test))
	}
	return interceptor(ctx, in, info, handler)
}

var _Tests_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Tests",
	HandlerType: (*TestsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTest",
			Handler:    _Tests_GetTest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sync.proto",
}
