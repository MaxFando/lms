package repository

import (
	"context"
	"github.com/MaxFando/lms/user-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	FindByName(ctx context.Context, name string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	UpdateRefreshToken(ctx context.Context, userID int64, token string) error
}