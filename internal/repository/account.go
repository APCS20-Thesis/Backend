package repository

import (
	"context"
	"errors"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepository interface {
	FindAccount(ctx context.Context, username string, password string) (*model.Account, error)
	CreateAccount(ctx context.Context, params *CreateAccountParams) error
	GetAccountInfo(ctx context.Context, accountUuid string) (*model.Account, error)
	GetAccountSetting(ctx context.Context, accountUuid string) (*model.Setting, error)
	UpdateAccountInfo(ctx context.Context, params *UpdateAccountInfoParams, accountUuid string) (*model.Account, error)
	UpdateAccountSetting(ctx context.Context, params *UpdateAccountSettingParams, accountUuid string) (*model.Setting, error)
}

type accountRepo struct {
	*gorm.DB
	TableName string
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepo{db, model.Account{}.TableName()}
}

func (r *accountRepo) FindAccount(ctx context.Context, username string, password string) (*model.Account, error) {
	var account model.Account
	err := r.WithContext(ctx).Table(r.TableName).
		Where("username = ?", username).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return &account, nil
}

type CreateAccountParams struct {
	Username  string
	Password  string
	FirstName string
	LastName  string
	Email     string
}

func (r *accountRepo) CreateAccount(ctx context.Context, params *CreateAccountParams) error {
	var account model.Account

	err := r.WithContext(ctx).Table(r.TableName).
		Where("username = ? OR email = ?", params.Username, params.Email).
		First(&account).Error

	if err == nil {
		if account.Username == params.Username {
			return status.Error(codes.AlreadyExists, "Username is already used")
		}
		return status.Error(codes.AlreadyExists, "Email is already used")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		account := model.Account{
			Username:  params.Username,
			Password:  params.Password,
			Email:     params.Email,
			FirstName: params.FirstName,
			LastName:  params.LastName,
		}

		createErr := r.WithContext(ctx).Table(r.TableName).Create(&account).Error
		if createErr != nil {
			return err
		}
		return nil
	}
	return err
}

func (r *accountRepo) GetAccountInfo(ctx context.Context, accountUuid string) (*model.Account, error) {
	var account model.Account

	err := r.WithContext(ctx).Table(r.TableName).
		Where("uuid = ?", accountUuid).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepo) GetAccountSetting(ctx context.Context, accountUuid string) (*model.Setting, error) {
	var setting model.Setting

	err := r.WithContext(ctx).Table(model.Setting{}.TableName()).
		Where("account_uuid = ?", accountUuid).
		First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

type UpdateAccountInfoParams struct {
	FirstName string
	LastName  string
	Phone     string
	Country   string
	Company   string
	Position  string
}

func (r *accountRepo) UpdateAccountInfo(ctx context.Context, params *UpdateAccountInfoParams, accountUuid string) (*model.Account, error) {
	account := &model.Account{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Phone:     params.Phone,
		Country:   params.Country,
		Company:   params.Company,
		Position:  params.Position,
	}

	updateErr := r.WithContext(ctx).Table(r.TableName).Clauses(clause.Returning{}).Where("uuid = ?", accountUuid).Updates(&account).Error
	if updateErr != nil {
		return nil, updateErr
	}
	return account, nil
}

type UpdateAccountSettingParams struct {
	NotifyCreateSource        api.Bool
	NotifyCreateDestination   api.Bool
	NotifyCreateMasterSegment api.Bool
	NotifyCreateSegment       api.Bool
}

func (r *accountRepo) UpdateAccountSetting(ctx context.Context, params *UpdateAccountSettingParams, accountUuid string) (*model.Setting, error) {
	setting := &model.Setting{
		NotifyCreateSource:        params.NotifyCreateSource,
		NotifyCreateDestination:   params.NotifyCreateDestination,
		NotifyCreateMasterSegment: params.NotifyCreateMasterSegment,
		NotifyCreateSegment:       params.NotifyCreateMasterSegment,
	}

	updateErr := r.WithContext(ctx).Table(model.Setting{}.TableName()).Clauses(clause.Returning{}).Where("account_uuid = ?", accountUuid).Updates(&setting).Error
	if updateErr != nil {
		return nil, updateErr
	}
	return setting, nil
}
