package clients

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	userpb "healthApi/api-gateway/genproto/user"
	"healthApi/api-gateway/internal/models"
	"healthApi/api-gateway/internal/utils"
)

type UserClient struct {
	client userpb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(serverAddr string) (UserClient, error) {
	conn, err := grpc.NewClient(":50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return UserClient{}, err
	}

	return UserClient{
		client: userpb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) Register(ctx context.Context, req models.UserRegisterRequest) (uint, error) {
	resp, err := c.client.Register(ctx, &userpb.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		return 0, utils.NewInternalServerError("Failed to create user : " + err.Error())
	}
	return uint(resp.UserId), nil
}

func (c *UserClient) Login(ctx context.Context, req models.UserLoginRequest) (uint, error) {
	resp, err := c.client.Login(ctx, &userpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return 0, utils.NewInternalServerError("Failed to fetch all exercises : " + err.Error())
	}
	return uint(resp.UserId), nil
}
