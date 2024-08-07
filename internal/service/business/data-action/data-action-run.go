package data_action

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

func (b business) ProcessNewDataActionRun(ctx context.Context, request *api.TriggerDataActionRunRequest, accountUuid string) error {
	logger := b.log.WithName("ProcessNewDataActionRun").WithValues("request", request)
	dataAction, err := b.repository.DataActionRepository.GetDataAction(ctx, request.Id)
	if err != nil {
		logger.Error(err, "cannot get data action")
		return err
	}
	if dataAction.AccountUuid.String() != accountUuid {
		return status.Error(codes.PermissionDenied, "user not have access to this data action")
	}

	err = b.db.Transaction(func(tx *gorm.DB) error {
		// trigger airflow
		dagRun, txErr := b.airflowAdapter.TriggerNewDagRun(ctx, dataAction.DagId, &airflow.TriggerNewDagRunRequest{})
		if txErr != nil {
			logger.Error(err, "cannot trigger new dag run")
			return txErr
		}

		// create new data action run
		_, txErr = b.repository.DataActionRunRepository.CreateDataActionRun(ctx, &repository.CreateDataActionRunParams{
			Tx:          tx,
			ActionId:    dataAction.ID,
			RunId:       dataAction.RunCount + 1,
			DagRunId:    dagRun.DagRunId,
			Status:      model.DataActionRunStatus_Processing,
			AccountUuid: uuid.MustParse(accountUuid),
		})
		if txErr != nil {
			logger.Error(err, "cannot create new data action run")
			return txErr
		}

		// increase data action run count by 1
		txErr = b.repository.DataActionRepository.UpdateDataAction(ctx, &repository.UpdateDataActionParams{
			Tx:       tx,
			ID:       dataAction.ID,
			RunCount: dataAction.RunCount + 1,
		})
		if txErr != nil {
			logger.Error(err, "cannot update data action run count")
			return txErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (b business) ProcessGetListDataActionRuns(ctx context.Context, request *api.GetListDataActionRunsRequest, accountUuid string) (*api.GetListDataActionRunsResponse, error) {
	queryDataActionRuns, err := b.repository.DataActionRunRepository.GetListDataActionRuns(ctx, &repository.GetListDataActionRunsParams{
		ActionTypes: request.Types,
		AccountUuid: uuid.MustParse(accountUuid),
		Page:        int(request.Page),
		PageSize:    int(request.PageSize),
	})
	if err != nil {
		b.log.WithName("ProcessGetListDataActionRuns").Error(err, "cannot get data action runs")
		return nil, err
	}
	var returnDataActionRuns []*api.DataActionRun
	for _, dataActionRun := range queryDataActionRuns.DataActionRuns {
		metadata, err := b.EnrichDataActionRunMetaData(ctx, dataActionRun)
		if err != nil {
			b.log.WithName("ProcessGetListDataActionRuns").Error(err, "cannot enrich data action run")
			return nil, err
		}
		returnDataActionRuns = append(returnDataActionRuns, &api.DataActionRun{
			Id:         dataActionRun.ID,
			ActionId:   dataActionRun.ActionId,
			ActionType: string(dataActionRun.ActionType),
			Status:     string(dataActionRun.Status),
			CreatedAt:  dataActionRun.CreatedAt.String(),
			UpdatedAt:  dataActionRun.UpdatedAt.String(),
			Metadata:   metadata,
		})
	}

	return &api.GetListDataActionRunsResponse{
		Code:    0,
		Message: "Success",
		Count:   queryDataActionRuns.Count,
		Results: returnDataActionRuns,
	}, nil
}

func (b business) EnrichDataActionRunMetaData(ctx context.Context, dataActionRun repository.DataActionRunWithExtraInfo) (*api.DataActionRun_MetaData, error) {
	switch dataActionRun.ActionType {
	case model.ActionType_ImportDataFromMySQL, model.ActionType_ImportDataFromFile, model.ActionType_ImportDataFromS3:
		sourceMapQuery, err := b.repository.SourceTableMapRepository.ListSourceTableMap(ctx, &repository.ListSourceTableMapParams{Ids: []int64{dataActionRun.ObjectId}})
		if err != nil {
			return nil, err
		}
		if len(sourceMapQuery.TableSourceMaps) == 0 {
			return &api.DataActionRun_MetaData{}, nil
		}
		return &api.DataActionRun_MetaData{ObjectReference: &api.DataActionRun_ObjectReference{
			Type: string(model.TargetTable_DataSource),
			Name: sourceMapQuery.TableSourceMaps[0].SourceName,
			Id:   sourceMapQuery.TableSourceMaps[0].SourceId,
		}}, nil
	case model.ActionType_ExportDataToS3CSV, model.ActionType_ExportGophish, model.ActionType_ExportToMySQL:
		if dataActionRun.TargetTable == string(model.TargetTable_DestTableMap) {
			dstMap, err := b.repository.DestTableMapRepository.ListDestinationTableMaps(ctx, &repository.ListDestinationTableMapsParams{
				Ids: []int64{dataActionRun.ObjectId},
			})
			if err != nil {
				return nil, err
			}
			if len(dstMap) == 0 {
				return &api.DataActionRun_MetaData{}, nil
			}
			return &api.DataActionRun_MetaData{ObjectReference: &api.DataActionRun_ObjectReference{
				Type: string(model.TargetTable_DataDestination),
				Name: dstMap[0].DestinationName,
				Id:   dstMap[0].DestinationId,
			}}, nil
		} else if dataActionRun.TargetTable == string(model.TargetTable_DestSegmentMap) {
			dstMap, err := b.repository.DestSegmentMapRepository.ListDestinationSegmentMaps(ctx, &repository.ListDestinationSegmentMapsParams{
				Ids: []int64{dataActionRun.ObjectId},
			})
			if err != nil {
				return nil, err
			}
			if len(dstMap) == 0 {
				return &api.DataActionRun_MetaData{}, nil
			}
			return &api.DataActionRun_MetaData{ObjectReference: &api.DataActionRun_ObjectReference{
				Type: string(model.TargetTable_DataDestination),
				Name: dstMap[0].DestinationName,
				Id:   dstMap[0].DestinationId,
			}}, nil
		} else if dataActionRun.TargetTable == string(model.TargetTable_DestMasterSegmentMap) {
			dstMap, err := b.repository.DestMasterSegmentMapRepository.ListDestinationMasterSegmentMaps(ctx, &repository.ListDestinationMasterSegmentMapsParams{
				Ids: []int64{dataActionRun.ObjectId},
			})
			if err != nil {
				return nil, err
			}
			if len(dstMap) == 0 {
				return &api.DataActionRun_MetaData{}, nil
			}
			return &api.DataActionRun_MetaData{ObjectReference: &api.DataActionRun_ObjectReference{
				Type: string(model.TargetTable_DataDestination),
				Name: dstMap[0].DestinationName,
				Id:   dstMap[0].DestinationId,
			}}, nil
		}
	case model.ActionType_CreateMasterSegment:
		masterSegment, err := b.repository.SegmentRepository.GetMasterSegment(ctx, dataActionRun.ObjectId)
		if err != nil {
			return nil, err
		}
		return &api.DataActionRun_MetaData{ObjectReference: &api.DataActionRun_ObjectReference{
			Type: string(model.TargetTable_MasterSegment),
			Name: masterSegment.Name,
			Id:   masterSegment.ID,
		}}, nil
	case model.ActionType_CreateSegment, model.ActionType_ApplyPredictModel:
		segment, err := b.repository.SegmentRepository.GetSegment(ctx, dataActionRun.ObjectId)
		if err != nil {
			return nil, err
		}
		return &api.DataActionRun_MetaData{
			ObjectReference: &api.DataActionRun_ObjectReference{
				Type: string(model.TargetTable_Segment),
				Name: segment.Name,
				Id:   segment.ID,
			},
			MasterSegmentId: segment.MasterSegmentId,
		}, nil
	case model.ActionType_TrainPredictModel:
		predictModel, err := b.repository.PredictModelRepository.GetPredictModel(ctx, dataActionRun.ObjectId)
		if err != nil {
			return nil, err
		}
		return &api.DataActionRun_MetaData{
			ObjectReference: &api.DataActionRun_ObjectReference{
				Type: string(model.TargetTable_PredictModel),
				Name: predictModel.Name,
				Id:   predictModel.ID,
			},
			MasterSegmentId: predictModel.MasterSegmentId,
		}, nil
	}

	return &api.DataActionRun_MetaData{}, nil
}

func (b business) ProcessGetTotalRunsPerDay(ctx context.Context, accountUuid string) (*api.GetDataActionRunsPerDayResponse, error) {
	const days = 15

	results, err := b.repository.DataActionRunRepository.GetTotalRunsPerDay(ctx, &repository.GetTotalRunsPerDayParams{AccountUuid: accountUuid})
	if err != nil {
		b.log.WithName("ProcessGetTotalRunsPerDay").Error(err, "cannot get data action runs per day", "uuid", accountUuid)
		return nil, err
	}

	// Generate dates for the last 15 days
	dates := make([]time.Time, 0, days)
	currentTime := time.Now()
	currentDate := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	for i := days - 1; i >= 0; i-- {
		dates = append(dates, currentDate.AddDate(0, 0, -i))
	}

	// Check and add missing dates to the data array
	fulfilledResult := make([]repository.TotalRunsPerDay, 0, days)
	arrIdx := 0
	for _, date := range dates {
		if arrIdx < len(results) && utils.CompareDay(date, results[arrIdx].Date) {
			fulfilledResult = append(fulfilledResult, results[arrIdx])
			arrIdx++
		} else {
			fulfilledResult = append(fulfilledResult, repository.TotalRunsPerDay{
				Date:  date,
				Total: 0,
			})
		}
	}

	runsPerDay := utils.Map(fulfilledResult, func(data repository.TotalRunsPerDay) *api.GetDataActionRunsPerDayResponse_TotalActionRunsPerDay {
		return &api.GetDataActionRunsPerDayResponse_TotalActionRunsPerDay{
			Date:  data.Date.String(),
			Total: int32(data.Total),
		}
	})

	return &api.GetDataActionRunsPerDayResponse{
		Code:    0,
		Message: "Success",
		Results: runsPerDay,
	}, nil
}

func (b business) ProcessGetDataRunsProportion(ctx context.Context, accountUuid string) (*api.GetDataRunsProportionResponse, error) {
	results, err := b.repository.DataActionRunRepository.GetTotalRunsPerType(ctx, &repository.GetTotalRunsPerTypeParams{AccountUuid: accountUuid})
	if err != nil {
		b.log.WithName("ProcessGetTotalRunsPerDay").Error(err, "cannot get data action runs per day", "uuid", accountUuid)
		return nil, err
	}

	typeCount := map[string]int32{
		"segmentation": 0,
		"source":       0,
		"destination":  0,
	}
	for _, each := range results {
		switch each.Type {
		case model.ActionType_ImportDataFromMySQL, model.ActionType_ImportDataFromFile, model.ActionType_ImportDataFromS3:
			typeCount["source"] = typeCount["source"] + int32(each.Total)
		case model.ActionType_ExportDataToS3CSV, model.ActionType_ExportToMySQL, model.ActionType_ExportGophish:
			typeCount["destination"] = typeCount["destination"] + int32(each.Total)
		case model.ActionType_CreateSegment, model.ActionType_CreateMasterSegment, model.ActionType_TrainPredictModel, model.ActionType_ApplyPredictModel:
			typeCount["segmentation"] = typeCount["segmentation"] + int32(each.Total)
		default:
		}
	}

	return &api.GetDataRunsProportionResponse{
		Code:    0,
		Message: "Success",
		Results: []*api.GetDataRunsProportionResponse_CategoryCount{
			{
				Category: "Source",
				Count:    typeCount["source"],
			},
			{
				Category: "Destination",
				Count:    typeCount["destination"],
			},
			{
				Category: "Segmentation",
				Count:    typeCount["segmentation"],
			},
		},
	}, nil
}
