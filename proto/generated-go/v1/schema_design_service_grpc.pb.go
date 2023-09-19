// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: v1/schema_design_service.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	SchemaDesignService_GetSchemaDesign_FullMethodName    = "/bytebase.v1.SchemaDesignService/GetSchemaDesign"
	SchemaDesignService_ListSchemaDesigns_FullMethodName  = "/bytebase.v1.SchemaDesignService/ListSchemaDesigns"
	SchemaDesignService_CreateSchemaDesign_FullMethodName = "/bytebase.v1.SchemaDesignService/CreateSchemaDesign"
	SchemaDesignService_UpdateSchemaDesign_FullMethodName = "/bytebase.v1.SchemaDesignService/UpdateSchemaDesign"
	SchemaDesignService_MergeSchemaDesign_FullMethodName  = "/bytebase.v1.SchemaDesignService/MergeSchemaDesign"
	SchemaDesignService_ParseSchemaString_FullMethodName  = "/bytebase.v1.SchemaDesignService/ParseSchemaString"
	SchemaDesignService_DeleteSchemaDesign_FullMethodName = "/bytebase.v1.SchemaDesignService/DeleteSchemaDesign"
	SchemaDesignService_DiffMetadata_FullMethodName       = "/bytebase.v1.SchemaDesignService/DiffMetadata"
)

// SchemaDesignServiceClient is the client API for SchemaDesignService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SchemaDesignServiceClient interface {
	GetSchemaDesign(ctx context.Context, in *GetSchemaDesignRequest, opts ...grpc.CallOption) (*SchemaDesign, error)
	ListSchemaDesigns(ctx context.Context, in *ListSchemaDesignsRequest, opts ...grpc.CallOption) (*ListSchemaDesignsResponse, error)
	CreateSchemaDesign(ctx context.Context, in *CreateSchemaDesignRequest, opts ...grpc.CallOption) (*SchemaDesign, error)
	UpdateSchemaDesign(ctx context.Context, in *UpdateSchemaDesignRequest, opts ...grpc.CallOption) (*SchemaDesign, error)
	MergeSchemaDesign(ctx context.Context, in *MergeSchemaDesignRequest, opts ...grpc.CallOption) (*SchemaDesign, error)
	ParseSchemaString(ctx context.Context, in *ParseSchemaStringRequest, opts ...grpc.CallOption) (*ParseSchemaStringResponse, error)
	DeleteSchemaDesign(ctx context.Context, in *DeleteSchemaDesignRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DiffMetadata(ctx context.Context, in *ParseSchemaStringRequest, opts ...grpc.CallOption) (*ParseSchemaStringResponse, error)
}

type schemaDesignServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSchemaDesignServiceClient(cc grpc.ClientConnInterface) SchemaDesignServiceClient {
	return &schemaDesignServiceClient{cc}
}

