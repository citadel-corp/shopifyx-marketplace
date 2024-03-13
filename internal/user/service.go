package user

import (
	"context"
	"fmt"
	"time"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/jwt"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/password"
)

type Service struct {
	Repository Repository
}

func (s *Service) Create(ctx context.Context, req CreateUserPayload) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	user := &User{
		Username:       req.Username,
		Name:           req.Name,
		HashedPassword: hashedPassword,
	}
	err = s.Repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*24, fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		Username:    req.Username,
		Name:        req.Name,
		AccessToken: accessToken,
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginPayload) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.Repository.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	match, err := password.Matches(req.Password, user.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrWrongPassword
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*24, fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		Username:    user.Username,
		Name:        user.Name,
		AccessToken: accessToken,
	}, nil
}
