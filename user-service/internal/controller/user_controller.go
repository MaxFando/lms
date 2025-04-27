package controller

import (
	"context"

	userservicev1 "github.com/MaxFando/lms/user-service/api/grpc/gen/go/user-service/v1"
	"github.com/MaxFando/lms/user-service/internal/service"
	"github.com/MaxFando/lms/user-service/internal/model"
)

type UserController struct {
	userservicev1.UnimplementedUserServiceServer
	svc *service.UserService
}

func NewUserController(svc *service.UserService) *UserController {
	return &UserController{svc: svc}
}

func (c *UserController) Register(ctx context.Context, req *userservicev1.RegisterRequest) (*userservicev1.RegisterResponse, error) {
	user, access, refresh, err := c.svc.Register(ctx, req.Name, req.Password)
	if err != nil {
		return nil, err
	}
	return &userservicev1.RegisterResponse{
		User: &userservicev1.User{
			Id:   user.ID,
			Name: user.Name,
			Role: user.Role,
		},
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (c *UserController) Login(ctx context.Context, req *userservicev1.LoginRequest) (*userservicev1.LoginResponse, error) {
	user, access, refresh, err := c.svc.Login(ctx, req.Name, req.Password)
	if err != nil {
		return nil, err
	}
	return &userservicev1.LoginResponse{
		User: &userservicev1.User{
			Id:   user.ID,
			Name: user.Name,
			Role: user.Role,
		},
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (c *UserController) GetUser(ctx context.Context, req *userservicev1.GetUserRequest) (*userservicev1.GetUserResponse, error) {
	user, err := c.svc.GetUser(ctx, req.Id)
	if err != nil || user == nil {
		return nil, err
	}
	return &userservicev1.GetUserResponse{
		User: &userservicev1.User{
			Id:   user.ID,
			Name: user.Name,
			Role: user.Role,
		},
	}, nil
}

func (c *UserController) ListUsers(ctx context.Context, req *userservicev1.ListUsersRequest) (*userservicev1.ListUsersResponse, error) {
	users, err := c.svc.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	resp := &userservicev1.ListUsersResponse{}
	for _, user := range users {
		resp.Users = append(resp.Users, &userservicev1.User{
			Id:   user.ID,
			Name: user.Name,
			Role: user.Role,
		})
	}
	return resp, nil
}