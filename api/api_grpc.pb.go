// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: api.proto

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

// CDPServiceClient is the client API for CDPService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CDPServiceClient interface {
	CheckHealth(ctx context.Context, in *CheckHealthRequest, opts ...grpc.CallOption) (*CommonResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*CommonResponse, error)
	GetAccountInfo(ctx context.Context, in *GetAccountInfoRequest, opts ...grpc.CallOption) (*GetAccountInfoResponse, error)
	GetListDataSources(ctx context.Context, in *GetListDataSourcesRequest, opts ...grpc.CallOption) (*GetListDataSourcesResponse, error)
	GetDataSource(ctx context.Context, in *GetDataSourceRequest, opts ...grpc.CallOption) (*GetDataSourceResponse, error)
	GetListDataTables(ctx context.Context, in *GetListDataTablesRequest, opts ...grpc.CallOption) (*GetListDataTablesResponse, error)
	GetDataTable(ctx context.Context, in *GetDataTableRequest, opts ...grpc.CallOption) (*GetDataTableResponse, error)
	GetConnection(ctx context.Context, in *GetConnectionRequest, opts ...grpc.CallOption) (*GetConnectionResponse, error)
	GetListConnections(ctx context.Context, in *GetListConnectionsRequest, opts ...grpc.CallOption) (*GetListConnectionsResponse, error)
	CreateConnection(ctx context.Context, in *CreateConnectionRequest, opts ...grpc.CallOption) (*CreateConnectionResponse, error)
	UpdateConnection(ctx context.Context, in *UpdateConnectionRequest, opts ...grpc.CallOption) (*UpdateConnectionResponse, error)
	DeleteConnection(ctx context.Context, in *DeleteConnectionRequest, opts ...grpc.CallOption) (*DeleteConnectionResponse, error)
	ExportDataTableToFile(ctx context.Context, in *ExportDataTableToFileRequest, opts ...grpc.CallOption) (*ExportDataTableToFileResponse, error)
	GetListFileExportRecords(ctx context.Context, in *GetListFileExportRecordsRequest, opts ...grpc.CallOption) (*GetListFileExportRecordsResponse, error)
}

type cDPServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCDPServiceClient(cc grpc.ClientConnInterface) CDPServiceClient {
	return &cDPServiceClient{cc}
}

