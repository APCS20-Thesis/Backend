package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AccountRepository interface {
	FindAccount(ctx context.Context, username string, password string) (*model.Account, error)
	CreateAccount(ctx context.Context, params *CreateAccountParams) error
	GetInfo(ctx context.Context, username string) (*model.Account, error)
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
	err := r.WithContext(ctx).Table(r.TableName).Where("username = ? AND password = ?", username, password).First(&account).Error
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
	var count int64
	err := r.WithContext(ctx).Table(r.TableName).Where("username = ?", params.Username).Count(&count).Error
	if err != nil {
		return err
	}
	if count != 0 {
		return status.Errorf(codes.AlreadyExists, "Username is already used")
	}
	account := &model.Account{Username: params.Username, Password: params.Password, Email: params.Email, FirstName: params.FirstName, LastName: params.LastName}

	err = r.WithContext(ctx).Create(&account).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *accountRepo) GetInfo(ctx context.Context, username string) (*model.Account, error) {
	var account model.Account
	err := r.WithContext(ctx).Table(r.TableName).Where("username = ?", username).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
