package data_table

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/adapter/airflow"
	"github.com/APCS20-Thesis/Backend/internal/adapter/query"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type Business interface {
	CreateDataTable(ctx context.Context, params *repository.CreateDataTableParams) (*model.DataTable, error)
	UpdateDataTable(ctx context.Context, params *repository.UpdateDataTableParams) error
	GetDataTable(ctx context.Context, request *api.GetDataTableRequest, accountUuid string) (*api.GetDataTableResponse, error)
	GetListDataTables(ctx context.Context, request *api.GetListDataTablesRequest, accountUuid string) ([]*api.GetListDataTablesResponse_DataTable, int64, error)
	GetListFileExportRecords(ctx context.Context, request *api.GetListFileExportRecordsRequest, accountUuid string) ([]*api.GetListFileExportRecordsResponse_FileExportRecord, error)
	GetQueryDataTable(ctx context.Context, request *api.GetQueryDataTableRequest, accountUuid string) (*api.GetQueryDataTableResponse, error)
	ProcessGetListDataActions(ctx context.Context, request *api.GetListDataActionsRequest, accountUuid string) (*api.GetListDataActionsResponse, error)
}

type business struct {
	log            logr.Logger
	repository     *repository.Repository
	airflowAdapter airflow.AirflowAdapter
	queryAdapter   query.QueryAdapter
}

func NewDataTableBusiness(log logr.Logger, repository *repository.Repository, airflowAdapter airflow.AirflowAdapter, queryAdapter query.QueryAdapter) Business {
	return &business{
		log:            log.WithName("DataTableBiz"),
		repository:     repository,
		airflowAdapter: airflowAdapter,
		queryAdapter:   queryAdapter,
	}
}

func (b business) ProcessGetListDataActions(ctx context.Context, request *api.GetListDataActionsRequest, accountUuid string) (*api.GetListDataActionsResponse, error) {
	dataActions, err := b.repository.DataActionRepository.GetListDataActions(ctx, &repository.GetListDataActionsParams{
		AccountUuid: uuid.MustParse(accountUuid),
		Page:        int(request.Page),
		PageSize:    int(request.PageSize),
	})
	if err != nil {
		b.log.WithName("list data actions").Error(err, "cannot get data actions")
		return nil, err
	}

	return &api.GetListDataActionsResponse{
		Code:    0,
		Message: "Success",
		Count:   50,
		Results: utils.Map(dataActions, func(modelDataAction model.DataAction) *api.DataAction {
			return &api.DataAction{
				Id:         modelDataAction.ID,
				ActionType: string(modelDataAction.ActionType),
				Status:     string(modelDataAction.Status),
				CreatedAt:  modelDataAction.CreatedAt.String(),
				UpdatedAt:  modelDataAction.UpdatedAt.String(),
			}
		}),
	}, nil

}
