package service

import (
	"context"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/MaxFando/lms/user-service/internal/model"
	"github.com/MaxFando/lms/user-service/internal/repository"
	"github.com/MaxFando/lms/user-service/internal/jwt"
)

type UserService struct {
	repo      repository.UserRepository
	jwt       jwt.JWTService
}

func NewUserService(repo repository.UserRepository, jwt jwt.JWTService) *UserService {
	return &UserService{repo: repo, jwt: jwt}
}

func (s *UserService) Register(ctx context.Context, name, password string) (*model.User, string, string, error) {
	existing, _ := s.repo.FindByName(ctx, name)
	if existing != nil {
		return nil, "", "", errors.New("user already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}
	user := &model.User{
		Name:     name,
		Password: string(hash),
		Role:     "USER",
	}
	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, "", "", err
	}
	user.ID = userID

	accessToken, refreshToken, err := s.jwt.GenerateTokens(user)
	if err != nil {
		return nil, "", "", err
	}
	_ = s.repo.UpdateRefreshToken(ctx, userID, refreshToken)
	return user, accessToken, refreshToken, nil
}

func (s *UserService) Login(ctx context.Context, name, password string) (*model.User, string, string, error) {
	user, err := s.repo.FindByName(ctx, name)
	if err != nil || user == nil {
		return nil, "", "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", errors.New("invalid password")
	}
	accessToken, refreshToken, err := s.jwt.GenerateTokens(user)
	if err != nil {
		return nil, "", "", err
	}
	_ = s.repo.UpdateRefreshToken(ctx, user.ID, refreshToken)
	return user, accessToken, refreshToken, nil
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.List(ctx)
}
