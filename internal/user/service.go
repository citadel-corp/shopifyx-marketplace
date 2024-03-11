package user

import "github.com/citadel-corp/shopifyx-marketplace/internal/common/password"

type Service struct {
	repository Repository
}

func (s *Service) Create(req CreateUserPayload) (*UserResponse, error) {
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
	err = s.repository.Create(user)
	if err != nil {
		return nil, err
	}
	// create access token with signed jwt
	return &UserResponse{
		Username:    req.Username,
		Name:        req.Name,
		AccessToken: "TODO",
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
