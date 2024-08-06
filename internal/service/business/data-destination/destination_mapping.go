package data_destination

import (
	"context"
	"encoding/json"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (b business) ProcessGetListDestinationMap(ctx context.Context, request *api.GetListDestinationMapRequest, accountUuid string) (*api.GetListDestinationMapResponse, error) {
	logger := b.log.WithName("ProcessGetListDestinationMap").WithValues("destId", request.DestinationId)

	// 1. Validate destination
	destination, err := b.repository.DataDestinationRepository.GetDataDestination(ctx, request.DestinationId)
	if err != nil {
		logger.Error(err, "cannot get destination information")
		return nil, err
	}
	if destination.AccountUuid.String() != accountUuid {
		return nil, status.Error(codes.PermissionDenied, "Only owner can access destination")
	}

	// 2. Get table maps
	tableMaps, err := b.GetDestinationTableMaps(ctx, request.DestinationId)
	if err != nil {
		logger.Error(err, "cannot get table maps")
		return nil, err
	}

	msSegmentMaps, err := b.GetDestinationMasterSegmentMaps(ctx, request.DestinationId)
	if err != nil {
		logger.Error(err, "cannot get master segment maps")
		return nil, err
	}

	segmentMaps, err := b.GetDestinationSegmentMaps(ctx, request.DestinationId)
	if err != nil {
		logger.Error(err, "cannot get segment maps")
		return nil, err
	}

	return &api.GetListDestinationMapResponse{
		Code:    0,
		Message: "Success",
		Count:   int64(len(tableMaps) + len(msSegmentMaps) + len(segmentMaps)),
		Results: append(append(tableMaps, msSegmentMaps...), segmentMaps...),
	}, nil
}

func (b business) GetDestinationTableMaps(ctx context.Context, destinationId int64) ([]*api.DestinationMappings, error) {
	tableMaps, err := b.repository.DestTableMapRepository.ListDestinationTableMaps(ctx, &repository.ListDestinationTableMapsParams{DestinationId: destinationId})
	if err != nil {
		return nil, err
	}

	mapType := model.DestTableMap{}.TableName()
	objectType := model.DataTable{}.TableName()
	destMappings := make([]*api.DestinationMappings, 0, len(tableMaps))
	for _, mapping := range tableMaps {
		var mappingOptions []*api.MappingOptionItem
		if mapping.MappingOptions.Valid && mapping.MappingOptions.RawMessage != nil {
			err = json.Unmarshal(mapping.MappingOptions.RawMessage, &mappingOptions)
			if err != nil {
				return nil, err
			}
		}
		destMappings = append(destMappings, &api.DestinationMappings{
			Id:           mapping.ID,
			Type:         mapType,
			ObjectType:   objectType,
			ObjectName:   mapping.TableName,
			ObjectId:     mapping.TableId,
			Mappings:     mappingOptions,
			DataActionId: mapping.DataActionId,
		})
	}

	return destMappings, nil
}

func (b business) GetDestinationSegmentMaps(ctx context.Context, destinationId int64) ([]*api.DestinationMappings, error) {
	segmentMaps, err := b.repository.DestSegmentMapRepository.ListDestinationSegmentMaps(ctx, &repository.ListDestinationSegmentMapsParams{DestinationId: destinationId})
	if err != nil {
		return nil, err
	}

	mapType := model.DestSegmentMap{}.TableName()
	objectType := model.Segment{}.TableName()
	destMappings := make([]*api.DestinationMappings, 0, len(segmentMaps))
	for _, mapping := range segmentMaps {
		var mappingOptions []*api.MappingOptionItem
		if mapping.MappingOptions.Valid && mapping.MappingOptions.RawMessage != nil {
			err = json.Unmarshal(mapping.MappingOptions.RawMessage, &mappingOptions)
			if err != nil {
				return nil, err
			}
		}
		destMappings = append(destMappings, &api.DestinationMappings{
			Id:           mapping.ID,
			Type:         mapType,
			ObjectType:   objectType,
			ObjectName:   mapping.SegmentName,
			ObjectId:     mapping.SegmentId,
			Mappings:     mappingOptions,
			DataActionId: mapping.DataActionId,
		})
	}

	return destMappings, nil
}

func (b business) GetDestinationMasterSegmentMaps(ctx context.Context, destinationId int64) ([]*api.DestinationMappings, error) {
	masterSegmentMaps, err := b.repository.DestMasterSegmentMapRepository.ListDestinationMasterSegmentMaps(ctx, &repository.ListDestinationMasterSegmentMapsParams{DestinationId: destinationId})
	if err != nil {
		return nil, err
	}

	mapType := model.DestMasterSegmentMap{}.TableName()
	objectType := model.MasterSegment{}.TableName()
	destMappings := make([]*api.DestinationMappings, 0, len(masterSegmentMaps))
	for _, mapping := range masterSegmentMaps {
		var mappingOptions []*api.MappingOptionItem
		if mapping.MappingOptions.Valid && mapping.MappingOptions.RawMessage != nil {
			err = json.Unmarshal(mapping.MappingOptions.RawMessage, &mappingOptions)
			if err != nil {
				return nil, err
			}
		}
		destMappings = append(destMappings, &api.DestinationMappings{
			Id:           mapping.ID,
			Type:         mapType,
			ObjectType:   objectType,
			ObjectName:   mapping.MasterSegmentName,
			ObjectId:     mapping.MasterSegmentId,
			Mappings:     mappingOptions,
			DataActionId: mapping.DataActionId,
		})
	}

	return destMappings, nil
}
