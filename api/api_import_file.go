package api

import (
	"context"
	"google.golang.org/grpc"
)

const (
	CDPService_ImportFile_FullMethodName = "/api.CDPServiceFile/ImportFile"
)

var CDPServiceFile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.CDPServiceFile",
	HandlerType: (*CDPServiceFileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ImportFile",
			Handler:    _CDPService_ImportFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func _CDPService_ImportFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceFileServer).ImportFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CDPService_ImportFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceFileServer).ImportFile(ctx, req.(*ImportFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func RegisterCDPServiceFileServer(s grpc.ServiceRegistrar, srv CDPServiceServer) {
	s.RegisterService(&CDPServiceFile_ServiceDesc, srv)
}
