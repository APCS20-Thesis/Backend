package api

import (
	"context"
	"google.golang.org/grpc"
)

const (
	CDPService_ImportCsv_FullMethodName = "/api.CDPServiceFile/ImportCsv"
)

var CDPServiceFile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.CDPServiceFile",
	HandlerType: (*CDPServiceFileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ImportCsv",
			Handler:    _CDPService_ImportCsv_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func _CDPService_ImportCsv_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportCsvRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CDPServiceFileServer).ImportCsv(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CDPService_ImportCsv_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CDPServiceFileServer).ImportCsv(ctx, req.(*ImportCsvRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func RegisterCDPServiceFileServer(s grpc.ServiceRegistrar, srv CDPServiceServer) {
	s.RegisterService(&CDPServiceFile_ServiceDesc, srv)
}