func (c *cDPServiceClient) CheckHealth(ctx context.Context, in *CheckHealthRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/CheckHealth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/SignUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) GetAccountInfo(ctx context.Context, in *GetAccountInfoRequest, opts ...grpc.CallOption) (*GetAccountInfoResponse, error) {
	out := new(GetAccountInfoResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/GetAccountInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) GetListDataSources(ctx context.Context, in *GetListDataSourcesRequest, opts ...grpc.CallOption) (*GetListDataSourcesResponse, error) {
	out := new(GetListDataSourcesResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/GetListDataSources", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) GetDataSource(ctx context.Context, in *GetDataSourceRequest, opts ...grpc.CallOption) (*GetDataSourceResponse, error) {
	out := new(GetDataSourceResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/GetDataSource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) GetListDataTables(ctx context.Context, in *GetListDataTablesRequest, opts ...grpc.CallOption) (*GetListDataTablesResponse, error) {
	out := new(GetListDataTablesResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/GetListDataTables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) GetDataTable(ctx context.Context, in *GetDataTableRequest, opts ...grpc.CallOption) (*GetDataTableResponse, error) {
	out := new(GetDataTableResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/GetDataTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) GetConnection(ctx context.Context, in *GetConnectionRequest, opts ...grpc.CallOption) (*GetConnectionResponse, error) {
	out := new(GetConnectionResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/GetConnection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) GetListConnections(ctx context.Context, in *GetListConnectionsRequest, opts ...grpc.CallOption) (*GetListConnectionsResponse, error) {
	out := new(GetListConnectionsResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/GetListConnections", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) CreateConnection(ctx context.Context, in *CreateConnectionRequest, opts ...grpc.CallOption) (*CreateConnectionResponse, error) {
	out := new(CreateConnectionResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/CreateConnection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) UpdateConnection(ctx context.Context, in *UpdateConnectionRequest, opts ...grpc.CallOption) (*UpdateConnectionResponse, error) {
	out := new(UpdateConnectionResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/UpdateConnection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) DeleteConnection(ctx context.Context, in *DeleteConnectionRequest, opts ...grpc.CallOption) (*DeleteConnectionResponse, error) {
	out := new(DeleteConnectionResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/DeleteConnection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) ExportDataTableToFile(ctx context.Context, in *ExportDataTableToFileRequest, opts ...grpc.CallOption) (*ExportDataTableToFileResponse, error) {
	out := new(ExportDataTableToFileResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/ExportDataTableToFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cDPServiceClient) GetListFileExportRecords(ctx context.Context, in *GetListFileExportRecordsRequest, opts ...grpc.CallOption) (*GetListFileExportRecordsResponse, error) {
	out := new(GetListFileExportRecordsResponse)
	err := c.cc.Invoke(ctx, "/api.CDPService/GetListFileExportRecords", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CDPServiceServer is the server API for CDPService service.
// All implementations must embed UnimplementedCDPServiceServer
// for forward compatibility
type CDPServiceServer interface {
	CheckHealth(context.Context, *CheckHealthRequest) (*CommonResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	SignUp(context.Context, *SignUpRequest) (*CommonResponse, error)
	GetAccountInfo(context.Context, *GetAccountInfoRequest) (*GetAccountInfoResponse, error)
	GetListDataSources(context.Context, *GetListDataSourcesRequest) (*GetListDataSourcesResponse, error)
	GetDataSource(context.Context, *GetDataSourceRequest) (*GetDataSourceResponse, error)
	GetListDataTables(context.Context, *GetListDataTablesRequest) (*GetListDataTablesResponse, error)
	GetDataTable(context.Context, *GetDataTableRequest) (*GetDataTableResponse, error)
	GetConnection(context.Context, *GetConnectionRequest) (*GetConnectionResponse, error)
	GetListConnections(context.Context, *GetListConnectionsRequest) (*GetListConnectionsResponse, error)
	CreateConnection(context.Context, *CreateConnectionRequest) (*CreateConnectionResponse, error)
	UpdateConnection(context.Context, *UpdateConnectionRequest) (*UpdateConnectionResponse, error)
	DeleteConnection(context.Context, *DeleteConnectionRequest) (*DeleteConnectionResponse, error)
	ExportDataTableToFile(context.Context, *ExportDataTableToFileRequest) (*ExportDataTableToFileResponse, error)
	GetListFileExportRecords(context.Context, *GetListFileExportRecordsRequest) (*GetListFileExportRecordsResponse, error)
	mustEmbedUnimplementedCDPServiceServer()
}

// UnimplementedCDPServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCDPServiceServer struct {
}

func (UnimplementedCDPServiceServer) CheckHealth(context.Context, *CheckHealthRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckHealth not implemented")
}
func (UnimplementedCDPServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedCDPServiceServer) SignUp(context.Context, *SignUpRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedCDPServiceServer) GetAccountInfo(context.Context, *GetAccountInfoRequest) (*GetAccountInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccountInfo not implemented")
}
func (UnimplementedCDPServiceServer) GetListDataSources(context.Context, *GetListDataSourcesRequest) (*GetListDataSourcesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListDataSources not implemented")
}
func (UnimplementedCDPServiceServer) GetDataSource(context.Context, *GetDataSourceRequest) (*GetDataSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDataSource not implemented")
}
func (UnimplementedCDPServiceServer) GetListDataTables(context.Context, *GetListDataTablesRequest) (*GetListDataTablesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListDataTables not implemented")
}
func (UnimplementedCDPServiceServer) GetDataTable(context.Context, *GetDataTableRequest) (*GetDataTableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDataTable not implemented")
}
func (UnimplementedCDPServiceServer) GetConnection(context.Context, *GetConnectionRequest) (*GetConnectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConnection not implemented")
}
func (UnimplementedCDPServiceServer) GetListConnections(context.Context, *GetListConnectionsRequest) (*GetListConnectionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListConnections not implemented")
}
func (UnimplementedCDPServiceServer) CreateConnection(context.Context, *CreateConnectionRequest) (*CreateConnectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateConnection not implemented")
}
func (UnimplementedCDPServiceServer) UpdateConnection(context.Context, *UpdateConnectionRequest) (*UpdateConnectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateConnection not implemented")
}
func (UnimplementedCDPServiceServer) DeleteConnection(context.Context, *DeleteConnectionRequest) (*DeleteConnectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteConnection not implemented")
}
func (UnimplementedCDPServiceServer) ExportDataTableToFile(context.Context, *ExportDataTableToFileRequest) (*ExportDataTableToFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportDataTableToFile not implemented")
}
func (UnimplementedCDPServiceServer) GetListFileExportRecords(context.Context, *GetListFileExportRecordsRequest) (*GetListFileExportRecordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListFileExportRecords not implemented")
}
func (UnimplementedCDPServiceServer) mustEmbedUnimplementedCDPServiceServer() {}

// UnsafeCDPServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CDPServiceServer will
// result in compilation errors.
type UnsafeCDPServiceServer interface {
	mustEmbedUnimplementedCDPServiceServer()
}

func RegisterCDPServiceServer(s grpc.ServiceRegistrar, srv CDPServiceServer) {
	s.RegisterService(&CDPService_ServiceDesc, srv)
}

func _CDPService_CheckHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckHealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).CheckHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/CheckHealth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).CheckHealth(ctx, req.(*CheckHealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/SignUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).SignUp(ctx, req.(*SignUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_GetAccountInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).GetAccountInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/GetAccountInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).GetAccountInfo(ctx, req.(*GetAccountInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_GetListDataSources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListDataSourcesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).GetListDataSources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/GetListDataSources",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).GetListDataSources(ctx, req.(*GetListDataSourcesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_GetDataSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDataSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).GetDataSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/GetDataSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).GetDataSource(ctx, req.(*GetDataSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_GetListDataTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListDataTablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).GetListDataTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/GetListDataTables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).GetListDataTables(ctx, req.(*GetListDataTablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_GetDataTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDataTableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).GetDataTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/GetDataTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).GetDataTable(ctx, req.(*GetDataTableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_GetConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).GetConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/GetConnection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).GetConnection(ctx, req.(*GetConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_GetListConnections_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListConnectionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).GetListConnections(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/GetListConnections",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).GetListConnections(ctx, req.(*GetListConnectionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_CreateConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).CreateConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/CreateConnection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).CreateConnection(ctx, req.(*CreateConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_UpdateConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).UpdateConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/UpdateConnection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).UpdateConnection(ctx, req.(*UpdateConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_DeleteConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).DeleteConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/DeleteConnection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).DeleteConnection(ctx, req.(*DeleteConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_ExportDataTableToFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportDataTableToFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).ExportDataTableToFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/ExportDataTableToFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).ExportDataTableToFile(ctx, req.(*ExportDataTableToFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CDPService_GetListFileExportRecords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListFileExportRecordsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceServer).GetListFileExportRecords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CDPService/GetListFileExportRecords",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceServer).GetListFileExportRecords(ctx, req.(*GetListFileExportRecordsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CDPService_ServiceDesc is the grpc.ServiceDesc for CDPService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CDPService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.CDPService",
	HandlerType: (*CDPServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckHealth",
			Handler:    _CDPService_CheckHealth_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _CDPService_Login_Handler,
		},
		{
			MethodName: "SignUp",
			Handler:    _CDPService_SignUp_Handler,
		},
		{
			MethodName: "GetAccountInfo",
			Handler:    _CDPService_GetAccountInfo_Handler,
		},
		{
			MethodName: "GetListDataSources",
			Handler:    _CDPService_GetListDataSources_Handler,
		},
		{
			MethodName: "GetDataSource",
			Handler:    _CDPService_GetDataSource_Handler,
		},
		{
			MethodName: "GetListDataTables",
			Handler:    _CDPService_GetListDataTables_Handler,
		},
		{
			MethodName: "GetDataTable",
			Handler:    _CDPService_GetDataTable_Handler,
		},
		{
			MethodName: "GetConnection",
			Handler:    _CDPService_GetConnection_Handler,
		},
		{
			MethodName: "GetListConnections",
			Handler:    _CDPService_GetListConnections_Handler,
		},
		{
			MethodName: "CreateConnection",
			Handler:    _CDPService_CreateConnection_Handler,
		},
		{
			MethodName: "UpdateConnection",
			Handler:    _CDPService_UpdateConnection_Handler,
		},
		{
			MethodName: "DeleteConnection",
			Handler:    _CDPService_DeleteConnection_Handler,
		},
		{
			MethodName: "ExportDataTableToFile",
			Handler:    _CDPService_ExportDataTableToFile_Handler,
		},
		{
			MethodName: "GetListFileExportRecords",
			Handler:    _CDPService_GetListFileExportRecords_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
