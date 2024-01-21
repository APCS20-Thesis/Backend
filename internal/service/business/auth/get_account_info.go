package auth

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
)

func (b *business) ProcessGetAccountInfo(ctx context.Context) (*api.Account, error) {
	account, err := b.repository.AccountRepository.GetAccountInfo(ctx)
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
	}, nil
}
