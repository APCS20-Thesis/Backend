package auth

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"github.com/APCS20-Thesis/Backend/utils"
)

func (b *business) ProcessGetAccountInfo(ctx context.Context, accountUuid string) (*api.Account, *api.Setting, error) {
	account, err := b.repository.AccountRepository.GetAccountInfo(ctx, accountUuid)
	if err != nil {
		b.log.WithName("ProcessGetAccountInfo").WithValues("Context", ctx).Error(err, "Fail to get info user")
		return nil, nil, err
	}
	setting, err := b.repository.AccountRepository.GetAccountSetting(ctx, accountUuid)
	if err != nil {
		b.log.WithName("ProcessGetAccountInfo").WithValues("Context", ctx).Error(err, "Fail to get setting of user")
		return nil, nil, err
	}

	return &api.Account{
			Uuid:      account.Uuid.String(),
			Username:  account.Username,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Email:     account.Email,
			Phone:     account.Phone,
			Country:   account.Country,
			Company:   account.Company,
			Position:  account.Position,
		}, &api.Setting{
			NotifyCreateSource:        setting.NotifyCreateSource,
			NotifyCreateDestination:   setting.NotifyCreateDestination,
			NotifyCreateMasterSegment: setting.NotifyCreateMasterSegment,
			NotifyCreateSegment:       setting.NotifyCreateSegment,
		}, nil
}
func (b *business) ProcessUpdateAccountInfo(ctx context.Context, request *api.UpdateAccountInfoRequest, accountUuid string) (*api.Account, error) {
	account, err := b.repository.AccountRepository.UpdateAccountInfo(ctx, &repository.UpdateAccountInfoParams{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Phone:     request.Phone,
		Country:   request.Country,
		Company:   request.Company,
		Position:  request.Position,
	}, accountUuid)
	if err != nil {
		b.log.WithName("ProcessGetAccountInfo").WithValues("Context", ctx).Error(err, "Fail to get info user")
		return nil, err
	}
	return &api.Account{
		Uuid:      account.Uuid.String(),
		Username:  account.Username,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Email:     account.Email,
		Phone:     account.Phone,
		Country:   account.Country,
		Company:   account.Company,
		Position:  account.Position,
	}, nil
}

func (b *business) ProcessUpdateAccountSetting(ctx context.Context, request *api.UpdateAccountSettingRequest, accountUuid string) (*api.Setting, error) {
	setting, err := b.repository.AccountRepository.UpdateAccountSetting(ctx, &repository.UpdateAccountSettingParams{
		NotifyCreateSource:        utils.ConvertWrappersBoolToBoolAdd(request.NotifyCreateSource),
		NotifyCreateDestination:   utils.ConvertWrappersBoolToBoolAdd(request.NotifyCreateDestination),
		NotifyCreateMasterSegment: utils.ConvertWrappersBoolToBoolAdd(request.NotifyCreateMasterSegment),
		NotifyCreateSegment:       utils.ConvertWrappersBoolToBoolAdd(request.NotifyCreateSegment),
	}, accountUuid)
	if err != nil {
		b.log.WithName("ProcessGetAccountInfo").WithValues("Context", ctx).Error(err, "Fail to get info user")
		return nil, err
	}
	return &api.Setting{
		NotifyCreateSource:        setting.NotifyCreateSource,
		NotifyCreateDestination:   setting.NotifyCreateDestination,
		NotifyCreateMasterSegment: setting.NotifyCreateMasterSegment,
		NotifyCreateSegment:       setting.NotifyCreateSegment,
	}, nil
}
