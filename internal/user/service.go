package user

import (
	"context"
	"fmt"
	"time"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/jwt"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/password"
)

type Service struct {
	repository Repository
}

func (s *Service) Create(ctx context.Context, req CreateUserPayload) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
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
	err = s.repository.Create(ctx, user)
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

func (s *Service) Login(req LoginPayload) (*UserResponse, error) {
	user, err := s.repository.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	match, err := password.Matches(req.Password, user.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrWrongUsernameOrPassword
	}
	// create access token with signed jwt
	return &UserResponse{
		Username:    user.Username,
		Name:        user.Name,
		AccessToken: "TODO",
	}, nil
}
