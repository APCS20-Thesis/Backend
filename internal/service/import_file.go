package service

import (
	"bytes"
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const maxFileSize = 1 << 20

// ...

func (s *Service) ImportFile(ctx context.Context, request *api.ImportFileRequest) (*api.ImportFileResponse, error) {
	s3Config := &aws.Config{
		Region: aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(
			s.config.S3StorageConfig.AccessKeyID,
			s.config.S3StorageConfig.SecretAccessKey,
			""),
	}
	s3Session, _ := session.NewSession(s3Config)

	uploader := s3manager.NewUploader(s3Session)
	input := &s3manager.UploadInput{
		Bucket:      aws.String("cdp-thesis-apcs"),                // bucket's name
		Key:         aws.String("datas/" + request.GetFileName()), // files destination location
		Body:        bytes.NewReader(request.GetFileContent()),    // content of the file
		ContentType: aws.String("csv"),                            // content type
	}
	_, err := uploader.UploadWithContext(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return &api.ImportFileResponse{Message: "Import Success", Code: 0}, nil
}
