package api

import (
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

type CDPServiceFileServer interface {
	ImportFile(context.Context, *ImportFileRequest) (*ImportFileResponse, error)
}
type UnimplementedCDPServiceFile struct {
}

func (UnimplementedCDPServiceFile) ImportFile(context.Context, *ImportFileRequest) (*ImportFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckHealth not implemented")
}
func request_CDPServiceFile_ImportFile_0(ctx context.Context, client CDPServiceFileClient, req *http.Request, pathParams map[string]string) (*ImportFileResponse, error) {
	content, header, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer content.Close()
	fileBytes, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}
	jsonMappingOptions := req.Form.Get("mapping_options")
	var mappingOptions map[string]string
	err = json.Unmarshal([]byte(jsonMappingOptions), &mappingOptions)
	if err != nil {
		return nil, err
	}
	form := req.Form
	fileType := form.Get("file_type")
	name := form.Get("name")
	description := form.Get("description")
	deltaTableName := form.Get("delta_table_name")
	var metadata runtime.ServerMetadata
	response, err := client.ImportFile(ctx,
		&ImportFileRequest{
			FileContent:    fileBytes,
			FileName:       header.Filename,
			FileSize:       header.Size,
			FileType:       fileType,
			MappingOptions: mappingOptions,
			Name:           name,
			Description:    description,
			DeltaTableName: deltaTableName,
		},
		grpc.Header(&metadata.HeaderMD),
		grpc.Trailer(&metadata.TrailerMD),
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type CDPServiceFileClient interface {
	ImportFile(ctx context.Context, in *ImportFileRequest, opts ...grpc.CallOption) (*ImportFileResponse, error)
}

type cDPServiceFileClient struct {
	cc grpc.ClientConnInterface
}

func NewCDPServiceFileClient(cc grpc.ClientConnInterface) CDPServiceFileClient {
	return &cDPServiceFileClient{cc}
}

func (c *cDPServiceFileClient) ImportFile(ctx context.Context, in *ImportFileRequest, opts ...grpc.CallOption) (*ImportFileResponse, error) {
	out := new(ImportFileResponse)
	err := c.cc.Invoke(ctx, CDPService_ImportFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func RegisterCDPServiceFileClient(ctx context.Context, mux *runtime.ServeMux, client CDPServiceFileClient) error {

	mux.Handle("POST", pattern_CDPServiceFile_ImportFile_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		err = req.ParseMultipartForm(32 << 20)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, err := request_CDPServiceFile_ImportFile_0(rctx, client, req, pathParams)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_CDPServiceFile_ImportFile_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	return nil
}

var (
	pattern_CDPServiceFile_ImportFile_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "v1", "data-source", "import-file"}, "", runtime.AssumeColonVerbOpt(true)))
)
var (
	forward_CDPServiceFile_ImportFile_0 = runtime.ForwardResponseMessage
)
