package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/MaxFando/lms/user-service/internal/jwt"
	"github.com/MaxFando/lms/user-service/internal/model"
	pubsubPkg "github.com/MaxFando/lms/user-service/internal/pubsub"
	"github.com/MaxFando/lms/user-service/internal/repository"
)

type UserService struct {
	repo   repository.UserRepository
	jwt    jwt.JWTService
	pubsub pubsubPkg.PubSub
}

func NewUserService(
	repo repository.UserRepository,
	jwtSvc jwt.JWTService,
	ps pubsubPkg.PubSub,
) *UserService {
	return &UserService{repo: repo, jwt: jwtSvc, pubsub: ps}
}

func (s *UserService) Register(ctx context.Context, name, password string) (*model.User, string, string, error) {
	existing, _ := s.repo.FindByName(ctx, name)
	if existing != nil {
		return nil, "", "", errors.New("пользователь уже существует")
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, "", "", err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	u := &model.User{
		Name:     name,
		Password: string(hash),
		Role:     "USER",
	}

	id, err := s.repo.CreateTx(ctx, tx, u)
	if err != nil {
		return nil, "", "", err
	}
	u.ID = id

	at, rt, err := s.jwt.GenerateTokens(u)
	if err != nil {
		return nil, "", "", err
	}

	if err = s.repo.UpdateRefreshTokenTx(ctx, tx, id, rt); err != nil {
		return nil, "", "", err
	}
	if err = tx.Commit(); err != nil {
		return nil, "", "", err
	}

	_ = s.pubsub.PublishUserCreated(ctx, id)
	return u, at, rt, nil
}

func (s *UserService) Login(ctx context.Context, name, password string) (*model.User, string, string, error) {
	u, err := s.repo.FindByName(ctx, name)
	if err != nil || u == nil {
		return nil, "", "", errors.New("пользователь не найден")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, "", "", errors.New("неверный пароль")
	}

	at, rt, err := s.jwt.GenerateTokens(u)
	if err != nil {
		return nil, "", "", err
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, "", "", err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = s.repo.UpdateRefreshTokenTx(ctx, tx, u.ID, rt); err != nil {
		return nil, "", "", err
	}
	if err = tx.Commit(); err != nil {
		return nil, "", "", err
	}

	_ = s.pubsub.PublishUserLoggedIn(ctx, u.ID)
	return u, at, rt, nil
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.List(ctx)
}
