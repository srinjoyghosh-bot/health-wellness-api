package server

import (
	"context"
	userpb "user-service/internal/genproto/user"
	"user-service/internal/models"
	"user-service/internal/service"
	"user-service/internal/utils"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	service service.UserService
}

func NewUserServer(service service.UserService) *UserServer {
	return &UserServer{service: service}
}

func (s *UserServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	user := models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	err := s.service.Create(&user)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to create user : " + err.Error())
	}

	return &userpb.RegisterResponse{UserId: int32(user.ID)}, nil

}

func (s *UserServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	user, err := s.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return nil, utils.NewUnauthorizedError("Unauthorized : " + err.Error())
	}
	return &userpb.LoginResponse{UserId: int32(user.ID)}, nil
}