func (c *schemaDesignServiceClient) GetSchemaDesign(ctx context.Context, in *GetSchemaDesignRequest, opts ...grpc.CallOption) (*SchemaDesign, error) {
	out := new(SchemaDesign)
	err := c.cc.Invoke(ctx, SchemaDesignService_GetSchemaDesign_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schemaDesignServiceClient) ListSchemaDesigns(ctx context.Context, in *ListSchemaDesignsRequest, opts ...grpc.CallOption) (*ListSchemaDesignsResponse, error) {
	out := new(ListSchemaDesignsResponse)
	err := c.cc.Invoke(ctx, SchemaDesignService_ListSchemaDesigns_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schemaDesignServiceClient) CreateSchemaDesign(ctx context.Context, in *CreateSchemaDesignRequest, opts ...grpc.CallOption) (*SchemaDesign, error) {
	out := new(SchemaDesign)
	err := c.cc.Invoke(ctx, SchemaDesignService_CreateSchemaDesign_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schemaDesignServiceClient) UpdateSchemaDesign(ctx context.Context, in *UpdateSchemaDesignRequest, opts ...grpc.CallOption) (*SchemaDesign, error) {
	out := new(SchemaDesign)
	err := c.cc.Invoke(ctx, SchemaDesignService_UpdateSchemaDesign_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schemaDesignServiceClient) MergeSchemaDesign(ctx context.Context, in *MergeSchemaDesignRequest, opts ...grpc.CallOption) (*SchemaDesign, error) {
	out := new(SchemaDesign)
	err := c.cc.Invoke(ctx, SchemaDesignService_MergeSchemaDesign_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schemaDesignServiceClient) ParseSchemaString(ctx context.Context, in *ParseSchemaStringRequest, opts ...grpc.CallOption) (*ParseSchemaStringResponse, error) {
	out := new(ParseSchemaStringResponse)
	err := c.cc.Invoke(ctx, SchemaDesignService_ParseSchemaString_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schemaDesignServiceClient) DeleteSchemaDesign(ctx context.Context, in *DeleteSchemaDesignRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, SchemaDesignService_DeleteSchemaDesign_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schemaDesignServiceClient) DiffMetadata(ctx context.Context, in *ParseSchemaStringRequest, opts ...grpc.CallOption) (*ParseSchemaStringResponse, error) {
	out := new(ParseSchemaStringResponse)
	err := c.cc.Invoke(ctx, SchemaDesignService_DiffMetadata_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SchemaDesignServiceServer is the server API for SchemaDesignService service.
// All implementations must embed UnimplementedSchemaDesignServiceServer
// for forward compatibility
type SchemaDesignServiceServer interface {
	GetSchemaDesign(context.Context, *GetSchemaDesignRequest) (*SchemaDesign, error)
	ListSchemaDesigns(context.Context, *ListSchemaDesignsRequest) (*ListSchemaDesignsResponse, error)
	CreateSchemaDesign(context.Context, *CreateSchemaDesignRequest) (*SchemaDesign, error)
	UpdateSchemaDesign(context.Context, *UpdateSchemaDesignRequest) (*SchemaDesign, error)
	MergeSchemaDesign(context.Context, *MergeSchemaDesignRequest) (*SchemaDesign, error)
	ParseSchemaString(context.Context, *ParseSchemaStringRequest) (*ParseSchemaStringResponse, error)
	DeleteSchemaDesign(context.Context, *DeleteSchemaDesignRequest) (*emptypb.Empty, error)
	DiffMetadata(context.Context, *ParseSchemaStringRequest) (*ParseSchemaStringResponse, error)
	mustEmbedUnimplementedSchemaDesignServiceServer()
}

// UnimplementedSchemaDesignServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSchemaDesignServiceServer struct {
}

func (UnimplementedSchemaDesignServiceServer) GetSchemaDesign(context.Context, *GetSchemaDesignRequest) (*SchemaDesign, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchemaDesign not implemented")
}
func (UnimplementedSchemaDesignServiceServer) ListSchemaDesigns(context.Context, *ListSchemaDesignsRequest) (*ListSchemaDesignsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSchemaDesigns not implemented")
}
func (UnimplementedSchemaDesignServiceServer) CreateSchemaDesign(context.Context, *CreateSchemaDesignRequest) (*SchemaDesign, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSchemaDesign not implemented")
}
func (UnimplementedSchemaDesignServiceServer) UpdateSchemaDesign(context.Context, *UpdateSchemaDesignRequest) (*SchemaDesign, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSchemaDesign not implemented")
}
func (UnimplementedSchemaDesignServiceServer) MergeSchemaDesign(context.Context, *MergeSchemaDesignRequest) (*SchemaDesign, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MergeSchemaDesign not implemented")
}
func (UnimplementedSchemaDesignServiceServer) ParseSchemaString(context.Context, *ParseSchemaStringRequest) (*ParseSchemaStringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ParseSchemaString not implemented")
}
func (UnimplementedSchemaDesignServiceServer) DeleteSchemaDesign(context.Context, *DeleteSchemaDesignRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSchemaDesign not implemented")
}
func (UnimplementedSchemaDesignServiceServer) DiffMetadata(context.Context, *ParseSchemaStringRequest) (*ParseSchemaStringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiffMetadata not implemented")
}
func (UnimplementedSchemaDesignServiceServer) mustEmbedUnimplementedSchemaDesignServiceServer() {}

// UnsafeSchemaDesignServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SchemaDesignServiceServer will
// result in compilation errors.
type UnsafeSchemaDesignServiceServer interface {
	mustEmbedUnimplementedSchemaDesignServiceServer()
}

func RegisterSchemaDesignServiceServer(s grpc.ServiceRegistrar, srv SchemaDesignServiceServer) {
	s.RegisterService(&SchemaDesignService_ServiceDesc, srv)
}

func _SchemaDesignService_GetSchemaDesign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSchemaDesignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchemaDesignServiceServer).GetSchemaDesign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchemaDesignService_GetSchemaDesign_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchemaDesignServiceServer).GetSchemaDesign(ctx, req.(*GetSchemaDesignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchemaDesignService_ListSchemaDesigns_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSchemaDesignsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchemaDesignServiceServer).ListSchemaDesigns(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchemaDesignService_ListSchemaDesigns_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchemaDesignServiceServer).ListSchemaDesigns(ctx, req.(*ListSchemaDesignsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchemaDesignService_CreateSchemaDesign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSchemaDesignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchemaDesignServiceServer).CreateSchemaDesign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchemaDesignService_CreateSchemaDesign_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchemaDesignServiceServer).CreateSchemaDesign(ctx, req.(*CreateSchemaDesignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchemaDesignService_UpdateSchemaDesign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSchemaDesignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchemaDesignServiceServer).UpdateSchemaDesign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchemaDesignService_UpdateSchemaDesign_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchemaDesignServiceServer).UpdateSchemaDesign(ctx, req.(*UpdateSchemaDesignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchemaDesignService_MergeSchemaDesign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MergeSchemaDesignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchemaDesignServiceServer).MergeSchemaDesign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchemaDesignService_MergeSchemaDesign_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchemaDesignServiceServer).MergeSchemaDesign(ctx, req.(*MergeSchemaDesignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchemaDesignService_ParseSchemaString_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseSchemaStringRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchemaDesignServiceServer).ParseSchemaString(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchemaDesignService_ParseSchemaString_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchemaDesignServiceServer).ParseSchemaString(ctx, req.(*ParseSchemaStringRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchemaDesignService_DeleteSchemaDesign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSchemaDesignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchemaDesignServiceServer).DeleteSchemaDesign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchemaDesignService_DeleteSchemaDesign_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchemaDesignServiceServer).DeleteSchemaDesign(ctx, req.(*DeleteSchemaDesignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchemaDesignService_DiffMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseSchemaStringRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchemaDesignServiceServer).DiffMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SchemaDesignService_DiffMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchemaDesignServiceServer).DiffMetadata(ctx, req.(*ParseSchemaStringRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SchemaDesignService_ServiceDesc is the grpc.ServiceDesc for SchemaDesignService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SchemaDesignService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bytebase.v1.SchemaDesignService",
	HandlerType: (*SchemaDesignServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSchemaDesign",
			Handler:    _SchemaDesignService_GetSchemaDesign_Handler,
		},
		{
			MethodName: "ListSchemaDesigns",
			Handler:    _SchemaDesignService_ListSchemaDesigns_Handler,
		},
		{
			MethodName: "CreateSchemaDesign",
			Handler:    _SchemaDesignService_CreateSchemaDesign_Handler,
		},
		{
			MethodName: "UpdateSchemaDesign",
			Handler:    _SchemaDesignService_UpdateSchemaDesign_Handler,
		},
		{
			MethodName: "MergeSchemaDesign",
			Handler:    _SchemaDesignService_MergeSchemaDesign_Handler,
		},
		{
			MethodName: "ParseSchemaString",
			Handler:    _SchemaDesignService_ParseSchemaString_Handler,
		},
		{
			MethodName: "DeleteSchemaDesign",
			Handler:    _SchemaDesignService_DeleteSchemaDesign_Handler,
		},
		{
			MethodName: "DiffMetadata",
			Handler:    _SchemaDesignService_DiffMetadata_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/schema_design_service.proto",
}
