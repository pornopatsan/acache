// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ACacheClient is the client API for ACache service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ACacheClient interface {
	Save(ctx context.Context, in *Item, opts ...grpc.CallOption) (*Response, error)
	Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (*ItemResponse, error)
	Remove(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Response, error)
}

type aCacheClient struct {
	cc grpc.ClientConnInterface
}

func NewACacheClient(cc grpc.ClientConnInterface) ACacheClient {
	return &aCacheClient{cc}
}

func (c *aCacheClient) Save(ctx context.Context, in *Item, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/src.ACache/Save", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aCacheClient) Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (*ItemResponse, error) {
	out := new(ItemResponse)
	err := c.cc.Invoke(ctx, "/src.ACache/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aCacheClient) Remove(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/src.ACache/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ACacheServer is the server API for ACache service.
// All implementations must embed UnimplementedACacheServer
// for forward compatibility
type ACacheServer interface {
	Save(context.Context, *Item) (*Response, error)
	Get(context.Context, *Key) (*ItemResponse, error)
	Remove(context.Context, *Key) (*Response, error)
	mustEmbedUnimplementedACacheServer()
}

// UnimplementedACacheServer must be embedded to have forward compatible implementations.
type UnimplementedACacheServer struct {
}

func (UnimplementedACacheServer) Save(context.Context, *Item) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedACacheServer) Get(context.Context, *Key) (*ItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedACacheServer) Remove(context.Context, *Key) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (UnimplementedACacheServer) mustEmbedUnimplementedACacheServer() {}

// UnsafeACacheServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ACacheServer will
// result in compilation errors.
type UnsafeACacheServer interface {
	mustEmbedUnimplementedACacheServer()
}

func RegisterACacheServer(s grpc.ServiceRegistrar, srv ACacheServer) {
	s.RegisterService(&ACache_ServiceDesc, srv)
}

func _ACache_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Item)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ACacheServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/src.ACache/Save",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ACacheServer).Save(ctx, req.(*Item))
	}
	return interceptor(ctx, in, info, handler)
}

func _ACache_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ACacheServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/src.ACache/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ACacheServer).Get(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _ACache_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ACacheServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/src.ACache/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ACacheServer).Remove(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

// ACache_ServiceDesc is the grpc.ServiceDesc for ACache service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ACache_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "src.ACache",
	HandlerType: (*ACacheServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Save",
			Handler:    _ACache_Save_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ACache_Get_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _ACache_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "item.proto",
}