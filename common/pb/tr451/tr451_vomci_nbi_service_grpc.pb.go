// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: tr451_vomci_nbi_service.proto

package tr451

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// VomciMessageNbiClient is the client API for VomciMessageNbi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VomciMessageNbiClient interface {
	Hello(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error)
	GetData(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error)
	ReplaceConfig(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error)
	UpdateConfigReplica(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error)
	UpdateConfigInstance(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error)
	RPC(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error)
	Action(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error)
	ListenForNotification(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (VomciMessageNbi_ListenForNotificationClient, error)
}

type vomciMessageNbiClient struct {
	cc grpc.ClientConnInterface
}

func NewVomciMessageNbiClient(cc grpc.ClientConnInterface) VomciMessageNbiClient {
	return &vomciMessageNbiClient{cc}
}

func (c *vomciMessageNbiClient) Hello(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error) {
	out := new(Msg)
	err := c.cc.Invoke(ctx, "/tr451_vomci_nbi_service.v1.VomciMessageNbi/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vomciMessageNbiClient) GetData(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error) {
	out := new(Msg)
	err := c.cc.Invoke(ctx, "/tr451_vomci_nbi_service.v1.VomciMessageNbi/GetData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vomciMessageNbiClient) ReplaceConfig(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error) {
	out := new(Msg)
	err := c.cc.Invoke(ctx, "/tr451_vomci_nbi_service.v1.VomciMessageNbi/ReplaceConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vomciMessageNbiClient) UpdateConfigReplica(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error) {
	out := new(Msg)
	err := c.cc.Invoke(ctx, "/tr451_vomci_nbi_service.v1.VomciMessageNbi/UpdateConfigReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vomciMessageNbiClient) UpdateConfigInstance(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error) {
	out := new(Msg)
	err := c.cc.Invoke(ctx, "/tr451_vomci_nbi_service.v1.VomciMessageNbi/UpdateConfigInstance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vomciMessageNbiClient) RPC(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error) {
	out := new(Msg)
	err := c.cc.Invoke(ctx, "/tr451_vomci_nbi_service.v1.VomciMessageNbi/RPC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vomciMessageNbiClient) Action(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Msg, error) {
	out := new(Msg)
	err := c.cc.Invoke(ctx, "/tr451_vomci_nbi_service.v1.VomciMessageNbi/Action", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vomciMessageNbiClient) ListenForNotification(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (VomciMessageNbi_ListenForNotificationClient, error) {
	stream, err := c.cc.NewStream(ctx, &VomciMessageNbi_ServiceDesc.Streams[0], "/tr451_vomci_nbi_service.v1.VomciMessageNbi/ListenForNotification", opts...)
	if err != nil {
		return nil, err
	}
	x := &vomciMessageNbiListenForNotificationClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type VomciMessageNbi_ListenForNotificationClient interface {
	Recv() (*Msg, error)
	grpc.ClientStream
}

type vomciMessageNbiListenForNotificationClient struct {
	grpc.ClientStream
}

func (x *vomciMessageNbiListenForNotificationClient) Recv() (*Msg, error) {
	m := new(Msg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// VomciMessageNbiServer is the server API for VomciMessageNbi service.
// All implementations must embed UnimplementedVomciMessageNbiServer
// for forward compatibility
type VomciMessageNbiServer interface {
	Hello(context.Context, *Msg) (*Msg, error)
	GetData(context.Context, *Msg) (*Msg, error)
	ReplaceConfig(context.Context, *Msg) (*Msg, error)
	UpdateConfigReplica(context.Context, *Msg) (*Msg, error)
	UpdateConfigInstance(context.Context, *Msg) (*Msg, error)
	RPC(context.Context, *Msg) (*Msg, error)
	Action(context.Context, *Msg) (*Msg, error)
	ListenForNotification(*empty.Empty, VomciMessageNbi_ListenForNotificationServer) error
	mustEmbedUnimplementedVomciMessageNbiServer()
}

// UnimplementedVomciMessageNbiServer must be embedded to have forward compatible implementations.
type UnimplementedVomciMessageNbiServer struct {
}

func (UnimplementedVomciMessageNbiServer) Hello(context.Context, *Msg) (*Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedVomciMessageNbiServer) GetData(context.Context, *Msg) (*Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetData not implemented")
}
func (UnimplementedVomciMessageNbiServer) ReplaceConfig(context.Context, *Msg) (*Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplaceConfig not implemented")
}
func (UnimplementedVomciMessageNbiServer) UpdateConfigReplica(context.Context, *Msg) (*Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateConfigReplica not implemented")
}
func (UnimplementedVomciMessageNbiServer) UpdateConfigInstance(context.Context, *Msg) (*Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateConfigInstance not implemented")
}
func (UnimplementedVomciMessageNbiServer) RPC(context.Context, *Msg) (*Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPC not implemented")
}
func (UnimplementedVomciMessageNbiServer) Action(context.Context, *Msg) (*Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Action not implemented")
}
func (UnimplementedVomciMessageNbiServer) ListenForNotification(*empty.Empty, VomciMessageNbi_ListenForNotificationServer) error {
	return status.Errorf(codes.Unimplemented, "method ListenForNotification not implemented")
}
func (UnimplementedVomciMessageNbiServer) mustEmbedUnimplementedVomciMessageNbiServer() {}

// UnsafeVomciMessageNbiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VomciMessageNbiServer will
// result in compilation errors.
type UnsafeVomciMessageNbiServer interface {
	mustEmbedUnimplementedVomciMessageNbiServer()
}

func RegisterVomciMessageNbiServer(s grpc.ServiceRegistrar, srv VomciMessageNbiServer) {
	s.RegisterService(&VomciMessageNbi_ServiceDesc, srv)
}

func _VomciMessageNbi_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VomciMessageNbiServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tr451_vomci_nbi_service.v1.VomciMessageNbi/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VomciMessageNbiServer).Hello(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

func _VomciMessageNbi_GetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VomciMessageNbiServer).GetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tr451_vomci_nbi_service.v1.VomciMessageNbi/GetData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VomciMessageNbiServer).GetData(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

func _VomciMessageNbi_ReplaceConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VomciMessageNbiServer).ReplaceConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tr451_vomci_nbi_service.v1.VomciMessageNbi/ReplaceConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VomciMessageNbiServer).ReplaceConfig(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

func _VomciMessageNbi_UpdateConfigReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VomciMessageNbiServer).UpdateConfigReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tr451_vomci_nbi_service.v1.VomciMessageNbi/UpdateConfigReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VomciMessageNbiServer).UpdateConfigReplica(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

func _VomciMessageNbi_UpdateConfigInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VomciMessageNbiServer).UpdateConfigInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tr451_vomci_nbi_service.v1.VomciMessageNbi/UpdateConfigInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VomciMessageNbiServer).UpdateConfigInstance(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

func _VomciMessageNbi_RPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VomciMessageNbiServer).RPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tr451_vomci_nbi_service.v1.VomciMessageNbi/RPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VomciMessageNbiServer).RPC(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

func _VomciMessageNbi_Action_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VomciMessageNbiServer).Action(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tr451_vomci_nbi_service.v1.VomciMessageNbi/Action",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VomciMessageNbiServer).Action(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

func _VomciMessageNbi_ListenForNotification_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(VomciMessageNbiServer).ListenForNotification(m, &vomciMessageNbiListenForNotificationServer{stream})
}

type VomciMessageNbi_ListenForNotificationServer interface {
	Send(*Msg) error
	grpc.ServerStream
}

type vomciMessageNbiListenForNotificationServer struct {
	grpc.ServerStream
}

func (x *vomciMessageNbiListenForNotificationServer) Send(m *Msg) error {
	return x.ServerStream.SendMsg(m)
}

// VomciMessageNbi_ServiceDesc is the grpc.ServiceDesc for VomciMessageNbi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VomciMessageNbi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tr451_vomci_nbi_service.v1.VomciMessageNbi",
	HandlerType: (*VomciMessageNbiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _VomciMessageNbi_Hello_Handler,
		},
		{
			MethodName: "GetData",
			Handler:    _VomciMessageNbi_GetData_Handler,
		},
		{
			MethodName: "ReplaceConfig",
			Handler:    _VomciMessageNbi_ReplaceConfig_Handler,
		},
		{
			MethodName: "UpdateConfigReplica",
			Handler:    _VomciMessageNbi_UpdateConfigReplica_Handler,
		},
		{
			MethodName: "UpdateConfigInstance",
			Handler:    _VomciMessageNbi_UpdateConfigInstance_Handler,
		},
		{
			MethodName: "RPC",
			Handler:    _VomciMessageNbi_RPC_Handler,
		},
		{
			MethodName: "Action",
			Handler:    _VomciMessageNbi_Action_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListenForNotification",
			Handler:       _VomciMessageNbi_ListenForNotification_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "tr451_vomci_nbi_service.proto",
}
