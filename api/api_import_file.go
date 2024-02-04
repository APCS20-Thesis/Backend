package api

type ImportFileRequest struct {
	FileContent []byte
	FileName    string
	FileSize    int64
	FileType    string
}

func (x *ImportFileRequest) GetFileContent() []byte {
	if x != nil {
		return x.FileContent
	}
	return x.FileContent
}

func (x *ImportFileRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return x.FileName
}

func (x *ImportFileRequest) GetFileSize() int64 {
	if x != nil {
		return x.FileSize
	}
	return x.FileSize
}

func (x *ImportFileRequest) GetFileType() string {
	if x != nil {
		return x.FileType
	}
	return x.FileType
}

type ImportFileResponse struct {
	Message string `json:"message,omitempty"`
	Code    int64  `json:"code,omitempty"`
}
