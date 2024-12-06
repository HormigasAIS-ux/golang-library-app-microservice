// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0
// source: author.proto

package author_grpc

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
	AuthorService_CreateAuthor_FullMethodName        = "/author_service.AuthorService/CreateAuthor"
	AuthorService_GetAuthorByUserUUID_FullMethodName = "/author_service.AuthorService/GetAuthorByUserUUID"
)

// AuthorServiceClient is the client API for AuthorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthorServiceClient interface {
	CreateAuthor(ctx context.Context, in *CreateAuthorReq, opts ...grpc.CallOption) (*CreateAuthorResp, error)
	GetAuthorByUserUUID(ctx context.Context, in *GetAuthorByUserUUIDReq, opts ...grpc.CallOption) (*GetAuthorByUserUUIDResp, error)
}

type authorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthorServiceClient(cc grpc.ClientConnInterface) AuthorServiceClient {
	return &authorServiceClient{cc}
}

func (c *authorServiceClient) CreateAuthor(ctx context.Context, in *CreateAuthorReq, opts ...grpc.CallOption) (*CreateAuthorResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAuthorResp)
	err := c.cc.Invoke(ctx, AuthorService_CreateAuthor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorServiceClient) GetAuthorByUserUUID(ctx context.Context, in *GetAuthorByUserUUIDReq, opts ...grpc.CallOption) (*GetAuthorByUserUUIDResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAuthorByUserUUIDResp)
	err := c.cc.Invoke(ctx, AuthorService_GetAuthorByUserUUID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorServiceServer is the server API for AuthorService service.
// All implementations must embed UnimplementedAuthorServiceServer
// for forward compatibility.
type AuthorServiceServer interface {
	CreateAuthor(context.Context, *CreateAuthorReq) (*CreateAuthorResp, error)
	GetAuthorByUserUUID(context.Context, *GetAuthorByUserUUIDReq) (*GetAuthorByUserUUIDResp, error)
	mustEmbedUnimplementedAuthorServiceServer()
}

// UnimplementedAuthorServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuthorServiceServer struct{}

func (UnimplementedAuthorServiceServer) CreateAuthor(context.Context, *CreateAuthorReq) (*CreateAuthorResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAuthor not implemented")
}
func (UnimplementedAuthorServiceServer) GetAuthorByUserUUID(context.Context, *GetAuthorByUserUUIDReq) (*GetAuthorByUserUUIDResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthorByUserUUID not implemented")
}
func (UnimplementedAuthorServiceServer) mustEmbedUnimplementedAuthorServiceServer() {}
func (UnimplementedAuthorServiceServer) testEmbeddedByValue()                       {}

// UnsafeAuthorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthorServiceServer will
// result in compilation errors.
type UnsafeAuthorServiceServer interface {
	mustEmbedUnimplementedAuthorServiceServer()
}

func RegisterAuthorServiceServer(s grpc.ServiceRegistrar, srv AuthorServiceServer) {
	// If the following call pancis, it indicates UnimplementedAuthorServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AuthorService_ServiceDesc, srv)
}

func _AuthorService_CreateAuthor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAuthorReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorServiceServer).CreateAuthor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthorService_CreateAuthor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorServiceServer).CreateAuthor(ctx, req.(*CreateAuthorReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorService_GetAuthorByUserUUID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuthorByUserUUIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorServiceServer).GetAuthorByUserUUID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthorService_GetAuthorByUserUUID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorServiceServer).GetAuthorByUserUUID(ctx, req.(*GetAuthorByUserUUIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthorService_ServiceDesc is the grpc.ServiceDesc for AuthorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "author_service.AuthorService",
	HandlerType: (*AuthorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAuthor",
			Handler:    _AuthorService_CreateAuthor_Handler,
		},
		{
			MethodName: "GetAuthorByUserUUID",
			Handler:    _AuthorService_GetAuthorByUserUUID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "author.proto",
}
