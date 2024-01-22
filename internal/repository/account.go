package repository

import (
	"context"
	"errors"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AccountRepository interface {
	FindAccount(ctx context.Context, username string, password string) (*model.Account, error)
	CreateAccount(ctx context.Context, params *CreateAccountParams) error
	GetAccountInfo(ctx context.Context, accountUuid string) (*model.Account, error)
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
		Where("username = ? AND password = ?", username, password).
		First(&account).Error
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
		account := &model.Account{
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
