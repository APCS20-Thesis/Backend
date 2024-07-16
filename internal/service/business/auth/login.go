package auth

import (
	"context"
	"errors"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (b *business) ProcessLogin(ctx context.Context, request *api.LoginRequest) (*model.Account, error) {
	account, err := b.repository.AccountRepository.FindAccount(ctx, request.Username, request.Password)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		b.log.WithName("ProcessLogin").
			WithValues("request", request).Error(err, "Wrong username or password.")
		return nil, status.Error(codes.InvalidArgument, "Wrong username or password.")
	}
	if err != nil {
		b.log.WithName("ProcessLogin").
			WithValues("request", request).Error(err, "Fail to get user data")
		return nil, status.Error(codes.Internal, "Fail to get account")
	}

	return account, nil
}
