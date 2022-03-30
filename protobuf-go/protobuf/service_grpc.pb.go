// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: service.proto

package protobuf

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

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	GetDBTokenAuthInfo(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*DBTokenResponse, error)
	RenewToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	CreateSession(ctx context.Context, in *SessionCreateRequest, opts ...grpc.CallOption) (*SessionCreateResponse, error)
	FinishSession(ctx context.Context, in *SessionFinishRequest, opts ...grpc.CallOption) (*SessionFinishResp, error)
	UploadReplayFile(ctx context.Context, in *ReplayRequest, opts ...grpc.CallOption) (*ReplayResponse, error)
	UploadCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error)
	DispatchTask(ctx context.Context, opts ...grpc.CallOption) (Service_DispatchTaskClient, error)
	ScanRemainReplays(ctx context.Context, in *RemainReplayRequest, opts ...grpc.CallOption) (*RemainReplayResponse, error)
	CreateCommandTicket(ctx context.Context, in *CommandConfirmRequest, opts ...grpc.CallOption) (*CommandConfirmResponse, error)
	CheckOrCreateAssetLoginTicket(ctx context.Context, in *AssetLoginTicketRequest, opts ...grpc.CallOption) (*AssetLoginTicketResponse, error)
	CheckTicketState(ctx context.Context, in *TicketRequest, opts ...grpc.CallOption) (*TicketStateResponse, error)
	CancelTicket(ctx context.Context, in *TicketRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	CreateForward(ctx context.Context, in *ForwardRequest, opts ...grpc.CallOption) (*ForwardResponse, error)
	DeleteForward(ctx context.Context, in *ForwardDeleteRequest, opts ...grpc.CallOption) (*StatusResponse, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) GetDBTokenAuthInfo(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*DBTokenResponse, error) {
	out := new(DBTokenResponse)
	err := c.cc.Invoke(ctx, "/message.Service/GetDBTokenAuthInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) RenewToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/message.Service/RenewToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) CreateSession(ctx context.Context, in *SessionCreateRequest, opts ...grpc.CallOption) (*SessionCreateResponse, error) {
	out := new(SessionCreateResponse)
	err := c.cc.Invoke(ctx, "/message.Service/CreateSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) FinishSession(ctx context.Context, in *SessionFinishRequest, opts ...grpc.CallOption) (*SessionFinishResp, error) {
	out := new(SessionFinishResp)
	err := c.cc.Invoke(ctx, "/message.Service/FinishSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) UploadReplayFile(ctx context.Context, in *ReplayRequest, opts ...grpc.CallOption) (*ReplayResponse, error) {
	out := new(ReplayResponse)
	err := c.cc.Invoke(ctx, "/message.Service/UploadReplayFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) UploadCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error) {
	out := new(CommandResponse)
	err := c.cc.Invoke(ctx, "/message.Service/UploadCommand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DispatchTask(ctx context.Context, opts ...grpc.CallOption) (Service_DispatchTaskClient, error) {
	stream, err := c.cc.NewStream(ctx, &Service_ServiceDesc.Streams[0], "/message.Service/DispatchTask", opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceDispatchTaskClient{stream}
	return x, nil
}

type Service_DispatchTaskClient interface {
	Send(*FinishedTaskRequest) error
	Recv() (*TaskResponse, error)
	grpc.ClientStream
}

type serviceDispatchTaskClient struct {
	grpc.ClientStream
}

func (x *serviceDispatchTaskClient) Send(m *FinishedTaskRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serviceDispatchTaskClient) Recv() (*TaskResponse, error) {
	m := new(TaskResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *serviceClient) ScanRemainReplays(ctx context.Context, in *RemainReplayRequest, opts ...grpc.CallOption) (*RemainReplayResponse, error) {
	out := new(RemainReplayResponse)
	err := c.cc.Invoke(ctx, "/message.Service/ScanRemainReplays", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) CreateCommandTicket(ctx context.Context, in *CommandConfirmRequest, opts ...grpc.CallOption) (*CommandConfirmResponse, error) {
	out := new(CommandConfirmResponse)
	err := c.cc.Invoke(ctx, "/message.Service/CreateCommandTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) CheckOrCreateAssetLoginTicket(ctx context.Context, in *AssetLoginTicketRequest, opts ...grpc.CallOption) (*AssetLoginTicketResponse, error) {
	out := new(AssetLoginTicketResponse)
	err := c.cc.Invoke(ctx, "/message.Service/CheckOrCreateAssetLoginTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) CheckTicketState(ctx context.Context, in *TicketRequest, opts ...grpc.CallOption) (*TicketStateResponse, error) {
	out := new(TicketStateResponse)
	err := c.cc.Invoke(ctx, "/message.Service/CheckTicketState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) CancelTicket(ctx context.Context, in *TicketRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/message.Service/CancelTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) CreateForward(ctx context.Context, in *ForwardRequest, opts ...grpc.CallOption) (*ForwardResponse, error) {
	out := new(ForwardResponse)
	err := c.cc.Invoke(ctx, "/message.Service/CreateForward", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DeleteForward(ctx context.Context, in *ForwardDeleteRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/message.Service/DeleteForward", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	GetDBTokenAuthInfo(context.Context, *TokenRequest) (*DBTokenResponse, error)
	RenewToken(context.Context, *TokenRequest) (*StatusResponse, error)
	CreateSession(context.Context, *SessionCreateRequest) (*SessionCreateResponse, error)
	FinishSession(context.Context, *SessionFinishRequest) (*SessionFinishResp, error)
	UploadReplayFile(context.Context, *ReplayRequest) (*ReplayResponse, error)
	UploadCommand(context.Context, *CommandRequest) (*CommandResponse, error)
	DispatchTask(Service_DispatchTaskServer) error
	ScanRemainReplays(context.Context, *RemainReplayRequest) (*RemainReplayResponse, error)
	CreateCommandTicket(context.Context, *CommandConfirmRequest) (*CommandConfirmResponse, error)
	CheckOrCreateAssetLoginTicket(context.Context, *AssetLoginTicketRequest) (*AssetLoginTicketResponse, error)
	CheckTicketState(context.Context, *TicketRequest) (*TicketStateResponse, error)
	CancelTicket(context.Context, *TicketRequest) (*StatusResponse, error)
	CreateForward(context.Context, *ForwardRequest) (*ForwardResponse, error)
	DeleteForward(context.Context, *ForwardDeleteRequest) (*StatusResponse, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) GetDBTokenAuthInfo(context.Context, *TokenRequest) (*DBTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDBTokenAuthInfo not implemented")
}
func (UnimplementedServiceServer) RenewToken(context.Context, *TokenRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenewToken not implemented")
}
func (UnimplementedServiceServer) CreateSession(context.Context, *SessionCreateRequest) (*SessionCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}
func (UnimplementedServiceServer) FinishSession(context.Context, *SessionFinishRequest) (*SessionFinishResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FinishSession not implemented")
}
func (UnimplementedServiceServer) UploadReplayFile(context.Context, *ReplayRequest) (*ReplayResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadReplayFile not implemented")
}
func (UnimplementedServiceServer) UploadCommand(context.Context, *CommandRequest) (*CommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadCommand not implemented")
}
func (UnimplementedServiceServer) DispatchTask(Service_DispatchTaskServer) error {
	return status.Errorf(codes.Unimplemented, "method DispatchTask not implemented")
}
func (UnimplementedServiceServer) ScanRemainReplays(context.Context, *RemainReplayRequest) (*RemainReplayResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScanRemainReplays not implemented")
}
func (UnimplementedServiceServer) CreateCommandTicket(context.Context, *CommandConfirmRequest) (*CommandConfirmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCommandTicket not implemented")
}
func (UnimplementedServiceServer) CheckOrCreateAssetLoginTicket(context.Context, *AssetLoginTicketRequest) (*AssetLoginTicketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckOrCreateAssetLoginTicket not implemented")
}
func (UnimplementedServiceServer) CheckTicketState(context.Context, *TicketRequest) (*TicketStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckTicketState not implemented")
}
func (UnimplementedServiceServer) CancelTicket(context.Context, *TicketRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelTicket not implemented")
}
func (UnimplementedServiceServer) CreateForward(context.Context, *ForwardRequest) (*ForwardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateForward not implemented")
}
func (UnimplementedServiceServer) DeleteForward(context.Context, *ForwardDeleteRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteForward not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_GetDBTokenAuthInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetDBTokenAuthInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/GetDBTokenAuthInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetDBTokenAuthInfo(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_RenewToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).RenewToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/RenewToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).RenewToken(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_CreateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CreateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/CreateSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CreateSession(ctx, req.(*SessionCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_FinishSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionFinishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).FinishSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/FinishSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).FinishSession(ctx, req.(*SessionFinishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_UploadReplayFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).UploadReplayFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/UploadReplayFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).UploadReplayFile(ctx, req.(*ReplayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_UploadCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).UploadCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/UploadCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).UploadCommand(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DispatchTask_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServiceServer).DispatchTask(&serviceDispatchTaskServer{stream})
}

type Service_DispatchTaskServer interface {
	Send(*TaskResponse) error
	Recv() (*FinishedTaskRequest, error)
	grpc.ServerStream
}

type serviceDispatchTaskServer struct {
	grpc.ServerStream
}

func (x *serviceDispatchTaskServer) Send(m *TaskResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serviceDispatchTaskServer) Recv() (*FinishedTaskRequest, error) {
	m := new(FinishedTaskRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Service_ScanRemainReplays_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemainReplayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ScanRemainReplays(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/ScanRemainReplays",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ScanRemainReplays(ctx, req.(*RemainReplayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_CreateCommandTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandConfirmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CreateCommandTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/CreateCommandTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CreateCommandTicket(ctx, req.(*CommandConfirmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_CheckOrCreateAssetLoginTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssetLoginTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CheckOrCreateAssetLoginTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/CheckOrCreateAssetLoginTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CheckOrCreateAssetLoginTicket(ctx, req.(*AssetLoginTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_CheckTicketState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CheckTicketState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/CheckTicketState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CheckTicketState(ctx, req.(*TicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_CancelTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CancelTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/CancelTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CancelTicket(ctx, req.(*TicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_CreateForward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CreateForward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/CreateForward",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CreateForward(ctx, req.(*ForwardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DeleteForward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DeleteForward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Service/DeleteForward",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DeleteForward(ctx, req.(*ForwardDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "message.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDBTokenAuthInfo",
			Handler:    _Service_GetDBTokenAuthInfo_Handler,
		},
		{
			MethodName: "RenewToken",
			Handler:    _Service_RenewToken_Handler,
		},
		{
			MethodName: "CreateSession",
			Handler:    _Service_CreateSession_Handler,
		},
		{
			MethodName: "FinishSession",
			Handler:    _Service_FinishSession_Handler,
		},
		{
			MethodName: "UploadReplayFile",
			Handler:    _Service_UploadReplayFile_Handler,
		},
		{
			MethodName: "UploadCommand",
			Handler:    _Service_UploadCommand_Handler,
		},
		{
			MethodName: "ScanRemainReplays",
			Handler:    _Service_ScanRemainReplays_Handler,
		},
		{
			MethodName: "CreateCommandTicket",
			Handler:    _Service_CreateCommandTicket_Handler,
		},
		{
			MethodName: "CheckOrCreateAssetLoginTicket",
			Handler:    _Service_CheckOrCreateAssetLoginTicket_Handler,
		},
		{
			MethodName: "CheckTicketState",
			Handler:    _Service_CheckTicketState_Handler,
		},
		{
			MethodName: "CancelTicket",
			Handler:    _Service_CancelTicket_Handler,
		},
		{
			MethodName: "CreateForward",
			Handler:    _Service_CreateForward_Handler,
		},
		{
			MethodName: "DeleteForward",
			Handler:    _Service_DeleteForward_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DispatchTask",
			Handler:       _Service_DispatchTask_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "service.proto",
}
