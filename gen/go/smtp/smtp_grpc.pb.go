// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: smtp.proto

package smtp

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SmtpPop3Service_SendEmail_FullMethodName          = "/proto.SmtpPop3Service/SendEmail"
	SmtpPop3Service_ForwardEmail_FullMethodName       = "/proto.SmtpPop3Service/ForwardEmail"
	SmtpPop3Service_ReplyEmail_FullMethodName         = "/proto.SmtpPop3Service/ReplyEmail"
	SmtpPop3Service_FetchEmailsViaPOP3_FullMethodName = "/proto.SmtpPop3Service/FetchEmailsViaPOP3"
)

// SmtpPop3ServiceClient is the client API for SmtpPop3Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SmtpPop3ServiceClient interface {
	SendEmail(ctx context.Context, in *SendEmailRequest, opts ...grpc.CallOption) (*SendEmailReply, error)
	ForwardEmail(ctx context.Context, in *ForwardEmailRequest, opts ...grpc.CallOption) (*ForwardEmailReply, error)
	ReplyEmail(ctx context.Context, in *ReplyEmailRequest, opts ...grpc.CallOption) (*ReplyEmailReply, error)
	FetchEmailsViaPOP3(ctx context.Context, in *FetchEmailsViaPOP3Request, opts ...grpc.CallOption) (*FetchEmailsViaPOP3Reply, error)
}

type smtpPop3ServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSmtpPop3ServiceClient(cc grpc.ClientConnInterface) SmtpPop3ServiceClient {
	return &smtpPop3ServiceClient{cc}
}

func (c *smtpPop3ServiceClient) SendEmail(ctx context.Context, in *SendEmailRequest, opts ...grpc.CallOption) (*SendEmailReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendEmailReply)
	err := c.cc.Invoke(ctx, SmtpPop3Service_SendEmail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *smtpPop3ServiceClient) ForwardEmail(ctx context.Context, in *ForwardEmailRequest, opts ...grpc.CallOption) (*ForwardEmailReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ForwardEmailReply)
	err := c.cc.Invoke(ctx, SmtpPop3Service_ForwardEmail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *smtpPop3ServiceClient) ReplyEmail(ctx context.Context, in *ReplyEmailRequest, opts ...grpc.CallOption) (*ReplyEmailReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReplyEmailReply)
	err := c.cc.Invoke(ctx, SmtpPop3Service_ReplyEmail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *smtpPop3ServiceClient) FetchEmailsViaPOP3(ctx context.Context, in *FetchEmailsViaPOP3Request, opts ...grpc.CallOption) (*FetchEmailsViaPOP3Reply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FetchEmailsViaPOP3Reply)
	err := c.cc.Invoke(ctx, SmtpPop3Service_FetchEmailsViaPOP3_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SmtpPop3ServiceServer is the server API for SmtpPop3Service service.
// All implementations must embed UnimplementedSmtpPop3ServiceServer
// for forward compatibility.
type SmtpPop3ServiceServer interface {
	SendEmail(context.Context, *SendEmailRequest) (*SendEmailReply, error)
	ForwardEmail(context.Context, *ForwardEmailRequest) (*ForwardEmailReply, error)
	ReplyEmail(context.Context, *ReplyEmailRequest) (*ReplyEmailReply, error)
	FetchEmailsViaPOP3(context.Context, *FetchEmailsViaPOP3Request) (*FetchEmailsViaPOP3Reply, error)
	mustEmbedUnimplementedSmtpPop3ServiceServer()
}

// UnimplementedSmtpPop3ServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSmtpPop3ServiceServer struct{}

func (UnimplementedSmtpPop3ServiceServer) SendEmail(context.Context, *SendEmailRequest) (*SendEmailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmail not implemented")
}
func (UnimplementedSmtpPop3ServiceServer) ForwardEmail(context.Context, *ForwardEmailRequest) (*ForwardEmailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForwardEmail not implemented")
}
func (UnimplementedSmtpPop3ServiceServer) ReplyEmail(context.Context, *ReplyEmailRequest) (*ReplyEmailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplyEmail not implemented")
}
func (UnimplementedSmtpPop3ServiceServer) FetchEmailsViaPOP3(context.Context, *FetchEmailsViaPOP3Request) (*FetchEmailsViaPOP3Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchEmailsViaPOP3 not implemented")
}
func (UnimplementedSmtpPop3ServiceServer) mustEmbedUnimplementedSmtpPop3ServiceServer() {}
func (UnimplementedSmtpPop3ServiceServer) testEmbeddedByValue()                         {}

// UnsafeSmtpPop3ServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SmtpPop3ServiceServer will
// result in compilation errors.
type UnsafeSmtpPop3ServiceServer interface {
	mustEmbedUnimplementedSmtpPop3ServiceServer()
}

func RegisterSmtpPop3ServiceServer(s grpc.ServiceRegistrar, srv SmtpPop3ServiceServer) {
	// If the following call pancis, it indicates UnimplementedSmtpPop3ServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SmtpPop3Service_ServiceDesc, srv)
}

func _SmtpPop3Service_SendEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SmtpPop3ServiceServer).SendEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SmtpPop3Service_SendEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SmtpPop3ServiceServer).SendEmail(ctx, req.(*SendEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SmtpPop3Service_ForwardEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SmtpPop3ServiceServer).ForwardEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SmtpPop3Service_ForwardEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SmtpPop3ServiceServer).ForwardEmail(ctx, req.(*ForwardEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SmtpPop3Service_ReplyEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplyEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SmtpPop3ServiceServer).ReplyEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SmtpPop3Service_ReplyEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SmtpPop3ServiceServer).ReplyEmail(ctx, req.(*ReplyEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SmtpPop3Service_FetchEmailsViaPOP3_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchEmailsViaPOP3Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SmtpPop3ServiceServer).FetchEmailsViaPOP3(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SmtpPop3Service_FetchEmailsViaPOP3_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SmtpPop3ServiceServer).FetchEmailsViaPOP3(ctx, req.(*FetchEmailsViaPOP3Request))
	}
	return interceptor(ctx, in, info, handler)
}

// SmtpPop3Service_ServiceDesc is the grpc.ServiceDesc for SmtpPop3Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SmtpPop3Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.SmtpPop3Service",
	HandlerType: (*SmtpPop3ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEmail",
			Handler:    _SmtpPop3Service_SendEmail_Handler,
		},
		{
			MethodName: "ForwardEmail",
			Handler:    _SmtpPop3Service_ForwardEmail_Handler,
		},
		{
			MethodName: "ReplyEmail",
			Handler:    _SmtpPop3Service_ReplyEmail_Handler,
		},
		{
			MethodName: "FetchEmailsViaPOP3",
			Handler:    _SmtpPop3Service_FetchEmailsViaPOP3_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "smtp.proto",
}