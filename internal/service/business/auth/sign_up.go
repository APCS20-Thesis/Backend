package auth

import (
	"context"
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/APCS20-Thesis/Backend/internal/repository"
	"google.golang.org/grpc/codes"
)

func (b *business) ProcessSignUp(ctx context.Context, request *api.SignUpRequest) (*api.CommonResponse, error) {
	err := b.repository.AccountRepository.
		CreateAccount(ctx, &repository.CreateAccountParams{
			Username:  request.Username,
			Password:  request.Password,
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Email:     request.Email,
		})
	if err != nil {
		b.log.WithName("ProcessSignUp").WithValues("request", request).Error(err, "Fail to create user data")
		return nil, err
	}

	return &api.CommonResponse{
		Code:    int32(codes.OK),
		Message: "Sign up success",
	}, nil
}
