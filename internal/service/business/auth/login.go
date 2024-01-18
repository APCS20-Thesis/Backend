package auth

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
)

func (b *business) ProcessLogin(ctx context.Context, request *api.LoginRequest) (*api.Account, error) {
	account, err := b.repository.AccountRepository.FindAccount(ctx, request.Username, request.Password)
	if err != nil {
		b.log.WithName("ProcessLogin").WithValues("request", request).Error(err, "Fail to get user data")
		return nil, err
	}

	return &api.Account{
		Username:  account.Username,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Email:     account.Email,
	}, nil
}
