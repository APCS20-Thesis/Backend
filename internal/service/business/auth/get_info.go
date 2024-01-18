package auth

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
)

func (b *business) ProcessGetInfo(ctx context.Context, username string) (*api.Account, error) {
	account, err := b.repository.AccountRepository.GetInfo(ctx, username)
	if err != nil {
		b.log.WithName("ProcessGetInfo").WithValues("username", username).Error(err, "Fail to get info user")
		return nil, err
	}

	return &api.Account{
		Username:  account.Username,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Email:     account.Email,
	}, nil
}
