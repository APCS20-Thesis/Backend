package api

import (
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

const (
	CDPService_ImportFile_FullMethodName = "/api.CDPServiceImportFile/ImportFile"
)

type CDPServiceFile interface {
	ImportFile(context.Context, *ImportFileRequest) (*ImportFileResponse, error)
}
type UnimplementedCDPServiceFile struct {
}

func (UnimplementedCDPServiceFile) ImportFile(context.Context, *ImportFileRequest) (*ImportFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckHealth not implemented")
}
func request_CDPServiceFile_ImportFile_0(ctx context.Context, service CDPServiceFile, req *http.Request, pathParams map[string]string) (*ImportFileResponse, error) {
	content, header, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer content.Close()
	fileBytes, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}
	form := req.Form
	fileType := form.Get("file_type")
	response, err := service.ImportFile(ctx,
		&ImportFileRequest{
			FileContent: fileBytes,
			FileName:    header.Filename,
			FileSize:    header.Size,
			FileType:    fileType,
		},
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func RegisterCDPServiceImportFile(ctx context.Context, mux *runtime.ServeMux, service CDPServiceFile) error {
	mux.Handle("POST", pattern_CDPServiceFile_ImportFile_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		err := req.ParseMultipartForm(32 << 20)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		resp, err := request_CDPServiceFile_ImportFile_0(ctx, service, req, pathParams)
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
	})
	return nil
}

var (
	pattern_CDPServiceFile_ImportFile_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "v1", "data-source", "import-file"}, "", runtime.AssumeColonVerbOpt(true)))
)
