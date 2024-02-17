package service

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Manager struct {
	S3Config *aws.Config
}

func NewS3Manager(region, accessKeyId, secretAccessKey string) *S3Manager {
	return &S3Manager{
		S3Config: &aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKeyId,
				secretAccessKey,
				""),
		},
	}
}

func (manager *S3Manager) S3Uploader(bucket string, key string, content []byte) error {
	s3Session, _ := session.NewSession(manager.S3Config)

	uploader := s3manager.NewUploader(s3Session)
	input := &s3manager.UploadInput{
		Bucket:      aws.String(bucket),       // bucket's name
		Key:         aws.String(key),          // files destination location
		Body:        bytes.NewReader(content), // content of the file
		ContentType: aws.String("csv"),        // content type
	}
	_, err := uploader.UploadWithContext(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}
