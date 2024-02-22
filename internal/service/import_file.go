package service

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/constants"
	"time"
)

func (s *Service) ImportFile(ctx context.Context, request *api.ImportFileRequest) (*api.ImportFileResponse, error) {
	accountUuid, err := GetAccountUuidFromCtx(ctx)
	if err != nil {
		s.log.WithName("ImportFile").
			WithValues("Context", ctx).
			Error(err, "Cannot get account_uuid from context")
		return nil, err
	}
	dateTime := GetDateTimeString()
	key := "data/" + accountUuid + "/" + dateTime + "_" + request.GetFileName()
	err = s.s3Manger.S3Uploader(
		constants.S3BucketName,
		key,
		request.GetFileContent())
	if err != nil {
		s.log.WithName("ImportFile").
			WithValues("Context", ctx).
			Error(err, "Cannot uploaded file to S3")
		return nil, err
	}
	filePath := s.config.S3StorageConfig.Host + "/" + constants.S3BucketName + "/" + key
	err = s.business.DataSourceBusiness.ProcessImportFile(ctx, request, accountUuid, dateTime, filePath)
	if err != nil {
		s.log.WithName("ImportFile").
			WithValues("Context", ctx).
			Error(err, "Failed to process import file")
		return nil, err
	}
	return &api.ImportFileResponse{Message: "Import Success", Code: 0}, nil
}

func GetDateTimeString() string {
	var currentTime time.Time
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		currentTime = time.Now()
	}
	currentTime = time.Now().In(location)

	return currentTime.Format("02012006150405")
}
