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
	"strconv"
)

type CDPServiceFileServer interface {
	ImportCsv(context.Context, *ImportCsvRequest) (*ImportCsvResponse, error)
}
type UnimplementedCDPServiceFile struct {
}

func (UnimplementedCDPServiceFile) ImportCsv(context.Context, *ImportCsvRequest) (*ImportCsvResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportCsv not implemented")
}
func request_CDPServiceFile_ImportCsv_0(ctx context.Context, client CDPServiceFileClient, req *http.Request, pathParams map[string]string) (*ImportCsvResponse, error) {
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
	var mappingOptions []*MappingOptionItem
	err = json.Unmarshal([]byte(jsonMappingOptions), &mappingOptions)
	if err != nil {
		return nil, err
	}

	jsonConfigurations := req.Form.Get("configurations")
	var configurations *ImportCsvConfigurations
	err = json.Unmarshal([]byte(jsonConfigurations), &configurations)
	if err != nil {
		return nil, err
	}

	form := req.Form
	name := form.Get("name")
	description := form.Get("description")
	newTableName := form.Get("new_table_name")
	writeMode := form.Get("write_mode")

	var tableId int64
	if len(form.Get("table_id")) > 0 {
		tableId, err = strconv.ParseInt(form.Get("table_id"), 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		tableId = 0
	}

	var metadata runtime.ServerMetadata
	response, err := client.ImportCsv(ctx,
		&ImportCsvRequest{
			FileContent:    fileBytes,
			FileName:       header.Filename,
			FileSize:       header.Size,
			MappingOptions: mappingOptions,
			Configurations: configurations,
			Name:           name,
			Description:    description,
			TableId:        tableId,
			NewTableName:   newTableName,
			WriteMode:      writeMode,
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
	ImportCsv(ctx context.Context, in *ImportCsvRequest, opts ...grpc.CallOption) (*ImportCsvResponse, error)
}

type cDPServiceFileClient struct {
	cc grpc.ClientConnInterface
}

func NewCDPServiceFileClient(cc grpc.ClientConnInterface) CDPServiceFileClient {
	return &cDPServiceFileClient{cc}
}

func (c *cDPServiceFileClient) ImportCsv(ctx context.Context, in *ImportCsvRequest, opts ...grpc.CallOption) (*ImportCsvResponse, error) {
	out := new(ImportCsvResponse)
	err := c.cc.Invoke(ctx, CDPService_ImportCsv_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func RegisterCDPServiceFileClient(ctx context.Context, mux *runtime.ServeMux, client CDPServiceFileClient) error {

	mux.Handle("POST", pattern_CDPServiceFile_ImportCsv_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
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
		resp, err := request_CDPServiceFile_ImportCsv_0(rctx, client, req, pathParams)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_CDPServiceFile_ImportCsv_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	return nil
}

var (
	pattern_CDPServiceFile_ImportCsv_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "v1", "data-source", "import-csv"}, "", runtime.AssumeColonVerbOpt(true)))
)
var (
	forward_CDPServiceFile_ImportCsv_0 = runtime.ForwardResponseMessage
)
